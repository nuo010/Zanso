/*
 最新版 Zanso 数据库结构
 业务层级：展册（collection） -> 分类（category） -> 资源（resource）
 说明：
 1. 本文件是全量建表脚本，适合新库初始化。
 2. 如果老库已经在线运行，请先备份，再按迁移脚本方式执行 rename/alter，不要直接 drop。
 3. 已包含默认角色、默认管理员用户和用户角色关联数据。
 默认管理员：
    login_name = admin
    password   = Admin@123456
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- 展册表：记录用户创建的展册
-- ----------------------------
DROP TABLE IF EXISTS `tbl_collection`;
CREATE TABLE `tbl_collection` (
  `id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键 ID',
  `user_id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '所属用户 ID',
  `name` varchar(160) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '展册名称',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '展册描述',
  `cover_url` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '封面地址',
  `visible` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否可见：1 可见，0 不可见',
  `status` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'draft' COMMENT '状态：draft 草稿，active 启用',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_tbl_collection_user_id` (`user_id`),
  KEY `idx_tbl_collection_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='展册表';

-- ----------------------------
-- 分类表：展册下面的分类节点
-- ----------------------------
DROP TABLE IF EXISTS `tbl_category`;
CREATE TABLE `tbl_category` (
  `id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键 ID',
  `user_id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '所属用户 ID',
  `collection_id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '所属展册 ID',
  `name` varchar(160) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类名称',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '分类描述',
  `cover_url` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '分类封面地址',
  `visible` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否可见：1 可见，0 不可见',
  `status` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'draft' COMMENT '状态：draft 草稿，active 启用',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_tbl_category_user_id` (`user_id`),
  KEY `idx_tbl_category_collection_id` (`collection_id`),
  KEY `idx_tbl_category_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分类表';

-- ----------------------------
-- 展册资源关系表：记录资源挂在展册还是挂在分类下
-- ----------------------------
DROP TABLE IF EXISTS `tbl_collection_resource_relation`;
CREATE TABLE `tbl_collection_resource_relation` (
  `id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键 ID',
  `user_id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '所属用户 ID',
  `collection_id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '所属展册 ID',
  `category_id` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '所属分类 ID，为空表示直接挂在展册下',
  `resource_id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '资源 ID',
  `resource_type` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '资源类型：image/video',
  `file_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '文件名',
  `file_size` bigint(20) DEFAULT NULL COMMENT '文件大小，单位字节',
  `mime_type` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'MIME 类型',
  `url` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '访问地址',
  `poster_url` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '视频封面访问地址',
  `sort` bigint(20) DEFAULT NULL COMMENT '排序值，越小越靠前',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_tbl_collection_resource_relation_user_id` (`user_id`),
  KEY `idx_tbl_collection_resource_relation_collection_id` (`collection_id`),
  KEY `idx_tbl_collection_resource_relation_category_id` (`category_id`),
  KEY `idx_tbl_collection_resource_relation_resource_id` (`resource_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='展册资源关系表';

-- ----------------------------
-- 资源表：文件基础信息
-- ----------------------------
DROP TABLE IF EXISTS `tbl_resource`;
CREATE TABLE `tbl_resource` (
  `id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键 ID',
  `user_id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '所属用户 ID',
  `resource_type` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '资源类型：image/video',
  `file_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '文件名',
  `file_ext` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '文件扩展名',
  `file_size` bigint(20) DEFAULT NULL COMMENT '文件大小，单位字节',
  `mime_type` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'MIME 类型',
  `storage_path` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '存储相对路径',
  `url` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '访问地址',
  `poster_url` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '视频封面访问地址',
  `status` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active' COMMENT '状态：active 启用',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_tbl_resource_user_id` (`user_id`),
  KEY `idx_tbl_resource_type` (`resource_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='资源表';

-- ----------------------------
-- 角色表
-- ----------------------------
DROP TABLE IF EXISTS `tbl_role`;
CREATE TABLE `tbl_role` (
  `id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键 ID',
  `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色名称',
  `code` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色编码',
  `description` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '角色描述',
  `status` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active' COMMENT '状态',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tbl_role_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- ----------------------------
-- 分享链接表：支持分享整个展册或单个分类
-- ----------------------------
DROP TABLE IF EXISTS `tbl_share_link`;
CREATE TABLE `tbl_share_link` (
  `id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键 ID',
  `user_id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '所属用户 ID',
  `collection_id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '展册 ID',
  `category_id` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '分类 ID，分享整个展册时为空',
  `target_type` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'collection' COMMENT '分享目标类型：collection/category',
  `share_code` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分享码',
  `title` varchar(160) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '分享标题',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '分享描述',
  `status` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active' COMMENT '状态',
  `view_count` bigint(20) DEFAULT NULL COMMENT '浏览次数',
  `expires_at` datetime(3) DEFAULT NULL COMMENT '过期时间，为空表示不过期',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tbl_share_link_share_code` (`share_code`),
  KEY `idx_tbl_share_link_user_id` (`user_id`),
  KEY `idx_tbl_share_link_collection_id` (`collection_id`),
  KEY `idx_tbl_share_link_category_id` (`category_id`),
  KEY `idx_tbl_share_link_target_type` (`target_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分享链接表';

-- ----------------------------
-- 分享访问日志表
-- ----------------------------
DROP TABLE IF EXISTS `tbl_share_view_log`;
CREATE TABLE `tbl_share_view_log` (
  `id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键 ID',
  `share_link_id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分享链接 ID',
  `collection_id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '展册 ID',
  `category_id` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '分类 ID',
  `target_type` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分享目标类型：collection/category',
  `user_id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '所属用户 ID',
  `viewer_ip` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '访问者 IP',
  `user_agent` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '访问者 UA',
  `referer` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '来源页面',
  `created_at` datetime(3) DEFAULT NULL COMMENT '访问时间',
  PRIMARY KEY (`id`),
  KEY `idx_tbl_share_view_log_share_link_id` (`share_link_id`),
  KEY `idx_tbl_share_view_log_collection_id` (`collection_id`),
  KEY `idx_tbl_share_view_log_category_id` (`category_id`),
  KEY `idx_tbl_share_view_log_user_id` (`user_id`),
  KEY `idx_tbl_share_view_log_created_at` (`created_at`),
  KEY `idx_tbl_share_view_log_target_type` (`target_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分享访问日志表';

-- ----------------------------
-- 公告表：系统管理员发布的平台公告
-- ----------------------------
DROP TABLE IF EXISTS `tbl_announcement`;
CREATE TABLE `tbl_announcement` (
  `id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键 ID',
  `title` varchar(160) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '公告标题',
  `content` text COLLATE utf8mb4_unicode_ci COMMENT '公告内容',
  `status` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'draft' COMMENT '状态：draft 草稿，active 发布',
  `created_by` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '发布人用户 ID',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_tbl_announcement_status` (`status`),
  KEY `idx_tbl_announcement_created_by` (`created_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='公告表';

-- ----------------------------
-- 用户表
-- ----------------------------
DROP TABLE IF EXISTS `tbl_user`;
CREATE TABLE `tbl_user` (
  `id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键 ID',
  `name` varchar(120) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名称',
  `login_name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '登录账号',
  `password_hash` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码哈希',
  `contact_name` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '联系人姓名',
  `contact_phone` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '联系人电话',
  `status` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active' COMMENT '状态',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tbl_user_login_name` (`login_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ----------------------------
-- 用户角色关系表
-- ----------------------------
DROP TABLE IF EXISTS `tbl_user_role`;
CREATE TABLE `tbl_user_role` (
  `id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键 ID',
  `user_id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户 ID',
  `role_id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色 ID',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tbl_user_role` (`user_id`,`role_id`),
  KEY `idx_tbl_user_role_role_id` (`role_id`),
  KEY `idx_tbl_user_role_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关系表';

-- ----------------------------
-- 初始化角色数据
-- ----------------------------
INSERT INTO `tbl_role` (`id`, `name`, `code`, `description`, `status`, `created_at`, `updated_at`) VALUES
('role_admin_00000000000000000001', '管理员', 'admin', '平台管理员角色', 'active', NOW(3), NOW(3)),
('role_user_000000000000000000001', '普通用户', 'user', '平台普通用户角色', 'active', NOW(3), NOW(3));

-- ----------------------------
-- 初始化默认管理员用户
-- 默认账号：admin
-- 默认密码：Admin@123456
-- ----------------------------
INSERT INTO `tbl_user` (`id`, `name`, `login_name`, `password_hash`, `contact_name`, `contact_phone`, `status`, `created_at`, `updated_at`) VALUES
('user_admin_000000000000000001', '系统管理员', 'admin', '$2a$10$QaOyM5YpLEOxvnXduE54duwkOrEAyt1gU9YJUdnKm5maFtTHw3ORu', '系统管理员', '13800000000', 'active', NOW(3), NOW(3));

-- ----------------------------
-- 初始化默认管理员角色关联
-- ----------------------------
INSERT INTO `tbl_user_role` (`id`, `user_id`, `role_id`, `created_at`) VALUES
('user_role_admin_0000000000001', 'user_admin_000000000000000001', 'role_admin_00000000000000000001', NOW(3));

-- ----------------------------
-- 初始化公告数据
-- ----------------------------
INSERT INTO `tbl_announcement` (`id`, `title`, `content`, `status`, `created_by`, `created_at`, `updated_at`) VALUES
('announcement_welcome_0000000001', '欢迎使用 Zanso 资源分享平台', '系统公告将展示在首页下方，管理员可以在公告管理中发布、隐藏或删除公告。', 'active', 'user_admin_000000000000000001', NOW(3), NOW(3));

SET FOREIGN_KEY_CHECKS = 1;
