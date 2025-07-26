// Package services 游戏管理的业务逻辑层
// 负责游戏的CRUD操作，包括业务验证和数据转换
//
// 业务说明：
// 1. 游戏管理是整个系统的核心功能，涉及复杂的关联数据处理
// 2. 一个游戏可以有多种类型、支持多个平台、支持多种语言
// 3. 一个游戏可以有多张封面图片和截图
// 4. 所有操作都需要保证数据一致性，使用事务处理
package services

import (
	"errors"                       // Go标准库，用于创建错误对象
	"go_server/global"             // 全局变量，包含数据库连接等
	"go_server/internal/admin/dto" // 数据传输对象，定义请求和响应结构
	"go_server/models"             // 数据模型，定义数据库表结构
	"gorm.io/gorm"                 // GORM ORM框架，用于数据库操作
	"time"                         // Go标准库，用于时间处理
)

// CreateGame 创建游戏
//
// 业务逻辑说明：
// 1. 数据验证：检查游戏名称是否重复（中文名必须唯一，英文名如果有也必须唯一）
// 2. 事务处理：由于涉及多个表的操作，必须使用数据库事务保证数据一致性
// 3. 关联数据：游戏创建后需要处理类型、平台、语言、封面、截图等关联数据
// 4. 审计信息：记录创建者信息，便于后续追踪和管理
//
// 参数说明：
// - req: 游戏创建请求，包含游戏的所有信息
// - adminName: 当前操作的管理员名称，用于审计
//
// 返回值：
// - error: 如果操作失败返回错误信息，成功返回nil
func CreateGame(req dto.GameReq, adminName string) error {
	// ========== 第一步：数据验证 ==========

	// 业务规则：游戏中文名称必须唯一
	// 原因：避免用户混淆，保证游戏标识的唯一性
	var count int64 // 声明int64类型变量，用于存储查询结果数量

	// GORM语法详解：
	// Model(&models.SysGame{}) - 指定要操作的数据库模型/表
	// Where("name_zh = ?", req.NameZh) - 添加WHERE条件，?是占位符，防止SQL注入
	// Count(&count) - 统计符合条件的记录数量，结果存入count变量
	global.DB.Model(&models.SysGame{}).Where("name_zh = ?", req.NameZh).Count(&count)

	// Go语法：if条件判断，当count大于0时表示已存在同名游戏
	if count > 0 {
		// errors.New()创建一个新的错误对象，包含指定的错误信息
		return errors.New("游戏中文名称已存在")
	}

	// 业务规则：如果提供了英文名称，英文名称也必须唯一
	// 原因：英文名称通常用于国际化展示，也需要保证唯一性
	// Go语法详解：
	// req.NameEn != nil - 检查指针是否不为空（NameEn是*string类型）
	// *req.NameEn != "" - 解引用指针并检查字符串是否不为空
	// && 是逻辑与操作符，两个条件都为true时整个表达式才为true
	if req.NameEn != nil && *req.NameEn != "" {
		// 重复使用count变量，Go允许变量重用
		global.DB.Model(&models.SysGame{}).Where("name_en = ?", *req.NameEn).Count(&count)
		if count > 0 {
			return errors.New("游戏英文名称已存在")
		}
	}

	// ========== 第二步：开启数据库事务 ==========

	// 业务原因：游戏创建涉及多个表的操作（游戏主表、类型关联表、平台关联表等）
	// 必须使用事务保证数据一致性，要么全部成功，要么全部失败

	// GORM语法：Begin()开启一个新的数据库事务
	// tx是事务对象，后续所有数据库操作都通过tx进行
	tx := global.DB.Begin()

	// Go语法：defer关键字确保函数返回前执行指定代码
	// 这是Go的资源清理模式，确保事务在异常情况下能够回滚
	defer func() {
		// recover()捕获panic异常，类似于其他语言的try-catch
		if r := recover(); r != nil {
			// 如果发生panic，回滚事务
			tx.Rollback()
		}
	}()

	// ========== 第三步：创建游戏主记录 ==========

	// Go语法：&models.SysGame{} 创建SysGame结构体的指针
	// 花括号{}内是结构体字段初始化，使用字段名: 值的格式
	// 在Go中，数据库操作通常使用指针，因为GORM需要修改结构体的字段（如自动设置ID）
	game := &models.SysGame{
		// 基础信息字段：直接从请求参数赋值
		NameZh:      req.NameZh,      // 游戏中文名称（必填）
		NameEn:      req.NameEn,      // 游戏英文名称（可选，指针类型）
		ReleaseDate: req.ReleaseDate, // 发布日期（可选，指针类型）
		Description: req.Description, // 游戏介绍（可选，指针类型）
		Rating:      req.Rating,      // 游戏评分（可选，指针类型）
		Size:        req.Size,        // 游戏大小（可选，指针类型）
		Price:       req.Price,       // 游戏价格（可选，指针类型）

		// 关联字段：外键关联到其他表
		PublisherID: req.PublisherID, // 厂商ID（可选，指针类型）
		StudioID:    req.StudioID,    // 工作室ID（可选，指针类型）

		// 其他字段
		ShutdownDate: req.ShutdownDate, // 停服日期（可选，指针类型）
		DemoVideo:    req.DemoVideo,    // 演示视频（可选，指针类型）
		Status:       req.Status,       // 状态（必填，int类型）

		// 审计字段：记录操作者信息
		CreateBy: adminName, // 创建者
		UpdateBy: adminName, // 更新者（创建时与创建者相同）
		// 注意：CreateTime和UpdateTime字段在数据库层面有默认值，这里不需要设置
	}

	// GORM语法：Create()方法插入新记录到数据库
	// tx.Create(game) 执行插入操作，game是指向结构体的指针
	// .Error 获取操作过程中的错误（如果有）
	// Go语法：:= 是短变量声明，等价于 var err error = tx.Create(game).Error
	if err := tx.Create(game).Error; err != nil {
		// 如果插入失败，回滚事务并返回错误
		tx.Rollback()
		return err
	}

	// 重要：执行Create后，GORM会自动将生成的主键ID设置到game.GameID字段
	// 这个ID将用于后续的关联数据插入

	// ========== 第四步：处理关联数据 ==========

	// 业务说明：游戏系统采用多对多关系设计
	// 一个游戏可以属于多个类型（如：既是RPG又是动作游戏）
	// 一个游戏可以支持多个平台（如：同时支持PC和PlayStation）
	// 一个游戏可以支持多种语言（如：同时支持中文和英文）

	// 处理游戏类型关联（多对多关系）
	// 业务逻辑：将游戏与选中的类型建立关联关系
	// 参数说明：tx(事务对象), game.GameID(刚创建的游戏ID), req.TypeCodes(类型编码数组)
	if err := createGameTypeRelations(tx, game.GameID, req.TypeCodes); err != nil {
		tx.Rollback() // 失败时回滚事务
		return err
	}

	// 处理游戏平台关联（多对多关系）
	// 业务逻辑：将游戏与支持的平台建立关联关系
	if err := createGamePlatformRelations(tx, game.GameID, req.PlatformCodes); err != nil {
		tx.Rollback()
		return err
	}

	// 处理游戏语言关联（多对多关系）
	// 业务逻辑：将游戏与支持的语言建立关联关系
	if err := createGameLanguageRelations(tx, game.GameID, req.LanguageCodes); err != nil {
		tx.Rollback()
		return err
	}

	// ========== 第五步：处理媒体文件 ==========

	// 处理游戏封面（一对多关系）
	// 业务逻辑：一个游戏可以有多张封面，其中一张可以设为主封面
	// 用途：游戏列表展示、详情页展示等
	if err := createGameCovers(tx, game.GameID, req.Covers); err != nil {
		tx.Rollback()
		return err
	}

	// 处理游戏截图（一对多关系）
	// 业务逻辑：一个游戏可以有多张截图，用于展示游戏画面
	// 用途：游戏详情页的画面展示
	if err := createGameScreenshots(tx, game.GameID, req.Screenshots); err != nil {
		tx.Rollback()
		return err
	}

	// ========== 第六步：提交事务 ==========

	// 业务说明：所有操作都成功后，提交事务使数据永久保存
	// GORM语法：Commit()提交事务，.Error获取提交过程中的错误
	// 如果提交失败，数据库会自动回滚所有操作
	return tx.Commit().Error
}

