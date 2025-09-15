import request from './request'

// 导入卡包
export function importDeck(data) {
  return request({
    url: '/import-export/decks',
    method: 'post',
    data,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 导出卡包
export function exportDeck(deckId, format) {
  return request({
    url: `/import-export/decks/${deckId}`,
    method: 'get',
    params: {
      format,
      download: 'true'
    },
    responseType: 'blob'
  })
}