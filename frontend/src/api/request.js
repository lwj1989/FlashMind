import axios from 'axios'

// 创建axios实例
const request = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  timeout: 5000
})

// 请求拦截器
request.interceptors.request.use(
  config => {
    // 在发送请求之前做些什么
    return config
  },
  error => {
    // 对请求错误做些什么
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  response => {
    // 对响应数据做点什么
    return response
  },
  error => {
    // 对响应错误做点什么
    console.error('响应错误:', error)
    
    // 如果是HTTP错误，提取错误信息
    if (error.response) {
      const { status, data } = error.response
      console.error(`HTTP错误: ${status}`, data)
      
      // 如果是404错误，返回更友好的错误信息
      if (status === 404) {
        error.message = '请求的资源不存在'
      }
    } else if (error.request) {
      // 请求已发出但没有收到响应
      error.message = '服务器无响应，请检查网络连接'
    } else {
      // 请求配置出错
      error.message = '请求配置错误'
    }
    
    return Promise.reject(error)
  }
)

export default request