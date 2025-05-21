import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export interface User {
  id: string
  tenant_id: string
  email: string
  name: string
  avatar_url?: string
  status: 'active' | 'inactive' | 'suspended'
  created_at: string
  updated_at: string
}

interface AuthState {
  user: User | null
  token: string | null
  refreshToken: string | null
  isAuthenticated: boolean
  login: (user: User, token: string, refreshToken: string) => void
  logout: () => void
  updateUser: (user: Partial<User>) => void
  updateToken: (token: string) => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      refreshToken: null,
      isAuthenticated: false,
      login: (user, token, refreshToken) => 
        set({ user, token, refreshToken, isAuthenticated: true }),
      logout: () => 
        set({ user: null, token: null, refreshToken: null, isAuthenticated: false }),
      updateUser: (userData) => 
        set((state) => ({
          user: state.user ? { ...state.user, ...userData } : null
        })),
      updateToken: (token) => 
        set({ token })
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({ 
        user: state.user, 
        token: state.token, 
        refreshToken: state.refreshToken,
        isAuthenticated: state.isAuthenticated
      })
    }
  )
)
