<template>
  <ElConfigProvider :locale="zhCn" :input="{ spellcheck: false }">
    <LoginModal ref="loginModalRef" />
    <ExpireModal ref="expireModalRef" />
    <RouterView />
  </ElConfigProvider>
</template>

<script setup lang="ts">
import { onMounted, ref, nextTick, watch } from 'vue'
import { ElConfigProvider, ElMessageBox } from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import settingApi from '@/api/modules/setting'

import LoginModal from './components/LoginModal/index.vue'
import ExpireModal from './components/ExpireModal/index.vue'
import { useUserStore } from '@/stores/modules/user'

import eventBus from '@/utils/event-bus'
import { shouldShowReminder, recordReminderShown, STORAGE_KEYS } from '@/utils/storage'
import { STORAGE_CONFIG } from '@/constants/storage'

import { EVENT_NAMES } from '@/constants/events'
import subscriptionApi from './api/modules/subscription'

const loginModalRef = ref<InstanceType<typeof LoginModal>>()
const expireModalRef = ref<InstanceType<typeof ExpireModal>>()

const userStore = useUserStore()

// 将插入样式和脚本的函数移到 onMounted 外部，提高可读性和复用性
const insertStyle = (content: string) => {
  if (!content) return

  const trimmed = content.trim()
  const type = 'style'
  const tag = `<${type}`
  try {
    if (trimmed.toLowerCase().startsWith(tag)) {
      document.head.insertAdjacentHTML('beforeend', trimmed)
    } else {
      const element = document.createElement(type)
      element.innerHTML = trimmed
      document.head.appendChild(element)
    }
  } catch (error) {
    console.error(`Failed to insert ${type} content:`, error)
  }
}

const insertScript = (content: string) => {
  if (!content) return

  const trimmed = content.trim()
  const node = document.createElement('div')
  node.innerHTML = trimmed
  const scripts = node.querySelectorAll('script')

  if (scripts.length) {
    scripts.forEach((script) => {
      const newScript = document.createElement('script')
      // Copy all attributes from the original script element
      Array.from(script.attributes).forEach((attr) => {
        newScript.setAttribute(attr.name, attr.value)
      })

      // If no src attribute, set type and content
      if (!script.src) {
        newScript.type = 'text/javascript'
        newScript.appendChild(document.createTextNode(script.innerHTML))
      }
      document.body.appendChild(newScript)
    })
  } else {
    const newScript = document.createElement('script')
    newScript.type = 'text/javascript'
    newScript.appendChild(document.createTextNode(content))
    document.body.appendChild(newScript)
  }
}

const checkSubscriptionExpire = async () => {
  await nextTick()

  // 检查是否需要显示过期提醒（使用通用方法）
  if (!shouldShowReminder(STORAGE_KEYS.EXPIRE_REMINDER, userStore.info.group_expire_time)) {
    return
  }

  const expireTime = userStore.info.group_expire_time
  const minExpireDay = STORAGE_CONFIG.MIN_EXPIRE_DAY

  // 情况: 已过期 | 即将过期（小于配置的最小天数）
  if (userStore.info.group_isexpired || userStore.info.group_expire_day < minExpireDay) {
    expireModalRef.value?.open({
      group_name: userStore.info.group_name,
      day: userStore.info.group_expire_day,
      expire_time: expireTime
    })
    recordReminderShown(STORAGE_KEYS.EXPIRE_REMINDER, expireTime)
  }
}

watch(
  () => userStore.is_login,
  () => {
    if (userStore.is_login) {
      setTimeout(() => {
        checkSubscriptionExpire()
      }, 1000)
    }
  },
  { immediate: true }
)

onMounted(() => {
  const search = new URLSearchParams(window.location.search)

  if (search.get('login_way') === 'wechat_login') {
    loginModalRef.value?.open({ way: 'wechat_login', openid: search.get('openid') || '', unionid: search.get('unionid') || '' })
  } else if (search.get('login_way') === 'wecom_login') {
    if (search.get('errorcode') === '-7') {
      ElMessageBox.confirm((window as any).$t('login.wecom_login_error'), (window as any).$t('common.tip'), {
        type: 'warning',
        showCancelButton: false
      })
    }
  }
  // 双层判断，防止支付宝支付成功后，页面刷新，导致搜索参数丢失
  if (search.get('method') === 'alipay.trade.page.pay.return' && search.get('out_trade_no')) {
    // 支付宝支付成功
    const order_id = search.get('out_trade_no')
    subscriptionApi.getOrderStatus({ order_id })
    // 删除支付宝支付成功后的搜索参数
    const delete_keys = [
      'charset',
      'method',
      'app_id',
      'utf-8',
      'out_trade_no',
      'sign',
      'total_amount',
      'trade_no',
      'auth_app_id',
      'version',
      'sign_type',
      'seller_id',
      'timestamp'
    ]
    delete_keys.forEach((key) => {
      search.delete(key)
    })
    window.location.search = search.toString()
  }

  settingApi.group
    .get('third_party_statistic')
    .then((res) => {
      const items = {
        style: res.data.find((item) => item.key === 'third_party_statistic_css')?.value || '',
        script: res.data.find((item) => item.key === 'third_party_statistic_header')?.value || ''
      }

      insertStyle(items.style)
      insertScript(items.script)
    })
    .catch((error) => {
      console.error('Failed to fetch third party statistic settings:', error)
    })

  // 保留原有的事件监听，用于处理登录成功后的情况
  eventBus.on(EVENT_NAMES.LOGIN_SUCCESS, () => {
    checkSubscriptionExpire()
  })
})
</script>
