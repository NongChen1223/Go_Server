-- ----------------------------
-- 前台用户表 - 面向C端用户的核心数据表
-- 用途：存储前台用户的基础信息，支持用户注册、登录、个人资料管理等功能
-- 场景：用户注册、登录验证、个人信息展示、用户列表查询等
-- 业务说明：此表为前台用户的主表，与后台管理员表(sys_admin_user)完全分离
-- 扩展性：user_type字段预留用于后续业务扩展，可区分不同类型的前台用户
-- ----------------------------
DROP TABLE IF EXISTS `sys_user`; -- 如果表存在则删除，确保重新创建时不会冲突
CREATE TABLE `sys_user`
(
    -- 主键字段：使用bigint确保足够的ID空间，AUTO_INCREMENT自动递增
    `user_id`         bigint(20)                                                    NOT NULL AUTO_INCREMENT COMMENT '用户ID - 主键，自动递增，唯一标识每个用户',

    -- 基础信息字段：用户的核心身份信息
    `user_name`       varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL COMMENT '用户账号 - 登录用的唯一用户名，不可重复',
    `nick_name`       varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL COMMENT '用户昵称 - 显示用的名称，可以重复',
    `user_type`       varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci   NULL DEFAULT '00' COMMENT '用户类型 - 预留字段，用于区分不同类型用户（00:普通用户，01:VIP用户等）',

    -- 联系方式字段：用于用户联系和找回密码等功能
    `email`           varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NULL DEFAULT '' COMMENT '用户邮箱 - 可用于登录和找回密码',
    `phone_number`    varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NULL DEFAULT '' COMMENT '手机号码 - 11位手机号，可用于短信验证',

    -- 个人信息字段：用户的个性化信息
    `sex`             tinyint(1)                                                    NULL DEFAULT NULL COMMENT '用户性别 - 0:男性 1:女性 2:未知，NULL表示未设置',
    `avatar`          varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '头像地址 - 存储头像图片的URL路径',

    -- 安全相关字段：用户认证和安全控制
    `password`        varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '密码 - 加密存储的用户密码',
    `status`          tinyint(1)                                                    NULL DEFAULT 1 COMMENT '帐号状态 - 0:停用(禁止登录) 1:正常(可以登录)',
    `del_flag`        char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci      NULL DEFAULT '0' COMMENT '删除标志 - 0:存在 2:已删除(软删除，数据仍保留)',

    -- 登录记录字段：记录用户最后登录信息
    `login_ip`        varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NULL DEFAULT '' COMMENT '最后登陆IP - 记录用户最后一次登录的IP地址',
    `login_date`      datetime                                                      NULL DEFAULT NULL COMMENT '最后登陆时间 - 记录用户最后一次登录的时间',

    -- 系统字段：记录数据的创建和修改信息
    `create_time`     datetime                                                      NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间 - 记录创建时自动设置为当前时间',
    `update_time`     datetime                                                      NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间 - 记录更新时自动更新为当前时间',
    `remark`          varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '备注 - 管理员可添加的用户备注信息',

    -- 主键定义：指定主键字段和索引类型
    PRIMARY KEY (`user_id`) USING BTREE COMMENT '主键索引 - 使用B-Tree索引，确保用户ID的唯一性和查询效率',

    -- 唯一索引：确保字段值的唯一性
    UNIQUE KEY `uk_user_name` (`user_name`) COMMENT '用户名唯一索引 - 确保用户名不重复，支持快速登录验证',

    -- 普通索引：提高常用查询字段的查询效率
    INDEX `idx_email` (`email`) COMMENT '邮箱索引 - 加速邮箱查询，支持邮箱登录和找回密码功能',
    INDEX `idx_phone` (`phone_number`) COMMENT '手机号索引 - 加速手机号查询，支持手机号登录和验证',
    INDEX `idx_status` (`status`) COMMENT '状态索引 - 加速按状态筛选用户，如查询所有正常用户',
    INDEX `idx_user_type` (`user_type`) COMMENT '用户类型索引 - 加速按用户类型查询，支持不同类型用户的分类管理'

) ENGINE = InnoDB                                                                 -- 使用InnoDB存储引擎，支持事务、外键约束和行级锁
  AUTO_INCREMENT = 1                                                              -- 自增起始值设为1
  CHARACTER SET = utf8mb4                                                         -- 使用utf8mb4字符集，支持emoji和特殊字符
  COLLATE = utf8mb4_general_ci                                                    -- 使用utf8mb4_general_ci排序规则，不区分大小写
  COMMENT = '前台用户表 - 存储C端用户的基础信息，支持用户注册登录和个人资料管理'    -- 表级注释，说明表的用途
  ROW_FORMAT = DYNAMIC;                                                           -- 使用动态行格式，节省存储空间

