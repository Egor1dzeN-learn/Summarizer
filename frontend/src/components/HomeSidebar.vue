<script setup lang="ts">
import type { MenuOption } from 'naive-ui'
import type { Component } from 'vue'
import { AddCircleOutline as NewIcon } from '@vicons/ionicons5'
import { PaperPlane as TelegramIcon } from '@vicons/ionicons5'
import { NIcon } from 'naive-ui'
import { computed, h, ref } from 'vue'
import { useChatsStore } from '@/stores/chats'
import { RouterLink, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()

const store = useChatsStore()
const authStore = useAuthStore()

const activeKey = computed(() =>
  router.currentRoute.value.name == 'chat'
    ? <string>router.currentRoute.value.params.id
    : <string>router.currentRoute.value.name,
)
const collapsed = ref(false)

function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

const logout = () => {
  authStore.logout()
  router.push({ name: 'login' })
}

const menuOptions = computed(
  () =>
    <MenuOption[]>[
      {
        label: () =>
          h(
            RouterLink,
            {
              to: { name: 'blank' },
            },
            { default: () => 'Новый чат' },
          ),
        key: 'blank',
        icon: renderIcon(NewIcon),
      },
      {
        type: 'group',
        key: 'history',
        label: 'История запросов',
        show: !collapsed.value,
        children: [
          {
            type: 'divider',
          },
          ...store.chats!.map((chat) => ({
            label: () =>
              h(
                RouterLink,
                {
                  to: {
                    name: 'chat',
                    params: {
                      id: chat.id,
                    },
                  },
                },
                { default: () => chat.title },
              ),
            key: chat.id.toString(),
          })),
        ],
      },
    ],
)
</script>

<template>
  <n-layout-sider
    bordered
    collapse-mode="width"
    :collapsed-width="64"
    :width="300"
    :collapsed="collapsed"
    show-trigger
    @collapse="collapsed = true"
    @expand="collapsed = false"
  >
    <n-flex vertical style="height: 100%">
      <n-flex align="center" style="padding: 24px 32px 16px 32px" size="small">
        <n-h6 v-if="!collapsed" style="margin: 0">Суммаризатор</n-h6>
      </n-flex>
      <n-menu
        :value="activeKey"
        :collapsed="collapsed"
        :collapsed-width="64"
        :options="menuOptions"
      />
      <div style="flex: 1" />
      <div style="padding: 4px">
        <n-button color="#40a7e3" @click="logout">
          <template #icon>
            <n-icon>
              <TelegramIcon />
            </n-icon>
          </template>
          {{ authStore.user?.name }}
        </n-button>
      </div>
    </n-flex>
  </n-layout-sider>
</template>

<style>
.n-menu-item-group-title {
  white-space: nowrap;
}
</style>
