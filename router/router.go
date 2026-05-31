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
		platform.POST("/users", service.CreateUser)
		platform.POST("/auth/login", service.LoginUser)
	}

	platformAuth := platform.Group("")
	platformAuth.Use(handler.UserAuthRequired())
	{
		platformAuth.POST("/auth/logout", service.LogoutUser)
		platformAuth.GET("/auth/profile", service.GetCurrentUserProfile)
		platformAuth.GET("/users", service.GetUserList)
		platformAuth.POST("/categories", service.CreateCategory)
		platformAuth.GET("/users/:userId/categories", service.GetUserCategoryList)
		platformAuth.GET("/categories/:id", service.GetCategoryDetail)
		platformAuth.POST("/categories/:id/resources", service.UploadCategoryResource)
		platformAuth.POST("/share-links", service.CreateShareLink)
		platform.GET("/share-links/:code", service.GetShareLinkDetail)
	}

	r.GET("/share/:code", service.GetShareLinkDetail)
	return r
}
