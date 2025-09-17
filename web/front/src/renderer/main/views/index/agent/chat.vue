<template>
  <!-- <section class="h-full overflow-hidden flex flex-col w-11/12 lg:w-4/5 pt-6 px-4 mx-auto box-border"> -->
  <section class="h-full overflow-hidden flex flex-col pt-6">
    <div class="relative flex-none flex items-center gap-4 px-4 w-11/12 mx-auto lg:w-4/5 mb-5">
      <ElBreadcrumb class="flex-1 w-0" :separator-icon="ArrowRight">
        <ElBreadcrumbItem v-if="navigationStore.homeNavigation.menu_path" :to="navigationStore.homeNavigation.menu_path" replace>
          <span class="text-regular leading-6 font-normal hover-text-theme">
            {{ navigationStore.homeNavigation.name }}
          </span>
        </ElBreadcrumbItem>
        <ElBreadcrumbItem v-if="navigationStore.agentNavigation.menu_path" :to="navigationStore.agentNavigation.menu_path" replace>
          <span class="text-regular leading-6 font-normal hover-text-theme">
            {{ navigationStore.agentNavigation.name }}
          </span>
        </ElBreadcrumbItem>
        <ElBreadcrumbItem>
          <span class="text-primary leading-6 inline-block truncate max-w-[10em] md:max-w-[30rem]" :title="detailData.name">
            {{ detailData.name }}
          </span>
        </ElBreadcrumbItem>
      </ElBreadcrumb>
      <ElButton v-if="currentAgent?.custom_config_obj?.agent_mode !== 'completion'" link @click="agentDetailRef?.showShare()">
        <SvgIcon class="mr-1" name="share-two" size="18" color="#4F5052" stroke />
        {{ $t('action.share') }}
      </ElButton>
      <ElButton link class="!ml-0" @click="agentDetailRef?.showUseCase()">
        <SvgIcon class="mr-1" name="layout-split" size="18" />
        {{ $t('chat.usage_guide') }}
      </ElButton>
    </div>
    <div ref="boxRef" class="flex-1 px-4 overflow-y-auto">
      <div class="w-11/12 lg:w-4/5 mx-auto flex-1 h-full">
        <AgentDetailView ref="agentDetailRef" hide-menu-header hide-footer show-recommend use-case-fixed show-history />
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ArrowRight } from '@element-plus/icons-vue'

import { onMounted, ref, computed, watch } from 'vue'

import SvgIcon from '@/components/SvgIcon.vue'
import AgentDetailView from '@/views/chat/index.vue'

import { useAgentStore } from '@/stores/modules/agent'
import { useNavigationStore } from '@/stores/modules/navigation'
import { useConversationStore } from '@/stores/modules/conversation'

const agentStore = useAgentStore()
const navigationStore = useNavigationStore()
const convStore = useConversationStore()

const agentDetailRef = ref<InstanceType<typeof AgentDetailView>>()
const boxRef = ref(null)

const currentAgent = computed(() => convStore.currentAgent)
const detailData = computed(() => {
  return agentDetailRef.value?.detailData || { name: '' }
})

watch(
  () => agentStore.boxHeight,
  (newData) => {
    if (newData) {
      boxRef.value.scrollTop = newData
    }
  }
)

onMounted(() => {
  navigationStore.fetchNavigations()
})
</script>

<style scoped></style>
