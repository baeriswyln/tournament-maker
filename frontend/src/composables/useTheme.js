import { ref, watchEffect } from 'vue'

const saved = localStorage.getItem('theme')
const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
const theme = ref(saved || (prefersDark ? 'dark' : 'light'))

function applyTheme(t) {
  const html = document.documentElement
  if (t === 'dark') {
    html.classList.add('dark')
    html.setAttribute('data-theme', 'dark')
  } else {
    html.classList.remove('dark')
    html.setAttribute('data-theme', 'light')
  }
}

// Apply immediately
applyTheme(theme.value)

export function useTheme() {
  function setTheme(t) {
    theme.value = t
    localStorage.setItem('theme', t)
    applyTheme(t)
  }

  function toggle() {
    setTheme(theme.value === 'dark' ? 'light' : 'dark')
  }

  return { theme, setTheme, toggle }
}
