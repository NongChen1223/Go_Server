package services

import (
	"errors"
	"go_server/global"
	"go_server/internal/admin/dto"
	"go_server/models"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

// GetMenuList 获取菜单列表（树形结构）
func GetMenuList(req dto.MenuListReq) ([]*dto.MenuListRes, error) {
	var menus []models.SysMenu
	query := global.DB.Model(&models.SysMenu{})

	// 条件查询
	if req.MenuName != "" {
		query = query.Where("menu_name LIKE ?", "%"+req.MenuName+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.MenuType != "" {
		query = query.Where("menu_type = ?", req.MenuType)
	}
	if req.ParentID != nil {
		query = query.Where("parent_id = ?", *req.ParentID)
	}
	if req.Visible != nil {
		query = query.Where("visible = ?", *req.Visible)
	}

	// 按排序字段和菜单ID排序
	err := query.Order("order_num ASC, menu_id ASC").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	// 转换为DTO并构建树形结构
	menuMap := make(map[uint64]*dto.MenuListRes)
	var rootMenus []*dto.MenuListRes

	// 先转换所有菜单为DTO
	for _, menu := range menus {
		menuRes := &dto.MenuListRes{
			MenuID:     menu.MenuID,
			ParentID:   menu.ParentID,
			Ancestors:  menu.Ancestors,
			MenuName:   menu.MenuName,
			OrderNum:   menu.OrderNum,
			Icon:       menu.Icon,
			Path:       menu.Path,
			Component:  menu.Component,
			Query:      menu.Query,
			MenuType:   menu.MenuType,
			Visible:    menu.Visible,
			Status:     menu.Status,
			Perms:      menu.Perms,
			IsFrame:    menu.IsFrame,
			IsCache:    menu.IsCache,
			CreateBy:   menu.CreateBy,
			CreateTime: menu.CreateTime,
			UpdateBy:   menu.UpdateBy,
			UpdateTime: menu.UpdateTime,
			Remark:     menu.Remark,
			Children:   []*dto.MenuListRes{},
		}
		menuMap[menu.MenuID] = menuRes
	}

	// 构建树形结构
	for _, menuRes := range menuMap {
		if menuRes.ParentID == 0 {
			// 顶级菜单
			rootMenus = append(rootMenus, menuRes)
		} else {
			// 子菜单，添加到父菜单的children中
			if parent, exists := menuMap[menuRes.ParentID]; exists {
				parent.Children = append(parent.Children, menuRes)
			}
		}
	}

	return rootMenus, nil
}

// GetMenuDetail 获取菜单详情
func GetMenuDetail(menuID uint64) (*dto.MenuDetailRes, error) {
	var menu models.SysMenu
	err := global.DB.First(&menu, menuID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("菜单不存在")
		}
		return nil, err
	}

	return &dto.MenuDetailRes{
		MenuID:     menu.MenuID,
		ParentID:   menu.ParentID,
		Ancestors:  menu.Ancestors,
		MenuName:   menu.MenuName,
		OrderNum:   menu.OrderNum,
		Icon:       menu.Icon,
		Path:       menu.Path,
		Component:  menu.Component,
		Query:      menu.Query,
		MenuType:   menu.MenuType,
		Visible:    menu.Visible,
		Status:     menu.Status,
		Perms:      menu.Perms,
		IsFrame:    menu.IsFrame,
		IsCache:    menu.IsCache,
		CreateBy:   menu.CreateBy,
		CreateTime: menu.CreateTime,
		UpdateBy:   menu.UpdateBy,
		UpdateTime: menu.UpdateTime,
		Remark:     menu.Remark,
	}, nil
}

// CreateMenu 创建菜单
func CreateMenu(req dto.MenuCreateReq, createBy string) error {
	// 检查父菜单是否存在
	if req.ParentID != 0 {
		var parentMenu models.SysMenu
		err := global.DB.First(&parentMenu, req.ParentID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("父菜单不存在")
			}
			return err
		}
	}

	// 检查菜单名称是否重复（同级菜单名称不能重复）
	var count int64
	err := global.DB.Model(&models.SysMenu{}).Where("parent_id = ? AND menu_name = ?", req.ParentID, req.MenuName).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("同级菜单名称不能重复")
	}

	// 构建祖级列表
	ancestors := "0"
	if req.ParentID != 0 {
		var parentMenu models.SysMenu
		global.DB.First(&parentMenu, req.ParentID)
		if parentMenu.Ancestors != "" {
			ancestors = parentMenu.Ancestors + "," + strconv.FormatUint(req.ParentID, 10)
		} else {
			ancestors = "0," + strconv.FormatUint(req.ParentID, 10)
		}
	}

	now := time.Now()
	menu := models.SysMenu{
		ParentID:   req.ParentID,
		Ancestors:  ancestors,
		MenuName:   req.MenuName,
		OrderNum:   req.OrderNum,
		Icon:       req.Icon,
		Path:       req.Path,
		Component:  req.Component,
		Query:      req.Query,
		MenuType:   req.MenuType,
		Visible:    req.Visible,
		Status:     req.Status,
		Perms:      req.Perms,
		IsFrame:    req.IsFrame,
		IsCache:    req.IsCache,
		CreateBy:   createBy,
		CreateTime: &now,
		UpdateBy:   createBy,
		UpdateTime: &now,
		Remark:     req.Remark,
	}

	return global.DB.Create(&menu).Error
}

