<template>
  <header class="h-16 bg-light-surface border-b border-light-border flex items-center justify-between px-6 z-10 sticky top-0">
    <div class="flex items-center space-x-6">
      <h2 class="text-lg font-medium text-light-text capitalize">{{ routeNameVn }}</h2>
      <div class="hidden md:flex items-center space-x-3 text-xs border-l border-light-border pl-6">
        <span class="text-light-muted flex items-center gap-1">
          Hàng đợi: 
          <span class="font-medium text-yellow-600 cursor-pointer hover:underline" @click="store.filterStatusTrigger = 'Pending'" title="Nhấp để click chọn các tài khoản Đang chờ bên dưới">
            {{ store.overview?.pending || 0 }} Đang chờ
          </span>,
          <span class="font-medium text-purple-600 cursor-pointer hover:underline" @click="store.filterStatusTrigger = 'Running'" title="Nhấp để click chọn các tài khoản Đang chạy bên dưới">
            {{ store.overview?.running || 0 }} Đang chạy
          </span>
        </span>
      </div>
    </div>
    <div class="flex items-center space-x-4">
      <div class="h-8 w-8 rounded-full bg-gradient-to-tr from-primary to-blue-400 p-[2px]">
        <div class="h-full w-full rounded-full bg-light-surface flex items-center justify-center">
          <UserIcon class="w-4 h-4 text-light-text" />
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { UserIcon } from '@heroicons/vue/24/outline'
import { useMainStore } from '../stores/main'

const route = useRoute()
const router = useRouter()
const store = useMainStore()

const routeNameVn = computed(() => {
  const name = route.name ? route.name.toString() : 'Overview'
  const map: Record<string, string> = {
    'Overview': 'Tổng quan',
    'Task Center': 'Trung tâm Tác vụ',
    'Accounts': 'Tài khoản (Clone)',
    'Queue': 'Biên chế Hàng đợi',
    'Logs': 'Nhật ký Hệ thống',
    'Settings': 'Cài đặt'
  }
  return map[name] || name
})

const currentMode = computed(() => {
  return 'Thực thi thật (Live)'
})
</script>
