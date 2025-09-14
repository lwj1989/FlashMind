<template>
  <div class="card-list">
    <div class="mb-6">
      <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div class="relative flex-1">
          <input 
            type="text" 
            v-model="searchQuery" 
            placeholder="搜索卡片..." 
            @input="handleSearch"
            class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
          <div class="absolute left-3 top-2.5 text-gray-400">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd" />
            </svg>
          </div>
        </div>
        
        <!-- 批量操作按钮 -->
        <div v-if="showBatchActions" class="flex items-center space-x-3">
          <button 
            @click="toggleSelectAll" 
            class="px-3 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors"
          >
            {{ selectedCards.length === cards.length && cards.length > 0 ? '取消全选' : '全选' }}
          </button>
          <button 
            @click="$emit('batch-delete', selectedCards)" 
            :disabled="selectedCards.length === 0"
            class="px-3 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            批量删除 ({{ selectedCards.length }})
          </button>
        </div>
      </div>
      
      <!-- 批量操作提示 -->
      <div v-if="showBatchActions && selectedCards.length > 0" class="mt-2 text-sm text-gray-600">
        已选择 {{ selectedCards.length }} 张卡片
      </div>
    </div>

    <div v-if="loading" class="text-center py-8">
      <div class="inline-block animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-600"></div>
      <p class="mt-2 text-gray-600">加载中...</p>
    </div>

    <div v-else-if="!cards || cards.length === 0" class="bg-white rounded-xl shadow p-8 text-center">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto text-gray-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
      </svg>
      <h3 class="text-xl font-medium text-gray-900 mb-2">暂无卡片</h3>
      <p class="text-gray-500 mb-4">创建您的第一张学习卡片开始使用</p>
      <button v-if="showCreateButton" @click="$emit('create')" class="bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700 transition-colors">
        创建卡片
      </button>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="card in (cards || [])" :key="card.id" 
           :class="[
             'bg-white rounded-xl shadow hover:shadow-lg transition-shadow',
             { 'ring-2 ring-indigo-500': showBatchActions && isCardSelected(card.id) }
           ]">
        <div class="p-6">
          <!-- 批量选择复选框 -->
          <div v-if="showBatchActions" class="flex items-center justify-between mb-4">
            <input 
              type="checkbox" 
              :checked="isCardSelected(card.id)"
              @change="toggleCardSelection(card.id)"
              class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
            />
            <span class="text-sm text-gray-500">ID: {{ card.id }}</span>
          </div>
          
          <div class="mb-4">
            <h3 class="text-lg font-semibold text-gray-800 mb-2">问题</h3>
            <p class="text-gray-600">{{ card.question }}</p>
          </div>
          
          <div class="mb-4">
            <h3 class="text-lg font-semibold text-gray-800 mb-2">答案</h3>
            <p class="text-gray-600">{{ card.answer }}</p>
          </div>
          
          <div v-if="card.deck_name || card.tag_name" class="mb-4">
            <div v-if="card.deck_name" class="text-sm text-gray-600 mb-1">
              <span class="font-medium">卡包：</span>{{ card.deck_name }}
            </div>
            <div v-if="card.tag_name" class="text-sm text-gray-600">
              <span class="font-medium">标签：</span>{{ card.tag_name }}
            </div>
          </div>
          
          <div class="flex justify-end space-x-2">
            <button @click="$emit('edit', card)" class="text-indigo-600 hover:text-indigo-800">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
              </svg>
            </button>
            <button @click="$emit('delete', card)" class="text-red-600 hover:text-red-800">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="pagination && pagination.total > (pagination.pageSize || 20)" class="mt-8 flex justify-center items-center space-x-4">
      <button 
        @click="handlePageChange((pagination.page || 1) - 1)" 
        :disabled="(pagination.page || 1) === 1"
        class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        上一页
      </button>
      <span class="text-gray-600">
        第 {{ pagination.page || 1 }} 页，共 {{ pagination.totalPages || 1 }} 页
      </span>
      <button 
        @click="handlePageChange((pagination.page || 1) + 1)" 
        :disabled="(pagination.page || 1) === (pagination.totalPages || 1)"
        class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        下一页
      </button>
    </div>
  </div>
</template>

<script>
import { ref, computed, watch } from 'vue'

export default {
  name: 'CardList',
  props: {
    cards: {
      type: Array,
      default: () => []
    },
    loading: {
      type: Boolean,
      default: false
    },
    title: {
      type: String,
      default: '卡片列表'
    },
    showCreateButton: {
      type: Boolean,
      default: true
    },
    pagination: {
      type: Object,
      default: () => ({
        page: 1,
        pageSize: 20,
        total: 0,
        totalPages: 0
      })
    },
    showBatchActions: {
      type: Boolean,
      default: false
    }
  },
  emits: ['search', 'page-change', 'create', 'edit', 'delete', 'batch-delete'],
  setup(props, { emit }) {
    const searchQuery = ref('')
    const selectedCards = ref([])
    let searchTimeout = null

    // 处理搜索
    const handleSearch = () => {
      clearTimeout(searchTimeout)
      searchTimeout = setTimeout(() => {
        emit('search', searchQuery.value)
      }, 500)
    }

    // 处理分页
    const handlePageChange = (page) => {
      if (page < 1) return
      emit('page-change', page)
    }

    // 检查卡片是否被选中
    const isCardSelected = (cardId) => {
      return selectedCards.value.includes(cardId)
    }

    // 切换卡片选择状态
    const toggleCardSelection = (cardId) => {
      const index = selectedCards.value.indexOf(cardId)
      if (index > -1) {
        selectedCards.value.splice(index, 1)
      } else {
        selectedCards.value.push(cardId)
      }
    }

    // 切换全选状态
    const toggleSelectAll = () => {
      if (selectedCards.value.length === props.cards.length && props.cards.length > 0) {
        selectedCards.value = []
      } else {
        selectedCards.value = props.cards.map(card => card.id)
      }
    }

    // 监听卡片列表变化，清理无效的选择
    watch(() => props.cards, (newCards) => {
      if (newCards) {
        const validIds = newCards.map(card => card.id)
        selectedCards.value = selectedCards.value.filter(id => validIds.includes(id))
      }
    })

    return {
      searchQuery,
      selectedCards,
      handleSearch,
      handlePageChange,
      isCardSelected,
      toggleCardSelection,
      toggleSelectAll
    }
  }
}
</script>