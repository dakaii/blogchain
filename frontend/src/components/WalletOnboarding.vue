<script setup lang="ts">
import { ref } from 'vue'
import { useBlockchainStore } from '@/stores/blockchain'
import { WalletIcon, KeyIcon, ArrowRightIcon, CheckIcon } from '@heroicons/vue/24/outline'

const store = useBlockchainStore()
const showMnemonicInput = ref(false)
const mnemonic = ref('')
const isConnecting = ref(false)

const steps = [
  { id: 1, name: 'Choose Wallet', complete: false },
  { id: 2, name: 'Connect', complete: false },
  { id: 3, name: 'Start Blogging', complete: false },
]

const connectKeplr = async () => {
  isConnecting.value = true
  try {
    await store.connectWithKeplr()
    steps[0].complete = true
    steps[1].complete = true
  } catch (error) {
    console.error('Keplr connection failed:', error)
  } finally {
    isConnecting.value = false
  }
}

const connectMnemonic = async () => {
  if (!mnemonic.value.trim()) return
  
  isConnecting.value = true
  try {
    await store.connectWallet(mnemonic.value)
    steps[0].complete = true
    steps[1].complete = true
    showMnemonicInput.value = false
    mnemonic.value = ''
  } catch (error) {
    console.error('Mnemonic connection failed:', error)
  } finally {
    isConnecting.value = false
  }
}
</script>

