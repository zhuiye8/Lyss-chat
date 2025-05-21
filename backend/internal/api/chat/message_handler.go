package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/your-org/lyss-chat-backend/internal/domain/chat"
	"github.com/your-org/lyss-chat-backend/internal/middleware"
	chatService "github.com/your-org/lyss-chat-backend/internal/service/chat"
	"github.com/your-org/lyss-chat-backend/internal/util"
	"github.com/your-org/lyss-chat-backend/pkg/config"
	"github.com/your-org/lyss-chat-backend/pkg/db"
	"github.com/your-org/lyss-chat-backend/pkg/logger"
)

// MessageHandler 表示消息处理器
type MessageHandler struct {
	service *chatService.Service
	logger  *logger.Logger
}

// NewMessageHandler 创建一个新的消息处理器
func NewMessageHandler(db *db.Postgres, cfg *config.Config, logger *logger.Logger) *MessageHandler {
	// 注意：这里需要传入 aiGraphs，但我们暂时传入 nil，后续会修复
	service := chatService.NewService(db, nil, cfg, logger)
	return &MessageHandler{
		service: service,
		logger:  logger,
	}
}

// ListMessages 处理获取消息列表请求
func (h *MessageHandler) ListMessages(w http.ResponseWriter, r *http.Request) {
	// 获取路径参数
	vars := mux.Vars(r)
	canvasID := vars["id"]
	if canvasID == "" {
		util.BadRequestError(w, "画布ID不能为空", nil)
		return
	}

	// 获取分页参数
	page := 1
	pageSize := 50
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
	messages, total, err := h.service.GetMessages(canvasID, page, pageSize)
	if err != nil {
		h.logger.Error("获取消息列表失败", err)
		util.InternalServerError(w, "获取消息列表失败")
		return
	}

	// 构建响应
	response := map[string]interface{}{
		"items":     messages,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}

	util.SuccessResponse(w, response, http.StatusOK)
}

// SendMessage 处理发送消息请求
func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	// 获取路径参数
	vars := mux.Vars(r)
	canvasID := vars["id"]
	if canvasID == "" {
		util.BadRequestError(w, "画布ID不能为空", nil)
		return
	}

	// 解析请求体
	var req chat.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.BadRequestError(w, "无效的请求体", nil)
		return
	}

	// 验证请求
	if req.Content == "" {
		util.BadRequestError(w, "消息内容不能为空", nil)
		return
	}

	// 获取用户ID
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		util.UnauthorizedError(w, "未认证")
		return
	}

	// 调用服务
	message, err := h.service.SendMessage(userID, canvasID, &req)
	if err != nil {
		h.logger.Error("发送消息失败", err)
		util.InternalServerError(w, "发送消息失败")
		return
	}

	util.SuccessResponse(w, message, http.StatusCreated)
}

// StreamMessage 处理流式发送消息请求
func (h *MessageHandler) StreamMessage(w http.ResponseWriter, r *http.Request) {
	// 获取路径参数
	vars := mux.Vars(r)
	canvasID := vars["id"]
	if canvasID == "" {
		util.BadRequestError(w, "画布ID不能为空", nil)
		return
	}

	// 解析请求体
	var req chat.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.BadRequestError(w, "无效的请求体", nil)
		return
	}

	// 验证请求
	if req.Content == "" {
		util.BadRequestError(w, "消息内容不能为空", nil)
		return
	}

	// 获取用户ID
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		util.UnauthorizedError(w, "未认证")
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	// 调用服务
	messageChan, err := h.service.StreamMessage(userID, canvasID, &req)
	if err != nil {
		h.logger.Error("流式发送消息失败", err)
		util.InternalServerError(w, "流式发送消息失败")
		return
	}

	// 发送用户消息确认
	userMessage := map[string]interface{}{
		"type":    "user",
		"content": req.Content,
	}
	userMessageJSON, _ := json.Marshal(userMessage)
	fmt.Fprintf(w, "data: %s\n\n", userMessageJSON)
	w.(http.Flusher).Flush()

	// 流式发送 AI 响应
	for message := range messageChan {
		// 构建事件数据
		event := map[string]interface{}{
			"type":    "assistant",
			"id":      message.ID,
			"content": message.Content,
		}
		
		// 序列化为 JSON
		eventJSON, err := json.Marshal(event)
		if err != nil {
			h.logger.Error("序列化事件失败", err)
			continue
		}
		
		// 发送事件
		fmt.Fprintf(w, "data: %s\n\n", eventJSON)
		w.(http.Flusher).Flush()
	}

	// 发送结束事件
	endEvent := map[string]interface{}{
		"type": "done",
	}
	endEventJSON, _ := json.Marshal(endEvent)
	fmt.Fprintf(w, "data: %s\n\n", endEventJSON)
	w.(http.Flusher).Flush()
}
