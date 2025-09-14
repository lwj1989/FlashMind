<template>
  <div class="min-h-screen bg-gray-50 p-4">
    <div class="max-w-6xl mx-auto">
      <header class="mb-8">
        <h1 class="text-3xl font-bold text-gray-800">标签管理</h1>
        <p class="text-gray-600 mt-2">管理所有卡包的标签</p>
      </header>

      <div class="bg-white rounded-xl shadow p-6 mb-6">
        <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
          <div class="flex-1">
            <div class="relative">
              <input 
                type="text" 
                v-model="searchQuery" 
                placeholder="搜索标签..." 
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
            <button @click="showCreateModal = true" class="bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700 transition-colors flex items-center">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clip-rule="evenodd" />
              </svg>
              新建标签
            </button>
          </div>
        </div>
      </div>

      <!-- 标签列表 -->
      <div v-if="loading" class="text-center py-8">
        <div class="inline-block animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-600"></div>
        <p class="mt-2 text-gray-600">加载中...</p>
      </div>

      <div v-else-if="filteredTags.length === 0" class="bg-white rounded-xl shadow p-8 text-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto text-gray-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
        </svg>
        <h3 class="text-xl font-medium text-gray-900 mb-2">暂无标签</h3>
        <p class="text-gray-500 mb-4">创建您的第一个标签开始使用</p>
        <button @click="showCreateModal = true" class="bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700 transition-colors">
          创建标签
        </button>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="tag in filteredTags" :key="tag.id" class="bg-white rounded-xl shadow hover:shadow-lg transition-shadow">
          <div class="p-6">
            <div class="flex justify-between items-start mb-4">
              <h3 class="text-lg font-semibold text-gray-800">{{ tag.name }}</h3>
              <div class="flex space-x-2">
                <button @click="editTag(tag)" class="text-gray-600 hover:text-gray-800">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
                  </svg>
                </button>
                <button @click="confirmDelete(tag)" class="text-red-600 hover:text-red-800">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
                  </svg>
                </button>
              </div>
            </div>
            

            
            <div v-if="tag.stats" class="grid grid-cols-2 gap-2 text-sm">
              <div class="bg-blue-50 p-2 rounded">
                <div class="text-blue-800 font-medium">{{ tag.stats.total_cards }}</div>
                <div class="text-blue-600">总卡片</div>
              </div>
              <div class="bg-yellow-50 p-2 rounded">
                <div class="text-yellow-800 font-medium">{{ tag.stats.due_cards }}</div>
                <div class="text-yellow-600">待复习</div>
              </div>
            </div>

            <div class="mt-4 space-y-2">
              <button @click="startTagStudy(tag)" class="w-full bg-green-600 text-white hover:bg-green-700 py-2 px-4 rounded-lg text-sm font-medium transition-colors">
                开始学习
              </button>
              <button @click="viewTagCards(tag)" class="w-full bg-indigo-50 text-indigo-700 hover:bg-indigo-100 py-2 px-4 rounded-lg text-sm font-medium transition-colors">
                查看相关卡片
              </button>
            </div>

            <div class="mt-4 text-xs text-gray-500">
              创建于 {{ formatDate(tag.created_at) }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建标签模态框 -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-xl shadow-lg w-full max-w-md">
        <div class="p-6">
          <h3 class="text-xl font-semibold text-gray-800 mb-4">创建新标签</h3>
          <form @submit.prevent="createTag">
            <div class="mb-4">
              <label for="tagName" class="block text-sm font-medium text-gray-700 mb-1">标签名称</label>
              <input 
                type="text" 
                id="tagName" 
                v-model="newTag.name" 
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

    <!-- 编辑标签模态框 -->
    <div v-if="showEditModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-xl shadow-lg w-full max-w-md">
        <div class="p-6">
          <h3 class="text-xl font-semibold text-gray-800 mb-4">编辑标签</h3>
          <form @submit.prevent="updateTag">
            <div class="mb-4">
              <label for="editTagName" class="block text-sm font-medium text-gray-700 mb-1">标签名称</label>
              <input 
                type="text" 
                id="editTagName" 
                v-model="editingTag.name" 
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
                required
              >
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
            确定要删除标签 "{{ deletingTag.name }}" 吗？此操作不可撤销。
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
              @click="deleteTag" 
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
import { getAllTags, createTag, updateTag, deleteTag } from '../api/tag'
import { ElMessage } from 'element-plus'

export default {
  name: 'TagsView',
  data() {
    return {
      tags: [],
      loading: false,
      showCreateModal: false,
      showEditModal: false,
      showDeleteModal: false,
      searchQuery: '',
      newTag: {
        name: ''
      },
      editingTag: {},
      deletingTag: {}
    }
  },
  computed: {
    filteredTags() {
      let result = this.tags
      
      // 按搜索词筛选
      if (this.searchQuery) {
        const query = this.searchQuery.toLowerCase()
        result = result.filter(tag => 
          tag.name.toLowerCase().includes(query)
        )
      }
      
      return result
    }
  },
  created() {
    this.fetchAllTags()
  },
  methods: {

    async fetchAllTags() {
      this.loading = true
      try {
        const response = await getAllTags()
        console.log('获取标签响应:', response)
        // 后端使用统一响应格式：{ code, message, data }
        // 当include_stats=true时，data是TagWithStats数组
        if (response.data && response.data.code === 'SUCCESS') {
          this.tags = response.data.data || []
        } else {
          // 兼容直接返回数组的情况
          this.tags = response.data || []
        }
      } catch (error) {
        console.error('获取标签失败:', error)
        ElMessage.error('获取标签失败，请稍后重试')
      } finally {
        this.loading = false
      }
    },
    async createTag() {
      try {
        const tagData = {
          name: this.newTag.name
        }
        const response = await createTag(tagData)
        console.log('创建标签响应:', response)
        this.showCreateModal = false
        this.newTag = { name: '' }
        ElMessage.success('标签创建成功！')
        this.fetchAllTags()
      } catch (error) {
        console.error('创建标签失败:', error)
        ElMessage.error('创建标签失败，请稍后重试')
      }
    },
    editTag(tag) {
      this.editingTag = { ...tag }
      this.showEditModal = true
    },
    async updateTag() {
      try {
        if (!this.editingTag || !this.editingTag.id) {
          ElMessage.error('标签信息不完整，请刷新页面重试')
          return
        }
        const tagData = {
          name: this.editingTag.name
        }
        const response = await updateTag(this.editingTag.id, tagData)
        console.log('更新标签响应:', response)
        this.showEditModal = false
        this.fetchAllTags()
        ElMessage.success('标签更新成功！')
      } catch (error) {
        console.error('更新标签失败:', error)
        ElMessage.error('更新标签失败，请稍后重试')
      }
    },
    confirmDelete(tag) {
      this.deletingTag = tag
      this.showDeleteModal = true
    },
    async deleteTag() {
      try {
        if (!this.deletingTag || !this.deletingTag.id) {
          ElMessage.error('标签信息不存在，请刷新页面重试')
          return
        }
        const response = await deleteTag(this.deletingTag.id)
        console.log('删除标签响应:', response)
        this.showDeleteModal = false
        ElMessage.success('标签删除成功！')
        this.fetchAllTags()
      } catch (error) {
        console.error('删除标签失败:', error)
        ElMessage.error('删除标签失败，请稍后重试')
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
    viewTagCards(tag) {
      if (!tag || !tag.id) {
        ElMessage.error('标签信息不存在，请刷新页面重试')
        return
      }
      this.$router.push(`/tags/${tag.id}/cards`)
    },
    startTagStudy(tag) {
      if (!tag || !tag.id) {
        ElMessage.error('标签信息不存在，请刷新页面重试')
        return
      }
      this.$router.push({
        path: '/study',
        query: { tagId: tag.id }
      })
    }
  },
  watch: {
    decks() {
      // 当卡包列表加载完成后，获取所有标签
      if (this.decks.length > 0) {
        this.fetchAllTags()
      }
    }
  }
}
</script>