<template>
  <div class="min-h-[60vh] flex items-center justify-center px-4">
    <div class="max-w-2xl w-full">
      <!-- Progress Steps -->
      <nav aria-label="Progress" class="mb-12">
        <ol class="flex items-center justify-center space-x-5">
          <li v-for="(step, index) in steps" :key="step.id" class="relative">
            <div v-if="index !== steps.length - 1" 
                 class="absolute top-5 left-5 -ml-px mt-0.5 h-0.5 w-full bg-gray-300 dark:bg-dark-border" 
                 :class="{ 'bg-accent-blue dark:bg-accent-blue': step.complete }" />
            
            <div class="group relative flex items-center">
              <span class="h-10 w-10 rounded-full flex items-center justify-center transition-all duration-300"
                    :class="{
                      'bg-accent-blue text-white': step.complete,
                      'bg-white dark:bg-dark-surface border-2 border-gray-300 dark:border-dark-border': !step.complete && !store.isConnected,
                      'ring-4 ring-accent-blue/20': store.isConnected && index === 2
                    }">
                <CheckIcon v-if="step.complete" class="w-5 h-5" />
                <span v-else>{{ step.id }}</span>
              </span>
              <span class="ml-3 text-sm font-medium text-gray-900 dark:text-gray-100">{{ step.name }}</span>
            </div>
          </li>
        </ol>
      </nav>

      <!-- Welcome Card -->
      <div class="bg-white dark:bg-dark-surface rounded-2xl shadow-xl border border-gray-200 dark:border-dark-border p-8 animate-slide-up">
        <div class="text-center mb-8">
          <h1 class="text-3xl font-bold text-gray-900 dark:text-gray-100 mb-3">
            Welcome to BlogChain
          </h1>
          <p class="text-gray-600 dark:text-gray-400">
            Connect your wallet to start sharing your thoughts on the blockchain
          </p>
        </div>

        <!-- Wallet Options -->
        <div class="space-y-4">
          <!-- Keplr Option -->
          <button
            @click="connectKeplr"
            :disabled="isConnecting"
            class="w-full group relative overflow-hidden rounded-xl p-6 text-left transition-all duration-300
                   bg-gradient-to-r from-accent-purple/10 to-accent-blue/10 
                   hover:from-accent-purple/20 hover:to-accent-blue/20
                   border border-gray-200 dark:border-dark-border hover:border-accent-blue/50
                   disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-4">
                <div class="p-3 bg-white dark:bg-dark-elevated rounded-lg shadow-sm">
                  <WalletIcon class="w-6 h-6 text-accent-purple" />
                </div>
                <div>
                  <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                    Connect with Keplr
                  </h3>
                  <p class="text-sm text-gray-600 dark:text-gray-400">
                    Use your Keplr browser extension
                  </p>
                </div>
              </div>
              <ArrowRightIcon class="w-5 h-5 text-gray-400 group-hover:text-accent-blue transition-colors" />
            </div>
          </button>

          <!-- Mnemonic Option -->
          <div class="relative">
            <button
              @click="showMnemonicInput = !showMnemonicInput"
              :disabled="isConnecting"
              class="w-full group relative overflow-hidden rounded-xl p-6 text-left transition-all duration-300
                     bg-gradient-to-r from-accent-orange/10 to-accent-pink/10 
                     hover:from-accent-orange/20 hover:to-accent-pink/20
                     border border-gray-200 dark:border-dark-border hover:border-accent-orange/50
                     disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <div class="flex items-center justify-between">
                <div class="flex items-center space-x-4">
                  <div class="p-3 bg-white dark:bg-dark-elevated rounded-lg shadow-sm">
                    <KeyIcon class="w-6 h-6 text-accent-orange" />
                  </div>
                  <div>
                    <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                      Import with Seed Phrase
                    </h3>
                    <p class="text-sm text-gray-600 dark:text-gray-400">
                      Use your 12 or 24 word mnemonic
                    </p>
                  </div>
                </div>
                <ArrowRightIcon 
                  class="w-5 h-5 text-gray-400 group-hover:text-accent-orange transition-all duration-300"
                  :class="{ 'rotate-90': showMnemonicInput }"
                />
              </div>
            </button>

            <!-- Mnemonic Input -->
            <Transition
              enter-active-class="transition-all duration-300"
              enter-from-class="opacity-0 -translate-y-2"
              enter-to-class="opacity-100 translate-y-0"
              leave-active-class="transition-all duration-300"
              leave-from-class="opacity-100 translate-y-0"
              leave-to-class="opacity-0 -translate-y-2"
            >
              <div v-if="showMnemonicInput" class="mt-4 p-6 bg-gray-50 dark:bg-dark-bg rounded-xl border border-gray-200 dark:border-dark-border">
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Enter your seed phrase
                </label>
                <textarea
                  v-model="mnemonic"
                  rows="3"
                  placeholder="Enter your 12 or 24 word mnemonic phrase..."
                  class="w-full rounded-lg border border-gray-300 dark:border-dark-border 
                         bg-white dark:bg-dark-surface px-4 py-3
                         text-gray-900 dark:text-gray-100 
                         placeholder-gray-400 dark:placeholder-gray-500
                         focus:outline-none focus:ring-2 focus:ring-accent-orange/50"
                />
                <div class="mt-4 flex justify-end space-x-3">
                  <button
                    @click="showMnemonicInput = false; mnemonic = ''"
                    class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100"
                  >
                    Cancel
                  </button>
                  <button
                    @click="connectMnemonic"
                    :disabled="!mnemonic.trim() || isConnecting"
                    class="px-6 py-2 bg-gradient-to-r from-accent-orange to-accent-pink 
                           text-white rounded-lg font-medium
                           hover:shadow-lg transform hover:-translate-y-0.5 transition-all duration-200
                           disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Connect Wallet
                  </button>
                </div>
              </div>
            </Transition>
          </div>
        </div>

        <!-- Test Account Info -->
        <div class="mt-8 p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800/50">
          <div class="flex items-start space-x-3">
            <div class="flex-shrink-0">
              <svg class="w-5 h-5 text-blue-600 dark:text-blue-400" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
              </svg>
            </div>
            <div class="text-sm text-blue-800 dark:text-blue-200">
              <p class="font-semibold mb-1">Test Account Available</p>
              <p class="text-xs font-mono break-all opacity-75">
                banner spread envelope side kite person disagree path silver will brother under couch edit food venture squirrel civil budget number acquire point work mass
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>