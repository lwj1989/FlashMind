<template>
  <div class="min-h-screen bg-gray-50 p-4">
    <div class="max-w-6xl mx-auto">
      <header class="mb-8">
        <div class="flex items-center mb-4">
          <router-link to="/decks" class="text-indigo-600 hover:text-indigo-800 flex items-center mr-4">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M9.707 16.707a1 1 0 01-1.414 0l-6-6a1 1 0 010-1.414l6-6a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l4.293 4.293a1 1 0 010 1.414z" clip-rule="evenodd" />
            </svg>
            返回卡包列表
          </router-link>
          <h1 v-if="deck" class="text-3xl font-bold text-gray-800">{{ deck.name }}</h1>
        </div>
      </header>

      <div v-if="loading" class="text-center py-8">
        <div class="inline-block animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-600"></div>
        <p class="mt-2 text-gray-600">加载中...</p>
      </div>

      <div v-else-if="deck" class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- 左侧：卡包信息和标签 -->
        <div class="lg:col-span-1 space-y-6">
          <!-- 卡包信息卡片 -->
          <div class="bg-white rounded-xl shadow p-6">
            <h2 class="text-xl font-semibold text-gray-800 mb-4">卡包信息</h2>
            <div class="space-y-3">
              <div class="flex justify-between">
                <span class="text-gray-600">状态</span>
                <span v-if="deck.archived" class="bg-gray-100 text-gray-800 text-xs px-2 py-1 rounded">已归档</span>
                <span v-else class="bg-green-100 text-green-800 text-xs px-2 py-1 rounded">使用中</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">创建时间</span>
                <span class="text-gray-800">{{ formatDate(deck.created_at) }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">更新时间</span>
                <span class="text-gray-800">{{ formatDate(deck.updated_at) }}</span>
              </div>
            </div>
            
            <div class="mt-6 space-y-3">
              <button @click="startStudy" class="w-full bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700 transition-colors">
                开始学习此卡包
              </button>
              <div class="flex space-x-3">
                <button @click="editDeck" class="flex-1 bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700 transition-colors">
                  编辑卡包
                </button>
                <button @click="confirmDelete" class="flex-1 bg-red-600 text-white px-4 py-2 rounded-lg hover:bg-red-700 transition-colors">
                  删除卡包
                </button>
              </div>
            </div>
          </div>

          <!-- 标签卡片 -->
          <div class="bg-white rounded-xl shadow p-6">
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-xl font-semibold text-gray-800">标签</h2>
              <button @click="showCreateTagModal = true" class="text-indigo-600 hover:text-indigo-800">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clip-rule="evenodd" />
                </svg>
              </button>
            </div>
            
            <div v-if="tagsLoading" class="text-center py-4">
              <div class="inline-block animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-indigo-600"></div>
            </div>
            
            <div v-else-if="tags.length === 0" class="text-center py-4 text-gray-500">
              暂无标签
            </div>
            
            <div v-else class="space-y-3">
              <div v-for="tag in tags" :key="tag.id" class="flex justify-between items-center group">
                <span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-indigo-100 text-indigo-800 hover:bg-indigo-200 transition-colors">
                  {{ tag.name }}
                </span>
                <div class="flex space-x-1 opacity-0 group-hover:opacity-100 transition-opacity">
                  <button @click="editTag(tag)" class="text-gray-400 hover:text-gray-600 p-1 rounded hover:bg-gray-100" title="编辑标签">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                      <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
                    </svg>
                  </button>
                  <button @click="confirmDeleteTag(tag)" class="text-red-400 hover:text-red-600 p-1 rounded hover:bg-red-50" title="删除标签">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                      <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧：卡片列表 -->
        <div class="lg:col-span-2">
          <div class="bg-white rounded-xl shadow p-6">
            <div class="flex justify-between items-center mb-6">
              <h2 class="text-xl font-semibold text-gray-800">卡片 ({{ filteredCards.length }})</h2>
              <button @click="showCreateCardModal = true" class="bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700 transition-colors flex items-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clip-rule="evenodd" />
                </svg>
                添加卡片
              </button>
            </div>

            <!-- 搜索框 -->
            <div class="mb-6">
              <div class="relative">
                <input 
                  v-model="searchQuery" 
                  type="text" 
                  placeholder="搜索卡片问题或答案..." 
                  class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
                >
                <div class="absolute left-3 top-2.5 text-gray-400">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd" />
                  </svg>
                </div>
              </div>
            </div>

            <!-- 卡片列表 -->
            <div class="text-center py-8 text-gray-500" v-if="filteredCards.length === 0">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto text-gray-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <p>{{ searchQuery ? '没有找到匹配的卡片' : '暂无卡片' }}</p>
              <button v-if="!searchQuery" @click="showCreateCardModal = true" class="mt-2 text-indigo-600 hover:text-indigo-800">
                添加第一张卡片
              </button>
            </div>
            
            <div v-else class="space-y-4">
              <div 
                v-for="card in filteredCards" 
                :key="card.id" 
                class="bg-gray-50 rounded-lg p-4 border border-gray-200 hover:shadow-md transition-shadow cursor-pointer"
                @click="viewCard(card)"
              >
                <div class="flex justify-between items-start">
                  <div class="flex-1 mr-4">
                    <div class="font-medium text-gray-900 mb-2 line-clamp-2" v-html="renderMarkdownPreview(card.front)"></div>
                    <div class="text-sm text-gray-600 line-clamp-3" v-html="renderMarkdownPreview(card.back)"></div>
                    
                    <!-- 标签显示 -->
                    <div v-if="card.tags && card.tags.length > 0" class="mt-3">
                      <div class="flex flex-wrap gap-1">
                        <span v-for="tag in card.tags" :key="tag.id" class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800">
                          {{ tag.name }}
                        </span>
                      </div>
                    </div>

                    <div class="mt-3 text-xs text-gray-500">
                      创建于 {{ formatDate(card.created_at) }}
                    </div>
                  </div>
                  
                  <!-- 操作按钮 -->
                  <div class="flex space-x-2" @click.stop>
                    <button 
                      @click="editCard(card)" 
                      class="text-gray-600 hover:text-gray-800 p-2 rounded hover:bg-gray-200 transition-colors"
                      title="编辑卡片"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                        <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
                      </svg>
                    </button>
                    <button 
                      @click="confirmDeleteCard(card)" 
                      class="text-red-600 hover:text-red-800 p-2 rounded hover:bg-red-50 transition-colors"
                      title="删除卡片"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
                      </svg>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 编辑卡包模态框 -->
    <div v-if="showEditDeckModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
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
                @click="showEditDeckModal = false" 
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

    <!-- 删除卡包确认模态框 -->
    <div v-if="showDeleteDeckModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-xl shadow-lg w-full max-w-md">
        <div class="p-6">
          <h3 class="text-xl font-semibold text-gray-800 mb-4">确认删除</h3>
          <p class="text-gray-600 mb-6">
            确定要删除卡包 "{{ deck.name }}" 吗？此操作不可撤销。
          </p>
          <div class="flex justify-end space-x-3">
            <button 
              type="button" 
              @click="showDeleteDeckModal = false" 
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

    <!-- 创建标签模态框 -->
    <div v-if="showCreateTagModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
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
                @click="showCreateTagModal = false" 
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
    <div v-if="showEditTagModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
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
                @click="showEditTagModal = false" 
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

    <!-- 删除标签确认模态框 -->
    <div v-if="showDeleteTagModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-xl shadow-lg w-full max-w-md">
        <div class="p-6">
          <h3 class="text-xl font-semibold text-gray-800 mb-4">确认删除</h3>
          <p class="text-gray-600 mb-6">
            确定要删除标签 "{{ deletingTag.name }}" 吗？此操作不可撤销。
          </p>
          <div class="flex justify-end space-x-3">
            <button 
              type="button" 
              @click="showDeleteTagModal = false" 
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
    
    <!-- 创建卡片模态框 -->
    <div v-if="showCreateCardModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md">
        <div class="p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">创建新卡片</h3>
          <form @submit.prevent="createCard">
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">正面</label>
              <textarea v-model="newCard.front" class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500" rows="3" required></textarea>
            </div>
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">背面</label>
              <textarea v-model="newCard.back" class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500" rows="3" required></textarea>
            </div>
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">标签</label>
              <select v-model="newCard.tag_ids" multiple class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500">
                <option v-for="tag in tags" :key="tag.id" :value="tag.id">{{ tag.name }}</option>
              </select>
              <p class="mt-1 text-xs text-gray-500">按住Ctrl/Cmd键可选择多个标签</p>
            </div>
            <div class="flex justify-end space-x-3">
              <button type="button" @click="showCreateCardModal = false" class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50">
                取消
              </button>
              <button type="submit" class="px-4 py-2 bg-indigo-600 text-white rounded-md text-sm font-medium hover:bg-indigo-700">
                创建
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
    
    <!-- 编辑卡片模态框 -->
    <div v-if="showEditCardModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md">
        <div class="p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">编辑卡片</h3>
          <form @submit.prevent="updateCard">
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">正面</label>
              <textarea v-model="editingCard.front" class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500" rows="3" required></textarea>
            </div>
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">背面</label>
              <textarea v-model="editingCard.back" class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500" rows="3" required></textarea>
            </div>
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">标签</label>
              <select v-model="editingCard.tag_ids" multiple class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500">
                <option v-for="tag in tags" :key="tag.id" :value="tag.id">{{ tag.name }}</option>
              </select>
              <p class="mt-1 text-xs text-gray-500">按住Ctrl/Cmd键可选择多个标签</p>
            </div>
            <div class="flex justify-end space-x-3">
              <button type="button" @click="showEditCardModal = false" class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50">
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
    
    <!-- 删除卡片确认模态框 -->
    <div v-if="showDeleteCardModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md">
        <div class="p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">确认删除</h3>
          <p class="text-gray-700 mb-6">确定要删除这张卡片吗？此操作无法撤销。</p>
          <div class="flex justify-end space-x-3">
            <button @click="showDeleteCardModal = false" class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50">
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
import api from '../services/api'
import { ElMessage } from 'element-plus'
import CardItem from '../components/CardItem.vue'

export default {
  name: 'DeckDetailView',
  components: {
    CardItem
  },
  data() {
    return {
      deck: null,
      tags: [],
      cards: [],
      loading: false,
      tagsLoading: false,
      searchQuery: '',
      showEditDeckModal: false,
      showDeleteDeckModal: false,
      showCreateTagModal: false,
      showEditTagModal: false,
      showDeleteTagModal: false,
      showCreateCardModal: false,
      showEditCardModal: false,
      showDeleteCardModal: false,
      editingDeck: {},
      newTag: {
        name: '',
        deck_id: null
      },
      editingTag: {},
      deletingTag: {},
      newCard: {
        front: '',
        back: '',
        deck_id: null,
        tag_ids: []
      },
      editingCard: {
        id: null,
        front: '',
        back: '',
        deck_id: null,
        tag_ids: []
      },
      deletingCard: null
    }
  },
  computed: {
    filteredCards() {
      if (!this.searchQuery) {
        return this.cards
      }
      const query = this.searchQuery.toLowerCase()
      return this.cards.filter(card => 
        card.front.toLowerCase().includes(query) || 
        card.back.toLowerCase().includes(query) ||
        (card.tags && card.tags.some(tag => tag.name.toLowerCase().includes(query)))
      )
    }
  },
  created() {
    const deckId = this.$route.params.id
    this.newCard.deck_id = parseInt(deckId)
    this.editingCard.deck_id = parseInt(deckId)
    this.newTag.deck_id = parseInt(deckId)
    this.fetchDeck()
    this.fetchTags()
    this.fetchCards()
  },
  methods: {
    async fetchDeck() {
      this.loading = true
      try {
        const deckId = this.$route.params.id
        const response = await api.getDeck(deckId)
        this.deck = response.data?.data || null
      } catch (error) {
        console.error('获取卡包详情失败:', error)
        ElMessage.error('获取卡包详情失败，请稍后重试')
      } finally {
        this.loading = false
      }
    },
    async fetchTags() {
      this.tagsLoading = true
      try {
        const deckId = this.$route.params.id
        const response = await api.getTags(deckId)
        this.tags = response.data?.data || []
      } catch (error) {
        console.error('获取标签失败:', error)
        ElMessage.error('获取标签失败，请稍后重试')
      } finally {
        this.tagsLoading = false
      }
    },
    async fetchCards() {
      try {
        const deckId = this.$route.params.id
        if (!deckId) {
          console.error('卡包ID为空')
          return
        }
        console.log('获取卡包下的卡片，deckId:', deckId)
        const response = await api.getCards(deckId)
        console.log('获取卡片响应:', response)
        this.cards = response.data?.data || []
      } catch (error) {
        console.error('获取卡片列表失败:', error)
        ElMessage.error('获取卡片列表失败，请稍后重试')
      }
    },
    editDeck() {
      if (!this.deck) {
        ElMessage.error('卡包信息不存在，请刷新页面重试')
        return
      }
      this.editingDeck = { ...this.deck }
      this.showEditDeckModal = true
    },
    async updateDeck() {
      try {
        if (!this.editingDeck || !this.editingDeck.id) {
          ElMessage.error('卡包信息不完整，请刷新页面重试')
          return
        }
        await api.updateDeck(this.editingDeck.id, {
          name: this.editingDeck.name,
          archived: this.editingDeck.archived
        })
        this.showEditDeckModal = false
        this.fetchDeck()
        ElMessage.success('卡包更新成功！')
      } catch (error) {
        console.error('更新卡包失败:', error)
        ElMessage.error('更新卡包失败，请稍后重试')
      }
    },
    confirmDelete() {
      this.showDeleteDeckModal = true
    },
    async deleteDeck() {
      try {
        if (!this.deck || !this.deck.id) {
          ElMessage.error('卡包信息不存在，请刷新页面重试')
          return
        }
        await api.deleteDeck(this.deck.id)
        ElMessage.success('卡包删除成功！')
        this.$router.push('/decks')
      } catch (error) {
        console.error('删除卡包失败:', error)
        ElMessage.error('删除卡包失败，请稍后重试')
      }
    },
    async createTag() {
      try {
        this.newTag.deck_id = parseInt(this.$route.params.id)
        await api.createTag(this.newTag)
        this.showCreateTagModal = false
        this.newTag = { name: '', deck_id: null }
        ElMessage.success('标签创建成功！')
        this.fetchTags()
      } catch (error) {
        console.error('创建标签失败:', error)
        ElMessage.error('创建标签失败，请稍后重试')
      }
    },
    editTag(tag) {
      this.editingTag = { ...tag }
      this.showEditTagModal = true
    },
    async updateTag() {
      try {
        if (!this.editingTag || !this.editingTag.id) {
          ElMessage.error('标签信息不完整，请刷新页面重试')
          return
        }
        await api.updateTag(this.editingTag.id, {
          name: this.editingTag.name
        })
        this.showEditTagModal = false
        this.fetchTags()
        ElMessage.success('标签更新成功！')
      } catch (error) {
        console.error('更新标签失败:', error)
        ElMessage.error('更新标签失败，请稍后重试')
      }
    },
    confirmDeleteTag(tag) {
      this.deletingTag = tag
      this.showDeleteTagModal = true
    },
    async deleteTag() {
      try {
        if (!this.deletingTag || !this.deletingTag.id) {
          ElMessage.error('标签信息不存在，请刷新页面重试')
          return
        }
        await api.deleteTag(this.deletingTag.id)
        this.showDeleteTagModal = false
        ElMessage.success('标签删除成功！')
        this.fetchTags()
      } catch (error) {
        console.error('删除标签失败:', error)
        ElMessage.error('删除标签失败，请稍后重试')
      }
    },
    async createCard() {
      try {
        await api.createCard(this.newCard)
        this.newCard = { front: '', back: '', deck_id: this.deck.id, tag_ids: [] }
        this.showCreateCardModal = false
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
        if (!this.editingCard || !this.editingCard.id) {
          ElMessage.error('卡片信息不完整，请刷新页面重试')
          return
        }
        await api.updateCard(this.editingCard.id, this.editingCard)
        this.showEditCardModal = false
        this.fetchCards()
        ElMessage.success('卡片更新成功！')
      } catch (error) {
        console.error('更新卡片失败:', error)
        ElMessage.error('更新卡片失败，请稍后重试')
      }
    },
    confirmDeleteCard(card) {
      this.deletingCard = card
      this.showDeleteCardModal = true
    },
    async deleteCard() {
      try {
        if (!this.deletingCard || !this.deletingCard.id) {
          ElMessage.error('卡片信息不存在，请刷新页面重试')
          return
        }
        await api.deleteCard(this.deletingCard.id)
        this.showDeleteCardModal = false
        ElMessage.success('卡片删除成功！')
        this.fetchCards()
      } catch (error) {
        console.error('删除卡片失败:', error)
        ElMessage.error('删除卡片失败，请稍后重试')
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
    startStudy() {
      // 跳转到学习页面，并传递卡包ID
      this.$router.push({
        path: '/study',
        query: { deckId: this.deck.id }
      })
    },
    renderMarkdownPreview(text) {
      if (!text) return ''
      // 简单的markdown预览，移除markdown语法显示纯文本
      return text
        .replace(/^#+\s+/gm, '') // 移除标题标记
        .replace(/\*\*(.*?)\*\*/g, '$1') // 移除粗体标记
        .replace(/\*(.*?)\*/g, '$1') // 移除斜体标记
        .replace(/`(.*?)`/g, '$1') // 移除行内代码标记
        .replace(/```[\s\S]*?```/g, '[代码块]') // 代码块用占位符
        .replace(/!\[.*?\]\(.*?\)/g, '[图片]') // 图片用占位符
        .replace(/\[.*?\]\(.*?\)/g, '[链接]') // 链接用占位符
    },
    viewCard(card) {
      // 触发查看卡片的弹窗
      this.$emit('view-card', card)
    }
  }
}
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>