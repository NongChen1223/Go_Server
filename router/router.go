package router

import (
	"github.com/gin-gonic/gin"
	"go_server/config"
	"go_server/middleware"
	"go_server/system/controllers"
)

/*
 *gin.Engine 表示该函数为返回一个指针类型
 */
func SetRouter() *gin.Engine {
	r := gin.Default()
	auth := r.Group("/api/auth")
	// 不需要身份验证的路由
	{
		auth.POST("login", controllers.SysUserLogin)       // 用户登录
		auth.POST("register", controllers.SysUserRegister) // 用户注册
	}
	api := r.Group("/api/user")
	api.Use(middleware.JWTAuthMiddleware()) // 使用 JWT 中间件进行身份验证
	// 需要身份验证的路由
	{
		auth.GET("info", controllers.SysUserInfo) // 获取用户信息
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
