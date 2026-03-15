import axios from 'axios'

const BASE_URL = import.meta.env.VITE_API_URL || 'https://stm.baeriswyl.dev/api'

const api = axios.create({
  baseURL: BASE_URL,
  headers: { 'Content-Type': 'application/json' }
})

// Attach token
api.interceptors.request.use(config => {
  const token = localStorage.getItem('access_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

// Auto-refresh on 401
api.interceptors.response.use(
  r => r,
  async error => {
    const original = error.config
    if (error.response?.status === 401 && !original._retry) {
      original._retry = true
      const refresh = localStorage.getItem('refresh_token')
      if (refresh) {
        try {
          const { data } = await axios.put(`${BASE_URL}/auth`, { refresh_token: refresh })
          localStorage.setItem('access_token', data.access_token)
          original.headers.Authorization = `Bearer ${data.access_token}`
          return api(original)
        } catch {
          localStorage.clear()
          window.location.href = '/login'
        }
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
