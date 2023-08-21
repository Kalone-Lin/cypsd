package src

import (
	"database/sql"
	"fmt"
)

var db *sql.DB
var Page = 200

func InitDb(dsn string) (err error) {
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("打开数据库失败,err:%v\n", err)
		return err
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("连接数据库失败,err:%v\n", err)
		return err
	}
	return nil
}
