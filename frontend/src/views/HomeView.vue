<template>
  <div class="min-h-screen bg-gradient-to-br from-indigo-50 to-purple-50">
    <div class="max-w-6xl mx-auto px-4 py-12">
      <!-- 标题部分 -->
      <div class="text-center mb-16">
        <h1 class="text-4xl md:text-5xl font-bold text-gray-800 mb-4">欢迎使用 FlashMind</h1>
        <p class="text-xl text-gray-600 max-w-2xl mx-auto">智能卡片学习系统，帮助您高效学习和记忆知识</p>
      </div>

      <!-- 功能卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8 mb-16">
        <div class="bg-white rounded-xl shadow-lg p-6 hover:shadow-xl transition-shadow">
          <div class="bg-indigo-100 w-16 h-16 rounded-full flex items-center justify-center mb-4">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
          </div>
          <h3 class="text-xl font-semibold text-gray-800 mb-2">卡包管理</h3>
          <p class="text-gray-600 mb-4">创建和管理您的学习卡包，按主题组织学习内容</p>
          <router-link to="/decks" class="text-indigo-600 hover:text-indigo-800 font-medium flex items-center">
            开始使用
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 ml-1" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10.293 5.293a1 1 0 011.414 0l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L12.586 11H5a1 1 0 110-2h7.586l-2.293-2.293a1 1 0 010-1.414z" clip-rule="evenodd" />
            </svg>
          </router-link>
        </div>

        <div class="bg-white rounded-xl shadow-lg p-6 hover:shadow-xl transition-shadow">
          <div class="bg-purple-100 w-16 h-16 rounded-full flex items-center justify-center mb-4">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-purple-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
            </svg>
          </div>
          <h3 class="text-xl font-semibold text-gray-800 mb-2">标签管理</h3>
          <p class="text-gray-600 mb-4">使用标签分类和标记您的卡片，便于查找和复习</p>
          <router-link to="/tags" class="text-purple-600 hover:text-purple-800 font-medium flex items-center">
            开始使用
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 ml-1" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10.293 5.293a1 1 0 011.414 0l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L12.586 11H5a1 1 0 110-2h7.586l-2.293-2.293a1 1 0 010-1.414z" clip-rule="evenodd" />
            </svg>
          </router-link>
        </div>

        <div class="bg-white rounded-xl shadow-lg p-6 hover:shadow-xl transition-shadow">
          <div class="bg-green-100 w-16 h-16 rounded-full flex items-center justify-center mb-4">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <h3 class="text-xl font-semibold text-gray-800 mb-2">智能复习</h3>
          <p class="text-gray-600 mb-4">基于间隔重复算法，科学安排复习计划</p>
          <router-link to="/decks" class="text-green-600 hover:text-green-800 font-medium flex items-center">
            开始使用
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 ml-1" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10.293 5.293a1 1 0 011.414 0l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L12.586 11H5a1 1 0 110-2h7.586l-2.293-2.293a1 1 0 010-1.414z" clip-rule="evenodd" />
            </svg>
          </router-link>
        </div>
      </div>

      <!-- 统计信息 -->
      <div class="bg-white rounded-xl shadow-lg p-8 mb-16" v-if="!loading">
        <h2 class="text-2xl font-bold text-gray-800 mb-6 text-center">系统概览</h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div class="bg-indigo-50 rounded-lg p-6 text-center">
            <div class="text-3xl font-bold text-indigo-700 mb-2">{{ stats.totalDecks }}</div>
            <div class="text-indigo-600">卡包总数</div>
          </div>
          <div class="bg-purple-50 rounded-lg p-6 text-center">
            <div class="text-3xl font-bold text-purple-700 mb-2">{{ stats.totalCards }}</div>
            <div class="text-purple-600">卡片总数</div>
          </div>
          <div class="bg-green-50 rounded-lg p-6 text-center">
            <div class="text-3xl font-bold text-green-700 mb-2">{{ stats.totalTags }}</div>
            <div class="text-green-600">标签总数</div>
          </div>
        </div>
      </div>

      <!-- 加载状态 -->
      <div v-else class="bg-white rounded-xl shadow-lg p-8 mb-16 text-center">
        <div class="inline-block animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-600"></div>
        <p class="mt-4 text-gray-600">加载统计数据中...</p>
      </div>

      <!-- 快速开始 -->
      <div class="text-center">
        <h2 class="text-2xl font-bold text-gray-800 mb-4">快速开始</h2>
        <p class="text-gray-600 mb-6">创建您的第一个卡包，开始高效学习之旅</p>
        <router-link to="/decks" class="bg-indigo-600 text-white px-6 py-3 rounded-lg hover:bg-indigo-700 transition-colors inline-flex items-center">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clip-rule="evenodd" />
          </svg>
          创建卡包
        </router-link>
      </div>
    </div>
  </div>
</template>

<script>
import api from '../services/api'

export default {
  name: 'HomeView',
  data() {
    return {
      loading: true,
      stats: {
        totalDecks: 0,
        totalCards: 0,
        totalTags: 0
      }
    }
  },
  created() {
    this.fetchStats()
  },
  methods: {
    async fetchStats() {
      try {
        // 获取卡包统计
        const decksResponse = await api.getDecks(true)
        const decks = decksResponse.data?.data?.decks || []
        
        this.stats.totalDecks = decks.length
        
        // 计算总卡片数
        let totalCards = 0
        for (const deck of decks) {
          if (deck.stats) {
            totalCards += deck.stats.total_cards || 0
          }
        }
        this.stats.totalCards = totalCards
        
        // 获取所有标签
        if (decks.length > 0) {
          const tagsPromises = decks.map(deck => api.getTags(deck.id))
          const tagsResponses = await Promise.all(tagsPromises)
          
          let totalTags = 0
          tagsResponses.forEach(response => {
            totalTags += response.data?.data?.length || 0
          })
          this.stats.totalTags = totalTags
        } else {
          this.stats.totalTags = 0
        }
        
        this.loading = false
      } catch (error) {
        console.error('获取统计数据失败:', error)
        this.loading = false
      }
    }
  }
}
</script>