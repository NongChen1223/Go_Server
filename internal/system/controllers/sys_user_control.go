package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_server/common"
	"go_server/constants"
	"go_server/internal/system/dto"
	"go_server/internal/system/services"
	"go_server/utils"
)

// SysUserRegister 用户注册
// @Summary 用户注册
// @Description 前台用户注册新账号
// @Tags System-用户认证
// @Accept json
// @Produce json
// @Param request body dto.RegisterUserReq true "注册请求参数"
// @Success 200 {object} common.Response{data=map[string]string} "注册成功，返回token"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 500 {object} common.Response "服务器错误"
// @Router /v1/api/auth/register [post]
func SysUserRegister(c *gin.Context) {
	var req dto.RegisterUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, constants.ErrorCode, utils.GetErrorMsg(req, err))
		return
	}
	token, err := services.Register(req)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}
	common.Success(c, gin.H{"token": token})
}

// SysUserLogin 用户登录
// @Summary 用户登录
// @Description 前台用户使用账号密码登录系统
// @Tags System-用户认证
// @Accept json
// @Produce json
// @Param request body dto.LoginReq true "登录请求参数"
// @Success 200 {object} common.Response{data=map[string]string} "登录成功，返回token"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "认证失败"
// @Router /v1/api/auth/login [post]
func SysUserLogin(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		// 打印详细错误信息以便调试
		fmt.Printf("登录参数绑定错误: %v\n", err)
		common.Error(c, constants.ErrorCode, utils.GetErrorMsg(req, err))
		return
	}
	token, err := services.Login(req)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}
	common.Success(c, gin.H{"token": token})
}

// SysUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的详细信息
// @Tags System-用户认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.Response{data=dto.UserInfoRes} "获取成功"
// @Failure 401 {object} common.Response "未授权"
// @Failure 500 {object} common.Response "服务器错误"
// @Router /v1/api/user/info [get]
func SysUserInfo(c *gin.Context) {
	println("获取用户信息")
	UserID := c.GetUint64("user_id")
	userInfo, err := services.UserInfo(UserID)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}
	common.Success(c, userInfo)
}
