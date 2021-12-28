package sign

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

func CreateSn(data map[string]interface{}, secret string, isDate bool) string {
	var keys []string
	var isFirst = true
	var str = ""
	for k, _ := range data {
		if k != "token" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, v := range keys {
		if isFirst {
			str = fmt.Sprintf("%v=%v", v, data[v])
			isFirst = false
		} else {
			str = str + fmt.Sprintf("&%v=%v", v, data[v])
		}
	}
	str = str + secret
	if isDate {
		str = str + time.Now().Format("2006-01-02")
	}
	token := strings.ToLower(ToMd5(str))
	return token
}

func CheckSn(data url.Values, sn string, secret string, isDate bool) bool {
	if sn == "" {
		return false
	}
	if strings.Compare(sn, CreateSnByFormData(data, secret, isDate)) == 0 {
		return true
	}
	return false
}

func CreateSnByFormData(data url.Values, secret string, isDate bool) string {
	var keys []string
	var isFirst = true
	var str = ""

	for k, _ := range data {
		if k != "token" {
			keys = append(keys, k)

		}

	}
	sort.Strings(keys)

	for _, v := range keys {
		if isFirst {

			str = fmt.Sprintf("%v=%v", v, data.Get(v))
			isFirst = false
		} else {
			str = str + fmt.Sprintf("&%v=%v", v, data.Get(v))

		}
	}
	str = str + secret
	if isDate {
		str = str + time.Now().Format("2006-01-02")
	}

	return strings.ToLower(ToMd5(str))
}
