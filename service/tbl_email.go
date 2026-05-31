package service

import (
	"github.com/gin-gonic/gin"
	"singo/db"
	"singo/model"
	"singo/result"
)

func GetEmail(c *gin.Context) {
	var tblEmailList []model.TblEmail
	var req = model.TblEmail{}
	c.Bind(&req)
	db.DB.Where("del_flag = ?", "0").Order("creat_time desc").Find(&tblEmailList)
	result.OkSetData(c, tblEmailList)
}
