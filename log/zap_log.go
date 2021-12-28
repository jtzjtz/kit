package log

import (
	"errors"
	"fmt"
	"github.com/jtzjtz/kit/file"
	"github.com/lestrrat/go-file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"time"
)

//全局日志实例（程序启动前调用InitLog(appName string, logDir string)）
var Logger *zap.SugaredLogger

func initLog(appName string, logDir string) (er error) {
	defer func() {
		if err := recover(); err != nil {
			er = err.(error)
		}
	}()
	if len(logDir) == 0 {
		logDir = "./"
	}
	if exists, _ := file.ExistsDir(logDir); !exists {
		if er = file.CreateDir(logDir); er != nil {
			return er
		}
	}
	if exists, _ := file.ExistsDir(logDir + "/backup"); !exists {
		if er = file.CreateDir(logDir + "/backup"); er != nil {
			return er
		}
	}
	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "time",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	// 实现两个判断日志等级的interface
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel
	})
	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel
	})

	// 获取 info、error日志文件的io.Writer 抽象 getWriter() 在下方实现
	accessWriter := getWriter(appName, "access", logDir)
	infoWriter := getWriter(appName, "info", logDir)
	errorWriter := getWriter(appName, "error", logDir)

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(accessWriter), debugLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
	)

	log := zap.New(core) // 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数, 有点小坑
	Logger = log.Sugar()
	if Logger == nil {
		return errors.New("初始化日志实例失败")
	}
	return nil
}

func getWriter(appName, logType, logDir string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每天分割一次日志
	hook, err := rotatelogs.New(
		logDir+"/backup/"+appName+"-"+logType+"%Y%m%d.log",
		rotatelogs.WithLinkName(fmt.Sprintf("%s/%s-%s.log", logDir, appName, logType)),
		rotatelogs.WithMaxAge(time.Hour*24*365*2),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	if err != nil {
		panic(err)
	}
	//ch :=make( chan os.Signal)
	//signal.Notify(ch, syscall.SIGHUP)
	//
	//go func(ch chan os.Signal) {
	//	<-ch
	//	hook.Rotate()
	//}(ch)
	return hook
}
