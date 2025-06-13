package dto

import (
	"go_server/utils"
	"time"
)

type LoginReq struct {
	Username  string `json:"username" binding:"required"`  // 用户名，必填
	Password  string `json:"password" binding:"required"`  // 密码，必填
	AutoLogin bool   `json:"autoLogin"`                    // 是否自动登录，可选
	GrantType string `json:"grantType" binding:"required"` // 授权类型，必填
	Type      string `json:"type" binding:"required"`      // 登录类型，必填
}

type LoginRes struct {
	Code        int    `json:"code"`
	AccessToken string `json:"access_token"` // 访问令牌
}

// GetMessages 自定义错误信息
func (loginReq LoginReq) GetMessages() utils.ValidatorMessages {
	return utils.ValidatorMessages{
		"Username.required":  "用户名不能为空",
		"Password.required":  "用户密码不能为空",
		"GrantType.required": "授权类型不能为空",
		"Type.required":      "登录类型不能为空",
	}
}

// 注册用户接口
type RegisterUserReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// 自定义错误信息
func (registerUserReq RegisterUserReq) GetMessages() utils.ValidatorMessages {
	return utils.ValidatorMessages{
		"Username.required": "用户名不能为空",
		"Password.required": "用户密码不能为空",
		"Email.required":    "邮箱不能为空",
		"Email.email":       "邮箱格式不正确",
	}
}

// 获取用户信息请求
type UserInfoRes struct {
	UserID      uint64     `json:"user_id"`
	Username    string     `json:"username"`
	NickName    string     `json:"nick_name"`
	UserType    string     `json:"user_type"`
	Email       string     `json:"email"`
	PhoneNumber string     `json:"phonenumber"`
	Sex         *int       `json:"sex"`
	Avatar      string     `json:"avatar"`
	LoginIP     string     `json:"login_ip"`
	LoginDate   *time.Time `json:"login_date"`
	CreateTime  *time.Time `json:"create_time"`
}
