<template>
  <span :class="[
    'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border',
    colorClasses[status?.toLowerCase()] || colorClasses['default']
  ]">
    {{ status }}
  </span>
</template>

<script setup lang="ts">
defineProps({
  status: {
    type: String,
    required: true
  }
})

const colorClasses: Record<string, string> = {
  'published': 'bg-green-500/10 text-green-400 border-green-500/20',
  'success': 'bg-green-500/10 text-green-400 border-green-500/20',
  'scheduled': 'bg-blue-500/10 text-blue-400 border-blue-500/20',
  'info': 'bg-blue-500/10 text-blue-400 border-blue-500/20',
  'draft': 'bg-gray-500/10 text-gray-400 border-gray-500/20',
  'error': 'bg-red-500/10 text-red-400 border-red-500/20',
  'failed': 'bg-red-500/10 text-red-400 border-red-500/20',
  'warning': 'bg-orange-500/10 text-orange-400 border-orange-500/20',
  'default': 'bg-gray-500/10 text-gray-400 border-gray-500/20'
}
</script>
