<template>
  <div class="min-h-screen bg-gray-50 p-4">
    <div class="max-w-6xl mx-auto">
      <header class="flex justify-between items-center mb-8">
        <h1 class="text-3xl font-bold text-gray-800">卡片管理</h1>
        <div class="flex items-center space-x-4">
          <button 
            @click="toggleBatchMode" 
            :class="[
              'px-4 py-2 rounded-lg transition-colors flex items-center',
              showBatchMode ? 'bg-gray-600 text-white hover:bg-gray-700' : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            ]"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M3 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd" />
            </svg>
            {{ showBatchMode ? '退出批量' : '批量管理' }}
          </button>
          <button 
            @click="openForm()" 
            class="bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700 transition-colors flex items-center"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clip-rule="evenodd" />
            </svg>
            创建卡片
          </button>
        </div>
      </header>

      <!-- 搜索过滤器 -->
      <div class="bg-white rounded-xl shadow p-6 mb-6">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">关键词搜索</label>
            <input 
              v-model="searchForm.keyword" 
              type="text" 
              placeholder="搜索问题或答案..."
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
              @keyup.enter="fetchCards"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">卡包</label>
            <select 
              v-model="searchForm.deck_id" 
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
              @change="fetchCards"
            >
              <option value="">全部卡包</option>
              <option v-for="deck in decks" :key="deck.id" :value="deck.id">{{ deck.name }}</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">标签</label>
            <select 
              v-model="searchForm.tag_id" 
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
              @change="fetchCards"
            >
              <option value="">全部标签</option>
              <option v-for="tag in tags" :key="tag.id" :value="tag.id">{{ tag.name }}</option>
            </select>
          </div>
          <div class="flex items-end">
            <button 
              @click="fetchCards" 
              class="w-full bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 transition-colors"
            >
              搜索
            </button>
          </div>
        </div>
      </div>

      <!-- 卡片列表表格 -->
      <div class="bg-white rounded-xl shadow overflow-hidden">
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th v-if="showBatchMode" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  <input 
                    type="checkbox" 
                    :checked="selectedCards.length === cardList.length && cardList.length > 0"
                    @change="toggleSelectAll"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">问题</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">卡包</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">标签</th>
                <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-if="loading" class="text-center">
                <td :colspan="showBatchMode ? 5 : 4" class="px-6 py-8 text-gray-500">
                  <div class="flex items-center justify-center">
                    <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-indigo-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    加载中...
                  </div>
                </td>
              </tr>
              <tr v-else-if="cardList.length === 0" class="text-center">
                <td :colspan="showBatchMode ? 5 : 4" class="px-6 py-8 text-gray-500">
                  暂无卡片数据
                </td>
              </tr>
              <tr v-else v-for="card in cardList" :key="card.id" class="hover:bg-gray-50">
                <td v-if="showBatchMode" class="px-6 py-4 whitespace-nowrap">
                  <input 
                    type="checkbox" 
                    :checked="selectedCards.includes(card.id)"
                    @change="toggleCardSelection(card.id)"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                </td>
                <td class="px-6 py-4">
                  <div 
                    class="text-sm text-gray-900 max-w-xs truncate cursor-pointer text-blue-600 hover:text-blue-800 hover:underline" 
                    :title="card.question"
                    @click="viewCardDetail(card)"
                  >
                    {{ card.question }}
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="text-sm text-gray-500">{{ getDeckName(card.deck_id) }}</div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span v-if="card.tag_name" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800">
                    {{ card.tag_name }}
                  </span>
                  <span v-else class="text-gray-400 text-sm">无标签</span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <button 
                    @click="openForm(card)" 
                    class="text-indigo-600 hover:text-indigo-900 mr-3 p-1"
                    title="编辑卡片"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                      <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
                    </svg>
                  </button>
                  <button 
                    @click="deleteCard(card)" 
                    class="text-red-600 hover:text-red-900 p-1"
                    title="删除卡片"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                      <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
                    </svg>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- 批量操作栏 -->
        <div v-if="showBatchMode && selectedCards.length > 0" class="bg-gray-50 px-6 py-3 border-t border-gray-200">
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-700">已选择 {{ selectedCards.length }} 张卡片</span>
            <button 
              @click="batchDeleteCards" 
              class="bg-red-600 text-white px-4 py-2 rounded-lg hover:bg-red-700 transition-colors"
            >
              批量删除
            </button>
          </div>
        </div>

        <!-- 分页 -->
        <div v-if="pagination.total > 0" class="bg-gray-50 px-6 py-3 border-t border-gray-200">
          <div class="flex items-center justify-between">
            <div class="text-sm text-gray-700">
              显示第 {{ (pagination.page - 1) * pagination.page_size + 1 }} 到 
              {{ Math.min(pagination.page * pagination.page_size, pagination.total) }} 条，
              共 {{ pagination.total }} 条记录
            </div>
            <nav class="flex items-center space-x-2">
              <button 
                @click="changePage(pagination.page - 1)"
                :disabled="pagination.page <= 1"
                class="px-3 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                上一页
              </button>
              
              <template v-for="page in getPageNumbers()" :key="page">
                <button 
                  v-if="typeof page === 'number'"
                  @click="changePage(page)"
                  :class="[
                    'px-3 py-2 border rounded-md',
                    page === pagination.page 
                      ? 'border-indigo-500 bg-indigo-600 text-white' 
                      : 'border-gray-300 text-gray-700 hover:bg-gray-50'
                  ]"
                >
                  {{ page }}
                </button>
                <span v-else class="px-3 py-2 text-gray-500">...</span>
              </template>
              
              <button 
                @click="changePage(pagination.page + 1)"
                :disabled="pagination.page >= pagination.total_pages"
                class="px-3 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                下一页
              </button>
            </nav>
          </div>
        </div>
      </div>

      <!-- 大尺寸卡片表单对话框 -->
      <div v-if="showForm" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50" @click="closeForm">
        <div class="bg-white rounded-xl shadow-lg w-full max-w-5xl max-h-[90vh] overflow-hidden" @click.stop>
          <div class="p-6 border-b border-gray-200">
            <h3 class="text-xl font-semibold text-gray-800">
              {{ form.id ? '编辑卡片' : '创建卡片' }}
            </h3>
          </div>
          <div class="p-6 overflow-y-auto max-h-[70vh]">
            <form @submit.prevent="handleSubmit">
              <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <!-- 基本信息 -->
                <div>
                  <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-2">卡包 *</label>
                    <select 
                      v-model="form.deck_id" 
                      required
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                    >
                      <option value="">请选择卡包</option>
                      <option v-for="deck in decks" :key="deck.id" :value="deck.id">{{ deck.name }}</option>
                    </select>
                  </div>
                  
                  <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-2">标签</label>
                    <select 
                      v-model="form.tag_id" 
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                    >
                      <option value="">无标签</option>
                      <option v-for="tag in tags" :key="tag.id" :value="tag.id">{{ tag.name }}</option>
                    </select>
                  </div>
                </div>

                <!-- 问题 -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">问题 *</label>
                  <div class="border border-gray-300 rounded-md overflow-hidden">
                    <div class="bg-gray-50 px-3 py-2 border-b border-gray-300">
                      <div class="flex space-x-2 text-sm">
                        <button 
                          type="button"
                          @click="questionMode = 'edit'"
                          :class="questionMode === 'edit' ? 'text-indigo-600 font-medium' : 'text-gray-600'"
                        >
                          编辑
                        </button>
                        <span class="text-gray-400">|</span>
                        <button 
                          type="button"
                          @click="questionMode = 'preview'"
                          :class="questionMode === 'preview' ? 'text-indigo-600 font-medium' : 'text-gray-600'"
                        >
                          预览
                        </button>
                      </div>
                    </div>
                    <div v-if="questionMode === 'edit'" class="p-3">
                      <textarea 
                        v-model="form.question" 
                        required
                        rows="8"
                        placeholder="输入问题，支持 Markdown 格式..."
                        class="w-full resize-none border-0 focus:ring-0 focus:outline-none"
                      ></textarea>
                    </div>
                    <div v-else class="p-3 bg-white min-h-[200px]">
                      <div class="markdown-content" v-html="renderMarkdown(form.question)"></div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 答案 -->
              <div class="mt-6">
                <label class="block text-sm font-medium text-gray-700 mb-2">答案 *</label>
                <div class="border border-gray-300 rounded-md overflow-hidden">
                  <div class="bg-gray-50 px-3 py-2 border-b border-gray-300">
                    <div class="flex space-x-2 text-sm">
                      <button 
                        type="button"
                        @click="answerMode = 'edit'"
                        :class="answerMode === 'edit' ? 'text-indigo-600 font-medium' : 'text-gray-600'"
                      >
                        编辑
                      </button>
                      <span class="text-gray-400">|</span>
                      <button 
                        type="button"
                        @click="answerMode = 'preview'"
                        :class="answerMode === 'preview' ? 'text-indigo-600 font-medium' : 'text-gray-600'"
                      >
                        预览
                      </button>
                    </div>
                  </div>
                  <div v-if="answerMode === 'edit'" class="p-3">
                    <textarea 
                      v-model="form.answer" 
                      required
                      rows="12"
                      placeholder="输入答案，支持 Markdown 格式。例如：
                      
