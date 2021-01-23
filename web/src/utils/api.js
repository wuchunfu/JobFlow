import request from '@/utils/request'

export const postRequest = (url, params) => {
  return request({
    method: 'post',
    url: url,
    data: params
  })
}
export const getRequest = (url, params) => {
  return request({
    method: 'get',
    url: url,
    params: params
  })
}
export const putRequest = (url, params) => {
  return request({
    method: 'put',
    url: url,
    data: params
  })
}
export const deleteRequest = (url) => {
  return request({
    method: 'delete',
    url: url
  })
}
export const uploadFileRequest = (url, params) => {
  return request({
    method: 'post',
    url: url,
    data: params,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
