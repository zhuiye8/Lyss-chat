package providers

import (
	"context"

	"github.com/cloudwego/eino/schema"
	"github.com/your-org/lyss-chat-backend/internal/domain/model"
	"github.com/your-org/lyss-chat-backend/pkg/logger"
)

// Provider 定义模型提供商接口
type Provider interface {
	// GetName 获取提供商名称
	GetName() string

	// GetModels 获取提供商支持的模型列表
	GetModels() []*model.Model

	// Call 调用模型
	Call(ctx context.Context, modelID string, messages []*schema.Message, params map[string]interface{}) (*schema.Message, error)

	// Stream 流式调用模型
	Stream(ctx context.Context, modelID string, messages []*schema.Message, params map[string]interface{}) (<-chan *schema.Message, error)
}

// ProviderFactory 定义提供商工厂接口
type ProviderFactory interface {
	// Create 创建提供商实例
	Create(apiKey string, logger *logger.Logger) (Provider, error)
}

// ProviderRegistry 管理所有模型提供商
type ProviderRegistry struct {
	factories map[string]ProviderFactory
	providers map[string]Provider
	logger    *logger.Logger
}

// NewProviderRegistry 创建一个新的提供商注册表
func NewProviderRegistry(logger *logger.Logger) *ProviderRegistry {
	return &ProviderRegistry{
		factories: make(map[string]ProviderFactory),
		providers: make(map[string]Provider),
		logger:    logger,
	}
}

// RegisterFactory 注册提供商工厂
func (r *ProviderRegistry) RegisterFactory(providerID string, factory ProviderFactory) {
	r.factories[providerID] = factory
}

// GetProvider 获取提供商实例
func (r *ProviderRegistry) GetProvider(providerID string, apiKey string) (Provider, error) {
	// 检查是否已经创建了提供商实例
	if provider, ok := r.providers[providerID]; ok {
		return provider, nil
	}

	// 获取提供商工厂
	factory, ok := r.factories[providerID]
	if !ok {
		return nil, ErrProviderNotFound
	}

	// 创建提供商实例
	provider, err := factory.Create(apiKey, r.logger)
	if err != nil {
		return nil, err
	}

	// 缓存提供商实例
	r.providers[providerID] = provider

	return provider, nil
}

// GetSupportedProviders 获取支持的提供商列表
func (r *ProviderRegistry) GetSupportedProviders() []string {
	providers := make([]string, 0, len(r.factories))
	for providerID := range r.factories {
		providers = append(providers, providerID)
	}
	return providers
}

// Errors
var (
	ErrProviderNotFound = NewError("provider not found")
	ErrModelNotFound    = NewError("model not found")
	ErrInvalidAPIKey    = NewError("invalid API key")
	ErrAPICallFailed    = NewError("API call failed")
)

// Error 定义提供商错误
type Error struct {
	Message string
}

// NewError 创建一个新的错误
func NewError(message string) *Error {
	return &Error{Message: message}
}

// Error 实现 error 接口
func (e *Error) Error() string {
	return e.Message
}
