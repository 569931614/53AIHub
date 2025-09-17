<template>
  <div class="h-full bg-white">
    <MainHeader sticky hide-user :sider-button="isSoftware">
      <template #before_suffix>
        <div class="text-base text-primary font-bold line-clamp-1 max-md:flex-1 max-md:text-center">
          {{ $t('profile.order_info') }}
        </div>
      </template>
      <template #after_suffix>
        <ElIcon class="text-primary cursor-pointer" size="18" @click="$router.back()">
          <Close />
        </ElIcon>
      </template>
    </MainHeader>

    <div class="px-10 py-8">
      <div class="flex items-center gap-4">
        <FilterDateRange
          v-model="searchParams.date"
          class="max-w-[300px]"
          :value-format="(date) => date.getTime()"
          size="large"
          @change="handleSearch"
        />
        <el-select v-model="searchParams.subscription" class="max-w-[200px]" size="large" @change="handleSearch">
          <template #prefix>{{ $t('order.subscription') }}：</template>
          <el-option v-for="item in subscription_list" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
        <el-select v-model="searchParams.pay_type" class="max-w-[180px]" size="large" @change="handleSearch">
          <template #prefix>{{ $t('order.pay_type') }}：</template>
          <el-option v-for="item in order_pay_type_list" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
        <el-input
          v-model="searchParams.keyword"
          class="max-w-[160px]"
          :prefix-icon="Search"
          size="large"
          :placeholder="$t('order.search_placeholder')"
          clearable
          @change="handleSearch"
        />
      </div>
      <!-- 内容区域 -->
      <TablePlus
        :data="tableData"
        :total="tableTotal"
        style="width: 100%"
        class="mt-4"
        header-row-class-name="rounded overflow-hidden"
        header-cell-class-name="!bg-[#F6F7F8] !h-[60px] !border-none"
        @page-size-change="handleSizeChange"
        @page-current-change="handleCurrentChange"
      >
        <ElTableColumn :label="$t('order.id')" min-width="160" prop="order_id" show-overflow-tooltip />
        <ElTableColumn :label="$t('order.subscription')" min-width="140" prop="subscription_name" show-overflow-tooltip>
          <template #default="{ row }">
            <span :class="{ 'text-[#9B9B9B]': !row.subscription_name }">
              {{ row.subscription_name }}*{{ row.duration }}{{ $t('subscription.' + row.time_unit) }}
            </span>
          </template>
        </ElTableColumn>
        <ElTableColumn :label="$t('order.status')" min-width="120" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="flex items-center gap-1">
              <ElTag
                class="!border-none"
                :type="
                  row.status == ORDER_STATUS.CANCEL
                    ? 'info'
                    : row.status == ORDER_STATUS.NOT_CONFIRM
                      ? 'danger'
                      : row.status == ORDER_STATUS.PENDING
                        ? 'primary'
                        : row.status == ORDER_STATUS.EXPIRED
                          ? 'warning'
                          : 'success'
                "
              >
                {{ ORDER_STATUS_LABEL_MAP.get(row.status) }}
              </ElTag>
            </div>
          </template>
        </ElTableColumn>
        <ElTableColumn :label="$t('order.amount')" min-width="140" prop="amount" show-overflow-tooltip>
          <template #default="{ row }">
            <span :class="{ 'text-[#9B9B9B]': !row.amount }">{{ row.currency || 'CNY' }}&nbsp;{{ ((+row.amount || 0) / 100).toFixed(2) }}</span>
          </template>
        </ElTableColumn>
        <ElTableColumn :label="$t('order.pay_type')" prop="pay_type" width="120" show-overflow-tooltip>
          <template #default="{ row }">
            {{ getPayTypeLabel(row.pay_type) }}
          </template>
        </ElTableColumn>
        <ElTableColumn :label="$t('order.create_time')" width="160" show-overflow-tooltip>
          <template #default="{ row }">
            <span :class="{ 'text-[#9B9B9B]': !row.created_time }">
              {{ row.created_time.slice(0, 16) || '--' }}
            </span>
          </template>
        </ElTableColumn>
        <ElTableColumn :label="$t('order.action')" width="170" fixed="right">
          <template #default="{ row }">
            <template v-if="row.pay_type == PAY_TYPE.ALIPAY && row.status == ORDER_STATUS.PENDING">
              <ElButton v-debounce class="text-[#5A6D9E] !bg-transparent" type="text" @click.stop="handleAlipay(row)">
                {{ $t('order.payment') }}
              </ElButton>
              <ElButton v-debounce class="text-[#5A6D9E] !bg-transparent" type="text" @click.stop="handleCancel(row)">
                {{ $t('order.cancel') }}
              </ElButton>
            </template>
            <span v-else class="text-[#9B9B9B]">--</span>
          </template>
        </ElTableColumn>
      </TablePlus>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Search, Close } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/modules/user'
