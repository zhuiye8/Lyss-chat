import { useState } from 'react'
import { ThoughtChain } from '@ant-design/x'
import { Button, Collapse } from 'antd'
import styles from './ThoughtChainDisplay.module.css'

export interface ThoughtItem {
  id: string
  title: string
  content: string
  status: 'pending' | 'success' | 'error'
  children?: ThoughtItem[]
}

interface ThoughtChainDisplayProps {
  thoughts: ThoughtItem[]
  loading?: boolean
}

const ThoughtChainDisplay = ({ thoughts, loading = false }: ThoughtChainDisplayProps) => {
  const [isExpanded, setIsExpanded] = useState(false)

  if (!thoughts || thoughts.length === 0) {
    return null
  }

  const renderThoughtItems = (items: ThoughtItem[]) => {
    return items.map(item => (
      <ThoughtChain
        key={item.id}
        title={item.title}
        content={item.content}
        status={item.status}
        loading={loading && item.status === 'pending'}
      >
        {item.children && item.children.length > 0 && renderThoughtItems(item.children)}
      </ThoughtChain>
    ))
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h3>AI 思维过程</h3>
        <Button 
          type="text" 
          onClick={() => setIsExpanded(!isExpanded)}
        >
          {isExpanded ? '收起' : '展开'}
        </Button>
      </div>
      
      <Collapse activeKey={isExpanded ? ['1'] : []}>
        <Collapse.Panel key="1" header="查看详细思考过程" className={styles.panel}>
          <div className={styles.thoughtChain}>
            {renderThoughtItems(thoughts)}
          </div>
        </Collapse.Panel>
      </Collapse>
    </div>
  )
}

export default ThoughtChainDisplay
