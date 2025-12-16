import type { User } from '../models'
import { getMe, login, logout } from '../services/auth'
import { defineStore } from 'pinia'

interface State {
  initialized: boolean
  user: User | null
}

export const useAuthStore = defineStore('auth', {
  state: (): State => ({
    initialized: false,
    user: null,
  }),
  getters: {
    isAuthenticated: (state) => !!state.user,
  },
  actions: {
    async fetchUser() {
      try {
        this.user = await getMe()
      } catch {
        this.user = null
      } finally {
        this.initialized = true
      }
    },
    async login(token: string) {
      try {
        this.user = await login(token)
      } catch {
        this.user = null
      }
    },
    async logout() {
      await logout()
      this.user = null
    },
  },
})
