package models

import (
	"time"
)

// SysUserFollow 用户关注表（关注其他用户）
type SysUserFollow struct {
	FollowID     uint64     `gorm:"primaryKey;autoIncrement;comment:'关注ID - 主键，自动递增'" json:"follow_id"`
	UserID       uint64     `gorm:"not null;uniqueIndex:uk_user_follow,priority:1;index:idx_user_id;comment:'用户ID - 关注者，关联sys_user表'" json:"user_id"`
	FollowUserID uint64     `gorm:"not null;uniqueIndex:uk_user_follow,priority:2;index:idx_follow_user_id;comment:'被关注用户ID - 被关注者，关联sys_user表'" json:"follow_user_id"`
	CreateTime   *time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'关注时间 - 自动填充'" json:"create_time"`
}

// TableName 指定表名
func (SysUserFollow) TableName() string {
	return "sys_user_follow"
}
