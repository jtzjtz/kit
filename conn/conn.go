package conn

import (
	"errors"
	"fmt"
)

var connNeedClosed []func() = make([]func(), 0)

//存储需要关闭的conn ，func() 函数内根据各自的规则写关闭代码
func PutNewConn(closeFunc func()) {
	connNeedClosed = append(connNeedClosed, closeFunc)
}
func CreateNewInit(createFunc func() error, closeFunc func() error) (er error) {
	retryCount := 0
	defer func() {
		if err := recover(); err != nil {
			er = errors.New("创建初始化失败")
		}
	}()
retry:
	er = createFunc()
	if er != nil {
		if retryCount < 3 {
			retryCount++
			fmt.Printf("CreateNewInit failed retrycount=%s error:", retryCount, er.Error())
			goto retry
		} else {
			return er
		}
	}
	if closeFunc != nil {
		er = closeFunc()
	}
	return er

}

//关闭所有存储的conn，建议在main函数中调用 defer  conn.close()
func Close() (er error) {

	defer func() {
		if err := recover(); err != nil {
			println("释放conn报错")
			er = errors.New("释放conn报错")
		}
	}()
	for _, closeFunc := range connNeedClosed {
		if closeFunc != nil {
			closeFunc()
		}
	}
	return nil

}
