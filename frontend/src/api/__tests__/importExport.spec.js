import { describe, it, expect, vi } from 'vitest'
import { importDeck, exportDeck, exportDeckCSV, importDeckCSV } from '@/api/importExport'
import request from '@/api/request'

// 模拟request模块
vi.mock('@/api/request', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn()
  }
}))

describe('ImportExport API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('importDeck', () => {
    it('应该正确导入JSON格式的卡包', async () => {
      const formData = new FormData()
      formData.append('file', new File(['{"name": "测试卡包", "cards": []}'], 'test.json', { type: 'application/json' }))
      
      const mockResponse = {
        data: {
          success: true,
          message: '导入成功'
        }
      }

      request.post.mockResolvedValue(mockResponse)

      const result = await importDeck(formData)

      expect(request.post).toHaveBeenCalledWith('/import/deck', formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
      expect(result).toEqual(mockResponse)
    })

    it('应该正确处理API错误', async () => {
      const formData = new FormData()
      formData.append('file', new File(['{"name": "测试卡包", "cards": []}'], 'test.json', { type: 'application/json' }))
      
      const mockError = new Error('API错误')
      request.post.mockRejectedValue(mockError)

      await expect(importDeck(formData)).rejects.toThrow('API错误')
    })
  })

  describe('exportDeck', () => {
    it('应该正确导出JSON格式的卡包', async () => {
      const deckId = 1
      const mockResponse = {
        data: {
          name: '测试卡包',
          cards: [
            { question: '问题1', answer: '答案1' },
            { question: '问题2', answer: '答案2' }
          ]
        }
      }

      request.get.mockResolvedValue(mockResponse)

      const result = await exportDeck(deckId)

      expect(request.get).toHaveBeenCalledWith(`/export/deck/${deckId}`)
      expect(result).toEqual(mockResponse)
    })

    it('应该正确处理API错误', async () => {
      const deckId = 1
      const mockError = new Error('API错误')
      request.get.mockRejectedValue(mockError)

      await expect(exportDeck(deckId)).rejects.toThrow('API错误')
    })
  })

  describe('exportDeckCSV', () => {
    it('应该正确导出CSV格式的卡包', async () => {
      const deckId = 1
      const mockResponse = {
        data: 'question,answer\n问题1,答案1\n问题2,答案2'
      }

      request.get.mockResolvedValue(mockResponse)

      const result = await exportDeckCSV(deckId)

      expect(request.get).toHaveBeenCalledWith(`/export/deck/${deckId}/csv`)
      expect(result).toEqual(mockResponse)
    })

    it('应该正确处理API错误', async () => {
      const deckId = 1
      const mockError = new Error('API错误')
      request.get.mockRejectedValue(mockError)

      await expect(exportDeckCSV(deckId)).rejects.toThrow('API错误')
    })
  })

  describe('importDeckCSV', () => {
    it('应该正确导入CSV格式的卡包', async () => {
      const formData = new FormData()
      formData.append('file', new File(['question,answer\n问题1,答案1'], 'test.csv', { type: 'text/csv' }))
      
      const mockResponse = {
        data: {
          success: true,
          message: '导入成功'
        }
      }

      request.post.mockResolvedValue(mockResponse)

      const result = await importDeckCSV(formData)

      expect(request.post).toHaveBeenCalledWith('/import/deck/csv', formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
      expect(result).toEqual(mockResponse)
    })

    it('应该正确处理API错误', async () => {
      const formData = new FormData()
      formData.append('file', new File(['question,answer\n问题1,答案1'], 'test.csv', { type: 'text/csv' }))
      
      const mockError = new Error('API错误')
      request.post.mockRejectedValue(mockError)

      await expect(importDeckCSV(formData)).rejects.toThrow('API错误')
    })
  })
})