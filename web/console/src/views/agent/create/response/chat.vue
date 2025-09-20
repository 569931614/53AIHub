<template>
  <div class="flex flex-col pt-7 relative">
    <div v-if="is_config_changed" class="absolute top-0 left-0 w-full h-full bg-black/70 z-10">
      <div class="flex flex-col items-center justify-center gap-6 w-full h-full box-border">
        <div class="text-base text-[#fff] text-center mx-8">
          {{ $t('debugger_config_change_confirm') }}
        </div>
        <ElButton v-debounce type="primary" size="large" @click="onRestart({ saveAction: true })">
          {{ $t('save_and_restart') }}
        </ElButton>
      </div>
    </div>
    <div class="flex items-center justify-between px-4 mb-2">
      <div class="text-base text-[#1D1E1F]">
        {{ $t('debug_preview') }}
      </div>
      <div class="flex-center gap-1 cursor-pointer" @click="onRestart">
        <ElIcon>
          <RefreshRight />
        </ElIcon>
        <span class="text-sm text-[#1D1E1F]">
          {{ $t('restart') }}
        </span>
      </div>
    </div>

    <x-bubble-list :messages="chat_list" class="flex-1 px-4 relative py-4" main-class="mx-5">
      <template #header>
        <ElEmpty v-if="showChatListEmpty" class="mt-10" :description="$t('chat.empty_desc')" />
        <x-bubble-assistant
          v-if="showWelcome"
          type="welcome"
          :content="agentFormStore.form_data.settings.opening_statement"
          :suggestions="agentFormStore.form_data.settings.suggested_questions"
          @suggestion="handleSuggestion"
        />
      </template>
      <template #item="{ message }">
        <x-bubble-user :content="message.question.content" :files="message.question.user_files">
          <template v-if="!message.answer.loading" #menu>
            <x-icon size="16" class="cursor-pointer" name="copy" @click="onCopy(message.question.content)" />
          </template>
        </x-bubble-user>
        <x-bubble-assistant
          :content="message.answer.content"
          :reasoning="message.answer.reasoning_content"
          :reasoning-expanded="message.answer.reasoning_expanded"
          :streaming="message.answer.loading"
          :always-show-menu="message_index === chat_list.length - 1"
        >
          <template v-if="!message.answer.loading" #menu>
            <x-icon size="16" class="cursor-pointer" name="copy" @click="onCopy(message.answer.content)" />
            <x-icon size="16" class="cursor-pointer" name="refresh" @click="onRestartGeneration(message)" />
          </template>
        </x-bubble-assistant>
      </template>

      <!-- <div class="flex flex-col space-y-4">
        <ElEmpty v-if="!chat_list.length" class="mt-10" :description="$t('chat.empty_desc')" />
        <template v-else>
          <template v-for="(item, item_index) in chat_list" :key="item_index">

          </template>
        </template>
      </div> -->
    </x-bubble-list>
    <div class="px-6 py-3">
      <x-sender
        :enable-upload="enable_upload"
        :accept-types="upload_accept"
        :http-request="httpRequest"
        :loading="chat_loading"
        allow-multiple
        enable-drag-upload
        :allow-send-with-files="allowSendWithFiles"
        @send="onSendConfirm"
        @stop="onStopGeneration"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { RefreshRight } from '@element-plus/icons-vue'
import { computed, nextTick, ref, watch } from 'vue'
import { useAgentFormStore } from '../store'

import { useConversationStore } from '@/stores'
import { copyToClip } from '@/utils/copy'
import { api_host } from '@/utils/config'
import { AGENT_TYPES } from '@/constants/platform/config'
import uploadApi from '@/api/modules/upload'

const agentFormStore = useAgentFormStore()
const conversationStore = useConversationStore()
const scroll_ref = ref()
const chat_list = ref([])
const conversationCreating = ref(false)

const chat_loading = computed(() => conversationCreating.value || chat_list.value.some(item => item.answer.loading))
const enable_upload = computed(() =>
  Boolean(
    agentFormStore.form_data.settings?.file_parse?.enable || agentFormStore.form_data.settings?.image_parse?.enable
  )
)
const upload_accept = computed(() => {
  let accept = ''
  if (agentFormStore.form_data.settings?.file_parse?.enable)
    accept += '.pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.csv,.txt,.html,.json,.xml,.md'
  if (agentFormStore.form_data.settings?.image_parse?.enable) accept += ',image/*'
  return accept
})
const allowSendWithFiles = computed(() => {
  return [AGENT_TYPES['53AI_AGENT'], AGENT_TYPES.FASTGPT_AGENT].includes(agentFormStore.agent_type)
})

const showWelcome = computed(() => {
  const { settings } = agentFormStore.form_data
  if (settings.opening_statement.replace(/\s/g, '')) return true
  if (settings.suggested_questions.length && settings.suggested_questions.some(item => item.content.replace(/\s/g, '')))
    return true
  return false
})
const showChatListEmpty = computed(() => {
  if (chat_list.value.length) return false
  if (showWelcome.value) return false
  return true
})

