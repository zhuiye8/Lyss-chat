# LYSS AI 平台 - 前端组件化架构 (V2.1 优化版)

**版本**: 2.1
**最后更新**: 2025年7月2日

---

## 1. 概述

V2.1 版本的前端架构核心是拥抱 `@ant-design/x` 这一官方推荐的、更灵活的 AI 对话组件库。我们不再使用已废弃的 `@ant-design/pro-chat`。

## 2. 核心业务组件

### 2.1. 对话页面 (`/chat`) 布局 (无变化)

```mermaid
graph TD
    subgraph ChatPage (src/pages/Chat)
        direction LR
        Sidebar[SideBar: 对话历史] --> MainContent
        subgraph MainContent
            direction TB
            ModelSelector[ModelSelector] --> ChatArea[ChatArea]
        end
    end
```

### 2.2. 模型选择器 (`ModelSelector.tsx`) (V2 API 对齐)

*   **核心职责**:
    1.  调用 **`GET /api/v1/chat/available-models`** API 获取当前用户可用的模型列表。
    2.  当用户切换模型时，通过 Zustand 更新全局选中的 `currentModelId`。
*   **数据获取**: 使用 React Query 的 `useQuery`。
    *   `queryKey`: `['chat', 'availableModels']`
*   **代码骨架**:
    ```tsx
    // src/components/Business/ModelSelector/index.tsx
    import { Select } from 'antd';
    import { useQuery } from '@tanstack/react-query';
    import { useChatStore } from '@/stores/chatStore';
    import { getAvailableModels } from '@/services/lyss/chatAPI'; // 对齐 V2 API

    const ModelSelector = () => {
      const { data: models, isLoading } = useQuery({
        queryKey: ['chat', 'availableModels'],
        queryFn: getAvailableModels,
      });
      const { currentModelId, setCurrentModelId } = useChatStore();

      return (
        <Select
          loading={isLoading}
          value={currentModelId}
          onChange={(value) => setCurrentModelId(value)}
          style={{ width: 280 }}
          placeholder="请选择模型"
        >
          {models?.map(model => (
            <Select.Option key={model.id} value={model.id}>
              {/* 可视化展示，包括模型名称和供应商信息 */}
              <span>{model.display_name || model.model_name}</span>
              <small style={{ marginLeft: 8, color: '#999' }}>
                by {model.provider.name} ({model.provider.scope})
              </small>
            </Select.Option>
          ))}
        </Select>
      );
    };
    ```

### 2.3. 对话区域 (`ChatArea.tsx`) (V2.1 核心重构)

*   **核心库**: `@ant-design/x`
*   **设计理念**: `@ant-design/x` 不是一个单一组件，而是一个 **Hooks + 原子组件** 的 SDK。我们必须使用其提供的 `useXChat` hook 来管理状态，并组合 `Chat`、`Bubble`、`Composer` 等原子组件来构建UI。这提供了无与伦比的灵活性。
*   **代码骨架 (正确示例)**:
    ```tsx
    // src/components/Business/ChatArea/index.tsx
    import { Chat, Bubble, Composer, useXChat } from '@ant-design/x';
    import { useChatStore } from '@/stores/chatStore';
    import { chatCompletion } from '@/services/lyss/chatAPI'; // 流式API服务
    import { Spin } from 'antd';

    // 定义消息的数据结构
    interface IChatMessage {
      content: string;
      role: 'user' | 'assistant' | 'system';
      id: string;
    }

    const ChatArea = () => {
      const { currentModelId } = useChatStore();

      // 使用 useXChat hook 管理聊天状态和逻辑
      const { messages, sendMessage, loading } = useXChat<IChatMessage>({
        // onSendMessage 是核心，它定义了消息发送时的行为
        onSendMessage: async (message) => {
          if (!currentModelId) {
            // 实际应有更友好的提示
            alert("Please select a model first!");
            return;
          }
          
          // 调用我们的流式API
          const stream = await chatCompletion({
            model_id: currentModelId,
            messages: [{ role: message.role, content: message.content }],
          });

          // 将流式响应传递给 x-chat 的流处理器
          return stream;
        },
      });

      // 渲染UI
      return (
        <Chat
          messages={messages}
          renderMessage={(msg) => <Bubble message={msg} />}
          composer={
            <Composer
              placeholder="请输入您的问题..."
              onSend={(value) => sendMessage({ role: 'user', content: value })}
            />
          }
          loading={loading && <Spin />}
        />
      );
    };
    ```

### 2.4. 供应商配置表单 (`ProviderForm.tsx`) (V2.1 优化)

*   **核心职责**: 动态生成表单以适应不同供应商的配置需求。
*   **实现**:
    1.  创建一个 `ProviderSchemaRegistry` 对象，用于存储每种 `provider_type` 对应的表单 schema (包含字段名、标签、校验规则等)。
    2.  在表单组件中，当用户选择一个 `provider_type` 时，从注册表中获取其 schema，然后动态渲染出对应的 `Form.Item` 列表。
    3.  提交时，将动态表单的值统一收集到 `config` 对象中，发送给后端。
*   **优势**: 此方案将供应商的“UI表现”也插件化了，新增供应商只需在前端注册其表单 Schema，无需修改核心表单组件代码。
