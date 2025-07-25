# 计费服务 (Billing Service) 开发文档

**版本**: 2.0  
**更新时间**: 2025-01-25  
**技术栈**: Go + Kratos + ClickHouse + 实时计费  
**状态**: 已确认

---

## 概述

计费服务是 Lyss AI Platform 的核心财务模块，负责精确计算和记录所有 AI 模型调用的成本。基于 Kratos 微服务框架构建，提供实时计费、成本分析、配额管理等功能，确保平台的商业可持续性。

### 核心职责

- **Token计费**: 精确计算输入/输出Token成本
- **实时扣费**: 即时从用户/群组余额扣除费用
- **成本统计**: 提供详细的使用报告和成本分析
- **配额管理**: 用户和群组的使用限额控制
- **财务对账**: 与供应商账单核对和成本分摊

---

## 技术架构

### 服务架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                    Billing Service                         │
├─────────────────────────────────────────────────────────────┤
│  HTTP Layer (Gin Router)                                   │
│  ├── /api/v1/billing/usage                                 │
│  ├── /api/v1/billing/quotas                                │
│  └── /api/v1/billing/reports                               │
├─────────────────────────────────────────────────────────────┤
│  Business Layer                                            │
│  ├── Billing Engine    (计费引擎)                          │
│  ├── Quota Manager     (配额管理)                          │
│  ├── Cost Calculator   (成本计算)                          │
│  └── Report Generator  (报告生成)                          │
├─────────────────────────────────────────────────────────────┤
│  Infrastructure Layer                                      │
│  ├── PostgreSQL       (计费数据)                           │
│  ├── Redis           (缓存配额)                            │
│  ├── ClickHouse      (分析数据)                            │
│  └── Message Queue   (异步处理)                            │
└─────────────────────────────────────────────────────────────┘
```

### 核心组件说明

#### 1. Billing Engine (计费引擎)
负责实时计费和扣费逻辑：

```go
type BillingEngine struct {
    pricingRepo  PricingRepository
    accountRepo  AccountRepository
    quotaManager *QuotaManager
    calculator   *CostCalculator
    eventBus     EventBus
}

type BillingEvent struct {
    ID           string    `json:"id"`
    UserID       string    `json:"user_id"`
    GroupID      string    `json:"group_id,omitempty"`
    SessionID    string    `json:"session_id"`
    RequestID    string    `json:"request_id"`
    ModelName    string    `json:"model_name"`
    ProviderID   string    `json:"provider_id"`
    InputTokens  int       `json:"input_tokens"`
    OutputTokens int       `json:"output_tokens"`
    InputCost    float64   `json:"input_cost"`
    OutputCost   float64   `json:"output_cost"`
    TotalCost    float64   `json:"total_cost"`
    Timestamp    time.Time `json:"timestamp"`
}
```

#### 2. Cost Calculator (成本计算器)
根据不同模型和供应商计算准确成本：

```go
type CostCalculator struct {
    pricingCache map[string]*ModelPricing
    exchangeRate *ExchangeRate
    markup       float64
}

type ModelPricing struct {
    ModelName     string  `json:"model_name"`
    ProviderID    string  `json:"provider_id"`
    InputPrice    float64 `json:"input_price"`    // 每1K token价格
    OutputPrice   float64 `json:"output_price"`   // 每1K token价格
    Currency      string  `json:"currency"`       // 币种
    EffectiveFrom time.Time `json:"effective_from"`
    EffectiveTo   *time.Time `json:"effective_to,omitempty"`
}

type CostResult struct {
    InputCost    float64 `json:"input_cost"`
    OutputCost   float64 `json:"output_cost"`
    TotalCost    float64 `json:"total_cost"`
    Currency     string  `json:"currency"`
    ExchangeRate float64 `json:"exchange_rate,omitempty"`
}
```

#### 3. Quota Manager (配额管理器)
管理用户和群组的使用配额：

```go
type QuotaManager struct {
    redis       *redis.Client
    quotaRepo   QuotaRepository
    alertSender AlertSender
}

type UserQuota struct {
    UserID          string    `json:"user_id"`
    DailyLimit      float64   `json:"daily_limit"`      // 日限额
    MonthlyLimit    float64   `json:"monthly_limit"`    // 月限额
    DailyUsed       float64   `json:"daily_used"`       // 日已用
    MonthlyUsed     float64   `json:"monthly_used"`     // 月已用
    Balance         float64   `json:"balance"`          // 账户余额
    LastResetDaily  time.Time `json:"last_reset_daily"`
    LastResetMonthly time.Time `json:"last_reset_monthly"`
    IsBlocked       bool      `json:"is_blocked"`
}

