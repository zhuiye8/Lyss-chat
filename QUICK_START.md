# 🚀 LYSS AI Platform - 快速启动指南

## 📋 问题修复说明

根据您遇到的Docker Compose问题，我们已经完成了以下修复：

### ✅ 已修复的问题
1. **移除过时的version字段** - Docker Compose新版本不再需要
2. **解决Poetry依赖问题** - 提供requirements.txt备选方案
3. **分离数据库服务** - 创建专用的数据库Docker配置
4. **本地开发支持** - 完善的本地运行脚本

## 🎯 推荐开发方式 (您要求的方案)

### 步骤 1: 启动数据库服务 (Docker)

```bash
# 方式1: 使用专用的数据库配置
cd /root/work/Lyss
docker-compose -f docker-compose.db.yml up -d

# 或方式2: 使用修复后的主配置文件 (backend已注释)
docker-compose up -d
```

### 步骤 2: 验证数据库服务状态

```bash
# 检查服务状态
docker-compose -f docker-compose.db.yml ps

# 查看日志 (如果有问题)
docker-compose -f docker-compose.db.yml logs
```

### 步骤 3: 本地运行后端 (便于调试)

```bash
cd /root/work/Lyss/backend

# 使用便捷启动脚本 (自动处理依赖和环境)
python run_dev.py

# 或手动方式
pip install -r requirements.txt
python -m uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload
```

## 🔧 配置说明

### 数据库连接配置
后端会自动连接到Docker中的数据库服务：
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379  
- **Qdrant**: localhost:6333

### 环境配置文件
- `backend/.env.local` - 本地开发配置
- `backend/.env.example` - 配置模板

## 📊 服务访问地址

启动成功后，您可以访问：

| 服务 | 地址 | 说明 |
|------|------|------|
| **后端API** | http://localhost:8000 | FastAPI应用 |
| **API文档** | http://localhost:8000/api/v1/docs | Swagger UI |
| **ReDoc文档** | http://localhost:8000/api/v1/redoc | 替代文档 |
| **PostgreSQL** | localhost:5432 | 数据库 |
| **Redis** | localhost:6379 | 缓存 |
| **Qdrant** | localhost:6333 | 向量数据库 |

## 🔐 默认管理员账户

```
邮箱: admin@lyss.ai
密码: admin123
```

## 🛠️ 开发调试优势

这种方式提供了以下调试优势：

1. **实时代码重载** - 修改代码立即生效
2. **完整错误堆栈** - 详细的Python错误信息
3. **IDE集成调试** - 可以设置断点和步进调试
4. **快速依赖更新** - 无需重建Docker镜像
5. **环境变量控制** - 灵活的配置管理

## 🔍 故障排除

### 数据库连接失败
```bash
# 检查数据库服务状态
docker-compose -f docker-compose.db.yml ps

# 重启数据库服务
docker-compose -f docker-compose.db.yml restart

# 查看详细日志
docker-compose -f docker-compose.db.yml logs db
```

### 依赖安装问题
```bash
# 清理并重装依赖
pip install --upgrade pip
pip install -r requirements.txt

# 或使用虚拟环境
python -m venv venv
source venv/bin/activate  # Linux/Mac
pip install -r requirements.txt
```

### 端口占用问题
```bash
# 检查端口占用
netstat -tulpn | grep :8000
netstat -tulpn | grep :5432

# 修改端口 (在.env.local中)
# 或停止占用端口的进程
```

## 📈 下一步开发

1. **测试API功能** - 通过 http://localhost:8000/api/v1/docs
2. **配置AI供应商** - 添加OpenAI、Anthropic等API密钥
3. **开发前端界面** - 按照FRONTEND_DEV_CHECKLIST.md
4. **添加业务逻辑** - 根据PRD.md需求扩展功能

## 🎉 成功指标

如果看到以下信息，说明启动成功：

```
🎉 Starting FastAPI development server
📍 Server: http://localhost:8000
📚 API Docs: http://localhost:8000/api/v1/docs
🔧 ReDoc: http://localhost:8000/api/v1/redoc
👤 Admin: admin@lyss.ai / admin123
```

现在您可以享受高效的本地开发体验！🚀