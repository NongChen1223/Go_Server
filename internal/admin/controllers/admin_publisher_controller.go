package controllers

import (
	"github.com/gin-gonic/gin"
	"go_server/common"
	"go_server/constants"
	"go_server/internal/admin/dto"
	"go_server/internal/admin/services"
	"go_server/utils"
	"strconv"
)

// CreatePublisher 创建游戏厂商
// 接收前端POST请求，创建新的游戏厂商
func CreatePublisher(c *gin.Context) {
	var req dto.PublisherReq
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

	// 调用服务层创建游戏厂商
	err := services.CreatePublisher(req, adminName)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	// 创建成功，返回成功响应
	common.Success(c, nil)
}

// UpdatePublisher 更新游戏厂商
// 接收前端PUT请求，更新指定的游戏厂商
func UpdatePublisher(c *gin.Context) {
	var req dto.PublisherReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, constants.ErrorCode, utils.GetErrorMsg(req, err))
		return
	}

	// 从URL路径参数中获取游戏厂商ID
	publisherIDStr := c.Param("id")
	publisherID, err := strconv.ParseUint(publisherIDStr, 10, 64)
	if err != nil {
		common.Error(c, constants.ErrorCode, "无效的厂商ID")
		return
	}

	// 将ID设置到请求结构体中
	req.PublisherID = publisherID

	// 获取当前管理员
	adminName := c.GetString("admin_name")
	if adminName == "" {
		adminName = "admin"
	}

	// 调用服务层更新游戏厂商
	err = services.UpdatePublisher(req, adminName)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, nil)
}

// DeletePublisher 删除游戏厂商
// 接收前端DELETE请求，删除指定的游戏厂商
func DeletePublisher(c *gin.Context) {
	// 从URL路径参数中获取要删除的游戏厂商ID
	publisherIDStr := c.Param("id")
	publisherID, err := strconv.ParseUint(publisherIDStr, 10, 64)
	if err != nil {
		common.Error(c, constants.ErrorCode, "无效的厂商ID")
		return
	}

	// 调用服务层删除游戏厂商
	err = services.DeletePublisher(publisherID)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	common.Success(c, nil)
}

// GetPublisherDetail 获取游戏厂商详情
// 接收前端GET请求，返回指定游戏厂商的详细信息
func GetPublisherDetail(c *gin.Context) {
	publisherIDStr := c.Param("id")
	publisherID, err := strconv.ParseUint(publisherIDStr, 10, 64)
	if err != nil {
		common.Error(c, constants.ErrorCode, "无效的厂商ID")
		return
	}

	// 调用服务层获取游戏厂商详情
	publisher, err := services.GetPublisherDetail(publisherID)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	// 返回游戏厂商详情数据
	common.Success(c, publisher)
}

// GetPublisherList 获取游戏厂商列表
// 接收前端GET请求，返回游戏厂商列表（支持分页和筛选）
func GetPublisherList(c *gin.Context) {
	var query dto.PublisherQuery
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

	// 调用服务层获取游戏厂商列表
	total, publishers, err := services.GetPublisherList(query)
	if err != nil {
		common.Error(c, constants.ErrorCode, err.Error())
		return
	}

	// 使用分页响应格式返回数据
	// 包含当前页、每页数量、总页数、总条数、数据列表
	common.SuccessWithPage(c, int64(query.PageNum), query.PageSize, total, publishers)
}
