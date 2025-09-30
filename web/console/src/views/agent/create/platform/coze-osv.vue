<template>
  <div :class="[showChannelConfig ? '' : 'py-7']">
    <template v-if="showChannelConfig">
      <div class="flex items-center justify-between mb-3">
        <div class="flex items-center gap-1">
          <h3 class="text-base text-[#1D1E1F]">
            {{ $t('agent_app.coze_agent_cn') }}
          </h3>

          <!-- <div class="flex-center text-[#9A9A9A] gap-1 ml-1" @click="handleOpenDialog">
            <svg-icon name="help" width="14" color="#999" />
            <span class="text-sm">{{ $t('how_get') }}</span>
          </div> -->
        </div>
      </div>
      <ElForm ref="channelFormRef" :model="channelForm" label-position="top" class="mt-3">
        <el-form-item :label="$t('module.website_info_name')">
          <el-select v-model="store.form_data.custom_config.provider_id" size="large">
            <el-option v-for="item in providers" :key="item.provider_id" :label="item.name" :value="item.provider_id" />
          </el-select>
        </el-form-item>

        <AgentType v-model="store.agent_type" :disabled="Boolean(store.agent_id)" :options="agentTypeOptions" />

        <el-form-item
          class="mb-9"
          prop="base_url"
          :label="$t('agent.coze.agent_link')"
          :rules="generateInputRules({ message: 'form_link_validator', validator: ['link'] })"
        >
          <el-input v-model="channelForm.base_url" size="large" />
        </el-form-item>
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
        <template v-if="store.agent_type === AGENT_TYPES.COZE_WORKFLOW_OSV">
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
    </ElForm>
  </div>
</template>

<script setup lang="ts">
import { inject, reactive, ref, watch, onMounted } from 'vue'
import AgentInfo from '../components/agent-info.vue'
import BaseConfig from '../components/base-config.vue'
import ExpandConfig from '../components/expand-config.vue'
import FieldInput from '../components/field-input.vue'
import UseScope from '../components/use-scope.vue'
import RelateApp from '../components/relate-agents.vue'
import AgentType from '../components/agent-type.vue'

import { channelApi } from '@/api/modules/channel'
import providersApi from '@/api/modules/providers/index'
import { transformProviderList } from '@/api/modules/providers/transform'
import { ProviderItem } from '@/api/modules/providers/types'

import { useAgentFormStore } from '../store'
import { generateInputRules } from '@/utils/form-rule'

import { AGENT_MODES, AGENT_TYPES, getAgentByAgentType, PROVIDER_VALUES } from '@/constants/platform/config'

const props = defineProps({
  showChannelConfig: {
    type: Boolean,
    default: false,
  },
})

const store = useAgentFormStore()

const channelFormRef = ref()

const agentTypeOptions = [
  {
    icon: 'agent',
    label: window.$t('agent.coze.agent_type_chat'),
    description: window.$t('agent.coze.agent_type_chat_desc'),
    value: AGENT_TYPES.COZE_AGENT_OSV,
  },
  {
    icon: 'completion-agent',
    label: window.$t('agent.coze.agent_type_workflow'),
    description: window.$t('agent.coze.agent_type_workflow_desc'),
    value: AGENT_TYPES.COZE_WORKFLOW_OSV,
  },
]

const channelInfo = inject('channelConfig') || {}
const providers = ref<ProviderItem[]>([])

const channelEditable = ref(false)
const channelForm = reactive({
  key: '',
  base_url: '',
  models: [],
  model: '',
  config: {
    agent_type: AGENT_MODES.COMPLETION,
  },
})

const agentFormRef = ref()

const onChannelSave = async () => {
  const valid = await channelFormRef.value.validate()
  if (!valid) return
  const agent = getAgentByAgentType(store.agent_type)

  if (agent.mode === 'completion') {
    try {
      const url = new URL(channelForm.base_url)
      const params = new URLSearchParams(url.search)
      channelForm.model = `workflow-${params.get('workflow_id')}` || ''
    } catch (error) {
      console.warn('Invalid URL format:', channelForm.base_url)
      channelForm.model = ''
    }
  } else {
    const model = channelForm.base_url.split('/').pop()
    channelForm.model = `bot-${model}`
  }

  const models = [channelForm.model]
  const name = 'coze_osv'
  const saveData = {
    channel_id: channelInfo.value.channel_id,
    key: channelForm.key,
    base_url: channelForm.base_url,
    config: channelForm.config,
    models,
    name,
  }
  const resultData = await channelApi.save({
    data: saveData,
  })
  Object.assign(channelInfo.value, resultData)
  if (!saveData.channel_id) saveData.channel_id = resultData.channel_id
  store.form_data.custom_config.channel_config = saveData
  store.form_data.model = models[0]
  channelEditable.value = true
}

const validateForm = async () => {
  channelFormRef.value && channelFormRef.value.validate()
  return agentFormRef.value && agentFormRef.value.validate()
}

const loadProviders = async () => {
  const list = await providersApi.list({
    providerType: PROVIDER_VALUES.COZE_OSV,
  })
  providers.value = transformProviderList(list)

  if (providers.value.length && !store.form_data.custom_config.provider_id) {
    store.form_data.custom_config.provider_id = providers.value[0].provider_id
  }
}

onMounted(() => {
  if (props.showChannelConfig) {
    loadProviders()
  }
})

watch(
  () => store.agent_data,
  ({ channel_config = {} } = {}) => {
    channelEditable.value = !!+channel_config.channel_id
    channelInfo.value.channel_id = +channel_config.channel_id || 0
    channelInfo.value.key = channelForm.key = channel_config.key || ''
    channelInfo.value.base_url = channelForm.base_url = channel_config.base_url || ''
    channelInfo.value.models = channelForm.models = channel_config.models || []
    channelInfo.value.model = channelForm.model = channelForm.models[0] || ''
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
