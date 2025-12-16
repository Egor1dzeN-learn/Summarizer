import axios from 'axios'
import { createDiscreteApi } from 'naive-ui'

const { message } = createDiscreteApi(['message'])
const NETWORK_ERROR_COOLDOWN_MS = 3000
let lastNetworkErrorAt = 0

const instance = axios.create({
  baseURL: '/api',
  withCredentials: true,
})

instance.interceptors.response.use(
  (response) => response,
  (error) => {
    const isNetworkError = error.code === 'ERR_NETWORK' || !error.response
    const now = Date.now()
    if (isNetworkError && now - lastNetworkErrorAt > NETWORK_ERROR_COOLDOWN_MS) {
      message.error('Не удалось подключиться к серверу. Проверьте сеть и попробуйте снова.')
      lastNetworkErrorAt = now
    }
    return Promise.reject(error)
  },
)

export default instance
