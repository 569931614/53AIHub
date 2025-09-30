<template>
  <ElDrawer
    v-model="visible"
    :title="$t(originData.label || '')"
    size="70%"
    destroy-on-close
    append-to-body
    :close-on-click-modal="false"
  >
    <div class="flex items-center justify-between gap-4 mb-4">
      <div>
        <ElInput
          v-model="filterForm.name"
          :prefix-icon="Search"
          :placeholder="$t('action_search')"
          size="large"
          clearable
          @change="refresh"
        />
      </div>
      <ElButton type="primary" size="large" @click="handleAddProvider">
        {{ $t('action_add') }}
      </ElButton>
    </div>

    <TablePlus
      header-row-class-name="rounded overflow-hidden"
      header-cell-class-name="!bg-[#F6F7F8] !h-[60px] !border-none"
      :data="tableData"
      :total="tableTotal"
      :loading="tableLoading"
      :page="1"
      :limit="tableTotal"
      :pagination="false"
    >
      <ElTableColumn prop="name" :label="$t('module.website_info_name')" min-width="100" show-overflow-tooltip>
      </ElTableColumn>

      <ElTableColumn
        v-if="[PROVIDER_VALUE['53AI'], PROVIDER_VALUE.COZE_OSV].includes(originData.provider_type)"
        prop="base_url"
        :label="$t('module.platform_tool_api_endpoint')"
        min-width="180"
        show-overflow-tooltip
      >
      </ElTableColumn>

      <ElTableColumn
        v-if="originData.provider_type === PROVIDER_VALUE.COZE_CN"
        prop="configs.client_id"
        :label="$t('module.platform_auth_client_id')"
        min-width="180"
        show-overflow-tooltip
      >
      </ElTableColumn>

      <ElTableColumn prop="created_time" width="180" :label="$t('add_time')"> </ElTableColumn>

      <ElTableColumn :label="$t('operation')" width="120" align="right" fixed="right">
        <template #default="{ row }">
          <ElButton type="primary" link @click="handleEditProvider(row)">
            {{ $t('action_edit') }}
          </ElButton>
          <ElButton type="primary" link @click="handleDeleteProvider(row)">
            {{ $t('action_delete') }}
          </ElButton>
        </template>
      </ElTableColumn>
    </TablePlus>
  </ElDrawer>

  <ProviderAuthorizeDialog ref="authorizeRef" @success="onSuccess" />
</template>

<script setup lang="ts">
import { Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { nextTick, reactive, ref } from 'vue'

import ProviderAuthorizeDialog from './provider-authorize-dialog.vue'

import { PROVIDER_VALUE } from '@/constants/platform/provider'
import { providerApi } from '@/api/modules/provider'

const emit = defineEmits<{
  (e: 'change'): void
}>()

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

const authorizeRef = ref()
const visible = ref(false)
const filterForm = reactive({
  name: '',
})

const tableData = ref<Agent.State[]>([])
const tableTotal = ref(0)
const tableLoading = ref(false)
const originData = ref<Record<string, any>>({})

const loadList = async () => {
  tableLoading.value = true
  // await loadSubscriptionList()
  try {
    const res = await providerApi.list({ params: { ...filterForm, providerType: originData.value.provider_type } })
    tableData.value = res
    tableTotal.value = res.length
  } finally {
    tableLoading.value = false
  }
}

const refresh = () => {
  return loadList()
}

const open = async ({ data = {}, type = PROVIDER_VALUE.DIFY } = {}) => {
  originData.value = data
  tableData.value = []
  await nextTick()
  visible.value = true
  refresh()
}

const close = () => {
  visible.value = false
}

const handleAddProvider = () => {
  authorizeRef.value.open({ data: originData.value })
}

const handleEditProvider = (data: AgentData) => {
  authorizeRef.value.open({ data })
}

const handleDeleteProvider = async (data: AgentData) => {
  await ElMessageBox.confirm(window.$t('module.platform_delete_confirm'))
  await providerApi.delete({ data: { provider_id: data.provider_id } })
  ElMessage.success(window.$t('action_delete_success'))
  loadList()
  emit('change')
}

const onSuccess = () => {
  loadList()
  emit('change')
}

defineExpose({
  open,
  close,
})
</script>

<style scoped></style>
