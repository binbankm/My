<template>
  <div class="min-h-screen bg-gray-100">
    <nav class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex">
            <!-- Mobile menu button -->
            <button
              @click="mobileMenuOpen = true"
              class="md:hidden inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-blue-500"
              aria-label="Open menu"
            >
              <Menu class="h-6 w-6" />
            </button>
            
            <div class="flex-shrink-0 flex items-center ml-2 md:ml-0">
              <h1 class="text-xl font-bold">ServerPanel</h1>
            </div>
            <div class="hidden md:ml-6 md:flex md:space-x-8">
              <router-link
                to="/"
                class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
                active-class="border-blue-500 text-gray-900"
              >
                <LayoutDashboard class="h-4 w-4 mr-2" />
                Dashboard
              </router-link>
              <router-link
                to="/containers"
                class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
                active-class="border-blue-500 text-gray-900"
              >
                <Container class="h-4 w-4 mr-2" />
                Containers
              </router-link>
              <router-link
                to="/files"
                class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
                active-class="border-blue-500 text-gray-900"
              >
                <FolderOpen class="h-4 w-4 mr-2" />
                Files
              </router-link>
              <router-link
                to="/database"
                class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
                active-class="border-blue-500 text-gray-900"
              >
                <Database class="h-4 w-4 mr-2" />
                Database
              </router-link>
              <router-link
                to="/settings"
                class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
                active-class="border-blue-500 text-gray-900"
              >
                <Settings class="h-4 w-4 mr-2" />
                Settings
              </router-link>
            </div>
          </div>
          <div class="flex items-center">
            <span class="text-sm text-gray-700 mr-4">{{ user?.username }}</span>
            <button
              @click="handleLogout"
              class="text-sm text-gray-500 hover:text-gray-700"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </nav>

    <!-- Mobile Drawer Menu -->
    <Sheet v-model:open="mobileMenuOpen" side="left">
      <div class="flex flex-col h-full">
        <div class="px-4 py-6">
          <h2 class="text-lg font-semibold text-gray-900 mb-4">Navigation</h2>
          <nav class="space-y-1">
            <router-link
              to="/"
              @click="mobileMenuOpen = false"
              class="flex items-center px-3 py-2 text-base font-medium rounded-md transition-colors"
              :class="$route.path === '/' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'"
            >
              <LayoutDashboard class="h-5 w-5 mr-3" />
              Dashboard
            </router-link>
            <router-link
              to="/containers"
              @click="mobileMenuOpen = false"
              class="flex items-center px-3 py-2 text-base font-medium rounded-md transition-colors"
              :class="$route.path === '/containers' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'"
            >
              <Container class="h-5 w-5 mr-3" />
              Containers
            </router-link>
            <router-link
              to="/files"
              @click="mobileMenuOpen = false"
              class="flex items-center px-3 py-2 text-base font-medium rounded-md transition-colors"
              :class="$route.path === '/files' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'"
            >
              <FolderOpen class="h-5 w-5 mr-3" />
              Files
            </router-link>
            <router-link
              to="/database"
              @click="mobileMenuOpen = false"
              class="flex items-center px-3 py-2 text-base font-medium rounded-md transition-colors"
              :class="$route.path === '/database' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'"
            >
              <Database class="h-5 w-5 mr-3" />
              Database
            </router-link>
            <router-link
              to="/settings"
              @click="mobileMenuOpen = false"
              class="flex items-center px-3 py-2 text-base font-medium rounded-md transition-colors"
              :class="$route.path === '/settings' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'"
            >
              <Settings class="h-5 w-5 mr-3" />
              Settings
            </router-link>
          </nav>
        </div>
        
        <div class="mt-auto border-t border-gray-200 px-4 py-4">
          <div class="flex items-center mb-3">
            <div class="flex-1">
              <p class="text-sm font-medium text-gray-900">{{ user?.username }}</p>
              <p class="text-xs text-gray-500">Administrator</p>
            </div>
          </div>
          <button
            @click="handleLogout"
            class="w-full flex items-center justify-center px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-md hover:bg-red-700 transition-colors"
          >
            <LogOut class="h-4 w-4 mr-2" />
            Logout
          </button>
        </div>
      </div>
    </Sheet>

    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../store/auth'
import { Sheet } from '../components/ui'
import { Menu, LayoutDashboard, Container, FolderOpen, Database, Settings, LogOut } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const user = computed(() => authStore.user)
const mobileMenuOpen = ref(false)

onMounted(async () => {
  if (!authStore.user) {
    await authStore.fetchUserInfo()
  }
})

const handleLogout = async () => {
  mobileMenuOpen.value = false
  await authStore.logout()
  router.push('/login')
}
</script>
