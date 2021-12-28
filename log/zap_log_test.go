package log

import (
	"testing"
	"time"
)

func initLogInstance() error {

	return InitLogger("kit", "/Users/jtz/logs", 4, 100)
}
func TestAddlog(t *testing.T) {
	if er := initLogInstance(); er != nil {
		t.Error(er.Error())
	}
	LogInfo(InfoEntity{Msg: "adf", TraceId: "32342342", Env: "testc"}, nil)

	for {
		//LogInfo(entity.LogInfo{Msg: "adf", TraceId: "32342342", Env: "testc"}, nil)
		LogError(ErrEntity{Msg: "sss"}, nil)
		time.Sleep(time.Second * 10)
	}

}
