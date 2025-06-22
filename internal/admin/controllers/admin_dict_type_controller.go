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
)

// CreateDictType 创建字典类型
// 接收前端POST请求，创建新的字典类型
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
// 接收前端PUT请求，更新指定的字典类型
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
// 接收前端DELETE请求，删除指定的字典类型
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
// 接收前端GET请求，返回指定字典类型的详细信息
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
// 接收前端GET请求，返回字典类型列表（支持分页和筛选）
func GetDictTypeList(c *gin.Context) {
	var query dto.DictTypeQuery
	fmt.Println("获取字典类型列表")
	// ShouldBindQuery 自动解析URL查询参数
	if err := c.ShouldBindQuery(&query); err != nil {
		common.Error(c, constants.ErrorCode, "查询参数错误")
		return
	}

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
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	// 使用分页响应格式返回数据
	// 包含当前页、每页数量、总页数、总条数、数据列表
	common.SuccessWithPage(c, int64(query.PageNum), query.PageSize, total, dictTypes)
}

// GetDictDataByType 根据字典类型获取字典数据
// 这是一个简单的接口，用于前端获取某个类型下的所有字典项
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
