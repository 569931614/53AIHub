<template>
  <div class="max-h-[360px] flex flex-col">
    <div class="flex-none flex items-center justify-between">
      <div class="flex-1 w-0">
        <GroupTabs v-model="filterForm.group_id" type="tabs-pure" :options="categoryList" @change="onGroupChange" />
      </div>
      <div class="flex-none flex-center gap-3 ml-8">
        <Search v-model="filterForm.keyword" placeholder="module.ai_toolbox_search_placeholder" />
      </div>
    </div>

    <div ref="scrollRef" class="flex-1 overflow-y-auto relative">
      <template v-for="item in searchList" :key="item.group_id">
        <div v-if="item.links.length" :id="`toolbox-group-${item.group_id}`" class="text-sm text-[#939499] mt-4">
          {{ item.group_name }}
        </div>
        <div v-if="item.links.length" class="grid grid-cols-3 gap-4 mt-3">
          <template v-for="(link, index) in item.links" :key="index">
            <div class="bg-white rounded border p-4 flex items-center gap-2 relative group">
              <img alt="AI 搜索" :src="link.logo" class="size-10" />
              <div class="flex-1 overflow-hidden">
                <div class="text-sm text-[#1D1E1F]">
                  {{ link.name }}
                </div>
                <div class="text-xs text-[#1D1E1F] text-opacity-60 truncate mt-1">
                  {{ link.description }}
                </div>
              </div>

              <div
                class="absolute inset-0 items-center justify-center bg-[#222326] bg-opacity-55 rounded hidden group-hover:flex gap-2"
              >
                <ElButton class="!mr-0" @click="handleVisit(link)"> {{ $t('action_visit') }} </ElButton
                ><ElButton type="primary" class="!ml-0" @click="handleAdd(link)">{{ $t('action_add') }}</ElButton>
              </div>
            </div>
          </template>
        </div>
      </template>
      <el-empty v-if="showEmpty" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, reactive, onMounted } from 'vue'
import aiLinkApi from '@/api/modules/ai-link'
import { useScroll } from '@/hooks/useScroll'

const { scrollRef, scrollTo, scrollToTop } = useScroll()

interface StoreState extends Category.State {
  links: AI_LINK.State[]
}

const emit = defineEmits(['add'])

const categoryList = ref<StoreState[]>([])
const filterForm = reactive({
  group_id: 0,
  keyword: '',
})

const searchList = computed(() => {
  return categoryList.value
    .filter(item => item.group_id > 0)
    .map(item => {
      return {
        ...item,
        links: item.links.filter(
          link => link.name.includes(filterForm.keyword) || link.description.includes(filterForm.keyword)
        ),
      }
    })
})

const showEmpty = computed(() => {
  return searchList.value.every(item => !item.links.length)
})

const loadList = () => {
  aiLinkApi.store().then(res => {
    const list = res.data.map((item, index) => {
      return {
        group_id: index + 1,
        group_name: item.group_name,
        links: item.links,
      }
    })
    list.unshift({
      group_id: 0,
      group_name: window.$t('all'),
    })
    categoryList.value = list
  })
}

const onGroupChange = e => {
  if (e.value > 0) {
    scrollTo(`#toolbox-group-${e.value}`)
  } else {
    scrollToTop()
  }
}

const handleAdd = (data?: AI_LINK.State) => {
  emit('add', { data })
}

const handleVisit = (link: AI_LINK.State) => {
  window.open(link.url, '_blank')
}
onMounted(() => {
  loadList()
})
</script>

<style scoped></style>