type GroupQuota struct {
    GroupID         string    `json:"group_id"`
    TotalLimit      float64   `json:"total_limit"`      // 群组总限额
    TotalUsed       float64   `json:"total_used"`       // 群组已用
    MemberLimits    map[string]float64 `json:"member_limits"` // 成员限额
    LastReset       time.Time `json:"last_reset"`
    IsActive        bool      `json:"is_active"`
}
```

---

## 数据模型设计

### 核心数据结构

#### 1. 使用记录表
```go
type UsageRecord struct {
    ID           int64     `gorm:"primaryKey;autoIncrement"`
    RequestID    string    `gorm:"column:request_id;size:100;not null;uniqueIndex"`
    UserID       string    `gorm:"column:user_id;size:50;not null;index"`
    GroupID      string    `gorm:"column:group_id;size:50;index"`
    SessionID    string    `gorm:"column:session_id;size:100;index"`
    ModelName    string    `gorm:"column:model_name;size:100;not null"`
    ProviderID   string    `gorm:"column:provider_id;size:50;not null"`
    InputTokens  int       `gorm:"column:input_tokens;not null"`
    OutputTokens int       `gorm:"column:output_tokens;not null"`
    InputCost    float64   `gorm:"column:input_cost;type:decimal(10,6);not null"`
    OutputCost   float64   `gorm:"column:output_cost;type:decimal(10,6);not null"`
    TotalCost    float64   `gorm:"column:total_cost;type:decimal(10,6);not null"`
    Currency     string    `gorm:"column:currency;size:10;default:'USD'"`
    BilledAt     time.Time `gorm:"column:billed_at;autoCreateTime;index"`
    CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}
