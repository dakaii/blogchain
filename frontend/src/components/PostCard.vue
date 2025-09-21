<script setup lang="ts">
import type { Post } from '@/services/blockchain';
import { useBlockchainStore } from '@/stores/blockchain';
import { computed, ref } from 'vue';
import { RouterLink } from 'vue-router';
import { HeartIcon, CalendarIcon, UserIcon, ChatBubbleBottomCenterTextIcon } from '@heroicons/vue/24/outline';
import { HeartIcon as HeartSolidIcon } from '@heroicons/vue/24/solid';

const props = defineProps<{
  post: Post
}>()

const store = useBlockchainStore()
const isLiking = ref(false)
const hasLiked = ref(false)

const truncatedBody = computed(() => {
  if (props.post.body.length > 150) {
    return props.post.body.substring(0, 150) + '...'
  }
  return props.post.body
})

const formattedDate = computed(() => {
  if (props.post.createdAt) {
    const date = new Date(parseInt(props.post.createdAt) * 1000)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const days = Math.floor(diff / (1000 * 60 * 60 * 24))
    
    if (days === 0) return 'Today'
    if (days === 1) return 'Yesterday'
    if (days < 7) return `${days} days ago`
    if (days < 30) return `${Math.floor(days / 7)} weeks ago`
    if (days < 365) return `${Math.floor(days / 30)} months ago`
    return date.toLocaleDateString()
  }
  return 'Unknown date'
})

const shortAddress = computed(() => {
  if (props.post.creator) {
    return `${props.post.creator.slice(0, 8)}...${props.post.creator.slice(-4)}`
  }
  return 'Anonymous'
})

async function handleLike() {
  if (!store.currentAddress) {
    window.$toast?.info('Connect Wallet', 'Please connect your wallet to like posts')
    return
  }
  
  isLiking.value = true
  const toastId = window.$toast?.pending('Liking Post', 'Broadcasting transaction...')
  
  try {
    await store.likePost(props.post.id)
    hasLiked.value = true
    window.$toast?.update(toastId, {
      type: 'success',
      title: 'Post Liked!',
      message: 'Your like has been recorded on the blockchain'
    })
  } catch (error) {
    window.$toast?.update(toastId, {
      type: 'error',
      title: 'Failed to Like',
      message: error.message || 'Transaction failed'
    })
  } finally {
    isLiking.value = false
  }
}
</script>

<template>
  <article class="group relative bg-white dark:bg-dark-surface rounded-2xl shadow-sm 
                  border border-gray-200 dark:border-dark-border
                  hover:shadow-xl hover:border-accent-blue/30 dark:hover:border-accent-blue/40
                  transition-all duration-300 overflow-hidden">
    <!-- Gradient accent on hover -->
    <div class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-accent-purple to-accent-blue 
                opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
    
    <div class="p-6">
      <!-- Header -->
      <div class="flex items-start justify-between mb-4">
        <div class="flex-1">
          <h2 class="text-xl font-bold text-gray-900 dark:text-gray-100 mb-2 
                     group-hover:text-accent-blue dark:group-hover:text-accent-blue 
                     transition-colors duration-200">
            {{ post.title }}
          </h2>
          
          <!-- Meta info -->
          <div class="flex items-center space-x-4 text-sm text-gray-500 dark:text-gray-400">
            <div class="flex items-center space-x-1">
              <UserIcon class="w-4 h-4" />
              <span class="font-mono">{{ shortAddress }}</span>
            </div>
            <div class="flex items-center space-x-1">
              <CalendarIcon class="w-4 h-4" />
              <span>{{ formattedDate }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Body -->
      <p class="text-gray-600 dark:text-gray-300 mb-4 line-clamp-3">
        {{ truncatedBody }}
      </p>

      <!-- Tags -->
      <div v-if="post.tags && post.tags.length" class="flex flex-wrap gap-2 mb-4">
        <span v-for="tag in post.tags" :key="tag" 
              class="px-3 py-1 text-xs font-medium rounded-full
                     bg-gradient-to-r from-accent-purple/10 to-accent-blue/10
                     dark:from-accent-purple/20 dark:to-accent-blue/20
                     text-accent-blue dark:text-accent-blue
                     border border-accent-blue/20 dark:border-accent-blue/30">
          #{{ tag }}
        </span>
      </div>

      <!-- Actions -->
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <!-- Like Button -->
          <button 
            @click="handleLike"
            :disabled="isLiking"
            class="flex items-center space-x-2 px-3 py-1.5 rounded-lg
                   transition-all duration-200 group/like"
            :class="{
              'bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400': hasLiked,
              'hover:bg-gray-100 dark:hover:bg-dark-elevated text-gray-600 dark:text-gray-400': !hasLiked && !isLiking,
              'opacity-50 cursor-wait': isLiking
            }"
            :aria-label="hasLiked ? 'Unlike post' : 'Like post'"
          >
            <HeartSolidIcon v-if="hasLiked" class="w-5 h-5 animate-pulse-slow" />
            <HeartIcon v-else class="w-5 h-5 group-hover/like:scale-110 transition-transform" />
            <span class="font-semibold">{{ post.likes || 0 }} {{ parseInt(post.likes || '0') === 1 ? 'like' : 'likes' }}</span>
          </button>

          <!-- Comments (placeholder for future) -->
          <button class="flex items-center space-x-2 px-3 py-1.5 rounded-lg
                         text-gray-600 dark:text-gray-400
                         hover:bg-gray-100 dark:hover:bg-dark-elevated
                         transition-all duration-200"
                  disabled>
            <ChatBubbleBottomCenterTextIcon class="w-5 h-5" />
            <span class="text-sm">0</span>
          </button>
        </div>

        <!-- Read More Link -->
        <RouterLink 
          :to="`/post/${post.id}`" 
          class="flex items-center space-x-1 text-sm font-medium
                 text-accent-blue hover:text-accent-purple
                 transition-colors duration-200 group/link"
        >
          <span>Read more</span>
          <svg class="w-4 h-4 group-hover/link:translate-x-1 transition-transform" 
               fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </RouterLink>
      </div>
    </div>
  </article>
</template>

<style scoped>
.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>