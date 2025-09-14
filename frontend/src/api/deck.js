import request from './request'

// 获取所有卡包
export function getDecks() {
  return request({
    url: '/decks',
    method: 'get'
  })
}

// 获取卡包详情
export function getDeck(id) {
  return request({
    url: `/decks/${id}`,
    method: 'get'
  })
}

// 创建卡包
export function createDeck(data) {
  return request({
    url: '/decks',
    method: 'post',
    data
  })
}

// 更新卡包
export function updateDeck(id, data) {
  return request({
    url: `/decks/${id}`,
    method: 'patch',
    data
  })
}

// 删除卡包
export function deleteDeck(id) {
  return request({
    url: `/decks/${id}`,
    method: 'delete'
  })
}

// 获取卡包统计信息
export function getDeckStats(id) {
  return request({
    url: `/decks/${id}/stats`,
    method: 'get'
  })
}