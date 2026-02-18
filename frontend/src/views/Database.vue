<template>
  <div>
    <h2 class="text-2xl font-bold mb-6">Database Management</h2>

    <div class="bg-white shadow-md rounded-lg p-6">
      <div class="mb-4">
        <button
          @click="fetchDatabases"
          class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
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
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Size</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="db in databases" :key="db.name">
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                {{ db.name }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ db.type }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ db.size }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                <button
                  @click="deleteDatabase(db.name)"
                  class="text-red-600 hover:text-red-900"
                >
                  Delete
                </button>
              </td>
            </tr>
          </tbody>
        </table>
        <div v-if="databases.length === 0" class="text-center py-8 text-gray-500">
          No databases found
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const databases = ref([])
const loading = ref(false)

const fetchDatabases = async () => {
  try {
    loading.value = true
    databases.value = await api.listDatabases()
  } catch (err) {
    alert('Failed to fetch databases: ' + (err.response?.data?.error || err.message))
  } finally {
    loading.value = false
  }
}

const deleteDatabase = async (name) => {
  if (!confirm(`Are you sure you want to delete database ${name}?`)) return
  try {
    await api.deleteDatabase(name)
    await fetchDatabases()
  } catch (err) {
    alert('Failed to delete database: ' + (err.response?.data?.error || err.message))
  }
}

onMounted(() => {
  fetchDatabases()
})
</script>
