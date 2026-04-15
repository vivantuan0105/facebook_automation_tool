const fs = require('fs');

function transformAccounts() {
  const filePath = 'frontend/src/views/Accounts.vue';
  let content = fs.readFileSync(filePath, 'utf8');

  // Wrap inside single div for transition compatibility
  if (!content.includes('<template>\n  <div>\n    <div class="space-y-6">') && !content.includes('<template>\r\n  <div>\r\n    <div class="space-y-6">')) {
    content = content.replace(/<template>\r?\n\s*<div class="space-y-6">/, '<template>\n  <div>\n    <div class="space-y-6">');
    content = content.replace(/<\/div>\r?\n<\/template>/, '  </div>\n  </div>\n</template>');
  }
  
  // Add Login Button to Template
  if (!content.includes('showLoginModal = true')) {
    content = content.replace(
      /<button\s*@click="showAddModal = true"/,
      '<button\n          @click="showLoginModal = true"\n          class="px-4 py-2 rounded-lg border border-indigo-200 bg-white text-indigo-600 hover:bg-indigo-50 font-medium shadow-sm transition-colors flex items-center gap-2 whitespace-nowrap"\n        >\n          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">\n            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-7.5a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 006 21h7.5a2.25 2.25 0 002.25-2.25V15" />\n            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 12H9m0 0l3-3m-3 3l3 3" />\n          </svg>\n          Login\n        </button>\n        <button \n          @click="showAddModal = true"'
    );
  }

  // Add the Login Modal HTML right before closing template
  if (!content.includes('<!-- Login Modal -->')) {
    const modalHTML = `
    <!-- Login Modal -->
    <div v-if="showLoginModal" @mousedown.self="showLoginModal = false" class="fixed inset-0 bg-slate-800/10 z-[100] pointer-events-auto transition-opacity">
      <div
        class="bg-white rounded-xl shadow-2xl w-full max-w-lg overflow-hidden absolute pointer-events-auto shadow-[0_0_20px_rgba(0,0,0,0.15)]"
        :style="{ top: loginModalY + 'px', left: loginModalX + 'px' }"
      >
        <div class="px-6 py-4 border-b border-gray-100 flex justify-between items-center bg-indigo-50 cursor-move" @mousedown="startDragLoginModal">
          <h3 class="text-lg font-bold text-gray-800">Login Facebook</h3>
          <button @click="showLoginModal = false" class="text-gray-400 hover:text-gray-600 transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div class="p-6 space-y-4">
          <div class="rounded-xl border border-slate-200 bg-slate-50 px-4 py-4">
            <div class="text-sm font-semibold text-slate-800">Luồng login phù hợp với app này</div>
            <ol class="mt-3 space-y-2 text-sm text-slate-600 list-decimal list-inside">
              <li>Mở trang login Facebook trong trình duyệt.</li>
              <li>Đăng nhập trên trình duyệt bằng tài khoản của bạn.</li>
              <li>Quay lại \`Thêm tài khoản\` và dán cookie/session vào hệ thống.</li>
            </ol>
          </div>

          <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
            <button
              @click="openFacebookLogin"
              class="inline-flex items-center justify-center gap-2 rounded-lg bg-indigo-600 px-4 py-2.5 text-sm font-medium text-white shadow-sm transition-colors hover:bg-indigo-700"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.5 6H17m0 0v3.5m0-3.5L10 13m-4 7h12a2 2 0 002-2V10a2 2 0 00-2-2h-3m-5 0H6a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
              Mở Facebook
            </button>
            <button
              @click="openAddAccountFromLogin"
              class="inline-flex items-center justify-center gap-2 rounded-lg border border-slate-200 bg-white px-4 py-2.5 text-sm font-medium text-slate-700 shadow-sm transition-colors hover:bg-slate-50"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clip-rule="evenodd" />
              </svg>
              Sang Thêm tài khoản
            </button>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Email / Số điện thoại / UID</label>
            <input v-model="loginIdentifier" type="text" placeholder="Nhập thông tin đăng nhập" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all">
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Mật khẩu</label>
            <input v-model="loginPassword" type="password" placeholder="Nhập mật khẩu" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all">
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Ghi chú</label>
            <input v-model="loginNote" type="text" placeholder="Ví dụ: Nick login mới" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all">
          </div>
        </div>

        <div class="px-6 py-4 bg-gray-50 border-t border-gray-100 flex justify-end gap-3">
          <button @click="showLoginModal = false" class="px-4 py-2 text-gray-600 font-medium hover:bg-gray-200 rounded-lg transition-colors">Đóng</button>
          <button @click="submitLogin" class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg shadow-sm transition-colors flex items-center gap-2" :disabled="isSubmittingLogin">
            <svg v-if="isSubmittingLogin" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Đăng nhập bằng Request
          </button>
        </div>
      </div>
    </div>
`;
    content = content.replace(/<\/template>/, modalHTML + '\n</template>');
  }

  // Update imports
  if (!content.includes('LoginAccount')) {
    content = content.replace(
      /import \{ GetAllAccounts.*\} from '\.\.\/\.\.\/wailsjs\/go\/app\/App'/,
      'import { GetAllAccounts, AddAccount, DeleteAccount, SelectPhotoDialog, ScanAccountData, GetAccountFriendsList, LoginAccount } from \'../../wailsjs/go/app/App\''
    );
  }
  if (!content.includes('BrowserOpenURL')) {
    content = content.replace(
      "import { useMainStore } from '../stores/main'",
      "import { BrowserOpenURL } from '../../wailsjs/runtime/runtime'\nimport { useMainStore } from '../stores/main'"
    );
  }

  // Add login logic state
  if (!content.includes('const loginIdentifier = ref(\'\')')) {
    const loginScript = `
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
  loginModalX.value = initialXLogin + (e.clientX - startXLogin)
  loginModalY.value = initialYLogin + (e.clientY - startYLogin)
}
function stopDragLoginModal() {
  isDraggingLogin = false
  document.removeEventListener('mousemove', onDragLoginModal)
  document.removeEventListener('mouseup', stopDragLoginModal)
}

function openFacebookLogin() {
  BrowserOpenURL('https://www.facebook.com/login')
}
function openAddAccountFromLogin() {
  showLoginModal.value = false
  showAddModal.value = true
}

async function submitLogin() {
  if (!loginIdentifier.value || !loginPassword.value) {
    alert("Vui lòng nhập định danh và mật khẩu!")
    return
  }
  
  isSubmittingLogin.value = true
  try {
    await LoginAccount(loginIdentifier.value, loginPassword.value, loginNote.value)
    alert("Đăng nhập thành công và đã thêm vào danh sách!")
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
`;
    content = content.replace('// Accounts Management', loginScript + '\n// Accounts Management');
  }

  fs.writeFileSync(filePath, content);
}

transformAccounts();
