<template>
  <Layout class="px-[60px] py-8">
    <Header :title="$t('module.operation_order')" />
    <div class="flex-1 flex flex-col bg-white p-6 mt-3 box-border max-h-[calc(100vh-100px)] overflow-auto">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3">
          <ElSelect v-model="filter_form.status" class="flex-none max-w-[180px]" size="large" @change="refresh">
            <template #prefix> {{ $t('order_status') }}： </template>
            <ElOption v-for="opt in order_status_options" :key="opt.value" :label="$t(opt.label)" :value="opt.value" />
          </ElSelect>

          <ElSelect v-model="filter_form.pay_type" class="flex-none max-w-[180px]" size="large" @change="refresh">
            <template #prefix> {{ $t('order_payment_type') }}： </template>
            <ElOption
              v-for="opt in order_payment_type_options"
              :key="opt.value"
              :label="$t(opt.label)"
              :value="opt.value"
            />
          </ElSelect>
          <div class="flex-none">
            <FilterDateRange
              v-model="filter_form.date"
              size="large"
              :value-format="(date: Date) => date.getTime()"
              @change="refresh"
            />
          </div>
        </div>
        <div class="flex-center gap-3">
          <ElInput
            v-model="filter_form.keyword"
            style="max-width: 268px"
            size="large"
            clearable
            :suffix-icon="Search"
            :placeholder="$t('module.operation_order_search_placeholder')"
            @change="refresh"
          />
          <ElButton class="min-w-[96px]" type="primary" size="large" @click="handleAdd">
            {{ $t('action_add') }}
          </ElButton>
        </div>
      </div>

      <div v-loading="tableLoading" class="flex-1 overflow-y-auto bg-white rounded-lg mt-4">
        <TablePlus
          :data="tableData"
          :total="tableTotal"
          style="width: 100%"
          header-row-class-name="rounded overflow-hidden"
          header-cell-class-name="!bg-[#F6F7F8] !h-[60px] !border-none"
          @page-size-change="handleSizeChange"
          @page-current-change="handleCurrentChange"
        >
          <ElTableColumn :label="$t('order_id')" min-width="160" prop="order_id" show-overflow-tooltip />
          <ElTableColumn
            :label="$t('order_subscription')"
            min-width="140"
            prop="subscription_name"
            show-overflow-tooltip
          >
            <template #default="{ row }">
              <span :class="{ 'text-[#9B9B9B]': !row.subscription_name }">
                {{ row.subscription_name }}*{{ row.duration }}{{ $t(row.time_unit) }}
              </span>
            </template>
          </ElTableColumn>
          <ElTableColumn :label="$t('order_amount')" min-width="140" prop="amount" show-overflow-tooltip>
            <template #default="{ row }">
              <span :class="{ 'text-[#9B9B9B]': !row.amount }">
                {{ row.currency || 'CNY' }}&nbsp;{{ ((+row.amount || 0) / 100).toFixed(2) }}
              </span>
            </template>
          </ElTableColumn>
          <ElTableColumn :label="$t('order_status')" min-width="120" show-overflow-tooltip>
            <template #default="{ row }">
              <div class="flex items-center gap-1">
                <ElTag
                  class="!border-none"
                  :type="
                    row.status == ORDER_STATUS.CANCELLED
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
                  {{ $t(ORDER_STATUS_LABEL_MAP.get(row.status) || '') }}
                </ElTag>
              </div>
            </template>
          </ElTableColumn>
          <ElTableColumn :label="$t('order_create_time')" width="160" show-overflow-tooltip>
            <template #default="{ row }">
              <span :class="{ 'text-[#9B9B9B]': !row.created_time }">
                {{ row.created_time.slice(0, 16) || '--' }}
              </span>
            </template>
          </ElTableColumn>
          <ElTableColumn :label="$t('user')" width="160" prop="nickname" show-overflow-tooltip>
            <template #default="{ row }">
              <span :class="{ 'text-[#9B9B9B]': !row.nickname }">
                {{ row.nickname || '--' }}
              </span>
            </template>
          </ElTableColumn>
          <ElTableColumn :label="$t('order_payment_type')" prop="pay_type" width="120" show-overflow-tooltip>
            <template #default="{ row }">
              {{ getPaymentLabel(row.pay_type) }}
            </template>
          </ElTableColumn>
          <ElTableColumn :label="$t('operation')" width="170" fixed="right" align="right">
            <template #default="{ row }">
              <template v-if="PAYMENT_TYPE.MANUAL === row.pay_type && row.status === ORDER_STATUS.NOT_CONFIRM">
                <ElButton class="text-[#5A6D9E] !bg-transparent" type="text" @click.stop="handleConfirm({ data: row })">
                  {{ $t('action_confirm_payment') }}
                </ElButton>
                <ElButton class="text-[#5A6D9E] !bg-transparent" type="text" @click.stop="handleAdd({ data: row })">
                  {{ $t('action_edit') }}
                </ElButton>
                <ElButton class="text-[#5A6D9E] !bg-transparent" type="text" @click.stop="handleDelete({ data: row })">
                  {{ $t('action_delete') }}
                </ElButton>
              </template>
              <span v-else class="text-[#9B9B9B]"> -- </span>
            </template>
          </ElTableColumn>
        </TablePlus>
      </div>
    </div>
  </Layout>
  <OrderAddDialog ref="add_ref" @success="loadList" />
