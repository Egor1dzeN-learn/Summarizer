<script setup lang="ts">
import { useChatsStore } from '@/stores/chats'
import { SendOutline as SendIcon } from '@vicons/ionicons5'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const store = useChatsStore()

const isInProgress = ref(false)

const prompt = ref('')
const doSubmit = async () => {
  isInProgress.value = true
  try {
    const id = (await store.createNewChat(prompt.value)).id
    router.push({
      name: 'chat',
      params: {
        id,
      },
    })
  } finally {
    isInProgress.value = false
  }
}
</script>

<template>
  <n-flex vertical align="center" justify="center" class="fill-height">
    <n-h1>Суммаризатор</n-h1>
    <n-spin :show="isInProgress" style="width: 100%">
      <n-flex vertical align="center" justify="center">
        <n-input
          v-model:value="prompt"
          type="textarea"
          size="large"
          round
          placeholder="Введите текст для обработки"
          style="max-width: 800px; margin: auto"
        >
          <template #suffix>
            <n-flex vertical justify="flex-end" class="fill-height" style="margin-bottom: 16px">
              <n-button tertiary circle type="primary" @click="doSubmit">
                <template #icon>
                  <n-icon><SendIcon /></n-icon>
                </template>
              </n-button>
            </n-flex>
          </template>
        </n-input>
      </n-flex>
    </n-spin>
  </n-flex>
</template>
