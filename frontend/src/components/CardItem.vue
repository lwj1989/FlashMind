<template>
  <div class="bg-white rounded-lg shadow p-4 mb-4 cursor-pointer hover:shadow-md transition-shadow duration-200" @click="viewCard">
    <div class="flex justify-between items-start mb-2">
      <div class="flex-1">
        <div class="text-sm text-gray-500 mb-1">{{ card.question }}</div>
        <div class="text-gray-800">{{ card.answer }}</div>
      </div>
      <div class="flex space-x-2 ml-4" @click.stop>
        <button @click="editCard" class="text-gray-600 hover:text-gray-800">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
            <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
          </svg>
        </button>
        <button @click="confirmDelete" class="text-red-600 hover:text-red-800">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
          </svg>
        </button>
      </div>
    </div>
    
    <div v-if="card.deck_name" class="mt-2 text-sm text-gray-600">
      卡包: {{ card.deck_name }}
    </div>
    
    <div v-if="card.tags && card.tags.length > 0" class="mt-2">
      <div class="flex flex-wrap gap-1">
        <span v-for="tag in card.tags" :key="tag.id" class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-indigo-100 text-indigo-800">
          {{ tag.name }}
        </span>
      </div>
    </div>
    
    <div class="mt-3 text-xs text-gray-500">
      创建于 {{ formatDate(card.created_at) }}
    </div>
  </div>
</template>

<script>
export default {
  name: 'CardItem',
  props: {
    card: {
      type: Object,
      required: true
    }
  },
  methods: {
    viewCard() {
      this.$router.push(`/cards/${this.card.id}`)
    },
    editCard() {
      this.$emit('edit', this.card)
    },
    confirmDelete() {
      this.$emit('delete', this.card)
    },
    formatDate(dateString) {
      const date = new Date(dateString)
      return date.toLocaleDateString()
    }
  }
}
</script>