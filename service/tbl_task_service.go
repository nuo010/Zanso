package service

import (
	"github.com/gin-gonic/gin"
	"zanso/db"
	"zanso/model"
	"zanso/result"
	"zanso/util"
)

func GetTask(c *gin.Context) {
	var req = model.TblTask{}
	c.Bind(&req)
	//allstudentjson, _ := json.Marshal(req)
	//util.Log().Println(string(allstudentjson))

	var tblTask []model.TblTask
	db.DB.Where("complete = ?", req.Complete).Order("creat_time desc").Find(&tblTask)
	result.OkSetData(c, tblTask)
}

func AddTask(c *gin.Context) {
	var task = model.TblTask{}
	c.Bind(&task)
	task.CreatTime = util.GetTime()
	task.UpdateTime = util.GetTime()
	task.Complete = "0"
	db.DB.Create(&task)
	result.Ok(c)

}
func UpDateTask(c *gin.Context) {
	var task = model.TblTask{}
	c.Bind(&task)
	db.DB.Model(&task).Update("complete", "1")
	result.Ok(c)
}
func DelTask(c *gin.Context) {
	var task = model.TblTask{}
	c.Bind(&task)
	db.DB.Delete(&model.TblTask{ID: task.ID})
	result.Ok(c)
}
