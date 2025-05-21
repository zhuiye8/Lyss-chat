package components

import (
	"context"

	"github.com/cloudwego/eino/schema"
	"github.com/your-org/lyss-chat-backend/internal/domain/model"
)

// ModelService 表示模型服务接口
type ModelService interface {
	GetModel(ctx context.Context, id string) (*model.Model, error)
	GetAPIKey(ctx context.Context, providerID string) (*model.APIKey, error)
	CallModel(ctx context.Context, modelID string, messages []*schema.Message) (*schema.Message, error)
	StreamModel(ctx context.Context, modelID string, messages []*schema.Message) (<-chan *schema.Message, error)
}

// ModelServiceAdapter 是模型服务的适配器
type ModelServiceAdapter struct {
	ModelService ModelService
}

// GetModel 实现 ModelProvider 接口的 GetModel 方法
func (a *ModelServiceAdapter) GetModel(ctx context.Context, modelID string) (*model.Model, error) {
	return a.ModelService.GetModel(ctx, modelID)
}

// GetAPIKey 实现 ModelProvider 接口的 GetAPIKey 方法
func (a *ModelServiceAdapter) GetAPIKey(ctx context.Context, providerID string) (*model.APIKey, error) {
	return a.ModelService.GetAPIKey(ctx, providerID)
}

// CallModel 实现 ModelProvider 接口的 CallModel 方法
func (a *ModelServiceAdapter) CallModel(ctx context.Context, modelID string, messages []*schema.Message) (*schema.Message, error) {
	return a.ModelService.CallModel(ctx, modelID, messages)
}

// StreamModel 实现 ModelProvider 接口的 StreamModel 方法
func (a *ModelServiceAdapter) StreamModel(ctx context.Context, modelID string, messages []*schema.Message) (<-chan *schema.Message, error) {
	return a.ModelService.StreamModel(ctx, modelID, messages)
}
