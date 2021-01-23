import {removeToken, setToken} from '@/utils/auth'
import {resetRouter} from '@/router'
import {postRequest} from '@/utils/api'

export default {
  // user login
  login({commit, state}, payload) {
    return new Promise((resolve, reject) => {
      const params = {
        username: payload.username.trim(),
        password: payload.password.trim()
      }
      postRequest('/sys/login', params).then(response => {
        // console.log(response)
        const {data} = response.data
        commit('SET_USER_ID', data.userId)
        commit('SET_TOKEN', data.token)
        setToken(data.token)
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  // user logout
  logout({commit}) {
    return new Promise((resolve, reject) => {
      postRequest('/sys/logout').then(response => {
        // console.log(response)
        removeToken() // must remove  token  first
        resetRouter()
        commit('RESET_STATE')
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  // remove token
  resetToken({commit}) {
    return new Promise(resolve => {
      removeToken() // must remove  token  first
      resetRouter()
      commit('RESET_STATE')
      resolve()
    })
  }
}
