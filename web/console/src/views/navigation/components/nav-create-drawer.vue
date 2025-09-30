<template>
  <ElDrawer
    v-model="visible"
    size="700px"
    :title="$t(isEditable ? 'action_edit' : 'action_create')"
    :close-on-click-modal="false"
    append-to-body
    destroy-on-close
    @close="close"
  >
    <ElForm ref="formRef" class="px-4" :model="formData" label-position="top">
      <h1 class="font-semibold text-[#1D1E1F] mb-6">{{ $t('basic_info') }}</h1>
      <ElFormItem :label="$t('type')">
        <ElRadioGroup v-model="formData.type" size="large" @change="handleTypeChange">
          <ElRadio
            v-for="value in [NAVIGATION_TYPE.SYSTEM, NAVIGATION_TYPE.EXTERNAL, NAVIGATION_TYPE.CUSTOM]"
            :key="value"
            :value="value"
            :disabled="formData.type === NAVIGATION_TYPE.SYSTEM || value === NAVIGATION_TYPE.SYSTEM || isEditable"
          >
            {{ $t(NAVIGATION_TYPE_LABEL_MAP.get(value) || '') }}
          </ElRadio>
        </ElRadioGroup>
      </ElFormItem>
      <ElFormItem :label="$t('icon')" prop="icon" :rules="[{ required: true, message: $t('form_select_placeholder') }]">
        <div class="size-12 border border-gray-200 rounded flex items-center justify-center">
          <ElImage class="size-6 overflow-hidden" :src="formData.icon" fit="contain" />
        </div>
        <ElPopover
          v-if="formData.type !== NAVIGATION_TYPE.SYSTEM"
          ref="iconPopoverRef"
          placement="bottom"
          :width="228"
          trigger="click"
          :show-arrow="false"
          popper-class="icon-popover"
        >
          <template #reference>
            <ElButton type="text" class="ml-4">{{ $t('action_modify') }}</ElButton>
          </template>
          <template #default>
            <div class="w-[228px] flex flex-wrap items-center">
              <ElButton
                v-for="item in iconList"
                :key="item"
                class="size-10 border-none icon-btn"
                @click="selectIcon(item)"
              >
                <ElImage class="size-6 overflow-hidden" :src="item" fit="contain" />
              </ElButton>
            </div>
          </template>
        </ElPopover>
      </ElFormItem>
      <ElFormItem :label="$t('name')" prop="name" :rules="[{ required: true, message: $t('form_input_placeholder') }]">
        <ElInput
          v-model="formData.name"
          size="large"
          :maxlength="20"
          show-word-limit
          :placeholder="$t('form_input_placeholder')"
        />
      </ElFormItem>
      <ElFormItem
        class="is-required"
        :label="$t('jump_path')"
        prop="jump_path"
        :rules="[
          ...generateInputRules({
            message: 'form_input_placeholder',
            validator: ['text', formData.type === NAVIGATION_TYPE.EXTERNAL ? 'url' : 'path'],
          }),
          {
            validator: (rule, value, callback) => {
              if (
                formData.type == NAVIGATION_TYPE.CUSTOM &&
                navigationList.some(item => item.jump_path === value && item.navigation_id !== originData.navigation_id)
              ) {
                return callback(new Error($t('form_path_same_tip')))
              }
              callback()
            },
            trigger: 'blur',
          },
        ]"
      >
        <ElInput
          v-if="formData.type === NAVIGATION_TYPE.SYSTEM"
          :model-value="domainUrl + formData.jump_path"
          size="large"
          :placeholder="$t('form_input_placeholder')"
          disabled
        />
        <ElInput
          v-else-if="formData.type === NAVIGATION_TYPE.EXTERNAL"
          v-model="formData.jump_path"
          size="large"
          :placeholder="$t('form_input_placeholder')"
        />
        <ElInput v-else v-model="formData.jump_path" size="large" :placeholder="$t('form_input_placeholder')">
          <template #prepend>{{ domainUrl }}</template>
        </ElInput>
      </ElFormItem>
      <ElFormItem :label="$t('open_method')" prop="target">
        <ElRadioGroup v-model="formData.target" size="large">
          <ElRadio v-for="value in [NAVIGATION_TARGET.SELF, NAVIGATION_TARGET.BLANK]" :key="value" :value="value">
            {{ $t(NAVIGATION_TARGET_LABEL_MAP.get(value) || '') }}
          </ElRadio>
        </ElRadioGroup>
      </ElFormItem>
      <ElDivider />
      <h1 class="font-semibold text-[#1D1E1F] mb-6">{{ $t('module.nav_seo_setting') }}</h1>
      <ElFormItem :label="$t('module.nav_seo_setting_title')">
        <ElInput
          v-model="formData.seo_title"
          maxlength="60"
          show-word-limit
          size="large"
          :placeholder="$t('form_input_placeholder')"
        />
      </ElFormItem>
      <ElFormItem :label="$t('module.nav_seo_setting_keywords')">
        <ElInput v-model="formData.seo_keywords" size="large" :placeholder="$t('form_input_placeholder')" />
        <div class="mt-2 text-xs text-[#999]">{{ $t('module.nav_seo_setting_keywords_tip') }}</div>
      </ElFormItem>
      <ElFormItem :label="$t('module.nav_seo_setting_description')">
        <ElInput
          v-model="formData.seo_description"
          type="textarea"
          :rows="5"
          maxlength="100"
          show-word-limit
          size="large"
          resize="none"
          :placeholder="$t('form_input_placeholder')"
        />
      </ElFormItem>
    </ElForm>
    <template #footer>
      <div class="flex border-t pt-5 justify-end w-full">
        <ElButton size="large" @click="close">
          {{ $t('action_cancel') }}
        </ElButton>
        <ElButton type="primary" size="large" :loading="isSubmitting" @click="handleSave">
          {{ $t('action_save') }}
        </ElButton>
      </div>
    </template>
  </ElDrawer>
