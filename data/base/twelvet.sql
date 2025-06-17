-- ----------------------------
-- Table structure for sys_user
-- ----------------------------
DROP TABLE IF EXISTS `sys_user`;
CREATE TABLE `sys_user`
(
    `user_id`      bigint(20)                                                    NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `user_name`    varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL COMMENT '用户账号',
    `nick_name`    varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL COMMENT '用户昵称',
    `user_type`    varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci   NULL DEFAULT '00' COMMENT '用户类型（00系统用户）',
    `email`        varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NULL DEFAULT '' COMMENT '用户邮箱',
    `phone_number` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NULL DEFAULT '' COMMENT '手机号码',
    `sex`          tinyint(1)                                                    NULL DEFAULT NULL COMMENT '用户性别（0男 1女 2未知）',
    `avatar`       varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '头像地址',
    `password`     varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '密码',
    `status`       tinyint(1)                                                    NULL DEFAULT NULL COMMENT '帐号状态（0正常 1停用）',
    `del_flag`     char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci      NULL DEFAULT '0' COMMENT '删除标志（0代表存在 2代表删除）',
    `login_ip`     varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NULL DEFAULT '' COMMENT '最后登陆IP',
    `login_date`   datetime                                                      NULL DEFAULT NULL COMMENT '最后登陆时间',
    `create_by`    varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NULL DEFAULT '' COMMENT '创建者',
    `create_time`  datetime                                                      NULL DEFAULT NULL COMMENT '创建时间',
    `update_by`    varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NULL DEFAULT '' COMMENT '更新者',
    `update_time`  datetime                                                      NULL DEFAULT NULL COMMENT '更新时间',
    `remark`       varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '备注',
    PRIMARY KEY (`user_id`) USING BTREE
) ENGINE = InnoDB
  AUTO_INCREMENT = 6
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT = '用户信息表'
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- 游戏基本信息表
-- 用途：存储游戏的核心信息
-- 场景：创建新游戏、查询游戏详情、游戏列表展示等
-- ----------------------------
CREATE TABLE `sys_game`
(
    `game_id`       bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '游戏ID - 主键，自动递增',
    `name_zh`       varchar(100) NOT NULL COMMENT '游戏中文名称 - 必填',
    `name_en`       varchar(100)          DEFAULT NULL COMMENT '游戏英文名称 - 可选',
    `release_date`  date                  DEFAULT NULL COMMENT '游戏发布日期 - 可选',
    `description`   text                  DEFAULT NULL COMMENT '游戏介绍 - 可选，存储长文本',
    `rating`        decimal(3, 1)         DEFAULT 0.0 COMMENT '游戏评分（最多一位小数）- 默认0分',
    `size`          varchar(50)           DEFAULT NULL COMMENT '游戏大小 - 可选，如"2.5GB"',
    `price`         decimal(10, 2)        DEFAULT NULL COMMENT '游戏价格 - 可选，保留两位小数',
    `publisher_id`  bigint(20)            DEFAULT NULL COMMENT '游戏厂商ID - 外键，关联sys_publisher表',
    `studio_id`     bigint(20)            DEFAULT NULL COMMENT '游戏工作室ID - 外键，关联sys_studio表',
    `shutdown_date` date                  DEFAULT NULL COMMENT '游戏停服日期 - 可选',
    `comment_count` int(11)      NOT NULL DEFAULT 0 COMMENT '游戏评论数量 - 默认0',
    `like_count`    int(11)      NOT NULL DEFAULT 0 COMMENT '游戏点赞数量 - 默认0',
    `demo_video`    varchar(255)          DEFAULT NULL COMMENT '游戏演示视频链接 - 可选',
    `create_time`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间 - 自动填充当前时间',
    `update_time`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间 - 自动更新为当前时间',
    `create_by`     varchar(64)           DEFAULT '' COMMENT '创建者 - 记录创建人',
    `update_by`     varchar(64)           DEFAULT '' COMMENT '更新者 - 记录最后修改人',
    `status`        tinyint(1)   NOT NULL DEFAULT 1 COMMENT '状态（0停用 1正常）- 控制游戏是否可见',
    PRIMARY KEY (`game_id`) COMMENT '主键索引 - 确保每条记录唯一标识，提高查询效率',
    INDEX `idx_publisher` (`publisher_id`) COMMENT '厂商索引 - 加速按厂商查询游戏',
    INDEX `idx_studio` (`studio_id`) COMMENT '工作室索引 - 加速按工作室查询游戏',
    INDEX `idx_rating` (`rating`) COMMENT '评分索引 - 加速按评分排序查询'
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='游戏基本信息表 - 存储游戏核心数据，作为其他游戏相关表的关联中心';

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
-- 字典类型表
-- 用途：存储系统中所有字典的类型定义
-- 场景：管理不同类型的字典数据，如游戏类型、平台类型等
-- 说明：字典表通常分为类型表和数据表两部分，类型表定义字典类别，数据表存储具体数据
-- ----------------------------
CREATE TABLE `sys_dict_type`
(
    `dict_id`     bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '字典类型ID - 主键，自动递增',
    `dict_name`   varchar(100) NOT NULL COMMENT '字典名称 - 如"游戏类型"、"平台类型"',
    `dict_type`   varchar(100) NOT NULL COMMENT '字典类型 - 唯一标识，如"game_type"、"platform_type"',
    `status`      tinyint(1)   NOT NULL DEFAULT 1 COMMENT '状态（0停用 1正常）- 控制该类字典是否可用',
    `remark`      varchar(500)          DEFAULT NULL COMMENT '备注 - 字典类型的说明',
    `create_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间 - 自动填充',
    `update_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间 - 自动更新',
    `create_by`   varchar(64)           DEFAULT '' COMMENT '创建者 - 记录创建人',
    `update_by`   varchar(64)           DEFAULT '' COMMENT '更新者 - 记录最后修改人',
    PRIMARY KEY (`dict_id`) COMMENT '主键索引 - 确保每条记录唯一标识',
    UNIQUE KEY `uk_dict_type` (`dict_type`) COMMENT '唯一索引 - 确保字典类型不重复'
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='字典类型表 - 定义系统中所有字典的类型，如游戏类型、平台类型等，是字典数据表的分类依据';

-- ----------------------------
-- 字典数据表
-- 用途：存储所有字典的具体数据项
-- 场景：前端下拉选择、标签展示等需要固定选项的场景
-- 说明：通过dict_type关联到字典类型表，实现不同类型字典数据的分组管理
-- ----------------------------
CREATE TABLE `sys_dict_data`
(
    `dict_code`   bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '字典编码 - 主键，自动递增',
    `dict_sort`   int(11)      NOT NULL DEFAULT 0 COMMENT '字典排序 - 控制同类字典的展示顺序',
    `dict_label`  varchar(100) NOT NULL COMMENT '字典标签 - 展示值，如"角色扮演"、"PC"',
    `dict_value`  varchar(100) NOT NULL COMMENT '字典键值 - 实际存储值，如"RPG"、"PC"',
    `dict_type`   varchar(100) NOT NULL COMMENT '字典类型 - 关联sys_dict_type表的dict_type',
    `is_default`  tinyint(1)   NOT NULL DEFAULT 0 COMMENT '是否默认（0否 1是）- 标记默认选中项',
    `status`      tinyint(1)   NOT NULL DEFAULT 1 COMMENT '状态（0停用 1正常）- 控制该字典项是否可用',
    `remark`      varchar(500)          DEFAULT NULL COMMENT '备注 - 字典项的说明',
    `create_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间 - 自动填充',
    `update_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间 - 自动更新',
    `create_by`   varchar(64)           DEFAULT '' COMMENT '创建者 - 记录创建人',
    `update_by`   varchar(64)           DEFAULT '' COMMENT '更新者 - 记录最后修改人',
    PRIMARY KEY (`dict_code`) COMMENT '主键索引 - 确保每条记录唯一标识',
    INDEX `idx_dict_type` (`dict_type`) COMMENT '字典类型索引 - 加速查询特定类型的所有字典项'
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='字典数据表 - 存储所有字典的具体数据项，如游戏类型中的"RPG"、"FPS"等，通过dict_type字段关联到字典类型';

-- ----------------------------
-- 初始化字典类型数据
-- ----------------------------
INSERT INTO `sys_dict_type` (`dict_name`, `dict_type`, `status`, `remark`)
VALUES ('游戏类型', 'game_type', 1, '游戏的类型分类，如RPG、FPS等'),
       ('平台类型', 'platform_type', 1, '游戏支持的平台类型，如PC、主机、手机等'),
       ('游戏语言', 'game_language', 1, '游戏支持的语言，如中文、英文等');

-- ----------------------------
-- 初始化字典数据
-- ----------------------------
INSERT INTO `sys_dict_data` (`dict_sort`, `dict_label`, `dict_value`, `dict_type`, `status`, `remark`)
VALUES
-- 游戏类型
(1, '角色扮演', 'RPG', 'game_type', 1, '角色扮演类游戏'),
(2, '第一人称射击', 'FPS', 'game_type', 1, '第一人称射击游戏'),
(3, '动作冒险', 'Action-Adventure', 'game_type', 1, '动作冒险类游戏'),
(4, '策略', 'Strategy', 'game_type', 1, '策略类游戏'),
(5, '模拟', 'Simulation', 'game_type', 1, '模拟类游戏'),
(6, '体育', 'Sports', 'game_type', 1, '体育类游戏'),
(7, '竞速', 'Racing', 'game_type', 1, '竞速类游戏'),
-- 平台类型
(1, 'PC', 'PC', 'platform_type', 1, 'PC平台'),
(2, 'PlayStation', 'PS', 'platform_type', 1, 'PlayStation平台'),
(3, 'Xbox', 'Xbox', 'platform_type', 1, 'Xbox平台'),
(4, 'Nintendo Switch', 'Switch', 'platform_type', 1, 'Nintendo Switch平台'),
(5, 'iOS', 'iOS', 'platform_type', 1, 'iOS平台'),
(6, 'Android', 'Android', 'platform_type', 1, 'Android平台'),
-- 游戏语言
(1, '简体中文', 'zh_CN', 'game_language', 1, '简体中文'),
(2, '繁体中文', 'zh_TW', 'game_language', 1, '繁体中文'),
(3, '英语', 'en_US', 'game_language', 1, '英语'),
(4, '日语', 'ja_JP', 'game_language', 1, '日语'),
(5, '韩语', 'ko_KR', 'game_language', 1, '韩语');

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