<template>
  <div class="relative w-full overflow-hidden">
    <div
      class="absolute top-0 left-1/2 -translate-x-1/2 w-[1920px] h-[700px] !max-w-none z-1"
      :style="{
        backgroundImage: `url(${$getPublicPath('/images/index/card_bg_v2.png')})`,
        backgroundSize: '1920px 700px',
        backgroundPosition: 'center -68px',
        backgroundRepeat: 'no-repeat'
      }"
    />
    <div class="relative z-2 w-full mt-[110px]">
      <h1
        class="mx-auto text-center text-[60px] font-bold"
        v-html="$t('index.banner_title', { name: `<span class='text-[#007AFF]'>AI</span>` })"
      ></h1>
      <p class="mx-auto text-center text-xl text-[#5E6881] mt-6">{{ $t('index.banner_desc') }}</p>
      <ElInput
        v-model="searchValue"
        :clearable="false"
        class="search-input relative left-1/2 -translate-x-1/2 max-w-[980px] w-[52%] h-[54px] mt-[76px]"
        style="--el-input-border-color: transparent"
        placeholder="输入您想查找的AI资源"
        @input="handleSearch"
      >
        <template #prefix>
          <ElIcon class="text-[#939499]">
            <Search />
          </ElIcon>
        </template>
        <template #suffix>
          <ElButton type="primary" size="large" class="h-[36px] w-[94px] rounded-[32px]" @click="handleSearch">
            {{ $t('action.search') }}
          </ElButton>
        </template>
      </ElInput>
      <div class="mx-auto text-center text-base text-[#757880] mt-[42px] flex items-center justify-center gap-3">
        <span>{{ $t('index.hot_search') }}:</span>
        <div class="hover-text-theme cursor-pointer" @click="handleHotSearch('网页抓取')">网页抓取</div>
        <div class="hover-text-theme cursor-pointer" @click="handleHotSearch('文件操作')">文件操作</div>
        <div class="hover-text-theme cursor-pointer" @click="handleHotSearch('邮件发送')">邮件发送</div>
        <div class="hover-text-theme cursor-pointer" @click="handleHotSearch('图片处理')">图片处理</div>
        <div class="hover-text-theme cursor-pointer" @click="handleHotSearch('浏览器自动化')">浏览器自动化</div>
      </div>
      <!-- 智能体 -->
      <div v-if="searchValue ? showAgentList.length > 0 : true" class="w-11/12 lg:w-4/5 py-6 md:py-8 lg:py-10 mx-auto box-border">
        <p v-if="searchValue" class="text-sm md:text-base mt-3 line-clamp-2 text-regular">
          {{ $t('module.agent') }}
        </p>
        <template v-else>
          <h2 class="text-xl md:text-2xl font-bold" style="color: var(--el-text-color-primary, #1d1e1f)">
            {{ $t('index.agent_recommend') }}
          </h2>
          <p class="text-sm md:text-base mt-3 line-clamp-2 text-regular">
            {{ $t('index.agent_recommend_desc') }}
          </p>
        </template>

        <!-- 功能卡片网格 -->
        <AgentList class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-5 mt-5 md:mt-8" :list="showAgentList"></AgentList>

        <router-link
          v-if="!searchValue && showAgentList.length > 0"
          class="block w-[240px] h-[40px] leading-[40px] border border-primary box-border text-center text-theme mt-[54px] rounded-[24px] mx-auto hover-bg-primary-light-9 transition-all duration-300"
          :to="{ path: '/index/agent' }"
        >
          {{ $t('action.view_more') }}
        </router-link>
      </div>

      <div
        class="w-full py-6 md:py-8 lg:py-10 mx-auto box-border"
        :style="{
          backgroundImage: !searchValue ? `url(${$getPublicPath('/images/index/card_bg_v3.png')})` : '',
          backgroundPosition: 'center center',
          backgroundRepeat: 'no-repeat'
        }"
      >
        <PromptView ref="promptRef" mode="index" hide-filter :show-limit="6" />
      </div>
      <!-- 工具箱 -->
      <div v-if="searchValue ? showToolkitList.length > 0 : true" class="w-11/12 lg:w-4/5 py-6 md:py-8 lg:py-10 mx-auto box-border">
        <p v-if="searchValue" class="text-sm md:text-base mt-3 line-clamp-2 text-regular">
          {{ $t('module.toolbox') }}
        </p>
        <template v-else>
          <h2 class="text-xl md:text-2xl font-bold" style="color: var(--el-text-color-primary, #1d1e1f)">
            {{ $t('index.toolbox_recommend') }}
          </h2>
          <p class="text-sm md:text-base mt-3 line-clamp-2 text-regular">
            {{ $t('index.toolbox_recommend_desc') }}
          </p>
        </template>

        <ToolkitList class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-5 mt-7" :list="showToolkitList" only-all />

        <router-link
          v-if="!searchValue && showToolkitList.length > 0"
          class="block w-[240px] h-[40px] leading-[40px] border border-primary box-border text-center text-theme mt-[54px] rounded-[24px] mx-auto hover-bg-primary-light-9 transition-all duration-300"
          :to="{ name: 'HomeToolkit' }"
        >
          {{ $t('action.view_more') }}
        </router-link>
      </div>
      <div class="w-full h-[100px]" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import { Search } from '@element-plus/icons-vue'
