<template>
  <div class="h-full flex flex-col">
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-2xl font-bold">Terminal</h2>
      <div class="flex gap-2">
        <button
          @click="reconnect"
          :disabled="connecting"
          class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {{ connecting ? 'Connecting...' : 'Reconnect' }}
        </button>
        <button
          @click="clearTerminal"
          class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
        >
          Clear
        </button>
      </div>
    </div>

    <!-- Connection Status -->
    <div v-if="error" class="mb-4 p-4 bg-red-50 border border-red-200 rounded">
      <p class="text-red-800">❌ {{ error }}</p>
    </div>

    <div v-else-if="!connected && !connecting" class="mb-4 p-4 bg-yellow-50 border border-yellow-200 rounded">
      <p class="text-yellow-800">⚠️ Terminal disconnected. Click "Reconnect" to start a new session.</p>
    </div>

    <div v-else-if="connecting" class="mb-4 p-4 bg-blue-50 border border-blue-200 rounded">
      <p class="text-blue-800">🔄 Connecting to terminal...</p>
    </div>

    <!-- Terminal Container -->
    <div class="flex-1 bg-black rounded-lg shadow-lg overflow-hidden">
      <div ref="terminalRef" class="w-full h-full"></div>
    </div>

    <!-- Terminal Info -->
    <div class="mt-4 text-sm text-gray-600">
      <p>💡 Tip: Use Ctrl+C to interrupt running commands. Type 'exit' to close the shell.</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import '@xterm/xterm/css/xterm.css'

const terminalRef = ref(null)
const terminal = ref(null)
const fitAddon = ref(null)
const ws = ref(null)
const connected = ref(false)
const connecting = ref(false)
const error = ref(null)

const initTerminal = () => {
  // Create terminal instance
  terminal.value = new Terminal({
    cursorBlink: true,
    fontSize: 14,
    fontFamily: 'Menlo, Monaco, "Courier New", monospace',
    theme: {
      background: '#000000',
      foreground: '#ffffff',
      cursor: '#ffffff',
      selection: 'rgba(255, 255, 255, 0.3)',
    },
    scrollback: 1000,
  })

  // Create fit addon
  fitAddon.value = new FitAddon()
  terminal.value.loadAddon(fitAddon.value)

  // Open terminal in DOM
  terminal.value.open(terminalRef.value)
  
  // Fit terminal to container
  fitAddon.value.fit()

  // Handle user input
  terminal.value.onData((data) => {
    if (ws.value && ws.value.readyState === WebSocket.OPEN) {
      ws.value.send(JSON.stringify({
        type: 'input',
        data: data
      }))
    }
  })

  // Handle terminal resize
  terminal.value.onResize(({ cols, rows }) => {
    if (ws.value && ws.value.readyState === WebSocket.OPEN) {
      ws.value.send(JSON.stringify({
        type: 'resize',
        cols: cols,
        rows: rows
      }))
    }
  })
}

const connectWebSocket = () => {
  connecting.value = true
  error.value = null

  const token = localStorage.getItem('token')
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const wsUrl = `${protocol}//${window.location.host}/api/terminal/ws?token=${token}`

  ws.value = new WebSocket(wsUrl)

  ws.value.onopen = () => {
    connecting.value = false
    connected.value = true
    error.value = null
    terminal.value?.writeln('\r\n✅ Connected to terminal\r\n')
  }

  ws.value.onmessage = (event) => {
    try {
      // Try to parse as JSON (for control messages)
      const msg = JSON.parse(event.data)
      if (msg.type === 'error') {
        error.value = msg.error
        terminal.value?.writeln(`\r\n❌ Error: ${msg.error}\r\n`)
      }
    } catch (e) {
      // Not JSON, treat as terminal output
      terminal.value?.write(event.data)
    }
  }

  ws.value.onerror = (err) => {
    console.error('WebSocket error:', err)
    error.value = 'WebSocket connection error'
    connecting.value = false
  }

  ws.value.onclose = () => {
    connected.value = false
    connecting.value = false
    terminal.value?.writeln('\r\n\r\n❌ Connection closed\r\n')
  }
}

const reconnect = async () => {
  // Close existing connection
  if (ws.value) {
    ws.value.close()
    ws.value = null
  }

  // Clear terminal
  terminal.value?.clear()

  // Wait a moment then reconnect
  await nextTick()
  connectWebSocket()
}

const clearTerminal = () => {
  terminal.value?.clear()
}

const handleResize = () => {
  if (fitAddon.value) {
    fitAddon.value.fit()
  }
}

onMounted(() => {
  initTerminal()
  connectWebSocket()

  // Handle window resize
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  
  if (ws.value) {
    ws.value.close()
  }
  
  if (terminal.value) {
    terminal.value.dispose()
  }
})
</script>

<style scoped>
/* Ensure terminal takes full height */
:deep(.xterm) {
  height: 100%;
  padding: 10px;
}

:deep(.xterm-viewport) {
  overflow-y: auto;
}
</style>
