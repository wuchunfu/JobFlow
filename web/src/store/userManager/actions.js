import {getRequest, postRequest} from '@/utils/api'

export default {
  // get user list
  getList({state}, payload) {
    return new Promise((resolve, reject) => {
      const params = {
        'page': payload.pageIndex,
        'pageSize': payload.pageSize,
        'username': payload.username
      }
      getRequest('/sys/user/list', params).then(response => {
        // console.log(response)
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  // get user info
  getTableInfo({state}, payload) {
    return new Promise((resolve, reject) => {
      getRequest(`/sys/user/info/${payload.userId}`).then(response => {
        // console.log(response)
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  // get user save or update
  getTableSaveOrUpdate({state}, payload) {
    return new Promise((resolve, reject) => {
      postRequest(`/sys/user/${!payload.userId ? 'save' : 'update'}`, payload).then(response => {
        // console.log(response)
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  // get user change password
  getTableChangePassword({state}, payload) {
    return new Promise((resolve, reject) => {
      postRequest('/sys/user/changePassword', payload).then(response => {
        // console.log(response)
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  // get user change login password
  getChangeLoginPassword({state}, payload) {
    return new Promise((resolve, reject) => {
      postRequest('/sys/user/changeLoginPassword', payload).then(response => {
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
        'userIds': payload.rowIds
      }
      postRequest('/sys/user/delete', params).then(response => {
        // console.log(response)
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  }
}
