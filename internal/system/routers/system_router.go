package routers

import (
	"github.com/gin-gonic/gin"
	"go_server/internal/system/controllers"
	"go_server/middleware"
)

// SetupSystemRoutes 设置系统路由
func SetupSystemRoutes(router *gin.RouterGroup) {
	// 前台API路由组，添加 /api 前缀
	api := router.Group("/api")
	{
		// 不需要身份验证的路由
		auth := api.Group("/auth")
		{
			auth.POST("login", controllers.SysUserLogin)       // 用户登录
			auth.POST("register", controllers.SysUserRegister) // 用户注册
		}

		// 需要身份验证的路由
		user := api.Group("/user")
		user.Use(middleware.JWTAuthMiddleware())
		{
			user.GET("info", controllers.SysUserInfo) // 获取用户信息
		}
	}
}
