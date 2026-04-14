<template>
  <div class="max-w-7xl mx-auto grid grid-cols-1 xl:grid-cols-3 gap-6">
    <!-- Form Area -->
    <div class="xl:col-span-2">
      <div class="bg-light-surface border border-light-border shadow-sm rounded-xl p-6 h-full flex flex-col">
        <h2 class="text-lg font-medium text-light-text mb-4 border-b border-light-border pb-2">Thiết lập Tác vụ mới</h2>
        
        <div class="space-y-5">
          <!-- Chọn Tài Khoản (Cookie) -->
          <div>
            <div class="flex justify-between items-center mb-1">
              <label class="text-sm font-medium text-light-text">Chọn Tài khoản (Clone)</label>
            </div>
            <CustomSelect 
              v-model="form.cookie" 
              :options="accountOptions"
              placeholder="--- Chọn tài khoản để chạy ---"
            />
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
            <!-- Task Type -->
            <div>
              <label class="block text-sm font-medium text-light-text mb-1">Loại tác vụ</label>
              <CustomSelect 
                v-model="form.taskType" 
                :options="taskTypeOptions" 
                @change="onTaskTypeChange"
              />
            </div>

            <!-- Target Mode -->
            <div>
              <label class="block text-sm font-medium text-light-text mb-1">Xác định mục tiêu</label>
              <CustomSelect 
                v-model="form.targetMode" 
                :options="targetModeOptions"
              />
            </div>

            <!-- Target Input -->
            <div class="md:col-span-2" v-if="form.targetMode === 'post_url'">
              <label class="block text-sm font-medium text-light-text mb-1">URL Bài viết</label>
              <input type="text" v-model="form.postUrl" class="w-full text-sm bg-light-bg border border-light-border rounded-lg p-2.5 focus:ring-primary focus:border-primary" placeholder="Ví dụ: https://www.facebook.com/permalink.php?story_fbid=pfbid...&id=..." />
            </div>

            <!-- Reaction Type / Message -->
            <div class="md:col-span-2" v-if="form.taskType === 'Like bài viết'">
              <label class="block text-sm font-medium text-light-text mb-1">Loại cảm xúc</label>
              <CustomSelect 
                v-model="form.reactionType" 
                :options="reactionTypeOptions" 
                placeholder="Chọn cảm xúc"
              />
            </div>
            <div class="md:col-span-2" v-if="['Comment bài viết', 'Đăng bài viết'].includes(form.taskType)">
              <label class="block text-sm font-medium text-light-text mb-1">{{ form.taskType === 'Đăng bài viết' ? 'Nội dung bài viết (Status)' : 'Nội dung bình luận' }}</label>
              <textarea v-model="form.message" class="w-full text-sm bg-light-bg border border-light-border rounded-lg p-2.5 focus:ring-primary focus:border-primary" rows="4" placeholder="Nhập nội dung tương tác vào đây..."></textarea>
            </div>
            
            <div class="md:col-span-2" v-if="form.taskType === 'Đăng bài viết'">
              <label class="block text-sm font-medium text-light-text mb-1">Đường dẫn file ảnh (Tùy chọn)</label>
              <div class="space-y-2">
                <div v-for="(path, index) in form.photoPaths" :key="index" :data-path="path" class="flex gap-2">
                  <input type="text" v-model="form.photoPaths[index]" class="flex-1 w-full text-sm bg-light-bg border border-light-border rounded-lg p-2.5 focus:ring-primary focus:border-primary" placeholder="Ví dụ: C:\images\cat.jpg" />
                  <button @click="removePhoto(index)" class="px-3 py-2 border border-red-200 bg-red-50 hover:bg-red-100 text-red-600 text-sm rounded-lg transition-colors">Xóa</button>
                </div>
                <button @click="onSelectPhoto" class="px-4 py-2 border border-light-border bg-gray-50 hover:bg-gray-100 text-light-text text-sm rounded-lg transition-colors inline-block">+ Chọn thư mục/file ảnh</button>
            </div>
          </div>
        </div>
      </div>

        <div class="mt-6 flex flex-wrap justify-between items-center gap-3 pt-4 border-t border-light-border">
          <button @click="resetForm" class="px-4 py-2 border border-slate-300 bg-white hover:bg-slate-50 text-slate-700 text-sm font-medium rounded-lg shadow-sm flex items-center gap-1.5 transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            Làm mới
          </button>
          
          <div class="flex flex-wrap gap-3">
            <button @click="addToQueue" class="px-4 py-2 border border-primary text-primary hover:bg-primary/5 text-sm font-medium rounded-lg disabled:opacity-50 transition-colors" :disabled="!!validationError || isSubmitting">Thêm Hàng Đợi</button>
            <button @click="runNow" class="px-4 py-2 bg-primary text-white hover:bg-primary-hover text-sm font-medium rounded-lg disabled:opacity-50 shadow-sm transition-colors" :disabled="!!validationError || isSubmitting">Chạy Ngay</button>
          </div>
        </div>
        
        <div v-if="validationError" class="mt-4 p-3 bg-red-50 text-red-600 text-sm rounded border border-red-100">
          Lỗi: {{ validationError }}
        </div>
      </div>
    </div>

    <!-- Summary / Preview panel -->
    <div>
      <div class="bg-light-surface border border-light-border shadow-sm rounded-xl p-6 h-full flex flex-col">
        <h3 class="text-sm font-semibold text-light-muted uppercase tracking-wider mb-4">Tóm tắt thiết lập</h3>
        <dl class="space-y-3 text-sm">
          <div class="flex justify-between border-b border-light-border pb-2">
            <dt class="text-light-muted">Loại tác vụ</dt>
            <dd class="font-medium text-light-text">{{ form.taskType || '-' }}</dd>
          </div>
          <div class="flex justify-between border-b border-light-border pb-2">
            <dt class="text-light-muted">Mục tiêu</dt>
            <dd class="font-medium text-light-text truncate max-w-[150px]">{{ form.targetMode === 'post_url' ? (form.postUrl || '-') : (form.postId || '-') }}</dd>
          </div>
          <div class="flex justify-between border-b border-light-border pb-2">
            <dt class="text-light-muted">Cảm xúc</dt>
            <dd class="font-medium text-light-text">{{ form.reactionType || '-' }}</dd>
          </div>
          <div class="flex justify-between border-b border-light-border pb-2">
            <dt class="text-light-muted">Hành động</dt>
            <dd class="font-medium text-green-600">Thực thi thật (Live Server)</dd>
          </div>
          <div class="flex justify-between">
            <dt class="text-light-muted">Tài khoản</dt>
            <dd class="font-medium text-light-text truncate w-32 text-right">{{ form.cookie ? selectedAccountName : 'Chưa chọn' }}</dd>
          </div>
        </dl>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useMainStore } from '../stores/main'
