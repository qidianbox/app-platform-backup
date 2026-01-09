# Terraform配置 - 阿里云ACK容器服务部署
# 版本：1.0
# 作者：Manus AI

terraform {
  required_version = ">= 1.0"
  required_providers {
    alicloud = {
      source  = "aliyun/alicloud"
      version = "~> 1.200"
    }
  }
}

# 阿里云Provider配置
provider "alicloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}

# 数据源：获取可用区
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

# ==================== VPC网络 ====================

# 创建VPC
resource "alicloud_vpc" "main" {
  vpc_name   = "${var.project_name}-vpc"
  cidr_block = "10.0.0.0/16"
  tags = {
    Project = var.project_name
    Env     = var.environment
  }
}

# 创建交换机 - 应用层
resource "alicloud_vswitch" "app" {
  vpc_id       = alicloud_vpc.main.id
  cidr_block   = "10.0.1.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = "${var.project_name}-app-vswitch"
  tags = {
    Project = var.project_name
    Layer   = "app"
  }
}

# 创建交换机 - 数据层
resource "alicloud_vswitch" "data" {
  vpc_id       = alicloud_vpc.main.id
  cidr_block   = "10.0.2.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = "${var.project_name}-data-vswitch"
  tags = {
    Project = var.project_name
    Layer   = "data"
  }
}

# ==================== 安全组 ====================

# ACK节点安全组
resource "alicloud_security_group" "ack" {
  name        = "${var.project_name}-ack-sg"
  vpc_id      = alicloud_vpc.main.id
  description = "Security group for ACK nodes"
  tags = {
    Project = var.project_name
  }
}

# 安全组规则 - 允许内部通信
resource "alicloud_security_group_rule" "internal" {
  type              = "ingress"
  ip_protocol       = "all"
  port_range        = "-1/-1"
  security_group_id = alicloud_security_group.ack.id
  cidr_ip           = "10.0.0.0/16"
}

# 安全组规则 - 允许HTTP
resource "alicloud_security_group_rule" "http" {
  type              = "ingress"
  ip_protocol       = "tcp"
  port_range        = "80/80"
  security_group_id = alicloud_security_group.ack.id
  cidr_ip           = "0.0.0.0/0"
}

# 安全组规则 - 允许HTTPS
resource "alicloud_security_group_rule" "https" {
  type              = "ingress"
  ip_protocol       = "tcp"
  port_range        = "443/443"
  security_group_id = alicloud_security_group.ack.id
  cidr_ip           = "0.0.0.0/0"
}

# ==================== RDS MySQL ====================

# 创建RDS实例
resource "alicloud_db_instance" "main" {
  engine               = "MySQL"
  engine_version       = "8.0"
  instance_type        = var.rds_instance_type
  instance_storage     = var.rds_storage
  instance_name        = "${var.project_name}-rds"
  vswitch_id           = alicloud_vswitch.data.id
  security_ips         = ["10.0.0.0/16"]
  instance_charge_type = "Postpaid"
  
  tags = {
    Project = var.project_name
    Env     = var.environment
  }
}

# 创建数据库
resource "alicloud_db_database" "main" {
  instance_id = alicloud_db_instance.main.id
  name        = var.db_name
  character_set = "utf8mb4"
}

# 创建数据库账号
resource "alicloud_db_account" "main" {
  db_instance_id   = alicloud_db_instance.main.id
  account_name     = var.db_user
  account_password = var.db_password
  account_type     = "Super"
}

# 授权数据库访问
resource "alicloud_db_account_privilege" "main" {
  instance_id  = alicloud_db_instance.main.id
  account_name = alicloud_db_account.main.account_name
  privilege    = "ReadWrite"
  db_names     = [alicloud_db_database.main.name]
}

# ==================== ACK集群 ====================

# 创建ACK托管集群
resource "alicloud_cs_managed_kubernetes" "main" {
  name                 = "${var.project_name}-ack"
  cluster_spec         = "ack.pro.small"
  version              = var.k8s_version
  
  # 网络配置
  pod_cidr             = "172.20.0.0/16"
  service_cidr         = "172.21.0.0/20"
  
  # Worker节点配置
  worker_vswitch_ids   = [alicloud_vswitch.app.id]
  
  # 其他配置
  new_nat_gateway      = true
  slb_internet_enabled = true
  
  # 组件配置
  addons {
    name = "terway-eniip"
  }
  addons {
    name = "csi-plugin"
  }
  addons {
    name = "csi-provisioner"
  }
  addons {
    name = "nginx-ingress-controller"
    disabled = true  # 使用ALB Ingress
  }
  
  tags = {
    Project = var.project_name
    Env     = var.environment
  }
}

# 创建节点池
resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                 = "${var.project_name}-nodepool"
  cluster_id           = alicloud_cs_managed_kubernetes.main.id
  vswitch_ids          = [alicloud_vswitch.app.id]
  
  # 节点配置
  instance_types       = [var.node_instance_type]
  system_disk_category = "cloud_essd"
  system_disk_size     = 120
  
  # 扩缩容配置
  desired_size         = var.node_desired_size
  
  # 自动扩缩容
  scaling_config {
    min_size = var.node_min_size
    max_size = var.node_max_size
  }
  
  # 安全组
  security_group_ids   = [alicloud_security_group.ack.id]
  
  tags = {
    Project = var.project_name
    Env     = var.environment
  }
}

# ==================== ACR容器镜像仓库 ====================

# 创建命名空间
resource "alicloud_cr_namespace" "main" {
  name               = var.project_name
  auto_create        = false
  default_visibility = "PRIVATE"
}

# 创建镜像仓库
resource "alicloud_cr_repo" "backend" {
  namespace = alicloud_cr_namespace.main.name
  name      = "backend"
  summary   = "Backend service image"
  repo_type = "PRIVATE"
}

# ==================== OSS存储桶 ====================

# 创建OSS存储桶（前端静态资源）
resource "alicloud_oss_bucket" "frontend" {
  bucket = "${var.project_name}-frontend-${var.region}"
  acl    = "public-read"
  
  website {
    index_document = "index.html"
    error_document = "index.html"
  }
  
  cors_rule {
    allowed_origins = ["*"]
    allowed_methods = ["GET", "HEAD"]
    allowed_headers = ["*"]
    max_age_seconds = 3600
  }
  
  tags = {
    Project = var.project_name
    Env     = var.environment
  }
}

# ==================== 输出 ====================

output "vpc_id" {
  value = alicloud_vpc.main.id
}

output "ack_cluster_id" {
  value = alicloud_cs_managed_kubernetes.main.id
}

output "rds_connection_string" {
  value = alicloud_db_instance.main.connection_string
}

output "rds_port" {
  value = alicloud_db_instance.main.port
}

output "acr_namespace" {
  value = alicloud_cr_namespace.main.name
}

output "oss_bucket" {
  value = alicloud_oss_bucket.frontend.bucket
}

output "oss_endpoint" {
  value = alicloud_oss_bucket.frontend.extranet_endpoint
}
