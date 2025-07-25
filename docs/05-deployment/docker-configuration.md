# Docker 配置和部署文档

**版本**: 2.0  
**更新时间**: 2025-01-25  
**技术栈**: Docker + Kubernetes + Prometheus + Grafana  
**状态**: 已确认

---

## 概述

本文档详细说明 Lyss AI Platform 的 Docker 容器化配置和部署方案。支持从本地开发到生产环境的完整部署流程。

### 部署架构

```
生产环境架构：
┌─────────────────────────────────────────┐
│  Kubernetes Cluster                    │
│  ├── Ingress Controller                │
│  ├── Gateway Service (3 replicas)      │
│  ├── User Service (2 replicas)         │
│  ├── Auth Service (2 replicas)         │
│  ├── Group Service (2 replicas)        │
│  ├── Credential Service (2 replicas)   │
│  └── Billing Service (2 replicas)      │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│  Infrastructure Services                │
│  ├── PostgreSQL Cluster                │
│  ├── Redis Cluster                     │
│  ├── Consul Cluster                    │
│  ├── Qdrant Vector DB                  │
│  └── Monitoring Stack                  │
└─────────────────────────────────────────┘
```

---

## Docker 镜像构建

### 1. 多阶段构建 Dockerfile

#### 微服务通用 Dockerfile 模板

**services/gateway-service/Dockerfile**：
```dockerfile
# === Build Stage ===
FROM golang:1.21-alpine AS builder

# 安装构建依赖
RUN apk add --no-cache git ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖 (利用 Docker 缓存层)
RUN go mod download

# 复制源代码
COPY . .

# 构建应用 (静态链接，减小镜像体积)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main cmd/main.go

# === Runtime Stage ===
FROM scratch

# 从 builder 阶段复制必要文件
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/main /main

# 复制配置文件
COPY --from=builder /app/configs /configs

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非 root 用户
USER 1000:1000

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
    CMD ["/main", "--health-check"]

# 暴露端口
EXPOSE 8000

# 启动应用
ENTRYPOINT ["/main"]
```

#### 前端应用 Dockerfile

**frontend/Dockerfile**：
```dockerfile
# === Build Stage ===
FROM node:18-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 package 文件
COPY package*.json ./

# 安装依赖
RUN npm ci --only=production

# 复制源代码
COPY . .

# 构建应用
RUN npm run build

# === Runtime Stage ===
FROM nginx:alpine

# 复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制 nginx 配置
COPY nginx.conf /etc/nginx/nginx.conf

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:80/health || exit 1

# 暴露端口
EXPOSE 80

# 启动 nginx
CMD ["nginx", "-g", "daemon off;"]
```

### 2. Docker Compose 配置

#### 生产环境 Docker Compose

