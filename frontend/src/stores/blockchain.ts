import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { BlockchainService, type Post } from '@/services/blockchain'

export const useBlockchainStore = defineStore('blockchain', () => {
  const service = new BlockchainService()
  const isConnected = ref(false)
  const currentAddress = ref('')
  const balance = ref<{ denom: string; amount: string }[]>([])
  const posts = ref<Post[]>([])
  const loading = ref(false)
  const error = ref('')

  const formattedBalance = computed(() => {
    const stakeBalance = balance.value.find(b => b.denom === 'stake')
    if (stakeBalance) {
      return `${(parseInt(stakeBalance.amount) / 1000000).toFixed(2)} STAKE`
    }
    return '0 STAKE'
  })

  async function connect() {
    try {
      loading.value = true
      error.value = ''
      await service.connect()
      isConnected.value = true
    } catch (e) {
      error.value = `Failed to connect: ${e}`
      console.error(e)
    } finally {
      loading.value = false
    }
  }

  async function connectWallet(mnemonic: string) {
    try {
      loading.value = true
      error.value = ''
      const { address } = await service.connectWithSigner(mnemonic)
      currentAddress.value = address
      isConnected.value = true
      await updateBalance()
    } catch (e) {
      error.value = `Failed to connect wallet: ${e}`
      console.error(e)
    } finally {
      loading.value = false
    }
  }

  async function connectWithKeplr() {
    try {
      loading.value = true
      error.value = ''
      const { address } = await service.connectWithKeplr()
      currentAddress.value = address
      isConnected.value = true
      await updateBalance()
    } catch (e) {
      error.value = `Failed to connect Keplr: ${e}`
      console.error(e)
    } finally {
      loading.value = false
    }
  }

  async function updateBalance() {
    if (!currentAddress.value) return
    try {
      const balances = await service.getBalance(currentAddress.value)
      balance.value = balances.map(coin => ({
        denom: coin.denom,
        amount: coin.amount
      }))
    } catch (e) {
      console.error('Failed to update balance:', e)
    }
  }

  async function fetchPosts() {
    try {
      loading.value = true
      error.value = ''
      const response = await service.getPosts({ limit: 50 })
      posts.value = response.posts || []
    } catch (e) {
      error.value = `Failed to fetch posts: ${e}`
      console.error(e)
    } finally {
      loading.value = false
    }
  }

  async function createPost(title: string, body: string, tags: string[]) {
    if (!currentAddress.value) {
      error.value = 'Please connect your wallet first'
      return
    }
    
    try {
      loading.value = true
      error.value = ''
      await service.createPost(currentAddress.value, title, body, tags)
      await fetchPosts()
      await updateBalance()
    } catch (e) {
      error.value = `Failed to create post: ${e}`
      console.error(e)
    } finally {
      loading.value = false
    }
  }

  async function likePost(postId: string) {
    if (!currentAddress.value) {
      error.value = 'Please connect your wallet first'
      return
    }
    
    try {
      loading.value = true
      error.value = ''
      await service.likePost(currentAddress.value, postId)
      await fetchPosts()
      await updateBalance()
    } catch (e) {
      error.value = `Failed to like post: ${e}`
      console.error(e)
    } finally {
      loading.value = false
    }
  }

  return {
    service,
    isConnected,
    currentAddress,
    balance,
    formattedBalance,
    posts,
    loading,
    error,
    connect,
    connectWallet,
    connectWithKeplr,
    updateBalance,
    fetchPosts,
    createPost,
    likePost
  }
})