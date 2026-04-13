<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h2 class="text-xl font-semibold text-white">System Logs</h2>
      <div class="flex space-x-2">
        <select class="bg-dark-surface border border-dark-border text-white text-sm rounded-lg focus:ring-primary focus:border-primary block p-2 transition-colors">
          <option>All Modules</option>
          <option>System</option>
          <option>Automation</option>
        </select>
        <button class="px-3 py-2 bg-dark-surface hover:bg-dark-border border border-dark-border text-white rounded-lg text-sm transition-colors" @click="store.fetchLogs()">
          Refresh
        </button>
      </div>
    </div>

    <div class="bg-dark-surface border border-dark-border rounded-xl shadow-sm overflow-hidden">
      <table class="w-full text-left text-sm text-dark-muted">
        <thead class="text-xs text-white uppercase bg-dark-bg border-b border-dark-border">
          <tr>
            <th class="px-6 py-4 font-medium w-48">Time</th>
            <th class="px-6 py-4 font-medium w-32">Status</th>
            <th class="px-6 py-4 font-medium w-32">Module</th>
            <th class="px-6 py-4 font-medium w-64">Action</th>
            <th class="px-6 py-4 font-medium">Details</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="store.logs.length === 0">
            <td colspan="5" class="px-6 py-8 text-center text-dark-muted">No logs recorded.</td>
          </tr>
          <tr v-for="log in store.logs" :key="log.id" class="border-b border-dark-border/50 hover:bg-dark-bg/50 transition-colors font-mono text-xs">
            <td class="px-6 py-4">{{ new Date(log.time).toLocaleString() }}</td>
            <td class="px-6 py-4">
              <span :class="{
                'text-green-400': log.status === 'Success',
                'text-red-400': log.status === 'Error',
                'text-yellow-400': log.status === 'Warning',
                'text-blue-400': log.status === 'Info'
              }">{{ log.status }}</span>
            </td>
            <td class="px-6 py-4">{{ log.module }}</td>
            <td class="px-6 py-4 text-white">{{ log.action }}</td>
            <td class="px-6 py-4 text-dark-muted truncate max-w-sm" :title="log.details">{{ log.details }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useMainStore } from '../stores/main'

const store = useMainStore()

onMounted(async () => {
  await store.fetchLogs()
})
</script>
