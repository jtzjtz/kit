package balance

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var cli naming_client.INamingClient

type NacosConfig struct {
	NacosIp     string
	NacosPort   int
	NamespaceId string
	GroupName   string
}

//创建nacos 客户端
func createNacosClient(conf NacosConfig) (naming_client.INamingClient, error) {
	sc := []constant.ServerConfig{
		{
			IpAddr: conf.NacosIp,
			Port:   uint64(conf.NacosPort),
		},
	}
	cc := constant.ClientConfig{
		NamespaceId:         conf.NamespaceId, //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "./tmp/nacos/log",
		CacheDir:            "./tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "error",
	}

	// a more graceful way to create naming client
	return clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

}