// createGameTypeRelations 创建游戏类型关联
//
// 业务说明：
// 游戏类型采用多对多关系设计，通过中间表sys_game_type_relation实现
// 这样设计的好处：
// 1. 一个游戏可以属于多个类型（如：既是RPG又是策略游戏）
// 2. 便于按类型筛选游戏
// 3. 类型数据统一管理，通过字典表维护
//
// 参数说明：
// - tx: 数据库事务对象，确保操作的原子性
// - gameID: 游戏ID，刚创建的游戏的主键
// - typeCodes: 类型编码数组，来自字典数据表的dict_code字段
//
// 技术实现：
// 使用for range循环遍历类型编码数组，为每个类型创建一条关联记录
func createGameTypeRelations(tx *gorm.DB, gameID uint64, typeCodes []uint64) error {
	// Go语法：for range循环遍历切片
	// range typeCodes 返回两个值：索引和元素值
	// 使用 _ 忽略索引，只使用元素值typeCode
	for _, typeCode := range typeCodes {
		// 创建关联记录
		// &models.SysGameTypeRelation{} 创建结构体指针
		relation := &models.SysGameTypeRelation{
			GameID:   gameID,   // 游戏ID，建立与游戏的关联
			DictCode: typeCode, // 字典编码，关联到字典数据表
		}

		// 插入关联记录到数据库
		// 如果插入失败，立即返回错误，事务会在上层函数中回滚
		if err := tx.Create(relation).Error; err != nil {
			return err
		}
	}
	// 所有关联记录都创建成功，返回nil表示无错误
	return nil
}

// createGamePlatformRelations 创建游戏平台关联
//
// 业务说明：
// 游戏平台关联表示游戏支持哪些平台运行
// 例如：一个游戏可能同时支持PC、PlayStation 5、Xbox Series X等多个平台
// 这对用户选择游戏很重要，用户会根据自己拥有的平台来筛选游戏
//
// 数据结构：与类型关联相同，都是多对多关系，通过中间表实现
func createGamePlatformRelations(tx *gorm.DB, gameID uint64, platformCodes []uint64) error {
	// 遍历平台编码数组，为每个平台创建关联记录
	for _, platformCode := range platformCodes {
		relation := &models.SysGamePlatformRelation{
			GameID:   gameID,       // 游戏ID
			DictCode: platformCode, // 平台字典编码（如：PC、PS5、Xbox等）
		}
		if err := tx.Create(relation).Error; err != nil {
			return err
		}
	}
	return nil
}

// createGameLanguageRelations 创建游戏语言关联
//
// 业务说明：
// 游戏语言关联表示游戏支持哪些语言
// 这对国际化很重要，用户可以根据语言偏好筛选游戏
// 例如：一个游戏可能支持中文、英文、日文等多种语言
//
// 实际应用场景：
// 1. 游戏详情页显示支持的语言列表
// 2. 用户可以按语言筛选游戏
// 3. 为不同语言用户推荐合适的游戏
func createGameLanguageRelations(tx *gorm.DB, gameID uint64, languageCodes []uint64) error {
	// 遍历语言编码数组，为每种语言创建关联记录
	for _, languageCode := range languageCodes {
		relation := &models.SysGameLanguageRelation{
			GameID:   gameID,       // 游戏ID
			DictCode: languageCode, // 语言字典编码（如：中文、英文、日文等）
		}
		if err := tx.Create(relation).Error; err != nil {
			return err
		}
	}
	return nil
}

