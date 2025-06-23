package dto

import (
	"go_server/common"
	"go_server/utils"
)

// DictTypeReq 字典类型请求结构体
// 用于创建和更新字典类型时接收前端传来的数据
type DictTypeReq struct {
	DictID              uint64 `json:"dict_id,omitempty"`            // 字典类型ID，更新时需要，创建时忽略
	DictName            string `json:"dict_name" binding:"required"` // 字典名称，必填
	DictType            string `json:"dict_type" binding:"required"` // 字典类型标识，必填，如"game_type"
	common.StatusEntity        // 嵌入状态字段
	common.RemarkEntity        // 嵌入备注字段
}

// DictTypeRes 字典类型响应结构体
// 用于返回给前端的字典类型数据
type DictTypeRes struct {
	DictID              uint64 `json:"dict_id"`   // 字典类型ID
	DictName            string `json:"dict_name"` // 字典名称
	DictType            string `json:"dict_type"` // 字典类型标识
	Status              int    `json:"status"`    // 状态（0停用 1正常）
	common.RemarkEntity        // 嵌入备注字段
	common.BaseEntity          // 嵌入基础审计字段
}

// DictTypeQuery 字典类型查询参数
// 用于列表查询时的筛选和分页参数
type DictTypeQuery struct {
	DictName            string `form:"dict_name"` // 字典名称，支持模糊查询
	DictType            string `form:"dict_type"` // 字典类型，支持模糊查询
	common.StatusEntity        // 嵌入状态字段
	common.PageQuery
}

// DictDataReq 字典数据请求结构体
// 用于创建和更新字典数据时接收前端传来的数据
type DictDataReq struct {
	DictCode            uint64 `json:"dict_code,omitempty"`            // 字典编码，更新时需要
	DictSort            int    `json:"dict_sort"`                      // 字典排序，数字越小越靠前
	DictLabel           string `json:"dict_label" binding:"required"`  // 字典标签，显示给用户看的文本
	DictValue           string `json:"dict_value" binding:"required"`  // 字典键值，程序中使用的值
	DictType            string `json:"dict_type" binding:"required"`   // 字典类型，关联字典类型表
	IsDefault           int    `json:"is_default" binding:"oneof=0 1"` // 是否默认选项
	common.StatusEntity        // 嵌入状态字段
	common.RemarkEntity        // 嵌入备注字段
}

// DictDataRes 字典数据响应结构体
// 用于返回给前端的字典数据
type DictDataRes struct {
	DictCode            uint64 `json:"dict_code"`  // 字典编码
	DictSort            int    `json:"dict_sort"`  // 字典排序
	DictLabel           string `json:"dict_label"` // 字典标签
	DictValue           string `json:"dict_value"` // 字典键值
	DictType            string `json:"dict_type"`  // 字典类型
	IsDefault           int    `json:"is_default"` // 是否默认
	common.StatusEntity        // 嵌入状态字段
	common.RemarkEntity        // 嵌入备注字段
	common.BaseEntity          // 嵌入基础审计字段
}

// DictDataQuery 字典数据查询参数
type DictDataQuery struct {
	DictType            string `form:"dict_type"`  // 字典类型筛选
	DictLabel           string `form:"dict_label"` // 字典标签，支持模糊查询
	common.StatusEntity        // 嵌入状态字段
	common.PageQuery
}

// GetMessages 自定义验证错误信息
// 当参数验证失败时，返回友好的中文错误信息
func (req DictTypeReq) GetMessages() utils.ValidatorMessages {
	return utils.ValidatorMessages{
		"DictName.required": "字典名称不能为空",
		"DictType.required": "字典类型不能为空",
		"Status.required":   "状态不能为空",
		"Status.oneof":      "状态只能是0或1",
	}
}

// GetMessages 字典数据的验证错误信息
func (req DictDataReq) GetMessages() utils.ValidatorMessages {
	return utils.ValidatorMessages{
		"DictLabel.required": "字典标签不能为空",
		"DictValue.required": "字典键值不能为空",
		"DictType.required":  "字典类型不能为空",
		"Status.required":    "状态不能为空",
		"Status.oneof":       "状态只能是0或1",
		"IsDefault.oneof":    "是否默认只能是0或1",
	}
}
