-- 用户表
CREATE TABLE IF NOT EXISTS `users` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `app_id` INT UNSIGNED NOT NULL COMMENT 'APP ID',
  `open_id` VARCHAR(255) NOT NULL COMMENT '用户唯一标识',
  `nickname` VARCHAR(255) DEFAULT NULL COMMENT '昵称',
  `avatar` VARCHAR(500) DEFAULT NULL COMMENT '头像URL',
  `phone` VARCHAR(20) DEFAULT NULL COMMENT '手机号',
  `email` VARCHAR(255) DEFAULT NULL COMMENT '邮箱',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 1-正常, 0-禁用',
  `last_login_at` DATETIME DEFAULT NULL COMMENT '最后登录时间',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  INDEX `idx_app_id` (`app_id`),
  INDEX `idx_open_id` (`open_id`),
  INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 消息表
CREATE TABLE IF NOT EXISTS `messages` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `app_id` INT UNSIGNED NOT NULL COMMENT 'APP ID',
  `user_id` INT UNSIGNED DEFAULT NULL COMMENT '接收用户ID，NULL表示系统消息',
  `title` VARCHAR(255) NOT NULL COMMENT '消息标题',
  `content` TEXT NOT NULL COMMENT '消息内容',
  `type` VARCHAR(50) NOT NULL DEFAULT 'system' COMMENT '消息类型: system, notification, alert',
  `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-未读, 1-已读',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  INDEX `idx_app_id` (`app_id`),
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_status` (`status`),
  INDEX `idx_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息表';

-- 推送记录表
CREATE TABLE IF NOT EXISTS `push_records` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `app_id` INT UNSIGNED NOT NULL COMMENT 'APP ID',
  `title` VARCHAR(255) NOT NULL COMMENT '推送标题',
  `content` TEXT NOT NULL COMMENT '推送内容',
  `target_type` VARCHAR(50) NOT NULL DEFAULT 'all' COMMENT '推送目标: all-全部用户, user-指定用户, tag-标签用户',
  `target_ids` TEXT DEFAULT NULL COMMENT '目标用户ID列表(JSON数组)',
  `status` VARCHAR(50) NOT NULL DEFAULT 'pending' COMMENT '状态: pending-待发送, sending-发送中, success-成功, failed-失败',
  `sent_count` INT NOT NULL DEFAULT 0 COMMENT '已发送数量',
  `success_count` INT NOT NULL DEFAULT 0 COMMENT '成功数量',
  `failed_count` INT NOT NULL DEFAULT 0 COMMENT '失败数量',
  `scheduled_at` DATETIME DEFAULT NULL COMMENT '定时发送时间',
  `sent_at` DATETIME DEFAULT NULL COMMENT '实际发送时间',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  INDEX `idx_app_id` (`app_id`),
  INDEX `idx_status` (`status`),
  INDEX `idx_scheduled_at` (`scheduled_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='推送记录表';

-- 事件表
CREATE TABLE IF NOT EXISTS `events` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `app_id` INT UNSIGNED NOT NULL COMMENT 'APP ID',
  `user_id` INT UNSIGNED DEFAULT NULL COMMENT '用户ID',
  `event_code` VARCHAR(100) NOT NULL COMMENT '事件代码',
  `event_name` VARCHAR(255) NOT NULL COMMENT '事件名称',
  `properties` JSON DEFAULT NULL COMMENT '事件属性(JSON)',
  `ip` VARCHAR(50) DEFAULT NULL COMMENT 'IP地址',
  `user_agent` VARCHAR(500) DEFAULT NULL COMMENT 'User Agent',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_app_id` (`app_id`),
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_event_code` (`event_code`),
  INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='事件表';

-- 事件定义表
CREATE TABLE IF NOT EXISTS `event_definitions` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `app_id` INT UNSIGNED NOT NULL COMMENT 'APP ID',
  `event_code` VARCHAR(100) NOT NULL COMMENT '事件代码',
  `event_name` VARCHAR(255) NOT NULL COMMENT '事件名称',
  `description` TEXT DEFAULT NULL COMMENT '事件描述',
  `properties_schema` JSON DEFAULT NULL COMMENT '属性Schema(JSON)',
  `is_active` TINYINT NOT NULL DEFAULT 1 COMMENT '是否启用',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  UNIQUE KEY `uk_app_event` (`app_id`, `event_code`),
  INDEX `idx_is_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='事件定义表';

-- 日志表
CREATE TABLE IF NOT EXISTS `logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `app_id` INT UNSIGNED NOT NULL COMMENT 'APP ID',
  `level` VARCHAR(20) NOT NULL DEFAULT 'info' COMMENT '日志级别: debug, info, warn, error, fatal',
  `module` VARCHAR(100) DEFAULT NULL COMMENT '模块名称',
  `message` TEXT NOT NULL COMMENT '日志消息',
  `context` JSON DEFAULT NULL COMMENT '上下文信息(JSON)',
  `ip` VARCHAR(50) DEFAULT NULL COMMENT 'IP地址',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_app_id` (`app_id`),
  INDEX `idx_level` (`level`),
  INDEX `idx_module` (`module`),
  INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='日志表';

-- 监控指标表
CREATE TABLE IF NOT EXISTS `monitor_metrics` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `app_id` INT UNSIGNED NOT NULL COMMENT 'APP ID',
  `metric_name` VARCHAR(100) NOT NULL COMMENT '指标名称',
  `metric_value` DECIMAL(20,4) NOT NULL COMMENT '指标值',
  `tags` JSON DEFAULT NULL COMMENT '标签(JSON)',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_app_id` (`app_id`),
  INDEX `idx_metric_name` (`metric_name`),
  INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='监控指标表';

