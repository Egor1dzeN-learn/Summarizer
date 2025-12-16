import api from '../api'
import type { User } from '../models'

export const getMe = async (): Promise<User> =>
  (await api.get('/me', { validateStatus: (status) => status == 200 })).data

export const login = async (token: string): Promise<User> =>
  (await api.post('/login', { token }, { validateStatus: (status) => status == 200 })).data

export const logout = async (): Promise<void> => await api.post('/logout')
