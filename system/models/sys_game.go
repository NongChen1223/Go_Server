package models

import (
	"time"
)

// SysGame 游戏基本信息表
type SysGame struct {
	GameID       uint64     `gorm:"primaryKey;autoIncrement;comment:'游戏ID - 主键，自动递增'" json:"game_id"`
	NameZh       string     `gorm:"size:100;not null;comment:'游戏中文名称 - 必填'" json:"name_zh"`
	NameEn       *string    `gorm:"size:100;comment:'游戏英文名称 - 可选'" json:"name_en"`
	ReleaseDate  *time.Time `gorm:"type:date;comment:'游戏发布日期 - 可选'" json:"release_date"`
	Description  *string    `gorm:"type:text;comment:'游戏介绍 - 可选，存储长文本'" json:"description"`
	Rating       *float64   `gorm:"type:decimal(3,1);default:0.0;comment:'游戏评分（最多一位小数）- 默认0分'" json:"rating"`
	Size         *string    `gorm:"size:50;comment:'游戏大小 - 可选，如\"2.5GB\"'" json:"size"`
	Price        *float64   `gorm:"type:decimal(10,2);comment:'游戏价格 - 可选，保留两位小数'" json:"price"`
	PublisherID  *uint64    `gorm:"comment:'游戏厂商ID - 外键，关联sys_publisher表'" json:"publisher_id"`
	StudioID     *uint64    `gorm:"comment:'游戏工作室ID - 外键，关联sys_studio表'" json:"studio_id"`
	ShutdownDate *time.Time `gorm:"type:date;comment:'游戏停服日期 - 可选'" json:"shutdown_date"`
	CommentCount int        `gorm:"not null;default:0;comment:'游戏评论数量 - 默认0'" json:"comment_count"`
	LikeCount    int        `gorm:"not null;default:0;comment:'游戏点赞数量 - 默认0'" json:"like_count"`
	DemoVideo    *string    `gorm:"size:255;comment:'游戏演示视频链接 - 可选'" json:"demo_video"`
	CreateTime   *time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'创建时间 - 自动填充当前时间'" json:"create_time"`
	UpdateTime   *time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'更新时间 - 自动更新为当前时间'" json:"update_time"`
	CreateBy     string     `gorm:"size:64;default:'';comment:'创建者 - 记录创建人'" json:"create_by"`
	UpdateBy     string     `gorm:"size:64;default:'';comment:'更新者 - 记录最后修改人'" json:"update_by"`
	Status       int        `gorm:"type:tinyint;not null;default:1;comment:'状态（0停用 1正常）- 控制游戏是否可见'" json:"status"`
}

// TableName 指定表名
func (SysGame) TableName() string {
	return "sys_game"
}
