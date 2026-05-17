import { watch, onUnmounted } from 'vue'
import { useFsStore } from '@/stores/fs'

export function useFileWatcher() {
  const fs = useFsStore()
  let eventSource: EventSource | null = null

  function connect() {
    disconnect()

    if (!fs.currentMount) return

    const token = localStorage.getItem('token')
    if (!token) return

    const params = new URLSearchParams({
      mount: fs.currentMount,
      path: fs.currentPath,
      token,
    })

    eventSource = new EventSource(`/api/fs/watch?${params.toString()}`)

    eventSource.addEventListener('change', (e: MessageEvent) => {
      try {
        const event = JSON.parse(e.data)
        fs.applyEvent(event)
      } catch {
        // ignore malformed events
      }
    })

    eventSource.onerror = () => {
      // EventSource will auto-reconnect on error
    }
  }

  function disconnect() {
    if (eventSource) {
      eventSource.close()
      eventSource = null
    }
  }

  // Reconnect when mount or path changes
  watch(
    () => [fs.currentMount, fs.currentPath],
    () => connect(),
    { immediate: true }
  )

  onUnmounted(() => {
    disconnect()
  })

  return { connect, disconnect }
}
