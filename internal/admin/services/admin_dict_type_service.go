// Package services 字典类型管理的业务逻辑层
// 负责字典类型的CRUD操作，包括业务验证和数据转换
package services

import (
	"errors"
	"go_server/global"
	"go_server/internal/admin/dto"
	"go_server/models"
	"gorm.io/gorm"
	"time"
)

// CreateDictType 创建字典类型
// 业务逻辑：检查重复性 -> 创建记录
func CreateDictType(req dto.DictTypeReq, adminName string) error {
	// 检查字典类型标识是否已存在
	// 字典类型标识是程序调用的唯一标识，不能重复
	var count int64
	global.DB.Model(&models.SysDictType{}).Where("dict_type = ?", req.DictType).Count(&count)
	if count > 0 {
		return errors.New("字典类型已存在")
	}

	// 创建字典类型记录
	dictType := &models.SysDictType{
		DictName: req.DictName,
		DictType: req.DictType,
		Status:   req.Status,
		Remark:   req.Remark,
		CreateBy: adminName, // 记录创建者，用于审计
		UpdateBy: adminName, // 创建时更新者与创建者相同
	}

	// 执行数据库插入操作
	return global.DB.Create(dictType).Error
}

// UpdateDictType 更新字典类型
// 业务逻辑：检查存在性 -> 检查重复性 -> 更新记录
func UpdateDictType(req dto.DictTypeReq, adminName string) error {
	// 检查要更新的字典类型是否存在
	var dictType models.SysDictType
	err := global.DB.First(&dictType, req.DictID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("字典类型不存在")
		}
		return err
	}

	// 如果修改了字典类型标识，需要检查新标识是否与其他记录冲突
	// 只有当新的dict_type与原来的不同时，才需要检查重复性
	if req.DictType != dictType.DictType {
		var count int64
		// 查找相同类型但不同ID的记录
		global.DB.Model(&models.SysDictType{}).Where("dict_type = ? AND dict_id != ?", req.DictType, req.DictID).Count(&count)
		if count > 0 {
			return errors.New("字典类型已存在")
		}
	}

	// 更新字典类型数据
	dictType.DictName = req.DictName
	dictType.DictType = req.DictType
	dictType.Status = req.Status
	dictType.Remark = req.Remark
	dictType.UpdateBy = adminName
	// 手动设置更新时间
	now := time.Now()
	dictType.UpdateTime = &now

	// 保存更新到数据库
	return global.DB.Save(&dictType).Error
}

// GetDictTypeDetail 获取字典类型详情
// 用于编辑页面回显数据或详情页面展示
func GetDictTypeDetail(dictID uint64) (*dto.DictTypeRes, error) {
	var dictType models.SysDictType
	// 根据主键查询单条记录
	err := global.DB.First(&dictType, dictID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("字典类型不存在")
		}
		return nil, err
	}

	// 将数据库模型转换为响应DTO
	// 使用逐个赋值的方式，更清晰明了
	res := &dto.DictTypeRes{
		DictID:   dictType.DictID,
		DictName: dictType.DictName,
		DictType: dictType.DictType,
		Status:   dictType.Status,
	}

	// 赋值嵌入的公共字段
	res.Remark = dictType.Remark
	res.CreateTime = dictType.CreateTime
	res.UpdateTime = dictType.UpdateTime
	res.CreateBy = dictType.CreateBy
	res.UpdateBy = dictType.UpdateBy

	return res, nil
}

// GetDictTypeList 获取字典类型列表（支持分页和条件查询）
// 这是管理后台最常用地查询接口
func GetDictTypeList(query dto.DictTypeQuery) (int64, []*dto.DictTypeRes, error) {
	var dictTypes []models.SysDictType
	var total int64

	// 创建查询构建器
	db := global.DB.Model(&models.SysDictType{})

	// 动态添加查询条件
	// 只有当查询参数不为空时才添加对应的WHERE条件
	if query.DictName != "" {
		db = db.Where("dict_name LIKE ?", "%"+query.DictName+"%") // 模糊查询
	}
	if query.DictType != "" {
		db = db.Where("dict_type LIKE ?", "%"+query.DictType+"%")
	}
	if query.Status != 0 {
		db = db.Where("status = ?", query.Status) // 精确查询
	}

	// 获取符合条件的总记录数（用于分页）
	db.Count(&total)

	// 分页查询
	if query.PageNum > 0 && query.PageSize > 0 {
		offset := (query.PageNum - 1) * query.PageSize // 分页计算公式
		db = db.Offset(offset).Limit(query.PageSize)
	}

	// 排序并执行查询
	err := db.Order("dict_id DESC").Find(&dictTypes).Error // 按ID降序，最新的在前面
	if err != nil {
		return 0, nil, err
	}

	// 转换为响应结构体数组
	var result []*dto.DictTypeRes
	for _, dictType := range dictTypes {
		// 创建单个响应对象
		res := &dto.DictTypeRes{
			DictID:   dictType.DictID,
			DictName: dictType.DictName,
			DictType: dictType.DictType,
			Status:   dictType.Status,
		}

		// 赋值嵌入的公共字段
		res.Remark = dictType.Remark
		res.CreateTime = dictType.CreateTime
		res.UpdateTime = dictType.UpdateTime
		res.CreateBy = dictType.CreateBy
		res.UpdateBy = dictType.UpdateBy

		// 添加到结果数组中
		result = append(result, res)
	}

	return total, result, nil
}

// DeleteDictType 删除字典类型
// 业务逻辑：检查存在性 -> 检查关联数据 -> 执行删除
func DeleteDictType(dictID uint64) error {
	// 检查字典类型是否存在
	var dictType models.SysDictType
	err := global.DB.First(&dictType, dictID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("字典类型不存在")
		}
		return err
	}

	// 检查是否有字典数据在使用此类型
	// 防止删除正在使用的字典类型，避免数据不一致
	var count int64
	global.DB.Model(&models.SysDictData{}).Where("dict_type = ?", dictType.DictType).Count(&count)
	if count > 0 {
		return errors.New("该字典类型下还有字典数据，无法删除")
	}

	// 执行物理删除
	return global.DB.Delete(&dictType).Error
}
