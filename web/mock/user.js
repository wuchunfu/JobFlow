const tokens = {
  admin: {
    token: 'admin-token'
  }
}

const users = {
  'admin-token': {
    roles: ['admin'],
    introduction: 'I am a super administrator',
    avatar: '/avatar.jpg',
    name: 'Super Admin'
  }
}

export default [
  // user login
  {
    url: '/vue-admin/user/login',
    type: 'post',
    response: config => {
      const {username} = config.body
      const token = tokens[username]
      if (!token) {
        return {
          code: 1,
          data: null,
          msg: 'Wrong account or password.'
        }
      }
      return {
        code: 0,
        data: token,
        msg: 'Account login succeeded.'
      }
    }
  },
  // get user info
  {
    url: '/vue-admin/user/info\.*',
    type: 'get',
    response: config => {
      const {token} = config.query
      const info = users[token]
      if (!info) {
        return {
          code: 1,
          data: null,
          msg: 'Account login failed, unable to get user details.'
        }
      }
      return {
        code: 0,
        data: info,
        msg: 'Get user details succeeded.'
      }
    }
  },
  // user logout
  {
    url: '/vue-admin/user/logout',
    type: 'post',
    response: () => {
      return {
        code: 0,
        data: null,
        msg: 'Account exit succeeded.'
      }
    }
  }
]
