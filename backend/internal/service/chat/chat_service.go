package chat

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/lyss-chat-backend/internal/domain/chat"
	"github.com/your-org/lyss-chat-backend/pkg/logger"
)

// CanvasRepository 定义画布仓储接口
type CanvasRepository interface {
	Create(canvas *chat.Canvas) error
	GetByID(id string) (*chat.Canvas, error)
	Update(canvas *chat.Canvas) error
	Delete(id string) error
	List(workspaceID string, canvasType *string, offset, limit int) ([]*chat.Canvas, int, error)
}

// MessageRepository 定义消息仓储接口
type MessageRepository interface {
	Create(message *chat.Message) error
	GetByID(id string) (*chat.Message, error)
	ListByCanvasID(canvasID string, offset, limit int) ([]*chat.Message, int, error)
}

// ChatService 实现聊天相关的业务逻辑
type ChatService struct {
	canvasRepo  CanvasRepository
	messageRepo MessageRepository
	logger      *logger.Logger
}

// NewChatService 创建一个新的 ChatService 实例
func NewChatService(canvasRepo CanvasRepository, messageRepo MessageRepository, logger *logger.Logger) *ChatService {
	return &ChatService{
		canvasRepo:  canvasRepo,
		messageRepo: messageRepo,
		logger:      logger,
	}
}

// CreateCanvas 创建一个新的画布
func (s *ChatService) CreateCanvas(ctx context.Context, workspaceID, title, description string, canvasType string, modelID *string, createdBy string) (*chat.Canvas, error) {
	// 验证画布类型
	if canvasType != chat.CanvasTypeChat && canvasType != chat.CanvasTypeCode {
		return nil, fmt.Errorf("无效的画布类型: %s", canvasType)
	}

	// 创建画布对象
	now := time.Now().UTC().Format(time.RFC3339)
	canvas := &chat.Canvas{
		ID:          uuid.New().String(),
		WorkspaceID: workspaceID,
		Title:       title,
		Description: description,
		Type:        canvasType,
		Status:      chat.CanvasStatusActive,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// 设置可选的模型 ID
	if modelID != nil {
		canvas.ModelID = *modelID
	}

	// 保存到数据库
	err := s.canvasRepo.Create(canvas)
	if err != nil {
		return nil, err
	}

	return canvas, nil
}

// GetCanvas 获取画布详情
func (s *ChatService) GetCanvas(ctx context.Context, id string) (*chat.Canvas, error) {
	return s.canvasRepo.GetByID(id)
}

// UpdateCanvas 更新画布
func (s *ChatService) UpdateCanvas(ctx context.Context, id, title, description string, status string, modelID *string) (*chat.Canvas, error) {
	// 获取现有画布
	canvas, err := s.canvasRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if title != "" {
		canvas.Title = title
	}
	if description != "" {
		canvas.Description = description
	}
	if status != "" {
		if status != chat.CanvasStatusActive && status != chat.CanvasStatusArchived {
			return nil, fmt.Errorf("无效的画布状态: %s", status)
		}
		canvas.Status = status
	}
	if modelID != nil {
		canvas.ModelID = *modelID
	}

	// 更新时间
	canvas.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	// 保存到数据库
	err = s.canvasRepo.Update(canvas)
	if err != nil {
		return nil, err
	}

	return canvas, nil
}

// DeleteCanvas 删除画布
func (s *ChatService) DeleteCanvas(ctx context.Context, id string) error {
	// 检查画布是否存在
	_, err := s.canvasRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 删除画布
	return s.canvasRepo.Delete(id)
}

// ListCanvases 获取画布列表
func (s *ChatService) ListCanvases(ctx context.Context, workspaceID string, canvasType *string, page, pageSize int) ([]*chat.Canvas, int, error) {
	// 计算偏移量
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	// 获取画布列表
	return s.canvasRepo.List(workspaceID, canvasType, offset, pageSize)
}

// CreateMessage 创建一个新的消息
func (s *ChatService) CreateMessage(ctx context.Context, canvasID, parentID, role, content string, metadata map[string]interface{}) (*chat.Message, error) {
	// 验证画布是否存在
	_, err := s.canvasRepo.GetByID(canvasID)
	if err != nil {
		return nil, err
	}

	// 验证角色
	if role != chat.MessageRoleUser && role != chat.MessageRoleAssistant && role != chat.MessageRoleSystem {
		return nil, fmt.Errorf("无效的消息角色: %s", role)
	}

	// 如果指定了父消息 ID，验证它是否存在
	if parentID != "" {
		_, err := s.messageRepo.GetByID(parentID)
		if err != nil {
			return nil, err
		}
	}

	// 创建消息对象
	message := &chat.Message{
		ID:        uuid.New().String(),
		CanvasID:  canvasID,
		ParentID:  parentID,
		Role:      role,
		Content:   content,
		Metadata:  metadata,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	// 保存到数据库
	err = s.messageRepo.Create(message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// GetMessage 获取消息详情
func (s *ChatService) GetMessage(ctx context.Context, id string) (*chat.Message, error) {
	return s.messageRepo.GetByID(id)
}

// ListMessages 获取画布下的消息列表
func (s *ChatService) ListMessages(ctx context.Context, canvasID string, page, pageSize int) ([]*chat.Message, int, error) {
	// 验证画布是否存在
	_, err := s.canvasRepo.GetByID(canvasID)
	if err != nil {
		return nil, 0, err
	}

	// 计算偏移量
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	// 获取消息列表
	return s.messageRepo.ListByCanvasID(canvasID, offset, pageSize)
}
