package common

import (
	"github.com/gin-gonic/gin"
	"go_server/constants"
	"go_server/system/dto"
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
