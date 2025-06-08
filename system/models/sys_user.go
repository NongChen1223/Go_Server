package models

import (
	"time"
)

type SysUser struct {
	UserID      uint64     `gorm:"primaryKey;autoIncrement;comment:'用户ID'" json:"user_id"` //primaryKey定义主键 autoIncrement定义自增
	Username    string     `gorm:"size:30;not null;comment:'用户账号'" json:"username"`
	NickName    string     `gorm:"size:30;not null;comment:'用户昵称'" json:"nick_name"`
	UserType    string     `gorm:"size:2;default:'00';comment:'用户类型（00系统用户）'" json:"user_type"`
	Email       string     `gorm:"size:50;default:'';comment:'用户邮箱'" json:"email"`
	PhoneNumber string     `gorm:"size:11;default:'';comment:'手机号码'" json:"phonenumber"`
	Sex         *int       `gorm:"type:tinyint;comment:'用户性别（0男 1女 2未知）'" json:"sex"`
	Avatar      string     `gorm:"size:100;default:'';comment:'头像地址'" json:"avatar"`
	Password    string     `gorm:"size:100;default:'';comment:'密码'" json:"password"`
	Status      *int       `gorm:"type:tinyint;comment:'帐号状态（0正常 1停用）'" json:"status"`
	DelFlag     string     `gorm:"size:1;default:'0';comment:'删除标志（0代表存在 2代表删除）'" json:"del_flag"`
	LoginIP     string     `gorm:"size:50;default:'';comment:'最后登陆IP'" json:"login_ip"`
	LoginDate   *time.Time `gorm:"comment:'最后登陆时间'" json:"login_date"`
	CreateBy    string     `gorm:"size:64;default:'';comment:'创建者'" json:"create_by"`
	CreateTime  *time.Time `gorm:"comment:'创建时间'" json:"create_time"`
	UpdateBy    string     `gorm:"size:64;default:'';comment:'更新者'" json:"update_by"`
	UpdateTime  *time.Time `gorm:"comment:'更新时间'" json:"update_time"`
	Remark      *string    `gorm:"size:500;comment:'备注'" json:"remark"`
}