**docker-compose.prod.yml**：
```yaml
version: '3.8'

services:
  # === 网关服务 ===
  gateway-service:
    build:
      context: ./services/gateway-service
      dockerfile: Dockerfile
    image: lyss/gateway-service:${VERSION:-latest}
    container_name: lyss-gateway
    ports:
      - "8000:8000"
    environment:
      - ENV=production
      - DB_HOST=postgres
      - REDIS_HOST=redis
      - CONSUL_HOST=consul
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
      consul:
        condition: service_healthy
    networks:
      - lyss-network
    restart: unless-stopped
    deploy:
      resources:
        limits:
          memory: 512M
          cpus: '0.5'
        reservations:
          memory: 256M
          cpus: '0.25'
    healthcheck:
      test: ["CMD", "/main", "--health-check"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 60s

  # === 用户服务 ===
  user-service:
    build:
      context: ./services/user-service
      dockerfile: Dockerfile
    image: lyss/user-service:${VERSION:-latest}
    environment:
      - ENV=production
      - DB_HOST=postgres
      - DB_NAME=lyss_user
      - CONSUL_HOST=consul
    depends_on:
      postgres:
        condition: service_healthy
      consul:
        condition: service_healthy
    networks:
      - lyss-network
    restart: unless-stopped
    deploy:
      replicas: 2
    healthcheck:
      test: ["CMD", "/main", "--health-check"]
      interval: 30s
      timeout: 10s
      retries: 3

  # === 认证服务 ===
  auth-service:
    build:
      context: ./services/auth-service
      dockerfile: Dockerfile
    image: lyss/auth-service:${VERSION:-latest}
    environment:
      - ENV=production
      - DB_HOST=postgres
      - DB_NAME=lyss_auth
      - REDIS_HOST=redis
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - lyss-network
    restart: unless-stopped
    deploy:
      replicas: 2

  # === 群组服务 ===
  group-service:
    build:
      context: ./services/group-service
      dockerfile: Dockerfile
    image: lyss/group-service:${VERSION:-latest}
    environment:
      - ENV=production
      - DB_HOST=postgres
      - DB_NAME=lyss_group
    depends_on:
      postgres:
        condition: service_healthy
    networks:
      - lyss-network
    restart: unless-stopped
    deploy:
      replicas: 2

  # === 凭证服务 ===
  credential-service:
    build:
      context: ./services/credential-service
      dockerfile: Dockerfile
    image: lyss/credential-service:${VERSION:-latest}
    environment:
      - ENV=production
      - DB_HOST=postgres
      - DB_NAME=lyss_credential
      - ENCRYPTION_KEY=${ENCRYPTION_KEY}
    depends_on:
      postgres:
        condition: service_healthy
    networks:
      - lyss-network
    restart: unless-stopped
    deploy:
      replicas: 2

  # === 计费服务 ===
  billing-service:
    build:
      context: ./services/billing-service
      dockerfile: Dockerfile
    image: lyss/billing-service:${VERSION:-latest}
    environment:
      - ENV=production
      - DB_HOST=postgres
      - DB_NAME=lyss_billing
      - REDIS_HOST=redis
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - lyss-network
    restart: unless-stopped
    deploy:
      replicas: 2

  # === 前端应用 ===
  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    image: lyss/frontend:${VERSION:-latest}
    container_name: lyss-frontend
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./certs:/etc/nginx/certs:ro
    networks:
      - lyss-network
    restart: unless-stopped
    depends_on:
      - gateway-service

  # === PostgreSQL 数据库 ===
  postgres:
    image: postgres:15-alpine
    container_name: lyss-postgres
    ports:
      - "5432:5432"
    environment:
      POSTGRES_DB: lyss_platform
      POSTGRES_USER: ${DB_USER:-lyss}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_INITDB_ARGS: "--encoding=UTF-8"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./infrastructure/postgres/init.sql:/docker-entrypoint-initdb.d/init.sql
      - ./infrastructure/postgres/postgresql.conf:/etc/postgresql/postgresql.conf
    networks:
      - lyss-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-lyss}"]
      interval: 10s
      timeout: 5s
      retries: 5
    deploy:
      resources:
        limits:
          memory: 2G
          cpus: '1.0'

  # === Redis 缓存 ===
  redis:
    image: redis:7-alpine
    container_name: lyss-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
      - ./infrastructure/redis/redis.conf:/usr/local/etc/redis/redis.conf
    command: redis-server /usr/local/etc/redis/redis.conf
    networks:
      - lyss-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  # === Consul 服务发现 ===
  consul:
    image: consul:1.17
    container_name: lyss-consul
    ports:
      - "8500:8500"
    volumes:
      - consul_data:/consul/data
      - ./infrastructure/consul/consul.hcl:/consul/config/consul.hcl
    command: >
      consul agent -server -bootstrap-expect=1 -ui 
      -bind=0.0.0.0 -client=0.0.0.0 
      -config-file=/consul/config/consul.hcl
    networks:
      - lyss-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "consul", "members"]
      interval: 10s
      timeout: 5s
      retries: 5

  # === Qdrant 向量数据库 ===
  qdrant:
    image: qdrant/qdrant:v1.7.4
    container_name: lyss-qdrant
    ports:
      - "6333:6333"
    volumes:
      - qdrant_data:/qdrant/storage
    environment:
      QDRANT__SERVICE__HTTP_PORT: 6333
    networks:
      - lyss-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:6333/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # === 监控服务 ===
  prometheus:
    image: prom/prometheus:v2.48.1
    container_name: lyss-prometheus
    ports:
      - "9090:9090"
    volumes:
      - prometheus_data:/prometheus
      - ./infrastructure/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--storage.tsdb.retention.time=15d'
      - '--web.enable-lifecycle'
    networks:
      - lyss-network
    restart: unless-stopped

  grafana:
    image: grafana/grafana:10.2.3
    container_name: lyss-grafana
    ports:
      - "3001:3000"
    volumes:
      - grafana_data:/var/lib/grafana
      - ./infrastructure/grafana:/etc/grafana/provisioning
    environment:
      GF_SECURITY_ADMIN_USER: ${GRAFANA_USER:-admin}
      GF_SECURITY_ADMIN_PASSWORD: ${GRAFANA_PASSWORD}
      GF_USERS_ALLOW_SIGN_UP: false
      GF_INSTALL_PLUGINS: grafana-piechart-panel
    networks:
      - lyss-network
    restart: unless-stopped
    depends_on:
      - prometheus

# === 数据卷 ===
volumes:
  postgres_data:
    driver: local
  redis_data:
    driver: local
  consul_data:
    driver: local
  qdrant_data:
    driver: local
  prometheus_data:
    driver: local
  grafana_data:
    driver: local

# === 网络配置 ===
networks:
  lyss-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
```

