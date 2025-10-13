<template>
  <div class="flex flex-col min-h-full">
    <header v-if="isShareMode" class="flex-none sticky top-0 z-10 bg-white" :class="[showRecommend ? 'w-4/6' : 'border-b']">
      <div class="h-[70px] flex items-center justify-between" :class="[showRecommend ? 'w-[95%]' : 'w-11/12 md:w-4/5 max-w-[800px] mx-auto']">
        <el-checkbox
          :model-value="state.selectAll"
          :label="state.selectAll ? $t('action.unselect_all') : $t('action.select_all')"
          :value="true"
          size="large"
          @click="handleSelectAll"
        />
        <div class="flex items-center gap-2">
          <div
            v-if="state.selectMessageIds.length"
            v-debounce
            class="h-8 flex items-center gap-1 px-2 rounded-md bg-[#F5F5F7] cursor-pointer hover:bg-[#E1E2E3]"
            @click.stop="handleCreateShare"
          >
            <svg-icon name="unlink" size="16" stroke />
            <span class="text-sm text-primary">复制链接</span>
          </div>
          <div class="size-8 flex-center rounded-md bg-[#F5F5F7] cursor-pointer hover:bg-[#E1E2E3]" @click="handleOpenShare()">
            <el-icon>
              <Close />
            </el-icon>
          </div>
        </div>
      </div>
    </header>
    <MainHeader v-else-if="!hideMenuHeader">
      <template #before_suffix>
        <div class="text-base text-primary line-clamp-1 max-md:flex-1 max-md:text-center" :title="currentConv.title || currentAgent.name || ''">
          {{ currentConv.title || currentAgent.name || '' }}
        </div>
      </template>
      <template #after_prefix>
        <span class="flex items-center gap-1 text-sm cursor-pointer md:hidden" @click="$router.back()">
          <svg-icon name="return" size="18" stroke></svg-icon>
        </span>
        <div
          v-if="state.messageList.length"
          v-tooltip="{ content: $t('action.share') }"
          class="h-6 px-1 rounded flex-center gap-1 cursor-pointer hover:bg-[#E1E2E3]"
          @click.stop="handleOpenShare()"
        >
          <svg-icon size="18" name="share-two" color="#4F5052" stroke></svg-icon>
        </div>
        <div
          v-tooltip="{ content: $t('chat.usage_guide') }"
          class="h-6 px-1 rounded flex-center gap-1 cursor-pointer hover:bg-[#E1E2E3]"
          @click="handleToggleGuide"
        >
          <div class="size-4">
            <svg-icon name="layout-split" size="18"></svg-icon>
          </div>
        </div>
      </template>
    </MainHeader>

    <!-- 消息列表区域 -->
    <div class="flex-1 flex mt-5">
      <x-bubble-list
        ref="bubbleListRef"
        :auto-scroll="false"
        class="flex-1"
        :messages="state.messageList"
        :main-class="showRecommend ? 'w-[95%]' : 'w-11/12 md:w-4/5 max-w-[800px] mx-auto'"
        enable-pull-up
        @pull-up="handleLoadListMore"
      >
        <template v-if="currentAgent.settings_obj && !state.messageList.length" #header>
          <div
            class="w-full mt-2 flex items-center gap-3 box-border p-6 rounded-xl overflow-hidden"
            :style="{
              background: `linear-gradient(90deg, rgba(243, 249, 254, 1) 0%, rgba(247, 243, 255, 1) 100%)`
            }"
          >
            <img class="flex-none size-10 rounded-full overflow-hidden" :src="currentAgent.logo" />
            <div class="flex-1 flex flex-col gap-1">
              <div class="text-xl font-semibold text-primary">{{ currentAgent.name }}</div>
              <div class="text-sm text-regular break-words whitespace-pre-wrap">
                {{ currentAgent.description }}
              </div>
            </div>
          </div>
          <div class="mt-2 mb-10">
            <AuthTagGroup label-position="top" :model-value="currentAgent.user_group_ids" />
          </div>
          <x-bubble-assistant
            v-if="showWelcome"
            type="welcome"
            :content="currentAgent.settings_obj.opening_statement"
            :suggestions="currentAgent.settings_obj.suggested_questions"
            @suggestion="handleSuggestion"
          ></x-bubble-assistant>
        </template>
        <template #item="{ message, index }">
          <div
            class="flex items-center gap-5 rounded-xl"
            :class="[isShareMode ? 'mb-4 px-3 py-4 bg-[#F5F5F5]' : '']"
            @click="handleSelectMessage(message)"
          >
            <el-checkbox v-if="isShareMode" :model-value="state.selectMessageIds.includes(message.id)" />
            <div class="flex-1">
              <!-- 用户消息气泡 -->
              <x-bubble-user :key="message.id + '_user'" :content="message.query" :class="[isShareMode ? '!mb-0' : '']" :files="message.user_files">
                <template v-if="state.displayMode === DISPLAY_MODE.CHAT" #menu>
                  <div
                    v-tooltip="{ content: $t('action.copy') }"
                    v-copy="message.query"
                    class="h-6 px-1 rounded flex-center cursor-pointer hover:bg-[#E1E2E3]"
                  >
                    <el-icon color="#9B9B9B">
                      <CopyDocument />
                    </el-icon>
                  </div>
                </template>
              </x-bubble-user>
            </div>
          </div>

          <!-- AI助手消息气泡 -->
          <div
            class="flex items-center gap-5 rounded-xl"
            :class="[isShareMode ? 'mb-4 px-3 py-4 bg-[#F5F5F5]' : '']"
            @click="handleSelectMessage(message)"
          >
            <el-checkbox v-if="isShareMode" :model-value="state.selectMessageIds.includes(message.id)" />
            <div class="flex-1">
              <x-bubble-assistant
                :key="message.id + '_assistant'"
                :content="message.answer"
                :class="[isShareMode ? '!mb-0' : '']"
                :reasoning="message.reasoning_content"
                :reasoning-expanded="message.reasoning_expanded"
                :streaming="message.loading"
                :always-show-menu="index === state.messageList.length - 1"
              >
                <template v-if="!message.loading && state.displayMode === DISPLAY_MODE.CHAT" #menu>
                  <div
                    v-tooltip="{ content: $t('action.copy') }"
                    v-copy="message.answer"
                    class="h-6 px-1 rounded flex-center cursor-pointer hover:bg-[#E1E2E3]"
                  >
                    <el-icon color="#9B9B9B">
                      <CopyDocument />
                    </el-icon>
                  </div>
                  <div
                    v-tooltip="{ content: $t('chat.regenerate') }"
                    class="h-6 px-1 rounded flex-center cursor-pointer hover:bg-[#E1E2E3]"
                    @click="handleRegenerate(message)"
                  >
                    <el-icon color="#9B9B9B">
                      <Refresh />
                    </el-icon>
                  </div>
                  <div class="h-6 px-1 rounded flex-center gap-1 cursor-pointer hover:bg-[#E1E2E3]" @click="handleOpenShare(message)">
                    <svg-icon size="18" name="share-two" color="#9B9B9B" stroke></svg-icon>
                    <span class="text-sm text-[#939499]">{{ $t('action.share') }}</span>
                  </div>
                  <div v-if="false" v-tooltip="{ content: $t('chat.like') }" class="h-6 px-1 rounded flex-center cursor-pointer hover:bg-[#E1E2E3]">
                    <svg-icon size="18" name="like" color="#9B9B9B"></svg-icon>
                  </div>
                  <div v-if="false" v-tooltip="{ content: $t('chat.like') }" class="h-6 px-1 rounded flex-center cursor-pointer hover:bg-[#E1E2E3]">
                    <svg-icon size="18" name="dislike" color="#9B9B9B"></svg-icon>
                  </div>
                </template>
              </x-bubble-assistant>
            </div>
          </div>
          <RelatedScene v-if="index === state.messageList.length - 1 && !message.loading && !isShareMode" :output="message.answer" />
        </template>
      </x-bubble-list>

      <div v-if="showRecommend" class="flex-none w-2/6 flex flex-col gap-4 pb-5" :class="[isShareMode ? '-mt-[70px]' : '']">
        <h2 class="flex-none text-base font-semibold text-regular">{{ $t('common.related_agent') }}</h2>
        <div v-if="currentAgent.agent_id" class="flex-1 overflow-y-auto flex flex-col gap-2.5">
          <template v-for="item in relatedAgentList" :key="item.agent_id">
            <div class="flex-none h-24 border rounded p-4 cursor-pointer hover:bg-[#F1F2F3]" @click="onSelectAgent(item)">
              <div class="flex items-center gap-2">
                <img class="size-6 rounded-full" :src="item.logo" />
                <span class="text-sm text-primary">{{ item.name }}</span>
              </div>
              <div class="text-sm text-regular line-clamp-2 mt-1.5" :title="item.description">
                {{ item.description || '--' }}
              </div>
            </div>
          </template>
        </div>
      </div>
    </div>

    <!-- 底部输入区域 -->
    <div v-if="!isShareMode" class="py-5 sticky bottom-0 bg-white" :class="[showRecommend ? 'w-4/6' : 'w-11/12 md:w-4/5 max-w-[800px] mx-auto']">
      <div class="flex gap-2 mb-2.5">
        <AgentTooltip @select="onSelectAgent">
          <div class="h-8 px-2 rounded-full flex-center gap-1.5 bg-[#F1F2F3] cursor-pointer hover:bg-[#E1E2E3]">
            <img class="size-4 rounded-full" :src="currentAgent.logo" alt="" />
            <span class="text-sm text-primary">{{ currentAgent.name }}</span>
            <div class="size-4 flex-center">
              <el-icon color="#333333">
                <ArrowDown />
              </el-icon>
            </div>
          </div>
        </AgentTooltip>

        <div v-if="false" class="h-8 px-2 rounded-full flex-center gap-1.5 bg-[#F1F2F3] cursor-pointer hover:bg-[#E1E2E3]">
          <div class="size-4">
            <svg-icon name="network"></svg-icon>
          </div>
          <span class="text-sm text-primary">{{ $t('chat.online_search') }}</span>
        </div>
        <div class="flex-1"></div>
        <div
          v-if="showHistory"
          v-permission
          class="h-8 px-2 rounded-full flex-center gap-1.5 bg-[#F1F2F3] text-sm text-primary cursor-pointer hover:bg-[#E1E2E3]"
          @click="handleHistory"
        >
          <div class="size-4">
            <svg-icon name="history"></svg-icon>
          </div>
          {{ $t('chat.history') }}
        </div>
        <div
          class="h-8 px-2 rounded-full flex-center gap-1.5 bg-[#F1F2F3] text-sm text-primary cursor-pointer hover:bg-[#E1E2E3]"
          @click="handleNewConversation"
        >
          <div class="size-4">
            <svg-icon name="plus"></svg-icon>
          </div>
          {{ $t('chat.new_conversation') }}
        </div>
      </div>
      <x-sender
        ref="senderRef"
        :http-request="httpRequest"
        :enable-upload="enable_upload"
        :accept-types="upload_accept"
        :loading="state.isStreaming"
        allow-multiple
        enable-drag-upload
        :allow-send-with-files="allowSendWithFiles"
        @send="handleSend"
        @stop="handleStop"
      ></x-sender>
      <div v-if="enterpriseStore.copyright.toLowerCase() !== 'true'" class="flex justify-center items-center my-2">
        <img :src="$getPublicPath('/images/chat/footer.png')" class="h-3" />
      </div>
    </div>
  </div>

  <!-- 右侧帮助面板 -->
  <Transition name="slide">
    <div
      v-if="state.showHelper"
      class="border-l bg-white left-0 right-0 top-0 bottom-0 z-10"
      :class="[useCaseFixed && state.showHelper ? 'fixed' : 'absolute']"
    >
      <div class="h-[70px] flex-center border-b relative">
        <h4 class="text-lg text-primary">{{ $t('chat.usage_guide') }}</h4>
        <div
          class="flex-center size-6 absolute right-2 top-1/2 -translate-y-1/2 rounded cursor-pointer hover:bg-[#ECEDEE]"
          @click="handleToggleGuide"
        >
          <el-icon>
            <Close />
          </el-icon>
        </div>
      </div>
      <Helper :agent="currentAgent"></Helper>
    </div>
  </Transition>
  <HistoryDrawer ref="historyRef" @new="handleNewConversation" />
