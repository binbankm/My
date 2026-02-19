<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-bold">System Logs</h2>
      <div class="flex gap-2">
        <button
          @click="fetchLogFiles"
          class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
        >
          Refresh
        </button>
      </div>
    </div>

    <!-- Log Files Sidebar and Viewer -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <!-- Log Files List -->
      <div class="md:col-span-1">
        <div class="bg-white shadow-md rounded-lg p-4">
          <h3 class="font-semibold mb-3">Log Files</h3>
          
          <div v-if="loadingFiles" class="text-center py-4">
            <p class="text-sm text-gray-600">Loading...</p>
          </div>

          <div v-else-if="logFiles.length === 0" class="text-center py-4">
            <p class="text-sm text-gray-600">No log files found</p>
          </div>

          <div v-else class="space-y-1">
            <button
              v-for="file in logFiles"
              :key="file.path"
              @click="selectLogFile(file)"
              :class="[
                'w-full text-left px-3 py-2 rounded text-sm transition-colors',
                selectedFile?.path === file.path
                  ? 'bg-blue-100 text-blue-700'
                  : 'hover:bg-gray-100'
              ]"
            >
              <div class="font-medium truncate">{{ file.name }}</div>
              <div class="text-xs text-gray-500">{{ formatSize(file.size) }}</div>
            </button>
          </div>

          <!-- System Logs Button -->
          <div class="mt-4 pt-4 border-t">
            <button
              @click="viewSystemLogs"
              :class="[
                'w-full text-left px-3 py-2 rounded text-sm transition-colors',
                viewingSystemLogs ? 'bg-blue-100 text-blue-700' : 'hover:bg-gray-100'
              ]"
            >
              <div class="font-medium">System Logs</div>
              <div class="text-xs text-gray-500">Live system logs</div>
            </button>
          </div>
        </div>

        <!-- Log Stats -->
        <div v-if="logStats" class="mt-4 bg-white shadow-md rounded-lg p-4">
          <h3 class="font-semibold mb-3 text-sm">Statistics</h3>
          <div class="space-y-2 text-xs">
            <div class="flex justify-between">
              <span class="text-gray-600">Total Files:</span>
              <span class="font-medium">{{ logStats.total_files || 0 }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600">Total Size:</span>
              <span class="font-medium">{{ formatSize(logStats.total_size || 0) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Log Content Viewer -->
      <div class="md:col-span-3">
        <div class="bg-white shadow-md rounded-lg p-4">
          <!-- Search Bar -->
          <div class="mb-4 flex gap-2">
            <input
              v-model="searchQuery"
              @keyup.enter="searchInLogs"
              type="text"
              placeholder="Search in logs..."
              class="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button
              @click="searchInLogs"
              class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              Search
            </button>
          </div>

          <!-- Log Actions -->
          <div v-if="selectedFile" class="mb-4 flex gap-2 text-sm">
            <a
              :href="getDownloadUrl(selectedFile.path)"
              class="px-3 py-1 bg-green-100 text-green-700 rounded hover:bg-green-200"
              download
            >
              Download
            </a>
            <button
              @click="clearLog"
              class="px-3 py-1 bg-red-100 text-red-700 rounded hover:bg-red-200"
            >
              Clear Log
            </button>
          </div>

          <!-- Log Content -->
          <div v-if="loadingContent" class="text-center py-12">
            <p class="text-gray-600">Loading log content...</p>
          </div>

          <div
            v-else-if="logContent"
            class="bg-gray-900 text-gray-100 p-4 rounded font-mono text-xs overflow-x-auto max-h-[600px] overflow-y-auto"
          >
            <pre class="whitespace-pre-wrap break-words">{{ logContent }}</pre>
          </div>

          <div v-else class="text-center py-12 text-gray-500">
            <p>Select a log file to view its content</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const logFiles = ref([])
const selectedFile = ref(null)
const logContent = ref('')
const searchQuery = ref('')
const loadingFiles = ref(false)
const loadingContent = ref(false)
const viewingSystemLogs = ref(false)
const logStats = ref(null)

const fetchLogFiles = async () => {
  loadingFiles.value = true
  try {
    const response = await api.listLogFiles()
    logFiles.value = response.data || []
  } catch (err) {
    console.error('Failed to fetch log files:', err)
    alert('Failed to load log files: ' + (err.response?.data?.error || err.message))
  } finally {
    loadingFiles.value = false
  }
}

const fetchLogStats = async () => {
  try {
    const response = await api.getLogStats()
    logStats.value = response.data
  } catch (err) {
    console.error('Failed to fetch log stats:', err)
  }
}

const selectLogFile = async (file) => {
  selectedFile.value = file
  viewingSystemLogs.value = false
  loadingContent.value = true
  
  try {
    const response = await api.readLogFile(file.path)
    logContent.value = response.data || response || ''
  } catch (err) {
    console.error('Failed to read log file:', err)
    logContent.value = `Error loading log file: ${err.response?.data?.error || err.message}`
  } finally {
    loadingContent.value = false
  }
}

const viewSystemLogs = async () => {
  selectedFile.value = null
  viewingSystemLogs.value = true
  loadingContent.value = true
  
  try {
    const response = await api.getSystemLogs()
    logContent.value = response.data || response || ''
  } catch (err) {
    console.error('Failed to fetch system logs:', err)
    logContent.value = `Error loading system logs: ${err.response?.data?.error || err.message}`
  } finally {
    loadingContent.value = false
  }
}

const searchInLogs = async () => {
  if (!searchQuery.value.trim()) {
    alert('Please enter a search query')
    return
  }

  loadingContent.value = true
  try {
    const response = await api.searchLogs(searchQuery.value)
    logContent.value = response.data || response || 'No results found'
    selectedFile.value = null
    viewingSystemLogs.value = false
  } catch (err) {
    console.error('Failed to search logs:', err)
    alert('Failed to search logs: ' + (err.response?.data?.error || err.message))
  } finally {
    loadingContent.value = false
  }
}

const clearLog = async () => {
  if (!selectedFile.value) return
  
  if (!confirm(`Are you sure you want to clear ${selectedFile.value.name}?`)) {
    return
  }
  
  try {
    await api.clearLogFile(selectedFile.value.path)
    logContent.value = ''
    alert('Log file cleared successfully')
  } catch (err) {
    console.error('Failed to clear log file:', err)
    alert('Failed to clear log file: ' + (err.response?.data?.error || err.message))
  }
}

const formatSize = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

const getDownloadUrl = (path) => {
  return `/api/logs/download?path=${encodeURIComponent(path)}`
}

onMounted(() => {
  fetchLogFiles()
  fetchLogStats()
})
</script>
