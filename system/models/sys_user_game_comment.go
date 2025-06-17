package models

import (
	"time"
)

// SysUserGameComment 用户游戏评论表
type SysUserGameComment struct {
	CommentID  uint64     `gorm:"primaryKey;autoIncrement;comment:'评论ID - 主键，自动递增'" json:"comment_id"`
	UserID     uint64     `gorm:"not null;index:idx_user_id;comment:'用户ID - 评论者，关联sys_user表'" json:"user_id"`
	GameID     uint64     `gorm:"not null;index:idx_game_id;comment:'游戏ID - 被评论的游戏，关联sys_game表'" json:"game_id"`
	Content    string     `gorm:"type:text;not null;comment:'评论内容 - 用户的评论文字'" json:"content"`
	ParentID   *uint64    `gorm:"index:idx_parent_id;comment:'父评论ID - 回复的评论ID，NULL表示一级评论'" json:"parent_id"`
	LikeCount  int        `gorm:"not null;default:0;comment:'点赞数 - 该评论获得的点赞数'" json:"like_count"`
	Status     int        `gorm:"type:tinyint;not null;default:1;comment:'状态（0隐藏 1显示）- 控制评论是否可见'" json:"status"`
	CreateTime *time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'评论时间 - 自动填充'" json:"create_time"`
	UpdateTime *time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'更新时间 - 自动更新'" json:"update_time"`
}

// TableName 指定表名
func (SysUserGameComment) TableName() string {
	return "sys_user_game_comment"
}
