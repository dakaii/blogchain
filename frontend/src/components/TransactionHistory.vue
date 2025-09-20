<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useBlockchainStore } from '@/stores/blockchain'

const store = useBlockchainStore()
const transactions = ref<any[]>([])
const loading = ref(false)

const formattedTransactions = computed(() => {
  return transactions.value.map(tx => {
    const messages = tx.tx?.body?.messages || []
    const timestamp = new Date(tx.timestamp).toLocaleString()
    const hash = tx.txhash
    
    // Determine transaction type based on message type
    const msgType = messages[0]?.['@type'] || ''
    let type = 'unknown'
    let amount = '0'
    let recipient = ''
    
    if (msgType.includes('MsgSend')) {
      const msg = messages[0]
      type = msg.from_address === store.currentAddress ? 'sent' : 'received'
      amount = msg.amount?.[0]?.amount || '0'
      recipient = type === 'sent' ? msg.to_address : msg.from_address
    } else if (msgType.includes('MsgCreatePost')) {
      type = 'create_post'
    } else if (msgType.includes('MsgLikePost')) {
      type = 'like_post'
    }
    
    return {
      hash: hash?.slice(0, 8) + '...' + hash?.slice(-8),
      fullHash: hash,
      type,
      amount,
      recipient,
      timestamp,
      success: tx.code === 0,
      height: tx.height,
      gas: tx.gas_used,
    }
  })
})

async function loadTransactions() {
  if (!store.currentAddress) return
  
  loading.value = true
  try {
    transactions.value = await store.service.getTransactionHistory(store.currentAddress)
  } catch (error) {
    console.error('Failed to load transactions:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (store.currentAddress) {
    loadTransactions()
  }
})

function getTypeColor(type: string) {
  switch (type) {
    case 'sent': return 'text-red-600 bg-red-100'
    case 'received': return 'text-green-600 bg-green-100'
    case 'create_post': return 'text-blue-600 bg-blue-100'
    case 'like_post': return 'text-pink-600 bg-pink-100'
    default: return 'text-gray-600 bg-gray-100'
  }
}

function getTypeIcon(type: string) {
  switch (type) {
    case 'sent': return '↑'
    case 'received': return '↓'
    case 'create_post': return '📝'
    case 'like_post': return '❤️'
    default: return '•'
  }
}

function formatAmount(amount: string) {
  return (parseInt(amount) / 1000000).toFixed(2)
}
</script>

<template>
  <div class="bg-white rounded-lg shadow-md p-6">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-xl font-semibold text-gray-800">Transaction History</h3>
      <button
        @click="loadTransactions"
        :disabled="loading"
        class="px-3 py-1 text-sm bg-gray-100 hover:bg-gray-200 rounded-lg transition"
      >
        {{ loading ? 'Loading...' : 'Refresh' }}
      </button>
    </div>
    
    <div v-if="!store.currentAddress" class="text-gray-500 text-center py-8">
      Connect your wallet to view transaction history
    </div>
    
    <div v-else-if="loading && transactions.length === 0" class="text-center py-8">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
    </div>
    
    <div v-else-if="transactions.length === 0" class="text-gray-500 text-center py-8">
      No transactions yet
    </div>
    
    <div v-else class="space-y-3">
      <div
        v-for="tx in formattedTransactions"
        :key="tx.fullHash"
        class="flex items-center justify-between p-3 border rounded-lg hover:bg-gray-50 transition"
      >
        <div class="flex items-center space-x-3">
          <div
            :class="getTypeColor(tx.type)"
            class="w-10 h-10 rounded-full flex items-center justify-center text-xl"
          >
            {{ getTypeIcon(tx.type) }}
          </div>
          <div>
            <div class="font-medium text-gray-800">
              {{ tx.hash }}
              <span
                v-if="!tx.success"
                class="ml-2 px-2 py-0.5 text-xs bg-red-100 text-red-600 rounded"
              >
                Failed
              </span>
            </div>
            <div class="text-sm text-gray-500">{{ tx.timestamp }}</div>
          </div>
        </div>
        
        <div class="text-right">
          <div v-if="tx.amount !== '0'" class="font-medium">
            <span v-if="tx.type === 'sent'" class="text-red-600">
              -{{ formatAmount(tx.amount) }} STAKE
            </span>
            <span v-else-if="tx.type === 'received'" class="text-green-600">
              +{{ formatAmount(tx.amount) }} STAKE
            </span>
          </div>
          <div class="text-xs text-gray-500">
            Block #{{ tx.height }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>