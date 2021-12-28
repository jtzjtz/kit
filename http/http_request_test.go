package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestPostJson(t *testing.T) {
	ctx := new(gin.Context)
	ctx.Request = new(http.Request)
	ctx.Request.Header = http.Header{}
	traceid := strconv.Itoa(int(time.Now().UnixNano()))
	ctx.Request.Header.Set(TRACEID, traceid)
	data := make(map[string]interface{})
	data["productSn"] = "SN163029560033161249194"
	extparams := new(ExtParams)
	extparams.Headers = map[string]string{}
	extparams.Headers["Authorization"] = "Bearer dd3pjjsjss"
	resbody, httpstatus, er := PostJson("https://api.ttjianbao.com/api/mall/product/getProductDetail", data, extparams, ctx)
	if er != nil || httpstatus != 200 || len(resbody) == 0 {
		t.Error(er)
	}
	println(traceid)
	println(resbody)
}
