import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/decks',
      name: 'decks',
      component: () => import('../views/DecksView.vue')
    },
    {
      path: '/cards',
      name: 'cards',
      component: () => import('../views/CardsView.vue')
    },
    {
      path: '/import-export',
      name: 'import-export',
      component: () => import('../views/ImportExportView.vue')
    },
    {
      path: '/decks/:id',
      name: 'deck-detail',
      component: () => import('../views/DeckDetailView.vue')
    },
    {
      path: '/cards/:id',
      name: 'card-detail',
      component: () => import('../views/CardDetailView.vue')
    },
    {
      path: '/tags',
      name: 'tags',
      component: () => import('../views/TagsView.vue')
    },
    {
      path: '/tags/:id/cards',
      name: 'tag-cards',
      component: () => import('../views/TagCardsView.vue')
    },
    {
      path: '/study',
      name: 'study',
      component: () => import('../views/StudyView.vue')
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('../views/SettingsView.vue')
    }
  ]
})

export default router