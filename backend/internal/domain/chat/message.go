package chat

import (
	"encoding/json"
	"time"
)

// Message 表示消息实体
type Message struct {
	ID         string          `json:"id" db:"id"`
	CanvasID   string          `json:"canvas_id" db:"canvas_id"`
	ParentID   *string         `json:"parent_id,omitempty" db:"parent_id"`
	Role       string          `json:"role" db:"role"`
	Content    string          `json:"content" db:"content"`
	Metadata   json.RawMessage `json:"metadata,omitempty" db:"metadata"`
	TokenCount *int            `json:"token_count,omitempty" db:"token_count"`
	CreatedBy  string          `json:"created_by" db:"created_by"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
}

// MessageRole 表示消息角色
const (
	MessageRoleUser      = "user"
	MessageRoleAssistant = "assistant"
	MessageRoleSystem    = "system"
)

// SendMessageRequest 表示发送消息的请求
type SendMessageRequest struct {
	Content  string          `json:"content" validate:"required"`
	ParentID *string         `json:"parent_id,omitempty" validate:"omitempty,uuid"`
	Metadata json.RawMessage `json:"metadata,omitempty"`
}

// MessageRepository 表示消息仓库接口
type MessageRepository interface {
	Create(message *Message) error
	GetByID(id string) (*Message, error)
	GetByCanvasID(canvasID string, offset, limit int) ([]*Message, int, error)
	GetConversation(messageID string, limit int) ([]*Message, error)
}

// MessageService 表示消息服务接口
type MessageService interface {
	Send(userID, canvasID string, req *SendMessageRequest) (*Message, error)
	GetByID(id string) (*Message, error)
	GetByCanvasID(canvasID string, page, pageSize int) ([]*Message, int, error)
	GetConversation(messageID string, limit int) ([]*Message, error)
	StreamResponse(userID, canvasID string, req *SendMessageRequest) (<-chan *Message, error)
}
