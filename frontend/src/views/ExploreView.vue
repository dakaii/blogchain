<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useBlockchainStore } from '@/stores/blockchain'
import PostCard from '@/components/PostCard.vue'

const store = useBlockchainStore()
const searchQuery = ref('')
const selectedTag = ref('')

onMounted(async () => {
  if (!store.posts.length) {
    await store.connect()
    await store.fetchPosts()
  }
})

const allTags = computed(() => {
  const tags = new Set<string>()
  store.posts.forEach(post => {
    if (post.tags) {
      post.tags.forEach(tag => tags.add(tag))
    }
  })
  return Array.from(tags)
})

const filteredPosts = computed(() => {
  let posts = [...store.posts]
  
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    posts = posts.filter(post => 
      post.title.toLowerCase().includes(query) ||
      post.body.toLowerCase().includes(query)
    )
  }
  
  if (selectedTag.value) {
    posts = posts.filter(post => 
      post.tags && post.tags.includes(selectedTag.value)
    )
  }
  
  return posts
})

const sortedPosts = computed(() => {
  return [...filteredPosts.value].sort((a, b) => {
    return parseInt(b.likes || '0') - parseInt(a.likes || '0')
  })
})
</script>

<template>
  <div>
    <div class="mb-8">
      <h1 class="text-4xl font-bold text-gray-800 mb-2">Explore Posts</h1>
      <p class="text-gray-600">Discover content from the BlogChain community</p>
    </div>
    
    <div class="mb-6 space-y-4">
      <div>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search posts..."
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>
      
      <div v-if="allTags.length" class="flex flex-wrap gap-2">
        <button
          @click="selectedTag = ''"
          :class="[
            'px-3 py-1 rounded-full text-sm transition',
            !selectedTag 
              ? 'bg-blue-600 text-white' 
              : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
          ]"
        >
          All
        </button>
        <button
          v-for="tag in allTags"
          :key="tag"
          @click="selectedTag = tag"
          :class="[
            'px-3 py-1 rounded-full text-sm transition',
            selectedTag === tag 
              ? 'bg-blue-600 text-white' 
              : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
          ]"
        >
          {{ tag }}
        </button>
      </div>
    </div>
    
    <div v-if="store.loading" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
    </div>
    
    <div v-else-if="sortedPosts.length === 0" class="text-center py-12">
      <p class="text-gray-500">No posts found matching your criteria.</p>
    </div>
    
    <div v-else class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
      <PostCard v-for="post in sortedPosts" :key="post.id" :post="post" />
    </div>
  </div>
</template>