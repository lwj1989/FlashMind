import request from './request'

// 获取系统统计信息
export function getSystemStats() {
  return request({
    url: '/system/stats',
    method: 'get'
  })
}

// 备份所有数据
export function backupData() {
  return request({
    url: '/system/backup',
    method: 'get',
    responseType: 'blob' // 重要：指定响应类型为blob
  })
}

// 恢复数据
export function restoreData(file, clearExisting = true) {
  const formData = new FormData()
  formData.append('file', file)
  
  return request({
    url: `/system/restore?clear_existing=${clearExisting}`,
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 清空所有数据
export function clearAllData() {
  return request({
    url: '/system/clear',
    method: 'delete'
  })
}
