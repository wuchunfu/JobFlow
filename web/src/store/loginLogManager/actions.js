import {getRequest} from '@/utils/api'

export default {
  // get login log list
  getList({state}, payload) {
    return new Promise((resolve, reject) => {
      const params = {
        'page': payload.pageIndex,
        'pageSize': payload.pageSize,
        'username': payload.username
      }
      getRequest('/sys/login/log/list', params).then(response => {
        // console.log(response)
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  }
}
