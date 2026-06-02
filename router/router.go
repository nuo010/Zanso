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
		platformAuth.POST("/users/:id/role", service.UpdateUserRole)
		platformAuth.POST("/collections", service.CreateCategory)
		platformAuth.GET("/collections", service.GetCurrentUserCategoryList)
		platformAuth.POST("/collections/:id/update", service.UpdateCategory)
		platformAuth.POST("/collections/:id/delete", service.DeleteCategory)
		platformAuth.POST("/categories", service.CreateCategoryItem)
		platformAuth.POST("/categories/:id/update", service.UpdateCategoryItem)
		platformAuth.POST("/categories/:id/delete", service.DeleteCategoryItem)
		platformAuth.POST("/categories/:id/resources/sort", service.UpdateCategoryResourceSort)
		platformAuth.GET("/dashboard/stats", service.GetDashboardStats)
		platformAuth.GET("/collections/:id", service.GetCategoryDetail)
		platformAuth.GET("/categories/:id", service.GetCategoryItemDetail)
		platformAuth.POST("/collections/:id/resources", service.UploadCategoryResource)
		platformAuth.POST("/resources/:id/delete", service.DeleteResource)
		platformAuth.POST("/share-links", service.CreateShareLink)
		platformAuth.GET("/share-links", service.GetShareLinkList)
		platformAuth.POST("/share-links/:id/delete", service.DeleteShareLink)
		platform.GET("/share-links/:code", service.GetShareLinkDetail)
	}

	r.GET("/share/:code", service.RenderSharePage)
	return r
}
