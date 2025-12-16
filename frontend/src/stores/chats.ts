import type { Chat } from '../models'
import { createNewChat, fetchChats, sendMessage, deleteChat } from '../services/chat'
import { defineStore } from 'pinia'

interface State {
  chats: Chat[] | null
  intervalId: number | null
}

export const useChatsStore = defineStore('chats', {
  state: (): State => ({
    chats: null,
    intervalId: null,
  }),
  actions: {
    async fetchData() {
      try {
        this.chats = await fetchChats()
      } catch {
        if (this.chats === null) {
          this.chats = []
        }
      }
    },
    startPolling(interval: number = 5000) {
      if (this.intervalId) {
        clearInterval(this.intervalId)
      }
      this.intervalId = setInterval(() => {
        this.fetchData()
      }, interval)
    },
    stopPolling() {
      if (this.intervalId) {
        clearInterval(this.intervalId)
        this.intervalId = null
      }
    },

    async createNewChat(prompt: string): Promise<Chat> {
      const chat = await createNewChat(prompt)
      this.chats = [...(this.chats ?? []), chat]
      return chat
    },
    async sendMessage(chatId: number, msg: string): Promise<void> {
      const entry = await sendMessage(chatId, msg)
      this.chats =
        this.chats?.map((chat) => {
          if (chat.id == chatId) {
            chat.entries = [...chat.entries, entry]
          }
          return chat
        }) ?? []
    },
    async deleteChat(chatId: number): Promise<void> {
      await deleteChat(chatId)
      this.chats = this.chats!.filter((chat) => chat.id != chatId)
    },
  },
})
