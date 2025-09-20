<script setup lang="ts">
import { ref, computed } from 'vue'
import { useBlockchainStore } from '@/stores/blockchain'

const store = useBlockchainStore()
const searchQuery = ref('')

const formatBalance = (amount: string, denom: string) => {
  const value = parseInt(amount) / 1000000
  return `${value.toFixed(2)} ${denom.toUpperCase()}`
}

const filteredBalances = computed(() => {
  if (!searchQuery.value) return store.balance
  
  return store.balance.filter(b => 
    b.denom.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

const totalValueUSD = computed(() => {
  // Mock USD values - in production, fetch from price API
  const prices: Record<string, number> = {
    'stake': 0.5,
    'token': 0.1,
  }
  
  return store.balance.reduce((total, b) => {
    const amount = parseInt(b.amount) / 1000000
    const price = prices[b.denom] || 0
    return total + (amount * price)
  }, 0).toFixed(2)
})

function getDenomIcon(denom: string) {
  switch (denom) {
    case 'stake': return '🪙'
    case 'token': return '🎫'
    default: return '💰'
  }
}
</script>

<template>
  <div class="bg-white rounded-lg shadow-md p-6">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-xl font-semibold text-gray-800">Assets</h3>
      <div v-if="store.balance.length > 0" class="relative">
        <input
          v-model="searchQuery"
          type="search"
          placeholder="Search assets..."
          class="pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <svg
          class="absolute left-3 top-2.5 w-5 h-5 text-gray-400"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
          />
        </svg>
      </div>
    </div>
    
    <div v-if="!store.currentAddress" class="text-gray-500 text-center py-8">
      Connect your wallet to view assets
    </div>
    
    <div v-else-if="store.balance.length === 0" class="text-gray-500 text-center py-8">
      No assets in wallet
    </div>
    
    <div v-else>
      <!-- Total Portfolio Value -->
      <div class="mb-6 p-4 bg-gradient-to-r from-blue-50 to-purple-50 rounded-lg">
        <div class="text-sm text-gray-600">Total Portfolio Value</div>
        <div class="text-2xl font-bold text-gray-800">${{ totalValueUSD }}</div>
      </div>
      
      <!-- Asset List -->
      <div class="space-y-3">
        <div
          v-for="balance in filteredBalances"
          :key="balance.denom"
          class="flex items-center justify-between p-4 border rounded-lg hover:bg-gray-50 transition"
        >
          <div class="flex items-center space-x-3">
            <div class="w-10 h-10 bg-gray-100 rounded-full flex items-center justify-center text-xl">
              {{ getDenomIcon(balance.denom) }}
            </div>
            <div>
              <div class="font-medium text-gray-800">
                {{ balance.denom.toUpperCase() }}
              </div>
              <div class="text-sm text-gray-500">
                Available
              </div>
            </div>
          </div>
          
          <div class="text-right">
            <div class="font-semibold text-gray-800">
              {{ formatBalance(balance.amount, balance.denom) }}
            </div>
            <div class="text-sm text-gray-500">
              ${{ ((parseInt(balance.amount) / 1000000) * (balance.denom === 'stake' ? 0.5 : 0.1)).toFixed(2) }}
            </div>
          </div>
        </div>
      </div>
      
      <div v-if="filteredBalances.length === 0 && searchQuery" class="text-center py-4 text-gray-500">
        No assets match your search
      </div>
    </div>
  </div>
</template>