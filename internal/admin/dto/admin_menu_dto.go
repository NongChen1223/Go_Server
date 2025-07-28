package dto

import "time"

// MenuListReq 菜单列表查询请求
type MenuListReq struct {
	MenuName string `form:"menu_name" example:"用户管理"`                             // 菜单名称
	Status   *int   `form:"status" example:"1"`                                   // 菜单状态（0停用 1正常）
	MenuType string `form:"menu_type" example:"C"`                                // 菜单类型（M目录 C菜单 F按钮）
	ParentID *int   `form:"parent_id" example:"0"`                                // 父菜单ID
	Visible  *int   `form:"visible" example:"1"`                                  // 菜单状态（0隐藏 1显示）
	Page     int    `form:"page,default=1" binding:"min=1" example:"1"`           // 页码
	Size     int    `form:"size,default=10" binding:"min=1,max=100" example:"10"` // 每页数量
}

// MenuListRes 菜单列表响应
type MenuListRes struct {
	MenuID     uint64         `json:"menu_id"`            // 菜单ID
	ParentID   uint64         `json:"parent_id"`          // 父菜单ID
	Ancestors  string         `json:"ancestors"`          // 祖级列表
	MenuName   string         `json:"menu_name"`          // 菜单名称
	OrderNum   int            `json:"order_num"`          // 显示顺序
	Icon       string         `json:"icon"`               // 菜单图标
	Path       string         `json:"path"`               // 路由地址
	Component  *string        `json:"component"`          // 组件路径
	Query      *string        `json:"query"`              // 路由参数
	MenuType   string         `json:"menu_type"`          // 菜单类型
	Visible    int            `json:"visible"`            // 菜单状态（0隐藏 1显示）
	Status     int            `json:"status"`             // 菜单状态（0停用 1正常）
	Perms      *string        `json:"perms"`              // 权限标识
	IsFrame    int            `json:"is_frame"`           // 是否为外链
	IsCache    int            `json:"is_cache"`           // 是否缓存
	CreateBy   string         `json:"create_by"`          // 创建者
	CreateTime *time.Time     `json:"create_time"`        // 创建时间
	UpdateBy   string         `json:"update_by"`          // 更新者
	UpdateTime *time.Time     `json:"update_time"`        // 更新时间
	Remark     string         `json:"remark"`             // 备注
	Children   []*MenuListRes `json:"children,omitempty"` // 子菜单列表
}

// MenuDetailRes 菜单详情响应
type MenuDetailRes struct {
	MenuID     uint64     `json:"menu_id"`     // 菜单ID
	ParentID   uint64     `json:"parent_id"`   // 父菜单ID
	Ancestors  string     `json:"ancestors"`   // 祖级列表
	MenuName   string     `json:"menu_name"`   // 菜单名称
	OrderNum   int        `json:"order_num"`   // 显示顺序
	Icon       string     `json:"icon"`        // 菜单图标
	Path       string     `json:"path"`        // 路由地址
	Component  *string    `json:"component"`   // 组件路径
	Query      *string    `json:"query"`       // 路由参数
	MenuType   string     `json:"menu_type"`   // 菜单类型
	Visible    int        `json:"visible"`     // 菜单状态（0隐藏 1显示）
	Status     int        `json:"status"`      // 菜单状态（0停用 1正常）
	Perms      *string    `json:"perms"`       // 权限标识
	IsFrame    int        `json:"is_frame"`    // 是否为外链
	IsCache    int        `json:"is_cache"`    // 是否缓存
	CreateBy   string     `json:"create_by"`   // 创建者
	CreateTime *time.Time `json:"create_time"` // 创建时间
	UpdateBy   string     `json:"update_by"`   // 更新者
	UpdateTime *time.Time `json:"update_time"` // 更新时间
	Remark     string     `json:"remark"`      // 备注
}

