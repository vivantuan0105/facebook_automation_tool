<template>
  <aside class="w-64 flex-shrink-0 bg-light-surface border-r border-light-border flex flex-col h-screen">
    <div class="h-16 flex items-center px-6 border-b border-light-border shrink-0">
      <h1 class="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-primary to-blue-500">Social Manager</h1>
    </div>
    
    <nav class="flex-1 py-6 px-3 space-y-1.5 overflow-y-auto">
      <router-link
        v-for="item in navItems"
        :key="item.path"
        :to="item.path"
        class="flex items-center px-3 py-2.5 rounded-lg text-sm font-medium transition-all duration-200"
        :class="[
          $route.path.startsWith(item.path)
            ? 'bg-primary/10 text-primary'
            : 'text-light-muted hover:bg-light-bg hover:text-light-text'
        ]"
      >
        <component :is="item.icon" class="w-5 h-5 mr-3" :class="[$route.path.startsWith(item.path) ? 'text-primary' : 'text-light-muted']" />
        {{ item.name }}
      </router-link>
    </nav>
    
    <div class="p-4 border-t border-light-border">
      <div class="bg-light-bg rounded-lg p-4">
        <p class="text-xs font-medium text-light-muted uppercase mb-2">Trạng thái hệ thống</p>
        <div class="flex items-center">
          <div class="w-2 h-2 rounded-full bg-green-500 mr-2"></div>
          <span class="text-sm font-medium text-light-text">Bình thường</span>
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { 
  HomeIcon, 
  DocumentTextIcon, 
  CpuChipIcon, 
  Cog8ToothIcon,
  QueueListIcon,
  UsersIcon
} from '@heroicons/vue/24/outline'

const navItems = [
  { name: 'Tổng quan', path: '/overview', icon: HomeIcon },
  { name: 'Tài khoản (Clone)', path: '/accounts', icon: UsersIcon },
  { name: 'Trung tâm tác vụ', path: '/task-center', icon: DocumentTextIcon },
  { name: 'Hàng đợi (Queue)', path: '/queue', icon: QueueListIcon },
  { name: 'Nhật ký (Logs)', path: '/logs', icon: CpuChipIcon },
  { name: 'Cài đặt', path: '/settings', icon: Cog8ToothIcon },
]
</script>
