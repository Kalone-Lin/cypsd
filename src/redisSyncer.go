package src

import "fmt"

func SolveSyncerInfo(action string) {
	for i := 0; i >= 0; i += Page {
		//1.查询数据
		var syncerInfoArr = make([]*Record, 0)
		QuerySyncerInfo(&syncerInfoArr, Page, i)
		if len(syncerInfoArr) == 0 {
			return
		}
		if action == "query" {
			fmt.Printf("query syncer info result: %v \n", syncerInfoArr)
			return
		}
		//2.判断是否需要加密，加密后更新数据
		UpdateSyncerInfo(&syncerInfoArr)
		fmt.Printf("[REDIS_SYNCER]===============已更新数量===============: %v \n", i+len(syncerInfoArr))
	}

}

func QuerySyncerInfo(syncerArr *[]*Record, limit int, offset int) error {
	queryStr := "select id,srcPassword,dstPassword from redis_syncer limit ? offset ?"
	stmt, err := db.Prepare(queryStr)
	if err != nil {
		fmt.Println("预处理失败,err", err)
		return err
	}
	defer stmt.Close()
	rows, err := stmt.Query(limit, offset)
	if err != nil {
		fmt.Println("查询失败,err", err)
		return err
	}
	defer rows.Close() //关闭连接
	for rows.Next() {
		record := &Record{}
		if err := rows.Scan(&record.Id, &record.SrcPassword, &record.DstPassword); err != nil {
			fmt.Printf("查询失败,err: %v \n", err)
			return err
		}
		*syncerArr = append(*syncerArr, record)
	}
	return nil
}

func UpdateSyncerInfo(syncerArr *[]*Record) error {
	for _, record := range *syncerArr {
		_, err := DecodeRedisPassword(record.SrcPassword)
		flag := false
		if err == nil {
			//未加密或密码为空
			//fmt.Printf("syncer src 已加密或密码为空===: %s, syncerId:%d \n", record.SrcPassword, record.Id)
		} else {
			srcEnPsd, err := EncodeRedisPassword(record.SrcPassword)
			if err != nil {
				return err
			}
			record.SrcPassword = srcEnPsd
			flag = true
		}

		_, err = DecodeRedisPassword(record.DstPassword)
		if err == nil {
			//未加密或密码为空
			//fmt.Printf("syncer dst 已加密或密码为空===: %s, syncerId:%d \n", record.DstPassword, record.Id)
		} else {
			dstEnPsd, err := EncodeRedisPassword(record.DstPassword)
			if err != nil {
				return err
			}
			record.DstPassword = dstEnPsd
			flag = true
		}

		if flag {
			sqlStr := "update redis_syncer set srcPassword=?,dstPassword=? where id=?"
			stmt, err := db.Prepare(sqlStr)
			if err != nil {
				return err
			}
			defer stmt.Close()
			if _, err := stmt.Exec(record.SrcPassword, record.DstPassword, record.Id); err != nil {
				fmt.Println("update syncer record err:%v \n", err)
				return err
			}
		}

	}
	return nil
}