// UpdateMenu 更新菜单
func UpdateMenu(req dto.MenuUpdateReq, updateBy string) error {
	// 检查菜单是否存在
	var menu models.SysMenu
	err := global.DB.First(&menu, req.MenuID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("菜单不存在")
		}
		return err
	}

	// 不能将自己设为父菜单
	if req.ParentID == req.MenuID {
		return errors.New("不能将自己设为父菜单")
	}

	// 检查父菜单是否存在
	if req.ParentID != 0 {
		var parentMenu models.SysMenu
		err := global.DB.First(&parentMenu, req.ParentID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("父菜单不存在")
			}
			return err
		}

		// 检查是否会形成循环引用
		if strings.Contains(parentMenu.Ancestors, strconv.FormatUint(req.MenuID, 10)) {
			return errors.New("不能将子菜单设为父菜单")
		}
	}

	// 检查菜单名称是否重复（同级菜单名称不能重复，排除自己）
	var count int64
	err = global.DB.Model(&models.SysMenu{}).Where("parent_id = ? AND menu_name = ? AND menu_id != ?", req.ParentID, req.MenuName, req.MenuID).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("同级菜单名称不能重复")
	}

	// 构建祖级列表
	ancestors := "0"
	if req.ParentID != 0 {
		var parentMenu models.SysMenu
		global.DB.First(&parentMenu, req.ParentID)
		if parentMenu.Ancestors != "" {
			ancestors = parentMenu.Ancestors + "," + strconv.FormatUint(req.ParentID, 10)
		} else {
			ancestors = "0," + strconv.FormatUint(req.ParentID, 10)
		}
	}

	now := time.Now()
	updates := map[string]interface{}{
		"parent_id":   req.ParentID,
		"ancestors":   ancestors,
		"menu_name":   req.MenuName,
		"order_num":   req.OrderNum,
		"icon":        req.Icon,
		"path":        req.Path,
		"component":   req.Component,
		"query":       req.Query,
		"menu_type":   req.MenuType,
		"visible":     req.Visible,
		"status":      req.Status,
		"perms":       req.Perms,
		"is_frame":    req.IsFrame,
		"is_cache":    req.IsCache,
		"update_by":   updateBy,
		"update_time": &now,
		"remark":      req.Remark,
	}

	return global.DB.Model(&menu).Updates(updates).Error
}

