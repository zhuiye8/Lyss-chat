package chat

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/your-org/lyss-chat-backend/internal/domain/chat"
	"github.com/your-org/lyss-chat-backend/internal/middleware"
	"github.com/your-org/lyss-chat-backend/internal/service/chat"
	"github.com/your-org/lyss-chat-backend/internal/util"
	"github.com/your-org/lyss-chat-backend/pkg/config"
	"github.com/your-org/lyss-chat-backend/pkg/db"
	"github.com/your-org/lyss-chat-backend/pkg/logger"
)

// CanvasHandler 表示画布处理器
type CanvasHandler struct {
	service *chat.Service
	logger  *logger.Logger
}

// NewCanvasHandler 创建一个新的画布处理器
func NewCanvasHandler(db *db.Postgres, cfg *config.Config, logger *logger.Logger) *CanvasHandler {
	// 注意：这里需要传入 aiGraphs，但我们暂时传入 nil，后续会修复
	service := chat.NewService(db, nil, cfg, logger)
	return &CanvasHandler{
		service: service,
		logger:  logger,
	}
}

// ListCanvases 处理获取画布列表请求
func (h *CanvasHandler) ListCanvases(w http.ResponseWriter, r *http.Request) {
	// 获取查询参数
	workspaceID := r.URL.Query().Get("workspace_id")
	if workspaceID == "" {
		util.BadRequestError(w, "工作区ID不能为空", nil)
		return
	}

	canvasType := r.URL.Query().Get("type")
	var canvasTypePtr *string
	if canvasType != "" {
		canvasTypePtr = &canvasType
	}

	// 获取分页参数
	page := 1
	pageSize := 20
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	// 调用服务
	canvases, total, err := h.service.ListCanvases(workspaceID, canvasTypePtr, page, pageSize)
	if err != nil {
		h.logger.Error("获取画布列表失败", err)
		util.InternalServerError(w, "获取画布列表失败")
		return
	}

	// 构建响应
	response := map[string]interface{}{
		"items":     canvases,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}

	util.SuccessResponse(w, response, http.StatusOK)
}

// GetCanvas 处理获取画布详情请求
func (h *CanvasHandler) GetCanvas(w http.ResponseWriter, r *http.Request) {
	// 获取路径参数
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		util.BadRequestError(w, "画布ID不能为空", nil)
		return
	}

	// 调用服务
	canvas, err := h.service.GetCanvas(id)
	if err != nil {
		h.logger.Error("获取画布详情失败", err)
		util.NotFoundError(w, "画布不存在")
		return
	}

	util.SuccessResponse(w, canvas, http.StatusOK)
}

// CreateCanvas 处理创建画布请求
func (h *CanvasHandler) CreateCanvas(w http.ResponseWriter, r *http.Request) {
	// 解析请求体
	var req chat.CreateCanvasRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.BadRequestError(w, "无效的请求体", nil)
		return
	}

	// 验证请求
	if req.Title == "" || req.WorkspaceID == "" {
		util.BadRequestError(w, "标题和工作区ID不能为空", nil)
		return
	}

	// 获取用户ID
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		util.UnauthorizedError(w, "未认证")
		return
	}

	// 调用服务
	canvas, err := h.service.CreateCanvas(userID, &req)
	if err != nil {
		h.logger.Error("创建画布失败", err)
		util.InternalServerError(w, "创建画布失败")
		return
	}

	util.SuccessResponse(w, canvas, http.StatusCreated)
}

// UpdateCanvas 处理更新画布请求
func (h *CanvasHandler) UpdateCanvas(w http.ResponseWriter, r *http.Request) {
	// 获取路径参数
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		util.BadRequestError(w, "画布ID不能为空", nil)
		return
	}

	// 解析请求体
	var req chat.UpdateCanvasRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.BadRequestError(w, "无效的请求体", nil)
		return
	}

	// 调用服务
	canvas, err := h.service.UpdateCanvas(id, &req)
	if err != nil {
		h.logger.Error("更新画布失败", err)
		util.NotFoundError(w, "画布不存在")
		return
	}

	util.SuccessResponse(w, canvas, http.StatusOK)
}

// DeleteCanvas 处理删除画布请求
func (h *CanvasHandler) DeleteCanvas(w http.ResponseWriter, r *http.Request) {
	// 获取路径参数
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		util.BadRequestError(w, "画布ID不能为空", nil)
		return
	}

	// 调用服务
	err := h.service.DeleteCanvas(id)
	if err != nil {
		h.logger.Error("删除画布失败", err)
		util.NotFoundError(w, "画布不存在")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