</template>

<script setup lang="ts">
import { ArrowDown, Close, CopyDocument, Refresh } from '@element-plus/icons-vue'
import { computed, defineAsyncComponent, onMounted, onBeforeUnmount, reactive, ref } from 'vue'
import MainHeader from '@/layout/header.vue'
import AuthTagGroup from '@/components/AuthTagGroup/index.vue'
import AgentTooltip from './components/agent-tooltip.vue'
import HistoryDrawer from './components/history.vue'
import RelatedScene from '@/components/RelatedScene/index.vue'

import { useAgentStore } from '@/stores/modules/agent'
import { useConversationStore } from '@/stores/modules/conversation'
import { useEnterpriseStore } from '@/stores/modules/enterprise'

import { API_HOST } from '@/api/host'
import chatApi from '@/api/modules/chat'
import conversationApi from '@/api/modules/conversation'
import uploadApi from '@/api/modules/upload'

import { checkPermission } from '@/utils/permission'
import sharesApi from '@/api/modules/share'
import { copyToClip } from '@/utils/copy'
import { isHashRouter } from '@/router'

const Helper = defineAsyncComponent(() => import('../helper.vue'))
// 扩展消息类型，添加动效相关属性
interface ExtendedMessage extends Conversation.Message {
  isNew?: boolean
}

