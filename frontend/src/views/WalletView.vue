<script setup lang="ts">
import { onMounted } from 'vue'
import { useBlockchainStore } from '@/stores/blockchain'
import AssetsDisplay from '@/components/AssetsDisplay.vue'
import TokenTransfer from '@/components/TokenTransfer.vue'
import TransactionHistory from '@/components/TransactionHistory.vue'

const store = useBlockchainStore()

onMounted(async () => {
  if (store.currentAddress) {
    await store.updateBalance()
  }
})
</script>

<template>
  <div>
    <div class="mb-8">
      <h1 class="text-4xl font-bold text-gray-800 mb-2">Wallet Dashboard</h1>
      <p class="text-gray-600">Manage your assets and transactions</p>
    </div>
    
    <div v-if="!store.currentAddress" class="bg-yellow-100 border-l-4 border-yellow-500 text-yellow-700 p-4 mb-6">
      <p class="font-medium">Wallet not connected</p>
      <p class="text-sm mt-1">Please connect your wallet to access dashboard features.</p>
    </div>
    
    <div v-else class="grid gap-6 lg:grid-cols-2">
      <!-- Left Column -->
      <div class="space-y-6">
        <AssetsDisplay />
        <TokenTransfer />
      </div>
      
      <!-- Right Column -->
      <div>
        <TransactionHistory />
      </div>
    </div>
  </div>
</template>