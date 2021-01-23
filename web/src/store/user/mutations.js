import {getToken} from '@/utils/auth'
import Cookies from 'js-cookie'

export default {
  SET_USER_ID: (state, userId) => {
    state.userId = userId
    Cookies.set('userId', userId)
  },
  RESET_STATE: (state) => {
    Object.assign(state, {
      token: getToken()
    })
  },
  SET_TOKEN: (state, token) => {
    state.token = token
  }
}