const DISPLAY_MODE = {
  CHAT: 'chat',
  SHARE: 'share'
}

const agentStore = useAgentStore()
const convStore = useConversationStore()
const enterpriseStore = useEnterpriseStore()

withDefaults(
  defineProps<{
    hideMenuHeader?: boolean
    showRecommend?: boolean
    useCaseFixed?: boolean
    showHistory?: boolean
  }>(),
  {
    hideMenuHeader: false,
    showRecommend: false,
    useCaseFixed: false,
    showHistory: false
  }
)

const abortController = ref<AbortController | null>(null)
const historyRef = ref<InstanceType<typeof HistoryDrawer> | null>(null)

const senderRef = ref(null)
const bubbleListRef = ref(null)
const containerHeight = ref(0)
let resizeObserver: ResizeObserver | null = null

const state = reactive<{
  offset: number
  limit: number
  showHelper: boolean
  messageList: ExtendedMessage[]
  isStreaming: boolean
  isLoadingMore: boolean
  hasMore: boolean
  displayMode: string
  selectMessageIds: number[]
  selectAll: boolean
}>({
  offset: 0,
  limit: 10,
  showHelper: false,
  messageList: [],
  isStreaming: false,
  isLoadingMore: false, // 添加加载更多状态标志
  hasMore: true, // 是否还有更多消息
  displayMode: DISPLAY_MODE.CHAT,
  selectMessageIds: [],
  selectAll: false
})

