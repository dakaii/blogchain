<script setup lang="ts">
import type { Post } from '@/services/blockchain';
import { useBlockchainStore } from '@/stores/blockchain';
import { computed } from 'vue';
import { RouterLink } from 'vue-router';

const props = defineProps<{
  post: Post
}>()

const store = useBlockchainStore()

const truncatedBody = computed(() => {
  if (props.post.body.length > 150) {
    return props.post.body.substring(0, 150) + '...'
  }
  return props.post.body
})

const formattedDate = computed(() => {
  if (props.post.createdAt) {
    return new Date(parseInt(props.post.createdAt) * 1000).toLocaleDateString()
  }
  return 'Unknown date'
})

async function handleLike() {
  await store.likePost(props.post.id)
}
</script>

<template>
  <div class="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
    <h3 class="text-xl font-semibold text-gray-800 mb-2">{{ post.title }}</h3>
    <p class="text-gray-600 mb-4">{{ truncatedBody }}</p>

    <div v-if="post.tags && post.tags.length" class="flex flex-wrap gap-2 mb-4">
      <span v-for="tag in post.tags" :key="tag" class="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full">
        {{ tag }}
      </span>
    </div>

    <div class="flex items-center justify-between text-sm text-gray-500">
      <div class="flex items-center space-x-4">
        <button @click="handleLike" class="flex items-center space-x-1 hover:text-red-500 transition"
          :disabled="!store.currentAddress">
          <span>❤️</span>
          <span>{{ post.likes || 0 }}</span>
        </button>
        <span>{{ formattedDate }}</span>
      </div>
      <RouterLink :to="`/post/${post.id}`" class="text-blue-600 hover:text-blue-800">
        Read more →
      </RouterLink>
    </div>
  </div>
</template>