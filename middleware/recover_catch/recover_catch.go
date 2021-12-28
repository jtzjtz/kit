package recover_catch

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jtzjtz/kit/log"
	"runtime/debug"
	"time"
)

type result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func RecoverCatchMiddleware(env, appName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				subject := fmt.Sprintf("【Error】项目：%s 环境：[%s] 异常捕获！", appName, env)
				reqUrl := c.Request.Method + "  " + c.Request.Host + c.Request.RequestURI
				reqUrl += " 请求body:" + getReqBody(c)
				bodyData := fmt.Sprintf("errmsg:%s logtime:%s url:%s useragent:%s clientip:%s debugStack:%s", err, time.Now().Format("2006/01/02 15:04:05.000"), reqUrl, c.Request.UserAgent(), c.GetHeader("X-Real-Ip"), string(debug.Stack()))
				errorData := log.ErrEntity{
					Msg:  subject,
					Data: bodyData,
					Env:  env,
				}
				log.LogError(errorData, c)
				c.Abort()
			}
		}()
		c.Next()
	}
}

//获取请求body
func getReqBody(ctx *gin.Context) string {
	requestBody := ""
	if ctx.Request.Method == "POST" {
		requestBodyByte, err := ctx.GetRawData()
		if err == nil {
			requestBody = string(requestBodyByte)
		}
	}
	return requestBody
}
