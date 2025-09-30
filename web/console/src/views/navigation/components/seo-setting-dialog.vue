<template>
  <ElDialog
    v-model="visible"
    :title="$t('module.nav_seo_setting')"
    :close-on-click-modal="false"
    width="720px"
    append-to-body
    @close="close"
  >
    <ElForm ref="formRef" :model="form" label-position="left" label-width="128px">
      <ElFormItem :label="$t('module.nav_seo_setting_title')">
        <ElInput
          v-model="form.title"
          maxlength="60"
          show-word-limit
          size="large"
          :placeholder="$t('form_input_placeholder')"
        />
      </ElFormItem>
      <ElFormItem :label="$t('module.nav_seo_setting_keywords')">
        <ElInput v-model="form.keywords" size="large" :placeholder="$t('form_input_placeholder')" />
        <div class="mt-2 text-xs text-[#999]">{{ $t('module.nav_seo_setting_keywords_tip') }}</div>
      </ElFormItem>
      <ElFormItem :label="$t('module.nav_seo_setting_description')">
        <ElInput
          v-model="form.description"
          type="textarea"
          :rows="5"
          size="large"
          resize="none"
          :placeholder="$t('form_input_placeholder')"
        />
      </ElFormItem>
    </ElForm>
    <template #footer>
      <div class="py-4 flex items-center justify-center">
        <ElButton class="w-[96px] h-[36px]" type="primary" @click="handleConfirm">{{ $t('action_confirm') }}</ElButton>
        <ElButton class="w-[96px] h-[36px] text-[#1D1E1F]" type="info" plain @click.stop="close">{{
          $t('action_cancel')
        }}</ElButton>
      </div>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'

// 类型定义
type SeoFormData = {
  title: string
  keywords: string
  description: string
}

type NavigationData = {
  seo_setting_info?: SeoFormData
}

// 组件引用
const formRef = ref()

// 响应式数据
const visible = ref(false)
const form = reactive<SeoFormData>({
  title: '',
  keywords: '',
  description: '',
})
const originData = ref<NavigationData>({})

// 默认SEO数据
const defaultSeoData = {
  title: '快速成一个产品研发创意的agent',
  keywords: '53AI，产品创意生成，智能体',
  description: '没有产品研发创意？使用最新的AI agent 助你快速完成产品创意工作',
}

// 工具方法
const resetForm = () => {
  Object.assign(form, {
    title: '',
    keywords: '',
    description: '',
  })
}

// 公共方法
const open = ({ data = {} }: { data?: NavigationData } = {}) => {
  const seoSettingInfo = data.seo_setting_info || {}
  form.title = seoSettingInfo.title || defaultSeoData.title
  form.keywords = seoSettingInfo.keywords || defaultSeoData.keywords
  form.description = seoSettingInfo.description || defaultSeoData.description
  originData.value = data
  visible.value = true
}

const close = () => {
  visible.value = false
  resetForm()
}

// 事件处理方法
const handleConfirm = () => {
  formRef.value?.validate((valid: boolean) => {
    if (!valid) return

    // 确保seo_setting_info对象存在
    if (!originData.value.seo_setting_info) {
      originData.value.seo_setting_info = {}
    }

    // 更新SEO设置信息
    originData.value.seo_setting_info.title = form.title
    originData.value.seo_setting_info.keywords = form.keywords
    originData.value.seo_setting_info.description = form.description

    ElMessage.success(window.$t('action_save_success'))
    close()
  })
}

// 暴露方法
defineExpose({
  open,
  close,
  reset: resetForm,
})
</script>

<style scoped lang="scss"></style>
