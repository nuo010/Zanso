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

	r.GET("/share/:code", service.GetShareLinkDetail)
	return r
}
