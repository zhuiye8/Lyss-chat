# 网关服务 (Gateway Service) 开发文档

**版本**: 2.0  
**更新时间**: 2025-01-25  
**技术栈**: Go + Kratos + 智能路由 + 会话管理  
**状态**: 已确认

---

## 概述

网关服务是 Lyss AI Platform 的核心入口，负责统一管理所有 AI 模型请求，实现路由转发、负载均衡、认证授权、限流熔断等关键功能。基于 Kratos 微服务框架构建，提供高性能、高可用的 API 网关能力。

### 核心职责

- **统一入口**: 所有AI请求的唯一入口点
- **智能路由**: 根据模型类型和负载情况选择最优供应商
- **会话管理**: 多轮对话上下文维护和状态管理
- **计费控制**: Token使用量计算和成本控制
- **监控告警**: 请求链路追踪和性能监控

---

## 技术架构

### 服务架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                    Gateway Service                          │
├─────────────────────────────────────────────────────────────┤
│  HTTP Layer (Kratos HTTP Server)                          │
│  ├── /api/v1/chat/completions                              │
│  ├── /api/v1/embeddings                                    │
│  └── /api/v1/models                                        │
├─────────────────────────────────────────────────────────────┤
│  Business Layer                                            │
│  ├── Router Manager    (路由选择)                          │
│  ├── Session Manager   (会话管理)                          │
│  ├── Load Balancer     (负载均衡)                          │
│  └── Context Compressor (上下文压缩)                       │
├─────────────────────────────────────────────────────────────┤
│  Infrastructure Layer                                      │
│  ├── Provider Adapters (供应商适配器)                      │
│  ├── Redis Client     (缓存/会话)                          │
│  ├── PostgreSQL       (配置/日志)                          │
│  └── Monitoring       (监控/追踪)                          │
└─────────────────────────────────────────────────────────────┘
```

### 核心组件说明

#### 1. Router Manager (路由管理器)
负责根据请求参数选择最优的AI供应商和模型：

```go
type RouterManager struct {
    providers    map[string]Provider
    loadBalancer *LoadBalancer
    healthCheck  *HealthChecker
    metrics      *Metrics
}

type RoutingDecision struct {
    ProviderID   string    `json:"provider_id"`
    ModelName    string    `json:"model_name"`
    Endpoint     string    `json:"endpoint"`
    Weight       int       `json:"weight"`
    LatencyMS    int       `json:"latency_ms"`
    SelectedAt   time.Time `json:"selected_at"`
}
```

#### 2. Session Manager (会话管理器)
维护用户多轮对话的上下文状态：

```go
type SessionManager struct {
    redis        *redis.Client
    compressor   *ContextCompressor
    mem0Client   *mem0.Client
    maxTokens    int
    ttl          time.Duration
}

type Session struct {
    ID           string                 `json:"id"`
    UserID       string                 `json:"user_id"`
    Messages     []ChatMessage          `json:"messages"`
    Context      string                 `json:"context"`
    TokenCount   int                    `json:"token_count"`
    LastActivity time.Time             `json:"last_activity"`
    Metadata     map[string]interface{} `json:"metadata"`
}
```

#### 3. Provider Adapters (供应商适配器)
统一不同AI供应商的接口差异：

```go
type Provider interface {
    Name() string
    ChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
    Embeddings(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error)
    Models(ctx context.Context) ([]Model, error)
    HealthCheck(ctx context.Context) error
}

// OpenAI适配器实现
type OpenAIProvider struct {
    client   *openai.Client
    config   ProviderConfig
    metrics  *ProviderMetrics
}

