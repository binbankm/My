<template>
  <div class="min-h-screen bg-gray-100">
    <!-- Top Navigation Bar (Mobile only) -->
    <nav class="md:hidden bg-white shadow-sm">
      <div class="px-4">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <!-- Mobile menu button -->
            <button
              @click="mobileMenuOpen = true"
              class="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-blue-500"
              aria-label="Open menu"
            >
              <Menu class="h-6 w-6" />
            </button>
            
            <div class="flex-shrink-0 flex items-center ml-2">
              <h1 class="text-xl font-bold">ServerPanel</h1>
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

    <!-- Desktop Layout with Sidebar -->
    <div class="flex md:h-screen">
      <!-- Desktop Sidebar -->
      <aside class="hidden md:flex md:flex-shrink-0">
        <div class="flex flex-col w-64 bg-white border-r border-gray-200">
          <!-- Logo/Header -->
          <div class="flex items-center h-16 px-6 border-b border-gray-200">
            <h1 class="text-xl font-bold text-gray-900">ServerPanel</h1>
          </div>

          <!-- Navigation -->
          <nav class="flex-1 px-4 py-6 space-y-1 overflow-y-auto">
            <router-link
              to="/"
              class="flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors"
              :class="$route.path === '/' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'"
            >
              <LayoutDashboard class="h-5 w-5 mr-3" />
              Dashboard
            </router-link>
            <router-link
              to="/containers"
              class="flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors"
              :class="$route.path === '/containers' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'"
            >
              <Container class="h-5 w-5 mr-3" />
              Containers
            </router-link>
            <router-link
              to="/files"
              class="flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors"
              :class="$route.path === '/files' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'"
            >
              <FolderOpen class="h-5 w-5 mr-3" />
              Files
            </router-link>
            <router-link
              to="/database"
              class="flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors"
              :class="$route.path === '/database' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'"
            >
              <Database class="h-5 w-5 mr-3" />
              Database
            </router-link>
            <router-link
              to="/terminal"
              class="flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors"
              :class="$route.path === '/terminal' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'"
            >
              <Terminal class="h-5 w-5 mr-3" />
              Terminal
            </router-link>
            <router-link
              to="/cron"
              class="flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors"
              :class="$route.path === '/cron' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'"
            >
              <Clock class="h-5 w-5 mr-3" />
              Cron Jobs
            </router-link>
            <router-link
              to="/logs"
              class="flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors"
              :class="$route.path === '/logs' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'"
            >
              <FileText class="h-5 w-5 mr-3" />
              Logs
            </router-link>
            <router-link
              to="/nginx"
              class="flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors"
              :class="$route.path === '/nginx' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'"
            >
              <Server class="h-5 w-5 mr-3" />
              Nginx
            </router-link>
            <router-link
              to="/backups"
              class="flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors"
              :class="$route.path === '/backups' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'"
            >
              <Archive class="h-5 w-5 mr-3" />
              Backups
            </router-link>
            <router-link
              to="/users"
              class="flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors"
              :class="$route.path === '/users' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'"
            >
              <Users class="h-5 w-5 mr-3" />
              Users
            </router-link>
            <router-link
              to="/settings"
              class="flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors"
              :class="$route.path === '/settings' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'"
            >
              <Settings class="h-5 w-5 mr-3" />
              Settings
            </router-link>
          </nav>

          <!-- User Info & Logout -->
          <div class="border-t border-gray-200 px-4 py-4">
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
      </aside>

      <!-- Main Content Area -->
      <div class="flex-1 overflow-auto">
        <main class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
          <router-view />
        </main>
      </div>
    </div>

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
              to="/terminal"
              @click="mobileMenuOpen = false"
              class="flex items-center px-3 py-2 text-base font-medium rounded-md transition-colors"
              :class="$route.path === '/terminal' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'"
            >
              <Terminal class="h-5 w-5 mr-3" />
              Terminal
            </router-link>
            <router-link
              to="/cron"
              @click="mobileMenuOpen = false"
              class="flex items-center px-3 py-2 text-base font-medium rounded-md transition-colors"
              :class="$route.path === '/cron' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'"
            >
              <Clock class="h-5 w-5 mr-3" />
              Cron Jobs
            </router-link>
            <router-link
              to="/logs"
              @click="mobileMenuOpen = false"
              class="flex items-center px-3 py-2 text-base font-medium rounded-md transition-colors"
              :class="$route.path === '/logs' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'"
            >
              <FileText class="h-5 w-5 mr-3" />
              Logs
            </router-link>
            <router-link
              to="/nginx"
              @click="mobileMenuOpen = false"
              class="flex items-center px-3 py-2 text-base font-medium rounded-md transition-colors"
              :class="$route.path === '/nginx' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'"
            >
              <Server class="h-5 w-5 mr-3" />
              Nginx
            </router-link>
            <router-link
              to="/backups"
              @click="mobileMenuOpen = false"
              class="flex items-center px-3 py-2 text-base font-medium rounded-md transition-colors"
              :class="$route.path === '/backups' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'"
            >
              <Archive class="h-5 w-5 mr-3" />
              Backups
            </router-link>
            <router-link
              to="/users"
              @click="mobileMenuOpen = false"
              class="flex items-center px-3 py-2 text-base font-medium rounded-md transition-colors"
              :class="$route.path === '/users' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50'"
            >
              <Users class="h-5 w-5 mr-3" />
              Users
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
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth'
import { Sheet } from '../components/ui'
import { Menu, LayoutDashboard, Container, FolderOpen, Database, Terminal, Settings, LogOut, Clock, FileText, Server, Archive, Users } from 'lucide-vue-next'

const router = useRouter()
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
