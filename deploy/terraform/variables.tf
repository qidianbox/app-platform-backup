# Terraform变量定义

# ==================== 阿里云认证 ====================

variable "access_key" {
  description = "阿里云AccessKey ID"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "阿里云AccessKey Secret"
  type        = string
  sensitive   = true
}

variable "region" {
  description = "部署区域"
  type        = string
  default     = "cn-hangzhou"
}

# ==================== 项目配置 ====================

variable "project_name" {
  description = "项目名称"
  type        = string
  default     = "app-platform"
}

variable "environment" {
  description = "环境标识"
  type        = string
  default     = "production"
}

# ==================== Kubernetes配置 ====================

variable "k8s_version" {
  description = "Kubernetes版本"
  type        = string
  default     = "1.28.3-aliyun.1"
}

variable "node_instance_type" {
  description = "节点实例规格"
  type        = string
  default     = "ecs.c6.xlarge"  # 4核8G
}

variable "node_min_size" {
  description = "节点池最小节点数"
  type        = number
  default     = 2
}

variable "node_max_size" {
  description = "节点池最大节点数"
  type        = number
  default     = 10
}

variable "node_desired_size" {
  description = "节点池期望节点数"
  type        = number
  default     = 2
}

# ==================== RDS配置 ====================

variable "rds_instance_type" {
  description = "RDS实例规格"
  type        = string
  default     = "mysql.n4.medium.1"  # 2核4G
}

variable "rds_storage" {
  description = "RDS存储空间(GB)"
  type        = number
  default     = 50
}

variable "db_name" {
  description = "数据库名称"
  type        = string
  default     = "app_platform"
}

variable "db_user" {
  description = "数据库用户名"
  type        = string
  default     = "app_platform"
}

variable "db_password" {
  description = "数据库密码"
  type        = string
  sensitive   = true
}

# ==================== 应用配置 ====================

variable "jwt_secret" {
  description = "JWT密钥"
  type        = string
  sensitive   = true
}
