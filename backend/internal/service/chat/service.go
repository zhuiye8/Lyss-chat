package chat

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/your-org/lyss-chat-backend/internal/ai/graphs"
	"github.com/your-org/lyss-chat-backend/internal/domain/chat"
	"github.com/your-org/lyss-chat-backend/internal/domain/model"
	"github.com/your-org/lyss-chat-backend/internal/repository/postgres"
	"github.com/your-org/lyss-chat-backend/pkg/config"
	"github.com/your-org/lyss-chat-backend/pkg/db"
	"github.com/your-org/lyss-chat-backend/pkg/logger"
)

// Service 表示聊天服务
type Service struct {
	canvasRepo  chat.CanvasRepository
	messageRepo chat.MessageRepository
	modelRepo   model.ModelRepository
	aiGraphs    *graphs.ChatGraphs
	logger      *logger.Logger
}

// NewService 创建一个新的聊天服务
func NewService(
	database *db.Postgres,
	aiGraphs *graphs.ChatGraphs,
	cfg *config.Config,
	logger *logger.Logger,
) *Service {
	canvasRepo := postgres.NewCanvasRepository(database)
	messageRepo := postgres.NewMessageRepository(database)
	modelRepo := postgres.NewModelRepository(database)

	return &Service{
		canvasRepo:  canvasRepo,
		messageRepo: messageRepo,
		modelRepo:   modelRepo,
		aiGraphs:    aiGraphs,
		logger:      logger,
	}
}

