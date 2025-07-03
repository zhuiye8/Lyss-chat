/**
 * 用户管理页面 (管理员)
 * 用户列表、角色管理、权限控制
 */

import React from 'react';
import { Card, Typography } from 'antd';

const { Title } = Typography;

const Users: React.FC = () => {
  return (
    <div style={{ padding: '24px' }}>
      <Card>
        <Title level={2}>用户管理</Title>
        <p>用户管理界面开发中...</p>
      </Card>
    </div>
  );
};

export default Users;