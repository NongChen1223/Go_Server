package dto

import "time"

// RoleListReq 角色列表查询请求
type RoleListReq struct {
	RoleName string `form:"role_name" example:"管理员"`                              // 角色名称
	RoleKey  string `form:"role_key" example:"admin"`                             // 角色权限字符串
	Status   *int   `form:"status" example:"1"`                                   // 角色状态（0停用 1正常）
	Page     int    `form:"page,default=1" binding:"min=1" example:"1"`           // 页码
	Size     int    `form:"size,default=10" binding:"min=1,max=100" example:"10"` // 每页数量
}

// RoleListRes 角色列表响应
type RoleListRes struct {
	RoleID            uint64     `json:"role_id"`             // 角色ID
	RoleName          string     `json:"role_name"`           // 角色名称
	RoleKey           string     `json:"role_key"`            // 角色权限字符串
	RoleSort          int        `json:"role_sort"`           // 显示顺序
	DataScope         int        `json:"data_scope"`          // 数据范围
	DataScopeName     string     `json:"data_scope_name"`     // 数据范围名称
	MenuCheckStrictly int        `json:"menu_check_strictly"` // 菜单树选择项是否关联显示
	DeptCheckStrictly int        `json:"dept_check_strictly"` // 部门树选择项是否关联显示
	Status            int        `json:"status"`              // 角色状态
	CreateBy          string     `json:"create_by"`           // 创建者
	CreateTime        *time.Time `json:"create_time"`         // 创建时间
	UpdateBy          string     `json:"update_by"`           // 更新者
	UpdateTime        *time.Time `json:"update_time"`         // 更新时间
	Remark            *string    `json:"remark"`              // 备注
}

// RoleDetailRes 角色详情响应
type RoleDetailRes struct {
	RoleID            uint64     `json:"role_id"`             // 角色ID
	RoleName          string     `json:"role_name"`           // 角色名称
	RoleKey           string     `json:"role_key"`            // 角色权限字符串
	RoleSort          int        `json:"role_sort"`           // 显示顺序
	DataScope         int        `json:"data_scope"`          // 数据范围
	DataScopeName     string     `json:"data_scope_name"`     // 数据范围名称
	MenuCheckStrictly int        `json:"menu_check_strictly"` // 菜单树选择项是否关联显示
	DeptCheckStrictly int        `json:"dept_check_strictly"` // 部门树选择项是否关联显示
	Status            int        `json:"status"`              // 角色状态
	CreateBy          string     `json:"create_by"`           // 创建者
	CreateTime        *time.Time `json:"create_time"`         // 创建时间
	UpdateBy          string     `json:"update_by"`           // 更新者
	UpdateTime        *time.Time `json:"update_time"`         // 更新时间
	Remark            *string    `json:"remark"`              // 备注
	MenuIDs           []uint64   `json:"menu_ids"`            // 菜单ID列表
}

// RoleCreateReq 创建角色请求
type RoleCreateReq struct {
	RoleName          string   `json:"role_name" binding:"required,min=1,max=30" example:"内容管理员"`                        // 角色名称
	RoleKey           string   `json:"role_key" binding:"required,min=1,max=100" example:"content"`                      // 角色权限字符串
	RoleSort          int      `json:"role_sort" binding:"min=0" example:"1"`                                            // 显示顺序
	DataScope         int      `json:"data_scope,omitempty" binding:"omitempty,oneof=1 2 3 4 5" swaggerignore:"true"`    // 数据范围（内部使用，前端无需传递）
	MenuCheckStrictly int      `json:"menu_check_strictly,omitempty" binding:"omitempty,oneof=0 1" swaggerignore:"true"` // 菜单树关联（内部使用，前端无需传递）
	DeptCheckStrictly int      `json:"dept_check_strictly,omitempty" binding:"omitempty,oneof=0 1" swaggerignore:"true"` // 部门树关联（内部使用，前端无需传递）
	Status            int      `json:"status" binding:"oneof=0 1" example:"1"`                                           // 角色状态（0停用 1正常）
	Remark            string   `json:"remark" binding:"max=500" example:"内容管理员角色"`                                       // 备注
	MenuIDs           []uint64 `json:"menu_ids" example:"1,2,3"`                                                         // 菜单ID列表
}

