package model

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/your-org/lyss-chat-backend/internal/ai/providers"
	"github.com/your-org/lyss-chat-backend/internal/domain/chat"
	"github.com/your-org/lyss-chat-backend/internal/domain/model"
	"github.com/your-org/lyss-chat-backend/pkg/logger"
)

// ModelRepository 定义模型仓储接口
type ModelRepository interface {
	Create(model *model.Model) error
	GetByID(id string) (*model.Model, error)
	Update(model *model.Model) error
	Delete(id string) error
	List(providerID *string, status *string, offset, limit int) ([]*model.Model, int, error)
}

// ProviderRepository 定义提供商仓储接口
type ProviderRepository interface {
	Create(provider *model.Provider) error
	GetByID(id string) (*model.Provider, error)
	Update(provider *model.Provider) error
	Delete(id string) error
	List(offset, limit int) ([]*model.Provider, int, error)
}

// APIKeyRepository 定义 API 密钥仓储接口
type APIKeyRepository interface {
	Create(apiKey *model.APIKey) error
	GetByID(id string) (*model.APIKey, error)
	GetByProviderID(providerID string) (*model.APIKey, error)
	Update(apiKey *model.APIKey) error
	Delete(id string) error
	List(providerID *string, offset, limit int) ([]*model.APIKey, int, error)
}

// MessageRepository 定义消息仓储接口
type MessageRepository interface {
	Create(message *chat.Message) error
	GetByID(id string) (*chat.Message, error)
	ListByCanvasID(canvasID string, offset, limit int) ([]*chat.Message, int, error)
}

// ModelService 实现模型相关的业务逻辑
type ModelService struct {
	modelRepo    ModelRepository
	providerRepo ProviderRepository
	apiKeyRepo   APIKeyRepository
	messageRepo  MessageRepository
	registry     *providers.ProviderRegistry
	logger       *logger.Logger
}

// NewModelService 创建一个新的 ModelService 实例
func NewModelService(
	modelRepo ModelRepository,
	providerRepo ProviderRepository,
	apiKeyRepo APIKeyRepository,
	messageRepo MessageRepository,
	registry *providers.ProviderRegistry,
	logger *logger.Logger,
) *ModelService {
	return &ModelService{
		modelRepo:    modelRepo,
		providerRepo: providerRepo,
		apiKeyRepo:   apiKeyRepo,
		messageRepo:  messageRepo,
		registry:     registry,
		logger:       logger,
	}
}

