<template>
  <div class="h-full bg-white">
    <MainHeader sticky :back="false" hide-user :sider-button="enterpriseStore.isSoftStyle">
      <template #before_suffix>
        <div class="text-base text-primary font-bold line-clamp-1 max-md:flex-1 max-md:text-center">
          {{ $t('action.setting') }}
        </div>
      </template>
      <template #after_suffix>
        <ElIcon class="text-primary cursor-pointer" size="18" @click="$router.back()">
          <Close />
        </ElIcon>
      </template>
    </MainHeader>

    <div class="hidden w-[250px] max-md:w-full flex-none py-5 bg-white">
      <div class="h-15 flex items-center gap-2 px-6 mb-2 text-primary max-md:hidden">
        <svg-icon size="16" name="setting"></svg-icon>
        <span class="text-[22px] max-md:text-lg">{{ $t('action.setting') }}</span>
      </div>
      <div class="flex flex-col py-2 gap-1 max-md:hidden">
        <template v-for="item in menus" :key="item.value">
          <div
            v-if="item.visible"
            class="h-10 flex items-center gap-3 px-6 cursor-pointer hover-bg-primary-light-9 hover-text-theme max-md:flex-1"
            :class="[page === item.value ? 'text-theme bg-primary-light-9' : 'text-regular']"
            @click="handleSelect(item)"
          >
            <svg-icon size="16" :name="item.icon || 'setting'"></svg-icon>
            <span class="text-sm">{{ item.label }}</span>
          </div>
        </template>
      </div>
      <el-tabs v-model="page" class="index-tabs md:!hidden mx-4">
        <template v-for="item in menus" :key="item.value">
          <el-tab-pane v-if="item.visible" :label="item.label" :name="item.value"></el-tab-pane>
        </template>
      </el-tabs>
    </div>

    <!-- 内容区域 -->
    <div class="flex-1 w-full md:w-4/5 lg:w-3/5 max-w-[600px] mx-auto box-border flex flex-col">
      <div class="flex-1 py-4 px-3 md:p-6 bg-[#FFFFFF] box-border overflow-y-auto">
        <UserInfo v-if="page === 'userinfo'"></UserInfo>
        <Password v-else-if="page === 'password'"></Password>
        <Toolbar v-else-if="page === 'toolbar'"></Toolbar>
        <Glider v-else-if="page === 'glider'"></Glider>
        <About v-else-if="page === 'about'"></About>
        <Common v-else-if="page === 'common'"></Common>

        <div
          v-if="page === 'userinfo'"
          class="h-11 mt-8 flex items-center justify-center bg-[#F8F8F9] gap-2 px-6 mb-2 rounded text-[#F84E55] cursor-pointer"
          @click="handleLogout"
        >
          <!-- <svg-icon size="14" name="quit"></svg-icon> -->
          <span class="text-sm">{{ $t('action.logout') }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Close } from '@element-plus/icons-vue'
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'

import MainHeader from '@/layout/header.vue'

import { useUserStore } from '@/stores/modules/user'
import { useEnterpriseStore } from '@/stores/modules/enterprise'

// import Toolbar from './toolbar.vue'
// import Glider from './glider.vue'
import About from './about.vue'
import Common from './common.vue'
import Password from './password.vue'
import UserInfo from './userinfo.vue'

const userStore = useUserStore()
const enterpriseStore = useEnterpriseStore()

const props = withDefaults(
  defineProps<{
    type?: string
  }>(),
  {
    type: 'glider'
  }
)

const route = useRoute()

const menus = computed(() => [
  {
    label: '划词菜单栏',
    value: 'glider',
    icon: 'mouse',
    visible: userStore.is_login && window.$isElectron
  },
  {
    label: '网页工具栏',
    value: 'toolbar',
    icon: 'web',
    visible: userStore.is_login && window.$isElectron
  },
  {
    label: '通用设置',
    value: 'common',
    icon: 'setting2',
    visible: window.$isElectron
  },
  {
    label: window.$t('profile.info'),
    value: 'userinfo',
    icon: 'person',
    visible: userStore.is_login
  },
  {
    label: window.$t('profile.change_password'),
    value: 'password',
    icon: 'lock',
    visible: userStore.is_login
  },
  {
    label: '关于',
    value: 'about',
    icon: 'info',
    visible: window.$isElectron
  }
])

const page = ref('')
const handleSelect = (item: any) => {
  page.value = item.value
}

const handleLogout = () => {
  userStore.logout()
  ElMessage.success(window.$t('status.logout_success'))
}

onMounted(() => {
  const visibleMenus = menus.value.filter((item) => item.visible).map((item) => item.value)
  const type = props.type || (route.query.type as string) || 'userinfo'
  page.value = visibleMenus.includes(type) ? type : visibleMenus[0]

  enterpriseStore.loadInfo()
})
</script>

<style scoped>
@media (width <= 768px) {
  .h-full {
    min-height: 100vh;
  }
}
</style>
