package sql

import (
	"fmt"
	"testing"
	"time"
)

func Test1EQ1(t *testing.T) {
	// 1=1
	sql := DefaultField
	t.Log(sql)
}
func TestAnd(t *testing.T) {
	// id = '1' AND name LIKE '%zhangsan%' AND age > 19
	sql := DefaultField.And("id").EQ("1").And("name").Like("zhangsan").And("age").GT(19)
	t.Log(sql)
}

func TestIn(t *testing.T) {
	// id IN (1, 2, 3)
	sql := DefaultField.And("id").In(1, 2, 3)
	t.Log(sql)
}

func TestInStr(t *testing.T) {
	// id IN ('a', 'b', 'c')
	sql := DefaultField.And("id").In("a", "b", "c")
	t.Log(sql)
}

func TestSqlInject(t *testing.T) {
	// full_name LIKE '%\'%f\') union select 1,2,3,4,5,6,7,user(),9,10,11,12,13,14#%\')%'
	sql := DefaultField.And("full_name").Like("'%f') union select 1,2,3,4,5,6,7,user(),9,10,11,12,13,14#%')")
	t.Log(sql)
}

func TestOR(t *testing.T) {
	// (name LIKE 'abc%' AND sex > 20 OR  name LIKE 'bcd%' AND sex < 20) AND job = 'xxx'
	sql := OR(DefaultField.And("name").RLike("abc").And("sex").GT(20), DefaultField.And("name").RLike("bcd").And("sex").LT(20)).And("job").EQ("xxx")
	t.Log(sql)
}

func TestOR2(t *testing.T) {
	// name = 'zhangsan' AND ( age < 20 OR  age > 25) AND email LIKE '%xxx@qq.com%'
	sql := DefaultField.And("name").EQ("zhangsan").And(fmt.Sprint(OR(DefaultField.And("age").LT(20), DefaultField.And("age").GT(25)).And("email").Like("xxx@qq.com")))
	t.Log(sql)
}

func TestTime(t *testing.T) {
	// create_time = "2020-12-30 15:35:49"
	sql := DefaultField.And("create_time").EQ(time.Now())
	t.Log(sql)
}
