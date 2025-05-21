import { useState } from 'react'
import { Outlet } from 'react-router-dom'
import { Layout, Menu, Avatar, Dropdown, Button } from 'antd'
import { 
  MenuFoldOutlined, 
  MenuUnfoldOutlined, 
  UserOutlined,
  LogoutOutlined,
  SettingOutlined
} from '@ant-design/icons'
import { useAuthStore } from '../../store/authStore'
import styles from './MainLayout.module.css'

const { Header, Sider, Content } = Layout

const MainLayout = () => {
  const [collapsed, setCollapsed] = useState(false)
  const { user, logout } = useAuthStore()

  const handleLogout = () => {
    logout()
  }

  const userMenu = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人信息'
    },
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '设置'
    },
    {
      type: 'divider'
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout
    }
  ]

  return (
    <Layout className={styles.layout}>
      <Sider 
        trigger={null} 
        collapsible 
        collapsed={collapsed}
        width={256}
        className={styles.sider}
      >
        <div className={styles.logo}>
          {!collapsed && <span>Lyss Chat</span>}
        </div>
        <Menu
          theme="dark"
          mode="inline"
          defaultSelectedKeys={['1']}
          items={[
            {
              key: '1',
              icon: <UserOutlined />,
              label: '聊天',
            }
          ]}
        />
      </Sider>
      <Layout>
        <Header className={styles.header}>
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
            className={styles.trigger}
          />
          <div className={styles.headerRight}>
            <Dropdown menu={{ items: userMenu }} placement="bottomRight">
              <div className={styles.userInfo}>
                <span className={styles.userName}>{user?.name}</span>
                <Avatar icon={<UserOutlined />} src={user?.avatar_url} />
              </div>
            </Dropdown>
          </div>
        </Header>
        <Content className={styles.content}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}

export default MainLayout
