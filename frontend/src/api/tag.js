import request from './request'

// 获取所有标签
export function getTags() {
  return request({
    url: '/tags',
    method: 'get'
  })
}

// 获取标签详情
export function getTag(id) {
  return request({
    url: `/tags/${id}`,
    method: 'get'
  })
}

// 创建标签
export function createTag(data) {
  return request({
    url: '/tags',
    method: 'post',
    data
  })
}

// 获取所有标签
export function getAllTags() {
  return request({
    url: '/tags?include_stats=true',
    method: 'get'
  })
}

// 更新标签
export function updateTag(id, data) {
  return request({
    url: `/tags/${id}`,
    method: 'patch',
    data
  })
}

// 删除标签
export function deleteTag(id) {
  return request({
    url: `/tags/${id}`,
    method: 'delete'
  })
}

// 获取卡包下的标签
export function getTagsByDeck(deckId) {
  return request({
    url: `/tags/deck/${deckId}`,
    method: 'get'
  })
}

// 获取标签统计信息
export function getTagStats(id) {
  return request({
    url: `/tags/${id}/stats`,
    method: 'get'
  })
}