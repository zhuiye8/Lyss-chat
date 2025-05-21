package middleware

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/your-org/lyss-chat-2.0/backend/pkg/logger"
)

// Logger 中间件记录请求日志
func Logger(log *logger.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// 创建自定义响应写入器以捕获状态码
			rw := &responseWriter{w, http.StatusOK}

			// 处理请求
			next.ServeHTTP(rw, r)

			// 计算请求处理时间
			duration := time.Since(start)

			// 记录请求日志
			log.Infof(
				"%s %s %d %s %s",
				r.Method,
				r.URL.Path,
				rw.statusCode,
				duration,
				r.RemoteAddr,
			)
		})
	}
}

// responseWriter 是一个自定义的响应写入器，用于捕获状态码
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader 重写 WriteHeader 方法以捕获状态码
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
