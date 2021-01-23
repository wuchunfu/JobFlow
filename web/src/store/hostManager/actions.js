import {getRequest, postRequest} from '@/utils/api'

export default {
  // get host list
  getList({state}, payload) {
    return new Promise((resolve, reject) => {
      const params = {
        'page': payload.pageIndex,
        'pageSize': payload.pageSize,
        'hostName': payload.hostName
      }
      getRequest('/sys/host/list', params).then(response => {
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
      getRequest(`/sys/host/info/${payload.hostId}`).then(response => {
        // console.log(response)
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  // get host save or update
  getTableSaveOrUpdate({state}, payload) {
    return new Promise((resolve, reject) => {
      postRequest(`/sys/host/${!payload.hostId ? 'save' : 'update'}`, payload).then(response => {
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
        'hostIds': payload.rowIds
      }
      postRequest('/sys/host/delete', params).then(response => {
        // console.log(response)
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  // delete user row
  getConnectTest({state}, payload) {
    return new Promise((resolve, reject) => {
      getRequest(`/sys/host/ping/${payload.hostId}`).then(response => {
        // console.log(response)
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  }
}
