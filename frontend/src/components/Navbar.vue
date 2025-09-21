<script setup lang="ts">
import { useBlockchainStore } from '@/stores/blockchain'
import { ref, computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { 
  HomeIcon, 
  PencilSquareIcon, 
  MagnifyingGlassIcon, 
  WalletIcon,
  Bars3Icon,
  XMarkIcon
} from '@heroicons/vue/24/outline'
import WalletModal from './WalletModal.vue'
import DarkModeToggle from './DarkModeToggle.vue'

const store = useBlockchainStore()
const route = useRoute()
const showWalletModal = ref(false)
const mobileMenuOpen = ref(false)

const navItems = [
  { path: '/', name: 'Home', icon: HomeIcon },
  { path: '/create', name: 'Create', icon: PencilSquareIcon },
  { path: '/explore', name: 'Explore', icon: MagnifyingGlassIcon },
  { path: '/wallet', name: 'Wallet', icon: WalletIcon },
]

const isActive = (path: string) => route.path === path
</script>

<template>
  <nav class="sticky top-0 z-40 backdrop-blur-xl bg-white/80 dark:bg-dark-surface/80 border-b border-gray-200 dark:border-dark-border">
    <div class="container mx-auto px-4">
      <div class="flex justify-between items-center h-16">
        <!-- Logo and Desktop Nav -->
        <div class="flex items-center space-x-8">
          <RouterLink to="/" class="flex items-center space-x-2 group">
            <div class="p-2 rounded-xl bg-gradient-to-br from-accent-purple to-accent-blue group-hover:shadow-lg transition-all duration-300">
              <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
              </svg>
            </div>
            <span class="text-xl font-bold bg-gradient-to-r from-accent-purple to-accent-blue bg-clip-text text-transparent">
              BlogChain
            </span>
          </RouterLink>
          
          <!-- Desktop Navigation -->
          <div class="hidden md:flex items-center space-x-1">
            <RouterLink 
              v-for="item in navItems" 
              :key="item.path"
              :to="item.path"
              :class="[
                'flex items-center space-x-2 px-4 py-2 rounded-lg transition-all duration-200',
                isActive(item.path) 
                  ? 'bg-accent-blue/10 text-accent-blue dark:bg-accent-blue/20' 
                  : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-dark-elevated'
              ]"
            >
              <component :is="item.icon" class="w-4 h-4" />
              <span class="font-medium">{{ item.name }}</span>
            </RouterLink>
          </div>
        </div>

        <!-- Right side items -->
        <div class="flex items-center space-x-3">
          <!-- Balance Display -->
          <div v-if="store.currentAddress" class="hidden sm:flex items-center space-x-3">
            <div class="text-right">
              <p class="text-xs text-gray-500 dark:text-gray-400">Balance</p>
              <p class="text-sm font-semibold text-gray-900 dark:text-gray-100">{{ store.formattedBalance }}</p>
            </div>
            <div class="px-3 py-1.5 bg-gradient-to-r from-accent-purple/10 to-accent-blue/10 
                        dark:from-accent-purple/20 dark:to-accent-blue/20 
                        rounded-lg border border-accent-blue/20 dark:border-accent-blue/30">
              <span class="text-sm font-mono text-gray-700 dark:text-gray-300">
                {{ store.currentAddress.slice(0, 8) }}...{{ store.currentAddress.slice(-4) }}
              </span>
            </div>
          </div>
          
          <!-- Connect Button -->
          <button v-else @click="showWalletModal = true"
            class="hidden sm:flex items-center space-x-2 px-4 py-2 
                   bg-gradient-to-r from-accent-purple to-accent-blue 
                   text-white rounded-lg font-medium
                   hover:shadow-lg transform hover:-translate-y-0.5 transition-all duration-200">
            <WalletIcon class="w-4 h-4" />
            <span>Connect</span>
          </button>

          <!-- Dark Mode Toggle -->
          <DarkModeToggle />

          <!-- Mobile menu button -->
          <button 
            @click="mobileMenuOpen = !mobileMenuOpen"
            class="md:hidden p-2 rounded-lg text-gray-700 dark:text-gray-300 
                   hover:bg-gray-100 dark:hover:bg-dark-elevated"
          >
            <Bars3Icon v-if="!mobileMenuOpen" class="w-6 h-6" />
            <XMarkIcon v-else class="w-6 h-6" />
          </button>
        </div>
      </div>
    </div>

    <!-- Mobile Menu -->
    <Transition
      enter-active-class="transition-all duration-300"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition-all duration-300"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-2"
    >
      <div v-if="mobileMenuOpen" class="md:hidden bg-white dark:bg-dark-surface border-t border-gray-200 dark:border-dark-border">
        <div class="container mx-auto px-4 py-4 space-y-2">
          <RouterLink 
            v-for="item in navItems" 
            :key="item.path"
            :to="item.path"
            @click="mobileMenuOpen = false"
            :class="[
              'flex items-center space-x-3 px-4 py-3 rounded-lg transition-all duration-200',
              isActive(item.path) 
                ? 'bg-accent-blue/10 text-accent-blue dark:bg-accent-blue/20' 
                : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-dark-elevated'
            ]"
          >
            <component :is="item.icon" class="w-5 h-5" />
            <span class="font-medium">{{ item.name }}</span>
          </RouterLink>

          <!-- Mobile Wallet Section -->
          <div class="pt-4 border-t border-gray-200 dark:border-dark-border">
            <div v-if="store.currentAddress" class="px-4 py-3">
              <p class="text-xs text-gray-500 dark:text-gray-400 mb-1">Connected Wallet</p>
              <p class="text-sm font-mono text-gray-700 dark:text-gray-300">
                {{ store.currentAddress.slice(0, 12) }}...{{ store.currentAddress.slice(-4) }}
              </p>
              <p class="text-sm font-semibold text-gray-900 dark:text-gray-100 mt-2">
                {{ store.formattedBalance }}
              </p>
            </div>
            <button v-else 
              @click="showWalletModal = true; mobileMenuOpen = false"
              class="w-full px-4 py-3 bg-gradient-to-r from-accent-purple to-accent-blue 
                     text-white rounded-lg font-medium">
              Connect Wallet
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </nav>

  <WalletModal v-if="showWalletModal" @close="showWalletModal = false" />
</template>