// createGameCovers 创建游戏封面
//
// 业务说明：
// 游戏封面是一对多关系，一个游戏可以有多张封面图片
// 封面的作用：
// 1. 主封面：在游戏列表中展示，吸引用户注意
// 2. 备用封面：在详情页或特殊场景下展示
// 3. 不同尺寸的封面：适配不同的显示场景
//
// 业务规则：
// - 一个游戏可以有多张封面，但只能有一张主封面
// - 封面有排序功能，控制展示顺序
// - 封面URL存储图片的访问路径
//
// 参数说明：
// - covers: 封面请求数组，包含URL、是否主封面、排序等信息
func createGameCovers(tx *gorm.DB, gameID uint64, covers []dto.GameCoverReq) error {
	// 遍历封面数组，为每张封面创建记录
	for _, cover := range covers {
		// 创建游戏封面记录
		// 注意：这里使用的是DTO中的GameCoverReq结构，需要转换为数据库模型
		gameCover := &models.SysGameCover{
			GameID:   gameID,         // 关联的游戏ID
			CoverURL: cover.CoverURL, // 封面图片URL
			IsMain:   cover.IsMain,   // 是否主封面（0否 1是）
			Sort:     cover.Sort,     // 排序值，数字越小越靠前
			// CreateTime字段在数据库层面有默认值，这里不需要设置
		}
		if err := tx.Create(gameCover).Error; err != nil {
			return err
		}
	}
	return nil
}

// createGameScreenshots 创建游戏截图
//
// 业务说明：
// 游戏截图也是一对多关系，一个游戏可以有多张截图
// 截图的作用：
// 1. 展示游戏的实际画面和玩法
// 2. 帮助用户了解游戏的视觉效果
// 3. 在游戏详情页形成图片画廊
//
// 业务规则：
// - 一个游戏可以有多张截图
// - 截图有排序功能，控制在详情页的展示顺序
// - 截图URL存储图片的访问路径
//
// 与封面的区别：
// - 封面主要用于列表展示和吸引用户
// - 截图主要用于详情展示和展现游戏内容
func createGameScreenshots(tx *gorm.DB, gameID uint64, screenshots []dto.GameScreenshotReq) error {
	// 遍历截图数组，为每张截图创建记录
	for _, screenshot := range screenshots {
		// 创建游戏截图记录
		gameScreenshot := &models.SysGameScreenshot{
			GameID:        gameID,                   // 关联的游戏ID
			ScreenshotURL: screenshot.ScreenshotURL, // 截图URL
			Sort:          screenshot.Sort,          // 排序值
			// CreateTime字段在数据库层面有默认值，这里不需要设置
		}
		if err := tx.Create(gameScreenshot).Error; err != nil {
			return err
		}
	}
	return nil
}

