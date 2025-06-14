package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go_server/config"
	"go_server/constants"
	"time"
)

type JWTError struct {
	Code    int
	Message string
}

func (e *JWTError) Error() string {
	return e.Message
}

// 密钥 (在生产环境中，应该将其放在环境变量中)
func getSecretKey() []byte {
	return []byte(config.AppConfig.JWT.PrivateKey)
}

// JWT 的声明结构
type MyClaims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

/**
 *  GenerateJWT
 *  @Description: 通过username生成一个Token
 *  @param username
 *  @return string
 *  @return error
 */
func GenerateJWT(user_id uint64) (string, error) {
	// 使用 UUID 生成唯一的 JWTID
	jwtID := uuid.New().String()
	// 设置声明
	claims := MyClaims{
		UserID: user_id, // 用户ID可以根据实际情况设置
		RegisteredClaims: jwt.RegisteredClaims{
			//Issuer:    "twelvet",                                                                                  // 发行者
			//Subject:   fmt.Sprintf("user-%d", user_id),                                                             // 用户标识
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.AppConfig.JWT.ExpirationHour) * time.Hour)), // n小时后过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                                                     // 当前时间
			ID:        jwtID,                                                                                              // 唯一标识
		},
	}
	// 创建 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecretKey())
}

/**
 * ParseJWT
 * @Description: 解析并校验 JWT Token，区分过期、无效、签名算法错误等情况，返回不同的响应码和信息
 * @param tokenString 待解析的 JWT 字符串
 * @return *MyClaims 解析出的自定义 claims
 * @return error 解析失败时返回自定义错误，包含 code 和 message
 */
func ParseJWT(tokenString string) (*MyClaims, error) {
	// 1. 解析 token，校验签名算法和密钥
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 校验签名算法是否为 HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// 签名算法错误
			return nil, &JWTError{Code: constants.ErrorCode, Message: constants.SignatureErrorMsg}
		}
		// 返回用于校验签名的密钥
		return getSecretKey(), nil
	})

	// 2. 解析失败，判断是否为 token 过期
	if err != nil {
		// 判断错误类型是否为 token 过期
		if errors.Is(err, jwt.ErrTokenExpired) {
			// Token 已过期
			return nil, &JWTError{Code: constants.ErrTokenExpired, Message: constants.TokenExpiredMsg}
		}
		// 其他解析错误
		return nil, &JWTError{Code: constants.ErrTokenInvalid, Message: constants.TokenAnalysis}
	}

	// 3. token 校验未通过（如签名不匹配等）
	if !token.Valid {
		return nil, &JWTError{Code: constants.ErrTokenInvalid, Message: constants.TokenInvalidMsg}
	}

	// 4. 提取自定义 claims 并返回
	if claims, ok := token.Claims.(*MyClaims); ok {
		return claims, nil
	}
	// 5. claims 解析失败
	return nil, &JWTError{Code: constants.ErrTokenInvalid, Message: constants.ClaimsAnalysis}
}
