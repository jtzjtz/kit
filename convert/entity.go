package convert

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

//模型间相互转换，json tag 需一致，支持弱类型解析
func EntityToEntity(oriEntity interface{}, toEntity interface{}) error {
	config := &mapstructure.DecoderConfig{
		TagName:          "json",
		Metadata:         nil,
		Result:           &toEntity,
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	decoderErr := decoder.Decode(oriEntity)
	if decoderErr != nil {
		return err
	}
	return nil
}

//convert to map with empty field deleted
func EntityToMap(entitySrc interface{}) map[string]interface{} {
	retMap := make(map[string]interface{})
	oriJson, err := json.Marshal(entitySrc)
	if err == nil {
		_ = json.Unmarshal(oriJson, &retMap)

	}
	return retMap
}

func EntityToMapWithEmpty(entitySrc interface{}, emptyFields []string) map[string]interface{} {
	mapResult := EntityToMap(entitySrc)
	emptyValueMap := MapToMapWithEmpty(entitySrc, emptyFields)

	for k, v := range emptyValueMap {
		mapResult[k] = v
	}

	return mapResult
}

func MapToMapWithEmpty(mapEntity interface{}, emptyFields []string) map[string]interface{} {
	var reflectValue reflect.Value
	var reflectValueSrc reflect.Value
	mapResult := make(map[string]interface{})

	if len(emptyFields) == 0 {
		return mapResult
	}
	if reflect.TypeOf(mapEntity).Kind() == reflect.Ptr {
		reflectValueSrc = reflect.ValueOf(mapEntity)
		if reflectValueSrc.IsNil() {
			return mapResult
		}
		reflectValue = reflectValueSrc.Elem()
		for _, fieldName := range emptyFields {
			if field := reflectValue.FieldByName(fieldName); field.IsValid() {
				value := field.Interface()
				mapResult[fieldName] = value
			}
		}

	} else if reflect.TypeOf(mapEntity).Kind() == reflect.Struct {
		reflectValue = reflect.ValueOf(mapEntity)
		for _, fieldName := range emptyFields {
			if field := reflectValue.FieldByName(fieldName); field.IsValid() {
				value := field.Interface()
				mapResult[fieldName] = value
			}
		}

	}
	return mapResult
}

// 判断变量是否为空
func Empty(val interface{}) bool {
	if val == nil {
		return true
	}

	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}
