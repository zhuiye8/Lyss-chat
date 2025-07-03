/**
 * LYSS AI Platform 用户注册页面
 * 提供用户注册功能，支持邮箱注册
 */

import React, { useState } from 'react';
import { Form, Input, Button, Card, Typography, Alert, message } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined } from '@ant-design/icons';
import { history } from '@umijs/max';
import { api } from '@/utils/request';
import { API_ENDPOINTS } from '@/config/api';
import { RegisterRequest } from '@/types/api';
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
  registerCard: {
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
  registerButton: {
    width: '100%',
    height: '40px',
    marginTop: '8px',
  },
  loginLink: {
    textAlign: 'center',
    marginTop: '16px',
  },
}));

const Register: React.FC = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string>('');
  const [success, setSuccess] = useState<string>('');
  const { styles } = useStyles();

  const handleSubmit = async (values: RegisterRequest) => {
    setLoading(true);
    setError('');
    setSuccess('');

    try {
      // 调用注册API
      await api.post(API_ENDPOINTS.AUTH.REGISTER, values);
      
      setSuccess('注册成功！请前往登录页面登录');
      message.success('注册成功！即将跳转到登录页面...');
      
      // 延迟2秒后跳转到登录页面
      setTimeout(() => {
        history.push('/user/login');
      }, 2000);
      
    } catch (error: any) {
      console.error('注册失败:', error);
      if (error.response?.status === 409) {
        setError('该邮箱已被注册，请使用其他邮箱或前往登录');
      } else if (error.response?.status === 400) {
        setError('注册信息不完整或格式错误，请检查后重试');
      } else {
        setError('注册失败，请检查网络连接或稍后重试');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.container}>
      <Card className={styles.registerCard}>
        <Title level={2} className={styles.title}>
          LYSS AI Platform
        </Title>
        <Text className={styles.subtitle}>
          智能AI聚合平台 - 用户注册
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

        {success && (
          <Alert
            message={success}
            type="success"
            showIcon
            style={{ marginBottom: '16px' }}
          />
        )}

        <Form
          form={form}
          name="register"
          onFinish={handleSubmit}
          autoComplete="off"
          size="large"
        >
          <Form.Item
            className={styles.formItem}
            name="email"
            rules={[
              { required: true, message: '请输入邮箱地址' },
              { type: 'email', message: '请输入有效的邮箱地址' },
            ]}
          >
            <Input
              prefix={<MailOutlined />}
              placeholder="邮箱地址"
              autoComplete="email"
            />
          </Form.Item>

          <Form.Item
            className={styles.formItem}
            name="first_name"
            rules={[
              { required: true, message: '请输入姓名' },
              { min: 2, message: '姓名至少2个字符' },
            ]}
          >
            <Input
              prefix={<UserOutlined />}
              placeholder="姓名"
              autoComplete="given-name"
            />
          </Form.Item>

          <Form.Item
            className={styles.formItem}
            name="last_name"
            rules={[
              { max: 50, message: '姓氏不能超过50个字符' },
            ]}
          >
            <Input
              prefix={<UserOutlined />}
              placeholder="姓氏（可选）"
              autoComplete="family-name"
            />
          </Form.Item>

          <Form.Item
            className={styles.formItem}
            name="password"
            rules={[
              { required: true, message: '请输入密码' },
              { min: 8, message: '密码至少8位字符' },
              { 
                pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]/,
                message: '密码必须包含大小写字母、数字和特殊字符'
              },
            ]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="密码"
              autoComplete="new-password"
            />
          </Form.Item>

          <Form.Item
            className={styles.formItem}
            name="confirmPassword"
            dependencies={['password']}
            rules={[
              { required: true, message: '请确认密码' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('password') === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error('两次输入的密码不一致'));
                },
              }),
            ]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="确认密码"
              autoComplete="new-password"
            />
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              loading={loading}
              disabled={success.length > 0}
              className={styles.registerButton}
            >
              {loading ? '注册中...' : '注册'}
            </Button>
          </Form.Item>
        </Form>

        <div className={styles.loginLink}>
          <Text>
            已有账户？{' '}
            <Button
              type="link"
              onClick={() => history.push('/user/login')}
              style={{ padding: 0 }}
            >
              立即登录
            </Button>
          </Text>
        </div>
      </Card>
    </div>
  );
};

export default Register;