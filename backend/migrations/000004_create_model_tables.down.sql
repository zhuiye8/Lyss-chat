-- 删除索引
DROP INDEX IF EXISTS idx_providers_tenant_id;
DROP INDEX IF EXISTS idx_providers_status;
DROP INDEX IF EXISTS idx_models_provider_id;
DROP INDEX IF EXISTS idx_models_status;
DROP INDEX IF EXISTS idx_api_keys_tenant_id;
DROP INDEX IF EXISTS idx_api_keys_provider_id;
DROP INDEX IF EXISTS idx_user_models_user_id;
DROP INDEX IF EXISTS idx_user_models_model_id;

-- 删除表
DROP TABLE IF EXISTS user_models;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS models;
DROP TABLE IF EXISTS providers;