// RoleUpdateReq 更新角色请求
type RoleUpdateReq struct {
	RoleID            uint64   `json:"role_id" binding:"required,min=1" example:"1"`                                     // 角色ID
	RoleName          string   `json:"role_name" binding:"required,min=1,max=30" example:"内容管理员"`                        // 角色名称
	RoleKey           string   `json:"role_key" binding:"required,min=1,max=100" example:"content"`                      // 角色权限字符串
	RoleSort          int      `json:"role_sort" binding:"min=0" example:"1"`                                            // 显示顺序
	DataScope         int      `json:"data_scope,omitempty" binding:"omitempty,oneof=1 2 3 4 5" swaggerignore:"true"`    // 数据范围（内部使用，前端无需传递）
	MenuCheckStrictly int      `json:"menu_check_strictly,omitempty" binding:"omitempty,oneof=0 1" swaggerignore:"true"` // 菜单树关联（内部使用，前端无需传递）
	DeptCheckStrictly int      `json:"dept_check_strictly,omitempty" binding:"omitempty,oneof=0 1" swaggerignore:"true"` // 部门树关联（内部使用，前端无需传递）
	Status            int      `json:"status" binding:"oneof=0 1" example:"1"`                                           // 角色状态（0停用 1正常）
	Remark            string   `json:"remark" binding:"max=500" example:"内容管理员角色"`                                       // 备注
	MenuIDs           []uint64 `json:"menu_ids"`                                                                         // 菜单ID列表
}

// RoleSelectRes 角色选择响应（用于下拉选择）
type RoleSelectRes struct {
	RoleID   uint64 `json:"role_id"`   // 角色ID
	RoleName string `json:"role_name"` // 角色名称
	RoleKey  string `json:"role_key"`  // 角色权限字符串
}

// RoleAuthReq 角色授权请求
type RoleAuthReq struct {
	RoleID  uint64   `json:"role_id" binding:"required,min=1" example:"1"` // 角色ID
	MenuIDs []uint64 `json:"menu_ids"`                                     // 菜单ID列表
}

// AdminRoleReq 管理员角色分配请求
type AdminRoleReq struct {
	AdminID uint64   `json:"admin_id" binding:"required,min=1" example:"1"` // 管理员ID
	RoleIDs []uint64 `json:"role_ids"`                                      // 角色ID列表
}

// AdminRoleRes 管理员角色响应
type AdminRoleRes struct {
	AdminID   uint64           `json:"admin_id"`   // 管理员ID
	AdminName string           `json:"admin_name"` // 管理员账号
	RealName  string           `json:"real_name"`  // 真实姓名
	Email     string           `json:"email"`      // 邮箱
	RoleIDs   []uint64         `json:"role_ids"`   // 角色ID列表
	Roles     []*RoleSelectRes `json:"roles"`      // 角色列表
}

// RoleMenuTreeRes 角色菜单树响应
type RoleMenuTreeRes struct {
	MenuID   uint64             `json:"menu_id"`            // 菜单ID
	ParentID uint64             `json:"parent_id"`          // 父菜单ID
	MenuName string             `json:"menu_name"`          // 菜单名称
	MenuType string             `json:"menu_type"`          // 菜单类型
	Checked  bool               `json:"checked"`            // 是否选中
	Children []*RoleMenuTreeRes `json:"children,omitempty"` // 子菜单列表
}

// RoleMenuAuthReq 角色菜单权限分配请求
type RoleMenuAuthReq struct {
	RoleID  uint64   `json:"role_id" binding:"required,min=1" example:"1"` // 角色ID
	MenuIDs []uint64 `json:"menu_ids" example:"1,2,3"`                     // 菜单ID列表
}
