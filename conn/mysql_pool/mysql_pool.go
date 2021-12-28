package mysql_pool

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func NewMySqlPool(o *Options) (*gorm.DB, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}
	dbConnStr := getDbConnStr(o.User, o.Pass, o.DataBase, o.Host, o.Port)
	DB, err := gorm.Open("mysql", dbConnStr)
	if err != nil {
		fmt.Printf("mysql connection err=%v\n", err)
	}

	DB.SingularTable(true)
	DB.LogMode(o.IsDebug)

	DB.DB().SetMaxOpenConns(o.MaxCap)
	DB.DB().SetMaxIdleConns(o.InitCap)

	return DB, nil
}

func getDbConnStr(user, pass, db, host, port string) string {
	return user + ":" + pass + "@(" + host + ":" + port + ")/" + db + "?charset=utf8&parseTime=True&loc=Local"
}