const httpRequest = async (dataFile: File) => {
  try {
    const res = await uploadApi.upload(dataFile)
    return {
      id: res.data.id,
      url: `${api_host}/api/preview/${res.data.preview_key || ''}`,
      size: res.data.size,
      name: res.data.file_name,
      mime_type: res.data.mime_type,
    }
  } catch (error) {
    return {}
  }
}

let conversation_id = 0
let active_chat_index = -1
let active_chat_data = {}
let abort_controller: any = null
const onSendConfirm = async (question: string, user_files?: any[], type = '') => {
  if (chat_loading.value) return
  user_files = user_files || []
  if (!agentFormStore.agent_data.agent_id) return ElMessage.warning(window.$t('agent_not_found'))
  if (!agentFormStore.agent_data.channel_type) await agentFormStore.saveAgentData({ hideToast: true })
  // return ElMessage.warning(window.$t('agent_channel_type_not_found'))

  if (abort_controller) abort_controller.abort()
  abort_controller = new AbortController()

  if (!conversation_id) {
    conversationCreating.value = true
    const { data = {} } = await conversationStore
      .save({ data: { agent_id: agentFormStore.agent_data.agent_id, title: question } })
      .finally(() => {
        conversationCreating.value = false
      })
    conversation_id = data.conversation_id
  }

  if (type !== 'regenerate')
    user_files =
      user_files?.map(item => ({
        type: 'image',
        content: `file_id:${item.id}`,
        filename: item.name,
        size: item.size,
        mime_type: item.mime_type,
        url: item.url,
      })) || []

  chat_list.value.push({
    question: {
      role: 'user',
      content: question,
      user_files,
    },
    answer: {
      loading: true,
      role: 'assistant',
      content: '',
      reasoning_expanded: true,
      reasoning_content: '',
    },
  })
  active_chat_index = chat_list.value.length - 1
  active_chat_data = chat_list.value[active_chat_index] || {}
  // Build multi-turn message history
  const historyMessages: any[] = []
  const toUserContent = (q: string, files: any[]) => {
    if (files && files.length) {
      return JSON.stringify([{ type: 'text', content: q }, ...(files || [])])
    }
    return q
  }
  // Include prior rounds (before the current one)
  const prior = chat_list.value.slice(0, active_chat_index)
  for (const item of prior) {
    historyMessages.push({
      role: 'user',
      content: toUserContent(item.question.content, item.question.user_files || []),
    })
    if (item.answer && item.answer.content) {
      historyMessages.push({
        role: item.answer.role || 'assistant',
        content: item.answer.content,
      })
    }
  }
  // Current user message
  const currentUserMessage = {
    role: 'user',
    content: toUserContent(question, user_files),
  }
  const messages = [...historyMessages, currentUserMessage]

  conversationStore
    .chat({
      data: {
        conversation_id,
        messages,
        agent_id: agentFormStore.agent_data.agent_id,
        agent_configs: agentFormStore.agent_data.configs,
      },
      hideError: true,
      onDownloadProgress: async ({ chunks = [], intact_content, intact_reasoning_content } = {}) => {
        active_chat_data.answer.content = intact_content || active_chat_data.answer.content || ''
        active_chat_data.answer.reasoning_content =
          intact_reasoning_content || active_chat_data.answer.reasoning_content || ''
        if (chunks[0] && chunks[0].role) active_chat_data.answer.role = chunks[0].role || ''
        await nextTick()
        if (scroll_ref.value) scroll_ref.value.scrollToBottom()
      },
      signal: abort_controller.signal,
    })
    .catch(err => {
      if (!active_chat_data.answer.content) active_chat_data.answer.content = err.message

      ElMessage.warning(
        err.message === 'Access token is invalid' ? window.$t('agent_app.check_agent_config_tip') : err.message
      )
    })
    .finally(() => {
      active_chat_data.answer.loading = false
      abort_controller = null
    })
  await nextTick()
  if (scroll_ref.value) scroll_ref.value.scrollToBottom()
}
const onStopGeneration = () => {
  if (abort_controller) {
    abort_controller.abort()
    abort_controller = null
    active_chat_data.answer.loading = false
  }
}
const onRestartGeneration = data => {
  // chat_list.value.splice(0, active_chat_index + 1)
  onSendConfirm(data.question.content, data.question.user_files, 'regenerate')
}
const onRestart = ({ saveAction = false } = {}) => {
  // 智能体 、模型选择暂时 搬到 设置里面了，所以这段暂时不用
  // if (saveAction)
  //   return emits('save', { restart: true })
  conversation_id = 0
  chat_list.value = []
  is_config_changed.value = false
}
const onCopy = async (text = '') => {
  await copyToClip(text)
  ElMessage.success(window.$t('action_copy_success'))
}
const handleSuggestion = (question: string) => {
  onSendConfirm(question)
}

const is_config_changed = ref(false)
watch(
  () => agentFormStore.form_data.custom_config,
  data => {
    is_config_changed.value = false
    if (conversation_id) is_config_changed.value = true
  },
  {
    deep: true,
  }
)

defineExpose({
  restart: onRestart,
  getIsConfigChanged: () => is_config_changed.value,
})
</script>

<style scoped></style>
