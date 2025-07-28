-- 网关服务数据库迁移 - 创建API请求日志表
-- 迁移版本: 001
-- 创建时间: 2025-07-28
-- 数据库: lyss_platform

BEGIN;

-- API 请求日志表
CREATE TABLE api_request_logs (
    id BIGSERIAL PRIMARY KEY,
    request_id VARCHAR(100) UNIQUE NOT NULL,
    user_id BIGINT,
    group_id BIGINT,
    credential_id BIGINT,
    model_name VARCHAR(100),
    provider_type VARCHAR(50),
    method VARCHAR(10),
    endpoint VARCHAR(500),
    request_size INTEGER,
    response_size INTEGER,
    status_code INTEGER,
    duration_ms INTEGER,
    error_message TEXT,
    request_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引设计
CREATE INDEX idx_api_request_logs_user_id ON api_request_logs(user_id);
CREATE INDEX idx_api_request_logs_request_time ON api_request_logs(request_time);
CREATE INDEX idx_api_request_logs_status_code ON api_request_logs(status_code);

-- 添加表注释
COMMENT ON TABLE api_request_logs IS 'API请求日志记录表';
COMMENT ON COLUMN api_request_logs.id IS '日志记录唯一标识符';
COMMENT ON COLUMN api_request_logs.request_id IS '请求唯一标识符';
COMMENT ON COLUMN api_request_logs.user_id IS '发起请求的用户ID';
COMMENT ON COLUMN api_request_logs.group_id IS '关联的群组ID';
COMMENT ON COLUMN api_request_logs.credential_id IS '使用的凭证ID';
COMMENT ON COLUMN api_request_logs.model_name IS '调用的模型名称';
COMMENT ON COLUMN api_request_logs.provider_type IS 'AI供应商类型';
COMMENT ON COLUMN api_request_logs.method IS 'HTTP请求方法';
COMMENT ON COLUMN api_request_logs.endpoint IS '请求端点';
COMMENT ON COLUMN api_request_logs.request_size IS '请求大小（字节）';
COMMENT ON COLUMN api_request_logs.response_size IS '响应大小（字节）';
COMMENT ON COLUMN api_request_logs.status_code IS 'HTTP状态码';
COMMENT ON COLUMN api_request_logs.duration_ms IS '请求处理时长（毫秒）';
COMMENT ON COLUMN api_request_logs.error_message IS '错误信息';
COMMENT ON COLUMN api_request_logs.request_time IS '请求时间';

COMMIT;