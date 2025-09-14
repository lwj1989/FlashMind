import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory, Router } from 'vue-router'
import CardForm from '@/components/CardForm.vue'
import { getDecks, getTagsByDeck } from '@/api/deck'

// 模拟API
vi.mock('@/api/deck', () => ({
  getDecks: vi.fn(),
  getTagsByDeck: vi.fn()
}))

describe('CardForm.vue', () => {
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

  it('应该正确渲染表单', () => {
    const wrapper = mount(CardForm, {
      global: {
        plugins: [router]
      },
      props: {
        visible: true,
        card: null
      }
    })

    expect(wrapper.find('form').exists()).toBe(true)
    expect(wrapper.find('select[placeholder="选择卡包"]').exists()).toBe(true)
    expect(wrapper.find('select[placeholder="选择标签"]').exists()).toBe(true)
    expect(wrapper.find('input[placeholder="输入问题"]').exists()).toBe(true)
    expect(wrapper.find('textarea[placeholder="输入答案"]').exists()).toBe(true)
  })

  it('应该加载卡包列表', async () => {
    const mockDecks = [
      { id: 1, name: '卡包1' },
      { id: 2, name: '卡包2' }
    ]
    getDecks.mockResolvedValue({ data: { decks: mockDecks } })

    const wrapper = mount(CardForm, {
      global: {
        plugins: [router]
      },
      props: {
        visible: true,
        card: null
      }
    })

    await wrapper.vm.$nextTick()

    expect(getDecks).toHaveBeenCalled()
    expect(wrapper.vm.decks).toEqual(mockDecks)
  })

  it('应该在卡包变更时加载标签', async () => {
    const mockTags = [
      { id: 1, name: '标签1' },
      { id: 2, name: '标签2' }
    ]
    getTagsByDeck.mockResolvedValue({ data: { tags: mockTags } })

    const wrapper = mount(CardForm, {
      global: {
        plugins: [router]
      },
      props: {
        visible: true,
        card: null
      }
    })

    // 设置卡包ID
    await wrapper.setData({ card: { deck_id: 1 } })

    await wrapper.vm.$nextTick()

    expect(getTagsByDeck).toHaveBeenCalledWith(1)
    expect(wrapper.vm.tags).toEqual(mockTags)
  })

  it('应该在表单提交时触发submit事件', async () => {
    const wrapper = mount(CardForm, {
      global: {
        plugins: [router]
      },
      props: {
        visible: true,
        card: null
      }
    })

    // 设置表单数据
    await wrapper.setData({
      card: {
        deck_id: 1,
        tag_id: 1,
        question: '测试问题',
        answer: '测试答案'
      }
    })

    // 提交表单
    await wrapper.find('form').trigger('submit.prevent')

    expect(wrapper.emitted('submit')).toBeTruthy()
    expect(wrapper.emitted('submit')[0][0]).toEqual({
      deck_id: 1,
      tag_id: 1,
      question: '测试问题',
      answer: '测试答案'
    })
  })

  it('应该在点击取消按钮时触发cancel事件', async () => {
    const wrapper = mount(CardForm, {
      global: {
        plugins: [router]
      },
      props: {
        visible: true,
        card: null
      }
    })

    await wrapper.find('button[type="button"]').trigger('click')

    expect(wrapper.emitted('cancel')).toBeTruthy()
  })

  it('应该正确编辑现有卡片', async () => {
    const existingCard = {
      id: 1,
      deck_id: 1,
      tag_id: 1,
      question: '现有问题',
      answer: '现有答案'
    }

    const wrapper = mount(CardForm, {
      global: {
        plugins: [router]
      },
      props: {
        visible: true,
        card: existingCard
      }
    })

    expect(wrapper.vm.card).toEqual(existingCard)
    expect(wrapper.vm.isEdit).toBe(true)
  })
})