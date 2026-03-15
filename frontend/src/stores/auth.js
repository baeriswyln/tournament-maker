import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api'
import {usersApi} from "../api/index.js";

export const useAuthStore = defineStore('auth', () => {
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))
  const accessToken = ref(localStorage.getItem('access_token') || null)
  const refreshToken = ref(localStorage.getItem('refresh_token') || null)

  const isAuthenticated = computed(() => !!accessToken.value)
  const isGuest = computed(() => isAuthenticated.value && !user.value?.username)
  const isAdmin = computed(() => isAuthenticated.value && user.value?.username)

  async function login(username, password) {
    const { data } = await authApi.login(username, password)
    _setTokens(data)
    const user_id = data.user_id

    if (data.access_token) {
      const {data} = await usersApi.get(user_id)
      setUser(data)
    }
  }

  async function register(payload) {
    const { data } = await authApi.register(payload)
    if (data.access_token) {
      _setTokens(data)
      localStorage.setItem('user', JSON.stringify({ id: data.id, visible_name: payload.visible_name, username: payload.username }))
      user.value = { id: data.id, visible_name: payload.visible_name, username: payload.username }
    }
    return data
  }

  function _setTokens(data) {
    accessToken.value = data.access_token
    refreshToken.value = data.refresh_token

    localStorage.setItem('access_token', data.access_token)
    if (data.refresh_token) localStorage.setItem('refresh_token', data.refresh_token)
  }

  function setUser(u) {
    user.value = u
    localStorage.setItem('user', JSON.stringify(u))
  }

  function logout() {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    localStorage.clear()
  }

  return { user, accessToken, refreshToken, isAuthenticated, isGuest, isAdmin, login, register, setUser, logout }
})
