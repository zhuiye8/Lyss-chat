package user

import (
	"time"
)

// User 表示用户实体
type User struct {
	ID        string    `json:"id" db:"id"`
	TenantID  string    `json:"tenant_id" db:"tenant_id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	Name      string    `json:"name" db:"name"`
	AvatarURL *string   `json:"avatar_url,omitempty" db:"avatar_url"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserStatus 表示用户状态
const (
	UserStatusActive    = "active"
	UserStatusInactive  = "inactive"
	UserStatusSuspended = "suspended"
)

// CreateUserRequest 表示创建用户的请求
type CreateUserRequest struct {
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=8"`
	Name     string  `json:"name" validate:"required"`
	TenantID string  `json:"tenant_id" validate:"required,uuid"`
	Status   *string `json:"status,omitempty"`
}

// UpdateUserRequest 表示更新用户的请求
type UpdateUserRequest struct {
	Name      *string `json:"name,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Status    *string `json:"status,omitempty"`
}

// LoginRequest 表示登录请求
type LoginRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	TenantID  string `json:"tenant_id" validate:"required,uuid"`
	IP        string `json:"-"` // 由服务器填充，不从客户端接收
	UserAgent string `json:"-"` // 由服务器填充，不从客户端接收
}

// RegisterRequest 表示注册请求
type RegisterRequest struct {
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=8"`
	Name      string  `json:"name" validate:"required"`
	TenantID  string  `json:"tenant_id" validate:"required,uuid"`
	Status    *string `json:"status,omitempty"`
	IP        string  `json:"-"` // 由服务器填充，不从客户端接收
	UserAgent string  `json:"-"` // 由服务器填充，不从客户端接收
}

// LoginResponse 表示登录响应
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	User         *User  `json:"user"`
}

// Repository 表示用户仓库接口
type Repository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email, tenantID string) (*User, error)
	Update(user *User) error
	Delete(id string) error
	List(tenantID string, offset, limit int) ([]*User, int, error)
}

// Service 表示用户服务接口
type Service interface {
	Create(req *CreateUserRequest) (*User, error)
	GetByID(id string) (*User, error)
	Update(id string, req *UpdateUserRequest) (*User, error)
	Delete(id string) error
	List(tenantID string, page, pageSize int) ([]*User, int, error)
	Login(req *LoginRequest) (*LoginResponse, error)
	RefreshToken(refreshToken string) (*LoginResponse, error)
}
