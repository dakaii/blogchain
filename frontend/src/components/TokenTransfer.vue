<script setup lang="ts">
import { ref, computed } from 'vue'
import { useBlockchainStore } from '@/stores/blockchain'

const store = useBlockchainStore()

// Tab state
const activeTab = ref<'send' | 'receive'>('send')

// Send form state
const recipientAddress = ref('')
const amount = ref('')
const selectedDenom = ref('stake')
const memo = ref('')
const sending = ref(false)
const sendError = ref('')
const sendSuccess = ref(false)

// Validation
const isValidAddress = computed(() => {
  return recipientAddress.value.startsWith('blogchain') && recipientAddress.value.length > 40
})

const isValidAmount = computed(() => {
  const num = parseFloat(amount.value)
  if (isNaN(num) || num <= 0) return false
  
  const balance = store.balance.find(b => b.denom === selectedDenom.value)
  if (!balance) return false
  
  const available = parseInt(balance.amount) / 1000000
  return num <= available
})

const availableBalance = computed(() => {
  const balance = store.balance.find(b => b.denom === selectedDenom.value)
  if (!balance) return '0'
  return (parseInt(balance.amount) / 1000000).toFixed(2)
})

async function sendTokens() {
  if (!isValidAddress.value || !isValidAmount.value) {
    sendError.value = 'Please check recipient address and amount'
    return
  }
  
  sending.value = true
  sendError.value = ''
  sendSuccess.value = false
  
  try {
    const amountInMicrounits = (parseFloat(amount.value) * 1000000).toString()
    await store.service.sendTokens(
      store.currentAddress,
      recipientAddress.value,
      amountInMicrounits,
      selectedDenom.value
    )
    
    sendSuccess.value = true
    // Reset form
    recipientAddress.value = ''
    amount.value = ''
    memo.value = ''
    
    // Refresh balance
    await store.updateBalance()
  } catch (error: any) {
    sendError.value = error.message || 'Failed to send tokens'
  } finally {
    sending.value = false
  }
}

function copyAddress() {
  navigator.clipboard.writeText(store.currentAddress)
}

function setMaxAmount() {
  amount.value = availableBalance.value
}
</script>

<template>
  <div class="bg-white rounded-lg shadow-md p-6">
    <!-- Tab Header -->
    <div class="flex border-b mb-6">
      <button
        @click="activeTab = 'send'"
        :class="[
          'px-4 py-2 font-medium transition',
          activeTab === 'send' 
            ? 'text-blue-600 border-b-2 border-blue-600' 
            : 'text-gray-600 hover:text-gray-800'
        ]"
      >
        Send
      </button>
      <button
        @click="activeTab = 'receive'"
        :class="[
          'px-4 py-2 font-medium transition ml-4',
          activeTab === 'receive' 
            ? 'text-blue-600 border-b-2 border-blue-600' 
            : 'text-gray-600 hover:text-gray-800'
        ]"
      >
        Receive
      </button>
    </div>
    
    <!-- Send Tab -->
    <div v-if="activeTab === 'send'">
      <div v-if="!store.currentAddress" class="text-gray-500 text-center py-8">
        Connect your wallet to send tokens
      </div>
      
      <form v-else @submit.prevent="sendTokens" class="space-y-4">
        <!-- Recipient Address -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Recipient Address
          </label>
          <input
            v-model="recipientAddress"
            type="text"
            placeholder="blogchain1..."
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            :class="{ 'border-red-500': recipientAddress && !isValidAddress }"
          />
          <p v-if="recipientAddress && !isValidAddress" class="text-red-500 text-sm mt-1">
            Invalid address format
          </p>
        </div>
        
        <!-- Amount and Token Selection -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Amount
          </label>
          <div class="flex space-x-2">
            <input
              v-model="amount"
              type="number"
              step="0.000001"
              placeholder="0.00"
              class="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              :class="{ 'border-red-500': amount && !isValidAmount }"
            />
            <select
              v-model="selectedDenom"
              class="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option v-for="balance in store.balance" :key="balance.denom" :value="balance.denom">
                {{ balance.denom.toUpperCase() }}
              </option>
            </select>
            <button
              type="button"
              @click="setMaxAmount"
              class="px-4 py-2 text-sm bg-gray-100 hover:bg-gray-200 rounded-lg transition"
            >
              Max
            </button>
          </div>
          <p class="text-sm text-gray-500 mt-1">
            Available: {{ availableBalance }} {{ selectedDenom.toUpperCase() }}
          </p>
          <p v-if="amount && !isValidAmount" class="text-red-500 text-sm mt-1">
            Insufficient balance
          </p>
        </div>
        
        <!-- Memo (Optional) -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Memo (Optional)
          </label>
          <input
            v-model="memo"
            type="text"
            placeholder="Add a note..."
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        
        <!-- Success/Error Messages -->
        <div v-if="sendSuccess" class="p-3 bg-green-100 text-green-700 rounded-lg">
          Transaction sent successfully!
        </div>
        <div v-if="sendError" class="p-3 bg-red-100 text-red-700 rounded-lg">
          {{ sendError }}
        </div>
        
        <!-- Submit Button -->
        <button
          type="submit"
          :disabled="sending || !isValidAddress || !isValidAmount"
          class="w-full px-4 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition"
        >
          {{ sending ? 'Sending...' : 'Send Tokens' }}
        </button>
      </form>
    </div>
    
    <!-- Receive Tab -->
    <div v-else>
      <div v-if="!store.currentAddress" class="text-gray-500 text-center py-8">
        Connect your wallet to receive tokens
      </div>
      
      <div v-else class="space-y-4">
        <!-- QR Code placeholder -->
        <div class="bg-gray-100 rounded-lg p-8 flex items-center justify-center">
          <div class="text-center">
            <div class="w-32 h-32 bg-white rounded-lg mb-4 mx-auto flex items-center justify-center">
              <!-- Simple QR placeholder -->
              <svg class="w-24 h-24" viewBox="0 0 24 24" fill="currentColor">
                <path d="M3 3h7v7H3V3m2 2v3h3V5H5m9-2h7v7h-7V3m2 2v3h3V5h-3M3 13h7v7H3v-7m2 2v3h3v-3H5m13 0h2v2h-2v-2m-2 2h2v2h-2v-2m2 2h2v2h-2v-2m0-4h2v2h-2v-2z"/>
              </svg>
            </div>
            <p class="text-sm text-gray-600">QR Code for your address</p>
          </div>
        </div>
        
        <!-- Address Display -->
        <div class="p-4 bg-gray-50 rounded-lg break-all">
          <p class="text-sm text-gray-600 mb-2">Your address:</p>
          <p class="font-mono text-sm">{{ store.currentAddress }}</p>
        </div>
        
        <!-- Copy Button -->
        <button
          @click="copyAddress"
          class="w-full px-4 py-3 bg-gray-200 hover:bg-gray-300 rounded-lg transition flex items-center justify-center space-x-2"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>
          </svg>
          <span>Copy Address</span>
        </button>
      </div>
    </div>
  </div>
</template>