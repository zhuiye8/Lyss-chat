package middleware

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/your-org/lyss-chat-backend/pkg/logger"
)

// Logger 创建一个日志中间件
func Logger(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// 包装 ResponseWriter 以捕获状态码
			wrapper := &responseWriterWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// 处理请求
			next.ServeHTTP(wrapper, r)

			// 记录请求信息
			duration := time.Since(start)
			log.Infof(
				"%s %s %d %s",
				r.Method,
				r.URL.Path,
				wrapper.statusCode,
				duration,
			)
		})
	}
}

// Recover 创建一个恢复中间件，用于捕获 panic
func Recover(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// 记录 panic 信息
					log.Error("服务器 panic", err)
					log.Debug(string(debug.Stack()))

					// 返回 500 错误
					http.Error(w, "内部服务器错误", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// CORS 创建一个 CORS 中间件
func CORS() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 设置 CORS 头
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// 处理预检请求
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// responseWriterWrapper 包装 http.ResponseWriter 以捕获状态码
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader 重写 WriteHeader 方法以捕获状态码
func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
