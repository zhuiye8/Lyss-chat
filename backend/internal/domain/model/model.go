package model

import (
	"encoding/json"
	"time"
)

// Model 表示模型实体
type Model struct {
	ID          string          `json:"id" db:"id"`
	ProviderID  string          `json:"provider_id" db:"provider_id"`
	ModelID     string          `json:"model_id" db:"model_id"`
	Name        string          `json:"name" db:"name"`
	Description *string         `json:"description,omitempty" db:"description"`
	Capabilities json.RawMessage `json:"capabilities" db:"capabilities"`
	Parameters  json.RawMessage `json:"parameters" db:"parameters"`
	Status      string          `json:"status" db:"status"`
	IsPublic    bool            `json:"is_public" db:"is_public"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

// ModelStatus 表示模型状态
const (
	ModelStatusActive   = "active"
	ModelStatusInactive = "inactive"
)

// CreateModelRequest 表示创建模型的请求
type CreateModelRequest struct {
	ProviderID   string          `json:"provider_id" validate:"required,uuid"`
	ModelID      string          `json:"model_id" validate:"required"`
	Name         string          `json:"name" validate:"required"`
	Description  *string         `json:"description,omitempty"`
	Capabilities json.RawMessage `json:"capabilities" validate:"required"`
	Parameters   json.RawMessage `json:"parameters" validate:"required"`
	IsPublic     bool            `json:"is_public"`
}

// UpdateModelRequest 表示更新模型的请求
type UpdateModelRequest struct {
	Name         *string         `json:"name,omitempty"`
	Description  *string         `json:"description,omitempty"`
	Capabilities *json.RawMessage `json:"capabilities,omitempty"`
	Parameters   *json.RawMessage `json:"parameters,omitempty"`
	Status       *string         `json:"status,omitempty" validate:"omitempty,oneof=active inactive"`
	IsPublic     *bool           `json:"is_public,omitempty"`
}

// ModelRepository 表示模型仓库接口
type ModelRepository interface {
	Create(model *Model) error
	GetByID(id string) (*Model, error)
	Update(model *Model) error
	Delete(id string) error
	List(providerID *string, status *string, offset, limit int) ([]*Model, int, error)
}

// ModelService 表示模型服务接口
type ModelService interface {
	Create(req *CreateModelRequest) (*Model, error)
	GetByID(id string) (*Model, error)
	Update(id string, req *UpdateModelRequest) (*Model, error)
	Delete(id string) error
	List(providerID *string, status *string, page, pageSize int) ([]*Model, int, error)
}
