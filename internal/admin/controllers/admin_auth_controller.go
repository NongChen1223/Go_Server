package controllers

import (
	"github.com/gin-gonic/gin"
	"go_server/common"
	"go_server/constants"
	"go_server/internal/admin/dto"
	"go_server/internal/admin/services"
)

// AdminLogin 管理员登录
func AdminLogin(c *gin.Context) {
	var req dto.AdminLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, constants.ErrorCode, "参数错误: "+err.Error())
		return
	}

	// 调用登录服务
	token, adminInfo, err := services.AdminLogin(req)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	// 返回登录成功信息
	common.Success(c, map[string]interface{}{
		"token":      token,
		"admin_info": adminInfo,
	})
}

// AdminRegister 管理员注册（创建管理员账号）
func AdminRegister(c *gin.Context) {
	var req dto.AdminRegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, constants.ErrorCode, "参数错误: "+err.Error())
		return
	}

	// 获取当前操作的管理员名称（如果是已登录的管理员创建）
	adminName := "system" // 默认系统创建
	if name, exists := c.Get("admin_name"); exists {
		adminName = name.(string)
	}

	// 调用注册服务
	err := services.AdminRegister(req, adminName)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, nil)
}
