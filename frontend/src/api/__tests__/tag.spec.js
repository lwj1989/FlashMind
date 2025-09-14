import { describe, it, expect, vi } from 'vitest'
import { getTags, createTag, updateTag, deleteTag, getTagsByDeck } from '@/api/tag'
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

describe('Tag API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getTags', () => {
    it('应该正确获取标签列表', async () => {
      const mockResponse = {
        data: {
          tags: [
            { id: 1, name: '标签1', deck_id: 1 },
            { id: 2, name: '标签2', deck_id: 1 }
          ]
        }
      }

      request.get.mockResolvedValue(mockResponse)

      const result = await getTags()

      expect(request.get).toHaveBeenCalledWith('/tags?include_stats=true')
      expect(result).toEqual(mockResponse)
    })

    it('应该正确处理API错误', async () => {
      const mockError = new Error('API错误')
      request.get.mockRejectedValue(mockError)

      await expect(getTags()).rejects.toThrow('API错误')
    })
  })

  describe('createTag', () => {
    it('应该正确创建标签', async () => {
      const mockTag = { name: '新标签', deck_id: 1 }
      const mockResponse = {
        data: {
          tag: { id: 1, name: '新标签', deck_id: 1 }
        }
      }

      request.post.mockResolvedValue(mockResponse)

      const result = await createTag(mockTag)

      expect(request.post).toHaveBeenCalledWith('/tags', mockTag)
      expect(result).toEqual(mockResponse)
    })

    it('应该正确处理API错误', async () => {
      const mockTag = { name: '新标签', deck_id: 1 }
      const mockError = new Error('API错误')
      request.post.mockRejectedValue(mockError)

      await expect(createTag(mockTag)).rejects.toThrow('API错误')
    })
  })

  describe('updateTag', () => {
    it('应该正确更新标签', async () => {
      const tagId = 1
      const mockTag = { name: '更新后的标签' }
      const mockResponse = {
        data: {
          tag: { id: 1, name: '更新后的标签', deck_id: 1 }
        }
      }

      request.patch.mockResolvedValue(mockResponse)

      const result = await updateTag(tagId, mockTag)

      expect(request.patch).toHaveBeenCalledWith(`/tags/${tagId}`, mockTag)
      expect(result).toEqual(mockResponse)
    })

    it('应该正确处理API错误', async () => {
      const tagId = 1
      const mockTag = { name: '更新后的标签' }
      const mockError = new Error('API错误')
      request.patch.mockRejectedValue(mockError)

      await expect(updateTag(tagId, mockTag)).rejects.toThrow('API错误')
    })
  })

  describe('deleteTag', () => {
    it('应该正确删除标签', async () => {
      const tagId = 1
      const mockResponse = {
        data: { success: true }
      }

      request.delete.mockResolvedValue(mockResponse)

      const result = await deleteTag(tagId)

      expect(request.delete).toHaveBeenCalledWith(`/tags/${tagId}`)
      expect(result).toEqual(mockResponse)
    })

    it('应该正确处理API错误', async () => {
      const tagId = 1
      const mockError = new Error('API错误')
      request.delete.mockRejectedValue(mockError)

      await expect(deleteTag(tagId)).rejects.toThrow('API错误')
    })
  })

  describe('getTagsByDeck', () => {
    it('应该正确获取指定卡包的标签列表', async () => {
      const deckId = 1
      const mockResponse = {
        data: {
          tags: [
            { id: 1, name: '标签1', deck_id: 1 },
            { id: 2, name: '标签2', deck_id: 1 }
          ]
        }
      }

      request.get.mockResolvedValue(mockResponse)

      const result = await getTagsByDeck(deckId)

      expect(request.get).toHaveBeenCalledWith(`/tags/deck/${deckId}`)
      expect(result).toEqual(mockResponse)
    })

    it('应该正确处理API错误', async () => {
      const deckId = 1
      const mockError = new Error('API错误')
      request.get.mockRejectedValue(mockError)

      await expect(getTagsByDeck(deckId)).rejects.toThrow('API错误')
    })
  })
})