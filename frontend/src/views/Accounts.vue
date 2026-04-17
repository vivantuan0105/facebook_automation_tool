<template>
  <div class="space-y-6">
    <!-- Header Controls -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center bg-white p-4 rounded-xl shadow-sm border border-gray-100 gap-4">
      <div class="flex items-center gap-4">
        <h2 class="text-lg font-bold text-gray-800 border-r border-gray-200 pr-4">Quản lý Tài Khoản</h2>
        <div class="flex items-center gap-2">
          <button @click="showTaskModal = true" :disabled="selectedUids.length === 0" class="px-3 py-1.5 bg-blue-50 border border-blue-200 text-blue-700 hover:bg-blue-100 rounded-lg text-sm font-medium transition-colors shadow-sm flex items-center gap-1.5 disabled:opacity-50 disabled:cursor-not-allowed">
            Thiết lập Tác vụ
          </button>
          <button @click="runSelected" :disabled="selectedUids.length === 0" class="px-3 py-1.5 bg-indigo-50 border border-indigo-200 text-indigo-700 hover:bg-indigo-100 rounded-lg text-sm font-medium transition-colors shadow-sm flex items-center gap-1.5 disabled:opacity-50 disabled:cursor-not-allowed">
            Chạy
          </button>
          <button @click="stopSelected" :disabled="selectedUids.length === 0" class="px-3 py-1.5 bg-amber-50 border border-amber-200 text-amber-700 hover:bg-amber-100 rounded-lg text-sm font-medium transition-colors shadow-sm flex items-center gap-1.5 disabled:opacity-50 disabled:cursor-not-allowed">
            Dừng
          </button>
        </div>
      </div>
      
      <div class="flex items-center gap-2">
        <button 
          @click="showAccountChoiceModal = true"
          class="bg-blue-600 hover:bg-blue-700 text-white px-5 py-2.5 rounded-lg font-medium shadow-sm transition-colors flex items-center gap-2 whitespace-nowrap">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clip-rule="evenodd" />
          </svg>
          Thêm tài khoản
        </button>
      </div>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-2">
      <div class="flex flex-wrap items-center gap-2">
        <button
          @click="activeSection = 'accounts'"
          :class="activeSection === 'accounts'
            ? 'bg-indigo-600 text-white shadow-sm'
            : 'text-slate-600 hover:text-indigo-700 hover:bg-indigo-50'"
          class="px-4 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
          Danh sách tài khoản
        </button>
        <button
          @click="activeSection = 'friends'"
          :class="activeSection === 'friends'
            ? 'bg-indigo-600 text-white shadow-sm'
            : 'text-slate-600 hover:text-indigo-700 hover:bg-indigo-50'"
          class="px-4 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
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
        <button
          @click="activeSection = 'profile'"
          :class="activeSection === 'profile'
            ? 'bg-indigo-600 text-white shadow-sm'
            : 'text-slate-600 hover:text-indigo-700 hover:bg-indigo-50'"
          class="px-4 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
          Chi tiết Profile
        </button>
      </div>
    </div>

    <!-- Account List -->
    <div v-if="activeSection === 'accounts'" class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
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
          <thead class="bg-indigo-50/60 text-slate-700 font-semibold border-b-2 border-indigo-100/50">
            <tr>
              <th class="py-2.5 px-4 w-12 border-r border-indigo-100/30 text-center">
                <input type="checkbox" v-model="selectAll" class="rounded border-gray-300 text-indigo-600 focus:ring-indigo-500 w-3.5 h-3.5 shadow-sm">
              </th>
              <th class="py-2.5 px-4 border-r border-indigo-100/30 text-left">UID / Tên gợi nhớ</th>
              <th class="py-2.5 px-4 border-r border-indigo-100/30 text-center">Bạn bè</th>
              <th class="py-2.5 px-4 border-r border-indigo-100/30 text-center">Bài viết</th>
              <th class="py-2.5 px-4 border-r border-indigo-100/30 text-center">Tác vụ</th>
              <th class="py-2.5 px-4 border-r border-indigo-100/30 text-center">Mục tiêu</th>
              <th class="py-2.5 px-4 border-r border-indigo-100/30 text-center">Trạng thái Tác vụ</th>
              <th class="py-2.5 px-4 border-r border-indigo-100/30 text-center">Trạng thái Acc</th>
              <th class="py-2.5 px-4 text-center">Thao tác</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-for="mem in mergedAccounts" :key="mem.acc.uid" class="hover:bg-slate-50 transition-colors">
              <td class="py-2 px-4 border-r border-slate-100 text-center">
                <input type="checkbox" v-model="selectedUids" :value="mem.acc.uid" class="rounded border-gray-300 text-indigo-600 focus:ring-indigo-500 w-3.5 h-3.5 shadow-sm">
              </td>
              <td class="py-2 px-4 border-r border-slate-100 text-left">
                <div class="font-medium text-slate-700">{{ mem.acc.name || 'Không xác định' }}</div>
                <div class="text-xs text-slate-500">{{ mem.acc.uid }}</div>
              </td>
              <td class="py-2 px-4 text-slate-600 border-r border-slate-100 text-center font-medium">
                <div class="flex items-center justify-center gap-1.5">
                  <span class="text-[11px] font-semibold px-1.5 text-slate-600">{{ mem.acc.friends || '-' }}</span>
                  <button
                    @click="openFriendsView(mem.acc.uid)"
                    :class="selectedFriendUid === mem.acc.uid
                      ? 'border-indigo-300 bg-indigo-100 text-indigo-700'
                      : 'border-indigo-200 bg-indigo-50/80 text-indigo-500 hover:border-indigo-300 hover:bg-indigo-100 hover:text-indigo-700'"
                    class="inline-flex h-6 w-6 items-center justify-center rounded border shadow-sm transition-colors"
                    title="Mở tab chi tiết bạn bè"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                    </svg>
                  </button>
                </div>
              </td>
              <td class="py-2 px-4 text-slate-600 border-r border-slate-100 text-center font-medium">
                <div class="flex items-center justify-center gap-1.5">
                  <span class="text-[11px] font-semibold px-1.5 text-slate-600">{{ mem.acc.posts || '-' }}</span>
                  <button
                    @click="openPostsView(mem.acc.uid)"
                    :class="selectedPostUid === mem.acc.uid
                      ? 'border-emerald-300 bg-emerald-100 text-emerald-700'
                      : 'border-emerald-200 bg-emerald-50/80 text-emerald-500 hover:border-emerald-300 hover:bg-emerald-100 hover:text-emerald-700'"
                    class="inline-flex h-6 w-6 items-center justify-center rounded border shadow-sm transition-colors"
                    title="Mở tab chi tiết bài viết"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
                    </svg>
                  </button>
                </div>
              </td>
              <td class="py-2 px-4 text-slate-600 border-r border-slate-100 text-center">
                {{ mem.task ? mem.task.taskType : '-' }}
              </td>
              <td class="py-2 px-4 text-slate-600 border-r border-slate-100 text-center max-w-[150px] truncate" :title="mem.task ? (mem.task.postUrl || mem.task.postId) : ''">
                {{ mem.task ? (mem.task.postUrl ? 'URL' : 'ID: ' + mem.task.postId) : '-' }}
              </td>
              <td class="py-2 px-4 border-r border-slate-100 text-center">
                <span v-if="mem.task" class="inline-block px-2 py-0.5 rounded text-[11px] font-medium border" :class="{
                  'bg-amber-50 text-amber-600 border-amber-200': mem.task.status === 'Pending',
                  'bg-indigo-50 text-indigo-600 border-indigo-200': mem.task.status === 'Running',
                  'bg-emerald-50 text-emerald-600 border-emerald-200': mem.task.status === 'Success',
                  'bg-rose-50 text-rose-600 border-rose-200': mem.task.status === 'Failed',
                  'bg-slate-100 text-slate-500 border-slate-200': mem.task.status === 'Paused'
                }">
                  {{ statusMap[mem.task.status] || mem.task.status }}
                </span>
                <span v-else class="text-gray-400">-</span>
              </td>
              <td class="py-2 px-4 border-r border-slate-100 text-center">
                <span class="inline-flex items-center justify-center gap-1.5 px-2 py-0.5 rounded-full text-[11px] font-medium"
                  :class="mem.acc.status === 'Live' ? 'bg-emerald-50 text-emerald-600 border border-emerald-200' : 'bg-rose-50 text-rose-600 border border-rose-200'">
                  <span class="w-1.5 h-1.5 rounded-full" :class="mem.acc.status === 'Live' ? 'bg-emerald-500' : 'bg-rose-500'"></span>
                  {{ mem.acc.status }}
                </span>
              </td>
              <td class="py-2 px-4 flex justify-center gap-2">
                <button @click="openProfileView(mem.acc.uid)" class="text-blue-500 hover:text-blue-700 p-1 rounded hover:bg-blue-50 transition-colors" title="Xem chi tiết Profile">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  </svg>
                </button>
                <button @click="deleteAccount(mem.acc.uid)" class="text-rose-500 hover:text-rose-700 p-1 rounded hover:bg-rose-50 transition-colors" title="Xóa tài khoản">
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

    <!-- Profile Detail Tab -->
    <div v-else-if="activeSection === 'profile'" class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
      <div class="border-b border-gray-100 px-6 py-4 flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <h3 class="text-lg font-bold text-slate-800">Hồ sơ cá nhân (Profile)</h3>
          <p class="text-sm text-slate-500">Xem chi tiết thông tin nhân khẩu học của tài khoản</p>
        </div>
      </div>
      <div class="p-6 bg-slate-50/50 min-h-[420px]">
        <div v-if="accounts.length === 0" class="rounded-xl border border-dashed border-slate-200 bg-white px-6 py-12 text-center text-slate-500">
          Chưa có tài khoản nào được lưu trên hệ thống
        </div>
        <div v-else class="space-y-4">
          <div class="rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
            <div class="max-w-[380px]">
              <label class="mb-2 block text-sm font-medium text-slate-700">Tài khoản hiển thị thông tin</label>
              <CustomSelect
                v-model="selectedProfileUid"
                :options="friendAccountOptions"
                placeholder="Chọn tài khoản để xem Profile"
              />
            </div>
          </div>
          
          <div v-if="!currentProfileData" class="rounded-xl border border-dashed border-slate-200 bg-white px-6 py-12 text-center text-slate-500">
            Vui lòng chọn tài khoản ở bộ lọc phía trên để xem chi tiết
          </div>
          <div v-else class="rounded-xl border border-slate-200 bg-white p-6 shadow-sm">
            <div class="flex items-center justify-between mb-6">
              <div class="flex items-center gap-4">
                <div class="h-16 w-16 bg-slate-100 rounded-full flex items-center justify-center border border-slate-200 overflow-hidden">
                   <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="text-slate-400"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                </div>
                <div>
                  <h4 class="text-xl font-bold text-slate-800">{{ currentProfileData.name || 'Chưa cập nhật' }}</h4>
                  <p class="text-sm text-slate-500 mt-0.5">UID: {{ currentProfileData.uid }}</p>
                </div>
              </div>
            </div>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-x-8 gap-y-6 pt-4 border-t border-slate-100">
              <div>
                <label class="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-1 block">Giới tính</label>
                <div class="font-medium text-slate-700">{{ currentProfileData.gender || 'Không có thông tin' }}</div>
              </div>
              <div>
                <label class="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-1 block">Followers</label>
                <div class="font-medium text-slate-700">{{ currentProfileData.followers || '0' }}</div>
              </div>
              <div>
                <label class="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-1 block">Ngày sinh</label>
                <div class="font-medium text-slate-700">{{ currentProfileData.birthday || 'Không có' }}</div>
              </div>
              <div>
                <label class="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-1 block">Năm sinh</label>
                <div class="font-medium text-slate-700">{{ currentProfileData.birthYear || 'Không có' }}</div>
              </div>
              <div class="md:col-span-2">
                <label class="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-1 block">Vị trí (Location)</label>
                <div class="font-medium text-slate-700">{{ currentProfileData.location || 'Không tiết lộ' }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Friends Detail Tab -->
    <div v-else-if="activeSection === 'friends'" class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
      <div class="border-b border-gray-100 px-6 py-4 flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <h3 class="text-lg font-bold text-slate-800">Chi tiết bạn bè theo tài khoản</h3>
          <p class="text-sm text-slate-500">
            {{ selectedFriendAccount ? `${selectedFriendRows.length} / ${selectedFriendTotal} bạn bè hiển thị` : `${accounts.length} tài khoản có thể chọn` }}
          </p>
        </div>

        <div class="w-full lg:w-80">
          <input
            v-model="friendSearch"
            type="text"
            placeholder="Tìm theo UID hoặc tên bạn bè"
            class="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-700 focus:border-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-100"
          >
        </div>
      </div>

      <div class="p-6 bg-slate-50/50 min-h-[420px]">
        <div v-if="accounts.length === 0" class="rounded-xl border border-dashed border-slate-200 bg-white px-6 py-12 text-center text-slate-500">
          Chưa có tài khoản để hiển thị chi tiết bạn bè.
        </div>
        <div v-else class="space-y-4">
          <div class="rounded-xl border border-slate-200 bg-white p-4 shadow-sm space-y-4">
            <div class="grid grid-cols-1 xl:grid-cols-[380px_280px] xl:justify-between gap-4 items-start">
              <div>
                <label class="mb-2 block text-sm font-medium text-slate-700">Tài khoản</label>
                <CustomSelect
                  v-model="selectedFriendUid"
                  :options="friendAccountOptions"
                  placeholder="Chọn tài khoản để xem bạn bè"
                  @change="handleFriendAccountChange"
                />
              </div>

              <div>
                <label class="mb-2 block text-sm font-medium text-slate-700">Tìm bạn bè</label>
                <input
                  v-model="friendSearch"
                  type="text"
                  placeholder="Tìm UID hoặc tên"
                  class="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-700 focus:border-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-100"
                >
              </div>
            </div>

            <div v-if="selectedFriendAccount" class="grid grid-cols-1 md:grid-cols-[minmax(0,1.1fr)_minmax(0,1fr)_160px] gap-3">
              <div class="rounded-xl border border-indigo-100 bg-indigo-50/70 px-4 py-3">
                <div class="text-xs font-semibold uppercase tracking-wide text-indigo-500">Tài khoản đang xem</div>
                <div class="mt-1 text-lg font-semibold text-slate-800 break-all">{{ selectedFriendAccount.name || 'Không xác định' }}</div>
              </div>
              <div class="rounded-xl border border-slate-200 bg-slate-50/70 px-4 py-3">
                <div class="text-xs font-semibold uppercase tracking-wide text-slate-500">UID tài khoản</div>
                <div class="mt-1 text-base font-semibold text-slate-800 break-all">{{ selectedFriendAccount.uid }}</div>
              </div>
              <div class="rounded-xl border border-emerald-100 bg-emerald-50/70 px-4 py-3">
                <div class="text-xs font-semibold uppercase tracking-wide text-emerald-500">Số bạn bè</div>
                <div class="mt-1 text-2xl font-bold text-slate-800">{{ selectedFriendTotal }}</div>
              </div>
            </div>
          </div>

          <div v-if="!selectedFriendAccount" class="rounded-xl border border-dashed border-slate-200 bg-white px-6 py-12 text-center text-slate-500">
            Chọn `tài khoản` để xổ danh sách bạn bè.
          </div>

          <div v-else-if="loadingFriendMap[selectedFriendAccount.uid]" class="flex min-h-[320px] items-center justify-center rounded-xl border border-slate-200 bg-white text-slate-500">
            <div class="text-center">
              <svg class="mx-auto h-8 w-8 animate-spin text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <p class="mt-3 text-sm">Đang tải danh sách bạn bè...</p>
            </div>
          </div>

          <div v-else-if="friendErrorMap[selectedFriendAccount.uid]" class="rounded-xl border border-amber-200 bg-amber-50 px-4 py-4 text-sm text-amber-700">
            {{ friendErrorMap[selectedFriendAccount.uid] }}
          </div>

          <div v-else-if="selectedFriendRows.length === 0" class="rounded-xl border border-dashed border-slate-200 bg-white px-6 py-12 text-center text-slate-500">
            {{ friendSearch ? 'Không có bạn bè nào khớp từ khóa.' : 'Tài khoản này chưa có dữ liệu bạn bè.' }}
          </div>

          <div v-else class="overflow-hidden rounded-xl border border-indigo-100 bg-white shadow-sm">
            <div class="border-b border-slate-100 px-4 py-3">
              <h4 class="font-semibold text-slate-800">Danh sách bạn bè của {{ selectedFriendAccount.name || selectedFriendAccount.uid }}</h4>
            </div>

            <div class="overflow-x-auto">
              <table class="w-full text-left text-sm">
                <thead class="bg-slate-50 text-slate-700">
                  <tr>
                    <th class="w-20 border-b border-slate-100 px-4 py-3 text-center font-semibold">STT</th>
                    <th class="w-72 border-b border-slate-100 px-4 py-3 font-semibold">UID</th>
                    <th class="border-b border-slate-100 px-4 py-3 font-semibold">Tên bạn bè</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100">
                  <tr
                    v-for="friend in selectedFriendRows"
                    :key="`${selectedFriendAccount.uid}-${friend.index}-${friend.raw}`"
                    class="transition-colors hover:bg-slate-50"
                  >
                    <td class="px-4 py-3 text-center">
                      <span class="inline-flex h-8 w-8 items-center justify-center rounded-full bg-indigo-100 font-bold text-indigo-600">
                        {{ friend.index }}
                      </span>
                    </td>
                    <td class="px-4 py-3 font-medium text-slate-600 break-all">
                      {{ friend.uid || 'Không tách được UID từ dữ liệu quét' }}
                    </td>
                    <td class="px-4 py-3 font-semibold text-slate-800 break-all">
                      {{ friend.name || friend.raw }}
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>


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
                      <th class="py-3 px-4 border-r border-slate-200 w-48">Post ID</th>
                      <th class="py-3 px-4 border-r border-slate-200 w-40">Thời gian</th>
                      <th class="py-3 px-4 w-full">Nội dung text</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-slate-100">
                    <tr v-for="post in selectedPostRows" :key="post.index" class="hover:bg-indigo-50/50 transition-colors">
                      <td class="py-2.5 px-4 text-center text-slate-400 font-mono">{{ post.index }}</td>
                      <td class="py-2.5 px-4 border-r border-slate-100 font-medium text-indigo-600 hover:underline cursor-pointer" @click="openPostUrl(post.url || post.id)" title="Click để mở bài viết trên Facebook">
                        {{ post.id }}
                      </td>
                      <td class="py-2.5 px-4 border-r border-slate-100 text-slate-500 text-sm whitespace-nowrap">
                        <div class="flex items-center gap-1.5" v-if="post.time">
                          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-clock"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                          {{ formatTimestamp(post.time) }}
                        </div>
                        <span v-else class="text-slate-300">-</span>
                      </td>
                      <td class="py-2.5 px-4 truncate max-w-lg text-slate-700" :title="post.text" @click="copyToClipboard(post.text)" style="cursor: copy;">{{ post.text }}</td>
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

    <!-- Task Setup Modal -->
    <div v-if="showTaskModal" @mousedown.self="showTaskModal = false" class="fixed inset-0 bg-slate-800/10 z-[100] pointer-events-auto transition-opacity">
      <div 
        class="bg-white rounded-xl shadow-2xl w-full max-w-2xl min-h-[520px] overflow-hidden shadow-[0_0_20px_rgba(0,0,0,0.15)] flex flex-col max-h-[90vh] absolute"
        :style="{ top: taskModalY + 'px', left: taskModalX + 'px' }"
      >
        <div class="px-6 py-4 border-b border-gray-100 flex justify-between items-center bg-blue-50 shrink-0 cursor-move" @mousedown="startDragTaskModal">
          <h3 class="text-lg font-bold text-gray-800">Thiết lập Tác vụ Hàng loạt</h3>
          <button @click="showTaskModal = false" class="text-gray-400 hover:text-gray-600 transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <div class="p-6 overflow-y-auto space-y-5 flex-1">
          <div v-if="taskError" class="mt-2 p-3 bg-red-50 text-red-600 text-sm rounded border border-red-100">
            Lỗi: {{ taskError }}
          </div>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1">Loại tác vụ</label>
              <CustomSelect v-model="form.taskType" :options="taskTypeOptions" @change="onTaskTypeChange"/>
            </div>

            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1">Xác định mục tiêu</label>
              <CustomSelect v-model="form.targetMode" :options="targetModeOptions"/>
            </div>

            <div class="md:col-span-2" v-if="form.targetMode === 'post_url'">
              <label class="block text-sm font-medium text-slate-700 mb-1">URL Bài viết</label>
              <input type="text" v-model="form.postUrl" class="w-full text-sm bg-white border border-gray-300 rounded-lg p-2.5 focus:ring-primary focus:border-primary" placeholder="Ví dụ: https://www.facebook.com/permalink.php?..." />
            </div>

            <div class="md:col-span-2" v-if="form.taskType === 'Like bài viết'">
              <label class="block text-sm font-medium text-slate-700 mb-1">Loại cảm xúc</label>
              <CustomSelect v-model="form.reactionType" :options="reactionTypeOptions" placeholder="Chọn cảm xúc"/>
            </div>

            <div class="md:col-span-2" v-if="['Comment bài viết', 'Đăng bài viết'].includes(form.taskType)">
              <label class="block text-sm font-medium text-slate-700 mb-1">{{ form.taskType === 'Đăng bài viết' ? 'Nội dung bài viết (Status)' : 'Nội dung bình luận' }}</label>
              <textarea v-model="form.message" class="w-full text-sm bg-white border border-gray-300 rounded-lg p-2.5 focus:ring-primary focus:border-primary min-h-[100px]" placeholder="Nhập nội dung... Hỗ trợ spin {A|B}"></textarea>
            </div>

            <div class="md:col-span-2" v-if="form.taskType === 'Đăng bài viết'">
              <label class="block text-sm font-medium text-slate-700 mb-1">Tệp đính kèm (Ảnh/Video)</label>
              <div class="space-y-3">
                <div v-if="form.photoPaths.length > 0" class="space-y-2">
                  <div v-for="(_, index) in form.photoPaths" :key="index" class="flex gap-2 items-center bg-gray-50 px-3 py-2.5 rounded-lg border border-gray-200">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
                    </svg>
                    <span class="flex-1 text-sm text-gray-700 truncate" :title="form.photoPaths[index]">{{ form.photoPaths[index].split('\\').pop() || form.photoPaths[index].split('/').pop() }}</span>
                    <button @click="removePhoto(index)" class="px-3 py-1 border border-red-200 bg-white hover:bg-red-50 text-red-600 text-sm rounded transition-colors whitespace-nowrap">Xóa</button>
                  </div>
                </div>
                <div>
                  <button @click="selectPhoto" class="w-full px-4 py-2.5 bg-slate-50 border border-slate-300 border-dashed hover:bg-slate-100 text-slate-700 text-sm font-medium rounded-lg transition-colors flex items-center justify-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" /></svg>
                    Chọn file Ảnh / Video từ máy tính
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="px-6 py-4 bg-gray-50 border-t border-gray-100 flex justify-end gap-3 shrink-0">
          <button @click="showTaskModal = false" class="px-4 py-2 text-gray-600 font-medium hover:bg-gray-200 rounded-lg transition-colors">Hủy</button>
          <button @click="applyTasks" :disabled="isSubmittingTask" class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg shadow-sm transition-colors flex items-center gap-2 disabled:opacity-70">
            Lưu Tác vụ cho {{selectedUids.length}} acc
          </button>
        </div>
      </div>
    </div>

    <!-- Tabbed Add Account Modal -->
    <div v-if="showAccountChoiceModal" @mousedown.self="showAccountChoiceModal = false" class="fixed inset-0 bg-slate-800/40 backdrop-blur-sm z-[100] pointer-events-auto transition-opacity flex items-center justify-center p-4">
      <div class="bg-white rounded-2xl shadow-[0_20px_50px_rgba(0,0,0,0.2)] w-full max-w-3xl overflow-hidden transform transition-all flex flex-col max-h-[90vh] animate-in fade-in zoom-in duration-200">
        <!-- Header -->
        <div class="px-6 py-5 border-b border-gray-100 flex justify-between items-center bg-gradient-to-r from-blue-50 to-indigo-50 shrink-0">
          <h3 class="text-xl font-bold text-gray-800 flex items-center gap-2">
            Thêm tài khoản
          </h3>
          <button @click="showAccountChoiceModal = false" class="text-gray-400 hover:text-gray-600 hover:bg-white rounded-full p-2 transition-all shadow-sm">
             <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <!-- Tabs -->
        <div class="flex border-b border-gray-200 bg-gray-50/50 shrink-0 px-6 pt-4 gap-2 overflow-x-auto">
          <button @click="addModalType = 'login'" :class="addModalType === 'login' ? 'border-blue-500 text-blue-700 bg-white shadow-sm' : 'border-transparent text-gray-500 hover:text-gray-700 hover:bg-gray-100/50'" class="px-5 py-3 border-b-2 font-bold text-sm transition-all rounded-t-xl flex items-center gap-2">
            Thêm bằng Tên đăng nhập/mật khẩu
          </button>
          <button @click="addModalType = 'cookie'" :class="addModalType === 'cookie' ? 'border-blue-500 text-blue-700 bg-white shadow-sm' : 'border-transparent text-gray-500 hover:text-gray-700 hover:bg-gray-100/50'" class="px-5 py-3 border-b-2 font-bold text-sm transition-all rounded-t-xl flex items-center gap-2">
            Thêm bằng Cookie
          </button>
          <button @click="addModalType = 'bulk'" :class="addModalType === 'bulk' ? 'border-blue-500 text-blue-700 bg-white shadow-sm' : 'border-transparent text-gray-500 hover:text-gray-700 hover:bg-gray-100/50'" class="px-5 py-3 border-b-2 font-bold text-sm transition-all rounded-t-xl flex items-center gap-2">
            Import hàng loạt
          </button>
        </div>

        <!-- Body (Scrollable) -->
        <div class="p-6 overflow-y-auto flex-1 bg-gray-50/30">
          
          <!-- Login Form Content -->
          <div v-if="addModalType === 'login'" class="space-y-4 animate-in fade-in slide-in-from-bottom-2 duration-300">
            <div v-if="loginError" class="bg-red-50 border border-red-200 text-red-600 text-sm p-3 rounded-lg flex items-start gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 shrink-0 mt-0.5" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
              </svg>
              {{ loginError }}
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Tài khoản (Tên đăng nhập / Email / Số điện thoại)</label>
              <input v-model="loginUsername" type="text" placeholder="Nhập tài khoản" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 outline-none transition-all">
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Mật khẩu</label>
              <input v-model="loginPassword" type="password" placeholder="Nhập mật khẩu" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 outline-none transition-all">
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Mã 2FA (Nếu có)</label>
              <input v-model="loginTwoFactorAuth" type="text" placeholder="Chuỗi khóa bảo mật 2FA hoặc mã 6 số" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 outline-none transition-all">
            </div>
          </div>

          <!-- Cookie Form Content -->
          <div v-if="addModalType === 'cookie'" class="space-y-4 animate-in fade-in slide-in-from-bottom-2 duration-300">
            <div v-if="addError" class="bg-red-50 border border-red-200 text-red-600 text-sm p-3 rounded-lg flex items-start gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 shrink-0 mt-0.5" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
              </svg>
              {{ addError }}
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Tên gợi nhớ</label>
              <input v-model="newAccName" type="text" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all">
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Chuỗi Cookie Facebook</label>
              <textarea v-model="newAccCookie" rows="6" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all font-mono text-sm"></textarea>
              <p class="text-xs text-gray-500 mt-1">Hệ thống sẽ tự động bóc tách c_user (UID) làm định danh.</p>
            </div>
          </div>

          <!-- Bulk Form Content -->
          <div v-if="addModalType === 'bulk'" class="space-y-4 animate-in fade-in slide-in-from-bottom-2 duration-300">
            <div v-if="bulkResult && bulkResult.error" class="bg-red-50 border border-red-200 text-red-600 text-sm p-3 rounded-lg flex items-start gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 shrink-0 mt-0.5" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
              </svg>
              {{ bulkResult.error }}
            </div>

            <div v-if="bulkResult && !bulkResult.error" class="bg-emerald-50 border border-emerald-200 p-4 rounded-lg">
               <p class="font-medium text-emerald-800 text-lg">Kết quả Import:</p>
               <p class="text-emerald-700 mt-1">✔ Thành công: <strong>{{ bulkResult.successCount }}</strong> / {{ bulkResult.totalProcessed }} dòng</p>
               <p v-if="bulkResult.failedCount > 0" class="text-red-600 mt-1">✘ Thất bại: <strong>{{ bulkResult.failedCount }}</strong></p>
               <div v-if="bulkResult.errors && bulkResult.errors.length > 0" class="mt-2 bg-white text-red-500 text-xs p-2 rounded max-h-24 overflow-y-auto border border-red-100">
                  <div v-for="(err, idx) in bulkResult.errors" :key="idx">{{ err }}</div>
               </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Danh sách Cookie (Định dạng <code class="bg-gray-100 px-1 py-0.5 rounded">UID|Pass|Cookie</code> hoặc chỉ <code class="bg-gray-100 px-1 py-0.5 rounded">Cookie</code>)</label>
              <textarea v-model="bulkCookies" rows="8" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all font-mono text-sm whitespace-pre"></textarea>
              <p class="text-xs text-gray-500 mt-2">Mỗi tài khoản trên một dòng. 
                  Hệ thống tự động tách và tìm chỗ chứa <strong>c_user=</strong> thành Cookie an toàn.</p>
            </div>
          </div>

        </div>
        
        <!-- Footer -->
        <div class="px-6 py-4 bg-gray-50 border-t border-gray-100 flex justify-end gap-3 shrink-0">
          
          <button v-if="addModalType === 'login'" @click="submitLogin" :disabled="isSubmittingLogin" class="px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg shadow-sm transition-colors flex items-center gap-2 disabled:opacity-70">
            <svg v-if="isSubmittingLogin" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Thêm
          </button>

          <button v-if="addModalType === 'cookie'" @click="submitAdd" :disabled="isSubmitting" class="px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg shadow-sm transition-colors flex items-center gap-2 disabled:opacity-70">
            <svg v-if="isSubmitting" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Thêm
          </button>

          <button v-if="addModalType === 'bulk'" @click="submitBulkAdd" :disabled="isSubmittingBulk" class="px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg shadow-sm transition-colors flex items-center gap-2 disabled:opacity-70">
            <svg v-if="isSubmittingBulk" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Import
          </button>
        </div>
      </div>
    </div>
    
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue'
import { GetAllAccounts, AddAccount, ImportMultipleAccounts, DeleteAccount, SelectPhotoDialog, GetAccountFriendsList, GetAccountPostsList, LoginWithPassword } from '../../wailsjs/go/app/App'
import { useMainStore } from '../stores/main'
import CustomSelect from '../components/CustomSelect.vue'

const store = useMainStore()
const accounts = ref<any[]>([])
const loading = ref(true)
let interval: any

// Accounts Management
const showAccountChoiceModal = ref(false)
const addModalType = ref<'login' | 'cookie' | 'bulk'>('login')
const isSubmitting = ref(false)
const newAccName = ref('')
const newAccCookie = ref('')
const addError = ref('')

// Bulk Add Management
const isSubmittingBulk = ref(false)
const bulkCookies = ref('')
const bulkResult = ref<any>(null)

// Login Management
const isSubmittingLogin = ref(false)
const loginUsername = ref('')
const loginPassword = ref('')
const loginTwoFactorAuth = ref('')
const loginError = ref('')

// Setup Tabs Navigation
const activeSection = ref<'accounts' | 'profile' | 'friends' | 'posts'>('accounts')

// Profile Detail View
const selectedProfileUid = ref<string>('')
const currentProfileData = computed(() => {
  return accounts.value.find(acc => acc.uid === selectedProfileUid.value) || null
})

const openProfileView = (uid: string) => {
  activeSection.value = 'profile'
  selectedProfileUid.value = uid
}

// Friends Detail Tab
const selectedFriendUid = ref<string>('')
const loadingFriendMap = ref<Record<string, boolean>>({})
const friendErrorMap = ref<Record<string, string>>({})
const friendsByAccount = ref<Record<string, string[]>>({})
const friendSearch = ref('')

const selectedFriendAccount = computed(() => {
  return accounts.value.find((acc) => acc.uid === selectedFriendUid.value) || null
})

const friendAccountOptions = computed(() =>
  accounts.value.map((acc) => ({
    value: acc.uid,
    label: `${acc.name || 'Không xác định'} | ${acc.uid}`,
  }))
)

const parseFriendEntries = (friends: string[]) => {
  return friends.map((friend, index) => {
    const [uidPart, ...nameParts] = friend.split('|')
    const uid = uidPart?.trim() || ''
    const name = nameParts.join('|').trim()

    return {
      index: index + 1,
      raw: friend,
      uid,
      name: name || friend,
    }
  })
}

const getFilteredFriends = (uid: string) => {
  const keyword = friendSearch.value.trim().toLowerCase()
  const parsed = parseFriendEntries(friendsByAccount.value[uid] || [])

  if (!keyword) {
    return parsed
  }

  return parsed.filter((friend) =>
    `${friend.raw} ${friend.uid} ${friend.name}`.toLowerCase().includes(keyword)
  )
}

const selectedFriendRows = computed(() => {
  if (!selectedFriendUid.value) return []
  return getFilteredFriends(selectedFriendUid.value)
})

const selectedFriendTotal = computed(() => {
  if (!selectedFriendUid.value) return 0
  return (friendsByAccount.value[selectedFriendUid.value] || []).length
})

const loadFriendsList = async (uid: string) => {
  friendErrorMap.value = { ...friendErrorMap.value, [uid]: '' }
  loadingFriendMap.value = { ...loadingFriendMap.value, [uid]: true }

  try {
    const list = await GetAccountFriendsList(uid)
    friendsByAccount.value = {
      ...friendsByAccount.value,
      [uid]: list || [],
    }
  } catch (err: any) {
    console.error(err)
    friendErrorMap.value = {
      ...friendErrorMap.value,
      [uid]: 'Chưa có dữ liệu bạn bè cho tài khoản này. Hãy quét bạn bè trước rồi mở lại tab chi tiết.'
    }
  } finally {
    loadingFriendMap.value = { ...loadingFriendMap.value, [uid]: false }
  }
}

const handleFriendAccountChange = async (uid: string) => {
  if (!uid) return
  selectedFriendUid.value = uid
  if (!friendsByAccount.value[uid] && !loadingFriendMap.value[uid]) {
    await loadFriendsList(uid)
  }
}

const openFriendsView = async (uid: string) => {
  activeSection.value = 'friends'
  selectedFriendUid.value = uid
  friendSearch.value = ''
  await handleFriendAccountChange(uid)
}

// Posts Detail Tab
const selectedPostUid = ref<string>('')
const loadingPostMap = ref<Record<string, boolean>>({})
const postErrorMap = ref<Record<string, string>>({})
const postsByAccount = ref<Record<string, string[]>>({})
const postSearch = ref('')

const selectedPostAccount = computed(() => {
  return accounts.value.find((acc) => acc.uid === selectedPostUid.value) || null
})

const parsePostEntries = (posts: string[]) => {
  return posts.map((post, index) => {
    // format: postID|url|timeX|content
    const parts = post.split('|')
    if (parts.length >= 4 && !post.startsWith('[')) {
      return {
        index: index + 1,
        id: parts[0],
        url: parts[1],
        time: parseInt(parts[2]),
        text: parts.slice(3).join('|')
      }
    }
    // format cu: [12345] content text...
    const match = post.match(/^\[(.*?)\] (.*)$/)
    if (match) {
      return {
        index: index + 1,
        id: match[1],
        url: '',
        time: 0,
        text: match[2]
      }
    }
    return {
      index: index + 1,
      id: '',
      url: '',
      time: 0,
      text: post
    }
  })
}

const formatTimestamp = (ts: number) => {
  if (!ts) return ''
  const d = new Date(ts * 1000)
  return `${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')} ${d.getDate().toString().padStart(2, '0')}/${(d.getMonth() + 1).toString().padStart(2, '0')}/${d.getFullYear()}`
}

const getFilteredPosts = (uid: string) => {
  const keyword = postSearch.value.trim().toLowerCase()
  const parsed = parsePostEntries(postsByAccount.value[uid] || [])

  if (!keyword) {
    return parsed
  }

  return parsed.filter((post) =>
    `${post.id} ${post.text}`.toLowerCase().includes(keyword)
  )
}

const selectedPostRows = computed(() => {
  if (!selectedPostUid.value) return []
  return getFilteredPosts(selectedPostUid.value)
})

const loadPostsList = async (uid: string) => {
  postErrorMap.value = { ...postErrorMap.value, [uid]: '' }
  loadingPostMap.value = { ...loadingPostMap.value, [uid]: true }

  try {
    const list = await GetAccountPostsList(uid)
    postsByAccount.value = {
      ...postsByAccount.value,
      [uid]: list || [],
    }
  } catch (err: any) {
    console.error(err)
    if (String(err).includes('cannot find the file') || String(err).includes('no such file')) {
        postErrorMap.value = {
          ...postErrorMap.value,
          [uid]: 'Chưa có dữ liệu bài viết. Hãy tạo Tác vụ Quét bài viết cho tài khoản này trước.'
        }
    } else {
        postErrorMap.value = {
          ...postErrorMap.value,
          [uid]: String(err)
        }
    }
  } finally {
    loadingPostMap.value = { ...loadingPostMap.value, [uid]: false }
  }
}

const handlePostAccountChange = async (uid: string) => {
  if (!uid) return
  selectedPostUid.value = uid
  if (!postsByAccount.value[uid] && !loadingPostMap.value[uid]) {
    await loadPostsList(uid)
  }
}

const openPostsView = async (uid: string) => {
  activeSection.value = 'posts'
  selectedPostUid.value = uid
  postSearch.value = ''
  await handlePostAccountChange(uid)
}

const openPostUrl = (postLink: string) => {
  if (postLink) {
    const url = postLink.startsWith('http') ? postLink : `https://www.facebook.com/${postLink}`
    window.open(url, '_blank')
  }
}

const copyToClipboard = (text: string) => {
  navigator.clipboard.writeText(text)
  // Simple toast placeholder
}

// Selected & Merged logic
const selectedUids = ref<string[]>([])
const clearedTaskIds = ref<Set<number>>(new Set())

watch(selectedUids, (newUids, oldUids) => {
  if (oldUids && oldUids.length > 0) {
    const removedUids = oldUids.filter(id => !newUids.includes(id))
    for (const uid of removedUids) {
      const mem = mergedAccounts.value.find(m => m.acc.uid === uid)
      if (mem && mem.realTask) {
        clearedTaskIds.value.add(mem.realTask.id)
      }
    }
  }
})

const mergedAccounts = computed(() => {
  return accounts.value.map(acc => {
    const accTasks = store.tasks.filter((t: any) => t.cookie === acc.cookie && !clearedTaskIds.value.has(t.id))
    accTasks.sort((a: any, b: any) => b.id - a.id)
    const activeTask = accTasks.find((t: any) => t.status !== 'Success' && t.status !== 'Failed') || accTasks[0]
    
    let displayTask = activeTask

    return { acc, task: displayTask, realTask: activeTask }
  })
})

const selectAll = computed({
  get: () => mergedAccounts.value.length > 0 && selectedUids.value.length === mergedAccounts.value.length,
  set: (v) => {
    if (v) selectedUids.value = mergedAccounts.value.map(a => a.acc.uid)
    else selectedUids.value = []
  }
})

watch(() => store.filterStatusTrigger, (status) => {
  if (status) {
    selectedUids.value = mergedAccounts.value
      .filter(m => m.realTask && m.realTask.status === status)
      .map(m => m.acc.uid)
    store.filterStatusTrigger = null
  }
})

// Task Form Modal Logic
const showTaskModal = ref(false)
const isSubmittingTask = ref(false)
const taskError = ref('')

const statusMap: Record<string, string> = {
  'Pending': 'Đang chờ',
  'Running': 'Đang chạy',
  'Success': 'Thành công',
  'Failed': 'Thất bại',
  'Paused': 'Tạm dừng'
}

const form = reactive({
  taskType: 'Like bài viết',
  execMode: 'live',
  targetMode: 'post_url',
  postUrl: '',
  postId: '',
  reactionType: 'Like',
  message: '',
  photoPaths: [] as string[],
})

const taskTypeOptions = [
  { value: 'Like bài viết', label: 'Thích (Reaction) bài viết' },
  { value: 'Comment bài viết', label: 'Bình luận bài viết' },
  { value: 'Đăng bài viết', label: 'Đăng bài viết (Post)' },
  { value: 'Quét thông tin', label: 'Quét thông tin Cá nhân (Scan Profile)' },
  { value: 'Quét bạn bè', label: 'Quét danh sách Bạn bè (Scan Friends)' },
  { value: 'Quét bài viết', label: 'Quét bài viết (Scan Timeline)' },
  { value: 'Reaction batch', label: 'Reaction hàng loạt (Sắp có)', disabled: true }
]

const targetModeOptions = computed(() => {
  if (form.taskType === 'Đăng bài viết') return [{ value: 'timeline', label: 'Đăng lên trang cá nhân (Timeline)' }]
  if (['Quét thông tin', 'Quét bạn bè'].includes(form.taskType)) {
    return [{ value: 'self_profile', label: 'Quét chính tài khoản này' }]
  }
  if (form.taskType === 'Quét bài viết') {
    return [
      { value: 'self_profile', label: 'Dùng chính tài khoản này (Quét chính mình)' },
      { value: 'post_url', label: 'Quét UID bất kỳ (Nhập vào URL bên dưới)' }
    ]
  }
  return [{ value: 'post_url', label: 'Dùng URL Bài viết (Khuyên dùng)' }]
})
const reactionTypeOptions = [
  { value: 'Like', label: 'Like 👍' },
  { value: 'Love', label: 'Love ❤️' },
  { value: 'Care', label: 'Care 🥰' },
  { value: 'Haha', label: 'Haha 😂' },
  { value: 'Wow', label: 'Wow 😮' },
  { value: 'Sad', label: 'Sad 😢' },
  { value: 'Angry', label: 'Angry 😡' },
]

const onTaskTypeChange = () => {
  if (form.taskType === 'Đăng bài viết') {
    form.targetMode = 'timeline'
    form.postUrl = ''
  } else if (['Quét thông tin', 'Quét bạn bè', 'Quét bài viết'].includes(form.taskType)) {
    form.targetMode = 'self_profile'
    form.postUrl = ''
  } else {
    form.targetMode = 'post_url'
  }
}


const removePhoto = (index: number) => {
  form.photoPaths.splice(index, 1)
}

const selectPhoto = async () => {
  try {
    const path = await SelectPhotoDialog()
    if (path) {
      form.photoPaths.push(path)
    }
  } catch (err) {
    console.error(err)
  }
}

const applyTasks = async () => {
  taskError.value = ''
  isSubmittingTask.value = true
  try {
    for (const uid of selectedUids.value) {
      const match = accounts.value.find(a => a.uid === uid)
      if (match) {
        // Validation check clone
        const f = { ...form, cookie: match.cookie }
        const err = await store.validateTask(f)
        if (err) {
          taskError.value = `Lỗi ở acc ${uid}: ${err}`;
          return;
        }
        await store.createTask(f)
      }
    }
    showTaskModal.value = false
  } finally {
    isSubmittingTask.value = false
  }
}

const runSelected = async () => {
  for (const uid of selectedUids.value) {
    const mem = mergedAccounts.value.find(m => m.acc.uid === uid)
    if (mem && mem.realTask && (mem.realTask.status === 'Pending' || mem.realTask.status === 'Paused' || mem.realTask.status === 'Failed')) {
      await store.runTaskNow(mem.realTask.id, 'live')
    }
  }
}

const stopSelected = async () => {
  for (const uid of selectedUids.value) {
    const mem = mergedAccounts.value.find(m => m.acc.uid === uid)
    if (mem && mem.realTask && mem.realTask.status === 'Running') {
      await store.pauseTask(mem.realTask.id)
    }
  }
}

async function fetchAccounts(silent = false) {
  const shouldShowLoading = !silent || loading.value
  if (shouldShowLoading) {
    loading.value = true
  }

  try {
    const data = await GetAllAccounts()
    accounts.value = data || []
    const availableUids = new Set(accounts.value.map((acc) => acc.uid))

    friendsByAccount.value = Object.fromEntries(
      Object.entries(friendsByAccount.value).filter(([uid]) => availableUids.has(uid))
    )
    loadingFriendMap.value = Object.fromEntries(
      Object.entries(loadingFriendMap.value).filter(([uid]) => availableUids.has(uid))
    )
    friendErrorMap.value = Object.fromEntries(
      Object.entries(friendErrorMap.value).filter(([uid]) => availableUids.has(uid))
    )

    if (selectedFriendUid.value && !availableUids.has(selectedFriendUid.value)) {
      selectedFriendUid.value = ''
    }

    if (!selectedFriendUid.value && accounts.value.length > 0 && activeSection.value === 'friends') {
      selectedFriendUid.value = accounts.value[0].uid
    }
  } catch (e: any) {
    console.error("Failed to load accounts:", e)
  } finally {
    if (shouldShowLoading) {
      loading.value = false
    }
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
    showAccountChoiceModal.value = false
    newAccName.value = ''
    newAccCookie.value = ''
  } catch (err: any) {
    addError.value = err.toString()
  } finally {
    isSubmitting.value = false
  }
}

async function submitBulkAdd() {
  if (!bulkCookies.value.trim()) {
    bulkResult.value = { error: "Vui lòng nhập danh sách cookie" }
    return
  }
  isSubmittingBulk.value = true
  bulkResult.value = null
  try {
    const result = await ImportMultipleAccounts(bulkCookies.value)
    bulkResult.value = result
    await fetchAccounts()
    if (result.failedCount === 0 && result.successCount > 0) {
      setTimeout(() => {
        showAccountChoiceModal.value = false
        bulkCookies.value = ''
        bulkResult.value = null
      }, 3000)
    }
  } catch (err: any) {
    bulkResult.value = { error: err.toString() }
  } finally {
    isSubmittingBulk.value = false
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


async function submitLogin() {
  if (!loginUsername.value.trim()) {
    loginError.value = 'Vui lòng nhập tài khoản (email hoặc số điện thoại)'
    return
  }
  if (!loginPassword.value) {
    loginError.value = 'Vui lòng nhập mật khẩu'
    return
  }

  isSubmittingLogin.value = true
  loginError.value = ''

  try {
    const result = await LoginWithPassword(
      loginUsername.value.trim(),
      loginPassword.value,
      loginTwoFactorAuth.value.trim()
    )

    if (result && result.status === 'Live' && result.uid) {
      // Đăng nhập thành công — cập nhật danh sách tài khoản
      await fetchAccounts()
      showAccountChoiceModal.value = false
      loginUsername.value = ''
      loginPassword.value = ''
      loginTwoFactorAuth.value = ''
    } else {
      // Các trường hợp đặc biệt
      const statusMsg: Record<string, string> = {
        'Checkpoint': 'Tài khoản bị Checkpoint. Hãy vào Facebook trên trình duyệt để xác minh bảo mật rồi thử lại.',
        'WrongPassword': 'Sai tài khoản hoặc mật khẩu. Kiểm tra lại thông tin đăng nhập.',
        'Failed': 'Đăng nhập không thành công. ' + (result?.rawResponse || ''),
      }
      loginError.value = statusMsg[result?.status] || 'Đăng nhập thất bại. Vui lòng thử lại.'
    }
  } catch (err: any) {
    // Backend trả lỗi qua error (go error truyền về JS)
    const msg = String(err)
    if (msg.includes('checkpoint') || msg.includes('Checkpoint')) {
      loginError.value = 'Tài khoản bị Checkpoint — vào Facebook trên trình duyệt để xác minh.'
    } else if (msg.includes('sai') || msg.includes('incorrect') || msg.includes('WrongPassword')) {
      loginError.value = 'Sai tài khoản hoặc mật khẩu. Hãy kiểm tra lại.'
    } else if (msg.includes('LSD') || msg.includes('publicKey') || msg.includes('bóc tách')) {
      loginError.value = 'Không thể kết nối đến Facebook. Kiểm tra mạng hoặc thử lại sau.'
    } else {
      loginError.value = msg || 'Đã xảy ra lỗi không xác định.'
    }
  } finally {
    isSubmittingLogin.value = false
  }
}

// Drag Task Modal Logic
const taskModalX = ref(300)
const taskModalY = ref(50)
let isDraggingTask = false
let startXTask = 0
let startYTask = 0
let initialXTask = 0
let initialYTask = 0

function startDragTaskModal(e: MouseEvent) {
  isDraggingTask = true
  startXTask = e.clientX
  startYTask = e.clientY
  initialXTask = taskModalX.value
  initialYTask = taskModalY.value
  document.addEventListener('mousemove', onDragTaskModal)
  document.addEventListener('mouseup', stopDragTaskModal)
}

function onDragTaskModal(e: MouseEvent) {
  if (!isDraggingTask) return
  const dx = e.clientX - startXTask
  const dy = e.clientY - startYTask
  taskModalX.value = initialXTask + dx
  taskModalY.value = initialYTask + dy
}

function stopDragTaskModal() {
  isDraggingTask = false
  document.removeEventListener('mousemove', onDragTaskModal)
  document.removeEventListener('mouseup', stopDragTaskModal)
}

onMounted(() => {
  void fetchAccounts()
  void store.fetchTasks()
  void store.fetchOverview()
  interval = setInterval(() => {
    void fetchAccounts(true)
    void store.fetchTasks()
    void store.fetchOverview()
  }, 2000)
})

onUnmounted(() => {
  clearInterval(interval)
})
</script>
