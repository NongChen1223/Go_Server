package models

import (
	"time"
)

// SysPublisher 游戏厂商表
type SysPublisher struct {
	PublisherID     uint64     `gorm:"primaryKey;autoIncrement;comment:'厂商ID - 主键，自动递增'" json:"publisher_id"`
	PublisherName   string     `gorm:"size:100;not null;comment:'厂商名称 - 如\"腾讯游戏\"、\"网易游戏\"'" json:"publisher_name"`
	PublisherCnName string     `gorm:"size:100;not null;default:'';comment:'厂商中文名称'" json:"publisher_cn_name"`
	LogoURL         *string    `gorm:"size:255;comment:'厂商LOGO - 存储logo图片路径'" json:"logo_url"`
	Description     *string    `gorm:"type:text;comment:'厂商介绍 - 存储长文本'" json:"description"`
	FoundedDate     *time.Time `gorm:"type:date;comment:'成立日期 - 厂商成立时间'" json:"founded_date"`
	Website         *string    `gorm:"size:255;comment:'官方网站 - 厂商官网链接'" json:"website"`
	Status          int        `gorm:"type:tinyint;not null;default:1;comment:'状态（0停用 1正常）- 控制厂商是否可见'" json:"status"`
	CreateTime      *time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'创建时间 - 自动填充'" json:"create_time"`
	UpdateTime      *time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'更新时间 - 自动更新'" json:"update_time"`
}

// TableName 指定表名
func (SysPublisher) TableName() string {
	return "sys_publisher"
}