// 阿里云适配器实现  
type AliyunProvider struct {
    client   *aliyun.Client
    config   ProviderConfig
    metrics  *ProviderMetrics
}
```

---

## 数据模型设计

### 核心数据结构

#### 1. 路由配置
```go
type RoutingRule struct {
    ID          int64     `gorm:"primaryKey;autoIncrement"`
    ModelName   string    `gorm:"column:model_name;size:100;not null;index"`
    ProviderID  string    `gorm:"column:provider_id;size:50;not null"`
    Priority    int       `gorm:"column:priority;default:1"`
    Weight      int       `gorm:"column:weight;default:100"`
    MaxRPS      int       `gorm:"column:max_rps;default:1000"`
    IsEnabled   bool      `gorm:"column:is_enabled;default:true"`
    CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
```

#### 2. 请求日志
```go
type RequestLog struct {
    ID           int64     `gorm:"primaryKey;autoIncrement"`
    RequestID    string    `gorm:"column:request_id;size:100;not null;uniqueIndex"`
    UserID       string    `gorm:"column:user_id;size:50;not null;index"`
    SessionID    string    `gorm:"column:session_id;size:100;index"`
    ModelName    string    `gorm:"column:model_name;size:100;not null"`
    ProviderID   string    `gorm:"column:provider_id;size:50;not null"`
    InputTokens  int       `gorm:"column:input_tokens;not null"`
    OutputTokens int       `gorm:"column:output_tokens;not null"`
    LatencyMS    int       `gorm:"column:latency_ms;not null"`
    Status       string    `gorm:"column:status;size:20;not null"`
    ErrorMessage string    `gorm:"column:error_message;type:text"`
    CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime;index"`
}
```

### 数据库Schema设计

```sql
-- 路由规则表
CREATE TABLE routing_rules (
    id BIGSERIAL PRIMARY KEY,
    model_name VARCHAR(100) NOT NULL,
    provider_id VARCHAR(50) NOT NULL,
    priority INTEGER DEFAULT 1,
    weight INTEGER DEFAULT 100,
    max_rps INTEGER DEFAULT 1000,
    is_enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_routing_rules_model ON routing_rules(model_name);
CREATE INDEX idx_routing_rules_provider ON routing_rules(provider_id);
CREATE INDEX idx_routing_rules_enabled ON routing_rules(is_enabled);

-- 请求日志表 (按月分区)
CREATE TABLE request_logs (
    id BIGSERIAL PRIMARY KEY,
    request_id VARCHAR(100) NOT NULL UNIQUE,
    user_id VARCHAR(50) NOT NULL,
    session_id VARCHAR(100),
    model_name VARCHAR(100) NOT NULL,
    provider_id VARCHAR(50) NOT NULL,
    input_tokens INTEGER NOT NULL,
    output_tokens INTEGER NOT NULL,
    latency_ms INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) PARTITION BY RANGE (created_at);

-- 创建月度分区表
CREATE TABLE request_logs_2025_01 PARTITION OF request_logs
FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

-- 创建索引
CREATE INDEX idx_request_logs_user_id ON request_logs(user_id);
CREATE INDEX idx_request_logs_session_id ON request_logs(session_id);
CREATE INDEX idx_request_logs_created_at ON request_logs(created_at);
```

---

## API接口设计

### 核心API端点

#### 1. Chat Completions API
兼容OpenAI API格式的聊天完成接口：

```go
// POST /api/v1/chat/completions
type ChatCompletionRequest struct {
    Model       string        `json:"model" binding:"required"`
    Messages    []ChatMessage `json:"messages" binding:"required"`
    Temperature *float32      `json:"temperature,omitempty"`
    MaxTokens   *int         `json:"max_tokens,omitempty"`
    Stream      bool         `json:"stream,omitempty"`
    SessionID   string       `json:"session_id,omitempty"`
    UserID      string       `json:"user_id" binding:"required"`
}

type ChatMessage struct {
    Role    string `json:"role" binding:"required,oneof=system user assistant"`
    Content string `json:"content" binding:"required"`
}

type ChatCompletionResponse struct {
    ID      string                 `json:"id"`
    Object  string                `json:"object"`
    Created int64                 `json:"created"`
    Model   string                `json:"model"`
    Choices []ChatCompletionChoice `json:"choices"`
    Usage   TokenUsage            `json:"usage"`
}
```

#### 2. Models API
获取可用模型列表：

```go
// GET /api/v1/models
type ModelsResponse struct {
    Object string  `json:"object"`
    Data   []Model `json:"data"`
}

type Model struct {
    ID         string   `json:"id"`
    Object     string   `json:"object"`
    Provider   string   `json:"provider"`
    MaxTokens  int      `json:"max_tokens"`
    InputPrice float64  `json:"input_price"`
    OutputPrice float64 `json:"output_price"`
}
```

#### 3. Session Management API
会话管理相关接口：

```go
// GET /api/v1/sessions/{session_id}
type SessionResponse struct {
    ID           string                 `json:"id"`
    UserID       string                 `json:"user_id"`
    MessageCount int                    `json:"message_count"`
    TokenCount   int                    `json:"token_count"`
    LastActivity time.Time             `json:"last_activity"`
    Status       string                `json:"status"`
}

// DELETE /api/v1/sessions/{session_id}
// 清除会话上下文
```

---

## 业务逻辑实现

### 路由选择算法

#### 智能路由逻辑
```go
func (rm *RouterManager) SelectProvider(ctx context.Context, req *RoutingRequest) (*RoutingDecision, error) {
    // 1. 获取模型对应的供应商列表
    providers, err := rm.getAvailableProviders(req.Model)
    if err != nil {
        return nil, err
    }

    // 2. 健康检查过滤
    healthyProviders := rm.filterHealthyProviders(providers)
    if len(healthyProviders) == 0 {
        return nil, ErrNoHealthyProvider
    }

    // 3. 负载均衡选择
    selected := rm.loadBalancer.Select(healthyProviders, req)
    
    // 4. 记录路由决策
    decision := &RoutingDecision{
        ProviderID: selected.ID,
        ModelName:  req.Model,
        Endpoint:   selected.Endpoint,
        Weight:     selected.Weight,
        SelectedAt: time.Now(),
    }

    // 5. 更新路由统计
    rm.metrics.RecordRouting(decision)
    
    return decision, nil
}
```

#### 负载均衡策略
```go
type LoadBalancerStrategy interface {
    Select(providers []Provider, req *RoutingRequest) Provider
}

// 加权轮询策略
type WeightedRoundRobin struct {
    counters map[string]int
    mutex    sync.RWMutex
}

func (w *WeightedRoundRobin) Select(providers []Provider, req *RoutingRequest) Provider {
    w.mutex.Lock()
    defer w.mutex.Unlock()

    totalWeight := 0
    for _, p := range providers {
        totalWeight += p.Weight
    }

    // 加权轮询算法实现
    for _, p := range providers {
        current := w.counters[p.ID]
        current += p.Weight
        w.counters[p.ID] = current
        
        if current >= totalWeight {
            w.counters[p.ID] = 0
            return p
        }
    }

    return providers[0] // 默认返回第一个
}
```

### 会话管理实现

#### 上下文压缩策略
```go
func (sm *SessionManager) CompressContext(ctx context.Context, session *Session) error {
    if session.TokenCount < sm.maxTokens {
        return nil // 无需压缩
    }

    // 1. 提取最近的重要消息
    recentMessages := session.Messages[len(session.Messages)-5:]
    
    // 2. 使用AI模型生成摘要
    summary, err := sm.generateSummary(ctx, session.Messages[:len(session.Messages)-5])
    if err != nil {
        return err
    }

    // 3. 构建新的上下文
    compressedMessages := []ChatMessage{
        {Role: "system", Content: fmt.Sprintf("Previous conversation summary: %s", summary)},
    }
    compressedMessages = append(compressedMessages, recentMessages...)

    // 4. 更新会话
    session.Messages = compressedMessages
    session.Context = summary
    session.TokenCount = sm.calculateTokens(compressedMessages)

    // 5. 保存到Redis
    return sm.saveSession(ctx, session)
}
```

#### 会话状态管理
```go
func (sm *SessionManager) GetOrCreateSession(ctx context.Context, userID, sessionID string) (*Session, error) {
    // 1. 尝试从Redis获取
    if sessionID != "" {
        if session, err := sm.getFromRedis(ctx, sessionID); err == nil {
            return session, nil
        }
    }

    // 2. 创建新会话
    session := &Session{
        ID:           generateSessionID(userID),
        UserID:       userID,
        Messages:     []ChatMessage{},
        TokenCount:   0,
        LastActivity: time.Now(),
        Metadata:     make(map[string]interface{}),
    }

    // 3. 保存到Redis
    if err := sm.saveSession(ctx, session); err != nil {
        return nil, err
    }

    return session, nil
}
```

### 供应商适配器实现

#### OpenAI适配器
```go
type OpenAIProvider struct {
    client  *openai.Client
    config  ProviderConfig
    metrics *ProviderMetrics
}

func (p *OpenAIProvider) ChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
    start := time.Now()
    defer func() {
        p.metrics.RecordLatency(time.Since(start))
    }()

    // 转换请求格式
    openaiReq := openai.ChatCompletionRequest{
        Model:       req.Model,
        Messages:    convertMessages(req.Messages),
        Temperature: req.Temperature,
        MaxTokens:   req.MaxTokens,
        Stream:      req.Stream,
    }

    // 调用OpenAI API
    resp, err := p.client.CreateChatCompletion(ctx, openaiReq)
    if err != nil {
        p.metrics.RecordError(err)
        return nil, fmt.Errorf("openai api error: %w", err)
    }

    // 转换响应格式
    return convertResponse(resp), nil
}
```

#### 阿里云适配器
```go
type AliyunProvider struct {
    client  *dashscope.Client
    config  ProviderConfig
    metrics *ProviderMetrics
}

func (p *AliyunProvider) ChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
    // 阿里云DashScope API调用实现
    dashReq := &dashscope.ChatRequest{
        Model: mapModelName(req.Model), // 模型名称映射
        Messages: convertToDashScopeMessages(req.Messages),
        Parameters: dashscope.Parameters{
            Temperature: req.Temperature,
            MaxTokens:   req.MaxTokens,
        },
    }

    resp, err := p.client.Chat(ctx, dashReq)
    if err != nil {
        return nil, fmt.Errorf("dashscope api error: %w", err)
    }

    return convertFromDashScopeResponse(resp), nil
}
```

---

## 服务配置

### 配置文件结构
```yaml
# config/gateway.yaml
server:
  http:
    addr: "0.0.0.0:8000"
    timeout: 30s
  grpc:
    addr: "0.0.0.0:9000"
    timeout: 30s

