package models

import (
	"time"
)

// SysUser 前台用户表
type SysUser struct {
	UserID      uint64     `gorm:"primaryKey;autoIncrement;comment:'用户ID'" json:"user_id"`
	UserName    string     `gorm:"size:30;not null;uniqueIndex:uk_user_name;comment:'用户账号'" json:"user_name"`
	NickName    string     `gorm:"size:30;not null;comment:'用户昵称'" json:"nick_name"`
	UserType    string     `gorm:"size:2;default:'00';comment:'用户类型（00系统用户）'" json:"user_type"`
	Email       string     `gorm:"size:50;default:'';index:idx_email;comment:'用户邮箱'" json:"email"`
	PhoneNumber string     `gorm:"size:11;default:'';index:idx_phone;comment:'手机号码'" json:"phone_number"`
	Sex         *int       `gorm:"type:tinyint;comment:'用户性别（0男 1女 2未知）'" json:"sex"`
	Avatar      string     `gorm:"size:100;default:'';comment:'头像地址'" json:"avatar"`
	Password    string     `gorm:"size:100;default:'';comment:'密码'" json:"password"`
	Status      int        `gorm:"type:tinyint;default:1;comment:'帐号状态（0停用 1正常）'" json:"status"`
	DelFlag     string     `gorm:"size:1;default:'0';comment:'删除标志（0存在 2删除）'" json:"del_flag"`
	LoginIP     string     `gorm:"size:50;default:'';comment:'最后登陆IP'" json:"login_ip"`
	LoginDate   *time.Time `gorm:"comment:'最后登陆时间'" json:"login_date"`
	CreateTime  *time.Time `gorm:"default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"create_time"`
	UpdateTime  *time.Time `gorm:"default:CURRENT_TIMESTAMP;comment:'更新时间'" json:"update_time"`
	Remark      *string    `gorm:"size:500;comment:'备注'" json:"remark"`
}

// 用户状态常量
const (
	UserStatusDisabled = 0 // 停用
	UserStatusEnabled  = 1 // 正常
)

// 删除标志常量
const (
	UserDelFlagExist   = "0" // 存在
	UserDelFlagDeleted = "2" // 删除
)

// IsEnabled 判断用户是否启用
func (u *SysUser) IsEnabled() bool {
	return u.Status == UserStatusEnabled && u.DelFlag == UserDelFlagExist
}

// TableName 指定表名
func (SysUser) TableName() string {
	return "sys_user"
}
