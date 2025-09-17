<template>
  <div
    v-if="!globalStore.siderVisible && !globalStore.siderCollapsed && !globalStore.isMdScreen"
    class="w-4 h-full absolute left-0 top-0 z-[3] hover:bg-gray-100/50 transition-colors"
    @mouseenter="globalStore.hoverSider(true)"
  />
  <div v-if="globalStore.isMdScreen && globalStore.siderCollapsed" class="absolute top-0 left-0 w-full h-full z-[10] bg-black/60"></div>
  <div
    ref="siderRef"
    class="w-[280px] border-r flex flex-col bg-[#fff] transition-all duration-300 ease-linear md:bg-[#FAFAFA]"
    :class="[
      {
        'absolute z-[11] left-0 -translate-x-full shadow-xl rounded-r': !globalStore.siderVisible,
        'translate-x-0': !globalStore.siderVisible && globalStore.siderCollapsed,
        'top-0 bottom-0': globalStore.isMdScreen
      },
      globalStore.isMdScreen ? 'top-0 bottom-0' : 'top-20 bottom-5'
    ]"
    @mouseleave="globalStore.hoverSider(false)"
  >
    <div class="px-3 h-[70px] flex items-center gap-2">
      <router-link :to="{ name: 'Home' }">
        <img class="size-8 rounded" :src="enterpriseStore.logo" alt="企业Logo" />
      </router-link>
      <router-link :to="{ name: 'Home' }" class="flex-1 text-base text-primary truncate">{{ enterpriseStore.display_name }}</router-link>

      <div
        v-tooltip="{
          content: globalStore.siderVisible ? $t('chat.collapse_side_bar') : $t('chat.expand_side_bar')
        }"
        class="flex-none size-7 rounded-md flex-center cursor-pointer hover:bg-[#ECEDEE]"
        @click="globalStore.toggleSider"
      >
        <svg-icon name="layout-left" size="20" color="#9A9A9A"></svg-icon>
      </div>
    </div>
    <div class="border-b mx-2"></div>
    <div class="flex-1 py-3 overflow-y-auto">
      <div class="px-5 pb-2">
        <el-input
          v-model="state.keyword"
          class="el-input--main"
          :placeholder="$t('action.search')"
          :prefix-icon="Search"
          style="--el-input-inner-height: 36px; --el-input-border-color: none"
          clearable
        ></el-input>
      </div>
      <div v-if="state.keyword" class="px-5">
        <template v-if="searchUsualAgents.length">
          <div class="h-9 flex items-center gap-2">
            <div class="text-sm text-[#939499]">
              {{ $t('module.chat') }}
            </div>
            <div class="flex-1 h-px bg-[#E6E8EB]"></div>
          </div>

          <div class="flex flex-col gap-1">
            <template v-for="item in searchUsualAgents" :key="item.agent_id">
              <div
                class="h-9 px-6 flex items-center gap-2 rounded cursor-pointer text-[#4F5052] hover:bg-[#ECEDEE]"
                @click="convStore.setCurrentState(item.agent_id, 0)"
              >
                <el-image class="flex-none size-6 rounded-full" :src="item.logo"></el-image>
                <div class="flex-1 text-sm truncate" v-html="item.formt_name"></div>
              </div>
            </template>
          </div>
        </template>
        <template v-if="searchAgents.length">
          <div class="h-9 flex items-center gap-2">
            <div class="text-sm text-[#939499] max-w-[100px] truncate">
              {{ $t('module.agent') }}
            </div>
            <div class="flex-1 h-px bg-[#E6E8EB]"></div>
          </div>

          <div class="flex flex-col gap-1">
            <template v-for="item in searchAgents" :key="item.agent_id">
              <div
                class="h-9 px-6 flex items-center gap-2 rounded cursor-pointer text-[#4F5052] hover:bg-[#ECEDEE]"
                @click="
                  () => {
                    convStore.pushUsualAgent(item), convStore.setCurrentState(item.agent_id, 0)
                  }
                "
              >
                <el-image class="flex-none size-6 rounded-full" :src="item.logo"></el-image>
                <div class="flex-1 text-sm truncate" v-html="item.formt_name"></div>
              </div>
            </template>
          </div>
        </template>
        <template v-if="searchToolBox.length">
          <div class="h-9 flex items-center gap-2">
            <div class="text-sm text-[#939499] max-w-[100px] truncate">
              {{ navigationStore.toolkitNavigation.name || $t('module.toolbox') }}
            </div>
            <div class="flex-1 h-px bg-[#E6E8EB]"></div>
          </div>

          <div class="flex flex-col gap-1">
            <template v-for="item in searchToolBox" :key="item.agent_id">
              <a :href="item.url" target="_blank" class="h-9 px-6 flex items-center gap-2 rounded cursor-pointer text-[#4F5052] hover:bg-[#ECEDEE]">
                <el-image class="flex-none size-6 rounded-full" :src="item.logo"></el-image>
                <div class="flex-1 text-sm truncate" v-html="item.name"></div>
              </a>
            </template>
          </div>
        </template>
        <template v-if="searchHistory.length">
          <div class="h-9 flex items-center gap-2">
            <div class="text-sm text-[#939499]">
              {{ $t('chat.history') }}
            </div>
            <div class="flex-1 h-px bg-[#E6E8EB]"></div>
          </div>

          <div class="flex flex-col gap-1">
            <template v-for="item in searchHistory" :key="item.agent_id">
              <div
                class="h-9 px-3 rounded-md flex items-center gap-2 mt-0.5 group cursor-pointer hover:bg-[#ECEDEE] overflow-hidden"
                @click="convStore.setCurrentState(item.agent_id, item.conversation_id)"
              >
                <span class="flex-1 text-sm text-[#1D1E1F] truncate">{{ item.title }}</span>
              </div>
            </template>
          </div>
        </template>
      </div>

      <div v-show="!state.keyword">
        <div v-if="userStore.is_login" class="px-2">
          <div
            class="h-9 px-2 rounded-md flex items-center gap-0.5 cursor-pointer text-primary hover:bg-[#ECEDEE]"
            @click="state.showUsualAgents = !state.showUsualAgents"
          >
            <div class="size-7 flex-center">
              <svg-icon size="18" name="chat"></svg-icon>
            </div>
            <p class="flex-1 text-base">{{ $t('module.chat') }}</p>

            <el-icon>
              <ArrowUp v-if="state.showUsualAgents" />
              <ArrowDown v-else />
            </el-icon>
          </div>
        </div>
        <div v-show="state.showUsualAgents && userStore.is_login" class="px-3">
          <transition-group name="list" tag="div" class="flex flex-col gap-1 mt-2">
            <template v-for="item in usualAgents" :key="item.agent_id">
              <div
                class="h-9 pl-6 pr-2 rounded-md flex items-center gap-2 cursor-pointer hover:bg-[#ECEDEE] group text-[#4F5052]"
                :class="[currentAgent.agent_id === item.agent_id && !currentConv.conversation_id ? 'bg-[#ECEDEE] ' : '']"
                @click="convStore.setCurrentState(item.agent_id, 0)"
              >
                <el-image class="flex-none size-6 rounded" :src="item.logo"></el-image>
                <div class="flex-1 text-sm truncate">
                  {{ item.name }}
                </div>

                <div
                  v-show="!state.sidebarCollapsed"
                  v-tooltip="{ content: item.is_fixed ? $t('action.unfixed') : $t('action.fixed') }"
                  class="size-7 flex-center"
                  @click.stop="convStore.toggleUsualAgentFixed(item)"
                >
                  <svg-icon v-if="item.is_fixed" name="top" color="#2563EB"></svg-icon>
                  <svg-icon v-else name="top" color="#4F5052" class="invisible group-hover:visible"></svg-icon>
                </div>
              </div>
            </template>
          </transition-group>
          <div
            v-if="convStore.usual_agents.length > SHOW_USUAL_AGENT_LEN"
            class="flex items-center gap-1 mt-3.5 ml-5 cursor-pointer text-[#9A9A9A]"
            @click="toggleAgent"
          >
            <span v-show="!state.sidebarCollapsed" class="text-sm">{{ state.agentCollapsed ? $t('action.collapse') : $t('action.expand') }}</span>
            <el-icon>
              <ArrowUp v-if="state.agentCollapsed" />
              <ArrowDown v-else />
            </el-icon>
          </div>
        </div>
        <div class="flex flex-col px-2 max-md:hidden">
          <template v-for="(item, index) in navigationStore.navigations" :key="index">
            <router-link
              v-if="item.type !== NAVIGATION_TYPE.EXTERNAL"
              :to="{ path: item.jump_path }"
              class="h-9 px-2 rounded-md flex items-center gap-0.5 mt-1.5 cursor-pointer text-[#4F5052] hover:bg-[#ECEDEE]"
              :class="[route.path === item.jump_path ? 'bg-[#ECEDEE]' : '']"
            >
              <div class="size-7 flex-center">
                <svg-icon size="18" :name="item.icon || 'app'" stroke></svg-icon>
              </div>
              <p class="flex-1 text-base text-[#000000] truncate">
                {{ item.name }}
              </p>
            </router-link>

            <a
              v-else
              :href="item.jump_path"
              :target="item.target === NAVIGATION_TARGET.BLANK ? '_blank' : '_self'"
              rel="noopener noreferrer"
              class="h-9 px-2 rounded-md flex items-center gap-0.5 mt-1.5 cursor-pointer text-[#4F5052] hover:bg-[#ECEDEE]"
            >
              <div class="size-7 flex-center">
                <svg-icon size="18" :name="item.icon || 'app'" stroke></svg-icon>
              </div>
              <p class="flex-1 text-base text-[#000000] truncate">
                {{ item.name }}
              </p>
            </a>
          </template>
          <div
            v-if="$isElectron"
            href="http://ziroom.hub.53ai.com/space"
            target="_blank"
            class="h-9 px-2 rounded-md flex items-center gap-0.5 cursor-pointer text-[#4F5052] hover:bg-[#ECEDEE]"
            @click="handleJumpToLibrary"
          >
            <div class="size-7 flex-center">
              <svg-icon size="18" name="tips"></svg-icon>
            </div>
            <p class="flex-1 text-base text-[#000000]">{{ $t('module.library') }}</p>
          </div>
        </div>
        <template v-if="userStore.is_login">
          <div class="mx-2 border-t mt-1.5"></div>
          <!-- 历史对话 -->
          <div class="px-2 mt-1.5">
            <div
              class="h-9 px-2 rounded-md flex items-center gap-0.5 cursor-pointer text-primary hover:bg-[#ECEDEE]"
              @click="state.showHistory = !state.showHistory"
            >
              <div class="size-7 flex-center">
                <svg-icon size="18" name="history"></svg-icon>
              </div>
              <p class="flex-1 text-base">{{ $t('chat.history') }}</p>
              <el-icon>
                <ArrowUp v-if="state.showHistory" />
                <ArrowDown v-else />
              </el-icon>
            </div>
          </div>
        </template>
        <div v-show="state.showHistory" class="px-3">
          <div v-if="state.isLoadingConv" class="flex-center mt-2">
            <el-icon class="animate-spin">
              <Loading></Loading>
            </el-icon>
          </div>
          <template v-for="conv in groupedConversations" :key="conv.key">
            <template v-if="conv.list.length">
              <div class="h-[30px] pl-6 pr-2 flex items-center text-xs text-[#9A9A9A] mb-1">
                {{ $t(conv.key) }}
              </div>
              <template v-for="item in conv.list" :key="item.conversation_id">
                <div
                  class="h-9 pl-6 pr-2 rounded-md flex items-center gap-2 mt-0.5 group cursor-pointer hover:bg-[#ECEDEE] overflow-hidden"
                  :class="[currentConv.conversation_id === item.conversation_id ? 'bg-[#F1F2F3]' : '']"
                  @click="convStore.setCurrentState(item.agent_id, item.conversation_id)"
                >
                  <span class="flex-1 text-sm text-[#1D1E1F] truncate">{{ item.title }}</span>
                  <el-dropdown @command="handleCommandConv($event, item)">
                    <div class="size-7 flex-center cursor-pointer invisible group-hover:visible max-md:visible" @click.stop>
                      <el-icon size="14">
                        <MoreFilled />
                      </el-icon>
                    </div>
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item command="edit">
                          <svg-icon name="edit" class="mr-1"></svg-icon>
                          {{ $t('action.rename') }}
                        </el-dropdown-item>
                        <el-dropdown-item command="del">
                          <span class="text-[#FA5151] flex-center">
                            <svg-icon name="del" class="mr-1"></svg-icon>
                            {{ $t('action.del') }}
                          </span>
                        </el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </div>
              </template>
            </template>
          </template>

          <div
            v-if="convStore.conversations.length > 5 && !state.expandHistory"
            class="h-9 pl-6 pr-2 flex items-center text-sm text-[#9A9A9A] mb-1 cursor-pointer"
            @click="state.expandHistory = !state.expandHistory"
          >
            {{ $t('action.view_more') }}
          </div>
        </div>
      </div>
    </div>
    <div v-if="userStore.is_login" class="h-14 px-4 flex items-center gap-2 border-t border-b md:hidden" @click="handleJumpToProfile">
      <el-image class="flex-none size-6" :src="userStore.info.avatar" :alt="userStore.info.nickname"></el-image>
      <div class="flex-1 flex items-center overflow-hidden">
        <p class="text-sm text-primary font-medium truncate">{{ userStore.info.nickname }}</p>
        <div
          v-if="!userStore.info.is_internal"
          class="flex-none h-6 flex-center gap-1 bg-[#F7F7F7] rounded-full px-2 text-xs text-placeholder whitespace-nowrap"
          :title="userStore.info.group_name"
          @click.stop="upgradeRef.open"
        >
          <img
            :src="
              !/\.png$/.test(userStore.info.group_icon)
                ? $getPublicPath(`/images/subscription/${userStore.info.group_icon}.png`)
                : userStore.info.group_icon
            "
            class="w-4 h-4 object-cover"
            alt="用户组图标"
          />
          <p class="max-w-[5em] truncate">{{ userStore.info.group_name }}</p>
        </div>
      </div>
      <el-icon size="14" color="#939499">
        <ArrowRight />
      </el-icon>
    </div>
  </div>

  <el-dialog v-model="state.editVisible" :title="$t('chat.edit_conversation')" width="480">
    <el-input v-model="convForm.title" v-trim size="large" :placeholder="$t('form.input_placeholder')"></el-input>

    <template #footer>
      <el-button @click="state.editVisible = false">{{ $t('action.cancel') }}</el-button>
      <el-button v-debounce type="primary" :disabled="!convForm.title.trim()" @click="handleEditConv">{{ $t('action.confirm') }}</el-button>
    </template>
  </el-dialog>

  <Upgrade ref="upgradeRef" />
