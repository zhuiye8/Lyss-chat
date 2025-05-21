import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button, Input, Modal, Form, message, Empty } from 'antd'
import { PlusOutlined, SearchOutlined } from '@ant-design/icons'
import { Conversations } from '@ant-design/x'
import { chatService } from '../../services/chatService'
import { useChatStore, Canvas } from '../../store/chatStore'
import styles from './ChatListPage.module.css'

// 临时使用固定的工作区ID，实际应该从状态或路由中获取
const WORKSPACE_ID = '123e4567-e89b-12d3-a456-426614174000'

const ChatListPage = () => {
  const [loading, setLoading] = useState(false)
  const [searchValue, setSearchValue] = useState('')
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [form] = Form.useForm()
  const navigate = useNavigate()
  const { canvases, setCanvases } = useChatStore()

  // 加载聊天列表
  const loadCanvases = async () => {
    try {
      setLoading(true)
      const response = await chatService.getCanvases(WORKSPACE_ID, {
        type: 'chat'
      })
      setCanvases(response.items)
    } catch (error) {
      message.error('加载聊天列表失败')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadCanvases()
  }, [])

  // 创建新聊天
  const handleCreateChat = async (values: { title: string }) => {
    try {
      const newCanvas = await chatService.createCanvas({
        title: values.title,
        workspace_id: WORKSPACE_ID,
        type: 'chat'
      })
      
      setCanvases([newCanvas, ...canvases])
      setIsModalVisible(false)
      form.resetFields()
      
      // 导航到新创建的聊天
      navigate(`/chat/${newCanvas.id}`)
    } catch (error) {
      message.error('创建聊天失败')
      console.error(error)
    }
  }

  // 过滤聊天列表
  const filteredCanvases = canvases.filter(canvas => 
    canvas.title.toLowerCase().includes(searchValue.toLowerCase())
  )

  // 转换为 Conversations 组件需要的格式
  const conversationItems = filteredCanvases.map(canvas => ({
    key: canvas.id,
    label: canvas.title,
    timestamp: new Date(canvas.created_at).getTime(),
    // 可以添加更多属性，如最后一条消息等
  }))

  // 处理聊天选择
  const handleConversationSelect = (key: string) => {
    navigate(`/chat/${key}`)
  }

  // 处理聊天操作菜单
  const handleMenuClick = async ({ key, domEvent }: { key: string, domEvent: React.MouseEvent }) => {
    domEvent.stopPropagation()
    
    const [action, canvasId] = key.split(':')
    const canvas = canvases.find(c => c.id === canvasId)
    
    if (!canvas) return
    
    switch (action) {
      case 'rename':
        // 实现重命名功能
        break
      case 'delete':
        try {
          await chatService.deleteCanvas(canvasId)
          setCanvases(canvases.filter(c => c.id !== canvasId))
          message.success('删除成功')
        } catch (error) {
          message.error('删除失败')
          console.error(error)
        }
        break
    }
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h2>聊天列表</h2>
        <Button 
          type="primary" 
          icon={<PlusOutlined />}
          onClick={() => setIsModalVisible(true)}
        >
          新建聊天
        </Button>
      </div>
      
      <div className={styles.search}>
        <Input
          placeholder="搜索聊天"
          prefix={<SearchOutlined />}
          value={searchValue}
          onChange={e => setSearchValue(e.target.value)}
        />
      </div>
      
      <div className={styles.list}>
        {conversationItems.length > 0 ? (
          <Conversations
            items={conversationItems}
            onActiveChange={handleConversationSelect}
            menu={{
              items: [
                { key: 'rename', label: '重命名' },
                { key: 'delete', label: '删除' }
              ],
              onClick: ({ key }, conversation) => {
                handleMenuClick({ 
                  key: `${key}:${conversation.key}`, 
                  domEvent: {} as React.MouseEvent 
                })
              }
            }}
          />
        ) : (
          <Empty 
            description={loading ? '加载中...' : '暂无聊天'} 
            image={Empty.PRESENTED_IMAGE_SIMPLE} 
          />
        )}
      </div>
      
      <Modal
        title="新建聊天"
        open={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        footer={null}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleCreateChat}
        >
          <Form.Item
            name="title"
            label="聊天标题"
            rules={[{ required: true, message: '请输入聊天标题' }]}
          >
            <Input placeholder="请输入聊天标题" />
          </Form.Item>
          
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              创建
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default ChatListPage
