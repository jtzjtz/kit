package access_log

/*日志采集*/
import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"github.com/jtzjtz/kit/http"
	"github.com/jtzjtz/kit/log"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var accessChannel chan *accessResponse

const (
	LOG_ROOT_PATH = "/data/accesslog/"
	LOG_EXTENSION = ".log"
)

type result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
type accessResponse struct {
	Ctx             *gin.Context
	RequestBody     string
	ResponseBody    string
	RequestTime     string
	RequestTimeUnix string
	HandleCount     int
}
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

type accessLogParams struct {
	channelBufferNum    int32  //channel 缓存区
	goroutineHandleNum  int32  //goruntine 协程启动数量,消费处理log日志信息
	logFullRes          bool   //是否记录全部响应数据
	timeOutSendMail     int32  //发送邮件的过期时间
	noNeedRouteSendMail string //不需要发送邮件地址路由地址
}

//设置默认值
var ap = accessLogParams{
	channelBufferNum:    5,
	goroutineHandleNum:  1,
	logFullRes:          false,
	timeOutSendMail:     0,
	noNeedRouteSendMail: "",
}

//设置参数
func SetAp(channelBufferNum, goroutineHandleNum int32, logFullRes bool) {
	ap.channelBufferNum = channelBufferNum
	ap.goroutineHandleNum = goroutineHandleNum
	ap.logFullRes = logFullRes
}

//设置请求超时处理参数
func SetTimeOutAp(timeOutSendMail int32, noNeedRouteSendMail string) {
	ap.timeOutSendMail = timeOutSendMail
	ap.noNeedRouteSendMail = noNeedRouteSendMail
}

//异步处理日志信息
func handleAccessLog(env, serviceName string, acsChannel chan *accessResponse) {
	for access := range acsChannel {
		ctx := access.Ctx

		if strings.Contains(ctx.Request.RequestURI, "ping") == true || strings.Contains(ctx.Request.RequestURI, "favicon") == true || ctx.Request.RequestURI == "/" {
			//过滤不需要记日志的url地址
			continue
		}
		scmResCode := 0
		scmResMsg := ""
		responseBody := access.ResponseBody
		if responseBody != "" {
			res := result{}
			err := json.Unmarshal([]byte(responseBody), &res)
			if err == nil {
				scmResCode = res.Code
				if ap.logFullRes == true {
					scmResMsg = responseBody
				} else {
					scmResMsg = res.Msg
				}
			}
		}
		requestBody := access.RequestBody
		requestTime, _ := strconv.ParseInt(access.RequestTimeUnix, 10, 64)

		requestUri := ctx.Request.RequestURI
		costTimeStr := fmt.Sprintf("%.3f", float64(time.Now().UnixNano()-requestTime)/1e9)
		//costTime, _ := strconv.ParseFloat(costTimeStr, 64)
		//if ap.timeOutSendMail > 0 && int32(costTime) > ap.timeOutSendMail && strings.Contains(ap.noNeedRouteSendMail, requestUri) == false {
		//	emailContent := fmt.Sprintf("请求地址：%s 请求方式：%s <br/>请求时间：%s<br/>请求耗时：%#v<br/> 请求参数:%s<br/> 响应状态码:%d<br/> 响应数据:%s<br/>", ctx.Request.RequestURI, ctx.Request.Method, access.RequestTime, costTime, requestBody, ctx.Writer.Status(), scmResMsg)
		//	println(emailContent) //发送邮件通知
		//}
		accessMap := make(map[string]string)

		accessMap["X-App-Info"] = ctx.GetHeader("X-App-Info")
		accessMap["X-SESSION-ID"] = ctx.GetHeader("X-SESSION-ID")

		accessMap["request_cost"] = costTimeStr //请求消耗时间
		accessMap["method"] = ctx.Request.Method
		if ctx.GetHeader("Authorization") != "" {
			accessMap["token"] = ctx.GetHeader("Authorization")
		} else {
			accessMap["token"] = ctx.GetHeader("X-Token")

		}
		accessMap["app_ver"] = ctx.GetHeader("X-App-Version")
		accessMap["request_length"] = strconv.Itoa(int(ctx.Request.ContentLength))
		accessMap["content_type"] = ctx.ContentType()
		accessMap["client_ip"] = ctx.GetHeader("X-Real-Ip")
		accessMap["request_timestamp"] = access.RequestTime
		accessMap["request_body"] = requestBody // post get 的参数
		accessMap["host"] = ctx.Request.Host
		accessMap["service_name"] = serviceName //微服务名字
		hostName, _ := os.Hostname()
		accessMap["server_name"] = hostName //当前服务器hostname
		accessMap["Forwarded-Host"] = ctx.GetHeader("X-Forwarded-For")
		accessMap["request_uri"] = requestUri
		accessMap["response_status"] = strconv.Itoa(ctx.Writer.Status()) //http status
		accessMap["user_agent"] = ctx.Request.UserAgent()                //用户端类型
		accessMap["http_referer"] = ctx.Request.Referer()
		accessMap["response_code"] = strconv.Itoa(scmResCode) //业务code
		accessMap["response_msg"] = scmResMsg                 //返回的数据
		accessMap["traceid"] = http.GetTraceId(ctx)
		accessMap["customerid"] = ctx.GetHeader("CustomerId")
		accessMap["env"] = env
		log.LogAccess(accessMap)

	}
}

//记录日志中间件
func AccessLogMiddleware(env, projectName string) gin.HandlerFunc {
	accessChannel = make(chan *accessResponse, ap.channelBufferNum)
	for i := 0; i < int(ap.goroutineHandleNum); i++ {
		go func(i int) {
			handleAccessLog(env, projectName, accessChannel)
		}(i)
	}

	return func(ctx *gin.Context) {
		requestTime := time.Now().String()
		requestTimeUnix := strconv.FormatInt(time.Now().UnixNano(), 10)
		requestBody := ""
		if ctx.Request.Method == "POST" {
			requestBodyByte, err := ctx.GetRawData()
			if err == nil {
				requestBody = string(requestBodyByte)
				ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBodyByte)) //把读过的字节流重新放到body
			}
		}
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw
		ctx.Next()
		accessChannel <- &accessResponse{Ctx: ctx, RequestBody: requestBody, ResponseBody: blw.body.String(), RequestTime: requestTime, RequestTimeUnix: requestTimeUnix}
	}
}
