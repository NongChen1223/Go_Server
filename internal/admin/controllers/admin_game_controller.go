package controllers

import (
	"github.com/gin-gonic/gin"
	"go_server/common"
	"go_server/constants"
	"go_server/internal/admin/dto"
	"go_server/internal/admin/services"
	"strconv"
)

// CreateGame 创建游戏
// @Summary 创建游戏
// @Description 创建新的游戏记录，包括基本信息和关联数据
// @Tags Admin-游戏管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.GameReq true "游戏创建请求参数"
// @Success 200 {object} common.Response "创建成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Router /v1/admin/games [post]
func CreateGame(c *gin.Context) {
	var req dto.GameReq
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

	// 调用创建服务
	err := services.CreateGame(req, adminName.(string))
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, nil)
}

// UpdateGame 更新游戏
// @Summary 更新游戏
// @Description 更新游戏信息，包括基本信息和关联数据
// @Tags Admin-游戏管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "游戏ID"
// @Param request body dto.GameReq true "游戏更新请求参数"
// @Success 200 {object} common.Response "更新成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Router /v1/admin/games/{id} [put]
func UpdateGame(c *gin.Context) {
	// 获取路径参数中的游戏ID
	gameIDStr := c.Param("id")
	gameID, err := strconv.ParseUint(gameIDStr, 10, 64)
	if err != nil {
		common.Error(c, constants.ErrorCode, "游戏ID格式错误")
		return
	}

	var req dto.GameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, constants.ErrorCode, "参数错误: "+err.Error())
		return
	}

	// 设置游戏ID
	req.GameID = gameID

	// 获取当前管理员信息
	adminName, exists := c.Get("admin_name")
	if !exists {
		common.Error(c, constants.ErrorCode, "获取管理员信息失败")
		return
	}

	// 调用更新服务
	err = services.UpdateGame(req, adminName.(string))
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, nil)
}

// GetGameDetail 获取游戏详情
// @Summary 获取游戏详情
// @Description 根据游戏ID获取游戏的详细信息，包括所有关联数据
// @Tags Admin-游戏管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "游戏ID"
// @Success 200 {object} common.Response{data=dto.GameRes} "获取成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Router /v1/admin/games/{id} [get]
func GetGameDetail(c *gin.Context) {
	// 获取路径参数中的游戏ID
	gameIDStr := c.Param("id")
	gameID, err := strconv.ParseUint(gameIDStr, 10, 64)
	if err != nil {
		common.Error(c, constants.ErrorCode, "游戏ID格式错误")
		return
	}

	// 调用获取详情服务
	game, err := services.GetGameDetail(gameID)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, game)
}

// GetGameList 获取游戏列表
// @Summary 获取游戏列表
// @Description 分页获取游戏列表，支持多种条件筛选
// @Tags Admin-游戏管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name_zh query string false "游戏中文名称（模糊查询）"
// @Param name_en query string false "游戏英文名称（模糊查询）"
// @Param publisher_id query int false "厂商ID（精确查询）"
// @Param studio_id query int false "工作室ID（精确查询）"
// @Param type_code query int false "游戏类型编码（精确查询）"
// @Param platform_code query int false "游戏平台编码（精确查询）"
// @Param language_code query int false "游戏语言编码（精确查询）"
// @Param min_rating query number false "最低评分"
// @Param max_rating query number false "最高评分"
// @Param min_price query number false "最低价格"
// @Param max_price query number false "最高价格"
// @Param status query int false "状态（0停用 1正常）"
// @Param page_num query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} common.PageResponse{records=[]dto.GameRes} "获取成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Router /v1/admin/games [get]
func GetGameList(c *gin.Context) {
	var query dto.GameQuery
	// ShouldBindQuery 自动解析URL查询参数
	if err := c.ShouldBindQuery(&query); err != nil {
		common.Error(c, constants.ErrorCode, "参数错误: "+err.Error())
		return
	}

	// 设置默认分页参数
	query.GetDefaultPage()

	// 调用获取列表服务
	total, games, err := services.GetGameList(query)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	// 返回分页数据
	common.SuccessWithPage(c, int64(query.PageNum), query.PageSize, total, games)
}

// DeleteGame 删除游戏
// @Summary 删除游戏
// @Description 根据游戏ID删除游戏及其所有关联数据
// @Tags Admin-游戏管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "游戏ID"
// @Success 200 {object} common.Response "删除成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Router /v1/admin/games/{id} [delete]
func DeleteGame(c *gin.Context) {
	// 获取路径参数中的游戏ID
	gameIDStr := c.Param("id")
	gameID, err := strconv.ParseUint(gameIDStr, 10, 64)
	if err != nil {
		common.Error(c, constants.ErrorCode, "游戏ID格式错误")
		return
	}

	// 调用删除服务
	err = services.DeleteGame(gameID)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, nil)
}
