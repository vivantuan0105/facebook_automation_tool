import { defineStore } from 'pinia'

// Call wails bindings by window to bypass TS until 'wails dev' creates types.
const apiCall = async (method: string, ...args: any[]) => {
  try {
    if ((window as any).go?.app?.App?.[method]) {
      return await (window as any).go.app.App[method](...args)
    }
    console.warn(`Wails API ${method} not injected yet. Mocking.`)
    return null;
  } catch (err) {
    console.error(err)
    throw err;
  }
}

export const useMainStore = defineStore('main', {
  state: () => ({
    stats: null as any,
    posts: [] as any[],
    tasks: [] as any[],
    logs: [] as any[],
    settings: null as any,
    loading: false
  }),
  actions: {
    async fetchDashboardStats() {
      const res = await apiCall('GetDashboardStats')
      if (res) this.stats = res
    },
    async fetchPosts() {
      const res = await apiCall('GetPosts')
      if (res) this.posts = res
    },
    async fetchTasks() {
      const res = await apiCall('GetTasks')
      if (res) this.tasks = res
    },
    async fetchLogs() {
      const res = await apiCall('GetLogs')
      if (res) this.logs = res
    },
    async fetchSettings() {
      const res = await apiCall('GetSettings')
      if (res) this.settings = res
    }
  }
})
