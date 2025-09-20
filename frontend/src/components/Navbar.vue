<script setup lang="ts">
import { useBlockchainStore } from '@/stores/blockchain'
import { ref } from 'vue'
import { RouterLink } from 'vue-router'
import WalletModal from './WalletModal.vue'

const store = useBlockchainStore()
const showWalletModal = ref(false)
</script>

<template>
  <nav class="bg-white shadow-lg">
    <div class="container mx-auto px-4">
      <div class="flex justify-between items-center py-4">
        <div class="flex items-center space-x-8">
          <RouterLink to="/" class="text-2xl font-bold text-gray-800 hover:text-gray-600">
            📝 BlogChain
          </RouterLink>
          <div class="hidden md:flex space-x-4">
            <RouterLink to="/" class="text-gray-700 hover:text-gray-900">
              Home
            </RouterLink>
            <RouterLink to="/create" class="text-gray-700 hover:text-gray-900">
              Create Post
            </RouterLink>
            <RouterLink to="/explore" class="text-gray-700 hover:text-gray-900">
              Explore
            </RouterLink>
            <RouterLink to="/wallet" class="text-gray-700 hover:text-gray-900">
              Wallet
            </RouterLink>
          </div>
        </div>

        <div class="flex items-center space-x-4">
          <div v-if="store.currentAddress" class="flex items-center space-x-2">
            <span class="text-sm text-gray-600">{{ store.formattedBalance }}</span>
            <span class="px-3 py-1 bg-green-100 text-green-800 rounded-full text-sm">
              {{ store.currentAddress.slice(0, 10) }}...{{ store.currentAddress.slice(-4) }}
            </span>
          </div>
          <button v-else @click="showWalletModal = true"
            class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition">
            Connect Wallet
          </button>
        </div>
      </div>
    </div>
  </nav>

  <WalletModal v-if="showWalletModal" @close="showWalletModal = false" />
</template>