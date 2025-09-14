import request from './request'

// 开始学习卡包
export function startDeckStudy(deckId, limit = 20) {
  return request({
    url: `/study/deck/${deckId}`,
    method: 'post',
    params: { limit }
  })
}

// 开始学习标签
export function startTagStudy(tagId, limit = 20) {
  return request({
    url: `/study/tag/${tagId}`,
    method: 'post',
    params: { limit }
  })
}

// 开始随机学习
export function startRandomStudy(limit = 10) {
  return request({
    url: '/study/random',
    method: 'post',
    params: { limit }
  })
}

// 获取到期复习卡片
export function getDueCards(limit = 20) {
  return request({
    url: '/study/due',
    method: 'get',
    params: { limit }
  })
}

// 提交复习结果
export function submitReview(cardId, result) {
  return request({
    url: `/study/review/${cardId}`,
    method: 'post',
    data: { result }
  })
}