</template>

<script setup lang="ts">
import { reactive, onMounted, onUnmounted, computed, ref, nextTick, watchEffect } from 'vue'
import { ArrowUp, ArrowDown, Search, ArrowRight, MoreFilled, Loading } from '@element-plus/icons-vue'
import { useRoute } from 'vue-router'
import { onClickOutside } from '@vueuse/core'
import { router } from '@/router'

import Upgrade from '@/components/Upgrade/index.vue'

import { useUserStore } from '@/stores/modules/user'
import { useAgentStore } from '@/stores/modules/agent'
import { useGlobalStore } from '@/stores/modules/global'
import { useConversationStore } from '@/stores/modules/conversation'
import { useEnterpriseStore } from '@/stores/modules/enterprise'
import { useLinksStore } from '@/stores/modules/links'
import { useNavigationStore } from '@/stores/modules/navigation'
import { NAVIGATION_TYPE, NAVIGATION_TARGET } from '@/constants/navigation'

import { EVENT_NAMES } from '@/constants/events'
import eventBus from '@/utils/event-bus'

const route = useRoute()

const userStore = useUserStore()
const agentStore = useAgentStore()
const globalStore = useGlobalStore()
const linksStore = useLinksStore()
const convStore = useConversationStore()
const enterpriseStore = useEnterpriseStore()
const navigationStore = useNavigationStore()
const SHOW_USUAL_AGENT_LEN = 4

