/**
 * LYSS AI Platform 应用入口配置
 * 包含布局配置、初始状态管理、权限控制等
 */

import React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { history, Link } from '@umijs/max';
import type { Settings as LayoutSettings } from '@ant-design/pro-components';
import type { RunTimeLayoutConfig } from '@umijs/max';
import { LinkOutlined } from '@ant-design/icons';

import { Footer, AvatarDropdown } from '@/components';
import { queryClient } from '@/config/reactQuery';
import { useAuthStore } from '@/stores/auth';
import { api } from '@/utils/request';
import { API_ENDPOINTS } from '@/config/api';
import { User } from '@/types/api';
import defaultSettings from '../config/defaultSettings';

const isDev = process.env.NODE_ENV === 'development';
const loginPath = '/user/login';

/**
 * 获取当前用户信息
 */
const fetchCurrentUser = async (): Promise<User | undefined> => {
  try {
    const response = await api.get<User>(API_ENDPOINTS.AUTH.CURRENT_USER);
    return response;
  } catch (error) {
    console.log('获取用户信息失败:', error);
    return undefined;
  }
};

/**
 * 应用初始状态
 */
export async function getInitialState(): Promise<{
  settings?: Partial<LayoutSettings>;
  currentUser?: User;
  loading?: boolean;
  fetchUserInfo?: () => Promise<User | undefined>;
}> {
  const fetchUserInfo = async () => {
    try {
      const user = await fetchCurrentUser();
      if (user) {
        // 更新认证状态
        useAuthStore.getState().setUser(user);
        return user;
      }
    } catch (error) {
      // 获取用户信息失败，跳转到登录页
      useAuthStore.getState().logout();
      history.push(loginPath);
    }
    return undefined;
  };

  // 如果不是登录页面，尝试获取用户信息
  const { location } = history;
  if (location.pathname !== loginPath && location.pathname !== '/user/register') {
    const currentUser = await fetchUserInfo();
    return {
      fetchUserInfo,
      currentUser,
      settings: defaultSettings as Partial<LayoutSettings>,
    };
  }

  return {
    fetchUserInfo,
    settings: defaultSettings as Partial<LayoutSettings>,
  };
}

/**
 * ProLayout 布局配置
 */
export const layout: RunTimeLayoutConfig = ({ initialState, setInitialState }) => {
  return {
    // 顶部操作区域
    actionsRender: () => [
      // 开发环境显示API文档链接
      ...(isDev ? [
        <Link key="api-docs" to="http://localhost:8000/api/v1/docs" target="_blank">
          <LinkOutlined />
          <span>API文档</span>
        </Link>
      ] : []),
    ],
    
    // 用户头像配置
    avatarProps: {
      src: initialState?.currentUser?.email ? 
        `https://api.dicebear.com/7.x/initials/svg?seed=${initialState.currentUser.email}` : 
        undefined,
      title: initialState?.currentUser?.email || '用户',
      render: (_, avatarChildren) => {
        return <AvatarDropdown>{avatarChildren}</AvatarDropdown>;
      },
    },
    
    // 页脚
    footerRender: () => <Footer />,
    
    // 路由变化时的权限检查
    onPageChange: () => {
      const { location } = history;
      const isAuthPage = ['/user/login', '/user/register'].includes(location.pathname);
      
      // 如果没有登录且不在认证页面，重定向到登录页
      if (!initialState?.currentUser && !isAuthPage) {
        history.push(loginPath);
      }
    },
    
    // 移除背景装饰图片
    bgLayoutImgList: [],
    
    // 开发环境显示的链接
    links: isDev ? [
      <Link key="backend-admin" to="http://localhost:8000/api/v1/docs" target="_blank">
        <LinkOutlined />
        <span>后端API文档</span>
      </Link>,
    ] : [],
    
    // 隐藏菜单头部
    menuHeaderRender: undefined,
    
    // 子组件渲染
    childrenRender: (children) => {
      return (
        <QueryClientProvider client={queryClient}>
          {children}
          {/* 开发环境显示React Query调试工具 */}
          {isDev && <ReactQueryDevtools initialIsOpen={false} />}
        </QueryClientProvider>
      );
    },
    
    // 应用默认设置
    ...initialState?.settings,
  };
};

/**
 * 网络请求配置
 * 这里使用我们自定义的axios配置
 */
export const request = {
  // 由于我们使用了自定义的axios实例，这里留空
  // 实际的请求配置在 src/utils/request.ts 中
};