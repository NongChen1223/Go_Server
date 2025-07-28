package services

import (
	"errors"
	"go_server/common"
	"go_server/global"
	"go_server/internal/admin/dto"
	"go_server/models"
	"gorm.io/gorm"
	"time"
)

// GetRoleList 获取角色列表
func GetRoleList(req dto.RoleListReq) (*common.PageResult, error) {
	var roles []models.SysRole
	var total int64

	query := global.DB.Model(&models.SysRole{}).Where("del_flag = ?", models.RoleDelFlagExist)

	// 条件查询
	if req.RoleName != "" {
		query = query.Where("role_name LIKE ?", "%"+req.RoleName+"%")
	}
	if req.RoleKey != "" {
		query = query.Where("role_key LIKE ?", "%"+req.RoleKey+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.Size
	err = query.Order("role_sort ASC, role_id ASC").Offset(offset).Limit(req.Size).Find(&roles).Error
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	var roleList []*dto.RoleListRes
	for _, role := range roles {
		roleRes := &dto.RoleListRes{
			RoleID:            role.RoleID,
			RoleName:          role.RoleName,
			RoleKey:           role.RoleKey,
			RoleSort:          role.RoleSort,
			DataScope:         role.DataScope,
			DataScopeName:     role.GetDataScopeName(),
			MenuCheckStrictly: role.MenuCheckStrictly,
			DeptCheckStrictly: role.DeptCheckStrictly,
			Status:            role.Status,
			CreateBy:          role.CreateBy,
			CreateTime:        role.CreateTime,
			UpdateBy:          role.UpdateBy,
			UpdateTime:        role.UpdateTime,
			Remark:            role.Remark,
		}
		roleList = append(roleList, roleRes)
	}

	return &common.PageResult{
		Current: req.Page,
		Size:    req.Size,
		Total:   int(total),
		Records: roleList,
	}, nil
}

// GetRoleDetail 获取角色详情
func GetRoleDetail(roleID uint64) (*dto.RoleDetailRes, error) {
	var role models.SysRole
	err := global.DB.Where("role_id = ? AND del_flag = ?", roleID, models.RoleDelFlagExist).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("角色不存在")
		}
		return nil, err
	}

	// 获取角色的菜单权限
	var menuIDs []uint64
	err = global.DB.Model(&models.SysRoleMenu{}).Where("role_id = ?", roleID).Pluck("menu_id", &menuIDs).Error
	if err != nil {
		return nil, err
	}

	return &dto.RoleDetailRes{
		RoleID:            role.RoleID,
		RoleName:          role.RoleName,
		RoleKey:           role.RoleKey,
		RoleSort:          role.RoleSort,
		DataScope:         role.DataScope,
		DataScopeName:     role.GetDataScopeName(),
		MenuCheckStrictly: role.MenuCheckStrictly,
		DeptCheckStrictly: role.DeptCheckStrictly,
		Status:            role.Status,
		CreateBy:          role.CreateBy,
		CreateTime:        role.CreateTime,
		UpdateBy:          role.UpdateBy,
		UpdateTime:        role.UpdateTime,
		Remark:            role.Remark,
		MenuIDs:           menuIDs,
	}, nil
}