```

#### 2. 账户余额表
```go
type Account struct {
    ID          int64     `gorm:"primaryKey;autoIncrement"`
    UserID      string    `gorm:"column:user_id;size:50;not null;uniqueIndex"`
    Balance     float64   `gorm:"column:balance;type:decimal(10,2);not null;default:0"`
    CreditLimit float64   `gorm:"column:credit_limit;type:decimal(10,2);default:0"`
    Currency    string    `gorm:"column:currency;size:10;default:'USD'"`
    Status      string    `gorm:"column:status;size:20;default:'active'"`
    CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
```

#### 3. 模型定价表
```go
type ModelPrice struct {
    ID            int64      `gorm:"primaryKey;autoIncrement"`
    ModelName     string     `gorm:"column:model_name;size:100;not null"`
    ProviderID    string     `gorm:"column:provider_id;size:50;not null"`
    InputPrice    float64    `gorm:"column:input_price;type:decimal(10,6);not null"`
    OutputPrice   float64    `gorm:"column:output_price;type:decimal(10,6);not null"`
    Currency      string     `gorm:"column:currency;size:10;default:'USD'"`
    EffectiveFrom time.Time  `gorm:"column:effective_from;not null"`
    EffectiveTo   *time.Time `gorm:"column:effective_to"`
    CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime"`
}
```

### 数据库Schema设计

```sql
-- 使用记录表 (按月分区)
CREATE TABLE usage_records (
    id BIGSERIAL PRIMARY KEY,
    request_id VARCHAR(100) NOT NULL UNIQUE,
    user_id VARCHAR(50) NOT NULL,
    group_id VARCHAR(50),
    session_id VARCHAR(100),
    model_name VARCHAR(100) NOT NULL,
    provider_id VARCHAR(50) NOT NULL,
    input_tokens INTEGER NOT NULL,
    output_tokens INTEGER NOT NULL,
    input_cost DECIMAL(10,6) NOT NULL,
    output_cost DECIMAL(10,6) NOT NULL,
    total_cost DECIMAL(10,6) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',
    billed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) PARTITION BY RANGE (billed_at);

-- 创建月度分区
CREATE TABLE usage_records_2025_01 PARTITION OF usage_records
FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

-- 账户余额表
CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL UNIQUE,
    balance DECIMAL(10,2) NOT NULL DEFAULT 0,
    credit_limit DECIMAL(10,2) DEFAULT 0,
    currency VARCHAR(10) DEFAULT 'USD',
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 模型定价表
CREATE TABLE model_prices (
    id BIGSERIAL PRIMARY KEY,
    model_name VARCHAR(100) NOT NULL,
    provider_id VARCHAR(50) NOT NULL,
    input_price DECIMAL(10,6) NOT NULL,
    output_price DECIMAL(10,6) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',
    effective_from TIMESTAMP NOT NULL,
    effective_to TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(model_name, provider_id, effective_from)
);

-- 用户配额表
CREATE TABLE user_quotas (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL UNIQUE,
    daily_limit DECIMAL(10,2) DEFAULT 100.00,
    monthly_limit DECIMAL(10,2) DEFAULT 1000.00,
    daily_used DECIMAL(10,2) DEFAULT 0,
    monthly_used DECIMAL(10,2) DEFAULT 0,
    last_reset_daily DATE DEFAULT CURRENT_DATE,
    last_reset_monthly DATE DEFAULT DATE_TRUNC('month', CURRENT_DATE),
    is_blocked BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_usage_records_user_id ON usage_records(user_id);
CREATE INDEX idx_usage_records_billed_at ON usage_records(billed_at);
CREATE INDEX idx_usage_records_group_id ON usage_records(group_id);
CREATE INDEX idx_model_prices_lookup ON model_prices(model_name, provider_id, effective_from);
CREATE INDEX idx_accounts_status ON accounts(status);
```

---

## API接口设计

### 核心API端点

#### 1. 计费处理API
```go
// POST /api/v1/billing/charge
type ChargeRequest struct {
    RequestID    string  `json:"request_id" binding:"required"`
    UserID       string  `json:"user_id" binding:"required"`
    GroupID      string  `json:"group_id,omitempty"`
    SessionID    string  `json:"session_id,omitempty"`
    ModelName    string  `json:"model_name" binding:"required"`
    ProviderID   string  `json:"provider_id" binding:"required"`
    InputTokens  int     `json:"input_tokens" binding:"min=0"`
    OutputTokens int     `json:"output_tokens" binding:"min=0"`
}

type ChargeResponse struct {
    RequestID   string  `json:"request_id"`
    InputCost   float64 `json:"input_cost"`
    OutputCost  float64 `json:"output_cost"`
    TotalCost   float64 `json:"total_cost"`
    Currency    string  `json:"currency"`
    Timestamp   int64   `json:"timestamp"`
    Success     bool    `json:"success"`
    Balance     float64 `json:"balance"`
}
```

#### 2. 使用统计API
```go
// GET /api/v1/billing/usage
type UsageQuery struct {
    UserID    string    `form:"user_id"`
    GroupID   string    `form:"group_id"`
    StartDate time.Time `form:"start_date"`
    EndDate   time.Time `form:"end_date"`
    ModelName string    `form:"model_name"`
    Page      int       `form:"page" binding:"min=1"`
    Size      int       `form:"size" binding:"min=1,max=100"`
}

type UsageResponse struct {
    Records    []UsageRecord `json:"records"`
    Summary    UsageSummary  `json:"summary"`
    Pagination Pagination    `json:"pagination"`
}

type UsageSummary struct {
    TotalRequests int     `json:"total_requests"`
    TotalTokens   int     `json:"total_tokens"`
    TotalCost     float64 `json:"total_cost"`
    ByModel       map[string]ModelUsage `json:"by_model"`
    ByProvider    map[string]ProviderUsage `json:"by_provider"`
}
```

#### 3. 配额管理API
```go
// GET /api/v1/billing/quotas/{user_id}
type QuotaResponse struct {
    UserID          string  `json:"user_id"`
    DailyLimit      float64 `json:"daily_limit"`
    MonthlyLimit    float64 `json:"monthly_limit"`
    DailyUsed       float64 `json:"daily_used"`
    MonthlyUsed     float64 `json:"monthly_used"`
    DailyRemaining  float64 `json:"daily_remaining"`
    MonthlyRemaining float64 `json:"monthly_remaining"`
    Balance         float64 `json:"balance"`
    IsBlocked       bool    `json:"is_blocked"`
}

// PUT /api/v1/billing/quotas/{user_id}
type UpdateQuotaRequest struct {
    DailyLimit   *float64 `json:"daily_limit,omitempty"`
    MonthlyLimit *float64 `json:"monthly_limit,omitempty"`
    IsBlocked    *bool    `json:"is_blocked,omitempty"`
}
```

#### 4. 账户管理API
```go
// GET /api/v1/billing/accounts/{user_id}
type AccountResponse struct {
    UserID      string  `json:"user_id"`
    Balance     float64 `json:"balance"`
    CreditLimit float64 `json:"credit_limit"`
    Currency    string  `json:"currency"`
    Status      string  `json:"status"`
}

// POST /api/v1/billing/accounts/{user_id}/recharge
type RechargeRequest struct {
    Amount      float64 `json:"amount" binding:"required,gt=0"`
    PaymentType string  `json:"payment_type" binding:"required"`
    OrderID     string  `json:"order_id,omitempty"`
}
```

---

## 业务逻辑实现

### 实时计费处理

#### 计费主流程
```go
func (be *BillingEngine) ProcessCharge(ctx context.Context, req *ChargeRequest) (*ChargeResponse, error) {
    // 1. 参数验证
    if err := be.validateChargeRequest(req); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }

    // 2. 检查配额
    if err := be.quotaManager.CheckQuota(ctx, req.UserID, req.GroupID); err != nil {
        return nil, fmt.Errorf("quota exceeded: %w", err)
    }

    // 3. 计算成本
    cost, err := be.calculator.Calculate(ctx, &CostRequest{
        ModelName:    req.ModelName,
        ProviderID:   req.ProviderID,
        InputTokens:  req.InputTokens,
        OutputTokens: req.OutputTokens,
    })
    if err != nil {
        return nil, fmt.Errorf("cost calculation failed: %w", err)
    }

    // 4. 检查余额
    account, err := be.accountRepo.GetByUserID(ctx, req.UserID)
    if err != nil {
        return nil, fmt.Errorf("account not found: %w", err)
    }

    if account.Balance+account.CreditLimit < cost.TotalCost {
        return nil, ErrInsufficientBalance
    }

    // 5. 执行扣费 (事务)
    if err := be.executeCharge(ctx, req, cost, account); err != nil {
        return nil, fmt.Errorf("charge execution failed: %w", err)
    }

    // 6. 发送事件
    event := &BillingEvent{
        ID:           generateEventID(),
        UserID:       req.UserID,
        GroupID:      req.GroupID,
        RequestID:    req.RequestID,
        ModelName:    req.ModelName,
        ProviderID:   req.ProviderID,
        InputTokens:  req.InputTokens,
        OutputTokens: req.OutputTokens,
        InputCost:    cost.InputCost,
        OutputCost:   cost.OutputCost,
        TotalCost:    cost.TotalCost,
        Timestamp:    time.Now(),
    }
    be.eventBus.Publish(ctx, "billing.charged", event)

    return &ChargeResponse{
        RequestID:  req.RequestID,
        InputCost:  cost.InputCost,
        OutputCost: cost.OutputCost,
        TotalCost:  cost.TotalCost,
        Currency:   cost.Currency,
        Timestamp:  time.Now().Unix(),
        Success:    true,
        Balance:    account.Balance - cost.TotalCost,
    }, nil
}
```

#### 事务性扣费操作
```go
func (be *BillingEngine) executeCharge(ctx context.Context, req *ChargeRequest, cost *CostResult, account *Account) error {
    return be.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // 1. 创建使用记录
        usage := &UsageRecord{
            RequestID:    req.RequestID,
            UserID:       req.UserID,
            GroupID:      req.GroupID,
            SessionID:    req.SessionID,
            ModelName:    req.ModelName,
            ProviderID:   req.ProviderID,
            InputTokens:  req.InputTokens,
            OutputTokens: req.OutputTokens,
            InputCost:    cost.InputCost,
            OutputCost:   cost.OutputCost,
            TotalCost:    cost.TotalCost,
            Currency:     cost.Currency,
        }
        
        if err := tx.Create(usage).Error; err != nil {
            return fmt.Errorf("failed to create usage record: %w", err)
        }

        // 2. 更新账户余额
        if err := tx.Model(account).Update("balance", account.Balance-cost.TotalCost).Error; err != nil {
            return fmt.Errorf("failed to update balance: %w", err)
        }

        // 3. 更新配额使用
        if err := be.quotaManager.UpdateUsage(ctx, tx, req.UserID, req.GroupID, cost.TotalCost); err != nil {
            return fmt.Errorf("failed to update quota: %w", err)
        }

        return nil
    })
}
```

### 成本计算实现

#### 动态定价计算
```go
func (cc *CostCalculator) Calculate(ctx context.Context, req *CostRequest) (*CostResult, error) {
    // 1. 获取模型定价
    pricing, err := cc.getModelPricing(ctx, req.ModelName, req.ProviderID)
    if err != nil {
        return nil, fmt.Errorf("pricing not found: %w", err)
    }

    // 2. 计算基础成本 (按1K token计费)
    inputCost := float64(req.InputTokens) / 1000.0 * pricing.InputPrice
    outputCost := float64(req.OutputTokens) / 1000.0 * pricing.OutputPrice

    // 3. 汇率转换
    if pricing.Currency != "USD" {
        rate, err := cc.exchangeRate.GetRate(pricing.Currency, "USD")
        if err != nil {
            return nil, fmt.Errorf("exchange rate error: %w", err)
        }
        inputCost *= rate
        outputCost *= rate
    }

    // 4. 应用加价
    inputCost *= (1 + cc.markup)
    outputCost *= (1 + cc.markup)

    // 5. 精度处理 (保留6位小数)
    inputCost = math.Round(inputCost*1000000) / 1000000
    outputCost = math.Round(outputCost*1000000) / 1000000
    totalCost := inputCost + outputCost

    return &CostResult{
        InputCost:  inputCost,
        OutputCost: outputCost,
        TotalCost:  totalCost,
        Currency:   "USD",
    }, nil
}
```

#### 定价缓存管理
```go
func (cc *CostCalculator) getModelPricing(ctx context.Context, modelName, providerID string) (*ModelPricing, error) {
    // 1. 检查内存缓存
    cacheKey := fmt.Sprintf("%s:%s", modelName, providerID)
    if pricing, exists := cc.pricingCache[cacheKey]; exists {
        if pricing.EffectiveTo == nil || time.Now().Before(*pricing.EffectiveTo) {
            return pricing, nil
        }
    }

    // 2. 从数据库获取当前有效定价
    var price ModelPrice
    err := cc.db.Where("model_name = ? AND provider_id = ? AND effective_from <= ? AND (effective_to IS NULL OR effective_to > ?)",
        modelName, providerID, time.Now(), time.Now()).
        Order("effective_from DESC").
        First(&price).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrPricingNotFound
        }
        return nil, err
    }

    // 3. 转换并缓存
    pricing := &ModelPricing{
        ModelName:     price.ModelName,
        ProviderID:    price.ProviderID,
        InputPrice:    price.InputPrice,
        OutputPrice:   price.OutputPrice,
        Currency:      price.Currency,
        EffectiveFrom: price.EffectiveFrom,
        EffectiveTo:   price.EffectiveTo,
    }

    cc.pricingCache[cacheKey] = pricing
    return pricing, nil
}
```

### 配额管理实现

#### 配额检查和更新
```go
func (qm *QuotaManager) CheckQuota(ctx context.Context, userID, groupID string) error {
    // 1. 检查用户配额
    quota, err := qm.getUserQuota(ctx, userID)
    if err != nil {
        return err
    }

    if quota.IsBlocked {
        return ErrUserBlocked
    }

    // 2. 检查日限额
    if quota.DailyUsed >= quota.DailyLimit {
        return ErrDailyQuotaExceeded
    }

    // 3. 检查月限额
    if quota.MonthlyUsed >= quota.MonthlyLimit {
        return ErrMonthlyQuotaExceeded
    }

    // 4. 检查群组配额 (如果适用)
    if groupID != "" {
        if err := qm.checkGroupQuota(ctx, groupID, userID); err != nil {
            return err
        }
    }

    return nil
}

