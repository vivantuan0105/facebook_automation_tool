import sys
import re

file_path = 'd:/thuctap/facebook_automation_tool/frontend/src/views/Accounts.vue'
with open(file_path, 'r', encoding='utf-8') as f:
    content = f.read()

# 1. Add chi tiết bài viết button
content = content.replace(
'''          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5V4H2v16h5m5-8V4m0 8l-4 4m4-4l4 4" />
          </svg>
          Chi tiết bạn bè
        </button>
      </div>''',
'''          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5V4H2v16h5m5-8V4m0 8l-4 4m4-4l4 4" />
          </svg>
          Chi tiết bạn bè
        </button>
        <button
          @click="activeSection = 'posts'"
          :class="activeSection === 'posts'
            ? 'bg-indigo-600 text-white shadow-sm'
            : 'text-slate-600 hover:text-indigo-700 hover:bg-indigo-50'"
          class="px-4 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
          </svg>
          Chi tiết bài viết
        </button>
      </div>'''
)

# 2. Add Posts and Friends buttons together
content = content.replace(
'''              <td class="py-2 px-4 text-slate-600 border-r border-slate-100 text-center font-medium">
                <div class="flex items-center justify-center gap-2">
                  <span>{{ mem.acc.friends || '-' }}</span>
                  <button
                    @click="openFriendsView(mem.acc.uid)"
                    :class="selectedFriendUid === mem.acc.uid
                      ? 'border-indigo-300 bg-indigo-100 text-indigo-700'
                      : 'border-indigo-200 bg-indigo-50/80 text-indigo-500 hover:border-indigo-300 hover:bg-indigo-100 hover:text-indigo-700'"
                    class="inline-flex h-8 w-8 items-center justify-center rounded-lg border shadow-sm transition-colors"
                    title="Mở tab chi tiết bạn bè"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                    </svg>
                  </button>
                </div>
              </td>''',
'''              <td class="py-2 px-4 text-slate-600 border-r border-slate-100 text-center font-medium">
                <div class="flex flex-col items-center justify-center gap-1">
                  <div class="flex items-center gap-2">
                    <span class="text-[10px] font-semibold px-1.5 bg-slate-100 rounded text-slate-500">B: {{ mem.acc.friends || '-' }}</span>
                    <button
                      @click="openFriendsView(mem.acc.uid)"
                      :class="selectedFriendUid === mem.acc.uid && activeSection === 'friends'
                        ? 'border-indigo-300 bg-indigo-100 text-indigo-700'
                        : 'border-indigo-200 bg-indigo-50/80 text-indigo-500 hover:border-indigo-300 hover:bg-indigo-100 hover:text-indigo-700'"
                      class="inline-flex h-5 w-5 items-center justify-center rounded border shadow-sm transition-colors"
                      title="Mở tab chi tiết bạn bè"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                      </svg>
                    </button>
                  </div>
                  <div class="flex items-center gap-2">
                    <span class="text-[10px] font-semibold px-1.5 bg-slate-100 rounded text-slate-500">P: Xem</span>
                    <button
                      @click="openPostsView(mem.acc.uid)"
                      :class="selectedPostUid === mem.acc.uid && activeSection === 'posts'
                        ? 'border-emerald-300 bg-emerald-100 text-emerald-700'
                        : 'border-emerald-200 bg-emerald-50/80 text-emerald-500 hover:border-emerald-300 hover:bg-emerald-100 hover:text-emerald-700'"
                      class="inline-flex h-5 w-5 items-center justify-center rounded border shadow-sm transition-colors"
                      title="Mở tab chi tiết bài viết"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
                      </svg>
                    </button>
                  </div>
                </div>
              </td>'''
)

