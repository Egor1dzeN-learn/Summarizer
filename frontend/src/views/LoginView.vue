<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { PaperPlane as TelegramIcon } from '@vicons/ionicons5'
import { useRoute, useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()
console.log(route.query)
if (route.query.key) {
  ;(async () => {
    await authStore.login(<string>route.query.key)
    router.push({ name: 'home' })
  })()
}
</script>

<template>
  <n-flex vertical align="center" justify="center" class="fill-height">
    <n-h1>Суммаризатор</n-h1>
    <n-button
      tag="a"
      href="https://t.me/summarizer_af41eedb_bot"
      target="_blank"
      color="#40a7e3"
      size="large"
    >
      <template #icon>
        <n-icon>
          <TelegramIcon />
        </n-icon>
      </template>
      Войти через Telegram
    </n-button>
  </n-flex>
</template>

<style scoped></style>
