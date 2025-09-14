import { describe, it, expect, vi } from 'vitest'
import { searchCards, createCard, getCard, updateCard, deleteCard, getCardsByDeck, getCardsByTag } from '@/api/card'
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

describe('Card API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('searchCards', () => {
    it('应该正确搜索卡片', async () => {
      const searchParams = { keyword: 'Go', deck_id: 1, tag_id: 1 }
      const mockResponse = {
        data: {
          cards: [
            { id: 1, question: 'Go问题', answer: 'Go答案', deck_id: 1, tag_id: 1 },
            { id: 2, question: 'Golang问题', answer: 'Golang答案', deck_id: 1, tag_id: 1 }
          ]
        }
      }

      request.get.mockResolvedValue(mockResponse)

      const result = await searchCards(searchParams)

      expect(request.get).toHaveBeenCalledWith('/cards', { params: searchParams })
      expect(result).toEqual(mockResponse)
    })

    it('应该正确处理不带参数的搜索', async () => {
      const mockResponse = {
        data: {
          cards: [
            { id: 1, question: '问题1', answer: '答案1', deck_id: 1, tag_id: 1 },
            { id: 2, question: '问题2', answer: '答案2', deck_id: 1, tag_id: 1 }
          ]
        }
      }

      request.get.mockResolvedValue(mockResponse)

      const result = await searchCards()

      expect(request.get).toHaveBeenCalledWith('/cards', { params: {} })
      expect(result).toEqual(mockResponse)
    })

    it('应该正确处理API错误', async () => {
      const searchParams = { keyword: 'Go' }
      const mockError = new Error('API错误')
      request.get.mockRejectedValue(mockError)

      await expect(searchCards(searchParams)).rejects.toThrow('API错误')
    })
  })

  describe('createCard', () => {
    it('应该正确创建卡片', async () => {
      const mockCard = { question: '新问题', answer: '新答案', deck_id: 1, tag_id: 1 }
      const mockResponse = {
        data: {
          card: { id: 1, question: '新问题', answer: '新答案', deck_id: 1, tag_id: 1 }
        }
      }

      request.post.mockResolvedValue(mockResponse)

      const result = await createCard(mockCard)

      expect(request.post).toHaveBeenCalledWith('/cards', mockCard)
      expect(result).toEqual(mockResponse)
    })

    it('应该正确处理API错误', async () => {
      const mockCard = { question: '新问题', answer: '新答案', deck_id: 1, tag_id: 1 }
      const mockError = new Error('API错误')
      request.post.mockRejectedValue(mockError)

      await expect(createCard(mockCard)).rejects.toThrow('API错误')
    })
  })

  describe('getCard', () => {
    it('应该正确获取卡片详情', async () => {
      const cardId = 1
      const mockResponse = {
        data: {
          card: { id: 1, question: '问题', answer: '答案', deck_id: 1, tag_id: 1 }
        }
      }

      request.get.mockResolvedValue(mockResponse)

      const result = await getCard(cardId)

      expect(request.get).toHaveBeenCalledWith(`/cards/${cardId}`)
      expect(result).toEqual(mockResponse)
    })

    it('应该正确处理API错误', async () => {
      const cardId = 1
      const mockError = new Error('API错误')
      request.get.mockRejectedValue(mockError)

      await expect(getCard(cardId)).rejects.toThrow('API错误')
    })
  })

  describe('updateCard', () => {
    it('应该正确更新卡片', async () => {
      const cardId = 1
      const mockCard = { question: '更新后的问题', answer: '更新后的答案' }
      const mockResponse = {
        data: {
          card: { id: 1, question: '更新后的问题', answer: '更新后的答案', deck_id: 1, tag_id: 1 }
        }
      }

      request.patch.mockResolvedValue(mockResponse)

      const result = await updateCard(cardId, mockCard)

      expect(request.patch).toHaveBeenCalledWith(`/cards/${cardId}`, mockCard)
      expect(result).toEqual(mockResponse)
    })

    it('应该正确处理API错误', async () => {
      const cardId = 1
      const mockCard = { question: '更新后的问题', answer: '更新后的答案' }
      const mockError = new Error('API错误')
      request.patch.mockRejectedValue(mockError)

      await expect(updateCard(cardId, mockCard)).rejects.toThrow('API错误')
    })
  })

  describe('deleteCard', () => {
    it('应该正确删除卡片', async () => {
      const cardId = 1
      const mockResponse = {
        data: { success: true }
      }

      request.delete.mockResolvedValue(mockResponse)

      const result = await deleteCard(cardId)

      expect(request.delete).toHaveBeenCalledWith(`/cards/${cardId}`)
      expect(result).toEqual(mockResponse)
    })

    it('应该正确处理API错误', async () => {
      const cardId = 1
      const mockError = new Error('API错误')
      request.delete.mockRejectedValue(mockError)

      await expect(deleteCard(cardId)).rejects.toThrow('API错误')
    })
  })

  describe('getCardsByDeck', () => {
    it('应该正确获取指定卡包的卡片列表', async () => {
      const deckId = 1
      const mockResponse = {
        data: {
          cards: [
            { id: 1, question: '问题1', answer: '答案1', deck_id: 1, tag_id: 1 },
            { id: 2, question: '问题2', answer: '答案2', deck_id: 1, tag_id: 2 }
          ]
        }
      }

      request.get.mockResolvedValue(mockResponse)

      const result = await getCardsByDeck(deckId)

      expect(request.get).toHaveBeenCalledWith(`/cards/deck/${deckId}`)
      expect(result).toEqual(mockResponse)
    })

    it('应该正确处理API错误', async () => {
      const deckId = 1
      const mockError = new Error('API错误')
      request.get.mockRejectedValue(mockError)

      await expect(getCardsByDeck(deckId)).rejects.toThrow('API错误')
    })
  })

  describe('getCardsByTag', () => {
    it('应该正确获取指定标签的卡片列表', async () => {
      const tagId = 1
      const mockResponse = {
        data: {
          cards: [
            { id: 1, question: '问题1', answer: '答案1', deck_id: 1, tag_id: 1 },
            { id: 2, question: '问题2', answer: '答案2', deck_id: 1, tag_id: 1 }
          ]
        }
      }

      request.get.mockResolvedValue(mockResponse)

      const result = await getCardsByTag(tagId)

      expect(request.get).toHaveBeenCalledWith(`/cards/tag/${tagId}`)
      expect(result).toEqual(mockResponse)
    })

    it('应该正确处理API错误', async () => {
      const tagId = 1
      const mockError = new Error('API错误')
      request.get.mockRejectedValue(mockError)

      await expect(getCardsByTag(tagId)).rejects.toThrow('API错误')
    })
  })
})