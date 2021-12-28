package http

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var TRACEID = "Eagleeye-Traceid"

type ExtParams struct {
	Headers    map[string]string //请求头
	RetryCount int               //失败重试次数
	TimeOut    time.Duration     //超时时间
}

//处理map到encode字符串
func handleMapToEncodeStr(params map[string]interface{}) string {
	paramStr := ""
	if params != nil && len(params) > 0 {
		i := 0
		for k, v := range params {
			val := fmt.Sprintf("%v", v)
			if i == 0 {
				paramStr += fmt.Sprintf("%v=%v", k, url.QueryEscape(val))
			} else {
				paramStr += fmt.Sprintf("&%v=%v", k, url.QueryEscape(val))
			}
			i++
		}
	}
	return paramStr
}

//发起GET请求
func GetRequest(ctx *gin.Context, reqUrl string, params map[string]interface{}, timeOut time.Duration, extParams ExtParams) (resBody string, resStatus int, resErr error) {
	if strings.Contains(reqUrl, "?") == false && params != nil && len(params) > 0 {
		reqUrl += "?"
	}
	reqUrl += handleMapToEncodeStr(params)
	return DoRequest(ctx, reqUrl, "", "GET", timeOut, extParams)
}

//发起POST请求
func PostRequest(ctx *gin.Context, reqUrl string, params map[string]interface{}, timeOut time.Duration, extParams ExtParams) (resBody string, resStatus int, resErr error) {
	if extParams.Headers == nil || len(extParams.Headers) == 0 {
		extParams.Headers = map[string]string{}
	}
	if _, ok := extParams.Headers["Content-Type"]; !ok {
		extParams.Headers["Content-Type"] = "application/x-www-form-urlencoded;charset=UTF-8"
	}
	return DoRequest(ctx, reqUrl, handleMapToEncodeStr(params), "POST", timeOut, extParams)
}

//POST发送json数据
func PostJsonRequest(ctx *gin.Context, reqUrl string, params map[string]interface{}, timeOut time.Duration, extParams ExtParams) (resBody string, resStatus int, resErr error) {
	if extParams.Headers == nil || len(extParams.Headers) == 0 {
		extParams.Headers = map[string]string{}
		extParams.Headers["Content-Type"] = "application/json;charset=UTF-8"
	} else {
		extParams.Headers["Content-Type"] = "application/json;charset=UTF-8"
	}
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return "", 0, errors.New("json解析错误:" + err.Error())
	}
	return DoRequest(ctx, reqUrl, string(jsonParams), "POST", timeOut, extParams)
}

//POST发送json数据
func PostJson(reqUrl string, params map[string]interface{}, extParams *ExtParams, ctx *gin.Context) (resBody string, resStatus int, resErr error) {
	var timeOut time.Duration
	if extParams != nil {
		if extParams.Headers == nil || len(extParams.Headers) == 0 {
			extParams.Headers = map[string]string{}
			extParams.Headers["Content-Type"] = "application/json;charset=UTF-8"
		} else {
			extParams.Headers["Content-Type"] = "application/json;charset=UTF-8"
		}
		timeOut = extParams.TimeOut
	}

	jsonParams, err := json.Marshal(params)
	if err != nil {
		return "", 0, errors.New("json解析错误:" + err.Error())
	}
	return DoRequest(ctx, reqUrl, string(jsonParams), "POST", timeOut, *extParams)
}

//发起Http请求
func DoRequest(ctx *gin.Context, reqUrl string, params string, method string, timeOut time.Duration, extParams ExtParams) (resBody string, resStatus int, resErr error) {
	reqTimeOut := timeOut
	retryCount := extParams.RetryCount
	if reqTimeOut <= 0 {
		reqTimeOut = 30 * time.Second //默认30秒超时
	}

	//添加常用http请求直接close链接，防止占用描述符链接资源
	if extParams.Headers == nil {
		extParams.Headers = map[string]string{
			"Connection": "close",
		}
	} else {
		if _, ok := extParams.Headers["Connection"]; !ok {
			extParams.Headers["Connection"] = "close"
		}
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		//MaxIdleConnsPerHost: 20,
		//DisableKeepAlives:false,
	}
	client := &http.Client{
		Timeout:   reqTimeOut,
		Transport: tr,
	}
	var res *http.Response
	var req *http.Request
	var err error
	for i := 0; i <= retryCount; i++ {
		req, err = http.NewRequest(method, reqUrl, strings.NewReader(params))
		if err != nil {
			if res != nil {
				_ = res.Body.Close()
			}
			continue
		}
		if ctx != nil {
			if ctx.GetHeader(TRACEID) != "" {
				req.Header.Set(TRACEID, ctx.GetHeader(TRACEID))
			}
		}
		if extParams.Headers != nil && len(extParams.Headers) > 0 {
			for k, v := range extParams.Headers {
				req.Header.Set(k, v)
			}
		}
		res, err = client.Do(req)
		if err != nil {
			if res != nil {
				_ = res.Body.Close()
			}
			continue
		} else {
			break
		}
	}
	statusCode := 0 //响应状态码
	if res != nil {
		statusCode = res.StatusCode
	}
	if err != nil {
		//请求异常
		return "", statusCode, err
	}

	defer func() {
		if res != nil {
			_ = res.Body.Close()
		}
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//解析数据异常
		return "", statusCode, err
	}

	return string(body), statusCode, nil
}
func GetTraceId(ctx *gin.Context) string {
	if len(ctx.GetHeader("TraceId")) > 0 {
		return ctx.GetHeader("TraceId")
	} else if len(ctx.GetHeader("Eagleeye-Traceid")) > 0 {
		return ctx.GetHeader("Eagleeye-Traceid")
	} else if len(ctx.GetHeader("X-B3-Traceid")) > 0 {
		return ctx.GetHeader("X-B3-Traceid")
	}
	return ""
}