const currentAgent = computed(() => convStore.currentAgent)
const currentConv = computed(() => convStore.currentConversation)
const custom_config_obj = computed(() => convStore.currentAgent.custom_config_obj)
const enable_upload = computed(() => Boolean(custom_config_obj.value?.file_parse?.enable || custom_config_obj.value?.image_parse?.enable))
const upload_accept = computed(() => {
  let accept = ''
  if (custom_config_obj.value?.file_parse?.enable) accept += '.pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.csv,.txt,.html,.json,.xml,.md'
  if (custom_config_obj.value?.image_parse?.enable) accept += ',image/*'
  return accept
})
const allowSendWithFiles = computed(() => {
  return ['53ai_agent', 'fastgpt_agent'].includes(custom_config_obj.value?.agent_type)
})

const relatedAgentList = computed(() => {
  return agentStore.agentList.filter((item) => item.agent_id !== currentAgent.value.agent_id).slice(0, 4)
})

const showWelcome = computed(() => {
  const settings = currentAgent.value.settings_obj
  if (settings.opening_statement && settings.opening_statement.replace(/\s/g, '')) return true
  if (
    settings.suggested_questions &&
    settings.suggested_questions.length &&
    settings.suggested_questions.some((item) => item.content.replace(/\s/g, ''))
  )
    return true
  return false
})
const handleHeightChange = (height) => {
  agentStore.setBoxHeight(height)
}