import { useRoute, useRouter } from 'vue-router'
import AgentList from '../agent/components/list.vue'
import ToolkitList from '../toolkit/components/list.vue'
import PromptView from '@/views/prompt/view.vue'

import { useLinksStore } from '@/stores/modules/links'
import { useAgentStore } from '@/stores/modules/agent'
import { useUserStore } from '@/stores/modules/user'

const route = useRoute()
const router = useRouter()
const linksStore = useLinksStore()
const agentStore = useAgentStore()
const userStore = useUserStore()
const promptRef = ref<InstanceType<typeof PromptView>>()
const searchValue = ref('')

const hasPermission = (userGroupIds: number[], itemGroupIds: number[]) => {
  if (!itemGroupIds || itemGroupIds.length === 0) return false
  return userGroupIds.some((groupId) => itemGroupIds.includes(groupId))
}

const showAgentList = computed(() => {
  const filterList = agentStore.agentList.filter((item) => item.user_group_ids.length > 0)
  if (searchValue.value) {
    return filterList.filter((item) => item.name.includes(searchValue.value))
  }
  return filterList.slice(0, 6)
})

const showToolkitList = computed(() => {
  const filterList = linksStore.links.filter(
    (item) => item.user_group_ids.length > 0 && hasPermission(userStore.info.group_ids || [], item.user_group_ids || [])
  )
  if (searchValue.value) {
    return filterList.filter((item) => item.name.includes(searchValue.value))
  }
  return filterList.slice(0, 6)
})

const handleHotSearch = (keyword: string) => {
  searchValue.value = keyword
  handleSearch()
}

const handleSearch = () => {
  promptRef.value?.search(searchValue.value)
}

onMounted(() => {
  linksStore.loadLinks()
  agentStore.loadAgentList()
})
</script>

<style scoped>
.search-input {
  /* background: rgba(255, 255, 255, 1); */
  background: transparent;
  box-shadow: 0 4px 12px rgb(90 188 230 / 12%);
}

.search-input::after {
  content: '';
  display: block;
  width: calc(100% + 6px);
  height: calc(100% + 6px);
  background: linear-gradient(90deg, rgb(125 232 220 / 100%) 0%, rgb(102 150 255 / 100%) 100%);
  position: absolute;
  z-index: -1;
  top: -3px;
  left: -3px;
  border-radius: 32px;
}

::v-deep(.search-input .el-input__wrapper) {
  box-shadow: none;
  border-radius: 32px;
}
</style>
