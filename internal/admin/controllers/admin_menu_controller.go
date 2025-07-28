package controllers

import (
	"github.com/gin-gonic/gin"
	"go_server/common"
	"go_server/constants"
	"go_server/internal/admin/dto"
	"go_server/internal/admin/services"
	"strconv"
)

// GetMenuList 获取菜单列表
// @Summary 获取菜单列表
// @Description 获取系统菜单列表，支持树形结构展示和条件查询
// @Tags [Admin]菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param menu_name query string false "菜单名称"
// @Param status query int false "菜单状态（0停用 1正常）"
// @Param menu_type query string false "菜单类型（M目录 C菜单 F按钮）"
// @Param parent_id query int false "父菜单ID"
// @Param visible query int false "菜单状态（0隐藏 1显示）"
// @Success 200 {object} common.Response{data=[]dto.MenuListRes} "获取成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Router /v1/admin/menus [get]
func GetMenuList(c *gin.Context) {
	var req dto.MenuListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Error(c, constants.ErrorCode, "参数错误: "+err.Error())
		return
	}

	// 调用服务获取菜单列表
	menuList, err := services.GetMenuList(req)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, menuList)
}

// GetMenuDetail 获取菜单详情
// @Summary 获取菜单详情
// @Description 根据菜单ID获取菜单的详细信息
// @Tags [Admin]菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "菜单ID"
// @Success 200 {object} common.Response{data=dto.MenuDetailRes} "获取成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 404 {object} common.Response "菜单不存在"
// @Router /v1/admin/menus/{id} [get]
func GetMenuDetail(c *gin.Context) {
	idStr := c.Param("id")
	menuID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, constants.ErrorCode, "菜单ID格式错误")
		return
	}

	// 调用服务获取菜单详情
	menuDetail, err := services.GetMenuDetail(menuID)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, menuDetail)
}

// CreateMenu 创建菜单
// @Summary 创建菜单
// @Description 创建新的系统菜单
// @Tags [Admin]菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.MenuCreateReq true "创建菜单请求参数"
// @Success 200 {object} common.Response "创建成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Router /v1/admin/menus [post]
func CreateMenu(c *gin.Context) {
	var req dto.MenuCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, constants.ErrorCode, "参数错误: "+err.Error())
		return
	}

	// 获取当前管理员信息
	adminName, exists := c.Get("admin_name")
	if !exists {
		common.Error(c, constants.ErrorCode, "获取管理员信息失败")
		return
	}

	// 调用服务创建菜单
	err := services.CreateMenu(req, adminName.(string))
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, "菜单创建成功")
}

// UpdateMenu 更新菜单
// @Summary 更新菜单
// @Description 更新系统菜单信息
// @Tags [Admin]菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "菜单ID"
// @Param request body dto.MenuUpdateReq true "更新菜单请求参数"
// @Success 200 {object} common.Response "更新成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 404 {object} common.Response "菜单不存在"
// @Router /v1/admin/menus/{id} [put]
func UpdateMenu(c *gin.Context) {
	idStr := c.Param("id")
	menuID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, constants.ErrorCode, "菜单ID格式错误")
		return
	}

	var req dto.MenuUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, constants.ErrorCode, "参数错误: "+err.Error())
		return
	}

	// 确保路径参数和请求体中的ID一致
	req.MenuID = menuID

	// 获取当前管理员信息
	adminName, exists := c.Get("admin_name")
	if !exists {
		common.Error(c, constants.ErrorCode, "获取管理员信息失败")
		return
	}

	// 调用服务更新菜单
	err = services.UpdateMenu(req, adminName.(string))
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, "菜单更新成功")
}

// DeleteMenu 删除菜单
// @Summary 删除菜单
// @Description 删除系统菜单
// @Tags [Admin]菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "菜单ID"
// @Success 200 {object} common.Response "删除成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 404 {object} common.Response "菜单不存在"
// @Router /v1/admin/menus/{id} [delete]
func DeleteMenu(c *gin.Context) {
	idStr := c.Param("id")
	menuID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, constants.ErrorCode, "菜单ID格式错误")
		return
	}

	// 调用服务删除菜单
	err = services.DeleteMenu(menuID)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, "菜单删除成功")
}

// GetMenuTree 获取菜单树
// @Summary 获取菜单树
// @Description 获取菜单树形结构，用于角色分配权限时的选择
// @Tags [Admin]菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.Response{data=[]dto.MenuTreeRes} "获取成功"
// @Failure 401 {object} common.Response "未授权"
// @Router /v1/admin/menus/tree [get]
func GetMenuTree(c *gin.Context) {
	// 调用服务获取菜单树
	menuTree, err := services.GetMenuTree()
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, menuTree)
}

// GetMenuRouters 获取菜单路由
// @Summary 获取菜单路由
// @Description 获取当前管理员的菜单路由信息，用于前端动态路由生成
// @Tags [Admin]菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.Response{data=[]dto.MenuRouterRes} "获取成功"
// @Failure 401 {object} common.Response "未授权"
// @Router /v1/admin/menus/routers [get]
func GetMenuRouters(c *gin.Context) {
	// 获取当前管理员ID
	adminID, exists := c.Get("admin_id")
	if !exists {
		common.Error(c, constants.ErrorCode, "获取管理员信息失败")
		return
	}

	// 调用服务获取菜单路由
	menuRouters, err := services.GetMenuRouters(adminID.(uint64))
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, menuRouters)
}
