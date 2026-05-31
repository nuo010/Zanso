package router

import (
	"zanso/handler"
	"zanso/service"

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
	r.Static("/uploads", "./uploads")

	// 新平台路由
	api := r.Group("/api")
	platform := api.Group("/platform")
	{
		platform.POST("/merchants", service.CreateMerchant)
		platform.POST("/auth/login", service.LoginMerchant)
	}

	platformAuth := platform.Group("")
	platformAuth.Use(handler.MerchantAuthRequired())
	{
		platformAuth.POST("/auth/logout", service.LogoutMerchant)
		platformAuth.GET("/auth/profile", service.GetCurrentMerchantProfile)
		platformAuth.GET("/merchants", service.GetMerchantList)
		platformAuth.POST("/products", service.CreateProduct)
		platformAuth.GET("/merchants/:merchantId/products", service.GetMerchantProductList)
		platformAuth.GET("/products/:id", service.GetProductDetail)
		platformAuth.POST("/products/:id/media", service.UploadProductMedia)
		platformAuth.POST("/share-links", service.CreateShareLink)
		platform.GET("/share-links/:code", service.GetShareLinkDetail)
	}

	legacy := api.Group("/legacy")
	{
		legacy.POST("getMemo", service.GetMemo)
		legacy.POST("addMemo", service.AddMemo)
		legacy.POST("delMemo", service.DelMemo)
		legacy.POST("getTask", service.GetTask)
		legacy.POST("addTask", service.AddTask)
		legacy.POST("upDateTask", service.UpDateTask)
		legacy.POST("delTask", service.DelTask)
		legacy.POST("uploadFile", service.UploadFile)
		legacy.POST("uploadFile2", service.UploadFile2)
		legacy.POST("getFileList", service.GetTblFileList)
		legacy.POST("getEmailList", service.GetEmail)
		legacy.POST("addWeightRecord", service.AddWeightRecord)
		legacy.POST("getWeightStatistics", service.GetWeightStatistics)
		legacy.POST("getWeightRecordList", service.GetWeightRecordList)
		legacy.POST("getQuotation", service.GetQuotation)
		legacy.POST("GetWeightRecord", service.GetWeightRecord)
		legacy.POST("GetWeightRecordPage", service.GetWeightRecordPage)
		legacy.GET("getdy", service.GetUrl)
	}

	r.GET("/share/:code", service.GetShareLinkDetail)
	return r
}