// UpdateGame 更新游戏
//
// 业务逻辑说明：
// 更新操作比创建更复杂，需要考虑以下几个方面：
// 1. 数据验证：确保要更新的记录存在
// 2. 唯一性检查：确保修改后的名称不与其他记录冲突
// 3. 关联数据处理：删除旧的关联数据，创建新的关联数据
// 4. 事务处理：保证所有操作的原子性
//
// 更新策略：
// 对于关联数据（类型、平台、语言、封面、截图），采用"删除重建"策略
// 虽然效率不是最高，但逻辑简单，数据一致性好
//
// 参数说明：
// - req: 更新请求，包含游戏ID和要更新的所有信息
// - adminName: 当前操作的管理员名称，用于审计
func UpdateGame(req dto.GameReq, adminName string) error {
	// ========== 第一步：验证游戏是否存在 ==========

	// 声明游戏模型变量，用于存储查询结果
	var game models.SysGame

	// GORM语法：First(&game, req.GameID) 根据主键查询单条记录
	// 等价于 SELECT * FROM sys_game WHERE game_id = req.GameID LIMIT 1
	// 查询结果会填充到game变量中
	err := global.DB.First(&game, req.GameID).Error
	if err != nil {
		// errors.Is() 用于判断错误类型，类似于其他语言的 instanceof
		// gorm.ErrRecordNotFound 是GORM预定义的"记录不存在"错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("游戏不存在")
		}
		// 如果是其他错误（如数据库连接错误），直接返回
		return err
	}

	// ========== 第二步：检查名称唯一性 ==========

	// 业务规则：只有当名称确实发生变化时，才需要检查唯一性
	// 这样可以避免不必要的数据库查询，提高性能

	// 检查中文名称是否发生变化
	// Go语法：!= 是不等于比较操作符
	if req.NameZh != game.NameZh {
		var count int64
		// 重要：查询条件中要排除当前记录本身（AND game_id != ?）
		// 否则会把自己也算进去，导致误判
		global.DB.Model(&models.SysGame{}).Where("name_zh = ? AND game_id != ?", req.NameZh, req.GameID).Count(&count)
		if count > 0 {
			return errors.New("游戏中文名称已存在")
		}
	}

	// 检查英文名称是否发生变化（更复杂，因为涉及指针比较）
	if req.NameEn != nil && *req.NameEn != "" {
		// 需要处理以下几种情况：
		// 1. 原来没有英文名称，现在添加了英文名称
		// 2. 原来有英文名称，现在修改了英文名称
		// 3. 原来有英文名称，现在没有变化（不需要检查）

		// Go语法：game.NameEn == nil 检查指针是否为空
		// *req.NameEn != *game.NameEn 解引用指针并比较字符串值
		if game.NameEn == nil || *req.NameEn != *game.NameEn {
			var count int64
			global.DB.Model(&models.SysGame{}).Where("name_en = ? AND game_id != ?", *req.NameEn, req.GameID).Count(&count)
			if count > 0 {
				return errors.New("游戏英文名称已存在")
			}
		}
	}

	// ========== 第三步：开启事务 ==========

	// 更新操作同样需要事务保护，因为涉及多个表的操作
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// ========== 第四步：更新游戏基本信息 ==========

	// 使用map[string]interface{}构建更新数据
	// 这种方式的优点：
	// 1. 可以动态构建更新字段
	// 2. 支持nil值的更新（如将某个字段设为NULL）
	// 3. 避免零值问题（Go的零值可能不是我们想要的）
	//
	// Go语法：map[string]interface{} 是一个映射类型
	// string是键的类型，interface{}是值的类型
	// interface{}可以存储任何类型的值，类似于其他语言的Object或Any
	updateData := map[string]interface{}{
		// 基础信息字段
		"name_zh":      req.NameZh,      // 游戏中文名称
		"name_en":      req.NameEn,      // 游戏英文名称（可能为nil）
		"release_date": req.ReleaseDate, // 发布日期（可能为nil）
		"description":  req.Description, // 游戏介绍（可能为nil）
		"rating":       req.Rating,      // 游戏评分（可能为nil）
		"size":         req.Size,        // 游戏大小（可能为nil）
		"price":        req.Price,       // 游戏价格（可能为nil）

		// 关联字段
		"publisher_id": req.PublisherID, // 厂商ID（可能为nil）
		"studio_id":    req.StudioID,    // 工作室ID（可能为nil）

		// 其他字段
		"shutdown_date": req.ShutdownDate, // 停服日期（可能为nil）
		"demo_video":    req.DemoVideo,    // 演示视频（可能为nil）
		"status":        req.Status,       // 状态

		// 审计字段
		"update_by":   adminName,  // 更新者
		"update_time": time.Now(), // 更新时间（手动设置当前时间）
	}

	// GORM语法：Updates()方法批量更新字段
	// Model(&game) 指定要更新的记录（通过主键识别）
	// Updates(updateData) 使用map中的数据更新对应字段
	// 注意：Updates会忽略零值，但我们使用map可以避免这个问题
	if err := tx.Model(&game).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return err
	}

	// ========== 第五步：更新关联数据（采用删除重建策略） ==========

	// 业务说明：
	// 对于多对多关系的更新，我们采用"删除重建"策略
	// 虽然不是最高效的方式，但逻辑简单，不容易出错
	//
	// 更高效的方式是：
	// 1. 比较新旧数据，找出要删除、新增、保留的记录
	// 2. 分别执行删除和新增操作
	// 但这种方式逻辑复杂，容易出错，对于大多数场景，删除重建已经足够

	// 删除原有的类型、平台、语言关联数据
	if err := deleteGameRelations(tx, req.GameID); err != nil {
		tx.Rollback()
		return err
	}

	// 重新创建类型关联数据
	// 使用之前定义的函数，保持代码复用
	if err := createGameTypeRelations(tx, req.GameID, req.TypeCodes); err != nil {
		tx.Rollback()
		return err
	}

	// 重新创建平台关联数据
	if err := createGamePlatformRelations(tx, req.GameID, req.PlatformCodes); err != nil {
		tx.Rollback()
		return err
	}

	// 重新创建语言关联数据
	if err := createGameLanguageRelations(tx, req.GameID, req.LanguageCodes); err != nil {
		tx.Rollback()
		return err
	}

	// ========== 第六步：更新媒体文件（封面和截图） ==========

	// 删除原有的封面和截图数据
	// 业务考虑：封面和截图的更新频率相对较低，删除重建是合理的
	if err := deleteGameCoversAndScreenshots(tx, req.GameID); err != nil {
		tx.Rollback()
		return err
	}

	// 重新创建封面数据
	if err := createGameCovers(tx, req.GameID, req.Covers); err != nil {
		tx.Rollback()
		return err
	}

	// 重新创建截图数据
	if err := createGameScreenshots(tx, req.GameID, req.Screenshots); err != nil {
		tx.Rollback()
		return err
	}

	// ========== 第七步：提交事务 ==========

	// 所有更新操作都成功后，提交事务
	return tx.Commit().Error
}

// deleteGameRelations 删除游戏关联数据
//
// 业务说明：
// 这个函数负责删除游戏的所有多对多关联数据
// 包括：游戏类型关联、游戏平台关联、游戏语言关联
//
// 使用场景：
// 1. 更新游戏时，删除旧的关联数据
// 2. 删除游戏时，清理所有关联数据
//
// 技术说明：
// 使用事务确保所有删除操作的原子性
// 如果任何一个删除操作失败，整个事务都会回滚
func deleteGameRelations(tx *gorm.DB, gameID uint64) error {
	// 删除游戏类型关联
	// GORM语法：Where().Delete() 执行条件删除
	// 等价于 DELETE FROM sys_game_type_relation WHERE game_id = ?
	// &models.SysGameTypeRelation{} 指定要删除的表
	if err := tx.Where("game_id = ?", gameID).Delete(&models.SysGameTypeRelation{}).Error; err != nil {
		return err
	}

	// 删除游戏平台关联
	// 同样的删除逻辑，针对不同的关联表
	if err := tx.Where("game_id = ?", gameID).Delete(&models.SysGamePlatformRelation{}).Error; err != nil {
		return err
	}

	// 删除游戏语言关联
	if err := tx.Where("game_id = ?", gameID).Delete(&models.SysGameLanguageRelation{}).Error; err != nil {
		return err
	}

	// 所有关联数据删除成功
	return nil
}

// deleteGameCoversAndScreenshots 删除游戏封面和截图
//
// 业务说明：
// 这个函数负责删除游戏的媒体文件记录
// 注意：这里只删除数据库记录，不删除实际的图片文件
//
// 实际项目中的考虑：
// 1. 图片文件的删除通常需要额外的处理（如从云存储删除）
// 2. 可能需要考虑图片的引用计数（多个地方使用同一张图片）
// 3. 可能需要异步删除图片文件，避免影响用户体验
//
// 当前实现：
// 只删除数据库记录，图片文件的清理可以通过定时任务处理
func deleteGameCoversAndScreenshots(tx *gorm.DB, gameID uint64) error {
	// 删除游戏封面记录
	// 注意：这不会删除实际的图片文件，只删除数据库中的记录
	if err := tx.Where("game_id = ?", gameID).Delete(&models.SysGameCover{}).Error; err != nil {
		return err
	}

	// 删除游戏截图记录
	if err := tx.Where("game_id = ?", gameID).Delete(&models.SysGameScreenshot{}).Error; err != nil {
		return err
	}

	return nil
}

