<template>
  <div>
    <h2 class="text-2xl font-bold mb-6">Settings</h2>

    <div class="space-y-6">
      <!-- System Information -->
      <div class="bg-white shadow-md rounded-lg p-6">
        <h3 class="text-lg font-medium mb-4">System Information</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700">Version</label>
            <p class="mt-1 text-sm text-gray-900">{{ settings.version || '1.0.0' }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700">Status</label>
            <p class="mt-1 text-sm text-green-600">Running</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700">Port</label>
            <p class="mt-1 text-sm text-gray-900">{{ settings.port || '8888' }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700">Environment</label>
            <p class="mt-1 text-sm text-gray-900">{{ settings.environment || 'production' }}</p>
          </div>
        </div>
      </div>

      <!-- Application Settings -->
      <div class="bg-white shadow-md rounded-lg p-6">
        <h3 class="text-lg font-medium mb-4">Application Settings</h3>
        
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Server Name
            </label>
            <input
              v-model="editableSettings.server_name"
              type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="My Server Panel"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Session Timeout (minutes)
            </label>
            <input
              v-model.number="editableSettings.session_timeout"
              type="number"
              min="5"
              max="1440"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="60"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Max Upload Size (MB)
            </label>
            <input
              v-model.number="editableSettings.max_upload_size"
              type="number"
              min="1"
              max="1024"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="100"
            />
          </div>

          <div class="flex items-center">
            <input
              v-model="editableSettings.enable_notifications"
              type="checkbox"
              id="notifications"
              class="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
            />
            <label for="notifications" class="ml-2 text-sm text-gray-700">
              Enable email notifications
            </label>
          </div>

          <div class="flex items-center">
            <input
              v-model="editableSettings.enable_auto_backup"
              type="checkbox"
              id="auto-backup"
              class="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
            />
            <label for="auto-backup" class="ml-2 text-sm text-gray-700">
              Enable automatic daily backups
            </label>
          </div>
        </div>

        <div class="mt-6 flex justify-end gap-2">
          <button
            @click="resetSettings"
            class="px-4 py-2 bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
          >
            Reset
          </button>
          <button
            @click="saveSettings"
            :disabled="saving"
            class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ saving ? 'Saving...' : 'Save Settings' }}
          </button>
        </div>
      </div>

      <!-- Security Settings -->
      <div class="bg-white shadow-md rounded-lg p-6">
        <h3 class="text-lg font-medium mb-4">Security</h3>
        
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Allowed IP Addresses
            </label>
            <textarea
              v-model="editableSettings.allowed_ips"
              rows="3"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono text-sm"
              placeholder="0.0.0.0/0 (allow all)&#10;192.168.1.0/24&#10;10.0.0.1"
            ></textarea>
            <p class="text-xs text-gray-500 mt-1">
              One IP/CIDR per line. Leave empty or use 0.0.0.0/0 to allow all.
            </p>
          </div>

          <div class="flex items-center">
            <input
              v-model="editableSettings.enforce_https"
              type="checkbox"
              id="https"
              class="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
            />
            <label for="https" class="ml-2 text-sm text-gray-700">
              Enforce HTTPS (redirect HTTP to HTTPS)
            </label>
          </div>

          <div class="flex items-center">
            <input
              v-model="editableSettings.two_factor_auth"
              type="checkbox"
              id="2fa"
              class="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
            />
            <label for="2fa" class="ml-2 text-sm text-gray-700">
              Require two-factor authentication
            </label>
          </div>
        </div>

        <div class="mt-6 flex justify-end">
          <button
            @click="saveSettings"
            :disabled="saving"
            class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ saving ? 'Saving...' : 'Save Security Settings' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const settings = ref({})
const editableSettings = ref({
  server_name: '',
  session_timeout: 60,
  max_upload_size: 100,
  enable_notifications: false,
  enable_auto_backup: true,
  allowed_ips: '',
  enforce_https: false,
  two_factor_auth: false,
})
const saving = ref(false)

const fetchSettings = async () => {
  try {
    const response = await api.getSettings()
    settings.value = response.data || response || {}
    
    // Populate editable settings
    editableSettings.value = {
      server_name: settings.value.server_name || '',
      session_timeout: settings.value.session_timeout || 60,
      max_upload_size: settings.value.max_upload_size || 100,
      enable_notifications: settings.value.enable_notifications || false,
      enable_auto_backup: settings.value.enable_auto_backup !== false,
      allowed_ips: settings.value.allowed_ips || '',
      enforce_https: settings.value.enforce_https || false,
      two_factor_auth: settings.value.two_factor_auth || false,
    }
  } catch (err) {
    console.error('Failed to fetch settings:', err)
  }
}

const saveSettings = async () => {
  saving.value = true
  try {
    await api.updateSettings(editableSettings.value)
    alert('Settings saved successfully!')
    await fetchSettings()
  } catch (err) {
    console.error('Failed to save settings:', err)
    alert('Failed to save settings: ' + (err.response?.data?.error || err.message))
  } finally {
    saving.value = false
  }
}

const resetSettings = () => {
  editableSettings.value = {
    server_name: settings.value.server_name || '',
    session_timeout: settings.value.session_timeout || 60,
    max_upload_size: settings.value.max_upload_size || 100,
    enable_notifications: settings.value.enable_notifications || false,
    enable_auto_backup: settings.value.enable_auto_backup !== false,
    allowed_ips: settings.value.allowed_ips || '',
    enforce_https: settings.value.enforce_https || false,
    two_factor_auth: settings.value.two_factor_auth || false,
  }
}

onMounted(() => {
  fetchSettings()
})
</script>
