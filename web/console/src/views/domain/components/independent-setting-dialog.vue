<template>
  <ElDialog
    v-model="visible"
    :title="$t('module.domain_independent')"
    :close-on-click-modal="false"
    width="700px"
    destroy-on-close
    append-to-body
    align-center
    @close="handleClose"
  >
    <ElForm ref="formRef" :model="form" label-position="top" @submit.prevent>
      <!-- 域名输入 -->
      <ElFormItem
        prop="domain"
        :rules="[
          {
            validator: validateDomain,
            trigger: 'blur',
          },
        ]"
      >
        <div class="flex items-center w-full">
          <ElInput
            v-model="form.domain"
            class="flex-1"
            :class="[shouldShowSubdirInput && 'has-subdir']"
            size="large"
            :maxlength="20"
            show-word-limit
            :placeholder="$t('module.domain_independent')"
          >
            <template #prepend> https:// </template>
          </ElInput>

          <!-- 子目录输入 -->
          <ElFormItem
            v-if="shouldShowSubdirInput"
            prop="subdir"
            :rules="generateInputRules({ message: 'form_input_placeholder' })"
          >
            <ElInput
              v-model="form.subdir"
              class="flex-none w-[250px] h-[42px] subdir-input"
              size="large"
              :maxlength="10"
              show-word-limit
              :placeholder="$t('form_input_placeholder')"
            >
              <template #prepend> / </template>
            </ElInput>
          </ElFormItem>
        </div>
      </ElFormItem>

      <!-- 子目录开关 -->
      <ElFormItem v-if="isCustomResolveType">
        <div class="flex items-center text-sm text-[#4F5052]">
          <span>{{ $t('module.use_subdirectories') }}</span>
          <ElTooltip :content="$t('module.use_subdirectories_tip')">
            <svg-icon class="text-[#A4AABA] ml-1" name="help" width="14" />
          </ElTooltip>
          <ElSwitch v-model="form.use_subdir" class="ml-2" size="small" />
        </div>
      </ElFormItem>

      <!-- 解析类型选择 -->
      <ElFormItem>
        <ElRadioGroup v-model="form.resolve_type" class="w-full">
          <ElRadio
            v-for="resolveType in resolveTypeOptions"
            :key="resolveType.value"
            class="flex-1 border py-6 px-4 rounded overflow-hidden"
            :class="[form.resolve_type === resolveType.value ? 'border-[#3664EF]' : '']"
            :value="resolveType.value"
          >
            {{ resolveType.label }}
          </ElRadio>
        </ElRadioGroup>
      </ElFormItem>

      <!-- CNAME 解析说明和配置 -->
      <template v-if="isCnameResolveType">
        <ul class="w-full flex flex-col gap-3 bg-[#F6F9FC] p-5 mb-6 box-border text-sm text-[#4F5052]">
          <li>{{ $t('module.domain_independent_cname_desc') }}</li>
          <li>{{ $t('module.domain_independent_cname_desc_1') }}</li>
          <li>{{ $t('module.domain_independent_cname_desc_2') }}</li>
          <li>{{ $t('module.domain_independent_cname_desc_3') }}</li>
        </ul>

        <!-- HTTPS 配置 -->
        <ElFormItem>
          <div class="flex items-center gap-2 text-sm text-[#4F5052]">
            <span>{{ $t('module.domain_independent_https') }}</span>
            <ElSwitch v-model="form.enable_https" size="small" />
            <template v-if="form.enable_https">
              <span class="ml-12">{{ $t('module.domain_independent_https_always') }}</span>
              <ElSwitch v-model="form.force_https" size="small" />
            </template>
          </div>
        </ElFormItem>
      </template>

      <!-- 自定义解析说明 -->
      <template v-else>
        <ul class="w-full flex flex-col gap-3 bg-[#F6F9FC] p-5 mb-6 box-border text-sm text-[#4F5052]">
          <li>{{ $t('module.domain_independent_self_desc_1') }}</li>
          <li>{{ $t('module.domain_independent_self_desc_2') }}</li>
          <ElDivider class="!my-2" />
          <li>{{ $t('module.domain_independent_self_desc_3', { site_id: enterpriseStore.info.eid }) }}</li>
        </ul>
      </template>
    </ElForm>

    <template #footer>
      <div class="py-4 flex items-center justify-center">
        <ElButton class="w-[96px] h-[36px]" type="primary" :loading="submitting" @click="handleConfirm">
          {{ $t('action_save') }}
        </ElButton>
        <ElButton class="w-[96px] h-[36px] text-[#1D1E1F]" type="info" plain @click="handleClose">
          {{ $t('action_cancel') }}
        </ElButton>
      </div>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { domainApi } from '@/api/modules/domain'
import { INDEPENDENT_RESOLVE_TYPE, INDEPENDENT_SSL_CERT_TYPE } from '@/constants/domain'
import { generateInputRules } from '@/utils/form-rule'
import { useEnterpriseStore } from '@/stores/modules/enterprise'

