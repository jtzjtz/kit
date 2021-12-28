package convert

import (
	"github.com/jtzjtz/kit/database"
	"reflect"
	"strings"
)

func getEqualCondition(name string, value interface{}) (database.SqlCondition, bool) {
	if name == "" || value == nil {
		return database.SqlCondition{}, false
	}

	return database.SqlCondition{
		QueryName: name,
		Predicate: database.SqlEqualPredicate,
		Value:     value,
	}, true
}

// map转数据库检索条件数组
func MapToSqlcondition(protoMegMap map[string]interface{}) []database.SqlCondition {
	conditions := []database.SqlCondition{}
	for name, value := range protoMegMap {
		if con, ok := getEqualCondition(name, value); ok {
			conditions = append(conditions, con)
		}
	}

	return conditions
}

func GetEntityPk(sct interface{}) (string, bool) {
	t := reflect.TypeOf(sct).Elem()
	for i := 0; i < t.NumField(); i++ {
		js := strings.Split(t.Field(i).Tag.Get("gorm"), ";")
		if len(js) > 1 && strings.TrimSpace(js[1]) == "primary_key" {
			cols := strings.Split(js[0], ":")
			return cols[1], true
		}
	}
	return "", false
}
