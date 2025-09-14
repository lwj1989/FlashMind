import request from './request'

// 创建卡片
export function createCard(data) {
  return request({
    url: '/cards',
    method: 'post',
    data
  })
}

// 获取卡片详情
export function getCard(id) {
  return request({
    url: `/cards/${id}`,
    method: 'get'
  })
}

// 更新卡片
export function updateCard(id, data) {
  return request({
    url: `/cards/${id}`,
    method: 'patch',
    data
  })
}

// 删除卡片
export function deleteCard(id) {
  return request({
    url: `/cards/${id}`,
    method: 'delete'
  })
}

// 获取卡包下的卡片
export function getCardsByDeck(deckId, page = 1, pageSize = 20) {
  return request({
    url: `/decks/${deckId}/cards`,
    method: 'get',
    params: {
      page,
      page_size: pageSize
    }
  })
}

// 获取标签下的卡片
export function getCardsByTag(tagId, page = 1, pageSize = 20) {
  return request({
    url: `/tags/${tagId}/cards`,
    method: 'get',
    params: {
      page,
      page_size: pageSize
    }
  })
}

// 搜索卡片
export function searchCards(params) {
  return request({
    url: '/cards',
    method: 'get',
    params
  })
}