<template>
  <ElDialog
    v-model="visible"
    :title="$t('action_setting') + (modelConfig?.label ? $t(modelConfig.label) : '')"
    :close-on-click-modal="false"
    width="600px"
    center
    destroy-on-close
    append-to-body
    @close="() => close()"
  >
    <ElForm ref="form_ref" :model="form" label-position="top" require-asterisk-position="right">
      <template v-for="config in modelSchemas" :key="config.prop">
        <ElFormItem
          v-if="config.type === 'input'"
          :label="$t(config.label)"
          :prop="config.prop"
          :required="config.required"
          :rules="generateInputRules({ message: 'form_input_placeholder' })"
        >
          <ElInput v-model="form[config.prop]" size="large" :placeholder="$t(config.placeholder || '')" />
        </ElFormItem>
        <ElFormItem
          v-else-if="config.type === 'url'"
          :label="$t(config.label)"
          :prop="config.prop"
          :required="config.required"
          :rules="generateInputRules({ message: 'form_input_placeholder', validator: ['text', 'link'] })"
        >
          <ElInput v-model="form[config.prop]" size="large" :placeholder="$t(config.placeholder || '')" />
        </ElFormItem>
        <ElFormItem
          v-else-if="config.type === 'select'"
          :label="$t(config.label)"
          :prop="config.prop"
          :required="config.required"
          :rules="generateInputRules({ message: 'form_select_placeholder' })"
        >
          <ul v-loading="loading" class="w-full max-h-[42vh] pr-1 flex flex-col gap-2 box-border overflow-auto">
            <div v-if="!modelOptions.length" class="text-[#9A9A9A] w-full text-center py-4">
              {{ $t('module.platform_model_models_empty') }}
            </div>
            <template v-for="opt in modelOptions" :key="opt.value">
              <div class="text-sm text-[#1D1E1F] text-opacity-60">{{ opt.model_type_name }}</div>
              <template v-for="item in opt.models" :key="item.model_id">
                <li class="w-full flex items-center gap-1.5">
                  <img v-if="item.icon" class="flex-none w-[20px] h-[20px] object-contain" :src="item.icon" />
                  <label class="text-sm text[#1D1E1F]">{{ item.label }}</label>
                  <div class="flex-1" />
                  <ElSwitch
                    :model-value="form.models.includes(item.value)"
                    size="small"
                    @change="handleModelChange(item.value)"
                  />
                </li>
              </template>
            </template>
          </ul>
        </ElFormItem>
      </template>
    </ElForm>
    <template #footer>
      <div class="pb-4 flex items-center justify-center">
        <ElButton v-debounce class="w-24 h-9" type="primary" @click="handleSave">
          {{ $t('action_save') }}
        </ElButton>
        <ElButton class="w-24 h-9 text-[#1D1E1F]" type="info" plain @click="close">
          {{ $t('action_cancel') }}
        </ElButton>
      </div>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
import { ref } from 'vue'

import type { FormInstance } from 'element-plus'
import { getFormConfig } from '@/constants/platform/model'
import type { FormConfig } from '@/constants/platform/model'

import { generateInputRules } from '@/utils/form-rule'

import type { ModelConfig } from '@/constants/platform/config'
import { getModelByChannelType } from '@/constants/platform/config'
import channelApi, { type ModelOption } from '@/api/modules/channel/index'
import { deepCopy } from '@/utils'

const props = withDefaults(
  defineProps<{
    modelList: ModelOption[]
  }>(),
  {
    modelList: () => [],
  }
)
const emit = defineEmits<{
  (e: 'success'): void
}>()
const form_ref = ref<FormInstance>()

const visible = ref(false)

const modelSchemas = ref<FormConfig[]>([])
const modelOptions = ref<ModelOption[]>([])
const modelConfig = ref<ModelConfig>({})

// 模型是否单选
const isSingleModel = ref(false)
const default_form = {
  type: 0,
  priority: 0,
  weight: 0,
  key: '',
  name: '',
  other: '',
  base_url: '',
  models: [],
  model_mapping: '',
  custom_config: '',
  config: {},
}
const form = ref<{
  [key: string]: any
}>({ ...default_form })

const loading = ref(false)

const clearForm = () => {
  form.value = { ...default_form }
}
const initForm = (data: any = {}) => {
  form.value.type = data.channel_type
  modelSchemas.value = getFormConfig(form.value.type)
  for (const schema of modelSchemas.value) {
    if (schema.type === 'select' && schema.prop === 'models') {
      form.value.models = (data.models || '').toString().split(',').filter(Boolean)
    } else if (schema.default) {
      form.value[schema.prop] = schema.default
    }
  }
}

const assignForm = (data: any = {}) => {
  if (!data.channel_id) return
  const models = Object.keys(data.custom_config).map(item => `${data.custom_config[item]}_${item}`)

  form.value.base_url = data.base_url
  form.value.key = data.key
  form.value.name = data.name
  form.value.other = data.other
  form.value.models = models

  form.value.weight = data.weight
  form.value.priority = data.priority
  form.value.config = data.config || {}
  form.value.model_mapping = data.model_mapping || ''
  form.value.channel_id = data.channel_id
}

const loadModelList = async () => {
  const models = props.modelList.find(item => item.channel_type === form.value.type)
  if (!models) return
  modelOptions.value = deepCopy(models.categories).map(item => {
    return {
      ...item,
      models: item.models.map(model => {
        return {
          ...model,
          value: `${item.model_type}_${model.model_id}`,
          label: model.model_id,
        }
      }),
    }
  })
}

const open = async (data: any = {}) => {
  modelConfig.value = getModelByChannelType(data.channel_type) || {}
  clearForm()
  initForm(data)
  assignForm(data)

  visible.value = true
  loadModelList()
}

const close = () => {
  visible.value = false
}

const handleModelChange = (value: string) => {
  if (form.value.models.includes(value)) form.value.models = form.value.models.filter(item => item !== value)
  else form.value.models.push(value)
}

const handleSave = () => {
  return form_ref.value?.validate().then(async () => {
    const data = deepCopy(form.value)
    const custom_config: Record<string, string> = {}
    const models: string[] = []
    if (isSingleModel.value) {
      const [model_type, model_id] = data.models.split('_')
      custom_config[model_id] = model_type
      models.push(model_id)
    } else {
      data.models.forEach((item: string) => {
        const [model_type, model_id] = item.split('_')
        custom_config[model_id] = model_type
        models.push(model_id)
      })
    }
    data.models = models.join(',')
    data.config = JSON.stringify(data.config)
    data.custom_config = JSON.stringify(custom_config)
    if (data.channel_id) {
      await channelApi.update(data.channel_id, data)
    } else {
      await channelApi.create(data)
    }

    ElMessage.success(window.$t('action_save_success'))
    emit('success')
    close()
  })
}

defineExpose({
  open,
  close,
})
</script>
