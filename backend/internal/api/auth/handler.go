package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/your-org/lyss-chat-backend/internal/domain/user"
	"github.com/your-org/lyss-chat-backend/internal/service/auth"
	"github.com/your-org/lyss-chat-backend/internal/util"
	"github.com/your-org/lyss-chat-backend/pkg/config"
	"github.com/your-org/lyss-chat-backend/pkg/db"
	"github.com/your-org/lyss-chat-backend/pkg/logger"
)

// Handler 表示认证处理器
type Handler struct {
	service *auth.Service
	logger  *logger.Logger
}

// NewHandler 创建一个新的认证处理器
func NewHandler(db *db.Postgres, redis *db.Redis, cfg *config.Config, logger *logger.Logger) *Handler {
	service := auth.NewService(db, redis, cfg, logger)
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Login 处理登录请求
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req user.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.BadRequestError(w, "无效的请求体", nil)
		return
	}

	// 验证请求
	if req.Email == "" || req.Password == "" || req.TenantID == "" {
		util.BadRequestError(w, "邮箱、密码和租户ID不能为空", nil)
		return
	}

	// 获取客户端 IP 和 User-Agent
	req.IP = r.RemoteAddr
	req.UserAgent = r.UserAgent()

	h.logger.Debug("登录请求", fmt.Sprintf("Email: %s, IP: %s, UserAgent: %s", req.Email, req.IP, req.UserAgent))

	// 调用服务
	resp, err := h.service.Login(&req)
	if err != nil {
		h.logger.Error("登录失败", err)
		util.UnauthorizedError(w, "邮箱或密码错误")
		return
	}

	util.SuccessResponse(w, resp, http.StatusOK)
}

// RefreshToken 处理刷新令牌请求
func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.BadRequestError(w, "无效的请求体", nil)
		return
	}

	// 验证请求
	if req.RefreshToken == "" {
		util.BadRequestError(w, "刷新令牌不能为空", nil)
		return
	}

	// 调用服务
	resp, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		h.logger.Error("刷新令牌失败", err)
		util.UnauthorizedError(w, "无效的刷新令牌")
		return
	}

	util.SuccessResponse(w, resp, http.StatusOK)
}

// Register 处理注册请求
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req user.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.BadRequestError(w, "无效的请求体", nil)
		return
	}

	// 验证请求
	if req.Email == "" || req.Password == "" || req.Name == "" || req.TenantID == "" {
		util.BadRequestError(w, "邮箱、密码、姓名和租户ID不能为空", nil)
		return
	}

	// 获取客户端 IP 和 User-Agent
	req.IP = r.RemoteAddr
	req.UserAgent = r.UserAgent()

	h.logger.Debug("注册请求", fmt.Sprintf("Email: %s, Name: %s, IP: %s", req.Email, req.Name, req.IP))

	// 调用服务
	newUser, err := h.service.Register(&req)
	if err != nil {
		h.logger.Error("注册失败", err)
		util.BadRequestError(w, fmt.Sprintf("注册失败: %s", err.Error()), nil)
		return
	}

	// 自动登录
	loginReq := &user.LoginRequest{
		Email:     req.Email,
		Password:  req.Password,
		TenantID:  req.TenantID,
		IP:        req.IP,
		UserAgent: req.UserAgent,
	}

	loginResp, err := h.service.Login(loginReq)
	if err != nil {
		h.logger.Error("注册后自动登录失败", err)
		// 返回用户信息，但不包含令牌
		util.SuccessResponse(w, newUser, http.StatusCreated)
		return
	}

	// 返回登录响应
	util.SuccessResponse(w, loginResp, http.StatusCreated)
}
