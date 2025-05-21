# Lyss-chat 2.0 前端 Windows 环境设置脚本

# 检查 Node.js 是否安装
function Check-NodeJS {
    try {
        $nodeVersion = node -v
        $npmVersion = npm -v
        
        Write-Host "Node.js 版本: $nodeVersion" -ForegroundColor Green
        Write-Host "npm 版本: $npmVersion" -ForegroundColor Green
        
        # 检查 Node.js 版本是否满足要求
        $versionString = $nodeVersion.Substring(1)  # 移除 'v' 前缀
        $versionParts = $versionString.Split('.')
        $majorVersion = [int]$versionParts[0]
        
        if ($majorVersion -lt 18) {
            Write-Host "警告: Node.js 版本低于推荐的 18.x 版本" -ForegroundColor Yellow
            return $false
        }
        
        return $true
    } catch {
        Write-Host "未检测到 Node.js，请安装 Node.js 18.x 或更高版本" -ForegroundColor Red
        return $false
    }
}

# 创建 .env.local 文件
function Create-EnvFile {
    $envPath = ".env.local"
    
    if (Test-Path $envPath) {
        Write-Host ".env.local 文件已存在" -ForegroundColor Green
    } else {
        Write-Host "创建 .env.local 文件..." -ForegroundColor Cyan
        
        $envContent = @"
VITE_API_BASE_URL=http://localhost:8000/v1
VITE_AUTH_DOMAIN=localhost
VITE_AUTH_CLIENT_ID=lyss-chat-local
"@
        
        $envContent | Out-File -FilePath $envPath -Encoding utf8
        
        if (Test-Path $envPath) {
            Write-Host ".env.local 文件创建成功" -ForegroundColor Green
        } else {
            Write-Host ".env.local 文件创建失败" -ForegroundColor Red
            exit 1
        }
    }
}

# 安装依赖
function Install-Dependencies {
    Write-Host "安装前端依赖..." -ForegroundColor Cyan
    
    npm install
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "依赖安装成功" -ForegroundColor Green
    } else {
        Write-Host "依赖安装失败" -ForegroundColor Red
        
        # 尝试清除缓存并重新安装
        Write-Host "尝试清除 npm 缓存并重新安装..." -ForegroundColor Yellow
        npm cache clean --force
        
        # 删除 node_modules 和 package-lock.json
        if (Test-Path "node_modules") {
            Remove-Item -Recurse -Force "node_modules"
        }
        
        if (Test-Path "package-lock.json") {
            Remove-Item -Force "package-lock.json"
        }
        
        # 重新安装
        npm install
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "依赖安装成功" -ForegroundColor Green
        } else {
            Write-Host "依赖安装失败，请手动解决问题" -ForegroundColor Red
            exit 1
        }
    }
}

# 主函数
function Main {
    Write-Host "Lyss-chat 2.0 前端 Windows 环境设置" -ForegroundColor Cyan
    
    # 检查 Node.js
    if (-not (Check-NodeJS)) {
        Write-Host "请安装 Node.js 18.x 或更高版本后重试" -ForegroundColor Yellow
        Write-Host "下载地址: https://nodejs.org/" -ForegroundColor Yellow
        exit 1
    }
    
    # 创建 .env.local 文件
    Create-EnvFile
    
    # 安装依赖
    Install-Dependencies
    
    Write-Host "前端环境设置完成！" -ForegroundColor Green
    Write-Host "可以使用以下命令启动前端开发服务器：" -ForegroundColor Cyan
    Write-Host "npm run dev" -ForegroundColor Yellow
}

# 执行主函数
Main
