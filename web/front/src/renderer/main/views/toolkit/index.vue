<template>
  <div class="">
    <MainHeader v-if="!hideHeader" sticky>
      <template #before_suffix>
        <div class="text-base text-primary font-bold line-clamp-1 max-md:flex-1 max-md:text-center">
          {{ navigationStore.toolkitNavigation.name || $t('module.toolbox') }}
        </div>
      </template>

      <template #after_prefix>
        <!-- <span class="flex items-center gap-1 text-sm cursor-pointer md:hidden" @click="$router.back()">
          <svg-icon name="return" size="18" stroke></svg-icon>
        </span> -->
      </template>
    </MainHeader>
    <div ref="scrollRef" v-loading="loading" class="box-border">
      <div class="w-11/12 lg:w-4/5 mx-auto min-h-[80vh]">
        <div
          ref="headerRef"
          class="left-0 z-[2] sticky top-0 bg-white flex md:flex-row flex-col-reverse gap-5 items-stretch md:items-center justify-between py-6 md:py-8 px-[16px] box-border rounded-[4px] overflow-hidden"
          :class="[navOffset ? '' : 'lg:py-10']"
          :style="{ top: navOffset ? navOffset + 'px' : '0' }"
          style="width: calc(100% + 32px); transform: translateX(-16px)"
        >
          <el-tabs
            v-model="state.group_id"
            class="index-tabs flex-1 overflow-hidden"
            style="--el-tabs-header-height: 36px"
            @tab-click="handleTabChange"
          >
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
          class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-5"
          :keyword="state.keyword"
          :list="links"
          :group-id="state.group_id"
          :only-all="onlyAll"
        />
      </div>
    </div>
    <MFooter v-if="!hideFooter" />
  </div>
</template>

<script setup lang="ts">
import { Search } from '@element-plus/icons-vue'
import { reactive, ref, computed, onMounted, inject } from 'vue'
import SearchInput from '@/components/Search/index.vue'
import MainHeader from '@/layout/header.vue'
import MFooter from '@/layout/m-footer.vue'
import ToolkitList from './components/list.vue'

import { useLinksStore } from '@/stores/modules/links'
import { useUserStore } from '@/stores/modules/user'
import { useNavigationStore } from '@/stores/modules/navigation'

const props = withDefaults(
  defineProps<{
    hideHeader?: boolean
    hideFooter?: boolean
    onlyAll?: boolean
    stickyOffset?: number
  }>(),
  {
    hideHeader: false,
    hideFooter: false,
    onlyAll: false,
    stickyOffset: 70
  }
)

const userStore = useUserStore()
const linksStore = useLinksStore()
const navigationStore = useNavigationStore()
const state = reactive({
  group_id: 0,
  keyword: ''
})

const navOffset = computed(() => {
  return props.hideHeader ? props.stickyOffset : props.stickyOffset
})

const loading = ref(false)

const categorys = computed(() => {
  const categorysList = linksStore.categorys
  return categorysList.filter((item) => {
    if (item.group_id === 0) return true
    // 产品要分类都可以展示出来
    return true
    // let list = linksStore.links.filter(row => row.group_id == item.group_id) || []
    // return !list || list.length === 0 ? false : true
  })
})

const links = computed(() => {
  return linksStore.links.filter((item) => {
    if (state.group_id === 0) return true
    return item.group_id === state.group_id
  })
})

const scrollRef = ref<HTMLElement>()
const mainScrollRef = inject('mainScrollRef')
const handleTabChange = () => {
  const viewId = `#group_${state.group_id}`
  const scrollContainer = mainScrollRef ? mainScrollRef.value : document.querySelector('main')
  const targetElement = scrollContainer.querySelector(viewId)
  if (targetElement) {
    const containerRect = scrollContainer.getBoundingClientRect()
    const targetRect = targetElement.getBoundingClientRect()
    const scrollTop = scrollContainer.scrollTop + targetRect.top - containerRect.top - (navOffset.value + 150)
    scrollContainer.scrollTo({
      top: Math.max(0, scrollTop), // 确保不会滚动到负值
      behavior: 'smooth'
    })
  }
}

const headerRef = ref<HTMLElement>()
onMounted(() => {
  loading.value = true
  linksStore.loadCategorys()
  linksStore.loadLinks().finally(() => {
    loading.value = false
  })
})
</script>

<style scoped></style>
