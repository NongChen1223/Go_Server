package middleware

import (
	"github.com/gin-gonic/gin"
	"go_server/common"
	"go_server/utils"
	"strings"
)

/**
 * JWTAuthMiddleware token验证中间件
 * 只要你的函数是 func(c *gin.Context) 这种形式（即 gin.HandlerFunc），就可以被 Gin 当作中间件或路由处理函数使用
 */
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		// 检查 Authorization 头是否存在
		if authHeader == "" {
			common.Error(c, common.ErrTokenEmpty, common.TokenEmptyMsg)
			// “中断请求流程”，让后面的代码不再执行 如果 token 校验失败，调用 c.Abort() 后，后面的 handler 不会再执行，直接返回响应
			c.Abort()
			return
		}
		//去掉请求头 Authorization 里的前缀 Bearer ，只保留后面的 token 字符串。
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claim, err := utils.ParseJWT(tokenStr)
		//token无效
		if err != nil {
			common.Error(c, common.ErrTokenInvalid, common.TokenInvalidMsg)
			c.Abort()
			return
		}
		/**
		将解析后的用户信息 用户ID和用户名存到 gin 的上下文（Context）里。
		这样后续的接口处理函数可以通过 c.Get("userID") 直接拿到当前登录用户的信息，无需重复解析 token。
		*/
		c.Set("userID", claim.UserID)
		c.Set("username", claim.Username)
		// 继续处理请求
		c.Next()
	}
}
