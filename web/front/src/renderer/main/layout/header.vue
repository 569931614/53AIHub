<template>
  <!--  v-if="needLogin ? userStore.is_login : true" -->
  <header class="flex-none h-[70px] border-b sticky top-0 z-10" :class="[type === 'homepage' ? 'nav-bg' : 'bg-white']">
    <div class="mx-auto px-4 flex items-center justify-between h-full" :class="mainClass">
      <div class="flex-1 flex items-center gap-2 overflow-hidden">
        <slot name="before_prefix"></slot>
        <div
          v-if="siderButton && !globalStore.siderVisible"
          v-tooltip="{
            content: $t('chat.expand_side_bar')
          }"
          class="flex-none size-7 rounded-md flex-center cursor-pointer hover:bg-[#ECEDEE]"
          @click="globalStore.toggleSider"
        >
          <svg-icon name="layout-left" size="20" color="#9A9A9A"></svg-icon>
        </div>
        <div
          v-if="back"
          v-tooltip="{
            content: $t('common.back')
          }"
          class="flex-none size-7 rounded-md flex-center cursor-pointer max-md:hidden hover:bg-[#ECEDEE]"
          @click="$router.back()"
        >
          <ElIcon class="text-regular cursor-pointer" size="14">
            <ArrowLeft />
          </ElIcon>
        </div>
        <slot name="before_suffix"></slot>
      </div>

      <div class="flex items-center gap-2">
        <slot name="after_prefix"></slot>
        <template v-if="!$isElectron && !hideUser">
          <el-tooltip v-if="userStore.is_login" effect="light" popper-class="el-popper--plain" placement="bottom-end">
            <template #content>
              <div class="w-[300px] p-4 flex items-center gap-2">
                <el-image class="flex-none size-10 rounded-full" :src="userStore.info.avatar"></el-image>
                <div class="flex-1 overflow-hidden">
                  <div class="w-full flex items-center gap-1 overflow-hidden">
                    <p class="flex-1 text-sm text-primary font-medium truncate">
                      {{ userStore.info.nickname }}
                    </p>
                    <div
                      v-if="!userStore.info.is_internal"
                      class="h-6 flex-center gap-1 bg-[#F7F7F7] rounded-full px-2 text-xs text-placeholder whitespace-nowrap"
                      :title="userStore.info.group_name"
                    >
                      <img
                        :src="
                          !/\.png$/.test(userStore.info.group_icon)
                            ? $getPublicPath(`/images/subscription/${userStore.info.group_icon}.png`)
                            : userStore.info.group_icon
                        "
                        class="w-4 h-4 object-cover"
                      />
                      <p class="max-w-[5em] truncate">{{ userStore.info.group_name }}</p>
                    </div>
                    <div
                      v-if="upgrade_visible && !userStore.info.is_internal"
                      class="flex items-center gap-1 ml-auto cursor-pointer hover:opacity-70 bg-[#F4F0FF] rounded-2xl h-6 px-2 box-border text-xs text-[#8E5EFF] whitespace-nowrap"
                      @click.stop="upgrade_ref.open"
                    >
                      <img :src="$getPublicPath(`/images/subscription/upgrade.png`)" class="w-4 h-4 object-cover" />
                      {{ $t('subscription.upgrade') }}
                    </div>
                  </div>
                  <div class="text-xs text-[#9A9A9A]">{{ userStore.info.email }}</div>
                </div>
              </div>
              <div class="flex flex-col gap-1.5 px-3 py-1.5 border-t">
                <router-link
                  :to="{ name: 'Profile' }"
                  class="h-8 px-3 flex items-center gap-2 rounded cursor-pointer hover:bg-[#ECEDEE] text-primary"
                >
                  <div class="flex-center size-6">
                    <svg-icon name="setting" size="16"></svg-icon>
                  </div>
                  <span class="text-sm">{{ $t('action.setting') }}</span>
                </router-link>
                <div
                  v-if="userStore.info.role > 1"
                  class="h-8 px-3 flex items-center gap-2 rounded cursor-pointer hover:bg-[#ECEDEE] text-primary"
                  @click="handleJumpToAdmin"
                >
                  <div class="flex-center size-6">
                    <svg-icon name="jump" size="16"></svg-icon>
                  </div>
                  <span class="text-sm">{{ $t('common.go_admin') }}</span>
                </div>
                <a
                  href="https://doc.53ai.com/%E5%85%A5%E9%97%A8/%E6%AC%A2%E8%BF%8E%E4%BD%BF%E7%94%A8.html"
                  target="_blank"
                  class="h-8 px-3 flex items-center gap-2 rounded cursor-pointer text-primary hover:bg-[#ECEDEE]"
                >
                  <div class="flex-center size-6">
                    <svg-icon name="toolkit" size="16" stroke></svg-icon>
                  </div>
                  <span class="text-sm">{{ $t('common.new_friend') }}</span>
                </a>
                <a
                  href="https://doc.53ai.com/%E7%A4%BE%E5%8C%BA/%E9%9C%80%E6%B1%82%E6%94%AF%E6%8C%81.html"
                  target="_blank"
                  class="h-8 px-3 flex items-center gap-2 rounded cursor-pointer text-primary hover:bg-[#ECEDEE]"
                >
                  <div class="flex-center size-6">
                    <svg-icon name="help" size="16"></svg-icon>
                  </div>
                  <span class="text-sm">{{ $t('common.help_feedback') }}</span>
                </a>
              </div>
              <div class="flex flex-col gap-1.5 px-3 py-1.5 border-t">
                <div v-debounce class="h-8 px-3 flex items-center gap-2 rounded cursor-pointer hover:bg-[#ECEDEE] text-primary" @click="handleLogout">
                  <div class="flex-center size-6">
                    <svg-icon name="quit" size="14"></svg-icon>
                  </div>
                  <span class="text-sm">{{ $t('login.quit') }}</span>
                </div>
              </div>
            </template>
            <div class="flex items-center gap-1.5 cursor-pointer max-md:hidden">
              <el-image class="size-[26px] rounded-full" :src="userStore.info.avatar"></el-image>
              <el-icon color="#A4A4AD">
                <ArrowDown />
              </el-icon>
            </div>
          </el-tooltip>
          <el-button v-else type="primary" @click="handleLogin">{{ $t('action.login') }}</el-button>
        </template>

        <slot name="after_suffix"></slot>
      </div>
    </div>

    <Upgrade ref="upgrade_ref" />
  </header>
