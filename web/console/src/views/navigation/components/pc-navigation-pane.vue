<template>
  <div
    class="w-full py-3 px-6 box-border rounded flex flex-row items-center bg-white"
    style="box-shadow: 0 5px 15px rgb(22 23 26 / 10%)"
  >
    <img
      class="flex-none w-10 h-10 object-contain"
      src="https://img.js.design/assets/img/61889010bb6d5d11bc62c2b8.png#029311e2b191cdf96801016808d17a46"
      alt=""
    />
    <ElTabs class="flex-1 w-0 ml-8">
      <ElTabPane v-for="item in nav_active_list" :key="item.name" :label="item.name" :name="item.name" />
    </ElTabs>
    <ElButton class="flex-none w-[128px] h-[36px] ml-8 border-none" type="primary" plain @click.stop="handleCreate">
      添加({{ nav_active_list.length }}/{{ nav_list.length }})
    </ElButton>
  </div>

  <div class="mt-5 rounded py-9 px-10 bg-white">
    <ElTable
      ref="table_ref"
      class="w-full rounded"
      :data="nav_list"
      :header-cell-style="{ background: '#F6F7F8', height: '60px' }"
      :cell-style="{ height: '66px' }"
      height="calc(100vh - 296px)"
    >
      <ElTableColumn width="40">
        <div class="pr-3 sort-icon cursor-move">
          <svg-icon name="drag" width="24px" height="32px" color="#a1a5af" />
        </div>
      </ElTableColumn>
      <ElTableColumn prop="name" :label="$t('module.nav_name')" min-width="120" show-overflow-tooltip />
      <ElTableColumn :label="$t('module.nav_type')" show-overflow-tooltip min-width="100">
        <template #default="{ row = {} }">
          {{ NAVIGATION_TYPE_LABEL_MAP.get(row.type) }}
        </template>
      </ElTableColumn>
      <ElTableColumn :label="$t('module.nav_target')" show-overflow-tooltip min-width="100">
        <template #default="{ row = {} }">
          {{ NAVIGATION_TARGET_LABEL_MAP.get(row.target) }}
        </template>
      </ElTableColumn>
      <ElTableColumn :label="$t('module.nav_visible')" width="120">
        <template #default="scope">
          <ElSwitch v-model="scope.row.visible" size="default" />
        </template>
      </ElTableColumn>
      <ElTableColumn :label="$t('module.nav_operation')" align="right" width="180" fixed="right">
        <template #default="{ row = {}, $index }">
          <ElButton
            class="h-auto !p-0 leading-none text-[#576D9C]"
            type="text"
            size="default"
            @click.stop="handleCreate({ data: row })"
          >
            {{ $t('action_edit') }}
          </ElButton>
          <ElButton
            class="h-auto !p-0 leading-none text-[#576D9C]"
            type="text"
            size="default"
            @click.stop="handlePreview({ data: row })"
          >
            {{ $t('action_preview') }}
          </ElButton>
          <ElButton
            class="h-auto !p-0 leading-none text-[#576D9C]"
            type="text"
            size="default"
            @click.stop="handleSeoSetting({ data: row })"
          >
            {{ $t('module.nav_operation_seo') }}
          </ElButton>
          <ElButton
            class="h-auto !p-0 leading-none text-[#576D9C]"
            type="text"
            size="default"
            @click.stop="handleDelete({ data: row, index: $index })"
          >
            {{ $t('action_delete') }}
          </ElButton>
        </template>
      </ElTableColumn>
    </ElTable>
  </div>

  <SeoSettingDialog ref="seo_setting_dialog_ref" />
  <NavCreateDialog ref="nav_create_dialog_ref" @confirm="handleCreateConfirm" />
</template>

<script setup lang="ts">
import { ref, computed, onUnmounted, watch, nextTick } from 'vue'
import Sortable from 'sortablejs'
import SeoSettingDialog from './seo-setting-dialog.vue'
import NavCreateDialog from './nav-create-dialog.vue'

import { NAVIGATION_TYPE_LABEL_MAP, NAVIGATION_TARGET_LABEL_MAP, NAVIGATION_INIT_DATA } from '@/constants/navigation'

// Types
type SortableEvent = {
  target: HTMLElement
  oldIndex: number
  newIndex: number
  from: HTMLElement
  to: HTMLElement
}

const table_ref = ref()
const seo_setting_dialog_ref = ref()
const nav_create_dialog_ref = ref()
const nav_list = ref(JSON.parse(JSON.stringify(NAVIGATION_INIT_DATA)))

const nav_active_list = computed(() => {
  return nav_list.value.filter(item => item.visible)
})

const handlePreview = ({ data = {} } = {}) => {
  if (data.url) window.open(data.url, '_blank')
}
const handleSeoSetting = ({ data = {} } = {}) => {
  seo_setting_dialog_ref.value.open({ data })
}
const handleDelete = async ({ index = -1 } = {}) => {
  await ElMessageBox.confirm(window.$t('module.action_delete_confirm'), window.$t('action_delete_tip'))
  if (index > 0) nav_list.value.splice(index, 1)
  ElMessage.success(window.$t('action_delete_success'))
}
const handleCreate = ({ data = {} } = {}) => {
  nav_create_dialog_ref.value.open({ data })
}
const handleCreateConfirm = async ({ data = {} } = {}) => {
  if (!data.id) {
    data.id = nav_list.value.length + 1
    const list = nav_list.value
    list.push(data)
    nav_list.value = []
    await nextTick()
    nav_list.value = list
  }
}
let _sortable_instance: Sortable | null = null
const initSortable = () => {
  if (!table_ref.value || !table_ref.value.$el || !nav_list.value.length) return
  const sortable_el = table_ref.value.$el.querySelector('.el-table__body tbody')
  _sortable_instance = Sortable.create(sortable_el, {
    onStart: (event: SortableEvent) => {
      const { target, oldIndex } = event
      target.children[oldIndex].style.background = '#ECF5FF'
    },
    onEnd: async (event: SortableEvent) => {
      const { target, newIndex: target_index, oldIndex: origin_index } = event
      if (target.children && target.children[target_index])
        target.children[target_index].style.background = 'transparent'
      if (target_index === origin_index) return
      const list = nav_list.value
      const origin_data = list.splice(origin_index, 1)[0]
      list.splice(target_index, 0, origin_data)
      nav_list.value = []
      await nextTick()
      nav_list.value = [...list].map((item, index) => ({ ...item, sort: index }))
      ElMessage.success(window.$t('action_sort_success'))
    },
  })
}
const destroySortable = () => {
  if (_sortable_instance) {
    _sortable_instance.destroy()
    _sortable_instance = null
  }
}

onUnmounted(() => {
  destroySortable()
})

watch(
  () => nav_list.value,
  async () => {
    destroySortable()
    await nextTick()
    initSortable()
  },
  { deep: true, immediate: true }
)
</script>

<style scoped lang="scss">
::v-deep(.el-tabs__content),
::v-deep(.el-tabs__nav-wrap::after) {
  display: none;
}

::v-deep(.el-tabs__header) {
  margin-bottom: 0;
}

::v-deep(.el-tabs__item) {
  pointer-events: none;
  font-size: 16px;
  font-weight: 400;
  color: rgb(29 30 31 / 100%);
}
</style>
