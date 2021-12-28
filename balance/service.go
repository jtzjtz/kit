package balance

import (
	"github.com/nacos-group/nacos-sdk-go/vo"
	"log"
)

//注册服务
func RegisterService(serviceName, groupName, serverIp string, serverPort, weight int, metaData map[string]string, nacosConf NacosConfig) error {
	var err error

	if cli == nil {
		cli, err = createNacosClient(nacosConf)
		if err != nil {
			return err
		}
	}
	//err = withAlive(serviceName, ip, port, ttl)

	_, err = cli.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          serverIp,
		Port:        uint64(serverPort),
		ServiceName: serviceName,
		GroupName:   groupName,
		Weight:      float64(weight),
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    metaData,
	})
	if err != nil {
		log.Println(err)
	}
	return err
}

// UnRegister remove service from nacos
func UnRegisterService(serviceName string, groupName string, serverIp string, serverPort int) {
	if cli != nil {
		cli.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          serverIp,
			Port:        uint64(serverPort),
			ServiceName: serviceName,
			GroupName:   groupName,
			Ephemeral:   true, //it must be true
		})
		log.Printf("group:%s service:%s %s:%v  取消注册", groupName, serviceName, serverIp, serverPort)
	}
}
