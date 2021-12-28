package log

import (
	"fmt"
	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/jtzjtz/kit/conn"
	"os"
	"testing"
	"time"
)

func initProducer() {

	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.AccessKeyID = "********"
	producerConfig.Endpoint = "cn-beijing-intranet.log.aliyuncs.com"
	producerConfig.AccessKeySecret = "*******"
	producerConfig.AllowLogLevel = ""
	producerConfig.LogFileName = "./data.log"

	conn.InitSLSProducer(producerConfig)
	conn.SetDefaultAccessLog("app-point-data", "app_data_server_test", "topic", "127.0.0.1")
	//conn.InitSLSProducer(producerConfig, "sq-project", "infolog", "topic", "127.0.0.1")
}
func TestAccessLog(t *testing.T) {

	initProducer()
	defer conn.Close()
	content := make(map[string]string)
	content["event_id"] = "adfadsadsf"
	content["service_name"] = "shequ-service"
	content["server_ip"] = "localIP"
	content["server_hostname"], _ = os.Hostname()
	content["create_time"] = time.Now().Format("2006-01-02 15:04:05")
	if err := AddAccessMap(content); err != nil {
		t.Error(err.Error())
	}

}
func TestCreateClient(t *testing.T) {
	AccessKeyID := "****"                                                           //阿里云访问密钥AccessKey ID。更多信息，请参见访问密钥。阿里云主账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM账号进行API访问或日常运维。
	AccessKeySecret := "*****"                                                      //阿里云访问密钥AccessKey Secret。
	Endpoint := "cn-beijing-intranet.log.aliyuncs.com"                              //日志服务的域名。更多信息，请参见服务入口。此处以杭州为例，其它地域请根据实际情况填写。
	Client := sls.CreateNormalInterface(Endpoint, AccessKeyID, AccessKeySecret, "") //创建Client。

	project, err := Client.CreateProject("sq-project", "社区project") //输入Project名称和描述。
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(project)
	err = Client.CreateLogStore("sq-project", "infolog", 2, 2, true, 16) //输入Project名称、Logstore名称、数据保存时长、Shard数量、开启自动分裂Shard功能和最大分裂数。如果数据保存时长配置为3650，表示永久保存。
	if err != nil {
		fmt.Println(err)
	}

}
