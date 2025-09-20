<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { useBlockchainStore } from '@/stores/blockchain'
import PostCard from '@/components/PostCard.vue'

const store = useBlockchainStore()

onMounted(async () => {
  await store.connect()
  await store.fetchPosts()
})
</script>

<template>
  <div>
    <div class="mb-8">
      <h1 class="text-4xl font-bold text-gray-800 mb-2">Welcome to BlogChain</h1>
      <p class="text-gray-600">A decentralized blogging platform powered by blockchain</p>
    </div>
    
    <div v-if="store.loading" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
    </div>
    
    <div v-else-if="store.posts.length === 0" class="text-center py-12">
      <p class="text-gray-500 mb-4">No posts yet. Be the first to create one!</p>
      <RouterLink to="/create" class="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
        Create First Post
      </RouterLink>
    </div>
    
    <div v-else class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
      <PostCard v-for="post in store.posts" :key="post.id" :post="post" />
    </div>
  </div>
</template>