// MenuCreateReq 创建菜单请求
type MenuCreateReq struct {
	ParentID  uint64  `json:"parent_id" binding:"min=0" example:"0"`                             // 父菜单ID
	MenuName  string  `json:"menu_name" binding:"required,min=1,max=50" example:"用户管理"`          // 菜单名称
	OrderNum  int     `json:"order_num" binding:"min=0" example:"1"`                             // 显示顺序
	Icon      string  `json:"icon" binding:"max=100" example:"user"`                             // 菜单图标
	Path      string  `json:"path" binding:"max=200" example:"/system/user"`                     // 路由地址
	Component *string `json:"component" binding:"omitempty,max=255" example:"system/user/index"` // 组件路径
	Query     *string `json:"query" binding:"omitempty,max=255" example:"userId=1"`              // 路由参数
	MenuType  string  `json:"menu_type" binding:"required,oneof=M C F" example:"C"`              // 菜单类型（M目录 C菜单 F按钮）
	Visible   int     `json:"visible" binding:"oneof=0 1" example:"1"`                           // 菜单状态（0隐藏 1显示）
	Status    int     `json:"status" binding:"oneof=0 1" example:"1"`                            // 菜单状态（0停用 1正常）
	Perms     *string `json:"perms" binding:"omitempty,max=100" example:"system:user:list"`      // 权限标识
	IsFrame   int     `json:"is_frame" binding:"oneof=0 1" example:"1"`                          // 是否为外链（0是 1否）
	IsCache   int     `json:"is_cache" binding:"oneof=0 1" example:"0"`                          // 是否缓存（0不缓存 1缓存）
	Remark    string  `json:"remark" binding:"max=500" example:"用户管理菜单"`                         // 备注
}

// MenuUpdateReq 更新菜单请求
type MenuUpdateReq struct {
	MenuID    uint64  `json:"menu_id" binding:"required,min=1" example:"1"`                      // 菜单ID
	ParentID  uint64  `json:"parent_id" binding:"min=0" example:"0"`                             // 父菜单ID
	MenuName  string  `json:"menu_name" binding:"required,min=1,max=50" example:"用户管理"`          // 菜单名称
	OrderNum  int     `json:"order_num" binding:"min=0" example:"1"`                             // 显示顺序
	Icon      string  `json:"icon" binding:"max=100" example:"user"`                             // 菜单图标
	Path      string  `json:"path" binding:"max=200" example:"/system/user"`                     // 路由地址
	Component *string `json:"component" binding:"omitempty,max=255" example:"system/user/index"` // 组件路径
	Query     *string `json:"query" binding:"omitempty,max=255" example:"userId=1"`              // 路由参数
	MenuType  string  `json:"menu_type" binding:"required,oneof=M C F" example:"C"`              // 菜单类型（M目录 C菜单 F按钮）
	Visible   int     `json:"visible" binding:"oneof=0 1" example:"1"`                           // 菜单状态（0隐藏 1显示）
	Status    int     `json:"status" binding:"oneof=0 1" example:"1"`                            // 菜单状态（0停用 1正常）
	Perms     *string `json:"perms" binding:"omitempty,max=100" example:"system:user:list"`      // 权限标识
	IsFrame   int     `json:"is_frame" binding:"oneof=0 1" example:"1"`                          // 是否为外链（0是 1否）
	IsCache   int     `json:"is_cache" binding:"oneof=0 1" example:"0"`                          // 是否缓存（0不缓存 1缓存）
	Remark    string  `json:"remark" binding:"max=500" example:"用户管理菜单"`                         // 备注
}

// MenuTreeRes 菜单树响应（用于角色分配权限时的树形选择）
type MenuTreeRes struct {
	MenuID   uint64         `json:"menu_id"`            // 菜单ID
	ParentID uint64         `json:"parent_id"`          // 父菜单ID
	MenuName string         `json:"menu_name"`          // 菜单名称
	MenuType string         `json:"menu_type"`          // 菜单类型
	Children []*MenuTreeRes `json:"children,omitempty"` // 子菜单列表
}

// MenuRouterRes 菜单路由响应（用于前端动态路由生成）
type MenuRouterRes struct {
	Name      string           `json:"name"`                // 路由名称
	Path      string           `json:"path"`                // 路由路径
	Hidden    bool             `json:"hidden"`              // 是否隐藏路由
	Redirect  string           `json:"redirect,omitempty"`  // 重定向地址
	Component string           `json:"component,omitempty"` // 组件路径
	Query     string           `json:"query,omitempty"`     // 路由参数
	Meta      *MenuRouterMeta  `json:"meta,omitempty"`      // 路由元信息
	Children  []*MenuRouterRes `json:"children,omitempty"`  // 子路由
}

// MenuRouterMeta 路由元信息
type MenuRouterMeta struct {
	Title   string `json:"title"`             // 路由标题
	Icon    string `json:"icon,omitempty"`    // 图标
	NoCache bool   `json:"noCache,omitempty"` // 是否不缓存
	Link    string `json:"link,omitempty"`    // 外链地址
}
