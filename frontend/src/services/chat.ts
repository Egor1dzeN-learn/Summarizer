import api from '../api'
import type { Chat, ChatEntry } from '../models'

export const fetchChats = async (): Promise<Chat[]> => (await api.get('/chats')).data

export const createNewChat = async (prompt: string): Promise<Chat> =>
  (await api.post('/chats', { prompt })).data

export const sendMessage = async (chatId: number, msg: string): Promise<ChatEntry> =>
  (await api.post(`/chats/${chatId}`, { msg })).data

export const deleteChat = async (chatId: number): Promise<void> =>
  await api.delete(`/chats/${chatId}`)
