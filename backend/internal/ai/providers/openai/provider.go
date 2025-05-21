package openai

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/cloudwego/eino/schema"
	"github.com/sashabaranov/go-openai"
	"github.com/your-org/lyss-chat-backend/internal/ai/providers"
	"github.com/your-org/lyss-chat-backend/internal/domain/model"
	"github.com/your-org/lyss-chat-backend/pkg/logger"
)

const (
	ProviderID = "openai"
)

// 支持的模型列表
var supportedModels = []*model.Model{
	{
		ID:         "gpt-3.5-turbo",
		ProviderID: ProviderID,
		Name:       "GPT-3.5 Turbo",
		Status:     model.ModelStatusActive,
		IsPublic:   true,
	},
	{
		ID:         "gpt-4",
		ProviderID: ProviderID,
		Name:       "GPT-4",
		Status:     model.ModelStatusActive,
		IsPublic:   true,
	},
	{
		ID:         "gpt-4-turbo",
		ProviderID: ProviderID,
		Name:       "GPT-4 Turbo",
		Status:     model.ModelStatusActive,
		IsPublic:   true,
	},
}

// Provider 实现 OpenAI 提供商
type Provider struct {
	client *openai.Client
	logger *logger.Logger
}

// Factory 实现 OpenAI 提供商工厂
type Factory struct{}

// Create 创建 OpenAI 提供商实例
func (f *Factory) Create(apiKey string, logger *logger.Logger) (providers.Provider, error) {
	if apiKey == "" {
		return nil, providers.ErrInvalidAPIKey
	}

	client := openai.NewClient(apiKey)
	return &Provider{
		client: client,
		logger: logger,
	}, nil
}

// GetName 获取提供商名称
func (p *Provider) GetName() string {
	return "OpenAI"
}

// GetModels 获取提供商支持的模型列表
func (p *Provider) GetModels() []*model.Model {
	return supportedModels
}

// Call 调用模型
func (p *Provider) Call(ctx context.Context, modelID string, messages []*schema.Message, params map[string]interface{}) (*schema.Message, error) {
	// 验证模型
	if !p.isModelSupported(modelID) {
		return nil, providers.ErrModelNotFound
	}

	// 转换消息格式
	openaiMessages := convertToOpenAIMessages(messages)

	// 设置请求参数
	req := openai.ChatCompletionRequest{
		Model:    modelID,
		Messages: openaiMessages,
	}

	// 应用可选参数
	if params != nil {
		if temp, ok := params["temperature"].(float64); ok {
			req.Temperature = float32(temp)
		}
		if topP, ok := params["top_p"].(float64); ok {
			req.TopP = float32(topP)
		}
		if maxTokens, ok := params["max_tokens"].(float64); ok {
			req.MaxTokens = int(maxTokens)
		}
		if presencePenalty, ok := params["presence_penalty"].(float64); ok {
			req.PresencePenalty = float32(presencePenalty)
		}
		if frequencyPenalty, ok := params["frequency_penalty"].(float64); ok {
			req.FrequencyPenalty = float32(frequencyPenalty)
		}
	}

	// 调用 API
	resp, err := p.client.CreateChatCompletion(ctx, req)
	if err != nil {
		p.logger.Error("OpenAI API 调用失败", err)
		return nil, fmt.Errorf("%w: %v", providers.ErrAPICallFailed, err)
	}

	// 检查响应
	if len(resp.Choices) == 0 {
		return nil, errors.New("OpenAI 返回了空响应")
	}

	// 转换响应
	return &schema.Message{
		Role:    resp.Choices[0].Message.Role,
		Content: resp.Choices[0].Message.Content,
	}, nil
}

// Stream 流式调用模型
func (p *Provider) Stream(ctx context.Context, modelID string, messages []*schema.Message, params map[string]interface{}) (<-chan *schema.Message, error) {
	// 验证模型
	if !p.isModelSupported(modelID) {
		return nil, providers.ErrModelNotFound
	}

	// 转换消息格式
	openaiMessages := convertToOpenAIMessages(messages)

	// 设置请求参数
	req := openai.ChatCompletionRequest{
		Model:    modelID,
		Messages: openaiMessages,
		Stream:   true,
	}

	// 应用可选参数
	if params != nil {
		if temp, ok := params["temperature"].(float64); ok {
			req.Temperature = float32(temp)
		}
		if topP, ok := params["top_p"].(float64); ok {
			req.TopP = float32(topP)
		}
		if maxTokens, ok := params["max_tokens"].(float64); ok {
			req.MaxTokens = int(maxTokens)
		}
		if presencePenalty, ok := params["presence_penalty"].(float64); ok {
			req.PresencePenalty = float32(presencePenalty)
		}
		if frequencyPenalty, ok := params["frequency_penalty"].(float64); ok {
			req.FrequencyPenalty = float32(frequencyPenalty)
		}
	}

	// 创建流式响应通道
	stream, err := p.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		p.logger.Error("OpenAI 流式 API 调用失败", err)
		return nil, fmt.Errorf("%w: %v", providers.ErrAPICallFailed, err)
	}

	// 创建输出通道
	outputChan := make(chan *schema.Message)

	// 启动 goroutine 处理流式响应
	go func() {
		defer close(outputChan)
		defer stream.Close()

		var fullContent strings.Builder

		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				// 流结束
				break
			}

			if err != nil {
				p.logger.Error("接收流式响应失败", err)
				return
			}

			// 检查是否有内容
			if len(response.Choices) == 0 || response.Choices[0].Delta.Content == "" {
				continue
			}

			// 累积内容
			fullContent.WriteString(response.Choices[0].Delta.Content)

			// 发送消息
			outputChan <- &schema.Message{
				Role:    "assistant",
				Content: fullContent.String(),
			}
		}
	}()

	return outputChan, nil
}

// isModelSupported 检查模型是否受支持
func (p *Provider) isModelSupported(modelID string) bool {
	for _, m := range supportedModels {
		if m.ID == modelID {
			return true
		}
	}
	return false
}

// convertToOpenAIMessages 将 schema.Message 转换为 OpenAI 消息
func convertToOpenAIMessages(messages []*schema.Message) []openai.ChatCompletionMessage {
	openaiMessages := make([]openai.ChatCompletionMessage, 0, len(messages))
	for _, msg := range messages {
		openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	return openaiMessages
}
