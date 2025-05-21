package chat

import (
	"time"
)

// Canvas 表示画布实体
type Canvas struct {
	ID          string     `json:"id" db:"id"`
	WorkspaceID string     `json:"workspace_id" db:"workspace_id"`
	Title       string     `json:"title" db:"title"`
	Description *string    `json:"description,omitempty" db:"description"`
	Type        string     `json:"type" db:"type"`
	Status      string     `json:"status" db:"status"`
	ModelID     *string    `json:"model_id,omitempty" db:"model_id"`
	CreatedBy   string     `json:"created_by" db:"created_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// CanvasType 表示画布类型
const (
	CanvasTypeChat = "chat"
	CanvasTypeCode = "code"
)

// CanvasStatus 表示画布状态
const (
	CanvasStatusActive   = "active"
	CanvasStatusArchived = "archived"
)

// CreateCanvasRequest 表示创建画布的请求
type CreateCanvasRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description *string `json:"description,omitempty"`
	WorkspaceID string  `json:"workspace_id" validate:"required,uuid"`
	Type        string  `json:"type" validate:"required,oneof=chat code"`
	ModelID     *string `json:"model_id,omitempty" validate:"omitempty,uuid"`
}

// UpdateCanvasRequest 表示更新画布的请求
type UpdateCanvasRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty" validate:"omitempty,oneof=active archived"`
	ModelID     *string `json:"model_id,omitempty" validate:"omitempty,uuid"`
}

// CanvasRepository 表示画布仓库接口
type CanvasRepository interface {
	Create(canvas *Canvas) error
	GetByID(id string) (*Canvas, error)
	Update(canvas *Canvas) error
	Delete(id string) error
	List(workspaceID string, canvasType *string, offset, limit int) ([]*Canvas, int, error)
}

// CanvasService 表示画布服务接口
type CanvasService interface {
	Create(userID string, req *CreateCanvasRequest) (*Canvas, error)
	GetByID(id string) (*Canvas, error)
	Update(id string, req *UpdateCanvasRequest) (*Canvas, error)
	Delete(id string) error
	List(workspaceID string, canvasType *string, page, pageSize int) ([]*Canvas, int, error)
}
