<template>
  <div class="mt-5 w-3/5">
    <!-- SMTP 配置表单 -->
    <ElForm ref="formRef" :model="formData" size="large" label-position="top">
      <!-- SMTP 服务器地址 -->
      <ElFormItem
        :label="$t('module.SMTP_server')"
        prop="smtp_host"
        :rules="generateInputRules({ message: 'form_input_placeholder' })"
      >
        <ElInput v-model="formData.smtp_host" :placeholder="$t('form_input_placeholder')" clearable />
      </ElFormItem>

      <!-- SMTP 端口 -->
      <ElFormItem
        :label="$t('module.SMTP_port')"
        prop="smtp_port"
        :rules="generateInputRules({ message: 'form_input_placeholder', validator: ['port'] })"
      >
        <ElInput v-model="formData.smtp_port" :placeholder="$t('form_input_placeholder')" clearable />
      </ElFormItem>

      <!-- 邮箱账号 -->
      <ElFormItem
        :label="$t('module.SMTP_email_account')"
        prop="smtp_username"
        :rules="generateInputRules({ message: 'form_input_placeholder', validator: ['email'] })"
      >
        <ElInput v-model="formData.smtp_username" :placeholder="$t('form_input_placeholder')" clearable />
      </ElFormItem>

      <!-- 邮箱密码 -->
      <ElFormItem
        :label="$t('module.SMTP_email_password')"
        prop="smtp_password"
        :rules="generateInputRules({ message: 'form_input_placeholder', validator: ['password'] })"
      >
        <ElInput
          v-model="formData.smtp_password"
          type="password"
          :placeholder="$t('form_input_placeholder')"
          clearable
          show-password
        />
      </ElFormItem>

      <!-- 发件人邮箱 -->
      <ElFormItem
        :label="$t('module.SMTP_addresser_email')"
        prop="smtp_from"
        :rules="generateInputRules({ message: 'form_input_placeholder', validator: ['email'] })"
      >
        <ElInput v-model="formData.smtp_from" :placeholder="$t('form_input_placeholder')" clearable />
      </ElFormItem>

      <!-- 启用 TLS/SSL -->
      <ElFormItem :label="$t('module.SMTP_openTLS')">
        <ElSwitch v-model="formData.smtp_is_ssl" />
      </ElFormItem>

      <!-- 收件人邮箱和测试发送 -->
      <ElFormItem
        :label="$t('module.SMTP_receiver_email')"
        prop="smtp_to"
        :rules="generateInputRules({ message: 'form_input_placeholder', validator: ['email'] })"
      >
        <div class="w-full flex gap-3">
          <ElInput v-model="formData.smtp_to" :placeholder="$t('form_input_placeholder')" clearable />
          <ElButton v-debounce type="primary" plain :disabled="countDown > 0" @click="handleSendEmail">
            {{ countDown > 0 ? `${countDown}s` : $t('module.SMTP_send_email') }}
          </ElButton>
        </div>
      </ElFormItem>
    </ElForm>

    <!-- 操作按钮 -->
    <div class="flex gap-3">
      <ElButton :loading="isSaving" type="primary" class="w-24 h-9" @click="handleSave">
        {{ $t('action.save') }}
      </ElButton>
      <ElButton class="w-24 h-9" @click="handleReset">
        {{ $t('action_reset') }}
      </ElButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { ElForm } from 'element-plus'
import { generateInputRules } from '@/utils/form-rule'
import { useEnterpriseStore } from '@/stores'

// 类型定义
interface SMTPFormData {
  smtp_host: string
  smtp_port: string
  smtp_username: string
  smtp_password: string
  smtp_from: string
  smtp_to: string
  smtp_is_ssl: boolean
}

interface SMTPConfig {
  content: string
  enabled: boolean
  type: 'smtp' | 'mobile'
}

interface TestEmailConfig {
  from: string
  host: string
  is_ssl: boolean
  password: string
  port: number
  to: string
  username: string
}

// 常量定义
const COUNTDOWN_DURATION = 60 // 发送邮件后的倒计时时长（秒）
const SMTP_TYPE = 'smtp' as const // SMTP 配置类型

// 状态管理
const store = useEnterpriseStore()

const formRef = ref<InstanceType<typeof ElForm>>() // 表单引用

// 响应式数据
const isSaving = ref<boolean>(false) // 保存状态
const countDown = ref<number>(0) // 倒计时

// 定时器
let timer: NodeJS.Timeout | null = null

// 表单数据
const formData = reactive<SMTPFormData>({
  smtp_host: '',
  smtp_port: '',
  smtp_username: '',
  smtp_password: '',
  smtp_from: '',
  smtp_to: '',
  smtp_is_ssl: true,
})

/**
 * 启动倒计时
 */
const startCountdown = (): void => {
  countDown.value = COUNTDOWN_DURATION
  timer = setInterval(() => {
    countDown.value--
    if (countDown.value <= 0) {
      clearInterval(timer as NodeJS.Timeout)
      timer = null
      countDown.value = 0
    }
  }, 1000)
}
/**
 * 发送测试邮件
 */
const handleSendEmail = async (): Promise<void> => {
  const valid = await formRef.value?.validate()
  if (!valid) return
  try {
    const config: TestEmailConfig = {
      from: formData.smtp_from,
      host: formData.smtp_host,
      is_ssl: formData.smtp_is_ssl,
      password: formData.smtp_password,
      port: Number(formData.smtp_port),
      to: formData.smtp_to,
      username: formData.smtp_username,
    }

    await store.sendTestEmail({ data: config })
    ElMessage.success(window.$t('action_send_success'))

    // 启动倒计时
    startCountdown()
  } catch (error) {
    console.error('发送测试邮件失败:', error)
  }
}

/**
 * 重置表单数据
 */
const handleReset = (): void => {
  Object.assign(formData, {
    smtp_host: '',
    smtp_port: '',
    smtp_username: '',
    smtp_password: '',
    smtp_from: '',
    smtp_to: '',
    smtp_is_ssl: true,
  })
}

/**
 * 保存 SMTP 配置
 */
const handleSave = async (): Promise<void> => {
  const valid = await formRef.value?.validate()
  if (!valid) return

  isSaving.value = true

  const requestConfig: SMTPConfig = {
    content: JSON.stringify(formData),
    enabled: true,
    type: SMTP_TYPE,
  }

  try {
    await store.saveSMTPInfo({ data: requestConfig })
    ElMessage.success(window.$t('action_save_success'))
  } catch (error) {
    console.error('保存SMTP配置失败:', error)
  } finally {
    isSaving.value = false
  }
}

/**
 * 获取表单数据并加载配置
 */
const getData = (): SMTPFormData => {
  store.loadSMTPDetail({ data: { type: SMTP_TYPE } }).then(data => {
    const res = JSON.parse(data.content)
    Object.assign(formData, {
      smtp_host: res.smtp_host || '',
      smtp_port: res.smtp_port || '',
      smtp_username: res.smtp_username || '',
      smtp_password: res.smtp_password || '',
      smtp_to: res.smtp_to || '',
      smtp_from: res.smtp_from || '',
      smtp_is_ssl: res.smtp_is_ssl ?? true,
    })
  })
  return { ...formData }
}

// 暴露给父组件的方法
defineExpose({
  getData,
})

/**
 * 组件挂载时加载配置数据
 */
onMounted((): void => {
  getData()
})

/**
 * 组件卸载时清理定时器
 */
onUnmounted((): void => {
  if (timer) {
    clearInterval(timer)
    timer = null
  }
})
</script>

<style scoped></style>
