<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-6xl mx-auto px-4 py-12">
      <!-- 学习选择界面 -->
      <div v-if="!studySession" class="space-y-8">
        <!-- 学习统计 -->
        <div class="bg-white rounded-xl shadow p-6">
          <h2 class="text-xl font-semibold text-gray-800 mb-6">学习统计</h2>
          <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div class="bg-gradient-to-br from-blue-50 to-blue-100 rounded-lg p-4 text-center">
              <div class="text-2xl font-bold text-blue-600">{{ studyStats.dueCards }}</div>
              <div class="text-sm text-gray-600">待复习</div>
            </div>
            <div class="bg-gradient-to-br from-green-50 to-green-100 rounded-lg p-4 text-center">
              <div class="text-2xl font-bold text-green-600">{{ studyStats.totalCards }}</div>
              <div class="text-sm text-gray-600">总卡片</div>
            </div>
            <div class="bg-gradient-to-br from-purple-50 to-purple-100 rounded-lg p-4 text-center">
              <div class="text-2xl font-bold text-purple-600">{{ studyStats.studiedToday }}</div>
              <div class="text-sm text-gray-600">今日已学</div>
            </div>
            <div class="bg-gradient-to-br from-orange-50 to-orange-100 rounded-lg p-4 text-center">
              <div class="text-2xl font-bold text-orange-600">{{ studyStats.streak }}</div>
              <div class="text-sm text-gray-600">连续天数</div>
            </div>
          </div>
        </div>
        
        <!-- 学习模式选择 -->
        <div class="bg-white rounded-xl shadow p-8">
          <h1 class="text-3xl font-bold text-gray-800 mb-8 text-center">开始学习</h1>
          
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <!-- 到期复习 -->
          <div class="bg-gradient-to-br from-blue-50 to-blue-100 rounded-xl p-6 border border-blue-200">
            <div class="text-center">
              <div class="text-blue-600 mb-4">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <h3 class="text-lg font-semibold text-gray-800 mb-2">到期复习</h3>
              <p class="text-gray-600 mb-4">复习到期的卡片</p>
              <button @click="startDueStudy" class="w-full bg-blue-600 text-white py-2 px-4 rounded-lg hover:bg-blue-700 transition-colors">
                开始复习
              </button>
            </div>
          </div>

          <!-- 卡包学习 -->
          <div class="bg-gradient-to-br from-green-50 to-green-100 rounded-xl p-6 border border-green-200">
            <div class="text-center">
              <div class="text-green-600 mb-4">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
              </div>
              <h3 class="text-lg font-semibold text-gray-800 mb-2">卡包学习</h3>
              <div class="mb-4">
                <select v-model="selectedDeckId" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500 text-sm">
                  <option value="">选择卡包</option>
                  <option v-for="deck in decks" :key="deck.id" :value="deck.id">
                    {{ deck.name }}
                  </option>
                </select>
              </div>
              <button @click="startDeckStudy" :disabled="!selectedDeckId" class="w-full bg-green-600 text-white py-2 px-4 rounded-lg hover:bg-green-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
                开始学习
              </button>
            </div>
          </div>

          <!-- 随机学习 -->
          <div class="bg-gradient-to-br from-purple-50 to-purple-100 rounded-xl p-6 border border-purple-200">
            <div class="text-center">
              <div class="text-purple-600 mb-4">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
              </div>
              <h3 class="text-lg font-semibold text-gray-800 mb-2">随机学习</h3>
              <div class="mb-4">
                <input v-model.number="randomLimit" type="number" min="5" max="50" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 text-sm" placeholder="卡片数量">
              </div>
              <button @click="startRandomStudy" class="w-full bg-purple-600 text-white py-2 px-4 rounded-lg hover:bg-purple-700 transition-colors">
                随机学习
              </button>
            </div>
          </div>
          </div>
        </div>

        <!-- 标签 -->
        <div v-if="tags.length > 0" class="mt-8">
          <h2 class="text-xl font-semibold text-gray-800 mb-4">标签</h2>
          <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
            <button 
              v-for="tag in tags" 
              :key="tag.id" 
              @click="startTagStudy(tag.id)" 
              class="bg-white border border-gray-300 hover:border-indigo-500 hover:bg-indigo-50 py-3 px-4 rounded-lg transition-colors flex items-center justify-center"
            >
              <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-sm font-medium bg-indigo-100 text-indigo-800">
                {{ tag.name }}
              </span>
            </button>
          </div>
        </div>
      </div>

      <!-- 学习界面 -->
      <div v-else class="space-y-6">
        <!-- 进度条 -->
        <div class="bg-white rounded-xl shadow p-6">
          <div class="flex justify-between items-center mb-2">
            <span class="text-sm font-medium text-gray-700">学习进度</span>
            <span class="text-sm text-gray-500">{{ Math.min(studySession.current + 1, studySession.total) }} / {{ studySession.total }}</span>
          </div>
          <div class="w-full bg-gray-200 rounded-full h-2">
            <div class="bg-indigo-600 h-2 rounded-full transition-all duration-300" :style="{ width: `${Math.min(((studySession.current + 1) / studySession.total) * 100, 100)}%` }"></div>
          </div>
        </div>

        <!-- 卡片内容 -->
        <div v-if="currentCard" class="bg-white rounded-xl shadow">
          <!-- 卡片头部 -->
          <div class="bg-indigo-600 text-white p-6 rounded-t-xl">
            <div class="flex justify-between items-start">
              <div>
                <h3 class="text-lg font-medium">{{ currentCard.deck_name }}</h3>
              </div>
              <button @click="exitStudy" class="text-indigo-200 hover:text-white">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>

          <!-- 问题 -->
          <div class="p-6 border-b border-gray-200">
            <h2 class="text-xl font-semibold text-gray-800 mb-4">问题</h2>
            <div class="prose prose-slate max-w-none mb-4">
              <div class="markdown-content" v-html="renderMarkdown(currentCard.question)"></div>
            </div>
            
            <!-- 标签显示在问题下方 -->
            <div v-if="currentCard.tag_name" class="mt-4">
              <button 
                @click="goToTagCards(currentCard.tag_id)" 
                class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-indigo-100 text-indigo-800 hover:bg-indigo-200 transition-colors cursor-pointer"
                :title="'点击查看标签「' + currentCard.tag_name + '」的所有卡片'"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M17.707 9.293a1 1 0 010 1.414l-7 7a1 1 0 01-1.414 0l-7-7A.997.997 0 012 10V5a3 3 0 013-3h5c.256 0 .512.098.707.293l7 7zM5 6a1 1 0 100-2 1 1 0 000 2z" clip-rule="evenodd" />
                </svg>
                {{ currentCard.tag_name }}
              </button>
            </div>
          </div>

          <!-- 答案（显示/隐藏） -->
          <div class="p-6">
            <div v-if="!showAnswer" class="text-center">
              <button @click="showAnswer = true" class="bg-indigo-600 text-white px-6 py-3 rounded-lg hover:bg-indigo-700 transition-colors">
                显示答案
              </button>
            </div>
            <div v-else>
              <h2 class="text-xl font-semibold text-gray-800 mb-4">答案</h2>
              <div class="prose prose-slate max-w-none mb-6">
                <div class="markdown-content" v-html="renderMarkdown(currentCard.answer)"></div>
              </div>
              
              <!-- 复习按钮 -->
              <div class="flex justify-center space-x-4">
                <button @click="submitReview(0)" class="bg-red-500 text-white px-6 py-3 rounded-lg hover:bg-red-600 transition-colors">
                  忘记了
                </button>
                <button @click="submitReview(1)" class="bg-yellow-500 text-white px-6 py-3 rounded-lg hover:bg-yellow-600 transition-colors">
                  模糊
                </button>
                <button @click="submitReview(2)" class="bg-green-500 text-white px-6 py-3 rounded-lg hover:bg-green-600 transition-colors">
                  记得
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- 学习完成 -->
        <div v-if="studyComplete" class="bg-white rounded-xl shadow p-8 text-center">
          <div class="text-green-600 mb-4">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <h3 class="text-2xl font-bold text-gray-800 mb-2">学习完成！</h3>
          <p class="text-gray-600 mb-6">恭喜你完成了这次学习，坚持学习是成功的关键！</p>
          <button @click="exitStudy" class="bg-indigo-600 text-white px-6 py-3 rounded-lg hover:bg-indigo-700 transition-colors">
            返回学习首页
          </button>
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

