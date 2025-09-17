<template>
  <Completion
    v-if="currentAgent?.custom_config_obj?.agent_mode === 'completion'"
    ref="chatRef"
    :use-case-fixed="useCaseFixed"
    :hide-menu-header="hideMenuHeader"
  />
  <Chat
    v-else
    ref="chatRef"
    :key="currentConv.virtual_id || currentConv.conversation_id"
    class="flex-1"
    :hide-menu-header="hideMenuHeader"
    :show-recommend="showRecommend"
    :use-case-fixed="useCaseFixed"
    :show-history="showHistory"
  />
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, defineAsyncComponent, computed, watch } from 'vue'
import { useRoute } from 'vue-router'

import { useAgentStore } from '@/stores/modules/agent'
import { useConversationStore } from '@/stores/modules/conversation'
import eventBus from '@/utils/event-bus'
import { EVENT_NAMES } from '@/constants/events'

const Chat = defineAsyncComponent(() => import('./chat/index.vue'))
const Completion = defineAsyncComponent(() => import('./completion/index.vue'))

const route = useRoute()
const agentStore = useAgentStore()
const convStore = useConversationStore()

const chatRef = ref<InstanceType<typeof Chat>>()

withDefaults(
  defineProps<{
    hideMenuHeader?: boolean
    hideFooter?: boolean
    showRecommend?: boolean
    useCaseFixed?: boolean
    showHistory?: boolean
  }>(),
  {
    hideMenuHeader: false,
    hideFooter: false,
    showRecommend: false,
    useCaseFixed: false,
    showHistory: false
  }
)

const currentConv = computed(() => convStore.currentConversation)
const currentAgent = computed(() => convStore.currentAgent)

// 监听路由参数变化
watch(
  () => route.query,
  async (newQuery) => {
    const agent_id = Number(newQuery.agent_id || '0')
    const conversation_id = newQuery.conversation_id || ''

    const agentList = await agentStore.loadAgentList()
    const agent = agentList.find((item) => item.agent_id === agent_id)

    if (agent) {
      convStore.pushUsualAgent(agent)
    } else {
      const usableAgent = convStore.usual_agents[0] || agentList[0]
      if (usableAgent) {
        convStore.pushUsualAgent(usableAgent)
      }
    }

    convStore.updateAgents(agentList)

    if (!conversation_id) {
      convStore.setCurrentState(agent_id, 0)
    } else {
      const convList = await convStore.loadConversations()
      const conversation = convList.find((item) => item.conversation_id === +conversation_id)
      if (conversation) {
        convStore.setCurrentState(conversation.agent_id, conversation.conversation_id)
      }
    }
  },
  { immediate: true }
)

onMounted(async () => {
  convStore.setBasePath(route.path.includes('/index') ? '/index/chat' : '/chat')
  agentStore.loadCategorys()

  convStore.loadConversations()
  eventBus.on(EVENT_NAMES.LOGIN_SUCCESS, () => {
    convStore.loadConversations()
  })
})

onUnmounted(() => {
  // 清理逻辑
  convStore.clearCurrentState()
})

defineExpose({
  detailData: currentAgent,
  showUseCase: () => {
    chatRef.value?.showUseCase()
  },
  hideUseCase: () => {
    chatRef.value?.hideUseCase()
  },
  showShare: () => {
    chatRef.value?.showShare()
  }
})
</script>

<style scoped>
.sidebar-transition {
  transition: width 0.3s ease-in-out;
}

.list-move,
.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
