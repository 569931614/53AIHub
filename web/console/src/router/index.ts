import type { App } from 'vue'
import type { RouteRecordRaw } from 'vue-router'
import { createRouter, createWebHashHistory } from 'vue-router'
import Layout from '@/layout/Layout.vue'
import { useEnterpriseStore, useUserStore } from '@/stores'
import { useDefaultUser } from '@/stores/modules/user'
import { useEnv } from '@/hooks/useEnv'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Main',
    component: Layout,
    redirect: () => '/index',
    children: [
      {
        path: '/index',
        name: 'Index',
        component: () => import('@/views/index.vue'),
      },
      // #ifndef KM
      {
        path: '/agent',
        name: 'Agent',
        component: () => import('@/views/agent/index.vue'),
      },
      {
        path: '/agent/create',
        name: 'AgentCreate',
        component: () => import('@/views/agent/create/index.vue'),
      },
      {
        path: '/prompt',
        name: 'Prompt',
        meta: {
          title: 'prompt.title',
        },
        component: () => import('@/views/prompt/index.vue'),
      },
      {
        path: '/prompt/create',
        name: 'PromptCreate',
        meta: {
          title: 'prompt.title',
        },
        component: () => import('@/views/prompt/create/index.vue'),
      },
      {
        path: '/toolbox',
        name: 'Toolbox',
        component: () => import('@/views/toolbox/index.vue'),
      },
      // #endif
      {
        path: '/order',
        name: 'Order',
        component: () => import('@/views/order/index.vue'),
      },
      {
        path: '/user/register',
        name: 'RegisterUser',
        meta: {
          title: 'register_user.title',
        },
        component: () => import('@/views/user/register/index.vue'),
      },
      {
        path: '/user/internal',
        name: 'InternalUser',
        meta: {
          title: 'internal_user.title',
        },
        component: () => import('@/views/user/internal/index.vue'),
      },
      {
        path: '/user/admin',
        name: 'AdminUser',
        meta: {
          title: 'admin_user.title',
        },
        component: () => import('@/views/user/admin/index.vue'),
      },
      {
        path: '/user/dialogue-record/:user_id',
        name: 'UserDialogueRecord',
        meta: {
          title: 'dialogue_record',
        },
        component: () => import('@/views/user/dialogue-record/index.vue'),
      },
      {
        path: '/info',
        name: 'Info',
        component: () => import('@/views/info/index.vue'),
      },
      {
        path: '/platform',
        name: 'Platform',
        // #ifdef KM
        component: () => import('@/views/platform/km.vue'),
        // #endif
        // #ifndef KM
        component: () => import('@/views/platform/index.vue'),
        // #endif
      },
      {
        path: '/navigation',
        name: 'Navigation',
        component: () => import('@/views/navigation/index.vue'),
      },
      {
        path: '/payment',
        name: 'Payment',
        component: () => import('@/views/payment/index.vue'),
      },
      {
        path: '/subscription',
        name: 'Subscription',
        component: () => import('@/views/subscription/index.vue'),
      },
      {
        path: '/domain',
        name: 'Domain',
        component: () => import('@/views/domain/index.vue'),
      },
      {
        path: '/statistics',
        name: 'Statistics',
        component: () => import('@/views/statistics/index.vue'),
      },
      {
        path: '/SMTP',
        name: 'SMTP',
        component: () => import('@/views/smtp/index.vue'),
      },
      {
        path: '/template-style',
        name: 'TemplateStyle',
        component: () => import('@/views/template-style/index.vue'),
      },
      {
        path: '/banner',
        name: 'Banner',
        component: () => import('@/views/banner/index.vue'),
      },
      {
        path: '/navigation',
        name: 'Navigation',
        meta: {
          title: 'navigation.title',
        },
        component: () => import('@/views/navigation/index.vue'),
      },
      {
        path: '/navigation/web-setting/:navigation_id',
        name: 'NavigationWebSetting',
        meta: {
          title: 'navigation.web_setting',
        },
        component: () => import('@/views/navigation/web-setting.vue'),
      },
      {
        path: '/sso',
        name: 'SSO',
        meta: {
          title: 'sso.title',
        },
        component: () => import('@/views/sso/index.vue'),
      },
      {
        path: '/system-log',
        name: 'SystemLog',
        component: () => import('@/views/system-log/index.vue'),
        meta: {
          title: 'module.system_log',
        },
      },
      // #ifdef KM
      {
        path: '/chunk',
        name: 'Chunk',
        component: () => import('@/views/chunk/index.vue'),
        meta: {
          title: 'module.chunk_setting',
        },
      },
      {
        path: '/model',
        name: 'Model',
        component: () => import('@/views/model/index.vue'),
        meta: {
          title: '模型设置',
        },
      },
      {
        path: '/space',
        name: 'Space',
        component: () => import('@/views/space/index.vue'),
        meta: {
          title: 'module.space',
        },
      },
      // #endif
    ],
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
  },
  {
    path: '/svglist',
    name: 'Svg',
    component: () => import('@/views/svg/index.vue'),
  },
  {
    path: '/500',
    name: 'Error500',
    component: () => import('@/views/exception/500/index.vue'),
  },
  {
    path: '/404',
    name: 'Error404',
    component: () => import('@/views/exception/404/index.vue'),
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    redirect: '/404',
  },
]

export const router = createRouter({
  // history: createWebHistory(base_path),
  history: createWebHashHistory(),
  routes,
  scrollBehavior: () => ({ left: 0, top: 0 }),
})

// 重写back方法，如果history.state.back为true，则调用原来的back方法，否则重定向到首页
const oldBack = router.back
router.back = () => {
  if (window.history.state.back) {
    oldBack()
  } else {
    router.push('/')
  }
}

export const gotoLogin = () => {
  const { isOpLocalEnv } = useEnv()
  // let login_url = sessionStorage.getItem('from_origin') || ''
  let login_url = ''
  if (!login_url) {
    login_url = `//${window.location.host}${window.location.search}`
    if (/(127.0.0.1)|(localhost)|(agenthubdev.cc)|(192.168.1.\d+)|/.test(login_url)) {
      // #ifndef KM
      login_url = `//${window.location.host}/console/saas-login/index.html${window.location.search}`
      // #endif

      // #ifdef KM
      login_url = `//${window.location.host}/console/km-login/index.html${window.location.search}`
      // #endif
    }
  }
  if (isOpLocalEnv.value) login_url = `${window.location.origin}/#/index`
  window.location.replace(login_url)
}
router.beforeEach(async (to, from, next) => {
  const user_store = useUserStore()
  const user_info = useDefaultUser()
  const enterprise_store = useEnterpriseStore()

  // const is_invalid_user = !user_info.access_token || !user_info.eid
  const is_invalid_user = !user_info.access_token
  if (is_invalid_user) user_store.logoff()
  if (!['/login', '/register'].includes(to.path) && is_invalid_user) {
    gotoLogin()
    return
  }
  // #ifndef KM
  if (['RegisterUser', 'InternalUser'].includes(to.name)) {
    if (!enterprise_store.info.eid) await enterprise_store.loadSelfInfo()
    if (
      (to.name === 'RegisterUser' &&
        !enterprise_store.info.is_independent &&
        !enterprise_store.info.is_industry) ||
      (to.name === 'InternalUser' &&
        !enterprise_store.info.is_enterprise &&
        !enterprise_store.info.is_industry)
    ) {
      next({ name: 'Error404' })
      return
    }
  }
  // #endif

  next()
})

export async function setupRouter(app: App) {
  app.use(router)
  await router.isReady()
}
