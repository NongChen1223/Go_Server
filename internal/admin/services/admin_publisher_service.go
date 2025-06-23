package services

import (
	"errors"
	"go_server/global"
	"go_server/internal/admin/dto"
	"go_server/models"
	"gorm.io/gorm"
	"time"
)

// 创建游戏厂商
func CreatePublisher(req dto.PublisherReq, adminName string) error {
	// 声明一个int64类型变量count用于存储查询结果数量
	var count int64

	// Model()方法指定要操作的数据库模型
	// Where()方法添加查询条件，第一个参数是条件语句，后面的参数是条件值
	// Count()方法统计符合条件的记录数量并存入count变量
	global.DB.Model(&models.SysPublisher{}).Where("publisher_name = ?", req.PublisherName).Count(&count)

	// 如果count大于0，表示已存在同名厂商
	if count > 0 {
		// errors.New()创建一个包含指定文本的error类型
		return errors.New("厂商名称已存在")
	}

	// 创建厂商记录
	// &符号表示获取结构体的指针，在Go中通常使用指针来操作数据库记录
	// 花括号{}内是结构体字段初始化，字段名: 值
	publisher := &models.SysPublisher{
		PublisherName: req.PublisherName, // 将请求中的厂商名称赋值给模型
		LogoURL:       req.LogoURL,       // LogoURL是指针类型(*string)，可以直接赋值nil或字符串指针
		Description:   req.Description,   // Description也是指针类型
		FoundedDate:   req.FoundedDate,   // FoundedDate是*time.Time类型
		Website:       req.Website,       // Website是*string类型
		Status:        req.Status,        // Status是int类型，表示厂商状态
		CreateBy:      adminName,         // 记录创建者
		UpdateBy:      adminName,         // 初始时更新者与创建者相同
	}

	// 执行数据库插入操作
	return global.DB.Create(publisher).Error
}

// UpdatePublisher 更新游戏厂商
// 业务逻辑：检查存在性 -> 检查重复性 -> 更新记录
func UpdatePublisher(req dto.PublisherReq, adminName string) error {
	// 检查要更新的游戏厂商是否存在
	var publisher models.SysPublisher
	err := global.DB.First(&publisher, req.PublisherID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("游戏厂商不存在")
		}
		return err
	}

	// 如果修改了厂商名称，需要检查新名称是否与其他记录冲突
	// 只有当新名称与原名称不同时才需要检查重复性
	if req.PublisherName != publisher.PublisherName {
		var count int64
		// 查找相同名称但不同ID的记录
		global.DB.Model(&models.SysPublisher{}).Where("publisher_name = ? AND publisher_id != ?", req.PublisherName, req.PublisherID).Count(&count)
		if count > 0 {
			return errors.New("厂商名称已存在")
		}
	}

	// 更新游戏厂商数据
	publisher.PublisherName = req.PublisherName
	publisher.LogoURL = req.LogoURL
	publisher.Description = req.Description
	publisher.FoundedDate = req.FoundedDate
	publisher.Website = req.Website
	publisher.Status = req.Status
	publisher.UpdateBy = adminName
	// 手动设置更新时间
	now := time.Now()
	publisher.UpdateTime = &now

	// 保存更新到数据库
	return global.DB.Save(&publisher).Error
}

// GetPublisherDetail 获取游戏厂商详情
// 用于编辑页面回显数据或详情页面展示
func GetPublisherDetail(publisherID uint64) (*dto.PublisherRes, error) {
	var publisher models.SysPublisher
	// 根据主键查询单条记录
	err := global.DB.First(&publisher, publisherID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("游戏厂商不存在")
		}
		return nil, err
	}

	// 将数据库模型转换为响应DTO
	// 使用逐个赋值的方式，更清晰明了
	res := &dto.PublisherRes{
		PublisherID:   publisher.PublisherID,
		PublisherName: publisher.PublisherName,
		LogoURL:       publisher.LogoURL,
		Description:   publisher.Description,
		FoundedDate:   publisher.FoundedDate,
		Website:       publisher.Website,
		Status:        publisher.Status,
	}

	// 赋值嵌入的公共字段
	res.CreateTime = publisher.CreateTime
	res.UpdateTime = publisher.UpdateTime
	res.CreateBy = publisher.CreateBy
	res.UpdateBy = publisher.UpdateBy

	return res, nil
}

// GetPublisherList 获取游戏厂商列表（支持分页和条件查询）
// 这是管理后台最常用的查询接口
func GetPublisherList(query dto.PublisherQuery) (int64, []*dto.PublisherRes, error) {
	var publishers []models.SysPublisher
	var total int64

	// 创建查询构建器
	db := global.DB.Model(&models.SysPublisher{})

	// 动态添加查询条件
	// 只有当查询参数不为空时才添加对应的WHERE条件
	if query.PublisherName != "" {
		db = db.Where("publisher_name LIKE ?", "%"+query.PublisherName+"%") // 模糊查询
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
	err := db.Order("publisher_id DESC").Find(&publishers).Error // 按ID降序，最新的在前面
	if err != nil {
		return 0, nil, err
	}

	// 转换为响应结构体数组
	var result []*dto.PublisherRes
	for _, publisher := range publishers {
		// 创建单个响应对象
		res := &dto.PublisherRes{
			PublisherID:   publisher.PublisherID,
			PublisherName: publisher.PublisherName,
			LogoURL:       publisher.LogoURL,
			Description:   publisher.Description,
			FoundedDate:   publisher.FoundedDate,
			Website:       publisher.Website,
			Status:        publisher.Status,
		}

		// 赋值嵌入的公共字段
		res.CreateTime = publisher.CreateTime
		res.UpdateTime = publisher.UpdateTime
		res.CreateBy = publisher.CreateBy
		res.UpdateBy = publisher.UpdateBy

		// 添加到结果数组中
		result = append(result, res)
	}

	return total, result, nil
}

// DeletePublisher 删除游戏厂商
// 业务逻辑：检查存在性 -> 检查关联数据 -> 执行删除
func DeletePublisher(publisherID uint64) error {
	// 检查游戏厂商是否存在
	var publisher models.SysPublisher
	err := global.DB.First(&publisher, publisherID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("游戏厂商不存在")
		}
		return err
	}

	// 检查是否有游戏在使用此厂商
	// 防止删除正在使用的厂商，避免数据不一致
	var count int64
	global.DB.Model(&models.SysGame{}).Where("publisher_id = ?", publisherID).Count(&count)
	if count > 0 {
		return errors.New("该厂商下还有游戏，无法删除")
	}

	// 执行物理删除
	return global.DB.Delete(&publisher).Error
}
