package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_server/common"
	"go_server/constants"
	"go_server/internal/admin/dto"
	"go_server/internal/admin/services"
	"go_server/utils"
	"strconv"
	"sync"
	"time"
)

// 防重复请求的缓存
var (
	requestCache = make(map[string]time.Time)
	cacheMutex   sync.RWMutex
)

// CreateDictType 创建字典类型
// @Summary 创建字典类型
// @Description 创建新的字典类型，用于系统配置管理
// @Tags Admin-字典管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.DictTypeReq true "字典类型信息"
// @Success 200 {object} common.Response "创建成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 500 {object} common.Response "服务器错误"
// @Router /v1/admin/dict/types [post]
func CreateDictType(c *gin.Context) {
	var req dto.DictTypeReq
	// ShouldBindJSON 自动解析JSON请求体并验证参数
	if err := c.ShouldBindJSON(&req); err != nil {
		// 使用自定义错误信息，提供友好的中文提示
		common.Error(c, constants.ErrorCode, utils.GetErrorMsg(req, err))
		return
	}

	// 从上下文中获取当前管理员信息
	// 这个信息是在认证中间件中设置的
	adminName := c.GetString("admin_name")
	if adminName == "" {
		adminName = "admin" // 默认管理员（兜底处理）
	}

	// 调用服务层创建字典类型
	err := services.CreateDictType(req, adminName)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	// 创建成功，返回成功响应
	common.Success(c, nil)
}

// UpdateDictType 更新字典类型
// @Summary 更新字典类型
// @Description 根据ID更新字典类型信息
// @Tags Admin-字典管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "字典类型ID"
// @Param request body dto.DictTypeReq true "字典类型信息"
// @Success 200 {object} common.Response "更新成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 404 {object} common.Response "字典类型不存在"
// @Router /v1/admin/dict/types/{id} [put]
func UpdateDictType(c *gin.Context) {
	var req dto.DictTypeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, constants.ErrorCode, utils.GetErrorMsg(req, err))
		return
	}

	// 从URL路径参数中获取字典类型ID
	dictIDStr := c.Param("id")
	dictID, err := strconv.ParseUint(dictIDStr, 10, 64)
	if err != nil {
		common.Error(c, constants.ErrorCode, "无效的字典类型ID")
		return
	}

	// 将ID设置到请求结构体中
	req.DictID = dictID

	// 获取当前管理员
	adminName := c.GetString("admin_name")
	if adminName == "" {
		adminName = "admin"
	}

	// 调用服务层更新字典类型
	err = services.UpdateDictType(req, adminName)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, nil)
}

// DeleteDictType 删除字典类型
// @Summary 删除字典类型
// @Description 根据ID删除字典类型
// @Tags Admin-字典管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "字典类型ID"
// @Success 200 {object} common.Response "删除成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 404 {object} common.Response "字典类型不存在"
// @Router /v1/admin/dict/types/{id} [delete]
func DeleteDictType(c *gin.Context) {
	// 从URL路径参数中获取要删除的字典类型ID
	dictIDStr := c.Param("id")
	dictID, err := strconv.ParseUint(dictIDStr, 10, 64)
	if err != nil {
		common.Error(c, constants.ErrorCode, "无效的字典类型ID")
		return
	}

	// 调用服务层删除字典类型
	err = services.DeleteDictType(dictID)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, nil)
}

// GetDictTypeDetail 获取字典类型详情
// @Summary 获取字典类型详情
// @Description 根据ID获取字典类型的详细信息
// @Tags Admin-字典管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "字典类型ID"
// @Success 200 {object} common.Response{data=dto.DictTypeRes} "获取成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 404 {object} common.Response "字典类型不存在"
// @Router /v1/admin/dict/types/{id} [get]
func GetDictTypeDetail(c *gin.Context) {
	dictIDStr := c.Param("id")
	dictID, err := strconv.ParseUint(dictIDStr, 10, 64)
	fmt.Println("字典类型ID", dictID)
	if err != nil {
		common.Error(c, constants.ErrorCode, "无效的字典类型ID")
		return
	}

	// 调用服务层获取字典类型详情
	dictType, err := services.GetDictTypeDetail(dictID)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	// 返回字典类型详情数据
	common.Success(c, dictType)
}