-- 告警表
CREATE TABLE IF NOT EXISTS `monitor_alerts` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `app_id` INT UNSIGNED NOT NULL COMMENT 'APP ID',
  `alert_name` VARCHAR(255) NOT NULL COMMENT '告警名称',
  `metric_name` VARCHAR(100) NOT NULL COMMENT '监控指标',
  `condition` VARCHAR(50) NOT NULL COMMENT '条件: gt, lt, eq, gte, lte',
  `threshold` DECIMAL(20,4) NOT NULL COMMENT '阈值',
  `status` VARCHAR(50) NOT NULL DEFAULT 'normal' COMMENT '状态: normal-正常, alerting-告警中',
  `last_alert_at` DATETIME DEFAULT NULL COMMENT '最后告警时间',
  `is_active` TINYINT NOT NULL DEFAULT 1 COMMENT '是否启用',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  INDEX `idx_app_id` (`app_id`),
  INDEX `idx_metric_name` (`metric_name`),
  INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='告警表';

-- 文件表
CREATE TABLE IF NOT EXISTS `files` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `app_id` INT UNSIGNED NOT NULL COMMENT 'APP ID',
  `filename` VARCHAR(255) NOT NULL COMMENT '文件名',
  `file_path` VARCHAR(500) NOT NULL COMMENT '文件路径',
  `file_size` BIGINT NOT NULL DEFAULT 0 COMMENT '文件大小(字节)',
  `mime_type` VARCHAR(100) DEFAULT NULL COMMENT 'MIME类型',
  `upload_by` INT UNSIGNED DEFAULT NULL COMMENT '上传者ID',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  INDEX `idx_app_id` (`app_id`),
  INDEX `idx_upload_by` (`upload_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件表';

-- 配置表
CREATE TABLE IF NOT EXISTS `configs` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `app_id` INT UNSIGNED NOT NULL COMMENT 'APP ID',
  `config_key` VARCHAR(255) NOT NULL COMMENT '配置键',
  `config_value` TEXT NOT NULL COMMENT '配置值',
  `description` TEXT DEFAULT NULL COMMENT '配置描述',
  `is_published` TINYINT NOT NULL DEFAULT 0 COMMENT '是否已发布',
  `published_at` DATETIME DEFAULT NULL COMMENT '发布时间',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  UNIQUE KEY `uk_app_key` (`app_id`, `config_key`),
  INDEX `idx_is_published` (`is_published`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='配置表';

-- 配置历史表
CREATE TABLE IF NOT EXISTS `config_history` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `config_id` INT UNSIGNED NOT NULL COMMENT '配置ID',
  `config_value` TEXT NOT NULL COMMENT '配置值',
  `operator_id` INT UNSIGNED DEFAULT NULL COMMENT '操作者ID',
  `operation` VARCHAR(50) NOT NULL COMMENT '操作类型: create, update, publish',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_config_id` (`config_id`),
  INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='配置历史表';

-- 版本表
CREATE TABLE IF NOT EXISTS `versions` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `app_id` INT UNSIGNED NOT NULL COMMENT 'APP ID',
  `version_name` VARCHAR(50) NOT NULL COMMENT '版本号',
  `version_code` INT NOT NULL COMMENT '版本代码',
  `description` TEXT DEFAULT NULL COMMENT '版本描述',
  `download_url` VARCHAR(500) DEFAULT NULL COMMENT '下载地址',
  `is_force_update` TINYINT NOT NULL DEFAULT 0 COMMENT '是否强制更新',
  `status` VARCHAR(50) NOT NULL DEFAULT 'draft' COMMENT '状态: draft-草稿, published-已发布, archived-已归档',
  `published_at` DATETIME DEFAULT NULL COMMENT '发布时间',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  INDEX `idx_app_id` (`app_id`),
  INDEX `idx_version_code` (`version_code`),
  INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='版本表';
