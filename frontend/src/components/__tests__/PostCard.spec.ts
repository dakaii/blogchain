import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import PostCard from '../PostCard.vue'
import { createTestingPinia } from '@pinia/testing'
import { useBlockchainStore } from '@/stores/blockchain'

describe('PostCard', () => {
  const mockPost = {
    id: '1',
    creator: 'blogchain1test...',
    title: 'Test Post Title',
    body: 'This is a test post body that should be truncated if it is too long...',
    tags: ['test', 'blockchain', 'vue'],
    createdAt: '1703001600',
    likes: '5'
  }

  it('renders post information correctly', () => {
    const wrapper = mount(PostCard, {
      props: {
        post: mockPost
      },
      global: {
        plugins: [createTestingPinia({
          createSpy: vi.fn
        })],
        stubs: {
          RouterLink: {
            template: '<a><slot /></a>'
          }
        }
      }
    })

    expect(wrapper.find('h2').text()).toBe(mockPost.title)
    expect(wrapper.text()).toContain('5 likes')
    expect(wrapper.text()).toContain('test')
    expect(wrapper.text()).toContain('blockchain')
    expect(wrapper.text()).toContain('vue')
  })

  it('handles like button click when connected', async () => {
    const wrapper = mount(PostCard, {
      props: {
        post: mockPost
      },
      global: {
        plugins: [createTestingPinia({
          createSpy: vi.fn,
          initialState: {
            blockchain: {
              currentAddress: 'blogchain1user...',
              isConnected: true
            }
          }
        })],
        stubs: {
          RouterLink: {
            template: '<a><slot /></a>'
          }
        }
      }
    })

    const store = useBlockchainStore()
    const likeButton = wrapper.find('button[aria-label="Like post"]')
    
    await likeButton.trigger('click')
    
    expect(store.likePost).toHaveBeenCalledWith(mockPost.id)
  })

  it('shows connect wallet message when not connected', async () => {
    const wrapper = mount(PostCard, {
      props: {
        post: mockPost
      },
      global: {
        plugins: [createTestingPinia({
          createSpy: vi.fn,
          initialState: {
            blockchain: {
              currentAddress: '',
              isConnected: false
            }
          }
        })],
        stubs: {
          RouterLink: {
            template: '<a><slot /></a>'
          }
        }
      }
    })

    const likeButton = wrapper.find('button[aria-label="Like post"]')
    await likeButton.trigger('click')
    
    // Should show alert or error message
    const store = useBlockchainStore()
    expect(store.likePost).not.toHaveBeenCalled()
  })

  it('formats date correctly', () => {
    const wrapper = mount(PostCard, {
      props: {
        post: mockPost
      },
      global: {
        plugins: [createTestingPinia({
          createSpy: vi.fn
        })],
        stubs: {
          RouterLink: {
            template: '<a><slot /></a>'
          }
        }
      }
    })

    // Should format Unix timestamp to readable date
    const dateText = wrapper.find('.text-gray-500').text()
    expect(dateText).toMatch(/\d{1,2}\/\d{1,2}\/\d{4}/)
  })

  it('truncates body text appropriately', () => {
    const longPost = {
      ...mockPost,
      body: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. '.repeat(10)
    }

    const wrapper = mount(PostCard, {
      props: {
        post: longPost
      },
      global: {
        plugins: [createTestingPinia({
          createSpy: vi.fn
        })],
        stubs: {
          RouterLink: {
            template: '<a><slot /></a>'
          }
        }
      }
    })

    const bodyText = wrapper.find('.text-gray-600').text()
    expect(bodyText.length).toBeLessThanOrEqual(153) // 150 chars + '...'
  })
})