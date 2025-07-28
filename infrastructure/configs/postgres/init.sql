-- Lyss AI Platform 数据库初始化脚本
-- 创建各微服务专用数据库

-- 创建数据库
CREATE DATABASE lyss_user;
CREATE DATABASE lyss_auth;
CREATE DATABASE lyss_group;
CREATE DATABASE lyss_credential;
CREATE DATABASE lyss_gateway;
CREATE DATABASE lyss_billing;

-- 创建各服务的专用用户
CREATE USER lyss_user_svc WITH ENCRYPTED PASSWORD 'user123';
CREATE USER lyss_auth_svc WITH ENCRYPTED PASSWORD 'auth123';
CREATE USER lyss_group_svc WITH ENCRYPTED PASSWORD 'group123';
CREATE USER lyss_credential_svc WITH ENCRYPTED PASSWORD 'credential123';
CREATE USER lyss_gateway_svc WITH ENCRYPTED PASSWORD 'gateway123';
CREATE USER lyss_billing_svc WITH ENCRYPTED PASSWORD 'billing123';

-- 授权数据库访问权限
GRANT ALL PRIVILEGES ON DATABASE lyss_user TO lyss_user_svc;
GRANT ALL PRIVILEGES ON DATABASE lyss_auth TO lyss_auth_svc;
GRANT ALL PRIVILEGES ON DATABASE lyss_group TO lyss_group_svc;
GRANT ALL PRIVILEGES ON DATABASE lyss_credential TO lyss_credential_svc;
GRANT ALL PRIVILEGES ON DATABASE lyss_gateway TO lyss_gateway_svc; 
GRANT ALL PRIVILEGES ON DATABASE lyss_billing TO lyss_billing_svc;

-- 为主数据库创建必要扩展
\c lyss_platform;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "btree_gin";

-- 为用户服务数据库创建扩展
\c lyss_user;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- 为认证服务数据库创建扩展
\c lyss_auth;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 为群组服务数据库创建扩展
\c lyss_group;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- 为凭证服务数据库创建扩展
\c lyss_credential;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 为网关服务数据库创建扩展
\c lyss_gateway;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- 为计费服务数据库创建扩展
\c lyss_billing;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 创建公共函数和类型
\c lyss_platform;

-- 创建状态枚举类型
CREATE TYPE user_status AS ENUM ('active', 'inactive', 'suspended');
CREATE TYPE group_type AS ENUM ('personal', 'team', 'organization');
CREATE TYPE credential_type AS ENUM ('openai', 'azure_openai', 'anthropic', 'google', 'aliyun');
CREATE TYPE billing_status AS ENUM ('active', 'suspended', 'expired');

-- 创建审计字段更新函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 创建ID生成函数
CREATE OR REPLACE FUNCTION generate_short_id(length INTEGER DEFAULT 8)
RETURNS TEXT AS $$
DECLARE
    chars TEXT := 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
    result TEXT := '';
    i INTEGER := 0;
BEGIN
    FOR i IN 1..length LOOP
        result := result || substr(chars, floor(random() * length(chars) + 1)::integer, 1);
    END LOOP;
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- 初始化完成标记
INSERT INTO information_schema.sql_features VALUES ('LYSS_INIT_COMPLETE', 'YES', '');

-- 显示初始化结果
SELECT 'Lyss AI Platform 数据库初始化完成' AS status;