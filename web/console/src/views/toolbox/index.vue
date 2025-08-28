<template>
  <Layout class="px-[60px] py-8">
    <Header :title="$t('module.ai_toolbox')">
      <template #right>
        <template v-if="isSort">
          <el-button type="default" class="min-w-[96px]" size="large" @click="handleSortCancel">
            {{ $t('action_cancel') }}
          </el-button>
          <el-button v-debounce type="primary" class="min-w-[96px] !ml-0" size="large" @click="handleSortSave">
            {{ $t('action_save') }}
          </el-button>
        </template>
        <template v-else>
          <el-button
            type="default"
            class="min-w-[96px]"
            size="large"
            :disabled="!!filterForm.keyword"
            @click="handleSortOpen"
          >
            {{ $t('action_sort') }}
          </el-button>
          <el-button type="primary" class="min-w-[96px] !ml-0" size="large" @click="handleStore">
            {{ $t('action_add') }}
          </el-button>
        </template>
      </template>
    </Header>
    <div class="flex-none flex items-center justify-between mt-5">
      <div class="flex-1 w-0">
        <GroupTabs
          ref="groupTabsRef"
          v-model="filterForm.group_id"
          :disabled="isSort"
          :group-type="GROUP_TYPE.AI_LINK"
          type="dropdown"
          @change="refresh"
          @get-options="onGetOptions"
        />
      </div>
      <div class="flex-none flex-center gap-3 ml-8">
        <ElInput
          v-model="filterForm.keyword"
          size="large"
          clearable
          :placeholder="$t('prompt.search_placeholder')"
          :suffix-icon="Search"
          @change="refresh"
        />
      </div>
    </div>

    <div v-loading="loading || saving" class="mt-6 flex-1 box-border pr-1 overflow-y-auto relative">
      <ElEmpty v-if="!aiLinkList.length" class="mt-10" :description="$t('no_data')" />
      <ul v-else>
        <li v-for="group in showGroupOptions" :key="group.group_id">
          <div
            v-show="!+filterForm.group_id || +filterForm.group_id < 0"
            class="text-[#4F5052] text-opacity-60 text-sm mb-4"
          >
            {{ group.group_name }}
          </div>
          <Sortable
            :model-value="group.children"
            identity="ai_link_id"
            class="grid grid-cols-4 gap-4 max-lg:grid-cols-3 max-md:grid-cols-2 max-sm:grid-cols-1 mb-9"
            :custom-sortable-id="`toolbox_group_${group.group_id}`"
            :disabled="!isSort"
            group="toolbox-group"
            force-render
            :loading="saving"
            @change="handleSortChange"
          >
            <template #item="{ item, index }">
              <div
                class="h-[72px] bg-white overflow-hidden group relative border rounded p-4 flex items-center gap-2 cursor-pointer"
                role="button"
                :aria-label="item.name"
              >
                <div
                  v-if="!isSort"
                  class="invisible group-hover:visible w-full h-full z-[2] absolute top-0 left-0 bg-black bg-opacity-40 flex-center gap-1.5"
                >
                  <ElButton type="default" @click="onUrlOpen({ data: item })">
                    {{ $t('action_visit') }}
                  </ElButton>
                  <ElButton type="primary" class="!ml-0" @click="onCommand('edit', item)">
                    {{ $t('action_edit') }}
                  </ElButton>
                  <ElButton type="default" class="!ml-0 !px-2" @click="onCommand('del', item)">
                    <ElIcon color="#FA5151" :size="16">
                      <Delete />
                    </ElIcon>
                  </ElButton>
                </div>
                <img class="w-10 h-10 object-cover rounded-full overflow-hidden" :src="item.logo" />
                <div class="flex-1 w-0">
                  <div class="flex items-center justify-between">
                    <div class="text-sm text-[#1D1E1F] font-semibold line-clamp-1">
                      {{ item.name }}
                    </div>
                  </div>
                  <div class="text-sm text-[#1D1E1F] text-opacity-60 line-clamp-1">
                    {{ item.description }}
                  </div>
                </div>
                <div v-show="isSort" class="sort-icon cursor-move">
                  <svg-icon name="drag" width="24px" height="42px" color="#a1a5af" />
                </div>
              </div>
            </template>
          </Sortable>
        </li>
      </ul>
    </div>

    <!-- <ToolBox class="mt-4" /> -->

    <createDrawer ref="drawerRef" @success="refresh" />
    <StoreDialog ref="storeRef" @add="onCreate" @success="refresh" />
  </Layout>
