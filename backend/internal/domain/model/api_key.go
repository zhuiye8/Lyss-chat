package model

import (
	"time"
)

// APIKey 表示 API 密钥实体
type APIKey struct {
	ID         string    `json:"id" db:"id"`
	ProviderID string    `json:"provider_id" db:"provider_id"`
	Name       string    `json:"name" db:"name"`
	Key        string    `json:"-" db:"key"`
	Status     string    `json:"status" db:"status"`
	CreatedBy  string    `json:"created_by" db:"created_by"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// APIKeyStatus 表示 API 密钥状态
const (
	APIKeyStatusActive   = "active"
	APIKeyStatusInactive = "inactive"
)

// CreateAPIKeyRequest 表示创建 API 密钥的请求
type CreateAPIKeyRequest struct {
	ProviderID string `json:"provider_id" validate:"required,uuid"`
	Name       string `json:"name" validate:"required"`
	Key        string `json:"key" validate:"required"`
}

// UpdateAPIKeyRequest 表示更新 API 密钥的请求
type UpdateAPIKeyRequest struct {
	Name   *string `json:"name,omitempty"`
	Key    *string `json:"key,omitempty"`
	Status *string `json:"status,omitempty" validate:"omitempty,oneof=active inactive"`
}

// APIKeyResponse 表示 API 密钥响应
type APIKeyResponse struct {
	ID         string    `json:"id"`
	ProviderID string    `json:"provider_id"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	CreatedBy  string    `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// APIKeyRepository 表示 API 密钥仓库接口
type APIKeyRepository interface {
	Create(apiKey *APIKey) error
	GetByID(id string) (*APIKey, error)
	GetByProviderID(providerID string) ([]*APIKey, error)
	Update(apiKey *APIKey) error
	Delete(id string) error
}

// APIKeyService 表示 API 密钥服务接口
type APIKeyService interface {
	Create(userID string, req *CreateAPIKeyRequest) (*APIKeyResponse, error)
	GetByID(id string) (*APIKeyResponse, error)
	GetByProviderID(providerID string) ([]*APIKeyResponse, error)
	Update(id string, req *UpdateAPIKeyRequest) (*APIKeyResponse, error)
	Delete(id string) error
}