func (qm *QuotaManager) UpdateUsage(ctx context.Context, tx *gorm.DB, userID, groupID string, cost float64) error {
    now := time.Now()
    today := now.Format("2006-01-02")
    thisMonth := now.Format("2006-01")

    // 1. 更新用户配额
    updates := map[string]interface{}{
        "daily_used":   gorm.Expr("daily_used + ?", cost),
        "monthly_used": gorm.Expr("monthly_used + ?", cost),
        "updated_at":   now,
    }

    // 检查是否需要重置日配额
    var quota UserQuota
    if err := tx.Where("user_id = ?", userID).First(&quota).Error; err != nil {
        return err
    }

    if quota.LastResetDaily.Format("2006-01-02") != today {
        updates["daily_used"] = cost
        updates["last_reset_daily"] = now
    }

    // 检查是否需要重置月配额
    if quota.LastResetMonthly.Format("2006-01") != thisMonth {
        updates["monthly_used"] = cost
        updates["last_reset_monthly"] = now
    }

    if err := tx.Model(&UserQuota{}).Where("user_id = ?", userID).Updates(updates).Error; err != nil {
        return err
    }

    // 2. 更新群组配额 (如果适用)
    if groupID != "" {
        if err := qm.updateGroupUsage(ctx, tx, groupID, userID, cost); err != nil {
            return err
        }
    }

    // 3. 更新Redis缓存
    qm.updateQuotaCache(ctx, userID, cost)

    return nil
}
```

---

## 报告生成系统

### 使用报告生成
```go
type ReportGenerator struct {
    db         *gorm.DB
    clickhouse *clickhouse.Conn
    cache      *redis.Client
}

