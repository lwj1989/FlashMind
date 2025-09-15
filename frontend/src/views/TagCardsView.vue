<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-6xl mx-auto px-4 py-12">
      <header class="mb-8">
        <div class="flex items-center mb-4">
          <button @click="goBack" class="mr-4 text-gray-600 hover:text-gray-800">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
          </button>
          <h1 class="text-3xl font-bold text-gray-800">标签卡片</h1>
        </div>
        <div class="flex items-center">
          <span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium mr-3" :style="{ backgroundColor: `${tag.color}20`, color: tag.color }">
            {{ tag.name }}
          </span>
          <span class="text-gray-600">卡包: {{ deck ? deck.name : '加载中...' }}</span>
        </div>
      </header>

      <div class="bg-white rounded-xl shadow p-6 mb-6">
        <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
          <div class="flex-1">
            <div class="relative">
              <input 
                type="text" 
                v-model="searchQuery" 
                placeholder="搜索卡片..." 
                class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
              >
              <div class="absolute left-3 top-2.5 text-gray-400">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd" />
                </svg>
              </div>
            </div>
          </div>
          <div class="flex space-x-3">
            <button @click="showCreateCardModal = true" class="bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700 transition-colors flex items-center">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clip-rule="evenodd" />
              </svg>
              添加卡片
            </button>
          </div>
        </div>
      </div>

      <!-- 卡片列表 -->
      <div v-if="loading" class="text-center py-8">
        <div class="inline-block animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-600"></div>
        <p class="mt-2 text-gray-600">加载中...</p>
      </div>

      <div v-else-if="filteredCards.length === 0" class="bg-white rounded-xl shadow p-8 text-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto text-gray-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <h3 class="text-xl font-medium text-gray-900 mb-2">暂无卡片</h3>
        <p class="text-gray-500 mb-4">为这个标签添加第一张卡片</p>
        <button @click="showCreateCardModal = true" class="bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700 transition-colors">
          添加卡片
        </button>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div 
          v-for="card in filteredCards" 
          :key="card.id" 
          class="bg-white rounded-xl shadow hover:shadow-lg transition-shadow cursor-pointer"
          @click="viewCard(card)"
        >
          <div class="p-6">
            <div class="flex justify-between items-start mb-4">
              <h3 class="text-lg font-semibold text-gray-800 line-clamp-2">{{ card.front }}</h3>
              <div class="flex space-x-2" @click.stop>
                <button @click="editCard(card)" class="text-gray-600 hover:text-gray-800">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
                  </svg>
                </button>
                <button @click="confirmDelete(card)" class="text-red-600 hover:text-red-800">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
                  </svg>
                </button>
              </div>
            </div>
            
            <div class="text-gray-600 mb-4 line-clamp-3">
              {{ card.back }}
            </div>
            
            <div v-if="card.tags && card.tags.length > 1" class="mt-3">
              <div class="flex flex-wrap gap-1">
                <span 
                  v-for="tag in card.tags.filter(t => t.id !== parseInt(tagId))" 
                  :key="tag.id" 
                  class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium" 
                  :style="{ backgroundColor: `${tag.color}20`, color: tag.color }"
                >
                  {{ tag.name }}
                </span>
              </div>
            </div>

            <div class="mt-4 text-xs text-gray-500">
              创建于 {{ formatDate(card.created_at) }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建卡片模态框 -->
    <div v-if="showCreateCardModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-xl shadow-lg w-full max-w-md">
        <div class="p-6">
          <h3 class="text-xl font-semibold text-gray-800 mb-4">创建新卡片</h3>
          <form @submit.prevent="createCard">
            <div class="mb-4">
              <label for="cardFront" class="block text-sm font-medium text-gray-700 mb-1">正面</label>
              <textarea 
                id="cardFront" 
                v-model="newCard.front" 
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
                rows="3"
                required
              ></textarea>
            </div>
            <div class="mb-4">
              <label for="cardBack" class="block text-sm font-medium text-gray-700 mb-1">背面</label>
              <textarea 
                id="cardBack" 
                v-model="newCard.back" 
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
                rows="3"
                required
              ></textarea>
            </div>
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">标签</label>
              <select v-model="newCard.tag_ids" multiple class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500">
                <option v-for="tag in tags" :key="tag.id" :value="tag.id">{{ tag.name }}</option>
              </select>
              <p class="mt-1 text-xs text-gray-500">按住Ctrl/Cmd键可选择多个标签</p>
            </div>
            <div class="flex justify-end space-x-3">
              <button 
                type="button" 
                @click="showCreateCardModal = false" 
                class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
              >
                取消
              </button>
              <button 
                type="submit" 
                class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700"
              >
                创建
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- 编辑卡片模态框 -->
    <div v-if="showEditCardModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-xl shadow-lg w-full max-w-md">
        <div class="p-6">
          <h3 class="text-xl font-semibold text-gray-800 mb-4">编辑卡片</h3>
          <form @submit.prevent="updateCard">
            <div class="mb-4">
              <label for="editCardFront" class="block text-sm font-medium text-gray-700 mb-1">正面</label>
              <textarea 
                id="editCardFront" 
                v-model="editingCard.front" 
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
                rows="3"
                required
              ></textarea>
            </div>
            <div class="mb-4">
              <label for="editCardBack" class="block text-sm font-medium text-gray-700 mb-1">背面</label>
              <textarea 
                id="editCardBack" 
                v-model="editingCard.back" 
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
                rows="3"
                required
              ></textarea>
            </div>
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">标签</label>
              <select v-model="editingCard.tag_ids" multiple class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500">
                <option v-for="tag in tags" :key="tag.id" :value="tag.id">{{ tag.name }}</option>
              </select>
              <p class="mt-1 text-xs text-gray-500">按住Ctrl/Cmd键可选择多个标签</p>
            </div>
            <div class="flex justify-end space-x-3">
              <button 
                type="button" 
                @click="showEditCardModal = false" 
                class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
              >
                取消
              </button>
              <button 
                type="submit" 
                class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700"
              >
                保存
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- 删除确认模态框 -->
    <div v-if="showDeleteCardModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-xl shadow-lg w-full max-w-md">
        <div class="p-6">
          <h3 class="text-xl font-semibold text-gray-800 mb-4">确认删除</h3>
          <p class="text-gray-600 mb-6">
            确定要删除这张卡片吗？此操作不可撤销。
          </p>
          <div class="flex justify-end space-x-3">
            <button 
              type="button" 
              @click="showDeleteCardModal = false" 
              class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
            >
              取消
            </button>
            <button 
              type="button" 
              @click="deleteCard" 
              class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
            >
              删除
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import api from '../services/api'
import { ElMessage } from 'element-plus'

export default {
  name: 'TagCardsView',
  data() {
    return {
      tag: {},
      deck: {},
      cards: [],
      tags: [],
      loading: true,
      searchQuery: '',
      showCreateCardModal: false,
      showEditCardModal: false,
      showDeleteCardModal: false,
      newCard: {
        front: '',
        back: '',
        deck_id: null,
        tag_ids: []
      },
      editingCard: {},
      deletingCard: {}
    }
  },
  computed: {
    tagId() {
      return this.$route.params.id
    },
    filteredCards() {
      if (!this.searchQuery) {
        return this.cards
      }
      
      const query = this.searchQuery.toLowerCase()
      return this.cards.filter(card => 
        card.front.toLowerCase().includes(query) || 
        card.back.toLowerCase().includes(query)
      )
    }
  },
  created() {
    this.fetchTag()
    this.fetchCards()
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
    async fetchTag() {
      try {
        const response = await api.getTag(this.tagId)
        this.tag = response.data || {}
        
        // 获取卡包信息
        if (this.tag.deck_id) {
          const deckResponse = await api.getDeck(this.tag.deck_id)
          this.deck = deckResponse.data || {}
          
          // 获取卡包的所有标签
          const tagsResponse = await api.getTags(this.tag.deck_id)
          this.tags = tagsResponse.data || []
          
          // 设置新卡片的默认值
          this.newCard.deck_id = this.tag.deck_id
          this.newCard.tag_ids = [parseInt(this.tagId)]
        }
      } catch (error) {
        console.error('获取标签信息失败:', error)
        ElMessage.error('获取标签信息失败，请稍后重试')
      }
    },
    async fetchCards() {
      this.loading = true
      try {
        const response = await api.getCardsByTag(this.tagId)
        this.cards = response.data || []
      } catch (error) {
        console.error('获取卡片列表失败:', error)
        ElMessage.error('获取卡片列表失败，请稍后重试')
      } finally {
        this.loading = false
      }
    },
    async createCard() {
      try {
        await api.createCard(this.newCard)
        this.showCreateCardModal = false
        this.newCard = { 
          front: '', 
          back: '', 
          deck_id: this.tag.deck_id, 
          tag_ids: [parseInt(this.tagId)] 
        }
        ElMessage.success('卡片创建成功！')
        this.fetchCards()
      } catch (error) {
        console.error('创建卡片失败:', error)
        ElMessage.error('创建卡片失败，请稍后重试')
      }
    },
    editCard(card) {
      this.editingCard = { 
        ...card, 
        tag_ids: card.tags ? card.tags.map(tag => tag.id) : [] 
      }
      this.showEditCardModal = true
    },
    async updateCard() {
      try {
        await api.updateCard(this.editingCard.id, this.editingCard)
        this.showEditCardModal = false
        ElMessage.success('卡片更新成功！')
        this.fetchCards()
      } catch (error) {
        console.error('更新卡片失败:', error)
        ElMessage.error('更新卡片失败，请稍后重试')
      }
    },
    confirmDelete(card) {
      this.deletingCard = card
      this.showDeleteCardModal = true
    },
    async deleteCard() {
      try {
        await api.deleteCard(this.deletingCard.id)
        this.showDeleteCardModal = false
        ElMessage.success('卡片删除成功！')
        this.fetchCards()
      } catch (error) {
        console.error('删除卡片失败:', error)
        ElMessage.error('删除卡片失败，请稍后重试')
      }
    },
    viewCard(card) {
      this.$router.push(`/cards/${card.id}`)
    },
    goBack() {
      this.$router.push('/tags')
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
        if (this.showCreateCardModal) {
          this.showCreateCardModal = false
        } else if (this.showEditCardModal) {
          this.showEditCardModal = false
        } else if (this.showDeleteCardModal) {
          this.showDeleteCardModal = false
        }
      }
    }
  }
}
</script>