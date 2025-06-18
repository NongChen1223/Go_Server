package models

import (
	"time"
)

// SysDictData 字典数据表
type SysDictData struct {
	DictCode   uint64     `gorm:"primaryKey;autoIncrement;comment:'字典编码 - 主键，自动递增'" json:"dict_code"`
	DictSort   int        `gorm:"not null;default:0;comment:'字典排序 - 控制同类字典的展示顺序'" json:"dict_sort"`
	DictLabel  string     `gorm:"size:100;not null;comment:'字典标签 - 展示值，如\"角色扮演\"、\"PC\"'" json:"dict_label"`
	DictValue  string     `gorm:"size:100;not null;comment:'字典键值 - 实际存储值，如\"RPG\"、\"PC\"'" json:"dict_value"`
	DictType   string     `gorm:"size:100;not null;index:idx_dict_type;comment:'字典类型 - 关联sys_dict_type表的dict_type'" json:"dict_type"`
	IsDefault  int        `gorm:"type:tinyint;not null;default:0;comment:'是否默认（0否 1是）- 标记默认选中项'" json:"is_default"`
	Status     int        `gorm:"type:tinyint;not null;default:1;comment:'状态（0停用 1正常）- 控制该字典项是否可用'" json:"status"`
	Remark     *string    `gorm:"size:500;comment:'备注 - 字典项的说明'" json:"remark"`
	CreateTime *time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'创建时间 - 自动填充'" json:"create_time"`
	UpdateTime *time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:'更新时间 - 自动更新'" json:"update_time"`
	CreateBy   string     `gorm:"size:64;default:'';comment:'创建者 - 记录创建人'" json:"create_by"`
	UpdateBy   string     `gorm:"size:64;default:'';comment:'更新者 - 记录最后修改人'" json:"update_by"`
}

// TableName 指定表名
func (SysDictData) TableName() string {
	return "sys_dict_data"
}
