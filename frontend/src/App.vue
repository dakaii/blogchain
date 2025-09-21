<script setup lang="ts">
import { RouterView } from 'vue-router'
import { onMounted, ref } from 'vue'
import Navbar from './components/Navbar.vue'
import TransactionToast from './components/TransactionToast.vue'
import { useDarkMode } from './composables/useDarkMode'

const toastRef = ref()
const { isDark } = useDarkMode()

// Provide toast methods globally
onMounted(() => {
  if (toastRef.value) {
    window.$toast = {
      success: (title: string, message: string, options?: any) => 
        toastRef.value.addToast({ type: 'success', title, message, ...options }),
      error: (title: string, message: string, options?: any) => 
        toastRef.value.addToast({ type: 'error', title, message, ...options }),
      pending: (title: string, message: string, options?: any) => 
        toastRef.value.addToast({ type: 'pending', title, message, ...options }),
      info: (title: string, message: string, options?: any) => 
        toastRef.value.addToast({ type: 'info', title, message, ...options }),
      update: toastRef.value.updateToast
    }
  }
})
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-dark-bg transition-colors duration-300">
    <Navbar />
    <main class="container mx-auto px-4 py-8">
      <RouterView v-slot="{ Component }">
        <Transition
          enter-active-class="transition-all duration-300"
          enter-from-class="opacity-0 translate-y-2"
          enter-to-class="opacity-100 translate-y-0"
          leave-active-class="transition-all duration-300"
          leave-from-class="opacity-100 translate-y-0"
          leave-to-class="opacity-0 -translate-y-2"
        >
          <component :is="Component" />
        </Transition>
      </RouterView>
    </main>
    <TransactionToast ref="toastRef" />
  </div>
</template>

<style>
/* Styles are in main.css */
</style>