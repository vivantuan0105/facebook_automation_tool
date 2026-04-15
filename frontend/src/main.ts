import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './style.css'

const app = createApp(App)

app.config.errorHandler = (err, _vm, _info) => {
  console.error('Vue Error:', err, _info)
  alert(`Vue Error: ${err}\nInfo: ${_info}`)
}
window.addEventListener('error', (event) => {
  alert(`JS Error: ${event.message}`)
})
window.addEventListener('unhandledrejection', (event) => {
  alert(`Unhandled Promise Rejection: ${event.reason}`)
})

const pinia = createPinia()

app.use(pinia)
app.use(router)

app.mount('#app')
