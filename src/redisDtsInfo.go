package src

import (
	"encoding/json"
	"fmt"
)

func SolveDtsInfo(action string) {
	for i := 0; i >= 0; i += Page {
		//1.查询数据
		var dtsInfoArr = make([]*DtsInfoRecord, 0)
		if err := QueryRedisDtsInfo(&dtsInfoArr, Page, i); err != nil {
			fmt.Printf("查询redis_dts_info数据库失败,err:%v\n", err)
			return
		}
		if len(dtsInfoArr) == 0 {
			return
		}
		if action == "query" {
			fmt.Printf("query dts info result: %v \n", dtsInfoArr)
			var sInfo []InsInfo
			json.Unmarshal(dtsInfoArr[0].SrcInfo_DB, &sInfo)
			fmt.Printf("query sInfo: %s", sInfo)
			return
		}
		//2.判断是否需要加密，加密后更新数据
		if err := UpdateDtsInfo(&dtsInfoArr); err != nil {
			fmt.Printf("加密更新失败: %v", err)
			return
		}
		fmt.Printf("[REDIS_DTS_INFO]===============已更新数量===============: %v \n", i+len(dtsInfoArr))
	}
}

func QueryRedisDtsInfo(dtsInfoArr *[]*DtsInfoRecord, limit int, offset int) error {
	queryStr := "select jobId,srcInfo,dstInfo,createTime,updateTime from redis_dts_info limit ? offset ?"
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
		record := &DtsInfoRecord{}
		if err := rows.Scan(&record.JobId, &record.SrcInfo_DB, &record.DstInfo_DB, &record.CreateTime, &record.UpdateTime); err != nil {
			fmt.Println("查询失败,err", err)
			return err
		}
		*dtsInfoArr = append(*dtsInfoArr, record)
	}
	return nil
}

func UpdateDtsInfo(dtsInfoArr *[]*DtsInfoRecord) error {
	for _, record := range *dtsInfoArr {
		flag := 0
		var srcInfo = make([]*InsInfo, 0)
		err := cypsd(&srcInfo, record.SrcInfo_DB, &flag)
		if err != nil {
			return err
		}
		var dstInfo = make([]*InsInfo, 0)
		err = cypsd(&dstInfo, record.DstInfo_DB, &flag)
		if err != nil {
			return err
		}
		if flag == 1 {
			maSrc := record.SrcInfo_DB
			if len(srcInfo) > 0 {
				maSrc, err = json.Marshal(srcInfo)
				if err != nil {
					return err
				}
			}
			maDst := record.DstInfo_DB
			if len(dstInfo) > 0 {
				maDst, err = json.Marshal(dstInfo)
				if err != nil {
					return err
				}
			}
			sqlStr := "update redis_dts_info set srcInfo=?,dstInfo=?,createTime=?,updateTime=? where jobId=?"
			stmt, err := db.Prepare(sqlStr)
			if err != nil {
				return err
			}
			defer stmt.Close()

			if _, err := stmt.Exec(maSrc, maDst, record.CreateTime, record.UpdateTime,
				record.JobId); err != nil {
				fmt.Println("update record err:%v \n", err)
				return err
			}
		}
	}
	return nil
}

func cypsd(insInfo *[]*InsInfo, infoByte []byte, flag *int) error {
	if err := json.Unmarshal(infoByte, &insInfo); err != nil {
		return nil
	}
	for _, info := range *insInfo {
		if len(info.Password) == 0 {
			continue
		}
		_, err := DecodeRedisPassword(info.Password)
		if err == nil {
			//fmt.Printf("已加密或密码为空===,SerialId:%s, psd:%s\n", info.SerialId, info.Password)
			continue
		}
		enPsd, err := EncodeRedisPassword(info.Password)
		if err != nil {
			//fmt.Printf("已加密或密码为空===,SerialId:%s,psd:%s\n", info.SerialId, enPsd)
			continue
		}
		info.Password = enPsd
		*flag = 1
	}
	return nil
}
