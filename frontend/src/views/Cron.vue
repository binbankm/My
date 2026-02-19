<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-bold">Cron Jobs</h2>
      <button
        @click="showCreateDialog = true"
        class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
      >
        + New Cron Job
      </button>
    </div>

    <!-- Cron Jobs List -->
    <div v-if="loading" class="text-center py-12">
      <p class="text-gray-600">Loading cron jobs...</p>
    </div>

    <div v-else-if="cronJobs.length === 0" class="text-center py-12 bg-gray-50 rounded-lg">
      <p class="text-gray-600">No cron jobs configured</p>
      <p class="text-sm text-gray-500 mt-2">Click "New Cron Job" to create one</p>
    </div>

    <div v-else class="space-y-4">
      <div
        v-for="job in cronJobs"
        :key="job.id"
        class="bg-white shadow-md rounded-lg p-4 border border-gray-200"
      >
        <div class="flex justify-between items-start">
          <div class="flex-1">
            <h3 class="text-lg font-semibold">{{ job.name }}</h3>
            <p class="text-sm text-gray-600 mt-1">{{ job.description }}</p>
            <div class="mt-2 flex items-center gap-4 text-sm">
              <span class="font-mono bg-gray-100 px-2 py-1 rounded">{{ job.schedule }}</span>
              <span :class="job.enabled ? 'text-green-600' : 'text-gray-400'">
                {{ job.enabled ? '● Enabled' : '○ Disabled' }}
              </span>
            </div>
            <div class="mt-2 text-xs text-gray-500">
              <p>Command: <code class="bg-gray-100 px-1 rounded">{{ job.command }}</code></p>
              <p v-if="job.last_run" class="mt-1">Last run: {{ formatDate(job.last_run) }}</p>
            </div>
          </div>
          <div class="flex gap-2">
            <button
              @click="editJob(job)"
              class="px-3 py-1 text-sm bg-gray-100 text-gray-700 rounded hover:bg-gray-200"
            >
              Edit
            </button>
            <button
              @click="deleteJob(job.id)"
              class="px-3 py-1 text-sm bg-red-100 text-red-700 rounded hover:bg-red-200"
            >
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit Dialog -->
    <div v-if="showCreateDialog || editingJob" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto">
        <h3 class="text-xl font-bold mb-4">{{ editingJob ? 'Edit' : 'Create' }} Cron Job</h3>
        
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Name</label>
            <input
              v-model="formData.name"
              type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Job name"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
            <input
              v-model="formData.description"
              type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Job description"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Schedule (Cron Expression)</label>
            <input
              v-model="formData.schedule"
              type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono"
              placeholder="0 0 * * *"
            />
            <p class="text-xs text-gray-500 mt-1">
              Format: minute hour day month weekday (e.g., "0 0 * * *" for daily at midnight)
            </p>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Command</label>
            <textarea
              v-model="formData.command"
              rows="3"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono"
              placeholder="/path/to/script.sh"
            ></textarea>
          </div>

          <div class="flex items-center">
            <input
              v-model="formData.enabled"
              type="checkbox"
              id="enabled"
              class="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
            />
            <label for="enabled" class="ml-2 text-sm text-gray-700">Enabled</label>
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
            @click="saveJob"
            :disabled="!formData.name || !formData.schedule || !formData.command"
            class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ editingJob ? 'Update' : 'Create' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const cronJobs = ref([])
const loading = ref(false)
const showCreateDialog = ref(false)
const editingJob = ref(null)
const formData = ref({
  name: '',
  description: '',
  schedule: '',
  command: '',
  enabled: true,
})

const fetchCronJobs = async () => {
  loading.value = true
  try {
    const response = await api.listCronJobs()
    cronJobs.value = response.data || []
  } catch (err) {
    console.error('Failed to fetch cron jobs:', err)
    alert('Failed to load cron jobs: ' + (err.response?.data?.error || err.message))
  } finally {
    loading.value = false
  }
}

const editJob = (job) => {
  editingJob.value = job
  formData.value = {
    name: job.name,
    description: job.description || '',
    schedule: job.schedule,
    command: job.command,
    enabled: job.enabled,
  }
}

const closeDialog = () => {
  showCreateDialog.value = false
  editingJob.value = null
  formData.value = {
    name: '',
    description: '',
    schedule: '',
    command: '',
    enabled: true,
  }
}

const saveJob = async () => {
  try {
    if (editingJob.value) {
      await api.updateCronJob(editingJob.value.id, formData.value)
    } else {
      await api.createCronJob(formData.value)
    }
    closeDialog()
    fetchCronJobs()
  } catch (err) {
    console.error('Failed to save cron job:', err)
    alert('Failed to save cron job: ' + (err.response?.data?.error || err.message))
  }
}

const deleteJob = async (id) => {
  if (!confirm('Are you sure you want to delete this cron job?')) {
    return
  }
  
  try {
    await api.deleteCronJob(id)
    fetchCronJobs()
  } catch (err) {
    console.error('Failed to delete cron job:', err)
    alert('Failed to delete cron job: ' + (err.response?.data?.error || err.message))
  }
}

const formatDate = (dateString) => {
  if (!dateString) return 'Never'
  return new Date(dateString).toLocaleString()
}

onMounted(() => {
  fetchCronJobs()
})
</script>
