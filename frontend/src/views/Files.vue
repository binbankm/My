<template>
  <div>
    <h2 class="text-2xl font-bold mb-6">File Manager</h2>

    <div class="bg-white shadow-md rounded-lg p-6">
      <div class="mb-4 flex items-center space-x-2">
        <span class="text-sm text-gray-600">Current Path:</span>
        <span class="text-sm font-mono">{{ currentPath }}</span>
        <button
          @click="navigateUp"
          class="ml-auto px-3 py-1 bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
        >
          Up
        </button>
        <button
          @click="fetchFiles"
          class="px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700"
        >
          Refresh
        </button>
      </div>

      <div v-if="loading" class="text-center py-8">Loading...</div>

      <div v-else>
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Size</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Modified</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="file in files" :key="file.path" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap">
                <button
                  v-if="file.isDir"
                  @click="navigateTo(file.path)"
                  class="text-blue-600 hover:text-blue-900 font-medium"
                >
                  📁 {{ file.name }}
                </button>
                <span v-else class="text-gray-900">📄 {{ file.name }}</span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ file.isDir ? '-' : formatBytes(file.size) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ file.modTime }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                <button
                  v-if="!file.isDir"
                  @click="downloadFile(file.path)"
                  class="text-blue-600 hover:text-blue-900 mr-3"
                >
                  Download
                </button>
                <button
                  @click="deleteFile(file.path)"
                  class="text-red-600 hover:text-red-900"
                >
                  Delete
                </button>
              </td>
            </tr>
          </tbody>
        </table>
        <div v-if="files.length === 0" class="text-center py-8 text-gray-500">
          Directory is empty
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const currentPath = ref('/')
const files = ref([])
const loading = ref(false)

const fetchFiles = async () => {
  try {
    loading.value = true
    files.value = await api.listFiles(currentPath.value)
  } catch (err) {
    alert('Failed to fetch files: ' + (err.response?.data?.error || err.message))
  } finally {
    loading.value = false
  }
}

const navigateTo = (path) => {
  currentPath.value = path
  fetchFiles()
}

const navigateUp = () => {
  const parts = currentPath.value.split('/').filter(Boolean)
  parts.pop()
  currentPath.value = parts.length ? '/' + parts.join('/') : '/'
  fetchFiles()
}

const downloadFile = (path) => {
  window.open(api.downloadFile(path), '_blank')
}

const deleteFile = async (path) => {
  if (!confirm(`Are you sure you want to delete ${path}?`)) return
  try {
    await api.deleteFile(path)
    await fetchFiles()
  } catch (err) {
    alert('Failed to delete file: ' + (err.response?.data?.error || err.message))
  }
}

const formatBytes = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

onMounted(() => {
  fetchFiles()
})
</script>
