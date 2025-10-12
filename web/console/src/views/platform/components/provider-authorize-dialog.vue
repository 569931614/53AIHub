<template>
  <ElDialog
    v-model="visible"
    :title="$t('action_authorize') + $t(originData.label || '')"
    :close-on-click-modal="false"
    width="720px"
    destroy-on-close
    append-to-body
    @close="close"
  >
    <ElForm ref="formRef" :model="form" :rules="formRules" label-position="top">
      <!-- 提示信息 -->
      <div class="w-full flex flex-col gap-3 bg-[#F6F9FC] p-5 mb-4 box-border text-sm text-[#4F5052]">
        <div class="whitespace-pre-wrap leading-7" v-html="guideHtml" />
        <ElIcon
          v-if="isCozeCN"
          ref="copyRef"
          class="cursor-pointer ml-1 mt-1 text-[#4F5052] hover:text-[#3664EF]"
          :size="14"
          @click="handleCopy(coze_auth_url)"
        >
          <CopyDocument />
        </ElIcon>
      </div>
      <el-form-item prop="name" :label="$t('module.website_info_name')">
        <ElInput v-model="form.name" size="large" :placeholder="$t('module.website_info_name_placeholder')" />
      </el-form-item>
      <template v-for="option in schemaOptions" :key="option.prop">
        <ElFormItem :label="option.label" :prop="option.prop">
          <ElInput
            :model-value="getFormValue(option.prop)"
            size="large"
            :placeholder="option.placeholder"
            @update:model-value="setFormValue(option.prop, $event)"
          />
        </ElFormItem>
      </template>
    </ElForm>

    <template #footer>
      <div v-if="isCozeCN" class="text-center text-sm text-[#9A9A9A]">
        {{ $t('platform_auth.coze_cn.tip_1') }}
      </div>
      <div class="py-4 flex items-center justify-center">
        <ElButton v-debounce class="w-[96px] h-[36px]" type="primary" @click="handleConfirm">
          {{ $t('action_confirm') }}
        </ElButton>
        <ElButton class="w-[96px] h-[36px] text-[#1D1E1F]" type="info" plain @click.stop="close">
          {{ $t('action_cancel') }}
        </ElButton>
      </div>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
import { computed, nextTick, reactive, ref, shallowRef } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { CopyDocument } from '@element-plus/icons-vue'
import { ElIcon } from 'element-plus'
import { useEnterpriseStore } from '@/stores'
import { copyToClip } from '@/utils/copy'
import { api_host } from '@/utils/config'
import { generateFormRules } from '@/utils/form-rule.v2'

import { PROVIDER_VALUE } from '@/constants/platform/provider'
import type { ProviderValueType } from '@/constants/platform/provider'

import providersApi from '@/api/modules/providers'
import type { ProviderCreateRequest } from '@/api/modules/providers/types'

// 类型定义
interface AuthForm extends Omit<ProviderCreateRequest, 'configs'> {
  configs: Record<string, string>
}

interface ProviderData {
  name?: string
  provider_type?: ProviderValueType
  provider_id?: number
  configs?: Record<string, string>
  base_url?: string
  is_authorized?: number
  access_token?: string
  label?: string
  id?: number
}

interface PlatformConfig {
  i18n_key: string
  tip: {
    url: string
    needRedirectUrl?: boolean
  }
  fields: {
    label: string
    prop: string
    placeholder: string
  }[]
  setFormData?: (form: AuthForm) => void
  needsConfirmation?: boolean
  getAuthUrl?: (form: AuthForm, redirectUrl: string, provider_id: number) => string
}

// 事件定义
const emits = defineEmits<{
  (e: 'success'): void
}>()

// 状态管理
const enterpriseStore = useEnterpriseStore()
const copyRef = shallowRef<InstanceType<typeof ElIcon> | null>(null)
const formRef = shallowRef<FormInstance>()
const visible = ref(false)
const form = reactive<AuthForm>({
  name: '',
  configs: {},
  base_url: '',
  access_token: '',
  provider_type: 0,
})
const originData = ref<ProviderData>({})

