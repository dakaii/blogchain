<script setup lang="ts">
import { useDarkMode } from '@/composables/useDarkMode'
import { SunIcon, MoonIcon } from '@heroicons/vue/24/outline'

const { isDark, toggleDarkMode } = useDarkMode()
</script>

<template>
  <button
    @click="toggleDarkMode"
    class="relative p-2 rounded-xl transition-all duration-300 
           bg-gray-100 hover:bg-gray-200 
           dark:bg-dark-surface dark:hover:bg-dark-elevated
           group"
    :aria-label="isDark ? 'Switch to light mode' : 'Switch to dark mode'"
  >
    <div class="relative w-6 h-6">
      <Transition
        enter-active-class="transition-all duration-300"
        enter-from-class="rotate-90 opacity-0"
        enter-to-class="rotate-0 opacity-100"
        leave-active-class="transition-all duration-300"
        leave-from-class="rotate-0 opacity-100"
        leave-to-class="-rotate-90 opacity-0"
      >
        <SunIcon
          v-if="isDark"
          class="absolute inset-0 w-6 h-6 text-yellow-500"
        />
        <MoonIcon
          v-else
          class="absolute inset-0 w-6 h-6 text-gray-700 dark:text-gray-300"
        />
      </Transition>
    </div>
  </button>
</template>