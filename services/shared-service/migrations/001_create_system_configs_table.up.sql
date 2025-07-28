-- 共享服务数据库迁移 - 创建系统配置表
-- 迁移版本: 001
-- 创建时间: 2025-07-28
-- 数据库: lyss_platform

BEGIN;

-- 系统配置表
CREATE TABLE system_configs (
    id BIGSERIAL PRIMARY KEY,
    config_key VARCHAR(100) UNIQUE NOT NULL,
    config_value JSONB,
    description TEXT,
    is_encrypted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引设计
CREATE INDEX idx_system_configs_config_key ON system_configs(config_key);

-- 添加表注释
COMMENT ON TABLE system_configs IS '系统配置表';
COMMENT ON COLUMN system_configs.id IS '配置记录唯一标识符';
COMMENT ON COLUMN system_configs.config_key IS '配置键名';
COMMENT ON COLUMN system_configs.config_value IS '配置值（JSON格式）';
COMMENT ON COLUMN system_configs.description IS '配置描述';
COMMENT ON COLUMN system_configs.is_encrypted IS '是否加密存储';
COMMENT ON COLUMN system_configs.created_at IS '创建时间';
COMMENT ON COLUMN system_configs.updated_at IS '更新时间';

COMMIT;