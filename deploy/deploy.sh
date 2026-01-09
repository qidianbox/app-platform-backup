#!/bin/bash
# APP中台管理系统 - 阿里云ACK自动化部署脚本
# 作者：Manus AI
# 版本：1.0

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查必要的工具
check_prerequisites() {
    log_info "检查必要的工具..."
    
    if ! command -v terraform &> /dev/null; then
        log_error "Terraform 未安装，请先安装 Terraform"
        exit 1
    fi
    
    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl 未安装，请先安装 kubectl"
        exit 1
    fi
    
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    log_info "所有必要工具已安装"
}

# 阶段一：部署基础设施
deploy_infrastructure() {
    log_info "========== 阶段一：部署基础设施 =========="
    
    cd terraform
    
    log_info "初始化 Terraform..."
    terraform init
    
    log_info "验证 Terraform 配置..."
    terraform validate
    
    log_info "规划基础设施变更..."
    terraform plan -out=tfplan
    
    log_info "应用基础设施变更（这可能需要15-20分钟）..."
    terraform apply tfplan
    
    log_info "基础设施部署完成！"
    
    # 保存输出
    terraform output -json > ../outputs.json
    
    cd ..
}

# 阶段二：配置kubectl
configure_kubectl() {
    log_info "========== 阶段二：配置kubectl =========="
    
    CLUSTER_ID=$(cat outputs.json | jq -r '.ack_cluster_id.value')
    REGION=$(cat terraform/terraform.tfvars | grep region | awk -F'"' '{print $2}')
    
    log_info "获取ACK集群kubeconfig..."
    aliyun cs GET /k8s/$CLUSTER_ID/user_config | jq -r '.config' > ~/.kube/config
    
    log_info "验证kubectl连接..."
    kubectl cluster-info
    
    log_info "kubectl配置完成！"
}

# 阶段三：构建并推送Docker镜像
build_and_push_image() {
    log_info "========== 阶段三：构建并推送Docker镜像 =========="
    
    REGION=$(cat terraform/terraform.tfvars | grep region | awk -F'"' '{print $2}')
    NAMESPACE=$(cat outputs.json | jq -r '.acr_namespace.value')
    
    IMAGE_TAG="registry.${REGION}.aliyuncs.com/${NAMESPACE}/backend:$(date +%Y%m%d%H%M%S)"
    
    log_info "构建Docker镜像..."
    cd ../backend
    docker build -t $IMAGE_TAG .
    
    log_info "登录ACR..."
    docker login --username=$ALICLOUD_ACCESS_KEY registry.${REGION}.aliyuncs.com
    
    log_info "推送镜像到ACR..."
    docker push $IMAGE_TAG
    
    log_info "镜像推送完成: $IMAGE_TAG"
    
    # 保存镜像标签
    echo $IMAGE_TAG > ../deploy/image_tag.txt
    
    cd ../deploy
}

# 阶段四：部署Kubernetes资源
deploy_kubernetes() {
    log_info "========== 阶段四：部署Kubernetes资源 =========="
    
    IMAGE_TAG=$(cat image_tag.txt)
    RDS_HOST=$(cat outputs.json | jq -r '.rds_connection_string.value')
    
    log_info "创建命名空间..."
    kubectl apply -f k8s/namespace.yaml
    
    log_info "创建ConfigMap..."
    kubectl apply -f k8s/configmap.yaml
    
    log_info "创建Secret..."
    # 替换变量
    envsubst < k8s/secret.yaml.template > k8s/secret.yaml
    kubectl apply -f k8s/secret.yaml
    rm k8s/secret.yaml  # 删除包含敏感信息的文件
    
    log_info "部署后端服务..."
    # 更新镜像标签
    sed -i "s|image:.*|image: $IMAGE_TAG|g" k8s/deployment.yaml
    kubectl apply -f k8s/deployment.yaml
    
    log_info "创建Service..."
    kubectl apply -f k8s/service.yaml
    
    log_info "配置Ingress..."
    kubectl apply -f k8s/ingress.yaml
    
    log_info "配置HPA自动扩缩..."
    kubectl apply -f k8s/hpa.yaml
    
    log_info "等待Pod就绪..."
    kubectl rollout status deployment/app-platform-backend -n app-platform --timeout=300s
    
    log_info "Kubernetes资源部署完成！"
}

# 阶段五：部署前端
deploy_frontend() {
    log_info "========== 阶段五：部署前端 =========="
    
    OSS_BUCKET=$(cat outputs.json | jq -r '.oss_bucket.value')
    OSS_ENDPOINT=$(cat outputs.json | jq -r '.oss_endpoint.value')
    
    log_info "构建前端..."
    cd ../frontend
    npm install
    npm run build
    
    log_info "上传到OSS..."
    aliyun oss cp -r dist/ oss://$OSS_BUCKET/ --force
    
    log_info "前端部署完成！"
    log_info "前端访问地址: http://$OSS_BUCKET.$OSS_ENDPOINT"
    
    cd ../deploy
}

# 阶段六：数据库迁移
migrate_database() {
    log_info "========== 阶段六：数据库迁移 =========="
    
    log_info "执行数据库迁移..."
    kubectl exec -it deployment/app-platform-backend -n app-platform -- /app/server migrate
    
    log_info "数据库迁移完成！"
}

# 显示部署结果
show_results() {
    log_info "========== 部署完成 =========="
    
    OSS_BUCKET=$(cat outputs.json | jq -r '.oss_bucket.value')
    OSS_ENDPOINT=$(cat outputs.json | jq -r '.oss_endpoint.value')
    
    echo ""
    echo "=================================================="
    echo "  APP中台管理系统部署成功！"
    echo "=================================================="
    echo ""
    echo "前端地址: http://$OSS_BUCKET.$OSS_ENDPOINT"
    echo ""
    echo "后端API: 请查看Ingress获取ALB地址"
    kubectl get ingress -n app-platform
    echo ""
    echo "Pod状态:"
    kubectl get pods -n app-platform
    echo ""
    echo "HPA状态:"
    kubectl get hpa -n app-platform
    echo ""
    echo "=================================================="
}

# 主函数
main() {
    log_info "开始部署 APP中台管理系统 到阿里云ACK..."
    
    check_prerequisites
    deploy_infrastructure
    configure_kubectl
    build_and_push_image
    deploy_kubernetes
    deploy_frontend
    migrate_database
    show_results
    
    log_info "所有部署步骤完成！"
}

# 执行主函数
main "$@"
