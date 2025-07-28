package models

import (
	"time"
)

// SysMenu 系统菜单表
type SysMenu struct {
	MenuID     uint64     `gorm:"primaryKey;autoIncrement;comment:'菜单ID'" json:"menu_id"`
	ParentID   uint64     `gorm:"default:0;index:idx_parent_id;comment:'父菜单ID'" json:"parent_id"`
	Ancestors  string     `gorm:"size:50;default:'';comment:'祖级列表'" json:"ancestors"`
	MenuName   string     `gorm:"size:50;not null;comment:'菜单名称'" json:"menu_name"`
	OrderNum   int        `gorm:"default:0;comment:'显示顺序'" json:"order_num"`
	Icon       string     `gorm:"size:100;default:'#';comment:'菜单图标'" json:"icon"`
	Path       string     `gorm:"size:200;default:'';comment:'路由地址'" json:"path"`
	Component  *string    `gorm:"size:255;comment:'组件路径'" json:"component"`
	Query      *string    `gorm:"size:255;comment:'路由参数'" json:"query"`
	MenuType   string     `gorm:"size:1;default:'M';index:idx_menu_type;comment:'菜单类型（M目录 C菜单 F按钮）'" json:"menu_type"`
	Visible    int        `gorm:"type:tinyint;default:1;index:idx_visible;comment:'菜单状态（0隐藏 1显示）'" json:"visible"`
	Status     int        `gorm:"type:tinyint;default:1;index:idx_status;comment:'菜单状态（0停用 1正常）'" json:"status"`
	Perms      *string    `gorm:"size:100;comment:'权限标识'" json:"perms"`
	IsFrame    int        `gorm:"type:tinyint;default:1;comment:'是否为外链（0是 1否）'" json:"is_frame"`
	IsCache    int        `gorm:"type:tinyint;default:0;comment:'是否缓存（0不缓存 1缓存）'" json:"is_cache"`
	CreateBy   string     `gorm:"size:64;default:'';comment:'创建者'" json:"create_by"`
	CreateTime *time.Time `gorm:"default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"create_time"`
	UpdateBy   string     `gorm:"size:64;default:'';comment:'更新者'" json:"update_by"`
	UpdateTime *time.Time `gorm:"default:CURRENT_TIMESTAMP;comment:'更新时间'" json:"update_time"`
	Remark     string     `gorm:"size:500;default:'';comment:'备注'" json:"remark"`

	// 关联字段（不存储到数据库）
	Children []*SysMenu `gorm:"-" json:"children,omitempty"` // 子菜单列表
}

// 菜单类型常量
const (
	MenuTypeDir    = "M" // 目录
	MenuTypeMenu   = "C" // 菜单
	MenuTypeButton = "F" // 按钮
)

// 菜单状态常量
const (
	MenuStatusDisabled = 0 // 停用
	MenuStatusEnabled  = 1 // 正常
)

// 菜单可见性常量
const (
	MenuVisibleHidden = 0 // 隐藏
	MenuVisibleShow   = 1 // 显示
)

// 外链常量
const (
	MenuIsFrame  = 0 // 是外链
	MenuNotFrame = 1 // 不是外链
)

// 缓存常量
const (
	MenuNotCache = 0 // 不缓存
	MenuIsCache  = 1 // 缓存
)

// IsEnabled 判断菜单是否启用
func (m *SysMenu) IsEnabled() bool {
	return m.Status == MenuStatusEnabled
}

// IsVisible 判断菜单是否可见
func (m *SysMenu) IsVisible() bool {
	return m.Visible == MenuVisibleShow
}

// IsDirectory 判断是否为目录
func (m *SysMenu) IsDirectory() bool {
	return m.MenuType == MenuTypeDir
}

// IsMenu 判断是否为菜单
func (m *SysMenu) IsMenu() bool {
	return m.MenuType == MenuTypeMenu
}

// IsButton 判断是否为按钮
func (m *SysMenu) IsButton() bool {
	return m.MenuType == MenuTypeButton
}

// IsExternalLink 判断是否为外链
func (m *SysMenu) IsExternalLink() bool {
	return m.IsFrame == MenuIsFrame
}

// HasPermission 判断是否有权限标识
func (m *SysMenu) HasPermission() bool {
	return m.Perms != nil && *m.Perms != ""
}

// GetPermission 获取权限标识
func (m *SysMenu) GetPermission() string {
	if m.Perms == nil {
		return ""
	}
	return *m.Perms
}

// IsTopLevel 判断是否为顶级菜单
func (m *SysMenu) IsTopLevel() bool {
	return m.ParentID == 0
}

// TableName 指定表名
func (SysMenu) TableName() string {
	return "sys_menu"
}