### 3. 构建脚本

**scripts/build-images.sh**：
```bash
#!/bin/bash

# Docker 镜像构建脚本
set -e

# 获取版本号
VERSION=${VERSION:-$(git rev-parse --short HEAD)}
REGISTRY=${REGISTRY:-"lyss"}

echo "🏗️  构建 Lyss AI Platform Docker 镜像 (版本: $VERSION)"

# 服务列表
SERVICES=("gateway-service" "user-service" "auth-service" "group-service" "credential-service" "billing-service")

# 构建微服务镜像
for SERVICE in "${SERVICES[@]}"; do
    echo "📦 构建 $SERVICE..."
    
    docker build \
        -t "$REGISTRY/$SERVICE:$VERSION" \
        -t "$REGISTRY/$SERVICE:latest" \
        --build-arg VERSION="$VERSION" \
        --build-arg BUILD_TIME="$(date -u +'%Y-%m-%dT%H:%M:%SZ')" \
        --build-arg GIT_COMMIT="$(git rev-parse HEAD)" \
        services/$SERVICE/
    
    echo "✅ $SERVICE 镜像构建完成"
done

# 构建前端镜像
echo "📦 构建前端应用..."
docker build \
    -t "$REGISTRY/frontend:$VERSION" \
    -t "$REGISTRY/frontend:latest" \
    --build-arg VERSION="$VERSION" \
    frontend/

echo "✅ 前端镜像构建完成"

# 显示构建的镜像
echo "📋 构建完成的镜像："
docker images | grep "$REGISTRY" | grep -E "($VERSION|latest)"

# 推送到镜像仓库 (可选)
if [ "$PUSH_IMAGES" = "true" ]; then
    echo "📤 推送镜像到仓库..."
    
    for SERVICE in "${SERVICES[@]}"; do
        docker push "$REGISTRY/$SERVICE:$VERSION"
        docker push "$REGISTRY/$SERVICE:latest"
    done
    
    docker push "$REGISTRY/frontend:$VERSION"
    docker push "$REGISTRY/frontend:latest"
    
    echo "✅ 所有镜像推送完成"
fi

echo "🎉 Docker 镜像构建完成！"
```

---

## Kubernetes 部署

### 1. Namespace 和 ConfigMap

**k8s/namespace.yaml**：
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: lyss-platform
  labels:
    name: lyss-platform
    version: v1
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: lyss-config
  namespace: lyss-platform
data:
  # 数据库配置
  DB_HOST: "postgres-service"
  DB_PORT: "5432"
  DB_USER: "lyss"
  
  # Redis 配置
  REDIS_HOST: "redis-service"
  REDIS_PORT: "6379"
  
  # Consul 配置
  CONSUL_HOST: "consul-service"
  CONSUL_PORT: "8500"
  
  # 应用配置
  ENV: "production"
  LOG_LEVEL: "info"
  METRICS_ENABLED: "true"
