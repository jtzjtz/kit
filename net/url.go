package net

import (
	"github.com/jtzjtz/kit/array"
	"net/url"
	"reflect"
	"strings"
)

//map转url

func MapToUrl(m map[string]string) string {
	vs := url.Values{}
	for k, v := range m {
		vs.Add(k, v)
	}
	return vs.Encode()
}

//url转map

func UrlToMap(u string) map[string][]string {
	m, _ := url.ParseQuery(u)
	return m
}

//从post数据中取要更新为空的字段 返回值为 更新为空的entity字段名字

func GetEmptyEntityFieldFromPost(postData url.Values, entity interface{}) []string {
	var emptyFields, result []string
	for k, v := range postData {
		vStr := strings.Join(v, "")
		if vStr == "" {
			emptyFields = append(emptyFields, k)
		} else if vStr == "0" {
			emptyFields = append(emptyFields, k)
		}
	}
	getType := reflect.TypeOf(entity)

	// 获取方法字段
	// 1. 先获取interface的reflect.Type，然后通过NumField进行遍历
	// 2. 再通过reflect.Type的Field获取其Field
	// 3. 最后通过Field的Interface()得到对应的value
	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		//	value := getValue.Field(i).Interface()
		formTag := field.Tag.Get("json")
		tags := strings.Split(formTag, ",")
		if len(tags) > 0 {
			if isHave, err := array.Contains(tags[0], emptyFields); err == nil && isHave == true {
				result = append(result, field.Name)
			}
		}

	}

	return result
}

// 从url数据中获取 SQL 的 Sort 和 Where，返回值为 Sort,Where

func GetSortAndWhereFormUrl(entity interface{}, values url.Values) (map[string]string, string) {
	var sort = map[string]string{}
	var where = ""
	var fieldWhereKey []string
	var fieldSortKey []string
	var filedSortMap = map[string]string{}
	for k, v := range values {
		vStr := strings.Join(v, "")
		if k == "sort" { // sort
			if strings.Index(vStr, "|") != -1 {
				for _, sv := range strings.Split(vStr, "|") {
					if strings.Index(sv, ",") != -1 {
						sortOneArr := strings.Split(sv, ",")
						if len(sortOneArr) == 2 && (sortOneArr[1] == "asc" || sortOneArr[1] == "desc") {
							fieldSortKey = append(fieldSortKey, sortOneArr[0])
							filedSortMap[sortOneArr[0]] = sortOneArr[1]
						}
					}
				}
			} else {
				if strings.Index(vStr, ",") != -1 {
					sortOneArr := strings.Split(vStr, ",")
					if len(sortOneArr) == 2 && (sortOneArr[1] == "asc" || sortOneArr[1] == "desc") {
						fieldSortKey = append(fieldSortKey, sortOneArr[0])
						filedSortMap[sortOneArr[0]] = sortOneArr[1]
					}
				}
			}
		}

		if k != "sort" && k != "page" && k != "page_num" { // where
			fieldWhereKey = append(fieldWhereKey, k)
		}
	}

	v := reflect.TypeOf(entity)
	if v.NumField() > 0 {
		var isHave = true
		var err error
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			formTag := ""
			formTagStr := field.Tag.Get("form")
			formTagArr := strings.Split(formTagStr, ",")
			if len(formTagArr) > 0 {
				formTag = formTagArr[0]
			}
			// 拼接 OrderBy
			isHave, err = array.Contains(formTag, fieldSortKey)
			if isHave && err == nil {
				sort[formTag] = filedSortMap[formTag]
			}
			// 拼接 Where
			isHave, err = array.Contains(formTag, fieldWhereKey)
			if isHave && err == nil {
				kind := field.Type.Kind()
				formTagValueStr := values.Get(formTag)
				var formTagValueArr []string
				if strings.Index(formTagValueStr, "|") != -1 {
					// in
					formTagValueArr = strings.Split(formTagValueStr, "|in")
					if len(formTagValueArr) == 2 && formTagValueArr[0] != "" {
						var inArr []string
						for _, v := range strings.Split(formTagValueArr[0], ",") {
							inArr = append(inArr, v)
						}
						where += formTag + " IN ('" + strings.Join(inArr, "','") + "') AND "
					}
					// like
					formTagValueArr = strings.Split(formTagValueStr, "|like")
					if len(formTagValueArr) == 2 && formTagValueArr[0] != "" {
						where += formTag + " LIKE '%" + formTagValueArr[0] + "%' AND "
					}
					// llike
					formTagValueArr = strings.Split(formTagValueStr, "|llike")
					if len(formTagValueArr) == 2 && formTagValueArr[0] != "" {
						where += formTag + " LIKE '%" + formTagValueArr[0] + "' AND "
					}
					// rlike
					formTagValueArr = strings.Split(formTagValueStr, "|rlike")
					if len(formTagValueArr) == 2 && formTagValueArr[0] != "" {
						where += formTag + " LIKE '" + formTagValueArr[0] + "%' AND "
					}
					// gt
					formTagValueArr = strings.Split(formTagValueStr, "|gt")
					if len(formTagValueArr) == 2 && formTagValueArr[0] != "" {
						where += formTag + " > '" + formTagValueArr[0] + "' AND "
					}
					// gte
					formTagValueArr = strings.Split(formTagValueStr, "|gte")
					if len(formTagValueArr) == 2 && formTagValueArr[0] != "" {
						where += formTag + " >= '" + formTagValueArr[0] + "' AND "
					}
					// lt
					formTagValueArr = strings.Split(formTagValueStr, "|lt")
					if len(formTagValueArr) == 2 && formTagValueArr[0] != "" {
						where += formTag + " < '" + formTagValueArr[0] + "' AND "
					}
					// lte
					formTagValueArr = strings.Split(formTagValueStr, "|lte")
					if len(formTagValueArr) == 2 && formTagValueArr[0] != "" {
						where += formTag + " <= '" + formTagValueArr[0] + "' AND "
					}
					// neq
					formTagValueArr = strings.Split(formTagValueStr, "|neq")
					if len(formTagValueArr) == 2 && formTagValueArr[0] != "" {
						where += formTag + " != '" + formTagValueArr[0] + "' AND "
					}
				} else {
					if kind == reflect.Int || kind == reflect.Float32 || kind == reflect.Float64 {
						where += formTag + " = " + formTagValueStr + " AND "
					}
					if kind == reflect.String {
						where += formTag + " = '" + formTagValueStr + "' AND "
					}
				}
			}
		}
		if where != "" {
			where = where[0 : len(where)-4]
		}
	}
	return sort, where
}
