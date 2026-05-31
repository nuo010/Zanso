package router

import (
	"singo/handler"
	"singo/service"

	"github.com/gin-gonic/gin"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	gin.ForceConsoleColor()
	r := gin.Default()

	// 跨域配置
	r.Use(handler.Cors())
	r.Use(handler.CustomLogger())
	r.Use(gin.Recovery())
	// 路由
	api := r.Group("/api")
	{
		api.POST("getMemo", service.GetMemo)
		api.POST("addMemo", service.AddMemo)
		api.POST("delMemo", service.DelMemo)
		api.POST("getTask", service.GetTask)
		api.POST("addTask", service.AddTask)
		api.POST("upDateTask", service.UpDateTask)
		api.POST("delTask", service.DelTask)
		api.POST("uploadFile", service.UploadFile)
		api.POST("uploadFile2", service.UploadFile2)
		api.POST("getFileList", service.GetTblFileList)
		api.POST("getEmailList", service.GetEmail)
		api.POST("addWeightRecord", service.AddWeightRecord)
		api.POST("getWeightStatistics", service.GetWeightStatistics)
		api.POST("getWeightRecordList", service.GetWeightRecordList)
		api.POST("getQuotation", service.GetQuotation)
		api.POST("GetWeightRecord", service.GetWeightRecord)
		api.POST("GetWeightRecordPage", service.GetWeightRecordPage)
		api.GET("getdy", service.GetUrl)
	}
	return r
}
