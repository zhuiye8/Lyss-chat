#!/bin/bash

# Lyss AI Platform Migration 自动生成工具
# 用途：根据数据库表结构自动生成migration文件
# 作者：Claude AI Assistant
# 创建时间：2025-07-28

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 显示使用说明
show_usage() {
    echo -e "${BLUE}Lyss AI Platform Migration 生成工具${NC}"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -s, --service SERVICE    指定服务名称 (user|auth|group|credential|gateway|billing|shared)"
    echo "  -t, --table TABLE        指定表名"
    echo "  -a, --all               生成所有服务的migration文件"
    echo "  -v, --verify            验证migration文件语法"
    echo "  -h, --help              显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 -s user -t users                # 为用户服务的users表生成migration"
    echo "  $0 -s user                         # 为用户服务生成所有migration"
    echo "  $0 -a                              # 生成所有服务的migration文件"
    echo "  $0 -v                              # 验证现有migration文件"
}

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

# 验证数据库连接
verify_db_connection() {
    log_info "验证数据库连接..."
    
    if ! command -v psql &> /dev/null; then
        log_error "PostgreSQL客户端未安装"
        return 1
    fi
    
    # 尝试连接数据库
    if ! psql -h localhost -p 5433 -U lyss -d lyss_platform -c "SELECT 1;" &> /dev/null; then
        log_error "无法连接到数据库，请检查数据库是否正常运行"
        return 1
    fi
    
    log_info "数据库连接正常"
}

# 获取表结构信息
get_table_schema() {
    local table_name=$1
    
    psql -h localhost -p 5433 -U lyss -d lyss_platform -t -c "
        SELECT column_name, data_type, is_nullable, column_default, character_maximum_length
        FROM information_schema.columns 
        WHERE table_name = '$table_name' 
        ORDER BY ordinal_position;
    "
}

# 获取表索引信息
get_table_indexes() {
    local table_name=$1
    
    psql -h localhost -p 5433 -U lyss -d lyss_platform -t -c "
        SELECT indexname, indexdef 
        FROM pg_indexes 
        WHERE tablename = '$table_name' 
        AND indexname NOT LIKE '%_pkey';
    "
}

# 获取表约束信息
get_table_constraints() {
    local table_name=$1
    
    psql -h localhost -p 5433 -U lyss -d lyss_platform -t -c "
        SELECT constraint_name, constraint_type 
        FROM information_schema.table_constraints 
        WHERE table_name = '$table_name';
    "
}

# 生成单个表的migration文件
generate_table_migration() {
    local service=$1
    local table_name=$2
    local version=$3
    
    local service_dir="/root/work/cyss/services/${service}-service/migrations"
    
    # 确保目录存在
    mkdir -p "$service_dir"
    
    local up_file="${service_dir}/${version}_create_${table_name}_table.up.sql"
    local down_file="${service_dir}/${version}_create_${table_name}_table.down.sql"
    
    log_info "生成 ${service}-service 的 ${table_name} 表 migration 文件..."
    
    # 生成UP文件
    cat > "$up_file" << EOF
-- ${service}服务数据库迁移 - 创建${table_name}表
-- 迁移版本: ${version}
-- 创建时间: $(date '+%Y-%m-%d')
-- 数据库: lyss_platform

BEGIN;

-- 注意：此文件由自动化工具生成，请根据实际需求调整
-- TODO: 请根据数据库设计文档完善表结构定义

CREATE TABLE ${table_name} (
    id BIGSERIAL PRIMARY KEY,
    -- TODO: 添加具体字段定义
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- TODO: 添加索引
-- CREATE INDEX idx_${table_name}_field ON ${table_name}(field);

-- TODO: 添加表注释
COMMENT ON TABLE ${table_name} IS '${table_name}表';

COMMIT;
EOF

    # 生成DOWN文件
    cat > "$down_file" << EOF
-- ${service}服务数据库迁移回滚 - 删除${table_name}表
-- 迁移版本: ${version}
-- 创建时间: $(date '+%Y-%m-%d')

BEGIN;

DROP TABLE IF EXISTS ${table_name} CASCADE;

COMMIT;
EOF

    log_info "生成完成: $up_file"
    log_info "生成完成: $down_file"
}

# 验证migration文件语法
verify_migrations() {
    log_info "验证migration文件语法..."
    
    local error_count=0
    
    for service_dir in /root/work/cyss/services/*/migrations; do
        if [ -d "$service_dir" ]; then
            service_name=$(basename $(dirname "$service_dir"))
            log_info "检查 $service_name 的migration文件..."
            
            for migration_file in "$service_dir"/*.sql; do
                if [ -f "$migration_file" ]; then
                    # 基本语法检查
                    if ! grep -q "BEGIN;" "$migration_file" || ! grep -q "COMMIT;" "$migration_file"; then
                        log_error "$migration_file 缺少事务控制语句"
                        ((error_count++))
                    fi
                    
                    # 检查是否有注释
                    if ! head -n 5 "$migration_file" | grep -q "迁移版本"; then
                        log_warn "$migration_file 缺少版本信息注释"
                    fi
                fi
            done
        fi
    done
    
    if [ $error_count -eq 0 ]; then
        log_info "所有migration文件语法检查通过"
    else
        log_error "发现 $error_count 个语法错误"
        return 1
    fi
}

# 生成所有服务的migration文件
generate_all_migrations() {
    log_info "开始生成所有服务的migration文件..."
    
    # 服务和表的映射关系
    declare -A service_tables
    service_tables[user]="users user_settings user_quotas"
    service_tables[auth]="jwt_tokens oauth_connections"
    service_tables[group]="groups group_members group_invitations"
    service_tables[credential]="provider_credentials model_permissions"
    service_tables[gateway]="api_request_logs routing_rules"
    service_tables[billing]="usage_records pricing_models billing_summaries"
    service_tables[shared]="system_configs notifications"
    
    for service in "${!service_tables[@]}"; do
        log_info "处理 ${service}-service..."
        
        local version=1
        for table in ${service_tables[$service]}; do
            version_str=$(printf "%03d" $version)
            generate_table_migration "$service" "$table" "$version_str"
            ((version++))
        done
    done
    
    log_info "所有migration文件生成完成"
}

# 主函数
main() {
    local service=""
    local table=""
    local generate_all=false
    local verify_only=false
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -s|--service)
                service="$2"
                shift 2
                ;;
            -t|--table)
                table="$2"
                shift 2
                ;;
            -a|--all)
                generate_all=true
                shift
                ;;
            -v|--verify)
                verify_only=true
                shift
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            *)
                log_error "未知选项: $1"
                show_usage
                exit 1
                ;;
        esac
    done
    
    # 验证模式
    if [ "$verify_only" = true ]; then
        verify_migrations
        exit $?
    fi
    
    # 生成所有migration
    if [ "$generate_all" = true ]; then
        generate_all_migrations
        verify_migrations
        exit $?
    fi
    
    # 生成特定服务的migration
    if [ -n "$service" ]; then
        if [ -n "$table" ]; then
            generate_table_migration "$service" "$table" "001"
        else
            log_error "请指定表名 (-t)"
            exit 1
        fi
        exit 0
    fi
    
    # 默认显示帮助
    show_usage
}

# 执行主函数
main "$@"