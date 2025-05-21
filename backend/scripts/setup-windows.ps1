# Lyss-chat 2.0 Windows 环境设置脚本

# 检查 Docker 是否运行
function Check-Docker {
    try {
        $dockerStatus = docker info
        Write-Host "Docker 正在运行" -ForegroundColor Green
        return $true
    } catch {
        Write-Host "Docker 未运行，请先启动 Docker Desktop" -ForegroundColor Red
        return $false
    }
}

# 检查容器状态
function Check-Container {
    param (
        [string]$containerName
    )
    
    $container = docker ps -f "name=$containerName" --format "{{.Names}}"
    
    if ($container -eq $containerName) {
        Write-Host "$containerName 容器正在运行" -ForegroundColor Green
        return $true
    } else {
        Write-Host "$containerName 容器未运行" -ForegroundColor Yellow
        return $false
    }
}

# 启动 Docker Compose
function Start-DockerCompose {
    Write-Host "启动 Docker Compose 服务..." -ForegroundColor Cyan
    docker-compose up -d
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Docker Compose 服务已启动" -ForegroundColor Green
    } else {
        Write-Host "Docker Compose 服务启动失败" -ForegroundColor Red
        exit 1
    }
}

# 运行数据库迁移
function Run-Migration {
    Write-Host "运行数据库迁移..." -ForegroundColor Cyan
    
    # 确保 migrations 目录存在
    if (-not (Test-Path -Path "migrations")) {
        Write-Host "未找到 migrations 目录" -ForegroundColor Red
        exit 1
    }
    
    # 尝试运行迁移
    try {
        go run cmd/migrate/main.go up
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "数据库迁移成功" -ForegroundColor Green
        } else {
            Write-Host "数据库迁移失败" -ForegroundColor Red
            exit 1
        }
    } catch {
        Write-Host "数据库迁移失败: $_" -ForegroundColor Red
        exit 1
    }
}

# 安装 Air
function Install-Air {
    Write-Host "检查 Air 是否已安装..." -ForegroundColor Cyan
    
    $airPath = Get-Command air -ErrorAction SilentlyContinue
    
    if ($airPath) {
        Write-Host "Air 已安装" -ForegroundColor Green
    } else {
        Write-Host "安装 Air..." -ForegroundColor Cyan
        go install github.com/cosmtrek/air@latest
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "Air 安装成功" -ForegroundColor Green
        } else {
            Write-Host "Air 安装失败" -ForegroundColor Red
            Write-Host "请手动安装 Air: go install github.com/cosmtrek/air@latest" -ForegroundColor Yellow
        }
    }
}

# 主函数
function Main {
    Write-Host "Lyss-chat 2.0 Windows 环境设置" -ForegroundColor Cyan
    
    # 检查 Docker
    if (-not (Check-Docker)) {
        exit 1
    }
    
    # 检查 PostgreSQL 容器
    $postgresRunning = Check-Container -containerName "lyss-chat-postgres"
    
    # 检查 Redis 容器
    $redisRunning = Check-Container -containerName "lyss-chat-redis"
    
    # 检查 MinIO 容器
    $minioRunning = Check-Container -containerName "lyss-chat-minio"
    
    # 如果任何容器未运行，启动 Docker Compose
    if (-not ($postgresRunning -and $redisRunning -and $minioRunning)) {
        Start-DockerCompose
    }
    
    # 运行数据库迁移
    Run-Migration
    
    # 安装 Air
    Install-Air
    
    Write-Host "环境设置完成！" -ForegroundColor Green
    Write-Host "可以使用以下命令启动后端服务：" -ForegroundColor Cyan
    Write-Host "cd backend" -ForegroundColor Yellow
    Write-Host "air" -ForegroundColor Yellow
}

# 执行主函数
Main
