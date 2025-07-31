package controllers

import (
	"github.com/gin-gonic/gin"
	"go_server/common"
	"go_server/constants"
	"go_server/internal/admin/dto"
	"go_server/internal/admin/services"
	"strconv"
)

// GetRoleList 获取角色列表
// @Summary 获取角色列表
// @Description 获取系统角色列表，支持分页和条件查询
// @Tags Admin-角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param role_name query string false "角色名称"
// @Param role_key query string false "角色权限字符串"
// @Param status query int false "角色状态（0停用 1正常）"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} common.Response{data=common.PageResult} "获取成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Router /v1/admin/roles [get]
func GetRoleList(c *gin.Context) {
	var req dto.RoleListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Error(c, constants.ErrorCode, "参数错误: "+err.Error())
		return
	}

	// 调用服务获取角色列表
	roleList, err := services.GetRoleList(req)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, roleList)
}

// GetRoleDetail 获取角色详情
// @Summary 获取角色详情
// @Description 根据角色ID获取角色的详细信息，包括菜单权限
// @Tags Admin-角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Success 200 {object} common.Response{data=dto.RoleDetailRes} "获取成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 404 {object} common.Response "角色不存在"
// @Router /v1/admin/roles/{id} [get]
func GetRoleDetail(c *gin.Context) {
	idStr := c.Param("id")
	roleID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, constants.ErrorCode, "角色ID格式错误")
		return
	}

	// 调用服务获取角色详情
	roleDetail, err := services.GetRoleDetail(roleID)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, roleDetail)
}

// CreateRole 创建角色
// @Summary 创建角色
// @Description 创建新的系统角色
// @Tags Admin-角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.RoleCreateReq true "创建角色请求参数"
// @Success 200 {object} common.Response "创建成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Router /v1/admin/roles [post]
func CreateRole(c *gin.Context) {
	var req dto.RoleCreateReq
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

	// 调用服务创建角色
	err := services.CreateRole(req, adminName.(string))
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, "角色创建成功")
}

// UpdateRole 更新角色
// @Summary 更新角色
// @Description 更新系统角色信息
// @Tags Admin-角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Param request body dto.RoleUpdateReq true "更新角色请求参数"
// @Success 200 {object} common.Response "更新成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 404 {object} common.Response "角色不存在"
// @Router /v1/admin/roles/{id} [put]
func UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	roleID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, constants.ErrorCode, "角色ID格式错误")
		return
	}

	var req dto.RoleUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, constants.ErrorCode, "参数错误: "+err.Error())
		return
	}

	// 确保路径参数和请求体中的ID一致
	req.RoleID = roleID

	// 获取当前管理员信息
	adminName, exists := c.Get("admin_name")
	if !exists {
		common.Error(c, constants.ErrorCode, "获取管理员信息失败")
		return
	}

	// 调用服务更新角色
	err = services.UpdateRole(req, adminName.(string))
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, "角色更新成功")
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Description 删除系统角色
// @Tags Admin-角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Success 200 {object} common.Response "删除成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 404 {object} common.Response "角色不存在"
// @Router /v1/admin/roles/{id} [delete]
func DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	roleID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, constants.ErrorCode, "角色ID格式错误")
		return
	}

	// 调用服务删除角色
	err = services.DeleteRole(roleID)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, "角色删除成功")
}

// GetRoleSelect 获取角色选择列表
// @Summary 获取角色选择列表
// @Description 获取可用的角色列表，用于下拉选择
// @Tags Admin-角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.Response{data=[]dto.RoleSelectRes} "获取成功"
// @Failure 401 {object} common.Response "未授权"
// @Router /v1/admin/roles/select [get]
func GetRoleSelect(c *gin.Context) {
	// 调用服务获取角色选择列表
	roleList, err := services.GetRoleSelect()
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, roleList)
}

// AuthRole 角色授权
// @Summary 角色授权
// @Description 为角色分配菜单权限
// @Tags Admin-角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.RoleAuthReq true "角色授权请求参数"
// @Success 200 {object} common.Response "授权成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 404 {object} common.Response "角色不存在"
// @Router /v1/admin/roles/auth [post]
func AuthRole(c *gin.Context) {
	var req dto.RoleAuthReq
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

	// 调用服务进行角色授权
	err := services.AuthRole(req, adminName.(string))
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, "角色授权成功")
}

// AssignAdminRole 分配管理员角色
// @Summary 分配管理员角色
// @Description 为管理员分配角色
// @Tags Admin-角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.AdminRoleReq true "管理员角色分配请求参数"
// @Success 200 {object} common.Response "分配成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 404 {object} common.Response "管理员不存在"
// @Router /v1/admin/roles/assign [post]
func AssignAdminRole(c *gin.Context) {
	var req dto.AdminRoleReq
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

	// 调用服务分配管理员角色
	err := services.AssignAdminRole(req, adminName.(string))
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, "管理员角色分配成功")
}

// GetAdminRole 获取管理员角色信息
// @Summary 获取管理员角色信息
// @Description 获取指定管理员的角色信息
// @Tags Admin-角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param admin_id path int true "管理员ID"
// @Success 200 {object} common.Response{data=dto.AdminRoleRes} "获取成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 404 {object} common.Response "管理员不存在"
// @Router /v1/admin/roles/admin/{admin_id} [get]
func GetAdminRole(c *gin.Context) {
	idStr := c.Param("admin_id")
	adminID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, constants.ErrorCode, "管理员ID格式错误")
		return
	}

	// 调用服务获取管理员角色信息
	adminRole, err := services.GetAdminRole(adminID)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, adminRole)
}

// GetRoleMenuTree 获取角色菜单权限树
// @Summary 获取角色菜单权限树
// @Description 获取指定角色的菜单权限树，显示已分配的菜单权限状态
// @Tags Admin-角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param role_id path int true "角色ID"
// @Success 200 {object} common.Response{data=[]dto.RoleMenuTreeRes} "获取成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 404 {object} common.Response "角色不存在"
// @Router /v1/admin/roles/menus/{role_id} [get]
func GetRoleMenuTree(c *gin.Context) {
	idStr := c.Param("role_id")
	roleID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, constants.ErrorCode, "角色ID格式错误")
		return
	}

	// 调用服务获取角色菜单权限树
	menuTree, err := services.GetRoleMenuTree(roleID)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, menuTree)
}

// AssignRoleMenus 分配角色菜单权限
// @Summary 分配角色菜单权限
// @Description 为指定角色分配菜单权限
// @Tags Admin-角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.RoleMenuAuthReq true "角色菜单权限分配请求参数"
// @Success 200 {object} common.Response "分配成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 404 {object} common.Response "角色不存在"
// @Router /v1/admin/roles/menus/assign [post]
func AssignRoleMenus(c *gin.Context) {
	var req dto.RoleMenuAuthReq
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

	// 调用服务分配角色菜单权限
	err := services.AssignRoleMenus(req, adminName.(string))
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, "角色菜单权限分配成功")
}
