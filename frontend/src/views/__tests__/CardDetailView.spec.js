import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory, Router } from 'vue-router'
import CardDetailView from '@/views/CardDetailView.vue'
import { getCard, updateCard, deleteCard } from '@/api/card'

// 模拟API
vi.mock('@/api/card', () => ({
  getCard: vi.fn(),
  updateCard: vi.fn(),
  deleteCard: vi.fn()
}))

describe('CardDetailView.vue', () => {
  let router: Router

  beforeEach(() => {
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', component: { template: '<div></div>' } },
        { path: '/cards/:id', component: CardDetailView }
      ]
    })

    vi.clearAllMocks()
  })

  it('应该正确加载卡片详情', async () => {
    const mockCard = {
      id: 1,
      question: '测试问题',
      answer: '测试答案',
      deck_id: 1,
      tag_id: 1,
      deck: { id: 1, name: '测试卡包' },
      tag: { id: 1, name: '测试标签' },
      review: {
        id: 1,
        card_id: 1,
        interval: 1,
        e_factor: 2.5,
        next_review: new Date(Date.now() + 86400000).toISOString()
      }
    }

    getCard.mockResolvedValue({ data: { card: mockCard } })

    const wrapper = mount(CardDetailView, {
      global: {
        plugins: [router]
      },
      props: {
        id: '1'
      }
    })

    await wrapper.vm.$nextTick()

    expect(getCard).toHaveBeenCalledWith(1)
    expect(wrapper.vm.card).toEqual(mockCard)
    expect(wrapper.text()).toContain('测试问题')
    expect(wrapper.text()).toContain('测试答案')
    expect(wrapper.text()).toContain('测试卡包')
    expect(wrapper.text()).toContain('测试标签')
  })

  it('应该正确显示复习信息', async () => {
    const mockCard = {
      id: 1,
      question: '测试问题',
      answer: '测试答案',
      deck_id: 1,
      tag_id: 1,
      deck: { id: 1, name: '测试卡包' },
      tag: { id: 1, name: '测试标签' },
      review: {
        id: 1,
        card_id: 1,
        interval: 1,
        e_factor: 2.5,
        next_review: new Date(Date.now() + 86400000).toISOString()
      }
    }

    getCard.mockResolvedValue({ data: { card: mockCard } })

    const wrapper = mount(CardDetailView, {
      global: {
        plugins: [router]
      },
      props: {
        id: '1'
      }
    })

    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('复习间隔: 1天')
    expect(wrapper.text()).toContain('难度系数: 2.5')
  })

  it('应该正确进入编辑模式', async () => {
    const mockCard = {
      id: 1,
      question: '测试问题',
      answer: '测试答案',
      deck_id: 1,
      tag_id: 1,
      deck: { id: 1, name: '测试卡包' },
      tag: { id: 1, name: '测试标签' },
      review: {
        id: 1,
        card_id: 1,
        interval: 1,
        e_factor: 2.5,
        next_review: new Date(Date.now() + 86400000).toISOString()
      }
    }

    getCard.mockResolvedValue({ data: { card: mockCard } })

    const wrapper = mount(CardDetailView, {
      global: {
        plugins: [router]
      },
      props: {
        id: '1'
      }
    })

    await wrapper.vm.$nextTick()

    // 初始状态不是编辑模式
    expect(wrapper.vm.isEditing).toBe(false)

    // 点击编辑按钮
    await wrapper.find('button.edit-btn').trigger('click')

    // 应该进入编辑模式
    expect(wrapper.vm.isEditing).toBe(true)
    expect(wrapper.vm.editingCard).toEqual(mockCard)
  })

  it('应该正确保存编辑', async () => {
    const mockCard = {
      id: 1,
      question: '测试问题',
      answer: '测试答案',
      deck_id: 1,
      tag_id: 1,
      deck: { id: 1, name: '测试卡包' },
      tag: { id: 1, name: '测试标签' },
      review: {
        id: 1,
        card_id: 1,
        interval: 1,
        e_factor: 2.5,
        next_review: new Date(Date.now() + 86400000).toISOString()
      }
    }

    const updatedCard = {
      ...mockCard,
      question: '更新后的问题',
      answer: '更新后的答案'
    }

    getCard.mockResolvedValue({ data: { card: mockCard } })
    updateCard.mockResolvedValue({ data: { card: updatedCard } })

    const wrapper = mount(CardDetailView, {
      global: {
        plugins: [router]
      },
      props: {
        id: '1'
      }
    })

    await wrapper.vm.$nextTick()

    // 进入编辑模式
    await wrapper.find('button.edit-btn').trigger('click')

    // 更新内容
    await wrapper.setData({
      editingCard: {
        ...mockCard,
        question: '更新后的问题',
        answer: '更新后的答案'
      }
    })

    // 保存编辑
    await wrapper.vm.saveEdit()

    expect(updateCard).toHaveBeenCalledWith(1, {
      question: '更新后的问题',
      answer: '更新后的答案',
      deck_id: 1,
      tag_id: 1
    })
    expect(wrapper.vm.isEditing).toBe(false)
    expect(wrapper.vm.card).toEqual(updatedCard)
  })

  it('应该正确取消编辑', async () => {
    const mockCard = {
      id: 1,
      question: '测试问题',
      answer: '测试答案',
      deck_id: 1,
      tag_id: 1,
      deck: { id: 1, name: '测试卡包' },
      tag: { id: 1, name: '测试标签' },
      review: {
        id: 1,
        card_id: 1,
        interval: 1,
        e_factor: 2.5,
        next_review: new Date(Date.now() + 86400000).toISOString()
      }
    }

    getCard.mockResolvedValue({ data: { card: mockCard } })

    const wrapper = mount(CardDetailView, {
      global: {
        plugins: [router]
      },
      props: {
        id: '1'
      }
    })

    await wrapper.vm.$nextTick()

    // 进入编辑模式
    await wrapper.find('button.edit-btn').trigger('click')

    // 修改内容
    await wrapper.setData({
      editingCard: {
        ...mockCard,
        question: '修改后的问题',
        answer: '修改后的答案'
      }
    })

    // 取消编辑
    await wrapper.vm.cancelEdit()

    expect(wrapper.vm.isEditing).toBe(false)
    // 内容应该恢复为原始值
    expect(wrapper.vm.card).toEqual(mockCard)
  })

  it('应该正确删除卡片', async () => {
    const mockCard = {
      id: 1,
      question: '测试问题',
      answer: '测试答案',
      deck_id: 1,
      tag_id: 1,
      deck: { id: 1, name: '测试卡包' },
      tag: { id: 1, name: '测试标签' },
      review: {
        id: 1,
        card_id: 1,
        interval: 1,
        e_factor: 2.5,
        next_review: new Date(Date.now() + 86400000).toISOString()
      }
    }

    getCard.mockResolvedValue({ data: { card: mockCard } })
    deleteCard.mockResolvedValue({ data: { success: true } })

    const wrapper = mount(CardDetailView, {
      global: {
        plugins: [router],
        mocks: {
          $router: { push: vi.fn() }
        }
      },
      props: {
        id: '1'
      }
    })

    await wrapper.vm.$nextTick()

    // 删除卡片
    await wrapper.vm.deleteCard()

    expect(deleteCard).toHaveBeenCalledWith(1)
    // 应该导航回卡片列表
    expect(wrapper.vm.$router.push).toHaveBeenCalledWith('/cards')
  })

  it('应该正确处理API错误', async () => {
    getCard.mockRejectedValue(new Error('API错误'))

    const wrapper = mount(CardDetailView, {
      global: {
        plugins: [router]
      },
      props: {
        id: '1'
      }
    })

    await wrapper.vm.$nextTick()

    expect(wrapper.vm.error).toBe('加载卡片详情失败，请重试')
    expect(wrapper.find('.error-message').exists()).toBe(true)
  })

  it('应该正确处理卡片不存在的情况', async () => {
    getCard.mockResolvedValue({ data: { card: null } })

    const wrapper = mount(CardDetailView, {
      global: {
        plugins: [router]
      },
      props: {
        id: '999'
      }
    })

    await wrapper.vm.$nextTick()

    expect(wrapper.vm.error).toBe('卡片不存在')
    expect(wrapper.find('.error-message').exists()).toBe(true)
  })
})