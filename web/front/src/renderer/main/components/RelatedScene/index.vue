<template>
  <div v-if="currentAgent.settings_obj?.relate_agents?.length">
    <el-divider v-if="isWorkflow">
      <span class="text-sm text-secondary">{{ $t('chat.completion_next_action') }}</span>
    </el-divider>
    <div v-else class="flex items-center">
      <SvgIcon name="related" stroke="true" class="text-secondary" />
      <p class="pl-2 text-sm text-secondary">{{ $t('chat.completion_scene') }}</p>
    </div>
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 mt-3" :class="[isWorkflow ? 'p-4' : '']">
      <template v-for="item in currentAgent.settings_obj?.relate_agents" :key="item.scene">
        <div
          class="p-4 flex items-center gap-2 border rounded-md cursor-pointer hover:shadow-lg transition-all duration-300"
          @click="handleNextAgent(item)"
        >
          <img class="size-10 rounded-md" :src="item.logo" />
          <div class="flex-1 min-w-0">
            <h6 class="text-sm truncate">{{ item.name }}</h6>
            <p class="text-xs text-secondary truncate">{{ item.description }}</p>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick } from 'vue'

import { useConversationStore } from '@/stores/modules/conversation'
import SvgIcon from '../SvgIcon.vue'

const convStore = useConversationStore()

const currentAgent = computed(() => convStore.currentAgent)

interface OutputField {
  id: string
  label: string
  value: string
  variable: string
}

const props = withDefaults(
  defineProps<{
    isWorkflow: boolean
    // 应用型需要传数组，对话型传字符串
    output: OutputField[] | string
  }>(),
  {
    isWorkflow: false,
    output: () => []
  }
)

const emit = defineEmits<{
  (e: 'initAgent'): void
}>()

const getParameter = (): OutputField[] => {
  if (props.isWorkflow) return props.output as OutputField[]
  // 跟后台设置统一变量名称
  return [
    {
      id: 'output',
      label: '',
      value: props.output as string,
      variable: 'text'
    }
  ]
}

const handleNextAgent = (item) => {
  const agent = convStore.findAgentByAgentId(item.agent_id)
  if (agent) {
    const currentAgentId = convStore.current_agentid
    const parameters = getParameter()
    convStore.setNextAgentPrepare({
      agent_id: item.agent_id,
      execution_rule: item.execution_rule,
      // 兼容旧数据
      is_workflow: typeof item.is_workflow === 'boolean' ? item.is_workflow : true,
      parameters: Object.keys(item.field_mapping).reduce((acc, key) => {
        acc[key] = item.field_mapping[key].replace(/\{\#(.*?)\#\}/g, (match, p1) => {
          return parameters.find((item) => item.variable === p1)?.value || ''
        })
        return acc
      }, {})
    })
    convStore.pushUsualAgent(agent)
    convStore.setCurrentState(item.agent_id, 0)
    if (item.agent_id === currentAgentId) {
      nextTick(() => {
        emit('initAgent')
      })
    }
  } else {
    ElMessage.warning(window.$t('chat.no_available_agent'))
  }
}
</script>

<style scoped></style>
