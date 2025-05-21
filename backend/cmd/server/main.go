package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// 健康检查响应
type HealthResponse struct {
	Status string `json:"status"`
}

// 健康检查处理器
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{Status: "ok"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Println("警告: 未找到 .env 文件，使用环境变量")
	}

	// 获取端口
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}

	// 创建路由器
	r := mux.NewRouter()

	// 注册健康检查路由
	r.HandleFunc("/v1/health", healthHandler).Methods("GET")

	// 创建 HTTP 服务器
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 启动服务器
	go func() {
		log.Printf("服务器启动在 :%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("服务器关闭中...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("服务器强制关闭: %v", err)
	}

	log.Println("服务器已关闭")
}
