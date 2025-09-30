<template>
  <Layout class="px-[60px] py-8">
    <!-- 页面标题 -->
    <Header :title="$t('module.SMTP')" />

    <div class="flex-1 flex flex-col gap-4 bg-white p-6 mt-3 box-border overflow-y-auto">
      <!-- 邮件日志配置区域 -->
      <div class="w-full px-6 py-4 border rounded hover:shadow">
        <div class="h-8 flex justify-between items-center">
          <p>{{ $t('module.SMTP_email_log') }}</p>
          <div class="flex items-center">
            <ElSwitch v-model="openEmail" @change="handleOpenEmail" />
            <span class="ml-2">{{ openEmail ? $t('action_enable') : $t('action_close') }}</span>
          </div>
        </div>

        <!-- 邮件配置表单 -->
        <EmailForm v-if="openEmail" ref="emialFormRef" />
      </div>

      <!-- 手机日志配置区域（暂未开放） -->
      <div class="w-full px-6 border rounded hover:shadow">
        <div class="h-16 flex justify-between items-center">
          <p>{{ $t('module.SMTP_mobile_log') }}</p>
          <div class="flex items-center">
            <ElSwitch v-model="openMobile" disabled @change="handleOpenMobile" @click="handleClick" />
            <span class="ml-2">{{ openMobile ? $t('action_enable') : $t('action_close') }}</span>
          </div>
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useEnterpriseStore } from '@/stores'
import EmailForm from './components/email-form.vue'

const SMTP_TYPE = {
  EMAIL: 'smtp',
  MOBILE: 'mobile',
} as const

type SMTPType = (typeof SMTP_TYPE)[keyof typeof SMTP_TYPE]

// 类型定义
interface SMTPConfig {
  content: string
  enabled: boolean
  type: SMTPType
}

// 状态管理
const store = useEnterpriseStore()

const emialFormRef = ref<InstanceType<typeof EmailForm>>() // 表单组件引用

// 响应式数据
const openEmail = ref<boolean>(false) // 邮件日志开关状态
const openMobile = ref<boolean>(false) // 手机日志开关状态（暂未开放）

/**
 * 处理邮件日志开关变化
 * 当关闭邮件日志时，保存当前配置并禁用
 */
const handleOpenEmail = async (): Promise<void> => {
  // 只有在关闭邮件日志时才需要保存配置
  if (!openEmail.value) {
    if (!emialFormRef.value) return

    const formData = emialFormRef.value.getData()
    if (!formData || Object.keys(formData).length === 0) return

    try {
      const requestConfig: SMTPConfig = {
        content: JSON.stringify(formData),
        enabled: false,
        type: SMTP_TYPE.EMAIL,
      }
      await store.saveSMTPInfo({ data: requestConfig })
    } catch (error) {
      console.error('保存SMTP配置失败:', error)
    }
  }
}

/**
 * 处理手机日志开关变化（暂未开放）
 */
const handleOpenMobile = (): void => {
  ElMessage.warning(window.$t('feature_coming_soon'))
}

/**
 * 处理手机日志开关点击（暂未开放）
 */
const handleClick = (): void => {
  ElMessage.warning(window.$t('feature_coming_soon'))
}

/**
 * 组件挂载时加载SMTP配置信息
 */
onMounted(async (): Promise<void> => {
  try {
    const data = await store.loadSMTPInfo()
    openEmail.value = data[0]?.enabled || false
  } catch (error) {
    console.error('加载SMTP配置失败:', error)
  }
})
</script>

<style scoped></style>
