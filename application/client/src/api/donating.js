import request from '@/utils/request'

export function queryDonatingList(data) {
  return request({
    url: '/queryDonatingList',
    method: 'post',
    data
  })
}

export function queryDonatingListByGrantee(data) {
  return request({
    url: '/queryDonatingListByGrantee',
    method: 'post',
    data
  })
}

export function updateDonating(data) {
  return request({
    url: '/updateDonating',
    method: 'post',
    data
  })
}

export function createDonating(data) {
  return request({
    url: '/createDonating',
    method: 'post',
    data
  })
}
