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

  // Cron Jobs
  listCronJobs: () => api.get('/cron'),
  getCronJob: (id) => api.get(`/cron/${id}`),
  createCronJob: (data) => api.post('/cron', data),
  updateCronJob: (id, data) => api.put(`/cron/${id}`, data),
  deleteCronJob: (id) => api.delete(`/cron/${id}`),

  // Logs
  listLogFiles: () => api.get('/logs/files'),
  readLogFile: (path) => api.get('/logs/read', { params: { path } }),
  searchLogs: (query) => api.get('/logs/search', { params: { query } }),
  getSystemLogs: () => api.get('/logs/system'),
  downloadLogFile: (path) => `/api/logs/download?path=${encodeURIComponent(path)}`,
  clearLogFile: (path) => api.post('/logs/clear', { path }),
  getLogStats: () => api.get('/logs/stats'),

  // Nginx
  listNginxSites: () => api.get('/nginx/sites'),
  getNginxSite: (name) => api.get(`/nginx/sites/${name}`),
  createNginxSite: (data) => api.post('/nginx/sites', data),
  updateNginxSite: (name, data) => api.put(`/nginx/sites/${name}`, data),
  deleteNginxSite: (name) => api.delete(`/nginx/sites/${name}`),
  enableNginxSite: (name) => api.post(`/nginx/sites/${name}/enable`),
  disableNginxSite: (name) => api.post(`/nginx/sites/${name}/disable`),
  testNginxConfig: () => api.post('/nginx/test'),
  reloadNginx: () => api.post('/nginx/reload'),
  getNginxStatus: () => api.get('/nginx/status'),

  // Backups
  listBackups: () => api.get('/backup'),
  createBackup: (data) => api.post('/backup', data),
  downloadBackup: (id) => `/api/backup/${id}/download`,
  deleteBackup: (id) => api.delete(`/backup/${id}`),
  restoreBackup: (id) => api.post(`/backup/${id}/restore`),
  getBackupStats: () => api.get('/backup/stats'),

  // Users
  listUsers: () => api.get('/users'),
  getUser: (id) => api.get(`/users/${id}`),
  createUser: (data) => api.post('/users', data),
  updateUser: (id, data) => api.put(`/users/${id}`, data),
  deleteUser: (id) => api.delete(`/users/${id}`),

  // Roles
  listRoles: () => api.get('/roles'),
  getRole: (id) => api.get(`/roles/${id}`),
  createRole: (data) => api.post('/roles', data),
  updateRole: (id, data) => api.put(`/roles/${id}`, data),
  deleteRole: (id) => api.delete(`/roles/${id}`),

  // Permissions
  listPermissions: () => api.get('/permissions'),
}
