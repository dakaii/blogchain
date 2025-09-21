import { describe, it, expect, beforeEach, vi } from 'vitest'
import { BlockchainService } from '../blockchain'

// Mock the cosmos modules
vi.mock('@cosmjs/stargate', () => ({
  StargateClient: {
    connect: vi.fn().mockResolvedValue({
      getAllBalances: vi.fn().mockResolvedValue([
        { denom: 'stake', amount: '1000000' },
        { denom: 'token', amount: '500000' }
      ])
    })
  },
  SigningStargateClient: {
    connectWithSigner: vi.fn().mockResolvedValue({
      signAndBroadcast: vi.fn().mockResolvedValue({
        code: 0,
        transactionHash: 'ABC123'
      })
    })
  },
  GasPrice: {
    fromString: vi.fn()
  }
}))

vi.mock('@cosmjs/proto-signing', () => ({
  DirectSecp256k1HdWallet: {
    fromMnemonic: vi.fn().mockResolvedValue({
      getAccounts: vi.fn().mockResolvedValue([
        { address: 'blogchain1test...', pubkey: new Uint8Array() }
      ])
    })
  }
}))

describe('BlockchainService', () => {
  let service: BlockchainService

  beforeEach(() => {
    service = new BlockchainService()
    vi.clearAllMocks()
  })

  describe('connect', () => {
    it('should connect to the blockchain', async () => {
      const client = await service.connect()
      expect(client).toBeDefined()
    })
  })

  describe('connectWithSigner', () => {
    it('should connect with mnemonic', async () => {
      const testMnemonic = 'test mnemonic phrase'
      const result = await service.connectWithSigner(testMnemonic)
      
      expect(result.client).toBeDefined()
      expect(result.address).toBe('blogchain1test...')
    })
  })

  describe('getBalance', () => {
    it('should fetch balance for an address', async () => {
      const address = 'blogchain1test...'
      const balances = await service.getBalance(address)
      
      expect(balances).toHaveLength(2)
      expect(balances[0].denom).toBe('stake')
      expect(balances[0].amount).toBe('1000000')
    })
  })

  describe('getPosts', () => {
    it('should fetch posts from API', async () => {
      global.fetch = vi.fn().mockResolvedValue({
        json: vi.fn().mockResolvedValue({
          posts: [
            { id: '1', title: 'Test Post 1' },
            { id: '2', title: 'Test Post 2' }
          ]
        })
      })

      const result = await service.getPosts({ limit: 10, offset: 0 })
      
      expect(fetch).toHaveBeenCalledWith('/api/blogchain/blog/v1/posts?pagination.limit=10&pagination.offset=0')
      expect(result.posts).toHaveLength(2)
    })
  })

  describe('getPost', () => {
    it('should fetch a single post', async () => {
      const mockPost = { id: '1', title: 'Test Post' }
      global.fetch = vi.fn().mockResolvedValue({
        json: vi.fn().mockResolvedValue(mockPost)
      })

      const result = await service.getPost('1')
      
      expect(fetch).toHaveBeenCalledWith('/api/blogchain/blog/v1/posts/1')
      expect(result).toEqual(mockPost)
    })
  })

  describe('createPost', () => {
    it('should throw error if signing client not initialized', async () => {
      await expect(
        service.createPost('creator', 'Title', 'Body', ['tag'])
      ).rejects.toThrow('Signing client not initialized')
    })

    it('should create a post when connected', async () => {
      await service.connectWithSigner('test mnemonic')
      
      const result = await service.createPost(
        'blogchain1test...',
        'Test Title',
        'Test Body',
        ['test', 'blockchain']
      )
      
      expect(result.code).toBe(0)
      expect(result.transactionHash).toBe('ABC123')
    })
  })

  describe('likePost', () => {
    it('should throw error if signing client not initialized', async () => {
      await expect(
        service.likePost('liker', '1')
      ).rejects.toThrow('Signing client not initialized')
    })
  })

  describe('sendTokens', () => {
    it('should throw error if signing client not initialized', async () => {
      await expect(
        service.sendTokens('from', 'to', '1000000')
      ).rejects.toThrow('Signing client not initialized')
    })

    it('should send tokens when connected', async () => {
      await service.connectWithSigner('test mnemonic')
      
      const result = await service.sendTokens(
        'blogchain1from...',
        'blogchain1to...',
        '1000000',
        'stake'
      )
      
      expect(result.code).toBe(0)
      expect(result.transactionHash).toBe('ABC123')
    })
  })
})