// Theme must be applied before mount to avoid flash
import './composables/useTheme.js'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'
import Aura from '@primeuix/themes/aura';

import router from './router'
import App from './App.vue'

import 'primeicons/primeicons.css'
import './assets/main.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(PrimeVue, {
  theme: {
    preset: Aura,
    options: {
      darkModeSelector: 'html.dark',
      cssLayer: {
        name: 'primevue',
        order: 'base primevue, theme, components, utilities'
      }
    }
  }
})
app.use(ToastService)
app.use(ConfirmationService)

app.mount('#app')
