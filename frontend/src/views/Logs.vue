<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h2 class="text-xl font-semibold text-light-text">Nhật ký Hệ thống (Logs)</h2>
      <div class="flex space-x-2">
        <select class="bg-light-surface border border-light-border text-light-text text-sm rounded-lg focus:ring-primary focus:border-primary block p-2 transition-colors">
          <option>Tất cả Module</option>
          <option>Task Center</option>
          <option>Hàng Đợi</option>
          <option>Hệ thống</option>
        </select>
        <button class="px-3 py-2 bg-light-surface hover:bg-light-border border border-light-border text-light-text rounded-lg text-sm transition-colors shadow-sm" @click="store.fetchLogs()">
          Làm mới
        </button>
      </div>
    </div>

    <div class="bg-light-surface border border-light-border rounded-xl shadow-sm overflow-hidden">
      <table class="w-full text-left text-sm text-light-text">
        <thead class="text-xs text-light-muted uppercase bg-light-bg border-b border-light-border">
          <tr>
            <th class="px-6 py-4 font-medium w-48">Thời gian</th>
            <th class="px-6 py-4 font-medium w-32">Trạng thái</th>
            <th class="px-6 py-4 font-medium w-32">Phân hệ</th>
            <th class="px-6 py-4 font-medium w-32">Hành động</th>
            <th class="px-6 py-4 font-medium">Chi tiết thực thi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="store.logs.length === 0">
            <td colspan="5" class="px-6 py-8 text-center text-light-muted">Chưa có bản ghi nhật ký.</td>
          </tr>
          <tr v-for="log in store.logs" :key="log.id" class="border-b border-light-border hover:bg-light-bg/50 transition-colors">
            <td class="px-6 py-4 text-light-muted">{{ new Date(log.time).toLocaleString() }}</td>
            <td class="px-6 py-4">
              <span class="px-2.5 py-1 rounded-full text-xs font-medium border" :class="{
                'bg-green-50 text-green-600 border-green-200': log.status === 'Success',
                'bg-red-50 text-red-600 border-red-200': log.status === 'Error' || log.status === 'Failed',
                'bg-yellow-50 text-yellow-600 border-yellow-200': log.status === 'Warning',
                'bg-blue-50 text-blue-600 border-blue-200': log.status === 'Info' || log.status === 'Running'
              }">{{ log.status }}</span>
            </td>
            <td class="px-6 py-4">{{ log.module }}</td>
            <td class="px-6 py-4 text-light-text font-medium">{{ log.action }}</td>
            <td class="px-6 py-4 text-light-muted truncate max-w-sm" :title="log.details">{{ log.details }}</td>
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

onMounted(() => {
  store.fetchLogs()
})
</script>