import MainHeader from '@/layout/header.vue'
import TablePlus from '@/components/TablePlus/index.vue'
import FilterDateRange from '@/components/Filter/date-range.vue'

import { subscriptionApi } from '@/api/modules/subscription'
import { ordersApi } from '@/api/modules/order'
import { getSimpleDateFormatString } from '@/utils/moment'

const ORDER_STATUS = {
  ALL: -1,
  NOT_CONFIRM: 1,
  PENDING: 2,
  PAID: 3,
  EXPIRED: 4,
  CANCEL: 5
}
const PAY_TYPE = {
  ALL: -1,
  WECHAT: 1,
  ALIPAY: 4,
  MANUAL: 2
}
const ORDER_STATUS_LABEL_MAP = new Map([
  [ORDER_STATUS.ALL, window.$t('order.status_list.all')],
  [ORDER_STATUS.NOT_CONFIRM, window.$t('order.status_list.not_confirm')],
  [ORDER_STATUS.PENDING, window.$t('order.status_list.pending')],
  [ORDER_STATUS.PAID, window.$t('order.status_list.paid')],
  [ORDER_STATUS.EXPIRED, window.$t('order.status_list.expired')],
  [ORDER_STATUS.CANCEL, window.$t('order.status_list.cancel')]
])
const order_pay_type_list = ref([
  { label: window.$t('common.all'), value: PAY_TYPE.ALL },
  { label: window.$t('subscription.wechat_pay'), value: PAY_TYPE.WECHAT },
  { label: window.$t('subscription.alipay'), value: PAY_TYPE.ALIPAY },
  { label: window.$t('subscription.manual_pay'), value: PAY_TYPE.MANUAL }
])

const subscription_list = computed(() => {
  return [{ label: window.$t('common.all'), value: 0 }].concat(
    userStore.subscriptions.map((item) => {
      return {
        label: item.group_name,
        value: item.group_id
      }
    })
  )
})

const userStore = useUserStore()

const tableData = ref([])
const tableTotal = ref(0)

const searchParams = ref({
  offset: 0,
  limit: 10,
  status: -1,
  pay_type: -1,
  keyword: '',
  date: [],
  subscription: 0
})

const getPayTypeLabel = (pay_type: number) => {
  return order_pay_type_list.value.find((item) => item.value === pay_type)?.label || ''
}

const loadList = async () => {
  const res = await ordersApi.list({
    ...searchParams.value,
    start_time: searchParams.value.date[0] ? new Date(searchParams.value.date[0]).getTime() : undefined,
    end_time: searchParams.value.date[1] ? new Date(searchParams.value.date[1]).getTime() : undefined
  })
  tableData.value = res.data.orders.map((item: Order.State) => ({
    ...item,
    created_time: getSimpleDateFormatString({
      date: item.created_time,
      format: 'YYYY-MM-DD hh:mm'
    })
  }))
  tableTotal.value = res.data.total
}

const handleSearch = () => {
  searchParams.value.offset = 0
  loadList()
}

const handleSizeChange = (size: number) => {
  searchParams.value.offset = 0
  searchParams.value.limit = size
  loadList()
}

const handleCurrentChange = (current: number) => {
  searchParams.value.offset = (current - 1) * searchParams.value.limit
  loadList()
}

const handleAlipay = (data: Order.State) => {
  return subscriptionApi.createOrder({
    params: {
      user_id: userStore.info.user_id,
      nickname: userStore.info.nickname,
      order_id: data.order_id,
      subscription_id: data.service_id,
      subscription_name: data.subscription_name,
      pay_type: data.pay_type,
      amount: Number(data.amount || 0),
      currency: data.currency,
      duration: 1,
      time_unit: data.time_unit
    },
    return_url: window.location.href
  })
}
const handleCancel = async (data: Order.State) => {
  await ElMessageBox.confirm(window.$t('order.cancel_confirm'))
  await ordersApi.close(data.order_id)
  loadList()
  ElMessage.success(window.$t('status.approve_cancel'))
}

onMounted(() => {
  loadList()
})
</script>
