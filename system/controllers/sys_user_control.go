package controllers

import (
	"github.com/gin-gonic/gin"
	"go_server/common"
	"go_server/system/dto"
	"go_server/system/services"
	"go_server/utils"
)

// SysUserRegister 用户注册
func SysUserRegister(c *gin.Context) {
	var req dto.RegisterUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, common.ErrorCode, utils.GetErrorMsg(req, err))
		return
	}
	token, err := services.Register(req)
	if err != nil {
		common.Error(c, common.ErrorCode, err.Error())
		return
	}
	common.Success(c, gin.H{"token": token})
}

// SysUserLogin 用户登录
func SysUserLogin(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, common.ErrorCode, utils.GetErrorMsg(req, err))
		return
	}
	token, err := services.Login(req)
	if err != nil {
		common.Error(c, common.ErrorCode, err.Error())
		return
	}
	common.Success(c, gin.H{"token": token})
}

// SysUserInfo 获取用户信息
func SysUserInfo(c *gin.Context) {
	println("获取用户信息")
	userID := c.GetUint64("userID")
	userInfo, err := services.UserInfo(userID)
	if err != nil {
		common.Error(c, common.ErrorCode, err.Error())
		return
	}
	common.Success(c, userInfo)
}
