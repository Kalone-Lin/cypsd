package main

import (
	"cypsd/src"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main_() {
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

func main__() {
	var array1 = []int{1, 2, 3}
	ans := permute(array1)
	fmt.Printf("%v \n", ans)
}

func permute(nums []int) (ans [][]int) {
	n := len(nums)
	path := make([]int, n)
	onPath := make([]bool, n)
	var dfs func(int)
	dfs = func(i int) {
		if i == n {
			ans = append(ans, append([]int(nil), path...))
			return
		}
		for j, on := range onPath {
			if !on {
				path[i] = nums[j]
				onPath[j] = true
				dfs(i + 1)
				onPath[j] = false
			}
		}
	}
	dfs(0)
	return
}

type ser struct {
	SerialId string
}

func main() {
	str := "{\"SerialId\":\"crs-9giiojif\"}"
	tmp := ser{}
	err := json.Unmarshal([]byte(str), &tmp)
	if err != nil {
		fmt.Printf("-----------goggogg-------------:", err)
	}
	fmt.Printf("-----------success-------------")
}
