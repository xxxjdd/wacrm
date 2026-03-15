import axios from 'axios'

const API_URL = 'https://api.dgxs.cn'

export const api = axios.create({
  baseURL: API_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add auth token to requests
api.interceptors.request.use((config) => {
  const stored = localStorage.getItem('wacrm-auth')
  if (stored) {
    const auth = JSON.parse(stored)
    if (auth.state?.token) {
      config.headers.Authorization = `Bearer ${auth.state.token}`
    }
  }
  return config
})

// Handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('wacrm-auth')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// Account API
export const accountsApi = {
  list: () => api.get('/api/accounts'),
  create: (data: any) => api.post('/api/accounts', data),
  delete: (id: number) => api.delete(`/api/accounts/${id}`),
  logout: (id: number) => api.post(`/api/accounts/${id}/logout`),
  getQR: (id: number) => api.get(`/api/accounts/${id}/qr`),
  verify: (id: number, sessionData: string) => api.post(`/api/accounts/${id}/verify`, { session_data: sessionData }),
}

// Customer API
export const customersApi = {
  list: (params?: any) => api.get('/api/customers', { params }),
  create: (data: any) => api.post('/api/customers', data),
  update: (id: number, data: any) => api.put(`/api/customers/${id}`, data),
  delete: (id: number) => api.delete(`/api/customers/${id}`),
  import: (data: any) => api.post('/api/customers/import', data),
}

// Message API
export const messagesApi = {
  list: (params?: any) => api.get('/api/messages', { params }),
  send: (data: any) => api.post('/api/messages/send', data),
  conversations: (params?: any) => api.get('/api/messages/conversations', { params }),
}

// Template API
export const templatesApi = {
  list: (params?: any) => api.get('/api/templates', { params }),
  create: (data: any) => api.post('/api/templates', data),
  update: (id: number, data: any) => api.put(`/api/templates/${id}`, data),
  delete: (id: number) => api.delete(`/api/templates/${id}`),
}

// Task API
export const tasksApi = {
  list: (params?: any) => api.get('/api/tasks', { params }),
  create: (data: any) => api.post('/api/tasks', data),
  update: (id: number, data: any) => api.put(`/api/tasks/${id}`, data),
  delete: (id: number) => api.delete(`/api/tasks/${id}`),
  run: (id: number) => api.post(`/api/tasks/${id}/run`),
}

// Stats API
export const statsApi = {
  overview: () => api.get('/api/stats/overview'),
  messages: (days?: number) => api.get('/api/stats/messages', { params: { days } }),
}
