package models

import (
	"time"
)

// SysRole 系统角色表
type SysRole struct {
	RoleID            uint64     `gorm:"primaryKey;autoIncrement;comment:'角色ID'" json:"role_id"`
	RoleName          string     `gorm:"size:30;not null;comment:'角色名称'" json:"role_name"`
	RoleKey           string     `gorm:"size:100;not null;uniqueIndex:uk_role_key;comment:'角色权限字符串'" json:"role_key"`
	RoleSort          int        `gorm:"default:0;comment:'显示顺序'" json:"role_sort"`
	DataScope         int        `gorm:"type:tinyint;default:1;comment:'数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限 5：仅本人数据权限）'" json:"data_scope"`
	MenuCheckStrictly int        `gorm:"type:tinyint;default:1;comment:'菜单树选择项是否关联显示'" json:"menu_check_strictly"`
	DeptCheckStrictly int        `gorm:"type:tinyint;default:1;comment:'部门树选择项是否关联显示'" json:"dept_check_strictly"`
	Status            int        `gorm:"type:tinyint;default:1;index:idx_status;comment:'角色状态（0停用 1正常）'" json:"status"`
	DelFlag           string     `gorm:"size:1;default:'0';comment:'删除标志（0存在 2删除）'" json:"del_flag"`
	CreateBy          string     `gorm:"size:64;default:'';comment:'创建者'" json:"create_by"`
	CreateTime        *time.Time `gorm:"default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"create_time"`
	UpdateBy          string     `gorm:"size:64;default:'';comment:'更新者'" json:"update_by"`
	UpdateTime        *time.Time `gorm:"default:CURRENT_TIMESTAMP;comment:'更新时间'" json:"update_time"`
	Remark            *string    `gorm:"size:500;comment:'备注'" json:"remark"`

	// 关联字段（不存储到数据库）
	MenuIDs []uint64   `gorm:"-" json:"menu_ids,omitempty"` // 菜单ID列表
	Menus   []*SysMenu `gorm:"-" json:"menus,omitempty"`    // 菜单列表
}

// 数据范围常量
const (
	DataScopeAll        = 1 // 全部数据权限
	DataScopeCustom     = 2 // 自定数据权限
	DataScopeDept       = 3 // 本部门数据权限
	DataScopeDeptAndSub = 4 // 本部门及以下数据权限
	DataScopeSelf       = 5 // 仅本人数据权限
)

// 角色状态常量
const (
	RoleStatusDisabled = 0 // 停用
	RoleStatusEnabled  = 1 // 正常
)

// 删除标志常量
const (
	RoleDelFlagExist   = "0" // 存在
	RoleDelFlagDeleted = "2" // 删除
)

// 树选择关联显示常量
const (
	CheckStrictlyFalse = 0 // 父子不互相关联显示
	CheckStrictlyTrue  = 1 // 父子互相关联显示
)

// IsEnabled 判断角色是否启用
func (r *SysRole) IsEnabled() bool {
	return r.Status == RoleStatusEnabled && r.DelFlag == RoleDelFlagExist
}

// IsAdmin 判断是否为超级管理员角色
func (r *SysRole) IsAdmin() bool {
	return r.RoleKey == "admin"
}

// HasAllDataScope 判断是否拥有全部数据权限
func (r *SysRole) HasAllDataScope() bool {
	return r.DataScope == DataScopeAll
}

// HasCustomDataScope 判断是否为自定义数据权限
func (r *SysRole) HasCustomDataScope() bool {
	return r.DataScope == DataScopeCustom
}

// HasDeptDataScope 判断是否为部门数据权限
func (r *SysRole) HasDeptDataScope() bool {
	return r.DataScope == DataScopeDept
}

// HasDeptAndSubDataScope 判断是否为部门及以下数据权限
func (r *SysRole) HasDeptAndSubDataScope() bool {
	return r.DataScope == DataScopeDeptAndSub
}

// HasSelfDataScope 判断是否仅本人数据权限
func (r *SysRole) HasSelfDataScope() bool {
	return r.DataScope == DataScopeSelf
}

// GetDataScopeName 获取数据范围名称
func (r *SysRole) GetDataScopeName() string {
	switch r.DataScope {
	case DataScopeAll:
		return "全部数据权限"
	case DataScopeCustom:
		return "自定数据权限"
	case DataScopeDept:
		return "本部门数据权限"
	case DataScopeDeptAndSub:
		return "本部门及以下数据权限"
	case DataScopeSelf:
		return "仅本人数据权限"
	default:
		return "未知"
	}
}

// TableName 指定表名
func (SysRole) TableName() string {
	return "sys_role"
}

// SysRoleMenu 角色和菜单关联表
type SysRoleMenu struct {
	RoleID uint64 `gorm:"primaryKey;comment:'角色ID'" json:"role_id"`
	MenuID uint64 `gorm:"primaryKey;comment:'菜单ID'" json:"menu_id"`
}

// TableName 指定表名
func (SysRoleMenu) TableName() string {
	return "sys_role_menu"
}

// SysAdminUserRole 管理员和角色关联表
type SysAdminUserRole struct {
	AdminID uint64 `gorm:"primaryKey;comment:'管理员ID'" json:"admin_id"`
	RoleID  uint64 `gorm:"primaryKey;comment:'角色ID'" json:"role_id"`
}

// TableName 指定表名
func (SysAdminUserRole) TableName() string {
	return "sys_admin_user_role"
}
