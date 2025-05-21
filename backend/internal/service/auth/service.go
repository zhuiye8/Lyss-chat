package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/your-org/lyss-chat-backend/internal/domain/user"
	"github.com/your-org/lyss-chat-backend/internal/repository/postgres"
	"github.com/your-org/lyss-chat-backend/pkg/config"
	"github.com/your-org/lyss-chat-backend/pkg/db"
	"github.com/your-org/lyss-chat-backend/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

// Service 表示认证服务
type Service struct {
	userRepo      user.Repository
	redis         *db.Redis
	cfg           *config.Config
	logger        *logger.Logger
	sessionManager *SessionManager
}

// NewService 创建一个新的认证服务
func NewService(db *db.Postgres, redis *db.Redis, cfg *config.Config, logger *logger.Logger) *Service {
	userRepo := postgres.NewUserRepository(db)

	// 创建 Redis 客户端适配器
	redisClient := &redisClientAdapter{redis: redis}

	// 创建会话管理器
	sessionManager := NewSessionManager(&redisClient, logger)

	return &Service{
		userRepo:      userRepo,
		redis:         redis,
		cfg:           cfg,
		logger:        logger,
		sessionManager: sessionManager,
	}
}

// redisClientAdapter 适配 Redis 客户端接口
type redisClientAdapter struct {
	redis *db.Redis
}

// Set 实现 RedisClient 接口的 Set 方法
func (a *redisClientAdapter) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return a.redis.Client.Set(ctx, key, value, expiration).Err()
}

// Get 实现 RedisClient 接口的 Get 方法
func (a *redisClientAdapter) Get(ctx context.Context, key string) (string, error) {
	return a.redis.Client.Get(ctx, key).Result()
}

// Del 实现 RedisClient 接口的 Del 方法
func (a *redisClientAdapter) Del(ctx context.Context, keys ...string) error {
	return a.redis.Client.Del(ctx, keys...).Err()
}
}

// Login 处理用户登录
func (s *Service) Login(req *user.LoginRequest) (*user.LoginResponse, error) {
	// 获取用户
	u, err := s.userRepo.GetByEmail(req.Email, req.TenantID)
	if err != nil {
		return nil, err
	}

	// 检查用户状态
	if u.Status != user.UserStatusActive {
		return nil, errors.New("用户未激活")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, err
	}

	// 生成令牌
	accessToken, err := s.generateAccessToken(u)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(u)
	if err != nil {
		return nil, err
	}

	// 创建上下文
	ctx := context.Background()

	// 存储刷新令牌
	refreshKey := fmt.Sprintf("refresh_token:%s", refreshToken)
	err = s.redis.Client.Set(ctx, refreshKey, u.ID, time.Duration(s.cfg.JWT.RefreshExpirationHours)*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	// 创建会话
	sessionID := refreshToken // 使用刷新令牌作为会话ID
	err = s.sessionManager.CreateSession(
		ctx,
		sessionID,
		u,
		req.IP,
		req.UserAgent,
		time.Duration(s.cfg.JWT.RefreshExpirationHours)*time.Hour,
	)
	if err != nil {
		s.logger.Error("创建会话失败", err)
		// 继续处理，不要因为会话创建失败而阻止登录
	}

	return &user.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.cfg.JWT.ExpirationHours * 3600,
		TokenType:    "Bearer",
		User:         u,
	}, nil
}

// RefreshToken 刷新访问令牌
func (s *Service) RefreshToken(refreshToken string) (*user.LoginResponse, error) {
	// 验证刷新令牌
	ctx := context.Background()
	refreshKey := fmt.Sprintf("refresh_token:%s", refreshToken)
	userID, err := s.redis.Client.Get(ctx, refreshKey).Result()
	if err != nil {
		return nil, errors.New("无效的刷新令牌")
	}

	// 获取用户
	u, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// 检查用户状态
	if u.Status != user.UserStatusActive {
		return nil, errors.New("用户未激活")
	}

	// 生成新的访问令牌
	accessToken, err := s.generateAccessToken(u)
	if err != nil {
		return nil, err
	}

	// 生成新的刷新令牌
	newRefreshToken, err := s.generateRefreshToken(u)
	if err != nil {
		return nil, err
	}

	// 删除旧的刷新令牌
	err = s.redis.Client.Del(ctx, refreshKey).Err()
	if err != nil {
		s.logger.Error("删除旧的刷新令牌失败", err)
	}

	// 存储新的刷新令牌
	newRefreshKey := fmt.Sprintf("refresh_token:%s", newRefreshToken)
	err = s.redis.Client.Set(ctx, newRefreshKey, u.ID, time.Duration(s.cfg.JWT.RefreshExpirationHours)*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	// 刷新会话
	sessionID := newRefreshToken
	// 尝试获取旧会话数据
	oldSessionData, err := s.sessionManager.GetSession(ctx, refreshToken)
	if err == nil {
		// 创建新会话，保留旧会话的 IP 和 UserAgent
		err = s.sessionManager.CreateSession(
			ctx,
			sessionID,
			u,
			oldSessionData.IP,
			oldSessionData.UserAgent,
			time.Duration(s.cfg.JWT.RefreshExpirationHours)*time.Hour,
		)
		if err != nil {
			s.logger.Error("刷新会话失败", err)
		}

		// 删除旧会话
		err = s.sessionManager.DeleteSession(ctx, refreshToken)
		if err != nil {
			s.logger.Error("删除旧会话失败", err)
		}
	}

	return &user.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    s.cfg.JWT.ExpirationHours * 3600,
		TokenType:    "Bearer",
		User:         u,
	}, nil
}

// Register 处理用户注册
func (s *Service) Register(req *user.RegisterRequest) (*user.User, error) {
	// 检查邮箱是否已存在
	existingUser, err := s.userRepo.GetByEmail(req.Email, req.TenantID)
	if err == nil && existingUser != nil {
		return nil, errors.New("邮箱已被注册")
	}

	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	status := user.UserStatusActive
	if req.Status != nil {
		status = *req.Status
	}

	newUser := &user.User{
		ID:        uuid.New().String(),
		TenantID:  req.TenantID,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Name:      req.Name,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 保存用户
	err = s.userRepo.Create(newUser)
	if err != nil {
		return nil, err
	}

	// 不返回密码
	newUser.Password = ""

	return newUser, nil
}

// generateAccessToken 生成访问令牌
func (s *Service) generateAccessToken(u *user.User) (string, error) {
	expirationTime := time.Now().Add(time.Duration(s.cfg.JWT.ExpirationHours) * time.Hour)
	claims := jwt.MapClaims{
		"user_id":   u.ID,
		"tenant_id": u.TenantID,
		"email":     u.Email,
		"exp":       expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.Secret))
}

// generateRefreshToken 生成刷新令牌
func (s *Service) generateRefreshToken(u *user.User) (string, error) {
	expirationTime := time.Now().Add(time.Duration(s.cfg.JWT.RefreshExpirationHours) * time.Hour)
	claims := jwt.MapClaims{
		"user_id":   u.ID,
		"tenant_id": u.TenantID,
		"exp":       expirationTime.Unix(),
		"type":      "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.Secret))
}
