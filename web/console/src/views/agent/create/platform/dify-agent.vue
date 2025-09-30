<template>
  <div :class="[showChannelConfig ? '' : 'py-7']">
    <template v-if="showChannelConfig">
      <div class="flex items-center justify-between mb-3">
        <div class="flex items-center gap-1">
          <h3 class="text-base text-[#1D1E1F]">
            {{ $t('dify') }}
          </h3>
          <ElPopover content="Right Top prompts info" placement="right-start" width="480">
            <template #reference>
              <div class="flex-center text-[#9A9A9A] gap-1 ml-1">
                <svg-icon name="help" width="14" color="#999" />
                <span class="text-sm">{{ $t('how_get') }}</span>
              </div>
            </template>
            <div
              class="whitespace-pre-wrap text-sm text-[#333] leading-6"
              v-html="
                $t('dify_agent_get_tip', {
                  url: `<a class='text-[#5A6D9E] underline' href='https://dify.ai/zh' target='_blank'>https://dify.ai/zh</a>`,
                })
              "
            />
          </ElPopover>
        </div>
      </div>
      <AgentType v-model="store.agent_type" :disabled="!!store.agent_id" :options="agentTypeOptions" />
      <ElForm ref="channelFormRef" :model="channelForm" label-position="top" class="mt-3">
        <div class="flex items-center gap-4">
          <ElFormItem
            class="flex-1"
            :label="$t('api_host')"
            prop="base_url"
            :rules="generateInputRules({ message: 'form_input_placeholder', validator: ['text', 'link'] })"
          >
            <ElInput v-model="channelForm.base_url" size="large" :placeholder="$t('form_input_placeholder')" />
          </ElFormItem>
          <ElFormItem
            class="flex-1"
            :label="$t('api_screet')"
            prop="key"
            :rules="generateInputRules({ message: 'form_input_placeholder' })"
          >
            <ElInput v-model="channelForm.key" size="large" :placeholder="$t('form_input_placeholder')" />
          </ElFormItem>
        </div>
      </ElForm>
    </template>

    <ElForm ref="agentFormRef" :model="store.form_data" label-width="104px" label-position="top">
      <template v-if="showChannelConfig">
        <div class="text-base text-[#1D1E1F] font-medium mt-6 mb-4">
          {{ $t('basic_info') }}
        </div>
        <AgentInfo v-model="store.form_data" />
      </template>
      <template v-else>
        <template v-if="store.agent_type === AGENT_TYPES.DIFY_WORKFLOW">
          <FieldInput
            v-model:list="store.form_data.settings.input_fields"
            :title="$t('agent.input_variable')"
            allow-update
            allow-add
            :update-request="inputUpdateRequest"
            type="input"
            :agent-type="store.agent_type"
          />
          <FieldInput
            v-model:list="store.form_data.settings.output_fields"
            :title="$t('agent.output_variable')"
            allow-add
            type="output"
            :agent-type="store.agent_type"
          />
          <RelateApp />
        </template>
        <template v-else>
          <BaseConfig />
          <RelateApp />
          <ExpandConfig />
        </template>
        <UseScope />
      </template>
    </ElForm>
  </div>
</template>

<script setup lang="ts">
import { inject, reactive, ref, watch } from 'vue'
import AgentInfo from '../components/agent-info.vue'
import BaseConfig from '../components/base-config.vue'
import UseScope from '../components/use-scope.vue'
import FieldInput from '../components/field-input.vue'
import AgentType from '../components/agent-type.vue'
import RelateApp from '../components/relate-agents.vue'

import { useAgentFormStore } from '../store'
import { generateInputRules } from '@/utils/form-rule'
import { generateRandomId } from '@/utils'
import md5 from '@/utils/md5'

import { AGENT_TYPES, getAgentByAgentType } from '@/constants/platform/config'

import { channelApi } from '@/api/modules/channel'
import { agentApi } from '@/api'

defineProps({
  showChannelConfig: {
    type: Boolean,
    default: false,
  },
})

const store = useAgentFormStore()

const channelInfo = inject('channelConfig') || {}
const channelFormRef = ref()
const channelEditable = ref(false)
const channelForm = reactive({
  key: '',
  base_url: '',
  models: [],
  config: {
    agent_type: 'chat',
  },
})
const agentFormRef = ref()

const agentTypeOptions = [
  {
    icon: 'agent',
    label: window.$t('agent.dify.agent_type_chat'),
    description: window.$t('agent.dify.agent_type_chat_desc'),
    value: AGENT_TYPES.DIFY_AGENT,
  },
  {
    icon: 'completion-agent',
    label: window.$t('agent.dify.agent_type_workflow'),
    description: window.$t('agent.dify.agent_type_workflow_desc'),
    value: AGENT_TYPES.DIFY_WORKFLOW,
  },
]

const inputUpdateRequest = () => {
  return agentApi.dify.workflow_field_list(store.form_data.custom_config.channel_config.channel_id).then(res => {
    return res.user_input_form
      .map(item => {
        const type = Object.keys(item)[0]
        const value = Object.values(item)[0] as any
        if (!type) return null
        return {
          id: generateRandomId(6, true),
          variable: value.variable,
          type: type === 'paragraph' ? 'textarea' : type === 'select' ? 'select' : 'text',
          label: value.label,
          desc: value.desc,
          required: value.required,
          multiple: value.multiple || false,
          options: (value.options || []).map((item: string) => ({
            id: generateRandomId(6, true),
            label: item,
          })),
          max_length: value.max_length || 0,
          show_word_limit: value.show_word_limit || false,
          is_system: true,
        }
      })
      .filter(Boolean)
  })
}

const onChannelSave = async () => {
  const valid = await channelFormRef.value.validate()
  if (!valid) return
  const agent = getAgentByAgentType(store.agent_type)
  const model =
    (agent && agent.mode === 'completion' ? 'workflow-' : '') + md5(`${channelForm.key}_${channelForm.base_url}`)
  const name = 'dify'
  const saveData = {
    channel_id: channelInfo.value.channel_id,
    key: channelForm.key,
    base_url: channelForm.base_url,
    config: channelForm.config,
    models: [model],
    name,
  }
  const resultData = await channelApi.save({
    data: saveData,
  })
  Object.assign(channelInfo.value, resultData)
  if (!saveData.channel_id) saveData.channel_id = resultData.channel_id
  store.form_data.custom_config.channel_config = saveData
  store.form_data.model = model
  ElMessage.success(window.$t('action_save_success'))
  channelEditable.value = true
}

const validateForm = async () => {
  channelFormRef.value && channelFormRef.value.validate()
  if (agentFormRef.value) await agentFormRef.value.validate()
  return true
}

watch(
  () => store.agent_data,
  ({ channel_config = {} } = {}) => {
    channelEditable.value = !!+channel_config.channel_id
    channelInfo.value.channel_id = +channel_config.channel_id || 0
    channelInfo.value.key = channelForm.key = channel_config.key || ''
    channelInfo.value.base_url = channelForm.base_url = channel_config.base_url || 'https://api.dify.ai/v1'
    channelInfo.value.models = channelForm.models = channel_config.models || []
    channelInfo.value.config = channelForm.config = {
      ...(channel_config.config || {}),
      agent_type: channel_config.config?.agent_type || 'chat',
    }
  },
  { immediate: true, deep: true }
)

defineExpose({
  validateForm,
  onChannelSave,
})
</script>

<style scoped></style>
