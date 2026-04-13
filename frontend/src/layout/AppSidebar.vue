<template>
  <aside class="w-64 bg-dark-surface border-r border-dark-border flex flex-col transition-all duration-300">
    <div class="h-16 flex items-center px-6 border-b border-dark-border select-none">
      <div class="w-8 h-8 rounded bg-primary flex items-center justify-center mr-3 font-bold text-white shadow-lg shadow-primary/30">SD</div>
      <h1 class="text-base font-semibold text-white tracking-wide">Social Manager</h1>
    </div>
    
    <nav class="flex-1 py-6 px-3 space-y-1 overflow-y-auto">
      <router-link 
        v-for="item in navItems" 
        :key="item.name"
        :to="item.path"
        class="flex items-center px-3 py-2.5 rounded-lg text-dark-muted hover:text-white hover:bg-dark-border/50 transition-colors group"
        active-class="bg-dark-border/80 text-white shadow-sm font-medium"
      >
        <component :is="item.icon" class="w-5 h-5 mr-3 group-hover:text-primary transition-colors" :class="{'text-primary': $route.path.includes(item.path)}" />
        {{ item.name }}
      </router-link>
    </nav>
    
    <div class="p-4 border-t border-dark-border">
      <div class="bg-dark-bg p-3 rounded-lg border border-dark-border flex items-center">
        <div class="w-2 h-2 rounded-full bg-green-500 mr-2 animate-pulse"></div>
        <span class="text-xs text-dark-muted font-medium">System Online</span>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { 
  HomeIcon, 
  DocumentTextIcon, 
  CpuChipIcon, 
  QueueListIcon, 
  Cog8ToothIcon 
} from '@heroicons/vue/24/outline'

const navItems = [
  { name: 'Dashboard', path: '/dashboard', icon: HomeIcon },
  { name: 'Posts', path: '/posts', icon: DocumentTextIcon },
  { name: 'Automation', path: '/automation', icon: CpuChipIcon },
  { name: 'Logs', path: '/logs', icon: QueueListIcon },
  { name: 'Settings', path: '/settings', icon: Cog8ToothIcon },
]
</script>
