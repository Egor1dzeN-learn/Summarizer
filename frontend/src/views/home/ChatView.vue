<script setup lang="ts">
import { useChatsStore } from '@/stores/chats'
import { computed, ref, watch } from 'vue'
import { SendOutline as SendIcon, TrashOutline as DeleteIcon } from '@vicons/ionicons5'
import { useRoute, useRouter } from 'vue-router'

const store = useChatsStore()

const router = useRouter()
const route = useRoute()
const params = computed(() => route.params)

const chat = computed(() => store.chats!.find((chat) => chat.id.toString() == params.value.id))

const containerRef = ref(null)
const prompt = ref('')

const scrollDown = () => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const el = <HTMLElement>(<any>containerRef.value)?.$el
  el.scrollTo({
    top: el.scrollHeight,
    behavior: 'smooth',
  })
}

const doSubmit = () => {
  store.sendMessage(chat.value!.id, prompt.value)
  prompt.value = ''
  scrollDown()
}
const doDelete = async () => {
  await store.deleteChat(chat.value!.id)
  router.push({ name: 'blank' })
}

watch(chat, () => scrollDown)
</script>

<template>
  <template v-if="chat">
    <n-flex vertical class="fill-height" align="center" justify="center">
      <n-flex vertical class="fill-height" style="width: 100%; max-width: 600px">
        <div style="display: flex; flex-direction: row; width: 100%; justify-content: end">
          <n-button type="error" ghost @click="doDelete"
            >Удалить чат
            <template #icon>
              <n-icon>
                <DeleteIcon />
              </n-icon>
            </template>
          </n-button>
        </div>
        <n-page-header :subtitle="chat.title" />
        <n-flex ref="containerRef" vertical style="flex: 1; overflow: auto">
          <n-text v-if="chat.text.length > chat.title.length">{{ chat.text }}</n-text>
          <template v-if="chat.entries && !!chat.entries.length">
            <template v-for="entry in chat.entries" v-bind:key="entry.id">
              <n-card :embedded="false" size="small" class="msg msg-sent">
                {{ entry.question }}
              </n-card>
              <n-card v-if="entry.answer" :embedded="true" size="small" class="msg msg-recv">
                {{ entry.answer }}
              </n-card>
              <n-skeleton
                v-else
                height="48px"
                style="min-height: 48px; margin-bottom: 20px"
                round
                class="msg msg-recv"
              />
            </template>
          </template>
          <n-flex vertical v-else align="center" justify="center">
            <img src="@/assets/chat.svg" style="width: 300px" />
            <n-h6>Вы можете задавать вопросы к тексту, используя форму ниже</n-h6>
          </n-flex>
          <br />
          <br />
          <br />
        </n-flex>
        <n-input
          v-model:value="prompt"
          type="textarea"
          size="large"
          round
          placeholder="Вопрос к тексту"
        >
          <template #suffix>
            <n-flex vertical justify="flex-end" class="fill-height" style="margin-bottom: 16px">
              <n-button tertiary circle type="primary" @click="doSubmit" @mousedown.prevent>
                <template #icon>
                  <n-icon><SendIcon /></n-icon>
                </template>
              </n-button>
            </n-flex>
          </template>
        </n-input>
      </n-flex>
    </n-flex>
  </template>

  <n-flex v-else vertical align="center" justify="center" style="height: 100%">
    <n-spin size="large" />
  </n-flex>
</template>

<style scoped>
.msg {
  border-radius: 24px;
  width: 40%;
}
.msg-recv {
  align-self: flex-start;
  border-top-left-radius: 0;
}
.msg-sent {
  align-self: flex-end;
  border-bottom-right-radius: 0;
}
</style>
