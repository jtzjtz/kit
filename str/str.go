package str

import "fmt"

//反转字符串
func ReverseString(s string) string {
	runes := []rune(s)

	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

//map转string
func MapToString(m map[string]interface{}, spilt_str string) string {
	str := ""
	for k, v := range m {
		str += fmt.Sprintf("%s:%s %s", k, v, spilt_str)
	}
	return str
}
