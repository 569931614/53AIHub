<template>
  <Layout class="px-15 py-8">
    <Header :title="$t('module.system_log')" />
    <div class="flex-1 flex flex-col bg-white p-6 mt-3 box-border max-h-[calc(100vh-100px)] overflow-auto">
      <div class="flex items-center gap-3">
        <div class="flex-none">
          <FilterDateRange
            v-model="selectDate"
            size="large"
            :value-format="(date: any) => new Date(date).getTime()"
            @change="onRefresh"
          />
        </div>
        <ElSelect
          v-model="filterForm.action"
          size="large"
          class="flex-none max-w-[180px]"
          clearable
          @change="onRefresh"
        >
          <template #prefix> {{ $t('system_log.log_action') }}: </template>
          <ElOption v-for="item in actions" :key="item.value" :label="item.text" :value="item.value" />
        </ElSelect>
        <ElSelect
          v-model="filterForm.module"
          size="large"
          class="flex-none max-w-[180px]"
          clearable
          @change="onRefresh"
        >
          <template #prefix> {{ $t('system_log.log_module') }}: </template>
          <ElOption v-for="item in modules" :key="item.value" :label="item.text" :value="item.value" />
        </ElSelect>

        <FilterUser v-model="userList" class="flex-none max-w-[180px]" type="user" @change="onRefresh" />
      </div>

      <div v-loading="tableLoading" class="flex-1 overflow-y-auto bg-white rounded-lg mt-4">
        <TablePlus
          :data="tableData"
          :total="tableTotal"
          style="width: 100%"
          header-row-class-name="rounded overflow-hidden"
          header-cell-class-name="!bg-[#F6F7F8] !h-[60px] !border-none"
          @page-size-change="handleSizeChange"
          @page-current-change="handleCurrentChange"
        >
          <ElTableColumn :label="$t('system_log.log_time')" min-width="160" prop="action_time" />
          <ElTableColumn :label="$t('system_log.log_action')" min-width="100" prop="action" show-overflow-tooltip>
            <template #default="{ row }">
              <span :class="{ 'text-[#9B9B9B]': !row.action }">
                {{ getActionLabel(row.action) }}
              </span>
            </template>
          </ElTableColumn>
          <ElTableColumn :label="$t('system_log.log_module')" min-width="100" prop="module" show-overflow-tooltip>
            <template #default="{ row }">
              <span :class="{ 'text-[#9B9B9B]': !row.module }">
                {{ getModuleLabel(row.module) }}
              </span>
            </template>
          </ElTableColumn>
          <ElTableColumn :label="$t('system_log.log_operator')" min-width="100" prop="nickname" show-overflow-tooltip>
            <template #default="{ row }">
              <span :class="{ 'text-[#9B9B9B]': !row.nickname }">
                {{ row.nickname }}
              </span>
            </template>
          </ElTableColumn>
          <ElTableColumn :label="$t('system_log.log_label')" prop="content" min-width="140" show-overflow-tooltip />
          <ElTableColumn :label="$t('system_log.log_ip')" min-width="120" show-overflow-tooltip>
            <template #default="{ row }">
              <span :class="{ 'text-[#9B9B9B]': !row.ip }">
                {{ row.ip }}
              </span>
            </template>
          </ElTableColumn>
        </TablePlus>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'

import FilterDateRange from '@/components/Filter/date-range.vue'
import FilterUser from '@/components/Filter/user.vue'

import { systemLogApi } from '@/api/modules/system-log/index'
import { transformSystemLogList, getDefaultSystemLogRequest } from '@/api/modules/system-log/transform'
import type { SystemLogListRequest, ActionItem, ModuleItem, SystemLogDisplayItem } from '@/api/modules/system-log/types'

type UserItem = {
  user_id: number
  nickname: string
}

// 响应式数据
const userList = ref<UserItem[]>([])
const selectDate = ref([])
const filterForm = reactive<SystemLogListRequest>(getDefaultSystemLogRequest())

const actions = ref<ActionItem[]>([])
const modules = ref<ModuleItem[]>([])
const tableData = ref<SystemLogDisplayItem[]>([])
const tableTotal = ref(0)
const tableLoading = ref(false)

// 数据加载方法
const loadList = async (): Promise<void> => {
  try {
    const params = {
      ...filterForm,
      start_time: selectDate.value[0],
      end_time: selectDate.value[1],
      user_id: userList.value.length ? userList.value.map(item => item.user_id).join(',') : undefined,
    }

    const data = await systemLogApi.list(params)
    tableData.value = transformSystemLogList(data.system_logs || [])
    tableTotal.value = data.count || 0
  } catch (error) {
    console.error('加载系统日志失败:', error)
    tableData.value = []
    tableTotal.value = 0
  }
}

const loadActions = async (): Promise<void> => {
  const data = await systemLogApi.actions()
  actions.value = data
}

const loadModules = async (): Promise<void> => {
  const data = await systemLogApi.modules()
  modules.value = data
}

// 工具方法
const getActionLabel = (action: number): string => {
  return actions.value.find(item => item.value === action)?.text || String(action)
}

const getModuleLabel = (module: number): string => {
  return modules.value.find(item => item.value === module)?.text || String(module)
}

// 事件处理方法
const onRefresh = (): void => {
  filterForm.offset = 0
  loadList()
}

const handleSizeChange = (size: number): void => {
  filterForm.limit = size
  onRefresh()
}

const handleCurrentChange = (page: number): void => {
  filterForm.offset = (page - 1) * filterForm.limit
  loadList()
}

// 生命周期
onMounted(async () => {
  tableLoading.value = true
  try {
    await Promise.all([loadList(), loadActions(), loadModules()])
  } finally {
    tableLoading.value = false
  }
})
</script>

<style scoped></style>
