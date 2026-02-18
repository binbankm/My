import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default {
  // Auth
  login: (data) => api.post('/auth/login', data),
  logout: () => api.post('/auth/logout'),
  getUserInfo: () => api.get('/auth/info'),

  // System
  getSystemInfo: () => api.get('/system/info'),
  getSystemStats: () => api.get('/system/stats'),

  // Docker
  listContainers: () => api.get('/docker/containers'),
  getContainer: (id) => api.get(`/docker/containers/${id}`),
  startContainer: (id) => api.post(`/docker/containers/${id}/start`),
  stopContainer: (id) => api.post(`/docker/containers/${id}/stop`),
  restartContainer: (id) => api.post(`/docker/containers/${id}/restart`),
  deleteContainer: (id) => api.delete(`/docker/containers/${id}`),
  listImages: () => api.get('/docker/images'),
  deleteImage: (id) => api.delete(`/docker/images/${id}`),

  // Files
  listFiles: (path) => api.get('/files', { params: { path } }),
  createFile: (data) => api.post('/files', data),
  updateFile: (data) => api.put('/files', data),
  deleteFile: (path) => api.delete('/files', { params: { path } }),
  downloadFile: (path) => `/api/files/download?path=${encodeURIComponent(path)}`,
  uploadFile: (formData) => api.post('/files/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  }),

  // Database
  listDatabases: () => api.get('/database'),
  createDatabase: (data) => api.post('/database', data),
  deleteDatabase: (name) => api.delete(`/database/${name}`),

  // Settings
  getSettings: () => api.get('/settings'),
  updateSettings: (data) => api.put('/settings', data),
}
