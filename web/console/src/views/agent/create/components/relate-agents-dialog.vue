<template>
  <ElDialog v-model="visible" :title="$t('action_add')" width="870px">
    <div class="h-[560px] flex flex-col">
      <div class="flex-none flex items-center justify-between">
        <div class="flex-1 w-0">
          <GroupTabs v-model="filterForm.group_id" :group-type="GROUP_TYPE.AGENT" @change="refresh" />
        </div>
        <div class="flex-none flex-center gap-3 ml-8">
          <Search
            v-model="filterForm.keyword"
            placeholder="module.ai_toolbox_search_placeholder"
            @change="refresh"
            @input="refresh"
          />
        </div>
      </div>

      <div v-loading="tableData.loading" class="flex-1 mt-4 flex flex-col gap-4 overflow-y-auto relative">
        <template v-for="item in tableData.list" :key="item.agent_id">
          <div class="bg-white rounded border p-4 flex items-center gap-2 relative group">
            <img alt="AI 搜索" :src="item.logo" class="size-10" />
            <div class="flex-1 overflow-hidden">
              <div class="text-sm text-[#1D1E1F]">
                {{ item.name }}
              </div>
              <div class="text-xs text-[#1D1E1F] text-opacity-60 truncate mt-1">
                {{ item.description }}
              </div>
            </div>
            <el-button v-if="isExist(item.agent_id)" type="primary" plain class="border-none" :disabled="true">{{
              $t('action.added')
            }}</el-button>
            <el-button v-else type="primary" class="invisible group-hover:visible" @click="handleAdd(item)">
              {{ $t('action_add') }}
            </el-button>
          </div>
        </template>

        <el-empty v-if="showEmpty" />
      </div>

      <Pagination
        :total="tableData.total"
        :page-size="filterForm.pageSize"
        :current-page="filterForm.page"
        @current-change="handleCurrentChange"
        @size-change="handleSizeChange"
      />
    </div>
  </ElDialog>
  <RelateAgentsSetting ref="relateAgentsSettingRef" />
</template>

<script setup lang="ts">
import { ref, computed, reactive } from 'vue';

import RelateAgentsSetting from './relate-agents-setting.vue';

import agentApi from '@/api/modules/agent';

import { GROUP_TYPE } from '@/constants/group';

const emit = defineEmits<{
  (e: 'select', item: Agent.State): void
}>()

const visible = ref(false)

const filterForm = reactive({
  group_id: 0,
  keyword: '',
  page: 1,
  pageSize: 20,
})

const tableData: {
  loading: boolean
  total: number
  list: Agent.State[]
} = reactive({
  loading: true,
  total: 0,
  list: [],
})

const showEmpty = computed(() => {
  return tableData.list.length === 0 && !tableData.loading
})

const existAgents = ref<Agent.RelateAgent[]>([])

const loadList = async () => {
  try {
    tableData.loading = true
    const { count = 0, agents = [] } = await agentApi.list({
      params: {
        group_id: filterForm.group_id,
        keyword: filterForm.keyword,
        offset: (filterForm.page - 1) * filterForm.pageSize,
        limit: filterForm.pageSize,
      },
    })
    tableData.total = count
    tableData.list = agents
  } finally {
    tableData.loading = false
  }
}

const refresh = () => {
  filterForm.page = 1
  loadList()
}

const handleCurrentChange = (page: number) => {
  filterForm.page = page
  loadList()
}

const handleSizeChange = (size: number) => {
  filterForm.pageSize = size
  loadList()
}

const open = (relateAgents: Agent.RelateAgent[]) => {
  existAgents.value = relateAgents
  visible.value = true
}

const close = () => {
  visible.value = false
}
const isExist = (agent_id: number) => {
  return (existAgents.value || []).some(item => item.agent_id === agent_id)
}

const handleAdd = (item: Agent.State) => {
  emit('select', item)
  close()
}

defineExpose({
  open,
  close,
})
</script>

<style></style>
