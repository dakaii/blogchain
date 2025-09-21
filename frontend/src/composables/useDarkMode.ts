import { ref, watch, computed } from 'vue'

// Dark mode composable with system preference detection
const isDark = ref(false)

export function useDarkMode() {
  // Initialize dark mode based on system preference or saved preference
  const initDarkMode = () => {
    const saved = localStorage.getItem('darkMode')
    if (saved !== null) {
      isDark.value = saved === 'true'
    } else {
      // Check system preference
      isDark.value = window.matchMedia('(prefers-color-scheme: dark)').matches
    }
    updateDocumentClass()
  }

  // Update document class for Tailwind dark mode
  const updateDocumentClass = () => {
    if (isDark.value) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }

  // Toggle dark mode
  const toggleDarkMode = () => {
    isDark.value = !isDark.value
    localStorage.setItem('darkMode', String(isDark.value))
    updateDocumentClass()
  }

  // Watch for system preference changes
  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  mediaQuery.addEventListener('change', (e) => {
    if (localStorage.getItem('darkMode') === null) {
      isDark.value = e.matches
      updateDocumentClass()
    }
  })

  // Initialize on load
  initDarkMode()

  return {
    isDark: computed(() => isDark.value),
    toggleDarkMode
  }
}