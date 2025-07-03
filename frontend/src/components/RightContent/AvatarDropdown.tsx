/**
 * LYSS AI Platform 用户头像下拉菜单组件
 * 提供用户信息显示、个人设置和登出功能
 */

import React, { useCallback } from 'react';
import { LogoutOutlined, SettingOutlined, UserOutlined } from '@ant-design/icons';
import { history, useModel } from '@umijs/max';
import { Spin, message } from 'antd';
import { createStyles } from 'antd-style';
import type { MenuInfo } from 'rc-menu/lib/interface';
import { useAuthStore } from '@/stores/auth';
import HeaderDropdown from '../HeaderDropdown';

export type GlobalHeaderRightProps = {
  menu?: boolean;
  children?: React.ReactNode;
};

// 用户名显示组件
export const AvatarName = () => {
  const { initialState } = useModel('@@initialState');
  const { currentUser } = initialState || {};
  
  // 显示用户邮箱或姓名
  const displayName = currentUser?.first_name 
    ? `${currentUser.first_name} ${currentUser.last_name || ''}`.trim()
    : currentUser?.email?.split('@')[0] || '用户';
    
  return <span className="anticon">{displayName}</span>;
};

const useStyles = createStyles(({ token }) => {
  return {
    action: {
      display: 'flex',
      height: '48px',
      marginLeft: 'auto',
      overflow: 'hidden',
      alignItems: 'center',
      padding: '0 8px',
      cursor: 'pointer',
      borderRadius: token.borderRadius,
      '&:hover': {
        backgroundColor: token.colorBgTextHover,
      },
    },
  };
});

export const AvatarDropdown: React.FC<GlobalHeaderRightProps> = ({ menu, children }) => {
  const { styles } = useStyles();
  const { initialState, setInitialState } = useModel('@@initialState');
  const { logout } = useAuthStore();

  /**
   * 处理用户登出
   */
  const handleLogout = useCallback(async () => {
    try {
      // 清理认证状态
      logout();
      
      // 清理初始状态
      setInitialState((s) => ({ ...s, currentUser: undefined }));
      
      message.success('已成功退出登录');
      
      // 跳转到登录页面，保存当前页面用于登录后跳转
      const { pathname, search } = window.location;
      if (pathname !== '/user/login') {
        const redirect = encodeURIComponent(pathname + search);
        history.replace(`/user/login?redirect=${redirect}`);
      }
    } catch (error) {
      console.error('登出失败:', error);
      message.error('登出失败，请重试');
    }
  }, [logout, setInitialState]);

  /**
   * 处理菜单点击事件
   */
  const onMenuClick = useCallback(
    (event: MenuInfo) => {
      const { key } = event;
      
      switch (key) {
        case 'logout':
          handleLogout();
          break;
        case 'settings':
          history.push('/settings');
          break;
        case 'center':
          // 个人中心暂时跳转到设置页
          history.push('/settings');
          break;
        default:
          break;
      }
    },
    [handleLogout],
  );

  // 加载状态
  const loading = (
    <span className={styles.action}>
      <Spin
        size="small"
        style={{
          marginLeft: 8,
          marginRight: 8,
        }}
      />
    </span>
  );

  // 如果没有初始状态，显示加载中
  if (!initialState) {
    return loading;
  }

  const { currentUser } = initialState;

  // 如果没有用户信息，显示加载中
  if (!currentUser || !currentUser.email) {
    return loading;
  }

  // 根据用户角色显示不同的菜单项
  const menuItems = [
    ...(menu
      ? [
          {
            key: 'center',
            icon: <UserOutlined />,
            label: '个人中心',
          },
          {
            key: 'settings',
            icon: <SettingOutlined />,
            label: '个人设置',
          },
          {
            type: 'divider' as const,
          },
        ]
      : []),
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
    },
  ];

  return (
    <HeaderDropdown
      menu={{
        selectedKeys: [],
        onClick: onMenuClick,
        items: menuItems,
      }}
    >
      {children}
    </HeaderDropdown>
  );
};