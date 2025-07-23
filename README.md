# 游戏管理系统 API

基于 Go + Gin + GORM 构建的游戏管理系统后端API，支持前台用户和后台管理双端功能。

## ✨ 特性

- 🚀 **模块化架构**：分离前台API、后台管理、系统功能
- 📚 **自动文档**：集成Swagger，自动生成API文档
- 🔐 **JWT认证**：支持用户和管理员双重认证体系
- 📊 **数据管理**：完整的CRUD操作和分页查询
- 🎮 **游戏管理**：游戏、厂商、分类等完整管理功能

## 📋 项目目录结构

```
Go_Server/
├── common/                     # 公共模块
│   ├── dto.go                 # 数据传输对象定义
│   └── response.go            # 统一响应格式
├── config/                     # 配置文件目录
│   ├── config.yml             # 主配置文件（数据库、JWT等）
│   ├── config.demo.yml        # 配置文件示例
│   ├── config.go              # 配置结构体定义
│   └── db.go                  # 数据库连接配置
├── constants/                  # 常量定义
│   ├── jwt.go                 # JWT相关常量
│   └── response_code.go       # 响应码常量
├── data/                       # 数据文件
│   └── base/                  # 基础数据
├── docs/                       # 生成的API文档
│   ├── docs.go                # Swagger文档定义
│   ├── swagger.json           # JSON格式文档
│   └── swagger.yaml           # YAML格式文档
├── global/                     # 全局变量
│   └── global.go              # 全局配置和变量
├── internal/                   # 内部业务模块
│   ├── admin/                 # 后台管理模块
│   │   ├── controllers/       # 控制器层
│   │   ├── dto/              # 管理员模块DTO
│   │   ├── routers/          # 路由配置
│   │   └── services/         # 业务逻辑层
│   ├── system/               # 系统模块
│   │   ├── controllers/      # 系统控制器
│   │   ├── dto/              # 系统模块DTO
│   │   ├── routers/          # 系统路由
│   │   └── services/         # 系统服务
│   └── web/                  # 前台模块
│       └── routes/           # 前台路由
├── middleware/                 # 中间件
│   └── auth_middleware.go     # 认证中间件
├── models/                     # 数据模型
│   ├── sys_admin_user.go      # 管理员用户模型
│   ├── sys_user.go            # 前台用户模型
│   ├── sys_dict_type.go       # 字典类型模型
│   ├── sys_dict_data.go       # 字典数据模型
│   ├── sys_publisher.go       # 厂商模型
│   ├── sys_game.go            # 游戏模型
│   ├── sys_game_cover.go      # 游戏封面模型
│   ├── sys_game_screenshot.go # 游戏截图模型
│   ├── sys_user_follow.go     # 用户关注模型
│   ├── sys_user_game_comment.go # 游戏评论模型
│   ├── sys_comment_like.go    # 评论点赞模型
│   ├── sys_game_type_relation.go     # 游戏类型关联
│   ├── sys_game_platform_relation.go # 游戏平台关联
│   └── sys_game_language_relation.go # 游戏语言关联
├── router/                     # 主路由
│   └── router.go              # 路由汇总配置
├── utils/                      # 工具函数
│   ├── jwt.go                 # JWT工具
│   ├── utils.go               # 通用工具函数
│   └── validator.go           # 数据验证工具
├── main.go                     # 程序入口
├── go.mod                      # Go模块定义
├── go.sum                      # 依赖版本锁定
└── README.md                   # 项目说明文档
```

## 🚀 快速开始

### 1. 环境要求

- **Go**: 1.19+
- **MySQL**: 5.7+
- **Git**: 最新版本
- **操作系统**: macOS/Linux/Windows

### 2. 项目初始化

#### 2.1 克隆项目
```bash
# 克隆项目到本地
git clone <repository-url>
cd Go_Server
```

#### 2.2 安装Go依赖
```bash
# 下载并安装项目依赖
go mod download
go mod tidy
```

#### 2.3 安装开发工具
```bash
# 安装Swagger文档生成工具
go install github.com/swaggo/swag/cmd/swag@latest

# 安装Swagger相关依赖
go get github.com/swaggo/files
go get github.com/swaggo/gin-swagger
go get github.com/swaggo/swag

# 验证安装
swag --version
```

### 3. 配置文件设置

#### 3.1 数据库配置

编辑 `config/config.yml` 文件：

```yaml
# 应用配置
App:
  Name: Go_server         # 应用程序名称
  Port: ":8801"          # 应用程序运行的端口

# 数据库配置
Database:
  Host: localhost        # 数据库主机地址
  Port: 3306            # 数据库端口
  User: root            # 数据库用户名
  Password: your_password    # 数据库密码
  Name: twelvet         # 数据库名称
  Charset: utf8         # 字符集
  ParseTime: True       # 是否解析时间段
  Loc: Local           # 时区
  MaxIdleCones: 10     # 空闲连接池中连接的最大数量
  MaxOpenCones: 100    # 数据库连接的最大数量

# JWT配置
JWT:
  PrivateKey: "Go_server"  # JWT秘钥（建议使用复杂密钥）
  ExpirationHour: 10       # JWT有效小时
```

#### 3.2 数据库初始化

```bash
# 1. 创建MySQL数据库
mysql -u root -p
CREATE DATABASE twelvet CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 2. 导入数据库结构（如果有SQL文件）
mysql -u root -p twelvet < database/twelvet.sql

# 3. 或者让程序自动创建表结构（如果配置了自动迁移）
```

