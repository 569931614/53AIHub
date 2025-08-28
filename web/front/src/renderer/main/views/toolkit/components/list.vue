<template>
  <div name="list" tag="div">
    <div v-if="!showList.length" class="col-span-full flex flex-col items-center justify-center">
      <el-empty :description="$t('common.no_data')" :image="$getPublicPath('/images/chat/completion_empty.png')" />
    </div>
    <template v-for="item in showList" :key="item.group_id">
      <h2 v-if="+item.group_id" :id="`group_${item.group_id}`" class="col-span-full text-placeholder">{{ item.group_name }}</h2>
      <div
        v-for="row in item.children"
        :key="row.id"
        class="min-h-20 bg-white rounded px-5 py-4 flex items-center gap-2 cursor-pointer border border-[#ECECEC] hover:shadow relative group"
        @click="handleAdd(row)"
      >
        <ElImage class="size-10 rounded-full" fit="contain" lazy :src="row.logo" />
        <div class="flex-1 overflow-hidden">
          <div
            class="text-base font-medium text-primary mb-1 mt-1 line-clamp-1"
            :title="row.name"
            v-html="row.name.replace(keyword, `<span class='text-theme'>${keyword}</span>`)"
          />
          <div
            class="text-sm text-regular text-opacity-60 line-clamp-1"
            :title="row.description"
            v-html="row.description.replace(keyword, `<span class='text-theme'>${keyword}</span>`)"
          />
        </div>

        <div
          v-if="row.has_share_account"
          class="absolute inset-0 items-center justify-center bg-[#222326] bg-opacity-55 rounded hidden group-hover:flex gap-2"
        >
          <ElButton class="!mr-0 hover:bg-white" @click.stop="handleVisit(row)">{{ $t('toolbox.account_access') }}</ElButton>
          <ElButton type="primary" class="!ml-0" @click.stop="handleAdd(row)">{{ $t('toolbox.direct_access') }}</ElButton>
        </div>
      </div>
    </template>
    <AccountDialog ref="dialogRef" />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useLinksStore } from '@/stores/modules/links'
import { useUserStore } from '@/stores/modules/user'
import AccountDialog from './account-dialog.vue'

const linksStore = useLinksStore()
const userStore = useUserStore()

const dialogRef = ref()

const props = withDefaults(
  defineProps<{
    list: Link.State[]
    keyword?: string
    onlyAll?: boolean
    groupId?: number
  }>(),
  {
    list: [],
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

const showList = computed(() => {
  // if (!props.keyword) return props.list
  // return props.list.filter(item => {
  //   return item.name.includes(props.keyword) || item.description.includes(props.keyword)
  // })
  const categorys = JSON.parse(JSON.stringify(linksStore.categorys || [])).filter((item) => (props.onlyAll ? item.group_id == 0 : item.group_id != 0))
  const list = JSON.parse(JSON.stringify(linksStore.links || []))
  categorys.forEach((item) => {
    item.children = list.filter((row) => {
      // 检查用户是否有权限访问
      const hasAccess = hasPermission(userStore.info.group_ids || [], row.user_group_ids || [])
      return hasAccess && row.group_id === item.group_id && row.user_group_ids.length > 0
    })
    if (props.onlyAll && item.group_id == 0) {
      item.children = props.list.filter((row) => {
        const hasAccess = hasPermission(userStore.info.group_ids || [], row.user_group_ids || [])
        return hasAccess && Array.isArray(row.user_group_ids) && row.user_group_ids.length > 0
      })
    }
    if (props.keyword) item.children = item.children.filter((row) => row.name.includes(props.keyword) || row.description.includes(props.keyword))
  })
  return categorys.filter((item) => item.children.length)
})

const handleAdd = (item: Link.State) => {
  if (item.has_share_account) return
  window.open(item.url, '_blank')
}

const handleVisit = (item: Link.State) => {
  dialogRef.value.open(item)
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
