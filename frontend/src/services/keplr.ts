import type { Window as KeplrWindow } from '@keplr-wallet/types'

declare global {
  interface Window extends KeplrWindow {}
}

export interface ChainInfo {
  chainId: string
  chainName: string
  rpc: string
  rest: string
  stakeCurrency: {
    coinDenom: string
    coinMinimalDenom: string
    coinDecimals: number
  }
  bip44: {
    coinType: number
  }
  bech32Config: {
    bech32PrefixAccAddr: string
    bech32PrefixAccPub: string
    bech32PrefixValAddr: string
    bech32PrefixValPub: string
    bech32PrefixConsAddr: string
    bech32PrefixConsPub: string
  }
  currencies: Array<{
    coinDenom: string
    coinMinimalDenom: string
    coinDecimals: number
  }>
  feeCurrencies: Array<{
    coinDenom: string
    coinMinimalDenom: string
    coinDecimals: number
    gasPriceStep: {
      low: number
      average: number
      high: number
    }
  }>
}

export class KeplrService {
  private chainInfo: ChainInfo = {
    chainId: 'blogchain',
    chainName: 'BlogChain',
    rpc: 'http://localhost:26657',
    rest: 'http://localhost:1317',
    stakeCurrency: {
      coinDenom: 'STAKE',
      coinMinimalDenom: 'stake',
      coinDecimals: 6,
    },
    bip44: {
      coinType: 118,
    },
    bech32Config: {
      bech32PrefixAccAddr: 'blogchain',
      bech32PrefixAccPub: 'blogchainpub',
      bech32PrefixValAddr: 'blogchainvaloper',
      bech32PrefixValPub: 'blogchainvaloperpub',
      bech32PrefixConsAddr: 'blogchainvalcons',
      bech32PrefixConsPub: 'blogchainvalconspub',
    },
    currencies: [
      {
        coinDenom: 'STAKE',
        coinMinimalDenom: 'stake',
        coinDecimals: 6,
      },
      {
        coinDenom: 'TOKEN',
        coinMinimalDenom: 'token',
        coinDecimals: 6,
      },
    ],
    feeCurrencies: [
      {
        coinDenom: 'STAKE',
        coinMinimalDenom: 'stake',
        coinDecimals: 6,
        gasPriceStep: {
          low: 0.01,
          average: 0.025,
          high: 0.04,
        },
      },
    ],
  }

  async isKeplrInstalled(): Promise<boolean> {
    return !!(window.keplr && window.getOfflineSigner)
  }

  async connectKeplr() {
    if (!await this.isKeplrInstalled()) {
      throw new Error('Keplr extension not found. Please install Keplr wallet.')
    }

    try {
      // Try to suggest the chain to Keplr
      await window.keplr!.experimentalSuggestChain(this.chainInfo as any)
      
      // Enable the chain
      await window.keplr!.enable(this.chainInfo.chainId)
      
      // Get the offline signer
      const offlineSigner = window.getOfflineSigner!(this.chainInfo.chainId)
      
      // Get account
      const accounts = await offlineSigner.getAccounts()
      
      return {
        offlineSigner,
        accounts,
        address: accounts[0].address,
      }
    } catch (error) {
      console.error('Failed to connect to Keplr:', error)
      throw error
    }
  }

  async getKeplrSigner() {
    if (!await this.isKeplrInstalled()) {
      throw new Error('Keplr extension not found')
    }
    
    await window.keplr!.enable(this.chainInfo.chainId)
    return window.getOfflineSigner!(this.chainInfo.chainId)
  }

  async getKey() {
    if (!await this.isKeplrInstalled()) {
      throw new Error('Keplr extension not found')
    }
    
    return await window.keplr!.getKey(this.chainInfo.chainId)
  }

  async signAndBroadcast(msgs: any[], fee: any, memo: string) {
    const key = await this.getKey()
    const offlineSigner = await this.getKeplrSigner()
    
    // This would need to be implemented with SigningStargateClient
    // Similar to what we have in blockchain.ts
    return { key, offlineSigner }
  }
}