```

**k8s/secrets.yaml**：
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: lyss-secrets
  namespace: lyss-platform
type: Opaque
data:
  # Base64 编码的密钥
  DB_PASSWORD: bHlzczEyMw==  # lyss123
  JWT_SECRET: c3VwZXItc2VjcmV0LWp3dC1rZXk=  # super-secret-jwt-key
  ENCRYPTION_KEY: ZW5jcnlwdGlvbi1rZXktMzItY2hhcnM=  # encryption-key-32-chars
  OPENAI_API_KEY: c2steW91ci1vcGVuYWkta2V5  # sk-your-openai-key
```

### 2. 数据库部署

**k8s/postgres.yaml**：
```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
  namespace: lyss-platform
spec:
  serviceName: postgres-service
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:15-alpine
        ports:
        - containerPort: 5432
        env:
        - name: POSTGRES_DB
          value: "lyss_platform"
        - name: POSTGRES_USER
          valueFrom:
            configMapKeyRef:
              name: lyss-config
              key: DB_USER
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: lyss-secrets
              key: DB_PASSWORD
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
        - name: postgres-init
          mountPath: /docker-entrypoint-initdb.d
        resources:
          requests:
            memory: "1Gi"
            cpu: "500m"
          limits:
            memory: "2Gi"
            cpu: "1000m"
        livenessProbe:
          exec:
            command:
            - pg_isready
            - -U
            - lyss
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          exec:
            command:
            - pg_isready
            - -U
            - lyss
          initialDelaySeconds: 10
          periodSeconds: 5
      volumes:
      - name: postgres-init
        configMap:
          name: postgres-init-scripts
  volumeClaimTemplates:
  - metadata:
      name: postgres-storage
    spec:
      accessModes: ["ReadWriteOnce"]
      resources:
        requests:
          storage: 20Gi
---
apiVersion: v1
kind: Service
metadata:
  name: postgres-service
  namespace: lyss-platform
spec:
  selector:
    app: postgres
  ports:
  - port: 5432
    targetPort: 5432
  type: ClusterIP
```

### 3. 微服务部署

**k8s/gateway-service.yaml**：
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gateway-service
  namespace: lyss-platform
  labels:
    app: gateway-service
    version: v1
spec:
  replicas: 3
  selector:
    matchLabels:
      app: gateway-service
  template:
    metadata:
      labels:
        app: gateway-service
        version: v1
    spec:
      containers:
      - name: gateway-service
        image: lyss/gateway-service:latest
        ports:
        - containerPort: 8000
        env:
        - name: ENV
          valueFrom:
            configMapKeyRef:
              name: lyss-config
              key: ENV
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: lyss-config
              key: DB_HOST
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: lyss-secrets
              key: DB_PASSWORD
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: lyss-secrets
              key: JWT_SECRET
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8000
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8000
          initialDelaySeconds: 10
          periodSeconds: 5
        volumeMounts:
        - name: config-volume
          mountPath: /configs
      volumes:
      - name: config-volume
        configMap:
          name: lyss-config
      imagePullSecrets:
      - name: lyss-registry-secret
---
apiVersion: v1
kind: Service
metadata:
  name: gateway-service
  namespace: lyss-platform
  labels:
    app: gateway-service
spec:
  selector:
    app: gateway-service
  ports:
  - name: http
    port: 8000
    targetPort: 8000
  type: ClusterIP
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: gateway-service-hpa
  namespace: lyss-platform
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: gateway-service
  minReplicas: 2
  maxReplicas: 10
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
```

### 4. Ingress 配置

**k8s/ingress.yaml**：
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: lyss-ingress
  namespace: lyss-platform
  annotations:
    kubernetes.io/ingress.class: "nginx"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/proxy-body-size: "50m"
    nginx.ingress.kubernetes.io/cors-allow-origin: "*"
    nginx.ingress.kubernetes.io/cors-allow-methods: "GET, POST, PUT, DELETE, OPTIONS"
    nginx.ingress.kubernetes.io/cors-allow-headers: "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,Authorization"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
spec:
  tls:
  - hosts:
    - api.lyss.ai
    - lyss.ai
    secretName: lyss-tls-secret
  rules:
  - host: api.lyss.ai
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: gateway-service
            port:
              number: 8000
  - host: lyss.ai
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: frontend-service
            port:
              number: 80
```

