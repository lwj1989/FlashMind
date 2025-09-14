<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-4xl mx-auto px-4 py-8">
      <div v-if="loading" class="flex justify-center items-center h-64">
        <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-500"></div>
      </div>
      
      <div v-else-if="card" class="bg-white rounded-lg shadow-lg overflow-hidden">
        <!-- 卡片头部 -->
        <div class="bg-indigo-600 text-white p-6">
          <div class="flex justify-between items-start">
            <div>
              <h1 class="text-2xl font-bold mb-2">卡片详情</h1>
              <p class="text-indigo-200">卡包: {{ deck ? deck.name : '加载中...' }}</p>
            </div>
            <div class="flex space-x-2">
              <button @click="editCard" class="bg-white text-indigo-600 hover:bg-indigo-50 px-4 py-2 rounded-md text-sm font-medium">
                编辑
              </button>
              <button @click="confirmDelete" class="bg-red-500 text-white hover:bg-red-600 px-4 py-2 rounded-md text-sm font-medium">
                删除
              </button>
              <button @click="goBack" class="bg-indigo-800 text-white hover:bg-indigo-900 px-4 py-2 rounded-md text-sm font-medium">
                返回
              </button>
            </div>
          </div>
        </div>
        
        <!-- 卡片内容 -->
        <div class="p-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
            <!-- 问题 -->
            <div class="bg-gray-50 rounded-lg p-6">
              <h2 class="text-lg font-medium text-gray-900 mb-4">问题</h2>
              <div class="prose max-w-none">
                <p class="whitespace-pre-wrap">{{ card.question }}</p>
              </div>
            </div>
            
            <!-- 答案 -->
            <div class="bg-gray-50 rounded-lg p-6">
              <h2 class="text-lg font-medium text-gray-900 mb-4">答案</h2>
              <div class="prose max-w-none">
                <p class="whitespace-pre-wrap">{{ card.answer }}</p>
              </div>
            </div>
          </div>
          
          <!-- 标签 -->
          <div class="mt-8">
            <h2 class="text-lg font-medium text-gray-900 mb-4">标签</h2>
            <div v-if="card.tag && card.tag.name" class="flex flex-wrap gap-2">
              <span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800">
                {{ card.tag.name }}
              </span>
            </div>
            <div v-else class="text-gray-500 italic">
              此卡片没有标签
            </div>
          </div>
          
          <!-- 元数据 -->
          <div class="mt-8 pt-6 border-t border-gray-200">
            <h2 class="text-lg font-medium text-gray-900 mb-4">元数据</h2>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
              <div>
                <p class="text-gray-500">创建时间</p>
                <p class="text-gray-900">{{ formatDate(card.created_at) }}</p>
              </div>
              <div>
                <p class="text-gray-500">更新时间</p>
                <p class="text-gray-900">{{ formatDate(card.updated_at) }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <div v-else class="bg-white rounded-lg shadow-lg p-8 text-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto text-gray-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <h3 class="text-lg font-medium text-gray-900 mb-2">卡片不存在</h3>
        <p class="text-gray-500 mb-4">找不到指定的卡片，可能已被删除或ID不正确。</p>
        <button @click="goBack" class="bg-indigo-600 text-white hover:bg-indigo-700 px-4 py-2 rounded-md text-sm font-medium">
          返回卡包
        </button>
      </div>
    </div>
    
    <!-- 编辑卡片模态框 -->
    <div v-if="showEditModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md">
        <div class="p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">编辑卡片</h3>
          <form @submit.prevent="updateCard">
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">问题</label>
              <textarea v-model="editingCard.question" class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500" rows="3" required></textarea>
            </div>
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">答案</label>
              <textarea v-model="editingCard.answer" class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500" rows="3" required></textarea>
            </div>
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">标签</label>
              <select v-model="editingCard.tag_id" class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500">
                <option value="">无标签</option>
                <option v-for="tag in tags" :key="tag.id" :value="tag.id">{{ tag.name }}</option>
              </select>
            </div>
            <div class="flex justify-end space-x-3">
              <button type="button" @click="showEditModal = false" class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50">
                取消
              </button>
              <button type="submit" class="px-4 py-2 bg-indigo-600 text-white rounded-md text-sm font-medium hover:bg-indigo-700">
                保存
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
    
    <!-- 删除确认模态框 -->
    <div v-if="showDeleteModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md">
        <div class="p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">确认删除</h3>
          <p class="text-gray-700 mb-6">确定要删除这张卡片吗？此操作无法撤销。</p>
          <div class="flex justify-end space-x-3">
            <button @click="showDeleteModal = false" class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50">
              取消
            </button>
            <button @click="deleteCard" class="px-4 py-2 bg-red-600 text-white rounded-md text-sm font-medium hover:bg-red-700">
              删除
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { getCard, updateCard, deleteCard } from '@/api/card'
import { getDeck } from '@/api/deck'
import { getTagsByDeck } from '@/api/tag'
import { ElMessage } from 'element-plus'

export default {
  name: 'CardDetailView',
  data() {
    return {
      card: null,
      deck: null,
      tags: [],
      loading: true,
      showEditModal: false,
      showDeleteModal: false,
      editingCard: {
        id: null,
        question: '',
        answer: '',
        deck_id: null,
        tag_id: null
      }
    }
  },
  created() {
    this.fetchCard()
  },
  mounted() {
    // 添加键盘事件监听器
    document.addEventListener('keydown', this.handleEscapeKey)
  },
  beforeUnmount() {
    // 移除键盘事件监听器
    document.removeEventListener('keydown', this.handleEscapeKey)
  },
  methods: {
    async fetchCard() {
      try {
        const cardId = this.$route.params.id
        const response = await getCard(cardId)
        console.log('获取卡片响应:', response)
        
        // 处理统一响应格式
        if (response.data && response.data.code === 'SUCCESS') {
          this.card = response.data.data
        } else {
          this.card = response.data
        }
        
        // 获取卡包信息
        if (this.card && this.card.deck_id) {
          const deckResponse = await getDeck(this.card.deck_id)
          if (deckResponse.data && deckResponse.data.code === 'SUCCESS') {
            this.deck = deckResponse.data.data
          } else {
            this.deck = deckResponse.data
          }
        }
        
        // 获取标签列表
        if (this.card && this.card.deck_id) {
          const tagsResponse = await getTagsByDeck(this.card.deck_id)
          if (tagsResponse.data && tagsResponse.data.code === 'SUCCESS') {
            this.tags = tagsResponse.data.data
          } else {
            this.tags = tagsResponse.data || []
          }
        }
        
        this.loading = false
      } catch (error) {
        console.error('获取卡片详情失败:', error)
        this.loading = false
      }
    },
    editCard() {
      if (!this.card) {
        ElMessage.error('卡片信息不存在，请刷新页面重试')
        return
      }
      this.editingCard = { 
        ...this.card
      }
      this.showEditModal = true
    },
    async updateCard() {
      try {
        if (!this.editingCard || !this.editingCard.id) {
          ElMessage.error('卡片信息不完整，请刷新页面重试')
          return
        }
        await updateCard(this.editingCard.id, this.editingCard)
        this.showEditModal = false
        this.fetchCard()
        ElMessage.success('卡片更新成功！')
      } catch (error) {
        console.error('更新卡片失败:', error)
        ElMessage.error('更新卡片失败，请稍后重试')
      }
    },
    confirmDelete() {
      this.showDeleteModal = true
    },
    async deleteCard() {
      try {
        if (!this.card) {
          ElMessage.error('卡片信息不存在，请刷新页面重试')
          return
        }
        await deleteCard(this.card.id)
        ElMessage.success('卡片删除成功！')
        this.$router.push(`/decks/${this.card.deck_id}`)
      } catch (error) {
        console.error('删除卡片失败:', error)
        ElMessage.error('删除卡片失败，请稍后重试')
      }
    },
    goBack() {
      if (this.card && this.card.deck_id) {
        this.$router.push(`/decks/${this.card.deck_id}`)
      } else {
        this.$router.push('/decks')
      }
    },
    formatDate(dateString) {
      if (!dateString) return ''
      const date = new Date(dateString)
      if (isNaN(date.getTime())) return ''
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
      })
    },
    // ESC 键关闭模态框
    handleEscapeKey(event) {
      if (event.key === 'Escape') {
        // 按优先级关闭模态框
        if (this.showEditModal) {
          this.showEditModal = false
        } else if (this.showDeleteModal) {
          this.showDeleteModal = false
        }
      }
    }
  }
}
</script>