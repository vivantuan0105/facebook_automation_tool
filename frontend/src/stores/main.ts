import { defineStore } from 'pinia'

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
    overview: null as any,
    tasks: [] as any[],
    executions: [] as any[],
    logs: [] as any[],
    settings: null as any,
    loading: false
  }),
  actions: {
    async fetchOverview() {
      const res = await apiCall('GetOverview')
      if (res) this.overview = res
    },
    async fetchTasks() {
      const res = await apiCall('GetReactionTasks')
      if (res) this.tasks = res
    },
    async fetchExecutions() {
      const res = await apiCall('GetExecutions')
      if (res) this.executions = res
    },
    async fetchLogs() {
      const res = await apiCall('GetLogs')
      if (res) this.logs = res
    },
    async fetchSettings() {
      const res = await apiCall('GetSettings')
      if (res) this.settings = res
    },
    async validateTask(input: any) {
      return await apiCall('ValidateReactionTask', input)
    },
    async createTask(input: any) {
      await apiCall('CreateReactionTask', input)
      await this.fetchTasks()
      await this.fetchOverview()
    },
    async runTaskNow(id: number, mode: string) {
      await apiCall('RunReactionTaskNow', id, mode)
      await this.fetchTasks()
      await this.fetchExecutions()
    },
    async pauseTask(id: number) {
      await apiCall('PauseReactionTask', id)
      await this.fetchTasks()
    },
    async retryTask(id: number) {
      await apiCall('RetryReactionTask', id)
      await this.fetchTasks()
    },
    async removeTask(id: number) {
      await apiCall('RemoveReactionTask', id)
      await this.fetchTasks()
    },
    async markCompleted(id: number) {
      await apiCall('MarkManualCompleted', id)
      await this.fetchTasks()
      await this.fetchExecutions()
    },
    async updateSettings(s: any) {
      await apiCall('UpdateSettings', s)
      await this.fetchSettings()
    }
  }
})
