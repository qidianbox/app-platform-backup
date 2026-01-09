# APP中台管理系统 - 阿里云ACK容器服务部署方案

**作者：Manus AI**  
**日期：2026年1月9日**  
**版本：2.0**

---

## 一、为什么选择ACK容器服务

阿里云ACK（Alibaba Cloud Container Service for Kubernetes）是企业级Kubernetes容器服务，相比传统ECS弹性伸缩具有显著优势：

| 特性 | ACK容器服务 | ECS弹性伸缩 |
|------|-------------|-------------|
| 扩容速度 | **秒级**（Pod调度） | 分钟级（ECS启动） |
| 资源利用率 | **高**（多服务共享节点） | 低（独立实例） |
| 滚动更新 | **原生支持** | 需要脚本 |
| 服务发现 | **原生支持** | 需要额外配置 |
| 健康检查 | **完善**（存活/就绪探针） | 基础 |
| 微服务架构 | **原生支持** | 一般 |
| CI/CD集成 | **原生支持** | 需要额外配置 |
| 行业标准 | **Kubernetes标准** | 阿里云专有 |

---

## 二、架构设计

### 2.1 整体架构图

```
                         ┌─────────────────┐
                         │    用户请求      │
                         └────────┬────────┘
                                  │
                         ┌────────▼────────┐
                         │   CDN (前端)     │ ← OSS静态资源
                         └────────┬────────┘
                                  │
                         ┌────────▼────────┐
                         │   ALB Ingress   │ ← 应用负载均衡
                         └────────┬────────┘
                                  │
              ┌───────────────────┼───────────────────┐
              │                   │                   │
       ┌──────▼──────┐     ┌──────▼──────┐     ┌──────▼──────┐
       │   Pod-1     │     │   Pod-2     │     │   Pod-N     │
       │  (后端)     │     │  (后端)     │     │  (后端)     │
       └──────┬──────┘     └──────┬──────┘     └──────┬──────┘
              │                   │                   │
              └───────────────────┼───────────────────┘
                                  │
                         ┌────────▼────────┐
                         │   RDS MySQL     │ ← 主从高可用
                         └─────────────────┘
```

### 2.2 Kubernetes资源规划

| 资源类型 | 名称 | 数量 | 说明 |
|----------|------|------|------|
| Namespace | app-platform | 1 | 应用命名空间 |
| Deployment | backend | 1 | 后端服务部署 |
| Service | backend-svc | 1 | 后端服务暴露 |
| Ingress | backend-ingress | 1 | ALB入口配置 |
| HPA | backend-hpa | 1 | 水平Pod自动扩缩 |
| ConfigMap | app-config | 1 | 应用配置 |
| Secret | app-secrets | 1 | 敏感信息 |

---

## 三、阿里云资源规划

### 3.1 ACK集群配置

| 配置项 | 推荐值 | 说明 |
|--------|--------|------|
| 集群类型 | **ACK托管版** | 免费托管Master节点 |
| Kubernetes版本 | 1.28+ | 最新稳定版 |
| 节点池规格 | ecs.c6.xlarge (4核8G) | 可运行多个Pod |
| 节点数量 | 2-5台（弹性） | 根据负载自动扩缩 |
| 容器运行时 | containerd | 推荐 |
| 网络插件 | Terway | 高性能网络 |

### 3.2 完整资源清单

| 资源类型 | 规格 | 数量 | 用途 | 预估月费用 |
|----------|------|------|------|------------|
| ACK托管集群 | 托管版 | 1 | Kubernetes集群 | **免费** |
| ECS节点池 | ecs.c6.xlarge (4核8G) | 2-5台 | 工作节点 | ¥600-1500 |
| RDS MySQL | mysql.n4.medium.1 (2核4G) | 1主1从 | 数据库 | ¥500-800 |
| ALB | 按规格计费 | 1 | 应用负载均衡 | ¥100 |
| OSS | 标准存储 | 按需 | 前端静态资源 | ¥10-50 |
| CDN | 按流量计费 | 按需 | 前端加速 | ¥50-200 |
| ACR | 基础版 | 1 | 容器镜像仓库 | **免费** |
| NAT网关 | 小型 | 1 | 出网访问 | ¥50 |

**预估总成本：¥1300-2700/月**（根据实际负载浮动）