.markdown-content :deep(blockquote) {
  border-left: 4px solid #3b82f6;
  background-color: #f8fafc;
  padding: 12px 16px;
  margin: 16px 0;
  border-radius: 0 8px 8px 0;
}

.markdown-content :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin: 16px 0;
}

.markdown-content :deep(th),
.markdown-content :deep(td) {
  border: 1px solid #e2e8f0;
  padding: 8px 12px;
  text-align: left;
}

.markdown-content :deep(th) {
  background-color: #f8fafc;
  font-weight: 600;
}

.markdown-content :deep(h1),
.markdown-content :deep(h2),
.markdown-content :deep(h3),
.markdown-content :deep(h4),
.markdown-content :deep(h5),
.markdown-content :deep(h6) {
  margin: 24px 0 16px 0;
  font-weight: 600;
  line-height: 1.25;
}

.markdown-content :deep(p) {
  margin: 12px 0;
}

.markdown-content :deep(ul),
.markdown-content :deep(ol) {
  margin: 12px 0;
  padding-left: 24px;
}

.markdown-content :deep(li) {
  margin: 4px 0;
}
</style>

<script>
import { ref, reactive, computed, onMounted, getCurrentInstance } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import { getDecks } from '@/api/deck'
import { getAllTags } from '@/api/tag'
import { startDeckStudy, startTagStudy, startRandomStudy, getDueCards, submitReview } from '@/api/study'