-- ----------------------------
-- 后台管理员用户表 - 面向内部管理人员的核心数据表
-- 用途：存储后台管理员的详细信息，包含安全控制和权限管理相关字段
-- 场景：管理员登录验证、权限控制、操作审计、部门管理等
-- 业务说明：此表专门用于后台管理系统，与前台用户表完全分离，安全级别更高
-- 安全特性：包含密码错误次数、账号锁定、登录统计等安全控制机制
-- ----------------------------
CREATE TABLE `sys_admin_user`
(
    -- 主键字段：管理员的唯一标识
    `admin_id`        bigint(20)                                                    NOT NULL AUTO_INCREMENT COMMENT '管理员ID - 主键，自动递增，唯一标识每个管理员',

    -- 身份信息字段：管理员的基本身份信息
    `admin_name`      varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL COMMENT '管理员账号 - 登录用的唯一账号，不可重复',
    `real_name`       varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL COMMENT '真实姓名 - 管理员的真实姓名，用于身份确认和显示',
    `email`           varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL COMMENT '邮箱 - 必填，用于接收系统通知和找回密码',
    `phone_number`    varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NULL DEFAULT '' COMMENT '手机号码 - 11位手机号，用于双重验证和紧急联系',
    `avatar`          varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '头像地址 - 管理员头像图片的URL路径',

    -- 安全认证字段：密码和账号安全控制
    `password`        varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码 - 必填，加密存储的管理员密码',
    `status`          tinyint(1)                                                    NULL DEFAULT 1 COMMENT '帐号状态 - 0:停用(禁止登录) 1:正常(可以登录)',
    `del_flag`        char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci      NULL DEFAULT '0' COMMENT '删除标志 - 0:存在 2:已删除(软删除，保留数据用于审计)',

    -- 登录记录字段：详细的登录信息统计
    `login_ip`        varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NULL DEFAULT '' COMMENT '最后登陆IP - 记录最后一次登录的IP地址，用于安全监控',
    `login_date`      datetime                                                      NULL DEFAULT NULL COMMENT '最后登陆时间 - 记录最后一次成功登录的时间',
    `login_count`     int(11)                                                       NULL DEFAULT 0 COMMENT '登录次数 - 累计登录次数，用于统计和分析',

    -- 安全控制字段：密码安全和账号锁定机制
    `last_pwd_time`   datetime                                                      NULL DEFAULT NULL COMMENT '最后修改密码时间 - 用于强制定期修改密码策略',
    `pwd_error_count` int(11)                                                       NULL DEFAULT 0 COMMENT '密码错误次数 - 连续密码错误次数，达到阈值后锁定账号',
    `lock_time`       datetime                                                      NULL DEFAULT NULL COMMENT '账号锁定时间 - 账号被锁定的时间，用于自动解锁判断',

    -- 权限管理字段：角色和部门信息
    `role_ids`        varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '角色ID列表 - 逗号分隔的角色ID，如"1,2,3"，用于权限控制',
    `department`      varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NULL DEFAULT '' COMMENT '所属部门 - 管理员所在的部门，用于组织架构管理',
    `position`        varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NULL DEFAULT '' COMMENT '职位 - 管理员的职位信息，如"系统管理员"、"运营专员"',

    -- 审计字段：记录数据的创建和修改信息
    `create_by`       varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NULL DEFAULT '' COMMENT '创建者 - 创建此管理员账号的人员，用于审计追踪',
    `create_time`     datetime                                                      NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间 - 账号创建时间，自动设置',
    `update_by`       varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NULL DEFAULT '' COMMENT '更新者 - 最后修改此账号信息的人员',
    `update_time`     datetime                                                      NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间 - 最后修改时间，自动更新',
    `remark`          varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '备注 - 管理员的备注信息，如特殊权限说明等',

    -- 主键定义：指定主键字段
    PRIMARY KEY (`admin_id`) USING BTREE COMMENT '主键索引 - 使用B-Tree索引，确保管理员ID的唯一性',

    -- 唯一索引：确保关键字段的唯一性
    UNIQUE KEY `uk_admin_name` (`admin_name`) COMMENT '管理员账号唯一索引 - 确保账号不重复，支持快速登录验证',
    UNIQUE KEY `uk_email` (`email`) COMMENT '邮箱唯一索引 - 确保邮箱不重复，支持邮箱找回密码',

    -- 普通索引：提高常用查询的效率
    INDEX `idx_status` (`status`) COMMENT '状态索引 - 加速按状态查询管理员，如查询所有正常状态的管理员',
    INDEX `idx_department` (`department`) COMMENT '部门索引 - 加速按部门查询管理员，支持部门管理功能'

) ENGINE = InnoDB                                                                 -- 使用InnoDB存储引擎，支持事务和外键约束
  AUTO_INCREMENT = 1                                                              -- 自增起始值
  CHARACTER SET = utf8mb4                                                         -- 字符集设置
  COLLATE = utf8mb4_general_ci                                                    -- 排序规则设置
  COMMENT = '后台管理员用户表 - 存储内部管理人员信息，包含完整的安全控制和权限管理机制'
  ROW_FORMAT = DYNAMIC;                                                           -- 动态行格式

-- ----------------------------
-- 游戏基本信息表 - 游戏数据的核心主表
-- 用途：存储游戏的所有核心信息，是整个游戏系统的数据中心
-- 场景：游戏录入、游戏详情展示、游戏列表查询、游戏搜索、游戏推荐等
-- 业务说明：此表是游戏相关功能的核心，其他表如封面、截图、评论等都会关联到此表
-- 数据特点：包含游戏的基础信息、统计数据、关联关系等完整信息
-- ----------------------------
CREATE TABLE `sys_game`
(
    -- 主键字段：游戏的唯一标识
    `game_id`       bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '游戏ID - 主键，自动递增，全局唯一标识每个游戏',

    -- 基础信息字段：游戏的核心描述信息
    `name_zh`       varchar(100) NOT NULL COMMENT '游戏中文名称 - 必填，游戏的中文名称，用于显示和搜索',
    `name_en`       varchar(100)          DEFAULT NULL COMMENT '游戏英文名称 - 可选，游戏的英文名称，支持国际化',
    `release_date`  date                  DEFAULT NULL COMMENT '游戏发布日期 - 可选，游戏正式发布的日期，格式：YYYY-MM-DD',
    `description`   text                  DEFAULT NULL COMMENT '游戏介绍 - 可选，游戏的详细介绍，支持长文本存储',

    -- 评价和规格字段：游戏的质量和技术信息
    `rating`        decimal(3, 1)         DEFAULT 0.0 COMMENT '游戏评分 - 评分范围0.0-10.0，保留一位小数，如8.5分',
    `size`          varchar(50)           DEFAULT NULL COMMENT '游戏大小 - 游戏安装包大小，如"2.5GB"、"1.2TB"',
    `price`         decimal(10, 2)        DEFAULT NULL COMMENT '游戏价格 - 游戏售价，保留两位小数，如99.99，NULL表示免费',

    -- 关联字段：与其他实体的关系
    `publisher_id`  bigint(20)            DEFAULT NULL COMMENT '游戏厂商ID - 外键，关联sys_publisher表，标识游戏发行商',
    `studio_id`     bigint(20)            DEFAULT NULL COMMENT '游戏工作室ID - 外键，关联sys_studio表，标识游戏开发商',

    -- 生命周期字段：游戏的运营状态
    `shutdown_date` date                  DEFAULT NULL COMMENT '游戏停服日期 - 可选，游戏停止运营的日期，NULL表示正常运营',

    -- 统计字段：游戏的用户互动数据
    `comment_count` int(11)      NOT NULL DEFAULT 0 COMMENT '游戏评论数量 - 该游戏的总评论数，用于热度排序',
    `like_count`    int(11)      NOT NULL DEFAULT 0 COMMENT '游戏点赞数量 - 该游戏的总点赞数，用于推荐算法',

    -- 多媒体字段：游戏的展示内容
    `demo_video`    varchar(255)          DEFAULT NULL COMMENT '游戏演示视频链接 - 游戏预告片或演示视频的URL',

    -- 系统字段：数据管理和审计信息
    `create_time`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间 - 记录创建时自动设置为当前时间',
    `update_time`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间 - 记录更新时自动更新为当前时间',
    `create_by`     varchar(64)           DEFAULT '' COMMENT '创建者 - 创建此游戏记录的管理员账号',
    `update_by`     varchar(64)           DEFAULT '' COMMENT '更新者 - 最后修改此游戏记录的管理员账号',
    `status`        tinyint(1)   NOT NULL DEFAULT 1 COMMENT '状态 - 0:停用(不显示) 1:正常(可显示)，控制游戏是否对外可见',

    -- 主键定义
    PRIMARY KEY (`game_id`) COMMENT '主键索引 - 确保游戏ID唯一性，提供最快的单条记录查询',

    -- 业务索引：提高常用查询的性能
    INDEX `idx_publisher` (`publisher_id`) COMMENT '厂商索引 - 加速"查询某厂商的所有游戏"等查询',
    INDEX `idx_studio` (`studio_id`) COMMENT '工作室索引 - 加速"查询某工作室开发的所有游戏"等查询',
    INDEX `idx_rating` (`rating`) COMMENT '评分索引 - 加速按评分排序查询，支持"高分游戏推荐"功能',
    INDEX `idx_status` (`status`) COMMENT '状态索引 - 加速按状态筛选，如只查询正常状态的游戏',
    INDEX `idx_release_date` (`release_date`) COMMENT '发布日期索引 - 加速按发布时间排序，支持"最新游戏"功能'

) ENGINE = InnoDB                                                                 -- InnoDB引擎支持事务和外键
  AUTO_INCREMENT = 1                                                              -- 自增起始值
  CHARACTER SET = utf8mb4                                                         -- 支持emoji和特殊字符
  COLLATE = utf8mb4_general_ci                                                    -- 不区分大小写排序
  COMMENT ='游戏基本信息表 - 游戏系统的核心数据表，存储游戏的完整信息，是其他游戏相关表的关联中心';

