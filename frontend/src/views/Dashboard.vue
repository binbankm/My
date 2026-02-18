<template>
  <div>
    <h2 class="text-2xl font-bold mb-6">Dashboard</h2>

    <!-- System Info -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
      <Card class="p-6">
        <div class="text-sm text-muted-foreground">Hostname</div>
        <div class="text-2xl font-bold">{{ systemInfo.hostname || 'Loading...' }}</div>
      </Card>
      <Card class="p-6">
        <div class="text-sm text-muted-foreground">OS</div>
        <div class="text-2xl font-bold">{{ systemInfo.platform || 'N/A' }}</div>
      </Card>
      <Card class="p-6">
        <div class="text-sm text-muted-foreground">CPU Cores</div>
        <div class="text-2xl font-bold">{{ systemInfo.cpuCores || 0 }}</div>
      </Card>
      <Card class="p-6">
        <div class="text-sm text-muted-foreground">Uptime</div>
        <div class="text-2xl font-bold">{{ formatUptime(systemInfo.uptime) }}</div>
      </Card>
    </div>

    <!-- System Stats -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">CPU Usage</CardTitle>
        </CardHeader>
        <CardContent>
          <div v-if="stats.cpu && stats.cpu.length">
            <div v-for="(cpu, index) in stats.cpu" :key="index" class="mb-2">
              <div class="flex justify-between text-sm mb-1">
                <span>Core {{ index }}</span>
                <span>{{ cpu.toFixed(1) }}%</span>
              </div>
              <div class="w-full bg-gray-200 rounded-full h-2">
                <div class="bg-blue-600 h-2 rounded-full" :style="{ width: cpu + '%' }"></div>
              </div>
            </div>
          </div>
          <div v-else class="text-muted-foreground">Loading...</div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Memory Usage</CardTitle>
        </CardHeader>
        <CardContent>
          <div v-if="stats.memory">
            <div class="flex justify-between text-sm mb-1">
              <span>Used: {{ formatBytes(stats.memory.used) }}</span>
              <span>Total: {{ formatBytes(stats.memory.total) }}</span>
            </div>
            <div class="w-full bg-gray-200 rounded-full h-2 mb-2">
              <div class="bg-green-600 h-2 rounded-full" :style="{ width: stats.memory.usedPercent + '%' }"></div>
            </div>
            <div class="text-sm text-muted-foreground">{{ stats.memory.usedPercent.toFixed(1) }}% used</div>
          </div>
          <div v-else class="text-muted-foreground">Loading...</div>
        </CardContent>
      </Card>
    </div>

    <!-- Disk Usage -->
    <Card>
      <CardHeader>
        <CardTitle class="text-lg">Disk Usage</CardTitle>
      </CardHeader>
      <CardContent>
        <div v-if="stats.disk && stats.disk.length">
          <div v-for="disk in stats.disk" :key="disk.path" class="mb-4">
            <div class="flex justify-between text-sm mb-1">
              <span>{{ disk.path }}</span>
              <span>{{ formatBytes(disk.used) }} / {{ formatBytes(disk.total) }}</span>
            </div>
            <div class="w-full bg-gray-200 rounded-full h-2">
              <div class="bg-yellow-600 h-2 rounded-full" :style="{ width: disk.usedPercent + '%' }"></div>
            </div>
          </div>
        </div>
        <div v-else class="text-muted-foreground">Loading...</div>
      </CardContent>
    </Card>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import api from '../api'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui'

const systemInfo = ref({})
const stats = ref({})
let intervalId = null

const fetchSystemInfo = async () => {
  try {
    systemInfo.value = await api.getSystemInfo()
  } catch (error) {
    console.error('Failed to fetch system info:', error)
  }
}

const fetchStats = async () => {
  try {
    stats.value = await api.getSystemStats()
  } catch (error) {
    console.error('Failed to fetch stats:', error)
  }
}

const formatBytes = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatUptime = (seconds) => {
  if (!seconds) return '0d 0h'
  const days = Math.floor(seconds / 86400)
  const hours = Math.floor((seconds % 86400) / 3600)
  return `${days}d ${hours}h`
}

onMounted(() => {
  fetchSystemInfo()
  fetchStats()
  intervalId = setInterval(fetchStats, 5000)
})

onUnmounted(() => {
  if (intervalId) {
    clearInterval(intervalId)
  }
})
</script>
