package sysConfig

import (
	"fmt"
	"zanso/db"
	"zanso/util"

	"github.com/spf13/viper"
)

// Init 初始化配置项
func Init() {
	// 从本地读取环境变量
	viper.SetConfigFile("config/config.yaml") // 设置配置文件名
	err := viper.ReadInConfig()               // 读取配置文件
	if err != nil {
		panic(fmt.Errorf("致命错误配置文件: %w", err))
	}
	util.Init()

	// 连接数据库
	db.Database()
	db.Redis()
	//db.MinioInit()
}
