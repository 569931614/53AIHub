<template>
  <Layout class="fixed top-0 left-0 w-screen h-screen z-[9] bg-[#F4F6FA]">
    <!-- 头部工具栏 -->
    <header
      class="w-full px-[56px] h-[70px] flex items-center gap-3 shadow box-border bg-white sticky top-0 left-0 right-0 z-[9]"
    >
      <SvgIcon name="web-edit" style="zoom: 1.2" color="#858585" width="24" />
      <div class="flex-1 flex flex-col gap-0.5">
        <span>{{ navigationDetail.name || $t(route.meta?.title as string) }}</span>
        <span class="text-xs text-[#9A9BA0]">
          {{ $t('last_edit') }}: {{ formatLastEditTime(navigationDetail.content_update_time) }}
        </span>
      </div>
      <ElButton class="!ml-0 !min-w-[96px] !h-9" type="info" plain size="large" @click="handleCancel">
        {{ $t('action_cancel') }}
      </ElButton>
      <ElButton class="!ml-0 !min-w-[96px] !h-9" type="primary" size="large" :loading="isSaving" @click="handleSave">
        {{ $t('action_save') }}
      </ElButton>
    </header>

    <!-- 主要内容区域 -->
    <div
      v-loading="isLoading"
      class="flex-1 flex flex-col w-5/6 max-w-[1084px] rounded box-border my-5 mx-auto bg-white"
    >
      <!-- 导航栏预览 -->
      <div class="w-full h-[76px] px-8 box-border flex items-center gap-4 border-b">
        <ElImage :src="enterpriseInfo.logo" class="flex-none size-10 rounded" fit="cover" />
        <h2 class="flex-none text-[#1D1E1F] font-semibold">
          {{ enterpriseInfo.display_name || '--' }}
        </h2>
        <ElMenu class="flex-1 w-0 overflow-hidden ml-2 !border-none" mode="horizontal">
          <ElMenuItem
            v-for="item in navigationList"
            :key="item.navigation_id"
            class="!cursor-auto !opacity-100"
            :index="item.jump_path"
            disabled
          >
            <span class="!text-base !text-[#1D1E1F]">{{ item.name }}</span>
          </ElMenuItem>
        </ElMenu>
      </div>

      <!-- 编辑器区域 -->
      <div class="flex-1 w-full p-2 box-border">
        <UEditor ref="ueditorRef" />
      </div>

      <!-- 版权信息 -->
      <div
        class="w-full h-[64px] px-[56px] box-border flex items-center bg-[#22252E] rounded-sm text-sm text-[#989A9D]"
      >
        {{ enterpriseInfo.copyright }}
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import UEditor from '@/components/UEditor/index.vue'
import eventBus from '@/utils/event-bus'
import { navigationApi } from '@/api/modules/navigation/index'
import { useEnterpriseStore } from '@/stores'
import type { NavigationItem } from '@/api/modules/navigation/types'

// 类型定义
type NavigationDetail = {
  name: string
  updated_time: string
  content_update_time: string
  content?: {
    html_content?: string
    updated_time?: string | number
  }
}

type NavigationDetailResponse = NavigationItem & {
  updated_time: string
  content_update_time: string
  content?: {
    html_content?: string
    updated_time?: string | number
  }
}

// 路由和状态管理
const route = useRoute()
const router = useRouter()
const enterpriseStore = useEnterpriseStore()

// 组件引用
const ueditorRef = ref()

// 响应式数据
const navigationDetail = ref<NavigationDetail>({} as NavigationDetail)
const navigationList = ref<NavigationItem[]>([])
const isLoading = ref(false)
const isSaving = ref(false)

// 计算属性
const enterpriseInfo = computed(() => enterpriseStore.info)

// 工具方法
const formatLastEditTime = (timestamp?: string | number) => {
  if (!timestamp) return ''
  return new Date(timestamp).toLocaleString().replace(/\//g, '-').slice(0, 15)
}

// 数据加载方法
const loadNavigationDetail = async () => {
  const { navigation_id } = route.params
  isLoading.value = true

  try {
    const data = (await navigationApi.detail(Number(navigation_id))) as NavigationDetailResponse

    // 格式化时间
    data.updated_time = formatLastEditTime(data.updated_time)
    const contentData = data.content || {}
    data.content_update_time = formatLastEditTime(contentData.updated_time)

    // 设置编辑器内容
    ueditorRef.value?.setValue(contentData.html_content || '')
    navigationDetail.value = data as NavigationDetail
  } catch (error) {
    console.error('加载导航详情失败:', error)
  } finally {
    isLoading.value = false
  }
}

const loadNavigationList = async () => {
  try {
    const { list = [] } = await navigationApi.list({})
    navigationList.value = list
  } catch (error) {
    console.error('加载导航列表失败:', error)
  }
}

const handleCancel = () => {
  router.replace({ name: 'Navigation' })
}
// 事件处理方法
const handleSave = async () => {
  try {
    const html = await ueditorRef.value?.getHtml()
    if (!html) return

    isSaving.value = true

    await navigationApi.saveContent({
      navigation_id: +route.params.navigation_id,
      html_content: html,
    })

    ElMessage.success(window.$t('action_save_success'))
    handleCancel()
  } catch (error) {
    console.error('保存内容失败:', error)
  } finally {
    isSaving.value = false
  }
}

// 生命周期
onMounted(async () => {
  isLoading.value = true
  try {
    await Promise.all([loadNavigationList(), loadNavigationDetail()])
  } finally {
    isLoading.value = false
  }
  eventBus.on('user-login-success', loadNavigationDetail)
})

onUnmounted(() => {
  eventBus.off('user-login-success', loadNavigationDetail)
})
</script>

<style scoped lang="scss"></style>
