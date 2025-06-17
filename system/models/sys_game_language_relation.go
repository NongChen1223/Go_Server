package models

// SysGameLanguageRelation 游戏与语言关联表（多对多）
type SysGameLanguageRelation struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement;comment:'ID - 主键，自动递增'" json:"id"`
	GameID   uint64 `gorm:"not null;uniqueIndex:uk_game_language,priority:1;index:idx_game_id;comment:'游戏ID - 关联sys_game表'" json:"game_id"`
	DictCode uint64 `gorm:"not null;uniqueIndex:uk_game_language,priority:2;index:idx_dict_code;comment:'字典编码 - 关联sys_dict_data表中game_language类型的数据'" json:"dict_code"`
}

// TableName 指定表名
func (SysGameLanguageRelation) TableName() string {
	return "sys_game_language_relation"
}