---

## 四、所需阿里云权限

### 4.1 RAM权限策略

| 权限策略 | 用途 | 必需程度 |
|----------|------|----------|
| AliyunCSFullAccess | 创建和管理ACK集群 | **必需** |
| AliyunECSFullAccess | 管理ECS节点 | **必需** |
| AliyunVPCFullAccess | 创建和管理VPC网络 | **必需** |
| AliyunRDSFullAccess | 创建和管理RDS数据库 | **必需** |
| AliyunOSSFullAccess | 创建和管理OSS存储桶 | **必需** |
| AliyunCDNFullAccess | 创建和管理CDN加速 | **必需** |
| AliyunALBFullAccess | 管理应用负载均衡 | **必需** |
| AliyunContainerRegistryFullAccess | 管理容器镜像仓库 | **必需** |
| AliyunNATGatewayFullAccess | 管理NAT网关 | **必需** |
| AliyunRAMReadOnlyAccess | 读取RAM配置 | 推荐 |

### 4.2 需要您提供的凭证

```bash
# 阿里云AccessKey（用于自动化部署）
ALICLOUD_ACCESS_KEY=<your-access-key-id>
ALICLOUD_SECRET_KEY=<your-access-key-secret>
ALICLOUD_REGION=cn-hangzhou  # 或其他区域

# 数据库配置
RDS_PASSWORD=<your-secure-password>  # 至少8位，包含大小写和数字

# 应用配置
JWT_SECRET=<random-32-char-string>   # 用于用户认证
```

---

## 五、自动化部署流程

### 5.1 部署阶段

**阶段一：基础设施创建（Terraform，约15分钟）**

| 步骤 | 操作 | 自动化 |
|------|------|--------|
| 1 | 创建VPC和交换机 | ✅ 全自动 |
| 2 | 创建ACK托管集群 | ✅ 全自动 |
| 3 | 创建节点池 | ✅ 全自动 |
| 4 | 创建RDS MySQL实例 | ✅ 全自动 |
| 5 | 创建ACR镜像仓库 | ✅ 全自动 |
| 6 | 创建OSS存储桶 | ✅ 全自动 |
| 7 | 配置ALB Ingress | ✅ 全自动 |

**阶段二：应用部署（kubectl，约10分钟）**

| 步骤 | 操作 | 自动化 |
|------|------|--------|
| 1 | 构建Docker镜像 | ✅ 全自动 |
| 2 | 推送镜像到ACR | ✅ 全自动 |
| 3 | 创建Namespace | ✅ 全自动 |
| 4 | 部署ConfigMap和Secret | ✅ 全自动 |
| 5 | 部署后端Deployment | ✅ 全自动 |
| 6 | 配置Service和Ingress | ✅ 全自动 |
| 7 | 配置HPA自动扩缩 | ✅ 全自动 |
| 8 | 数据库迁移 | ✅ 全自动 |

**阶段三：前端部署（约5分钟）**

| 步骤 | 操作 | 自动化 |
|------|------|--------|
| 1 | 构建前端静态文件 | ✅ 全自动 |
| 2 | 上传到OSS | ✅ 全自动 |
| 3 | 配置CDN加速 | ✅ 全自动 |
| 4 | 刷新CDN缓存 | ✅ 全自动 |

**总部署时间：约30分钟**

### 5.2 弹性扩缩配置

**Pod级别自动扩缩（HPA）**

| 指标 | 扩容阈值 | 缩容阈值 | 最小Pod | 最大Pod |
|------|----------|----------|---------|---------|
| CPU使用率 | > 70% | < 30% | 2 | 20 |
| 内存使用率 | > 80% | < 40% | 2 | 20 |

**节点级别自动扩缩（Cluster Autoscaler）**

| 配置项 | 值 | 说明 |
|--------|-----|------|
| 最小节点数 | 2 | 保证高可用 |
| 最大节点数 | 10 | 成本控制 |
| 扩容触发 | Pod无法调度 | 自动添加节点 |
| 缩容触发 | 节点利用率<50% | 自动移除节点 |

---

## 六、我可以帮您自动完成的操作

### 6.1 完全自动化清单

