package dto

import (
	"go_server/common"
	"go_server/config"
	"go_server/utils"
	"time"
)

// GameReq 游戏请求结构体
// 用于创建和更新游戏时接收前端传来的数据
type GameReq struct {
	GameID       uint64             `json:"game_id,omitempty" example:"1"`                      // 游戏ID，更新时需要，创建时忽略
	NameZh       string             `json:"name_zh" binding:"required" example:"原神"`            // 游戏中文名称，必填
	NameEn       *string            `json:"name_en" example:"Genshin Impact"`                   // 游戏英文名称，可选
	ReleaseDate  *config.CustomTime `json:"release_date" example:"2020-09-28"`                  // 游戏发布日期，可选，格式：YYYY-MM-DD
	Description  *string            `json:"description" example:"开放世界冒险游戏"`                     // 游戏介绍，可选
	Rating       *float64           `json:"rating" example:"9.5"`                               // 游戏评分，可选
	Size         *string            `json:"size" example:"15GB"`                                // 游戏大小，可选
	Price        *float64           `json:"price" example:"0"`                                  // 游戏价格，可选
	PublisherID  *uint64            `json:"publisher_id" example:"1"`                           // 游戏厂商ID，可选
	StudioID     *uint64            `json:"studio_id" example:"1"`                              // 游戏工作室ID，可选
	ShutdownDate *config.CustomTime `json:"shutdown_date" example:"2030-12-31"`                 // 游戏停服日期，可选，格式：YYYY-MM-DD
	DemoVideo    *string            `json:"demo_video" example:"https://example.com/video.mp4"` // 游戏演示视频链接，可选
	Status       int                `json:"status" binding:"required,oneof=0 1" example:"1"`    // 状态，必填，只能是0或1

	// 关联数据
	TypeCodes     []uint64 `json:"type_codes"`     // 游戏类型编码数组
	PlatformCodes []uint64 `json:"platform_codes"` // 游戏平台编码数组
	LanguageCodes []uint64 `json:"language_codes"` // 游戏语言编码数组

	// 封面和截图
	Covers      []GameCoverReq      `json:"covers"`      // 游戏封面数组
	Screenshots []GameScreenshotReq `json:"screenshots"` // 游戏截图数组
}

// GameCoverReq 游戏封面请求结构体
type GameCoverReq struct {
	CoverID  uint64 `json:"cover_id,omitempty"`           // 封面ID，更新时需要
	CoverURL string `json:"cover_url" binding:"required"` // 封面图片URL，必填
	IsMain   int    `json:"is_main" binding:"oneof=0 1"`  // 是否主封面，必填
	Sort     int    `json:"sort"`                         // 排序，可选
}

// GameScreenshotReq 游戏截图请求结构体
type GameScreenshotReq struct {
	ScreenshotID  uint64 `json:"screenshot_id,omitempty"`           // 截图ID，更新时需要
	ScreenshotURL string `json:"screenshot_url" binding:"required"` // 截图URL，必填
	Sort          int    `json:"sort"`                              // 排序，可选
}

// GameRes 游戏响应结构体
// 用于返回游戏数据
type GameRes struct {
	GameID       uint64             `json:"game_id"`       // 游戏ID
	NameZh       string             `json:"name_zh"`       // 游戏中文名称
	NameEn       *string            `json:"name_en"`       // 游戏英文名称
	ReleaseDate  *config.CustomTime `json:"release_date"`  // 游戏发布日期
	Description  *string            `json:"description"`   // 游戏介绍
	Rating       *float64           `json:"rating"`        // 游戏评分
	Size         *string            `json:"size"`          // 游戏大小
	Price        *float64           `json:"price"`         // 游戏价格
	PublisherID  *uint64            `json:"publisher_id"`  // 游戏厂商ID
	StudioID     *uint64            `json:"studio_id"`     // 游戏工作室ID
	ShutdownDate *config.CustomTime `json:"shutdown_date"` // 游戏停服日期
	CommentCount int                `json:"comment_count"` // 游戏评论数量
	LikeCount    int                `json:"like_count"`    // 游戏点赞数量
	DemoVideo    *string            `json:"demo_video"`    // 游戏演示视频链接
	Status       int                `json:"status"`        // 状态

	// 关联数据
	PublisherName *string             `json:"publisher_name"` // 厂商名称
	StudioName    *string             `json:"studio_name"`    // 工作室名称
	Types         []GameTypeInfo      `json:"types"`          // 游戏类型信息
	Platforms     []GamePlatformInfo  `json:"platforms"`      // 游戏平台信息
	Languages     []GameLanguageInfo  `json:"languages"`      // 游戏语言信息
	Covers        []GameCoverRes      `json:"covers"`         // 游戏封面
	Screenshots   []GameScreenshotRes `json:"screenshots"`    // 游戏截图

	common.BaseEntity // 嵌入基础审计字段
}

