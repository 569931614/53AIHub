<template>
  <Layout class="px-[60px] py-8">
    <!-- 页面标题 -->
    <Header :title="$t('module.statistics')" />

    <!-- 主要内容区域 -->
    <div class="flex-1 flex flex-col bg-white p-6 mt-3 box-border">
      <!-- 可滚动内容区域 -->
      <div class="flex-1 max-h-[calc(100vh-240px)] overflow-auto">
        <!-- 页面标题和描述 -->
        <h1 class="font-semibold text-[#1D1E1F]">
          {{ $t('module.statistics_header_title') }}
        </h1>
        <div class="text-[#9A9A9A] text-sm mt-4">
          {{ $t('module.statistics_header_desc') }}
        </div>
        <!-- 头部统计代码输入 -->
        <div class="text-[#9A9A9A] text-sm mt-6">
          {{ $t('module.statistics_textarea_label_1') }}
        </div>
        <ElInput
          v-model="head.value"
          v-loading="loading"
          class="mt-3 !w-[600px]"
          style="--el-input-bg-color: #f7f8fa"
          type="textarea"
          resize="none"
          :placeholder="$t('module.statistics_textarea_label_1_example')"
          :rows="8"
        />
        <!-- CSS样式代码输入 -->
        <div class="text-[#9A9A9A] text-sm mt-6">
          {{ $t('module.statistics_textarea_label_2') }}
        </div>
        <ElInput
          v-model="css.value"
          v-loading="loading"
          class="mt-3 !w-[600px]"
          style="--el-input-bg-color: #f7f8fa"
          type="textarea"
          resize="none"
          :rows="8"
          :placeholder="$t('module.statistics_textarea_label_2_example')"
        />
      </div>

      <!-- 底部操作区域 -->
      <ElDivider />
      <ElButton v-debounce class="h-[36px] w-[96px]" type="primary" size="large" @click="handleSave">
        {{ $t('action_save') }}
      </ElButton>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useSettingStore } from '@/stores/modules/setting'

// 设置存储实例
const settingStore = useSettingStore()

// 加载状态
const loading = ref(false)

// 统计配置项常量
const STATISTICS_KEYS = {
  HEAD: 'third_party_statistic_header',
  CSS: 'third_party_statistic_css',
} as const

// 头部统计代码配置
const head = reactive({
  setting_id: 0,
  key: STATISTICS_KEYS.HEAD,
  value: '',
})

// CSS样式代码配置
const css = reactive({
  setting_id: 0,
  key: STATISTICS_KEYS.CSS,
  value: '',
})

/**
 * 初始化页面数据
 * 从服务器加载现有的统计配置
 */
const initializeData = async () => {
  loading.value = true

  try {
    const settingsData = await settingStore.loadListData()

    // 查找并更新头部统计配置
    const headSetting = settingsData.find(item => item.key === STATISTICS_KEYS.HEAD)
    if (headSetting) {
      Object.assign(head, headSetting)
    }
    // 查找并更新CSS样式配置
    const cssSetting = settingsData.find(item => item.key === STATISTICS_KEYS.CSS)
    if (cssSetting) {
      Object.assign(css, cssSetting)
    }
  } finally {
    loading.value = false
  }
}

/**
 * 保存统计配置
 * 同时保存头部统计代码和CSS样式代码
 */
const handleSave = async () => {
  // 并行保存两个配置项
  const [headResult, cssResult] = await Promise.all([
    settingStore.save(head.setting_id, {
      key: head.key,
      value: head.value,
    }),
    settingStore.save(css.setting_id, {
      key: css.key,
      value: css.value,
    }),
  ])
  // 更新设置ID
  head.setting_id = headResult.setting_id || 0
  css.setting_id = cssResult.setting_id || 0

  // 显示成功消息
  ElMessage.success(window.$t('action_save_success'))
}

// 组件挂载时初始化数据
onMounted(initializeData)
</script>

<style scoped></style>