# 3. Add Posts Template after Friends
posts_template = '''
    <!-- Posts Detail Tab -->
    <div v-else-if="activeSection === 'posts'" class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
      <div class="border-b border-gray-100 px-6 py-4 flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <h3 class="text-lg font-bold text-slate-800">Chi tiết bài viết theo tài khoản</h3>
          <p class="text-sm text-slate-500">
            {{ selectedPostAccount ? `${selectedPostRows.length} bài viết hiển thị` : `${accounts.length} tài khoản có thể chọn` }}
          </p>
        </div>

        <div class="w-full lg:w-80">
          <input
            v-model="postSearch"
            type="text"
            placeholder="Tìm theo ID bài viết hoặc nội dung"
            class="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-700 focus:border-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-100"
          >
        </div>
      </div>

      <div class="p-6 bg-slate-50/50 min-h-[420px]">
        <div v-if="accounts.length === 0" class="rounded-xl border border-dashed border-slate-200 bg-white px-6 py-12 text-center text-slate-500">
          Chưa có tài khoản để hiển thị chi tiết bài viết.
        </div>
        <div v-else class="space-y-4">
          <div class="rounded-xl border border-slate-200 bg-white p-4 shadow-sm space-y-4">
            <div class="grid grid-cols-1 xl:grid-cols-[380px_280px] xl:justify-between gap-4 items-start">
              <div>
                <label class="mb-2 block text-sm font-medium text-slate-700">Tài khoản</label>
                <CustomSelect
                  v-model="selectedPostUid"
                  :options="friendAccountOptions"
                  placeholder="Chọn tài khoản để xem bài viết"
                  @change="handlePostAccountChange"
                />
              </div>

              <div>
                <label class="mb-2 block text-sm font-medium text-slate-700">Tìm bài viết</label>
                <input
                  v-model="postSearch"
                  type="text"
                  placeholder="Tìm theo text/ID"
                  class="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-700 focus:border-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-100"
                >
              </div>
            </div>

            <div v-if="selectedPostAccount" class="grid grid-cols-1 md:grid-cols-[minmax(0,1.1fr)_minmax(0,1fr)_160px] gap-3">
              <div class="rounded-xl border border-indigo-100 bg-indigo-50/70 px-4 py-3">
                <p class="text-xs font-semibold text-indigo-400 uppercase tracking-wider mb-1">Đang xem bài viết của ID</p>
                <p class="font-bold text-indigo-900 truncate">{{ selectedPostUid }}</p>
              </div>
              <div class="rounded-xl border border-emerald-100 bg-emerald-50/70 px-4 py-3">
                <p class="text-xs font-semibold text-emerald-500 uppercase tracking-wider mb-1">Tên tài khoản</p>
                <p class="font-bold text-emerald-900 truncate">{{ selectedPostAccount?.name || 'N/A' }}</p>
              </div>
            </div>
          </div>

          <div v-if="selectedPostAccount">
            <div v-if="loadingPostMap[selectedPostUid]" class="p-12 text-center text-slate-500 bg-white rounded-xl border border-slate-200">
              <div class="animate-spin inline-block w-8 h-8 border-4 border-indigo-500 border-t-transparent rounded-full mb-3"></div>
              <p>Đang tải danh sách bài viết từ ổ đĩa...</p>
            </div>
            
            <div v-else-if="postErrorMap[selectedPostUid]" class="p-8 text-center bg-rose-50 border border-rose-200 rounded-xl">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 mx-auto text-rose-400 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              <p class="text-rose-700 font-medium">{{ postErrorMap[selectedPostUid] }}</p>
            </div>

            <div v-else-if="selectedPostRows.length === 0" class="p-12 text-center bg-white border border-slate-200 rounded-xl text-slate-500">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto text-slate-300 mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
              </svg>
              <p class="text-lg font-medium text-slate-700">Không có bài viết nào</p>
              <p class="text-sm mt-1">Chưa quét được bài viết hoặc không khớp từ khóa tìm kiếm.</p>
            </div>

            <div v-else class="bg-white border border-slate-200 rounded-xl shadow-sm overflow-hidden flex flex-col max-h-[600px]">
              <div class="overflow-y-auto flex-1">
                <table class="w-full text-left text-sm whitespace-nowrap">
                  <thead class="bg-slate-50 text-slate-600 font-semibold sticky top-0 shadow-sm z-10 border-b border-slate-200">
                    <tr>
                      <th class="py-3 px-4 w-16 text-center border-r border-slate-200">#</th>
                      <th class="py-3 px-4 border-r border-slate-200">Post ID</th>
                      <th class="py-3 px-4 w-full">Nội dung text</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-slate-100">
                    <tr v-for="post in selectedPostRows" :key="post.index" class="hover:bg-indigo-50/50 transition-colors">
                      <td class="py-2.5 px-4 text-center text-slate-400 font-mono">{{ post.index }}</td>
                      <td class="py-2.5 px-4 border-r border-slate-100 font-medium text-indigo-600 hover:underline cursor-pointer" @click="copyToClipboard(post.id)">
                        {{ post.id }}
                      </td>
                      <td class="py-2.5 px-4 truncate max-w-lg text-slate-700" :title="post.text">{{ post.text }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
              <div class="border-t border-slate-100 bg-slate-50 px-4 py-3 text-sm text-slate-600 flex justify-between items-center">
                <span>Đang hiển thị <strong>{{ selectedPostRows.length }}</strong> bài viết</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
'''

content = content.replace('    <!-- Task Setup Modal -->', posts_template + '\n    <!-- Task Setup Modal -->')

with open(file_path, 'w', encoding='utf-8') as f:
    f.write(content)
