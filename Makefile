# Lyss AI Platform 开发工具集
# 版本: 1.0
# 更新时间: 2025-01-25

.PHONY: help setup health-check status-check start-all stop-all clean

# 默认目标
help:
	@echo "Lyss AI Platform 开发命令"
	@echo ""
	@echo "🚀 环境管理:"
	@echo "  make setup              - 初始化开发环境"
	@echo "  make health-check       - 检查所有服务健康状态"
	@echo "  make start-all          - 启动所有服务"
	@echo "  make stop-all           - 停止所有服务"
	@echo "  make clean              - 清理临时文件"
	@echo ""
	@echo "📊 状态管理:"
	@echo "  make status-check       - 检查项目开发状态"
	@echo "  make status-update      - 更新项目状态信息"
	@echo "  make progress-report    - 生成开发进度报告"
	@echo "  make issue-list         - 查看已知问题列表"
	@echo ""
	@echo "🧪 测试相关:"
	@echo "  make test-unit          - 运行单元测试"
	@echo "  make test-integration   - 运行集成测试"
	@echo "  make test-e2e           - 运行端到端测试"
	@echo "  make test-all           - 运行所有测试"
	@echo ""
	@echo "🏗️ 构建部署:"
	@echo "  make docker-build-all   - 构建所有Docker镜像"
	@echo "  make k8s-deploy-test    - Kubernetes测试部署"
	@echo ""
	@echo "🔧 开发工具:"
	@echo "  make code-format        - 代码格式化"
	@echo "  make code-lint          - 代码静态检查"
	@echo "  make generate-docs      - 生成API文档"

# 环境管理
setup:
	@echo "🚀 初始化Lyss AI Platform开发环境..."
	@echo "⚠️  此命令需要在阶段0-开发环境准备阶段实现"
	@echo "📖 请参考: docs/05-deployment/development-setup.md"

health-check:
	@echo "🏥 检查服务健康状态..."
	@echo "⚠️  此命令需要在环境搭建完成后实现"
	@echo "检查项目:"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - Redis: localhost:6379" 
	@echo "  - ClickHouse: localhost:8123"
	@echo "  - Prometheus: localhost:9090"
	@echo "  - Grafana: localhost:3000"

start-all:
	@echo "🚀 启动所有服务..."
	@echo "⚠️  此命令需要在相应服务开发完成后实现"

stop-all:
	@echo "🛑 停止所有服务..."
	@echo "⚠️  此命令需要在相应服务开发完成后实现"

clean:
	@echo "🧹 清理临时文件..."
	@find . -name "*.log" -delete
	@find . -name "*.tmp" -delete
	@find . -name "__pycache__" -type d -exec rm -rf {} + 2>/dev/null || true
	@echo "✅ 清理完成"

# 状态管理  
status-check:
	@echo "📊 Lyss AI Platform 项目状态检查"
	@echo "=================================="
	@echo ""
	@echo "📋 当前阶段信息:"
	@head -10 status/phase-progress.md | tail -5
	@echo ""
	@echo "🏗️ 基础设施状态:"
	@if [ -f "status/service-status.json" ]; then \
		echo "  - PostgreSQL: $$(cat status/service-status.json | grep -A2 'postgresql' | grep 'status' | cut -d'"' -f4)"; \
		echo "  - Redis: $$(cat status/service-status.json | grep -A2 'redis' | grep 'status' | cut -d'"' -f4)"; \
		echo "  - ClickHouse: $$(cat status/service-status.json | grep -A2 'clickhouse' | grep 'status' | cut -d'"' -f4)"; \
	else \
		echo "  ⚠️  service-status.json 文件不存在"; \
	fi
	@echo ""
	@echo "🔍 已知问题数量:"
	@if [ -f "status/known-issues.md" ]; then \
		echo "  - 总问题数: $$(grep -c "### 问题 #" status/known-issues.md || echo 0)"; \
	else \
		echo "  ⚠️  known-issues.md 文件不存在"; \
	fi

status-update:
	@echo "📝 更新项目状态信息..."
	@echo "当前时间: $$(date '+%Y-%m-%d %H:%M:%S')" >> status/update-log.txt
	@echo "✅ 状态更新完成"

progress-report:
	@echo "📈 生成开发进度报告..."
	@echo "===================="
	@if [ -f "status/phase-progress.md" ]; then \
		echo "📊 阶段完成情况:"; \
		grep -E "完成|已完成|待开始" status/phase-progress.md | head -5; \
	fi
	@echo ""
	@if [ -f "status/service-status.json" ]; then \
		echo "🏗️ 服务开发状态:"; \
		echo "  - 已完成: $$(grep -c '"status": "completed"' status/service-status.json || echo 0)"; \
		echo "  - 进行中: $$(grep -c '"status": "in_progress"' status/service-status.json || echo 0)"; \
		echo "  - 未开始: $$(grep -c '"status": "not_started"' status/service-status.json || echo 0)"; \
	fi

