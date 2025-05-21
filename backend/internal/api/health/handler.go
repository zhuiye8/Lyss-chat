package health

import (
	"net/http"

	"github.com/your-org/lyss-chat-2.0/backend/internal/util"
	"github.com/your-org/lyss-chat-2.0/backend/pkg/db"
	"github.com/your-org/lyss-chat-2.0/backend/pkg/logger"
)

// Handler 表示健康检查处理器
type Handler struct {
	db     *db.Postgres
	redis  *db.Redis
	minio  *db.MinIO
	logger *logger.Logger
}

// NewHandler 创建一个新的健康检查处理器
func NewHandler(db *db.Postgres, redis *db.Redis, minio *db.MinIO, logger *logger.Logger) *Handler {
	return &Handler{
		db:     db,
		redis:  redis,
		minio:  minio,
		logger: logger,
	}
}

// Health 处理健康检查请求
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	// 检查数据库连接
	if err := h.db.Ping(); err != nil {
		h.logger.Error("数据库连接失败", err)
		util.InternalServerError(w, "数据库连接失败")
		return
	}

	// 检查 Redis 连接
	if err := h.redis.Ping(r.Context()).Err(); err != nil {
		h.logger.Error("Redis 连接失败", err)
		util.InternalServerError(w, "Redis 连接失败")
		return
	}

	// 返回成功响应
	util.SuccessResponse(w, map[string]string{"status": "ok"}, http.StatusOK)
}
