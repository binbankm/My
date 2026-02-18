import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || null)

  const login = async (username, password) => {
    const response = await api.login({ username, password })
    token.value = response.token
    user.value = response.user
    localStorage.setItem('token', response.token)
    return response
  }

  const logout = async () => {
    try {
      await api.logout()
    } finally {
      token.value = null
      user.value = null
      localStorage.removeItem('token')
    }
  }

  const fetchUserInfo = async () => {
    if (token.value) {
      user.value = await api.getUserInfo()
    }
  }

  return {
    user,
    token,
    login,
    logout,
    fetchUserInfo,
  }
})
