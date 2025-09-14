<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <nav class="bg-white shadow-md">
      <div class="max-w-6xl mx-auto px-4">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <router-link to="/" class="flex items-center space-x-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
              </svg>
              <span class="text-xl font-bold text-gray-800">FlashMind</span>
            </router-link>
          </div>
          <div class="flex items-center space-x-4">
            <router-link 
              to="/" 
              class="px-3 py-2 rounded-md text-sm font-medium text-gray-700 hover:text-indigo-600 hover:bg-gray-50"
              active-class="text-indigo-600 bg-gray-50"
            >
              首页
            </router-link>
            <router-link 
              to="/study" 
              class="px-3 py-2 rounded-md text-sm font-medium text-gray-700 hover:text-indigo-600 hover:bg-gray-50"
              active-class="text-indigo-600 bg-gray-50"
            >
              开始学习
            </router-link>
            <router-link 
              to="/decks" 
              class="px-3 py-2 rounded-md text-sm font-medium text-gray-700 hover:text-indigo-600 hover:bg-gray-50"
              active-class="text-indigo-600 bg-gray-50"
            >
              卡包管理
            </router-link>
            <router-link 
              to="/tags" 
              class="px-3 py-2 rounded-md text-sm font-medium text-gray-700 hover:text-indigo-600 hover:bg-gray-50"
              active-class="text-indigo-600 bg-gray-50"
            >
              标签管理
            </router-link>
            <router-link 
              to="/cards" 
              class="px-3 py-2 rounded-md text-sm font-medium text-gray-700 hover:text-indigo-600 hover:bg-gray-50"
              active-class="text-indigo-600 bg-gray-50"
            >
              卡片管理
            </router-link>
            <router-link 
              to="/import-export" 
              class="px-3 py-2 rounded-md text-sm font-medium text-gray-700 hover:text-indigo-600 hover:bg-gray-50"
              active-class="text-indigo-600 bg-gray-50"
            >
              导入导出
            </router-link>
            <router-link 
              to="/settings" 
              class="px-3 py-2 rounded-md text-sm font-medium text-gray-700 hover:text-indigo-600 hover:bg-gray-50"
              active-class="text-indigo-600 bg-gray-50"
            >
              设置
            </router-link>
          </div>
        </div>
      </div>
    </nav>

    <main class="flex-grow">
      <router-view />
    </main>

    <footer class="bg-white border-t border-gray-200 py-6">
      <div class="max-w-6xl mx-auto px-4">
        <p class="text-center text-gray-500 text-sm">
          © {{ new Date().getFullYear() }} FlashMind - 智能卡片学习系统
        </p>
      </div>
    </footer>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>
