import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory, Router } from 'vue-router'
import CardsView from '@/views/CardsView.vue'
import { searchCards } from '@/api/card'

// 模拟API
vi.mock('@/api/card', () => ({
  searchCards: vi.fn()
}))

describe('CardsView.vue', () => {
  let router: Router

  beforeEach(() => {
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', component: { template: '<div></div>' } },
        { path: '/cards', component: CardsView }
      ]
    })

    vi.clearAllMocks()
  })

  it('应该正确加载卡片列表', async () => {
    const mockCards = [
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

    searchCards.mockResolvedValue({ data: { cards: mockCards } })

    const wrapper = mount(CardsView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    expect(searchCards).toHaveBeenCalled()
    expect(wrapper.vm.cards).toEqual(mockCards)
  })

  it('应该正确处理搜索功能', async () => {
    const mockCards = [
      {
        id: 1,
        question: 'Go问题',
        answer: 'Go答案',
        deck: { id: 1, name: 'Go卡包' },
        tag: { id: 1, name: 'Go标签' }
      }
    ]

    searchCards.mockResolvedValue({ data: { cards: mockCards } })

    const wrapper = mount(CardsView, {
      global: {
        plugins: [router]
      }
    })

    // 设置搜索关键词
    await wrapper.setData({ searchKeyword: 'Go' })
    await wrapper.vm.searchCards()

    expect(searchCards).toHaveBeenCalledWith({ keyword: 'Go' })
    expect(wrapper.vm.cards).toEqual(mockCards)
  })

  it('应该正确显示创建卡片表单', async () => {
    const wrapper = mount(CardsView, {
      global: {
        plugins: [router]
      }
    })

    // 初始状态不显示表单
    expect(wrapper.vm.showCardForm).toBe(false)
    expect(wrapper.findComponent({ name: 'CardForm' }).exists()).toBe(false)

    // 点击创建按钮
    await wrapper.find('button.create-btn').trigger('click')

    // 应该显示表单
    expect(wrapper.vm.showCardForm).toBe(true)
    expect(wrapper.findComponent({ name: 'CardForm' }).exists()).toBe(true)
  })

  it('应该在创建卡片成功后刷新列表', async () => {
    const mockCards = [
      {
        id: 1,
        question: '问题1',
        answer: '答案1',
        deck: { id: 1, name: '卡包1' },
        tag: { id: 1, name: '标签1' }
      }
    ]

    searchCards.mockResolvedValue({ data: { cards: mockCards } })

    const wrapper = mount(CardsView, {
      global: {
        plugins: [router]
      }
    })

    // 显示表单
    await wrapper.setData({ showCardForm: true })

    // 触发表单提交事件
    await wrapper.findComponent({ name: 'CardForm' }).vm.$emit('submit', {
      deck_id: 1,
      tag_id: 1,
      question: '新问题',
      answer: '新答案'
    })

    // 表单应该关闭
    expect(wrapper.vm.showCardForm).toBe(false)

    // 列表应该刷新
    expect(searchCards).toHaveBeenCalledTimes(2)
  })

  it('应该在卡片删除后刷新列表', async () => {
    const mockCards = [
      {
        id: 1,
        question: '问题1',
        answer: '答案1',
        deck: { id: 1, name: '卡包1' },
        tag: { id: 1, name: '标签1' }
      }
    ]

    searchCards.mockResolvedValue({ data: { cards: mockCards } })

    const wrapper = mount(CardsView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    // 触发卡片删除事件
    await wrapper.findComponent({ name: 'CardList' }).vm.$emit('deleted', 1)

    // 列表应该刷新
    expect(searchCards).toHaveBeenCalledTimes(2)
  })

  it('应该正确处理API错误', async () => {
    searchCards.mockRejectedValue(new Error('API错误'))

    const wrapper = mount(CardsView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    expect(wrapper.vm.error).toBe('加载卡片失败，请重试')
    expect(wrapper.find('.error-message').exists()).toBe(true)
  })
})