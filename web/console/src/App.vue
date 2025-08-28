<template>
  <ElConfigProvider :locale="locale === 'en' ? en : locale === 'ja' ? ja : locale === 'zh-tw' ? zhTw : zhCn">
    <RouterView />
    <UserLoginDialog />
  </ElConfigProvider>
</template>

<script setup lang="ts">
import { ElConfigProvider } from 'element-plus'
import { useI18n } from 'vue-i18n'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import zhTw from 'element-plus/es/locale/lang/zh-tw'
import en from 'element-plus/es/locale/lang/en'
import ja from 'element-plus/es/locale/lang/ja'

import { onMounted } from 'vue'
import eventBus from '@/utils/event-bus'
import { gotoLogin } from '@/router'
// import { useRouter } from 'vue-router'
import { useEnterpriseStore, useUserStore } from '@/stores'

import settingApi from '@/api/modules/setting'
import { useEnv } from '@/hooks/useEnv'

const { isOpLocalEnv } = useEnv()
const { locale } = useI18n()
// const router = useRouter()
const enterprise_store = useEnterpriseStore()
const user_store = useUserStore()

eventBus.on('user-login-expired', async () => {
  await user_store.logoff()
  gotoLogin()
})
enterprise_store.loadSelfInfo()
user_store.loadSelfInfo()

const insertScript = (content: string) => {
  if (!content) return

  const trimmed = content.trim()
  const node = document.createElement('div')
  node.innerHTML = trimmed
  const scripts = node.querySelectorAll('script')

  if (scripts.length) {
    scripts.forEach(script => {
      const newScript = document.createElement('script')
      // Copy all attributes from the original script element
      Array.from(script.attributes).forEach(attr => {
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

onMounted(() => {
  if (isOpLocalEnv.value) {
    settingApi.detail('third_party_statistic').then(res => {
      const items = {
        script: res.data.find(item => item.key === 'third_party_statistic_header')?.value || '',
      }
      insertScript(items.script)
    })
  }
})
</script>

<style>
/* 临时添加，后续hub-ui-x 要支持 */
.markdown-body img {
  display: inline;
}
</style>