### 4. 启动项目

#### 4.1 开发模式启动
```bash
# 直接运行
go run main.go

# 或者先编译再运行
go build -o go_server
./go_server
```

#### 4.2 生产模式启动
```bash
# 设置生产环境变量
export GIN_MODE=release
go run main.go
```

### 5. 验证安装

#### 5.1 检查服务状态
```bash
# 健康检查接口
curl http://localhost:8801/health

# 预期响应
{
  "status": "ok",
  "message": "服务运行正常"
}
```

#### 5.2 访问API文档
- **Swagger UI**: http://localhost:8801/swagger/index.html
- **JSON格式**: http://localhost:8801/swagger/doc.json

## 📖 API文档管理

### 1. 生成Swagger文档

#### 1.1 基本生成
```bash
# 生成文档到docs目录
swag init -g main.go --output ./docs --parseDependency --parseInternal
```

#### 1.2 查看生成结果
```bash
# 检查生成的文件
ls -la docs/
# docs.go swagger.json swagger.yaml

# 查看接口数量
cat docs/swagger.json | jq '.paths | keys | length'

# 查看接口列表
cat docs/swagger.json | jq '.paths | keys'
```

### 2. Swagger注解规范

#### 2.1 主文档注解（main.go）
```go
// @title           游戏管理系统 API
// @version         1.0
// @description     这是一个游戏管理系统的API文档
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
```

#### 2.2 接口注解示例
```go
// @Summary 管理员登录
// @Description 管理员使用账号密码登录系统，返回JWT token
// @Tags 管理员认证
// @Accept json
// @Produce json
// @Param request body dto.AdminLoginReq true "登录请求参数"
// @Success 200 {object} common.Response{data=map[string]interface{}} "登录成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "认证失败"
// @Router /v1/admin/login [post]
func AdminLogin(c *gin.Context) {
    // 实现代码
}
```

## 🛠️ 开发指南

### 1. 添加新接口

#### 1.1 创建控制器
```go
// internal/admin/controllers/new_controller.go
package controllers

import (
    "github.com/gin-gonic/gin"
    "go_server/common"
)

// NewFunction 新功能
// @Summary 新功能描述
// @Description 详细的功能说明
// @Tags 功能分组
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.NewReq true "请求参数"
// @Success 200 {object} common.Response "成功响应"
// @Failure 400 {object} common.Response "参数错误"
// @Router /v1/admin/new [post]
func NewFunction(c *gin.Context) {
    // 实现逻辑
    common.Success(c, "操作成功")
}
```

#### 1.2 添加路由
```go
// internal/admin/routers/admin_router.go
authAdmin.POST("/new", controllers.NewFunction)
```

#### 1.3 重新生成文档
```bash
swag init -g main.go --output ./docs --parseDependency --parseInternal
```

### 2. 数据模型管理

#### 2.1 创建模型
```go
// models/new_model.go
package models

import (
    "gorm.io/gorm"
    "time"
)

type NewModel struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    Name      string         `gorm:"size:100;not null" json:"name"`
    Status    int            `gorm:"default:1" json:"status"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

#### 2.2 创建DTO
```go
// internal/admin/dto/new_dto.go
package dto

type NewReq struct {
    Name   string `json:"name" binding:"required" example:"示例名称"`
    Status int    `json:"status" example:"1"`
}

type NewRes struct {
    ID     uint   `json:"id"`
    Name   string `json:"name"`
    Status int    `json:"status"`
}
```

### 3. 常用开发命令

```bash
# 项目管理
go mod download      # 下载依赖
go mod tidy         # 整理依赖

# 开发调试
go run main.go      # 启动服务
go build           # 编译项目

# 文档管理
swag init -g main.go --output ./docs --parseDependency --parseInternal

# 测试相关
go test ./...      # 运行测试
go fmt ./...       # 格式化代码
```

## 🔧 故障排除

### 1. 常见问题

#### 1.1 swag命令未找到
```bash
# 重新安装swag工具
go install github.com/swaggo/swag/cmd/swag@latest

# 检查GOPATH
echo $GOPATH
export PATH=$PATH:$GOPATH/bin

# 或使用完整路径
/Users/$(whoami)/go/bin/swag init -g main.go --output ./docs
```

#### 1.2 数据库连接失败
```bash
# 检查MySQL服务状态
brew services list | grep mysql  # macOS
systemctl status mysql           # Linux

# 测试数据库连接
mysql -h localhost -u root -p -e "SELECT 1"

# 检查配置文件
cat config/config.yml | grep -A 10 "Database:"
```

#### 1.3 端口被占用
```bash
# 查看端口占用
lsof -i :8801

# 杀死占用进程
kill -9 <PID>

# 或修改配置文件中的端口
```

#### 1.4 Swagger文档为空
```bash
# 检查注解语法
grep -r "@Summary" internal/

# 重新生成文档
rm -rf docs/
swag init -g main.go --output ./docs --parseDependency --parseInternal

# 检查生成结果
cat docs/swagger.json | jq '.paths'
```

## 📚 相关文档

- [Go官方文档](https://golang.org/doc/)
- [Gin框架文档](https://gin-gonic.com/docs/)
- [GORM文档](https://gorm.io/docs/)
- [Swagger注解文档](https://github.com/swaggo/swag)

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。
