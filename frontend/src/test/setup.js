import { vi } from 'vitest'
import { config } from '@vue/test-utils'

// 全局模拟
vi.mock('axios', () => ({
  default: {
    create: () => ({
      get: vi.fn(),
      post: vi.fn(),
      patch: vi.fn(),
      delete: vi.fn(),
      interceptors: {
        request: { use: vi.fn() },
        response: { use: vi.fn() }
      }
    })
  }
}))

// 全局配置
config.global.stubs = {
  'router-link': true,
  'router-view': true
}