</template>

<script setup lang="ts">
import { Search } from '@element-plus/icons-vue'
import { onMounted, reactive, ref } from 'vue'
import FilterDateRange from '@/components/Filter/date-range.vue'
import OrderAddDialog from './components/order-add-dialog.vue'
import { ORDER_STATUS, ORDER_STATUS_LABEL_MAP } from '@/constants/order'

import { orderApi } from '@/api'
import { PAYMENT_TYPE, PAYMENT_TYPE_LABEL_MAP } from '@/constants/payment'
import type { PaymentType } from '@/constants/payment'

const add_ref = ref()
const filter_form = reactive({
  status: ORDER_STATUS.ALL,
  pay_type: PAYMENT_TYPE.ALL,
  keyword: '',
  offset: 0,
  limit: 10,
  date: [],
})
const order_status_options = [
  ORDER_STATUS.ALL,
  ORDER_STATUS.NOT_CONFIRM,
  ORDER_STATUS.PENDING,
  ORDER_STATUS.PAID,
  ORDER_STATUS.EXPIRED,
  ORDER_STATUS.CANCELLED,
].map(value => ({
  value,
  label: ORDER_STATUS_LABEL_MAP.get(value) || ORDER_STATUS_LABEL_MAP.get(ORDER_STATUS.ALL),
}))
/**
 * 安全获取支付方式标签
 */
const getPaymentLabel = (payType: PaymentType) => {
  const label = PAYMENT_TYPE_LABEL_MAP.get(payType)
  return label || ''
}

const order_payment_type_options = [
  PAYMENT_TYPE.ALL,
  PAYMENT_TYPE.WECHAT,
  PAYMENT_TYPE.ALIPAY,
  PAYMENT_TYPE.MANUAL,
].map(value => ({
  value,
  label: getPaymentLabel(value),
}))
const tableData = ref([])
const tableTotal = ref(0)
const tableLoading = ref(false)

const loadList = async () => {
  tableLoading.value = true
  const { total = 0, list = [] } = await orderApi
    .list({
      params: {
        ...filter_form,
        start_time: filter_form.date[0],
        end_time: filter_form.date[1],
      },
    })
    .finally(() => {
      tableLoading.value = false
    })
  tableTotal.value = total
  tableData.value = list
}

const refresh = () => {
  filter_form.offset = 0
  loadList()
}
const handleSizeChange = (size: number) => {
  filter_form.limit = size
  refresh()
}
const handleCurrentChange = (page: number) => {
  filter_form.offset = (page - 1) * filter_form.limit
  loadList()
}

const handleAdd = ({ data = {} } = {}) => {
  add_ref.value.open({ data })
}
const handleDelete = async ({ data = {} } = {}) => {
  await ElMessageBox.confirm(window.$t('module.operation_order_delete_confirm'), window.$t('tip'))
  await orderApi.delete_order({ id: data.id })
  ElMessage.success(window.$t('action_delete_success'))
  loadList()
}
const handleConfirm = async ({ data = {} } = {}) => {
  await ElMessageBox.confirm(window.$t('order.confirm_tip'), window.$t('tip'))
  await orderApi.confirm_order({ id: data.id })
  ElMessage.success(window.$t('action_confirm_success'))
  data.status = ORDER_STATUS.PAID
}

onMounted(async () => {
  refresh()
})
</script>

<style scoped></style>
