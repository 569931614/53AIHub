<template>
  <div v-loading="state.isLoading" class="h-full flex flex-col bg-white relative overflow-hidden">
    <div class="flex-none h-[110px] flex items-center justify-center relative px-20 max-md:h-20 max-md:px-10">
      <div class="absolute left-4 md:left-11 flex items-center gap-2">
        <img :src="enterpriseStore.logo" :title="enterpriseStore.display_name" class="h-[34px]" />
        <h1 class="text-base font-bold max-md:hidden">{{ enterpriseStore.display_name }}</h1>
      </div>
      <h2 v-if="state.agent" class="text-lg md:text-xl text-[#1D1E1F] text-center px-4 md:px-0 max-w-[calc(100%-120px)] md:max-w-none truncate">
        {{ state.user?.nickname }}与{{ state.agent?.name || '--' }}的对话
      </h2>
    </div>

    <template v-if="state.agent">
      <!-- 消息列表区域 -->
      <div class="flex-1 overflow-hidden">
        <x-bubble-list
          ref="bubbleListRef"
          :auto-scroll="false"
          class="flex-1"
          :messages="state.messageList"
          :main-class="'w-11/12 md:w-4/5 max-w-[800px] mx-auto mt-5'"
          enable-pull-up
        >
          <template #item="{ message, index }">
            <!-- 用户消息气泡 -->
            <x-bubble-user :key="message.id + '_user'" :content="message.query" :files="message.user_files"></x-bubble-user>

            <!-- AI助手消息气泡 -->
            <x-bubble-assistant
              :key="message.id + '_assistant'"
              :content="message.answer"
              :reasoning="message.reasoning_content"
              :reasoning-expanded="message.reasoning_expanded"
              :streaming="message.loading"
              :always-show-menu="index === state.messageList.length - 1"
            ></x-bubble-assistant>

            <div v-if="state.messageList.length - 1 === index" class="text-xs text-[#939499]">{{ $t('common.ai_generated') }}</div>
          </template>
        </x-bubble-list>
      </div>

      <div class="flex-none h-[110px] flex items-center justify-center relative">
        <div class="h-12 flex items-center gap-2 bg-[#2563EB] rounded-full px-8 cursor-pointer hover:bg-[#1D5ECD]" @click="handleOpenAgent">
          <img :src="state.agent?.logo" :title="state.agent?.name" class="h-5 rounded-full" />
          <span class="text-sm text-white truncate max-w-80">跟“{{ state.agent?.name || '--' }}”聊一聊</span>
        </div>
      </div>
    </template>
    <template v-else-if="!state.isLoading">
      <div class="flex-center flex-col gap-2">
        <el-empty :image="$getPublicPath('/images/chat/completion_empty.png')" :description="$t('chat.no_available_agent_desc')" />
        <el-button class="rounded-full" size="large" type="primary" @click="handleBackHome">{{ $t('common.back_home') }}</el-button>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'

import { useEnterpriseStore } from '@/stores/modules/enterprise'

import sharesApi from '@/api/modules/share'
import type { ShareFindReponse } from '@/api/modules/share/types'

import { router } from '@/router'

const route = useRoute()
const enterpriseStore = useEnterpriseStore()

const bubbleListRef = ref(null)

const state = reactive<{
  user: ShareFindReponse['user'] | null
  agent: ShareFindReponse['agent'] | null
  messageList: ShareFindReponse['messages']
  isLoading: boolean
}>({
  user: null,
  agent: null,
  messageList: [],
  isLoading: false
})

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
  }
}

const handleOpenAgent = () => {
  router.replace(`${enterpriseStore.isSoftStyle ? '' : '/index'}/chat?agent_id=${state.agent?.agent_id}`)
}

const handleBackHome = () => {
  router.replace('/')
}

onMounted(async () => {
  state.isLoading = true
  const shareId = route.query.share_id as string
  try {
    const res = await sharesApi.find(shareId)
    state.messageList = res.messages.map(messageUtils.formatMessage) as ShareFindReponse['messages']
    state.user = res.user
    state.agent = res.agent
  } finally {
    state.isLoading = false
  }
})

onBeforeUnmount(() => {})
</script>

<style scoped></style>
