package main

import (
	"fmt"
	"go_server/config"
	"go_server/router"
)

// @title           游戏管理系统 API
// @version         1.0
// @description     这是一个游戏管理系统的API文档，包含前端用户、后台管理、游戏管理等功能
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8801
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @tag.name Admin
// @tag.description 后台管理相关接口，包含管理员认证、角色管理、菜单管理、字典管理、游戏管理等功能

// @tag.name System
// @tag.description 系统功能相关接口，包含前台用户认证、用户管理等功能

// @tag.name Web
// @tag.description 前台网站相关接口，包含首页展示、内容管理等功能

func main() {
	fmt.Println("Initializing application...") // 提前打印初始化信息
	config.InitConfig()
	fmt.Println("App start...", config.AppConfig.App.Port)
	router.InitRouter()
}
