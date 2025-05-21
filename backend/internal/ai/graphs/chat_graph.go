package graphs

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/your-org/lyss-chat-backend/internal/ai/components"
	"github.com/your-org/lyss-chat-backend/pkg/logger"
)

// 聊天模板
const chatTemplate = `你是一个有用的AI助手。请根据用户的问题提供准确、有帮助的回答。
如果你不知道答案，请诚实地说你不知道，不要编造信息。
`

// ChatGraph 表示聊天图形
type ChatGraph struct {
	graph *compose.CompiledGraph[map[string]any, *schema.Message]
}

// NewChatGraph 创建一个新的聊天图形
func NewChatGraph(ctx context.Context, model components.ChatModelAdapter, logger *logger.Logger) (*ChatGraph, error) {
	// 创建图形
	graph := compose.NewGraph[map[string]any, *schema.Message]()

	// 添加聊天模板节点
	err := graph.AddChatTemplateNode("node_template", chatTemplate)
	if err != nil {
		return nil, err
	}

	// 添加聊天模型节点
	err = graph.AddChatModelNode("node_model", &model)
	if err != nil {
		return nil, err
	}

	// 添加边
	err = graph.AddEdge(compose.START, "node_template")
	if err != nil {
		return nil, err
	}

	err = graph.AddEdge("node_template", "node_model")
	if err != nil {
		return nil, err
	}

	// 编译图形
	compiled, err := graph.Compile(ctx)
	if err != nil {
		return nil, err
	}

	return &ChatGraph{
		graph: compiled,
	}, nil
}

// Invoke 调用聊天图形
func (g *ChatGraph) Invoke(ctx context.Context, input map[string]any) (*schema.Message, error) {
	return g.graph.Invoke(ctx, input)
}

// Stream 流式调用聊天图形
func (g *ChatGraph) Stream(ctx context.Context, input map[string]any) (<-chan *schema.Message, error) {
	return g.graph.Stream(ctx, input)
}
