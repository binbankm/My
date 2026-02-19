<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-bold">Nginx Management</h2>
      <div class="flex gap-2">
        <button
          @click="testConfig"
          class="px-4 py-2 bg-yellow-600 text-white rounded hover:bg-yellow-700"
        >
          Test Config
        </button>
        <button
          @click="reloadNginx"
          class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
        >
          Reload Nginx
        </button>
        <button
          @click="showCreateDialog = true"
          class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
        >
          + New Site
        </button>
      </div>
    </div>

    <!-- Nginx Status -->
    <div v-if="nginxStatus" class="mb-4 bg-white shadow-md rounded-lg p-4">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="font-semibold">Nginx Status</h3>
          <p class="text-sm text-gray-600 mt-1">{{ nginxStatus.version || 'Unknown version' }}</p>
        </div>
        <div :class="nginxStatus.running ? 'text-green-600' : 'text-red-600'">
          {{ nginxStatus.running ? '● Running' : '○ Stopped' }}
        </div>
      </div>
    </div>

    <!-- Sites List -->
    <div v-if="loading" class="text-center py-12">
      <p class="text-gray-600">Loading sites...</p>
    </div>

    <div v-else-if="sites.length === 0" class="text-center py-12 bg-gray-50 rounded-lg">
      <p class="text-gray-600">No Nginx sites configured</p>
      <p class="text-sm text-gray-500 mt-2">Click "New Site" to create one</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="site in sites"
        :key="site.name"
        class="bg-white shadow-md rounded-lg p-4 border border-gray-200"
      >
        <div class="flex justify-between items-start mb-3">
          <h3 class="text-lg font-semibold">{{ site.name }}</h3>
          <span
            :class="[
              'px-2 py-1 rounded text-xs font-medium',
              site.enabled ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
            ]"
          >
            {{ site.enabled ? 'Enabled' : 'Disabled' }}
          </span>
        </div>

        <div class="space-y-2 text-sm mb-4">
          <div>
            <span class="text-gray-600">Domain:</span>
            <span class="ml-2 font-mono">{{ site.server_name || 'N/A' }}</span>
          </div>
          <div v-if="site.listen">
            <span class="text-gray-600">Port:</span>
            <span class="ml-2 font-mono">{{ site.listen }}</span>
          </div>
          <div v-if="site.root">
            <span class="text-gray-600">Root:</span>
            <span class="ml-2 font-mono text-xs">{{ site.root }}</span>
          </div>
        </div>

        <div class="flex gap-2">
          <button
            v-if="!site.enabled"
            @click="enableSite(site.name)"
            class="flex-1 px-3 py-1 text-sm bg-green-100 text-green-700 rounded hover:bg-green-200"
          >
            Enable
          </button>
          <button
            v-else
            @click="disableSite(site.name)"
            class="flex-1 px-3 py-1 text-sm bg-yellow-100 text-yellow-700 rounded hover:bg-yellow-200"
          >
            Disable
          </button>
          <button
            @click="editSite(site)"
            class="flex-1 px-3 py-1 text-sm bg-gray-100 text-gray-700 rounded hover:bg-gray-200"
          >
            Edit
          </button>
          <button
            @click="deleteSite(site.name)"
            class="px-3 py-1 text-sm bg-red-100 text-red-700 rounded hover:bg-red-200"
          >
            Delete
          </button>
        </div>
      </div>
    </div>

    <!-- Create/Edit Dialog -->
    <div v-if="showCreateDialog || editingSite" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-full max-w-3xl max-h-[90vh] overflow-y-auto">
        <h3 class="text-xl font-bold mb-4">{{ editingSite ? 'Edit' : 'Create' }} Nginx Site</h3>
        
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Site Name</label>
            <input
              v-model="formData.name"
              :disabled="!!editingSite"
              type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100"
              placeholder="example.com"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Server Name (Domain)</label>
            <input
              v-model="formData.server_name"
              type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="example.com www.example.com"
            />
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Listen Port</label>
              <input
                v-model="formData.listen"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="80"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Root Directory</label>
              <input
                v-model="formData.root"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="/var/www/html"
              />
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Configuration</label>
            <textarea
              v-model="formData.config"
              rows="12"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono text-sm"
              placeholder="location / {&#10;    try_files $uri $uri/ =404;&#10;}"
            ></textarea>
            <p class="text-xs text-gray-500 mt-1">
              Enter the server block configuration (without server { })
            </p>
          </div>

          <div class="flex items-center">
            <input
              v-model="formData.enabled"
              type="checkbox"
              id="enabled"
              class="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
            />
            <label for="enabled" class="ml-2 text-sm text-gray-700">Enable site immediately</label>
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
            @click="saveSite"
            :disabled="!formData.name || !formData.server_name"
            class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ editingSite ? 'Update' : 'Create' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const sites = ref([])
