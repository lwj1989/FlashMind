import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory, Router } from 'vue-router'
import ImportExportView from '@/views/ImportExportView.vue'
import { getDecks } from '@/api/deck'
import { importDeck, exportDeck, exportDeckCSV, importDeckCSV } from '@/api/importExport'

// 模拟API
vi.mock('@/api/deck', () => ({
  getDecks: vi.fn()
}))

vi.mock('@/api/importExport', () => ({
  importDeck: vi.fn(),
  exportDeck: vi.fn(),
  exportDeckCSV: vi.fn(),
  importDeckCSV: vi.fn()
}))

// 模拟文件读取器
global.FileReader = vi.fn().mockImplementation(() => ({
  readAsText: vi.fn(),
  onload: vi.fn(),
  onerror: vi.fn()
}))

describe('ImportExportView.vue', () => {
  let router: Router

  beforeEach(() => {
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', component: { template: '<div></div>' } },
        { path: '/import-export', component: ImportExportView }
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

    const wrapper = mount(ImportExportView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    expect(getDecks).toHaveBeenCalled()
    expect(wrapper.vm.decks).toEqual(mockDecks)
  })

  it('应该正确导入JSON文件', async () => {
    const mockDecks = [
      { id: 1, name: '卡包1', archived: false }
    ]

    const mockFile = new File(['{"name": "测试卡包", "cards": []}'], 'test.json', { type: 'application/json' })

    getDecks.mockResolvedValue({ data: { decks: mockDecks } })
    importDeck.mockResolvedValue({ data: { success: true } })

    const wrapper = mount(ImportExportView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    // 设置文件
    await wrapper.setData({ selectedFile: mockFile })

    // 导入文件
    await wrapper.vm.uploadFile()

    expect(importDeck).toHaveBeenCalled()
    expect(wrapper.vm.message).toBe('导入成功')
  })

  it('应该正确导出JSON文件', async () => {
    const mockDecks = [
      { id: 1, name: '卡包1', archived: false }
    ]

    const mockExportData = {
      name: '卡包1',
      cards: [
        { question: '问题1', answer: '答案1' }
      ]
    }

    getDecks.mockResolvedValue({ data: { decks: mockDecks } })
    exportDeck.mockResolvedValue({ data: mockExportData })

    const wrapper = mount(ImportExportView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    // 设置选中的卡包
    await wrapper.setData({ selectedDeckId: 1 })

    // 导出文件
    await wrapper.vm.exportDeck()

    expect(exportDeck).toHaveBeenCalledWith(1)
  })

  it('应该正确导入CSV文件', async () => {
    const mockDecks = [
      { id: 1, name: '卡包1', archived: false }
    ]

    const mockFile = new File(['question,answer\n问题1,答案1'], 'test.csv', { type: 'text/csv' })

    getDecks.mockResolvedValue({ data: { decks: mockDecks } })
    importDeckCSV.mockResolvedValue({ data: { success: true } })

    const wrapper = mount(ImportExportView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    // 设置文件
    await wrapper.setData({ selectedFile: mockFile })

    // 导入文件
    await wrapper.vm.uploadCSVFile()

    expect(importDeckCSV).toHaveBeenCalled()
    expect(wrapper.vm.message).toBe('导入成功')
  })

  it('应该正确导出CSV文件', async () => {
    const mockDecks = [
      { id: 1, name: '卡包1', archived: false }
    ]

    const mockCSVData = 'question,answer\n问题1,答案1'

    getDecks.mockResolvedValue({ data: { decks: mockDecks } })
    exportDeckCSV.mockResolvedValue({ data: mockCSVData })

    const wrapper = mount(ImportExportView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    // 设置选中的卡包
    await wrapper.setData({ selectedDeckId: 1 })

    // 导出文件
    await wrapper.vm.exportDeckCSV()

    expect(exportDeckCSV).toHaveBeenCalledWith(1)
  })

  it('应该正确处理文件选择', async () => {
    const mockDecks = [
      { id: 1, name: '卡包1', archived: false }
    ]

    getDecks.mockResolvedValue({ data: { decks: mockDecks } })

    const wrapper = mount(ImportExportView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    // 模拟文件选择
    const mockFile = new File(['test'], 'test.json', { type: 'application/json' })
    const fileInput = wrapper.find('input[type="file"]')
    Object.defineProperty(fileInput.element, 'files', {
      value: [mockFile]
    })

    await fileInput.trigger('change')

    expect(wrapper.vm.selectedFile).toBe(mockFile)
  })

  it('应该正确处理API错误', async () => {
    getDecks.mockRejectedValue(new Error('API错误'))

    const wrapper = mount(ImportExportView, {
      global: {
        plugins: [router]
      }
    })

    await wrapper.vm.$nextTick()

    expect(wrapper.vm.error).toBe('加载卡包失败，请重试')
    expect(wrapper.find('.error-message').exists()).toBe(true)
  })
})