// GetGameDetail 获取游戏详情
//
// 业务说明：
// 这个函数用于获取游戏的完整详细信息，包括：
// 1. 游戏基本信息
// 2. 关联的厂商和工作室信息
// 3. 游戏类型、平台、语言信息
// 4. 游戏封面和截图
//
// 使用场景：
// 1. 管理后台的游戏详情页面
// 2. 编辑游戏时的数据回显
// 3. 前台游戏详情页面的数据展示
//
// 性能考虑：
// 这个函数会执行多次数据库查询，在高并发场景下可能需要优化
// 可以考虑使用JOIN查询或缓存来提高性能
//
// 参数说明：
// - gameID: 游戏ID，用于查询指定的游戏
//
// 返回值：
// - *dto.GameRes: 游戏详情响应对象，包含所有相关信息
// - error: 错误信息，如果游戏不存在或查询失败
func GetGameDetail(gameID uint64) (*dto.GameRes, error) {
	// ========== 第一步：查询游戏基本信息 ==========

	var game models.SysGame
	// GORM语法：First(&game, gameID) 根据主键查询单条记录
	// 这是GORM的简化写法，等价于 First(&game, "game_id = ?", gameID)
	err := global.DB.First(&game, gameID).Error
	if err != nil {
		// 检查是否是"记录不存在"错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("游戏不存在")
		}
		// 其他错误（如数据库连接错误）直接返回
		return nil, err
	}

	// ========== 第二步：构建响应对象 ==========

	// 将数据库模型转换为响应DTO
	// 这种手动赋值的方式虽然代码较多，但清晰明了，便于维护
	// 也可以使用反射或第三方库自动转换，但会牺牲性能和可读性
	res := &dto.GameRes{
		// 基础信息字段：直接从数据库模型复制
		GameID:       game.GameID,       // 游戏ID
		NameZh:       game.NameZh,       // 中文名称
		NameEn:       game.NameEn,       // 英文名称（可能为nil）
		ReleaseDate:  game.ReleaseDate,  // 发布日期（可能为nil）
		Description:  game.Description,  // 游戏介绍（可能为nil）
		Rating:       game.Rating,       // 游戏评分（可能为nil）
		Size:         game.Size,         // 游戏大小（可能为nil）
		Price:        game.Price,        // 游戏价格（可能为nil）
		PublisherID:  game.PublisherID,  // 厂商ID（可能为nil）
		StudioID:     game.StudioID,     // 工作室ID（可能为nil）
		ShutdownDate: game.ShutdownDate, // 停服日期（可能为nil）
		CommentCount: game.CommentCount, // 评论数量
		LikeCount:    game.LikeCount,    // 点赞数量
		DemoVideo:    game.DemoVideo,    // 演示视频（可能为nil）
		Status:       game.Status,       // 状态
	}

	// 赋值嵌入的公共字段（审计信息）
	// 这些字段来自common.BaseEntity，用于记录数据的创建和修改信息
	res.CreateTime = game.CreateTime // 创建时间
	res.UpdateTime = game.UpdateTime // 更新时间
	res.CreateBy = game.CreateBy     // 创建者
	res.UpdateBy = game.UpdateBy     // 更新者

	// ========== 第三步：获取关联信息 ==========

	// 获取厂商名称
	// 业务说明：如果游戏关联了厂商，需要获取厂商的名称用于显示
	// Go语法：game.PublisherID != nil 检查指针是否不为空
	if game.PublisherID != nil {
		var publisher models.SysPublisher
		// *game.PublisherID 解引用指针获取实际的ID值
		// err == nil 表示查询成功（没有错误）
		if err := global.DB.First(&publisher, *game.PublisherID).Error; err == nil {
			// &publisher.PublisherName 获取字符串的指针
			// 因为res.PublisherName是*string类型（指针类型）
			res.PublisherName = &publisher.PublisherName
		}
		// 如果查询失败（如厂商被删除），PublisherName保持为nil
		// 前端可以根据nil值判断是否显示厂商信息
	}

	// 获取工作室名称
	// 业务说明：工作室功能暂未实现，预留接口
	// TODO: 实现工作室模型后取消注释
	// if game.StudioID != nil {
	//     var studio models.SysStudio
	//     if err := global.DB.First(&studio, *game.StudioID).Error; err == nil {
	//         res.StudioName = &studio.StudioName
	//     }
	// }

	// ========== 第四步：获取多对多关联数据 ==========

	// 获取游戏类型信息
	// 这些函数会执行JOIN查询，获取字典表中的详细信息
	res.Types = getGameTypes(gameID)

	// 获取游戏平台信息
	res.Platforms = getGamePlatforms(gameID)

	// 获取游戏语言信息
	res.Languages = getGameLanguages(gameID)

	// ========== 第五步：获取媒体文件信息 ==========

	// 获取游戏封面
	// 返回所有封面，按主封面优先、排序值升序排列
	res.Covers = getGameCovers(gameID)

	// 获取游戏截图
	// 返回所有截图，按排序值升序排列
	res.Screenshots = getGameScreenshots(gameID)

	// 返回完整的游戏详情信息
	return res, nil
}

