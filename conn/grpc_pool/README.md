## gRPC 连接池

#### 配置

- InitTargets 服务地址
- InitCap 初始化容量
- MaxCap 最大容量
- DialTimeout 连接超时时长
- IdleTimeout 空闲超时时长

#### 导入

```go
import "github.com/jtzjtz/kit/conn/grpc_pool"
```
#### 初始化

```go
func main() {
    options := &grpc_pool.Options{
        InitTargets:  []string{"127.0.0.1:8001"},
        InitCap:      50,
        MaxCap:       100,
        DialTimeout:  time.Second * 600,
        IdleTimeout:  time.Second * 30,
    }
    
    p, err := grpc_pool.NewGRPCPool(options, grpc.WithInsecure())
    
    if err != nil {
        log.Printf("%#v\n", err)
        return
    }
    
    if p == nil {
        log.Printf("p= %#v\n", p)
        return
    }
    
    defer p.Close()
    
    conn, err := p.Get()
    if err != nil {
        log.Printf("%#v\n", err)
        return
    }
    
    defer p.Put(conn)
    
    //todo
    //conn.DoSomething()
    
    log.Printf("len=%d\n", p.IdleCount())
}
```

#### 依赖

- Goole gRPC：google.golang.org/grpc

#### 压测

Docker 环境：

- 内存 1G 
- CPU 单核 

wrk 压测工具：

- -c 跟服务器建立并保持的TCP连接数量
- -d 压测时间
- -t 使用多少个线程进行压测


调用 helloWorld 服务：

```
wrk -c500 -d30s -t4 http://127.0.0.1:8001/test/grpc_pool


Running 30s test @ http://10.70.30.106:8080/hello
  4 threads and 500 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    81.82ms  135.89ms 720.77ms   94.59%
    Req/Sec     0.99k   452.95     2.18k    79.17%
  4807 requests in 6.87s, 352.08KB read
  Socket errors: connect 0, read 500, write 0, timeout 0
Requests/sec:    700.04
Transfer/sec:     51.27KB
```

