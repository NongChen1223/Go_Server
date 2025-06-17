package models

import (
	"time"
)

// SysDictType 字典类型表
type SysDictType struct {
	DictID     uint64     `gorm:"primaryKey;autoIncrement;comment:'字典类型ID - 主键，自动递增'" json:"dict_id"`
	DictName   string     `gorm:"size:100;not null;comment:'字典名称 - 如\"游戏类型\"、\"平台类型\"'" json:"dict_name"`
	DictType   string     `gorm:"size:100;not null;uniqueIndex:uk_dict_type;comment:'字典类型 - 唯一标识，如\"game_type\"、\"platform_type\"'" json:"dict_type"`
	Status     int        `gorm:"type:tinyint;not null;default:1;comment:'状态（0停用 1正常）- 控制该类字典是否可用'" json:"status"`
	Remark     *string    `gorm:"size:500;comment:'备注 - 字典类型的说明'" json:"remark"`
	CreateTime *time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'创建时间 - 自动填充'" json:"create_time"`
	UpdateTime *time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'更新时间 - 自动更新'" json:"update_time"`
	CreateBy   string     `gorm:"size:64;default:'';comment:'创建者 - 记录创建人'" json:"create_by"`
	UpdateBy   string     `gorm:"size:64;default:'';comment:'更新者 - 记录最后修改人'" json:"update_by"`
}

// TableName 指定表名
func (SysDictType) TableName() string {
	return "sys_dict_type"
}