// getGameTypes 获取游戏类型信息
//
// 业务说明：
// 通过JOIN查询获取游戏的所有类型信息
// 不仅返回类型编码，还返回类型的显示名称和值
//
// 技术实现：
// 使用原生SQL查询，因为涉及多表JOIN，原生SQL更直观
// 也可以用GORM的关联查询，但对于这种简单的JOIN，原生SQL更高效
//
// 查询逻辑：
// 1. 从关联表中找到该游戏的所有类型编码
// 2. JOIN字典数据表获取类型的详细信息
// 3. 只返回状态为正常(status=1)的类型
// 4. 按字典排序字段排序
func getGameTypes(gameID uint64) []dto.GameTypeInfo {
	// 声明结果切片，用于存储查询结果
	var types []dto.GameTypeInfo

	// 原生SQL查询
	// 使用反引号(`)定义多行字符串，保持SQL的可读性
	query := `
		SELECT d.dict_code, d.dict_label, d.dict_value
		FROM sys_game_type_relation r
		JOIN sys_dict_data d ON r.dict_code = d.dict_code
		WHERE r.game_id = ? AND d.status = 1
		ORDER BY d.dict_sort
	`

	// GORM语法：Raw()执行原生SQL查询
	// Scan(&types)将查询结果扫描到结构体切片中
	// GORM会自动根据字段名映射到结构体字段
	global.DB.Raw(query, gameID).Scan(&types)
	return types
}

// getGamePlatforms 获取游戏平台信息
//
// 业务说明：
// 获取游戏支持的所有平台信息
// 平台信息对用户很重要，用户会根据自己拥有的设备选择游戏
//
// 查询逻辑：与getGameTypes相同，只是查询的表不同
func getGamePlatforms(gameID uint64) []dto.GamePlatformInfo {
	var platforms []dto.GamePlatformInfo

	// 查询游戏平台关联表和字典数据表
	query := `
		SELECT d.dict_code, d.dict_label, d.dict_value
		FROM sys_game_platform_relation r
		JOIN sys_dict_data d ON r.dict_code = d.dict_code
		WHERE r.game_id = ? AND d.status = 1
		ORDER BY d.dict_sort
	`

	global.DB.Raw(query, gameID).Scan(&platforms)
	return platforms
}

// getGameLanguages 获取游戏语言信息
//
// 业务说明：
// 获取游戏支持的所有语言信息
// 语言支持对国际化用户很重要
//
// 实际应用：
// 1. 在游戏详情页显示支持的语言
// 2. 用户可以按语言筛选游戏
// 3. 为不同语言用户推荐合适的游戏
func getGameLanguages(gameID uint64) []dto.GameLanguageInfo {
	var languages []dto.GameLanguageInfo

	// 查询游戏语言关联表和字典数据表
	query := `
		SELECT d.dict_code, d.dict_label, d.dict_value
		FROM sys_game_language_relation r
		JOIN sys_dict_data d ON r.dict_code = d.dict_code
		WHERE r.game_id = ? AND d.status = 1
		ORDER BY d.dict_sort
	`

	global.DB.Raw(query, gameID).Scan(&languages)
	return languages
}

// getGameCovers 获取游戏封面
//
// 业务说明：
// 获取游戏的所有封面图片信息
// 封面图片用于游戏的视觉展示，是吸引用户的重要元素
//
// 排序规则：
// 1. 主封面优先（is_main DESC）：主封面排在最前面
// 2. 按排序值升序（sort ASC）：排序值小的在前面
//
// 使用场景：
// 1. 游戏列表：通常只显示主封面
// 2. 游戏详情：显示所有封面，用户可以切换查看
// 3. 轮播图：按顺序展示多张封面
func getGameCovers(gameID uint64) []dto.GameCoverRes {
	// 查询数据库中的封面记录
	var covers []models.SysGameCover

	// GORM语法：Where().Order().Find() 链式调用
	// Where("game_id = ?", gameID) - 查询条件
	// Order("is_main DESC, sort ASC") - 排序规则，多个字段用逗号分隔
	// Find(&covers) - 查询多条记录
	global.DB.Where("game_id = ?", gameID).Order("is_main DESC, sort ASC").Find(&covers)

	// 将数据库模型转换为响应DTO
	// 使用切片存储转换后的结果
	var result []dto.GameCoverRes

	// Go语法：for range循环遍历切片
	for _, cover := range covers {
		// append()函数向切片添加元素
		// 这里创建新的DTO对象并添加到结果切片中
		result = append(result, dto.GameCoverRes{
			CoverID:    cover.CoverID,    // 封面ID
			CoverURL:   cover.CoverURL,   // 封面图片URL
			IsMain:     cover.IsMain,     // 是否主封面
			Sort:       cover.Sort,       // 排序值
			CreateTime: cover.CreateTime, // 创建时间
		})
	}
	return result
}

// getGameScreenshots 获取游戏截图
//
// 业务说明：
// 获取游戏的所有截图信息
// 截图用于展示游戏的实际画面，帮助用户了解游戏内容
//
// 排序规则：
// 按排序值升序排列，确保截图按指定顺序展示
//
// 使用场景：
// 1. 游戏详情页：以画廊形式展示所有截图
// 2. 预览功能：用户可以快速浏览游戏画面
// 3. 推荐系统：选择优质截图用于推荐展示
//
// 与封面的区别：
// 截图没有"主截图"的概念，所有截图地位相等
// 主要通过排序值控制展示顺序
func getGameScreenshots(gameID uint64) []dto.GameScreenshotRes {
	// 查询数据库中的截图记录
	var screenshots []models.SysGameScreenshot

	// 按排序值升序查询所有截图
	// 注意：这里只有一个排序字段，不像封面有两个排序字段
	global.DB.Where("game_id = ?", gameID).Order("sort ASC").Find(&screenshots)

	// 转换为响应DTO
	var result []dto.GameScreenshotRes
	for _, screenshot := range screenshots {
		result = append(result, dto.GameScreenshotRes{
			ScreenshotID:  screenshot.ScreenshotID,  // 截图ID
			ScreenshotURL: screenshot.ScreenshotURL, // 截图URL
			Sort:          screenshot.Sort,          // 排序值
			CreateTime:    screenshot.CreateTime,    // 创建时间
		})
	}
	return result
}