const siderRef = ref()
const upgradeRef = ref()

const userGroups = userStore.info.group_ids || []

const state = reactive({
  keyword: '',
  sidebarCollapsed: false,
  agentCollapsed: false,
  showUsualAgents: true,
  showConversations: true,
  editVisible: false,
  showHistory: false,
  expandHistory: false,
  isLoadingConv: true
})

const convForm = reactive({
  conversation_id: 0,
  title: ''
})

const usualAgents = computed(() => {
  const filterList = convStore.usual_agents.filter((item) => item.user_group_ids.length > 0)
  const agents = state.agentCollapsed ? filterList : filterList.slice(0, SHOW_USUAL_AGENT_LEN)
  return agents
})
const currentAgent = computed(() => convStore.currentAgent)
const currentConv = computed(() => convStore.currentConversation)
const searchUsualAgents = computed(() => {
  if (!state.keyword) return []
  return convStore.usual_agents
    .filter((agent) => agent.name.toLowerCase().includes(state.keyword.toLowerCase()))
    .map((item) => {
      return {
        ...item,
        formt_name: item.name.replace(new RegExp(state.keyword, 'g'), `<span class="text-theme">${state.keyword}</span>`)
      }
    })
})
const searchAgents = computed(() => {
  if (!state.keyword) return []
  return agentStore.agentList
    .filter((item) => {
      const keywordMatch = item.name.toLowerCase().includes(state.keyword.toLowerCase())
      const hasCommonGroup = item.user_group_ids?.some((groupId) => userGroups.includes(groupId))
      return keywordMatch && hasCommonGroup
    })
    .map((item) => {
      return {
        ...item,
        formt_name: item.name.replace(new RegExp(state.keyword, 'g'), `<span class="text-theme">${state.keyword}</span>`)
      }
    })
})
const searchToolBox = computed(() => {
  if (!state.keyword) return []
  return linksStore.links
    .filter((item) => {
      const keywordMatch = item.name.toLowerCase().includes(state.keyword.toLowerCase())
      const hasCommonGroup = item.user_group_ids?.some((groupId) => userGroups.includes(groupId))
      return keywordMatch && hasCommonGroup
    })
    .map((item) => {
      return {
        ...item,
        name: item.name.replace(new RegExp(state.keyword, 'g'), `<span class="text-theme">${state.keyword}</span>`)
      }
    })
})
const searchHistory = computed(() => {
  if (!state.keyword) return []
  return convStore.conversations
    .filter((item) => item.title.toLowerCase().includes(state.keyword.toLowerCase()))
    .map((item) => {
      return {
        ...item,
        name: item.title.replace(new RegExp(state.keyword, 'g'), `<span class="text-theme">${state.keyword}</span>`)
      }
    })
})
// 添加对话分组逻辑
const groupedConversations = computed(() => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  const yesterday = new Date(today)
  yesterday.setDate(yesterday.getDate() - 1)

  const lastWeekStart = new Date(today)
  lastWeekStart.setDate(lastWeekStart.getDate() - 7)

  const historyList = state.expandHistory ? convStore.conversations : convStore.conversations.slice(0, 5)

  return [
    {
      key: 'time.today',
      list: historyList.filter((conv) => {
        const convDate = new Date(conv.updated_time)
        return convDate >= today
      })
    },
    {
      key: 'time.yesterday',
      list: historyList.filter((conv) => {
        const convDate = new Date(conv.updated_time)
        return convDate >= yesterday && convDate < today
      })
    },
    {
      key: 'time.week_last',
      list: historyList.filter((conv) => {
        const convDate = new Date(conv.updated_time)
        return convDate >= lastWeekStart && convDate < yesterday
      })
    },
    {
      key: 'time.week_ago',
      list: historyList.filter((conv) => {
        const convDate = new Date(conv.updated_time)
        return convDate < lastWeekStart
      })
    }
  ]
})

