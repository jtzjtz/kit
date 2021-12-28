package config

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/logger"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"log"
	"os"
)

type IDataConfig interface {
	//重载接口
	Reload() error
	//检查接口
	IsLoad() bool
}

type NacosConfig struct {
	Address             string
	Port                uint64
	Scheme              string
	ContextPath         string
	NameSpaceId         string
	TimeoutMs           uint64
	NotLoadCacheAtStart bool
	LogDir              string
	CacheDir            string
	RotateTime          string
	MaxAge              int64
	LogLevel            string
}
type Client struct {
	conn        config_client.IConfigClient
	cacheDIr    string
	isConnected bool
}

//初始化nacos客户端
func InitClient(conf NacosConfig) (Client, error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(
			conf.Address,
			conf.Port,
			constant.WithScheme(conf.Scheme),
			constant.WithContextPath(conf.ContextPath)),
	}

	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(conf.NameSpaceId),
		constant.WithTimeoutMs(conf.TimeoutMs),
		constant.WithNotLoadCacheAtStart(conf.NotLoadCacheAtStart),
		constant.WithLogDir(conf.LogDir),
		constant.WithCacheDir(conf.CacheDir),
		constant.WithRotateTime(conf.RotateTime),
		constant.WithMaxAge(conf.MaxAge),
		constant.WithLogLevel(conf.LogLevel),
	)
	nacosClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		log.Fatal("nacos init err", err)
		return Client{isConnected: false}, err

	}
	c := Client{conn: nacosClient, cacheDIr: conf.CacheDir, isConnected: true}
	return c, err
}

//获取配置信息 返回字符串
func (t *Client) GetConfig(dataid string, group string) (string, error) {
	if t.isConnected == false {
		return "", errors.New("nacos 连接失败")
	}
	data, er := t.conn.GetConfig(vo.ConfigParam{
		DataId: dataid,
		Group:  group,
	})
	if er != nil {
		log.Fatal("nacos get config err", er)
		return "", er

	}
	return data, nil
}

//获取配置信息 转为targetConfig isWatch 是否监听配置变更
func (t *Client) GetDataConfig(dataid string, group string, targetConfig IDataConfig, isWatch bool) error {
	data, err := t.GetConfig(dataid, group)
	if err != nil {
		return err
	}

	if saveError := saveConfig(data, dataid, group, targetConfig, t.cacheDIr); saveError != nil {
		return saveError
	}
	if isWatch {
		return t.WatchConfig(dataid, group, targetConfig, true)
	}
	return nil
}

//监听
func (t *Client) WatchConfig(dataid string, group string, targetConfig IDataConfig, isReload bool) error {
	err := t.conn.ListenConfig(vo.ConfigParam{
		DataId: dataid,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			//是否重载变量
			shouldReload := isReload
			//因为开启监听也会调用OnChange函数，所以加此判断
			//if  !model.IsLoad(){
			//	shouldReload=false  //如果model 没有被加载过,则不需要做重载
			//}
			err := saveConfig(data, dataid, group, targetConfig, t.cacheDIr)
			if err != nil {
				return
			}

			//重载
			if shouldReload { //重载关键代码
				err := targetConfig.Reload()
				if err != nil {
					logger.Error(err)
					return
				} else {
					logger.Info(dataid, " 重载完成")
				}
			}

		},
	})
	if err != nil {
		logger.Error("listen config error,dataid:", dataid, err)
	}
	return err
}

func saveConfig(data string, dataid string, group string, targetConfig IDataConfig, confDir string) error {
	//线上配置缓存文件
	cacheFile := fmt.Sprintf("%s/%s-%s.yaml", confDir, group, dataid)
	file, err := os.OpenFile(cacheFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer file.Close()
	//写入缓存文件
	_, err = file.WriteString(data)
	if err != nil {
		logger.Error(err)
		return err
	}
	//从缓存文件加载配置
	err = config.LoadFile(cacheFile)
	if err != nil {
		logger.Error(err)
		return err
	}
	//接收
	err = config.Scan(targetConfig)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
