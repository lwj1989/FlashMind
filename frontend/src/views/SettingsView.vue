<template>
  <div class="min-h-screen bg-gray-50 p-4">
    <div class="max-w-6xl mx-auto">
      <header class="mb-8">
        <h1 class="text-3xl font-bold text-gray-800">系统设置</h1>
        <p class="text-gray-600 mt-2">管理您的数据和系统配置</p>
      </header>

      <!-- 数据管理 -->
      <div class="bg-white rounded-xl shadow p-6 mb-6">
        <h2 class="text-xl font-semibold text-gray-800 mb-6">数据管理</h2>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <!-- 备份数据 -->
          <div class="bg-blue-50 border border-blue-200 rounded-lg p-6">
            <div class="flex items-center mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10" />
              </svg>
              <h3 class="text-lg font-semibold text-gray-800 ml-3">备份数据</h3>
            </div>
            <p class="text-gray-600 mb-4">将所有数据导出为备份文件</p>
            <button 
              @click="backupAllData" 
              :disabled="isBackingUp"
              class="w-full bg-blue-600 text-white py-2 px-4 rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {{ isBackingUp ? '备份中...' : '开始备份' }}
            </button>
          </div>

          <!-- 恢复数据 -->
          <div class="bg-green-50 border border-green-200 rounded-lg p-6">
            <div class="flex items-center mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
              </svg>
              <h3 class="text-lg font-semibold text-gray-800 ml-3">恢复数据</h3>
            </div>
            <p class="text-gray-600 mb-4">从备份文件恢复所有数据</p>
            <input 
              ref="fileInput"
              type="file" 
              accept=".json,.zip"
              @change="handleFileSelect"
              class="hidden"
            />
            <button 
              @click="$refs.fileInput.click()"
              :disabled="isRestoring"
              class="w-full bg-green-600 text-white py-2 px-4 rounded-lg hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {{ isRestoring ? '恢复中...' : '选择备份文件' }}
            </button>
          </div>

          <!-- 清空数据 -->
          <div class="bg-red-50 border border-red-200 rounded-lg p-6">
            <div class="flex items-center mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
              <h3 class="text-lg font-semibold text-gray-800 ml-3">清空数据</h3>
            </div>
            <p class="text-gray-600 mb-4">删除所有卡片、卡包和标签</p>
            <button 
              @click="showClearConfirm = true"
              :disabled="isClearing"
              class="w-full bg-red-600 text-white py-2 px-4 rounded-lg hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {{ isClearing ? '清空中...' : '清空数据' }}
            </button>
          </div>
        </div>
      </div>

      <!-- 系统信息 -->
      <div class="bg-white rounded-xl shadow p-6">
        <h2 class="text-xl font-semibold text-gray-800 mb-6">系统信息</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div class="space-y-4">
            <div class="flex justify-between">
              <span class="text-gray-600">总卡包数量</span>
              <span class="font-semibold">{{ stats.totalDecks }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600">总卡片数量</span>
              <span class="font-semibold">{{ stats.totalCards }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600">总标签数量</span>
              <span class="font-semibold">{{ stats.totalTags }}</span>
            </div>
          </div>
          <div class="space-y-4">
            <div class="flex justify-between">
              <span class="text-gray-600">今日学习卡片</span>
              <span class="font-semibold">{{ stats.todayStudied }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600">待复习卡片</span>
              <span class="font-semibold">{{ stats.dueCards }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600">系统版本</span>
              <span class="font-semibold">v1.0.0</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 确认清空数据对话框 -->
      <div v-if="showClearConfirm" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50" @click="showClearConfirm = false">
        <div class="bg-white rounded-xl shadow-lg w-full max-w-md" @click.stop>
          <div class="p-6">
            <h3 class="text-xl font-semibold text-red-600 mb-4">⚠️ 危险操作</h3>
            <p class="text-gray-600 mb-6">
              此操作将<strong>永久删除</strong>所有数据，包括：
            </p>
            <ul class="text-gray-600 mb-6 list-disc list-inside">
              <li>所有卡包</li>
              <li>所有卡片</li>
              <li>所有标签</li>
              <li>所有学习记录</li>
            </ul>
            <p class="text-red-600 font-semibold mb-6">此操作不可撤销！</p>
            
            <div class="mb-6">
              <label class="block text-gray-700 text-sm font-medium mb-2">
                请输入 "DELETE" 以确认删除：
              </label>
              <input 
                v-model="confirmText"
                type="text" 
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-red-500"
                placeholder="输入 DELETE"
              />
            </div>
            
            <div class="flex justify-end space-x-3">
              <button 
                @click="showClearConfirm = false" 
                class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
              >
                取消
              </button>
              <button 
                @click="clearAllDataAction" 
                :disabled="confirmText !== 'DELETE' || isClearing"
                class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {{ isClearing ? '清空中...' : '确认清空' }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 恢复数据确认对话框 -->
      <div v-if="showRestoreConfirm" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50" @click="showRestoreConfirm = false">
        <div class="bg-white rounded-xl shadow-lg w-full max-w-md" @click.stop>
          <div class="p-6">
            <h3 class="text-xl font-semibold text-yellow-600 mb-4">⚠️ 恢复数据</h3>
            <p class="text-gray-600 mb-6">
              恢复数据将会：
            </p>
            <ul class="text-gray-600 mb-6 list-disc list-inside">
              <li>清空现有的所有数据</li>
              <li>导入备份文件中的数据</li>
              <li>可能覆盖当前的学习进度</li>
            </ul>
            <p class="text-yellow-600 font-semibold mb-6">建议先备份当前数据！</p>
            
            <div class="flex justify-end space-x-3">
              <button 
                @click="showRestoreConfirm = false" 
                class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
              >
                取消
              </button>
              <button 
                @click="confirmRestoreData" 
                :disabled="isRestoring"
                class="px-4 py-2 bg-yellow-600 text-white rounded-lg hover:bg-yellow-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {{ isRestoring ? '恢复中...' : '确认恢复' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getDecks } from '@/api/deck'
import { getAllTags } from '@/api/tag'
import { searchCards } from '@/api/card'
import { getDueCards } from '@/api/study'
import { getSystemStats, backupData, restoreData, clearAllData } from '@/api/system'

export default {
  name: 'SettingsView',
  setup() {
    const isBackingUp = ref(false)
    const isRestoring = ref(false)
    const isClearing = ref(false)
    const showClearConfirm = ref(false)
    const showRestoreConfirm = ref(false)
    const confirmText = ref('')
    const selectedFile = ref(null)
    
    const stats = reactive({
      totalDecks: 0,
      totalCards: 0,
      totalTags: 0,
      todayStudied: 0,
      dueCards: 0
    })

    // 加载统计信息
    const loadStats = async () => {
      try {
        // 使用系统统计API
        const response = await getSystemStats()
        if (response.data && response.data.code === 'SUCCESS') {
          const data = response.data.data
          stats.totalDecks = data.total_decks || 0
          stats.totalCards = data.total_cards || 0
          stats.totalTags = data.total_tags || 0
        }

        // 获取到期卡片数量
        const dueResponse = await getDueCards(1000)
        stats.dueCards = dueResponse.data.data?.queue?.length || 0

        // 今日学习数量 (暂时模拟)
        stats.todayStudied = 0
      } catch (error) {
        console.error('加载统计信息失败:', error)
      }
    }

    // 备份所有数据
    const backupAllData = async () => {
      isBackingUp.value = true
      try {
        const response = await backupData()
        
        // 创建下载链接
        const blob = new Blob([response.data], { type: 'application/json' })
        const url = URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = url
        link.download = `flashmind_backup_${new Date().toISOString().split('T')[0]}.json`
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        URL.revokeObjectURL(url)

        ElMessage.success('数据备份成功')
      } catch (error) {
        console.error('备份失败:', error)
        ElMessage.error('备份失败: ' + error.message)
      } finally {
        isBackingUp.value = false
      }
    }

    // 处理文件选择
    const handleFileSelect = (event) => {
      const file = event.target.files[0]
      if (file) {
        selectedFile.value = file
        showRestoreConfirm.value = true
      }
    }

    // 确认恢复数据
    const confirmRestoreData = async () => {
      if (!selectedFile.value) return

      isRestoring.value = true
      showRestoreConfirm.value = false
      
      try {
        const response = await restoreData(selectedFile.value, true)
        
        if (response.data && response.data.code === 'SUCCESS') {
          ElMessage.success(response.data.message || '数据恢复成功')
        } else {
          throw new Error(response.data?.message || '恢复失败')
        }
        
        loadStats() // 重新加载统计信息
      } catch (error) {
        console.error('恢复失败:', error)
        ElMessage.error('恢复失败: ' + (error.response?.data?.message || error.message))
      } finally {
        isRestoring.value = false
        selectedFile.value = null
        // 重置文件输入
        const fileInput = document.querySelector('input[type="file"]')
        if (fileInput) fileInput.value = ''
      }
    }

    // 清空所有数据
    const clearAllDataAction = async () => {
      if (confirmText.value !== 'DELETE') return

      isClearing.value = true
      showClearConfirm.value = false
      
      try {
        const response = await clearAllData()
        
        if (response.data && response.data.code === 'SUCCESS') {
          ElMessage.success(response.data.message || '数据清空成功')
        } else {
          throw new Error(response.data?.message || '清空失败')
        }
        
        loadStats() // 重新加载统计信息
      } catch (error) {
        console.error('清空失败:', error)
        ElMessage.error('清空失败: ' + (error.response?.data?.message || error.message))
      } finally {
        isClearing.value = false
        confirmText.value = ''
      }
    }

    onMounted(() => {
      loadStats()
    })

    return {
      isBackingUp,
      isRestoring,
      isClearing,
      showClearConfirm,
      showRestoreConfirm,
      confirmText,
      stats,
      backupAllData,
      handleFileSelect,
      confirmRestoreData,
      clearAllDataAction
    }
  }
}
</script>
