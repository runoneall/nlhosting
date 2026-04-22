import { createApp } from 'vue'
import App from '@/App.vue'
const app = createApp(App)

import { createPinia } from 'pinia'
app.use(createPinia())

import router from '@/router'
app.use(router)

import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
app.use(ElementPlus)

app.mount('#app')
