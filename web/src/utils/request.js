import axios from 'axios'
import {Message} from 'element-ui'

// create an axios instance
const service = axios.create({
  baseURL: process.env.VUE_APP_BASE_API, // url = base url + request url
  withCredentials: true, // send cookies when cross-domain requests
  timeout: 5000, // request timeout
  headers: {
    'Content-Type': 'application/json; charset=utf-8'
  }
})

// request interceptor
service.interceptors.request.use(request => {
    return request
  }, error => {
    // do something with request error
    console.log(error) // for debug
    return Promise.reject(error)
  }
)

// response interceptor
service.interceptors.response.use(response => {
    const res = response.data
    if (res !== '' && res !== undefined && res.code === 401) {
      Message({
        message: res.msg || 'Invalid token, 401',
        type: 'error',
        duration: 5 * 1000
      })
    } else {
      return response
    }
  }, error => {
    Message({
      message: error.msg || 'Refused to connect',
      type: 'error',
      duration: 5 * 1000
    })
    return Promise.reject(error)
  }
)

export default service
