<template>
  <div
    class="left-0 z-[2] sticky top-0 bg-white flex md:flex-row flex-col-reverse gap-5 items-stretch md:items-center justify-between py-6 md:py-8 box-border rounded overflow-hidden"
    :style="{ top: stickyOffset ? stickyOffset + 'px' : '0' }"
  >
    <el-tabs v-model="state.group_id" class="index-tabs flex-1 overflow-hidden" style="--el-tabs-header-height: 36px" @tab-click="handleTabChange">
      <el-tab-pane v-for="item in categorys" :key="item.group_id" :label="item.group_name" :name="item.group_id" />
    </el-tabs>
    <div>
      <SearchInput v-model="state.keyword" class="hidden md:flex" :placeholder="$t('action.search') + $t('module.prompt')" />
      <ElInput
        v-model="state.keyword"
        size="large"
        class="w-full md:hidden el-input--main"
        :placeholder="$t('toolbox.search_placeholder')"
        :prefix-icon="Search"
      />
    </div>
  </div>
  <ToolkitList
    class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-5"
    :keyword="state.keyword"
    :list="links"
    :group-id="state.group_id"
  />
</template>

<script setup lang="ts">
import { reactive, computed, onMounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import SearchInput from '@/components/Search/index.vue'
import ToolkitList from './list.vue'
import { scrollToElement } from '@/utils/scroll'
import { useLinksStore } from '@/stores/modules/links'

const linksStore = useLinksStore()

const props = defineProps<{
  stickyOffset?: number
}>()

const state = reactive({
  group_id: 0,
  keyword: ''
})

const categorys = computed(() => {
  return linksStore.categorys
})

const links = computed(() => {
  return linksStore.links.filter((item) => {
    if (state.group_id === 0) return true
    return item.group_id === state.group_id
  })
})

const handleTabChange = () => {
  scrollToElement(`#group_${state.group_id}`, (props.stickyOffset || 0) + 150)
}
onMounted(() => {
  linksStore.loadCategorys()
  linksStore.loadLinks()
})
</script>

<style scoped></style>
