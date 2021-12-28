### 日志中间件
#### 参数
- apolloEnv 服务启动env
- projectName 项目名称 如api-oms
- emailTaskId DTQ处理邮件发送的任务ID
- notifyUser 异常通知收件人邮箱
- notifyCc  异常通知抄件人邮箱，多个用逗号隔开

#### 调用

```go
    import "github.com/gin-gonic/gin"
    import "github.com/jtzjtz/kit/middleware/access_log"

    func main() {
        r := gin.Default()
        access_log.SetAp(30,1,false)
        r.Use(access_log.AccessLogMiddleware(env,projectName )
    }
```
