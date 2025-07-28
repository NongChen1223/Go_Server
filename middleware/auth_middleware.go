package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_server/common"
	"go_server/constants"
	"go_server/global"
	"go_server/models"
	"go_server/utils"
	"net/http"
	"strings"
)

/**
 * JWTAuthMiddleware token验证中间件
 * 只要你的函数是 func(c *gin.Context) 这种形式（即 gin.HandlerFunc），就可以被 Gin 当作中间件或路由处理函数使用
 */
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的 Authorization 字段
		authHeader := c.Request.Header.Get("Authorization")
		fmt.Println("中间件token验证", authHeader)
		// 1. 检查 Authorization 头是否存在
		if authHeader == "" {
			// Token 不能为空，返回自定义错误码和信息
			common.Error(c, constants.ErrTokenEmpty, constants.TokenEmptyMsg)
			c.Abort() // 中断请求流程
			return
		}
		// 2. 去掉前缀 "Bearer "，只保留 token 字符串
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Println("token字符串", tokenStr)
		// 3. 解析并校验 token
		claim, err := utils.ParseJWT(tokenStr)
		fmt.Printf("中间件解析token: claim=%+v, err=%v\n", claim, err)
		// 4. 校验 token 是否有效
		if err != nil {
			// 判断是否为自定义 JWTError 类型
			if jwtErr, ok := err.(*utils.JWTError); ok {
				// 返回对应的 code 和 message
				common.Error(c, jwtErr.Code, jwtErr.Message)
			} else {
				// 兜底处理，返回通用 token 错误
				common.Error(c, constants.ErrTokenInvalid, err.Error())
			}
			c.Abort()
			return
		}
		/**
		 * 5. 将解析后的用户信息（如用户ID）存到 gin 的上下文（Context）里
		 * 这样后续的接口处理函数可以通过 c.Get("user_id") 直接拿到当前登录用户的信息，无需重复解析 token
		 */
		c.Set("user_id", claim.UserID)
		// 6. 继续处理请求
		c.Next()
	}
}

// AdminAuthMiddleware 管理员认证中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的 Authorization 字段
		authHeader := c.Request.Header.Get("Authorization")

		// 检查 Authorization 头是否存在
		if authHeader == "" {
			common.Error(c, constants.ErrTokenEmpty, "管理员token不能为空")
			c.Abort()
			return
		}

		// 去掉前缀 "Bearer "，只保留 token 字符串
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析并校验 token
		claim, err := utils.ParseJWT(tokenStr)
		if err != nil {
			if jwtErr, ok := err.(*utils.JWTError); ok {
				common.Error(c, jwtErr.Code, jwtErr.Message)
			} else {
				common.Error(c, constants.ErrTokenInvalid, "管理员token无效")
			}
			c.Abort()
			return
		}

		// 验证管理员是否存在且状态正常
		var admin models.SysAdminUser
		err = global.DB.Where("admin_id = ? AND status = ? AND del_flag = ?",
			claim.UserID, models.AdminStatusEnabled, models.AdminDelFlagExist).First(&admin).Error
		if err != nil {
			common.Error(c, constants.ErrTokenInvalid, "管理员账号不存在或已被禁用")
			c.Abort()
			return
		}

		// 将管理员信息存到上下文中
		c.Set("admin_id", admin.AdminID)
		c.Set("admin_name", admin.AdminName)

		// 继续处理请求
		c.Next()
	}
}

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许的源
		c.Header("Access-Control-Allow-Origin", "*")
		// 设置允许的请求方法
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// 设置允许的请求头
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		// 设置允许凭证
		c.Header("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// 继续处理请求
		c.Next()
	}
}
