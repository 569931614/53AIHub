<template>
  <el-tooltip
    ref="tooltipRef"
    effect="light"
    trigger="click"
    :popper-class="['el-popper--plain', { 'el-popper--m': isSmScreen }]"
    :placement="isSmScreen ? 'bottom' : 'top-start'"
  >
    <template #content>
      <div class="p-5 w-[596px] max-md:w-full">
        <div class="flex items-center justify-between relative">
          <div class="absolute right-0">
            <div
              v-tooltip="{ content: $t('action.close') }"
              class="size-5 flex-center rounded cursor-pointer hover:bg-[#E1E2E3]"
              @click="handleClose"
            >
              <el-icon>
                <Close />
              </el-icon>
            </div>
          </div>
          <h4 class="text-lg text-primary w-full max-md:text-center">{{ $t('action.find') }}</h4>
        </div>
        <el-input
          v-model="state.keyword"
          size="large"
          :placeholder="$t('action.search') + $t('module.agent')"
          :prefix-icon="Search"
          class="mt-4 el-input--main"
        ></el-input>
        <el-tabs v-model="state.group_id" class="index-tabs mt-4">
          <template v-for="item in agentStore.categorys" :key="item.group_id">
            <el-tab-pane :label="item.group_name" :name="item.group_id"></el-tab-pane>
          </template>
        </el-tabs>
        <div class="h-[300px] overflow-y-auto mt-5">
          <div class="grid gap-4 grid-cols-2 max-md:grid-cols-1">
            <template v-for="item in showAgentList" :key="item.agent_id">
              <div class="flex items-center p-3 bg-[#F8F9FA] rounded-lg cursor-pointer" @click="handleSelect(item)">
                <el-avatar :size="36" :src="item.logo" class="mr-2" />
                <div class="flex-1 overflow-hidden">
                  <h3 class="text-base font-medium text-primary mb-1">{{ item.name }}</h3>
                  <p class="text-sm text-secondary truncate">
                    {{ item.description || '- -' }}
                  </p>
                </div>
              </div>
            </template>
          </div>
          <el-empty
            v-if="showAgentList.length === 0"
            :description="$t('common.no_data')"
            :image-size="100"
            :image="$getPublicPath('/images/chat/completion_empty.png')"
          />
        </div>
      </div>
    </template>

    <slot></slot>
  </el-tooltip>
</template>

<script setup lang="ts">
import { Close, Search } from '@element-plus/icons-vue'
import { computed, reactive, ref } from 'vue'

import { useAgentStore } from '@/stores/modules/agent'
import { useBasicLayout } from '@/hooks/useBasicLayout'

const { isSmScreen } = useBasicLayout()

const emits = defineEmits<{
  (e: 'select', data: Agent.State): void
}>()

const agentStore = useAgentStore()

const tooltipRef = ref()

const state: {
  keyword: string
  group_id: number
} = reactive({
  keyword: '',
  group_id: 0
})

const showAgentList = computed(() => {
  const filterList = agentStore.agentList.filter((item) => item.user_group_ids.length > 0)
  // if (!state.keyword) {
  //   return state.group_id === 0 ? agentStore.agentList : agentStore.agentList.filter((item) => item.group_id === state.group_id)
  // }
  if (!state.keyword) {
    return state.group_id === 0 ? filterList : filterList.filter((item) => item.group_id === state.group_id)
  }
  const keyword = state.keyword.toLowerCase().trim()
  // return agentStore.agentList.filter((item) => {
  //   const matchKeyword = item.name.toLowerCase().includes(keyword) || (item.description && item.description.toLowerCase().includes(keyword))
  //   return (state.group_id === 0 || item.group_id === state.group_id) && matchKeyword
  // })
  return filterList.filter((item) => {
    const matchKeyword = item.name.toLowerCase().includes(keyword) || (item.description && item.description.toLowerCase().includes(keyword))
    return (state.group_id === 0 || item.group_id === state.group_id) && matchKeyword
  })
})

const handleClose = () => {
  tooltipRef.value.hide()
}
const handleSelect = (item: Agent.State) => {
  emits('select', item)
  handleClose()
}
</script>

<style scoped></style>
