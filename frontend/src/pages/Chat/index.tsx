/**
 * LYSS AI Platform AI对话页面
 * 使用@ant-design/x提供智能对话界面
 */

import React, { useState, useEffect } from 'react';
import { 
  Card, 
  Typography, 
  Select, 
  Button, 
  Space,
  Spin,
  message,
  Row,
  Col,
  Tag,
  Tooltip,
  Alert
} from 'antd';
import { 
  SendOutlined, 
  RobotOutlined, 
  UserOutlined,
  SettingOutlined,
  ClearOutlined,
  HistoryOutlined
} from '@ant-design/icons';
import { Sender, XStream } from '@ant-design/x';
import { useAuthStore } from '@/stores/auth';
import { api } from '@/utils/request';
import { API_ENDPOINTS } from '@/config/api';
import { ChatMessage, ChatRequest, Provider, Model } from '@/types/api';
import { createStyles } from 'antd-style';

const { Title, Text } = Typography;

const useStyles = createStyles(({ token }) => ({
  container: {
    padding: '24px',
    minHeight: 'calc(100vh - 64px)',
    background: token.colorBgContainer,
  },
  chatHeader: {
    marginBottom: '16px',
    padding: '16px',
    background: token.colorBgElevated,
    borderRadius: token.borderRadiusLG,
    border: `1px solid ${token.colorBorder}`,
  },
  chatContainer: {
    height: 'calc(100vh - 280px)',
    display: 'flex',
    flexDirection: 'column',
  },
  messageList: {
    flex: 1,
    overflowY: 'auto',
    padding: '16px',
    background: token.colorBgContainer,
    borderRadius: token.borderRadiusLG,
    border: `1px solid ${token.colorBorder}`,
    marginBottom: '16px',
  },
  messageItem: {
    marginBottom: '16px',
    display: 'flex',
    gap: '12px',
  },
  messageAvatar: {
    width: '32px',
    height: '32px',
    borderRadius: '50%',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    fontSize: '16px',
    flexShrink: 0,
  },
  userAvatar: {
    background: token.colorPrimary,
    color: token.colorWhite,
  },
  assistantAvatar: {
    background: token.colorSuccess,
    color: token.colorWhite,
  },
  messageContent: {
    flex: 1,
    padding: '12px 16px',
    borderRadius: token.borderRadiusLG,
    wordBreak: 'break-word',
    whiteSpace: 'pre-wrap',
  },
  userMessage: {
    background: token.colorPrimaryBg,
    border: `1px solid ${token.colorPrimaryBorder}`,
    marginLeft: 'auto',
    marginRight: '0',
    maxWidth: '80%',
  },
  assistantMessage: {
    background: token.colorBgElevated,
    border: `1px solid ${token.colorBorder}`,
    maxWidth: '80%',
  },
  inputContainer: {
    background: token.colorBgElevated,
    borderRadius: token.borderRadiusLG,
    border: `1px solid ${token.colorBorder}`,
    padding: '12px',
  },
  modelSelector: {
    marginBottom: '12px',
  },
  emptyState: {
    textAlign: 'center',
    padding: '40px',
    color: token.colorTextSecondary,
  },
  loadingMessage: {
    padding: '12px 16px',
    borderRadius: token.borderRadiusLG,
    background: token.colorBgElevated,
    border: `1px solid ${token.colorBorder}`,
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
  },
}));

