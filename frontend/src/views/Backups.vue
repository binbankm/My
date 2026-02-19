<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-bold">Backups</h2>
      <button
        @click="showCreateDialog = true"
        class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
      >
        + Create Backup
      </button>
    </div>

    <!-- Backup Stats -->
    <div v-if="backupStats" class="mb-6 grid grid-cols-1 md:grid-cols-3 gap-4">
      <div class="bg-white shadow-md rounded-lg p-4">
        <div class="text-sm text-gray-600">Total Backups</div>
        <div class="text-2xl font-bold mt-1">{{ backupStats.total_backups || 0 }}</div>
      </div>
      <div class="bg-white shadow-md rounded-lg p-4">
        <div class="text-sm text-gray-600">Total Size</div>
        <div class="text-2xl font-bold mt-1">{{ formatSize(backupStats.total_size || 0) }}</div>
      </div>
      <div class="bg-white shadow-md rounded-lg p-4">
        <div class="text-sm text-gray-600">Latest Backup</div>
        <div class="text-lg font-semibold mt-1">{{ formatDate(backupStats.latest_backup) }}</div>
      </div>
    </div>

    <!-- Backups List -->
    <div v-if="loading" class="text-center py-12">
      <p class="text-gray-600">Loading backups...</p>
    </div>

    <div v-else-if="backups.length === 0" class="text-center py-12 bg-gray-50 rounded-lg">
      <p class="text-gray-600">No backups found</p>
      <p class="text-sm text-gray-500 mt-2">Click "Create Backup" to create your first backup</p>
    </div>

    <div v-else class="space-y-4">
      <div
        v-for="backup in backups"
        :key="backup.id"
        class="bg-white shadow-md rounded-lg p-4 border border-gray-200"
      >
        <div class="flex justify-between items-start">
          <div class="flex-1">
            <h3 class="text-lg font-semibold">{{ backup.name }}</h3>
            <p v-if="backup.description" class="text-sm text-gray-600 mt-1">{{ backup.description }}</p>
            <div class="mt-2 flex items-center gap-4 text-sm">
              <span class="text-gray-600">
                <span class="font-medium">Size:</span> {{ formatSize(backup.size) }}
              </span>
              <span class="text-gray-600">
                <span class="font-medium">Created:</span> {{ formatDate(backup.created_at) }}
              </span>
            </div>
            <div v-if="backup.path" class="mt-2 text-xs text-gray-500">
              <code class="bg-gray-100 px-2 py-1 rounded">{{ backup.path }}</code>
            </div>
          </div>
          <div class="flex gap-2">
            <a
              :href="api.downloadBackup(backup.id)"
              class="px-3 py-1 text-sm bg-green-100 text-green-700 rounded hover:bg-green-200"
              download
            >
              Download
            </a>
            <button
              @click="restoreBackup(backup)"
              class="px-3 py-1 text-sm bg-blue-100 text-blue-700 rounded hover:bg-blue-200"
            >
              Restore
            </button>
            <button
              @click="deleteBackup(backup.id)"
              class="px-3 py-1 text-sm bg-red-100 text-red-700 rounded hover:bg-red-200"
            >
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Dialog -->
    <div v-if="showCreateDialog" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-full max-w-2xl">
        <h3 class="text-xl font-bold mb-4">Create Backup</h3>
        
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Backup Name</label>
            <input
              v-model="formData.name"
              type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="backup-2024-01-15"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
            <input
              v-model="formData.description"
              type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Before system upgrade"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Backup Type</label>
            <select
              v-model="formData.type"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="full">Full System Backup</option>
              <option value="database">Database Only</option>
              <option value="files">Files Only</option>
              <option value="config">Configuration Only</option>
            </select>
          </div>

          <div v-if="formData.type === 'files'">
            <label class="block text-sm font-medium text-gray-700 mb-1">Paths to Backup</label>
            <textarea
              v-model="formData.paths"
              rows="3"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono"
              placeholder="/var/www&#10;/etc/nginx&#10;/home/user"
            ></textarea>
            <p class="text-xs text-gray-500 mt-1">One path per line</p>
          </div>
        </div>

        <div class="mt-6 flex justify-end gap-2">
          <button
            @click="closeDialog"
            class="px-4 py-2 bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
          >
            Cancel
          </button>
          <button
            @click="createBackup"
            :disabled="!formData.name || creating"
            class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ creating ? 'Creating...' : 'Create Backup' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const backups = ref([])
const loading = ref(false)
const creating = ref(false)
const showCreateDialog = ref(false)
const backupStats = ref(null)
const formData = ref({
  name: '',
  description: '',
  type: 'full',
  paths: '',
})

const fetchBackups = async () => {
  loading.value = true
  try {
    const response = await api.listBackups()
    backups.value = response.data || []
  } catch (err) {
    console.error('Failed to fetch backups:', err)
    alert('Failed to load backups: ' + (err.response?.data?.error || err.message))
  } finally {
    loading.value = false
  }
}

const fetchBackupStats = async () => {
  try {
    const response = await api.getBackupStats()
    backupStats.value = response.data || response
  } catch (err) {
    console.error('Failed to fetch backup stats:', err)
  }
}

const closeDialog = () => {
  showCreateDialog.value = false
  formData.value = {
    name: '',
    description: '',
    type: 'full',
    paths: '',
  }
}

const createBackup = async () => {
  creating.value = true
  try {
    const data = { ...formData.value }
    if (data.type === 'files' && data.paths) {
      data.paths = data.paths.split('\n').filter(p => p.trim())
    }
    
    await api.createBackup(data)
    closeDialog()
    alert('Backup created successfully!')
    fetchBackups()
    fetchBackupStats()
  } catch (err) {
    console.error('Failed to create backup:', err)
    alert('Failed to create backup: ' + (err.response?.data?.error || err.message))
  } finally {
    creating.value = false
  }
}

const restoreBackup = async (backup) => {
  if (!confirm(`Are you sure you want to restore backup "${backup.name}"? This will overwrite current data.`)) {
    return
  }
  
  try {
    await api.restoreBackup(backup.id)
    alert('Backup restored successfully! Please restart the application.')
  } catch (err) {
    console.error('Failed to restore backup:', err)
    alert('Failed to restore backup: ' + (err.response?.data?.error || err.message))
  }
}

const deleteBackup = async (id) => {
  if (!confirm('Are you sure you want to delete this backup?')) {
    return
  }
  
  try {
    await api.deleteBackup(id)
    fetchBackups()
    fetchBackupStats()
  } catch (err) {
    console.error('Failed to delete backup:', err)
    alert('Failed to delete backup: ' + (err.response?.data?.error || err.message))
  }
}

const formatSize = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

const formatDate = (dateString) => {
  if (!dateString) return 'Never'
  return new Date(dateString).toLocaleString()
}

onMounted(() => {
  fetchBackups()
  fetchBackupStats()
})
</script>
