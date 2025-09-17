<template>
  <Layout class="px-[60px] py-8">
    <Header :title="$t('module.agent')">
      <template #right>
        <el-button
          v-version="{ module: VERSION_MODULE.AGENT, count: all_total, content: $t('version.agent_limit') }"
          type="primary"
          size="large"
          @click="add_visible = true"
        >
          + {{ $t('action_add') }}
        </el-button>
      </template>
    </Header>

    <div class="flex-1 overflow-y-auto bg-white rounded-lg px-10 py-6 mt-4">
      <div class="flex items-center justify-between mb-4">
        <div class="flex flex-1 w-0 gap-2">
          <ElSelect
            v-model="filter_form.group"
            size="large"
            :placeholder="$t('all')"
            class="flex-none max-w-[150px]"
            @change="refresh"
          >
            <template #prefix> {{ $t('internal_user.group.title') }}: </template>
            <ElOption v-if="filter_form.group" :label="$t('all_group')" :value="''" class="min-w-[204px]" />
            <ElOption
              v-for="item in groupList"
              :key="item.group_id"
              :label="item.group_name"
              :value="item.group_id"
              class="min-w-[204px] bg-[#ffffff]"
            />

            <!-- 固定底部操作栏 -->
            <template #footer>
              <div class="flex justify-between items-center mx-2 px-1 pt-2 sticky bottom-0 bg-white z-10">
                <div class="cursor-pointer text-[#5A6D9E] text-sm hover:opacity-80" @click="openDialog">
                  {{ $t('group_management') }}
                </div>
              </div>
            </template>
          </ElSelect>

          <!-- 一个智能体可以在多个分组中 -->
          <!-- <GroupTabs
          ref="groupTabsRef"
          v-model="filter_form.group"
          type="dropdown"
          :group-type="GROUP_TYPE.AGENT"
          :is-agent-group="true"
          @change="refresh"
          @get-options="refresh"
          @open="openDialog"
        /> -->

          <ElSelect
            v-model="filter_form.platform"
            size="large"
            class="flex-none max-w-[160px]"
            clearable
            :placeholder="$t('all')"
            @change="refresh"
          >
            <template #prefix> {{ $t('module.platform_v2') }}: </template>
            <ElOption v-if="filter_form.platform" :label="$t('module.all_platform')" :value="''" />
            <ElOption
              v-for="item in channels"
              :key="item.label"
              :label="$t(item.label)"
              :value="item.label === 'provider_platform.prompt' ? '1,3,44,36' : item.channelType"
            />
          </ElSelect>

          <ElSelect
            v-model="filter_form.type"
            size="large"
            class="flex-none max-w-[160px]"
            clearable
            :placeholder="$t('all')"
            @change="refresh"
          >
            <template #prefix> {{ $t('type') }}: </template>
            <ElOption v-if="filter_form.type" :label="$t('all_type')" :value="''" />
            <ElOption :label="$t('agent_type_completion_v2')" value="1" />
            <ElOption :label="$t('agent_type_chat_v2')" value="0" />
          </ElSelect>
        </div>
        <div class="flex-none flex-center gap-3 ml-8">
          <!-- <Search v-model="filter_form.keyword" placeholder="module.agent_search_placeholder" @change="refresh" /> -->
          <ElInput
            v-model="filter_form.keyword"
            size="large"
            clearable
            :suffix-icon="Search"
            :placeholder="$t('agent.name_v2')"
            @change="refresh"
          />
        </div>
      </div>

      <TablePlus
        header-row-class-name="rounded overflow-hidden"
        header-cell-class-name="!bg-[#F6F7F8] !h-[60px] !border-none"
        :data="table_data"
        :total="table_total"
        :loading="table_loading"
        :page="filter_form.page"
        :limit="filter_form.page_size"
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
                  {{ row.description || '' }}
                </div>
              </div>
            </div>
          </template>
        </ElTableColumn>
        <ElTableColumn :label="$t('module.platform_v2')" min-width="140" show-overflow-tooltip>
          <template #default="{ row = {} }">
            {{ $t(`agent_app.${row.agent_type}`) || '--' }}
          </template>
        </ElTableColumn>
        <ElTableColumn :label="$t('type')" min-width="140" show-overflow-tooltip>
          <template #default="{ row = {} }">
            {{ row.backend_agent_type === 0 ? $t('agent_type_chat_v2') : $t('agent_type_completion_v2') }}
          </template>
        </ElTableColumn>
        <ElTableColumn :label="$t('usage_range')" min-width="140" show-overflow-tooltip>
          <template #default="{ row }">
            <span :class="!row.user_group_names.length ? 'text-[#999]' : ''">
              {{ row.user_group_names.join('、') || '--' }}
            </span>
          </template>
        </ElTableColumn>
        <ElTableColumn :label="$t('action_enable')" min-width="80">
          <template #default="{ row }">
            <ElSwitch v-model="row.enable" @change="onAgentStatusChange({ data: row })" />
          </template>
        </ElTableColumn>
        <ElTableColumn :label="$t('operation')" width="120" align="right" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="onAgentAdd(row.agent_type, row)">
              {{ $t('action_edit') }}
            </el-button>
            <el-button type="primary" link @click="onAgentDelete({ data: row })">
              {{ $t('action_delete') }}
            </el-button>
          </template>
        </ElTableColumn>
      </TablePlus>
    </div>
  </Layout>

  <el-drawer v-model="add_visible" :title="$t('action_add')" size="650px" @opened="onAddOpened">
    <ul class="w-full min-h-[300px] overflow-y-auto">
      <li v-for="(item, itemIndex) in filteredAgentOptions" :key="itemIndex">
        <h4 class="text-sm text-[#939499]">
          {{ $t(item.title) }}
        </h4>
        <ul class="flex flex-col gap-3 pt-4 pb-6">
          <li
            v-for="row in item.filteredChildren"
            :key="row.value"
            class="h-[72px] px-6 rounded flex items-center gap-3 bg-[#F8F9FA] cursor-pointer hover:shadow"
          >
            <img class="flex-none size-10 rounded-lg" :src="row.icon" alt="" />
            <div class="flex-1 text-base text-[#1D1E1F] truncate">
              {{ $t(row.label) }}
            </div>
            <ElButton type="primary" plain class="border-none" @click="handleAgentPrepare(row)">
              {{ $t('action_add') }}
            </ElButton>
          </li>
        </ul>
      </li>
    </ul>
  </el-drawer>

  <CreateDrawer ref="createDrawerRef" @success="row => onAgentAdd(row.agent_type, row, true)" />
  <GroupDialog ref="dialogRef" :group-type="GROUP_TYPE.AGENT" @change="handleGroupChange" />
