package service

import (
	"context"
	"encoding/json"
	"fmt"
	"singo/db"
	"strings"
	"time"
)

func GetFilebeat() {
	fmt.Println("GetFilebeat")
	ctx := context.Background()
	// 每次拉取的数据数量
	batchSize := int64(100)
	listKey := "filebeat"

	for {
		// 获取列表长度
		listLen, err := db.RedisClient.LLen(ctx, listKey).Result()
		if err != nil {
			fmt.Println("获取list长度错误:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if listLen == 0 {
			fmt.Println("list长度为空，等待中...")
			time.Sleep(5 * time.Second)
			continue
		}

		// 计算要拉取的范围
		start := int64(0)
		end := start + batchSize - 1
		if end >= listLen {
			end = listLen - 1
		}

		// 使用 LRANGE 获取批量数据
		data, err := db.RedisClient.LRange(ctx, listKey, start, end).Result()
		if err != nil {
			fmt.Println("批量获取数据错误:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// 消费数据
		for _, item := range data {
			if strings.Contains(item, "single") {
				var result map[string]interface{}
				err := json.Unmarshal([]byte(item), &result)
				if err != nil {
					fmt.Println("Error:", err)
					continue
				}
				i, ok := result["message"]
				if !ok {
					continue
				}
				var innerData map[string]interface{}
				if err := json.Unmarshal([]byte(i.(string)), &innerData); err != nil {
					continue
				}
				fmt.Println("数据内容:", innerData["log"])
			}
		}

		// 删除已消费的数据
		_, err = db.RedisClient.LTrim(ctx, listKey, end+1, -1).Result()
		if err != nil {
			fmt.Println("删除list数据错误:", err)
			time.Sleep(5 * time.Second)
		}
	}
}