| 操作 | 说明 | 预计时间 |
|------|------|----------|
| ✅ 创建Terraform配置 | 基础设施即代码 | 5分钟 |
| ✅ 创建Dockerfile | 后端容器化 | 2分钟 |
| ✅ 创建K8s YAML配置 | Deployment/Service/Ingress/HPA | 5分钟 |
| ✅ 执行Terraform部署 | 创建ACK集群和所有云资源 | 15分钟 |
| ✅ 构建并推送Docker镜像 | 后端镜像到ACR | 3分钟 |
| ✅ 部署K8s资源 | 应用到集群 | 2分钟 |
| ✅ 配置HPA自动扩缩 | Pod弹性伸缩 | 1分钟 |
| ✅ 部署前端到OSS | 静态资源上传 | 2分钟 |
| ✅ 配置CDN和域名 | 加速和解析 | 3分钟 |
| ✅ 数据库初始化 | 表结构和初始数据 | 2分钟 |

### 6.2 需要您手动操作

| 操作 | 原因 | 操作指南 |
|------|------|----------|
| 创建RAM用户和AccessKey | 安全考虑 | 见第四节 |
| 域名备案 | 中国大陆法规 | 阿里云备案系统 |
| SSL证书申请 | 需要域名验证 | 可使用阿里云免费证书 |

---

## 七、Kubernetes配置文件预览

### 7.1 Deployment配置

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-platform-backend
  namespace: app-platform
spec:
  replicas: 2
  selector:
    matchLabels:
      app: backend
  template:
    metadata:
      labels:
        app: backend
    spec:
      containers:
      - name: backend
        image: registry.cn-hangzhou.aliyuncs.com/app-platform/backend:latest
        ports:
        - containerPort: 8080
        resources:
          requests:
            cpu: "250m"
            memory: "256Mi"
          limits:
            cpu: "1000m"
            memory: "1Gi"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        envFrom:
        - configMapRef:
            name: app-config
        - secretRef:
            name: app-secrets
```

### 7.2 HPA自动扩缩配置

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: backend-hpa
  namespace: app-platform
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: app-platform-backend
  minReplicas: 2
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
      - type: Percent
        value: 100
        periodSeconds: 15
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 50
        periodSeconds: 60
```

---

## 八、成本对比

| 方案 | 基础配置月费用 | 生产配置月费用 | 优势 |
|------|----------------|----------------|------|
| ECS弹性伸缩 | ¥850 | ¥2900 | 简单直接 |
| **ACK容器服务** | **¥1300** | **¥2700** | 更灵活、更标准化 |

ACK方案虽然基础配置略贵，但在生产环境下由于更高的资源利用率，实际成本可能更低。

---

## 九、下一步行动

如果您决定使用ACK容器服务部署，请按以下步骤操作：

1. **创建RAM用户**：登录阿里云控制台，创建具有上述权限的RAM用户
2. **获取AccessKey**：保存AccessKey ID和Secret
3. **提供凭证**：将以下信息提供给我：
   - AccessKey ID
   - AccessKey Secret
   - 部署区域（如cn-hangzhou）
   - RDS密码
4. **确认配置**：确认节点规格、数据库配置
5. **开始部署**：我将自动执行所有部署操作

**预计总部署时间：30分钟**

---

## 附录：CI/CD流水线配置

部署完成后，我还可以为您配置GitHub Actions自动化流水线，实现代码提交后自动部署：

```yaml
# .github/workflows/deploy.yml
name: Deploy to ACK

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Login to ACR
      uses: aliyun/acr-login@v1
      with:
        region-id: cn-hangzhou
        access-key-id: ${{ secrets.ALICLOUD_ACCESS_KEY }}
        access-key-secret: ${{ secrets.ALICLOUD_SECRET_KEY }}
    
    - name: Build and Push Image
      run: |
        docker build -t registry.cn-hangzhou.aliyuncs.com/app-platform/backend:${{ github.sha }} .
        docker push registry.cn-hangzhou.aliyuncs.com/app-platform/backend:${{ github.sha }}
    
    - name: Deploy to ACK
      uses: aliyun/ack-deploy@v1
      with:
        cluster-id: ${{ secrets.ACK_CLUSTER_ID }}
        image: registry.cn-hangzhou.aliyuncs.com/app-platform/backend:${{ github.sha }}
```

---

**文档结束**

如有任何问题，请随时联系。
