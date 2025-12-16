<script setup lang="ts">
import type { GlobalThemeOverrides } from 'naive-ui'
const themeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#2472b9',
    primaryColorHover: '#40b8e8',
  },
}
</script>

<template>
  <n-config-provider abstract :theme-overrides="themeOverrides">
    <router-view />
  </n-config-provider>
</template>

<style scoped></style>
