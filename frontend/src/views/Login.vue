<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <Card class="max-w-md w-full">
      <CardHeader>
        <CardTitle class="text-center">ServerPanel</CardTitle>
        <p class="mt-2 text-center text-sm text-muted-foreground">
          Linux Server Management Panel
        </p>
      </CardHeader>
      <CardContent>
        <form class="space-y-4" @submit.prevent="handleLogin">
          <div class="space-y-2">
            <Label for="username">Username</Label>
            <Input
              id="username"
              v-model="username"
              type="text"
              placeholder="Enter username"
              required
            />
          </div>
          
          <div class="space-y-2">
            <Label for="password">Password</Label>
            <Input
              id="password"
              v-model="password"
              type="password"
              placeholder="Enter password"
              required
            />
          </div>

          <div v-if="error" class="text-sm text-destructive text-center">
            {{ error }}
          </div>

          <Button
            type="submit"
            :disabled="loading"
            class="w-full"
          >
            {{ loading ? 'Logging in...' : 'Sign in' }}
          </Button>
        </form>
      </CardContent>
    </Card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth'
import { Card, CardHeader, CardTitle, CardContent, Button, Input, Label } from '@/components/ui'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('admin')
const password = ref('')
const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  try {
    loading.value = true
    error.value = ''
    await authStore.login(username.value, password.value)
    router.push('/')
  } catch (err) {
    error.value = err.response?.data?.error || 'Login failed'
  } finally {
    loading.value = false
  }
}
</script>
