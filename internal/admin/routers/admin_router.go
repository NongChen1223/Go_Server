package routers

import (
	"github.com/gin-gonic/gin"
	"go_server/internal/admin/controllers"
	"go_server/middleware"
)

// SetupAdminRoutes 设置后台管理路由
func SetupAdminRoutes(router *gin.RouterGroup) {
	// 后台管理路由组
	admin := router.Group("/admin")
	{
		// 不需要认证的路由
		admin.POST("/login", controllers.AdminLogin)       // 管理员登录
		admin.POST("/register", controllers.AdminRegister) // 管理员注册

		// 需要管理员认证的路由
		authAdmin := admin.Group("")
		authAdmin.Use(middleware.AdminAuthMiddleware()) // 管理员权限中间件
		{
			authAdmin.POST("/logout", controllers.AdminLogout) // 管理员登出

			// 字典管理路由组
			// 字典管理路由组
			dictGroup := authAdmin.Group("/dict")
			{
				// 字典类型管理
				dictTypeGroup := dictGroup.Group("/types")
				{
					dictTypeGroup.GET("", controllers.GetDictTypeList)       // 获取字典类型列表
					dictTypeGroup.GET("/:id", controllers.GetDictTypeDetail) // 获取字典类型详情
					dictTypeGroup.POST("", controllers.CreateDictType)       // 创建字典类型
					dictTypeGroup.PUT("/:id", controllers.UpdateDictType)    // 更新字典类型
					dictTypeGroup.DELETE("/:id", controllers.DeleteDictType) // 删除字典类型
				}

				// 字典数据管理
				dictDataGroup := dictGroup.Group("/data")
				{
					dictDataGroup.GET("/:type", controllers.GetDictDataByType) // 根据类型获取字典数据
				}
			}

			// 游戏厂商管理路由组
			publisherGroup := authAdmin.Group("/publishers")
			{
				publisherGroup.GET("", controllers.GetPublisherList)       // 获取厂商列表
				publisherGroup.GET("/:id", controllers.GetPublisherDetail) // 获取厂商详情
				publisherGroup.POST("", controllers.CreatePublisher)       // 创建厂商
				publisherGroup.PUT("/:id", controllers.UpdatePublisher)    // 更新厂商
				publisherGroup.DELETE("/:id", controllers.DeletePublisher) // 删除厂商
			}

			// 游戏管理路由组
			gameGroup := authAdmin.Group("/games")
			{
				gameGroup.GET("", controllers.GetGameList)       // 获取游戏列表
				gameGroup.GET("/:id", controllers.GetGameDetail) // 获取游戏详情
				gameGroup.POST("", controllers.CreateGame)       // 创建游戏
				gameGroup.PUT("/:id", controllers.UpdateGame)    // 更新游戏
				gameGroup.DELETE("/:id", controllers.DeleteGame) // 删除游戏
			}
		}
	}
}
