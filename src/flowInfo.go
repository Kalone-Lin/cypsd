package src

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

func SolveFlowInfo(action string) {
	for i := 0; i >= 0; i += Page {
		//1.查询数据
		var flowInfoArr = make([]*Flow, 0)
		QueryFlowInfo(&flowInfoArr, Page, i)
		if action == "query" {
			fmt.Printf("query flow info result: %v \n", flowInfoArr[0].Context)
			return
		}
		if len(flowInfoArr) == 0 {
			return
		}
		//2.判断是否需要加密，加密后更新数据
		if err := UpdateFlowInfo(&flowInfoArr); err != nil {
			fmt.Printf("flow加密更新失败: %v \n", flowInfoArr[0].Context)
			return
		}
		fmt.Printf("[FLOW]===============已更新数量===============: %d \n", i+len(flowInfoArr))
	}
}

func QueryFlowInfo(flowArr *[]*Flow, limit int, offset int) error {
	queryStr := "select id,context from flow limit ? offset ?"
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
		record := &Flow{}
		if err := rows.Scan(&record.Id, &record.Context); err != nil {
			fmt.Printf("查询失败,err: %v \n", err)
			return err
		}
		*flowArr = append(*flowArr, record)
	}
	return nil
}

func UpdateFlowInfo(flowArr *[]*Flow) error {
	for _, flow := range *flowArr {
		conMap := make(map[string]interface{})
		if err := json.Unmarshal(flow.Context, &conMap); err != nil {
			fmt.Printf("error------: %v \n", err)
			return err
		}

		dstInfoStr := conMap["dstInfo"]
		flag := 0
		dstRes, ok := dstInfoStr.([]interface{})
		dstInfoRes := make([]*DstInfo, 0)
		if dstInfoStr != nil && ok {
			for _, record := range dstRes {
				dst := &DstInfo{}
				if err := mapstructure.Decode(record, dst); err != nil {
					return err
				}
				_, err := DecodeRedisPassword(dst.Password)
				if err == nil {
					//fmt.Printf("flow dst已加密或密码为空======: %s, flowId:%d \n", dst.Password, flow.Id)
					dstInfoRes = append(dstInfoRes, dst)
					continue
				}
				redisPassword, err := EncodeRedisPassword(dst.Password)
				if err != nil {
					return err
				}
				dst.Password = redisPassword
				dstInfoRes = append(dstInfoRes, dst)
				flag = 1
			}
		}

		opts := conMap["opts"]
		optsRes := make([]*RedisOpt, 0)
		if opts != nil {
			for _, opt := range opts.([]interface{}) {
				optRes := &RedisOpt{}
				if err := mapstructure.Decode(opt, optRes); err != nil {
					return err
				}
				_, err := DecodeRedisPassword(optRes.Password)
				if err == nil {
					//fmt.Printf("flow opt已加密或密码为空======: %s , flowId:%d \n", password, flow.Id)
					optsRes = append(optsRes, optRes)
					continue
				}
				redisPassword, err := EncodeRedisPassword(optRes.Password)
				if err != nil {
					return err
				}
				optRes.Password = redisPassword
				optsRes = append(optsRes, optRes)
				flag = 1
			}
		}

		dstCcInfoStr := conMap["dstCcInfo"]
		dstCcInfo := &DstCcInfo{}
		if dstCcInfoStr != nil {
			if err := mapstructure.Decode(dstCcInfoStr, dstCcInfo); err != nil {
				return err
			}
			_, err := DecodeRedisPassword(dstCcInfo.RedisPassword)
			if err == nil {
				//fmt.Printf("flow redisPsd已加密或密码为空======: %s , flowId:%d \n", redisPsd, flow.Id)
			} else {
				redisEnPsd, err := EncodeRedisPassword(dstCcInfo.RedisPassword)
				if err != nil {
					return err
				}
				dstCcInfo.RedisPassword = redisEnPsd
			}

			_, err = DecodeRedisPassword(dstCcInfo.UserPassword)
			if err == nil {
				//fmt.Printf("flow userPsd已加密或密码为空======: %s , flowId:%d \n", userPsd, flow.Id)
			} else {
				userEnpsd, err := EncodeRedisPassword(dstCcInfo.UserPassword)
				if err != nil {
					return err
				}
				dstCcInfo.UserPassword = userEnpsd
			}

			_, err = DecodeRedisPassword(dstCcInfo.TendisPassword)
			if err == nil {
				//fmt.Printf("flow tendisPsd已加密或密码为空======: %s , flowId:%d \n", tendisPsd, flow.Id)
			} else {
				tendisEnpsd, err := EncodeRedisPassword(dstCcInfo.TendisPassword)
				if err != nil {
					return err
				}
				dstCcInfo.TendisPassword = tendisEnpsd
			}
			flag = 1
		}

		if flag == 1 {
			if ok {
				conMap["dstInfo"] = dstInfoRes
			}
			conMap["opts"] = optsRes
			conMap["dstCcInfo"] = dstCcInfo
			dstByte, err := json.Marshal(conMap)
			if err != nil {
				return err
			}
			sqlStr := "update flow set context=? where id=?"
			stmt, err := db.Prepare(sqlStr)
			if err != nil {
				return err
			}
			defer stmt.Close()
			if _, err := stmt.Exec(dstByte, flow.Id); err != nil {
				fmt.Println("update flow record err:%v \n", err)
				return err
			}
		}
	}
	return nil
}
