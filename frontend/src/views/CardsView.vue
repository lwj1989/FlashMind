<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-6xl mx-auto px-4 py-12">
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

      <!-- 批量操作栏 -->
      <div v-if="showBatchMode && selectedCards.length > 0" class="bg-white rounded-xl shadow mb-4">
        <div class="bg-indigo-50 px-6 py-3 border-l-4 border-indigo-400">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-indigo-800">已选择 {{ selectedCards.length }} 张卡片</span>
            <div class="flex space-x-3">
              <button 
                @click="showBatchMoveModal = true" 
                class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors text-sm"
              >
                批量移动
              </button>
              <button 
                @click="showBatchTagModal = true" 
                class="bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700 transition-colors text-sm"
              >
                批量标签
              </button>
              <button 
                @click="batchDeleteCards" 
                class="bg-red-600 text-white px-4 py-2 rounded-lg hover:bg-red-700 transition-colors text-sm"
              >
                批量删除
              </button>
            </div>
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
                  <button v-if="card.tag_name" 
                    @click="goToTagCards(card.tag_id)" 
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800 hover:bg-indigo-200 transition-colors cursor-pointer"
                    :title="'点击查看标签「' + card.tag_name + '」的所有卡片'"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3 mr-1" viewBox="0 0 20 20" fill="currentColor">
                      <path fill-rule="evenodd" d="M17.707 9.293a1 1 0 010 1.414l-7 7a1 1 0 01-1.414 0l-7-7A.997.997 0 012 10V5a3 3 0 013-3h5c.256 0 .512.098.707.293l7 7zM5 6a1 1 0 100-2 1 1 0 000 2z" clip-rule="evenodd" />
                    </svg>
                    {{ card.tag_name }}
                  </button>
                  <span v-else class="text-gray-400 text-sm">无标签</span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <button 
                    @click="openForm(card)" 
                    class="text-indigo-600 hover:text-emerald-900 mr-3 p-1"
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
                <div class="markdown-content text-blue-900 mb-4" v-html="renderMarkdown(selectedCardDetail.question)"></div>
                
                <!-- 标签显示在问题下方 -->
                <div v-if="selectedCardDetail.tag_name" class="mt-4">
                  <button 
                    @click="goToTagCards(selectedCardDetail.tag_id)" 
                    class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-200 text-blue-800 hover:bg-blue-300 transition-colors cursor-pointer"
                    :title="'点击查看标签「' + selectedCardDetail.tag_name + '」的所有卡片'"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" viewBox="0 0 20 20" fill="currentColor">
                      <path fill-rule="evenodd" d="M17.707 9.293a1 1 0 010 1.414l-7 7a1 1 0 01-1.414 0l-7-7A.997.997 0 012 10V5a3 3 0 013-3h5c.256 0 .512.098.707.293l7 7zM5 6a1 1 0 100-2 1 1 0 000 2z" clip-rule="evenodd" />
                    </svg>
                    {{ selectedCardDetail.tag_name }}
                  </button>
                </div>
              </div>

              <!-- 答案部分 -->
              <div class="bg-green-50 rounded-lg p-6">
                <h4 class="text-lg font-semibold text-green-800 mb-4">答案</h4>
                <div class="markdown-content text-green-900" v-html="renderMarkdown(selectedCardDetail.answer)"></div>
              </div>
            </div>

            <!-- 卡片信息 -->
            <div class="mt-6 bg-gray-50 rounded-lg p-4">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                <div>
                  <span class="text-gray-600">所属卡包：</span>
                  <span class="font-medium">{{ getDeckName(selectedCardDetail.deck_id) }}</span>
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

      <!-- 批量移动卡片模态框 -->
      <div v-if="showBatchMoveModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
        <div class="bg-white rounded-xl shadow-lg w-full max-w-md">
          <div class="p-6">
            <h3 class="text-xl font-semibold text-gray-800 mb-4">批量移动卡片</h3>
            <p class="text-gray-600 mb-4">将选中的 {{ selectedCards.length }} 张卡片移动到指定卡包：</p>
            
            <form @submit.prevent="confirmBatchMove">
              <div class="mb-4">
                <label class="block text-sm font-medium text-gray-700 mb-2">目标卡包</label>
                <select 
                  v-model="batchMoveTargetDeck" 
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">请选择目标卡包</option>
                  <option v-for="deck in decks" :key="deck.id" :value="deck.id">{{ deck.name }}</option>
                </select>
              </div>
              
              <div class="flex justify-end space-x-3">
                <button 
                  type="button" 
                  @click="closeBatchMoveModal"
                  class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
                >
                  取消
                </button>
                <button 
                  type="submit"
                  :disabled="batchMoving"
                  class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
                >
                  {{ batchMoving ? '移动中...' : '确认移动' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>

      <!-- 批量修改标签模态框 -->
      <div v-if="showBatchTagModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
        <div class="bg-white rounded-xl shadow-lg w-full max-w-md">
          <div class="p-6">
            <h3 class="text-xl font-semibold text-gray-800 mb-4">批量修改标签</h3>
            <p class="text-gray-600 mb-4">为选中的 {{ selectedCards.length }} 张卡片设置标签：</p>
            
            <form @submit.prevent="confirmBatchTag">
              <div class="mb-4">
                <label class="block text-sm font-medium text-gray-700 mb-2">操作类型</label>
                <div class="space-y-2">
                  <label class="flex items-center">
                    <input 
                      type="radio" 
                      v-model="batchTagAction" 
                      value="set"
                      class="mr-2"
                    />
                    <span>设置标签（覆盖原有标签）</span>
                  </label>
                  <label class="flex items-center">
                    <input 
                      type="radio" 
                      v-model="batchTagAction" 
                      value="add"
                      class="mr-2"
                    />
                    <span>添加标签（保留原有标签）</span>
                  </label>
                  <label class="flex items-center">
                    <input 
                      type="radio" 
                      v-model="batchTagAction" 
                      value="remove"
                      class="mr-2"
                    />
                    <span>移除标签</span>
                  </label>
                </div>
              </div>
              
              <div class="mb-4">
                <label class="block text-sm font-medium text-gray-700 mb-2">选择标签</label>
                <select 
                  v-model="batchTagTarget" 
                  :required="batchTagAction !== 'clear'"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                >
                  <option value="">请选择标签</option>
                  <option v-for="tag in tags" :key="tag.id" :value="tag.id">{{ tag.name }}</option>
                </select>
              </div>
              
              <div class="flex justify-end space-x-3">
                <button 
                  type="button" 
                  @click="closeBatchTagModal"
                  class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
                >
                  取消
                </button>
                <button 
                  type="submit"
                  :disabled="batchTagging"
                  class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-50"
                >
                  {{ batchTagging ? '处理中...' : '确认操作' }}
                </button>
              </div>
            </form>
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

.markdown-content :deep(h1),
.markdown-content :deep(h2),
.markdown-content :deep(h3),
.markdown-content :deep(h4),
.markdown-content :deep(h5),
.markdown-content :deep(h6) {
  font-weight: 600;
  margin-top: 1.5rem;
  margin-bottom: 0.75rem;
}

.markdown-content :deep(h1) { font-size: 1.75rem; }
.markdown-content :deep(h2) { font-size: 1.5rem; }
.markdown-content :deep(h3) { font-size: 1.25rem; }
.markdown-content :deep(h4) { font-size: 1.125rem; }

.markdown-content :deep(p) {
  margin-bottom: 1rem;
}

.markdown-content :deep(ul),
.markdown-content :deep(ol) {
  margin-bottom: 1rem;
  padding-left: 1.5rem;
}

.markdown-content :deep(li) {
  margin-bottom: 0.25rem;
}

.markdown-content :deep(blockquote) {
  border-left: 4px solid #e5e7eb;
  padding-left: 1rem;
  margin: 1rem 0;
  color: #6b7280;
  font-style: italic;
}

.markdown-content :deep(pre) {
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 8px;
  margin: 16px 0;
  overflow: hidden;
  font-family: 'Fira Code', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.5;
  color: #e2e8f0;
  position: relative;
}

.markdown-content :deep(pre .code-header) {
  background-color: #0f172a;
  padding: 8px 16px;
  border-bottom: 1px solid #334155;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
}

.markdown-content :deep(pre .code-lang) {
  color: #64748b;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.markdown-content :deep(pre .copy-btn) {
  background: none;
  border: none;
  color: #64748b;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  transition: all 0.2s;
}

.markdown-content :deep(pre .copy-btn:hover) {
  color: #e2e8f0;
  background-color: #334155;
}

.markdown-content :deep(pre .code-content) {
  padding: 16px;
  overflow-x: auto;
}

.markdown-content :deep(code) {
  background-color: #f1f5f9;
  color: #e11d48;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Fira Code', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
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

.markdown-content :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 1rem 0;
}

.markdown-content :deep(th),
.markdown-content :deep(td) {
  border: 1px solid #e5e7eb;
  padding: 0.5rem;
  text-align: left;
}

.markdown-content :deep(th) {
  background-color: #f9fafb;
  font-weight: 600;
}

.markdown-content :deep(a) {
  color: #3b82f6;
  text-decoration: underline;
}

.markdown-content :deep(a:hover) {
  color: #1d4ed8;
}
</style>

<script>
import { ref, reactive, onMounted, onUnmounted, getCurrentInstance } from 'vue'
import { ElMessage } from 'element-plus'
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import { searchCards, createCard, updateCard, deleteCard as deleteCardAPI } from '@/api/card'
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
    
    // 批量操作相关
    const showBatchMoveModal = ref(false)
    const showBatchTagModal = ref(false)
    const batchMoveTargetDeck = ref('')
    const batchMoving = ref(false)
    const batchTagAction = ref('set')
    const batchTagTarget = ref('')
    const batchTagging = ref(false)

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

    // 代码格式化函数
    const formatCode = (code, lang) => {
      // 基本的代码格式化逻辑
      let formatted = code
      
      // 移除首尾空白行
      formatted = formatted.replace(/^\s*\n+|\n+\s*$/g, '')
      
      // 统一缩进 - 找到最小缩进并规范化
      const lines = formatted.split('\n')
      const nonEmptyLines = lines.filter(line => line.trim().length > 0)
      
      if (nonEmptyLines.length > 0) {
        // 找到最小缩进
        const minIndent = Math.min(...nonEmptyLines.map(line => {
          const match = line.match(/^(\s*)/)
          return match ? match[1].length : 0
        }))
        
        // 移除多余的缩进
        const normalizedLines = lines.map(line => {
          if (line.trim().length === 0) return ''
          return line.substring(minIndent)
        })
        
        formatted = normalizedLines.join('\n')
      }
      
      return formatted
    }

    // 渲染 Markdown
    const renderMarkdown = (text) => {
      if (!text) return ''
      
      marked.setOptions({
        breaks: true,
        gfm: true,
        sanitize: false,
        highlight: function(code, lang) {
          // 先格式化代码
          const formattedCode = formatCode(code, lang)
          
          // 尝试语言检测和高亮
          if (lang && hljs.getLanguage(lang)) {
            try {
              const result = hljs.highlight(formattedCode, { language: lang })
              return `<div class="code-header">
                        <span class="code-lang">${lang.toUpperCase()}</span>
                        <button class="copy-btn" onclick="copyCode(this)" title="复制代码">
                          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                            <path d="m5 15-6-6v-6h6l6 6"></path>
                          </svg>
                        </button>
                      </div>
                      <div class="code-content">${result.value}</div>`
            } catch (err) {
              console.warn('代码高亮失败:', err)
            }
          }
          
          // 自动语言检测
          try {
            const result = hljs.highlightAuto(formattedCode)
            const detectedLang = result.language || 'text'
            return `<div class="code-header">
                      <span class="code-lang">${detectedLang.toUpperCase()}</span>
                      <button class="copy-btn" onclick="copyCode(this)" title="复制代码">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                          <path d="m5 15-6-6v-6h6l6 6"></path>
                        </svg>
                      </button>
                    </div>
                    <div class="code-content">${result.value}</div>`
          } catch (err) {
            return `<div class="code-header">
                      <span class="code-lang">TEXT</span>
                      <button class="copy-btn" onclick="copyCode(this)" title="复制代码">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                          <path d="m5 15-6-6v-6h6l6 6"></path>
                        </svg>
                      </button>
                    </div>
                    <div class="code-content">${formattedCode}</div>`
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
        const cardData = {
          deck_id: parseInt(form.deck_id),
          question: form.question,
          answer: form.answer
        }
        
        // 只有当tag_id存在且不为空时才添加
        if (form.tag_id && form.tag_id !== '') {
          cardData.tag_id = parseInt(form.tag_id)
        } else {
          cardData.tag_id = null
        }
        
        console.log('发送卡片数据:', cardData) // 调试信息
        
        if (form.id) {
          await updateCard(form.id, cardData)
          ElMessage.success('卡片更新成功')
        } else {
          await createCard(cardData)
          ElMessage.success('卡片创建成功')
        }
        closeForm()
        fetchCards()
      } catch (error) {
        console.error('保存卡片失败:', error)
        console.error('错误详情:', error.response?.data) // 更详细的错误信息
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
          await deleteCardAPI(deletingCard.value.id)
          ElMessage.success('卡片删除成功')
        } else {
          const deletePromises = batchDeletingCards.value.map(id => deleteCardAPI(id))
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

    // 跳转到标签卡片页面
    const goToTagCards = (tagId) => {
      if (tagId) {
        closeCardDetailModal() // 关闭当前模态框
        const route = getCurrentInstance().proxy.$route
        if (route.query.tag_id !== tagId.toString()) {
          // 更新查询参数并重新获取卡片
          searchForm.tag_id = tagId
          fetchCards()
        }
      }
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

    // 批量移动卡片相关方法
    const closeBatchMoveModal = () => {
      showBatchMoveModal.value = false
      batchMoveTargetDeck.value = ''
    }

    const confirmBatchMove = async () => {
      if (!batchMoveTargetDeck.value) {
        ElMessage.warning('请选择目标卡包')
        return
      }

      if (selectedCards.value.length === 0) {
        ElMessage.warning('请先选择要移动的卡片')
        return
      }

      batchMoving.value = true
      try {
        // 使用单个更新API进行批量移动
        const promises = selectedCards.value.map(cardId => {
          const card = cardList.value.find(c => c.id === cardId)
          if (!card) return Promise.resolve()
          
          // 如果卡片已经在目标卡包中，跳过
          if (card.deck_id && card.deck_id.toString() === batchMoveTargetDeck.value.toString()) {
            return Promise.resolve()
          }
          
          // 构建更新数据，只修改卡包ID，保持其他信息不变
          const updateData = {
            deck_id: batchMoveTargetDeck.value,
            tag_id: card.tag_id || null,
            question: card.question,
            answer: card.answer
          }
          
          return updateCard(cardId, updateData)
        })
        
        await Promise.all(promises)
        ElMessage.success(`成功移动 ${selectedCards.value.length} 张卡片`)
        
        // 清理状态
        selectedCards.value = []
        showBatchMode.value = false
        closeBatchMoveModal()
        fetchCards()
      } catch (error) {
        console.error('批量移动卡片失败:', error)
        ElMessage.error('批量移动卡片失败')
      } finally {
        batchMoving.value = false
      }
    }

    // 批量修改标签相关方法
    const closeBatchTagModal = () => {
      showBatchTagModal.value = false
      batchTagAction.value = 'set'
      batchTagTarget.value = ''
    }

    const confirmBatchTag = async () => {
      if (!batchTagTarget.value && batchTagAction.value !== 'clear') {
        ElMessage.warning('请选择标签')
        return
      }

      if (selectedCards.value.length === 0) {
        ElMessage.warning('请先选择要修改标签的卡片')
        return
      }

      batchTagging.value = true
      try {
        // TODO: 这里需要调用批量修改标签的API
        // 根据不同的操作类型进行处理
        let actionText = ''
        switch (batchTagAction.value) {
          case 'set':
            actionText = '设置标签'
            break
          case 'add':
            actionText = '添加标签'
            break
          case 'remove':
            actionText = '移除标签'
            break
        }

        // 使用单个更新API进行批量标签操作
        const promises = selectedCards.value.map(cardId => {
          const card = cardList.value.find(c => c.id === cardId)
          if (!card) return Promise.resolve()
          
          let newTagId = null
          if (batchTagAction.value === 'set') {
            newTagId = batchTagTarget.value
          } else if (batchTagAction.value === 'add') {
            // 添加标签逻辑（如果卡片没有标签，设置新标签；如果已有标签，保持不变）
            newTagId = card.tag_id || batchTagTarget.value
          } else if (batchTagAction.value === 'remove') {
            // 移除标签逻辑：只有当卡片的标签是目标标签时才移除
            if (card.tag_id && card.tag_id.toString() === batchTagTarget.value.toString()) {
              newTagId = null
            } else {
              // 如果不是目标标签，保持原标签不变
              return Promise.resolve()
            }
          }
          
          const updateData = {
            deck_id: card.deck_id,
            tag_id: newTagId,
            question: card.question,
            answer: card.answer
          }
          
          return updateCard(cardId, updateData)
        })
        
        await Promise.all(promises)
        ElMessage.success(`成功为 ${selectedCards.value.length} 张卡片${actionText}`)
        
        // 清理状态
        selectedCards.value = []
        showBatchMode.value = false
        closeBatchTagModal()
        fetchCards()
      } catch (error) {
        console.error('批量修改标签失败:', error)
        ElMessage.error('批量修改标签失败')
      } finally {
        batchTagging.value = false
      }
    }

    // ESC 键关闭模态框
    const handleEscapeKey = (event) => {
      if (event.key === 'Escape') {
        // 按优先级关闭模态框
        if (showForm.value) {
          closeForm()
        } else if (showCardDetailModal.value) {
          closeCardDetailModal()
        } else if (showBatchMoveModal.value) {
          closeBatchMoveModal()
        } else if (showBatchTagModal.value) {
          closeBatchTagModal()
        } else if (deletingCard.value || batchDeletingCards.value.length > 0) {
          cancelDelete()
        }
      }
    }

    onMounted(() => {
      // 添加键盘事件监听器
      document.addEventListener('keydown', handleEscapeKey)
      
      // 添加复制代码功能
      window.copyCode = function(button) {
        const codeContent = button.parentElement.nextElementSibling
        const code = codeContent.textContent || codeContent.innerText
        
        if (navigator.clipboard) {
          navigator.clipboard.writeText(code).then(() => {
            // 更新按钮图标为已复制状态
            const originalHTML = button.innerHTML
            button.innerHTML = `<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                  <polyline points="20,6 9,17 4,12"></polyline>
                                </svg>`
            button.style.color = '#22c55e'
            
            setTimeout(() => {
              button.innerHTML = originalHTML
              button.style.color = ''
            }, 2000)
          }).catch(err => {
            console.error('复制失败:', err)
          })
        } else {
          // 降级方案
          const textArea = document.createElement('textarea')
          textArea.value = code
          document.body.appendChild(textArea)
          textArea.select()
          try {
            document.execCommand('copy')
            const originalHTML = button.innerHTML
            button.innerHTML = `<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                  <polyline points="20,6 9,17 4,12"></polyline>
                                </svg>`
            button.style.color = '#22c55e'
            
            setTimeout(() => {
              button.innerHTML = originalHTML
              button.style.color = ''
            }, 2000)
          } catch (err) {
            console.error('复制失败:', err)
          }
          document.body.removeChild(textArea)
        }
      }
      
      fetchDecks()
      fetchTags()
      
      // 检查路由查询参数，如果有tag_id或deck_id参数，自动搜索
      const route = getCurrentInstance().proxy.$route
      if (route.query.tag_id) {
        searchForm.tag_id = route.query.tag_id
      }
      if (route.query.deck_id) {
        searchForm.deck_id = route.query.deck_id
      }
      
      fetchCards()
    })

    onUnmounted(() => {
      // 移除键盘事件监听器
      document.removeEventListener('keydown', handleEscapeKey)
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
      showBatchMoveModal,
      showBatchTagModal,
      batchMoveTargetDeck,
      batchMoving,
      batchTagAction,
      batchTagTarget,
      batchTagging,
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
      goToTagCards,
      editCardFromDetail,
      formatDate,
      closeBatchMoveModal,
      confirmBatchMove,
      closeBatchTagModal,
      confirmBatchTag
    }
  }
}
</script>
