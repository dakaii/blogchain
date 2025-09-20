<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useBlockchainStore } from '@/stores/blockchain'

const router = useRouter()
const store = useBlockchainStore()

const title = ref('')
const body = ref('')
const tags = ref('')
const submitting = ref(false)

async function handleSubmit() {
  if (!store.currentAddress) {
    store.error = 'Please connect your wallet first'
    return
  }
  
  if (!title.value.trim() || !body.value.trim()) {
    alert('Please fill in all required fields')
    return
  }
  
  submitting.value = true
  const tagList = tags.value.split(',').map(t => t.trim()).filter(t => t)
  
  await store.createPost(title.value, body.value, tagList)
  
  submitting.value = false
  
  if (!store.error) {
    router.push('/')
  }
}
</script>

<template>
  <div class="max-w-4xl mx-auto">
    <h1 class="text-3xl font-bold text-gray-800 mb-8">Create New Post</h1>
    
    <div v-if="!store.currentAddress" class="bg-yellow-100 border-l-4 border-yellow-500 text-yellow-700 p-4 mb-6">
      <p>Please connect your wallet to create a post.</p>
    </div>
    
    <form @submit.prevent="handleSubmit" class="space-y-6">
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">
          Title *
        </label>
        <input
          v-model="title"
          type="text"
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="Enter your post title..."
        />
      </div>
      
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">
          Content *
        </label>
        <textarea
          v-model="body"
          rows="10"
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="Write your post content..."
        ></textarea>
      </div>
      
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">
          Tags (comma-separated)
        </label>
        <input
          v-model="tags"
          type="text"
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="blockchain, web3, technology..."
        />
      </div>
      
      <div v-if="store.error" class="text-red-600">
        {{ store.error }}
      </div>
      
      <div class="flex space-x-4">
        <button
          type="submit"
          :disabled="submitting || !store.currentAddress"
          class="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {{ submitting ? 'Publishing...' : 'Publish Post' }}
        </button>
        <RouterLink
          to="/"
          class="px-6 py-3 bg-gray-200 text-gray-800 rounded-lg hover:bg-gray-300"
        >
          Cancel
        </RouterLink>
      </div>
    </form>
  </div>
</template>