// CreateProvider 创建一个新的提供商
func (s *ModelService) CreateProvider(ctx context.Context, name, description string) (*model.Provider, error) {
	// 创建提供商对象
	now := time.Now().UTC().Format(time.RFC3339)
	provider := &model.Provider{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Status:      model.ProviderStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// 保存到数据库
	err := s.providerRepo.Create(provider)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// GetProvider 获取提供商详情
func (s *ModelService) GetProvider(ctx context.Context, id string) (*model.Provider, error) {
	return s.providerRepo.GetByID(id)
}

// UpdateProvider 更新提供商
func (s *ModelService) UpdateProvider(ctx context.Context, id, name, description string, status string) (*model.Provider, error) {
	// 获取现有提供商
	provider, err := s.providerRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if name != "" {
		provider.Name = name
	}
	if description != "" {
		provider.Description = description
	}
	if status != "" {
		if status != model.ProviderStatusActive && status != model.ProviderStatusInactive {
			return nil, fmt.Errorf("无效的提供商状态: %s", status)
		}
		provider.Status = status
	}

	// 更新时间
	provider.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	// 保存到数据库
	err = s.providerRepo.Update(provider)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// DeleteProvider 删除提供商
func (s *ModelService) DeleteProvider(ctx context.Context, id string) error {
	// 检查提供商是否存在
	_, err := s.providerRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 删除提供商
	return s.providerRepo.Delete(id)
}

// ListProviders 获取提供商列表
func (s *ModelService) ListProviders(ctx context.Context, page, pageSize int) ([]*model.Provider, int, error) {
	// 计算偏移量
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	// 获取提供商列表
	return s.providerRepo.List(offset, pageSize)
}

// CreateAPIKey 创建一个新的 API 密钥
func (s *ModelService) CreateAPIKey(ctx context.Context, providerID, name, key string) (*model.APIKey, error) {
	// 验证提供商是否存在
	_, err := s.providerRepo.GetByID(providerID)
	if err != nil {
		return nil, err
	}

	// 创建 API 密钥对象
	now := time.Now().UTC().Format(time.RFC3339)
	apiKey := &model.APIKey{
		ID:         uuid.New().String(),
		ProviderID: providerID,
		Name:       name,
		Key:        key,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// 保存到数据库
	err = s.apiKeyRepo.Create(apiKey)
	if err != nil {
		return nil, err
	}

	return apiKey, nil
}

// GetAPIKey 获取 API 密钥详情
func (s *ModelService) GetAPIKey(ctx context.Context, id string) (*model.APIKey, error) {
	return s.apiKeyRepo.GetByID(id)
}

// GetAPIKeyByProviderID 根据提供商 ID 获取 API 密钥
func (s *ModelService) GetAPIKeyByProviderID(ctx context.Context, providerID string) (*model.APIKey, error) {
	return s.apiKeyRepo.GetByProviderID(providerID)
}

// UpdateAPIKey 更新 API 密钥
func (s *ModelService) UpdateAPIKey(ctx context.Context, id, name, key string) (*model.APIKey, error) {
	// 获取现有 API 密钥
	apiKey, err := s.apiKeyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if name != "" {
		apiKey.Name = name
	}
	if key != "" {
		apiKey.Key = key
	}

	// 更新时间
	apiKey.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	// 保存到数据库
	err = s.apiKeyRepo.Update(apiKey)
	if err != nil {
		return nil, err
	}

	return apiKey, nil
}

// DeleteAPIKey 删除 API 密钥
func (s *ModelService) DeleteAPIKey(ctx context.Context, id string) error {
	// 检查 API 密钥是否存在
	_, err := s.apiKeyRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 删除 API 密钥
	return s.apiKeyRepo.Delete(id)
}

// ListAPIKeys 获取 API 密钥列表
func (s *ModelService) ListAPIKeys(ctx context.Context, providerID *string, page, pageSize int) ([]*model.APIKey, int, error) {
	// 计算偏移量
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	// 获取 API 密钥列表
	return s.apiKeyRepo.List(providerID, offset, pageSize)
}

// CreateModel 创建一个新的模型
func (s *ModelService) CreateModel(ctx context.Context, providerID, name string, isPublic bool) (*model.Model, error) {
	// 验证提供商是否存在
	_, err := s.providerRepo.GetByID(providerID)
	if err != nil {
		return nil, err
	}

	// 创建模型对象
	now := time.Now().UTC().Format(time.RFC3339)
	model := &model.Model{
		ID:         uuid.New().String(),
		ProviderID: providerID,
		Name:       name,
		Status:     model.ModelStatusActive,
		IsPublic:   isPublic,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// 保存到数据库
	err = s.modelRepo.Create(model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// GetModel 获取模型详情
func (s *ModelService) GetModel(ctx context.Context, id string) (*model.Model, error) {
	return s.modelRepo.GetByID(id)
}

// UpdateModel 更新模型
func (s *ModelService) UpdateModel(ctx context.Context, id, name string, status string, isPublic *bool) (*model.Model, error) {
	// 获取现有模型
	model, err := s.modelRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if name != "" {
		model.Name = name
	}
	if status != "" {
		if status != model.ModelStatusActive && status != model.ModelStatusInactive {
			return nil, fmt.Errorf("无效的模型状态: %s", status)
		}
		model.Status = status
	}
	if isPublic != nil {
		model.IsPublic = *isPublic
	}

	// 更新时间
	model.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	// 保存到数据库
	err = s.modelRepo.Update(model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// DeleteModel 删除模型
func (s *ModelService) DeleteModel(ctx context.Context, id string) error {
	// 检查模型是否存在
	_, err := s.modelRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 删除模型
	return s.modelRepo.Delete(id)
}

// ListModels 获取模型列表
func (s *ModelService) ListModels(ctx context.Context, providerID *string, status *string, page, pageSize int) ([]*model.Model, int, error) {
	// 计算偏移量
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	// 获取模型列表
	return s.modelRepo.List(providerID, status, offset, pageSize)
}

// CallModel 调用模型
func (s *ModelService) CallModel(ctx context.Context, modelID string, messages []*schema.Message) (*schema.Message, error) {
	// 获取模型
	model, err := s.modelRepo.GetByID(modelID)
	if err != nil {
		return nil, err
	}

	// 获取 API 密钥
	apiKey, err := s.apiKeyRepo.GetByProviderID(model.ProviderID)
	if err != nil {
		return nil, err
	}

	// 获取提供商
	provider, err := s.registry.GetProvider(model.ProviderID, apiKey.Key)
	if err != nil {
		return nil, err
	}

	// 调用模型
	return provider.Call(ctx, modelID, messages, nil)
}

// StreamModel 流式调用模型
func (s *ModelService) StreamModel(ctx context.Context, modelID string, messages []*schema.Message) (<-chan *schema.Message, error) {
	// 获取模型
	model, err := s.modelRepo.GetByID(modelID)
	if err != nil {
		return nil, err
	}

	// 获取 API 密钥
	apiKey, err := s.apiKeyRepo.GetByProviderID(model.ProviderID)
	if err != nil {
		return nil, err
	}

	// 获取提供商
	provider, err := s.registry.GetProvider(model.ProviderID, apiKey.Key)
	if err != nil {
		return nil, err
	}

	// 流式调用模型
	return provider.Stream(ctx, modelID, messages, nil)
}
