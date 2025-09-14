import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory, Router } from 'vue-router'
import CardList from '@/components/CardList.vue'
import { searchCards } from '@/api/card'

// 模拟API
vi.mock('@/api/card', () => ({
  searchCards: vi.fn()
}))

describe('CardList.vue', () => {
  let router: Router

  beforeEach(() => {
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', component: { template: '<div></div>' } }
      ]
    })

    vi.clearAllMocks()
  })

  it('应该正确渲染卡片列表', () => {
    const cards = [
      {
        id: 1,
        question: '问题1',
        answer: '答案1',
        deck: { id: 1, name: '卡包1' },
        tag: { id: 1, name: '标签1' }
      },
      {
        id: 2,
        question: '问题2',
        answer: '答案2',
        deck: { id: 2, name: '卡包2' },
        tag: { id: 2, name: '标签2' }
      }
    ]

    const wrapper = mount(CardList, {
      global: {
        plugins: [router]
      },
      props: {
        cards
      }
    })

    const cardItems = wrapper.findAllComponents({ name: 'CardItem' })
    expect(cardItems.length).toBe(2)
    expect(cardItems[0].props('card')).toEqual(cards[0])
    expect(cardItems[1].props('card')).toEqual(cards[1])
  })

  it('应该显示空列表消息', () => {
    const wrapper = mount(CardList, {
      global: {
        plugins: [router]
      },
      props: {
        cards: []
      }
    })

    expect(wrapper.text()).toContain('没有找到卡片')
  })

  it('应该在卡片被删除时触发deleted事件', async () => {
    const cards = [
      {
        id: 1,
        question: '问题1',
        answer: '答案1',
        deck: { id: 1, name: '卡包1' },
        tag: { id: 1, name: '标签1' }
      }
    ]

    const wrapper = mount(CardList, {
      global: {
        plugins: [router]
      },
      props: {
        cards
      }
    })

    // 触发CardItem的deleted事件
    await wrapper.findComponent({ name: 'CardItem' }).vm.$emit('deleted', 1)

    expect(wrapper.emitted('deleted')).toBeTruthy()
    expect(wrapper.emitted('deleted')[0][0]).toBe(1)
  })

  it('应该正确显示加载状态', () => {
    const wrapper = mount(CardList, {
      global: {
        plugins: [router]
      },
      props: {
        cards: [],
        loading: true
      }
    })

    expect(wrapper.find('.loading-indicator').exists()).toBe(true)
  })

  it('应该正确显示错误状态', () => {
    const wrapper = mount(CardList, {
      global: {
        plugins: [router]
      },
      props: {
        cards: [],
        error: '加载失败'
      }
    })

    expect(wrapper.text()).toContain('加载失败')
  })
})