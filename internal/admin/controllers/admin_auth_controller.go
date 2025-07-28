package controllers

import (
	"github.com/gin-gonic/gin"
	"go_server/common"
	"go_server/constants"
	"go_server/internal/admin/dto"
	"go_server/internal/admin/services"
)

// AdminLogin 管理员登录
// @Summary 管理员登录
// @Description 管理员使用账号密码登录系统，返回JWT token
// @Tags Admin-管理员认证
// @Accept json
// @Produce json
// @Param request body dto.AdminLoginReq true "登录请求参数"
// @Success 200 {object} common.Response{data=map[string]interface{}} "登录成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "认证失败"
// @Router /v1/admin/login [post]
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
// @Summary 创建管理员账号
// @Description 创建新的管理员账号（需要管理员权限）
// @Tags Admin-管理员认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.AdminRegisterReq true "注册请求参数"
// @Success 200 {object} common.Response "创建成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Router /v1/admin/register [post]
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

// AdminLogout 管理员登出
// @Summary 管理员登出
// @Description 管理员登出系统，清除token
// @Tags Admin-管理员认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.Response "登出成功"
// @Failure 401 {object} common.Response "未授权"
// @Router /v1/admin/logout [post]
func AdminLogout(c *gin.Context) {
	// 这里可以实现token黑名单逻辑
	// 目前简单返回成功，实际项目中可以将token加入黑名单
	common.Success(c, gin.H{"message": "登出成功"})
}
