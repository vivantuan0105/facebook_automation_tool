<template>
  <div class="space-y-6">
    <!-- Header Controls -->
    <div class="flex justify-between items-center bg-white p-4 rounded-xl shadow-sm border border-gray-100">
      <div>
        <h2 class="text-lg font-bold text-gray-800">Quản lý Tài Khoản (Clone)</h2>
        <p class="text-sm text-gray-500 mt-1">Danh sách các tài khoản Facebook đang được biên chế để thả Like.</p>
      </div>
      <button 
        @click="showAddModal = true"
        class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium shadow-sm transition-colors flex items-center gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clip-rule="evenodd" />
        </svg>
        Thêm tài khoản
      </button>
    </div>

    <!-- Account List -->
    <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
      <div v-if="loading" class="p-8 text-center text-gray-500">
        <div class="animate-spin inline-block w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full mb-2"></div>
        <p>Đang tải danh sách tài khoản...</p>
      </div>
      
      <div v-else-if="accounts.length === 0" class="p-12 text-center text-gray-500">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto text-gray-300 mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
        </svg>
        <p class="text-lg font-medium text-gray-700">Chưa có tài khoản nào</p>
        <p class="text-sm mt-1">Vui lòng bấm Thêm tài khoản để đưa nick vào hệ thống.</p>
      </div>

      <div v-else class="overflow-x-auto w-full">
        <table class="w-full text-left text-sm whitespace-nowrap">
          <thead class="bg-gray-50 text-gray-600 font-medium border-b border-gray-200">
            <tr>
              <th class="py-3 px-4 w-12 border-r border-gray-200/60">#</th>
              <th class="py-3 px-4 border-r border-gray-200/60">Tên gợi nhớ</th>
              <th class="py-3 px-4 border-r border-gray-200/60">UID (c_user)</th>
              <th class="py-3 px-4 border-r border-gray-200/60">Ngày thêm</th>
              <th class="py-3 px-4 border-r border-gray-200/60">Trạng thái</th>
              <th class="py-3 px-4 text-right">Thao tác</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-for="(acc, index) in accounts" :key="acc.uid" class="hover:bg-gray-50 transition-colors">
              <td class="py-3 px-4 text-gray-500 border-r border-gray-100">{{ index + 1 }}</td>
              <td class="py-3 px-4 font-medium text-gray-800 border-r border-gray-100">
                {{ acc.name || 'Không xác định' }}
              </td>
              <td class="py-3 px-4 border-r border-gray-100">
                <span class="bg-gray-100 text-gray-700 px-2 py-1 rounded font-mono text-xs border border-gray-200">
                  {{ acc.uid }}
                </span>
              </td>
              <td class="py-3 px-4 text-gray-500 border-r border-gray-100">{{ acc.createdAt }}</td>
              <td class="py-3 px-4 border-r border-gray-100">
                <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium"
                  :class="acc.status === 'Live' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'">
                  <span class="w-1.5 h-1.5 rounded-full" :class="acc.status === 'Live' ? 'bg-green-500' : 'bg-red-500'"></span>
                  {{ acc.status }}
                </span>
              </td>
              <td class="py-3 px-4 text-right flex justify-end gap-2">
                <button @click="deleteAccount(acc.uid)" class="text-red-500 hover:text-red-700 p-1 rounded hover:bg-red-50 transition-colors" title="Xóa tài khoản">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Add Account Modal -->
    <div v-if="showAddModal" class="fixed inset-0 bg-gray-900/10 backdrop-blur-sm z-50 pointer-events-none">
      <div 
        ref="addModalRef"
        class="bg-white rounded-xl shadow-2xl w-full max-w-lg overflow-hidden absolute pointer-events-auto shadow-[0_0_20px_rgba(0,0,0,0.15)]"
        :style="{ top: addModalY + 'px', left: addModalX + 'px' }"
      >
        <div 
          class="px-6 py-4 border-b border-gray-100 flex justify-between items-center bg-blue-50 cursor-move"
          @mousedown="startDragAddModal"
        >
          <h3 class="text-lg font-bold text-gray-800">Thêm Tài Khoản (Cookie)</h3>
          <button @click="showAddModal = false" class="text-gray-400 hover:text-gray-600 transition-colors">
             <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <div class="p-6 space-y-4">
          <div v-if="addError" class="bg-red-50 border border-red-200 text-red-600 text-sm p-3 rounded-lg flex items-start gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 shrink-0 mt-0.5" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
            </svg>
            {{ addError }}
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Tên gợi nhớ</label>
            <input v-model="newAccName" type="text" placeholder="Ví dụ: Clone Trà Đá 01" 
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all">
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Chuỗi Cookie Facebook</label>
            <textarea v-model="newAccCookie" rows="4" placeholder="c_user=1000...; xs=...;" 
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all font-mono text-xs"></textarea>
            <p class="text-xs text-gray-500 mt-1">Hệ thống sẽ tự động bóc tách c_user (UID) làm định danh.</p>
          </div>
        </div>

        <div class="px-6 py-4 bg-gray-50 border-t border-gray-100 flex justify-end gap-3">
          <button @click="showAddModal = false" class="px-4 py-2 text-gray-600 font-medium hover:bg-gray-200 rounded-lg transition-colors">Hủy</button>
          <button @click="submitAdd" :disabled="isSubmitting" 
            class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg shadow-sm transition-colors flex items-center gap-2 disabled:opacity-70">
            <svg v-if="isSubmitting" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Lưu tài khoản
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetAllAccounts, AddAccount, DeleteAccount } from '../../wailsjs/go/app/App'

