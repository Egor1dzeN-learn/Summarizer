<script setup lang="ts">
import { useChatsStore } from '@/stores/chats'
import { onBeforeUnmount, onMounted } from 'vue'

const store = useChatsStore()
onMounted(() => {
  store.fetchData()
  store.startPolling()
})
onBeforeUnmount(() => store.stopPolling())
</script>

<template>
  <n-layout v-if="store.chats" has-sider style="height: 100%">
    <home-sidebar />
    <n-layout>
      <n-layout-content style="height: 100%" content-style="padding: 24px;">
        <router-view />
      </n-layout-content>
    </n-layout>
  </n-layout>
  <n-flex v-else vertical align="center" justify="center" style="height: 100%">
    <n-spin size="large" />
  </n-flex>
</template>
