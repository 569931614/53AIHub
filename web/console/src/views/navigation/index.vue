<template>
  <Layout class="px-[60px] py-8">
    <Header :title="$t(route.meta?.title as string)" />
    <div class="flex-1 flex flex-col bg-white px-10 py-8 mt-3 box-border max-h-[calc(100vh-100px)] overflow-auto">
      <!-- 搜索区域 -->
      <div class="flex items-center justify-between">
        <ElInput
          v-model="searchKeyword"
          style="max-width: 268px"
          size="large"
          clearable
          :suffix-icon="Search"
          :placeholder="$t('navigation.search_placeholder')"
          @change="handleSearch"
        />
      </div>

      <!-- 表格区域 -->
      <div v-loading="isLoading" class="flex-1 overflow-y-auto bg-white rounded-lg mt-4">
        <TablePlus
          v-model:data="navigationList"
          :sortable="!searchKeyword"
          :total="totalCount"
          style="width: 100%"
          header-row-class-name="rounded overflow-hidden"
          header-cell-class-name="!bg-[#F6F7F8] !h-[60px] !border-none"
          :pagination="false"
          @page-size-change="handlePageSizeChange"
          @page-current-change="handlePageChange"
          @sortable-change="handleSortChange"
        >
          <!-- 名称列 -->
          <ElTableColumn :label="$t('name')" min-width="160" prop="name" show-overflow-tooltip>
            <template #default="scope">
              <div class="flex items-center">
                <ElImage
                  v-if="scope?.row?.icon"
                  class="size-[18px] mr-2 overflow-hidden"
                  :src="scope?.row?.icon"
                  fit="contain"
                />
                {{ scope?.row?.name }}
              </div>
            </template>
          </ElTableColumn>

          <!-- 类型列 -->
          <ElTableColumn :label="$t('type')" width="100" prop="type" show-overflow-tooltip>
            <template #default="{ row }">
              <span :class="{ 'text-[#9B9B9B]': !row.type_label }">
                {{ $t(row.type_label) || '--' }}
              </span>
            </template>
          </ElTableColumn>

          <!-- 跳转路径列 -->
          <ElTableColumn :label="$t('jump_path')" min-width="160" prop="jump_path" show-overflow-tooltip>
            <template #default="{ row }">
              <span :class="{ 'text-[#9B9B9B]': !row.jump_path }">
                {{ row.jump_path || '--' }}
              </span>
            </template>
          </ElTableColumn>

          <!-- 打开方式列 -->
          <ElTableColumn :label="$t('open_method')" width="100" prop="target" show-overflow-tooltip>
            <template #default="{ row }">
              <span :class="{ 'text-[#9B9B9B]': !row.target_label }">
                {{ $t(row.target_label) || '--' }}
              </span>
            </template>
          </ElTableColumn>

          <!-- 状态开关列 -->
          <ElTableColumn :label="$t('navigation_is_open')" width="140" prop="status" show-overflow-tooltip>
            <template #default="{ row }">
              <ElSwitch v-model="row.status" :active-value="1" :inactive-value="0" @change="handleStatusChange(row)" />
            </template>
          </ElTableColumn>

          <!-- 操作列 -->
          <ElTableColumn :label="$t('operation')" width="100" fixed="right" align="right">
            <template #default="{ row }">
              <div class="flex items-center justify-end gap-4">
                <!-- 页面编辑按钮 -->
                <ElTooltip v-if="row.type === NAVIGATION_TYPE.CUSTOM" placement="top" :content="$t('page_edit')">
                  <span>
                    <SvgIcon
                      class="cursor-pointer"
                      name="editor"
                      color="#5A6D9E"
                      width="16"
                      @click="handlePageEdit(row)"
                    />
                  </span>
                </ElTooltip>

                <!-- 设置按钮 -->
                <ElTooltip placement="top" :content="$t('action_setting')">
                  <span>
                    <SvgIcon
                      class="cursor-pointer"
                      name="web-setting"
                      color="#5A6D9E"
                      width="16"
                      @click="handleEdit(row)"
                    />
                  </span>
                </ElTooltip>

                <!-- 删除按钮 -->
                <ElTooltip placement="top" :content="$t('action_delete')">
                  <ElIcon
                    :class="[
                      row.type === NAVIGATION_TYPE.SYSTEM
                        ? 'text-[#BDC5D8] cursor-not-allowed'
                        : 'text-[#5A6D9E] cursor-pointer',
                    ]"
                    @click="handleDelete(row)"
                  >
                    <Delete />
                  </ElIcon>
                </ElTooltip>
              </div>
            </template>
          </ElTableColumn>
        </TablePlus>

        <!-- 添加按钮 -->
        <ElButton
          v-if="!searchKeyword"
          class="mt-4 !border-none"
          type="primary"
          size="large"
          plain
          :disabled="isAddDisabled"
          @click="handleAdd()"
        >
          + {{ $t('action_add') }}（{{ navigationList.length }}/{{ MAX_ITEMS }}）
        </ElButton>
      </div>
    </div>

    <NavCreateDrawer ref="navCreateRef" @success="loadNavigationData" />
  </Layout>
