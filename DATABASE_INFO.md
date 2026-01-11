# 数据库信息

## 数据库配置

- **数据库类型**: MySQL 8.0
- **主机地址**: rm-bp13s51058fu3r061.mysql.rds.aliyuncs.com
- **端口**: 3306
- **数据库名**: app_platform
- **用户名**: app_platform
- **密码**: App@Platform123
- **字符集**: utf8mb4

## 连接配置

```yaml
database:
  driver: mysql
  host: rm-bp13s51058fu3r061.mysql.rds.aliyuncs.com
  port: 3306
  username: app_platform
  password: App@Platform123
  database: app_platform
  charset: utf8mb4
  parse_time: true
  loc: Local
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600
```

## 数据库说明

由于数据库服务器位于阿里云RDS，无法从外网直接访问导出数据。如需备份数据库，请通过以下方式：

1. **通过阿里云控制台**：登录阿里云RDS控制台，使用备份功能导出数据
2. **通过内网访问**：在阿里云ECS实例上执行mysqldump命令
3. **使用数据传输服务**：使用阿里云DTS服务进行数据迁移和备份

## 数据库表结构

主要数据表包括：
- `apps` - APP应用表
- `users` - 用户表
- `modules` - 模块表
- `configs` - 配置表
- `logs` - 日志表
- `audit_logs` - 审计日志表
- `messages` - 消息表
- `push_records` - 推送记录表
- `versions` - 版本管理表

详细表结构请参考后端代码中的数据库迁移文件。

## 备份建议

1. 定期通过阿里云RDS控制台创建自动备份
2. 重要数据变更前手动创建备份快照
3. 定期导出数据到本地存储
4. 使用阿里云OSS存储备份文件

## 注意事项

⚠️ **安全提示**：
- 数据库密码已在配置文件中明文存储，建议使用环境变量或密钥管理服务
- 生产环境应限制数据库访问IP白名单
- 定期更换数据库密码
- 启用SSL连接加密