// GetDictTypeList 获取字典类型列表
// @Summary 获取字典类型列表
// @Description 分页获取字典类型列表，支持按名称和类型筛选
// @Tags Admin-字典管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param dict_name query string false "字典名称（模糊查询）"
// @Param dict_type query string false "字典类型（模糊查询）"
// @Param status query int false "状态（0停用 1正常）"
// @Param page_num query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} common.PageResponse{records=[]dto.DictTypeRes} "获取成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Router /v1/admin/dict/types [get]
func GetDictTypeList(c *gin.Context) {
	var query dto.DictTypeQuery

	// 防重复请求检查
	clientIP := c.ClientIP()
	requestKey := fmt.Sprintf("dict_type_list_%s", clientIP)

	cacheMutex.RLock()
	lastRequestTime, exists := requestCache[requestKey]
	cacheMutex.RUnlock()

	now := time.Now()
	if exists && now.Sub(lastRequestTime) < 2*time.Second {
		fmt.Printf("检测到重复请求，忽略。客户端IP: %s, 上次请求时间: %s\n", clientIP, lastRequestTime.Format("15:04:05"))
		common.Error(c, constants.ErrorCode, "请求过于频繁，请稍后再试")
		return
	}

	// 更新请求时间
	cacheMutex.Lock()
	requestCache[requestKey] = now
	cacheMutex.Unlock()

	fmt.Printf("=== 字典类型列表请求 ===\n")
	fmt.Printf("客户端IP: %s\n", clientIP)
	fmt.Printf("请求时间: %s\n", now.Format("2006-01-02 15:04:05"))

	// ShouldBindQuery 自动解析URL查询参数
	if err := c.ShouldBindQuery(&query); err != nil {
		fmt.Printf("参数绑定错误: %v\n", err)
		common.Error(c, constants.ErrorCode, "查询参数错误")
		return
	}
	fmt.Printf("查询参数: %+v\n", query)

	// 设置默认分页参数
	// 如果前端没有传分页参数，使用默认值
	if query.PageNum <= 0 {
		query.PageNum = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	// 调用服务层获取字典类型列表
	total, dictTypes, err := services.GetDictTypeList(query)
	if err != nil {
		fmt.Printf("服务层错误: %v\n", err)
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	fmt.Printf("查询结果: total=%d, count=%d\n", total, len(dictTypes))

	// 添加调试响应头
	c.Header("X-Request-Time", time.Now().Format("2006-01-02 15:04:05"))
	c.Header("X-Total-Records", fmt.Sprintf("%d", total))
	c.Header("X-Debug-Info", "dict-type-list-success")

	// 使用分页响应格式返回数据
	// 包含当前页、每页数量、总页数、总条数、数据列表
	fmt.Printf("返回成功响应，总记录数: %d\n", total)
	common.SuccessWithPage(c, int64(query.PageNum), query.PageSize, total, dictTypes)
}

// GetDictDataByType 根据字典类型获取字典数据
// @Summary 根据类型获取字典数据
// @Description 根据字典类型标识获取该类型下的所有字典数据
// @Tags Admin-字典管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type path string true "字典类型标识"
// @Success 200 {object} common.Response{data=[]interface{}} "获取成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Router /v1/admin/dict/data/{type} [get]
func GetDictDataByType(c *gin.Context) {
	dictType := c.Param("type")
	if dictType == "" {
		common.Error(c, constants.ErrorCode, "字典类型不能为空")
		return
	}

	// TODO: 这里可以实现获取字典数据的逻辑
	// 目前返回空数组作为占位
	common.Success(c, []interface{}{})
}
