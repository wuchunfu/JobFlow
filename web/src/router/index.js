import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

const originalPush = Router.prototype.push
Router.prototype.push = function push(location, onResolve, onReject) {
  if (onResolve || onReject) {
    return originalPush.call(this, location, onResolve, onReject)
  }
  return originalPush.call(this, location).catch(err => err)
}

const originalReplace = Router.prototype.replace
Router.prototype.replace = function replace(location, onResolve, onReject) {
  if (onResolve || onReject) {
    return originalReplace.call(this, location, onResolve, onReject)
  }
  return originalReplace.call(this, location).catch(err => err)
}

/* Layout */
import Layout from '@/layout'

/**
 * Note: sub-menu only appear when route children.length >= 1
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'             the icon show in the sidebar
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes = [
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },
  {
    path: '/404',
    component: () => import('@/views/404'),
    hidden: true
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [{
      path: 'dashboard',
      name: '看板',
      component: () => import('@/views/dashboard/index'),
      meta: {title: '看板', icon: 'dashboard'}
    }]
  },
  {
    path: '/taskManager',
    component: Layout,
    redirect: '/taskManager/taskList',
    name: '任务管理',
    meta: {title: '任务管理', icon: 'peoples'},
    children: [
      {
        path: 'taskList',
        name: '任务列表',
        component: () => import('@/views/taskManager/list/index'),
        meta: {title: '任务列表', icon: 'user'}
      },
      {
        path: 'taskLog',
        name: '任务日志',
        component: () => import('@/views/taskLogManager/list/index'),
        meta: {title: '任务日志', icon: 'user'}
      }
    ]
  },
  {
    path: '/hostManager',
    component: Layout,
    redirect: '/hostManager/hostList',
    name: '主机管理',
    meta: {title: '主机管理', icon: 'peoples'},
    children: [
      {
        path: 'hostList',
        name: '主机列表',
        component: () => import('@/views/hostManager/list/index'),
        meta: {title: '主机列表', icon: 'user'}
      }
    ]
  },
  {
    path: '/userManager',
    component: Layout,
    redirect: '/userManager/userList',
    name: '用户管理',
    meta: {title: '用户管理', icon: 'peoples'},
    children: [
      {
        path: 'userList',
        name: '用户列表',
        component: () => import('@/views/userManager/list/index'),
        meta: {title: '用户列表', icon: 'user'}
      }
    ]
  },
  {
    path: '/systemManager',
    component: Layout,
    redirect: '/systemManager/loginLogList',
    name: '系统管理',
    meta: {title: '系统管理', icon: 'peoples'},
    children: [
      {
        path: 'loginLogList',
        name: '登陆日志列表',
        component: () => import('@/views/loginLogManager/list/index'),
        meta: {title: '登陆日志列表', icon: 'user'}
      },
      {
        path: 'systemLog',
        name: '系统日志',
        component: () => import('@/views/systemLog/index'),
        meta: {title: '系统日志', icon: 'user'}
      }
    ]
  },
  // 404 page must be placed at the end !!!
  {path: '*', redirect: '/404', hidden: true}
]

const createRouter = () => new Router({
  mode: 'history', // require service support
  scrollBehavior: () => ({y: 0}),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
