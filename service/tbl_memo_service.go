package service

import (
	"github.com/gin-gonic/gin"
	"zanso/db"
	"zanso/model"
	"zanso/result"
	"zanso/util"
)

func GetMemo(c *gin.Context) {
	var tblMemo []model.TblMemo
	var req = model.TblMemo{}
	c.Bind(&req)
	if req.TokenKey != "liguanglong9527" {
		result.ErrSetMsg(c, "tokenKey错误")
		return
	}
	db.DB.Where("name LIKE ?", "%"+req.Name+"%").Find(&tblMemo)
	result.OkSetData(c, tblMemo)

}
func AddMemo(c *gin.Context) {
	var tblMemo = model.TblMemo{}
	c.ShouldBind(&tblMemo)
	tblMemo.CreatedAt = util.GetTime()
	tblMemo.UpdatedAt = util.GetTime()
	if tblMemo.TokenKey != "liguanglong9527" {
		result.ErrSetMsg(c, "tokenKey错误")
		return
	}

	db.DB.Save(&tblMemo)
	result.Ok(c)
}
func UpdateMemo(c *gin.Context) {
	var tblMemo = model.TblMemo{}
	c.Bind(&tblMemo)
	tblMemo.CreatedAt = util.GetTime()
	tblMemo.UpdatedAt = util.GetTime()
	db.DB.Save(&tblMemo)
	result.Ok(c)
}
func DelMemo(c *gin.Context) {
	var tblMemo = model.TblMemo{}
	c.Bind(&tblMemo)
	if tblMemo.TokenKey != "liguanglong9527" {
		result.ErrSetMsg(c, "tokenKey错误")
		return
	}
	db.DB.Delete(&tblMemo)
	result.Ok(c)
}