// CreateCanvas 创建一个新的画布
func (s *Service) CreateCanvas(userID string, req *chat.CreateCanvasRequest) (*chat.Canvas, error) {
	// 设置默认状态
	status := chat.CanvasStatusActive

	// 创建画布
	canvas := &chat.Canvas{
		ID:          uuid.New().String(),
		WorkspaceID: req.WorkspaceID,
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Status:      status,
		ModelID:     req.ModelID,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 保存画布
	err := s.canvasRepo.Create(canvas)
	if err != nil {
		return nil, fmt.Errorf("创建画布失败: %w", err)
	}

	// 创建系统消息
	systemMessage := &chat.Message{
		ID:        uuid.New().String(),
		CanvasID:  canvas.ID,
		Role:      chat.MessageRoleSystem,
		Content:   "欢迎使用 Lyss Chat！我是您的 AI 助手，有什么可以帮您的吗？",
		CreatedBy: userID,
		CreatedAt: time.Now(),
	}

	// 保存系统消息
	err = s.messageRepo.Create(systemMessage)
	if err != nil {
		s.logger.Error("创建系统消息失败", err)
		// 继续处理，不要因为系统消息创建失败而阻止画布创建
	}

	return canvas, nil
}

// GetCanvas 获取画布
func (s *Service) GetCanvas(id string) (*chat.Canvas, error) {
	return s.canvasRepo.GetByID(id)
}

// UpdateCanvas 更新画布
func (s *Service) UpdateCanvas(id string, req *chat.UpdateCanvasRequest) (*chat.Canvas, error) {
	// 获取画布
	canvas, err := s.canvasRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Title != nil {
		canvas.Title = *req.Title
	}
	if req.Description != nil {
		canvas.Description = req.Description
	}
	if req.Status != nil {
		canvas.Status = *req.Status
	}
	if req.ModelID != nil {
		canvas.ModelID = req.ModelID
	}

	// 保存画布
	err = s.canvasRepo.Update(canvas)
	if err != nil {
		return nil, fmt.Errorf("更新画布失败: %w", err)
	}

	return canvas, nil
}

// DeleteCanvas 删除画布
func (s *Service) DeleteCanvas(id string) error {
	return s.canvasRepo.Delete(id)
}

// ListCanvases 列出画布
func (s *Service) ListCanvases(workspaceID string, canvasType *string, page, pageSize int) ([]*chat.Canvas, int, error) {
	offset := (page - 1) * pageSize
	return s.canvasRepo.List(workspaceID, canvasType, offset, pageSize)
}

// SendMessage 发送消息
func (s *Service) SendMessage(userID, canvasID string, req *chat.SendMessageRequest) (*chat.Message, error) {
	// 获取画布
	canvas, err := s.canvasRepo.GetByID(canvasID)
	if err != nil {
		return nil, err
	}

	// 创建用户消息
	userMessage := &chat.Message{
		ID:        uuid.New().String(),
		CanvasID:  canvasID,
		ParentID:  req.ParentID,
		Role:      chat.MessageRoleUser,
		Content:   req.Content,
		Metadata:  req.Metadata,
		CreatedBy: userID,
		CreatedAt: time.Now(),
	}

	// 保存用户消息
	err = s.messageRepo.Create(userMessage)
	if err != nil {
		return nil, fmt.Errorf("保存用户消息失败: %w", err)
	}

	// 获取对话历史
	var history []*chat.Message
	if req.ParentID != nil {
		history, err = s.messageRepo.GetConversation(*req.ParentID, 10)
		if err != nil {
			s.logger.Error("获取对话历史失败", err)
			// 继续处理，使用空历史
			history = []*chat.Message{}
		}
	}

	// 添加当前用户消息
	history = append(history, userMessage)

	// 转换为 Eino 消息格式
	einoMessages := make([]*schema.Message, 0, len(history))
	for _, msg := range history {
		einoMessages = append(einoMessages, &schema.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// 调用 AI 模型
	ctx := context.Background()
	var modelID string
	if canvas.ModelID != nil {
		modelID = *canvas.ModelID
	} else {
		// 使用默认模型
		modelID = "default"
	}

	// 准备输入
	input := map[string]any{
		"messages": einoMessages,
		"model_id": modelID,
	}

	// 调用 AI 图形
	aiResponse, err := s.aiGraphs.Chat.Invoke(ctx, input)
	if err != nil {
		s.logger.Error("调用 AI 模型失败", err)
		return nil, fmt.Errorf("调用 AI 模型失败: %w", err)
	}

	// 创建 AI 响应消息
	aiMessage := &chat.Message{
		ID:        uuid.New().String(),
		CanvasID:  canvasID,
		ParentID:  &userMessage.ID,
		Role:      chat.MessageRoleAssistant,
		Content:   aiResponse.Content,
		CreatedBy: userID,
		CreatedAt: time.Now(),
	}

	// 保存 AI 响应消息
	err = s.messageRepo.Create(aiMessage)
	if err != nil {
		s.logger.Error("保存 AI 响应消息失败", err)
		return nil, fmt.Errorf("保存 AI 响应消息失败: %w", err)
	}

	return aiMessage, nil
}

// GetMessages 获取消息
func (s *Service) GetMessages(canvasID string, page, pageSize int) ([]*chat.Message, int, error) {
	offset := (page - 1) * pageSize
	return s.messageRepo.GetByCanvasID(canvasID, offset, pageSize)
}

// StreamMessage 流式发送消息
func (s *Service) StreamMessage(userID, canvasID string, req *chat.SendMessageRequest) (<-chan *chat.Message, error) {
	// 获取画布
	canvas, err := s.canvasRepo.GetByID(canvasID)
	if err != nil {
		return nil, err
	}

	// 创建用户消息
	userMessage := &chat.Message{
		ID:        uuid.New().String(),
		CanvasID:  canvasID,
		ParentID:  req.ParentID,
		Role:      chat.MessageRoleUser,
		Content:   req.Content,
		Metadata:  req.Metadata,
		CreatedBy: userID,
		CreatedAt: time.Now(),
	}

	// 保存用户消息
	err = s.messageRepo.Create(userMessage)
	if err != nil {
		return nil, fmt.Errorf("保存用户消息失败: %w", err)
	}

	// 获取对话历史
	var history []*chat.Message
	if req.ParentID != nil {
		history, err = s.messageRepo.GetConversation(*req.ParentID, 10)
		if err != nil {
			s.logger.Error("获取对话历史失败", err)
			// 继续处理，使用空历史
			history = []*chat.Message{}
		}
	}

	// 添加当前用户消息
	history = append(history, userMessage)

	// 转换为 Eino 消息格式
	einoMessages := make([]*schema.Message, 0, len(history))
	for _, msg := range history {
		einoMessages = append(einoMessages, &schema.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// 调用 AI 模型
	ctx := context.Background()
	var modelID string
	if canvas.ModelID != nil {
		modelID = *canvas.ModelID
	} else {
		// 使用默认模型
		modelID = "default"
	}

	// 准备输入
	input := map[string]any{
		"messages": einoMessages,
		"model_id": modelID,
	}

	// 创建结果通道
	resultChan := make(chan *chat.Message)

	// 启动 goroutine 处理流式响应
	go func() {
		defer close(resultChan)

		// 调用 AI 图形流式接口
		aiResponseChan, err := s.aiGraphs.Chat.Stream(ctx, input)
		if err != nil {
			s.logger.Error("调用 AI 模型流式接口失败", err)
			return
		}

		// 创建 AI 响应消息
		aiMessage := &chat.Message{
			ID:        uuid.New().String(),
			CanvasID:  canvasID,
			ParentID:  &userMessage.ID,
			Role:      chat.MessageRoleAssistant,
			Content:   "",
			CreatedBy: userID,
			CreatedAt: time.Now(),
		}

		// 处理流式响应
		var fullContent string
		for chunk := range aiResponseChan {
			fullContent += chunk.Content
			
			// 更新消息内容
			aiMessage.Content = fullContent
			
			// 发送到结果通道
			resultChan <- &chat.Message{
				ID:        aiMessage.ID,
				CanvasID:  aiMessage.CanvasID,
				ParentID:  aiMessage.ParentID,
				Role:      aiMessage.Role,
				Content:   aiMessage.Content,
				CreatedBy: aiMessage.CreatedBy,
				CreatedAt: aiMessage.CreatedAt,
			}
		}

		// 保存完整的 AI 响应消息
		aiMessage.Content = fullContent
		err = s.messageRepo.Create(aiMessage)
		if err != nil {
			s.logger.Error("保存 AI 响应消息失败", err)
		}
	}()

	return resultChan, nil
}
