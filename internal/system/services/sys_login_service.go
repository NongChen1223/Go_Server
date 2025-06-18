package services

import (
	"errors"
	"go_server/global"
	"go_server/internal/system/dto"
	"go_server/models"
	"go_server/utils"
)

type GinJWTMiddleware struct {
}

// Login 用户登录
func Login(req dto.LoginReq) (string, error) {
	var user models.SysUser
	err := global.DB.Where("user_name = ?", req.UserName).First(&user).Error
	if err != nil {
		return "", errors.New("用户不存在")
	}
	ok, err := utils.CompareHashAndPassword(user.Password, req.Password)
	if err != nil || !ok {
		return "", errors.New("密码错误")
	}
	token, err := utils.GenerateJWT(user.UserID)
	if err != nil {
		return "", errors.New("生成token失败")
	}
	return token, nil
}

// Register 用户注册
func Register(req dto.RegisterUserReq) (string, error) {
	var count int64
	global.DB.Model(&models.SysUser{}).Where("user_name = ?", req.UserName).Count(&count)
	if count > 0 {
		return "", errors.New("该用户名已存在")
	}
	global.DB.Model(&models.SysUser{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		return "", errors.New("该邮箱已经注册")
	}
	hashedPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		return "", errors.New("密码加密失败")
	}
	user := &models.SysUser{
		UserName: req.UserName,
		Password: hashedPwd,
		Email:    req.Email,
	}
	if err := global.DB.Create(user).Error; err != nil {
		return "", errors.New("注册失败，请稍后再试")
	}
	token, err := utils.GenerateJWT(user.UserID)
	if err != nil {
		return "", errors.New("生成token失败")
	}
	return token, nil
}

// UserInfo 获取用户信息
func UserInfo(UserID uint64) (*dto.UserInfoRes, error) {
	var user models.SysUser
	if err := global.DB.First(&user, UserID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	res := &dto.UserInfoRes{
		UserID:      user.UserID,
		UserName:    user.UserName,
		NickName:    user.NickName,
		UserType:    user.UserType,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Sex:         user.Sex,
		Avatar:      user.Avatar,
		LoginIP:     user.LoginIP,
		LoginDate:   user.LoginDate,
		CreateTime:  user.CreateTime,
	}
	return res, nil
}
