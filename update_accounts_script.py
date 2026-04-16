import sys
import re

file_path = 'd:/thuctap/facebook_automation_tool/frontend/src/views/Accounts.vue'
with open(file_path, 'r', encoding='utf-8') as f:
    content = f.read()

# 1. Update imports
content = content.replace(
    '''import { GetAllAccounts, AddAccount, DeleteAccount, SelectPhotoDialog, ScanAccountData, GetAccountFriendsList, LoginWithPassword } from '../../wailsjs/go/app/App'''',
    '''import { GetAllAccounts, AddAccount, DeleteAccount, SelectPhotoDialog, ScanAccountData, GetAccountFriendsList, GetAccountPostsList, LoginWithPassword } from '../../wailsjs/go/app/App''''
)

# 2. Update activeSection type
content = content.replace(
    '''const activeSection = ref<'accounts' | 'friends'>('accounts')''',
    '''const activeSection = ref<'accounts' | 'friends' | 'posts'>('accounts')'''
)

# 3. Add Posts Logic after friends logic
posts_logic = '''
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
    // format: [12345] content text...
    const match = post.match(/^\[(.*?)\] (.*)$/)
    if (match) {
      return {
        index: index + 1,
        id: match[1],
        text: match[2]
      }
    }
    return {
      index: index + 1,
      id: '',
      text: post
    }
  })
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

const openPostUrl = (postId: string) => {
  if (postId) {
    // Basic formatting to a post URL
    const url = `https://www.facebook.com/${postId}`
    window.open(url, '_blank')
  }
}

const copyToClipboard = (text: string) => {
  navigator.clipboard.writeText(text)
  alert('Đã copy ID: ' + text)
}
'''

content = content.replace(
    '// Selected & Merged logic',
    posts_logic + '\n// Selected & Merged logic'
)

# 4. Also add "Quét bài viết" to TaskType options
content = content.replace(
    "{ value: 'Quét bạn bè', label: 'Quét danh sách Bạn bè (Scan Friends)' },",
    "{ value: 'Quét bạn bè', label: 'Quét danh sách Bạn bè (Scan Friends)' },\n  { value: 'Quét bài viết', label: 'Quét Timeline Bài viết (Scan Posts)' },"
)

# 5. Fix targetModeOptions setting logic
content = content.replace(
    "if (form.taskType === 'Quét thông tin' || form.taskType === 'Quét bạn bè') {",
    "if (form.taskType === 'Quét thông tin' || form.taskType === 'Quét bạn bè' || form.taskType === 'Quét bài viết') {"
)
# Note: targetModeOptions targetModeOptions setup handler in template change
content = content.replace(
    "if (form.taskType === 'Quét thông tin' || form.taskType === 'Quét bạn bè') return [{ value: 'self_profile', label: 'Quét danh sách bạn bè của tài khoản này' }]",
    "if (form.taskType === 'Quét thông tin' || form.taskType === 'Quét bạn bè' || form.taskType === 'Quét bài viết') return [{ value: 'self_profile', label: 'Dùng chính tài khoản này (Quét chính mình)' }, { value: 'post_url', label: 'Quét bài viết của một ID khác (Nhập UID)' }]"
)

with open(file_path, 'w', encoding='utf-8') as f:
    f.write(content)
