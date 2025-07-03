/**
 * LYSS AI Platform 登录页面
 * 提供用户登录功能，支持邮箱密码登录
 */

import React, { useState } from 'react';
import { Form, Input, Button, Card, Typography, Alert, message } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { history } from '@umijs/max';
import { useAuthStore } from '@/stores/auth';
import { api } from '@/utils/request';
import { API_ENDPOINTS } from '@/config/api';
import { LoginRequest, LoginResponse } from '@/types/api';
import { createStyles } from 'antd-style';

const { Title, Text } = Typography;

const useStyles = createStyles(({ token }) => ({
  container: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    minHeight: '100vh',
    background: `linear-gradient(135deg, ${token.colorPrimary}15 0%, ${token.colorPrimary}05 100%)`,
    padding: '20px',
  },
  loginCard: {
    width: '100%',
    maxWidth: '400px',
    padding: '32px',
    boxShadow: token.boxShadowTertiary,
    borderRadius: token.borderRadiusLG,
  },
  title: {
    textAlign: 'center',
    marginBottom: '8px',
  },
  subtitle: {
    textAlign: 'center',
    color: token.colorTextSecondary,
    marginBottom: '32px',
  },
  formItem: {
    marginBottom: '16px',
  },
  loginButton: {
    width: '100%',
    height: '40px',
    marginTop: '8px',
  },
  registerLink: {
    textAlign: 'center',
    marginTop: '16px',
  },
}));

const Login: React.FC = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string>('');
  const { styles } = useStyles();
  const { setUser, setToken } = useAuthStore();

  const handleSubmit = async (values: LoginRequest) => {
    setLoading(true);
    setError('');

    try {
      // 调用登录API
      const response = await api.post<LoginResponse>(API_ENDPOINTS.AUTH.LOGIN, values);
      
      // 保存token
      setToken(response.access_token);
      
      // 获取用户信息
      const userInfo = await api.get(API_ENDPOINTS.AUTH.CURRENT_USER);
      setUser(userInfo);
      
      message.success('登录成功！');
      
      // 跳转到主页
      const urlParams = new URL(window.location.href).searchParams;
      const redirect = urlParams.get('redirect') || '/chat';
      history.push(redirect);
      
    } catch (error: any) {
      console.error('登录失败:', error);
      if (error.response?.status === 401) {
        setError('邮箱或密码错误，请重试');
      } else {
        setError('登录失败，请检查网络连接或稍后重试');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.container}>
      <Card className={styles.loginCard}>
        <Title level={2} className={styles.title}>
          LYSS AI Platform
        </Title>
        <Text className={styles.subtitle}>
          智能AI聚合平台 - 用户登录
        </Text>

        {error && (
          <Alert
            message={error}
            type="error"
            showIcon
            closable
            onClose={() => setError('')}
            style={{ marginBottom: '16px' }}
          />
        )}

        <Form
          form={form}
          name="login"
          onFinish={handleSubmit}
          autoComplete="off"
          size="large"
        >
          <Form.Item
            className={styles.formItem}
            name="username"
            rules={[
              { required: true, message: '请输入邮箱地址' },
              { type: 'email', message: '请输入有效的邮箱地址' },
            ]}
          >
            <Input
              prefix={<UserOutlined />}
              placeholder="邮箱地址"
              autoComplete="email"
            />
          </Form.Item>

          <Form.Item
            className={styles.formItem}
            name="password"
            rules={[
              { required: true, message: '请输入密码' },
              { min: 6, message: '密码至少6位字符' },
            ]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="密码"
              autoComplete="current-password"
            />
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              loading={loading}
              className={styles.loginButton}
            >
              {loading ? '登录中...' : '登录'}
            </Button>
          </Form.Item>
        </Form>

        <div className={styles.registerLink}>
          <Text>
            还没有账户？{' '}
            <Button
              type="link"
              onClick={() => history.push('/user/register')}
              style={{ padding: 0 }}
            >
              立即注册
            </Button>
          </Text>
        </div>
      </Card>
    </div>
  );
};

export default Login;