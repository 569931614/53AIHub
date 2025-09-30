<template>
  <div :class="[showChannelConfig ? '' : 'py-7']">
    <ElForm ref="agentFormRef" :model="store.form_data" label-width="104px" label-position="top">
      <template v-if="showChannelConfig">
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center gap-1">
            <h3 class="text-base text-[#1D1E1F]">
              {{ $t('53ai') }}
            </h3>
          </div>
        </div>
        <el-form-item :label="$t('module.website_info_name')">
          <el-select v-model="store.form_data.custom_config.provider_id" size="large" @change="onProviderChange">
            <el-option v-for="item in providers" :key="item.provider_id" :label="item.name" :value="item.provider_id" />
          </el-select>
        </el-form-item>
        <AgentType
          v-model="store.agent_type"
          :disabled="Boolean(store.agent_id)"
          :options="agentTypeOptions"
          @change="onAgentTypeChange"
        />

        <el-form-item
          :label="$t('select_agent')"
          prop="custom_config.chat53ai_agent_id"
          :rules="generateInputRules({ message: 'form_select_placeholder' })"
        >
          <SelectPlus
            v-model="store.form_data.custom_config.chat53ai_agent_id"
            :use-i18n="false"
            size="large"
            :options="bots"
            @change="onBotChange"
          />
        </el-form-item>

        <div class="text-base text-[#1D1E1F] font-medium mt-6 mb-4">
          {{ $t('basic_info') }}
        </div>
        <AgentInfo v-model="store.form_data" />
      </template>
      <template v-else>
        <template v-if="store.agent_type === AGENT_TYPES['53AI_WORKFLOW']">
          <FieldInput
            v-model:list="store.form_data.settings.input_fields"
            :title="$t('agent.input_variable')"
            allow-update
            :update-request="inputUpdateRequest"
            type="input"
            :agent-type="store.agent_type"
          />
          <FieldInput
            v-model:list="store.form_data.settings.output_fields"
            :title="$t('agent.output_variable')"
            type="output"
            allow-add
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
import { ref, onMounted } from 'vue'
import AgentType from '../components/agent-type.vue'
import AgentInfo from '../components/agent-info.vue'
import RelateApp from '../components/relate-agents.vue'
import BaseConfig from '../components/base-config.vue'
import ExpandConfig from '../components/expand-config.vue'
import UseScope from '../components/use-scope.vue'
import FieldInput from '../components/field-input.vue'

import { useAgentFormStore } from '../store'
import { generateInputRules } from '@/utils/form-rule'
import { AGENT_TYPES, PROVIDER_VALUES } from '@/constants/platform/config'

import providersApi from '@/api/modules/providers/index'
import { transformProviderList } from '@/api/modules/providers/transform'
import { ProviderItem } from '@/api/modules/providers/types'
import agentApi, { BotItem53aiItem, transform53aiBotItem } from '@/api/modules/agent'

const props = defineProps({
  showChannelConfig: {
    type: Boolean,
    default: false,
  },
})

const store = useAgentFormStore()

const agentFormRef = ref()

const agentTypeOptions = [
  {
    icon: 'agent',
    label: window.$t('agent.53ai.agent_type_chat'),
    description: window.$t('agent.53ai.agent_type_chat_desc'),
    value: AGENT_TYPES['53AI_AGENT'],
  },
  {
    icon: 'completion-agent',
    label: window.$t('agent.53ai.agent_type_workflow'),
    description: window.$t('agent.53ai.agent_type_workflow_desc'),
    value: AGENT_TYPES['53AI_WORKFLOW'],
  },
]
const providers = ref<ProviderItem[]>([])
const bots = ref<BotItem53aiItem[]>([])

// 53ai 的 agent 选择后，需要设置 那边的开场白和建议问题
const onBotChange = (data: { value: string; option: any }) => {
  if (!store.agent_id) {
    store.form_data.logo = data.option.logo
    store.form_data.name = data.option.name
    store.form_data.description = data.option.description || ''
  }
  store.form_data.custom_config.chat53ai_agent_id = data.value
  store.form_data.settings.opening_statement = data.option.opening_statement
  store.form_data.settings.suggested_questions = data.option.suggested_questions.map(item => {
    return {
      id: Math.random().toString(36).substring(2, 15),
      content: item,
    }
  })
}

const inputUpdateRequest = async () => {
  const res = await agentApi.chat53ai.workflow_field_list(store.form_data.custom_config.chat53ai_agent_id)
  return res.user_input_form.map(item => {
    const value: any = Object.values(item)[0]
    return {
      id: value.id,
      variable: value.variable,
      type: value.type_53ai,
      label: value.label,
      desc: value.desc,
      required: value.required,
      multiple: value.multiple,
      options: value.options_53ai,
      max_length: value.max_length,
      show_word_limit: value.showWordLimit,
      date_format: value.mode || '',
      file_type: value.docType,
      file_accept: value.accept,
      file_size: value.size,
      file_limit: value.limit,
      is_system: true,
    }
  })
}

const load53aiBots = async () => {
  const list = await agentApi.chat53ai.bots_list({
    provider_id: store.form_data.custom_config.provider_id,
  })
  bots.value = list.map(transform53aiBotItem)
}
const load53aiWorkflows = async () => {
  const list = await agentApi.chat53ai.workflow_list({
    provider_id: store.form_data.custom_config.provider_id,
  })
  bots.value = list.map(transform53aiBotItem)
}

const loadApp = async () => {
  if (store.agent_type === AGENT_TYPES['53AI_AGENT']) {
    load53aiBots()
  } else {
    load53aiWorkflows()
  }
}

const onProviderChange = () => {
  loadApp()
  store.form_data.custom_config.chat53ai_agent_id = ''
}
const loadProviders = async () => {
  const list = await providersApi.list({
    providerType: PROVIDER_VALUES['53AI'],
  })
  providers.value = transformProviderList(list)

  const customConfig = store.form_data.custom_config
  if (providers.value.length && !customConfig.provider_id) {
    customConfig.provider_id = providers.value[0].provider_id
  }
  loadApp()
}

const onAgentTypeChange = () => {
  store.form_data.custom_config.chat53ai_agent_id = ''
  loadApp()
}

onMounted(() => {
  if (props.showChannelConfig) {
    loadProviders()
  }
})

defineExpose({
  validateForm() {
    return agentFormRef.value && agentFormRef.value.validate()
  },
})
</script>

<style scoped></style>
