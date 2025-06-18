package models

import (
	"time"
)

// SysAdminUser 后台管理员用户表
type SysAdminUser struct {
	AdminID       uint64     `gorm:"primaryKey;autoIncrement;comment:'管理员ID'" json:"admin_id"`
	AdminName     string     `gorm:"size:30;not null;uniqueIndex:uk_admin_name;comment:'管理员账号'" json:"admin_name"`
	RealName      string     `gorm:"size:30;not null;comment:'真实姓名'" json:"real_name"`
	Email         string     `gorm:"size:50;not null;uniqueIndex:uk_email;comment:'邮箱'" json:"email"`
	PhoneNumber   string     `gorm:"size:11;default:'';comment:'手机号码'" json:"phone_number"`
	Avatar        string     `gorm:"size:100;default:'';comment:'头像地址'" json:"avatar"`
	Password      string     `gorm:"size:100;not null;comment:'密码'" json:"password"`
	Status        int        `gorm:"type:tinyint;default:1;comment:'帐号状态（0停用 1正常）'" json:"status"`
	DelFlag       string     `gorm:"size:1;default:'0';comment:'删除标志（0存在 2删除）'" json:"del_flag"`
	LoginIP       string     `gorm:"size:50;default:'';comment:'最后登陆IP'" json:"login_ip"`
	LoginDate     *time.Time `gorm:"comment:'最后登陆时间'" json:"login_date"`
	LoginCount    int        `gorm:"default:0;comment:'登录次数'" json:"login_count"`
	LastPwdTime   *time.Time `gorm:"comment:'最后修改密码时间'" json:"last_pwd_time"`
	PwdErrorCount int        `gorm:"default:0;comment:'密码错误次数'" json:"pwd_error_count"`
	LockTime      *time.Time `gorm:"comment:'账号锁定时间'" json:"lock_time"`
	RoleIds       string     `gorm:"size:200;default:'';comment:'角色ID列表，逗号分隔'" json:"role_ids"`
	Department    string     `gorm:"size:50;default:'';comment:'所属部门'" json:"department"`
	Position      string     `gorm:"size:50;default:'';comment:'职位'" json:"position"`
	CreateBy      string     `gorm:"size:64;default:'';comment:'创建者'" json:"create_by"`
	CreateTime    *time.Time `gorm:"default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"create_time"`
	UpdateBy      string     `gorm:"size:64;default:'';comment:'更新者'" json:"update_by"`
	UpdateTime    *time.Time `gorm:"default:CURRENT_TIMESTAMP;comment:'更新时间'" json:"update_time"`
	Remark        *string    `gorm:"size:500;comment:'备注'" json:"remark"`
}

// 管理员状态常量
const (
	AdminStatusDisabled = 0 // 停用
	AdminStatusEnabled  = 1 // 正常
)

// 删除标志常量
const (
	AdminDelFlagExist   = "0" // 存在
	AdminDelFlagDeleted = "2" // 删除
)

// IsEnabled 判断管理员是否启用
func (a *SysAdminUser) IsEnabled() bool {
	return a.Status == AdminStatusEnabled && a.DelFlag == AdminDelFlagExist
}

// IsLocked 判断管理员是否被锁定
func (a *SysAdminUser) IsLocked() bool {
	if a.LockTime == nil {
		return false
	}
	// 锁定时间超过30分钟自动解锁
	return time.Since(*a.LockTime) < 30*time.Minute
}

// ShouldChangePassword 判断是否需要修改密码（超过90天）
func (a *SysAdminUser) ShouldChangePassword() bool {
	if a.LastPwdTime == nil {
		return true // 从未修改过密码，需要修改
	}
	// 超过90天需要修改密码
	return time.Since(*a.LastPwdTime) > 90*24*time.Hour
}

// CanLogin 判断是否可以登录
func (a *SysAdminUser) CanLogin() bool {
	return a.IsEnabled() && !a.IsLocked()
}

// TableName 指定表名
func (SysAdminUser) TableName() string {
	return "sys_admin_user"
}
