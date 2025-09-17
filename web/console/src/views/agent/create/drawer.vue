<template>
  <ElDrawer
    v-model="visible"
    :title="editable ? $t('action_edit') : $t('action_add')"
    size="840px"
    destroy-on-close
    append-to-body
    :close-on-click-modal="false"
  >
    <AgentForm ref="agentFormRef" :agent-type="agentType" show-channel-config />
    <template #footer>
      <div class="flex border-t pt-5 justify-end w-full">
        <ElButton size="large" @click="handleClose">
          {{ $t('action_cancel') }}
        </ElButton>
        <ElButton v-debounce type="primary" size="large" @click="handleSave">
          {{ $t('action_confirm') }}
        </ElButton>
      </div>
    </template>
  </ElDrawer>
</template>

<script setup lang="ts">
import { nextTick, provide, ref } from 'vue'
import AgentForm from './platform/index.vue'

import { useAgentFormStore } from './store'

import { AGENT_TYPES, getAgentByAgentType } from '@/constants/platform/config'
import type { AgentType } from '@/constants/platform/config'

interface OpenParams {
  agent_type?: AgentType
  data?: {
    channel_config?: ChannelConfig
    label?: string
    value?: string
  }
  agent_id?: number
  group_id?: number
  cache?: boolean
}

const emits = defineEmits(['success', 'cancel'])
const agentFormStore = useAgentFormStore()

const agentFormRef = ref()
const visible = ref(false)
const agentType = ref<AgentType>(AGENT_TYPES.DIFY_AGENT)
const editable = ref(false)
const channelConfig = ref({})

provide('channelConfig', channelConfig)

// 加载渠道列表
const loadBotsList = () => {
  switch (agentType.value) {
    case AGENT_TYPES.COZE_AGENT_CN:
      agentFormStore.loadCozeWorkspaceOptions()
      break
    case AGENT_TYPES.APP_BUILDER:
      agentFormStore.loadAppBuilderBotOptions()
      break
    case AGENT_TYPES['53AI_AGENT']:
    case AGENT_TYPES['53AI_WORKFLOW']:
      agentFormStore.load53aiAppOptions()
      break
  }
}

// 编辑页过来时，不重置状态，直接打开
async function open({ agent_type, data = {}, agent_id, group_id = 0, cache = false }: OpenParams = {}) {
  // Update form state
  agentType.value = agent_type || AGENT_TYPES.DIFY_AGENT
  editable.value = !!+agent_id
  // Update channel config
  channelConfig.value = data.channel_config || {}
  channelConfig.value.name = channelConfig.value.name || data.label || ''

  if (!channelConfig.value.channel_type && data.value) channelConfig.value.channel_type = data.value

  if (cache) {
    loadBotsList()
  } else {
    agentFormStore.resetState()
    await nextTick()
    agentFormStore.agent_id = +agent_id || 0
    agentFormStore.agent_type = agentFormStore.form_data.custom_config.agent_type = agentType.value
    agentFormStore.form_data.logo = getAgentByAgentType(agentType.value).icon
    agentFormStore.form_data.group_id = group_id || 0
    loadBotsList()
    await agentFormStore.loadDetailData()
  }
  // Load necessary data
  agentFormStore.loadGroupOptions()
  agentFormStore.loadChannelOptions()

  visible.value = true
}

const handleClose = () => {
  visible.value = false
  emits('cancel')
}

const handleSave = async () => {
  const comp_ref = agentFormRef.value
  if (comp_ref && comp_ref.validateForm) {
    const valid = await comp_ref.validateForm()
    if (!valid) return Promise.reject()
    await comp_ref.onChannelSave?.()
    await agentFormStore.saveAgentData()
    emits('success', {
      agent_id: agentFormStore.agent_id,
      agent_type: agentFormStore.agent_type,
      action: editable.value ? 'update' : 'create',
    })
  }

  handleClose()
}

defineExpose({ open, close: handleClose })
</script>

<style scoped></style>
