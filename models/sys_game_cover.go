package models

import (
	"time"
)

// SysGameCover 游戏封面表（支持多张）
type SysGameCover struct {
	CoverID    uint64     `gorm:"primaryKey;autoIncrement;comment:'封面ID - 主键，自动递增'" json:"cover_id"`
	GameID     uint64     `gorm:"not null;comment:'游戏ID - 外键，关联sys_game表'" json:"game_id"`
	CoverURL   string     `gorm:"size:255;not null;comment:'封面图片URL - 存储图片路径'" json:"cover_url"`
	IsMain     int        `gorm:"type:tinyint;not null;default:0;comment:'是否主封面（0否 1是）- 标记主要展示图'" json:"is_main"`
	Sort       int        `gorm:"not null;default:0;comment:'排序 - 控制多张封面的展示顺序'" json:"sort"`
	CreateTime *time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'创建时间 - 自动填充'" json:"create_time"`
}

// TableName 指定表名
func (SysGameCover) TableName() string {
	return "sys_game_cover"
}
