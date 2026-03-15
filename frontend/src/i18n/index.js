import { ref, computed } from 'vue'
import en from './en.js'
import fr from './fr.js'
import de from './de.js'

const locales = { en, fr, de }
export const availableLocales = [
  { code: 'en', label: 'English', flag: '🇬🇧' },
  { code: 'fr', label: 'Français', flag: '🇫🇷' },
  { code: 'de', label: 'Deutsch', flag: '🇩🇪' },
]

// Detect browser/saved language
function detectLocale() {
  const saved = localStorage.getItem('locale')
  if (saved && locales[saved]) return saved
  const browser = navigator.language?.slice(0, 2)
  if (browser && locales[browser]) return browser
  return 'en'
}

const currentLocale = ref(detectLocale())

export function useI18n() {
  function t(key, params = {}) {
    const messages = locales[currentLocale.value] || locales.en
    let str = messages[key] ?? locales.en[key] ?? key
    // Replace {param} placeholders
    Object.entries(params).forEach(([k, v]) => {
      str = str.replace(new RegExp(`\\{${k}\\}`, 'g'), v)
    })
    return str
  }

  function setLocale(code) {
    if (!locales[code]) return
    currentLocale.value = code
    localStorage.setItem('locale', code)
    document.documentElement.lang = code
  }

  const locale = computed(() => currentLocale.value)

  return { t, setLocale, locale, availableLocales }
}

// Set initial lang attr
document.documentElement.lang = currentLocale.value
