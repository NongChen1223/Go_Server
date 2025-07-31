package dto

import "time"

// AdminLoginReq 管理员登录请求
type AdminLoginReq struct {
	AdminName string `json:"admin_name" binding:"required" example:"admin"` // 管理员账号
	Password  string `json:"password" binding:"required" example:"123456"`  // 密码
}

// AdminLoginRes 管理员登录响应
type AdminLoginRes struct {
	AdminID    uint64     `json:"admin_id"`   // 管理员ID
	AdminName  string     `json:"admin_name"` // 管理员账号
	RealName   string     `json:"real_name"`  // 真实姓名
	Email      string     `json:"email"`      // 邮箱
	Avatar     string     `json:"avatar"`     // 头像
	Department string     `json:"department"` // 部门
	Position   string     `json:"position"`   // 职位
	LoginDate  *time.Time `json:"login_date"` // 最后登录时间
}

// AdminRegisterReq 管理员注册请求
type AdminRegisterReq struct {
	AdminName   string `json:"admin_name" binding:"required,min=3,max=30" example:"admin"`    // 管理员账号
	RealName    string `json:"real_name" binding:"required,min=2,max=30" example:"张三"`        // 真实姓名
	Email       string `json:"email" binding:"required,email" example:"admin@example.com"`    // 邮箱
	Password    string `json:"password" binding:"required,min=6,max=20" example:"123456"`     // 密码
	PhoneNumber string `json:"phone_number" binding:"omitempty,len=11" example:"13800138000"` // 手机号
	Department  string `json:"department" binding:"omitempty,max=50" example:"技术部"`           // 部门
	Position    string `json:"position" binding:"omitempty,max=50" example:"系统管理员"`           // 职位
	Remark      string `json:"remark" binding:"omitempty,max=500" example:"系统管理员账号"`          // 备注
}

// AdminLogoutRes 管理员登出响应
type AdminLogoutRes struct {
	Message string `json:"message" example:"登出成功"`
}

// AdminInfoRes 管理员信息响应
type AdminInfoRes struct {
	AdminInfo *AdminInfoDetail   `json:"admin_info"` // 管理员基本信息
	Roles     []*AdminRoleInfo   `json:"roles"`      // 角色信息
	Menus     []*AdminMenuRouter `json:"menus"`      // 菜单路由
	Perms     []string           `json:"perms"`      // 权限列表
}

// AdminInfoDetail 管理员详细信息
type AdminInfoDetail struct {
	AdminID   uint64 `json:"admin_id"`   // 管理员ID
	AdminName string `json:"admin_name"` // 管理员账号
	RealName  string `json:"real_name"`  // 真实姓名
	Email     string `json:"email"`      // 邮箱
	Status    int    `json:"status"`     // 状态
}

// AdminRoleInfo 管理员角色信息
type AdminRoleInfo struct {
	RoleID   uint64 `json:"role_id"`   // 角色ID
	RoleName string `json:"role_name"` // 角色名称
	RoleKey  string `json:"role_key"`  // 角色权限字符串
}

// AdminMenuRouter 管理员菜单路由
type AdminMenuRouter struct {
	Name      string             `json:"name"`                // 路由名称
	Path      string             `json:"path"`                // 路由路径
	Hidden    bool               `json:"hidden"`              // 是否隐藏路由
	Redirect  string             `json:"redirect,omitempty"`  // 重定向地址
	Component string             `json:"component,omitempty"` // 组件路径
	Query     string             `json:"query,omitempty"`     // 路由参数
	Meta      *AdminRouterMeta   `json:"meta,omitempty"`      // 路由元信息
	Children  []*AdminMenuRouter `json:"children,omitempty"`  // 子路由
}

// AdminRouterMeta 管理员路由元信息
type AdminRouterMeta struct {
	Title   string `json:"title"`             // 路由标题
	Icon    string `json:"icon,omitempty"`    // 图标
	NoCache bool   `json:"noCache,omitempty"` // 是否不缓存
	Link    string `json:"link,omitempty"`    // 外链地址
}
