package dto

import (
	"go_server/utils"
	"time"
)

// 登录用户接口
type LoginReq struct {
	UserName string `json:"user_name" binding:"required"` // 用户名，必填
	Password string `json:"password" binding:"required"`  // 密码，必填
}

// 登录响应接口
type LoginRes struct {
	Code        int    `json:"code"`
	AccessToken string `json:"access_token"` // 访问令牌
}

// 登录用户接口 自定义错误信息
func (loginReq LoginReq) GetMessages() utils.ValidatorMessages {
	return utils.ValidatorMessages{
		"UserName.required": "用户名不能为空",
		"Password.required": "用户密码不能为空",
	}
}

// 注册用户接口
type RegisterUserReq struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// 注册用户接口 自定义错误信息
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
	UserName    string     `json:"user_name"`
	NickName    string     `json:"nick_name"`
	UserType    string     `json:"user_type"`
	Email       string     `json:"email"`
	PhoneNumber string     `json:"phone_number"`
	Sex         *int       `json:"sex"`
	Avatar      string     `json:"avatar"`
	LoginIP     string     `json:"login_ip"`
	LoginDate   *time.Time `json:"login_date"`
	CreateTime  *time.Time `json:"create_time"`
}
