<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { useBlockchainStore } from '@/stores/blockchain'
import type { Post } from '@/services/blockchain'

const route = useRoute()
const store = useBlockchainStore()
const post = ref<Post | null>(null)
const loading = ref(true)

const formattedDate = computed(() => {
  if (post.value?.createdAt) {
    return new Date(parseInt(post.value.createdAt) * 1000).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })
  }
  return 'Unknown date'
})

onMounted(async () => {
  try {
    await store.connect()
    const response = await store.service.getPost(route.params.id as string)
    post.value = response.post
  } catch (error) {
    console.error('Failed to load post:', error)
  } finally {
    loading.value = false
  }
})

async function handleLike() {
  if (!post.value) return
  await store.likePost(post.value.id)
  // Refresh the post to get updated like count
  const response = await store.service.getPost(post.value.id)
  post.value = response.post
}
</script>

<template>
  <div class="max-w-4xl mx-auto">
    <div v-if="loading" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
    </div>
    
    <div v-else-if="!post" class="text-center py-12">
      <p class="text-gray-500 mb-4">Post not found</p>
      <RouterLink to="/" class="text-blue-600 hover:text-blue-800">
        ← Back to Home
      </RouterLink>
    </div>
    
    <article v-else class="bg-white rounded-lg shadow-lg p-8">
      <header class="mb-6">
        <h1 class="text-4xl font-bold text-gray-800 mb-4">{{ post.title }}</h1>
        
        <div class="flex items-center justify-between text-gray-600">
          <div class="flex items-center space-x-4">
            <span>{{ formattedDate }}</span>
            <span class="text-sm">
              by {{ post.creator.slice(0, 10) }}...{{ post.creator.slice(-4) }}
            </span>
          </div>
          
          <button
            @click="handleLike"
            class="flex items-center space-x-2 px-4 py-2 bg-red-50 text-red-600 rounded-lg hover:bg-red-100 transition"
            :disabled="!store.currentAddress"
          >
            <span>❤️</span>
            <span>{{ post.likes || 0 }} likes</span>
          </button>
        </div>
      </header>
      
      <div v-if="post.tags && post.tags.length" class="flex flex-wrap gap-2 mb-6">
        <span
          v-for="tag in post.tags"
          :key="tag"
          class="px-3 py-1 bg-blue-100 text-blue-800 text-sm rounded-full"
        >
          {{ tag }}
        </span>
      </div>
      
      <div class="prose prose-lg max-w-none">
        <p class="whitespace-pre-wrap">{{ post.body }}</p>
      </div>
      
      <footer class="mt-8 pt-6 border-t">
        <RouterLink to="/" class="text-blue-600 hover:text-blue-800">
          ← Back to Home
        </RouterLink>
      </footer>
    </article>
  </div>
</template>