const accounts = ref<any[]>([])
const loading = ref(true)

const showAddModal = ref(false)
const isSubmitting = ref(false)
const newAccName = ref('')
const newAccCookie = ref('')
const addError = ref('')

// Drag Add Modal Logic
const addModalX = ref(100)
const addModalY = ref(50)
let isDragging = false
let startX = 0
let startY = 0
let initialX = 0
let initialY = 0

function startDragAddModal(e: MouseEvent) {
  isDragging = true
  startX = e.clientX
  startY = e.clientY
  initialX = addModalX.value
  initialY = addModalY.value
  document.addEventListener('mousemove', onDragAddModal)
  document.addEventListener('mouseup', stopDragAddModal)
}

function onDragAddModal(e: MouseEvent) {
  if (!isDragging) return
  const dx = e.clientX - startX
  const dy = e.clientY - startY
  addModalX.value = initialX + dx
  addModalY.value = initialY + dy
}

function stopDragAddModal() {
  isDragging = false
  document.removeEventListener('mousemove', onDragAddModal)
  document.removeEventListener('mouseup', stopDragAddModal)
}

async function fetchAccounts() {
  loading.value = true
  try {
    const data = await GetAllAccounts()
    accounts.value = data || []
  } catch (e: any) {
    console.error("Failed to load accounts:", e)
  } finally {
    loading.value = false
  }
}

async function submitAdd() {
  if (!newAccCookie.value) {
    addError.value = "Vui lòng nhập Cookie"
    return
  }
  
  isSubmitting.value = true
  addError.value = ""
  try {
    await AddAccount(newAccName.value, newAccCookie.value)
    await fetchAccounts()
    showAddModal.value = false
    newAccName.value = ''
    newAccCookie.value = ''
  } catch (err: any) {
    addError.value = err.toString()
  } finally {
    isSubmitting.value = false
  }
}

async function deleteAccount(uid: string) {
  if (confirm(`Bạn có chắc muốn xóa tài khoản có UID ${uid}? Mọi dữ liệu sẽ bị xóa.`)) {
    try {
      await DeleteAccount(uid)
      await fetchAccounts()
    } catch (err) {
      alert("Lỗi khi xóa: " + err)
    }
  }
}

onMounted(() => {
  fetchAccounts()
})
</script>

<style scoped>
.animate-fade-in-up {
  animation: fadeInUp 0.3s ease-out;
}
@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
