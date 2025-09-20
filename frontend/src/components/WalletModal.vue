<script setup lang="ts">
import { useBlockchainStore } from '@/stores/blockchain'
import { ref } from 'vue'

const emit = defineEmits(['close'])
const store = useBlockchainStore()
const mnemonic = ref('')
const loading = ref(false)
const connectionType = ref<'keplr' | 'mnemonic'>('mnemonic')

async function connectWallet() {
  loading.value = true
  
  if (connectionType.value === 'keplr') {
    await store.connectWithKeplr()
  } else {
    if (!mnemonic.value.trim()) {
      alert('Please enter a mnemonic')
      loading.value = false
      return
    }
    await store.connectWallet(mnemonic.value)
  }
  
  loading.value = false

  if (!store.error) {
    emit('close')
  }
}

function useTestAccount() {
  mnemonic.value = 'banner spread envelope side kite person disagree path silver will brother under couch edit food venture squirrel civil budget number acquire point work mass'
}
</script>

<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg p-6 w-full max-w-md">
      <h2 class="text-2xl font-bold mb-4">Connect Wallet</h2>

      <div class="space-y-4">
        <!-- Connection Type Selector -->
        <div class="flex space-x-2 p-1 bg-gray-100 rounded-lg">
          <button
            @click="connectionType = 'keplr'"
            :class="[
              'flex-1 px-3 py-2 rounded-md transition',
              connectionType === 'keplr' 
                ? 'bg-white shadow text-blue-600' 
                : 'text-gray-600 hover:text-gray-800'
            ]"
          >
            Keplr Wallet
          </button>
          <button
            @click="connectionType = 'mnemonic'"
            :class="[
              'flex-1 px-3 py-2 rounded-md transition',
              connectionType === 'mnemonic' 
                ? 'bg-white shadow text-blue-600' 
                : 'text-gray-600 hover:text-gray-800'
            ]"
          >
            Mnemonic
          </button>
        </div>
        
        <!-- Keplr Connection -->
        <div v-if="connectionType === 'keplr'" class="text-center py-4">
          <p class="text-gray-600 mb-4">
            Connect with your Keplr wallet extension
          </p>
          <p class="text-sm text-gray-500">
            Make sure Keplr extension is installed
          </p>
        </div>
        
        <!-- Mnemonic Input -->
        <div v-else>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Mnemonic Phrase
          </label>
          <textarea v-model="mnemonic"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            rows="3" placeholder="Enter your 24-word mnemonic phrase..."></textarea>
          
          <button @click="useTestAccount" class="w-full px-4 py-2 text-sm text-gray-600 hover:text-gray-800 mt-2">
            Use Test Account (Alice)
          </button>
        </div>

        <div v-if="store.error" class="text-red-600 text-sm">
          {{ store.error }}
        </div>

        <div class="flex space-x-2">
          <button @click="connectWallet" :disabled="loading"
            class="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50">
            {{ loading ? 'Connecting...' : 'Connect' }}
          </button>
          <button @click="emit('close')"
            class="flex-1 px-4 py-2 bg-gray-200 text-gray-800 rounded-lg hover:bg-gray-300">
            Cancel
          </button>
        </div>
      </div>
    </div>
  </div>
</template>