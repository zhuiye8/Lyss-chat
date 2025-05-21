package graphs

import (
	"context"

	"github.com/your-org/lyss-chat-backend/internal/ai/components"
	"github.com/your-org/lyss-chat-backend/internal/service/model"
	"github.com/your-org/lyss-chat-backend/pkg/logger"
)

// ChatGraphs 包含所有聊天图形
type ChatGraphs struct {
	Chat *ChatGraph
}

// NewChatGraphs 创建一个新的聊天图形集合
func NewChatGraphs(ctx context.Context, modelService *model.Service, logger *logger.Logger) (*ChatGraphs, error) {
	// 创建模型提供商适配器
	modelProvider := &components.ModelServiceAdapter{
		ModelService: modelService,
	}

	// 创建聊天模型适配器
	chatModel := components.NewChatModelAdapter(modelProvider, "", logger)

	// 创建聊天图形
	chatGraph, err := NewChatGraph(ctx, *chatModel, logger)
	if err != nil {
		return nil, err
	}

	return &ChatGraphs{
		Chat: chatGraph,
	}, nil
}
