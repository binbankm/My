<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-bold">User Management</h2>
      <button
        @click="showCreateDialog = true"
        class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
      >
        + New User
      </button>
    </div>

    <!-- Users List -->
    <div v-if="loading" class="text-center py-12">
      <p class="text-gray-600">Loading users...</p>
    </div>

    <div v-else-if="users.length === 0" class="text-center py-12 bg-gray-50 rounded-lg">
      <p class="text-gray-600">No users found</p>
    </div>

    <div v-else class="bg-white shadow-md rounded-lg overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Username
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Email
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Role
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Status
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Created
            </th>
            <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Actions
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="user in users" :key="user.id">
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm font-medium text-gray-900">{{ user.username }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm text-gray-500">{{ user.email || 'N/A' }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span class="px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full bg-blue-100 text-blue-800">
                {{ user.role || 'User' }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span
                :class="[
                  'px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full',
                  user.active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                ]"
              >
                {{ user.active ? 'Active' : 'Inactive' }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ formatDate(user.created_at) }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <button
                @click="editUser(user)"
                class="text-blue-600 hover:text-blue-900 mr-3"
              >
                Edit
              </button>
              <button
                @click="deleteUser(user.id)"
                class="text-red-600 hover:text-red-900"
              >
                Delete
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create/Edit Dialog -->
    <div v-if="showCreateDialog || editingUser" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-full max-w-2xl">
        <h3 class="text-xl font-bold mb-4">{{ editingUser ? 'Edit' : 'Create' }} User</h3>
        
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Username</label>
            <input
              v-model="formData.username"
              :disabled="!!editingUser"
              type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100"
              placeholder="username"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
            <input
              v-model="formData.email"
              type="email"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="user@example.com"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
            <input
              v-model="formData.password"
              type="password"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              :placeholder="editingUser ? 'Leave empty to keep current password' : 'Password'"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Role</label>
            <select
              v-model="formData.role"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option v-for="role in roles" :key="role.id" :value="role.name">
                {{ role.name }}
              </option>
            </select>
          </div>

          <div class="flex items-center">
            <input
              v-model="formData.active"
              type="checkbox"
              id="active"
              class="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
            />
            <label for="active" class="ml-2 text-sm text-gray-700">Active</label>
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
            @click="saveUser"
            :disabled="!formData.username || (!editingUser && !formData.password)"
            class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ editingUser ? 'Update' : 'Create' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const users = ref([])
const roles = ref([])
const loading = ref(false)
const showCreateDialog = ref(false)
const editingUser = ref(null)
const formData = ref({
  username: '',
  email: '',
  password: '',
  role: 'User',
  active: true,
})

const fetchUsers = async () => {
  loading.value = true
  try {
    const response = await api.listUsers()
    users.value = response.data || []
  } catch (err) {
    console.error('Failed to fetch users:', err)
    alert('Failed to load users: ' + (err.response?.data?.error || err.message))
  } finally {
    loading.value = false
  }
}

const fetchRoles = async () => {
  try {
    const response = await api.listRoles()
    roles.value = response.data || []
    if (roles.value.length === 0) {
      // Default roles if none exist
      roles.value = [
        { id: 1, name: 'Admin' },
        { id: 2, name: 'User' },
        { id: 3, name: 'Viewer' }
      ]
    }
  } catch (err) {
    console.error('Failed to fetch roles:', err)
    // Set default roles on error
    roles.value = [
      { id: 1, name: 'Admin' },
      { id: 2, name: 'User' },
      { id: 3, name: 'Viewer' }
    ]
  }
}

const editUser = (user) => {
  editingUser.value = user
  formData.value = {
    username: user.username,
    email: user.email || '',
    password: '',
    role: user.role || 'User',
    active: user.active !== undefined ? user.active : true,
  }
}

const closeDialog = () => {
  showCreateDialog.value = false
  editingUser.value = null
  formData.value = {
    username: '',
    email: '',
    password: '',
    role: 'User',
    active: true,
  }
}

const saveUser = async () => {
  try {
    const data = { ...formData.value }
    // Don't send password if empty when editing
    if (editingUser.value && !data.password) {
      delete data.password
    }
    
    if (editingUser.value) {
      await api.updateUser(editingUser.value.id, data)
    } else {
      await api.createUser(data)
    }
    closeDialog()
    fetchUsers()
  } catch (err) {
    console.error('Failed to save user:', err)
    alert('Failed to save user: ' + (err.response?.data?.error || err.message))
  }
}

const deleteUser = async (id) => {
  if (!confirm('Are you sure you want to delete this user?')) {
    return
  }
  
  try {
    await api.deleteUser(id)
    fetchUsers()
  } catch (err) {
    console.error('Failed to delete user:', err)
    alert('Failed to delete user: ' + (err.response?.data?.error || err.message))
  }
}

const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  return new Date(dateString).toLocaleDateString()
}

onMounted(() => {
  fetchUsers()
  fetchRoles()
})
</script>
