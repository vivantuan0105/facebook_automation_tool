<template>
  <div class="bg-light-surface rounded-xl border border-light-border p-5 hover:border-primary/50 transition-colors shadow-sm relative overflow-hidden group">
    <div class="absolute -right-4 -top-4 w-24 h-24 bg-primary/5 rounded-full blur-2xl group-hover:bg-primary/10 transition-all"></div>
    <div class="flex items-center justify-between relative z-10">
      <div>
        <p class="text-light-muted text-sm font-medium mb-1">{{ title }}</p>
        <h3 class="text-2xl font-bold text-light-text">{{ value }}</h3>
      </div>
      <div :class="['p-3 rounded-lg flex items-center justify-center', iconBgClass]">
        <component :is="icon" :class="['w-6 h-6', iconColorClass]" />
      </div>
    </div>
    <div v-if="trend" class="mt-4 flex items-center text-xs relative z-10">
      <span :class="trend > 0 ? 'text-green-400' : 'text-red-400'" class="font-medium mr-1 flex items-center">
        <ArrowUpIcon v-if="trend > 0" class="w-3 h-3 mr-0.5" />
        <ArrowDownIcon v-else class="w-3 h-3 mr-0.5" />
        {{ Math.abs(trend) }}%
      </span>
      <span class="text-light-muted">from last week</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ArrowUpIcon, ArrowDownIcon } from '@heroicons/vue/20/solid'

defineProps({
  title: String,
  value: [String, Number],
  icon: Function,
  iconBgClass: {
    type: String,
    default: 'bg-primary/10'
  },
  iconColorClass: {
    type: String,
    default: 'text-primary'
  },
  trend: Number
})
</script>

