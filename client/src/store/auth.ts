import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { api } from '../api/client'

interface User {
  id: number
  username: string
  email: string
  nickname?: string
  role: string
  avatar?: string
}

interface AuthState {
  token: string | null
  user: User | null
  isLoading: boolean
  login: (username: string, password: string) => Promise<void>
  register: (username: string, email: string, password: string) => Promise<void>
  logout: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      user: null,
      isLoading: false,

      login: async (username: string, password: string) => {
        set({ isLoading: true })
        try {
          const response = await api.post('/api/auth/login', { username, password })
          set({ 
            token: response.data.token, 
            user: response.data.user,
            isLoading: false 
          })
        } catch (error) {
          set({ isLoading: false })
          throw error
        }
      },

      register: async (username: string, email: string, password: string) => {
        set({ isLoading: true })
        try {
          const response = await api.post('/api/auth/register', { username, email, password })
          set({ 
            token: response.data.token, 
            user: response.data.user,
            isLoading: false 
          })
        } catch (error) {
          set({ isLoading: false })
          throw error
        }
      },

      logout: () => {
        set({ token: null, user: null })
      },
    }),
    {
      name: 'wacrm-auth',
    }
  )
)
