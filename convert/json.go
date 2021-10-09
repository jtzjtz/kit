package convert

import (
	"encoding/json"
	"github.com/jtzjtz/kit/error_n"
)

//实体转json
func JsonEncode(data interface{}) (jsonStr string) {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return ""
	} else {
		return string(dataByte)
	}
}

// json bytes 转实体
func JsonDecodeBytes(strBytes []byte, toEntity interface{}) error {
	if len(strBytes) == 0 {
		return error_n.NewError("jsonStr 不能为空")
	}
	var data interface{}
	err := json.Unmarshal(strBytes, &data)
	if err != nil {
		return err
	} else {
		return EntityToEntity(data, &toEntity)
	}
}

//json string 转实体
func JsonDecode(jsonStr string, toEntity interface{}) error {

	return JsonDecodeBytes([]byte(jsonStr), toEntity)

}
