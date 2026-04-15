const fs = require('fs');

const path = 'frontend/src/views/Accounts.vue';
let content = fs.readFileSync(path, 'utf8');

// 1. Root wrapper
content = content.replace(/<template>\r?\n\s*<div class="space-y-6">/, '<template>\n  <div>\n    <div class="space-y-6">');
content = content.replace(/<\/div>\r?\n<\/template>/, '  </div>\n  </div>\n</template>');

// 2. Head button replacement
const buttonReplacement = `
      <div class="flex items-center gap-3">
        <button
          @click="showLoginModal = true"
          class="px-4 py-2 rounded-lg border border-indigo-200 bg-white text-indigo-600 hover:bg-indigo-50 font-medium shadow-sm transition-colors flex items-center gap-2 whitespace-nowrap"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-7.5a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 006 21h7.5a2.25 2.25 0 002.25-2.25V15" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 12H9m0 0l3-3m-3 3l3 3" />
          </svg>
          Login Facebook
        </button>

        <button 
          @click="showAddModal = true"
          class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium shadow-sm transition-colors flex items-center gap-2 whitespace-nowrap">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clip-rule="evenodd" />
          </svg>
          Thêm tài khoản
        </button>
      </div>
    </div>
`.trim();

content = content.replace(/<button\s*@click="showAddModal = true"\s*class="bg-blue-600[\s\S]*?Thêm tài khoản\s*<\/button>\r?\n\s*<\/div>/, buttonReplacement);

// 3. Login Modal String (Without Luồng Login Info and Open Buttons)
const loginModal = `
    <!-- Login Modal -->
    <div v-if="showLoginModal" @mousedown.self="showLoginModal = false" class="fixed inset-0 bg-slate-800/10 z-[100] pointer-events-auto transition-opacity">
      <div
        class="bg-white rounded-xl shadow-2xl w-full max-w-sm overflow-hidden absolute pointer-events-auto shadow-[0_0_20px_rgba(0,0,0,0.15)]"
        :style="{ top: loginModalY + 'px', left: loginModalX + 'px' }"
      >
        <div class="px-6 py-4 border-b border-gray-100 flex justify-between items-center bg-indigo-50 cursor-move" @mousedown="startDragLoginModal">
          <h3 class="text-lg font-bold text-gray-800">Login Facebook (Request)</h3>
          <button @click="showLoginModal = false" class="text-gray-400 hover:text-gray-600 transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <div class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Tài khoản (UID/Email)</label>
            <input v-model="loginIdentifier" type="text" placeholder="Nhập tài khoản" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all">
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Mật khẩu</label>
            <input v-model="loginPassword" type="password" placeholder="Nhập mật khẩu" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all" @keyup.enter="submitLogin">
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Ghi chú (Tùy chọn)</label>
            <input v-model="loginNote" type="text" placeholder="Ví dụ: Nick login mới" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all">
          </div>
        </div>

        <div class="px-6 py-4 bg-gray-50 border-t border-gray-100 flex justify-end gap-3">
          <button @click="showLoginModal = false" class="px-4 py-2 text-gray-600 font-medium hover:bg-gray-200 rounded-lg transition-colors">Hủy</button>
          <button @click="submitLogin" :disabled="isSubmittingLogin" class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg shadow-sm transition-colors flex items-center gap-2 disabled:opacity-70">
            <svg v-if="isSubmittingLogin" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Đăng nhập ngay
          </button>
        </div>
      </div>
    </div>
</template>
`;

content = content.replace(/<\/template>/, loginModal);

// 4. Update imports
content = content.replace(
  /import { GetAllAccounts, AddAccount, DeleteAccount, SelectPhotoDialog, ScanAccountData, GetAccountFriendsList } from '\.\.\/\.\.\/wailsjs\/go\/app\/App'/,
  "import { GetAllAccounts, AddAccount, DeleteAccount, SelectPhotoDialog, ScanAccountData, GetAccountFriendsList, LoginAccount } from '../../wailsjs/go/app/App'"
);

// 5. Appending logic state
const logicState = `
// Login Management
const showLoginModal = ref(false)
const loginIdentifier = ref('')
const loginPassword = ref('')
const loginNote = ref('')
const isSubmittingLogin = ref(false)

const loginModalX = ref(180)
const loginModalY = ref(80)
let isDraggingLogin = false
let startXLogin = 0
let startYLogin = 0
let initialXLogin = 0
let initialYLogin = 0

function startDragLoginModal(e: MouseEvent) {
  isDraggingLogin = true
  startXLogin = e.clientX
  startYLogin = e.clientY
  initialXLogin = loginModalX.value
  initialYLogin = loginModalY.value
  document.addEventListener('mousemove', onDragLoginModal)
  document.addEventListener('mouseup', stopDragLoginModal)
}

function onDragLoginModal(e: MouseEvent) {
  if (!isDraggingLogin) return
  const dx = e.clientX - startXLogin
  const dy = e.clientY - startYLogin
  loginModalX.value = initialXLogin + dx
  loginModalY.value = initialYLogin + dy
}

function stopDragLoginModal() {
  isDraggingLogin = false
  document.removeEventListener('mousemove', onDragLoginModal)
  document.removeEventListener('mouseup', stopDragLoginModal)
}

async function submitLogin() {
  if (!loginIdentifier.value || !loginPassword.value) {
    alert("Vui lòng nhập định danh và mật khẩu!")
    return
  }
  
  isSubmittingLogin.value = true
  try {
    await LoginAccount(loginIdentifier.value, loginPassword.value, loginNote.value)
    alert("Đăng nhập thành công và đã tự động lưu Cookie!")
    showLoginModal.value = false
    loginIdentifier.value = ''
    loginPassword.value = ''
    loginNote.value = ''
    await fetchAccounts()
  } catch (err: any) {
    alert("Lỗi đăng nhập: " + err.toString())
  } finally {
    isSubmittingLogin.value = false
  }
}

// Accounts Management`;

content = content.replace(/\/\/\s*Accounts Management/, logicState);

fs.writeFileSync(path, content);
