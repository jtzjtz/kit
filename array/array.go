package array

import (
	"errors"
	"reflect"
)

//判断数组或切片中是否存在某元素
func Contains(find interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == find {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(find)).IsValid() {
			return true, nil
		}
	}

	return false, errors.New("not in array")
}
func In_Array(val interface{}, array interface{}) bool {
	exists := false
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				exists = true
			}
		}
	}
	return exists
}

//求数组的并集、交集、差集----int-array
//求并集
func UnionArray(slice1, slice2 []int32) []int32 {
	m := make(map[int32]int)
	for _, v := range slice1 {
		m[v]++
	}
	for _, v := range slice2 {
		times, _ := m[v]
		if times == 0 {
			slice1 = append(slice1, v)
		}
	}
	return slice1
}

//求交集
func IntersectArray(slice1, slice2 []int32) []int32 {
	m := make(map[int32]int)
	nn := make([]int32, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

//求差集 (slice1去掉slice2中存在的元素)
func DifferenceArray(slice1, slice2 []int32) []int32 {
	m := make(map[int32]int)
	nn := make([]int32, 0)
	inter := IntersectArray(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range slice1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}
