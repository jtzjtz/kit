package log

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jtzjtz/kit/http"
	"runtime"
	"strings"
	"time"
)

type InfoEntity struct {
	Data    interface{} `json:"data"` //额外保存的数据
	Msg     string      `json:"msg"`
	TraceId string      `json:"trace_id"`
	Env     string      `json:"env"` //环境
	Caller  string      `json:"caller"`
}
type ErrEntity struct {
	Data    interface{} `json:"data"` //额外保存的数据
	Msg     string      `json:"msg"`
	TraceId string      `json:"trace_id"`
	Env     string      `json:"env"` //环境
	Caller  string      `json:"caller"`
}

//初始化日志服务 handlerCnt：4 处理日志协程数量  logChannelCnt：1000 日志通道缓存大小

func InitLogger(appName, logDir string, handlerCnt, logChannelCnt int) error {
	er := initLog(appName, logDir)
	initHandler(handlerCnt, logChannelCnt)
	return er
}

//添加info日志到本地 ,ctx 没有传nil

func LogInfo(info InfoEntity, ctx *gin.Context) error {
	if Logger == nil || logChannel == nil {
		return errors.New("日志服务未初始化")
	}
	if info.TraceId == "" && ctx != nil {
		info.TraceId = http.GetTraceId(ctx)
	}
	info.Caller = caller(3)
	pushChannel(logEntity{LogLevel: INFOLEVEL, LogData: info})
	return nil
}
func LogInfoData(data interface{}, ctx *gin.Context) error {
	return LogInfo(InfoEntity{Data: data}, ctx)
}

//记录info msg ,ctx 没有传nil

func LogInfoMsg(msg string, ctx *gin.Context) error {
	return LogInfo(InfoEntity{Msg: msg}, ctx)
}
func LogInfoMsgData(msg string, data interface{}, ctx *gin.Context) error {
	return LogInfo(InfoEntity{Msg: msg, Data: data}, ctx)
}

//添加error日志到本地,ctx 没有传nil

func LogError(logError ErrEntity, ctx *gin.Context) error {
	if Logger == nil || logChannel == nil {
		return errors.New("日志服务未初始化")
	}
	if logError.TraceId == "" && ctx != nil {
		logError.TraceId = http.GetTraceId(ctx)

	}
	logError.Caller = caller(3)

	pushChannel(logEntity{LogLevel: ERRORLEVEL, LogData: logError})
	return nil
}

//记录error msg ,ctx 没有传nil

func LogErrorMsg(msg string, ctx *gin.Context) error {
	return LogError(ErrEntity{Msg: msg}, ctx)
}
func LogErrorData(data interface{}, ctx *gin.Context) error {
	return LogError(ErrEntity{Data: data}, ctx)
}
func LogErrorMsgData(msg string, data interface{}, ctx *gin.Context) error {
	return LogError(ErrEntity{Data: data, Msg: msg}, ctx)
}

//记录accesslog

func LogAccess(data map[string]string) error {
	if Logger == nil || logChannel == nil {
		return errors.New("日志服务未初始化")
	}
	//Logger.Debugw("access", "logdata", data)
	pushChannel(logEntity{LogLevel: ACCESSLEVEL, LogData: data})
	return nil
}

func pushChannel(entity logEntity) {
	c := time.After(time.Millisecond * 500)
	select {
	case <-c:
		println("push log channel timeout 500ms")
	case logChannel <- entity:
	}
}
func caller(skip int) string {
	var sb strings.Builder
	layer := 1
	for ; skip < 6; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if ok {
			pcName := runtime.FuncForPC(pc).Name() //获取函数名
			sb.WriteString(fmt.Sprintf("【%d】 %s [%s:%d]\n", layer, file, pcName, line))
			layer++
		}

	}
	return sb.String()

}
