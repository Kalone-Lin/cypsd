package main

import (
	"cypsd/src"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var (
		dsn       = ""
		tableName = ""
		action    = ""
		inpage    = 0
	)
	if action == "quit" {
		return
	}
	fmt.Scan(&dsn, &tableName, &action, &inpage)
	if inpage > 0 {
		src.Page = inpage
	}
	fmt.Printf("dsn: %s, tableName: %s, action: %s \n", dsn, tableName, action)
	if err := src.InitDb(dsn); err != nil {
		fmt.Printf("初始化数据库失败,err:%v\n", err)
		return
	}
	switch tableName {
	case "redis_dts_info":
		src.SolveDtsInfo(action)
	case "flow":
		src.SolveFlowInfo(action)
	case "redis_syncer":
		src.SolveSyncerInfo(action)
	case "all":
		src.SolveDtsInfo(action)
		src.SolveFlowInfo(action)
		src.SolveSyncerInfo(action)
	default:

	}

	fmt.Printf("===============finished==================\n")
}
