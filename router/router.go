package router

import (
	"github.com/gin-gonic/gin"
	"go_server/config"
	adminRouters "go_server/internal/admin/routers"
	systemRouters "go_server/internal/system/routers"
)

/*
 * SetRouter 设置路由
 * gin.Engine 表示该函数为返回一个指针类型
 */
func SetRouter() *gin.Engine {
	r := gin.Default()

	// 创建 v1 版本的路由组
	v1 := r.Group("/v1")
	{
		// 集成各端路由
		adminRouters.SetupAdminRoutes(v1)   // 后台管理路由
		systemRouters.SetupSystemRoutes(v1) // 系统路由
	}

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "服务运行正常",
		})
	})

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
