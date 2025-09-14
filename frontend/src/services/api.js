import axios from 'axios'

const apiClient = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  headers: {
    'Content-Type': 'application/json'
  }
})

export default {
  // 卡包相关API
  getDecks(includeStats = false) {
    return apiClient.get(`/decks${includeStats ? '?include_stats=true' : ''}`)
  },
  
  getDeck(id) {
    return apiClient.get(`/decks/${id}`)
  },
  
  createDeck(deck) {
    return apiClient.post('/decks', deck)
  },
  
  updateDeck(id, deck) {
    return apiClient.patch(`/decks/${id}`, deck)
  },
  
  deleteDeck(id) {
    return apiClient.delete(`/decks/${id}`)
  },
  
  getDeckStats(id) {
    return apiClient.get(`/decks/${id}/stats`)
  },
  
  // 标签相关API
  getAllTags() {
    return apiClient.get('/tags')
  },
  
  getTags(deckId) {
    return apiClient.get(`/tags/deck/${deckId}`)
  },
  
  getTag(id) {
    return apiClient.get(`/tags/${id}`)
  },
  
  createTag(tag) {
    return apiClient.post('/tags', tag)
  },
  
  updateTag(id, tag) {
    return apiClient.patch(`/tags/${id}`, tag)
  },
  
  deleteTag(id) {
    return apiClient.delete(`/tags/${id}`)
  },
  
  getTagStats(id) {
    return apiClient.get(`/tags/${id}/stats`)
  },
  
  // 卡片相关API
  getCards(deckId) {
    return apiClient.get(`/decks/${deckId}/cards`)
  },
  
  getCard(id) {
    return apiClient.get(`/cards/${id}`)
  },
  
  createCard(card) {
    return apiClient.post('/cards', card)
  },
  
  updateCard(id, card) {
    return apiClient.patch(`/cards/${id}`, card)
  },
  
  deleteCard(id) {
    return apiClient.delete(`/cards/${id}`)
  },
  
  getCardsByTag(tagId) {
    return apiClient.get(`/cards/tag/${tagId}`)
  }
}