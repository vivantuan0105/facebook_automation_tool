<template>
  <header class="h-16 bg-dark-surface border-b border-dark-border flex items-center justify-between px-6 z-10 sticky top-0">
    <div class="flex items-center">
      <h2 class="text-lg font-medium text-white capitalize">{{ routeName }}</h2>
    </div>
    <div class="flex items-center space-x-4">
      <button class="p-2 rounded-full text-dark-muted hover:text-white hover:bg-dark-border transition-colors">
        <BellIcon class="w-5 h-5" />
      </button>
      <div class="h-8 w-8 rounded-full bg-gradient-to-tr from-primary to-blue-400 p-[2px]">
        <div class="h-full w-full rounded-full bg-dark-surface flex items-center justify-center">
          <UserIcon class="w-4 h-4 text-white" />
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { BellIcon, UserIcon } from '@heroicons/vue/24/outline'

const route = useRoute()
const routeName = computed(() => {
  return route.name ? route.name.toString() : 'Dashboard'
})
</script>
