## Redis 连接池

#### 配置

#### 导入

```go
import "github.com/jtzjtz/kit/conn/redis_pool"
```

#### 初始化

```go
func main()  {
   options := &redis_pool.Options{
       Host:               "127.0.0.1:6379",
       PassWord:           "",
       Database:           0,
       InitCap:            10,
       MaxCap:             100,
       IsWait:             true,
       IdleTimeout:        5 * time.Second,
       DialConnectTimeout: 5 * time.Second,
       DialReadTimeout:    5 * time.Second,
       DialWriteTimeout:   5 * time.Second,
   }
   p, err := redis_pool.NewRedisPool(options)

   redis := p.Get()

   r, err := redis.Do("ping")
   if err != nil {
      log.Println("[ERROR] ping redis fail", err)
   }
   
   fmt.Println(r)
}
```

#### 依赖

- redigo：github.com/gomodule/redigo

#### 压测

Docker 环境：

- 内存 1G 
- CPU 单核 

wrk 压测工具：

- -c 跟服务器建立并保持的TCP连接数量
- -d 压测时间
- -t 使用多少个线程进行压测


调用 redis：

```
wrk -c100 -d30s -t4 http://127.0.0.1:8001/


```

