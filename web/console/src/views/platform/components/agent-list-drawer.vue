<template>
  <ElDrawer
    v-model="visible"
    :title="drawerTitle"
    size="70%"
    destroy-on-close
    append-to-body
    :close-on-click-modal="false"
  >
    <div class="flex items-center justify-between gap-4 mb-4">
      <div>
        <ElInput
          v-model="filterForm.keyword"
          :prefix-icon="Search"
          :placeholder="$t('action_search')"
          size="large"
          @change="refresh"
        />
      </div>
      <ElButton
        v-version="{ module: VERSION_MODULE.AGENT, count: all_total, content: $t('version.agent_limit') }"
        type="primary"
        size="large"
        @click="onAgentCreate"
      >
        {{ $t('action_add') }}
      </ElButton>
    </div>

    <TablePlus
      header-row-class-name="rounded overflow-hidden"
      header-cell-class-name="!bg-[#F6F7F8] !h-[60px] !border-none"
      :data="tableData"
      :total="tableTotal"
      :loading="tableLoading"
      :page="filterForm.offset + 1"
      :limit="filterForm.limit"
      @page-size-change="onTableSizeChange"
      @page-current-change="onTableCurrentChange"
    >
      <ElTableColumn prop="date" :label="$t('module.agent')" min-width="180" show-overflow-tooltip>
        <template #default="{ row }">
          <div class="flex items-center gap-2 w-full">
            <img class="flex-none w-8 h-8 rounded-full overflow-hidden" :src="row.logo" alt="" />
            <div class="flex-1 w-0 text-sm flex flex-col">
              <div class="text-[#2563EB] truncate">
                {{ row.name || '--' }}
              </div>
              <div v-show="row.description" class="text-xs text-[#808080] truncate">
                {{ row.description }}
              </div>
            </div>
          </div>
        </template>
      </ElTableColumn>

      <ElTableColumn :label="$t('usage_range')" min-width="140" show-overflow-tooltip>
        <template #default="{ row }">
          <span :class="{ 'text-[#999]': !row.user_group_names.length }">
            {{ row.user_group_names.join('、') || '--' }}
          </span>
        </template>
      </ElTableColumn>

      <ElTableColumn :label="$t('sort')" width="80" show-overflow-tooltip>
        <template #default="{ row }">
          {{ row.sort }}
        </template>
      </ElTableColumn>

      <ElTableColumn :label="$t('action_enable')" width="80">
        <template #default="{ row }">
          <ElSwitch v-model="row.enable" @change="onAgentStatusChange({ data: row })" />
        </template>
      </ElTableColumn>

      <ElTableColumn :label="$t('operation')" width="120" align="right" fixed="right">
        <template #default="{ row }">
          <ElButton type="primary" link @click="onAgentCreate({ data: row })">
            {{ $t('action_edit') }}
          </ElButton>
          <ElButton type="primary" link @click="onAgentDelete({ data: row })">
            {{ $t('action_delete') }}
          </ElButton>
        </template>
      </ElTableColumn>
    </TablePlus>
  </ElDrawer>

  <AgentCreateDrawer ref="agentCreateRef" @success="onAgentCreateSuccess" />
</template>

<script setup lang="ts">
import { Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, nextTick, reactive, ref } from 'vue'
import AgentCreateDrawer from '@/views/agent/create/drawer.vue'

import { PROVIDER_VALUE } from '@/constants/platform/provider'
import { AGENT_TYPE, agentApi } from '@/api/modules/agent'
import { subscriptionApi } from '@/api/modules/subscription'
import { getAgentByChannelType } from '@/constants/platform/config'

import { VERSION_MODULE } from '@/constants/enterprise'

interface AgentData {
  agent_id?: number
  name?: string
  description?: string
  logo?: string
  enable?: boolean
  sort?: number
  user_group_ids: number[]
  user_group_names: string[]
  channel_config?: {
    channel_type?: string
  }
}

