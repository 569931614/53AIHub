<template>
  <div v-loading="loading" class="w-full">
    <div v-if="isInMobile" :style="{ height, width }" class="flex flex-col justify-center items-center">
      <div class="w-[220px] h-[220px] border relative rounded-lg overflow-hidden">
        <img :src="$getPublicPath('/images/login/wecom_login.png')" alt="" />
        <div class="absolute inset-0 bg-white bg-opacity-90 flex justify-center items-center">
          <ElButton type="danger" @click="handleLogin"> {{ $t('login.immediate_login') }} </ElButton>
        </div>
      </div>
      <p class="text-sm text-regular mt-3 text-opacity-60">{{ $t('login.login_by_wechat') }}</p>
    </div>
    <iframe
      v-else
      ref="iframeRef"
      class="-translate-x-1.5 scale-[1] overflow-hidden"
      :style="{ height, width }"
      scrolling="no"
      :src="WECHAT_LOGIN_URL"
      frameborder="0"
      @load="handleLoad"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue';
import { useBasicLayout } from '@/hooks/useBasicLayout';
import { api_host } from '@/utils/config';

const props = withDefaults(
  defineProps<{
    width?: string
    height?: string
  }>(),
  {
    width: '100%',
    height: '280px',
  }
)

const emits = defineEmits(['oauthSuccess'])

const { isInMobile } = useBasicLayout()

// const WECHAT_LOGIN_URL = `https://work.wescrm.com/wechat_oauth_login.html?plain=1&height=280&appid=wxbe904d4182458106&suiteid=53aihub&api=${encodeURIComponent(`${api_host}/api/saas/wechat/redirect`)}&redirect_url=${encodeURIComponent(`${location.origin}/oauth_login.html`)}`
const WECHAT_LOGIN_URL = `https://work.wescrm.com/wechat_oauth_login.html?plain=1&height=280&appid=wxbe904d4182458106&suiteid=53aihub&api=${encodeURIComponent(`${api_host}/api/saas/wechat/redirect`)}&redirect_url=${encodeURIComponent(`${location.origin}/console/oauth_login.html`)}`

const iframeRef = ref()
const loading = ref(false)

const mobileWechatUrl = ref('')

let timer: any = null

const messageHandler = (res: any = {}) => {
  const data = res.data || {}
  const params = data.params || {}
  const from = data.from || params.from || ''
  const action = data.action || params.action || ''
  const state = data.state || params.state || ''
  const openid = data.openid || params.openid || ''
  const nickname = data.nickname || params.nickname || ''
  const unionid = data.unionid || params.unionid || ''
  const access_token = data.access_token || params.access_token || ''
  // if (!['53hub', '53aihub', 'wescrm'].includes(from))
  //   return
  clearInterval(timer)

  if (openid || access_token) {
    emits('oauthSuccess', { openid, nickname, unionid, access_token })
  }
}

onMounted(() => {
  if (isInMobile.value) {
  } else {
    loading.value = true
    timer = setInterval(() => {
      const contentWindow =
        (iframeRef.value.contentWindow && iframeRef.value.contentWindow[0]) || iframeRef.value.contentWindow
      if (contentWindow) {
        let oauthLoginData = null
        try {
          oauthLoginData = contentWindow.sessionStorage.getItem('oauth_login_data')
          oauthLoginData = JSON.parse(oauthLoginData)
        } catch (error) {
          oauthLoginData = null
        }
        if (oauthLoginData) messageHandler({ data: oauthLoginData })
      }
    }, 2000)
  }
  // window.addEventListener('message', messageHandler)
})
onBeforeUnmount(() => {
  // window.removeEventListener('message', messageHandler)
  clearInterval(timer)
})

const handleLogin = () => {
  // const redirect_url = QYY_HOST + `/v4/xbot/hubredirect?appid=${ process.env.VITE_GLOB_OFFICIALID }&state=wechat_redirect&redirecturl=${encodeURIComponent(location.origin)}`
  const redirect_url =
    'https://api.ibos.cn' +
    `/v4/xbot/hubredirect?appid=${import.meta.env.VITE_GLOB_OFFICIALID}&state=wechat_redirect&redirecturl=${encodeURIComponent(`${location.origin}/?login_way=wechat_login`)}`
  window.location.href = redirect_url
}

const handleLoad = () => {
  loading.value = false
}
</script>

<style scoped></style>
