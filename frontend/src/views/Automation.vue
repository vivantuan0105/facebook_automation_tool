<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h2 class="text-xl font-semibold text-white">Automation Tasks</h2>
      <button class="px-4 py-2 bg-primary hover:bg-primary-hover text-white rounded-lg text-sm font-medium transition-colors flex items-center shadow-sm">
        <PlusIcon class="w-4 h-4 mr-2" />
        Create Task
      </button>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-if="store.tasks.length === 0" class="col-span-full text-center text-dark-muted py-12">
        No automation tasks configured.
      </div>
      <div v-for="task in store.tasks" :key="task.id" 
           class="bg-dark-surface border border-dark-border rounded-xl p-5 hover:border-primary/30 transition-colors shadow-sm group">
        <div class="flex justify-between items-start mb-4">
          <div class="flex items-center">
            <div class="w-10 h-10 rounded-lg bg-dark-bg border border-dark-border flex items-center justify-center mr-3 group-hover:border-primary/50 transition-colors">
              <CommandLineIcon class="w-5 h-5 text-dark-muted group-hover:text-primary transition-colors" />
            </div>
            <div>
              <h3 class="text-white font-medium">{{ task.name }}</h3>
              <span class="text-xs text-primary">{{ task.type }}</span>
            </div>
          </div>
          <StatusBadge :status="task.enabled ? 'Success' : 'Draft'" />
        </div>
        <p class="text-sm text-dark-muted mb-6 h-10">{{ task.description }}</p>
        <div class="text-xs text-dark-muted flex justify-between items-center border-t border-dark-border pt-4">
          <span>Schedule: <strong class="text-white font-medium">{{ task.schedule }}</strong></span>
          <button class="px-3 py-1 bg-dark-bg hover:bg-dark-border border border-dark-border text-white rounded transition-colors" @click="runTask(task.id)">
            Run Now
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useMainStore } from '../stores/main'
import StatusBadge from '../components/StatusBadge.vue'
import { PlusIcon, CommandLineIcon } from '@heroicons/vue/24/outline'

const store = useMainStore()

onMounted(async () => {
  await store.fetchTasks()
})

const runTask = async (id: number) => {
  try {
    if ((window as any).go?.app?.App?.RunTaskOnce) {
      await (window as any).go.app.App.RunTaskOnce(id)
      alert('Task triggered successfully.')
    } else {
      alert(`Mock trigger task ID: ${id}`)
    }
  } catch(e) {
    console.error(e)
  }
}
</script>
