<template>
  <div class="h-full flex">
    <Sider v-if="showSidebar"></Sider>
    <main ref="mainRef" class="flex-1 overflow-y-auto">
      <RouterView v-slot="{ Component, route }">
        <component :is="Component" :key="route.path" />
      </RouterView>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, provide, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useEnterpriseStore } from '@/stores/modules/enterprise'
import { useAgentStore } from '@/stores/modules/agent'

import Sider from './sider.vue'

const mainRef = ref<HTMLElement>()
const route = useRoute()
const enterpriseStore = useEnterpriseStore()
const agentStore = useAgentStore()

provide('mainRef', mainRef)

const showSidebar = computed(() => {
  const routes = ['Chat', 'Agent', 'Toolkit', 'Prompt', 'PromptDetail', 'Profile', 'Order']
  return routes.includes(route.name as string) && enterpriseStore.isSoftStyle
})

watch(
  () => agentStore.boxHeight,
  (newHeight) => {
    mainRef.value.scrollTop = newHeight
  }
)

onMounted(() => {})
</script>

<style></style>
