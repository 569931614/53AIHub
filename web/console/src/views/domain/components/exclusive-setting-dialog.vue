<template>
  <ElDialog
    v-model="visible"
    :title="$t('module.domain_exclusive')"
    :close-on-click-modal="false"
    width="700px"
    destroy-on-close
    append-to-body
    align-center
    @close="handleClose"
  >
    <ElForm ref="formRef" :model="form" label-position="top" @submit.prevent>
      <ElFormItem
        prop="domain"
        :rules="[
          {
            validator: validateDomain,
            trigger: 'blur',
          },
        ]"
      >
        <ElInput
          v-model="form.domain"
          size="large"
          :maxlength="20"
          show-word-limit
          :placeholder="$t('module.domain_exclusive')"
        >
          <template #prepend> https:// </template>
          <template #append> {{ isDevEnv ? '.hub' : '' }}.53ai.com </template>
        </ElInput>
      </ElFormItem>
    </ElForm>

    <template #footer>
      <div class="py-4 flex items-center justify-center">
        <ElButton class="w-24 h-9" type="primary" :loading="submitting" @click="handleConfirm">
          {{ $t('action_save') }}
        </ElButton>
        <ElButton class="w-24 h-9 text-[#1D1E1F]" type="info" plain @click="handleClose">
          {{ $t('action_cancel') }}
        </ElButton>
      </div>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { domainApi } from '@/api/modules/domain'
import { useEnv } from '@/hooks/useEnv'

// 类型定义
type DomainData = {
  id?: number
  domain?: string
  [key: string]: unknown
}

type FormData = {
  domain: string
}

type ValidationRule = {
  field: string
  fullField: string
  type: string
}

type ValidationCallback = (error?: Error) => void

// Hooks
const { isDevEnv } = useEnv()

// Emits
const emits = defineEmits<{
  (e: 'success'): void
}>()

// 组件引用
const formRef = ref<InstanceType<typeof ElForm>>()

// 响应式数据
const visible = ref(false)
const submitting = ref(false)
const originData = ref<DomainData>({})

const form = reactive<FormData>({
  domain: '',
})

// 计算属性
const domainSuffix = computed(() => {
  return `${isDevEnv.value ? '.hub' : ''}.53ai.com`
})

// 方法定义
const validateDomain = (rule: ValidationRule, value: string, callback: ValidationCallback): void => {
  const trimmedValue = (value || '').trim()

  if (!trimmedValue) {
    callback(new Error(window.$t('form_input_placeholder')))
    return
  }

  // 域名格式验证
  if (!/^[a-z0-9-]{5,20}$/.test(trimmedValue)) {
    callback(new Error(window.$t('module.domain_exclusive_validator_1')))
    return
  }

  // 不能以连字符开头或结尾
  if (trimmedValue.startsWith('-') || trimmedValue.endsWith('-')) {
    callback(new Error(window.$t('module.domain_exclusive_validator_2')))
    return
  }

  callback()
}

const resetForm = () => {
  form.domain = ''
}

const populateForm = (data: DomainData) => {
  const domain = data.domain || ''
  // 移除后缀，只保留前缀部分
  form.domain = domain
    .replace(/^https?:\/\//, '')
    .replace(/\.hub\.53ai\.com$/, '')
    .replace(/\.53ai\.com$/, '')
}

const buildDomainUrl = (): string => {
  return `${form.domain}${domainSuffix.value}`
}

const open = ({ data = {} }: { data?: DomainData } = {}) => {
  populateForm(data)
  originData.value = data
  visible.value = true
}

const handleClose = () => {
  visible.value = false
  resetForm()
}

const handleConfirm = async () => {
  const isValid = await formRef.value?.validate()
  if (!isValid) return
  try {
    submitting.value = true

    const domainUrl = buildDomainUrl()
    const requestData = { domain: domainUrl }

    if (originData.value.id) {
      await domainApi.updateExclusive(originData.value.id, requestData)
    } else {
      await domainApi.createExclusive(requestData)
    }

    ElMessage.success(window.$t('action_save_success'))
    emits('success')
    handleClose()
  } catch (error) {
    console.error('保存独占域名失败:', error)
    ElMessage.error(window.$t('action_save_failed'))
  } finally {
    submitting.value = false
  }
}

// 暴露方法
defineExpose({
  open,
  close: handleClose,
  reset: resetForm,
})
</script>

<style scoped lang="scss">
// 可以在这里添加组件特定的样式
</style>
