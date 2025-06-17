package models

// SysGameTypeRelation 游戏与类型关联表（多对多）
type SysGameTypeRelation struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement;comment:'ID - 主键，自动递增'" json:"id"`
	GameID   uint64 `gorm:"not null;uniqueIndex:uk_game_type,priority:1;index:idx_game_id;comment:'游戏ID - 关联sys_game表'" json:"game_id"`
	DictCode uint64 `gorm:"not null;uniqueIndex:uk_game_type,priority:2;index:idx_dict_code;comment:'字典编码 - 关联sys_dict_data表中game_type类型的数据'" json:"dict_code"`
}

// TableName 指定表名
func (SysGameTypeRelation) TableName() string {
	return "sys_game_type_relation"
}