export default {
  name: 'StudyView',
  setup() {
    const studySession = ref(null)
    const showAnswer = ref(false)
    const decks = ref([])
    const tags = ref([])
    const selectedDeckId = ref('')
    const randomLimit = ref(10)
    
    // 学习统计数据
    const studyStats = reactive({
      dueCards: 0,
      totalCards: 0,
      studiedToday: 0,
      streak: 0
    })

    const currentCard = computed(() => {
      if (!studySession.value || studySession.value.current >= studySession.value.queue.length) {
        return null
      }
      return studySession.value.queue[studySession.value.current]
    })

    const studyComplete = computed(() => {
      return studySession.value && studySession.value.current >= studySession.value.total
    })

    // 获取卡包列表
    const fetchDecks = async () => {
      try {
        const response = await getDecks()
        if (response.data && response.data.code === 'SUCCESS') {
          decks.value = response.data.data.decks || response.data.data || []
        } else {
          decks.value = response.data.decks || response.data || []
        }
      } catch (error) {
        console.error('获取卡包失败:', error)
      }
    }

    // 获取标签列表
    const fetchTags = async () => {
      try {
        const response = await getAllTags()
        if (response.data && response.data.code === 'SUCCESS') {
          tags.value = response.data.data || []
        } else {
          tags.value = response.data || []
        }
      } catch (error) {
        console.error('获取标签失败:', error)
      }
    }

    // 开始到期复习
    const startDueStudy = async () => {
      try {
        const response = await getDueCards(20)
        const session = response.data.code === 'SUCCESS' ? response.data.data : response.data
        if (session.queue.length === 0) {
          ElMessage.info('当前没有到期的卡片需要复习')
          return
        }
        studySession.value = session
        showAnswer.value = false
      } catch (error) {
        console.error('开始复习失败:', error)
        ElMessage.error('开始复习失败')
      }
    }

    // 开始卡包学习
    const startDeckStudySession = async () => {
      if (!selectedDeckId.value) {
        ElMessage.warning('请选择一个卡包')
        return
      }
      try {
        const response = await startDeckStudy(selectedDeckId.value, 20)
        const session = response.data.code === 'SUCCESS' ? response.data.data : response.data
        if (session.queue.length === 0) {
          ElMessage.info('该卡包中没有卡片')
          return
        }
        studySession.value = session
        showAnswer.value = false
      } catch (error) {
        console.error('开始学习失败:', error)
        ElMessage.error('开始学习失败')
      }
    }

    // 开始标签学习
    const startTagStudySession = async (tagId) => {
      try {
        const response = await startTagStudy(tagId, 20)
        const session = response.data.code === 'SUCCESS' ? response.data.data : response.data
        if (session.queue.length === 0) {
          ElMessage.info('该标签下没有卡片')
          return
        }
        studySession.value = session
        showAnswer.value = false
      } catch (error) {
        console.error('开始学习失败:', error)
        ElMessage.error('开始学习失败')
      }
    }

    // 开始随机学习
    const startRandomStudySession = async () => {
      if (randomLimit.value < 5 || randomLimit.value > 50) {
        ElMessage.warning('请输入5-50之间的数字')
        return
      }
      try {
        const response = await startRandomStudy(randomLimit.value)
        const session = response.data.code === 'SUCCESS' ? response.data.data : response.data
        if (session.queue.length === 0) {
          ElMessage.info('没有可学习的卡片')
          return
        }
        studySession.value = session
        showAnswer.value = false
      } catch (error) {
        console.error('开始随机学习失败:', error)
        ElMessage.error('开始随机学习失败')
      }
    }

    // 提交复习结果
    const submitReviewResult = async (result) => {
      if (!currentCard.value) return
      
      try {
        await submitReview(currentCard.value.card_id, result)
        
        // 移动到下一张卡片
        studySession.value.current++
        studySession.value.completed++
        showAnswer.value = false
        
        // 提示消息
        const messages = ['下次再来复习这张卡片', '继续加油！', '很好，继续保持！']
        ElMessage.success(messages[result])
        
      } catch (error) {
        console.error('提交复习结果失败:', error)
        ElMessage.error('提交复习结果失败')
      }
    }

    // 退出学习
    const exitStudy = () => {
      studySession.value = null
      showAnswer.value = false
      selectedDeckId.value = ''
    }

    // 跳转到标签卡片页面
    const goToTagCards = (tagId) => {
      if (tagId) {
        getCurrentInstance().proxy.$router.push({
          path: '/cards',
          query: { tag_id: tagId }
        })
      }
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
      
      // 配置 marked 选项
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
        // 如果 markdown 解析失败，回退到纯文本
        return `<p>${text.replace(/\n/g, '<br>')}</p>`
      }
    }

    onMounted(() => {
      fetchDecks()
      fetchTags()
      
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
      
      // 检查查询参数，自动开始学习
      const route = useRoute()
      if (route.query.deckId) {
        selectedDeckId.value = parseInt(route.query.deckId)
        // 等待数据加载完成后开始学习
        setTimeout(() => {
          startDeckStudySession()
        }, 500)
      } else if (route.query.tagId) {
        setTimeout(() => {
          startTagStudySession(parseInt(route.query.tagId))
        }, 500)
      }
    })

    return {
      studySession,
      showAnswer,
      decks,
      tags,
      selectedDeckId,
      randomLimit,
      currentCard,
      studyComplete,
      studyStats,
      startDueStudy,
      startDeckStudy: startDeckStudySession,
      startTagStudy: startTagStudySession,
      startRandomStudy: startRandomStudySession,
      submitReview: submitReviewResult,
      exitStudy,
      goToTagCards,
      renderMarkdown
    }
  }
}
</script>

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