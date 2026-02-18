<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-bold">Docker Containers</h2>
      <button
        @click="fetchContainers"
        class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
      >
        Refresh
      </button>
    </div>

    <div v-if="loading" class="text-center py-8">Loading...</div>

    <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
      {{ error }}
    </div>

    <div v-else class="bg-white shadow-md rounded-lg overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Name
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Image
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Status
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Ports
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Actions
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="container in containers" :key="container.Id">
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
              {{ container.Names[0]?.replace('/', '') || 'N/A' }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ container.Image }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span
                :class="[
                  'px-2 inline-flex text-xs leading-5 font-semibold rounded-full',
                  container.State === 'running' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                ]"
              >
                {{ container.State }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ formatPorts(container.Ports) }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium space-x-2">
              <button
                v-if="container.State !== 'running'"
                @click="startContainer(container.Id)"
                class="text-green-600 hover:text-green-900"
              >
                Start
              </button>
              <button
                v-if="container.State === 'running'"
                @click="stopContainer(container.Id)"
                class="text-yellow-600 hover:text-yellow-900"
              >
                Stop
              </button>
              <button
                @click="restartContainer(container.Id)"
                class="text-blue-600 hover:text-blue-900"
              >
                Restart
              </button>
              <button
                @click="deleteContainer(container.Id)"
                class="text-red-600 hover:text-red-900"
              >
                Delete
              </button>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="containers.length === 0" class="text-center py-8 text-gray-500">
        No containers found
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const containers = ref([])
const loading = ref(false)
const error = ref('')

const fetchContainers = async () => {
  try {
    loading.value = true
    error.value = ''
    containers.value = await api.listContainers()
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to fetch containers'
  } finally {
    loading.value = false
  }
}

const startContainer = async (id) => {
  try {
    await api.startContainer(id)
    await fetchContainers()
  } catch (err) {
    alert('Failed to start container: ' + (err.response?.data?.error || err.message))
  }
}

const stopContainer = async (id) => {
  try {
    await api.stopContainer(id)
    await fetchContainers()
  } catch (err) {
    alert('Failed to stop container: ' + (err.response?.data?.error || err.message))
  }
}

const restartContainer = async (id) => {
  try {
    await api.restartContainer(id)
    await fetchContainers()
  } catch (err) {
    alert('Failed to restart container: ' + (err.response?.data?.error || err.message))
  }
}

const deleteContainer = async (id) => {
  if (!confirm('Are you sure you want to delete this container?')) return
  try {
    await api.deleteContainer(id)
    await fetchContainers()
  } catch (err) {
    alert('Failed to delete container: ' + (err.response?.data?.error || err.message))
  }
}

const formatPorts = (ports) => {
  if (!ports || ports.length === 0) return 'None'
  return ports.map(p => p.PublicPort ? `${p.PublicPort}:${p.PrivatePort}` : p.PrivatePort).join(', ')
}

onMounted(() => {
  fetchContainers()
})
</script>