func (rg *ReportGenerator) GenerateUsageReport(ctx context.Context, req *ReportRequest) (*UsageReport, error) {
    // 1. 构建查询条件
    query := rg.buildQuery(req)
    
    // 2. 查询聚合数据
    var records []UsageRecord
    if err := query.Find(&records).Error; err != nil {
        return nil, err
    }

    // 3. 数据聚合计算
    summary := rg.calculateSummary(records)
    
    // 4. 生成图表数据
    charts := rg.generateCharts(records, req.TimeRange)

    return &UsageReport{
        Period:    req.Period,
        Summary:   summary,
        Records:   records,
        Charts:    charts,
        Generated: time.Now(),
    }, nil
}

func (rg *ReportGenerator) calculateSummary(records []UsageRecord) *UsageSummary {
    summary := &UsageSummary{
        ByModel:    make(map[string]ModelUsage),
        ByProvider: make(map[string]ProviderUsage),
    }

    for _, record := range records {
        summary.TotalRequests++
        summary.TotalTokens += record.InputTokens + record.OutputTokens
        summary.TotalCost += record.TotalCost

        // 按模型统计
        modelKey := record.ModelName
        modelUsage := summary.ByModel[modelKey]
        modelUsage.Requests++
        modelUsage.Tokens += record.InputTokens + record.OutputTokens
        modelUsage.Cost += record.TotalCost
        summary.ByModel[modelKey] = modelUsage

        // 按供应商统计
        providerKey := record.ProviderID
        providerUsage := summary.ByProvider[providerKey]
        providerUsage.Requests++
        providerUsage.Tokens += record.InputTokens + record.OutputTokens
        providerUsage.Cost += record.TotalCost
        summary.ByProvider[providerKey] = providerUsage
    }

    return summary
}
```

### 实时数据聚合
```go
// 使用ClickHouse进行高性能数据分析
func (rg *ReportGenerator) GetRealTimeMetrics(ctx context.Context, userID string) (*RealTimeMetrics, error) {
    query := `
    SELECT 
        countIf(billed_at >= today()) as today_requests,
        sumIf(total_cost, billed_at >= today()) as today_cost,
        countIf(billed_at >= toStartOfMonth(now())) as month_requests,
        sumIf(total_cost, billed_at >= toStartOfMonth(now())) as month_cost,
        uniqIf(model_name, billed_at >= today()) as models_used_today
    FROM usage_records 
    WHERE user_id = ?
        AND billed_at >= toStartOfMonth(now())
    `

    var metrics RealTimeMetrics
    if err := rg.clickhouse.QueryRow(ctx, query, userID).Scan(
        &metrics.TodayRequests,
        &metrics.TodayCost,
        &metrics.MonthRequests,
        &metrics.MonthCost,
        &metrics.ModelsUsedToday,
    ); err != nil {
        return nil, err
    }

    return &metrics, nil
}
```

---

## 服务配置

### 配置文件结构
```yaml
# config/billing.yaml
server:
  http:
    addr: "0.0.0.0:8003"
    timeout: 30s
  grpc:
    addr: "0.0.0.0:9003"
    timeout: 30s

