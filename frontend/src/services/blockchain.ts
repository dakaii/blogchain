import { SigningStargateClient, StargateClient } from '@cosmjs/stargate'
import { DirectSecp256k1HdWallet } from '@cosmjs/proto-signing'
import type { OfflineDirectSigner } from '@cosmjs/proto-signing'
import { GasPrice } from '@cosmjs/stargate'
import { KeplrService } from './keplr'

export interface Post {
  id: string
  creator: string
  title: string
  body: string
  tags: string[]
  createdAt: string
  likes: string
}

export class BlockchainService {
  private client: StargateClient | null = null
  private signingClient: SigningStargateClient | null = null
  private chainId = 'blogchain'
  private rpcUrl = 'http://localhost:26657'
  private apiUrl = '/api'
  private keplr = new KeplrService()

  async connect() {
    this.client = await StargateClient.connect(this.rpcUrl)
    return this.client
  }

  async connectWithSigner(mnemonic: string) {
    const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic, {
      prefix: 'blogchain'
    })
    
    const [account] = await wallet.getAccounts()
    
    this.signingClient = await SigningStargateClient.connectWithSigner(
      this.rpcUrl,
      wallet,
      {
        gasPrice: GasPrice.fromString('0.025stake')
      }
    )
    
    return { client: this.signingClient, address: account.address }
  }

  async connectWithKeplr() {
    const { offlineSigner, address } = await this.keplr.connectKeplr()
    
    this.signingClient = await SigningStargateClient.connectWithSigner(
      this.rpcUrl,
      offlineSigner as OfflineDirectSigner,
      {
        gasPrice: GasPrice.fromString('0.025stake')
      }
    )
    
    return { client: this.signingClient, address }
  }

  async getBalance(address: string) {
    if (!this.client) await this.connect()
    return await this.client!.getAllBalances(address)
  }

  async getPosts(pagination?: { limit?: number; offset?: number }) {
    const params = new URLSearchParams()
    if (pagination?.limit) params.append('pagination.limit', pagination.limit.toString())
    if (pagination?.offset) params.append('pagination.offset', pagination.offset.toString())
    
    const response = await fetch(`${this.apiUrl}/blogchain/blog/v1/posts?${params}`)
    const data = await response.json()
    return data
  }

  async getPost(id: string) {
    const response = await fetch(`${this.apiUrl}/blogchain/blog/v1/posts/${id}`)
    const data = await response.json()
    return data
  }

  async createPost(creator: string, title: string, body: string, tags: string[] = []) {
    if (!this.signingClient) throw new Error('Signing client not initialized')
    
    const msg = {
      typeUrl: '/blogchain.blog.v1.MsgCreatePost',
      value: {
        creator,
        title,
        body,
        tags
      }
    }
    
    const fee = {
      amount: [{ denom: 'stake', amount: '5000' }],
      gas: '200000'
    }
    
    const result = await this.signingClient.signAndBroadcast(
      creator,
      [msg],
      fee,
      'Creating blog post'
    )
    
    return result
  }

  async likePost(liker: string, postId: string) {
    if (!this.signingClient) throw new Error('Signing client not initialized')
    
    const msg = {
      typeUrl: '/blogchain.blog.v1.MsgLikePost',
      value: {
        liker,
        postId
      }
    }
    
    const fee = {
      amount: [{ denom: 'stake', amount: '2500' }],
      gas: '100000'
    }
    
    const result = await this.signingClient.signAndBroadcast(
      liker,
      [msg],
      fee,
      'Liking post'
    )
    
    return result
  }

  async getTransactionHistory(address: string, limit = 10) {
    const response = await fetch(
      `${this.apiUrl}/cosmos/tx/v1beta1/txs?events=message.sender='${address}'&order_by=ORDER_BY_DESC&limit=${limit}`
    )
    const data = await response.json()
    return data.tx_responses || []
  }

  async sendTokens(fromAddress: string, toAddress: string, amount: string, denom = 'stake') {
    if (!this.signingClient) throw new Error('Signing client not initialized')
    
    const msg = {
      typeUrl: '/cosmos.bank.v1beta1.MsgSend',
      value: {
        fromAddress,
        toAddress,
        amount: [{ denom, amount }]
      }
    }
    
    const fee = {
      amount: [{ denom: 'stake', amount: '5000' }],
      gas: '200000'
    }
    
    const result = await this.signingClient.signAndBroadcast(
      fromAddress,
      [msg],
      fee,
      'Token transfer'
    )
    
    return result
  }

  async isKeplrAvailable() {
    return await this.keplr.isKeplrInstalled()
  }
}