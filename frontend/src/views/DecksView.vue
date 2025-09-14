<template>
  <div class="min-h-screen bg-gray-50 p-4">
    <div class="max-w-6xl mx-auto">
      <header class="flex justify-between items-center mb-8">
        <h1 class="text-3xl font-bold text-gray-800">卡包管理</h1>
        <button @click="showCreateModal = true" class="bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700 transition-colors flex items-center">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clip-rule="evenodd" />
          </svg>
          新建卡包
        </button>
      </header>

      <!-- 卡包列表 -->
      <div v-if="loading" class="text-center py-8">
        <div class="inline-block animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-600"></div>
        <p class="mt-2 text-gray-600">加载中...</p>
      </div>

      <div v-else-if="decks.length === 0" class="bg-white rounded-xl shadow p-8 text-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto text-gray-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
        </svg>
        <h3 class="text-xl font-medium text-gray-900 mb-2">暂无卡包</h3>
        <p class="text-gray-500 mb-4">创建您的第一个学习卡包开始使用</p>
        <button @click="showCreateModal = true" class="bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700 transition-colors">
          创建卡包
        </button>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="deck in decks" :key="deck.id" class="bg-white rounded-xl shadow hover:shadow-lg transition-shadow">
          <div class="p-6">
            <div class="flex justify-between items-start mb-4">
              <h3 class="text-lg font-semibold text-gray-800">{{ deck.name }}</h3>
              <div class="flex space-x-2">
                <span v-if="deck.archived" class="bg-gray-100 text-gray-800 text-xs px-2 py-1 rounded">已归档</span>
                <button @click="editDeck(deck)" class="text-gray-600 hover:text-gray-800" title="编辑卡包">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
                  </svg>
                </button>
                <button @click="confirmDelete(deck)" class="text-red-600 hover:text-red-800" title="删除卡包">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
                  </svg>
                </button>
              </div>
            </div>
            
            <div v-if="deck.stats" class="grid grid-cols-2 gap-2 mb-4 text-sm">
              <div class="bg-blue-50 p-2 rounded">
                <div class="text-blue-800 font-medium">{{ deck.stats.total_cards }}</div>
                <div class="text-blue-600">总卡片</div>
              </div>
              <div class="bg-yellow-50 p-2 rounded">
                <div class="text-yellow-800 font-medium">{{ deck.stats.due_cards }}</div>
                <div class="text-yellow-600">待复习</div>
              </div>
              <div class="bg-green-50 p-2 rounded">
                <div class="text-green-800 font-medium">{{ deck.stats.tag_count }}</div>
                <div class="text-green-600">标签数</div>
              </div>
              <div class="bg-purple-50 p-2 rounded">
                <div class="text-purple-800 font-medium">{{ deck.stats.today_studied }}</div>
                <div class="text-purple-600">今日学习</div>
              </div>
            </div>

            <div class="mt-4 space-y-2">
              <button @click="startDeckStudy(deck)" class="w-full bg-green-600 text-white hover:bg-green-700 py-2 px-4 rounded-lg text-sm font-medium transition-colors">
                开始学习
              </button>
              <button @click="viewDeckCards(deck)" class="w-full bg-indigo-50 text-indigo-700 hover:bg-indigo-100 py-2 px-4 rounded-lg text-sm font-medium transition-colors">
                查看相关卡片
              </button>
            </div>

            <div class="mt-4 text-xs text-gray-500">
              创建于 {{ formatDate(deck.created_at) }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建卡包模态框 -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-xl shadow-lg w-full max-w-md">
        <div class="p-6">
          <h3 class="text-xl font-semibold text-gray-800 mb-4">创建新卡包</h3>
          <form @submit.prevent="createDeck">
            <div class="mb-4">
              <label for="deckName" class="block text-sm font-medium text-gray-700 mb-1">卡包名称</label>
              <input 
                type="text" 
                id="deckName" 
                v-model="newDeck.name" 
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
                required
              >
            </div>
            <div class="flex justify-end space-x-3">
              <button 
                type="button" 
                @click="showCreateModal = false" 
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

    <!-- 编辑卡包模态框 -->
    <div v-if="showEditModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-xl shadow-lg w-full max-w-md">
        <div class="p-6">
          <h3 class="text-xl font-semibold text-gray-800 mb-4">编辑卡包</h3>
          <form @submit.prevent="updateDeck">
            <div class="mb-4">
              <label for="editDeckName" class="block text-sm font-medium text-gray-700 mb-1">卡包名称</label>
              <input 
                type="text" 
                id="editDeckName" 
                v-model="editingDeck.name" 
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
                required
              >
            </div>
            <div class="mb-4">
              <div class="flex items-center">
                <input 
                  type="checkbox" 
                  id="editDeckArchived" 
                  v-model="editingDeck.archived" 
                  class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                >
                <label for="editDeckArchived" class="ml-2 block text-sm text-gray-700">
                  归档此卡包
                </label>
              </div>
            </div>
            <div class="flex justify-end space-x-3">
              <button 
                type="button" 
                @click="showEditModal = false" 
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
    <div v-if="showDeleteModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-xl shadow-lg w-full max-w-md">
        <div class="p-6">
          <h3 class="text-xl font-semibold text-gray-800 mb-4">确认删除</h3>
          <p class="text-gray-600 mb-6">
            确定要删除卡包 "{{ deletingDeck.name }}" 吗？此操作不可撤销。
          </p>
          <div class="flex justify-end space-x-3">
            <button 
              type="button" 
              @click="showDeleteModal = false" 
              class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
            >
              取消
            </button>
            <button 
              type="button" 
              @click="deleteDeck" 
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
  name: 'DecksView',
  data() {
    return {
      decks: [],
      loading: false,
      showCreateModal: false,
      showEditModal: false,
      showDeleteModal: false,
      newDeck: {
        name: ''
      },
      editingDeck: {},
      deletingDeck: {}
    }
  },
  created() {
    this.fetchDecks()
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
    async fetchDecks() {
      this.loading = true
      try {
        const response = await api.getDecks(true)
        this.decks = response.data?.data?.decks || []
      } catch (error) {
        console.error('获取卡包失败:', error)
        ElMessage.error('获取卡包失败，请稍后重试')
      } finally {
        this.loading = false
      }
    },
    async createDeck() {
      this.loading = true
      try {
        await api.createDeck(this.newDeck)
        this.newDeck = { name: '' }
        this.showCreateModal = false
        ElMessage.success('卡包创建成功！')
        this.fetchDecks()
      } catch (error) {
        console.error('创建卡包失败:', error)
        ElMessage.error('创建卡包失败，请稍后重试')
      } finally {
        this.loading = false
      }
    },
    editDeck(deck) {
      this.editingDeck = { ...deck }
      this.showEditModal = true
    },
    async updateDeck() {
      try {
        await api.updateDeck(this.editingDeck.id, {
          name: this.editingDeck.name,
          archived: this.editingDeck.archived
        })
        this.showEditModal = false
        this.fetchDecks()
        ElMessage.success('卡包更新成功！')
      } catch (error) {
        console.error('更新卡包失败:', error)
        ElMessage.error('更新卡包失败，请稍后重试')
      }
    },
    confirmDelete(deck) {
      this.deletingDeck = deck
      this.showDeleteModal = true
    },
    async deleteDeck() {
      try {
        if (!this.deletingDeck || !this.deletingDeck.id) {
          ElMessage.error('卡包信息不存在，请刷新页面重试')
          return
        }
        await api.deleteDeck(this.deletingDeck.id)
        this.showDeleteModal = false
        ElMessage.success('卡包删除成功！')
        this.fetchDecks()
      } catch (error) {
        console.error('删除卡包失败:', error)
        ElMessage.error('删除卡包失败，请稍后重试')
      }
    },
    startDeckStudy(deck) {
      if (!deck || !deck.id) {
        ElMessage.error('卡包信息不存在，请刷新页面重试')
        return
      }
      // 跳转到学习页面，并传递卡包ID
      this.$router.push({
        path: '/study',
        query: { deckId: deck.id }
      })
    },
    viewDeckCards(deck) {
      if (!deck || !deck.id) {
        ElMessage.error('卡包信息不存在，请刷新页面重试')
        return
      }
      // 跳转到卡片管理页面并搜索该卡包
      this.$router.push({
        path: '/cards',
        query: { deck_id: deck.id }
      })
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
        if (this.showCreateModal) {
          this.showCreateModal = false
        } else if (this.showEditModal) {
          this.showEditModal = false
        } else if (this.showDeleteModal) {
          this.showDeleteModal = false
        }
      }
    }
  }
}
</script>