data:
  database:
    driver: postgres
    source: "host=localhost user=lyss password=lyss123 dbname=lyss_billing port=5432 sslmode=disable"
  redis:
    addr: "localhost:6379"
    password: ""
    db: 2
  clickhouse:
    addr: "localhost:9000"
    database: "lyss_analytics"
    username: "default"
    password: ""

billing:
  currency: "USD"
  markup: 0.20  # 20% 加价
  precision: 6   # 成本精度
  quota_check_interval: "1m"
  
pricing:
  cache_ttl: "1h"
  update_interval: "5m"
  
quota:
  default_daily_limit: 100.00
  default_monthly_limit: 1000.00
  warning_threshold: 0.8  # 80% 时发送警告
  
reports:
  max_records: 10000
  cache_ttl: "10m"
  
exchange_rates:
  provider: "fixer.io"
  api_key: "${FIXER_API_KEY}"
  update_interval: "1h"
```

### 环境变量配置
```bash
# .env
BILLING_HTTP_ADDR=:8003
BILLING_GRPC_ADDR=:9003

# 定价API
FIXER_API_KEY=your_fixer_api_key
OPENAI_PRICING_URL=https://api.openai.com/v1/models
ALIYUN_PRICING_URL=https://dashscope.aliyuncs.com/api/v1/pricing

# ClickHouse 配置
CLICKHOUSE_ADDR=localhost:9000
CLICKHOUSE_DATABASE=lyss_analytics
CLICKHOUSE_USERNAME=default
CLICKHOUSE_PASSWORD=

