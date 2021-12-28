package sql

import (
	"fmt"
	"strings"
	"time"
)

var replace = map[string]string{"\\": "\\\\", "'": `\'`, "\\0": "\\\\0", "\n": "\\n", "\r": "\\r", `"`: `\"`, "\x1a": "\\Z"}

func mysqlRealEscapeString(v interface{}) interface{} {
	vstr, ok := v.(string)
	if !ok {
		t, istime := v.(time.Time)
		if !istime {
			return v
		}

		vstr = t.Format("2006-01-02 15:04:05")
	}

	for key, escape := range replace {
		vstr = strings.Replace(vstr, key, escape, -1)
	}

	return fmt.Sprintf("'%s'", vstr)
}

func mysqlRealEscapeStringLike(v interface{}) interface{} {
	vstr, ok := v.(string)
	if !ok {
		return v
	}

	for key, escape := range replace {
		vstr = strings.Replace(vstr, key, escape, -1)
	}

	return fmt.Sprintf("%s", vstr)
}

// OR .
func OR(sqls ...Fielder) Fielder {
	strSQL := make([]string, len(sqls))
	for i, sql := range sqls {
		strSQL[i] = fmt.Sprint(sql)
	}

	return Field(fmt.Sprintf("(%s)", strings.Join(strSQL, " OR ")))
}

// Relation 表示字段与值的关系
type Relation interface {
	In(...interface{}) Fielder
	Like(interface{}) Fielder
	LLike(interface{}) Fielder
	RLike(interface{}) Fielder
	GT(interface{}) Fielder
	GTE(interface{}) Fielder
	LT(interface{}) Fielder
	LTE(interface{}) Fielder
	EQ(interface{}) Fielder
	NEQ(interface{}) Fielder
	String() string
}

// Fielder 字段接口
type Fielder interface {
	And(f string) Relation
}

var _ Relation = Field("")
var _ Fielder = Field("")

// Field Fielder 接口的实现
type Field string

// And 。
func (ff Field) And(f string) Relation {
	return Field(fmt.Sprintf("%s AND %s", ff, f))
}

// In .
func (ff Field) In(vs ...interface{}) Fielder {
	v := make([]string, len(vs))

	for index, item := range vs {
		v[index] = fmt.Sprintf("%v", mysqlRealEscapeString(item))
	}
	sql := fmt.Sprintf("%s IN (%s)", ff, strings.Join(v, ", "))
	sql = strings.TrimPrefix(sql, "1=1 AND")

	return Field(sql)
}

// Like .
func (ff Field) Like(v interface{}) Fielder {
	sql := fmt.Sprintf("%s LIKE '%%%s%%'", ff, mysqlRealEscapeStringLike(v))
	sql = strings.TrimPrefix(sql, "1=1 AND")
	return Field(sql)
}

// LLike .
func (ff Field) LLike(v interface{}) Fielder {
	sql := fmt.Sprintf("%s LIKE '%%%s'", ff, mysqlRealEscapeStringLike(v))
	sql = strings.TrimPrefix(sql, "1=1 AND")
	return Field(sql)
}

// RLike .
func (ff Field) RLike(v interface{}) Fielder {
	sql := fmt.Sprintf("%s LIKE '%s%%'", ff, mysqlRealEscapeStringLike(v))
	sql = strings.TrimPrefix(sql, "1=1 AND")
	return Field(sql)
}

// GT .
func (ff Field) GT(v interface{}) Fielder {
	sql := fmt.Sprintf("%s > %v", ff, mysqlRealEscapeString(v))
	sql = strings.TrimPrefix(sql, "1=1 AND")
	return Field(sql)
}

// GTE .
func (ff Field) GTE(v interface{}) Fielder {
	sql := fmt.Sprintf("%s >= %v", ff, mysqlRealEscapeString(v))
	sql = strings.TrimPrefix(sql, "1=1 AND")
	return Field(sql)
}

// LT .
func (ff Field) LT(v interface{}) Fielder {
	sql := fmt.Sprintf("%s < %v", ff, mysqlRealEscapeString(v))
	sql = strings.TrimPrefix(sql, "1=1 AND")
	return Field(sql)
}

// LTE .
func (ff Field) LTE(v interface{}) Fielder {
	sql := fmt.Sprintf("%s <= %v", ff, mysqlRealEscapeString(v))
	sql = strings.TrimPrefix(sql, "1=1 AND")
	return Field(sql)
}

// EQ .
func (ff Field) EQ(v interface{}) Fielder {
	sql := fmt.Sprintf("%s = %v", ff, mysqlRealEscapeString(v))
	sql = strings.TrimPrefix(sql, "1=1 AND")
	return Field(sql)
}

// NEQ .
func (ff Field) NEQ(v interface{}) Fielder {
	sql := fmt.Sprintf("%s != %v", ff, mysqlRealEscapeString(v))
	sql = strings.TrimPrefix(sql, "1=1 AND")
	return Field(sql)
}

func (ff Field) String() string {
	return string(ff)
}

var (
	// DefaultField .
	DefaultField Fielder = Field("1=1")
)