# 标题
- 列表项
`inline code`
```javascript
// 代码块
console.log('Hello World');
```"
                      class="w-full resize-none border-0 focus:ring-0 focus:outline-none"
                    ></textarea>
                  </div>
                  <div v-else class="p-3 bg-white min-h-[300px]">
                    <div class="markdown-content" v-html="renderMarkdown(form.answer)"></div>
                  </div>
                </div>
              </div>

              <!-- 表单按钮 -->
              <div class="mt-6 flex justify-end space-x-3">
                <button 
                  type="button" 
                  @click="closeForm"
                  class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
                >
                  取消
                </button>
                <button 
                  type="submit"
                  class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700"
                >
                  {{ form.id ? '更新' : '创建' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>

      <!-- 删除确认对话框 -->
      <div v-if="deletingCard || batchDeletingCards.length > 0" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
        <div class="bg-white rounded-xl shadow-lg w-full max-w-md">
          <div class="p-6">
            <h3 class="text-xl font-semibold text-red-600 mb-4">确认删除</h3>
            <p class="text-gray-600 mb-6">
              {{ deletingCard 
                ? `确定要删除卡片"${deletingCard.question}"吗？` 
                : `确定要删除选中的 ${batchDeletingCards.length} 张卡片吗？`
              }}
              <br>此操作不可撤销。
            </p>
            <div class="flex justify-end space-x-3">
              <button 
                @click="cancelDelete" 
                :disabled="isDeleting"
                class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 disabled:opacity-50"
              >
                取消
              </button>
              <button 
                @click="confirmDelete" 
                :disabled="isDeleting"
                class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50"
              >
                {{ isDeleting ? '删除中...' : '确认删除' }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 卡片详情查看模态框 -->
      <div v-if="showCardDetailModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50" @click="closeCardDetailModal">
        <div class="bg-white rounded-xl shadow-lg w-full max-w-4xl max-h-[90vh] overflow-hidden" @click.stop>
          <div class="p-6 border-b border-gray-200">
            <div class="flex justify-between items-center">
              <h3 class="text-xl font-semibold text-gray-800">卡片详情</h3>
              <button @click="closeCardDetailModal" class="text-gray-400 hover:text-gray-600">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                </svg>
              </button>
            </div>
          </div>
          
          <div class="p-6 overflow-y-auto max-h-[70vh]" v-if="selectedCardDetail">
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <!-- 问题部分 -->
              <div class="bg-blue-50 rounded-lg p-6">
                <h4 class="text-lg font-semibold text-blue-800 mb-4">问题</h4>
                <div class="markdown-content text-blue-900" v-html="renderMarkdown(selectedCardDetail.question)"></div>
              </div>

              <!-- 答案部分 -->
              <div class="bg-green-50 rounded-lg p-6">
                <h4 class="text-lg font-semibold text-green-800 mb-4">答案</h4>
                <div class="markdown-content text-green-900" v-html="renderMarkdown(selectedCardDetail.answer)"></div>
              </div>
            </div>

            <!-- 卡片信息 -->
            <div class="mt-6 bg-gray-50 rounded-lg p-4">
              <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
                <div>
                  <span class="text-gray-600">所属卡包：</span>
                  <span class="font-medium">{{ getDeckName(selectedCardDetail.deck_id) }}</span>
                </div>
                <div>
                  <span class="text-gray-600">标签：</span>
                  <span v-if="selectedCardDetail.tag_name" class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800">
                    {{ selectedCardDetail.tag_name }}
                  </span>
                  <span v-else class="text-gray-400">无标签</span>
                </div>
                <div>
                  <span class="text-gray-600">创建时间：</span>
                  <span class="font-medium">{{ formatDate(selectedCardDetail.created_at) }}</span>
                </div>
              </div>
            </div>

            <!-- 操作按钮 -->
            <div class="mt-6 flex justify-end space-x-3">
              <button 
                @click="editCardFromDetail" 
                class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
              >
                编辑卡片
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.markdown-content {
  color: #374151;
  line-height: 1.7;
}

.markdown-content :deep(pre) {
  background-color: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 16px;
  margin: 16px 0;
  overflow-x: auto;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.5;
}

.markdown-content :deep(code) {
  background-color: #f1f5f9;
  color: #e11d48;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 0.875em;
  font-weight: 500;
}

.markdown-content :deep(pre code) {
  background-color: transparent;
  color: inherit;
  padding: 0;
  border-radius: 0;
  font-weight: normal;
}
</style>

<script>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import { searchCards, createCard, updateCard, deleteCard } from '@/api/card'
import { getDecks } from '@/api/deck'
import { getAllTags } from '@/api/tag'

export default {
  name: 'CardsView',
  setup() {
    const loading = ref(false)
    const showForm = ref(false)
    const showBatchMode = ref(false)
    const selectedCards = ref([])
    const cardList = ref([])
    const decks = ref([])
    const tags = ref([])
    const deletingCard = ref(null)
    const batchDeletingCards = ref([])
    const isDeleting = ref(false)
    const questionMode = ref('edit')
    const answerMode = ref('edit')
    const showCardDetailModal = ref(false)
    const selectedCardDetail = ref(null)

    const searchForm = reactive({
      keyword: '',
      deck_id: '',
      tag_id: ''
    })

    const pagination = reactive({
      page: 1,
      page_size: 20,
      total: 0,
      total_pages: 0
    })

    const form = reactive({
      id: null,
      deck_id: '',
      tag_id: '',
      question: '',
      answer: ''
    })

    const defaultCard = {
      id: null,
      deck_id: '',
      tag_id: '',
      question: '',
      answer: ''
    }

    // 渲染 Markdown
    const renderMarkdown = (text) => {
      if (!text) return ''
      
      marked.setOptions({
        breaks: true,
        gfm: true,
        sanitize: false,
        highlight: function(code, lang) {
          if (lang && hljs.getLanguage(lang)) {
            try {
              return hljs.highlight(code, { language: lang }).value
            } catch (err) {
              console.warn('代码高亮失败:', err)
            }
          }
          try {
            return hljs.highlightAuto(code).value
          } catch (err) {
            return code
          }
        }
      })
      
      try {
        return marked(text)
      } catch (error) {
        console.error('Markdown 渲染错误:', error)
        return `<p>${text.replace(/\n/g, '<br>')}</p>`
      }
    }

    // 获取卡包名称
    const getDeckName = (deckId) => {
      const deck = decks.value.find(d => d.id === deckId)
      return deck ? deck.name : '未知卡包'
    }

    // 加载卡片
    const fetchCards = async () => {
      loading.value = true
      try {
        const params = {
          page: pagination.page,
          page_size: pagination.page_size,
          ...searchForm
        }
        
        // 清理空值
        Object.keys(params).forEach(key => {
          if (params[key] === '' || params[key] === null || params[key] === undefined) {
            delete params[key]
          }
        })

        const response = await searchCards(params)
        
        if (response.data && response.data.code === 'SUCCESS') {
          const data = response.data.data
          cardList.value = data.cards || []
          pagination.total = data.total || 0
          pagination.total_pages = data.total_pages || 0
        }
      } catch (error) {
        console.error('获取卡片列表失败:', error)
        ElMessage.error('获取卡片列表失败')
      } finally {
        loading.value = false
      }
    }

    // 加载卡包
    const fetchDecks = async () => {
      try {
        const response = await getDecks()
        if (response.data && response.data.code === 'SUCCESS') {
          decks.value = response.data.data.decks || []
        }
      } catch (error) {
        console.error('获取卡包列表失败:', error)
      }
    }

    // 加载标签
    const fetchTags = async () => {
      try {
        const response = await getAllTags()
        if (response.data && response.data.code === 'SUCCESS') {
          tags.value = response.data.data || []
        }
      } catch (error) {
        console.error('获取标签列表失败:', error)
      }
    }

    // 打开表单
    const openForm = (card = null) => {
      if (card) {
        Object.assign(form, card)
        questionMode.value = 'edit'
        answerMode.value = 'edit'
      } else {
        Object.assign(form, defaultCard)
        questionMode.value = 'edit'
        answerMode.value = 'edit'
      }
      showForm.value = true
    }

    // 关闭表单
    const closeForm = () => {
      showForm.value = false
      Object.assign(form, defaultCard)
    }

    // 提交表单
    const handleSubmit = async () => {
      try {
        if (form.id) {
          await updateCard(form.id, form.deck_id, form.tag_id || null, form.question, form.answer)
          ElMessage.success('卡片更新成功')
        } else {
          await createCard(form.deck_id, form.tag_id || null, form.question, form.answer)
          ElMessage.success('卡片创建成功')
        }
        closeForm()
        fetchCards()
      } catch (error) {
        console.error('保存卡片失败:', error)
        ElMessage.error('保存卡片失败')
      }
    }

    // 删除卡片
    const deleteCard = (card) => {
      deletingCard.value = card
    }

    // 确认删除
    const confirmDelete = async () => {
      if (!deletingCard.value && batchDeletingCards.value.length === 0) return
      
      isDeleting.value = true
      try {
        if (deletingCard.value) {
          await deleteCard(deletingCard.value.id)
          ElMessage.success('卡片删除成功')
        } else {
          const deletePromises = batchDeletingCards.value.map(id => deleteCard(id))
          await Promise.all(deletePromises)
          ElMessage.success(`成功删除 ${batchDeletingCards.value.length} 张卡片`)
        }
        fetchCards()
      } catch (error) {
        console.error('删除卡片失败:', error)
        ElMessage.error('删除卡片失败')
      } finally {
        deletingCard.value = null
        batchDeletingCards.value = []
        selectedCards.value = []
        isDeleting.value = false
      }
    }

    // 取消删除
    const cancelDelete = () => {
      deletingCard.value = null
      batchDeletingCards.value = []
    }

    // 切换批量模式
    const toggleBatchMode = () => {
      showBatchMode.value = !showBatchMode.value
      if (!showBatchMode.value) {
        selectedCards.value = []
      }
    }

    // 切换卡片选择
    const toggleCardSelection = (cardId) => {
      const index = selectedCards.value.indexOf(cardId)
      if (index > -1) {
        selectedCards.value.splice(index, 1)
      } else {
        selectedCards.value.push(cardId)
      }
    }

    // 全选/取消全选
    const toggleSelectAll = () => {
      if (selectedCards.value.length === cardList.value.length) {
        selectedCards.value = []
      } else {
        selectedCards.value = cardList.value.map(card => card.id)
      }
    }

    // 批量删除
    const batchDeleteCards = () => {
      if (selectedCards.value.length === 0) {
        ElMessage.warning('请先选择要删除的卡片')
        return
      }
      batchDeletingCards.value = [...selectedCards.value]
    }

    // 分页
    const changePage = (page) => {
      if (page < 1 || page > pagination.total_pages) return
      pagination.page = page
      fetchCards()
    }

    // 获取分页数字
    const getPageNumbers = () => {
      const pages = []
      const totalPages = pagination.total_pages
      const currentPage = pagination.page
      
      if (totalPages <= 7) {
        for (let i = 1; i <= totalPages; i++) {
          pages.push(i)
        }
      } else {
        if (currentPage <= 4) {
          for (let i = 1; i <= 5; i++) {
            pages.push(i)
          }
          pages.push('...')
          pages.push(totalPages)
        } else if (currentPage >= totalPages - 3) {
          pages.push(1)
          pages.push('...')
          for (let i = totalPages - 4; i <= totalPages; i++) {
            pages.push(i)
          }
        } else {
          pages.push(1)
          pages.push('...')
          for (let i = currentPage - 1; i <= currentPage + 1; i++) {
            pages.push(i)
          }
          pages.push('...')
          pages.push(totalPages)
        }
      }
      
      return pages
    }

    // 查看卡片详情
    const viewCardDetail = (card) => {
      selectedCardDetail.value = card
      showCardDetailModal.value = true
    }

    // 关闭卡片详情模态框
    const closeCardDetailModal = () => {
      showCardDetailModal.value = false
      selectedCardDetail.value = null
    }

    // 从详情页面编辑卡片
    const editCardFromDetail = () => {
      if (selectedCardDetail.value) {
        openForm(selectedCardDetail.value)
        closeCardDetailModal()
      }
    }

    // 格式化日期
    const formatDate = (dateString) => {
      if (!dateString) return ''
      const date = new Date(dateString)
      if (isNaN(date.getTime())) return ''
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
      })
    }

    onMounted(() => {
      fetchDecks()
      fetchTags()
      fetchCards()
    })

    return {
      loading,
      showForm,
      showBatchMode,
      selectedCards,
      cardList,
      decks,
      tags,
      deletingCard,
      batchDeletingCards,
      isDeleting,
      searchForm,
      pagination,
      form,
      questionMode,
      answerMode,
      showCardDetailModal,
      selectedCardDetail,
      renderMarkdown,
      getDeckName,
      fetchCards,
      openForm,
      closeForm,
      handleSubmit,
      deleteCard,
      confirmDelete,
      cancelDelete,
      toggleBatchMode,
      toggleCardSelection,
      toggleSelectAll,
      batchDeleteCards,
      changePage,
      getPageNumbers,
      viewCardDetail,
      closeCardDetailModal,
      editCardFromDetail,
      formatDate
    }
  }
}
</script>
