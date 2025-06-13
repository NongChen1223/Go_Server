package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go_server/config"
	"time"
)

// 密钥 (在生产环境中，应该将其放在环境变量中)
func getSecretKey() []byte {
	return []byte(config.AppConfig.JWT.PrivateKey)
}

// JWT 的声明结构
type MyClaims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

/**
 *  GenerateJWT
 *  @Description: 通过username生成一个Token
 *  @param username
 *  @return string
 *  @return error
 */
func GenerateJWT(userid uint64) (string, error) {
	// 使用 UUID 生成唯一的 JWTID
	jwtID := uuid.New().String()
	// 设置声明
	claims := MyClaims{
		UserID: userid, // 用户ID可以根据实际情况设置
		RegisteredClaims: jwt.RegisteredClaims{
			//Issuer:    "twelvet",                                                                                  // 发行者
			//Subject:   fmt.Sprintf("user-%d", userid),                                                             // 用户标识
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(1 * config.AppConfig.JWT.ExpirationHour))), // 1小时后过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                                             // 当前时间
			ID:        jwtID,                                                                                      // 唯一标识
		},
	}
	// 创建 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecretKey())
}

/**
 * ParseJWT
 *  @Description: 解析Token
 *  @param token
 *  @return *MyClaims
 *  @return error
 */
func ParseJWT(tokenString string) (*MyClaims, error) {
	// 解析 token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// 返回密钥，用于验证签名
		return getSecretKey(), nil
	})
	// 如果 token 无效或解析失败
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("无效的Token")
	}
	// 将 Claims 提取并返回
	if claims, ok := token.Claims.(*MyClaims); ok {
		return claims, nil
	}
	return nil, fmt.Errorf("无法解析claims")
}
