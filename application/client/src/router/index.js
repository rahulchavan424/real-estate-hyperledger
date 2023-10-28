import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/**
 * Note: Sub-menus only appear when route children.length >= 1.
 * For details, see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true - If set to true, the item will not be shown in the sidebar (default is false).
 * alwaysShow: true - If set to true, the root menu will always be shown.
 *                      If not set alwaysShow, when an item has more than one children route,
 *                      it will be in nested mode; otherwise, the root menu will not be shown.
 * redirect: noRedirect - If set to noRedirect, there will be no redirect in the breadcrumb.
 * name: 'router-name' - The name is used by <keep-alive> (must be set!!!)
 * meta: {
 *   roles: ['admin', 'editor'] - Control the page roles (you can set multiple roles).
 *   title: 'title' - The name shown in the sidebar and breadcrumb (recommended to set).
 *   icon: 'svg-name' - The icon shown in the sidebar.
 *   breadcrumb: false - If set to false, the item will be hidden in the breadcrumb (default is true).
 *   activeMenu: '/example/list' - If set to a path, the sidebar will highlight the path you set.
 * }
 */

/**
 * constantRoutes - A base page that does not have permission requirements.
 * All roles can access these routes.
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
    redirect: '/realestate',
    children: [
      {
        path: 'realestate',
        name: 'Realestate',
        component: () => import('@/views/realestate/list/index'),
        meta: {
          title: 'Real Estate Information',
          icon: 'realestate'
        }
      }
    ]
  }
]

/**
 * asyncRoutes - Routes that need to be dynamically loaded based on user roles.
 */
export const asyncRoutes = [
  {
    path: '/selling',
    component: Layout,
    redirect: '/selling/all',
    name: 'Selling',
    alwaysShow: true,
    meta: {
      title: 'Sales',
      icon: 'selling'
    },
    children: [
      {
        path: 'all',
        name: 'SellingAll',
        component: () => import('@/views/selling/all/index'),
        meta: {
          title: 'All Sales',
          icon: 'sellingAll'
        }
      },
      {
        path: 'me',
        name: 'SellingMe',
        component: () => import('@/views/selling/me/index'),
        meta: {
          roles: ['editor'],
          title: 'Initiated by Me',
          icon: 'sellingMe'
        }
      },
      {
        path: 'buy',
        name: 'SellingBuy',
        component: () => import('@/views/selling/buy/index'),
        meta: {
          roles: ['editor'],
          title: 'Purchased by Me',
          icon: 'sellingBuy'
        }
      }
    ]
  },
  {
    path: '/donating',
    component: Layout,
    redirect: '/donating/all',
    name: 'Donating',
    alwaysShow: true,
    meta: {
      title: 'Donations',
      icon: 'donating'
    },
    children: [
      {
        path: 'all',
        name: 'DonatingAll',
        component: () => import('@/views/donating/all/index'),
        meta: {
          title: 'All Donations',
          icon: 'donatingAll'
        }
      },
      {
        path: 'donor',
        name: 'DonatingDonor',
        component: () => import('@/views/donating/donor/index'),
        meta: {
          roles: ['editor'],
          title: 'Initiated by Me',
          icon: 'donatingDonor'
        }
      },
      {
        path: 'grantee',
        name: 'DonatingGrantee',
        component: () => import('@/views/donating/grantee/index'),
        meta: {
          roles: ['editor'],
          title: 'Received by Me',
          icon: 'donatingGrantee'
        }
      }
    ]
  },
  {
    path: '/addRealestate',
    component: Layout,
    meta: {
      roles: ['admin']
    },
    children: [
      {
        path: '/addRealestate',
        name: 'AddRealestate',
        component: () => import('@/views/realestate/add/index'),
        meta: {
          title: 'Add Real Estate',
          icon: 'addRealestate'
        }
      }
    ]
  },
  {
    path: '*',
    redirect: '/404',
    hidden: true
  }
]

const createRouter = () => new Router({
  base: '/',
  // mode: 'history', // require service support
  scrollBehavior: () => ({
    y: 0
  }),
  routes: constantRoutes
})

const router = createRouter()

// Details: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // Reset the router
}

export default router