// 点击外部区域关闭侧边栏
onClickOutside(siderRef, () => {
  if (!globalStore.siderVisible && globalStore.siderCollapsed) globalStore.siderCollapsed = false
})

watchEffect(() => {
  globalStore.siderVisible = !globalStore.isMdScreen
})

const toggleAgent = () => {
  state.agentCollapsed = !state.agentCollapsed
}

const handleEditConv = () => {
  return convStore.editConversation(convForm).then(() => {
    state.editVisible = false
  })
}
const delConversation = async (conv: Conversation.Info) => {
  await ElMessageBox.confirm(window.$t('chat.conversation_confirm_delete'), window.$t('action.del'), {
    confirmButtonText: window.$t('action.del'),
    cancelButtonText: window.$t('action.cancel'),
    type: 'warning'
  })
  convStore.delConversation(conv)
}
const handleCommandConv = (event: string, conv: Conversation.Info) => {
  if (event === 'del') {
    delConversation(conv)
  } else if (event === 'edit') {
    convForm.conversation_id = conv.conversation_id
    convForm.title = conv.title
    state.editVisible = true
  }
}

const handleJumpToLibrary = () => {
  if (window.$isElectron) {
    window.$chat53ai.$win({
      type: 'new-tab',
      data: JSON.stringify({
        title: 'AI知识库',
        closable: true,
        ability: false,
        urls: ['http://ziroom.hub.53ai.com/space']
      })
    })
  } else {
    window.open('http://ziroom.hub.53ai.com/space', '_blank')
  }
}

const handleJumpToProfile = () => {
  router.push({ name: 'Profile' })
  globalStore.toggleSider()
}

const loadConversations = () => {
  state.isLoadingConv = true
  convStore.loadConversations().finally(() => {
    state.isLoadingConv = false
  })
}

onMounted(async () => {
  agentStore.loadCategorys()
  navigationStore.fetchNavigations()

  await agentStore.loadAgentList()
  loadConversations()
  eventBus.on(EVENT_NAMES.LOGIN_SUCCESS, () => {
    loadConversations()
  })
  convStore.updateAgents(agentStore.agentList)
  nextTick(() => {
    linksStore.loadLinks()
  })
})

onUnmounted(() => {
  // 清理逻辑
  convStore.clearCurrentState()
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