issue-list:
	@echo "🐛 已知问题列表"
	@echo "=============="
	@if [ -f "status/known-issues.md" ]; then \
		grep -E "### 问题 #|优先级.*高|解决状态.*待解决" status/known-issues.md | head -10; \
	else \
		echo "✅ 暂无已知问题"; \
	fi

# 测试相关
test-unit:
	@echo "🧪 运行单元测试..."
	@echo "⚠️  具体测试命令需要在服务开发阶段实现"

test-integration:
	@echo "🔗 运行集成测试..."
	@echo "⚠️  具体测试命令需要在集成阶段实现"

test-e2e:
	@echo "🎯 运行端到端测试..."
	@echo "⚠️  具体测试命令需要在系统集成阶段实现"

test-all: test-unit test-integration test-e2e
	@echo "✅ 所有测试完成"

# 构建部署
docker-build-all:
	@echo "🐳 构建所有Docker镜像..."
	@echo "⚠️  具体构建命令需要在生产部署准备阶段实现"

k8s-deploy-test:
	@echo "☸️  Kubernetes测试部署..."
	@echo "⚠️  具体部署命令需要在生产部署准备阶段实现"

# 开发工具
code-format:
	@echo "🎨 代码格式化..."
	@if command -v gofmt >/dev/null 2>&1; then \
		find . -name "*.go" -exec gofmt -w {} \;; \
		echo "✅ Go代码格式化完成"; \
	else \
		echo "⚠️  gofmt 未安装，跳过Go代码格式化"; \
	fi

code-lint:
	@echo "🔍 代码静态检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "⚠️  golangci-lint 未安装，请先安装: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

generate-docs:
	@echo "📚 生成API文档..."
	@echo "⚠️  具体文档生成命令需要在服务开发阶段实现"

# 数据库相关 (将在阶段1实现)
db-migrate:
	@echo "🗄️  数据库迁移..."
	@echo "⚠️  此命令需要在阶段1-数据库实施阶段实现"

db-ping:
	@echo "🏓 测试数据库连接..."
	@echo "⚠️  此命令需要在阶段1-数据库实施阶段实现"

db-seed:
	@echo "🌱 导入测试数据..."
	@echo "⚠️  此命令需要在阶段1-数据库实施阶段实现"

# 服务相关命令模板 (将在对应阶段实现)
user-service-start:
	@echo "👤 启动用户服务..."
	@echo "⚠️  此命令需要在阶段2-用户服务开发阶段实现"

user-service-test:
	@echo "🧪 用户服务测试..."
	@echo "⚠️  此命令需要在阶段2-用户服务开发阶段实现"

auth-service-start:
	@echo "🔐 启动认证服务..."
	@echo "⚠️  此命令需要在阶段3-认证服务开发阶段实现"

auth-service-test:
	@echo "🧪 认证服务测试..."
	@echo "⚠️  此命令需要在阶段3-认证服务开发阶段实现"

# 故障恢复相关
recovery-check:
	@echo "🚨 检查会话中断恢复状态..."
	@echo "📋 当前阶段状态:"
	@if [ -f "status/phase-progress.md" ]; then \
		grep -E "当前阶段|阶段状态|负责会话" status/phase-progress.md | head -3; \
	fi
	@echo ""
	@echo "🔍 最近更新:"
	@if [ -f "status/update-log.txt" ]; then \
		tail -5 status/update-log.txt; \
	else \
		echo "  ⚠️  无更新记录"; \
	fi

rollback-to-checkpoint:
	@echo "⏪ 回滚到最近检查点..."
	@echo "⚠️  具体回滚逻辑需要根据当前阶段实现"
	@echo "📖 请查看对应阶段的故障恢复指南"

# 调试相关
debug-service-%:
	@echo "🔧 调试服务: $*"
	@echo "⚠️  具体调试命令需要在服务实现时定义"

isolate-service-%:
	@echo "🚧 隔离故障服务: $*"
	@echo "⚠️  具体隔离逻辑需要在服务实现时定义"

infrastructure-recover:
	@echo "🏗️  基础设施恢复..."
	@echo "⚠️  具体恢复逻辑需要在环境搭建完成后实现"

# 监控相关
logs-tail:
	@echo "📜 查看实时日志..."
	@echo "⚠️  具体日志命令需要在服务实现时定义"

monitoring:
	@echo "📊 打开监控面板..."
	@echo "🌐 Grafana: http://localhost:3000"  
	@echo "🌐 Prometheus: http://localhost:9090"
	@echo "⚠️  请确保监控服务已启动"

# 项目初始化 (仅在第一次执行)
init-project:
	@echo "🎯 初始化项目结构..."
	@mkdir -p cmd services internal pkg configs migrations status logs
	@echo "📁 项目目录结构创建完成"
	@echo "📝 创建初始配置文件..."
	@touch status/update-log.txt
	@echo "✅ 项目初始化完成"

# 版本信息
version:
	@echo "Lyss AI Platform 开发工具集"
	@echo "版本: 1.0"
	@echo "更新时间: 2025-01-25"
	@echo "支持的开发阶段: 0-9 (共10个阶段)"