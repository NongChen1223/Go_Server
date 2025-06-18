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
