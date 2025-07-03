/**
 * 个人设置页面
 * 用户个人信息管理、偏好设置
 */

import React from 'react';
import { Card, Typography } from 'antd';

const { Title } = Typography;

const Settings: React.FC = () => {
  return (
    <div style={{ padding: '24px' }}>
      <Card>
        <Title level={2}>个人设置</Title>
        <p>个人设置界面开发中...</p>
      </Card>
    </div>
  );
};

export default Settings;