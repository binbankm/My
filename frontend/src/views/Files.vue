<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-bold">File Manager</h2>
    </div>

    <div class="bg-white shadow-md rounded-lg p-6">
      <!-- Breadcrumb Navigation -->
      <div class="mb-4 flex items-center space-x-2">
        <span class="text-sm text-gray-600">Path:</span>
        <nav class="flex flex-1" aria-label="Breadcrumb">
          <ol class="inline-flex items-center space-x-1">
            <li v-for="(part, index) in pathParts" :key="index" class="inline-flex items-center">
              <button
                v-if="index < pathParts.length - 1"
                @click="navigateToPart(index)"
                class="text-sm text-blue-600 hover:text-blue-900"
              >
                {{ part || 'root' }}
              </button>
              <span v-else class="text-sm font-medium text-gray-700">
                {{ part || 'root' }}
              </span>
              <span v-if="index < pathParts.length - 1" class="mx-1 text-gray-400">/</span>
            </li>
          </ol>
        </nav>
      </div>

      <!-- Action Buttons -->
      <div class="mb-4 flex flex-wrap gap-2">
        <button
          @click="navigateUp"
          :disabled="currentPath === '/'"
          class="px-4 py-2 bg-gray-200 text-gray-700 rounded hover:bg-gray-300 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          ⬆️ Up
        </button>
        <button
          @click="fetchFiles"
          class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
        >
          🔄 Refresh
        </button>
        <button
          @click="showCreateModal = true; createIsDir = true"
          class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
        >
          📁 New Folder
        </button>
        <button
          @click="showCreateModal = true; createIsDir = false"
          class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
        >
          📄 New File
        </button>
        <label class="px-4 py-2 bg-purple-600 text-white rounded hover:bg-purple-700 cursor-pointer">
          ⬆️ Upload File
          <input
            type="file"
            @change="handleFileUpload"
            class="hidden"
            ref="fileInput"
          />
        </label>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="text-center py-8">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-gray-300 border-t-blue-600"></div>
        <p class="mt-2 text-gray-600">Loading...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-8">
        <p class="text-red-600">❌ {{ error }}</p>
        <button
          @click="fetchFiles"
          class="mt-4 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
        >
          Try Again
        </button>
      </div>

      <!-- File List -->
      <div v-else>
        <div class="overflow-x-auto">
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
                  <span v-else class="text-gray-900">
                    {{ getFileIcon(file.name) }} {{ file.name }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {{ file.isDir ? '-' : formatBytes(file.size) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {{ file.modTime }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium space-x-2">
                  <button
                    v-if="!file.isDir && isEditableFile(file.name)"
                    @click="editFile(file)"
                    class="text-green-600 hover:text-green-900"
                  >
                    Edit
                  </button>
                  <button
                    v-if="!file.isDir"
                    @click="downloadFile(file.path)"
                    class="text-blue-600 hover:text-blue-900"
                  >
                    Download
                  </button>
                  <button
                    @click="deleteFile(file)"
                    class="text-red-600 hover:text-red-900"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        
        <!-- Empty State -->
        <div v-if="files.length === 0" class="text-center py-12">
          <div class="text-6xl mb-4">📂</div>
          <p class="text-gray-500 text-lg">This directory is empty</p>
          <p class="text-gray-400 text-sm mt-2">Create a new file or folder to get started</p>
        </div>
      </div>
    </div>

    <!-- Create File/Folder Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-full max-w-md">
        <h3 class="text-lg font-bold mb-4">
          {{ createIsDir ? 'Create New Folder' : 'Create New File' }}
        </h3>
        <input
          v-model="newName"
          type="text"
          :placeholder="createIsDir ? 'Folder name' : 'File name'"
          class="w-full px-3 py-2 border border-gray-300 rounded mb-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
          @keyup.enter="createFileOrFolder"
        />
        <div class="flex justify-end space-x-2">
          <button
            @click="showCreateModal = false; newName = ''"
            class="px-4 py-2 bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
          >
            Cancel
          </button>
          <button
            @click="createFileOrFolder"
            class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Create
          </button>
        </div>
      </div>
    </div>

    <!-- Edit File Modal -->
    <div v-if="showEditModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-full max-w-4xl max-h-[90vh] flex flex-col">
        <h3 class="text-lg font-bold mb-4">
          Edit File: {{ editingFile?.name }}
        </h3>
        <textarea
          v-model="fileContent"
          class="flex-1 w-full px-3 py-2 border border-gray-300 rounded mb-4 font-mono text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          rows="20"
        ></textarea>
        <div class="flex justify-end space-x-2">
          <button
            @click="showEditModal = false; editingFile = null; fileContent = ''"
            class="px-4 py-2 bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
          >
            Cancel
          </button>
          <button
            @click="saveFile"
            class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Save
          </button>
        </div>
      </div>
    </div>

    <!-- Upload Progress -->
    <div v-if="uploading" class="fixed bottom-4 right-4 bg-white rounded-lg shadow-lg p-4 w-80">
      <p class="font-medium mb-2">Uploading file...</p>
      <div class="w-full bg-gray-200 rounded-full h-2">
        <div class="bg-blue-600 h-2 rounded-full transition-all" :style="{ width: uploadProgress + '%' }"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '../api'

const currentPath = ref('/')
const files = ref([])
const loading = ref(false)
const error = ref(null)
const showCreateModal = ref(false)
const showEditModal = ref(false)
const createIsDir = ref(true)
const newName = ref('')
const editingFile = ref(null)
const fileContent = ref('')
const uploading = ref(false)
const uploadProgress = ref(0)
const fileInput = ref(null)

const pathParts = computed(() => {
  const parts = currentPath.value.split('/').filter(Boolean)
  return ['', ...parts]
})

const fetchFiles = async () => {
  try {
    loading.value = true
    error.value = null
    files.value = await api.listFiles(currentPath.value)
  } catch (err) {
    error.value = err.response?.data?.error || err.message || 'Failed to load files'
    files.value = []
  } finally {
    loading.value = false
  }
}

const navigateTo = (path) => {
  currentPath.value = path
  fetchFiles()
}

const navigateToPart = (index) => {
  if (index === 0) {
    currentPath.value = '/'
  } else {
    const parts = currentPath.value.split('/').filter(Boolean).slice(0, index)
    currentPath.value = '/' + parts.join('/')
  }
  fetchFiles()
}

const navigateUp = () => {
  if (currentPath.value === '/') return
  const parts = currentPath.value.split('/').filter(Boolean)
  parts.pop()
  currentPath.value = parts.length ? '/' + parts.join('/') : '/'
  fetchFiles()
}

const createFileOrFolder = async () => {
  if (!newName.value.trim()) {
    alert('Please enter a name')
    return
  }
  
  try {
    const path = currentPath.value === '/' 
      ? '/' + newName.value 
      : currentPath.value + '/' + newName.value
    
    await api.createFile({ path, isDir: createIsDir.value })
    showCreateModal.value = false
    newName.value = ''
    await fetchFiles()
  } catch (err) {
    alert('Failed to create: ' + (err.response?.data?.error || err.message))
  }
}

const editFile = async (file) => {
  try {
    editingFile.value = file
    // Try to fetch file content - we need to add this to the API
    const response = await fetch(api.downloadFile(file.path))
    const content = await response.text()
    fileContent.value = content
    showEditModal.value = true
  } catch (err) {
    alert('Failed to load file: ' + (err.response?.data?.error || err.message))
  }
}

const saveFile = async () => {
  try {
    await api.updateFile({
      path: editingFile.value.path,
      content: fileContent.value
    })
    showEditModal.value = false
    editingFile.value = null
    fileContent.value = ''
    await fetchFiles()
  } catch (err) {
    alert('Failed to save file: ' + (err.response?.data?.error || err.message))
  }
}

const handleFileUpload = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  try {
    uploading.value = true
    uploadProgress.value = 0
    
    const formData = new FormData()
    formData.append('file', file)
    formData.append('path', currentPath.value)

    await api.uploadFile(formData)
    
    uploadProgress.value = 100
    setTimeout(() => {
      uploading.value = false
      uploadProgress.value = 0
    }, 500)
    
    // Reset input
    if (fileInput.value) {
      fileInput.value.value = ''
    }
    
    await fetchFiles()
  } catch (err) {
    uploading.value = false
    uploadProgress.value = 0
    alert('Failed to upload file: ' + (err.response?.data?.error || err.message))
  }
}

const downloadFile = (path) => {
  window.open(api.downloadFile(path), '_blank')
}

const deleteFile = async (file) => {
  const type = file.isDir ? 'folder' : 'file'
  if (!confirm(`Are you sure you want to delete this ${type}?\n${file.name}`)) return
  
  try {
    await api.deleteFile(file.path)
    await fetchFiles()
  } catch (err) {
    alert('Failed to delete: ' + (err.response?.data?.error || err.message))
  }
}

const isEditableFile = (name) => {
  const editableExts = ['.txt', '.md', '.js', '.json', '.html', '.css', '.yml', '.yaml', 
                        '.conf', '.config', '.sh', '.env', '.xml', '.log', '.ini', '.vue',
                        '.go', '.py', '.php', '.java', '.c', '.cpp', '.h', '.ts', '.jsx', '.tsx']
  return editableExts.some(ext => name.toLowerCase().endsWith(ext))
}

const getFileIcon = (name) => {
  const ext = name.split('.').pop()?.toLowerCase()
  const iconMap = {
    js: '📜', ts: '📘', vue: '💚', html: '🌐', css: '🎨',
    json: '📋', xml: '📋', yml: '⚙️', yaml: '⚙️',
    md: '📝', txt: '📄', log: '📊',
    jpg: '🖼️', jpeg: '🖼️', png: '🖼️', gif: '🖼️', svg: '🎨',
    pdf: '📕', doc: '📘', docx: '📘',
    zip: '📦', tar: '📦', gz: '📦',
    sh: '⚡', py: '🐍', go: '🔵', php: '🐘',
  }
  return iconMap[ext] || '📄'
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
