package log

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
	"time"
)

func TestLogInfo(t *testing.T) {

	InitLogger("kit", "/Users/jtz/logs", 4, 100)
	var info InfoEntity
	info.Msg = "测试info日志"
	info.Data = "infokk;k;"
	ctx := new(gin.Context)
	ctx.Request = new(http.Request)
	ctx.Request.Header = http.Header{}
	ctx.Request.Header.Set("TraceId", time.Now().String())
	LogInfo(info, ctx)
	time.Sleep(1 * time.Second)

}

func TestLogError(t *testing.T) {

	InitLogger("kit", "/Users/jtz/logs", 4, 100)
	var errorInfo ErrEntity
	errorInfo.Msg = "测试ERROR日志"
	errorInfo.TraceId = "3534535345345"
	//errorInfo.Data = nil

	LogError(errorInfo, nil)

	time.Sleep(1 * time.Second)
}