const isShareMode = computed(() => {
  return state.displayMode === DISPLAY_MODE.SHARE
})

const handleToggleGuide = () => {
  state.showHelper = !state.showHelper
}

const onSelectAgent = (agent: Agent.State) => {
  convStore.pushUsualAgent(agent)
  convStore.setCurrentState(agent.agent_id, 0)
}

const handleNewConversation = () => {
  convStore.setCurrentState(currentAgent.value.agent_id, 0)
}

const httpRequest = (dataFile: File) => {
  return new Promise((resolve, reject) => {
    const isPermission = checkPermission({
      groupIds: currentAgent.value.user_group_ids,
      onClick: async () => {
        try {
          const res = await uploadApi.upload(dataFile)
          resolve({
            id: res.data.id,
            url: `${API_HOST}/api/preview/${res.data.preview_key || ''}`,
            size: res.data.size,
            name: res.data.file_name,
            mime_type: res.data.mime_type
          })
        } catch (error) {
          reject(error)
        }
      }
    })
    if (!isPermission) {
      reject(new Error(window.$t('authority.login_not_permission')))
    }
  })
}

// 处理流式数据的函数
const processStreamData = (e: any, processedLength: number): number => {
  const lastMessage = state.messageList[state.messageList.length - 1]
  if (!e.event?.target) return processedLength

  const fullResponse = e.event.target.response || ''
  // 只处理新增的部分
  const newChunk = fullResponse.substring(processedLength)
  const newProcessedLength = fullResponse.length

  try {
    // 处理SSE格式的数据
    const lines = newChunk.split('\n').filter((line) => line.trim() !== '' && line.trim() !== 'data: [DONE]')
    console.log(lines)
    for (const line of lines) {
      if (line.startsWith('data: ')) {
        try {
          const data = JSON.parse(line.slice(6))
          const { message_id } = data
          const content = data.choices?.[0]?.delta?.content
          const reasoning_content = data.choices?.[0]?.delta?.reasoning_content
          if (content) {
            // 将新内容追加到当前消息
            lastMessage.answer += content
          }
          if (reasoning_content) {
            // 将新内容追加到当前消息
            lastMessage.reasoning_content += reasoning_content
          }
          if (message_id) {
            lastMessage.id = message_id
          }
        } catch (err) {
          console.error('解析JSON失败:', err)
        }
      }
    }
  } catch (err) {
    console.error('处理流数据失败:', err)
  }

  return newProcessedLength
}

// 工具函数
const messageUtils = {
  // 格式化消息
  formatMessage: (item: any): ExtendedMessage => {
    const data = {
      ...item,
      query: ''
    }
    const { content } = JSON.parse(item.message)[0]
    try {
      const arr = JSON.parse(content)
      const query = arr.find((item) => item.type === 'text')?.content
      data.query = query
      data.user_files = arr.filter((item) => item.type === 'image')
    } catch (error) {
      data.query = content
    }
    return data
  },

  // 格式化文件
  formatFiles: (user_files: any[]): Conversation.UserFile[] =>
    user_files?.map((item) => ({
      type: 'image',
      content: `file_id:${item.id}`,
      filename: item.name,
      size: item.size,
      mime_type: item.mime_type,
      url: item.url
    })) || [],

  // 创建新消息
  createNewMessage: (query: string, agent_id: number, conversation_id: number, user_files: Conversation.UserFile[]): ExtendedMessage => ({
    id: Date.now(),
    query,
    answer: '',
    loading: true,
    agent_id,
    conversation_id,
    reasoning_content: '',
    reasoning_expanded: true,
    user_files
  })
}

// 加载消息列表
const loadMessages = async (conversation_id: number, offset: number, limit: number) => {
  try {
    const res = await conversationApi.messasges(conversation_id, { offset, limit })
    const list = res.data.messages.map(messageUtils.formatMessage)
    return {
      messages: list,
      hasMore: list.length === limit
    }
  } catch (err) {
    console.error('加载消息失败:', err)
    return { messages: [], hasMore: false }
  }
}

