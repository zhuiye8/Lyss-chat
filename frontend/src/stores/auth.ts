/**
 * 用户认证状态管理
 * 使用Zustand管理用户登录状态、用户信息等
 */

import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { User } from '@/types/api';
import { setToken, removeToken } from '@/utils/request';

interface AuthState {
  // 状态
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  
  // 操作方法
  setUser: (user: User) => void;
  setToken: (token: string) => void;
  logout: () => void;
  setLoading: (loading: boolean) => void;
  
  // 权限检查
  hasRole: (role: string | string[]) => boolean;
  canAdmin: () => boolean;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      // 初始状态
      user: null,
      isAuthenticated: false,
      isLoading: false,
      
      // 设置用户信息
      setUser: (user: User) => {
        set({ 
          user, 
          isAuthenticated: true,
          isLoading: false 
        });
      },
      
      // 设置token并更新认证状态
      setToken: (token: string) => {
        setToken(token);
        set({ isAuthenticated: true });
      },
      
      // 登出
      logout: () => {
        removeToken();
        set({ 
          user: null, 
          isAuthenticated: false,
          isLoading: false 
        });
      },
      
      // 设置加载状态
      setLoading: (loading: boolean) => {
        set({ isLoading: loading });
      },
      
      // 检查用户是否有指定角色
      hasRole: (role: string | string[]) => {
        const { user } = get();
        if (!user) return false;
        
        if (Array.isArray(role)) {
          return role.includes(user.role);
        }
        return user.role === role;
      },
      
      // 检查是否为管理员
      canAdmin: () => {
        const { hasRole } = get();
        return hasRole(['admin', 'super_admin']);
      },
    }),
    {
      name: 'auth-storage', // localStorage key
      partialize: (state) => ({
        user: state.user,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);