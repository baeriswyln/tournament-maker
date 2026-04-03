import axios from 'axios'

const BASE_URL = import.meta.env.VITE_API_URL || 'https://stm.baeriswyl.dev/api'

const api = axios.create({
  baseURL: BASE_URL,
  headers: { 'Content-Type': 'application/json' }
})

// Decode JWT exp claim (seconds since epoch), returns null if unreadable
function tokenExp(token) {
  try {
    return JSON.parse(atob(token.split('.')[1])).exp ?? null
  } catch {
    return null
  }
}

// True when the token expires within the next 30 seconds
function expiringSoon(token) {
  const exp = tokenExp(token)
  return exp !== null && exp - Date.now() / 1000 < 30
}

// Singleton promise so concurrent requests share one refresh call
let refreshPromise = null

async function proactiveRefresh() {
  const refresh = localStorage.getItem('refresh_token')
  if (!refresh) return null
  if (!refreshPromise) {
    refreshPromise = axios.put(`${BASE_URL}/auth`, { refresh_token: refresh })
      .then(({ data }) => {
        localStorage.setItem('access_token', data.access_token)
        if (data.refresh_token) localStorage.setItem('refresh_token', data.refresh_token)
        return data.access_token
      })
      .catch(() => {
        localStorage.clear()
        window.dispatchEvent(new CustomEvent('session-expired'))
        window.location.href = '/login'
        return null
      })
      .finally(() => { refreshPromise = null })
  }
  return refreshPromise
}

// Attach token — refresh proactively if expiring soon
api.interceptors.request.use(async config => {
  let token = localStorage.getItem('access_token')
  if (token && expiringSoon(token)) {
    token = await proactiveRefresh() ?? token
  }
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

// Reactive refresh on 401 (fallback for clock skew / edge cases)
api.interceptors.response.use(
  r => r,
  async error => {
    const original = error.config
    if (error.response?.status === 401 && !original._retry) {
      original._retry = true
      const newToken = await proactiveRefresh()
      if (newToken) {
        original.headers.Authorization = `Bearer ${newToken}`
        return api(original)
      }
    }
    return Promise.reject(error)
  }
)

// ── Auth ──────────────────────────────────────────────
export const authApi = {
  login: (username, password) => api.post('/auth', { username, password }),
  refresh: (refresh_token) => api.put('/auth', { refresh_token }),
  register: (data) => api.post('/register', data),
}

// ── Users ─────────────────────────────────────────────
export const usersApi = {
  list: () => api.get('/users'),
  get: (id) => api.get(`/users/${id}`),
  create: (data) => api.post('/users', data),
  update: (id, data) => api.put(`/users/${id}`, data),
  delete: (id) => api.delete(`/users/${id}`),
}

// ── Tournaments ───────────────────────────────────────
export const tournamentsApi = {
  list: () => api.get('/tournaments'),
  get: (id) => api.get(`/tournaments/${id}`),
  create: (data) => api.post('/tournaments', data),
  update: (id, data) => api.put(`/tournaments/${id}`, data),
  delete: (id) => api.delete(`/tournaments/${id}`),
  reset: (id) => api.delete(`/tournaments/${id}/reset`),
  scoreboard: (id) => api.get(`/tournaments/${id}/scoreboard`),
}

// ── Participants ──────────────────────────────────────
export const participantsApi = {
  list: (tid) => api.get(`/tournaments/${tid}/participants`),
  add: (tid, ids) => api.post(`/tournaments/${tid}/participants`, ids),
  remove: (tid, id) => api.delete(`/tournaments/${tid}/participants`, { data: { id } }),
}

// ── Encounters ────────────────────────────────────────
export const encountersApi = {
  list: (tid, userId) => api.get(`/tournaments/${tid}/encounters`, { params: userId ? { user_id: userId } : {} }),
  create: (tid, encounters) => api.post(`/tournaments/${tid}/encounters`, encounters),
  updateScore: (tid, encounter) => api.put(`/tournaments/${tid}/encounters`, encounter),
  deleteScore: (tid, data) => api.delete(`/tournaments/${tid}/encounters`, { data }),
}

export default api
