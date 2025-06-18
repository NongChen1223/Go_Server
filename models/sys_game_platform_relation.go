package models

// SysGamePlatformRelation 游戏与平台关联表（多对多）
type SysGamePlatformRelation struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement;comment:'ID - 主键，自动递增'" json:"id"`
	GameID   uint64 `gorm:"not null;uniqueIndex:uk_game_platform,priority:1;index:idx_game_id;comment:'游戏ID - 关联sys_game表'" json:"game_id"`
	DictCode uint64 `gorm:"not null;uniqueIndex:uk_game_platform,priority:2;index:idx_dict_code;comment:'字典编码 - 关联sys_dict_data表中platform_type类型的数据'" json:"dict_code"`
}

// TableName 指定表名
func (SysGamePlatformRelation) TableName() string {
	return "sys_game_platform_relation"
}
