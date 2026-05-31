package service

import (
	"encoding/json"
	"errors"
	"log"
	"singo/db"
	"singo/model"
	"singo/result"
	"singo/util"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Series struct {
	Name string     `json:"name"`
	Data []*float64 `json:"data"` // Num 是 int 类型
}

func AddWeightRecord(c *gin.Context) {
	var req = model.TblWeightRecord{}
	err := c.Bind(&req)
	if err != nil {
		result.ErrSetMsg(c, "绑定参数失败")
		return
	}
	if req.UserId == "" {
		result.ErrSetMsg(c, "用户ID不能为空")
		return
	}
	if req.Num <= 0 {
		result.ErrSetMsg(c, "体重数据无效")
		return
	}

	req.ID = util.GetUuid()
	now := util.GetTime()

	year, month, day := now.Date()
	dateKey := time.Date(year, month, day, 0, 0, 0, 0, now.Location())

	// 5. 查询该用户今天是否已有记录
	var existing model.TblWeightRecord
	err = db.DB.Where("user_id = ? AND creat_time >= ?", req.UserId, dateKey).
		Order("creat_time DESC").
		First(&existing).Error

	if err == nil {
		// 找到了今天已有的记录，删除它
		delErr := db.DB.Delete(&existing).Error
		if delErr != nil {
			log.Printf("删除用户 %s 今日旧记录失败: %v", req.UserId, delErr)
			result.ErrSetMsg(c, "删除旧记录失败")
			return
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 数据库其他错误
		log.Printf("查询用户 %s 今日记录失败: %v", req.UserId, err)
		result.ErrSetMsg(c, "服务器错误")
		return
	}
	req.CreatTime = util.GetTime()
	dbRes := db.DB.Create(&req)
	if dbRes.Error != nil {
		result.ErrSetMsg(c, "插入数据库失败")
		return
	}
	result.Ok(c)
}

func GetWeightStatistics(c *gin.Context) {
	//var weightRecordList []model.TblWeightRecord

	var req = model.TblWeightRecord{}
	err := c.Bind(&req)
	if err != nil {
		result.ErrSetMsg(c, "参数错误")
		return
	}

	// 获取当前时间 & 计算本周一 00:00:00
	now := time.Now() // 假设是 2025-11-09 (Sunday) in CST (UTC+8)

	// 正确获取本地当天 00:00:00
	loc := now.Location()
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, loc)

	weekday := int(now.Weekday())        // Sunday = 0
	daysSinceMonday := (weekday + 6) % 7 // Sunday → 6
	thisMonday := today.AddDate(0, 0, -daysSinceMonday)
	lastMonday := thisMonday.AddDate(0, 0, -7)
	nextMonday := thisMonday.AddDate(0, 0, 7)

	// 查询过去两周的数据（从上周一 到 本周日）
	var records []model.TblWeightRecord
	err = db.DB.Where("user_id = ? AND creat_time >= ? AND creat_time < ?",
		req.UserId, lastMonday, nextMonday).
		Order("creat_time ASC").
		Find(&records).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 没有记录也合法
			records = []model.TblWeightRecord{}
		} else {
			log.Printf("查询体重记录失败: %v", err)
			result.ErrSetMsg(c, "查询失败")
			return
		}
	}
	// 初始化两个切片：上周和本周各7天，初始为0（表示无数据）
	lastWeekData := make([]*float64, 7) // 上周一 ~ 上周日
	thisWeekData := make([]*float64, 7) // 本周一 ~ 本周日

	// 按日期聚合：map[年月日] -> 当天所有 Num 值
	type dayKey struct{ Year, Month, Day int }
	dayValues := make(map[dayKey][]*float64)

	for _, r := range records {
		date := r.CreatTime
		key := dayKey{date.Year(), int(date.Month()), date.Day()}
		if r.Num == 0 {
			dayValues[key] = append(dayValues[key], nil)
		} else {
			dayValues[key] = append(dayValues[key], &r.Num)
		}

	}

	// 辅助函数：计算某天应填入的值（取最大值，也可改为平均值）
	//getValueForDay := func(values []*float64) *float64 {
	//	if len(values) == 0 {
	//		return nil // 前端可改为 null 或不显示
	//	}
	//	// 取最大值作为代表（也可取平均）
	//	max := values[0]
	//	for _, v := range values {
	//		if *v > max {
	//			max = v
	//		}
	//	}
	//	return &max
	//	// 平均值版本：
	//	// sum := 0
	//	// for _, v := range values { sum += v }
	//	// return sum / len(values)
	//}

	// 填充上周数据（7天完整）
	for i := 0; i < 7; i++ {
		date := lastMonday.AddDate(0, 0, i)
		key := dayKey{date.Year(), int(date.Month()), date.Day()}
		if vals, exists := dayValues[key]; exists {
			lastWeekData[i] = vals[0]
		}
		// 否则保持为 0
	}

	// 填充本周数据（未来日期不填）
	for i := 0; i < 7; i++ {
		date := thisMonday.AddDate(0, 0, i)
		if date.After(now) {
			thisWeekData[i] = nil // 未来日期
			continue
		}
		key := dayKey{date.Year(), int(date.Month()), date.Day()}
		if vals, exists := dayValues[key]; exists {
			thisWeekData[i] = vals[0]
		}
	}

	// 构造响应
	series := []Series{
		{
			Name: "本周",
			Data: thisWeekData,
		},
		{
			Name: "上周",
			Data: lastWeekData,
		},
	}
	jsonData, _ := json.Marshal(series)
	log.Printf("DEBUG JSON: %s", jsonData)
	//db.DB.Where("user_id = ?", req.UserId).Order("creat_time desc").Find(&weightRecordList)
	result.OkSetData(c, series)
}

