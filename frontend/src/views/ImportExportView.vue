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
          <h2 class="text-xl font-semibold text-gray-800 mb-4">批量导入卡包</h2>
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
                     multiple
                     style="display: none">
              
              <div class="upload-content">
                <div class="upload-icon">📁</div>
                <p>拖放文件到此处或 <span @click="triggerFileInput" class="upload-link">点击选择</span></p>
                <p class="upload-hint">支持同时选择多个 JSON、CSV 和 TXT 文件</p>
                <p class="upload-hint text-indigo-600">每个文件将作为独立的卡包导入</p>
              </div>
            </div>
            
            <!-- 多文件显示 -->
            <div v-if="selectedFiles.length > 0" class="files-info mt-4">
              <div class="mb-4">
                <h3 class="text-lg font-medium text-gray-800 mb-2">已选择 {{ selectedFiles.length }} 个文件</h3>
                <div class="max-h-48 overflow-y-auto space-y-2">
                  <div v-for="(file, index) in selectedFiles" :key="index" 
                       class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                    <div class="flex-1">
                      <p class="font-medium text-gray-800">{{ getFileName(file.name) }}</p>
                      <p class="text-sm text-gray-600">{{ file.name }} ({{ formatFileSize(file.size) }})</p>
                    </div>
                    <button @click="removeFile(index)" 
                            class="text-red-500 hover:text-red-700 p-1" 
                            title="移除文件">
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                      </svg>
                    </button>
                  </div>
                </div>
              </div>
              
              <!-- 批量导入进度 -->
              <div v-if="batchImporting" class="mb-4">
                <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
                  <div class="flex items-center mb-2">
                    <svg class="animate-spin h-4 w-4 text-blue-600 mr-2" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <span class="text-sm font-medium text-blue-800">
                      正在导入 {{ currentImportIndex + 1 }} / {{ selectedFiles.length }}：{{ selectedFiles[currentImportIndex]?.name }}
                    </span>
                  </div>
                  <div class="w-full bg-blue-200 rounded-full h-2">
                    <div class="bg-blue-600 h-2 rounded-full transition-all duration-300" 
                         :style="{ width: `${((currentImportIndex + 1) / selectedFiles.length) * 100}%` }"></div>
                  </div>
                </div>
              </div>
              
              <!-- 导入结果 -->
              <div v-if="batchImportResults.length > 0" class="mb-4">
                <h4 class="text-md font-medium text-gray-800 mb-2">导入结果</h4>
                <div class="max-h-32 overflow-y-auto space-y-1">
                  <div v-for="(result, index) in batchImportResults" :key="index" 
                       class="flex items-center text-sm"
                       :class="result.success ? 'text-green-600' : 'text-red-600'">
                    <svg v-if="result.success" class="h-4 w-4 mr-2" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                    </svg>
                    <svg v-else class="h-4 w-4 mr-2" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                    </svg>
                    <span>{{ result.fileName }}: {{ result.message }}</span>
                  </div>
                </div>
              </div>
              
              <div class="flex space-x-3 mt-4">
                <button @click="batchUploadFiles" 
                        :disabled="batchImporting || selectedFiles.length === 0" 
                        class="bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700 transition-colors flex-1 disabled:opacity-50">
                  {{ batchImporting ? '导入中...' : `批量导入 ${selectedFiles.length} 个文件` }}
                </button>
                <button @click="clearFiles" 
                        :disabled="batchImporting"
                        class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 flex-1 disabled:opacity-50">
                  清空列表
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
    const selectedFiles = ref([])
    const isDragOver = ref(false)
    const uploading = ref(false)
    const exporting = ref(false)
    const selectedDeckId = ref('')
    const exportFormat = ref('json')
    const decks = ref([])
    const importResult = ref(null)
    const deckName = ref('')
    
    // 批量导入相关
    const batchImporting = ref(false)
    const currentImportIndex = ref(0)
    const batchImportResults = ref([])
    
    // 保持向后兼容的单文件引用
    const selectedFile = ref(null)

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
      const files = Array.from(event.target.files)
      if (files.length > 0) {
        validateAndSetFiles(files)
      }
    }

    // 处理拖放
    const handleDrop = (event) => {
      isDragOver.value = false
      const files = Array.from(event.dataTransfer.files)
      if (files.length > 0) {
        validateAndSetFiles(files)
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
    const validateAndSetFiles = (files) => {
      const validTypes = ['application/json', 'text/csv', 'text/plain']
      const validExtensions = ['.json', '.csv', '.txt']
      
      const validFiles = []
      const invalidFiles = []
      
      files.forEach(file => {
        const isValidType = validTypes.includes(file.type) || 
                           validExtensions.some(ext => file.name.toLowerCase().endsWith(ext))
        
        if (isValidType) {
          validFiles.push(file)
        } else {
          invalidFiles.push(file.name)
        }
      })
      
      if (invalidFiles.length > 0) {
        ElMessage.error(`以下文件格式不支持: ${invalidFiles.join(', ')}。请上传 JSON、CSV 或 TXT 格式的文件`)
      }
      
      if (validFiles.length > 0) {
        // 添加到现有文件列表中，避免重复
        validFiles.forEach(file => {
          const exists = selectedFiles.value.some(existingFile => 
            existingFile.name === file.name && existingFile.size === file.size
          )
          if (!exists) {
            selectedFiles.value.push(file)
          }
        })
        
        ElMessage.success(`成功添加 ${validFiles.length} 个文件`)
      }
    }

    // 获取文件名（去掉扩展名）
    const getFileName = (fullName) => {
      const lastDotIndex = fullName.lastIndexOf('.')
      return lastDotIndex > 0 ? fullName.substring(0, lastDotIndex) : fullName
    }

    // 移除单个文件
    const removeFile = (index) => {
      selectedFiles.value.splice(index, 1)
    }

    // 清除所有文件
    const clearFiles = () => {
      selectedFiles.value = []
      batchImportResults.value = []
      if (fileInput.value) {
        fileInput.value.value = ''
      }
    }

    // 兼容性：清除选择的文件
    const clearFile = () => {
      clearFiles()
    }

    // 批量上传文件
    const batchUploadFiles = async () => {
      if (selectedFiles.value.length === 0) return
      
      batchImporting.value = true
      currentImportIndex.value = 0
      batchImportResults.value = []
      
      let successCount = 0
      let failureCount = 0
      
      for (let i = 0; i < selectedFiles.value.length; i++) {
        currentImportIndex.value = i
        const file = selectedFiles.value[i]
        const fileName = getFileName(file.name)
        
        try {
          const formData = new FormData()
          formData.append('file', file)
          formData.append('deck_name', fileName)
          
          const response = await importDeck(formData)
          console.log('导入响应:', response)
          
          let responseData = response.data
          if (response.data && response.data.code === 'SUCCESS') {
            responseData = response.data.data
          }
          
          batchImportResults.value.push({
            fileName: fileName,
            success: true,
            message: `成功导入 ${responseData.card_count || 0} 张卡片`,
            data: responseData
          })
          
          successCount++
        } catch (error) {
          console.error(`导入文件 ${file.name} 失败:`, error)
          
          batchImportResults.value.push({
            fileName: fileName,
            success: false,
            message: error.response?.data?.message || '导入失败',
            error: error
          })
          
          failureCount++
        }
        
        // 小延迟，避免过快请求
        await new Promise(resolve => setTimeout(resolve, 300))
      }
      
      batchImporting.value = false
      
      // 显示最终结果
      if (successCount > 0 && failureCount === 0) {
        ElMessage.success(`批量导入完成！成功导入 ${successCount} 个卡包`)
      } else if (successCount > 0 && failureCount > 0) {
        ElMessage.warning(`批量导入完成！成功 ${successCount} 个，失败 ${failureCount} 个`)
      } else {
        ElMessage.error(`批量导入失败！${failureCount} 个文件导入失败`)
      }
      
      // 如果全部成功，清空文件列表
      if (failureCount === 0) {
        setTimeout(() => {
          clearFiles()
        }, 3000)
      }
    }

    // 兼容性：单文件上传
    const uploadFile = async () => {
      if (selectedFiles.value.length === 1) {
        await batchUploadFiles()
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
      selectedFiles,
      isDragOver,
      uploading,
      exporting,
      selectedDeckId,
      exportFormat,
      decks,
      importResult,
      deckName,
      // 批量导入相关
      batchImporting,
      currentImportIndex,
      batchImportResults,
      // 函数
      triggerFileInput,
      handleFileChange,
      handleDrop,
      handleDragOver,
      handleDragLeave,
      clearFile,
      clearFiles,
      removeFile,
      getFileName,
      uploadFile,
      batchUploadFiles,
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