</template>

<script setup name="Agent" lang="ts">
import { nextTick, onMounted, onUnmounted, reactive, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Search } from '@element-plus/icons-vue'

import { ElMessage } from 'element-plus'
import CreateDrawer from '@/views/agent/create/drawer.vue'
import GroupDialog from '@/components/GroupDialog/index.vue'

import eventBus from '@/utils/event-bus'

import { agentApi } from '@/api/modules/agent'
import { providerApi } from '@/api/modules/provider'
import { subscriptionApi } from '@/api/modules/subscription'
import { groupApi } from '@/api/modules/group'
import { getProvidersByAuth, getProviderByAgentId, AgentType, channels } from '@/constants/platform/config'
import { AGENT_APP_OPTIONS } from '@/constants/platform/agent'
import { VERSION_MODULE } from '@/constants/enterprise'
import { GROUP_TYPE } from '@/constants/group'

interface SubscriptionItem {
  group_id: number
  group_name: string
}

interface ProviderItem {
  provider_type: string
  [key: string]: any
}

const router = useRouter()
const createDrawerRef = ref()
const dialogRef = ref()

const filter_form = reactive({
  group: '',
  platform: '',
  type: '',
  keyword: '',
  page: 1,
  page_size: 10,
})
const all_total = ref(0)
const table_data = ref<Agent.State[]>([])
const table_total = ref(0)
const table_loading = ref(false)
const add_visible = ref(false)
const subscriptionList = ref<SubscriptionItem[]>([])
const groupList = ref<SubscriptionItem[]>([])

const auth_providers = ref<ProviderItem[]>([])

const openDialog = () => {
  dialogRef.value?.open()
}

const filteredAgentOptions = computed(() => {
  const authMap = new Map(auth_providers.value.map(provider => [provider.provider_type, provider.is_auth]))
  return AGENT_APP_OPTIONS.map(item => {
    const filteredChildren = item.children.filter(row => {
      const provider = getProviderByAgentId(row.value)

      if (!provider?.auth) {
        return true
      }
      // 需要授权的应用，检查是否已授权
      return authMap.get(provider.id) === true
    })

    return {
      ...item,
      filteredChildren,
    }
  }).filter(item => {
    return item.filteredChildren.length > 0
  })
})

const loadSubscriptionList = async () => {
  if (!subscriptionList.value.length)
    subscriptionList.value = await subscriptionApi.list({ params: { offset: 0, limit: 1000 } })
}

