package log

import (
	"encoding/json"
	"errors"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/gin-gonic/gin"
	"github.com/gogo/protobuf/proto"
	"github.com/jtzjtz/kit/conn"
	"os"

	"time"
)

//添加access日志到SLS（程序运行前，请初始化conn.InitSLSProducer）
func AddAccessMap(data map[string]string) error {
	logConfig := conn.GetDefaultAccessLogConfig()
	if len(logConfig.Project) == 0 || len(logConfig.Logstore) == 0 {
		return errors.New("未设置logstore")
	}
	return conn.SendMap(data, logConfig)
}

//添加info日志到SLS（程序运行前，请初始化conn.InitSLSProducer）
func AddInfoMap(data map[string]string, project ...conn.LogProjectConf) error {
	logConfig := conn.GetDefaultInfoLogConfig()
	if len(project) > 0 {
		logConfig = project[0]
	}
	if len(logConfig.Project) == 0 || len(logConfig.Logstore) == 0 {
		return errors.New("")
	}
	return conn.SendMap(data, logConfig)
}

//添加info日志到SLS（程序运行前，请初始化conn.InitSLSProducer）
func AddInfo(info InfoEntity, ctx *gin.Context, project ...conn.LogProjectConf) error {
	logConfig := conn.GetDefaultInfoLogConfig()
	if len(project) > 0 {
		logConfig = project[0]
	}
	if len(logConfig.Project) == 0 || len(logConfig.Logstore) == 0 {
		return errors.New("")
	}

	content := []*sls.LogContent{}
	hostname, _ := os.Hostname()
	if ctx != nil {
		content = append(content, []*sls.LogContent{

			&sls.LogContent{
				Key:   proto.String("customer_id"),
				Value: proto.String(ctx.Param("customer_id")),
			},
			&sls.LogContent{
				Key:   proto.String("request_id"),
				Value: proto.String(ctx.GetHeader("Request_Id")),
			}, &sls.LogContent{
				Key:   proto.String("host"),
				Value: proto.String(ctx.Request.Host),
			}, &sls.LogContent{
				Key:   proto.String("uri"),
				Value: proto.String(ctx.Request.RequestURI),
			},
		}...)
	}

	content = append(content, []*sls.LogContent{
		&sls.LogContent{
			Key:   proto.String("server_name"),
			Value: proto.String(hostname),
		},
		&sls.LogContent{
			Key:   proto.String("level"),
			Value: proto.String("INFO"),
		}, &sls.LogContent{
			Key:   proto.String("time"),
			Value: proto.String(time.Now().Format("2006-01-02 15:04:05")),
		}, &sls.LogContent{
			Key:   proto.String("msg"),
			Value: proto.String(info.Msg),
		}, &sls.LogContent{
			Key:   proto.String("env"),
			Value: proto.String(info.Env),
		},
	}...)
	byteData, byteEr := json.Marshal(info.Data)
	if byteEr == nil {
		content = append(content, &sls.LogContent{
			Key:   proto.String("data"),
			Value: proto.String(string(byteData)),
		})
	}
	if info.TraceId != "" {
		content = append(content, &sls.LogContent{
			Key:   proto.String("trace_id"),
			Value: proto.String(info.TraceId),
		})

	} else {
		content = append(content, &sls.LogContent{
			Key:   proto.String("trace_id"),
			Value: proto.String(ctx.GetHeader("Trace_Id")),
		})
	}

	log := &sls.Log{
		Time:     proto.Uint32(uint32(time.Now().Unix())),
		Contents: content,
	}
	return conn.SendLog(log, logConfig)
}

//添加error日志到SLS（程序运行前，请初始化conn.InitSLSProducer）
func AddError(logError ErrEntity, ctx *gin.Context, project ...conn.LogProjectConf) error {
	logConfig := conn.GetDefaultErrorLogConfig()
	if len(project) > 0 {
		logConfig = project[0]
	}
	if len(logConfig.Project) == 0 || len(logConfig.Logstore) == 0 {
		return errors.New("")
	}

	content := []*sls.LogContent{}
	hostname, _ := os.Hostname()
	if ctx != nil {
		content = append(content, []*sls.LogContent{
			&sls.LogContent{
				Key:   proto.String("customer_id"),
				Value: proto.String(ctx.Param("customer_id")),
			},
			&sls.LogContent{
				Key:   proto.String("request_id"),
				Value: proto.String(ctx.GetHeader("Request_Id")),
			}, &sls.LogContent{
				Key:   proto.String("host"),
				Value: proto.String(ctx.Request.Host),
			}, &sls.LogContent{
				Key:   proto.String("uri"),
				Value: proto.String(ctx.Request.RequestURI),
			},
		}...)
	}

	content = append(content, []*sls.LogContent{
		&sls.LogContent{
			Key:   proto.String("server_name"),
			Value: proto.String(hostname),
		},
		&sls.LogContent{
			Key:   proto.String("level"),
			Value: proto.String("ERROR"),
		}, &sls.LogContent{
			Key:   proto.String("time"),
			Value: proto.String(time.Now().Format("2006-01-02 15:04:05")),
		}, &sls.LogContent{
			Key:   proto.String("msg"),
			Value: proto.String(logError.Msg),
		}, &sls.LogContent{
			Key:   proto.String("env"),
			Value: proto.String(logError.Env),
		},
	}...)
	byteData, byteEr := json.Marshal(logError.Data)
	if byteEr == nil {
		content = append(content, &sls.LogContent{
			Key:   proto.String("data"),
			Value: proto.String(string(byteData)),
		})
	}
	if logError.TraceId != "" {
		content = append(content, &sls.LogContent{
			Key:   proto.String("trace_id"),
			Value: proto.String(logError.TraceId),
		})

	} else {
		content = append(content, &sls.LogContent{
			Key:   proto.String("trace_id"),
			Value: proto.String(ctx.GetHeader("Trace_Id")),
		})
	}

	log := &sls.Log{
		Time:     proto.Uint32(uint32(time.Now().Unix())),
		Contents: content,
	}
	return conn.SendLog(log, logConfig)
}

//添加msg日志到SLS（程序运行前，请初始化conn.InitSLSProducer）
func AddLogStr(msg string, ctx *gin.Context) error {
	return AddInfo(InfoEntity{Msg: msg}, ctx)
}

//添加error日志到SLS（程序运行前，请初始化conn.InitSLSProducer）
func AddErrorMap(data map[string]string, project ...conn.LogProjectConf) error {
	logConfig := conn.GetDefaultErrorLogConfig()
	if len(project) > 0 {
		logConfig = project[0]
	}
	if len(logConfig.Project) == 0 || len(logConfig.Logstore) == 0 {
		return errors.New("")
	}
	return conn.SendMap(data, logConfig)
}
