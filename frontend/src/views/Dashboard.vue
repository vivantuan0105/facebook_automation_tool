<template>
  <div class="space-y-6">
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <StatCard 
        title="Total Posts" 
        :value="store.stats?.total_posts || 0" 
        :icon="DocumentTextIcon" 
        :trend="5"
      />
      <StatCard 
        title="Active Automations" 
        :value="store.stats?.total_tasks || 0" 
        :icon="CpuChipIcon"
        iconBgClass="bg-purple-500/10"
        iconColorClass="text-purple-500" 
      />
      <StatCard 
        title="Successful Runs" 
        :value="store.stats?.successful_runs || 0" 
        :icon="CheckCircleIcon" 
        iconBgClass="bg-green-500/10"
        iconColorClass="text-green-500" 
        :trend="12"
      />
      <StatCard 
        title="Recent Errors" 
        :value="store.stats?.recent_errors || 0" 
        :icon="ExclamationTriangleIcon" 
        iconBgClass="bg-red-500/10"
        iconColorClass="text-red-500" 
        :trend="-2"
      />
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="lg:col-span-2 bg-dark-surface rounded-xl border border-dark-border p-6 shadow-sm">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-medium text-white">Recent Activity</h3>
          <button class="text-primary text-sm hover:underline">View All</button>
        </div>
        <div class="space-y-4">
          <div v-if="store.logs.length === 0" class="text-dark-muted text-center py-8">No activity yet</div>
          <div v-for="log in store.logs.slice(0, 5)" :key="log.id" class="flex items-start bg-dark-bg p-3 rounded-lg border border-dark-border">
            <div class="w-8 h-8 rounded-full flex items-center justify-center mr-3 mt-0.5" 
                 :class="log.status === 'Success' ? 'bg-green-500/10 text-green-500' : 'bg-red-500/10 text-red-500'">
              <CheckCircleIcon v-if="log.status === 'Success'" class="w-5 h-5" />
              <ExclamationTriangleIcon v-else class="w-5 h-5" />
            </div>
            <div>
              <p class="text-sm font-medium text-white">{{ log.action }}</p>
              <p class="text-xs text-dark-muted mt-1">{{ new Date(log.time).toLocaleString() }} - {{ log.details }}</p>
            </div>
          </div>
        </div>
      </div>

      <div class="space-y-6">
        <div class="bg-dark-surface rounded-xl border border-dark-border p-6 shadow-sm">
          <h3 class="text-lg font-medium text-white mb-4">System Status</h3>
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <span class="text-sm text-dark-muted">Database</span>
              <StatusBadge status="success" />
            </div>
            <div class="flex items-center justify-between">
              <span class="text-sm text-dark-muted">Scheduler</span>
              <StatusBadge status="success" />
            </div>
            <div class="flex items-center justify-between">
              <span class="text-sm text-dark-muted">Worker Threads</span>
              <StatusBadge status="info" />
            </div>
          </div>
        </div>

        <div class="bg-dark-surface rounded-xl border border-dark-border p-6 shadow-sm">
          <h3 class="text-lg font-medium text-white mb-4">Account Connection</h3>
          <div v-if="store.stats?.connection_status === 'connected'" class="text-center p-4 border border-green-500/30 bg-green-500/5 rounded-lg">
            <div class="w-12 h-12 bg-green-500/20 rounded-full flex items-center justify-center mx-auto mb-3">
              <UserIcon class="w-6 h-6 text-green-400" />
            </div>
            <p class="text-sm font-medium text-white mb-1">Connected as</p>
            <p class="text-xs text-dark-muted">{{ store.stats?.connected_account }}</p>
          </div>
          <div v-else class="text-center p-4 border border-dark-border bg-dark-bg rounded-lg">
            <p class="text-sm text-dark-muted mb-3">No account connected</p>
            <button class="px-4 py-2 bg-primary hover:bg-primary-hover text-white rounded text-sm transition-colors">
              Connect Account
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useMainStore } from '../stores/main'
import StatCard from '../components/StatCard.vue'
import StatusBadge from '../components/StatusBadge.vue'
import { 
  DocumentTextIcon, 
  CpuChipIcon, 
  CheckCircleIcon, 
  ExclamationTriangleIcon,
  UserIcon
} from '@heroicons/vue/24/outline'

const store = useMainStore()

onMounted(async () => {
  await store.fetchDashboardStats()
  await store.fetchLogs()
})
</script>
