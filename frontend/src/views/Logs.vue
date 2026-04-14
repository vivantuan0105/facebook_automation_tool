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
      <div class="overflow-x-auto w-full">
        <table class="w-full text-left text-sm text-light-text whitespace-nowrap">
        <thead class="bg-indigo-50/60 text-slate-700 font-semibold border-b-2 border-indigo-100/50 text-xs uppercase">
          <tr>
            <th class="px-6 py-2.5 w-48 border-r border-indigo-100/30 text-center">Thời gian</th>
            <th class="px-6 py-2.5 w-32 border-r border-indigo-100/30 text-center">Trạng thái</th>
            <th class="px-6 py-2.5 w-32 border-r border-indigo-100/30 text-center">Phân hệ</th>
            <th class="px-6 py-2.5 w-32 border-r border-indigo-100/30 text-center">Hành động</th>
            <th class="px-6 py-2.5 text-center">Chi tiết thực thi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="store.logs.length === 0">
            <td colspan="5" class="px-6 py-8 text-center text-light-muted">Chưa có bản ghi nhật ký.</td>
          </tr>
          <tr v-for="log in store.logs" :key="log.id" class="border-b border-slate-100 hover:bg-slate-50 transition-colors">
            <td class="px-6 py-2 text-slate-500 border-r border-slate-100 text-center">{{ new Date(log.time).toLocaleString() }}</td>
            <td class="px-6 py-2 border-r border-slate-100 text-center">
              <span class="inline-block px-2.5 py-0.5 rounded-full text-[11px] font-medium border" :class="{
                'bg-emerald-50 text-emerald-600 border-emerald-200': log.status === 'Success',
                'bg-rose-50 text-rose-600 border-rose-200': log.status === 'Error' || log.status === 'Failed',
                'bg-amber-50 text-amber-600 border-amber-200': log.status === 'Warning',
                'bg-indigo-50 text-indigo-600 border-indigo-200': log.status === 'Info' || log.status === 'Running'
              }">{{ log.status }}</span>
            </td>
            <td class="px-6 py-2 border-r border-slate-100 text-slate-600 text-left">{{ log.module }}</td>
            <td class="px-6 py-2 text-slate-700 font-medium border-r border-slate-100 text-left">{{ log.action }}</td>
            <td class="px-6 py-2 text-slate-500 truncate max-w-sm text-left" :title="log.details">{{ log.details }}</td>
          </tr>
        </tbody>
      </table>
      </div>
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
