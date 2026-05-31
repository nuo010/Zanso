package service

import (
	"github.com/gin-gonic/gin"
	"zanso/db"
	"zanso/model"
	"zanso/result"
)

func GetQuotation(c *gin.Context) {
	var quotation model.TblQuotation // 注意拼写：Quotation

	// 使用 RAW SQL 实现随机取一条（兼容性好）
	err := db.DB.Order("RAND()").Limit(1).Find(&quotation).Error
	// 如果是 PostgreSQL，把 RAND() 改成 RANDOM()

	if err != nil {
		result.ErrSetMsg(c, "查询错误")
		return
	}

	result.OkSetData(c, quotation)
}
