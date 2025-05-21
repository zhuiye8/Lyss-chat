package model

import (
	"time"
)

// Provider 表示提供商实体
type Provider struct {
	ID          string     `json:"id" db:"id"`
	TenantID    string     `json:"tenant_id" db:"tenant_id"`
	Code        string     `json:"code" db:"code"`
	Name        string     `json:"name" db:"name"`
	Description *string    `json:"description,omitempty" db:"description"`
	BaseURL     *string    `json:"base_url,omitempty" db:"base_url"`
	Status      string     `json:"status" db:"status"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// ProviderStatus 表示提供商状态
const (
	ProviderStatusActive   = "active"
	ProviderStatusInactive = "inactive"
)

// CreateProviderRequest 表示创建提供商的请求
type CreateProviderRequest struct {
	TenantID    string  `json:"tenant_id" validate:"required,uuid"`
	Code        string  `json:"code" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
	BaseURL     *string `json:"base_url,omitempty"`
}

// UpdateProviderRequest 表示更新提供商的请求
type UpdateProviderRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	BaseURL     *string `json:"base_url,omitempty"`
	Status      *string `json:"status,omitempty" validate:"omitempty,oneof=active inactive"`
}

// ProviderRepository 表示提供商仓库接口
type ProviderRepository interface {
	Create(provider *Provider) error
	GetByID(id string) (*Provider, error)
	GetByCode(code, tenantID string) (*Provider, error)
	Update(provider *Provider) error
	Delete(id string) error
	List(tenantID string, offset, limit int) ([]*Provider, int, error)
}

// ProviderService 表示提供商服务接口
type ProviderService interface {
	Create(req *CreateProviderRequest) (*Provider, error)
	GetByID(id string) (*Provider, error)
	Update(id string, req *UpdateProviderRequest) (*Provider, error)
	Delete(id string) error
	List(tenantID string, page, pageSize int) ([]*Provider, int, error)
}
