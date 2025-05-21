package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/your-org/lyss-chat-2.0/backend/pkg/config"
)

func main() {
	// 解析命令行参数
	var command string
	flag.StringVar(&command, "command", "up", "迁移命令: up, down, version")
	flag.Parse()

	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Println("警告: 未找到 .env 文件，使用环境变量")
	}

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 构建数据库连接字符串
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	// 获取迁移文件路径
	migrationsPath, err := getMigrationsPath()
	if err != nil {
		log.Fatalf("获取迁移文件路径失败: %v", err)
	}

	// 创建迁移实例
	sourceURL := fmt.Sprintf("file://%s", strings.ReplaceAll(migrationsPath, "\\", "/"))
	m, err := migrate.New(
		sourceURL,
		dbURL,
	)
	if err != nil {
		log.Fatalf("创建迁移实例失败: %v", err)
	}
	defer m.Close()

	// 执行迁移命令
	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("迁移失败: %v", err)
		}
		log.Println("迁移成功")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("回滚失败: %v", err)
		}
		log.Println("回滚成功")
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("获取版本失败: %v", err)
		}
		log.Printf("当前版本: %d, 是否有未完成的迁移: %v", version, dirty)
	default:
		log.Fatalf("未知命令: %s", command)
	}
}

// getMigrationsPath 获取迁移文件路径
func getMigrationsPath() (string, error) {
	// 尝试从当前目录获取
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 尝试找到 migrations 目录
	migrationsPath := filepath.Join(currentDir, "migrations")
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		// 如果当前目录下没有 migrations 目录，尝试上一级目录
		migrationsPath = filepath.Join(filepath.Dir(currentDir), "migrations")
		if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
			return "", fmt.Errorf("未找到 migrations 目录")
		}
	}

	return migrationsPath, nil
}