</template>

<script setup lang="ts">
import { onMounted, ref, onUnmounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import HubUiX from 'hub-ui-x'
import { ArrowDown, ArrowLeft } from '@element-plus/icons-vue'

import Upgrade from '@/components/Upgrade/index.vue'

import { useUserStore } from '@/stores/modules/user'
import { useEnterpriseStore } from '@/stores/modules/enterprise'
import { useGlobalStore } from '@/stores/modules/global'

import eventBus from '@/utils/event-bus'

import { EVENT_NAMES } from '@/constants/events'

import { ADMIN_URL } from '@/api/host'

import { checkPermission } from '@/utils/permission'

withDefaults(
  defineProps<{
    needLogin?: boolean
    mainClass?: any
    siderButton?: boolean
    type?: 'homepage' | 'chat'
    sticky?: boolean
    hideUser?: boolean
    back?: boolean
  }>(),
  {
    needLogin: true,
    mainClass: '',
    siderButton: true,
    back: true
  }
)

const { locale } = useI18n()
const userStore = useUserStore()
const globalStore = useGlobalStore()
const enterpriseStore = useEnterpriseStore()

const upgrade_ref = ref(null)
const upgrade_visible = ref(false)

const handleLogin = async () => {
  await checkPermission()
}

const handleJumpToAdmin = () => {
  const url = `${ADMIN_URL}?access_token=${userStore.info.access_token}&eid=${userStore.info.eid}&from_origin=${encodeURIComponent(window.location.origin)}`
  console.info('adminUrl: ', url)
  window.open(url, '_blank')
}

const handleLogout = () => {
  userStore.logout()
  ElMessage.success(window.$t('status.logout_success'))
}

const validateUpgrade = async () => {
  if (userStore.info.access_token) {
    await nextTick()
    if (upgrade_ref.value) upgrade_visible.value = await upgrade_ref.value.validateUpgrade()
  }
}

const updateLocal = (language) => {
  locale.value = language
  HubUiX.setLang(language)
}

// 初始化
enterpriseStore.loadFromStorage().then((res) => {
  updateLocal(res.language)
})
onMounted(async () => {
  // 获取最新的企业信息
  await enterpriseStore.loadInfo().then((res) => {
    updateLocal(res.language)
  })
  await validateUpgrade()
  eventBus.on(EVENT_NAMES.LOGIN_SUCCESS, validateUpgrade)

  if (window.$chat53ai) {
    window.$chat53ai.$on('agenthub:service', (_event, { type, data } = {}) => {
      console.log(type, data)
      if (type === 'login') {
        handleLogin()
      }
    })
  }
})

onUnmounted(() => {
  eventBus.off(EVENT_NAMES.LOGIN_SUCCESS, validateUpgrade)
})
</script>

<style></style>
