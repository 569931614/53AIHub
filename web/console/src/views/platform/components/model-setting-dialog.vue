<template>
  <ElDialog
    v-model="model_visible"
    :title="$t('module.platform_model_models_edit')"
    :close-on-click-modal="false"
    width="600px"
    destroy-on-close
    append-to-body
  >
    <ElForm ref="model_form_ref" :model="model_form" label-position="top">
      <ElFormItem
        :label="$t('module.platform_model_models_id')"
        prop="id"
        class="is-required"
        :rules="generateInputRules({ message: 'form_input_placeholder' })"
      >
        <ElInput v-model="model_form.id" size="large" disabled :placeholder="$t('form_input_placeholder')" />
      </ElFormItem>
      <ElFormItem :label="$t('module.platform_model_models_name')" prop="name">
        <ElInput v-model="model_form.name" size="large" :placeholder="$t('form_input_placeholder')" />
      </ElFormItem>
    </ElForm>
    <template #footer>
      <div class="pb-4 flex items-center justify-center">
        <ElButton
          class="w-24 h-9"
          type="primary"
          :loading="submitting || loading"
          @click="() => onSave({ action: 'model_edit' })"
        >
          {{ $t('action_save') }}
        </ElButton>
        <ElButton class="w-24 h-9 text-[#1D1E1F]" type="info" plain @click="() => close({ action: 'model_edit' })">
          {{ $t('action_cancel') }}
        </ElButton>
      </div>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useChannelStore } from '@/stores'
import { generateInputRules } from '@/utils/form-rule'

interface ModelData {
  id?: string
  value?: string
  label?: string
  icon?: string
  models?: string[]
  config?: {
    model_alias_map?: Record<string, string>
  }
  channel_id?: string
  key?: string
  base_url?: string
  type?: string
  other?: string
  name?: string
  organization_id?: string
}

const emits = defineEmits<{
  (e: 'success', result: { action: 'add' | 'edit' | 'model_edit'; data: any }): void
}>()

const channel_store = useChannelStore()
const model_form_ref = ref()
const loading = ref(false)
const submitting = ref(false)

const origin_data = ref<ModelData>({})
const model_visible = ref(false)
const model_form = reactive({
  id: '',
  name: '',
  channel_type: 0,
})

const open = async ({ data = {} as ModelData } = {}) => {
  model_form.id = data.id || data.value || ''
  model_form.name = data.label || data.name || ''
  model_form.channel_type = data.channel_type
  model_visible.value = true
  origin_data.value = data
}

const close = () => {
  model_visible.value = false
}

const onSave = () => {
  model_form_ref.value.validate(async (valid: boolean) => {
    if (!valid) return
    let { id = '', name = '', channel_type } = model_form
    name = name.trim()
    const config = origin_data.value.config || {}
    config.model_alias_map = {
      ...(config.model_alias_map || {}),
      [id]: name,
    }
    if (!config.model_alias_map[id]) delete config.model_alias_map[id]
    if (!Object.keys(config.model_alias_map).length) delete config.model_alias_map
    const data = {
      channel_id: origin_data.value.channel_id,
      config,
      key: origin_data.value.key,
      base_url: origin_data.value.base_url,
      models: origin_data.value.models,
      name: origin_data.value.name,
      other: origin_data.value.other,
      type: channel_type,
      custom_config: origin_data.value.custom_config ? JSON.stringify(origin_data.value.custom_config) : '',
    }
    submitting.value = true
    await channel_store.save({ data }).finally(() => {
      submitting.value = false
    })
    emits('success', { action: 'model_edit', data: { id, name: name || id } })
    ElMessage.success(window.$t('action_save_success'))
    close()
  })
}

defineExpose({
  open,
  close,
})
</script>

<style scoped></style>
