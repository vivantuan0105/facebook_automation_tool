<template>
  <div class="space-y-6 max-w-6xl mx-auto">
    <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-4">
      <StatCard title="Tổng số tác vụ" :value="store.overview?.totalTasks || 0" iconBgClass="bg-blue-500/10" iconColorClass="text-blue-500" />
      <StatCard title="Đang chờ" :value="store.overview?.pending || 0" iconBgClass="bg-yellow-500/10" iconColorClass="text-yellow-500" />
      <StatCard title="Đang chạy" :value="store.overview?.running || 0" iconBgClass="bg-purple-500/10" iconColorClass="text-purple-500" />
      <StatCard title="Thành công" :value="store.overview?.success || 0" iconBgClass="bg-green-500/10" iconColorClass="text-green-500" />
      <StatCard title="Thất bại" :value="store.overview?.failed || 0" iconBgClass="bg-red-500/10" iconColorClass="text-red-500" />
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="lg:col-span-2 space-y-6">
        <div class="bg-light-surface rounded-xl border border-light-border p-6 shadow-sm">
          <h3 class="text-lg font-medium text-light-text mb-4">Danh sách hàng đợi (Preview)</h3>
          <div v-if="store.tasks.length === 0" class="text-light-muted text-center py-4">Chưa có tác vụ nào trong hàng đợi</div>
          <div class="overflow-x-auto text-sm w-full" v-else>
            <table class="w-full text-left whitespace-nowrap">
              <thead class="bg-indigo-50/60 text-slate-700 font-semibold border-b-2 border-indigo-100/50 text-xs uppercase">
                <tr>
                  <th class="py-2.5 px-3 border-r border-indigo-100/30 text-center">Loại tác vụ</th>
                  <th class="py-2.5 px-3 border-r border-indigo-100/30 text-center">Mục tiêu</th>
                  <th class="py-2.5 px-3 text-center">Trạng thái</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="t in store.tasks.slice(0, 5)" :key="t.id" class="border-b border-slate-100 hover:bg-slate-50 transition-colors">
                  <td class="py-1.5 px-3 text-slate-700 border-r border-slate-100 text-left">{{ t.taskType }}</td>
                  <td class="py-1.5 px-3 text-slate-500 truncate max-w-[150px] border-r border-slate-100 text-left">{{ t.postUrl || t.postId }}</td>
                  <td class="py-1.5 px-3 text-center">
                    <span class="inline-block px-1.5 py-0.5 rounded text-[11px] font-medium" :class="{
                      'bg-amber-50 text-amber-600': t.status==='Pending',
                      'bg-emerald-50 text-emerald-600': t.status==='Success',
                      'bg-indigo-50 text-indigo-600': t.status==='Running' || t.status==='WaitingConfirmation',
                      'bg-rose-50 text-rose-600': t.status==='Failed'
                    }">{{ t.status }}</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div class="bg-light-surface rounded-xl border border-light-border p-6 shadow-sm">
          <h3 class="text-lg font-medium text-light-text mb-4">Hướng dẫn hoạt động</h3>
          <ol class="list-decimal list-inside text-sm text-light-muted space-y-2">
            <li>Mở <router-link class="text-primary hover:underline" to="/task-center">Trung tâm tác vụ</router-link>.</li>
            <li>Nhập chuỗi Cookie phiên làm việc (giữ cục bộ, ẩn tự động).</li>
            <li>Chọn phương thức xác định mục tiêu (Post URL hoặc Post ID) và điền vào.</li>
            <li>Định cấu hình Advanced Context (Relay/GraphQL vars) nếu dùng Mock chuyên sâu.</li>
            <li>Chọn Chế độ thực thi (Giả lập để test an toàn, hoặc Thủ công để được hướng dẫn qua trình duyệt ngoài).</li>
            <li>Xác thực và đưa vào Hàng đợi hoặc Chạy ngay.</li>
          </ol>
        </div>
      </div>

      <div class="bg-light-surface rounded-xl border border-light-border p-6 shadow-sm flex flex-col h-full">
        <h3 class="text-lg font-medium text-light-text mb-4">Hoạt động gần đây</h3>
        <div class="space-y-4 flex-1 overflow-y-auto pr-2">
          <div v-if="store.logs.length === 0" class="text-light-muted text-center py-4">Chưa có hoạt động</div>
          <div v-for="log in store.logs.slice(0, 6)" :key="log.id" class="text-sm pb-3 border-b border-light-border last:border-0">
            <p class="font-medium text-light-text">{{ log.action }}</p>
            <p class="text-xs text-light-muted">{{ log.details }}</p>
            <p class="text-[10px] text-light-muted/70 mt-1">{{ new Date(log.time).toLocaleString() }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useMainStore } from '../stores/main'
import StatCard from '../components/StatCard.vue'

const store = useMainStore()
let interval: any;

const refreshData = async () => {
  await store.fetchOverview()
  await store.fetchTasks()
  await store.fetchLogs()
}

onMounted(() => {
  refreshData()
  interval = setInterval(refreshData, 3000)
})

onUnmounted(() => {
  clearInterval(interval)
})
</script>
