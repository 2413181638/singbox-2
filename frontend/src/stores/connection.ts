import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '@/utils/api'

export interface ConnectionStatus {
  isRunning: boolean
  stats?: {
    up: number
    down: number
  }
}

export const useConnectionStore = defineStore('connection', () => {
  const isRunning = ref(false)
  const loading = ref(false)
  const stats = ref<any>({})

  const startConnection = async () => {
    loading.value = true
    try {
      await api.startConnection()
      isRunning.value = true
    } finally {
      loading.value = false
    }
  }

  const stopConnection = async () => {
    loading.value = true
    try {
      await api.stopConnection()
      isRunning.value = false
    } finally {
      loading.value = false
    }
  }

  const updateStatus = async () => {
    try {
      const status = await api.getConnectionStatus()
      isRunning.value = status.isRunning
      stats.value = status.stats || {}
    } catch (error) {
      console.error('Failed to update connection status:', error)
    }
  }

  return {
    isRunning,
    loading,
    stats,
    startConnection,
    stopConnection,
    updateStatus
  }
})