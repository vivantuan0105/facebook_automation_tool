<template>
  <div class="space-y-6 max-w-4xl">
    <h2 class="text-xl font-semibold text-white">Settings</h2>

    <div class="bg-dark-surface border border-dark-border rounded-xl shadow-sm">
      <div class="p-6 border-b border-dark-border">
        <h3 class="text-lg font-medium text-white mb-1">General Application Settings</h3>
        <p class="text-sm text-dark-muted mb-6">Manage system-wide preferences</p>
        
        <div class="space-y-5" v-if="localSettings">
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
            <label class="text-sm font-medium text-white">Application Name</label>
            <div class="md:col-span-2">
              <input type="text" v-model="localSettings.app_name" class="w-full bg-dark-bg border border-dark-border text-white text-sm rounded-lg focus:ring-primary focus:border-primary block p-2.5 transition-colors">
            </div>
          </div>
          
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
            <label class="text-sm font-medium text-white">Theme</label>
            <div class="md:col-span-2">
              <select v-model="localSettings.theme" class="w-full bg-dark-bg border border-dark-border text-white text-sm rounded-lg focus:ring-primary focus:border-primary block p-2.5 transition-colors">
                <option value="dark">Dark Theme</option>
                <option value="light">Light Theme</option>
              </select>
            </div>
          </div>
          
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
            <label class="text-sm font-medium text-white">Log Level</label>
            <div class="md:col-span-2">
              <select v-model="localSettings.log_level" class="w-full bg-dark-bg border border-dark-border text-white text-sm rounded-lg focus:ring-primary focus:border-primary block p-2.5 transition-colors">
                <option value="debug">Debug</option>
                <option value="info">Info</option>
                <option value="warning">Warning</option>
                <option value="error">Error</option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
            <label class="text-sm font-medium text-white">Auto Start</label>
            <div class="md:col-span-2 flex items-center">
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="localSettings.auto_start" class="sr-only peer">
                <div class="w-11 h-6 bg-dark-border peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
                <span class="ml-3 text-sm font-medium text-dark-muted">Launch app on system startup</span>
              </label>
            </div>
          </div>
        </div>
        <div v-else class="text-dark-muted py-4">Loading settings...</div>
      </div>
      
      <div class="p-6 bg-dark-bg/50 rounded-b-xl flex justify-end">
        <button class="px-5 py-2.5 bg-primary hover:bg-primary-hover text-white rounded-lg text-sm font-medium transition-colors" @click="saveSettings">
          Save Changes
        </button>
      </div>
    </div>

    <!-- API Connection Section -->
    <div class="bg-dark-surface border border-dark-border rounded-xl shadow-sm">
      <div class="p-6">
        <h3 class="text-lg font-medium text-white mb-1">Official API Connection</h3>
        <p class="text-sm text-dark-muted mb-6">Manage official platform integration</p>
        
        <div class="p-5 border border-dark-border rounded-lg bg-dark-bg">
          <div class="flex items-center justify-between mb-4">
            <div>
              <p class="text-white font-medium">Status</p>
              <p class="text-sm text-dark-muted">{{ localSettings?.connection_status === 'connected' ? 'Connected successfully' : 'Not connected' }}</p>
            </div>
            <span v-if="localSettings?.connection_status === 'connected'" class="bg-green-500/10 text-green-400 text-xs font-medium px-2.5 py-0.5 rounded border border-green-500/20">Active Session</span>
            <span v-else class="bg-red-500/10 text-red-400 text-xs font-medium px-2.5 py-0.5 rounded border border-red-500/20">Disconnected</span>
          </div>
          
          <div v-if="localSettings?.connection_status === 'connected'" class="mb-5 p-4 border border-dark-border rounded bg-dark-surface">
            <p class="text-sm text-dark-muted mb-1">Connected Account</p>
            <p class="text-white font-medium">{{ localSettings?.connected_account }}</p>
          </div>
          
          <div class="flex space-x-3 mt-2">
            <button v-if="localSettings?.connection_status !== 'connected'" class="px-4 py-2 bg-white text-black hover:bg-gray-200 rounded text-sm font-medium transition-colors" @click="mockConnect">
              Connect via Official OAuth
            </button>
            <button v-else class="px-4 py-2 bg-red-500/10 hover:bg-red-500/20 border border-red-500/20 text-red-400 rounded text-sm font-medium transition-colors" @click="mockDisconnect">
              Disconnect Account
            </button>
          </div>
          <p class="text-xs text-dark-muted mt-4 italic">* This application only utilizes official APIs and standard authentication flows.</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useMainStore } from '../stores/main'

const store = useMainStore()
const localSettings = ref<any>(null)

onMounted(async () => {
  await store.fetchSettings()
  if (store.settings) {
    localSettings.value = JSON.parse(JSON.stringify(store.settings))
  }
})

const saveSettings = async () => {
  try {
    if ((window as any).go?.app?.App?.UpdateSettings) {
      await (window as any).go.app.App.UpdateSettings(localSettings.value)
      alert('Settings saved successfully.')
      await store.fetchSettings()
    } else {
      alert('Mock: Settings saved.')
    }
  } catch (e) {
    console.error(e)
  }
}

const mockConnect = async () => {
  try {
    if ((window as any).go?.app?.App?.ConnectAccount) {
      await (window as any).go.app.App.ConnectAccount('dummy_token')
      await store.fetchSettings()
      localSettings.value = JSON.parse(JSON.stringify(store.settings))
    } else {
      localSettings.value.connection_status = 'connected'
      localSettings.value.connected_account = 'mock_page_demo'
    }
  } catch (e) {
    console.error(e)
  }
}

const mockDisconnect = async () => {
  try {
    if ((window as any).go?.app?.App?.DisconnectAccount) {
      await (window as any).go.app.App.DisconnectAccount()
      await store.fetchSettings()
      localSettings.value = JSON.parse(JSON.stringify(store.settings))
    } else {
      localSettings.value.connection_status = 'disconnected'
      localSettings.value.connected_account = ''
    }
  } catch (e) {
    console.error(e)
  }
}
</script>
