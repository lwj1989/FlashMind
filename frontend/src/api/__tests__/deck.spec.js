import { describe, it, expect, vi } from 'vitest'
import { getDecks, createDeck, updateDeck, deleteDeck } from '@/api/deck'
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

describe('Deck API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getDecks', () => {
    it('应该正确获取卡包列表', async () => {
      const mockResponse = {
        data: {
          decks: [
            { id: 1, name: '卡包1', archived: false },
            { id: 2, name: '卡包2', archived: false }
          ]
        }
      }

      request.get.mockResolvedValue(mockResponse)

      const result = await getDecks()

      expect(request.get).toHaveBeenCalledWith('/decks?include_stats=true')
      expect(result).toEqual(mockResponse)
    })

    it('应该正确处理API错误', async () => {
      const mockError = new Error('API错误')
      request.get.mockRejectedValue(mockError)

      await expect(getDecks()).rejects.toThrow('API错误')
    })
  })

  describe('createDeck', () => {
    it('应该正确创建卡包', async () => {
      const mockDeck = { name: '新卡包' }
      const mockResponse = {
        data: {
          deck: { id: 1, name: '新卡包', archived: false }
        }
      }

      request.post.mockResolvedValue(mockResponse)

      const result = await createDeck(mockDeck)

      expect(request.post).toHaveBeenCalledWith('/decks', mockDeck)
      expect(result).toEqual(mockResponse)
    })

    it('应该正确处理API错误', async () => {
      const mockDeck = { name: '新卡包' }
      const mockError = new Error('API错误')
      request.post.mockRejectedValue(mockError)

      await expect(createDeck(mockDeck)).rejects.toThrow('API错误')
    })
  })

  describe('updateDeck', () => {
    it('应该正确更新卡包', async () => {
      const deckId = 1
      const mockDeck = { name: '更新后的卡包' }
      const mockResponse = {
        data: {
          deck: { id: 1, name: '更新后的卡包', archived: false }
        }
      }

      request.patch.mockResolvedValue(mockResponse)

      const result = await updateDeck(deckId, mockDeck)

      expect(request.patch).toHaveBeenCalledWith(`/decks/${deckId}`, mockDeck)
      expect(result).toEqual(mockResponse)
    })

    it('应该正确处理API错误', async () => {
      const deckId = 1
      const mockDeck = { name: '更新后的卡包' }
      const mockError = new Error('API错误')
      request.patch.mockRejectedValue(mockError)

      await expect(updateDeck(deckId, mockDeck)).rejects.toThrow('API错误')
    })
  })

  describe('deleteDeck', () => {
    it('应该正确删除卡包', async () => {
      const deckId = 1
      const mockResponse = {
        data: { success: true }
      }

      request.delete.mockResolvedValue(mockResponse)

      const result = await deleteDeck(deckId)

      expect(request.delete).toHaveBeenCalledWith(`/decks/${deckId}`)
      expect(result).toEqual(mockResponse)
    })

    it('应该正确处理API错误', async () => {
      const deckId = 1
      const mockError = new Error('API错误')
      request.delete.mockRejectedValue(mockError)

      await expect(deleteDeck(deckId)).rejects.toThrow('API错误')
    })
  })
})