func GetWeightRecordList(c *gin.Context) {
	var req = model.TblWeightRecord{}
	c.Bind(&req)
	var tblWeightRecord []model.TblWeightRecord
	err := db.DB.Raw(`
        SELECT * 
        FROM tbl_weight_record 
        WHERE DATE(creat_time) = CURDATE() 
        ORDER BY creat_time DESC
    `).Scan(&tblWeightRecord).Error

	if err != nil {
		// 可选：记录日志
		log.Printf("查询体重记录失败: %v", err)
		result.ErrSetMsg(c, "查询失败")
		return
	}
	result.OkSetData(c, tblWeightRecord)
}

func GetWeightRecord(c *gin.Context) {
	//var weightRecordList []model.TblWeightRecord

	var req = model.TblWeightRecord{}
	err := c.Bind(&req)
	if err != nil {
		result.ErrSetMsg(c, "参数错误")
		return
	}
	var tblWeightRecord []model.TblWeightRecord
	err = db.DB.Raw(`
        SELECT * 
        FROM tbl_weight_record 
        WHERE user_id = ?
        ORDER BY creat_time DESC
    `, req.UserId).Scan(&tblWeightRecord).Error

	if err != nil {
		// 可选：记录日志
		log.Printf("查询体重记录失败: %v", err)
		result.ErrSetMsg(c, "查询失败")
		return
	}
	result.OkSetData(c, tblWeightRecord)
}

func GetWeightRecordPage(c *gin.Context) {
	// 获取请求参数，包括页码和每页条数
	page := c.DefaultQuery("page", "1")           // 默认页码为1
	pageSize := c.DefaultQuery("page_size", "10") // 默认每页10条记录

	// 转换为整数
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		result.ErrSetMsg(c, "无效的页码")
		return
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt <= 0 {
		result.ErrSetMsg(c, "无效的每页条数")
		return
	}

	// 计算偏移量
	offset := (pageInt - 1) * pageSizeInt

	// 绑定请求参数
	var req = model.TblWeightRecord{}
	err = c.Bind(&req)
	if err != nil {
		result.ErrSetMsg(c, "参数错误")
		return
	}

	var tblWeightRecord []model.TblWeightRecord
	// 查询分页数据
	err = db.DB.Raw(`
        SELECT * 
        FROM tbl_weight_record 
        WHERE user_id = ?
        ORDER BY creat_time DESC
        LIMIT ? OFFSET ?
    `, req.UserId, pageSizeInt, offset).Scan(&tblWeightRecord).Error

	if err != nil {
		// 可选：记录日志
		log.Printf("查询体重记录失败: %v", err)
		result.ErrSetMsg(c, "查询失败")
		return
	}

	// 获取总记录数，用于分页
	var total int64
	err = db.DB.Raw(`
        SELECT COUNT(*) 
        FROM tbl_weight_record 
        WHERE user_id = ?
    `, req.UserId).Scan(&total).Error

	if err != nil {
		log.Printf("查询总记录数失败: %v", err)
		result.ErrSetMsg(c, "查询失败")
		return
	}

	// 返回分页数据
	result.OkSetData(c, map[string]interface{}{
		"total":   total,
		"records": tblWeightRecord,
	})
}
