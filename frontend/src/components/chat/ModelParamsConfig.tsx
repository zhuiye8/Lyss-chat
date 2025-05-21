import { useState } from 'react'
import { Form, InputNumber, Slider, Switch, Button, Collapse, Tooltip } from 'antd'
import { QuestionCircleOutlined } from '@ant-design/icons'
import styles from './ModelParamsConfig.module.css'

export interface ModelParams {
  temperature: number
  top_p: number
  max_tokens: number
  presence_penalty: number
  frequency_penalty: number
  stream: boolean
}

interface ModelParamsConfigProps {
  initialParams: ModelParams
  onParamsChange: (params: ModelParams) => void
}

const defaultParams: ModelParams = {
  temperature: 0.7,
  top_p: 1,
  max_tokens: 2000,
  presence_penalty: 0,
  frequency_penalty: 0,
  stream: true
}

const ModelParamsConfig = ({ 
  initialParams = defaultParams, 
  onParamsChange 
}: ModelParamsConfigProps) => {
  const [form] = Form.useForm()
  const [isExpanded, setIsExpanded] = useState(false)

  const handleValuesChange = (_: any, allValues: ModelParams) => {
    onParamsChange(allValues)
  }

  return (
    <div className={styles.container}>
      <Collapse 
        activeKey={isExpanded ? ['1'] : []}
        onChange={() => setIsExpanded(!isExpanded)}
      >
        <Collapse.Panel 
          key="1" 
          header="模型参数配置" 
          className={styles.panel}
        >
          <Form
            form={form}
            layout="vertical"
            initialValues={initialParams}
            onValuesChange={handleValuesChange}
          >
            <Form.Item
              label={
                <span>
                  温度 
                  <Tooltip title="控制生成文本的随机性。较高的值会使输出更加随机，较低的值会使输出更加确定。">
                    <QuestionCircleOutlined className={styles.infoIcon} />
                  </Tooltip>
                </span>
              }
              name="temperature"
            >
              <div className={styles.sliderWithInput}>
                <Slider 
                  min={0} 
                  max={2} 
                  step={0.1} 
                  style={{ flex: 1 }} 
                />
                <InputNumber 
                  min={0} 
                  max={2} 
                  step={0.1} 
                  style={{ width: 60, marginLeft: 16 }} 
                />
              </div>
            </Form.Item>

            <Form.Item
              label={
                <span>
                  Top P 
                  <Tooltip title="控制生成文本的多样性。较低的值会使输出更加确定，较高的值会使输出更加多样。">
                    <QuestionCircleOutlined className={styles.infoIcon} />
                  </Tooltip>
                </span>
              }
              name="top_p"
            >
              <div className={styles.sliderWithInput}>
                <Slider 
                  min={0} 
                  max={1} 
                  step={0.05} 
                  style={{ flex: 1 }} 
                />
                <InputNumber 
                  min={0} 
                  max={1} 
                  step={0.05} 
                  style={{ width: 60, marginLeft: 16 }} 
                />
              </div>
            </Form.Item>

            <Form.Item
              label={
                <span>
                  最大令牌数 
                  <Tooltip title="生成文本的最大长度。">
                    <QuestionCircleOutlined className={styles.infoIcon} />
                  </Tooltip>
                </span>
              }
              name="max_tokens"
            >
              <div className={styles.sliderWithInput}>
                <Slider 
                  min={100} 
                  max={4000} 
                  step={100} 
                  style={{ flex: 1 }} 
                />
                <InputNumber 
                  min={100} 
                  max={4000} 
                  step={100} 
                  style={{ width: 80, marginLeft: 16 }} 
                />
              </div>
            </Form.Item>

            <Form.Item
              label={
                <span>
                  存在惩罚 
                  <Tooltip title="控制模型重复已经生成的内容的倾向。较高的值会减少重复。">
                    <QuestionCircleOutlined className={styles.infoIcon} />
                  </Tooltip>
                </span>
              }
              name="presence_penalty"
            >
              <div className={styles.sliderWithInput}>
                <Slider 
                  min={-2} 
                  max={2} 
                  step={0.1} 
                  style={{ flex: 1 }} 
                />
                <InputNumber 
                  min={-2} 
                  max={2} 
                  step={0.1} 
                  style={{ width: 60, marginLeft: 16 }} 
                />
              </div>
            </Form.Item>

            <Form.Item
              label={
                <span>
                  频率惩罚 
                  <Tooltip title="控制模型重复同一词语的倾向。较高的值会减少重复。">
                    <QuestionCircleOutlined className={styles.infoIcon} />
                  </Tooltip>
                </span>
              }
              name="frequency_penalty"
            >
              <div className={styles.sliderWithInput}>
                <Slider 
                  min={-2} 
                  max={2} 
                  step={0.1} 
                  style={{ flex: 1 }} 
                />
                <InputNumber 
                  min={-2} 
                  max={2} 
                  step={0.1} 
                  style={{ width: 60, marginLeft: 16 }} 
                />
              </div>
            </Form.Item>

            <Form.Item
              label={
                <span>
                  流式输出 
                  <Tooltip title="启用流式输出可以实时显示生成的文本。">
                    <QuestionCircleOutlined className={styles.infoIcon} />
                  </Tooltip>
                </span>
              }
              name="stream"
              valuePropName="checked"
            >
              <Switch />
            </Form.Item>

            <Form.Item>
              <Button 
                type="default" 
                onClick={() => form.resetFields()}
              >
                重置为默认值
              </Button>
            </Form.Item>
          </Form>
        </Collapse.Panel>
      </Collapse>
    </div>
  )
}

export default ModelParamsConfig