data:
  database:
    driver: postgres
    source: "host=localhost user=lyss password=lyss123 dbname=lyss_gateway port=5432 sslmode=disable"
  redis:
    addr: "localhost:6379"
    password: ""
    db: 0
    dial_timeout: 5s
    read_timeout: 3s
    write_timeout: 3s

routing:
  default_provider: "openai"
  max_retries: 3
  timeout: 30s
  load_balancer: "weighted_round_robin"

session:
  max_tokens: 4000
  ttl: 3600s
  compression_threshold: 3000
  summary_model: "gpt-3.5-turbo"

providers:
  openai:
    api_key: "${OPENAI_API_KEY}"
    base_url: "https://api.openai.com/v1"
    timeout: 30s
    max_rps: 100
  
  aliyun:
    api_key: "${ALIYUN_API_KEY}"
    base_url: "https://dashscope.aliyuncs.com/api/v1"
    timeout: 30s
    max_rps: 200

monitoring:
  prometheus:
    enabled: true
    addr: ":9090"
  jaeger:
    endpoint: "http://localhost:14268/api/traces"
```

### 环境变量配置
```bash
# .env
GATEWAY_HTTP_ADDR=:8000
GATEWAY_GRPC_ADDR=:9000

DB_HOST=localhost
DB_PORT=5432
DB_USER=lyss
DB_PASSWORD=lyss123
DB_NAME=lyss_gateway

REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

OPENAI_API_KEY=sk-...
ALIYUN_API_KEY=sk-...

LOG_LEVEL=info
LOG_FORMAT=json
```

---

## 部署配置

### Docker配置
```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway cmd/gateway/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/gateway .
COPY --from=builder /app/config ./config

EXPOSE 8000 9000
CMD ["./gateway", "-config", "config/gateway.yaml"]
```

### Docker Compose配置
```yaml
# docker-compose.yml
version: '3.8'

services:
  gateway:
    build: .
    ports:
      - "8000:8000"
      - "9000:9000"
    environment:
      - OPENAI_API_KEY=${OPENAI_API_KEY}
      - ALIYUN_API_KEY=${ALIYUN_API_KEY}
    depends_on:
      - postgres
      - redis
    volumes:
      - ./config:/app/config
    networks:
      - lyss-network

networks:
  lyss-network:
    driver: bridge
```

### Kubernetes配置
```yaml
# k8s/gateway-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gateway-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: gateway-service
  template:
    metadata:
      labels:
        app: gateway-service
    spec:
      containers:
      - name: gateway
        image: lyss/gateway:latest
        ports:
        - containerPort: 8000
        - containerPort: 9000
        env:
        - name: OPENAI_API_KEY
          valueFrom:
            secretKeyRef:
              name: api-secrets
              key: openai-key
        resources:
          requests:
            memory: "256Mi"
            cpu: "200m"
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
          initialDelaySeconds: 5
          periodSeconds: 5
