import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory, Router } from 'vue-router'
import TagCardsView from '@/views/TagCardsView.vue'
import { getTag, getCardsByTag } from '@/api/tag'

// 模拟API
vi.mock('@/api/tag', () => ({
  getTag: vi.fn(),
  getCardsByTag: vi.fn()
}))

describe('TagCardsView.vue', () => {
  let router: Router

  beforeEach(() => {
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', component: { template: '<div></div>' } },
        { path: '/tags/:id', component: TagCardsView }
      ]
    })

    vi.clearAllMocks()
  })

  it('应该正确加载标签详情和卡片列表', async () => {
    const mockTag = {
      id: 1,
      name: '测试标签',
      deck_id: 1,
      deck: { id: 1, name: '测试卡包' },
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
        tag_id: 1
      }
    ]

    getTag.mockResolvedValue({ data: { tag: mockTag } })
    getCardsByTag.mockResolvedValue({ data: { cards: mockCards } })

    const wrapper = mount(TagCardsView, {
      global: {
        plugins: [router]
      },
      props: {
        id: '1'
      }
    })

    await wrapper.vm.$nextTick()

    expect(getTag).toHaveBeenCalledWith(1)
    expect(getCardsByTag).toHaveBeenCalledWith(1)
    expect(wrapper.vm.tag).toEqual(mockTag)
    expect(wrapper.vm.cards).toEqual(mockCards)
  })

  it('应该正确显示标签统计信息', async () => {
    const mockTag = {
      id: 1,
      name: '测试标签',
      deck_id: 1,
      deck: { id: 1, name: '测试卡包' },
      stats: {
        total_cards: 10,
        due_cards: 2,
        new_cards: 3
      }
    }

    const mockCards = []

    getTag.mockResolvedValue({ data: { tag: mockTag } })
    getCardsByTag.mockResolvedValue({ data: { cards: mockCards } })

    const wrapper = mount(TagCardsView, {
      global: {
        plugins: [router]
      },
      props: {
        id: '1'
      }
    })

    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('测试标签')
    expect(wrapper.text()).toContain('测试卡包')
    expect(wrapper.text()).toContain('总卡片数: 10')
    expect(wrapper.text()).toContain('待复习: 2')
    expect(wrapper.text()).toContain('新卡片: 3')
  })

  it('应该正确处理API错误', async () => {
    getTag.mockRejectedValue(new Error('API错误'))

    const wrapper = mount(TagCardsView, {
      global: {
        plugins: [router]
      },
      props: {
        id: '1'
      }
    })

    await wrapper.vm.$nextTick()

    expect(wrapper.vm.error).toBe('加载标签详情失败，请重试')
    expect(wrapper.find('.error-message').exists()).toBe(true)
  })

  it('应该正确处理标签不存在的情况', async () => {
    getTag.mockResolvedValue({ data: { tag: null } })

    const wrapper = mount(TagCardsView, {
      global: {
        plugins: [router]
      },
      props: {
        id: '999'
      }
    })

    await wrapper.vm.$nextTick()

    expect(wrapper.vm.error).toBe('标签不存在')
    expect(wrapper.find('.error-message').exists()).toBe(true)
  })

  it('应该在卡片删除后刷新列表', async () => {
    const mockTag = {
      id: 1,
      name: '测试标签',
      deck_id: 1,
      deck: { id: 1, name: '测试卡包' },
      stats: {
        total_cards: 2,
        due_cards: 0,
        new_cards: 0
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
        tag_id: 1
      }
    ]

    getTag.mockResolvedValue({ data: { tag: mockTag } })
    getCardsByTag.mockResolvedValue({ data: { cards: mockCards } })

    const wrapper = mount(TagCardsView, {
      global: {
        plugins: [router]
      },
      props: {
        id: '1'
      }
    })

    await wrapper.vm.$nextTick()

    // 触发卡片删除事件
    await wrapper.findComponent({ name: 'CardList' }).vm.$emit('deleted', 1)

    // 列表应该刷新
    expect(getCardsByTag).toHaveBeenCalledTimes(2)
  })
})