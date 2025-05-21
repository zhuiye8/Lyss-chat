# Lyss-chat 2.0 Windows 环境搭建指南

本文档提供了在 Windows 环境下搭建 Lyss-chat 2.0 开发环境的详细步骤。

## 前提条件

确保您的系统已安装以下软件：

1. **Git**：用于克隆和管理代码仓库
2. **Docker Desktop for Windows**：用于运行容器化服务
3. **Go 1.24.0 或更高版本**：用于后端开发
4. **Node.js 18.x 或更高版本**：用于前端开发
5. **Visual Studio Code**（推荐）：用于代码编辑

## 克隆项目

1. 打开 PowerShell 或 Command Prompt
2. 执行以下命令：

```powershell
mkdir -p d:\work\test-demo
cd d:\work\test-demo
git clone https://github.com/your-org/lyss-chat-2.0.git
cd lyss-chat-2.0
```

## 后端环境搭建

### 1. 启动 Docker 服务

确保 Docker Desktop 正在运行，然后执行：

```powershell
cd d:\work\test-demo\lyss-chat-2.0\backend
docker-compose up -d
```

验证容器是否正在运行：

```powershell
docker ps -f "name=lyss-chat"
```

应该看到三个容器：`lyss-chat-postgres`、`lyss-chat-redis` 和 `lyss-chat-minio`。

### 2. 安装 Go 依赖

```powershell
cd d:\work\test-demo\lyss-chat-2.0\backend
go mod download
```

### 3. 安装 Air（用于热重载）

```powershell
go install github.com/cosmtrek/air@latest
```

### 4. 运行数据库迁移

```powershell
cd d:\work\test-demo\lyss-chat-2.0\backend
go run cmd/migrate/main.go up
```

如果遇到问题，可以使用我们提供的 PowerShell 脚本：

```powershell
cd d:\work\test-demo\lyss-chat-2.0\backend
.\scripts\setup-windows.ps1
```

### 5. 启动后端服务

```powershell
cd d:\work\test-demo\lyss-chat-2.0\backend
air
```

后端服务将在 http://localhost:8000 上运行。

## 前端环境搭建

### 1. 创建环境变量文件

在前端项目根目录创建 `.env.local` 文件：

```powershell
cd d:\work\test-demo\lyss-chat-2.0\frontend
```

创建文件 `.env.local` 并添加以下内容：

```
VITE_API_BASE_URL=http://localhost:8000/v1
VITE_AUTH_DOMAIN=localhost
VITE_AUTH_CLIENT_ID=lyss-chat-local
```

### 2. 安装依赖

```powershell
cd d:\work\test-demo\lyss-chat-2.0\frontend
npm install
```

### 3. 启动前端开发服务器

```powershell
cd d:\work\test-demo\lyss-chat-2.0\frontend
npm run dev
```

如果遇到问题，可以使用我们提供的 PowerShell 脚本：

```powershell
cd d:\work\test-demo\lyss-chat-2.0\frontend
.\scripts\setup-windows.ps1
```

前端服务将在 http://localhost:5173 上运行。

## 验证环境

### 验证后端

1. 打开浏览器，访问 http://localhost:8000/v1/health
2. 应该看到类似 `{"data":{"status":"ok"}}` 的响应

### 验证前端

1. 打开浏览器，访问 http://localhost:5173
2. 应该看到 Lyss-chat 2.0 的登录页面

## Windows 环境下的常用命令替代

在 Windows PowerShell 中，一些 Linux/macOS 命令需要替换：

| Linux/macOS 命令 | Windows PowerShell 替代 |
|-----------------|------------------------|
| `docker ps \| grep postgres` | `docker ps -f "name=postgres"` 或 `docker ps \| findstr postgres` |
| `rm -rf node_modules` | `Remove-Item -Recurse -Force node_modules` |
| `cat file.txt` | `Get-Content file.txt` |
| `mkdir -p dir/subdir` | `New-Item -ItemType Directory -Force -Path dir/subdir` |
| `touch file.txt` | `New-Item -ItemType File -Path file.txt` |

## 常见问题解决

### Docker 容器无法启动

1. 检查 Docker Desktop 是否正在运行
2. 检查端口是否被占用：
   ```powershell
   netstat -ano | findstr :5432  # 检查 PostgreSQL 端口
   netstat -ano | findstr :6379  # 检查 Redis 端口
   netstat -ano | findstr :9000  # 检查 MinIO 端口
   ```
3. 如果端口被占用，可以修改 `docker-compose.yml` 文件中的端口映射

### 数据库迁移失败

1. 确保 PostgreSQL 容器正在运行
2. 检查 `.env` 文件中的数据库配置是否正确
3. 尝试手动连接数据库：
   ```powershell
   docker exec -it lyss-chat-postgres psql -U postgres -d lyss_chat
   ```

### 前端依赖安装失败

1. 清除 npm 缓存：
   ```powershell
   npm cache clean --force
   ```
2. 删除 `node_modules` 目录和 `package-lock.json` 文件：
   ```powershell
   Remove-Item -Recurse -Force node_modules
   Remove-Item -Force package-lock.json
   ```
3. 重新安装依赖：
   ```powershell
   npm install
   ```

### 前端无法连接后端 API

1. 确保后端服务正在运行
2. 检查 `.env.local` 文件中的 `VITE_API_BASE_URL` 是否正确
3. 检查浏览器控制台是否有 CORS 错误
4. 使用浏览器开发者工具检查网络请求，查看具体错误信息
