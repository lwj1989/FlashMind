<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-6xl mx-auto px-4 py-12">
      <header class="mb-8">
        <h1 class="text-3xl font-bold text-gray-800">导入导出</h1>
        <p class="text-gray-600 mt-2">管理卡包的导入和导出</p>
      </header>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- 导入卡包部分 -->
        <div class="bg-white rounded-xl shadow p-6">
          <h2 class="text-xl font-semibold text-gray-800 mb-4">导入卡包</h2>
          <div class="import-container">
            <div class="upload-area" 
                 :class="{ 'drag-over': isDragOver }"
                 @drop.prevent="handleDrop"
                 @dragover.prevent="handleDragOver"
                 @dragleave.prevent="handleDragLeave">
              <input type="file" 
                     ref="fileInput" 
                     @change="handleFileChange" 
                     accept=".json,.csv,.txt"
                     style="display: none">
              
              <div class="upload-content">
                <div class="upload-icon">📁</div>
                <p>拖放文件到此处或 <span @click="triggerFileInput" class="upload-link">点击上传</span></p>
                <p class="upload-hint">支持 JSON、CSV 和 TXT 格式</p>
              </div>
            </div>
            
            <div v-if="selectedFile" class="file-info mt-4">
              <p class="text-gray-700">已选择文件: {{ selectedFile.name }}</p>
              <p class="text-gray-700">文件大小: {{ formatFileSize(selectedFile.size) }}</p>
              
              <div class="mt-4">
                <label for="deck-name" class="block text-sm font-medium text-gray-700 mb-1">卡包名称</label>
                <input 
                  id="deck-name" 
                  type="text" 
                  v-model="deckName" 
                  placeholder="请输入卡包名称"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  required
                />
              </div>
              
              <div class="flex space-x-3 mt-4">
                <button @click="uploadFile" :disabled="uploading || !deckName.trim()" class="bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700 transition-colors flex-1 disabled:opacity-50">
                  {{ uploading ? '上传中...' : '导入' }}
                </button>
                <button @click="clearFile" class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 flex-1">
                  取消
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- 导出卡包部分 -->
        <div class="bg-white rounded-xl shadow p-6">
          <h2 class="text-xl font-semibold text-gray-800 mb-4">导出卡包</h2>
          <div class="export-container">
            <div class="mb-4">
              <label for="deck-select" class="block text-sm font-medium text-gray-700 mb-1">选择卡包</label>
              <select id="deck-select" v-model="selectedDeckId" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500">
                <option value="">请选择卡包</option>
                <option v-for="deck in decks" :key="deck.id" :value="deck.id">
                  {{ deck.name }}
                </option>
              </select>
            </div>
            
            <div class="mb-6">
              <label class="block text-sm font-medium text-gray-700 mb-2">导出格式</label>
              <div class="flex space-x-4">
                <label class="flex items-center">
                  <input type="radio" v-model="exportFormat" value="json" class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded">
                  <span class="ml-2 text-gray-700">JSON</span>
                </label>
                <label class="flex items-center">
                  <input type="radio" v-model="exportFormat" value="csv" class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded">
                  <span class="ml-2 text-gray-700">CSV</span>
                </label>
                <label class="flex items-center">
                  <input type="radio" v-model="exportFormat" value="txt" class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded">
                  <span class="ml-2 text-gray-700">TXT</span>
                </label>
              </div>
            </div>
            
            <button @click="exportDeck" :disabled="!selectedDeckId || exporting" class="w-full bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700 transition-colors">
              {{ exporting ? '导出中...' : '导出' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 导入结果对话框 -->
    <div v-if="importResult" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50" @click="closeImportResult">
      <div class="bg-white rounded-xl shadow-lg w-full max-w-md" @click.stop>
        <div class="p-6">
          <h3 class="text-xl font-semibold text-gray-800 mb-4">导入结果</h3>
          <div class="result-content">
            <p v-if="importResult.success" class="text-green-600 font-medium mb-4">
              导入成功！
            </p>
            <p v-else class="text-red-600 font-medium mb-4">
              导入失败：{{ importResult.message }}
            </p>
            
            <div v-if="importResult.data" class="bg-gray-50 p-4 rounded-lg mb-4">
              <p class="text-gray-700">卡包名称: {{ importResult.data.deck_name }}</p>
              <p class="text-gray-700">导入卡片数量: {{ importResult.data.card_count }}</p>
            </div>
          </div>
          <div class="flex justify-end">
            <button @click="closeImportResult" class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700">
              确定
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getDecks } from '@/api/deck'
import { importDeck, exportDeck } from '@/api/importExport'

export default {
  name: 'ImportExportView',
  setup() {
    const fileInput = ref(null)
    const selectedFile = ref(null)
    const isDragOver = ref(false)
    const uploading = ref(false)
    const exporting = ref(false)
    const selectedDeckId = ref('')
    const exportFormat = ref('json')
    const decks = ref([])
    const importResult = ref(null)
    const deckName = ref('')

    // 获取卡包列表
    const fetchDecks = async () => {
      try {
        const response = await getDecks()
        console.log('获取卡包响应:', response)
        if (response.data && response.data.code === 'SUCCESS') {
          decks.value = response.data.data.decks || response.data.data || []
        } else {
          decks.value = response.data.decks || response.data || []
        }
      } catch (error) {
        console.error('获取卡包列表失败:', error)
        ElMessage.error('获取卡包列表失败')
      }
    }

    // 触发文件选择
    const triggerFileInput = () => {
      fileInput.value.click()
    }

    // 处理文件选择
    const handleFileChange = (event) => {
      const file = event.target.files[0]
      if (file) {
        validateAndSetFile(file)
      }
    }

    // 处理拖放
    const handleDrop = (event) => {
      isDragOver.value = false
      const file = event.dataTransfer.files[0]
      if (file) {
        validateAndSetFile(file)
      }
    }

    // 处理拖拽悬停
    const handleDragOver = () => {
      isDragOver.value = true
    }

    // 处理拖拽离开
    const handleDragLeave = () => {
      isDragOver.value = false
    }

    // 验证并设置文件
    const validateAndSetFile = (file) => {
      const validTypes = ['application/json', 'text/csv', 'text/plain']
      const validExtensions = ['.json', '.csv', '.txt']
      
      const isValidType = validTypes.includes(file.type) || 
                         validExtensions.some(ext => file.name.toLowerCase().endsWith(ext))
      
      if (!isValidType) {
        ElMessage.error('请上传 JSON、CSV 或 TXT 格式的文件')
        return
      }
      
      selectedFile.value = file
      
      // 自动填充卡包名称（去掉文件扩展名）
      const fileName = file.name
      const lastDotIndex = fileName.lastIndexOf('.')
      const nameWithoutExtension = lastDotIndex > 0 ? fileName.substring(0, lastDotIndex) : fileName
      deckName.value = nameWithoutExtension
    }

    // 清除选择的文件
    const clearFile = () => {
      selectedFile.value = null
      deckName.value = ''
      if (fileInput.value) {
        fileInput.value.value = ''
      }
    }

    // 上传文件
    const uploadFile = async () => {
      if (!selectedFile.value || !deckName.value.trim()) return
      
      uploading.value = true
      try {
        const formData = new FormData()
        formData.append('file', selectedFile.value)
        formData.append('deck_name', deckName.value.trim())
        const response = await importDeck(formData)
        console.log('导入响应:', response)
        let responseData = response.data
        if (response.data && response.data.code === 'SUCCESS') {
          responseData = response.data.data
        }
        importResult.value = {
          success: true,
          data: responseData
        }
        ElMessage.success('导入成功')
        clearFile()
      } catch (error) {
        console.error('导入失败:', error)
        importResult.value = {
          success: false,
          message: error.response?.data?.message || '导入失败'
        }
        ElMessage.error('导入失败')
      } finally {
        uploading.value = false
      }
    }

    // 导出卡包
    const exportDeckFile = async () => {
      if (!selectedDeckId.value) return
      
      exporting.value = true
      try {
        const response = await exportDeck(selectedDeckId.value, exportFormat.value)
        
        // 创建下载链接
        const url = window.URL.createObjectURL(new Blob([response.data]))
        const link = document.createElement('a')
        link.href = url
        
        // 获取卡包名称作为文件名
        const deck = decks.value.find(d => d.id === selectedDeckId.value)
        const filename = deck ? `${deck.name}.${exportFormat.value}` : `deck.${exportFormat.value}`
        
        link.setAttribute('download', filename)
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        
        ElMessage.success('导出成功')
      } catch (error) {
        console.error('导出失败:', error)
        ElMessage.error('导出失败')
      } finally {
        exporting.value = false
      }
    }

    // 关闭导入结果对话框
    const closeImportResult = () => {
      importResult.value = null
    }

    // 格式化文件大小
    const formatFileSize = (bytes) => {
      if (bytes === 0) return '0 Bytes'
      
      const k = 1024
      const sizes = ['Bytes', 'KB', 'MB', 'GB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
    }

    // ESC 键关闭模态框
    const handleEscapeKey = (event) => {
      if (event.key === 'Escape') {
        if (importResult.value) {
          importResult.value = null
        }
      }
    }

    onMounted(() => {
      // 添加键盘事件监听器
      document.addEventListener('keydown', handleEscapeKey)
      fetchDecks()
    })

    onUnmounted(() => {
      // 移除键盘事件监听器
      document.removeEventListener('keydown', handleEscapeKey)
    })

    return {
      fileInput,
      selectedFile,
      isDragOver,
      uploading,
      exporting,
      selectedDeckId,
      exportFormat,
      decks,
      importResult,
      deckName,
      triggerFileInput,
      handleFileChange,
      handleDrop,
      handleDragOver,
      handleDragLeave,
      clearFile,
      uploadFile,
      exportDeck: exportDeckFile,
      closeImportResult,
      formatFileSize
    }
  }
}
</script>

<style scoped>
.import-export-view {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.view-header {
  margin-bottom: 30px;
}

.content-container {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.section {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.section h2 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #333;
}

/* 导入部分样式 */
.import-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.upload-area {
  border: 2px dashed #ccc;
  border-radius: 8px;
  padding: 40px;
  text-align: center;
  transition: border-color 0.3s;
}

.upload-area.drag-over {
  border-color: #4CAF50;
  background-color: rgba(76, 175, 80, 0.05);
}

.upload-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.upload-icon {
  font-size: 48px;
}

.upload-link {
  color: #4CAF50;
  cursor: pointer;
  text-decoration: underline;
}

.upload-hint {
  color: #666;
  font-size: 14px;
}

.file-info {
  background-color: #f5f5f5;
  border-radius: 8px;
  padding: 15px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

/* 导出部分样式 */
.export-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-weight: bold;
  color: #333;
}

.form-control {
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
}

.radio-group {
  display: flex;
  gap: 20px;
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 5px;
  cursor: pointer;
}

/* 按钮样式 */
.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  transition: background-color 0.3s;
}

.btn-primary {
  background-color: #4CAF50;
  color: white;
}

.btn-primary:hover {
  background-color: #45a049;
}

.btn-primary:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.btn-secondary {
  background-color: #9e9e9e;
  color: white;
}

.btn-secondary:hover {
  background-color: #757575;
}

/* 模态框样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
  padding: 20px;
  max-width: 500px;
  width: 90%;
}

.modal-content h3 {
  margin-top: 0;
  margin-bottom: 15px;
  text-align: center;
}

.result-content {
  margin-bottom: 20px;
}

.success-message {
  color: #4CAF50;
  font-weight: bold;
  text-align: center;
}

.error-message {
  color: #f44336;
  font-weight: bold;
  text-align: center;
}

.result-details {
  margin-top: 15px;
  padding: 10px;
  background-color: #f5f5f5;
  border-radius: 4px;
}

.form-actions {
  display: flex;
  justify-content: center;
}
</style>