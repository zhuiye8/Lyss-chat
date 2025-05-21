package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config 表示应用程序配置
type Config struct {
	LogLevel string        `json:"log_level"`
	Server   ServerConfig  `json:"server"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig   `json:"redis"`
	MinIO    MinIOConfig   `json:"minio"`
	JWT      JWTConfig     `json:"jwt"`
}

// ServerConfig 表示服务器配置
type ServerConfig struct {
	Port         int `json:"port"`
	ReadTimeout  int `json:"read_timeout"`
	WriteTimeout int `json:"write_timeout"`
	IdleTimeout  int `json:"idle_timeout"`
}

// DatabaseConfig 表示数据库配置
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

// RedisConfig 表示 Redis 配置
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// MinIOConfig 表示 MinIO 配置
type MinIOConfig struct {
	Endpoint  string `json:"endpoint"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
	UseSSL    bool   `json:"use_ssl"`
}

// JWTConfig 表示 JWT 配置
type JWTConfig struct {
	Secret           string `json:"secret"`
	ExpirationHours  int    `json:"expiration_hours"`
	RefreshExpirationHours int `json:"refresh_expiration_hours"`
}

// Load 从配置文件加载配置
func Load() (*Config, error) {
	// 默认配置
	config := &Config{
		LogLevel: "info",
		Server: ServerConfig{
			Port:         8000,
			ReadTimeout:  60,
			WriteTimeout: 60,
			IdleTimeout:  120,
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			DBName:   "lyss_chat",
			SSLMode:  "disable",
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
		},
		MinIO: MinIOConfig{
			Endpoint:  "localhost:9000",
			AccessKey: "minioadmin",
			SecretKey: "minioadmin",
			Bucket:    "lyss-chat",
			UseSSL:    false,
		},
		JWT: JWTConfig{
			Secret:           "your-secret-key",
			ExpirationHours:  24,
			RefreshExpirationHours: 168,
		},
	}

	// 尝试从配置文件加载
	configFile := getConfigFile()
	if configFile != "" {
		file, err := os.Open(configFile)
		if err == nil {
			defer file.Close()
			decoder := json.NewDecoder(file)
			err = decoder.Decode(config)
			if err != nil {
				return nil, err
			}
		}
	}

	// 从环境变量覆盖配置
	loadFromEnv(config)

	return config, nil
}

// getConfigFile 获取配置文件路径
func getConfigFile() string {
	// 按优先级尝试不同的配置文件位置
	configPaths := []string{
		"configs/config.json",
		"../configs/config.json",
		"../../configs/config.json",
	}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			absPath, _ := filepath.Abs(path)
			return absPath
		}
	}

	return ""
}

// loadFromEnv 从环境变量加载配置
func loadFromEnv(config *Config) {
	// 服务器配置
	if port := os.Getenv("APP_PORT"); port != "" {
		var p int
		if _, err := fmt.Sscanf(port, "%d", &p); err == nil {
			config.Server.Port = p
		}
	}

	// 数据库配置
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		var p int
		if _, err := fmt.Sscanf(port, "%d", &p); err == nil {
			config.Database.Port = p
		}
	}
	if user := os.Getenv("DB_USER"); user != "" {
		config.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		config.Database.DBName = dbName
	}
	if sslMode := os.Getenv("DB_SSL_MODE"); sslMode != "" {
		config.Database.SSLMode = sslMode
	}

	// Redis 配置
	if host := os.Getenv("REDIS_HOST"); host != "" {
		config.Redis.Host = host
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		var p int
		if _, err := fmt.Sscanf(port, "%d", &p); err == nil {
			config.Redis.Port = p
		}
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		config.Redis.Password = password
	}
	if db := os.Getenv("REDIS_DB"); db != "" {
		var d int
		if _, err := fmt.Sscanf(db, "%d", &d); err == nil {
			config.Redis.DB = d
		}
	}

	// MinIO 配置
	if endpoint := os.Getenv("MINIO_ENDPOINT"); endpoint != "" {
		config.MinIO.Endpoint = endpoint
	}
	if accessKey := os.Getenv("MINIO_ACCESS_KEY"); accessKey != "" {
		config.MinIO.AccessKey = accessKey
	}
	if secretKey := os.Getenv("MINIO_SECRET_KEY"); secretKey != "" {
		config.MinIO.SecretKey = secretKey
	}
	if bucket := os.Getenv("MINIO_BUCKET"); bucket != "" {
		config.MinIO.Bucket = bucket
	}
	if useSSL := os.Getenv("MINIO_USE_SSL"); strings.ToLower(useSSL) == "true" {
		config.MinIO.UseSSL = true
	}

	// JWT 配置
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		config.JWT.Secret = secret
	}
	if expiration := os.Getenv("JWT_EXPIRATION"); expiration != "" {
		var e int
		if _, err := fmt.Sscanf(expiration, "%d", &e); err == nil {
			config.JWT.ExpirationHours = e
		}
	}
	if refreshExpiration := os.Getenv("JWT_REFRESH_EXPIRATION"); refreshExpiration != "" {
		var e int
		if _, err := fmt.Sscanf(refreshExpiration, "%d", &e); err == nil {
			config.JWT.RefreshExpirationHours = e
		}
	}
}