const loadGroupList = async () => {
  if (!groupList.value.length) groupList.value = await groupApi.list({ params: { group_type: GROUP_TYPE.AGENT } })
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
const loadListData = async () => {
  table_loading.value = true
  await loadSubscriptionList()
  await loadGroupList()

  try {
    const { count = 0, agents = [] } = await agentApi.list({
      params: {
        group_id: filter_form.group,
        channel_types: filter_form.platform,
        agent_types: filter_form.type,
        keyword: filter_form.keyword,
        offset: (filter_form.page - 1) * filter_form.page_size,
        limit: filter_form.page_size,
      },
    })

    table_total.value = count
    table_data.value = []
    await nextTick()

    table_data.value = agents.map((item: Partial<Agent.State> = {}) => {
      const agent = item as Agent.State
      agent.user_group_ids = agent.user_group_ids || []
      agent.user_group_names = agent.user_group_ids
        .map(value => {
          const subscription = subscriptionList.value.find(row => row.group_id === value)
          return subscription?.group_name || ''
        })
        .filter(group_name => !!group_name)

      return agent
    })
  } finally {
    table_loading.value = false
  }

  loadAllTotal()
}

const refresh = async () => {
  filter_form.page = 1
  await loadListData()
}

const handleGroupChange = (result: { value: Group[] }) => {
  groupList.value = result.value
  filter_form.group = ''
  refresh()
}

const onTableSizeChange = (size: number) => {
  filter_form.page_size = size
  refresh()
}

const onTableCurrentChange = (current: number) => {
  filter_form.page = current
  loadListData()
}

const onAgentDelete = async ({ data: { agent_id } }: { data: { agent_id: number } }) => {
  await ElMessageBox.confirm(window.$t('agent_delete_confirm'), window.$t('action_delete'))
  await agentApi.delete({ data: { agent_id } })
  ElMessage.success(window.$t('action_delete_success'))
  loadListData()
}

const loadProviderList = async () => {
  const list = (await providerApi.list()) as ProviderItem[]

  auth_providers.value = getProvidersByAuth(true).map(item => {
    const provider_type = item.id
    return {
      ...item,
      provider_type,
      is_auth: !!list.find(row => row.provider_type === provider_type),
    }
  })
}

const checkAuth = (value: AgentType): Promise<void> => {
  return new Promise<void>((resolve, reject) => {
    const provider = getProviderByAgentId(value)
    const auth_provider = auth_providers.value.find(row => row.provider_type === provider.id)

    if (auth_provider && !auth_provider?.is_auth) {
      reject(new Error('Authentication required'))
      return ElMessageBox.confirm(
        window.$t('auth_required', { provider_name: window.$t(provider.label) }),
        window.$t('tip'),
        {
          confirmButtonText: window.$t('action_go'),
          cancelButtonText: window.$t('action_cancel'),
          type: 'warning',
        }
      ).then(() => {
        router.push({
          name: 'Platform',
        })
      })
    }

    resolve()
  })
}

const handleAgentPrepare = async (data: { value: string; channel_type: number }) => {
  await checkAuth(data.value)
  createDrawerRef.value.open({
    agent_type: data.value,
    group_id: +filter_form.group_id > 0 ? filter_form.group_id : undefined,
    data: { channel_config: { channel_type: data.channel_type } },
  })

  add_visible.value = false
}

const onAgentAdd = async (value: string, data: Partial<Agent.State> = {}, is_new = false) => {
  await checkAuth(value)
  loadListData()
  await router.push({
    name: 'AgentCreate',
    query: data.agent_id
      ? {
          type: data.agent_type,
          agent_id: data.agent_id,
          is_new,
        }
      : {
          type: value,
          group_id: +filter_form.group_id > 0 ? filter_form.group_id : undefined,
          is_new,
        },
  })
}

const onAgentStatusChange = async ({ data: { agent_id, enable } }: { data: { agent_id: number; enable: boolean } }) => {
  await agentApi.updateStatus({ data: { agent_id, enable } })
  ElMessage.success(window.$t(enable ? 'action_enable_success' : 'action_disable_success'))
}

// 添加onAddOpened函数修复错误
const onAddOpened = () => {
  // 抽屉打开时的处理逻辑
}

onMounted(() => {
  refresh()
  loadProviderList()
  eventBus.on('user-login-success', refresh)
  eventBus.on('agent-change', loadListData)
})

onUnmounted(() => {
  eventBus.off('user-login-success', refresh)
  eventBus.off('agent-change', loadListData)
})
</script>

<style lang="scss" scoped>
::v-deep(.el-table__cell) {
  padding: 14px 0;
}

::v-deep(.el-select__placeholder) {
  color: black !important;
}
</style>
