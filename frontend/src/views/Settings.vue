<template>
  <div class="space-y-6 max-w-4xl mx-auto">
    <div class="bg-light-surface border border-light-border rounded-xl shadow-sm">
      <div class="p-6 border-b border-light-border">
        <h3 class="text-lg font-medium text-light-text mb-1">Cài đặt Ứng dụng</h3>
        <p class="text-sm text-light-muted mb-6">Quản lý cấu hình toàn cục của hệ thống</p>
        
        <div class="space-y-5" v-if="localSettings">
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
            <label class="text-sm font-medium text-light-text">Tên Ứng dụng</label>
            <div class="md:col-span-2">
              <input type="text" v-model="localSettings.appName" class="w-full bg-light-bg border border-light-border text-light-text text-sm rounded-lg focus:ring-primary focus:border-primary block p-2.5 transition-colors">
            </div>
          </div>
          
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
            <label class="text-sm font-medium text-light-text">Giao diện (Theme)</label>
            <div class="md:col-span-2">
              <select v-model="localSettings.theme" class="w-full bg-light-bg border border-light-border text-light-text text-sm rounded-lg focus:ring-primary focus:border-primary block p-2.5 transition-colors">
                <option value="light">Nền Sáng (Light Theme)</option>
                <option value="dark" disabled>Nền Tối (Sắp mở)</option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
            <label class="text-sm font-medium text-light-text">Bảo mật Cookie</label>
            <div class="md:col-span-2 flex flex-col space-y-3">
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="localSettings.maskCookieByDefault" class="sr-only peer">
                <div class="w-11 h-6 bg-light-border peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
                <span class="ml-3 text-sm font-medium text-light-muted">Luôn mã hóa ẩn Cookie khi nhập liệu</span>
              </label>

              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="localSettings.persistToJson" class="sr-only peer">
                <div class="w-11 h-6 bg-light-border peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
                <span class="ml-3 text-sm font-medium text-light-muted">Tự động sao lưu cấu hình vào file JSON nội bộ</span>
              </label>
            </div>
          </div>
        </div>
        <div v-else class="text-light-muted py-4">Đang tải cài đặt...</div>
      </div>
      
      <div class="p-6 bg-light-bg/50 rounded-b-xl flex justify-end">
        <button class="px-5 py-2.5 bg-primary hover:bg-primary-hover text-white shadow-sm rounded-lg text-sm font-medium transition-colors" @click="saveSettings">
          Lưu Cài đặt
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useMainStore } from '../stores/main'

const store = useMainStore()
const localSettings = ref<any>(null)

onMounted(async () => {
  await store.fetchSettings()
  if (store.settings) {
    localSettings.value = JSON.parse(JSON.stringify(store.settings))
  }
})

const saveSettings = async () => {
  try {
    if (localSettings.value) {
      await store.updateSettings(localSettings.value)
      alert('Đã lưu cấu hình thành công.')
      await store.fetchSettings()
    }
  } catch (e) {
    console.error(e)
  }
}
</script>