// GetGameList 获取游戏列表（支持分页和条件查询）
//
// 业务说明：
// 这是管理后台最常用的查询接口，支持多种筛选条件和分页
// 主要用于游戏管理页面的列表展示和搜索功能
//
// 查询特点：
// 1. 支持多种筛选条件的组合
// 2. 支持模糊查询和精确查询
// 3. 支持范围查询（如价格区间、评分区间）
// 4. 支持关联查询（如按类型、平台筛选）
// 5. 支持分页，避免一次性加载大量数据
//
// 性能考虑：
// 1. 使用动态查询构建，只添加有值的查询条件
// 2. 先统计总数，再查询分页数据
// 3. 列表查询不加载完整的关联数据，提高查询速度
//
// 参数说明：
// - query: 查询条件对象，包含各种筛选参数和分页参数
//
// 返回值：
// - int64: 符合条件的总记录数
// - []*dto.GameRes: 当前页的游戏列表
// - error: 错误信息
func GetGameList(query dto.GameQuery) (int64, []*dto.GameRes, error) {
	// 声明变量存储查询结果
	var games []models.SysGame // 游戏模型切片
	var total int64            // 总记录数

	// ========== 第一步：创建查询构建器 ==========

	// GORM语法：Model()指定要查询的表
	// 这里创建一个查询构建器，后续可以链式添加查询条件
	db := global.DB.Model(&models.SysGame{})

	// ========== 第二步：动态添加查询条件 ==========

	// 业务规则：只有当查询参数不为空时，才添加对应的查询条件
	// 这样可以实现灵活的组合查询

	// 游戏名称模糊查询
	// Go语法：query.NameZh != "" 检查字符串是否不为空
	if query.NameZh != "" {
		// LIKE查询用于模糊匹配，%是SQL通配符
		// "%"+query.NameZh+"%" 表示包含指定字符串的记录
		db = db.Where("name_zh LIKE ?", "%"+query.NameZh+"%")
	}
	if query.NameEn != "" {
		db = db.Where("name_en LIKE ?", "%"+query.NameEn+"%")
	}

	// 关联ID精确查询
	// 注意：这些字段是指针类型，需要检查是否为nil
	if query.PublisherID != nil {
		// *query.PublisherID 解引用指针获取实际值
		db = db.Where("publisher_id = ?", *query.PublisherID)
	}
	if query.StudioID != nil {
		db = db.Where("studio_id = ?", *query.StudioID)
	}

	// 状态查询
	// 注意：Status是int类型，0是有效值（停用状态）
	// 但在查询参数中，0通常表示"不筛选"，所以这里用 != 0 判断
	if query.Status != 0 {
		db = db.Where("status = ?", query.Status)
	}

	// 评分范围查询
	// 支持设置最低评分和最高评分，实现区间筛选
	if query.MinRating != nil {
		db = db.Where("rating >= ?", *query.MinRating) // 大于等于最低评分
	}
	if query.MaxRating != nil {
		db = db.Where("rating <= ?", *query.MaxRating) // 小于等于最高评分
	}

	// 价格范围查询
	// 与评分查询类似，支持价格区间筛选
	if query.MinPrice != nil {
		db = db.Where("price >= ?", *query.MinPrice)
	}
	if query.MaxPrice != nil {
		db = db.Where("price <= ?", *query.MaxPrice)
	}

	// ========== 第三步：处理关联查询 ==========

	// 业务说明：类型、平台、语言是多对多关系，需要通过子查询实现筛选
	// 使用IN子查询的方式，查找在关联表中存在指定关系的游戏

	// 按游戏类型筛选
	if query.TypeCode != nil {
		// 子查询：从类型关联表中查找包含指定类型的游戏ID
		db = db.Where("game_id IN (SELECT game_id FROM sys_game_type_relation WHERE dict_code = ?)", *query.TypeCode)
	}

	// 按游戏平台筛选
	if query.PlatformCode != nil {
		db = db.Where("game_id IN (SELECT game_id FROM sys_game_platform_relation WHERE dict_code = ?)", *query.PlatformCode)
	}

	// 按游戏语言筛选
	if query.LanguageCode != nil {
		db = db.Where("game_id IN (SELECT game_id FROM sys_game_language_relation WHERE dict_code = ?)", *query.LanguageCode)
	}

	// ========== 第四步：统计总记录数 ==========

	// 业务说明：分页查询需要知道总记录数，用于计算总页数
	// 注意：Count()操作不受Offset()和Limit()影响，统计的是所有符合条件的记录
	// GORM语法：Count(&total) 统计记录数并存入total变量
	db.Count(&total)

	// ========== 第五步：分页查询 ==========

	// 业务说明：分页是列表查询的重要功能，避免一次性加载大量数据
	// 分页参数验证：只有当页码和页面大小都大于0时才进行分页
	if query.PageNum > 0 && query.PageSize > 0 {
		// 计算偏移量（跳过的记录数）
		// 公式：(页码 - 1) × 每页数量
		// 例如：第2页，每页10条，偏移量 = (2-1) × 10 = 10
		offset := (query.PageNum - 1) * query.PageSize

		// GORM语法：Offset()设置偏移量，Limit()设置查询数量
		// 等价于SQL的 LIMIT query.PageSize OFFSET offset
		db = db.Offset(offset).Limit(query.PageSize)
	}

	// ========== 第六步：执行查询 ==========

	// 排序规则：按创建时间倒序，最新创建的游戏排在前面
	// 这符合管理后台的使用习惯，管理员通常关注最新的数据
	// GORM语法：Order()设置排序，Find()查询多条记录
	err := db.Order("create_time DESC").Find(&games).Error
	if err != nil {
		// 查询失败，返回错误
		return 0, nil, err
	}

	// ========== 第七步：转换为响应DTO ==========

	// 业务说明：列表查询不需要加载完整的关联数据，只加载必要信息
	// 这样可以提高查询性能，减少数据传输量

	// 声明结果切片，存储转换后的DTO对象
	var result []*dto.GameRes

	// 遍历查询结果，逐个转换为DTO
	for _, game := range games {
		// 创建游戏响应DTO
		// 注意：这里使用指针类型，因为result切片存储的是指针
		gameRes := &dto.GameRes{
			// 基础信息字段：直接复制
			GameID:       game.GameID,
			NameZh:       game.NameZh,
			NameEn:       game.NameEn,
			ReleaseDate:  game.ReleaseDate,
			Description:  game.Description,
			Rating:       game.Rating,
			Size:         game.Size,
			Price:        game.Price,
			PublisherID:  game.PublisherID,
			StudioID:     game.StudioID,
			ShutdownDate: game.ShutdownDate,
			CommentCount: game.CommentCount,
			LikeCount:    game.LikeCount,
			DemoVideo:    game.DemoVideo,
			Status:       game.Status,
		}

		// 赋值审计字段
		gameRes.CreateTime = game.CreateTime
		gameRes.UpdateTime = game.UpdateTime
		gameRes.CreateBy = game.CreateBy
		gameRes.UpdateBy = game.UpdateBy

		// ========== 第八步：获取必要的关联信息 ==========

		// 获取厂商名称（用于列表显示）
		// 业务考虑：厂商名称对用户很重要，需要在列表中显示
		if game.PublisherID != nil {
			var publisher models.SysPublisher
			if err := global.DB.First(&publisher, *game.PublisherID).Error; err == nil {
				gameRes.PublisherName = &publisher.PublisherName
			}
			// 如果查询厂商失败（如厂商被删除），PublisherName保持为nil
		}

		// 获取主封面（用于列表展示）
		// 业务考虑：列表中需要显示游戏封面，但只需要主封面即可
		// 性能考虑：不加载所有封面，只加载主封面，减少数据量
		var mainCover models.SysGameCover
		if err := global.DB.Where("game_id = ? AND is_main = 1", game.GameID).First(&mainCover).Error; err == nil {
			// 创建封面DTO切片，只包含主封面
			gameRes.Covers = []dto.GameCoverRes{
				{
					CoverID:    mainCover.CoverID,
					CoverURL:   mainCover.CoverURL,
					IsMain:     mainCover.IsMain,
					Sort:       mainCover.Sort,
					CreateTime: mainCover.CreateTime,
				},
			}
		}
		// 如果没有主封面，Covers字段保持为nil（空切片）

		// 获取游戏类型（用于列表展示）
		// 业务考虑：游戏类型对用户筛选和了解游戏很重要
		// 性能考虑：虽然会增加查询次数，但类型信息通常不多，影响可控
		// 优化建议：在高并发场景下，可以考虑使用JOIN查询或缓存优化
		gameRes.Types = getGameTypes(game.GameID)

		// 将转换后的DTO添加到结果切片中
		result = append(result, gameRes)
	}

	// ========== 第九步：返回查询结果 ==========

	// 返回三个值：
	// 1. total: 总记录数，用于前端计算分页信息
	// 2. result: 当前页的游戏列表
	// 3. nil: 表示没有错误
	return total, result, nil
}

