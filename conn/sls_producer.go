package conn

import (
	"errors"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"time"
)

type LogProjectConf struct {
	Project  string
	Logstore string
	Topic    string
	Source   string
}

var SLSProducer *producer.Producer
var accessProject, accesslogstore, _topic, _source string
var defaultAccessLog, defaultInfoLog, defaultErrorLog LogProjectConf

//初始化SLSProducer，设置默认project logstore，topic,source
func InitSLSProducer(producerConfig *producer.ProducerConfig) {
	SLSProducer = producer.InitProducer(producerConfig)
	SLSProducer.Start()
	PutNewConn(func() { //程序关闭时调用conn.Close()
		SLSProducer.SafeClose()
	})
}
func SetDefaultAccessLog(project, logstore, topic, source string) {
	defaultAccessLog.Source = source
	defaultAccessLog.Topic = topic
	defaultAccessLog.Logstore = logstore
	defaultAccessLog.Project = project
}
func GetDefaultAccessLogConfig() LogProjectConf {
	return defaultAccessLog
}

func SetDefaultInfoLog(project, logstore, topic, source string) {
	defaultInfoLog.Source = source
	defaultInfoLog.Topic = topic
	defaultInfoLog.Logstore = logstore
	defaultInfoLog.Project = project
}
func GetDefaultInfoLogConfig() LogProjectConf {
	return defaultInfoLog
}
func SetDefaultErrorLog(project, logstore, topic, source string) {
	defaultErrorLog.Source = source
	defaultErrorLog.Topic = topic
	defaultErrorLog.Logstore = logstore
	defaultErrorLog.Project = project
}
func GetDefaultErrorLogConfig() LogProjectConf {
	return defaultErrorLog
}

//发送map类型数据到logstore
func SendMap(data map[string]string, project LogProjectConf) error {
	log := producer.GenerateLog(uint32(time.Now().Unix()), data)
	if SLSProducer == nil || len(project.Project) == 0 || len(project.Logstore) == 0 {
		return errors.New("SLSProducer未初始化")
	}
	return SLSProducer.SendLog(project.Project, project.Logstore, project.Topic, project.Source, log)
}

//发送log类型数据到logstore
func SendLog(log *sls.Log, project LogProjectConf) error {
	if SLSProducer == nil || len(project.Project) == 0 || len(project.Logstore) == 0 {
		return errors.New("SLSProducer未初始化")
	}
	return SLSProducer.SendLog(project.Project, project.Logstore, project.Topic, project.Source, log)
}
