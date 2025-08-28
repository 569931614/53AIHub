<template>
  <Layout class="px-[60px] py-8">
    <Header back :title="title || agentFormStore.form_data.name" class="mb-5">
      <template #title_prefix>
        <el-image v-if="agentFormStore.form_data.logo" :src="agentFormStore.form_data.logo" class="w-8 rounded" />
        <div v-else class="size-8 rounded" />
      </template>
      <template #title_suffix>
        <el-button class="!size-5 !p-0 flex-center" @click="handleEdit">
          <svg-icon name="edit" width="14" />
        </el-button>
      </template>
    </Header>
    <el-tabs v-model="active_tab_name" class="flex-1 agent-tabs el-tabs--full">
      <el-tab-pane :label="$t('app_config')" name="first" lazy>
        <div id="app-config-full-screen-hook" v-loading="agentFormStore.loading" class="relative h-full flex bg-white">
          <div class="flex-1 flex flex-col overflow-hidden">
            <AgentForm
              ref="agentFormRef"
              class="flex-1 py-7 px-4 overflow-y-auto"
              :agent-type="agentFormStore.agent_type"
            />

            <div class="border-t px-4 py-5">
              <el-button type="primary" size="large" :loading="agentFormStore.saving" @click="onSave">
                {{ $t('action_save') }}
              </el-button>
            </div>
          </div>
          <div class="flex-none w-px border-r" />
          <template v-if="agentFormStore.agent_option_data.mode === 'chat'">
            <ChatView ref="chatRef" class="flex-1 overflow-hidden" @save="onSave" />
          </template>
          <template v-else-if="agentFormStore.agent_option_data.mode === 'completion'">
            <CompletionView class="flex-1 py-7 overflow-hidden" />
          </template>
        </div>
      </el-tab-pane>
      <el-tab-pane :label="$t('usage_guide')" name="second" lazy>
        <GuideView>
          <template #footer>
            <div class="border-t px-4 py-5 mt-8 sticky bottom-0 left-0 right-0 bg-white z-10">
              <el-button type="primary" size="large" :loading="agentFormStore.saving" @click="onSave">
                {{ $t('action_save') }}
              </el-button>
            </div>
          </template>
        </GuideView>
      </el-tab-pane>
      <el-tab-pane v-if="!isWorkEnv" :label="$t('dialogue_record')" name="third" lazy>
        <DialogueRecordView type="agent" :related-id="agentFormStore.agent_id" />
      </el-tab-pane>
    </el-tabs>
    <InfoDrawer ref="infoDrawerRef" @success="eventBus.emit('agent-change')" @cancel="handleCancel" />
  </Layout>
</template>

<script setup lang="ts">
import { nextTick, onMounted, provide, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import AgentForm from './platform/index.vue'
import ChatView from './response/chat.vue'
import CompletionView from './response/completion.vue'
import GuideView from './guide.vue'
import InfoDrawer from './drawer.vue'
import { useAgentFormStore } from './store'
import DialogueRecordView from '@/components/DialogueRecord/index.vue'

import { AGENT_TYPE } from '@/constants/platform/agent'
import eventBus from '@/utils/event-bus'
import { useEnv } from '@/hooks/useEnv'
import { getAgentByAgentType } from '@/constants/platform/config'
import type { AgentType } from '@/constants/platform/config'

const route = useRoute()
const router = useRouter()
const agentFormStore = useAgentFormStore()
const { isWorkEnv } = useEnv()

const chatRef = ref<InstanceType<typeof ChatView>>()
const agentFormRef = ref<InstanceType<typeof AgentForm>>()
const infoDrawerRef = ref<InstanceType<typeof InfoDrawer>>()

const active_tab_name = ref('first')
const channelConfig = ref({})
const title = ref()
provide('channelConfig', channelConfig)

const onSave = async ({ restart = false } = {}) => {
  if (agentFormStore.saving) return
  const comp_ref = agentFormRef.value
  let agent_id = ''
  if (comp_ref && comp_ref.save) {
    agentFormStore.saving = true
    const { data = {} } = await comp_ref.save().catch(() => {
      agentFormStore.saving = false
    })
    agent_id = data.agent_id
  } else if (comp_ref && comp_ref.validateForm) {
    const valid = await comp_ref.validateForm()
    if (!valid) return Promise.reject()
    await agentFormStore.saveAgentData().catch(() => {
      agentFormStore.saving = false
    })
    agentFormStore.saving = true
    agent_id = agentFormStore.agent_id
  }
  eventBus.emit('agent-change', { agent_id })
  if (agent_id) await router.replace({ name: 'AgentCreate', query: { type: agentFormStore.agent_type, agent_id } })
  agentFormStore.saving = false
  if (chatRef.value) {
    if (restart || chatRef.value?.getIsConfigChanged()) chatRef.value?.restart()
  }
}

const handleEdit = () => {
  title.value = agentFormStore.form_data.name
  if (infoDrawerRef.value) {
    infoDrawerRef.value.open({
      agent_type: agentFormStore.agent_type,
      agent_id: agentFormStore.agent_id,
      data: {
        channel_config: channelConfig.value,
      },
      cache: true,
    })
  }
}

const handleCancel = () => {
  agentFormStore.form_data.name = title.value
}

const handleAgentChange = () => {
  // 重新加载数据以更新标题
  title.value = agentFormStore.form_data.name
}

onMounted(async () => {
  agentFormStore.resetState()
  eventBus.on('agent-change', handleAgentChange)
  await nextTick()
  const agent_type = (route.query.type as string) || 'prompt'
  agentFormStore.is_new = route.query.is_new === 'true'
  agentFormStore.agent_id = Number(route.query.agent_id as string)
  agentFormStore.agent_type = agent_type
  // 加载详情
  await agentFormStore.loadDetailData()
  // 加载渠道
  agentFormStore.loadChannelOptions()

  if (agent_type !== AGENT_TYPE.PROMPT) {
    const data = getAgentByAgentType(agent_type as AgentType)
    channelConfig.value.name = data.channelType
    channelConfig.value.channel_type = data.channelValue
  }
})
</script>

<style scoped>
.agent-tabs :deep(.el-tabs__nav-wrap) {
  /* padding: 0 40px; */
}
</style>