# 告警配置
ALERT_WEBHOOK_URL=https://hooks.slack.com/services/...
ALERT_EMAIL_SMTP=smtp.gmail.com:587
ALERT_EMAIL_USER=alerts@lyss.ai
ALERT_EMAIL_PASSWORD=your_password
```

---

## 监控和告警

### 关键指标监控
```go
var (
    billingTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "billing_charges_total",
            Help: "Total number of billing charges processed",
        },
        []string{"user_id", "model", "provider", "status"},
    )

    billingAmount = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "billing_amount_total",
            Help: "Total billing amount in USD",
        },
        []string{"user_id", "model", "provider"},
    )

    quotaViolations = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "quota_violations_total", 
            Help: "Total number of quota violations",
        },
        []string{"user_id", "type"}, // type: daily, monthly, balance
    )

    billingLatency = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "billing_processing_duration_seconds",
            Help: "Billing processing duration in seconds",
        },
        []string{"operation"},
    )
)
```

### 智能告警系统
```go
type AlertManager struct {
    rules   []AlertRule
    sender  AlertSender
    history map[string]time.Time // 防止重复告警
}

type AlertRule struct {
    Name      string
    Condition func(ctx context.Context, event *BillingEvent) bool
    Template  string
    Cooldown  time.Duration
}

func (am *AlertManager) setupDefaultRules() {
    am.rules = []AlertRule{
        {
            Name: "quota_80_percent",
            Condition: func(ctx context.Context, event *BillingEvent) bool {
                quota, _ := am.quotaManager.GetUserQuota(ctx, event.UserID)
                return quota.DailyUsed/quota.DailyLimit >= 0.8
            },
            Template: "用户 {{.UserID}} 日配额已使用 {{.UsagePercent}}%",
            Cooldown: time.Hour,
        },
        {
            Name: "high_cost_request",
            Condition: func(ctx context.Context, event *BillingEvent) bool {
                return event.TotalCost > 10.0 // 单次请求超过10美元
            },
            Template: "高成本请求告警: 用户 {{.UserID}} 单次请求花费 ${{.TotalCost}}",
            Cooldown: time.Minute * 5,
        },
        {
            Name: "balance_low",
            Condition: func(ctx context.Context, event *BillingEvent) bool {
                account, _ := am.accountRepo.GetByUserID(ctx, event.UserID)
                return account.Balance < 10.0 // 余额低于10美元
            },
            Template: "余额不足告警: 用户 {{.UserID}} 余额仅剩 ${{.Balance}}",
            Cooldown: time.Hour * 6,
        },
    }
}
```

---

## 数据分析和优化

### 成本分析
```go
type CostAnalyzer struct {
    db         *gorm.DB
    clickhouse *clickhouse.Conn
}

func (ca *CostAnalyzer) AnalyzeCostTrends(ctx context.Context, period string) (*CostAnalysis, error) {
    // 使用ClickHouse进行高性能分析
    query := `
    SELECT 
        toDate(billed_at) as date,
        model_name,
        provider_id,
        sum(total_cost) as daily_cost,
        count(*) as request_count,
        sum(input_tokens + output_tokens) as total_tokens,
        avg(total_cost) as avg_cost_per_request
    FROM usage_records 
    WHERE billed_at >= subtractDays(now(), ?)
    GROUP BY date, model_name, provider_id
    ORDER BY date DESC, daily_cost DESC
    `

    days := map[string]int{
        "week":  7,
        "month": 30,
        "quarter": 90,
    }[period]

    rows, err := ca.clickhouse.Query(ctx, query, days)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var trends []CostTrend
    var totalCost float64
    
    for rows.Next() {
        var trend CostTrend
        if err := rows.Scan(&trend.Date, &trend.ModelName, &trend.ProviderID, 
                          &trend.DailyCost, &trend.RequestCount, &trend.TotalTokens, 
                          &trend.AvgCostPerRequest); err != nil {
            return nil, err
        }
        trends = append(trends, trend)
        totalCost += trend.DailyCost
    }

    return &CostAnalysis{
        Period:     period,
        TotalCost:  totalCost,
        Trends:     trends,
        Generated:  time.Now(),
    }, nil
}
```

### 异常检测
```go
func (ca *CostAnalyzer) DetectAnomalies(ctx context.Context) ([]Anomaly, error) {
    // 检测异常高成本用户
    var anomalies []Anomaly

    // 1. 检测成本异常用户 (今日成本超过历史平均值3倍)
    query := `
    WITH user_daily_avg AS (
        SELECT user_id, avg(daily_cost) as avg_cost
        FROM (
            SELECT user_id, toDate(billed_at) as date, sum(total_cost) as daily_cost
            FROM usage_records 
            WHERE billed_at >= subtractDays(now(), 30)
            GROUP BY user_id, date
        ) GROUP BY user_id
    ),
    today_cost AS (
        SELECT user_id, sum(total_cost) as today_cost
        FROM usage_records
        WHERE toDate(billed_at) = today()
        GROUP BY user_id
    )
    SELECT t.user_id, t.today_cost, a.avg_cost
    FROM today_cost t
    JOIN user_daily_avg a ON t.user_id = a.user_id
    WHERE t.today_cost > a.avg_cost * 3
    `

    rows, err := ca.clickhouse.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var userID string
        var todayCost, avgCost float64
        if err := rows.Scan(&userID, &todayCost, &avgCost); err != nil {
            continue
        }

        anomalies = append(anomalies, Anomaly{
            Type:        "high_cost_user",
            UserID:      userID,
            Description: fmt.Sprintf("用户今日消费 $%.2f，超过历史平均值 $%.2f", todayCost, avgCost),
            Severity:    "warning",
            DetectedAt:  time.Now(),
        })
    }

    return anomalies, nil
}
```

---

## 测试策略

### 单元测试
```go
func TestBillingEngine_ProcessCharge(t *testing.T) {
    // 准备测试数据
    engine := setupTestBillingEngine()
    req := &ChargeRequest{
        RequestID:    "test-req-001",
        UserID:       "user-001",
        ModelName:    "gpt-3.5-turbo",
        ProviderID:   "openai",
        InputTokens:  100,
        OutputTokens: 50,
    }

    // 执行测试
    resp, err := engine.ProcessCharge(context.Background(), req)

    // 验证结果
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, req.RequestID, resp.RequestID)
    assert.Greater(t, resp.TotalCost, 0.0)
    assert.True(t, resp.Success)
}