import { GetAllAccounts, SelectPhotoDialog } from '../../wailsjs/go/app/App'
import CustomSelect from '../components/CustomSelect.vue'

const store = useMainStore()
const router = useRouter()

const accounts = ref<any[]>([])

const taskTypeOptions = [
  { value: 'Like bài viết', label: 'Thích (Reaction) bài viết' },
  { value: 'Comment bài viết', label: 'Bình luận bài viết' },
  { value: 'Đăng bài viết', label: 'Đăng bài viết (Post)' },
  { value: 'Reaction batch', label: 'Reaction hàng loạt (Sắp có)', disabled: true }
]

const targetModeOptions = computed(() => {
  if (form.taskType === 'Đăng bài viết') {
    return [{ value: 'timeline', label: 'Đăng lên trang cá nhân (Timeline)' }]
  }
  return [{ value: 'post_url', label: 'Dùng URL Bài viết (Khuyên dùng)' }]
})

const reactionTypeOptions = ['Like', 'Love', 'Care', 'Haha', 'Wow', 'Sad', 'Angry'].map(r => ({ value: r, label: r }))

const accountOptions = computed(() => {
  return accounts.value.map(acc => ({
    value: acc.cookie,
    label: `${acc.name || acc.uid} (${acc.status})`
  }))
})

const selectedAccountName = computed(() => {
  const acc = accounts.value.find(a => a.cookie === form.cookie)
  return acc ? (acc.name || acc.uid) : '...'
})

onMounted(async () => {
  try {
    const data = await GetAllAccounts()
    accounts.value = data || []
  } catch(e) {
    console.error(e)
  }
})

const form = reactive({
  cookie: '',
  taskType: 'Like bài viết',
  execMode: 'live',
  targetMode: 'post_url',
  postUrl: '',
  postId: '',
  reactionType: 'Like',
  message: '',
  photoPaths: [] as string[],
})



const validationError = ref('')
const validationSuccess = ref(false)

watch(form, () => {
  validationError.value = ''
  validationSuccess.value = false
}, { deep: true })

const onTaskTypeChange = () => {
  if (form.taskType === 'Đăng bài viết') {
    form.targetMode = 'timeline'
    form.postUrl = ''
  } else {
    form.targetMode = 'post_url'
  }
}

const onSelectPhoto = async () => {
  try {
    const path = await SelectPhotoDialog()
    if (path) {
      if(!form.photoPaths) form.photoPaths = []
      form.photoPaths.push(path)
    }
  } catch(e) {
    console.error("Lỗi chọn file:", e)
  }
}

const removePhoto = (index: number) => {
  form.photoPaths.splice(index, 1)
}

const isSubmitting = ref(false)

const validateForm = async () => {
  validationError.value = ''
  const err = await store.validateTask(form)
  if (err) {
    validationError.value = err
  }
}

const addToQueue = async () => {
  if (isSubmitting.value) return
  isSubmitting.value = true
  await validateForm()
  if (!validationError.value) {
    await store.createTask(form)
    router.push('/queue')
  }
  isSubmitting.value = false
}

const runNow = async () => {
  if (isSubmitting.value) return
  isSubmitting.value = true
  await validateForm()
  if (!validationError.value) {
    await store.createTask(form)
    const newestTask = store.tasks[0]
    await store.runTaskNow(newestTask.id, form.execMode)
    router.push('/queue')
  }
  isSubmitting.value = false
}

const resetForm = () => {
  form.cookie = ''
  form.postUrl = ''
  form.postId = ''
  form.reactionType = 'Like'
  form.photoPaths = []

  validationError.value = ''
  validationSuccess.value = false
}
</script>