const handleLoadListMore = async (done: () => void): Promise<void> => {
  if (state.isLoadingMore || !state.hasMore) return done()

  const { conversation_id } = currentConv.value
  if (!conversation_id) return done()

  state.isLoadingMore = true
  state.offset += state.limit

  try {
    const { messages, hasMore } = await loadMessages(conversation_id, state.offset, state.limit)
    state.hasMore = hasMore
    state.messageList.unshift(...messages)
  } catch (err) {
    state.offset = Math.max(0, state.offset - state.limit)
  } finally {
    state.isLoadingMore = false
  }
  return done()
}

const loadList = async () => {
  const { conversation_id } = currentConv.value
  if (!conversation_id) return

  if (!enterpriseStore.template_style_info.style_type) {
    await enterpriseStore.initTemplateStyle()
  }

  state.isLoadingMore = true
  state.offset = 0
  state.hasMore = true

  try {
    const { messages, hasMore } = await loadMessages(conversation_id, state.offset, state.limit)
    state.hasMore = hasMore
    state.messageList = messages
  } finally {
    state.isLoadingMore = false
  }
}

const buildHistoryMessages = (maxPairs: number = 6) => {
  const history: Array<{ role: 'user' | 'assistant'; content: string }> = []
  // 仅使用已完成的历史（包含 answer 的）
  const finished = state.messageList.filter((m) => m.query && (m.answer !== undefined))
  const start = Math.max(0, finished.length - maxPairs)
  for (let i = start; i < finished.length; i++) {
    const m = finished[i]
    const hasFiles = Array.isArray(m.user_files) && m.user_files.length > 0
    const userContent = hasFiles
      ? JSON.stringify([{ type: 'text', content: m.query }, ...m.user_files])
      : m.query
    history.push({ role: 'user', content: userContent })
    if (m.answer) history.push({ role: 'assistant', content: m.answer })
  }
  return history
}

const getHistoryPairs = (): number => {
  try {
    // Only prompt type carries context
    if (custom_config_obj.value?.agent_type !== 'prompt') return 0
    const configs = JSON.parse(currentAgent.value.configs || '{}')
    const chat = configs.chat || {}
    const completion_params = configs.completion_params || {}
    const v = Number(chat.history_pairs ?? completion_params.history_pairs)
    if (Number.isFinite(v) && v > 0) return Math.floor(v)
  } catch (e) {}
  return 6
}


const sendMessage = async (query: string, user_files: any[]) => {
  if (state.isStreaming) return

  const { agent_id } = currentAgent.value
  const { conversation_id } = currentConv.value

  // 先构建历史对话，避免把当前占位消息带入请求
  const history = buildHistoryMessages(getHistoryPairs())

  // 创建新消息（用于界面展示与流式追加）
  const newMessage = messageUtils.createNewMessage(query, agent_id, conversation_id, user_files)
  state.messageList.push(newMessage)

  const configs = JSON.parse(currentAgent.value.configs || '{}')
  const completion_params = configs.completion_params || {}
  state.isStreaming = true
  abortController.value = new AbortController()
  let processedLength = 0

  let content = query
  if (user_files.length > 0) {
    content = JSON.stringify([{ type: 'text', content: query }, ...user_files])
  }

  try {
    await chatApi.completions(
      {
        conversation_id,
        model: `agent-${agent_id}`,
        messages: [...history, { content, role: 'user' }],
        frequency_penalty: 0,
        presence_penalty: 0,
        stream: true,
        temperature: 0,
        top_p: 0,
        ...completion_params
      },
      {
        responseType: 'stream',
        onDownloadProgress: (e) => {
          processedLength = processStreamData(e, processedLength)
        },
        signal: abortController.value.signal
      }
    )
  } catch (err: any) {
    if (err.message !== 'canceled') {
      console.log(err)
      const lastMessage = state.messageList[state.messageList.length - 1]
      if (lastMessage && !lastMessage.answer) {
        lastMessage.answer = err.response?.data || window.$t('response_code.network_error')
      }
    }
  } finally {
    const lastMessage = state.messageList[state.messageList.length - 1]
    if (lastMessage) {
      lastMessage.loading = false
    }

    state.isStreaming = false
    abortController.value = null
  }
}