// DeleteGame 删除游戏
//
// 业务说明：
// 游戏删除是一个复杂的操作，需要考虑数据完整性和业务规则
// 删除策略：物理删除（真正从数据库中删除记录）
//
// 业务规则：
// 1. 只能删除存在的游戏
// 2. 不能删除有用户评论的游戏（保护用户数据）
// 3. 删除游戏时必须同时删除所有关联数据
// 4. 使用事务确保删除操作的原子性
//
// 删除顺序：
// 1. 删除多对多关联数据（类型、平台、语言）
// 2. 删除一对多关联数据（封面、截图）
// 3. 删除游戏主记录
//
// 注意事项：
// 1. 这里只删除数据库记录，不删除实际的图片文件
// 2. 实际项目中可能需要考虑软删除（标记删除而不是物理删除）
// 3. 可能需要记录删除日志，便于审计和恢复
//
// 参数说明：
// - gameID: 要删除的游戏ID
//
// 返回值：
// - error: 删除失败时返回错误信息，成功返回nil
func DeleteGame(gameID uint64) error {
	// ========== 第一步：验证游戏是否存在 ==========

	var game models.SysGame
	err := global.DB.First(&game, gameID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("游戏不存在")
		}
		// 其他错误（如数据库连接错误）
		return err
	}

	// ========== 第二步：业务规则检查 ==========

	// 检查是否有用户评论
	// 业务规则：有用户评论的游戏不能删除，保护用户数据
	// 实际项目中可能还需要检查其他关联数据，如：
	// - 用户收藏
	// - 用户评分
	// - 购买记录等
	var commentCount int64
	global.DB.Model(&models.SysUserGameComment{}).Where("game_id = ?", gameID).Count(&commentCount)
	if commentCount > 0 {
		return errors.New("该游戏下还有用户评论，无法删除")
	}

	// ========== 第三步：开启事务删除 ==========

	// 删除操作涉及多个表，必须使用事务保证数据一致性
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// ========== 第四步：删除关联数据 ==========

	// 删除多对多关联数据（类型、平台、语言）
	// 必须先删除关联数据，再删除主记录，避免外键约束错误
	if err := deleteGameRelations(tx, gameID); err != nil {
		tx.Rollback()
		return err
	}

	// 删除一对多关联数据（封面和截图）
	if err := deleteGameCoversAndScreenshots(tx, gameID); err != nil {
		tx.Rollback()
		return err
	}

	// ========== 第五步：删除游戏主记录 ==========

	// 最后删除游戏主记录
	// 业务说明：必须在删除所有关联数据后再删除主记录
	// GORM语法：Delete(&game) 删除指定的记录
	// 由于game变量包含主键信息，GORM会根据主键删除对应记录
	if err := tx.Delete(&game).Error; err != nil {
		tx.Rollback()
		return err
	}

	// ========== 第六步：提交事务 ==========

	// 所有删除操作都成功后，提交事务
	// 如果提交失败，数据库会自动回滚所有操作
	return tx.Commit().Error
}
