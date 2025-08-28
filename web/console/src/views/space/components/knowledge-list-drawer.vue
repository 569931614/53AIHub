<template>
  <el-drawer v-model="visible" :title="$t('knowledge.name')" size="80%">
    <TablePlus
      :data="tableData"
      :total="tableTotal"
      style="width: 100%"
      header-row-class-name="rounded overflow-hidden"
      header-cell-class-name="!bg-[#F6F7F8] !h-[60px] !border-none"
      @page-size-change="handleSizeChange"
      @page-current-change="handleCurrentChange"
    >
      <ElTableColumn :label="$t('common.name')" min-width="160" prop="name" show-overflow-tooltip />
      <ElTableColumn
        :label="$t('knowledge.short_name')"
        min-width="140"
        prop="description"
        show-overflow-tooltip
      ></ElTableColumn>
      <ElTableColumn :label="$t('common.member')" width="160" prop="nickname" show-overflow-tooltip>
        <template #default="{ row }">
          <span :class="{ 'text-[#9B9B9B]': !row.nickname }">
            {{ row.nickname || '--' }}
          </span>
        </template>
      </ElTableColumn>
      <ElTableColumn :label="$t('common.creator')" prop="capacity" width="120" show-overflow-tooltip>
        <template #default="{ row }">
          {{ row.capacity || '--' }}
        </template>
      </ElTableColumn>
      <ElTableColumn :label="$t('created_time')" min-width="140" prop="created_time" show-overflow-tooltip>
        <template #default="{ row }">
          <span :class="{ 'text-[#9B9B9B]': !row.created_time }">
            {{ row.created_time }}
          </span>
        </template>
      </ElTableColumn>
      <ElTableColumn :label="$t('operation')" width="170" fixed="right" align="right">
        <template #default="{ row }">
          <ElButton class="text-[#5A6D9E] !bg-transparent" type="text" @click.stop="handleEdit(row)">
            {{ $t('action.manage') }}
          </ElButton>
          <ElTooltip :content="$t('space.delete_tip', { total: 111 })" placement="top" trigger="hover">
            <span>
              <ElButton
                :disabled="11 > 0"
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
  </el-drawer>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { librariesApi } from '@/api/modules/libraries/index'
import type { SpaceItem } from '@/api/modules/spaces/types'
import type { LibraryItem } from '@/api/modules/libraries/types'

const visible = ref(false)
const space_id = ref(0)
const tableData = ref<LibraryItem[]>([])
const tableTotal = ref(0)

const loadList = () => {
  librariesApi
    .list({
      space_id: space_id.value,
      offset: 0,
      limit: 10,
    })
    .then(res => {
      tableData.value = res.list
      tableTotal.value = res.total
    })
}

const open = (data: SpaceItem) => {
  space_id.value = data.id
  loadList()
  visible.value = true
}

const close = () => {
  visible.value = false
}

const handleSizeChange = (size: number) => {
  console.log(size)
}

const handleCurrentChange = (current: number) => {
  console.log(current)
}

const handleEdit = (row: LibraryItem) => {
  console.log('Edit:', row)
}

const handleDelete = (row: LibraryItem) => {
  console.log('Delete:', row)
}

defineExpose({
  open,
  close,
})
</script>
