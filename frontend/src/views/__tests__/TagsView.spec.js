import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory, Router } from 'vue-router'
import TagsView from '@/views/TagsView.vue'
import { getTags, createTag, updateTag, deleteTag } from '@/api/tag'

// 模拟API
vi.mock('@/api/tag', () => ({
  getTags: vi.fn(),
  createTag: vi.fn(),
  updateTag: vi.fn(),
  deleteTag: vi.fn()
}))

describe('TagsView.vue', () => {
  let router: Router

  beforeEach(() => {
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', component: { template: '<div></div>' } },
        { path: '/tags', component: TagsView }
      ]
    })

    vi.clearAllMocks()
  })

  it('应该正确加载标签列表', async () => {
    const mockTags = [
      { id: 1, name: '标签1', deck_id: 1 },
      { id: 2, name: '标签2', deck_id: 1 }
    ]

    getTags.mockResolvedValue({ data: { tags: mockTags } })

    const wrapper = mount(TagsView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    expect(getTags).toHaveBeenCalled()
    expect(wrapper.vm.tags).toEqual(mockTags)
  })

  it('应该正确创建新标签', async () => {
    const mockTags = [
      { id: 1, name: '标签1', deck_id: 1 }
    ]

    const newTag = { id: 2, name: '新标签', deck_id: 1 }

    getTags.mockResolvedValue({ data: { tags: mockTags } })
    createTag.mockResolvedValue({ data: { tag: newTag } })

    const wrapper = mount(TagsView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    // 设置新标签名称和卡包ID
    await wrapper.setData({ 
      newTagName: '新标签',
      newTagDeckId: 1
    })

    // 提交创建表单
    await wrapper.vm.createTag()

    expect(createTag).toHaveBeenCalledWith({ name: '新标签', deck_id: 1 })
    expect(wrapper.vm.newTagName).toBe('')

    // 列表应该刷新
    expect(getTags).toHaveBeenCalledTimes(2)
  })

  it('应该正确更新标签', async () => {
    const mockTags = [
      { id: 1, name: '标签1', deck_id: 1 }
    ]

    getTags.mockResolvedValue({ data: { tags: mockTags } })
    updateTag.mockResolvedValue({ data: { success: true } })

    const wrapper = mount(TagsView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    // 开始编辑
    await wrapper.setData({ editingTag: { id: 1, name: '标签1', deck_id: 1 } })

    // 更新名称
    await wrapper.setData({ editingTagName: '更新后的标签' })

    // 提交更新
    await wrapper.vm.updateTag()

    expect(updateTag).toHaveBeenCalledWith(1, { name: '更新后的标签' })
    expect(wrapper.vm.editingTag).toBe(null)

    // 列表应该刷新
    expect(getTags).toHaveBeenCalledTimes(2)
  })

  it('应该正确删除标签', async () => {
    const mockTags = [
      { id: 1, name: '标签1', deck_id: 1 }
    ]

    getTags.mockResolvedValue({ data: { tags: mockTags } })
    deleteTag.mockResolvedValue({ data: { success: true } })

    const wrapper = mount(TagsView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    // 删除标签
    await wrapper.vm.deleteTag(1)

    expect(deleteTag).toHaveBeenCalledWith(1)

    // 列表应该刷新
    expect(getTags).toHaveBeenCalledTimes(2)
  })

  it('应该正确处理API错误', async () => {
    getTags.mockRejectedValue(new Error('API错误'))

    const wrapper = mount(TagsView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    expect(wrapper.vm.error).toBe('加载标签失败，请重试')
    expect(wrapper.find('.error-message').exists()).toBe(true)
  })
})