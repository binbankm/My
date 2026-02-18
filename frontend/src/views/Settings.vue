<template>
  <div>
    <h2 class="text-2xl font-bold mb-6">Settings</h2>

    <div class="bg-white shadow-md rounded-lg p-6">
      <div class="space-y-4">
        <div>
          <h3 class="text-lg font-medium mb-2">System Information</h3>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700">Version</label>
              <p class="mt-1 text-sm text-gray-900">1.0.0</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700">Status</label>
              <p class="mt-1 text-sm text-green-600">Running</p>
            </div>
          </div>
        </div>

        <hr />

        <div>
          <h3 class="text-lg font-medium mb-2">Application Settings</h3>
          <p class="text-sm text-gray-500">Configure your server panel settings here.</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const settings = ref([])

const fetchSettings = async () => {
  try {
    settings.value = await api.getSettings()
  } catch (err) {
    console.error('Failed to fetch settings:', err)
  }
}

onMounted(() => {
  fetchSettings()
})
</script>
