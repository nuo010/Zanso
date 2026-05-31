package main

import (
	"fmt"
	"singo/router"
	"singo/sysConfig"
	"singo/util"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	//"singo/service"
	_ "time/tzdata"
)

func main() {
	//daemonize()
	// 从配置文件读取配置
	sysConfig.Init()
	util.BuildLogger("debug")
	util.Log().Info("服务启动")
	util.Log().Debug("调试信息")
	gin.SetMode(gin.DebugMode)
	// 装载路由
	r := router.NewRouter()
	//go service.GetFilebeat()
	err := r.Run(":" + viper.GetString("server.port"))
	if err != nil {
		fmt.Println("服务器启动失败！")
	}

	defer func() {
		err := util.CloseLogger()
		if err != nil {
			fmt.Println(err)
		}
	}()

}
