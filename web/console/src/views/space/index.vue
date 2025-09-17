<template>
  <Layout class="px-[60px] py-8">
    <Header :title="$t('module.space')" />
    <div class="flex-1 flex flex-col bg-white p-6 mt-3 box-border max-h-[calc(100vh-100px)] overflow-auto">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3">
          <ElInput
            v-model="filter_form.name"
            style="max-width: 268px"
            size="large"
            clearable
            :prefix-icon="Search"
            :placeholder="$t('space.search_placeholder')"
            @change="refresh"
          />
        </div>
        <div class="flex-center gap-3">
          <ElButton type="primary" size="large" @click="handleAdd"> + {{ $t('action_add') }} </ElButton>
        </div>
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
          <ElTableColumn :label="$t('space.name')" min-width="160" prop="name" show-overflow-tooltip>
            <template #default="{ row }">
              <div class="flex items-center gap-2">
                <img :src="row.icon" class="size-8" />
                <span>
                  {{ row.name }}
                </span>
              </div>
            </template>
          </ElTableColumn>

          <ElTableColumn
            :label="$t('knowledge.name')"
            min-width="140"
            prop="library_count"
            show-overflow-tooltip
          ></ElTableColumn>
          <ElTableColumn :label="$t('created_time')" min-width="160" prop="created_time" show-overflow-tooltip>
            <template #default="{ row }">
              {{ row.created_time }}
            </template>
          </ElTableColumn>

          <ElTableColumn :label="$t('common.creator')" width="160" prop="owner_info.nickname" show-overflow-tooltip>
            <template #default="{ row }">
              <template v-if="row.is_default">
                {{ $t('space.system') }}
              </template>
              <template v-else>
                {{ row.owner_info.nickname }}
              </template>
            </template>
          </ElTableColumn>
          <!-- <ElTableColumn :label="$t('space.capacity')" prop="capacity" width="120" show-overflow-tooltip>
            <template #default="{ row }">
              {{ row.capacity || '--' }}
            </template>
          </ElTableColumn> -->
          <ElTableColumn :label="$t('operation')" width="170" fixed="right" align="right">
            <template #default="{ row }">
              <ElButton class="text-[#5A6D9E] !bg-transparent" type="text" @click.stop="handleEdit(row)">
                {{ $t('action_edit') }}
              </ElButton>
              <!-- <ElButton class="text-[#5A6D9E] !bg-transparent" type="text" @click.stop="handleManage(row)">
                {{ $t('knowledge.name') }}
              </ElButton> -->
              <ElTooltip
                v-if="row.is_default"
                :content="$t('space.system_space_delete_tip')"
                placement="top"
                trigger="hover"
              >
                <span class="ml-2">
                  <ElButton disabled class="text-[#5A6D9E] !bg-transparent" type="text" @click.stop="handleDelete(row)">
                    {{ $t('action_delete') }}
                  </ElButton>
                </span>
              </ElTooltip>
              <ElTooltip
                v-else
                :disabled="row.library_count === 0"
                :content="$t('space.delete_tip', { total: row.library_count })"
                placement="top"
                trigger="hover"
              >
                <span class="ml-2">
                  <ElButton
                    :disabled="row.library_count > 0"
                    class="text-[#5A6D9E] !bg-transparent"
                    type="text"
                    @click.stop="handleDelete(row)"
                  >
                    {{ $t('action_delete') }}
                  </ElButton>
                </span>
              </ElTooltip>
            </template>
          </ElTableColumn>
        </TablePlus>
      </div>
    </div>
    <InfoSaveDialog ref="infoSaveDialogRef" @refresh="refresh" />
    <KnowledgeListDrawer ref="knowledgeListDrawerRef" @refresh="refresh" />
  </Layout>
</template>

<script setup lang="ts">
import { Search } from '@element-plus/icons-vue'
import { onMounted, reactive, ref } from 'vue'
import spacesApi from '@/api/modules/spaces/index'
import type { SpaceItem, SpaceDisplayItem } from '@/api/modules/spaces/types'
import { transformSpaceList } from '@/api/modules/spaces/transform'

import InfoSaveDialog from './components/info-save-dialog.vue'
import KnowledgeListDrawer from './components/knowledge-list-drawer.vue'

const infoSaveDialogRef = ref<InstanceType<typeof InfoSaveDialog>>()
const knowledgeListDrawerRef = ref<InstanceType<typeof KnowledgeListDrawer>>()
const filter_form = reactive({
  name: '',
  offset: 0,
  limit: 10,
})

const tableData = ref<SpaceDisplayItem[]>([])
const tableTotal = ref(0)
const tableLoading = ref(false)

const loadList = async () => {
  tableLoading.value = true
  const res = await spacesApi.list({ ...filter_form }).finally(() => {
    tableLoading.value = false
  })
  tableTotal.value = res.count
  tableData.value = transformSpaceList(res.spaces)
}

const refresh = () => {
  filter_form.offset = 0
  loadList()
}
const handleSizeChange = (size: number) => {
  filter_form.limit = size
  refresh()
}
const handleCurrentChange = (page: number) => {
  filter_form.offset = (page - 1) * filter_form.limit
  loadList()
}

const handleAdd = () => {
  infoSaveDialogRef.value?.open()
}
const handleEdit = (data: SpaceItem) => {
  infoSaveDialogRef.value?.open(data)
}
// const handleManage = (data: SpaceItem) => {
//   knowledgeListDrawerRef.value?.open(data)
// }
const handleDelete = async (data: SpaceItem) => {
  await ElMessageBox.prompt(
    window.$t('space.delete_confirm_tip', { name: data.name }),
    window.$t('space.delete_confirm'),
    {
      inputPattern: new RegExp(data.name.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')),
      dangerouslyUseHTMLString: true,
      buttonSize: 'large',
      cancelButtonClass: 'w-24',
      confirmButtonClass: 'w-24 !bg-[#FA5151] !border-none',
      inputPlaceholder: window.$t('space.delete_confirm_tip_placeholder'),
      inputErrorMessage: window.$t('space.delete_confirm_tip_placeholder'),
    }
  ).then(async ({ value }) => {
    if (value === data.name) {
      await spacesApi.delete(data.id)
      ElMessage.success(window.$t('action_delete_success'))
      refresh()
    }
  })
}

onMounted(async () => {
  refresh()
})
</script>

<style scoped></style>
