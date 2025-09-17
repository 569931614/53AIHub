<template>
  <div :class="[showChannelConfig ? '' : 'py-7']">
    <ElForm ref="form_ref" :model="store.form_data" label-width="104px" label-position="top">
      <template v-if="showChannelConfig">
        <div class="text-base text-[#1D1E1F] font-medium mb-3">
          {{ $t('agent_app.coze_agent_cn') }}
        </div>
        <AgentType
          v-model="store.agent_type"
          :disabled="store.agent_id"
          :options="agentTypeOptions"
          @change="handleAgentTypeChange"
        />

        <el-form-item
          v-if="store.agent_type === AGENT_TYPES.COZE_WORKFLOW_CN"
          class="mb-9"
          prop="custom_config.coze_bot_url"
          :label="$t('agent.coze.workflow_link')"
          :rules="generateInputRules({ message: 'form_link_validator', validator: ['link'] })"
        >
          <el-input v-model="store.form_data.custom_config.coze_bot_url" size="large" />
        </el-form-item>
        <div v-else-if="store.agent_type === AGENT_TYPES.COZE_AGENT_CN" class="flex items-center gap-4 mb-9">
          <ElFormItem
            class="flex-1 mb-0"
            :label="$t('agent.coze.workspace')"
            prop="custom_config.coze_workspace_id"
            :rules="generateInputRules({ message: 'form_select_placeholder' })"
          >
            <SelectPlus
              v-model="store.form_data.custom_config.coze_workspace_id"
              size="large"
              :options="store.coze_workspace_options"
            />
          </ElFormItem>
          <ElFormItem
            class="flex-1 mb-0"
            :label="$t('agent.name')"
            prop="custom_config.coze_bot_id"
            :rules="generateInputRules({ message: 'form_select_placeholder' })"
          >
            <SelectPlus
              v-model="store.form_data.custom_config.coze_bot_id"
              size="large"
              :options="store.coze_bot_options"
              @change="onBotChange"
            />
          </ElFormItem>
        </div>
        <div class="text-base text-[#1D1E1F] font-medium mb-4">
          {{ $t('agent.base_info') }}
        </div>
        <AgentInfo v-model="store.form_data" />
      </template>
      <template v-if="!showChannelConfig">
        <template v-if="store.agent_type === AGENT_TYPES.COZE_WORKFLOW_CN">
          <FieldInput
            v-model:list="store.form_data.settings.input_fields"
            :title="$t('agent.input_variable')"
            allow-add
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
      <!-- <div class="border-t mb-4" />
      <LimitConfig /> -->
    </ElForm>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import AgentType from '../components/agent-type.vue'
import AgentInfo from '../components/agent-info.vue'
import BaseConfig from '../components/base-config.vue'
import ExpandConfig from '../components/expand-config.vue'
import UseScope from '../components/use-scope.vue'
import FieldInput from '../components/field-input.vue'
import RelateApp from '../components/relate-agents.vue'

import { useAgentFormStore } from '../store'
// import LimitConfig from '../components/limit-config.vue'
import { generateInputRules } from '@/utils/form-rule'

import { AGENT_TYPES } from '@/constants/platform/config'

defineProps({
  showChannelConfig: {
    type: Boolean,
    default: false,
  },
})

const store = useAgentFormStore()

const agentTypeOptions = [
  {
    icon: 'agent',
    label: window.$t('agent.coze.agent_type_chat'),
    description: window.$t('agent.coze.agent_type_chat_desc'),
    value: AGENT_TYPES.COZE_AGENT_CN,
  },
  {
    icon: 'completion-agent',
    label: window.$t('agent.coze.agent_type_workflow'),
    description: window.$t('agent.coze.agent_type_workflow_desc'),
    value: AGENT_TYPES.COZE_WORKFLOW_CN,
  },
]

const form_ref = ref()

const validateForm = async () => form_ref.value.validate()

const handleAgentTypeChange = (value: string) => {
  if (value === AGENT_TYPES.COZE_AGENT_CN) {
    store.form_data.custom_config.coze_workspace_id = store.coze_workspace_options[0].value
    store.form_data.custom_config.coze_bot_id = store.coze_bot_options[0].value
  } else {
    store.form_data.custom_config.coze_workspace_id = ''
    store.form_data.custom_config.coze_bot_id = ''
  }
}

const onBotChange = (data: { value: string; option: any }) => {
  if (!store.agent_id) {
    store.form_data.logo = data.option.icon
    store.form_data.name = data.option.label
    store.form_data.description = data.option.description || ''
  }
  // store.form_data.custom_config.coze_bot_id = value.value
}

defineExpose({
  validateForm,
})

watch(
  () => store.form_data.custom_config.coze_workspace_id,
  workspace_id => {
    if (workspace_id) {
      store.loadCozeBotOptions(workspace_id)
    }
  },
  { immediate: true }
)
</script>

<style scoped></style>
