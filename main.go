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

func main() {
	fmt.Println("Initializing application...") // 提前打印初始化信息
	config.InitConfig()
	fmt.Println("App start...", config.AppConfig.App.Port)
	router.InitRouter()
}
