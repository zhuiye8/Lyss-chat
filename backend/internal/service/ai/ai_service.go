package ai

import (
	"context"
	"time"

	"github.com/cloudwego/eino"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/your-org/lyss-chat-backend/internal/ai/components"
	"github.com/your-org/lyss-chat-backend/internal/domain/chat"
	"github.com/your-org/lyss-chat-backend/pkg/logger"
)

// ModelService 定义模型服务接口
type ModelService interface {
	GetModel(ctx context.Context, id string) (*model.Model, error)
	GetAPIKey(ctx context.Context, providerID string) (*model.APIKey, error)
	CallModel(ctx context.Context, modelID string, messages []*schema.Message) (*schema.Message, error)
	StreamModel(ctx context.Context, modelID string, messages []*schema.Message) (<-chan *schema.Message, error)
}

// MessageRepository 定义消息仓储接口
type MessageRepository interface {
	Create(message *chat.Message) error
	GetByID(id string) (*chat.Message, error)
	ListByCanvasID(canvasID string, offset, limit int) ([]*chat.Message, int, error)
}

// CanvasRepository 定义画布仓储接口
type CanvasRepository interface {
	GetByID(id string) (*chat.Canvas, error)
}

// AIService 实现 AI 相关的业务逻辑
type AIService struct {
	modelService   ModelService
	messageRepo    MessageRepository
	canvasRepo     CanvasRepository
	einoClient     *eino.Client
	defaultModelID string
	logger         *logger.Logger
}

// NewAIService 创建一个新的 AIService 实例
func NewAIService(
	modelService ModelService,
	messageRepo MessageRepository,
	canvasRepo CanvasRepository,
	einoClient *eino.Client,
	defaultModelID string,
	logger *logger.Logger,
) *AIService {
	return &AIService{
		modelService:   modelService,
		messageRepo:    messageRepo,
		canvasRepo:     canvasRepo,
		einoClient:     einoClient,
		defaultModelID: defaultModelID,
		logger:         logger,
	}
}

// GenerateResponse 生成 AI 响应
func (s *AIService) GenerateResponse(ctx context.Context, canvasID string, messages []*chat.Message, modelID string, stream bool) (*chat.Message, error) {
	// 获取画布
	canvas, err := s.canvasRepo.GetByID(canvasID)
	if err != nil {
		return nil, err
	}

	// 如果未指定模型 ID，使用画布的模型 ID 或默认模型 ID
	if modelID == "" {
		if canvas.ModelID != "" {
			modelID = canvas.ModelID
		} else {
			modelID = s.defaultModelID
		}
	}

	// 转换消息格式
	einoMessages := convertToEinoMessages(messages)

	// 创建模型适配器
	modelAdapter := components.NewChatModelAdapter(
		&components.ModelServiceAdapter{ModelService: s.modelService},
		modelID,
		s.logger,
	)

	// 调用模型
	var response *schema.Message
	if stream {
		// 流式调用
		responseChan, err := modelAdapter.Stream(ctx, einoMessages)
		if err != nil {
			return nil, err
		}

		// 获取最后一个响应
		var lastResponse *schema.Message
		for resp := range responseChan {
			lastResponse = resp
		}
		response = lastResponse
	} else {
		// 非流式调用
		response, err = modelAdapter.Call(ctx, einoMessages)
		if err != nil {
			return nil, err
		}
	}

	// 创建 AI 消息
	aiMessage := &chat.Message{
		ID:        uuid.New().String(),
		CanvasID:  canvasID,
		Role:      chat.MessageRoleAssistant,
		Content:   response.Content,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	// 保存到数据库
	err = s.messageRepo.Create(aiMessage)
	if err != nil {
		return nil, err
	}

	return aiMessage, nil
}

// StreamResponse 流式生成 AI 响应
func (s *AIService) StreamResponse(ctx context.Context, canvasID string, messages []*chat.Message, modelID string) (<-chan *schema.Message, error) {
	// 获取画布
	canvas, err := s.canvasRepo.GetByID(canvasID)
	if err != nil {
		return nil, err
	}

	// 如果未指定模型 ID，使用画布的模型 ID 或默认模型 ID
	if modelID == "" {
		if canvas.ModelID != "" {
			modelID = canvas.ModelID
		} else {
			modelID = s.defaultModelID
		}
	}

	// 转换消息格式
	einoMessages := convertToEinoMessages(messages)

	// 创建模型适配器
	modelAdapter := components.NewChatModelAdapter(
		&components.ModelServiceAdapter{ModelService: s.modelService},
		modelID,
		s.logger,
	)

	// 流式调用模型
	responseChan, err := modelAdapter.Stream(ctx, einoMessages)
	if err != nil {
		return nil, err
	}

	// 创建输出通道
	outputChan := make(chan *schema.Message)

	// 启动 goroutine 处理流式响应
	go func() {
		defer close(outputChan)

		var lastResponse *schema.Message
		for resp := range responseChan {
			// 发送响应
			outputChan <- resp
			lastResponse = resp
		}

		// 保存最终响应到数据库
		if lastResponse != nil {
			aiMessage := &chat.Message{
				ID:        uuid.New().String(),
				CanvasID:  canvasID,
				Role:      chat.MessageRoleAssistant,
				Content:   lastResponse.Content,
				CreatedAt: time.Now().UTC().Format(time.RFC3339),
			}

			err := s.messageRepo.Create(aiMessage)
			if err != nil {
				s.logger.Error("保存 AI 消息失败", err)
			}
		}
	}()

	return outputChan, nil
}

// convertToEinoMessages 将 chat.Message 转换为 schema.Message
func convertToEinoMessages(messages []*chat.Message) []*schema.Message {
	einoMessages := make([]*schema.Message, 0, len(messages))
	for _, msg := range messages {
		einoMessages = append(einoMessages, &schema.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	return einoMessages
}
