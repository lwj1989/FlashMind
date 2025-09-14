<template>
  <form @submit.prevent="handleSubmit">
    <div class="mb-4">
      <label for="deck" class="block text-sm font-medium text-gray-700 mb-1">卡包</label>
      <select id="deck" v-model="card.deck_id" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500" required>
        <option value="">请选择卡包</option>
        <option v-for="deck in decks" :key="deck.id" :value="deck.id">
          {{ deck.name }}
        </option>
      </select>
    </div>

    <div class="mb-4">
      <label for="tag" class="block text-sm font-medium text-gray-700 mb-1">标签（可选）</label>
      <select id="tag" v-model="card.tag_id" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500">
        <option value="">无标签</option>
        <option v-for="tag in filteredTags" :key="tag.id" :value="tag.id">
          {{ tag.name }}
        </option>
      </select>
    </div>

    <div class="mb-4">
      <label for="question" class="block text-sm font-medium text-gray-700 mb-1">问题</label>
      <textarea id="question" v-model="card.question" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500" required rows="3"></textarea>
    </div>

    <div class="mb-6">
      <label for="answer" class="block text-sm font-medium text-gray-700 mb-1">答案</label>
      <textarea id="answer" v-model="card.answer" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500" required rows="3"></textarea>
    </div>

    <div class="flex justify-end space-x-3">
      <button type="button" @click="$emit('cancel')" class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50">
        取消
      </button>
      <button type="submit" class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700">
        {{ isEdit ? '更新' : '创建' }}
      </button>
    </div>
  </form>
</template>

<script>
import { ref, computed, watch } from 'vue'
import { getDecks } from '@/api/deck'
import { getTagsByDeck } from '@/api/tag'

export default {
  name: 'CardForm',
  props: {
    initialCard: {
      type: Object,
      default: () => ({
        deck_id: '',
        tag_id: null,
        question: '',
        answer: ''
      })
    },
    isEdit: {
      type: Boolean,
      default: false
    }
  },
  emits: ['submit', 'cancel'],
  setup(props, { emit }) {
    const card = ref({ ...props.initialCard })
    const decks = ref([])
    const tags = ref([])

    // 获取所有卡包
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
        console.error('获取卡包失败:', error)
      }
    }

    // 获取标签
    const fetchTags = async (deckId) => {
      if (!deckId) {
        tags.value = []
        return
      }
      try {
        const response = await getTagsByDeck(deckId)
        console.log('获取标签响应:', response)
        if (response.data && response.data.code === 'SUCCESS') {
          tags.value = response.data.data || []
        } else {
          tags.value = response.data || []
        }
      } catch (error) {
        console.error('获取标签失败:', error)
      }
    }

    // 过滤标签
    const filteredTags = computed(() => {
      return tags.value
    })

    // 监听卡包变化，获取对应标签
    watch(() => card.value.deck_id, (newDeckId) => {
      if (newDeckId) {
        fetchTags(newDeckId)
      } else {
        tags.value = []
        card.value.tag_id = null
      }
    })

    // 提交表单
    const handleSubmit = () => {
      emit('submit', { ...card.value })
    }

    // 初始化
    fetchDecks()
    if (card.value.deck_id) {
      fetchTags(card.value.deck_id)
    }

    return {
      card,
      decks,
      tags,
      filteredTags,
      handleSubmit
    }
  }
}
</script>