```

---

## 监控和运维

### 健康检查端点
```go
// GET /health - 服务健康检查
func (s *GatewayService) HealthCheck(ctx context.Context, req *emptypb.Empty) (*HealthResponse, error) {
    checks := map[string]bool{
        "database": s.db.Ping() == nil,
        "redis":    s.redis.Ping(ctx).Err() == nil,
        "providers": s.checkProviders(ctx),
    }

    healthy := true
    for _, ok := range checks {
        if !ok {
            healthy = false
            break
        }
    }

    return &HealthResponse{
        Status: map[string]bool(checks),
        Healthy: healthy,
        Timestamp: time.Now().Unix(),
    }, nil
}
```

### Prometheus指标
```go
var (
    requestTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "gateway_requests_total",
            Help: "Total number of requests processed",
        },
        []string{"provider", "model", "status"},
    )

    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "gateway_request_duration_seconds",
            Help: "Request duration in seconds",
        },
        []string{"provider", "model"},
    )

    tokenUsage = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "gateway_tokens_used_total",
            Help: "Total number of tokens used",
        },
        []string{"provider", "model", "type"},
    )
)
```

### 链路追踪
```go
func (s *GatewayService) ChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
    // 创建追踪span
    span, ctx := opentracing.StartSpanFromContext(ctx, "gateway.chat_completion")
    defer span.Finish()

    // 添加追踪标签
    span.SetTag("user_id", req.UserID)
    span.SetTag("model", req.Model)
    span.SetTag("session_id", req.SessionID)

    // 业务逻辑处理...
    response, err := s.processChatRequest(ctx, req)
    
    if err != nil {
        span.SetTag("error", true)
        span.LogFields(log.Error(err))
    }

    return response, err
}
```

---

## 测试策略

### 单元测试
```go
func TestRouterManager_SelectProvider(t *testing.T) {
    rm := &RouterManager{
        providers: map[string]Provider{
            "openai": &MockProvider{ID: "openai", Weight: 100},
            "aliyun": &MockProvider{ID: "aliyun", Weight: 50},
        },
        loadBalancer: NewWeightedRoundRobin(),
    }

    req := &RoutingRequest{Model: "gpt-3.5-turbo"}
    decision, err := rm.SelectProvider(context.Background(), req)

    assert.NoError(t, err)
    assert.NotEmpty(t, decision.ProviderID)
    assert.Equal(t, "gpt-3.5-turbo", decision.ModelName)
}
```

### 集成测试
```go
func TestGatewayService_ChatCompletion_Integration(t *testing.T) {
    // 启动测试服务器
    srv := setupTestServer(t)
    defer srv.Shutdown()

    // 创建测试请求
    reqBody := ChatCompletionRequest{
        Model: "gpt-3.5-turbo",
        Messages: []ChatMessage{
            {Role: "user", Content: "Hello, world!"},
        },
        UserID: "test_user",
    }

    // 发送HTTP请求
    resp, err := http.Post(srv.URL+"/api/v1/chat/completions", 
        "application/json", bytes.NewBuffer(marshal(reqBody)))
    
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    // 验证响应
    var chatResp ChatCompletionResponse
    json.NewDecoder(resp.Body).Decode(&chatResp)
    assert.NotEmpty(t, chatResp.Choices)
    assert.Greater(t, chatResp.Usage.TotalTokens, 0)
}
```

### 性能测试
```bash
# 使用wrk进行压力测试
wrk -t12 -c400 -d30s --script=load_test.lua http://localhost:8000/api/v1/chat/completions
```

---

## 最佳实践

### 错误处理
- 统一错误码和错误信息格式
- 区分业务错误和系统错误
- 实现优雅降级和熔断机制
- 详细的错误日志和监控告警

### 性能优化
- 连接池管理和复用
- 请求缓存和结果缓存
- 异步处理和批量操作
- 资源限制和背压机制

### 安全考虑
- API密钥安全存储和轮换
- 请求限流和DDoS防护
- 输入验证和SQL注入防护
- 敏感数据加密存储

### 可观测性
- 结构化日志和日志聚合
- 全链路追踪和性能分析
- 业务指标和技术指标监控
- 告警规则和故障恢复

---

## 开发指南

### 本地开发环境
```bash
# 1. 启动基础设施
docker-compose up -d postgres redis consul

# 2. 安装依赖
go mod tidy

# 3. 运行数据库迁移
make migrate

# 4. 启动服务
make run-gateway

# 5. 运行测试
make test
```

### 代码结构
```
gateway-service/
├── cmd/
│   └── gateway/
│       └── main.go
├── internal/
│   ├── biz/           # 业务逻辑
│   ├── data/          # 数据访问
│   ├── server/        # 服务器配置  
│   └── service/       # 服务接口
├── api/
│   └── gateway/       # Protobuf定义
├── config/
│   └── gateway.yaml   # 配置文件
├── deployments/       # 部署文件
└── test/             # 测试文件
```

---

*本文档详细描述了网关服务的设计和实现方案，为AI开发者提供完整的开发指导。*

**最后更新**: 2025-01-25  
**下次检查**: 2025-01-26