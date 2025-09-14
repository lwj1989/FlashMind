import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory, Router } from 'vue-router'
import DeckDetailView from '@/views/DeckDetailView.vue'
import { getDeck, getCardsByDeck } from '@/api/deck'

// 模拟API
vi.mock('@/api/deck', () => ({
  getDeck: vi.fn(),
  getCardsByDeck: vi.fn()
}))

describe('DeckDetailView.vue', () => {
  let router: Router

  beforeEach(() => {
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', component: { template: '<div></div>' } },
        { path: '/decks/:id', component: DeckDetailView }
      ]
    })

    vi.clearAllMocks()
  })

  it('应该正确加载卡包详情', async () => {
    const mockDeck = {
      id: 1,
      name: '测试卡包',
      archived: false,
      stats: {
        total_cards: 10,
        due_cards: 2,
        new_cards: 3
      }
    }

    const mockCards = [
      {
        id: 1,
        question: '问题1',
        answer: '答案1',
        deck_id: 1,
        tag_id: 1
      },
      {
        id: 2,
        question: '问题2',
        answer: '答案2',
        deck_id: 1,
        tag_id: 2
      }
    ]

    getDeck.mockResolvedValue({ data: { deck: mockDeck } })
    getCardsByDeck.mockResolvedValue({ data: { cards: mockCards } })

    const wrapper = mount(DeckDetailView, {
      global: {
        plugins: [router]
      },
      props: {
        id: '1'
      }
    })

    await wrapper.vm.$nextTick()

    expect(getDeck).toHaveBeenCalledWith(1)
    expect(getCardsByDeck).toHaveBeenCalledWith(1)
    expect(wrapper.vm.deck).toEqual(mockDeck)
    expect(wrapper.vm.cards).toEqual(mockCards)
  })

  it('应该正确显示卡包统计信息', async () => {
    const mockDeck = {
      id: 1,
      name: '测试卡包',
      archived: false,
      stats: {
        total_cards: 10,
        due_cards: 2,
        new_cards: 3
      }
    }

    const mockCards = []

    getDeck.mockResolvedValue({ data: { deck: mockDeck } })
    getCardsByDeck.mockResolvedValue({ data: { cards: mockCards } })

    const wrapper = mount(DeckDetailView, {
      global: {
        plugins: [router]
      },
      props: {
        id: '1'
      }
    })

    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('总卡片数: 10')
    expect(wrapper.text()).toContain('待复习: 2')
    expect(wrapper.text()).toContain('新卡片: 3')
  })

  it('应该正确处理API错误', async () => {
    getDeck.mockRejectedValue(new Error('API错误'))

    const wrapper = mount(DeckDetailView, {
      global: {
        plugins: [router]
      },
      props: {
        id: '1'
      }
    })

    await wrapper.vm.$nextTick()

    expect(wrapper.vm.error).toBe('加载卡包详情失败，请重试')
    expect(wrapper.find('.error-message').exists()).toBe(true)
  })

  it('应该正确处理卡包不存在的情况', async () => {
    getDeck.mockResolvedValue({ data: { deck: null } })

    const wrapper = mount(DeckDetailView, {
      global: {
        plugins: [router]
      },
      props: {
        id: '999'
      }
    })

    await wrapper.vm.$nextTick()

    expect(wrapper.vm.error).toBe('卡包不存在')
    expect(wrapper.find('.error-message').exists()).toBe(true)
  })
})