</template>

<script setup lang="ts">
import { Delete, Search } from '@element-plus/icons-vue'
import { computed, onMounted, provide, reactive, ref } from 'vue'
import StoreDialog from './components/store-dialog.vue'
import createDrawer from './components/create-drawer.vue'

import { aiLinkApi } from '@/api/modules/ai-link'

import { GROUP_TYPE } from '@/constants/group'

const groupTabsRef = ref()
const storeRef = ref()
const drawerRef = ref()
const filterForm = reactive({
  group_id: [],
  keyword: '',
})
const aiLinkList = ref([])
const loading = ref(false)

const showGroupOptions = computed(() => {
  let options = groupOptions.value.filter(item => +item.group_id > 0 && item.children && item.children.length)
  if (filterForm.group_id.length > 0) options = options.filter(item => filterForm.group_id.includes(item.group_id))
  return options
})

const refresh = async () => {
  await loadListData()
}

const loadListData = async () => {
  loading.value = true
  const data = await aiLinkApi.list({ params: filterForm }).finally(() => {
    loading.value = false
  })
  aiLinkList.value = data
  setAiLinkToGroup()
}
const handleStore = () => {
  storeRef.value.open()
}
const onCreate = ({ data = {} } = {}) => {
  if (filterForm.group_id.length > 0 && filterForm.group_id[0] !== '-1') {
    data.group_id = filterForm.group_id[0]
  }
  drawerRef.value.open({ data })
}
const groupOptions = ref([])
provide('groupOptions', groupOptions)

const onGetOptions = ({ options = [] } = {}) => {
  const allOption = { group_id: '-1', group_name: $t('all_group'), children: [] }
  groupOptions.value = [
    allOption,
    ...options.map(item => ({
      ...item,
      children: [],
    })),
  ]
  setAiLinkToGroup()
}
const setAiLinkToGroup = () => {
  groupOptions.value.forEach(item => {
    item.children = aiLinkList.value.filter(aiLink => aiLink.group_id === item.group_id)
    item.children.sort((a, b) => b.sort - a.sort)
  })
}
const onUrlOpen = ({ data = {} } = {}) => {
  window.open(data.url, '_blank')
}
const onCommand = async (command, data = {}) => {
  switch (command) {
    case 'edit':
      onCreate({ data })
      break
    case 'del':
      await ElMessageBox.confirm(window.$t('action_delete_confirm'), window.$t('action_delete_tip'))
      await aiLinkApi.delete({ data: { ai_link_id: data.ai_link_id } })
      ElMessage.success(window.$t('action_delete_success'))
      refresh()
      break
  }
}
const isSort = ref(false)
const saving = ref(false)
const handleSortOpen = () => {
  isSort.value = true
}
const handleSortCancel = () => {
  isSort.value = false
  setAiLinkToGroup()
}
const handleSortSave = async () => {
  saving.value = true
  const items = groupOptions.value.reduce((result, item = {}) => {
    return [
      ...result,
      ...item.children.map((child, index) => ({
        group_id: item.group_id,
        id: child.ai_link_id,
        // name: child.name,
        sort: item.children.length - index,
      })),
    ]
  }, [])
  await aiLinkApi.sort({ items }).finally(() => {
    saving.value = false
  })
  ElMessage.success(window.$t('action_save_success'))
  handleSortCancel()
  refresh()
}
const timer = null
const handleSortChange = (result = {}) => {
  const { action, originSortableId = '', targetSortableId = '', value = [] } = result
  const originGroupId = originSortableId.split('_').pop() || ''
  const targetGroupId = targetSortableId.split('_').pop() || ''
  const originGroupOption = groupOptions.value.find(item => item.group_id === +originGroupId)
  const targetGroupOption = groupOptions.value.find(item => item.group_id === +targetGroupId)
  switch (action) {
    case 'add':
      if (targetGroupOption) targetGroupOption.children = [...value]
      break
    case 'remove':
      if (originGroupOption) originGroupOption.children = [...value]
      break
    case 'sort':
      if (targetGroupOption) targetGroupOption.children = [...value]
      break
  }
}

onMounted(() => {
  refresh()
})
</script>

<style scoped></style>
