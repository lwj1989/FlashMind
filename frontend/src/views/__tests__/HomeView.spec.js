import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory, Router } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import { getDecks } from '@/api/deck'

// 模拟API
vi.mock('@/api/deck', () => ({
  getDecks: vi.fn()
}))

describe('HomeView.vue', () => {
  let router: Router

  beforeEach(() => {
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', component: HomeView },
        { path: '/decks', component: { template: '<div></div>' } },
        { path: '/cards', component: { template: '<div></div>' } },
        { path: '/tags', component: { template: '<div></div>' } },
        { path: '/import-export', component: { template: '<div></div>' } }
      ]
    })

    vi.clearAllMocks()
  })

  it('应该正确加载卡包列表', async () => {
    const mockDecks = [
      { id: 1, name: '卡包1', archived: false, stats: { total_cards: 10, due_cards: 2 } },
      { id: 2, name: '卡包2', archived: false, stats: { total_cards: 5, due_cards: 1 } }
    ]

    getDecks.mockResolvedValue({ data: { decks: mockDecks } })

    const wrapper = mount(HomeView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    expect(getDecks).toHaveBeenCalled()
    expect(wrapper.vm.decks).toEqual(mockDecks)
    expect(wrapper.text()).toContain('卡包1')
    expect(wrapper.text()).toContain('卡包2')
  })

  it('应该正确显示卡包统计信息', async () => {
    const mockDecks = [
      { id: 1, name: '卡包1', archived: false, stats: { total_cards: 10, due_cards: 2 } }
    ]

    getDecks.mockResolvedValue({ data: { decks: mockDecks } })

    const wrapper = mount(HomeView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('10张卡片')
    expect(wrapper.text()).toContain('2张待复习')
  })

  it('应该正确导航到卡包详情', async () => {
    const mockDecks = [
      { id: 1, name: '卡包1', archived: false, stats: { total_cards: 10, due_cards: 2 } }
    ]

    getDecks.mockResolvedValue({ data: { decks: mockDecks } })

    const wrapper = mount(HomeView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    // 点击卡包
    await wrapper.find('.deck-item').trigger('click')

    expect(wrapper.vm.$router.currentRoute.value.path).toBe('/decks/1')
  })

  it('应该正确导航到卡包管理页面', async () => {
    const mockDecks = []

    getDecks.mockResolvedValue({ data: { decks: mockDecks } })

    const wrapper = mount(HomeView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    // 点击管理卡包按钮
    await wrapper.find('button.manage-decks-btn').trigger('click')

    expect(wrapper.vm.$router.currentRoute.value.path).toBe('/decks')
  })

  it('应该正确导航到卡片管理页面', async () => {
    const mockDecks = []

    getDecks.mockResolvedValue({ data: { decks: mockDecks } })

    const wrapper = mount(HomeView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    // 点击管理卡片按钮
    await wrapper.find('button.manage-cards-btn').trigger('click')

    expect(wrapper.vm.$router.currentRoute.value.path).toBe('/cards')
  })

  it('应该正确导航到标签管理页面', async () => {
    const mockDecks = []

    getDecks.mockResolvedValue({ data: { decks: mockDecks } })

    const wrapper = mount(HomeView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    // 点击管理标签按钮
    await wrapper.find('button.manage-tags-btn').trigger('click')

    expect(wrapper.vm.$router.currentRoute.value.path).toBe('/tags')
  })

  it('应该正确导航到导入导出页面', async () => {
    const mockDecks = []

    getDecks.mockResolvedValue({ data: { decks: mockDecks } })

    const wrapper = mount(HomeView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    // 点击导入导出按钮
    await wrapper.find('button.import-export-btn').trigger('click')

    expect(wrapper.vm.$router.currentRoute.value.path).toBe('/import-export')
  })

  it('应该正确处理API错误', async () => {
    getDecks.mockRejectedValue(new Error('API错误'))

    const wrapper = mount(HomeView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    expect(wrapper.vm.error).toBe('加载卡包列表失败，请重试')
    expect(wrapper.find('.error-message').exists()).toBe(true)
  })

  it('应该正确显示空状态', async () => {
    const mockDecks = []

    getDecks.mockResolvedValue({ data: { decks: mockDecks } })

    const wrapper = mount(HomeView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('还没有卡包')
    expect(wrapper.text()).toContain('创建你的第一个卡包开始学习吧')
  })
})