---

## 监控和日志

### 1. Prometheus 监控配置

**infrastructure/prometheus/prometheus.yml**：
```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "alert_rules.yml"

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - alertmanager:9093

scrape_configs:
  # Prometheus 自监控
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  # 微服务监控
  - job_name: 'gateway-service'
    static_configs:
      - targets: ['gateway-service:8000']
    metrics_path: '/metrics'
    scrape_interval: 30s

  - job_name: 'user-service'
    consul_sd_configs:
      - server: 'consul:8500'
        services: ['user-service']
    relabel_configs:
      - source_labels: [__meta_consul_service]
        target_label: job

  # 基础设施监控
  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres-exporter:9187']

  - job_name: 'redis'
    static_configs:
      - targets: ['redis-exporter:9121']

  # Kubernetes 集群监控
  - job_name: 'kubernetes-pods'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
        action: keep
        regex: true
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
        action: replace
        target_label: __metrics_path__
        regex: (.+)
```

### 2. 告警规则配置

**infrastructure/prometheus/alert_rules.yml**：
```yaml
groups:
- name: lyss-platform-alerts
  rules:
  # 服务可用性告警
  - alert: ServiceDown
    expr: up == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "服务 {{ $labels.job }} 不可用"
      description: "服务 {{ $labels.job }} 在实例 {{ $labels.instance }} 上已宕机超过1分钟"

  # 高CPU使用率告警
  - alert: HighCPUUsage
    expr: (100 - (avg by (instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)) > 80
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "高CPU使用率 ({{ $labels.instance }})"
      description: "CPU使用率已超过80%，当前值: {{ $value }}%"

  # 高内存使用率告警
  - alert: HighMemoryUsage
    expr: (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100 > 85
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "高内存使用率 ({{ $labels.instance }})"
      description: "内存使用率已超过85%，当前值: {{ $value }}%"

  # 数据库连接告警
  - alert: DatabaseConnectionFailure
    expr: pg_up == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "数据库连接失败"
      description: "PostgreSQL 数据库连接失败，服务可能受影响"

  # API响应时间告警
  - alert: HighAPILatency
    expr: histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le, method, route)) > 1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "API响应时间过高"
      description: "{{ $labels.method }} {{ $labels.route }} 的95%分位响应时间超过1秒: {{ $value }}s"

  # 错误率告警
  - alert: HighErrorRate
    expr: sum(rate(http_requests_total{status=~"5.."}[5m])) by (service) / sum(rate(http_requests_total[5m])) by (service) > 0.05
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "高错误率 ({{ $labels.service }})"
      description: "服务 {{ $labels.service }} 的错误率超过5%: {{ $value | humanizePercentage }}"
```

### 3. 日志收集配置

**k8s/fluentd-config.yaml**：
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: fluentd-config
  namespace: lyss-platform
