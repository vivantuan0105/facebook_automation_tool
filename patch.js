const fs = require('fs');
let content = fs.readFileSync('frontend/src/views/Accounts.vue', 'utf8');

content = content.replace('<template>\n  <div class="space-y-6">', '<template>\n  <div>\n    <div class="space-y-6">');
content = content.replace('  </div>\n</template>', '  </div>\n\n    <!-- Login Modal -->\n    <div v-if="showLoginModal" @mousedown.self="showLoginModal = false" class="fixed inset-0 bg-slate-800/10 z-[100] pointer-events-auto transition-opacity">\n      <div\n        class="bg-white rounded-xl shadow-2xl w-full max-w-lg overflow-hidden absolute pointer-events-auto shadow-[0_0_20px_rgba(0,0,0,0.15)]"\n        :style="{ top: loginModalY + \'px\', left: loginModalX + \'px\' }"\n      >\n        <div class="px-6 py-4 border-b border-gray-100 flex justify-between items-center bg-indigo-50 cursor-move" @mousedown="startDragLoginModal">\n          <h3 class="text-lg font-bold text-gray-800">Login Facebook</h3>\n          <button @click="showLoginModal = false" class="text-gray-400 hover:text-gray-600 transition-colors">\n            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">\n              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />\n            </svg>\n          </button>\n        </div>\n\n        <div class="p-6 space-y-4">\n          <div class="rounded-xl border border-slate-200 bg-slate-50 px-4 py-4">\n            <div class="text-sm font-semibold text-slate-800">Luồng login trực tiếp</div>\n            <ol class="mt-3 space-y-2 text-sm text-slate-600 list-decimal list-inside">\n              <li>Mở trang login Facebook trong trình duyệt.</li>\n              <li>Nếu có xác thực, hãy xác minh trên Cửa sổ bật lên.</li>\n            </ol>\n          </div>\n\n          <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">\n            <button\n              @click="openFacebookLogin"\n              class="inline-flex items-center justify-center gap-2 rounded-lg bg-indigo-600 px-4 py-2.5 text-sm font-medium text-white shadow-sm transition-colors hover:bg-indigo-700"\n            >\n              Mở Facebook\n            </button>\n            <button\n              @click="openAddAccountFromLogin"\n              class="inline-flex items-center justify-center gap-2 rounded-lg border border-slate-200 bg-white px-4 py-2.5 text-sm font-medium text-slate-700 shadow-sm transition-colors hover:bg-slate-50"\n            >\n              Sang Thêm tài khoản\n            </button>\n          </div>\n\n          <div>\n            <label class="block text-sm font-medium text-gray-700 mb-1">Email / Số điện thoại / UID</label>\n            <input v-model="loginIdentifier" type="text" placeholder="Nhập thông tin" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all">\n          </div>\n          <div>\n            <label class="block text-sm font-medium text-gray-700 mb-1">Mật khẩu</label>\n            <input v-model="loginPassword" type="password" placeholder="Nhập mật khẩu" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all">\n          </div>\n          <div>\n            <label class="block text-sm font-medium text-gray-700 mb-1">Ghi chú</label>\n            <input v-model="loginNote" type="text" placeholder="Ví dụ: Nick login mới" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all">\n          </div>\n        </div>\n\n        <div class="px-6 py-4 bg-gray-50 border-t border-gray-100 flex justify-end gap-3">\n          <button @click="showLoginModal = false" class="px-4 py-2 text-gray-600 font-medium hover:bg-gray-200 rounded-lg transition-colors">Đóng</button>\n          <button @click="submitLogin" class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg shadow-sm transition-colors flex items-center gap-2" :disabled="isSubmittingLogin">\n            Đăng nhập bằng Request\n          </button>\n        </div>\n      </div>\n    </div>\n    \n  </div>\n</template>');

content = content.replace("import { GetAllAccounts, AddAccount, DeleteAccount, SelectPhotoDialog, ScanAccountData, GetAccountFriendsList } from '../../wailsjs/go/app/App'", "import { GetAllAccounts, AddAccount, DeleteAccount, SelectPhotoDialog, ScanAccountData, GetAccountFriendsList, LoginAccount } from '../../wailsjs/go/app/App'\nimport { BrowserOpenURL } from '../../wailsjs/runtime/runtime'");

let loginState = `
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
  isDraggingLogin = true; startXLogin = e.clientX; startYLogin = e.clientY;
  initialXLogin = loginModalX.value; initialYLogin = loginModalY.value;
  document.addEventListener('mousemove', onDragLoginModal)
  document.addEventListener('mouseup', stopDragLoginModal)
}
function onDragLoginModal(e: MouseEvent) {
  if (!isDraggingLogin) return;
  loginModalX.value = initialXLogin + (e.clientX - startXLogin)
  loginModalY.value = initialYLogin + (e.clientY - startYLogin)
}
function stopDragLoginModal() {
  isDraggingLogin = false;
  document.removeEventListener('mousemove', onDragLoginModal); document.removeEventListener('mouseup', stopDragLoginModal)
}

function openFacebookLogin() { BrowserOpenURL('https://www.facebook.com/login') }
function openAddAccountFromLogin() { showLoginModal.value = false; showAddModal.value = true }

async function submitLogin() {
  if (!loginIdentifier.value || !loginPassword.value) { alert("Vui lòng nhập định danh và mật khẩu!"); return; }
  isSubmittingLogin.value = true;
  try {
    await LoginAccount(loginIdentifier.value, loginPassword.value, loginNote.value);
    alert("Đăng nhập thành công và đã thêm vào danh sách!");
    showLoginModal.value = false; loginIdentifier.value = ''; loginPassword.value = ''; loginNote.value = '';
    await fetchAccounts();
  } catch (err: any) { alert("Lỗi đăng nhập: " + err); } finally { isSubmittingLogin.value = false; }
}
`;

content = content.replace('// Accounts Management', loginState + '\n// Accounts Management');

content = content.replace('<button \n          @click="showAddModal = true"', '<button\n          @click="showLoginModal = true"\n          class="px-4 py-2 rounded-lg border border-indigo-200 bg-white text-indigo-600 hover:bg-indigo-50 font-medium shadow-sm transition-colors flex items-center gap-2 whitespace-nowrap"\n        >\n          Login\n        </button>\n\n        <button \n          @click="showAddModal = true"');

fs.writeFileSync('frontend/src/views/Accounts.vue', content);