</template>

<script setup lang="ts">
import { reactive, ref, computed, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { navigationApi } from '@/api/modules/navigation'
import { generateInputRules } from '@/utils/form-rule'
import { useEnterpriseStore } from '@/stores'
import { useEnv } from '@/hooks/useEnv'
import { img_host } from '@/utils/config'

import {
  NAVIGATION_TYPE,
  NAVIGATION_TARGET,
  NAVIGATION_TYPE_LABEL_MAP,
  NAVIGATION_TARGET_LABEL_MAP,
  type NavigationType,
  type NavigationTarget,
} from '@/constants/navigation'
import type { NavigationItem, CreateNavigationData, UpdateNavigationData } from '@/api/modules/navigation/types'

// 类型定义
type NavigationFormData = {
  type: number
  name: string
  icon: string
  jump_path: string
  target: number
  seo_title: string
  seo_keywords: string
  seo_description: string
}

// 事件定义
const emits = defineEmits<{
  success: [data: { data: NavigationFormData }]
}>()

// 状态管理
const enterpriseStore = useEnterpriseStore()
const router = useRouter()
const { isOpLocalEnv } = useEnv()

// 组件引用
const formRef = ref()

// 响应式数据
const visible = ref(false)
const isSubmitting = ref(false)
const originData = ref<Partial<NavigationItem>>({})
const navigationList = ref<NavigationItem[]>([])

// 表单数据
const formData = reactive<NavigationFormData>({
  type: NAVIGATION_TYPE.EXTERNAL,
  name: '',
  icon: '',
  jump_path: '',
  target: NAVIGATION_TARGET.SELF,
  seo_title: '',
  seo_keywords: '',
  seo_description: '',
})

// 计算属性
const enterpriseInfo = computed(() => enterpriseStore.info)
const domainUrl = computed(() => `${isOpLocalEnv.value ? window.location.origin : enterpriseInfo.value.domain}/#`)
const isEditable = computed(() => !!originData.value.navigation_id)

// 工具方法
const resetForm = () => {
  Object.assign(formData, {
    type: NAVIGATION_TYPE.EXTERNAL,
    name: '',
    icon: '',
    jump_path: '',
    target: NAVIGATION_TARGET.SELF,
    seo_title: '',
    seo_keywords: '',
    seo_description: '',
  })
}

const iconPopoverRef = ref()
const iconList: string[] = []
for (let i = 1; i <= 12; i++) {
  iconList.push(`${img_host}/navigation/icon${i}.png`)
}

const selectIcon = (icon: string) => {
  formData.icon = icon
  iconPopoverRef.value.hide()
}

// 公共方法
const open = async ({
  data = {} as Partial<NavigationItem>,
  navigationList: _navigationList = [],
}: {
  data?: Partial<NavigationItem>
  navigationList?: NavigationItem[]
} = {}) => {
  resetForm()
  await nextTick()

  const config = (data.config as any) || {}
  formData.type = +(data.type || NAVIGATION_TYPE.EXTERNAL)
  formData.name = data.name || ''
  formData.icon = (data as any).icon || `${img_host}/navigation/icon5.png`
  formData.jump_path = data.jump_path || ''
  formData.target = (data.target || config.target || NAVIGATION_TARGET.SELF) as any
  formData.seo_title = config.seo_title || ''
  formData.seo_keywords = config.seo_keywords || ''
  formData.seo_description = config.seo_description || ''
  originData.value = data
  navigationList.value = _navigationList
  visible.value = true
}

const close = () => {
  visible.value = false
}

// 事件处理方法
const handleTypeChange = () => {
  formData.jump_path = ''
  formRef.value?.clearValidate('jump_path')
}

const handleSave = async () => {
  const valid = await formRef.value?.validate()
  if (!valid) return

  isSubmitting.value = true

  try {
    const saveData: CreateNavigationData | UpdateNavigationData = {
      ...(originData.value.navigation_id ? { navigation_id: String(originData.value.navigation_id) } : {}),
      type: formData.type as NavigationType,
      name: formData.name,
      jump_path: formData.jump_path,
      sort: originData.value.sort || 9999 - navigationList.value.length,
      config: {
        target: formData.target as NavigationTarget,
        seo_title: formData.seo_title,
        seo_keywords: formData.seo_keywords.replace(/，/g, ','),
        seo_description: formData.seo_description,
      },
      icon: formData.icon,
    }

    const result = await navigationApi.save(saveData)
    ElMessage.success(window.$t('action_save_success'))
    emits('success', { data: formData })
    close()

    // 如果是新建的自定义导航，跳转到页面编辑
    if (!isEditable.value && formData.type === NAVIGATION_TYPE.CUSTOM) {
      const navigationId = (result as any)?.navigation_id || originData.value.navigation_id
      if (navigationId) {
        router.push({
          name: 'NavigationWebSetting',
          params: {
            navigation_id: String(navigationId),
          },
        })
      }
    }
  } catch (error) {
    console.error('保存导航失败:', error)
  } finally {
    isSubmitting.value = false
  }
}

// 暴露方法
defineExpose({
  open,
  close,
  reset: resetForm,
})
</script>

<style scoped lang="scss">
.icon-btn {
  margin: 0 12px 12px 0;

  --el-button-hover-bg-color: rgb(243 244 245 / 100%);
}
</style>

<style lang="scss">
.icon-popover {
  margin-left: 36px;
}
</style>