data:
  fluent.conf: |
    <source>
      @type tail
      path /var/log/containers/*lyss*.log
      pos_file /var/log/fluentd-containers.log.pos
      tag kubernetes.*
      format json
      time_format %Y-%m-%dT%H:%M:%S.%NZ
    </source>

    <filter kubernetes.**>
      @type kubernetes_metadata
    </filter>

    <match kubernetes.**>
      @type elasticsearch
      host elasticsearch-service
      port 9200
      index_name lyss-platform
      type_name fluentd
      logstash_format true
      logstash_prefix lyss-platform
      flush_interval 10s
    </match>
```

---

## 部署脚本

### 1. 生产环境部署脚本

**scripts/deploy-production.sh**：
```bash
#!/bin/bash

# 生产环境部署脚本
set -e

# 配置变量
VERSION=${VERSION:-latest}
NAMESPACE=${NAMESPACE:-lyss-platform}
KUBECONFIG=${KUBECONFIG:-~/.kube/config}

echo "🚀 开始部署 Lyss AI Platform 到生产环境"
echo "版本: $VERSION"
echo "命名空间: $NAMESPACE"

# 检查必要工具
command -v kubectl >/dev/null 2>&1 || { echo "❌ kubectl 未安装"; exit 1; }
command -v helm >/dev/null 2>&1 || { echo "❌ helm 未安装"; exit 1; }

# 检查集群连接
echo "🔍 检查 Kubernetes 集群连接..."
kubectl cluster-info || { echo "❌ 无法连接到 Kubernetes 集群"; exit 1; }

# 创建命名空间
echo "📦 创建命名空间..."
kubectl create namespace $NAMESPACE --dry-run=client -o yaml | kubectl apply -f -

# 部署密钥
echo "🔐 部署密钥配置..."
if [ -f ".env.production" ]; then
    kubectl create secret generic lyss-secrets \
        --from-env-file=.env.production \
        --namespace=$NAMESPACE \
        --dry-run=client -o yaml | kubectl apply -f -
else
    echo "⚠️  生产环境配置文件 .env.production 不存在"
    exit 1
fi

# 部署配置映射
echo "⚙️  部署配置映射..."
kubectl apply -f k8s/configmap.yaml

# 部署基础设施服务
echo "🏗️  部署基础设施服务..."
kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/redis.yaml
kubectl apply -f k8s/consul.yaml

# 等待基础设施就绪
echo "⏳ 等待基础设施服务就绪..."
kubectl wait --for=condition=ready pod -l app=postgres --timeout=300s -n $NAMESPACE
kubectl wait --for=condition=ready pod -l app=redis --timeout=300s -n $NAMESPACE
kubectl wait --for=condition=ready pod -l app=consul --timeout=300s -n $NAMESPACE

# 运行数据库迁移
echo "🗄️  运行数据库迁移..."
kubectl apply -f k8s/migration-job.yaml
kubectl wait --for=condition=complete job/db-migration --timeout=600s -n $NAMESPACE

# 部署微服务
echo "🚀 部署微服务..."
for service in gateway-service user-service auth-service group-service credential-service billing-service; do
    echo "部署 $service..."
    envsubst < k8s/${service}.yaml | kubectl apply -f -
done

# 部署前端
echo "🌐 部署前端应用..."
envsubst < k8s/frontend.yaml | kubectl apply -f -

# 部署 Ingress
echo "🌍 部署 Ingress..."
kubectl apply -f k8s/ingress.yaml

# 等待所有服务就绪
echo "⏳ 等待所有服务就绪..."
kubectl wait --for=condition=available deployment --all --timeout=600s -n $NAMESPACE

# 检查部署状态
echo "📊 检查部署状态..."
kubectl get pods -n $NAMESPACE
kubectl get services -n $NAMESPACE
kubectl get ingress -n $NAMESPACE

# 运行健康检查
echo "🏥 运行健康检查..."
./scripts/health-check.sh

echo "🎉 Lyss AI Platform 部署完成！"
echo "📋 访问地址："
echo "  - API: https://api.lyss.ai"
echo "  - 前端: https://lyss.ai"
echo "  - 监控: https://monitor.lyss.ai"

# 显示有用的命令
echo "💡 有用的命令："
echo "  - 查看日志: kubectl logs -f deployment/gateway-service -n $NAMESPACE"
echo "  - 查看状态: kubectl get pods -n $NAMESPACE"
echo "  - 扩缩容: kubectl scale deployment gateway-service --replicas=5 -n $NAMESPACE"
```

### 2. 健康检查脚本

**scripts/health-check.sh**：
```bash
#!/bin/bash

# 健康检查脚本
set -e

NAMESPACE=${NAMESPACE:-lyss-platform}
API_BASE=${API_BASE:-https://api.lyss.ai}

echo "🏥 开始健康检查..."

# 检查 Kubernetes 资源状态
echo "📊 检查 Kubernetes 资源状态..."
kubectl get pods -n $NAMESPACE
kubectl get services -n $NAMESPACE

# 检查服务健康状态
SERVICES=("gateway-service" "user-service" "auth-service" "group-service" "credential-service" "billing-service")

for SERVICE in "${SERVICES[@]}"; do
    echo "🔍 检查 $SERVICE 健康状态..."
    
    # 检查 Pod 状态
    POD_STATUS=$(kubectl get pods -l app=$SERVICE -n $NAMESPACE -o jsonpath='{.items[0].status.phase}')
    if [ "$POD_STATUS" != "Running" ]; then
        echo "❌ $SERVICE Pod 状态异常: $POD_STATUS"
        kubectl describe pod -l app=$SERVICE -n $NAMESPACE
        exit 1
    fi
    
    # 检查健康端点
    if [ "$SERVICE" = "gateway-service" ]; then
        HEALTH_URL="$API_BASE/health"
    else
        SERVICE_IP=$(kubectl get service $SERVICE -n $NAMESPACE -o jsonpath='{.spec.clusterIP}')
        HEALTH_URL="http://$SERVICE_IP:8080/health"
    fi
    
    if curl -f -s "$HEALTH_URL" >/dev/null; then
        echo "✅ $SERVICE 健康检查通过"
    else
        echo "❌ $SERVICE 健康检查失败"
        exit 1
    fi
done

# 检查数据库连接
echo "🗄️  检查数据库连接..."
DB_POD=$(kubectl get pods -l app=postgres -n $NAMESPACE -o jsonpath='{.items[0].metadata.name}')
if kubectl exec $DB_POD -n $NAMESPACE -- pg_isready -U lyss >/dev/null 2>&1; then
    echo "✅ 数据库连接正常"
else
    echo "❌ 数据库连接失败"
    exit 1
fi

# 检查 Redis 连接
echo "🔴 检查 Redis 连接..."
REDIS_POD=$(kubectl get pods -l app=redis -n $NAMESPACE -o jsonpath='{.items[0].metadata.name}')
if kubectl exec $REDIS_POD -n $NAMESPACE -- redis-cli ping | grep -q PONG; then
    echo "✅ Redis 连接正常"
else
    echo "❌ Redis 连接失败"
    exit 1
fi

# API 功能测试
echo "🧪 API 功能测试..."

# 健康检查端点
if curl -f -s "$API_BASE/health" | grep -q "ok"; then
    echo "✅ API 健康检查通过"
else
    echo "❌ API 健康检查失败"
    exit 1
fi

# 用户注册测试
REGISTER_RESPONSE=$(curl -s -X POST "$API_BASE/auth/register" \
    -H "Content-Type: application/json" \
    -d '{
        "username": "healthcheck",
        "email": "healthcheck@lyss.ai",
        "password": "test123",
        "full_name": "Health Check User"
    }')

if echo "$REGISTER_RESPONSE" | grep -q "success\|already exists"; then
    echo "✅ 用户注册接口正常"
else
    echo "❌ 用户注册接口异常"
    echo "$REGISTER_RESPONSE"
    exit 1
fi

# 检查监控指标
echo "📊 检查监控指标..."
if curl -f -s "$API_BASE/metrics" >/dev/null; then
    echo "✅ 监控指标端点正常"
else
    echo "⚠️  监控指标端点异常 (非致命)"
fi

echo "🎉 所有健康检查通过！"
echo "✅ Lyss AI Platform 运行正常"
```

### 3. 回滚脚本

**scripts/rollback.sh**：
```bash
#!/bin/bash

# 回滚脚本
set -e

NAMESPACE=${NAMESPACE:-lyss-platform}
REVISION=${REVISION:-1}

echo "🔄 开始回滚 Lyss AI Platform"
echo "目标版本: $REVISION"

# 确认回滚操作
read -p "确认要回滚到版本 $REVISION 吗? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ 回滚操作已取消"
    exit 1
fi

# 回滚所有部署
DEPLOYMENTS=$(kubectl get deployments -n $NAMESPACE -o jsonpath='{.items[*].metadata.name}')

for DEPLOYMENT in $DEPLOYMENTS; do
    echo "🔄 回滚 $DEPLOYMENT..."
    kubectl rollout undo deployment/$DEPLOYMENT --to-revision=$REVISION -n $NAMESPACE
done

# 等待回滚完成
echo "⏳ 等待回滚完成..."
for DEPLOYMENT in $DEPLOYMENTS; do
    kubectl rollout status deployment/$DEPLOYMENT -n $NAMESPACE --timeout=300s
done

# 验证回滚结果
echo "🔍 验证回滚结果..."
./scripts/health-check.sh

echo "✅ 回滚完成！"
```

---

## 安全配置

### 1. 网络策略

**k8s/network-policy.yaml**：
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: lyss-network-policy
  namespace: lyss-platform
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
  ingress:
  # 允许同命名空间内的通信
  - from:
    - namespaceSelector:
        matchLabels:
          name: lyss-platform
    ports:
    - protocol: TCP
      port: 8000
    - protocol: TCP
      port: 5432
    - protocol: TCP
      port: 6379
  
  # 允许从 ingress 控制器访问
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    ports:
    - protocol: TCP
      port: 8000
      
  egress:
  # 允许 DNS 查询
  - to: []
    ports:
    - protocol: UDP
      port: 53
    - protocol: TCP
      port: 53
      
  # 允许访问外部 AI API
  - to: []
    ports:
    - protocol: TCP
      port: 443
```

### 2. Pod 安全策略

**k8s/pod-security-policy.yaml**：
```yaml
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: lyss-psp
spec:
  # 不允许特权容器
  privileged: false
  
  # 不允许使用主机网络
  hostNetwork: false
  hostIPC: false
  hostPID: false
  
  # 运行用户限制
  runAsUser:
    rule: MustRunAsNonRoot
  
  # 文件系统权限
  fsGroup:
    rule: RunAsAny
    
  # 卷类型限制
  volumes:
  - 'configMap'
  - 'emptyDir'
  - 'projected'
  - 'secret'
  - 'downwardAPI'
  - 'persistentVolumeClaim'
  
  # 禁用特权升级
  allowPrivilegeEscalation: false
  
  # 安全上下文
  securityContext:
    readOnlyRootFilesystem: true
    runAsNonRoot: true
    runAsUser: 1000
```

---

## 优化和调优

### 1. 资源优化

**k8s/resource-quotas.yaml**：
```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: lyss-resource-quota
  namespace: lyss-platform
spec:
  hard:
    requests.cpu: "4"
    requests.memory: 8Gi
    limits.cpu: "8"
    limits.memory: 16Gi
    persistentvolumeclaims: "10"
    services: "20"
    secrets: "10"
    configmaps: "10"
```

### 2. 节点亲和性

**优化示例 - gateway-service 部署**：
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gateway-service
spec:
  template:
    spec:
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: node-type
                operator: In
                values: ["compute"]
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values: ["gateway-service"]
              topologyKey: kubernetes.io/hostname
```

---

## 总结

本 Docker 配置和部署文档提供了 Lyss AI Platform 从开发到生产的完整容器化解决方案：

### 核心特性

1. **多阶段构建**: 优化镜像大小和安全性
2. **Kubernetes 原生**: 完整的 K8s 资源配置
3. **监控完备**: Prometheus + Grafana 监控栈
4. **安全加固**: 网络策略和安全上下文
5. **自动化部署**: 一键部署和健康检查脚本

### 部署流程

```bash
# 构建镜像
./scripts/build-images.sh

# 部署到生产环境  
./scripts/deploy-production.sh

# 健康检查
./scripts/health-check.sh

# 回滚 (如需要)
./scripts/rollback.sh
```

### 监控访问

- **Prometheus**: http://monitor.lyss.ai:9090
- **Grafana**: http://monitor.lyss.ai:3001
- **API**: https://api.lyss.ai
- **前端**: https://lyss.ai

现在 Lyss AI Platform 已具备完整的生产环境部署能力！

---

*本文档将随着平台发展持续更新，确保部署配置的最佳实践。*

**最后更新**: 2025-01-25  
**下次检查**: 2025-02-15