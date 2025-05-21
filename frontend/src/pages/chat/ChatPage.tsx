import { useState, useEffect, useRef } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Spin, message, Select, Button } from 'antd'
import { ArrowLeftOutlined, SettingOutlined } from '@ant-design/icons'
import { Bubble, Sender, useXChat, XProvider, XStream, ThoughtChain } from '@ant-design/x'
import { chatService } from '../../services/chatService'
import { modelService, Model } from '../../services/modelService'
import { useChatStore, Message } from '../../store/chatStore'
import { useAuthStore } from '../../store/authStore'
import ThoughtChainDisplay, { ThoughtItem } from '../../components/chat/ThoughtChainDisplay'
import ModelParamsConfig, { ModelParams } from '../../components/chat/ModelParamsConfig'
import styles from './ChatPage.module.css'

const ChatPage = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { user } = useAuthStore()
  const {
    currentCanvas,
    messages,
    setCurrentCanvas,
    setMessages,
    addMessage,
    updateMessage,
    setLoading,
    loading
  } = useChatStore()

  const [models, setModels] = useState<Model[]>([])
  const [selectedModel, setSelectedModel] = useState<string>('')
  const [sendLoading, setSendLoading] = useState(false)
  const [thoughts, setThoughts] = useState<ThoughtItem[]>([])
  const [showModelParams, setShowModelParams] = useState(false)
  const [modelParams, setModelParams] = useState<ModelParams>({
    temperature: 0.7,
    top_p: 1,
    max_tokens: 2000,
    presence_penalty: 0,
    frequency_penalty: 0,
    stream: true
  })
  const messagesEndRef = useRef<HTMLDivElement>(null)

  // 加载画布信息
  const loadCanvas = async () => {
    if (!id) return

    try {
      setLoading(true)
      const canvas = await chatService.getCanvas(id)
      setCurrentCanvas(canvas)

      // 如果画布有关联的模型，设置为选中的模型
      if (canvas.model_id) {
        setSelectedModel(canvas.model_id)
      }
    } catch (error) {
      message.error('加载聊天信息失败')
      console.error(error)
      navigate('/')
    } finally {
      setLoading(false)
    }
  }

  // 加载消息历史
  const loadMessages = async () => {
    if (!id) return

    try {
      setLoading(true)
      const response = await chatService.getMessages(id)
      setMessages(response.items)
    } catch (error) {
      message.error('加载消息历史失败')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  // 加载可用模型
  const loadModels = async () => {
    try {
      const response = await modelService.getModels({
        status: 'active'
      })
      setModels(response.items)

      // 如果没有选中的模型且有可用模型，选择第一个
      if (!selectedModel && response.items.length > 0) {
        setSelectedModel(response.items[0].id)
      }
    } catch (error) {
      message.error('加载模型列表失败')
      console.error(error)
    }
  }

  useEffect(() => {
    loadCanvas()
    loadMessages()
    loadModels()

    return () => {
      // 清理工作
      setCurrentCanvas(null)
      setMessages([])
    }
  }, [id])

  // 滚动到底部
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  // 发送消息
  const handleSendMessage = async (content: string) => {
    if (!id || !content.trim() || !user) return

    try {
      setSendLoading(true)

      // 添加用户消息到列表
      const userMessage: Message = {
        id: `temp-${Date.now()}`,
        canvas_id: id,
        role: 'user',
        content,
        created_at: new Date().toISOString()
      }
      addMessage(userMessage)

      // 添加临时的 AI 消息
      const aiMessage: Message = {
        id: `temp-ai-${Date.now()}`,
        canvas_id: id,
        role: 'assistant',
        content: '思考中...',
        created_at: new Date().toISOString()
      }
      addMessage(aiMessage)

      // 使用 XStream 组件处理流式响应
      const streamUrl = `/api/v1/canvases/${id}/messages/stream`;
      const requestData = {
        content,
        model_id: selectedModel,
        params: {
          ...modelParams,
          // 确保流式输出
          stream: true
        }
      };

      // 创建 XStream 实例
      const stream = new XStream({
        url: streamUrl,
        method: 'POST',
        data: requestData,
        onChunk: (chunk) => {
          // 更新 AI 消息内容
          updateMessage(aiMessage.id, chunk.content || '');

          // 处理思维链数据
          if (chunk.thoughts) {
            try {
              const thoughtData = JSON.parse(chunk.thoughts);
              if (Array.isArray(thoughtData)) {
                setThoughts(thoughtData.map(item => ({
                  id: item.id || `thought-${Date.now()}-${Math.random()}`,
                  title: item.title || '思考步骤',
                  content: item.content || '',
                  status: item.status || 'success',
                  children: item.children || []
                })));
              }
            } catch (error) {
              console.error('解析思维链数据失败', error);
            }
          }
        },
        onComplete: () => {
          setSendLoading(false);
        },
        onError: (error) => {
          message.error('获取 AI 响应失败');
          console.error(error);
          setSendLoading(false);
        }
      });

      // 启动流式请求
      stream.start();

    } catch (error) {
      message.error('发送消息失败');
      console.error(error);
      setSendLoading(false);
    }
  }

  // 更新选中的模型
  const handleModelChange = async (modelId: string) => {
    if (!id || !currentCanvas) return

    setSelectedModel(modelId)

    try {
      await chatService.updateCanvas(id, { model_id: modelId })
      message.success('模型已更新')
    } catch (error) {
      message.error('更新模型失败')
      console.error(error)
    }
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Button
          icon={<ArrowLeftOutlined />}
          type="text"
          onClick={() => navigate('/')}
        />
        <h2>{currentCanvas?.title || '加载中...'}</h2>
        <div className={styles.headerRight}>
          <Select
            placeholder="选择模型"
            value={selectedModel}
            onChange={handleModelChange}
            style={{ width: 200 }}
            options={models.map(model => ({
              label: model.name,
              value: model.id
            }))}
          />
          <Button
            icon={<SettingOutlined />}
            type="text"
            onClick={() => setShowModelParams(!showModelParams)}
            style={{ marginLeft: 8 }}
          />
        </div>
      </div>

      {showModelParams && (
        <ModelParamsConfig
          initialParams={modelParams}
          onParamsChange={setModelParams}
        />
      )}

      <div className={styles.content}>
        {loading ? (
          <div className={styles.loading}>
            <Spin size="large" />
          </div>
        ) : (
          <div className={styles.messages}>
            {messages.map(message => (
              <Bubble
                key={message.id}
                content={message.content}
                placement={message.role === 'user' ? 'end' : 'start'}
                variant={message.role === 'user' ? 'filled' : 'outlined'}
                typing={message.role === 'assistant' && message.content === '思考中...'}
              />
            ))}
            <div ref={messagesEndRef} />
          </div>
        )}
      </div>

      {thoughts.length > 0 && (
        <ThoughtChainDisplay thoughts={thoughts} loading={sendLoading} />
      )}

      <div className={styles.footer}>
        <Sender
          onSubmit={handleSendMessage}
          disabled={sendLoading}
          loading={sendLoading}
          placeholder="输入消息..."
        />
      </div>
    </div>
  )
}

export default ChatPage
