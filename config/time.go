package config

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// CustomTime 自定义时间类型，支持多种格式
type CustomTime struct {
	time.Time
}

// UnmarshalJSON 自定义 JSON 解析，支持多种时间格式
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	// 移除引号
	timeStr := strings.Trim(string(data), `"`)
	if timeStr == "null" || timeStr == "" {
		return nil
	}

	// 支持的时间格式列表
	formats := []string{
		"2006-01-02",                // YYYY-MM-DD
		"2006-01-02 15:04:05",       // YYYY-MM-DD HH:MM:SS
		"2006-01-02T15:04:05Z",      // RFC3339 UTC
		"2006-01-02T15:04:05+08:00", // RFC3339 with timezone
		"2006-01-02T15:04:05.000Z",  // RFC3339 with milliseconds
		"2006/01/02",                // YYYY/MM/DD
		"01/02/2006",                // MM/DD/YYYY
	}

	// 尝试解析各种格式
	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			ct.Time = t
			return nil
		}
	}

	return fmt.Errorf("无法解析时间格式: %s，支持格式: YYYY-MM-DD, YYYY-MM-DD HH:MM:SS 等", timeStr)
}

// MarshalJSON 自定义 JSON 序列化，统一输出为 YYYY-MM-DD 格式
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(ct.Time.Format("2006-01-02"))
}

// CustomDateTime 自定义日期时间类型，输出包含时间
type CustomDateTime struct {
	time.Time
}

// UnmarshalJSON 自定义 JSON 解析
func (cdt *CustomDateTime) UnmarshalJSON(data []byte) error {
	// 移除引号
	timeStr := strings.Trim(string(data), `"`)
	if timeStr == "null" || timeStr == "" {
		return nil
	}

	// 支持的时间格式列表
	formats := []string{
		"2006-01-02 15:04:05",       // YYYY-MM-DD HH:MM:SS
		"2006-01-02T15:04:05Z",      // RFC3339 UTC
		"2006-01-02T15:04:05+08:00", // RFC3339 with timezone
		"2006-01-02T15:04:05.000Z",  // RFC3339 with milliseconds
		"2006-01-02",                // YYYY-MM-DD (自动补充时间为 00:00:00)
	}

	// 尝试解析各种格式
	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			cdt.Time = t
			return nil
		}
	}

	return fmt.Errorf("无法解析时间格式: %s", timeStr)
}

// MarshalJSON 自定义 JSON 序列化，输出为 YYYY-MM-DD HH:MM:SS 格式
func (cdt CustomDateTime) MarshalJSON() ([]byte, error) {
	if cdt.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(cdt.Time.Format("2006-01-02 15:04:05"))
}
