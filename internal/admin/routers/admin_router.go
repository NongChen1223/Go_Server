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
		}
	}
}