// 平台配置
const PLATFORM_CONFIGS: Record<ProviderValueType, PlatformConfig> = {
  [PROVIDER_VALUE.APP_BUILDER]: {
    i18n_key: 'platform_auth.app_builder.tip',
    tip: { url: 'https://qianfan.cloud.baidu.com/appbuilder' },
    fields: [
      {
        label: window.$t('module.platform_tool_api_key'),
        prop: 'access_token',
        placeholder: window.$t('module.platform_tool_api_key_placeholder'),
      },
    ],
  },
  [PROVIDER_VALUE.COZE_CN]: {
    i18n_key: 'platform_auth.coze_cn.tip',
    tip: {
      url: 'https://www.coze.cn/open/oauth/apps',
      needRedirectUrl: true,
    },
    fields: [
      {
        label: window.$t('module.platform_auth_client_id'),
        prop: 'configs.client_id',
        placeholder: window.$t('module.platform_auth_client_id_placeholder'),
      },
      {
        label: window.$t('module.platform_auth_client_secret'),
        prop: 'configs.client_secret',
        placeholder: window.$t('module.platform_auth_client_secret_placeholder'),
      },
    ],
    needsConfirmation: true,
    setFormData: (form: AuthForm) => {
      form.access_token = ''
    },
    getAuthUrl: (form: AuthForm, redirectUrl: string, provider_id: number) =>
      `https://www.coze.cn/api/permission/oauth2/authorize?response_type=code&client_id=${form.configs.client_id}&redirect_uri=${encodeURIComponent(redirectUrl)}&state=${encodeURIComponent(`provider_id=${provider_id.toString()}`)}`,
  },
  [PROVIDER_VALUE.COZE_OSV]: {
    i18n_key: 'platform_auth.coze_osv.tip',
    tip: { url: 'https://www.53ai.com/' },
    fields: [
      {
        label: window.$t('module.platform_tool_api_endpoint'),
        prop: 'base_url',
        placeholder: window.$t('module.platform_model_base_url_placeholder'),
      },
      {
        label: window.$t('module.platform_tool_token'),
        prop: 'access_token',
        placeholder: window.$t('module.platform_tool_token_placeholder'),
      },
    ],
  },
  [PROVIDER_VALUE['53AI']]: {
    i18n_key: 'platform_auth.53ai.tip',
    tip: { url: 'https://www.53ai.com/' },
    fields: [
      {
        label: window.$t('module.platform_auth_url'),
        prop: 'base_url',
        placeholder: window.$t('module.platform_model_base_url_placeholder_53ai'),
      },
      {
        label: window.$t('module.platform_auth_secret'),
        prop: 'access_token',
        placeholder: window.$t('module.platform_tool_api_key_placeholder'),
      },
    ],
    setFormData: (form: AuthForm) => {
      form.base_url = form.base_url.trim() || 'https://api.53ai.com'
      form.access_token = form.access_token.trim()
    },
  },
}

// 计算属性
const currentConfig = computed(() => {
  const provider_type = originData.value.provider_type
  return typeof provider_type === 'number' ? PLATFORM_CONFIGS[provider_type] : null
})

const isCozeCN = computed(() => originData.value.provider_type === PROVIDER_VALUE.COZE_CN)

const coze_auth_url = computed(() => {
  const enterprise_info = enterpriseStore.info
  return `${api_host}/api/callback/cozecn/auth/${enterprise_info.eid}`
})

const guideHtml = computed(() => {
  const config = currentConfig.value
  if (!config) return ''

  const tipParams: Record<string, string> = {
    url: `<a class='text-[#5A6D9E]' href='${config.tip.url}' target='_blank'>${config.tip.url}</a>`,
  }

  if (config.tip.needRedirectUrl) {
    tipParams.redirect_url = `<span class='text-[#F04F4D]'>${coze_auth_url.value}</span><span class='copy-hook'></span>`
    tipParams.client_id = `<span class='text-[#F04F4D]'>${window.$t('module.platform_auth_client_id')}</span>`
    tipParams.client_secret = `<span class='text-[#F04F4D]'>${window.$t('module.platform_auth_client_secret')}</span>`
  }

  return window.$t(config.i18n_key, tipParams)
})