// 类型定义
type DomainConfig = {
  resolve_type?: number
  enable_https?: boolean | number
  force_https?: boolean | number
  ssl_cert_type?: number
  ssl_certificate?: string
  ssl_private_key?: string
  subdir?: string
  use_subdir?: boolean | number
  [key: string]: unknown
}

type DomainData = {
  id?: number
  domain?: string
  config?: DomainConfig
  [key: string]: unknown
}

type FormData = {
  domain: string
  resolve_type: number
  enable_https: boolean
  force_https: boolean
  ssl_cert_type: number
  ssl_certificate: string
  ssl_private_key: string
  subdir: string
  use_subdir: boolean
}

type ValidationRule = {
  field: string
  fullField: string
  type: string
}

type ValidationCallback = (error?: Error) => void

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
  resolve_type: INDEPENDENT_RESOLVE_TYPE.CNAME,
  enable_https: false,
  force_https: false,
  ssl_cert_type: INDEPENDENT_SSL_CERT_TYPE['53AI'],
  ssl_certificate: '',
  ssl_private_key: '',
  subdir: '',
  use_subdir: false,
})

// Store
const enterpriseStore = useEnterpriseStore()

// 计算属性
const isCnameResolveType = computed(() => form.resolve_type === INDEPENDENT_RESOLVE_TYPE.CNAME)
const isCustomResolveType = computed(() => form.resolve_type === INDEPENDENT_RESOLVE_TYPE.CUSTOM)
const shouldShowSubdirInput = computed(() => isCustomResolveType.value && form.use_subdir)

const resolveTypeOptions = computed(() => [
  {
    value: INDEPENDENT_RESOLVE_TYPE.CNAME,
    label: window.$t('module.domain_independent_cname'),
  },
  {
    value: INDEPENDENT_RESOLVE_TYPE.CUSTOM,
    label: window.$t('module.domain_independent_self'),
  },
])

// 方法定义
const validateDomain = (rule: ValidationRule, value: string, callback: ValidationCallback): void => {
  const trimmedValue = (value || '').trim()
  if (trimmedValue) {
    callback()
  } else {
    callback(new Error(window.$t('form_input_placeholder')))
  }
  // 这里可以添加更严格的域名验证逻辑
  // if (!/^[a-z0-9-]{5,}$/.test(trimmedValue)) {
  //   return callback(new Error('域名格式不正确'))
  // }
}

const resetForm = () => {
  Object.assign(form, {
    domain: '',
    resolve_type: INDEPENDENT_RESOLVE_TYPE.CNAME,
    enable_https: false,
    force_https: false,
    ssl_cert_type: INDEPENDENT_SSL_CERT_TYPE['53AI'],
    ssl_certificate: '',
    ssl_private_key: '',
    subdir: 'chat',
    use_subdir: false,
  })
}

const populateForm = (data: DomainData) => {
  const config = data.config || {}

  form.domain = (data.domain || '').trim().replace(/^https?:\/\//, '')
  form.resolve_type = Number(config.resolve_type) || INDEPENDENT_RESOLVE_TYPE.CNAME
  form.enable_https = Boolean(Number(config.enable_https))
  form.force_https = Boolean(Number(config.force_https))
  form.ssl_cert_type = Number(config.ssl_cert_type) || INDEPENDENT_SSL_CERT_TYPE['53AI']
  form.ssl_certificate = config.ssl_certificate || ''
  form.ssl_private_key = config.ssl_private_key || ''
  form.subdir = config.subdir || 'chat'
  form.use_subdir = Boolean(Number(config.use_subdir))
}

const buildConfigData = (): DomainConfig => ({
  resolve_type: form.resolve_type,
  enable_https: form.enable_https,
  force_https: form.force_https,
  ssl_cert_type: form.ssl_cert_type,
  ssl_certificate: form.ssl_certificate,
  ssl_private_key: form.ssl_private_key,
  subdir: form.subdir,
  use_subdir: form.use_subdir,
})

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

    const requestData = {
      domain: form.domain,
      config: buildConfigData(),
    }

    if (originData.value.id) {
      await domainApi.updateIndependent(originData.value.id, requestData)
    } else {
      await domainApi.createIndependent(requestData)
    }

    ElMessage.success(window.$t('action_save_success'))
    emits('success')
    handleClose()
  } catch (error) {
    console.error('保存独立域名失败:', error)
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
:deep(.has-subdir .el-input__wrapper) {
  border-top-right-radius: 0 !important;
  border-bottom-right-radius: 0 !important;
  box-shadow: none !important;
  border-style: solid;
  border-color: #dcdfe6;
  border-top-width: 1px;
  border-bottom-width: 1px;
  border-left-width: 1px;
  box-sizing: border-box;
}

:deep(.subdir-input .el-input-group__prepend) {
  border-radius: 0 !important;
  padding: 0 8px !important;
}
</style>
