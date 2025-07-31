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
			authAdmin.GET("/info", controllers.GetAdminInfo)   // 获取管理员信息

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

			// 菜单管理路由组
			menuGroup := authAdmin.Group("/menus")
			{
				menuGroup.GET("", controllers.GetMenuList)            // 获取菜单列表
				menuGroup.GET("/tree", controllers.GetMenuTree)       // 获取菜单树
				menuGroup.GET("/routers", controllers.GetMenuRouters) // 获取菜单路由
				menuGroup.GET("/:id", controllers.GetMenuDetail)      // 获取菜单详情
				menuGroup.POST("", controllers.CreateMenu)            // 创建菜单
				menuGroup.PUT("/:id", controllers.UpdateMenu)         // 更新菜单
				menuGroup.DELETE("/:id", controllers.DeleteMenu)      // 删除菜单
			}

			// 角色管理路由组
			roleGroup := authAdmin.Group("/roles")
			{
				roleGroup.GET("", controllers.GetRoleList)                    // 获取角色列表
				roleGroup.GET("/select", controllers.GetRoleSelect)           // 获取角色选择列表
				roleGroup.GET("/admin/:admin_id", controllers.GetAdminRole)   // 获取管理员角色信息
				roleGroup.GET("/menus/:role_id", controllers.GetRoleMenuTree) // 获取角色菜单权限树
				roleGroup.GET("/:id", controllers.GetRoleDetail)              // 获取角色详情
				roleGroup.POST("", controllers.CreateRole)                    // 创建角色
				roleGroup.POST("/auth", controllers.AuthRole)                 // 角色授权（旧接口，保持兼容）
				roleGroup.POST("/menus/assign", controllers.AssignRoleMenus)  // 分配角色菜单权限
				roleGroup.POST("/assign", controllers.AssignAdminRole)        // 分配管理员角色
				roleGroup.PUT("/:id", controllers.UpdateRole)                 // 更新角色
				roleGroup.DELETE("/:id", controllers.DeleteRole)              // 删除角色
			}
		}
	}
}
