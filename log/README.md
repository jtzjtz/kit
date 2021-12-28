## 普通文件日志记录
#### 参数
- path 存储的文件路径
- fileName 保存的文件名称
- data 保存的数据（最终会以json的格式保存在文件中）
- isAsync 是否异步保存

#### 调用

```go
    import "github.com/jtzjtz/kit/log"
    import "github.com/jtzjtz/kit/conn"


    func init() {
    	
    }
    func main() {
       //如需记录日志异常通知
       res,err := log.AddLog(path string,fileName string, data map[string]interface{},isAsync bool) (bool, error)
       defer conn.Close()

    }
```