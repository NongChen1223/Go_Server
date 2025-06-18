package models

import (
	"time"
)

// SysGameScreenshot 游戏截图表
type SysGameScreenshot struct {
	ScreenshotID  uint64     `gorm:"primaryKey;autoIncrement;comment:'截图ID - 主键，自动递增'" json:"screenshot_id"`
	GameID        uint64     `gorm:"not null;comment:'游戏ID - 外键，关联sys_game表'" json:"game_id"`
	ScreenshotURL string     `gorm:"size:255;not null;comment:'截图URL - 存储图片路径'" json:"screenshot_url"`
	Sort          int        `gorm:"not null;default:0;comment:'排序 - 控制多张截图的展示顺序'" json:"sort"`
	CreateTime    *time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'创建时间 - 自动填充'" json:"create_time"`
}

// TableName 指定表名
func (SysGameScreenshot) TableName() string {
	return "sys_game_screenshot"
}
