import { Routes, Route, Navigate } from 'react-router-dom'
import { XProvider } from '@ant-design/x'
import { ConfigProvider } from 'antd'
import zhCN from 'antd/lib/locale/zh_CN'
import { useAuthStore } from './store/authStore'
import MainLayout from './components/layout/MainLayout'
import LoginPage from './pages/auth/LoginPage'
import RegisterPage from './pages/auth/RegisterPage'
import ChatListPage from './pages/chat/ChatListPage'
import ChatPage from './pages/chat/ChatPage'
import NotFoundPage from './pages/NotFoundPage'

// 受保护的路由组件
const ProtectedRoute = ({ children }: { children: React.ReactNode }) => {
  const { isAuthenticated } = useAuthStore()
  
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }
  
  return <>{children}</>
}

function App() {
  return (
    <XProvider>
      <ConfigProvider locale={zhCN}>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          
          <Route path="/" element={
            <ProtectedRoute>
              <MainLayout />
            </ProtectedRoute>
          }>
            <Route index element={<ChatListPage />} />
            <Route path="chat/:id" element={<ChatPage />} />
          </Route>
          
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </ConfigProvider>
    </XProvider>
  )
}

export default App