const schemaOptions = computed(() => currentConfig.value?.fields || [])

const formRules = computed(() => {
  const rules: FormRules = {
    name: generateFormRules({ validator: ['required'] }),
  }

  currentConfig.value?.fields.forEach(field => {
    rules[field.prop] = generateFormRules({ validator: ['required'], message: field.placeholder })
  })

  return rules
})

// 表单值处理 - 简化逻辑
const getFormValue = (prop: string) => {
  const isConfigProp = prop.startsWith('configs.')
  const key = isConfigProp ? prop.replace('configs.', '') : prop
  return isConfigProp ? form.configs[key] || '' : (form as any)[key] || ''
}

const setFormValue = (prop: string, value: string) => {
  const isConfigProp = prop.startsWith('configs.')
  const key = isConfigProp ? prop.replace('configs.', '') : prop

  if (isConfigProp) {
    form.configs[key] = value
  } else {
    ;(form as any)[key] = value
  }
}

// 方法
const reset = () => {
  Object.assign(form, {
    name: '',
    configs: {},
    base_url: '',
    access_token: '',
    provider_type: 0,
  })
}

const open = async ({ data = {} as ProviderData } = {}) => {
  reset()
  originData.value = data

  // 填充表单数据
  Object.assign(form, {
    name: data.name || '',
    provider_type: data.provider_type || 0,
    base_url: data.base_url || '',
    access_token: data.access_token || '',
  })
  Object.assign(form.configs, data.configs || {})

  // 应用平台特定的表单数据处理
  currentConfig.value?.setFormData?.(form)

  visible.value = true

  // 处理 Coze CN 的复制功能
  if (isCozeCN.value) {
    await nextTick()
    const copy_hook_el = formRef.value?.$el.querySelector('.copy-hook')
    if (copyRef.value?.$el && copy_hook_el) {
      copy_hook_el.appendChild(copyRef.value.$el)
    }
  }
}

const close = () => {
  visible.value = false
  reset()
}

const handleCopy = async (text: string) => {
  if (!text) return
  await copyToClip(text)
  ElMessage.success(window.$t('action_copy_success'))
}

const handleAuthorization = async (auth_url: string, provider_type: ProviderValueType) => {
  const auth_window = window.open(auth_url, '_blank', 'width=1000,height=800')

  window.addEventListener('message', ({ data = {} }) => {
    if (data.provider_type === provider_type) {
      auth_window?.close()
      ElMessage.success(window.$t('action_authorize_success'))
      emits('success')
    }
  })
}

const handleConfirm = async () => {
  const valid = await formRef.value?.validate()
  if (!valid) return

  const config = currentConfig.value
  if (!config) return

  const data: ProviderCreateRequest = {
    name: form.name,
    provider_type: form.provider_type,
    configs: JSON.stringify(form.configs),
    base_url: form.base_url,
    access_token: form.access_token,
  }

  // 检查是否需要确认
  if (config.needsConfirmation) {
    await ElMessageBox.confirm(window.$t('module.platform_auth_coze_confirm'), window.$t('tip'))
  }
  let provider_id = originData.value.provider_id
  // 保存数据
  if (provider_id) {
    await providersApi.update(provider_id, data)
  } else {
    const result = await providersApi.create(data)
    provider_id = result.provider_id
  }

  // 处理授权流程
  if (config.getAuthUrl) {
    const auth_url = config.getAuthUrl(form, coze_auth_url.value, provider_id)
    console.log(auth_url)
    await handleAuthorization(auth_url, originData.value.provider_type as ProviderValueType)
  }

  ElMessage.success(window.$t('action_save_success'))
  emits('success')
  close()
}

// 暴露方法
defineExpose({
  open,
  close,
  reset,
})
</script>

<style scoped>
.el-dialog {
  --el-dialog-padding-primary: 24px;
}
</style>
