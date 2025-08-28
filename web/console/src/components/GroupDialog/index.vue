<template>
  <ElDialog
    v-model="visible"
    :title="$t('group_management')"
    width="680px"
    align-center
    destroy-on-close
    :close-on-click-modal="false"
  >
    <div class="text-[#182B50] text-opacity-60 text-sm pb-4">
      {{ $t('display_order') }}
    </div>
    <Sortable
      v-model="options"
      v-loading="submitting"
      identity="group_id"
      custom-sortable-id="group_sortable"
      class="w-full flex flex-col gap-4 max-h-[60vh] overflow-y-auto"
    >
      <template #item="{ item, index }">
        <div class="flex items-center">
          <div class="pr-3 sort-icon cursor-move">
            <svg-icon name="drag" width="16px" height="32px" color="#a1a5af" />
          </div>
          <div class="flex-1">
            <el-input
              v-model="item.group_name"
              size="large"
              :placeholder="$t('form_input_placeholder')"
              class="w-full"
              maxlength="10"
              show-word-limit
            />
          </div>
          <el-icon class="ml-4 cursor-pointer" color="rgba(24, 43, 80, 0.4)" @click="onGroupRemove({ index })">
            <Delete />
          </el-icon>
        </div>
      </template>
    </Sortable>
    <ElButton type="primary" link class="mt-4 ml-5" @click="onGroupAdd"> +{{ $t('action_add') }} </ElButton>

    <template #footer>
      <div class="py-4 flex items-center justify-center">
        <ElButton class="w-[96px] h-[36px]" type="primary" :loading="submitting" @click="onSave">
          {{ $t('action_save') }}
        </ElButton>
        <ElButton class="w-[96px] h-[36px] text-[#1D1E1F]" type="info" plain @click.stop="close">
          {{ $t('action_cancel') }}
        </ElButton>
      </div>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
import { Delete } from '@element-plus/icons-vue'

import { ref, onMounted, nextTick } from 'vue'
import type { Group } from '@/api/modules/group'
import { groupApi, DEFAULT_GROUP_DATA } from '@/api/modules/group'
import { deepCopy } from '@/utils'
import type { GroupType } from '@/constants/group'

const props = withDefaults(
  defineProps<{
    groupType: GroupType
    beforeRemove?: (args: { data: Group }) => boolean | Promise<boolean>
  }>(),
  {
    groupType: 1,
  }
)

const emits = defineEmits<{
  (e: 'change', result: { value: Group[] }): any
}>()

const visible = ref(false)
const options = ref<Group[]>([])
const originalOptions = ref<Group[]>([])
const submitting = ref(false)
let deleted_group_list: Group[] = []

onMounted(() => {
  refresh()
})

const refresh = async () => {
  const list = await groupApi.list({ params: { group_type: props.groupType } })
  options.value = [...list]
  originalOptions.value = [...list]
  if (!options.value.length) options.value = [deepCopy(DEFAULT_GROUP_DATA)]
  else emits('change', { value: options.value })
}
const open = async ({ value = [] } = {}) => {
  submitting.value = false
  options.value = []
  await nextTick()
  options.value = [...originalOptions.value]
  if (deleted_group_list.length) await refresh()
  deleted_group_list = []
  if (value && value.length) options.value = [...value]
  visible.value = true
}
const close = async () => {
  visible.value = false
}
const onGroupAdd = async () => {
  options.value.push({
    ...deepCopy(DEFAULT_GROUP_DATA),
    group_id: -Date.now(),
    sort: options.value.length,
  })
  await nextTick()
  const scrollerEl = document.querySelector('#group_sortable')
  scrollerEl.scrollTop = scrollerEl.scrollHeight
}
const onGroupRemove = async ({ index }: { index: number }) => {
  const data = options.value[index]
  if (options.value.filter(item => item.group_id).length === 1) return ElMessage.warning(window.$t('group_min_one'))
  const intercept = props.beforeRemove ? props.beforeRemove({ data }) : true
  if (intercept === false) return
  if (intercept.then && intercept.catch) {
    const res = await intercept
    if (res === false) return
  }
  await ElMessageBox.confirm(window.$t('group_delete_confirm'), window.$t('action_delete'))
  deleted_group_list.push(data)
  options.value.splice(index, 1)
}
const deleteGroup = async () => {
  for (let i = 0; i < deleted_group_list.length; i++) {
    const item = deleted_group_list[i]
    if (item.group_id > 0) await groupApi.delete({ data: { group_id: item.group_id } })
  }
}
const onSave = async () => {
  if (submitting.value) return
  const list = options.value.filter(item => item.group_name.trim())
  if (!list.length) return ElMessage.warning(window.$t('group_not_empty'))
  submitting.value = true
  await deleteGroup().catch(() => {
    submitting.value = false
  })
  await groupApi
    .save({
      data: {
        group_type: props.groupType,
        groups: list.map((item, index) => ({
          group_name: item.group_name.trim(),
          group_id: item.group_id > 0 ? item.group_id : 0,
          sort: list.length - index,
        })),
      },
    })
    .catch(() => {
      submitting.value = false
    })
  deleted_group_list = []
  ElMessage.success(window.$t('action_save_success'))
  refresh()
  close()
}

defineExpose({
  open,
  close,
})
</script>

<style scoped></style>
