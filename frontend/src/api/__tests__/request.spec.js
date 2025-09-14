import { describe, it, expect, vi, beforeEach } from 'vitest'
import axios from 'axios'
import request from '@/api/request'

// 模拟axios
vi.mock('axios', () => ({
  create: vi.fn(() => ({
    get: vi.fn(),
    post: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn(),
    interceptors: {
      request: { use: vi.fn() },
      response: { use: vi.fn() }
    }
  }))
}))

describe('Request Configuration', () => {
  let mockAxiosInstance

  beforeEach(() => {
    mockAxiosInstance = {
      get: vi.fn(),
      post: vi.fn(),
      patch: vi.fn(),
      delete: vi.fn(),
      interceptors: {
        request: { use: vi.fn() },
        response: { use: vi.fn() }
      }
    }

    axios.create.mockReturnValue(mockAxiosInstance)
  })

  it('应该正确配置axios实例', () => {
    // 重新导入request模块以触发axios.create
    vi.resetModules()
    const request = require('@/api/request').default

    expect(axios.create).toHaveBeenCalledWith({
      baseURL: 'http://localhost:8080',
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json'
      }
    })
  })

  it('应该正确配置请求拦截器', () => {
    // 重新导入request模块以触发拦截器配置
    vi.resetModules()
    const request = require('@/api/request').default

    expect(mockAxiosInstance.interceptors.request.use).toHaveBeenCalled()
  })

  it('应该正确配置响应拦截器', () => {
    // 重新导入request模块以触发拦截器配置
    vi.resetModules()
    const request = require('@/api/request').default

    expect(mockAxiosInstance.interceptors.response.use).toHaveBeenCalled()
  })

  it('请求拦截器应该正确添加token', () => {
    // 模拟localStorage中的token
    const mockToken = 'test-token'
    vi.spyOn(localStorage.__proto__, 'getItem').mockReturnValue(mockToken)

    // 重新导入request模块以触发拦截器配置
    vi.resetModules()
    const request = require('@/api/request').default

    // 获取请求拦截器的回调函数
    const requestInterceptor = mockAxiosInstance.interceptors.request.use.mock.calls[0][0]

    // 测试请求拦截器
    const config = { headers: {} }
    const result = requestInterceptor(config)

    expect(result.headers.Authorization).toBe(`Bearer ${mockToken}`)
  })

  it('请求拦截器应该正确处理没有token的情况', () => {
    // 模拟localStorage中没有token
    vi.spyOn(localStorage.__proto__, 'getItem').mockReturnValue(null)

    // 重新导入request模块以触发拦截器配置
    vi.resetModules()
    const request = require('@/api/request').default

    // 获取请求拦截器的回调函数
    const requestInterceptor = mockAxiosInstance.interceptors.request.use.mock.calls[0][0]

    // 测试请求拦截器
    const config = { headers: {} }
    const result = requestInterceptor(config)

    expect(result.headers.Authorization).toBeUndefined()
  })

  it('响应拦截器应该正确处理成功响应', () => {
    // 重新导入request模块以触发拦截器配置
    vi.resetModules()
    const request = require('@/api/request').default

    // 获取响应拦截器的成功回调函数
    const successInterceptor = mockAxiosInstance.interceptors.response.use.mock.calls[0][0]

    // 测试响应拦截器
    const response = { data: { success: true } }
    const result = successInterceptor(response)

    expect(result).toEqual(response)
  })

  it('响应拦截器应该正确处理错误响应', () => {
    // 重新导入request模块以触发拦截器配置
    vi.resetModules()
    const request = require('@/api/request').default

    // 获取响应拦截器的错误回调函数
    const errorInterceptor = mockAxiosInstance.interceptors.response.use.mock.calls[0][1]

    // 测试响应拦截器
    const error = new Error('网络错误')
    const result = errorInterceptor(error)

    expect(Promise.reject).toHaveBeenCalledWith(error)
  })

  it('响应拦截器应该正确处理401错误', () => {
    // 重新导入request模块以触发拦截器配置
    vi.resetModules()
    const request = require('@/api/request').default

    // 获取响应拦截器的错误回调函数
    const errorInterceptor = mockAxiosInstance.interceptors.response.use.mock.calls[0][1]

    // 模拟401错误
    const error = {
      response: {
        status: 401
      }
    }

    // 模拟window.location
    const mockLocation = { href: '' }
    Object.defineProperty(window, 'location', {
      value: mockLocation,
      writable: true
    })

    // 测试响应拦截器
    errorInterceptor(error)

    expect(mockLocation.href).toBe('/login')
  })
})