</template>

<script setup lang="ts">
import { Search, Delete } from '@element-plus/icons-vue'
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import NavCreateDrawer from './components/nav-create-drawer.vue'
import eventBus from '@/utils/event-bus'
import { navigationApi } from '@/api/modules/navigation/index'
import { transformNavigationList } from '@/api/modules/navigation/transform'
import { NAVIGATION_TYPE, NAVIGATION_CONSTANTS } from '@/constants/navigation'
import type {
  NavigationItem,
  RawNavigationItem,
  NavigationListResponse,
  UpdateNavigationSortData,
} from '@/api/modules/navigation/types'
import { sleep } from '@/utils'

// 常量定义
const MAX_ITEMS = NAVIGATION_CONSTANTS.MAX_ITEMS

// 路由和组件引用
const route = useRoute()
const router = useRouter()
const navCreateRef = ref()

// 响应式数据
const navigationList = ref<NavigationItem[]>([])
const totalCount = ref(0)
const isLoading = ref(false)
const searchKeyword = ref('')

// 分页参数
const paginationParams = ref({
  offset: 0,
  limit: 10,
})

// 计算属性
const isAddDisabled = computed(() => navigationList.value.length >= MAX_ITEMS)

// 数据加载方法
const loadNavigationData = async () => {
  isLoading.value = true
  try {
    const params = {
      ...paginationParams.value,
      keyword: searchKeyword.value,
    }

    const rawData = await navigationApi.list(params)

    // 转换数据
    const transformedList = transformNavigationList(rawData as RawNavigationItem[])
    const response: NavigationListResponse = {
      total: transformedList.length,
      list: transformedList.filter(item => item.name.includes(searchKeyword.value)),
    }

    // 如果没有数据，初始化默认数据
    if (!response.list.length) {
      await navigationApi.init()
      await sleep(1.5)
      await loadNavigationData()
      return
    }

    totalCount.value = response.total
    navigationList.value = response.list
  } catch (error) {
    console.error('加载导航数据失败:', error)
  } finally {
    isLoading.value = false
  }
}

// 搜索处理
const handleSearch = () => {
  paginationParams.value.offset = 0
  loadNavigationData()
}

// 分页处理
const handlePageSizeChange = (size: number) => {
  paginationParams.value.limit = size
  paginationParams.value.offset = 0
  loadNavigationData()
}

const handlePageChange = (page: number) => {
  paginationParams.value.offset = (page - 1) * paginationParams.value.limit
  loadNavigationData()
}

// 状态切换
const handleStatusChange = async (row: NavigationItem) => {
  try {
    await navigationApi.updateStatus({
      navigation_id: row.navigation_id,
      status: row.status ? 1 : 0,
    })
    ElMessage.success(window.$t(row.status ? 'action_enable_success' : 'action_disable_success'))
  } catch (error) {
    console.error('更新状态失败:', error)
    // 恢复原状态
    row.status = row.status ? 0 : 1
  }
}

// 添加导航
const handleAdd = (data: NavigationItem | null = null) => {
  navCreateRef.value.open({
    data: data || {},
    navigationList: navigationList.value,
  })
}

// 编辑导航
const handleEdit = (data: NavigationItem) => {
  handleAdd(data)
}

// 页面编辑
const handlePageEdit = (row: NavigationItem) => {
  router.push({
    name: 'NavigationWebSetting',
    params: {
      navigation_id: row.navigation_id,
    },
  })
}

// 删除导航
const handleDelete = async (row: NavigationItem) => {
  if (row.type === NAVIGATION_TYPE.SYSTEM) return

  try {
    await ElMessageBox.confirm(window.$t('navigation.delete_confirm'), window.$t('tip'))
    await navigationApi.delete(row.navigation_id)
    ElMessage.success(window.$t('action_delete_success'))
    loadNavigationData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除导航失败:', error)
    }
  }
}

// 排序处理
const handleSortChange = async ({ data = [] }: { data: NavigationItem[] }) => {
  try {
    const sortList: UpdateNavigationSortData = data.map((item, index) => ({
      navigation_id: item.navigation_id,
      sort: 9999 - index,
    }))
    await navigationApi.updateSort(sortList)
    ElMessage.success(window.$t('action_sort_success'))
  } catch (error) {
    console.error('更新排序失败:', error)
  }
}

// 生命周期
onMounted(() => {
  loadNavigationData()
  eventBus.on('user-login-success', loadNavigationData)
})

onUnmounted(() => {
  eventBus.off('user-login-success', loadNavigationData)
})
</script>

<style scoped lang="scss"></style>
