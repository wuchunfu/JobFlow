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
      getRequest('/sys/task/list', params).then(response => {
        console.log(response)
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  // get host list
  getHostAllList({state}) {
    return new Promise((resolve, reject) => {
      getRequest('/sys/task/hostAllList').then(response => {
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
  // get host save or update
  getTableSaveOrUpdate({state}, payload) {
    return new Promise((resolve, reject) => {
      postRequest(`/sys/task/${!payload.taskId ? 'save' : 'update'}`, payload).then(response => {
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
  },
  // delete user row
  getRun({state}, payload) {
    return new Promise((resolve, reject) => {
      getRequest(`/sys/task/run/${payload.taskId}`).then(response => {
        // console.log(response)
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  // delete user row
  getStop({state}, payload) {
    return new Promise((resolve, reject) => {
      getRequest(`/sys/task/stop/${payload.taskId}`).then(response => {
        // console.log(response)
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  // delete user row
  getEnable({state}, payload) {
    return new Promise((resolve, reject) => {
      getRequest(`/sys/task/enable/${payload.taskId}`).then(response => {
        // console.log(response)
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  // delete user row
  getDisable({state}, payload) {
    return new Promise((resolve, reject) => {
      getRequest(`/sys/task/disable/${payload.taskId}`).then(response => {
        // console.log(response)
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  }
}