func TestCostCalculator_Calculate(t *testing.T) {
    calculator := &CostCalculator{
        pricingCache: map[string]*ModelPricing{
            "gpt-3.5-turbo:openai": {
                InputPrice:  0.0015,
                OutputPrice: 0.002,
                Currency:    "USD",
            },
        },
        markup: 0.2,
    }

    req := &CostRequest{
        ModelName:    "gpt-3.5-turbo",
        ProviderID:   "openai", 
        InputTokens:  1000,
        OutputTokens: 500,
    }

    result, err := calculator.Calculate(context.Background(), req)

    assert.NoError(t, err)
    assert.Equal(t, 0.0015*1.2, result.InputCost)   // (1000/1000) * 0.0015 * 1.2
    assert.Equal(t, 0.001*1.2, result.OutputCost)   // (500/1000) * 0.002 * 1.2
    assert.Equal(t, result.InputCost+result.OutputCost, result.TotalCost)
}
```

### 集成测试
```go
func TestBillingService_Integration(t *testing.T) {
    // 启动测试服务
    srv := setupTestBillingService(t)
    defer srv.Shutdown()

    // 准备测试账户
    setupTestAccount(t, "user-001", 100.0)

    // 测试计费请求
    reqBody := ChargeRequest{
        RequestID:    "integration-test-001",
        UserID:       "user-001",
        ModelName:    "gpt-3.5-turbo",
        ProviderID:   "openai",
        InputTokens:  500,
        OutputTokens: 300,
    }

    resp, err := http.Post(srv.URL+"/api/v1/billing/charge", 
        "application/json", bytes.NewBuffer(marshal(reqBody)))

    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    // 验证响应
    var chargeResp ChargeResponse
    json.NewDecoder(resp.Body).Decode(&chargeResp)
    assert.True(t, chargeResp.Success)
    assert.Greater(t, chargeResp.TotalCost, 0.0)

    // 验证余额更新
    account := getTestAccount(t, "user-001")
    expectedBalance := 100.0 - chargeResp.TotalCost
    assert.InDelta(t, expectedBalance, account.Balance, 0.000001)
}
```

---

## 最佳实践

### 数据一致性
- 使用数据库事务确保扣费操作的原子性
- 实现分布式锁防止并发扣费问题
- 定期数据校验和对账
- 异常情况的补偿机制

### 性能优化
- 热点数据Redis缓存
- 定价信息内存缓存
- 批量数据操作和异步处理
- 数据库查询优化和索引设计

### 安全考虑
- 敏感信息加密存储
- API接口鉴权和授权
- 防刷和恶意使用检测
- 审计日志完整记录

### 可扩展性
- 水平扩展和负载均衡
- 数据分片和分区策略
- 事件驱动架构
- 微服务间松耦合设计

---

*本文档详细描述了计费服务的完整设计和实现方案，确保平台的财务管理准确可靠。*

**最后更新**: 2025-01-25  
**下次检查**: 2025-01-26