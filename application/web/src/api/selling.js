import request from '@/utils/request'

export function querySellingList(data) {
  return request({
    url: '/querySellingList',
    method: 'post',
    data
  })
}

export function querySellingListByBuyer(data) {
  return request({
    url: '/querySellingListByBuyer',
    method: 'post',
    data
  })
}

export function createSellingByBuy(data) {
  return request({
    url: '/createSellingByBuy',
    method: 'post',
    data
  })
}

export function updateSelling(data) {
  return request({
    url: '/updateSelling',
    method: 'post',
    data
  })
}

export function createSelling(data) {
  return request({
    url: '/createSelling',
    method: 'post',
    data
  })
}
