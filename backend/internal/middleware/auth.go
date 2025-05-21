package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/your-org/lyss-chat-2.0/backend/pkg/config"
)

// contextKey 是用于上下文的键类型
type contextKey string

// UserIDKey 是用户 ID 的上下文键
const UserIDKey contextKey = "user_id"

// TenantIDKey 是租户 ID 的上下文键
const TenantIDKey contextKey = "tenant_id"

// Auth 创建一个认证中间件
func Auth(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 从请求头获取令牌
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "未提供认证令牌", http.StatusUnauthorized)
				return
			}

			// 解析令牌
			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
			claims, err := validateToken(tokenString, cfg.JWT.Secret)
			if err != nil {
				http.Error(w, "无效的认证令牌", http.StatusUnauthorized)
				return
			}

			// 将用户 ID 和租户 ID 添加到上下文
			userID, ok := claims["user_id"].(string)
			if !ok {
				http.Error(w, "无效的认证令牌", http.StatusUnauthorized)
				return
			}

			tenantID, ok := claims["tenant_id"].(string)
			if !ok {
				http.Error(w, "无效的认证令牌", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, TenantIDKey, tenantID)

			// 处理请求
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// validateToken 验证 JWT 令牌
func validateToken(tokenString, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名算法")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("无效的令牌")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("无效的令牌声明")
	}

	return claims, nil
}

// GetUserID 从上下文获取用户 ID
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

// GetTenantID 从上下文获取租户 ID
func GetTenantID(ctx context.Context) (string, bool) {
	tenantID, ok := ctx.Value(TenantIDKey).(string)
	return tenantID, ok
}
