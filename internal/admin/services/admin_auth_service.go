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

// GetAdminInfo 获取管理员信息
func GetAdminInfo(adminID uint64) (*dto.AdminInfoRes, error) {
	// 获取管理员基本信息
	var admin models.SysAdminUser
	err := global.DB.Where("admin_id = ? AND del_flag = ?", adminID, models.AdminDelFlagExist).First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("管理员不存在")
		}
		return nil, err
	}

	// 获取管理员的角色信息
	var roles []models.SysRole
	query := `
		SELECT r.* FROM sys_role r
		LEFT JOIN sys_admin_user_role aur ON r.role_id = aur.role_id
		WHERE aur.admin_id = ? AND r.status = 1 AND r.del_flag = '0'
		ORDER BY r.role_sort ASC
	`
	err = global.DB.Raw(query, adminID).Scan(&roles).Error
	if err != nil {
		return nil, err
	}

	// 获取管理员的菜单路由
	var menus []models.SysMenu
	menuQuery := `
		SELECT DISTINCT m.* FROM sys_menu m
		LEFT JOIN sys_role_menu rm ON m.menu_id = rm.menu_id
		LEFT JOIN sys_admin_user_role aur ON rm.role_id = aur.role_id
		LEFT JOIN sys_role r ON aur.role_id = r.role_id
		WHERE aur.admin_id = ? AND m.status = 1 AND m.visible = 1 AND r.status = 1 AND r.del_flag = '0'
		AND m.menu_type IN ('M', 'C')
		ORDER BY m.order_num ASC, m.menu_id ASC
	`
	err = global.DB.Raw(menuQuery, adminID).Scan(&menus).Error
	if err != nil {
		return nil, err
	}

	// 获取管理员的权限列表
	var perms []string
	permQuery := `
		SELECT DISTINCT m.perms FROM sys_menu m
		LEFT JOIN sys_role_menu rm ON m.menu_id = rm.menu_id
		LEFT JOIN sys_admin_user_role aur ON rm.role_id = aur.role_id
		LEFT JOIN sys_role r ON aur.role_id = r.role_id
		WHERE aur.admin_id = ? AND m.status = 1 AND r.status = 1 AND r.del_flag = '0'
		AND m.perms IS NOT NULL AND m.perms != ''
	`
	err = global.DB.Raw(permQuery, adminID).Pluck("perms", &perms).Error
	if err != nil {
		return nil, err
	}

	// 构建响应数据
	adminInfo := &dto.AdminInfoDetail{
		AdminID:   admin.AdminID,
		AdminName: admin.AdminName,
		RealName:  admin.RealName,
		Email:     admin.Email,
		Status:    admin.Status,
	}

	// 转换角色信息
	var roleInfos []*dto.AdminRoleInfo
	for _, role := range roles {
		roleInfos = append(roleInfos, &dto.AdminRoleInfo{
			RoleID:   role.RoleID,
			RoleName: role.RoleName,
			RoleKey:  role.RoleKey,
		})
	}

	// 构建菜单路由树
	menuRouters := buildAdminMenuRouters(menus)

	return &dto.AdminInfoRes{
		AdminInfo: adminInfo,
		Roles:     roleInfos,
		Menus:     menuRouters,
		Perms:     perms,
	}, nil
}

// buildAdminMenuRouters 构建管理员菜单路由树
func buildAdminMenuRouters(menus []models.SysMenu) []*dto.AdminMenuRouter {
	// 转换为路由DTO并构建树形结构
	routerMap := make(map[uint64]*dto.AdminMenuRouter)
	var rootRouters []*dto.AdminMenuRouter

	// 先转换所有菜单为路由DTO
	for _, menu := range menus {
		router := &dto.AdminMenuRouter{
			Name:      menu.MenuName,
			Path:      menu.Path,
			Hidden:    menu.Visible == models.MenuVisibleHidden,
			Component: getAdminComponent(menu),
			Children:  []*dto.AdminMenuRouter{},
		}

		// 设置路由参数
		if menu.Query != nil && *menu.Query != "" {
			router.Query = *menu.Query
		}

		// 设置路由元信息
		router.Meta = &dto.AdminRouterMeta{
			Title:   menu.MenuName,
			Icon:    menu.Icon,
			NoCache: menu.IsCache == models.MenuNotCache,
		}

		// 如果是外链
		if menu.IsFrame == models.MenuIsFrame {
			router.Meta.Link = menu.Path
		}

		routerMap[menu.MenuID] = router
	}

	// 构建树形结构
	for _, router := range routerMap {
		// 通过菜单ID找到对应的菜单信息
		var menu models.SysMenu
		for _, m := range menus {
			if routerMap[m.MenuID] == router {
				menu = m
				break
			}
		}

		if menu.ParentID == 0 {
			// 顶级路由
			rootRouters = append(rootRouters, router)
		} else {
			// 子路由，添加到父路由的children中
			if parent, exists := routerMap[menu.ParentID]; exists {
				parent.Children = append(parent.Children, router)
			}
		}
	}

	return rootRouters
}

// getAdminComponent 获取管理员路由组件路径
func getAdminComponent(menu models.SysMenu) string {
	if menu.Component != nil && *menu.Component != "" {
		return *menu.Component
	}

	// 如果是目录且没有指定组件，返回Layout
	if menu.MenuType == models.MenuTypeDir {
		return "Layout"
	}

	return ""
}
