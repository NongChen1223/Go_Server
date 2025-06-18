package models

import (
	"time"
)

// SysCommentLike 评论点赞表
type SysCommentLike struct {
	LikeID     uint64     `gorm:"primaryKey;autoIncrement;comment:'点赞ID - 主键，自动递增'" json:"like_id"`
	CommentID  uint64     `gorm:"not null;uniqueIndex:uk_user_comment,priority:2;index:idx_comment_id;comment:'评论ID - 被点赞的评论，关联sys_user_game_comment表'" json:"comment_id"`
	UserID     uint64     `gorm:"not null;uniqueIndex:uk_user_comment,priority:1;index:idx_user_id;comment:'用户ID - 点赞者，关联sys_user表'" json:"user_id"`
	CreateTime *time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'点赞时间 - 自动填充'" json:"create_time"`
}

// TableName 指定表名
func (SysCommentLike) TableName() string {
	return "sys_comment_like"
}
