package services

import (
	"errors"
	"go_server/global"
	"go_server/internal/admin/dto"
	"go_server/models"
	"go_server/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

// AdminLogin 管理员登录
func AdminLogin(req dto.AdminLoginReq) (string, *dto.AdminLoginRes, error) {
	// 查找管理员
	var admin models.SysAdminUser
	err := global.DB.Where("admin_name = ? AND status = ? AND del_flag = ?",
		req.AdminName, models.AdminStatusEnabled, models.AdminDelFlagExist).First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("管理员账号不存在或已被禁用")
		}
		return "", nil, err
	}

	// 检查账号是否被锁定
	if admin.IsLocked() {
		return "", nil, errors.New("账号已被锁定，请30分钟后重试")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password))
	if err != nil {
		// 密码错误，增加错误次数
		admin.PwdErrorCount++
		if admin.PwdErrorCount >= 5 {
			// 错误次数达到5次，锁定账号
			now := time.Now()
			admin.LockTime = &now
		}
		global.DB.Save(&admin)
		return "", nil, errors.New("密码错误")
	}

	// 登录成功，重置错误次数和锁定时间
	admin.PwdErrorCount = 0
	admin.LockTime = nil
	admin.LoginCount++
	now := time.Now()
	admin.LoginDate = &now
	global.DB.Save(&admin)

	// 生成JWT token
	token, err := utils.GenerateJWT(admin.AdminID)
	if err != nil {
		return "", nil, errors.New("生成token失败")
	}

	// 构造返回数据
	adminInfo := &dto.AdminLoginRes{
		AdminID:    admin.AdminID,
		AdminName:  admin.AdminName,
		RealName:   admin.RealName,
		Email:      admin.Email,
		Avatar:     admin.Avatar,
		Department: admin.Department,
		Position:   admin.Position,
		LoginDate:  admin.LoginDate,
	}

	return token, adminInfo, nil
}

// AdminRegister 管理员注册（创建管理员账号）
func AdminRegister(req dto.AdminRegisterReq, creatorName string) error {
	// 检查管理员账号是否已存在
	var count int64
	global.DB.Model(&models.SysAdminUser{}).Where("admin_name = ?", req.AdminName).Count(&count)
	if count > 0 {
		return errors.New("管理员账号已存在")
	}

	// 检查邮箱是否已存在
	global.DB.Model(&models.SysAdminUser{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		return errors.New("邮箱已被使用")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 创建管理员记录
	admin := &models.SysAdminUser{
		AdminName:   req.AdminName,
		RealName:    req.RealName,
		Email:       req.Email,
		Password:    string(hashedPassword),
		PhoneNumber: req.PhoneNumber,
		Department:  req.Department,
		Position:    req.Position,
		Status:      models.AdminStatusEnabled,
		DelFlag:     models.AdminDelFlagExist,
		CreateBy:    creatorName,
		UpdateBy:    creatorName,
		Remark:      &req.Remark,
	}

	return global.DB.Create(admin).Error
}