const handleSend = (question: string, user_files: any[]) => {
  checkPermission({
    groupIds: currentAgent.value.user_group_ids,
    onClick: async () => {
      const { agent_id } = currentAgent.value
      if (!agent_id) return ElMessage.warning(window.$t('chat.no_available_agent'))
      if (!currentConv.value.conversation_id) {
        try {
          const conversation = await convStore.createConversation(agent_id, question)
          convStore.addConversation({ ...conversation, virtual_id: currentConv.value.virtual_id })
          convStore.setCurrentState(conversation.agent_id, conversation.conversation_id)
        } catch (err) {
          console.error('创建对话失败:', err)
        }
      }
      await sendMessage(question, messageUtils.formatFiles(user_files))
      return true
    }
  })
}

const handleHistory = () => {
  historyRef.value?.open()
}

const handleStop = () => {
  // 取消当前请求
  if (abortController.value) {
    abortController.value.abort()
    abortController.value = null
  }
  state.isStreaming = false
}

const handleRegenerate = async (message: Conversation.Message) => {
  await sendMessage(message.query, message.user_files)
}

const handleSuggestion = (suggestion: string) => {
  handleSend(suggestion, [])
}

const handleOpenShare = (message?: Conversation.Message) => {
  state.selectAll = false
  state.selectMessageIds = []
  if (message) {
    state.displayMode = DISPLAY_MODE.SHARE
  } else {
    state.displayMode = DISPLAY_MODE.SHARE === state.displayMode ? DISPLAY_MODE.CHAT : DISPLAY_MODE.SHARE
  }
}

const handleSelectAll = () => {
  if (state.displayMode === DISPLAY_MODE.SHARE) {
    state.selectMessageIds = state.selectAll ? [] : state.messageList.map((item) => item.id)
    state.selectAll = !state.selectAll
  }
}
const handleSelectMessage = (message: Conversation.Message) => {
  if (state.displayMode === DISPLAY_MODE.SHARE) {
    if (state.selectMessageIds.includes(message.id)) {
      state.selectMessageIds = state.selectMessageIds.filter((id) => id !== message.id)
      state.selectAll = false
    } else {
      state.selectMessageIds.push(message.id)
    }
  }
}

const handleCreateShare = () => {
  return sharesApi
    .create({
      message_ids: state.selectMessageIds,
      conversation_id: currentConv.value.conversation_id,
      select_all: state.selectAll
    })
    .then((res) => {
      const link = `${window.location.origin}${isHashRouter ? '#' : ''}/share/chat?share_id=${res.share_id}`
      copyToClip(link).then(() => {
        ElMessage.success(window.$t('chat.completion_share_link'))
      })
      state.displayMode = DISPLAY_MODE.CHAT
    })
}

onMounted(() => {
  loadList()
  if (bubbleListRef.value) {
    resizeObserver = new ResizeObserver((entries) => {
      for (const entry of entries) {
        const { height } = entry.contentRect
        containerHeight.value = height
        handleHeightChange(height)
      }
    })

    // 注意：有些组件可能需要访问 $el，有些则直接是 DOM 元素
    const targetElement = bubbleListRef.value.$el || bubbleListRef.value
    resizeObserver.observe(targetElement)
  }

  const prepare = convStore.next_agent_prepare
  if (prepare.agent_id) {
    const question = prepare.parameters.input
    senderRef.value?.setPrompt && senderRef.value?.setPrompt(question)
    if (prepare.execution_rule === 'auto') {
      senderRef.value?.setPrompt('')
      handleSend(question, [])
    }
    convStore.setNextAgentPrepare({})
  }
})

onBeforeUnmount(() => {
  if (resizeObserver) {
    resizeObserver.disconnect()
  }
})

defineExpose({
  showUseCase: () => {
    state.showHelper = true
  },
  hideUseCase: () => {
    state.showHelper = false
  },
  showShare: () => {
    handleOpenShare()
  }
})
</script>

<style scoped></style>
