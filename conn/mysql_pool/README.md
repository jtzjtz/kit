## MySQL 连接池

#### 配置

- InitCap 初始化数量
- MaxCap 最大连接数
- Debug 是否开启 Debug
- User 用户名
- Pass 密码
- Host IP 地址
- Port 端口
- DataBase 数据库名称

#### 导入

```go
import "github.com/jtzjtz/kit/conn/mysql_pool"
```

#### 初始化

```go
func main() {
    options := &mysql_pool.Options{
        InitCap:      50,
        MaxCap:       100,
        IsDebug:      true,
        User:         "user",
        Pass:         "***",
        Host:         "127.0.0.1",
        Port:         "3306",
        DataBase:     "dbname",
    }
    
    conn, err := mysql_pool.NewMySqlPool(options)
    
    if err != nil {
        log.Printf("%#v\n", err)
        return
    }
    
    if conn == nil {
        log.Printf("conn= %#v\n", conn)
        return
    }
    
    //todo
    //conn.DoSomething()
}
```

#### 依赖

- gorm：github.com/jinzhu/gorm

#### 压测

Docker 环境：

- 内存 1G 
- CPU 单核 

wrk 压测工具：

- -c 跟服务器建立并保持的TCP连接数量
- -d 压测时间
- -t 使用多少个线程进行压测


调用 MySQL 链接：

```
wrk -c100 -d30s -t4 http://127.0.0.1:8001/test/mysql_pool


```