const loading = ref(false)
const showCreateDialog = ref(false)
const editingSite = ref(null)
const nginxStatus = ref(null)
const formData = ref({
  name: '',
  server_name: '',
  listen: '80',
  root: '/var/www/html',
  config: '',
  enabled: true,
})

const fetchSites = async () => {
  loading.value = true
  try {
    const response = await api.listNginxSites()
    sites.value = response.data || []
  } catch (err) {
    console.error('Failed to fetch sites:', err)
    alert('Failed to load Nginx sites: ' + (err.response?.data?.error || err.message))
  } finally {
    loading.value = false
  }
}

const fetchNginxStatus = async () => {
  try {
    const response = await api.getNginxStatus()
    nginxStatus.value = response.data || response
  } catch (err) {
    console.error('Failed to fetch Nginx status:', err)
  }
}

const editSite = async (site) => {
  try {
    const response = await api.getNginxSite(site.name)
    const siteData = response.data || response
    editingSite.value = site
    formData.value = {
      name: siteData.name || site.name,
      server_name: siteData.server_name || '',
      listen: siteData.listen || '80',
      root: siteData.root || '/var/www/html',
      config: siteData.config || '',
      enabled: siteData.enabled !== undefined ? siteData.enabled : site.enabled,
    }
  } catch (err) {
    console.error('Failed to fetch site details:', err)
    alert('Failed to load site details: ' + (err.response?.data?.error || err.message))
  }
}

const closeDialog = () => {
  showCreateDialog.value = false
  editingSite.value = null
  formData.value = {
    name: '',
    server_name: '',
    listen: '80',
    root: '/var/www/html',
    config: '',
    enabled: true,
  }
}

const saveSite = async () => {
  try {
    if (editingSite.value) {
      await api.updateNginxSite(editingSite.value.name, formData.value)
    } else {
      await api.createNginxSite(formData.value)
    }
    closeDialog()
    fetchSites()
  } catch (err) {
    console.error('Failed to save site:', err)
    alert('Failed to save site: ' + (err.response?.data?.error || err.message))
  }
}

const deleteSite = async (name) => {
  if (!confirm(`Are you sure you want to delete site "${name}"?`)) {
    return
  }
  
  try {
    await api.deleteNginxSite(name)
    fetchSites()
  } catch (err) {
    console.error('Failed to delete site:', err)
    alert('Failed to delete site: ' + (err.response?.data?.error || err.message))
  }
}

const enableSite = async (name) => {
  try {
    await api.enableNginxSite(name)
    fetchSites()
  } catch (err) {
    console.error('Failed to enable site:', err)
    alert('Failed to enable site: ' + (err.response?.data?.error || err.message))
  }
}

const disableSite = async (name) => {
  try {
    await api.disableNginxSite(name)
    fetchSites()
  } catch (err) {
    console.error('Failed to disable site:', err)
    alert('Failed to disable site: ' + (err.response?.data?.error || err.message))
  }
}

const testConfig = async () => {
  try {
    const response = await api.testNginxConfig()
    alert(response.message || 'Configuration test successful!')
  } catch (err) {
    console.error('Failed to test config:', err)
    alert('Configuration test failed: ' + (err.response?.data?.error || err.message))
  }
}

const reloadNginx = async () => {
  try {
    const response = await api.reloadNginx()
    alert(response.message || 'Nginx reloaded successfully!')
    fetchNginxStatus()
  } catch (err) {
    console.error('Failed to reload Nginx:', err)
    alert('Failed to reload Nginx: ' + (err.response?.data?.error || err.message))
  }
}

onMounted(() => {
  fetchSites()
  fetchNginxStatus()
})
</script>
