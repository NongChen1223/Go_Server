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
