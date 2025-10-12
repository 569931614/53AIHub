<template>
  <div name="list" tag="div">
    <div v-if="!showList.length" class="col-span-full flex flex-col items-center justify-center">
      <el-empty :description="$t('common.no_data')" />
    </div>
    <template v-for="item in showList" :key="item">
      <h2 v-if="item.group_id !== null" :id="`group_${item.group_id}`" class="col-span-full text-placeholder">{{ item.group_name }}</h2>
      <div
        v-for="row in item.children"
        :key="row.id"
        class="min-h-20 bg-white rounded px-5 py-4 flex items-center gap-2 cursor-pointer border border-[#ECECEC] hover:shadow relative group"
        @click="handleCardClick(row)"
      >
        <ElImage class="size-10 rounded-full" fit="contain" lazy :src="row.logo" />
        <div class="flex-1 overflow-hidden">
          <div class="text-base font-medium text-primary mb-1 mt-1 line-clamp-1" :title="row.name" v-html="highlightText(row.name, keyword)" />
          <div class="text-sm text-regular text-opacity-60 line-clamp-1" :title="row.description" v-html="highlightText(row.description, keyword)" />
        </div>

        <!-- PC端悬停覆盖层 -->
        <div
          v-if="row.has_share_account"
          class="absolute inset-0 items-center justify-center bg-[#222326] bg-opacity-55 rounded hidden md:group-hover:flex gap-2"
        >
          <ElButton class="!mr-0 hover:bg-white" @click.stop="handleVisit(row)">{{ $t('toolbox.account_access') }}</ElButton>
          <ElButton type="primary" class="!ml-0" @click.stop="handleTo(row)">{{ $t('toolbox.direct_access') }}</ElButton>
        </div>
      </div>
    </template>

    <!-- 移动端底部弹窗 -->
    <div v-if="showMobileModal" class="fixed inset-0 z-50 md:hidden" @click="closeMobileModal">
      <!-- 背景遮罩 -->
      <div class="absolute inset-0 bg-black bg-opacity-50" />

      <!-- 底部弹窗 -->
      <div class="absolute bottom-0 left-0 right-0 bg-white rounded-t-2xl transform transition-transform duration-300 ease-out">
        <div class="flex flex-col">
          <button class="h-12 font-medium text-base hover:bg-gray-200 transition-colors" @click="handleTo(selectedItem!)">
            {{ $t('toolbox.direct_access') }}
          </button>
          <button class="h-12 font-medium text-base hover:bg-gray-200 transition-colors" @click="handleVisit(selectedItem!)">
            {{ $t('toolbox.account_access') }}
          </button>
          <button class="h-12 border-t font-medium text-base hover:bg-gray-50 transition-colors" @click="closeMobileModal">
            {{ $t('action.cancel') }}
          </button>
        </div>
      </div>
    </div>

    <AccountDialog ref="dialogRef" />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useLinksStore } from '@/stores/modules/links'
import { useUserStore } from '@/stores/modules/user'
import AccountDialog from './account-dialog.vue'
import { useBasicLayout } from '@/hooks/useBasicLayout'

const linksStore = useLinksStore()
const userStore = useUserStore()
const { isSmScreen } = useBasicLayout()

const dialogRef = ref()
const showMobileModal = ref(false)
const selectedItem = ref<Link.State | null>(null)

const props = withDefaults(
  defineProps<{
    list: Link.State[]
    keyword?: string
    onlyAll?: boolean
    groupId?: number
  }>(),
  {
    list: () => [],
    keyword: '',
    onlyAll: false,
    groupId: 0
  }
)

// 检查用户是否有权限访问
const hasPermission = (userGroupIds: number[], itemGroupIds: number[]) => {
  if (!itemGroupIds || itemGroupIds.length === 0) return false
  return userGroupIds.some((groupId) => itemGroupIds.includes(groupId))
}

// 高亮文本函数（不区分大小写）
const highlightText = (text: string, keyword: string) => {
  if (!text || !keyword) return text

  const regex = new RegExp(keyword, 'gi')
  return text.replace(regex, (match) => `<span class='text-theme'>${match}</span>`)
}

const showList = computed(() => {
  const categories = linksStore.categorys || []

  const links = linksStore.links || []

  let list = categories
    .map((category) => {
      const children = links.filter((link) => {
        const hasAccess = hasPermission(userStore.info.group_ids || [], (link as any).user_group_ids || [])
        return hasAccess && link.group_id === category.group_id && (link as any).user_group_ids?.length > 0
      })

      // 搜索过滤
      const filteredChildren = props.keyword
        ? children.filter((link) => {
            const keyword = props.keyword.toLowerCase()
            return link.name?.toLowerCase().includes(keyword) || link.description?.toLowerCase().includes(keyword)
          })
        : children

      return props.onlyAll ? { group_id: null, group_name: null, children: filteredChildren } : { ...category, children: filteredChildren }
    })
    .filter((category) => category.children.length > 0)
  // 首页最多展示6项
  if (props.onlyAll) {
    const children: any = []
    list.forEach((item: any) => children.push(...item.children))
    list = [
      {
        group_id: null,
        group_name: null,
        children: children.slice(0, 6)
      }
    ]
  }
  return list
})

const handleCardClick = (item: Link.State) => {
  if (isSmScreen.value && item.has_share_account) {
    selectedItem.value = item
    showMobileModal.value = true
  } else {
    window.open(item.url, '_blank')
  }
}

const handleTo = (item: Link.State) => {
  closeMobileModal()
  window.open(item.url, '_blank')
}

const handleVisit = (item: Link.State) => {
  closeMobileModal()
  dialogRef.value.open(item)
}

const closeMobileModal = () => {
  showMobileModal.value = false
  selectedItem.value = null
}
</script>

<style scoped>
.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateY(30px);
}
</style>
