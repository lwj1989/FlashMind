import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory, Router } from 'vue-router'
import CardItem from '@/components/CardItem.vue'
import { deleteCard } from '@/api/card'

// 模拟API
vi.mock('@/api/card', () => ({
  deleteCard: vi.fn()
}))

// 模拟路由
const mockRouter = {
  push: vi.fn()
}

describe('CardItem.vue', () => {
  let router: Router

  beforeEach(() => {
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', component: { template: '<div></div>' } },
        { path: '/cards/:id', component: { template: '<div></div>' } }
      ]
    })

    vi.clearAllMocks()
  })

  it('应该正确渲染卡片信息', () => {
    const card = {
      id: 1,
      question: '测试问题',
      answer: '测试答案',
      deck: { id: 1, name: '测试卡包' },
      tag: { id: 1, name: '测试标签' }
    }

    const wrapper = mount(CardItem, {
      global: {
        plugins: [router],
        mocks: {
          $router: mockRouter
        }
      },
      props: {
        card
      }
    })

    expect(wrapper.text()).toContain('测试问题')
    expect(wrapper.text()).toContain('测试卡包')
    expect(wrapper.text()).toContain('测试标签')
  })

  it('应该在点击卡片时导航到详情页', async () => {
    const card = {
      id: 1,
      question: '测试问题',
      answer: '测试答案',
      deck: { id: 1, name: '测试卡包' },
      tag: { id: 1, name: '测试标签' }
    }

    const wrapper = mount(CardItem, {
      global: {
        plugins: [router],
        mocks: {
          $router: mockRouter
        }
      },
      props: {
        card
      }
    })

    await wrapper.find('.card').trigger('click')

    expect(mockRouter.push).toHaveBeenCalledWith(`/cards/${card.id}`)
  })

  it('应该在点击删除按钮时触发删除操作', async () => {
    const card = {
      id: 1,
      question: '测试问题',
      answer: '测试答案',
      deck: { id: 1, name: '测试卡包' },
      tag: { id: 1, name: '测试标签' }
    }

    deleteCard.mockResolvedValue({ data: { success: true } })

    const wrapper = mount(CardItem, {
      global: {
        plugins: [router],
        mocks: {
          $router: mockRouter
        }
      },
      props: {
        card
      }
    })

    // 点击删除按钮
    await wrapper.find('button.delete-btn').trigger('click')

    expect(deleteCard).toHaveBeenCalledWith(card.id)
    expect(wrapper.emitted('deleted')).toBeTruthy()
    expect(wrapper.emitted('deleted')[0][0]).toBe(card.id)
  })

  it('应该显示卡片的复习状态', () => {
    const card = {
      id: 1,
      question: '测试问题',
      answer: '测试答案',
      deck: { id: 1, name: '测试卡包' },
      tag: { id: 1, name: '测试标签' },
      review: {
        next_review: new Date(Date.now() + 86400000).toISOString(), // 明天
        interval: 1,
        e_factor: 2.5
      }
    }

    const wrapper = mount(CardItem, {
      global: {
        plugins: [router],
        mocks: {
          $router: mockRouter
        }
      },
      props: {
        card
      }
    })

    expect(wrapper.text()).toContain('复习间隔: 1天')
    expect(wrapper.text()).toContain('难度系数: 2.5')
  })

  it('应该正确处理没有复习信息的卡片', () => {
    const card = {
      id: 1,
      question: '测试问题',
      answer: '测试答案',
      deck: { id: 1, name: '测试卡包' },
      tag: { id: 1, name: '测试标签' },
      review: null
    }

    const wrapper = mount(CardItem, {
      global: {
        plugins: [router],
        mocks: {
          $router: mockRouter
        }
      },
      props: {
        card
      }
    })

    expect(wrapper.text()).toContain('新卡片')
  })
})