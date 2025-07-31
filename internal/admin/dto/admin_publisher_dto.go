package dto

import (
	"go_server/common"
	"go_server/utils"
	"time"
)

// PublisherReq 游戏厂商请求结构体
// 用于创建和更新游戏厂商时接收前端传来的数据
type PublisherReq struct {
	PublisherID   uint64     `json:"publisher_id,omitempty" example:"1"`              // 厂商ID，更新时需要，创建时忽略
	PublisherName string     `json:"publisher_name" binding:"required" example:"米哈游"` // 厂商名称，必填
	LogoURL       *string    `json:"logo_url" example:"https://example.com/logo.png"` // 厂商LOGO，可选
	Description   *string    `json:"description" example:"知名游戏开发商"`                   // 厂商介绍，可选
	FoundedDate   *time.Time `json:"founded_date" example:"2012-02-13T00:00:00Z"`     // 成立日期，可选
	Website       *string    `json:"website" example:"https://www.mihoyo.com"`        // 官方网站，可选
	Status        int        `json:"status" binding:"required,oneof=0 1" example:"1"` // 状态，必填，只能是0或1
}

// PublisherRes 游戏厂商响应结构体
// 用于返回游戏厂商数据
type PublisherRes struct {
	PublisherID   uint64     `json:"publisher_id"`   // 厂商ID
	PublisherName string     `json:"publisher_name"` // 厂商名称
	LogoURL       *string    `json:"logo_url"`       // 厂商LOGO
	Description   *string    `json:"description"`    // 厂商介绍
	FoundedDate   *time.Time `json:"founded_date"`   // 成立日期
	Website       *string    `json:"website"`        // 官方网站
	Status        int        `json:"status"`         // 状态
	common.BaseEntity
}

// PublisherQuery 游戏厂商查询参数
// 用于查询游戏厂商列表时的筛选条件
type PublisherQuery struct {
	PublisherName       string `form:"publisher_name"` // 厂商名称，支持模糊查询
	PublisherID         string `form:"publisher_id"`   // 厂商ID 精确查询
	common.StatusEntity        // 嵌入状态字段
	common.PageQuery           // 嵌入分页字段

}

// GetMessages 自定义验证错误信息
// 当参数验证失败时，返回友好的中文错误信息
func (req PublisherReq) GetMessages() utils.ValidatorMessages {
	return utils.ValidatorMessages{
		"PublisherName.required": "厂商名称不能为空",
		"Status.required":        "状态不能为空",
		"Status.oneof":           "状态只能是0或1",
	}
}
