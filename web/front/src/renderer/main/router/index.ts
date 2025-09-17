import type { App } from 'vue'
import type { RouteRecordRaw } from 'vue-router'
import { createRouter, createWebHashHistory, createWebHistory } from 'vue-router'
import { checkPermission } from '@/utils/permission'
import useEnv from '@/hooks/useEnv'

import { useUserStore } from '@/stores/modules/user'

const { isOpLocalEnv } = useEnv()

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    redirect: () => {
      return window.innerWidth < 768 || window.electron ? '/agent' : '/index-redirect?from_home=1'
    },
    // redirect: '/index',
    component: () => import('@/layout/index.vue'),
    // component: () => import('@/views/index/agent/index.vue'),
    children: [
      {
        path: '/chat',
        name: 'Chat',
        component: () => import('@/views/chat/index.vue')
      },
      {
        path: '/agent',
        name: 'Agent',
        component: () => import('@/views/agent/index.vue')
      },
      {
        path: '/toolkit',
        name: 'Toolkit',
        component: () => import('@/views/toolkit/index.vue')
      },
      {
        path: '/profile',
        name: 'Profile',
        component: () => import('@/views/profile/index.vue'),
        meta: {
          pass: true,
          auth: !window.$isElectron
        }
      },
      {
        path: '/toolbox',
        name: 'Toolbox',
        component: () => import('@/views/desktop/tools/index.vue')
      },
      {
        path: '/prompt',
        name: 'Prompt',
        component: () => import('@/views/prompt/index.vue')
      },
      {
        path: '/prompt/:prompt_id',
        name: 'PromptDetail',
        component: () => import('@/views/prompt/detail/index.vue')
      },
      {
        path: '/order',
        name: 'Order',
        component: () => import('@/views/order/index.vue'),
        meta: {
          pass: true,
          auth: true
        }
      }
    ]
  },
  {
    path: '/index',
    name: 'Index',
    component: () => import('@/views/index/layout.vue'),
    children: [
      {
        path: '/index-redirect',
        name: 'HomeRedirect',
        component: () => import('@/views/index/redirect.vue')
      },
      {
        path: '/index',
        name: 'HomeIndex',
        component: () => import('@/views/index/index.vue')
      },
      {
        path: '/index/agent',
        name: 'HomeAgent',
        component: () => import('@/views/index/agent/index.vue')
      },
      {
        path: '/index/chat',
        name: 'HomeChat',
        component: () => import('@/views/index/agent/chat.vue')
      },
      {
        path: '/index/prompt',
        name: 'HomePrompt',
        component: () => import('@/views/index/prompt/index.vue')
      },
      {
        path: '/index/prompt/:prompt_id',
        name: 'HomePromptDetail',
        component: () => import('@/views/index/prompt/detail.vue')
      },
      {
        path: '/index/toolkit',
        name: 'HomeToolkit',
        component: () => import('@/views/index/toolkit.vue')
      }
    ]
  },

  {
    path: '/share/chat',
    name: 'ShareChat',
    component: () => import('@/views/share/chat.vue'),
    meta: {
      pass: true
    }
  },
  {
    path: '/guide',
    name: 'Guide',
    component: () => import('@/components/Lead/index.vue')
  },

  {
    path: '/svglist',
    name: 'Svg',
    component: () => import('@/views/svglist/index.vue')
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    redirect: '/404'
  },
  {
    path: '/index/:pathMatch(.*)*',
    name: 'IndexNotFound',
    redirect: () => {
      const redirect = window.location.hash.replace('#', '')
      return `/index-redirect?redirect=${redirect}`
    }
  }
]

export const isHashRouter = true
export const isHistoryRouter = false

export const router = createRouter({
  // electron环境 需要使用hash跳转
  history: isHashRouter ? createWebHashHistory() : createWebHistory(),
  routes,
  scrollBehavior: () => ({ left: 0, top: 0 })
})

// 添加路由守卫
router.beforeEach(async (to, _from, next) => {
  const isGuidePath = to.path === '/guide'
  if (isGuidePath && _from.path === '/') {
    next()
    return
  }

  if (isOpLocalEnv.value) {
    try {
      const { system } = await import('@/api/modules/system')
      const res = await system.init()
      if (!res.data && !isGuidePath) {
        next('/guide')
        return
      }
      if (res.data && isGuidePath) {
        next('/index')
        return
      }
    } catch (error) {
      console.log(error)
    }
  }

  const userStore = useUserStore()
  const isLoggedIn = localStorage.getItem('access_token') // 或其他判断用户是否登录的方法
  if (isLoggedIn) {
    userStore.getUserInfo()
  }

  const { useEnterpriseStore } = await import('@/stores/modules/enterprise')
  const enterpriseStore = useEnterpriseStore()
  if (!enterpriseStore.display_name) await enterpriseStore.loadInfo()
  const isWebsite = !enterpriseStore.isSoftStyle
  const isIndex = to.path.startsWith('/index')
  const isPass = to.meta.pass

  if (isPass || window.electron) {
  } else if (isIndex && !isWebsite) {
    next('/agent')
  } else if (!isIndex && isWebsite) {
    next('/')
  }
  if (to.meta.auth) {
    // 使用统一的认证检查
    if (!isLoggedIn) {
      await checkPermission()
      next('/') // 重定向到首页
      return
    }
  }
  next()
})

export async function setupRouter(app: App) {
  app.use(router)
  await router.isReady()
}
