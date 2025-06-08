package router

import (
	"github.com/gin-gonic/gin"
	"go_server/config"
	"go_server/system/controllers"
)

/*
 *gin.Engine 表示该函数为返回一个指针类型
 */
func SetRouter() *gin.Engine {
	r := gin.Default()

	auth := r.Group("/api/auth")
	{
		auth.POST("login", controllers.SysUserLogin)
		auth.POST("register", controllers.SysUserRegister)
	}
	return r
}

// 初始化路由
func InitRouter() {
	r := SetRouter()
	port := config.AppConfig.App.Port
	if port == "" {
		port = ":8080"
	}
	r.Run(port)
}
