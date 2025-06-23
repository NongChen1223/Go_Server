package common

import "time"

// TimeEntity 时间字段组合
// 包含创建时间和更新时间
type TimeEntity struct {
	CreateTime *time.Time `json:"create_time"` // 创建时间
	UpdateTime *time.Time `json:"update_time"` // 更新时间
}

// OperatorEntity 操作者字段组合
// 包含创建者和更新者
type OperatorEntity struct {
	CreateBy string `json:"create_by"` // 创建者
	UpdateBy string `json:"update_by"` // 更新者
}

// BaseEntity 基础实体字段
// 包含所有模型通用的审计字段
type BaseEntity struct {
	TimeEntity     // 嵌入时间字段
	OperatorEntity // 嵌入操作者字段
}

// StatusEntity 状态实体字段
// 包含状态字段，适用于大多数需要启用/禁用功能的实体
type StatusEntity struct {
	Status int `json:"status"` // 状态（0停用 1正常）
}

// RemarkEntity 备注实体字段
// 包含备注字段，适用于需要添加备注说明的实体
type RemarkEntity struct {
	Remark *string `json:"remark"` // 备注
}

// SortEntity 排序实体字段
// 包含排序字段，适用于需要自定义排序的实体
type SortEntity struct {
	Sort int `json:"sort"` // 排序，数字越小越靠前
}

// PageQuery 通用分页查询参数
// 所有需要分页的查询都可以嵌入此结构体
type PageQuery struct {
	PageNum  int `form:"page_num" json:"page_num"`   // 页码，从1开始
	PageSize int `form:"page_size" json:"page_size"` // 每页数量
}

// GetDefaultPage 获取默认分页参数
func (q *PageQuery) GetDefaultPage() {
	if q.PageNum <= 0 {
		q.PageNum = 1 // 默认第1页
	}
	if q.PageSize <= 0 {
		q.PageSize = 10 // 默认每页10条
	}
	// 可以设置最大页面大小限制，防止查询过大数据量
	if q.PageSize > 100 {
		q.PageSize = 100
	}
}

// GetOffset 获取SQL查询的偏移量
func (q *PageQuery) GetOffset() int {
	return (q.PageNum - 1) * q.PageSize
}
