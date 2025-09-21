<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { CheckCircleIcon, XCircleIcon, ArrowPathIcon, InformationCircleIcon } from '@heroicons/vue/24/outline'

interface Toast {
  id: string
  type: 'success' | 'error' | 'pending' | 'info'
  title: string
  message: string
  txHash?: string
  duration?: number
}

const toasts = ref<Toast[]>([])

const addToast = (toast: Omit<Toast, 'id'>) => {
  const id = Math.random().toString(36).substr(2, 9)
  const newToast = { ...toast, id }
  toasts.value.push(newToast)
  
  if (toast.type !== 'pending' && (!toast.duration || toast.duration > 0)) {
    setTimeout(() => removeToast(id), toast.duration || 5000)
  }
  
  return id
}

const removeToast = (id: string) => {
  const index = toasts.value.findIndex(t => t.id === id)
  if (index !== -1) {
    toasts.value.splice(index, 1)
  }
}

const updateToast = (id: string, updates: Partial<Toast>) => {
  const toast = toasts.value.find(t => t.id === id)
  if (toast) {
    Object.assign(toast, updates)
    if (updates.type && updates.type !== 'pending') {
      setTimeout(() => removeToast(id), 5000)
    }
  }
}

// Expose methods for parent components
defineExpose({
  addToast,
  removeToast,
  updateToast
})
</script>

<template>
  <div class="fixed top-4 right-4 z-50 space-y-2 max-w-md">
    <TransitionGroup
      enter-active-class="transition-all duration-300"
      enter-from-class="translate-x-full opacity-0"
      enter-to-class="translate-x-0 opacity-100"
      leave-active-class="transition-all duration-300"
      leave-from-class="translate-x-0 opacity-100"
      leave-to-class="translate-x-full opacity-0"
    >
      <div
        v-for="toast in toasts"
        :key="toast.id"
        :class="[
          'relative overflow-hidden rounded-xl shadow-lg p-4 min-w-[320px]',
          'backdrop-blur-xl border',
          {
            'bg-green-500/10 border-green-500/30 dark:bg-green-500/20 dark:border-green-500/40': toast.type === 'success',
            'bg-red-500/10 border-red-500/30 dark:bg-red-500/20 dark:border-red-500/40': toast.type === 'error',
            'bg-blue-500/10 border-blue-500/30 dark:bg-blue-500/20 dark:border-blue-500/40': toast.type === 'pending',
            'bg-gray-500/10 border-gray-500/30 dark:bg-gray-500/20 dark:border-gray-500/40': toast.type === 'info',
          }
        ]"
      >
        <!-- Background animation for pending state -->
        <div
          v-if="toast.type === 'pending'"
          class="absolute inset-0 bg-gradient-to-r from-transparent via-white/5 to-transparent animate-pulse"
        />
        
        <div class="relative flex items-start gap-3">
          <!-- Icon -->
          <div class="flex-shrink-0">
            <CheckCircleIcon
              v-if="toast.type === 'success'"
              class="w-5 h-5 text-green-500"
            />
            <XCircleIcon
              v-else-if="toast.type === 'error'"
              class="w-5 h-5 text-red-500"
            />
            <ArrowPathIcon
              v-else-if="toast.type === 'pending'"
              class="w-5 h-5 text-blue-500 animate-spin"
            />
            <InformationCircleIcon
              v-else
              class="w-5 h-5 text-gray-500"
            />
          </div>
          
          <!-- Content -->
          <div class="flex-1 min-w-0">
            <p class="text-sm font-semibold text-gray-900 dark:text-gray-100">
              {{ toast.title }}
            </p>
            <p class="mt-1 text-sm text-gray-600 dark:text-gray-400">
              {{ toast.message }}
            </p>
            <a
              v-if="toast.txHash"
              :href="`https://explorer.example.com/tx/${toast.txHash}`"
              target="_blank"
              class="mt-2 inline-flex items-center text-xs text-blue-500 hover:text-blue-600 dark:text-blue-400 dark:hover:text-blue-300"
            >
              View transaction →
            </a>
          </div>
          
          <!-- Close button -->
          <button
            v-if="toast.type !== 'pending'"
            @click="removeToast(toast.id)"
            class="flex-shrink-0 ml-2 text-gray-400 hover:text-gray-500 dark:text-gray-500 dark:hover:text-gray-400"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
    </TransitionGroup>
  </div>
</template>