const Chat: React.FC = () => {
  const { styles } = useStyles();
  const { user } = useAuthStore();
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [loading, setLoading] = useState(false);
  const [providers, setProviders] = useState<Provider[]>([]);
  const [models, setModels] = useState<Model[]>([]);
  const [selectedProvider, setSelectedProvider] = useState<string>('');
  const [selectedModel, setSelectedModel] = useState<string>('');
  const [loadingModels, setLoadingModels] = useState(false);

  // 获取可用的提供商
  useEffect(() => {
    const fetchProviders = async () => {
      try {
        const response = await api.get<Provider[]>(API_ENDPOINTS.CHAT.PROVIDERS);
        setProviders(response);
        if (response.length > 0) {
          setSelectedProvider(response[0].id);
        }
      } catch (error) {
        console.error('获取提供商失败:', error);
        message.error('获取提供商失败');
      }
    };

    fetchProviders();
  }, []);

  // 根据选择的提供商获取模型
  useEffect(() => {
    if (selectedProvider) {
      const fetchModels = async () => {
        setLoadingModels(true);
        try {
          const response = await api.get<Model[]>(
            `${API_ENDPOINTS.CHAT.MODELS}?provider_id=${selectedProvider}`
          );
          setModels(response);
          if (response.length > 0) {
            setSelectedModel(response[0].id);
          }
        } catch (error) {
          console.error('获取模型失败:', error);
          message.error('获取模型失败');
        } finally {
          setLoadingModels(false);
        }
      };

      fetchModels();
    }
  }, [selectedProvider]);

  // 处理发送消息
  const handleSendMessage = async (content: string) => {
    if (!content.trim() || !selectedModel) {
      message.warning('请输入消息内容并选择模型');
      return;
    }

    const userMessage: ChatMessage = {
      id: Date.now().toString(),
      content,
      role: 'user',
      timestamp: new Date().toISOString(),
    };

    setMessages(prev => [...prev, userMessage]);
    setLoading(true);

    try {
      const chatRequest: ChatRequest = {
        model_id: selectedModel,
        message: content,
        conversation_id: 'default', // 简化版本，使用默认会话ID
      };

      // 使用XStream处理流式响应
      const stream = new XStream();
      
      // 添加临时的加载消息
      const loadingMessage: ChatMessage = {
        id: (Date.now() + 1).toString(),
        content: '',
        role: 'assistant',
        timestamp: new Date().toISOString(),
      };
      
      setMessages(prev => [...prev, loadingMessage]);

      // 发送请求并处理流式响应
      const response = await api.post(API_ENDPOINTS.CHAT.SEND, chatRequest);
      
      // 更新助手消息
      const assistantMessage: ChatMessage = {
        id: (Date.now() + 1).toString(),
        content: response.message || '抱歉，我没有收到有效的回复',
        role: 'assistant',
        timestamp: new Date().toISOString(),
      };

      setMessages(prev => prev.slice(0, -1).concat(assistantMessage));
      
    } catch (error: any) {
      console.error('发送消息失败:', error);
      
      // 移除加载消息
      setMessages(prev => prev.slice(0, -1));
      
      // 添加错误消息
      const errorMessage: ChatMessage = {
        id: (Date.now() + 1).toString(),
        content: '抱歉，消息发送失败。请检查网络连接或稍后重试。',
        role: 'assistant',
        timestamp: new Date().toISOString(),
      };
      
      setMessages(prev => [...prev, errorMessage]);
      message.error('消息发送失败');
    } finally {
      setLoading(false);
    }
  };

  // 清空对话历史
  const handleClearMessages = () => {
    setMessages([]);
    message.success('对话历史已清空');
  };

  // 渲染消息列表
  const renderMessages = () => {
    if (messages.length === 0) {
      return (
        <div className={styles.emptyState}>
          <RobotOutlined style={{ fontSize: '48px', marginBottom: '16px' }} />
          <Text>开始与AI对话吧！选择一个模型并发送消息</Text>
        </div>
      );
    }

    return messages.map((message) => (
      <div key={message.id} className={styles.messageItem}>
        <div 
          className={`${styles.messageAvatar} ${
            message.role === 'user' ? styles.userAvatar : styles.assistantAvatar
          }`}
        >
          {message.role === 'user' ? <UserOutlined /> : <RobotOutlined />}
        </div>
        <div 
          className={`${styles.messageContent} ${
            message.role === 'user' ? styles.userMessage : styles.assistantMessage
          }`}
        >
          {message.content}
        </div>
      </div>
    ));
  };

  return (
    <div className={styles.container}>
      <div className={styles.chatHeader}>
        <Row justify="space-between" align="middle">
          <Col>
            <Title level={3} style={{ margin: 0 }}>
              AI智能对话
            </Title>
            <Text type="secondary">
              与AI助手进行自然对话，获取智能回答
            </Text>
          </Col>
          <Col>
            <Space>
              <Tooltip title="对话历史">
                <Button icon={<HistoryOutlined />} type="text" />
              </Tooltip>
              <Tooltip title="清空对话">
                <Button 
                  icon={<ClearOutlined />} 
                  type="text" 
                  onClick={handleClearMessages}
                  disabled={messages.length === 0}
                />
              </Tooltip>
              <Tooltip title="设置">
                <Button icon={<SettingOutlined />} type="text" />
              </Tooltip>
            </Space>
          </Col>
        </Row>
      </div>

      {/* 模型选择器 */}
      <Card className={styles.modelSelector}>
        <Row gutter={16} align="middle">
          <Col span={8}>
            <Text strong>选择提供商:</Text>
            <Select
              value={selectedProvider}
              onChange={setSelectedProvider}
              style={{ width: '100%', marginTop: '8px' }}
              placeholder="选择AI提供商"
            >
              {providers.map(provider => (
                <Select.Option key={provider.id} value={provider.id}>
                  {provider.name}
                </Select.Option>
              ))}
            </Select>
          </Col>
          <Col span={8}>
            <Text strong>选择模型:</Text>
            <Select
              value={selectedModel}
              onChange={setSelectedModel}
              style={{ width: '100%', marginTop: '8px' }}
              placeholder="选择AI模型"
              loading={loadingModels}
              disabled={!selectedProvider}
            >
              {models.map(model => (
                <Select.Option key={model.id} value={model.id}>
                  {model.name}
                </Select.Option>
              ))}
            </Select>
          </Col>
          <Col span={8}>
            <Text strong>状态:</Text>
            <div style={{ marginTop: '8px' }}>
              {selectedModel ? (
                <Tag color="success">已选择模型</Tag>
              ) : (
                <Tag color="default">请选择模型</Tag>
              )}
            </div>
          </Col>
        </Row>
      </Card>

      {/* 对话界面 */}
      <div className={styles.chatContainer}>
        <div className={styles.messageList}>
          {renderMessages()}
          {loading && (
            <div className={styles.loadingMessage}>
              <Spin size="small" />
              <Text>AI正在思考中...</Text>
            </div>
          )}
        </div>

        {/* 输入框 */}
        <div className={styles.inputContainer}>
          <Sender
            placeholder="输入消息..."
            loading={loading}
            onSubmit={handleSendMessage}
            disabled={!selectedModel}
            suffix={
              <Button
                type="primary"
                icon={<SendOutlined />}
                disabled={!selectedModel || loading}
              >
                发送
              </Button>
            }
          />
          {!selectedModel && (
            <Alert
              message="请先选择AI提供商和模型"
              type="warning"
              showIcon
              style={{ marginTop: '8px' }}
            />
          )}
        </div>
      </div>
    </div>
  );
};

export default Chat;