// CreateRole 创建角色
func CreateRole(req dto.RoleCreateReq, createBy string) error {
	// 检查角色名称是否重复
	var count int64
	err := global.DB.Model(&models.SysRole{}).Where("role_name = ? AND del_flag = ?", req.RoleName, models.RoleDelFlagExist).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("角色名称已存在")
	}

	// 检查角色权限字符串是否重复
	err = global.DB.Model(&models.SysRole{}).Where("role_key = ? AND del_flag = ?", req.RoleKey, models.RoleDelFlagExist).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("角色权限字符串已存在")
	}

	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()
	role := models.SysRole{
		RoleName:          req.RoleName,
		RoleKey:           req.RoleKey,
		RoleSort:          req.RoleSort,
		DataScope:         req.DataScope,
		MenuCheckStrictly: req.MenuCheckStrictly,
		DeptCheckStrictly: req.DeptCheckStrictly,
		Status:            req.Status,
		DelFlag:           models.RoleDelFlagExist,
		CreateBy:          createBy,
		CreateTime:        &now,
		UpdateBy:          createBy,
		UpdateTime:        &now,
	}

	if req.Remark != "" {
		role.Remark = &req.Remark
	}

	// 创建角色
	err = tx.Create(&role).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 分配菜单权限
	if len(req.MenuIDs) > 0 {
		var roleMenus []models.SysRoleMenu
		for _, menuID := range req.MenuIDs {
			roleMenus = append(roleMenus, models.SysRoleMenu{
				RoleID: role.RoleID,
				MenuID: menuID,
			})
		}
		err = tx.Create(&roleMenus).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// UpdateRole 更新角色
func UpdateRole(req dto.RoleUpdateReq, updateBy string) error {
	// 检查角色是否存在
	var role models.SysRole
	err := global.DB.Where("role_id = ? AND del_flag = ?", req.RoleID, models.RoleDelFlagExist).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		return err
	}

	// 超级管理员角色不允许修改
	if role.IsAdmin() {
		return errors.New("超级管理员角色不允许修改")
	}

	// 检查角色名称是否重复（排除自己）
	var count int64
	err = global.DB.Model(&models.SysRole{}).Where("role_name = ? AND role_id != ? AND del_flag = ?", req.RoleName, req.RoleID, models.RoleDelFlagExist).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("角色名称已存在")
	}

	// 检查角色权限字符串是否重复（排除自己）
	err = global.DB.Model(&models.SysRole{}).Where("role_key = ? AND role_id != ? AND del_flag = ?", req.RoleKey, req.RoleID, models.RoleDelFlagExist).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("角色权限字符串已存在")
	}

	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()
	updates := map[string]interface{}{
		"role_name":           req.RoleName,
		"role_key":            req.RoleKey,
		"role_sort":           req.RoleSort,
		"data_scope":          req.DataScope,
		"menu_check_strictly": req.MenuCheckStrictly,
		"dept_check_strictly": req.DeptCheckStrictly,
		"status":              req.Status,
		"update_by":           updateBy,
		"update_time":         &now,
	}

	if req.Remark != "" {
		updates["remark"] = req.Remark
	} else {
		updates["remark"] = nil
	}

	// 更新角色
	err = tx.Model(&role).Updates(updates).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 删除原有的菜单权限
	err = tx.Where("role_id = ?", req.RoleID).Delete(&models.SysRoleMenu{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 重新分配菜单权限
	if len(req.MenuIDs) > 0 {
		var roleMenus []models.SysRoleMenu
		for _, menuID := range req.MenuIDs {
			roleMenus = append(roleMenus, models.SysRoleMenu{
				RoleID: req.RoleID,
				MenuID: menuID,
			})
		}
		err = tx.Create(&roleMenus).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// DeleteRole 删除角色
func DeleteRole(roleID uint64) error {
	// 检查角色是否存在
	var role models.SysRole
	err := global.DB.Where("role_id = ? AND del_flag = ?", roleID, models.RoleDelFlagExist).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		return err
	}

	// 超级管理员角色不允许删除
	if role.IsAdmin() {
		return errors.New("超级管理员角色不允许删除")
	}

	// 检查是否有管理员在使用此角色
	var adminRoleCount int64
	err = global.DB.Model(&models.SysAdminUserRole{}).Where("role_id = ?", roleID).Count(&adminRoleCount).Error
	if err != nil {
		return err
	}
	if adminRoleCount > 0 {
		return errors.New("角色已分配给管理员，不允许删除")
	}

	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 软删除角色
	now := time.Now()
	err = tx.Model(&role).Updates(map[string]interface{}{
		"del_flag":    models.RoleDelFlagDeleted,
		"update_time": &now,
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 删除角色菜单关联
	err = tx.Where("role_id = ?", roleID).Delete(&models.SysRoleMenu{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetRoleSelect 获取角色选择列表
func GetRoleSelect() ([]*dto.RoleSelectRes, error) {
	var roles []models.SysRole
	err := global.DB.Where("status = ? AND del_flag = ?", models.RoleStatusEnabled, models.RoleDelFlagExist).
		Order("role_sort ASC, role_id ASC").Find(&roles).Error
	if err != nil {
		return nil, err
	}

	var roleList []*dto.RoleSelectRes
	for _, role := range roles {
		roleRes := &dto.RoleSelectRes{
			RoleID:   role.RoleID,
			RoleName: role.RoleName,
			RoleKey:  role.RoleKey,
		}
		roleList = append(roleList, roleRes)
	}

	return roleList, nil
}

// AuthRole 角色授权
func AuthRole(req dto.RoleAuthReq, updateBy string) error {
	// 检查角色是否存在
	var role models.SysRole
	err := global.DB.Where("role_id = ? AND del_flag = ?", req.RoleID, models.RoleDelFlagExist).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		return err
	}

	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除原有的菜单权限
	err = tx.Where("role_id = ?", req.RoleID).Delete(&models.SysRoleMenu{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 重新分配菜单权限
	if len(req.MenuIDs) > 0 {
		var roleMenus []models.SysRoleMenu
		for _, menuID := range req.MenuIDs {
			roleMenus = append(roleMenus, models.SysRoleMenu{
				RoleID: req.RoleID,
				MenuID: menuID,
			})
		}
		err = tx.Create(&roleMenus).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// 更新角色的更新时间
	now := time.Now()
	err = tx.Model(&role).Updates(map[string]interface{}{
		"update_by":   updateBy,
		"update_time": &now,
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// AssignAdminRole 分配管理员角色
func AssignAdminRole(req dto.AdminRoleReq, updateBy string) error {
	// 检查管理员是否存在
	var admin models.SysAdminUser
	err := global.DB.Where("admin_id = ? AND del_flag = ?", req.AdminID, models.AdminDelFlagExist).First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("管理员不存在")
		}
		return err
	}

	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除原有的角色关联
	err = tx.Where("admin_id = ?", req.AdminID).Delete(&models.SysAdminUserRole{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 重新分配角色
	if len(req.RoleIDs) > 0 {
		var adminRoles []models.SysAdminUserRole
		for _, roleID := range req.RoleIDs {
			adminRoles = append(adminRoles, models.SysAdminUserRole{
				AdminID: req.AdminID,
				RoleID:  roleID,
			})
		}
		err = tx.Create(&adminRoles).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// 更新管理员的更新时间
	now := time.Now()
	err = tx.Model(&admin).Updates(map[string]interface{}{
		"update_by":   updateBy,
		"update_time": &now,
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetAdminRole 获取管理员角色信息
func GetAdminRole(adminID uint64) (*dto.AdminRoleRes, error) {
	// 获取管理员信息
	var admin models.SysAdminUser
	err := global.DB.Where("admin_id = ? AND del_flag = ?", adminID, models.AdminDelFlagExist).First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("管理员不存在")
		}
		return nil, err
	}

	// 获取管理员的角色ID列表
	var roleIDs []uint64
	err = global.DB.Model(&models.SysAdminUserRole{}).Where("admin_id = ?", adminID).Pluck("role_id", &roleIDs).Error
	if err != nil {
		return nil, err
	}

	// 获取角色详细信息
	var roles []*dto.RoleSelectRes
	if len(roleIDs) > 0 {
		var roleModels []models.SysRole
		err = global.DB.Where("role_id IN ? AND del_flag = ?", roleIDs, models.RoleDelFlagExist).Find(&roleModels).Error
		if err != nil {
			return nil, err
		}

		for _, role := range roleModels {
			roles = append(roles, &dto.RoleSelectRes{
				RoleID:   role.RoleID,
				RoleName: role.RoleName,
				RoleKey:  role.RoleKey,
			})
		}
	}

	return &dto.AdminRoleRes{
		AdminID:   admin.AdminID,
		AdminName: admin.AdminName,
		RealName:  admin.RealName,
		Email:     admin.Email,
		RoleIDs:   roleIDs,
		Roles:     roles,
	}, nil
}
