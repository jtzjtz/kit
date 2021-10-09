package convert

import (
	"fmt"
	"testing"
)

func TestJsonDecode(t *testing.T) {
	a := make(map[string]string)
	er := JsonDecode("{\"aa\":2,\"bb\":3}", a)

	if er != nil {
		t.Error(er)
	}
	fmt.Println(a)
}
