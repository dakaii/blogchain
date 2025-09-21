<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { RouterLink } from 'vue-router'
import { useBlockchainStore } from '@/stores/blockchain'
import PostCard from '@/components/PostCard.vue'
import WalletOnboarding from '@/components/WalletOnboarding.vue'
import { SparklesIcon, RocketLaunchIcon, ShieldCheckIcon } from '@heroicons/vue/24/outline'

const store = useBlockchainStore()

const features = [
  {
    icon: SparklesIcon,
    title: 'Decentralized Content',
    description: 'Your posts are stored forever on the blockchain, censorship-resistant and immutable.'
  },
  {
    icon: ShieldCheckIcon,
    title: 'True Ownership',
    description: 'You own your content and interactions. No platform can delete or modify your posts.'
  },
  {
    icon: RocketLaunchIcon,
    title: 'Web3 Native',
    description: 'Built on Cosmos SDK with native wallet integration and on-chain interactions.'
  }
]

const latestPosts = computed(() => {
  return [...store.posts].sort((a, b) => {
    const timeA = parseInt(a.createdAt || '0')
    const timeB = parseInt(b.createdAt || '0')
    return timeB - timeA
  }).slice(0, 6)
})

onMounted(async () => {
  await store.connect()
  await store.fetchPosts()
})
</script>

<template>
  <div>
    <!-- Hero Section -->
    <div v-if="!store.currentAddress" class="mb-12">
      <WalletOnboarding />
    </div>
    
    <!-- Connected User Welcome -->
    <div v-else class="mb-12">
      <div class="relative overflow-hidden rounded-3xl bg-gradient-to-br from-accent-purple via-accent-blue to-accent-green p-8 md:p-12">
        <div class="absolute inset-0 bg-black/20" />
        <div class="relative z-10">
          <h1 class="text-4xl md:text-5xl font-bold text-white mb-4">
            Welcome back to BlogChain
          </h1>
          <p class="text-xl text-white/90 mb-8 max-w-2xl">
            Share your thoughts on the decentralized web. Your content, your ownership, forever on-chain.
          </p>
          <div class="flex flex-wrap gap-4">
            <RouterLink to="/create" 
              class="inline-flex items-center space-x-2 px-6 py-3 
                     bg-white text-gray-900 rounded-xl font-semibold
                     hover:shadow-xl transform hover:-translate-y-0.5 transition-all duration-200">
              <PencilSquareIcon class="w-5 h-5" />
              <span>Create New Post</span>
            </RouterLink>
            <RouterLink to="/explore"
              class="inline-flex items-center space-x-2 px-6 py-3 
                     bg-white/10 backdrop-blur-sm text-white rounded-xl font-semibold
                     border border-white/20 hover:bg-white/20 transition-all duration-200">
              <MagnifyingGlassIcon class="w-5 h-5" />
              <span>Explore Posts</span>
            </RouterLink>
          </div>
        </div>
        
        <!-- Decorative elements -->
        <div class="absolute -top-24 -right-24 w-96 h-96 bg-white/10 rounded-full blur-3xl" />
        <div class="absolute -bottom-24 -left-24 w-96 h-96 bg-white/10 rounded-full blur-3xl" />
      </div>
    </div>

    <!-- Features Section -->
    <div class="grid md:grid-cols-3 gap-6 mb-12">
      <div v-for="feature in features" :key="feature.title"
        class="p-6 bg-white dark:bg-dark-surface rounded-2xl border border-gray-200 dark:border-dark-border
               hover:shadow-lg hover:border-accent-blue/30 transition-all duration-300">
        <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-accent-purple to-accent-blue 
                    flex items-center justify-center mb-4">
          <component :is="feature.icon" class="w-6 h-6 text-white" />
        </div>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-2">
          {{ feature.title }}
        </h3>
        <p class="text-sm text-gray-600 dark:text-gray-400">
          {{ feature.description }}
        </p>
      </div>
    </div>

    <!-- Latest Posts Section -->
    <div class="mb-8">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-gray-100">
          Latest Posts
        </h2>
        <RouterLink to="/explore" 
          class="text-accent-blue hover:text-accent-purple transition-colors duration-200">
          View all →
        </RouterLink>
      </div>

      <div v-if="store.loading" class="flex justify-center py-12">
        <div class="relative">
          <div class="animate-spin rounded-full h-12 w-12 border-4 border-gray-200 dark:border-dark-border" />
          <div class="absolute inset-0 animate-spin rounded-full h-12 w-12 border-t-4 border-accent-blue" />
        </div>
      </div>
      
      <div v-else-if="latestPosts.length === 0" 
        class="text-center py-16 bg-gray-50 dark:bg-dark-bg rounded-2xl border-2 border-dashed border-gray-300 dark:border-dark-border">
        <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-gray-100 dark:bg-dark-surface flex items-center justify-center">
          <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
        </div>
        <p class="text-gray-500 dark:text-gray-400 mb-4">No posts yet. Be the first to create one!</p>
        <RouterLink to="/create" 
          class="inline-flex items-center space-x-2 px-6 py-3 
                 bg-gradient-to-r from-accent-purple to-accent-blue 
                 text-white rounded-xl font-semibold
                 hover:shadow-lg transform hover:-translate-y-0.5 transition-all duration-200">
          <PencilSquareIcon class="w-5 h-5" />
          <span>Create First Post</span>
        </RouterLink>
      </div>
      
      <div v-else class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <PostCard v-for="post in latestPosts" :key="post.id" :post="post" />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { PencilSquareIcon, MagnifyingGlassIcon } from '@heroicons/vue/24/outline'
export default {
  components: {
    PencilSquareIcon,
    MagnifyingGlassIcon
  }
}
</script>