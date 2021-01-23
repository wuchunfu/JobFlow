import {getRequest, postRequest} from '@/utils/api'

export default {
  // get host list
  getList({state}, payload) {
    return new Promise((resolve, reject) => {
      const params = {
        'page': payload.pageIndex,
        'pageSize': payload.pageSize,
        'taskName': payload.taskName
      }
      getRequest('/sys/task/log/list', params).then(response => {
        // console.log(response)
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  // get host info
  getTableInfo({state}, payload) {
    return new Promise((resolve, reject) => {
      getRequest(`/sys/task/info/${payload.taskId}`).then(response => {
        // console.log(response)
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  // delete user row
  getDeleteRow({state}, payload) {
    return new Promise((resolve, reject) => {
      const params = {
        'taskIds': payload.rowIds
      }
      postRequest('/sys/task/delete', params).then(response => {
        // console.log(response)
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  }
}