// GameTypeInfo 游戏类型信息
type GameTypeInfo struct {
	DictCode  uint64 `json:"dict_code"`  // 字典编码
	DictLabel string `json:"dict_label"` // 字典标签
	DictValue string `json:"dict_value"` // 字典值
}

// GamePlatformInfo 游戏平台信息
type GamePlatformInfo struct {
	DictCode  uint64 `json:"dict_code"`  // 字典编码
	DictLabel string `json:"dict_label"` // 字典标签
	DictValue string `json:"dict_value"` // 字典值
}

// GameLanguageInfo 游戏语言信息
type GameLanguageInfo struct {
	DictCode  uint64 `json:"dict_code"`  // 字典编码
	DictLabel string `json:"dict_label"` // 字典标签
	DictValue string `json:"dict_value"` // 字典值
}

// GameCoverRes 游戏封面响应结构体
type GameCoverRes struct {
	CoverID    uint64     `json:"cover_id"`    // 封面ID
	CoverURL   string     `json:"cover_url"`   // 封面图片URL
	IsMain     int        `json:"is_main"`     // 是否主封面
	Sort       int        `json:"sort"`        // 排序
	CreateTime *time.Time `json:"create_time"` // 创建时间
}

// GameScreenshotRes 游戏截图响应结构体
type GameScreenshotRes struct {
	ScreenshotID  uint64     `json:"screenshot_id"`  // 截图ID
	ScreenshotURL string     `json:"screenshot_url"` // 截图URL
	Sort          int        `json:"sort"`           // 排序
	CreateTime    *time.Time `json:"create_time"`    // 创建时间
}

// GameQuery 游戏查询参数
// 用于查询游戏列表时的筛选条件
type GameQuery struct {
	NameZh              string   `form:"name_zh"`       // 游戏中文名称，支持模糊查询
	NameEn              string   `form:"name_en"`       // 游戏英文名称，支持模糊查询
	PublisherID         *uint64  `form:"publisher_id"`  // 厂商ID，精确查询
	StudioID            *uint64  `form:"studio_id"`     // 工作室ID，精确查询
	TypeCode            *uint64  `form:"type_code"`     // 游戏类型编码，精确查询
	PlatformCode        *uint64  `form:"platform_code"` // 游戏平台编码，精确查询
	LanguageCode        *uint64  `form:"language_code"` // 游戏语言编码，精确查询
	MinRating           *float64 `form:"min_rating"`    // 最低评分
	MaxRating           *float64 `form:"max_rating"`    // 最高评分
	MinPrice            *float64 `form:"min_price"`     // 最低价格
	MaxPrice            *float64 `form:"max_price"`     // 最高价格
	common.StatusEntity          // 嵌入状态字段
	common.PageQuery             // 嵌入分页字段
}

// GetMessages 自定义验证错误信息
// 当参数验证失败时，返回友好的中文错误信息
func (req GameReq) GetMessages() utils.ValidatorMessages {
	return utils.ValidatorMessages{
		"NameZh.required": "游戏中文名称不能为空",
		"Status.required": "状态不能为空",
		"Status.oneof":    "状态只能是0或1",
	}
}

// GetMessages 游戏封面的验证错误信息
func (req GameCoverReq) GetMessages() utils.ValidatorMessages {
	return utils.ValidatorMessages{
		"CoverURL.required": "封面图片URL不能为空",
		"IsMain.oneof":      "是否主封面只能是0或1",
	}
}

// GetMessages 游戏截图的验证错误信息
func (req GameScreenshotReq) GetMessages() utils.ValidatorMessages {
	return utils.ValidatorMessages{
		"ScreenshotURL.required": "截图URL不能为空",
	}
}