interface SubscriptionData {
  group_id: number
  group_name: string
}

const emits = defineEmits(['change'])

const agentCreateRef = ref()
const visible = ref(false)
const filterForm = reactive({
  channel_types: PROVIDER_VALUE.DIFY as string,
  keyword: '',
  offset: 0,
  limit: 10,
})
const drawerTitle = computed(() => {
  const agent = getAgentByChannelType(filterForm.channel_types)
  return agent && agent.label ? window.$t(agent.label) : ''
})
const tableData = ref<Agent.State[]>([])
const tableTotal = ref(0)
const tableLoading = ref(false)
const originData = ref<Record<string, any>>({})
const subscriptionList = ref<SubscriptionData[]>([])
const all_total = ref(0)

const loadSubscriptionList = async () => {
  if (!subscriptionList.value.length) {
    const response = await subscriptionApi.list({ params: { offset: 0, limit: 1000 } })
    subscriptionList.value = response
  }
}

const loadAllTotal = async () => {
  const { count = 0 } = await agentApi.list({
    params: {
      group_id: '-1',
      keyword: '',
      offset: 0,
      limit: 1,
    },
  })
  all_total.value = count
}

const loadListData = async ({ channel_types = filterForm.channel_types } = {}) => {
  tableLoading.value = true
  await loadSubscriptionList()
  loadAllTotal()
  try {
    const { count = 0, agents = [] } = await agentApi.list({ params: { ...filterForm, channel_types } })

    tableData.value = agents.map((item: Agent.State) => ({
      ...item,
      user_group_ids: item.user_group_ids || [],
      user_group_names: (item.user_group_ids || [])
        .map(value => (subscriptionList.value.find(row => row.group_id === value) || {}).group_name)
        .filter(Boolean),
    }))

    tableTotal.value = count
    originData.value.agent_total = count

    return { count, agents }
  } finally {
    tableLoading.value = false
  }
}

const refresh = () => {
  filterForm.offset = 0
  return loadListData()
}

const open = async ({ data = {}, type = PROVIDER_VALUE.DIFY } = {}) => {
  filterForm.channel_types = type
  originData.value = data
  tableData.value = []
  await nextTick()
  visible.value = true
  refresh()
}

const close = () => {
  visible.value = false
}

const onAgentCreate = ({ data = {} }: { data?: Partial<AgentData> } = {}) => {
  const agent = getAgentByChannelType(filterForm.channel_types)
  const agent_type = agent.name || AGENT_TYPE.DIFY_AGENT
  agentCreateRef.value.open({ agent_type, data: originData.value, agent_id: data.agent_id })
}

const onAgentStatusChange = async ({ data }: { data: AgentData }) => {
  await agentApi.updateStatus({ data: { agent_id: data.agent_id, enable: data.enable } })
  ElMessage.success(window.$t(data.enable ? 'action_enable_success' : 'action_disable_success'))
}

const onAgentDelete = async ({ data }: { data: AgentData }) => {
  await ElMessageBox.confirm(window.$t('agent_delete_confirm'), window.$t('action_delete'))
  await agentApi.delete({ data: { agent_id: data.agent_id } })
  ElMessage.success(window.$t('action_delete_success'))
  loadListData()
}

const onTableSizeChange = (size: number) => {
  filterForm.limit = size
  refresh()
}

const onTableCurrentChange = (page: number) => {
  filterForm.offset = (page - 1) * filterForm.limit
  loadListData()
}

const onAgentCreateSuccess = ({ action = 'update' } = {}) => {
  return action === 'create' ? refresh() : loadListData()
}

defineExpose({
  open,
  close,
  loadListData,
  create: ({ data = {}, type = PROVIDER_VALUE.DIFY } = {}) => {
    filterForm.channel_types = type
    originData.value = data
    onAgentCreate()
  },
})
</script>

<style scoped></style>
