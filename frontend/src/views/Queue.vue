<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center h-10">
      <h2 class="text-xl font-semibold text-light-text">Hàng đợi tác vụ</h2>
      <div class="flex gap-2" v-if="selectedTasks.length > 0">
        <button @click="runSelected" class="px-3 py-1.5 bg-indigo-50 border border-indigo-200 text-indigo-700 hover:bg-indigo-100 rounded-lg text-sm font-medium transition-colors shadow-sm flex items-center gap-1.5">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM9.555 7.168A1 1 0 008 8v4a1 1 0 001.555.832l3-2a1 1 0 000-1.664l-3-2z" clip-rule="evenodd" />
          </svg>
          Chạy đã chọn ({{selectedTasks.length}})
        </button>
        <button @click="removeSelected" class="px-3 py-1.5 bg-rose-50 border border-rose-200 text-rose-700 hover:bg-rose-100 rounded-lg text-sm font-medium transition-colors shadow-sm flex items-center gap-1.5">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
          </svg>
          Xoá
        </button>
      </div>
    </div>

    <div class="bg-light-surface border border-light-border rounded-xl shadow-sm overflow-hidden">
      <div v-if="store.tasks.length === 0" class="text-light-muted text-center py-8">
        Chưa có tác vụ nào. Vui lòng chuyển sang Trung tâm tác vụ để tạo mới.
      </div>
      <div class="overflow-x-auto w-full" v-else>
        <table class="w-full text-left text-sm text-light-text whitespace-nowrap">
          <thead class="bg-indigo-50/60 text-slate-700 font-semibold border-b-2 border-indigo-100/50 text-xs uppercase">
            <tr>
              <th class="px-4 py-2.5 w-10 border-r border-indigo-100/30 text-center">
                <input type="checkbox" v-model="selectAll" class="rounded border-gray-300 text-indigo-600 focus:ring-indigo-500 w-3.5 h-3.5 shadow-sm">
              </th>
              <th class="px-4 py-2.5 border-r border-indigo-100/30 text-center">ID</th>
              <th class="px-4 py-2.5 border-r border-indigo-100/30 text-center">Loại</th>
              <th class="px-4 py-2.5 border-r border-indigo-100/30 text-center">Hành động</th>
              <th class="px-4 py-2.5 border-r border-indigo-100/30 text-center">Mục tiêu</th>
              <th class="px-4 py-2.5 border-r border-indigo-100/30 text-center">Khởi tạo lúc</th>
              <th class="px-4 py-2.5 text-center">Trạng thái</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="task in store.tasks" :key="task.id" class="border-b border-slate-100 hover:bg-slate-50 transition-colors">
              <td class="px-4 py-2 border-r border-slate-100 text-center">
                <input type="checkbox" v-model="selectedTasks" :value="task.id" class="rounded border-gray-300 text-indigo-600 focus:ring-indigo-500 w-3.5 h-3.5 shadow-sm">
              </td>
              <td class="px-4 py-2 border-r border-slate-100 text-slate-500 text-center">{{ task.id }}</td>
              <td class="px-4 py-2 font-medium border-r border-slate-100 text-slate-700 text-left">{{ task.taskType }}</td>
              <td class="px-4 py-2 border-r border-slate-100 text-slate-600 text-left">{{ task.reactionType }}</td>
              <td class="px-4 py-2 truncate max-w-xs text-blue-600 hover:underline cursor-pointer border-r border-slate-100 text-center" :title="task.postUrl || task.postId">
                {{ task.postUrl ? 'URL' : 'ID: ' + task.postId }}
              </td>
              <td class="px-4 py-2 text-slate-500 border-r border-slate-100 text-center">{{ new Date(task.createdAt).toLocaleString() }}</td>
              <td class="px-4 py-2 text-center">
                <span class="inline-block px-2 py-0.5 rounded text-[11px] font-medium border" :class="{
                  'bg-amber-50 text-amber-600 border-amber-200': task.status === 'Pending',
                  'bg-indigo-50 text-indigo-600 border-indigo-200': task.status === 'Running',
                  'bg-emerald-50 text-emerald-600 border-emerald-200': task.status === 'Success',
                  'bg-rose-50 text-rose-600 border-rose-200': task.status === 'Failed',
                  'bg-slate-100 text-slate-500 border-slate-200': task.status === 'Paused'
                }">
                  {{ task.status }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useMainStore } from '../stores/main'

const store = useMainStore()
let interval: any

const selectedTasks = ref<number[]>([])

const selectAll = computed({
  get: () => {
    return store.tasks.length > 0 && selectedTasks.value.length === store.tasks.length
  },
  set: (val) => {
    if (val) {
      selectedTasks.value = store.tasks.map((t: any) => t.id)
    } else {
      selectedTasks.value = []
    }
  }
})

onMounted(() => {
  store.fetchTasks()
  interval = setInterval(() => {
    store.fetchTasks()
  }, 2000)
})

onUnmounted(() => {
  clearInterval(interval)
})


const runSelected = async () => {
  for (const id of selectedTasks.value) {
    const task = store.tasks.find((t: any) => t.id === id)
    if (task && (task.status === 'Pending' || task.status === 'Paused' || task.status === 'Failed')) {
      await store.runTaskNow(id, 'live')
    }
  }
}

const removeSelected = async () => {
  for (const id of selectedTasks.value) {
    await store.removeTask(id)
  }
  selectedTasks.value = []
}
</script>
