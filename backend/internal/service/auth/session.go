package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/your-org/lyss-chat-backend/internal/domain/user"
	"github.com/your-org/lyss-chat-backend/pkg/logger"
)

// SessionManager 管理用户会话
type SessionManager struct {
	redis  *RedisClient
	logger *logger.Logger
}

// RedisClient 是 Redis 客户端接口
type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
}

// NewSessionManager 创建一个新的会话管理器
func NewSessionManager(redis *RedisClient, logger *logger.Logger) *SessionManager {
	return &SessionManager{
		redis:  redis,
		logger: logger,
	}
}

// SessionData 表示会话数据
type SessionData struct {
	UserID    string    `json:"user_id"`
	TenantID  string    `json:"tenant_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	LastLogin time.Time `json:"last_login"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
}

// CreateSession 创建一个新的会话
func (m *SessionManager) CreateSession(ctx context.Context, sessionID string, u *user.User, ip, userAgent string, expiration time.Duration) error {
	sessionData := &SessionData{
		UserID:    u.ID,
		TenantID:  u.TenantID,
		Email:     u.Email,
		Name:      u.Name,
		LastLogin: time.Now(),
		IP:        ip,
		UserAgent: userAgent,
	}

	// 序列化会话数据
	data, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}

	// 存储会话数据
	key := fmt.Sprintf("session:%s", sessionID)
	return (*m.redis).Set(ctx, key, string(data), expiration)
}

// GetSession 获取会话数据
func (m *SessionManager) GetSession(ctx context.Context, sessionID string) (*SessionData, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	data, err := (*m.redis).Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var sessionData SessionData
	if err := json.Unmarshal([]byte(data), &sessionData); err != nil {
		return nil, err
	}

	return &sessionData, nil
}

// DeleteSession 删除会话
func (m *SessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return (*m.redis).Del(ctx, key)
}

// RefreshSession 刷新会话过期时间
func (m *SessionManager) RefreshSession(ctx context.Context, sessionID string, expiration time.Duration) error {
	// 获取当前会话数据
	sessionData, err := m.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	// 更新最后登录时间
	sessionData.LastLogin = time.Now()

	// 序列化会话数据
	data, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}

	// 重新存储会话数据
	key := fmt.Sprintf("session:%s", sessionID)
	return (*m.redis).Set(ctx, key, string(data), expiration)
}

// ListActiveSessions 列出用户的所有活跃会话
func (m *SessionManager) ListActiveSessions(ctx context.Context, userID string) ([]*SessionData, error) {
	// 注意：这是一个简化实现，实际上需要使用 Redis 的 SCAN 命令
	// 在真实环境中，应该使用 Redis 的 SCAN 命令或者维护一个用户到会话的映射
	m.logger.Warn("ListActiveSessions 是一个简化实现，实际应用中应使用 Redis 的 SCAN 命令")
	return []*SessionData{}, nil
}
