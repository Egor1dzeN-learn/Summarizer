export interface User {
  name: string
}

export interface Chat {
  id: number
  title: string
  text: string
  entries: ChatEntry[]
}

export interface ChatEntry {
  id: number
  chat_id: number

  question: string
  answer?: string
}
