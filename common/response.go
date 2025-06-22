package common

import (
	"github.com/gin-gonic/gin"
	"go_server/constants"
	"go_server/internal/system/dto"
	"net/http"
)

type Response struct {
	// 代码
	Code int `json:"code" example:"200"`
	// 数据集
	Data interface{} `json:"data,omitempty"`
	// 消息
	Msg string `json:"msg"`
}

// PageResponse 分页响应结构体
type PageResponse struct {
	Current int64       `json:"current"` // 当前页码
	Size    int         `json:"size"`    // 每页条数
	Pages   int64       `json:"pages"`   // 总页数
	Total   int64       `json:"total"`   // 总条数
	Records interface{} `json:"records"` // 数据列表
}

// 通用响应函数
func Respond(c *gin.Context, httpCode, code int, msg string, data interface{}) {
	c.JSON(httpCode, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

func LoginSuccess(c *gin.Context, res dto.LoginRes) {
	c.JSON(http.StatusOK, res)
}

// 成功响应（带数据）
func Success(c *gin.Context, data interface{}) {
	Respond(c, http.StatusOK, constants.SuccessCode, "操作成功", data)
}

// 成功响应（带自定义消息）
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	Respond(c, http.StatusOK, constants.SuccessCode, msg, data)
}

// 错误响应（无数据）
func Error(c *gin.Context, code int, msg string) {
	Respond(c, http.StatusOK, code, msg, nil)
}

// 错误响应（HTTP 状态码和业务码分离）
func ErrorWithHTTPCode(c *gin.Context, httpCode, code int, msg string) {
	Respond(c, httpCode, code, msg, nil)
}

// SuccessWithPage 分页成功响应
// 参数：
//   - current: 当前页码
//   - size: 每页条数
//   - total: 总条数
//   - records: 数据列表
func SuccessWithPage(c *gin.Context, current int64, size int, total int64, records interface{}) {
	// 计算总页数
	pages := (total + int64(size) - 1) / int64(size) // 向上取整
	if total == 0 {
		pages = 0
	}

	pageData := PageResponse{
		Current: current,
		Size:    size,
		Pages:   pages,
		Total:   total,
		Records: records,
	}

	Respond(c, http.StatusOK, constants.SuccessCode, "操作成功", pageData)
}

// SuccessWithPageAndMsg 分页成功响应（带自定义消息）
func SuccessWithPageAndMsg(c *gin.Context, msg string, current int64, size int, total int64, records interface{}) {
	// 计算总页数
	pages := (total + int64(size) - 1) / int64(size) // 向上取整
	if total == 0 {
		pages = 0
	}

	pageData := PageResponse{
		Current: current,
		Size:    size,
		Pages:   pages,
		Total:   total,
		Records: records,
	}

	Respond(c, http.StatusOK, constants.SuccessCode, msg, pageData)
}
