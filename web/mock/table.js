import Mock from 'mockjs'

const listData = Mock.mock({
  'items|10': [{
    id: '@id',
    rowId: '@id',
    title: '@sentence(10, 20)',
    author: 'admin',
    pageViews: '@integer(300, 5000)',
    'status|1': ['published', 'draft', 'deleted'],
    createTime: '@datetime'
  }]
})

const updateData = Mock.mock({
  'items|1': [{
    id: '@id',
    rowId: '@id',
    title: '@sentence(10, 20)',
    author: 'admin',
    pageViews: '@integer(300, 5000)',
    'status|1': ['published', 'draft', 'deleted'],
    createTime: '@datetime'
  }]
})

export default [
  {
    url: '/vue-admin/table/list',
    type: 'get',
    response: config => {
      const {page, limit, title} = config.query
      console.log(page, limit, title)
      const items = listData.items
      return {
        code: 0,
        data: {
          total: items.length,
          items: items
        },
        msg: 'Get data successful.'
      }
    }
  },
  {
    url: '/vue-admin/table/info',
    type: 'get',
    response: () => {
      const items = updateData.items
      return {
        code: 0,
        data: {
          total: items.length,
          items: items
        },
        msg: 'Get data successful.'
      }
    }
  },
  {
    url: '/vue-admin/table/save',
    type: 'post',
    response: config => {
      const items = config.body
      return {
        code: 0,
        data: {
          total: items.length,
          items: items
        },
        msg: 'Save data successful.'
      }
    }
  },
  {
    url: '/vue-admin/table/update',
    type: 'post',
    response: config => {
      const items = config.body
      return {
        code: 0,
        data: {
          total: items.length,
          items: items
        },
        msg: 'Update data successful.'
      }
    }
  },
  {
    url: '/vue-admin/table/delete',
    type: 'post',
    response: config => {
      const items = config.body
      return {
        code: 0,
        data: {
          total: items.length,
          items: items
        },
        msg: 'Delete data successful.'
      }
    }
  }
]
