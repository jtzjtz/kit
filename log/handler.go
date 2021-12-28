package log

import (
	"sync"
)

type LogLevel int8

const (
	ACCESSLEVEL LogLevel = iota //accesslog级别

	INFOLEVEL //info日志级别

	ERRORLEVEL //error 日志级别
)

type logEntity struct {
	LogLevel LogLevel
	LogData  interface{}
}

var logChannel chan logEntity
var openFileGroup sync.Map

//初始化异步日志通道
func initHandler(handlerCnt, channelCnt int) {
	logChannel = make(chan logEntity, channelCnt)
	for i := 0; i < handlerCnt; i++ {
		go handleLog()
	}
}

//异步写日志方法
func handleLog() {
	defer func() {
		if err := recover(); err != nil {
			println("handle log err", err)
		}
	}()

	for logEntity := range logChannel {

		if Logger == nil {
			println("Logger 日志实例为空")

		}
		switch logEntity.LogLevel {

		case ACCESSLEVEL:
			Logger.Debugw("access", "logdata", logEntity.LogData)
			break
		case INFOLEVEL:
			Logger.Infow("info", "logdata", logEntity.LogData)
			break
		case ERRORLEVEL:
			Logger.Errorw("error", "logdata", logEntity.LogData)
			break
		default:
			println("logEntity.LogLevel 类型不匹配")
		}
	}
}
