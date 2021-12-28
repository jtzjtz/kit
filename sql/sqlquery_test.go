package sql

import (
	"testing"
	"time"
)

func TestQuery(t *testing.T) {
	sql, err := Query("")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sql)
}

func TestQuery1(t *testing.T) {
	sql, err := Query("id = ? AND name LIKE ? AND age > ?", "1", "%zhangsan%", 19)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sql)
}

func TestQuery2(t *testing.T) {
	// full_name LIKE '%\'%f\') union select 1,2,3,4,5,6,7,user(),9,10,11,12,13,14#%\')%'
	sql, err := Query("full_name LIKE ?", "'%f') union select 1,2,3,4,5,6,7,user(),9,10,11,12,13,14#%')")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sql)
}

func TestQuery3(t *testing.T) {
	// create_time = '2021-01-04 14:08:35'
	sql, err := Query("create_time = ?", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sql)
}

func TestQuery4(t *testing.T) {
	// name = 'zhangsan' AND ( age < 20 OR  age > 25) AND email LIKE '%xxx@qq.com%'
	sql, err := Query("name = ? AND ( age < ? OR  age > ?) AND email LIKE ?", "zhangsan", 20, int32(25), "%xxx@qq.com")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sql)
}

func TestQuery5(t *testing.T) {
	// name = 'zhangsan' AND ( age < 20 OR  age > 25) AND email LIKE '%xxx@qq.com%'
	sql, err := Query("price = ?", 20.11)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sql)
}

func TestQuery6(t *testing.T) {
	// full_name LIKE '%\'%f\') union select 1,2,3,4,5,6,7,user(),9,10,11,12,13,14#%\')%'
	sql, err := Query("full_name LIKE ?", []byte("'%f') union select 1,2,3,4,5,6,7,user(),9,10,11,12,13,14#%')"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sql)
}