// DeleteMenu 删除菜单
func DeleteMenu(menuID uint64) error {
	// 检查菜单是否存在
	var menu models.SysMenu
	err := global.DB.First(&menu, menuID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("菜单不存在")
		}
		return err
	}

	// 检查是否有子菜单
	var childCount int64
	err = global.DB.Model(&models.SysMenu{}).Where("parent_id = ?", menuID).Count(&childCount).Error
	if err != nil {
		return err
	}
	if childCount > 0 {
		return errors.New("存在子菜单，不允许删除")
	}

	// 检查是否有角色在使用此菜单
	var roleMenuCount int64
	err = global.DB.Model(&models.SysRoleMenu{}).Where("menu_id = ?", menuID).Count(&roleMenuCount).Error
	if err != nil {
		return err
	}
	if roleMenuCount > 0 {
		return errors.New("菜单已分配给角色，不允许删除")
	}

	return global.DB.Delete(&menu).Error
}

// GetMenuTree 获取菜单树（用于角色分配权限）
func GetMenuTree() ([]*dto.MenuTreeRes, error) {
	var menus []models.SysMenu
	err := global.DB.Where("status = ?", models.MenuStatusEnabled).Order("order_num ASC, menu_id ASC").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	// 转换为DTO并构建树形结构
	menuMap := make(map[uint64]*dto.MenuTreeRes)
	var rootMenus []*dto.MenuTreeRes

	// 先转换所有菜单为DTO
	for _, menu := range menus {
		menuRes := &dto.MenuTreeRes{
			MenuID:   menu.MenuID,
			ParentID: menu.ParentID,
			MenuName: menu.MenuName,
			MenuType: menu.MenuType,
			Children: []*dto.MenuTreeRes{},
		}
		menuMap[menu.MenuID] = menuRes
	}

	// 构建树形结构
	for _, menuRes := range menuMap {
		if menuRes.ParentID == 0 {
			// 顶级菜单
			rootMenus = append(rootMenus, menuRes)
		} else {
			// 子菜单，添加到父菜单的children中
			if parent, exists := menuMap[menuRes.ParentID]; exists {
				parent.Children = append(parent.Children, menuRes)
			}
		}
	}

	return rootMenus, nil
}

// GetMenuRouters 获取菜单路由（用于前端动态路由生成）
func GetMenuRouters(adminID uint64) ([]*dto.MenuRouterRes, error) {
	// 获取管理员的所有菜单权限
	var menus []models.SysMenu
	query := `
		SELECT DISTINCT m.* FROM sys_menu m
		LEFT JOIN sys_role_menu rm ON m.menu_id = rm.menu_id
		LEFT JOIN sys_admin_user_role aur ON rm.role_id = aur.role_id
		LEFT JOIN sys_role r ON aur.role_id = r.role_id
		WHERE aur.admin_id = ? AND m.status = 1 AND m.visible = 1 AND r.status = 1 AND r.del_flag = '0'
		AND m.menu_type IN ('M', 'C')
		ORDER BY m.order_num ASC, m.menu_id ASC
	`
	err := global.DB.Raw(query, adminID).Scan(&menus).Error
	if err != nil {
		return nil, err
	}

	// 转换为路由DTO并构建树形结构
	routerMap := make(map[uint64]*dto.MenuRouterRes)
	var rootRouters []*dto.MenuRouterRes

	// 先转换所有菜单为路由DTO
	for _, menu := range menus {
		router := &dto.MenuRouterRes{
			Name:      menu.MenuName,
			Path:      menu.Path,
			Hidden:    menu.Visible == models.MenuVisibleHidden,
			Component: getComponent(menu),
			Children:  []*dto.MenuRouterRes{},
		}

		// 设置路由参数
		if menu.Query != nil && *menu.Query != "" {
			router.Query = *menu.Query
		}

		// 设置路由元信息
		router.Meta = &dto.MenuRouterMeta{
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

	return rootRouters, nil
}

// getComponent 获取组件路径
func getComponent(menu models.SysMenu) string {
	if menu.Component != nil && *menu.Component != "" {
		return *menu.Component
	}

	// 如果是目录且没有指定组件，返回Layout
	if menu.MenuType == models.MenuTypeDir {
		return "Layout"
	}

	return ""
}
