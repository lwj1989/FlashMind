import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory, Router } from 'vue-router'
import DecksView from '@/views/DecksView.vue'
import { getDecks, createDeck, updateDeck, deleteDeck } from '@/api/deck'

// 模拟API
vi.mock('@/api/deck', () => ({
  getDecks: vi.fn(),
  createDeck: vi.fn(),
  updateDeck: vi.fn(),
  deleteDeck: vi.fn()
}))

describe('DecksView.vue', () => {
  let router: Router

  beforeEach(() => {
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', component: { template: '<div></div>' } },
        { path: '/decks', component: DecksView }
      ]
    })

    vi.clearAllMocks()
  })

  it('应该正确加载卡包列表', async () => {
    const mockDecks = [
      { id: 1, name: '卡包1', archived: false },
      { id: 2, name: '卡包2', archived: false }
    ]

    getDecks.mockResolvedValue({ data: { decks: mockDecks } })

    const wrapper = mount(DecksView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    expect(getDecks).toHaveBeenCalled()
    expect(wrapper.vm.decks).toEqual(mockDecks)
  })

  it('应该正确创建新卡包', async () => {
    const mockDecks = [
      { id: 1, name: '卡包1', archived: false }
    ]

    const newDeck = { id: 2, name: '新卡包', archived: false }

    getDecks.mockResolvedValue({ data: { decks: mockDecks } })
    createDeck.mockResolvedValue({ data: { deck: newDeck } })

    const wrapper = mount(DecksView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    // 设置新卡包名称
    await wrapper.setData({ newDeckName: '新卡包' })

    // 提交创建表单
    await wrapper.vm.createDeck()

    expect(createDeck).toHaveBeenCalledWith({ name: '新卡包' })
    expect(wrapper.vm.newDeckName).toBe('')

    // 列表应该刷新
    expect(getDecks).toHaveBeenCalledTimes(2)
  })

  it('应该正确更新卡包', async () => {
    const mockDecks = [
      { id: 1, name: '卡包1', archived: false }
    ]

    getDecks.mockResolvedValue({ data: { decks: mockDecks } })
    updateDeck.mockResolvedValue({ data: { success: true } })

    const wrapper = mount(DecksView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    // 开始编辑
    await wrapper.setData({ editingDeck: { id: 1, name: '卡包1' } })

    // 更新名称
    await wrapper.setData({ editingDeckName: '更新后的卡包' })

    // 提交更新
    await wrapper.vm.updateDeck()

    expect(updateDeck).toHaveBeenCalledWith(1, { name: '更新后的卡包' })
    expect(wrapper.vm.editingDeck).toBe(null)

    // 列表应该刷新
    expect(getDecks).toHaveBeenCalledTimes(2)
  })

  it('应该正确删除卡包', async () => {
    const mockDecks = [
      { id: 1, name: '卡包1', archived: false }
    ]

    getDecks.mockResolvedValue({ data: { decks: mockDecks } })
    deleteDeck.mockResolvedValue({ data: { success: true } })

    const wrapper = mount(DecksView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    // 删除卡包
    await wrapper.vm.deleteDeck(1)

    expect(deleteDeck).toHaveBeenCalledWith(1)

    // 列表应该刷新
    expect(getDecks).toHaveBeenCalledTimes(2)
  })

  it('应该正确切换卡包归档状态', async () => {
    const mockDecks = [
      { id: 1, name: '卡包1', archived: false }
    ]

    const updatedDeck = { id: 1, name: '卡包1', archived: true }

    getDecks.mockResolvedValue({ data: { decks: mockDecks } })
    updateDeck.mockResolvedValue({ data: { deck: updatedDeck } })

    const wrapper = mount(DecksView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    // 切换归档状态
    await wrapper.vm.toggleArchive(1)

    expect(updateDeck).toHaveBeenCalledWith(1, { archived: true })

    // 列表应该刷新
    expect(getDecks).toHaveBeenCalledTimes(2)
  })

  it('应该正确处理API错误', async () => {
    getDecks.mockRejectedValue(new Error('API错误'))

    const wrapper = mount(DecksView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    expect(wrapper.vm.error).toBe('加载卡包失败，请重试')
    expect(wrapper.find('.error-message').exists()).toBe(true)
  })
})