-- ----------------------------
-- 游戏封面表（支持多张）
-- 用途：存储游戏的封面图片，一个游戏可以有多张封面
-- 场景：游戏详情页展示封面、游戏列表展示缩略图等
-- ----------------------------
CREATE TABLE `sys_game_cover`
(
    `cover_id`    bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '封面ID - 主键，自动递增',
    `game_id`     bigint(20)   NOT NULL COMMENT '游戏ID - 外键，关联sys_game表',
    `cover_url`   varchar(255) NOT NULL COMMENT '封面图片URL - 存储图片路径',
    `is_main`     tinyint(1)   NOT NULL DEFAULT 0 COMMENT '是否主封面（0否 1是）- 标记主要展示图',
    `sort`        int(11)      NOT NULL DEFAULT 0 COMMENT '排序 - 控制多张封面的展示顺序',
    `create_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间 - 自动填充',
    PRIMARY KEY (`cover_id`) COMMENT '主键索引 - 确保每条记录唯一标识',
    INDEX `idx_game_id` (`game_id`) COMMENT '游戏ID索引 - 加速查询特定游戏的所有封面'
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='游戏封面表 - 存储游戏封面图片，支持一个游戏多张封面，通过is_main标记主封面';

-- ----------------------------
-- 游戏截图表
-- 用途：存储游戏的截图，一个游戏可以有多张截图
-- 场景：游戏详情页展示游戏画面、游戏预览等
-- ----------------------------
CREATE TABLE `sys_game_screenshot`
(
    `screenshot_id`  bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '截图ID - 主键，自动递增',
    `game_id`        bigint(20)   NOT NULL COMMENT '游戏ID - 外键，关联sys_game表',
    `screenshot_url` varchar(255) NOT NULL COMMENT '截图URL - 存储图片路径',
    `sort`           int(11)      NOT NULL DEFAULT 0 COMMENT '排序 - 控制多张截图的展示顺序',
    `create_time`    datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间 - 自动填充',
    PRIMARY KEY (`screenshot_id`) COMMENT '主键索引 - 确保每条记录唯一标识',
    INDEX `idx_game_id` (`game_id`) COMMENT '游戏ID索引 - 加速查询特定游戏的所有截图'
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='游戏截图表 - 存储游戏截图，支持一个游戏多张截图，用于展示游戏实际画面';

-- ----------------------------
-- 字典类型表 - 系统字典的分类管理表
-- 用途：定义和管理系统中所有字典数据的分类类型
-- 场景：字典分类管理、下拉选项分组、系统配置管理等
-- 业务说明：字典系统采用两层结构，此表定义字典的类别，sys_dict_data表存储具体的字典项
-- 设计模式：这是典型的"字典表设计模式"，将可配置的枚举值从代码中分离到数据库
-- 使用示例：游戏类型(game_type)包含RPG、FPS等选项，平台类型(platform_type)包含PC、手机等选项
-- ----------------------------
CREATE TABLE `sys_dict_type`
(
    -- 主键字段：字典类型的唯一标识
    `dict_id`     bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '字典类型ID - 主键，自动递增，唯一标识每个字典类型',

    -- 显示信息字段：用于前端展示的字典类型信息
    `dict_name`   varchar(100) NOT NULL COMMENT '字典名称 - 显示名称，如"游戏类型"、"平台类型"，用于管理界面展示',
    `dict_type`   varchar(100) NOT NULL COMMENT '字典类型标识 - 唯一标识符，如"game_type"、"platform_type"，用于程序调用',

    -- 控制字段：字典类型的状态管理
    `status`      tinyint(1)   NOT NULL DEFAULT 1 COMMENT '状态 - 0:停用(该类字典不可用) 1:正常(该类字典可用)',
    `remark`      varchar(500)          DEFAULT NULL COMMENT '备注 - 字典类型的详细说明，如使用场景、注意事项等',

    -- 系统字段：数据管理和审计信息
    `create_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间 - 记录创建时自动设置',
    `update_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间 - 记录更新时自动更新',
    `create_by`   varchar(64)           DEFAULT '' COMMENT '创建者 - 创建此字典类型的管理员账号',
    `update_by`   varchar(64)           DEFAULT '' COMMENT '更新者 - 最后修改此字典类型的管理员账号',

    -- 主键定义
    PRIMARY KEY (`dict_id`) COMMENT '主键索引 - 确保字典类型ID的唯一性',

    -- 唯一索引：确保字典类型标识的唯一性
    UNIQUE KEY `uk_dict_type` (`dict_type`) COMMENT '字典类型唯一索引 - 确保dict_type字段不重复，保证程序调用的准确性'

) ENGINE = InnoDB                                                                 -- InnoDB引擎
  AUTO_INCREMENT = 1                                                              -- 自增起始值
  CHARACTER SET = utf8mb4                                                         -- 字符集
  COLLATE = utf8mb4_general_ci                                                    -- 排序规则
  COMMENT ='字典类型表 - 定义系统中所有字典的分类类型，是字典数据的分组依据，支持系统配置的灵活管理';

-- ----------------------------
-- 字典数据表 - 系统字典的具体数据存储表
-- 用途：存储所有字典类型下的具体数据项，是系统配置的核心数据表
-- 场景：前端下拉选择、单选多选组件、标签展示、筛选条件等所有需要固定选项的场景
-- 业务说明：与sys_dict_type表配合使用，通过dict_type字段关联，实现字典数据的分类管理
-- 数据特点：支持排序、默认值、状态控制等功能，满足复杂的业务需求
-- 使用示例：游戏类型下有"RPG"、"FPS"等选项，平台类型下有"PC"、"手机"等选项
-- ----------------------------
CREATE TABLE `sys_dict_data`
(
    -- 主键字段：字典数据项的唯一标识
    `dict_code`   bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '字典编码 - 主键，自动递增，唯一标识每个字典数据项',

    -- 排序字段：控制字典项的显示顺序
    `dict_sort`   int(11)      NOT NULL DEFAULT 0 COMMENT '字典排序 - 数字越小越靠前，控制同类字典项的展示顺序',

    -- 核心数据字段：字典项的显示和存储值
    `dict_label`  varchar(100) NOT NULL COMMENT '字典标签 - 显示给用户看的文本，如"角色扮演"、"个人电脑"',
    `dict_value`  varchar(100) NOT NULL COMMENT '字典键值 - 程序中实际使用的值，如"RPG"、"PC"',

    -- 分类字段：关联到字典类型表
    `dict_type`   varchar(100) NOT NULL COMMENT '字典类型 - 关联sys_dict_type表的dict_type字段，实现分类管理',

    -- 控制字段：字典项的特殊属性和状态
    `is_default`  tinyint(1)   NOT NULL DEFAULT 0 COMMENT '是否默认 - 0:否 1:是，标记该选项是否为默认选中项',
    `status`      tinyint(1)   NOT NULL DEFAULT 1 COMMENT '状态 - 0:停用(不显示) 1:正常(可显示)，控制该字典项是否可用',
    `remark`      varchar(500)          DEFAULT NULL COMMENT '备注 - 字典项的详细说明，如使用场景、特殊含义等',

    -- 系统字段：数据管理和审计信息
    `create_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间 - 记录创建时自动设置',
    `update_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间 - 记录更新时自动更新',
    `create_by`   varchar(64)           DEFAULT '' COMMENT '创建者 - 创建此字典项的管理员账号',
    `update_by`   varchar(64)           DEFAULT '' COMMENT '更新者 - 最后修改此字典项的管理员账号',

    -- 主键定义
    PRIMARY KEY (`dict_code`) COMMENT '主键索引 - 确保字典编码的唯一性',

    -- 业务索引：提高查询性能
    INDEX `idx_dict_type` (`dict_type`) COMMENT '字典类型索引 - 加速"查询某类型下所有字典项"的查询，这是最常用的查询方式',
    INDEX `idx_dict_type_sort` (`dict_type`, `dict_sort`) COMMENT '类型排序复合索引 - 加速按类型查询并按排序字段排序的操作',
    INDEX `idx_status` (`status`) COMMENT '状态索引 - 加速按状态筛选字典项'

) ENGINE = InnoDB                                                                 -- InnoDB引擎
  AUTO_INCREMENT = 1                                                              -- 自增起始值
  CHARACTER SET = utf8mb4                                                         -- 字符集
  COLLATE = utf8mb4_general_ci                                                    -- 排序规则
  COMMENT ='字典数据表 - 存储系统中所有字典的具体数据项，支持排序、默认值、状态控制等功能，是系统配置的核心数据表';

-- ----------------------------
-- 初始化字典类型数据
-- 说明：为系统预置基础的字典类型，这些是游戏系统必需的基础分类
-- 用途：定义游戏相关的各种分类维度，支持后续的字典数据录入
-- 注意：这些数据是系统的基础数据，删除后会影响相关功能的正常使用
-- ----------------------------
INSERT INTO `sys_dict_type` (`dict_name`, `dict_type`, `status`, `remark`)
VALUES
    ('游戏类型', 'game_type', 1, '游戏的类型分类，如RPG、FPS等，用于游戏分类和筛选功能'),
    ('平台类型', 'platform_type', 1, '游戏支持的平台类型，如PC、主机、手机等，用于平台筛选功能'),
    ('游戏语言', 'game_language', 1, '游戏支持的语言，如中文、英文等，用于语言筛选功能');

-- ----------------------------
-- 初始化字典数据
-- 说明：为每个字典类型预置常用的字典项，这些是游戏行业的标准分类
-- 数据来源：基于游戏行业的通用分类标准和主流平台
-- 排序说明：dict_sort字段控制显示顺序，数字越小越靠前
-- 扩展性：可以根据业务需要随时添加新的字典项
-- ----------------------------
INSERT INTO `sys_dict_data` (`dict_sort`, `dict_label`, `dict_value`, `dict_type`, `status`, `remark`)
VALUES
-- ========== 游戏类型数据 ==========
-- 按照游戏类型的流行程度和重要性排序
(1, '角色扮演', 'RPG', 'game_type', 1, '角色扮演类游戏 - Role Playing Game，如《最终幻想》系列'),
(2, '第一人称射击', 'FPS', 'game_type', 1, '第一人称射击游戏 - First Person Shooter，如《使命召唤》系列'),
(3, '动作冒险', 'Action-Adventure', 'game_type', 1, '动作冒险类游戏 - 结合动作和冒险元素，如《塞尔达传说》系列'),
(4, '策略', 'Strategy', 'game_type', 1, '策略类游戏 - 需要策略思考和规划，如《文明》系列'),
(5, '模拟', 'Simulation', 'game_type', 1, '模拟类游戏 - 模拟现实场景，如《模拟人生》系列'),
(6, '体育', 'Sports', 'game_type', 1, '体育类游戏 - 体育运动模拟，如《FIFA》系列'),
(7, '竞速', 'Racing', 'game_type', 1, '竞速类游戏 - 赛车竞速类，如《极品飞车》系列'),

-- ========== 平台类型数据 ==========
-- 按照平台的市场份额和重要性排序
(1, 'PC', 'PC', 'platform_type', 1, 'PC平台 - 个人电脑，包括Windows、Mac、Linux'),
(2, 'PlayStation', 'PS', 'platform_type', 1, 'PlayStation平台 - 索尼游戏主机，包括PS4、PS5等'),
(3, 'Xbox', 'Xbox', 'platform_type', 1, 'Xbox平台 - 微软游戏主机，包括Xbox One、Xbox Series等'),
(4, 'Nintendo Switch', 'Switch', 'platform_type', 1, 'Nintendo Switch平台 - 任天堂便携式游戏主机'),
(5, 'iOS', 'iOS', 'platform_type', 1, 'iOS平台 - 苹果手机和平板设备'),
(6, 'Android', 'Android', 'platform_type', 1, 'Android平台 - 安卓手机和平板设备'),

-- ========== 游戏语言数据 ==========
-- 按照语言的使用人数和重要性排序
(1, '简体中文', 'zh_CN', 'game_language', 1, '简体中文 - 中国大陆使用的中文'),
(2, '繁体中文', 'zh_TW', 'game_language', 1, '繁体中文 - 台湾、香港等地区使用的中文'),
(3, '英语', 'en_US', 'game_language', 1, '英语 - 国际通用语言'),
(4, '日语', 'ja_JP', 'game_language', 1, '日语 - 日本语言，游戏大国的主要语言'),
(5, '韩语', 'ko_KR', 'game_language', 1, '韩语 - 韩国语言，手游强国的主要语言');

-- ----------------------------
-- 游戏与类型关联表（多对多）
-- 用途：实现游戏和类型的多对多关系
-- 场景：一个游戏可以属于多个类型，一个类型可以包含多个游戏
-- ----------------------------
CREATE TABLE `sys_game_type_relation`
(
    `id`        bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID - 主键，自动递增',
    `game_id`   bigint(20) NOT NULL COMMENT '游戏ID - 关联sys_game表',
    `dict_code` bigint(20) NOT NULL COMMENT '字典编码 - 关联sys_dict_data表中game_type类型的数据',
    PRIMARY KEY (`id`) COMMENT '主键索引 - 确保每条记录唯一标识',
    UNIQUE KEY `uk_game_type` (`game_id`, `dict_code`) COMMENT '唯一索引 - 确保一个游戏不会重复关联同一个类型',
    INDEX `idx_game_id` (`game_id`) COMMENT '游戏ID索引 - 加速查询特定游戏的所有类型',
    INDEX `idx_dict_code` (`dict_code`) COMMENT '字典编码索引 - 加速查询特定类型的所有游戏'
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='游戏与类型关联表 - 实现游戏和类型的多对多关系，一个游戏可以有多个类型，如RPG、动作等';

-- ----------------------------
-- 游戏与平台关联表（多对多）
-- 用途：实现游戏和平台的多对多关系
-- 场景：一个游戏可以支持多个平台，一个平台可以有多个游戏
-- ----------------------------
CREATE TABLE `sys_game_platform_relation`
(
    `id`        bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID - 主键，自动递增',
    `game_id`   bigint(20) NOT NULL COMMENT '游戏ID - 关联sys_game表',
    `dict_code` bigint(20) NOT NULL COMMENT '字典编码 - 关联sys_dict_data表中platform_type类型的数据',
    PRIMARY KEY (`id`) COMMENT '主键索引 - 确保每条记录唯一标识',
    UNIQUE KEY `uk_game_platform` (`game_id`, `dict_code`) COMMENT '唯一索引 - 确保一个游戏不会重复关联同一个平台',
    INDEX `idx_game_id` (`game_id`) COMMENT '游戏ID索引 - 加速查询特定游戏支持的所有平台',
    INDEX `idx_dict_code` (`dict_code`) COMMENT '字典编码索引 - 加速查询特定平台上的所有游戏'
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='游戏与平台关联表 - 实现游戏和平台的多对多关系，记录游戏支持的平台，如PC、PS5等';

-- ----------------------------
-- 游戏与语言关联表（多对多）
-- 用途：实现游戏和语言的多对多关系
-- 场景：一个游戏可以支持多种语言，一种语言可以被多个游戏支持
-- ----------------------------
CREATE TABLE `sys_game_language_relation`
(
    `id`        bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID - 主键，自动递增',
    `game_id`   bigint(20) NOT NULL COMMENT '游戏ID - 关联sys_game表',
    `dict_code` bigint(20) NOT NULL COMMENT '字典编码 - 关联sys_dict_data表中game_language类型的数据',
    PRIMARY KEY (`id`) COMMENT '主键索引 - 确保每条记录唯一标识',
    UNIQUE KEY `uk_game_language` (`game_id`, `dict_code`) COMMENT '唯一索引 - 确保一个游戏不会重复关联同一种语言',
    INDEX `idx_game_id` (`game_id`) COMMENT '游戏ID索引 - 加速查询特定游戏支持的所有语言',
    INDEX `idx_dict_code` (`dict_code`) COMMENT '字典编码索引 - 加速查询支持特定语言的所有游戏'
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='游戏与语言关联表 - 实现游戏和语言的多对多关系，记录游戏支持的语言，如中文、英文等';

-- ----------------------------
-- 游戏厂商表
-- 用途：存储游戏厂商信息
-- 场景：厂商详情页、游戏详情页显示厂商信息等
-- ----------------------------
CREATE TABLE `sys_publisher`
(
    `publisher_id`   bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '厂商ID - 主键，自动递增',
    `publisher_name` varchar(100) NOT NULL COMMENT '厂商名称 - 如"腾讯游戏"、"网易游戏"',
    `logo_url`       varchar(255)          DEFAULT NULL COMMENT '厂商LOGO - 存储logo图片路径',
    `description`    text                  DEFAULT NULL COMMENT '厂商介绍 - 存储长文本',
    `founded_date`   date                  DEFAULT NULL COMMENT '成立日期 - 厂商成立时间',
    `website`        varchar(255)          DEFAULT NULL COMMENT '官方网站 - 厂商官网链接',
    `status`         tinyint(1)   NOT NULL DEFAULT 1 COMMENT '状态（0停用 1正常）- 控制厂商是否可见',
    `create_time`    datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间 - 自动填充',
    `update_time`    datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间 - 自动更新',
    PRIMARY KEY (`publisher_id`) COMMENT '主键索引 - 确保每条记录唯一标识'
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='游戏厂商表 - 存储游戏厂商信息，包括厂商名称、LOGO、介绍等基本信息';

-- ----------------------------
-- 用户关注表（关注其他用户）
-- 用途：实现用户之间的关注关系
-- 场景：社交功能、关注列表等
-- ----------------------------
CREATE TABLE `sys_user_follow`
(
    `follow_id`   bigint(20) NOT NULL AUTO_INCREMENT COMMENT '关注ID - 主键，自动递增',
    `user_id`     bigint(20) NOT NULL COMMENT '用户ID - 关注者，关联sys_user表',
    `follow_user_id` bigint(20) NOT NULL COMMENT '被关注用户ID - 被关注者，关联sys_user表',
    `create_time` datetime   NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '关注时间 - 自动填充',
    PRIMARY KEY (`follow_id`) COMMENT '主键索引 - 确保每条记录唯一标识',
    UNIQUE KEY `uk_user_follow` (`user_id`, `follow_user_id`) COMMENT '唯一索引 - 确保不会重复关注',
    INDEX `idx_user_id` (`user_id`) COMMENT '用户ID索引 - 加速查询用户关注的所有人',
    INDEX `idx_follow_user_id` (`follow_user_id`) COMMENT '被关注用户ID索引 - 加速查询关注特定用户的所有人'
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='用户关注表 - 记录用户之间的关注关系，实现社交功能';

-- ----------------------------
-- 用户游戏评论表
-- 用途：存储用户对游戏的评论
-- 场景：游戏详情页评论区、用户评论历史等
-- ----------------------------
CREATE TABLE `sys_user_game_comment`
(
    `comment_id`  bigint(20)  NOT NULL AUTO_INCREMENT COMMENT '评论ID - 主键，自动递增',
    `user_id`     bigint(20)  NOT NULL COMMENT '用户ID - 评论者，关联sys_user表',
    `game_id`     bigint(20)  NOT NULL COMMENT '游戏ID - 被评论的游戏，关联sys_game表',
    `content`     text        NOT NULL COMMENT '评论内容 - 用户的评论文字',
    `parent_id`   bigint(20)  DEFAULT NULL COMMENT '父评论ID - 回复的评论ID，NULL表示一级评论',
    `like_count`  int(11)     NOT NULL DEFAULT 0 COMMENT '点赞数 - 该评论获得的点赞数',
    `status`      tinyint(1)  NOT NULL DEFAULT 1 COMMENT '状态（0隐藏 1显示）- 控制评论是否可见',
    `create_time` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '评论时间 - 自动填充',
    `update_time` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间 - 自动更新',
    PRIMARY KEY (`comment_id`) COMMENT '主键索引 - 确保每条记录唯一标识',
    INDEX `idx_game_id` (`game_id`) COMMENT '游戏ID索引 - 加速查询特定游戏的所有评论',
    INDEX `idx_user_id` (`user_id`) COMMENT '用户ID索引 - 加速查询特定用户的所有评论',
    INDEX `idx_parent_id` (`parent_id`) COMMENT '父评论ID索引 - 加速查询评论的所有回复'
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='用户游戏评论表 - 存储用户对游戏的评论和回复，支持多级评论';

-- ----------------------------
-- 评论点赞表
-- 用途：记录用户对评论的点赞
-- 场景：评论点赞功能、热门评论排序等
-- ----------------------------
CREATE TABLE `sys_comment_like`
(
    `like_id`     bigint(20) NOT NULL AUTO_INCREMENT COMMENT '点赞ID - 主键，自动递增',
    `comment_id`  bigint(20) NOT NULL COMMENT '评论ID - 被点赞的评论，关联sys_user_game_comment表',
    `user_id`     bigint(20) NOT NULL COMMENT '用户ID - 点赞者，关联sys_user表',
    `create_time` datetime   NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '点赞时间 - 自动填充',
    PRIMARY KEY (`like_id`) COMMENT '主键索引 - 确保每条记录唯一标识',
    UNIQUE KEY `uk_user_comment` (`user_id`, `comment_id`) COMMENT '唯一索引 - 确保用户不会重复点赞同一评论',
    INDEX `idx_comment_id` (`comment_id`) COMMENT '评论ID索引 - 加速查询特定评论的所有点赞',
    INDEX `idx_user_id` (`user_id`) COMMENT '用户ID索引 - 加速查询用户点赞的所有评论'
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='评论点赞表 - 记录用户对评论的点赞，支持点赞和取消点赞';

-- ----------------------------
-- 系统菜单表 - 后台管理系统的菜单权限核心表
-- 用途：存储后台管理系统的所有菜单项，包括页面菜单、按钮权限等
-- 场景：后台菜单展示、权限控制、路由生成、按钮权限验证等
-- 业务说明：采用树形结构设计，支持多级菜单，通过parent_id实现父子关系
-- 权限模型：结合角色表实现RBAC权限控制，精确到按钮级别的权限管理
-- 菜单类型：支持目录(M)、菜单(C)、按钮(F)三种类型，满足不同的权限控制需求
-- ----------------------------
CREATE TABLE `sys_menu`
(
    -- 主键字段：菜单的唯一标识
    `menu_id`     bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '菜单ID - 主键，自动递增，唯一标识每个菜单项',

    -- 层级关系字段：实现树形菜单结构
    `parent_id`   bigint(20)   NOT NULL DEFAULT 0 COMMENT '父菜单ID - 0表示顶级菜单，其他值表示父菜单的menu_id，构建树形结构',
    `ancestors`   varchar(50)  NOT NULL DEFAULT '' COMMENT '祖级列表 - 所有父级菜单ID的路径，如"0,1,2"，用于快速查询所有子菜单',

    -- 显示信息字段：菜单的展示相关信息
    `menu_name`   varchar(50)  NOT NULL COMMENT '菜单名称 - 显示在界面上的菜单文字，如"用户管理"、"系统设置"',
    `order_num`   int(11)      NOT NULL DEFAULT 0 COMMENT '显示顺序 - 数字越小越靠前，控制同级菜单的排序',
    `icon`        varchar(100)          DEFAULT '#' COMMENT '菜单图标 - 菜单项的图标标识，如"user"、"setting"，#表示无图标',

    -- 路由信息字段：前端路由和组件相关
    `path`        varchar(200)          DEFAULT '' COMMENT '路由地址 - 前端路由路径，如"/system/user"，空字符串表示不是路由菜单',
    `component`   varchar(255)          DEFAULT NULL COMMENT '组件路径 - 前端组件的路径，如"system/user/index"，NULL表示不是页面组件',
    `query`       varchar(255)          DEFAULT NULL COMMENT '路由参数 - 路由跳转时携带的参数，如"userId=1&status=1"',

    -- 权限控制字段：菜单类型和权限标识
    `menu_type`   char(1)      NOT NULL DEFAULT 'M' COMMENT '菜单类型 - M:目录(只做分组) C:菜单(对应页面) F:按钮(页面内按钮权限)',
    `visible`     tinyint(1)   NOT NULL DEFAULT 1 COMMENT '菜单状态 - 0:隐藏(不显示在菜单中) 1:显示(正常显示)',
    `status`      tinyint(1)   NOT NULL DEFAULT 1 COMMENT '菜单状态 - 0:停用(禁止访问) 1:正常(可以访问)',
    `perms`       varchar(100)          DEFAULT NULL COMMENT '权限标识 - 权限字符串，如"system:user:list"，用于后端权限验证',

    -- 外部链接字段：支持外部链接菜单
    `is_frame`    tinyint(1)   NOT NULL DEFAULT 1 COMMENT '是否为外链 - 0:是外链(跳转到外部网站) 1:不是外链(内部路由)',

    -- 缓存控制字段：前端页面缓存控制
    `is_cache`    tinyint(1)   NOT NULL DEFAULT 0 COMMENT '是否缓存 - 0:不缓存(每次重新加载) 1:缓存(保持页面状态)',

    -- 系统字段：数据管理和审计信息
    `create_by`   varchar(64)           DEFAULT '' COMMENT '创建者 - 创建此菜单的管理员账号',
    `create_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间 - 菜单创建时间，自动设置',
    `update_by`   varchar(64)           DEFAULT '' COMMENT '更新者 - 最后修改此菜单的管理员账号',
    `update_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间 - 最后修改时间，自动更新',
    `remark`      varchar(500)          DEFAULT '' COMMENT '备注 - 菜单的详细说明，如使用场景、特殊权限等',

    -- 主键定义
    PRIMARY KEY (`menu_id`) COMMENT '主键索引 - 确保菜单ID的唯一性',

    -- 业务索引：提高查询性能
    INDEX `idx_parent_id` (`parent_id`) COMMENT '父菜单ID索引 - 加速查询某菜单的所有子菜单，构建菜单树时使用',
    INDEX `idx_menu_type` (`menu_type`) COMMENT '菜单类型索引 - 加速按类型查询菜单，如只查询页面菜单或按钮权限',
    INDEX `idx_status` (`status`) COMMENT '状态索引 - 加速按状态筛选菜单，如只查询正常状态的菜单',
    INDEX `idx_visible` (`visible`) COMMENT '可见性索引 - 加速按可见性筛选菜单，构建前端菜单时使用'

) ENGINE = InnoDB                                                                 -- InnoDB引擎支持事务
  AUTO_INCREMENT = 1                                                              -- 自增起始值
  CHARACTER SET = utf8mb4                                                         -- 字符集
  COLLATE = utf8mb4_general_ci                                                    -- 排序规则
  COMMENT ='系统菜单表 - 存储后台管理系统的菜单权限信息，支持树形结构和RBAC权限控制，精确到按钮级别的权限管理';

-- ----------------------------
-- 系统角色表 - 后台权限管理的核心表
-- 用途：定义系统中的各种角色，如超级管理员、普通管理员、运营人员等
-- 场景：角色管理、权限分配、用户角色绑定、权限验证等
-- 业务说明：RBAC权限模型的核心组件，通过角色来组织权限，用户通过角色获得权限
-- 权限继承：支持角色层级和权限继承，可以设置角色的数据权限范围
-- 扩展性：预留了数据权限字段，支持行级数据权限控制
-- ----------------------------
CREATE TABLE `sys_role`
(
    -- 主键字段：角色的唯一标识
    `role_id`             bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '角色ID - 主键，自动递增，唯一标识每个角色',

    -- 基础信息字段：角色的基本信息
    `role_name`           varchar(30)  NOT NULL COMMENT '角色名称 - 角色的显示名称，如"超级管理员"、"内容管理员"',
    `role_key`            varchar(100) NOT NULL COMMENT '角色权限字符串 - 角色的唯一标识符，如"admin"、"editor"，用于程序中的权限判断',
    `role_sort`           int(11)      NOT NULL DEFAULT 0 COMMENT '显示顺序 - 数字越小越靠前，控制角色列表的排序',

    -- 权限范围字段：数据权限控制
    `data_scope`          tinyint(1)   NOT NULL DEFAULT 1 COMMENT '数据范围 - 1:全部数据权限 2:自定数据权限 3:本部门数据权限 4:本部门及以下数据权限 5:仅本人数据权限',
    `menu_check_strictly` tinyint(1)   NOT NULL DEFAULT 1 COMMENT '菜单树选择项是否关联显示 - 0:父子不互相关联显示 1:父子互相关联显示',
    `dept_check_strictly` tinyint(1)   NOT NULL DEFAULT 1 COMMENT '部门树选择项是否关联显示 - 0:父子不互相关联显示 1:父子互相关联显示',

    -- 状态控制字段：角色的启用状态
    `status`              tinyint(1)   NOT NULL DEFAULT 1 COMMENT '角色状态 - 0:停用(该角色不可用) 1:正常(该角色可用)',
    `del_flag`            char(1)               DEFAULT '0' COMMENT '删除标志 - 0:存在 2:已删除(软删除，保留数据用于审计)',

    -- 系统字段：数据管理和审计信息
    `create_by`           varchar(64)           DEFAULT '' COMMENT '创建者 - 创建此角色的管理员账号',
    `create_time`         datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间 - 角色创建时间，自动设置',
    `update_by`           varchar(64)           DEFAULT '' COMMENT '更新者 - 最后修改此角色的管理员账号',
    `update_time`         datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间 - 最后修改时间，自动更新',
    `remark`              varchar(500)          DEFAULT NULL COMMENT '备注 - 角色的详细说明，如职责范围、特殊权限等',

    -- 主键定义
    PRIMARY KEY (`role_id`) COMMENT '主键索引 - 确保角色ID的唯一性',

    -- 唯一索引：确保角色标识的唯一性
    UNIQUE KEY `uk_role_key` (`role_key`) COMMENT '角色权限字符串唯一索引 - 确保role_key不重复，保证权限判断的准确性',

    -- 业务索引：提高查询性能
    INDEX `idx_status` (`status`) COMMENT '状态索引 - 加速按状态查询角色，如只查询正常状态的角色'

) ENGINE = InnoDB                                                                 -- InnoDB引擎
  AUTO_INCREMENT = 1                                                              -- 自增起始值
  CHARACTER SET = utf8mb4                                                         -- 字符集
  COLLATE = utf8mb4_general_ci                                                    -- 排序规则
  COMMENT ='系统角色表 - 定义系统中的各种角色，是RBAC权限模型的核心组件，支持数据权限和菜单权限的精细化控制';

-- ----------------------------
-- 角色和菜单关联表 - 实现角色与菜单权限的多对多关系
-- 用途：建立角色和菜单之间的关联关系，实现基于角色的权限控制
-- 场景：角色权限分配、权限验证、菜单权限查询等
-- 业务说明：RBAC权限模型的关键组件，通过此表实现"角色拥有哪些菜单权限"的映射
-- 权限粒度：支持页面级和按钮级权限控制，精确到每个操作按钮
-- 数据特点：多对多关系，一个角色可以有多个菜单权限，一个菜单权限可以分配给多个角色
-- ----------------------------
CREATE TABLE `sys_role_menu`
(
    -- 主键字段：关联关系的唯一标识
    `role_id` bigint(20) NOT NULL COMMENT '角色ID - 关联sys_role表，标识哪个角色',
    `menu_id` bigint(20) NOT NULL COMMENT '菜单ID - 关联sys_menu表，标识哪个菜单权限',

    -- 联合主键定义：确保同一角色不会重复分配同一菜单权限
    PRIMARY KEY (`role_id`, `menu_id`) COMMENT '联合主键 - 确保角色和菜单的关联关系唯一，避免重复分配权限',

    -- 外键索引：提高关联查询性能
    INDEX `idx_role_id` (`role_id`) COMMENT '角色ID索引 - 加速"查询某角色拥有的所有菜单权限"的查询',
    INDEX `idx_menu_id` (`menu_id`) COMMENT '菜单ID索引 - 加速"查询某菜单权限分配给了哪些角色"的查询'

) ENGINE = InnoDB                                                                 -- InnoDB引擎支持外键约束
  CHARACTER SET = utf8mb4                                                         -- 字符集
  COLLATE = utf8mb4_general_ci                                                    -- 排序规则
  COMMENT ='角色和菜单关联表 - 实现角色与菜单权限的多对多关系，是RBAC权限控制的核心映射表';

-- ----------------------------
-- 管理员和角色关联表 - 实现管理员与角色的多对多关系
-- 用途：建立管理员和角色之间的关联关系，实现用户权限的分配
-- 场景：用户角色分配、权限验证、用户权限查询等
-- 业务说明：RBAC权限模型的用户层，通过此表实现"用户拥有哪些角色"的映射
-- 权限继承：管理员通过角色间接获得菜单权限，支持一个用户拥有多个角色
-- 灵活性：支持动态角色分配，可以随时调整用户的角色和权限
-- ----------------------------
CREATE TABLE `sys_admin_user_role`
(
    -- 主键字段：关联关系的唯一标识
    `admin_id` bigint(20) NOT NULL COMMENT '管理员ID - 关联sys_admin_user表，标识哪个管理员',
    `role_id`  bigint(20) NOT NULL COMMENT '角色ID - 关联sys_role表，标识哪个角色',

    -- 联合主键定义：确保同一管理员不会重复分配同一角色
    PRIMARY KEY (`admin_id`, `role_id`) COMMENT '联合主键 - 确保管理员和角色的关联关系唯一，避免重复分配角色',

    -- 外键索引：提高关联查询性能
    INDEX `idx_admin_id` (`admin_id`) COMMENT '管理员ID索引 - 加速"查询某管理员拥有的所有角色"的查询',
    INDEX `idx_role_id` (`role_id`) COMMENT '角色ID索引 - 加速"查询某角色分配给了哪些管理员"的查询'

) ENGINE = InnoDB                                                                 -- InnoDB引擎支持外键约束
  CHARACTER SET = utf8mb4                                                         -- 字符集
  COLLATE = utf8mb4_general_ci                                                    -- 排序规则
  COMMENT ='管理员和角色关联表 - 实现管理员与角色的多对多关系，通过角色为管理员分配权限';

-- ----------------------------
-- 初始化系统菜单数据
-- 说明：为后台管理系统预置基础的菜单结构，这些是系统必需的基础菜单
-- 用途：定义后台管理系统的菜单树结构，包括系统管理、游戏管理等模块
-- 菜单层级：采用三级菜单结构 - 一级目录、二级菜单、三级按钮权限
-- 权限标识：perms字段采用"模块:功能:操作"的格式，如"system:user:list"
-- 注意：这些数据是系统的基础数据，删除后会影响后台管理功能的正常使用
-- ----------------------------
INSERT INTO `sys_menu` (`menu_id`, `parent_id`, `ancestors`, `menu_name`, `order_num`, `path`, `component`, `menu_type`, `visible`, `status`, `perms`, `icon`, `create_by`, `remark`)
VALUES
-- ========== 一级目录菜单 ==========
(1, 0, '0', '系统管理', 1, 'system', NULL, 'M', 1, 1, '', 'system', 'admin', '系统管理目录'),
(2, 0, '0', '游戏管理', 2, 'game', NULL, 'M', 1, 1, '', 'game', 'admin', '游戏管理目录'),

-- ========== 系统管理二级菜单 ==========
(100, 1, '0,1', '用户管理', 1, 'user', 'system/user/index', 'C', 1, 1, 'system:user:list', 'user', 'admin', '管理员用户管理菜单'),
(101, 1, '0,1', '角色管理', 2, 'role', 'system/role/index', 'C', 1, 1, 'system:role:list', 'peoples', 'admin', '角色管理菜单'),
(102, 1, '0,1', '菜单管理', 3, 'menu', 'system/menu/index', 'C', 1, 1, 'system:menu:list', 'tree-table', 'admin', '菜单管理菜单'),
(103, 1, '0,1', '字典管理', 4, 'dict', 'system/dict/index', 'C', 1, 1, 'system:dict:list', 'dict', 'admin', '字典管理菜单'),

-- ========== 游戏管理二级菜单 ==========
(200, 2, '0,2', '游戏列表', 1, 'list', 'game/list/index', 'C', 1, 1, 'game:list:list', 'list', 'admin', '游戏列表管理菜单'),
(201, 2, '0,2', '厂商管理', 2, 'publisher', 'game/publisher/index', 'C', 1, 1, 'game:publisher:list', 'company', 'admin', '游戏厂商管理菜单'),

-- ========== 用户管理按钮权限 ==========
(1000, 100, '0,1,100', '用户查询', 1, '', '', 'F', 1, 1, 'system:user:query', '#', 'admin', ''),
(1001, 100, '0,1,100', '用户新增', 2, '', '', 'F', 1, 1, 'system:user:add', '#', 'admin', ''),
(1002, 100, '0,1,100', '用户修改', 3, '', '', 'F', 1, 1, 'system:user:edit', '#', 'admin', ''),
(1003, 100, '0,1,100', '用户删除', 4, '', '', 'F', 1, 1, 'system:user:remove', '#', 'admin', ''),

-- ========== 角色管理按钮权限 ==========
(1010, 101, '0,1,101', '角色查询', 1, '', '', 'F', 1, 1, 'system:role:query', '#', 'admin', ''),
(1011, 101, '0,1,101', '角色新增', 2, '', '', 'F', 1, 1, 'system:role:add', '#', 'admin', ''),
(1012, 101, '0,1,101', '角色修改', 3, '', '', 'F', 1, 1, 'system:role:edit', '#', 'admin', ''),
(1013, 101, '0,1,101', '角色删除', 4, '', '', 'F', 1, 1, 'system:role:remove', '#', 'admin', ''),

-- ========== 菜单管理按钮权限 ==========
(1020, 102, '0,1,102', '菜单查询', 1, '', '', 'F', 1, 1, 'system:menu:query', '#', 'admin', ''),
(1021, 102, '0,1,102', '菜单新增', 2, '', '', 'F', 1, 1, 'system:menu:add', '#', 'admin', ''),
(1022, 102, '0,1,102', '菜单修改', 3, '', '', 'F', 1, 1, 'system:menu:edit', '#', 'admin', ''),
(1023, 102, '0,1,102', '菜单删除', 4, '', '', 'F', 1, 1, 'system:menu:remove', '#', 'admin', '');

-- ----------------------------
-- 初始化系统角色数据
-- 说明：为系统预置基础的角色，这些是系统运行必需的基础角色
-- 角色层级：超级管理员拥有所有权限，普通管理员拥有部分权限
-- 权限范围：通过data_scope字段控制数据权限范围
-- 注意：超级管理员角色不可删除，是系统的最高权限角色
-- ----------------------------
INSERT INTO `sys_role` (`role_id`, `role_name`, `role_key`, `role_sort`, `data_scope`, `status`, `create_by`, `remark`)
VALUES
(1, '超级管理员', 'admin', 1, 1, 1, 'admin', '超级管理员角色，拥有系统所有权限，不可删除'),
(2, '普通管理员', 'common', 2, 2, 1, 'admin', '普通管理员角色，拥有基础的管理权限');

-- ----------------------------
-- 初始化角色菜单关联数据
-- 说明：为预置角色分配对应的菜单权限，建立角色与权限的映射关系
-- 权限分配：超级管理员拥有所有菜单权限，普通管理员拥有基础的查询和管理权限
-- 权限粒度：精确到按钮级别，可以控制每个操作按钮的显示和访问权限
-- 扩展性：新增菜单后，需要在此处为对应角色分配权限
-- ----------------------------
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`)
VALUES
-- ========== 超级管理员权限(拥有所有权限) ==========
-- 一级目录权限
(1, 1), (1, 2),
-- 系统管理菜单权限
(1, 100), (1, 101), (1, 102), (1, 103),
-- 游戏管理菜单权限
(1, 200), (1, 201),
-- 用户管理按钮权限
(1, 1000), (1, 1001), (1, 1002), (1, 1003),
-- 角色管理按钮权限
(1, 1010), (1, 1011), (1, 1012), (1, 1013),
-- 菜单管理按钮权限
(1, 1020), (1, 1021), (1, 1022), (1, 1023),

-- ========== 普通管理员权限(基础权限) ==========
-- 游戏管理目录和菜单权限
(2, 2), (2, 200), (2, 201),
-- 系统管理目录权限(仅查看)
(2, 1), (2, 103),
-- 基础查询权限
(2, 1000), (2, 1010), (2, 1020);

-- ----------------------------
-- 初始化管理员角色关联数据
-- 说明：为默认管理员账号分配超级管理员角色
-- 权限继承：通过角色关联，管理员账号继承角色的所有菜单权限
-- 多角色支持：一个管理员可以拥有多个角色，权限取并集
-- ----------------------------
INSERT INTO `sys_admin_user_role` (`admin_id`, `role_id`)
VALUES
(1, 1); -- 为admin账号分配超级管理员角色