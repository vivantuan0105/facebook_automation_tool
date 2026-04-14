<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h2 class="text-xl font-semibold text-light-text">Hàng đợi tác vụ</h2>
      <button @click="store.fetchTasks()" class="px-3 py-2 bg-light-surface hover:bg-light-border border border-light-border text-light-text rounded-lg text-sm transition-colors shadow-sm">
        Làm mới
      </button>
    </div>

    <div class="bg-light-surface border border-light-border rounded-xl shadow-sm overflow-hidden">
      <div v-if="store.tasks.length === 0" class="text-light-muted text-center py-8">
        Chưa có tác vụ nào. Vui lòng chuyển sang Trung tâm tác vụ để tạo mới.
      </div>
      <div class="overflow-x-auto" v-else>
        <table class="w-full text-left text-sm text-light-text">
          <thead class="text-xs text-light-muted uppercase bg-light-bg border-b border-light-border">
            <tr>
              <th class="px-4 py-3 font-medium">ID</th>
              <th class="px-4 py-3 font-medium">Loại</th>
              <th class="px-4 py-3 font-medium">Hành động</th>
              <th class="px-4 py-3 font-medium">Mục tiêu</th>
              <th class="px-4 py-3 font-medium">Khởi tạo lúc</th>
              <th class="px-4 py-3 font-medium">Trạng thái</th>
              <th class="px-4 py-3 font-medium text-right">Tuỳ chọn</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="task in store.tasks" :key="task.id" class="border-b border-light-border hover:bg-light-bg/50 transition-colors">
              <td class="px-4 py-3">{{ task.id }}</td>
              <td class="px-4 py-3 font-medium">{{ task.taskType }}</td>
              <td class="px-4 py-3">{{ task.reactionType }}</td>
              <td class="px-4 py-3 truncate max-w-xs text-primary hover:underline cursor-pointer" :title="task.postUrl || task.postId">
                {{ task.postUrl ? 'URL' : 'ID: ' + task.postId }}
              </td>
              <td class="px-4 py-3 text-light-muted">{{ new Date(task.createdAt).toLocaleString() }}</td>
              <td class="px-4 py-3">
                <span class="px-2 py-1 rounded text-xs font-medium border" :class="{
                  'bg-yellow-50 text-yellow-600 border-yellow-200': task.status === 'Pending',
                  'bg-purple-50 text-purple-600 border-purple-200': task.status === 'Running',
                  'bg-green-50 text-green-600 border-green-200': task.status === 'Success',
                  'bg-red-50 text-red-600 border-red-200': task.status === 'Failed',
                  'bg-gray-100 text-gray-500 border-gray-200': task.status === 'Paused'
                }">
                  {{ task.status }}
                </span>
              </td>
              <td class="px-4 py-3 text-right space-x-2">
                <button v-if="task.status === 'Pending' || task.status === 'Paused'" @click="runTask(task.id)" class="text-primary hover:underline text-xs">Chạy</button>
                <button v-if="task.status === 'Running'" @click="store.pauseTask(task.id)" class="text-yellow-600 hover:underline text-xs">Dừng</button>
                <button v-if="task.status === 'Failed'" @click="store.retryTask(task.id)" class="text-blue-600 hover:underline text-xs">Thử lại</button>
                <button @click="store.removeTask(task.id)" class="text-red-500 hover:underline text-xs">Xóa</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useMainStore } from '../stores/main'

const store = useMainStore()
let interval: any

onMounted(() => {
  store.fetchTasks()
  interval = setInterval(() => {
    store.fetchTasks()
  }, 2000)
})

onUnmounted(() => {
  clearInterval(interval)
})

const runTask = async (id: number) => {
  await store.runTaskNow(id, 'live')
}
</script>
