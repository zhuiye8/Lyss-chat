-- 创建 providers 表
CREATE TABLE IF NOT EXISTS providers (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    base_url VARCHAR(1024),
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

-- 创建 models 表
CREATE TABLE IF NOT EXISTS models (
    id UUID PRIMARY KEY,
    provider_id UUID NOT NULL REFERENCES providers(id) ON DELETE CASCADE,
    model_id VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    capabilities JSONB,
    parameters JSONB,
    max_tokens INTEGER,
    token_price_input DECIMAL(10, 6),
    token_price_output DECIMAL(10, 6),
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(provider_id, model_id)
);

-- 创建 api_keys 表
CREATE TABLE IF NOT EXISTS api_keys (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    provider_id UUID NOT NULL REFERENCES providers(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    encrypted_key TEXT NOT NULL,
    iv VARCHAR(255) NOT NULL,
    salt VARCHAR(255) NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 创建 user_models 表
CREATE TABLE IF NOT EXISTS user_models (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    model_id UUID NOT NULL REFERENCES models(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, model_id)
);

-- 创建索引
CREATE INDEX idx_providers_tenant_id ON providers(tenant_id);
CREATE INDEX idx_providers_status ON providers(status);
CREATE INDEX idx_models_provider_id ON models(provider_id);
CREATE INDEX idx_models_status ON models(status);
CREATE INDEX idx_api_keys_tenant_id ON api_keys(tenant_id);
CREATE INDEX idx_api_keys_provider_id ON api_keys(provider_id);
CREATE INDEX idx_user_models_user_id ON user_models(user_id);
CREATE INDEX idx_user_models_model_id ON user_models(model_id);
