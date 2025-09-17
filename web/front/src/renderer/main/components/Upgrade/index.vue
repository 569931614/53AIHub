<template>
  <div>
    <ElDialog
      v-model="visible"
      class="max-w-[1200px] rounded-md !p-0 overflow-hidden w-[90%] max-md:mt-12 max-md:w-full"
      header-class="hidden"
      footer-class="hidden"
      :show-close="false"
      append-to-body
      center
      destroy-on-close
    >
      <ElContainer v-loading="loading" class="h-[74vh] max-md:h-[90vh] md:!flex-row">
        <!-- <ElScrollbar class="w-full"> -->
        <ElMain class="!flex flex-col gap-4 relative !overflow-scroll !py-8 !pl-12 max-md:!pl-5 !pr-2">
          <ElIcon class="cursor-pointer absolute top-6 right-5 text-[#666] md:hidden" :size="20">
            <Close @click.stop="close" />
          </ElIcon>
          <h1 class="flex-none text-xl font-semibold text-[#000]">
            {{ $t('subscription.version_title') }}
          </h1>
          <div class="flex-1 relative overflow-y-hidden">
            <ElButton
              v-show="scroll_left > 0"
              class="absolute left-0 top-1/2 -translate-y-1/2 z-10 bg-[#53585f] text-[#fff] !border-none !outline-none"
              :icon="ArrowLeft"
              size="default"
              circle
              @click.stop="handleScrollLeft"
            />
            <ElScrollbar
              ref="version_scrollbar_ref"
              class="relative version-scrollbar"
              view-class="relative flex flex-nowrap gap-4 transition-all duration-300"
            >
              <div
                v-for="(item, item_index) in subscription_options"
                :key="item_index"
                class="flex-none rounded-md px-3.5 pt-4 pb-10 box-border border cursor-pointer min-w-[180px] h-max min-h-full max-md:(max-h-[430px] overflow-y-auto)"
                :class="[
                  subscription_options.length > 2 ? 'w-[31%] max-w-[240px]' : subscription_options.length > 1 ? 'w-[47%] max-w-[340px]' : 'w-full',
                  item.group_id == active_group_id
                    ? `shadow-[0_0_20px_rgba(6, 114, 255, 0.2)] border-[#2F74FF] bg-cover bg-center`
                    : 'border-[#EAEDF7] bg-[#F5F6FC]'
                ]"
                :style="{
                  backgroundImage: item.group_id == active_group_id ? `url(${$getPublicPath('/images/subscription/version-bg.png')})` : 'none'
                }"
                @click.stop="handleVersionSelect({ data: item })"
              >
                <header class="w-full flex items-center gap-3">
                  <img
                    :src="!/\.png$/.test(item.logo) ? $getPublicPath(`/images/subscription/${item.logo}.png`) : item.logo"
                    class="flex-none w-12 h-12 rounded-full overflow-auto"
                  />
                  <h2
                    class="text-lg font-semibold text-[#000] truncate"
                    :class="subscription_options.length <= 1 ? 'flex-none' : 'flex-1'"
                    :title="item.group_name"
                  >
                    {{ item.group_name }}
                  </h2>
                  <div v-if="subscription_options.length <= 1" class="flex-1" />
                  <div v-if="subscription_options.length <= 1">
                    <div class="flex items-center gap-1">
                      <span class="text-lg font-bold text-[#000]">
                        {{ item.month_info.currency_symbol }}
                        {{ item.month_info.amount }}
                      </span>
                      <span class="text-xs text-[#333]">/ {{ $t(`subscription.${item.month_info.time_unit}`) }}</span>
                    </div>
                    <div class="text-xs text-[#9A9A9A]">
                      {{
                        $t(`subscription.credit_month_amount`, {
                          amount: ` ${item.credit_month_info.amount} `
                        })
                      }}
                    </div>
                  </div>
                </header>
                <ElDivider class="!my-4 !border-[#E7ECF7]" />
                <template v-if="subscription_options.length > 1">
                  <div class="flex items-center gap-1">
                    <span class="text-lg font-bold text-[#000]">
                      {{ item.month_info.currency_symbol }}
                      {{ item.month_info.amount }}
                    </span>
                    <span class="text-xs text-[#333]">/ {{ $t(`subscription.${item.month_info.time_unit}`) }}</span>
                  </div>
                  <div class="text-xs text-[#9A9A9A]">
                    {{
                      $t(`subscription.credit_month_amount`, {
                        amount: ` ${item.credit_month_info.amount} `
                      })
                    }}
                  </div>
                </template>
                <div :class="[subscription_options.length > 1 ? '' : 'flex flex-row']">
                  <div v-if="item.agents && item.agents.length" :class="[subscription_options.length > 1 ? 'mt-6' : ' flex-1']">
                    <h2 class="text-sm font-semibold text-[#1D1E1F]">
                      {{ $t(`subscription.agent_bots_title`) }}
                    </h2>
                    <ul class="flex flex-wrap gap-3.5 mt-4">
                      <li v-for="(row, row_index) in item.agents" :key="row_index" class="flex items-center gap-2 w-full" :title="row.name">
                        <img :src="row.logo || ''" :alt="row.name" class="flex-none w-4 h-4 rounded-full overflow-hidden" />
                        <div class="flex-1 text-sm text-[#4F5052] truncate">
                          {{ row.name }}
                        </div>
                      </li>
                    </ul>
                  </div>
                  <div :class="[subscription_options.length > 1 ? 'mt-5' : 'flex-1']">
                    <h2 class="text-sm font-semibold text-[#1D1E1F]">
                      {{ $t(`subscription.ai_assistant_title`) }}
                    </h2>
                    <div class="mt-4 flex flex-wrap items-center gap-2">
                      <ElButton
                        class="!p-1.5"
                        type="default"
                        plain
                        size="default"
                        :disabled="!item.ai_enabled"
                        @click.stop="handleAIAssistantOpen({ type: 'windows' })"
                      >
                        <img
                          :src="$getPublicPath(`/images/windows.png`)"
                          class="w-4 h-4 object-cover mr-1"
                          :class="{ 'opacity-50': !item.ai_enabled }"
                        />
                        <span class="text-xs">Windows</span>
                      </ElButton>
                      <ElButton
                        class="!p-1.5 !ml-0"
                        type="default"
                        plain
                        size="default"
                        :disabled="!item.ai_enabled"
                        @click.stop="handleAIAssistantOpen({ type: 'ios' })"
                      >
                        <img :src="$getPublicPath(`/images/ios.png`)" class="w-4 h-4 object-cover mr-1" :class="{ 'opacity-50': !item.ai_enabled }" />
                        <span class="text-xs">macOS</span>
                      </ElButton>
                      <ElButton
                        class="!p-1.5 !ml-0"
                        type="default"
                        plain
                        size="default"
                        :disabled="!item.ai_enabled"
                        @click.stop="handleAIAssistantOpen({ type: 'chrome' })"
                      >
                        <img
                          :src="$getPublicPath(`/images/chrome.png`)"
                          class="w-4 h-4 object-cover mr-1"
                          :class="{ 'opacity-50': !item.ai_enabled }"
                        />
                        <span class="text-xs">Google</span>
                      </ElButton>
                    </div>
                  </div>
                </div>
              </div>
            </ElScrollbar>
            <ElButton
              v-show="scroll_left_limit > 0 && scroll_left < scroll_left_limit"
              class="absolute right-0 top-1/2 -translate-y-1/2 z-10 bg-[#53585f] text-[#fff] !border-none !outline-none"
              :icon="ArrowRight"
              size="default"
              circle
              @click.stop="handleScrollRight"
            />
          </div>
          <template v-if="+(active_subscription_info.month_info || {}).amount || +(active_subscription_info.year_info || {}).amount">
            <h1 class="flex-none text-xl font-semibold text-[#000]">
              {{ $t('subscription.time_title') }}
            </h1>
            <ul class="flex-none w-full flex items-center gap-5 overflow-x-auto overflow-y-hidden">
              <template v-for="key in ['month', 'year']" :key="key">
                <li
                  v-if="+(active_subscription_info[`${key}_info`] || {}).amount"
                  class="h-[128px] w-[36%] max-w-[280px] p-5 box-border border rounded-md min-w-[150px]"
                  :class="[
                    key == active_time_unit
                      ? `shadow-[0_0_20px_rgba(6, 114, 255, 0.2)] border-[#2F74FF] bg-[url('https://img.js.design/assets/img/67e3c887985b9b806ee10e09.png#681a27dfb07db7652079a5900c960b9e')] bg-cover bg-center`
                      : 'border-[#EAEDF7] bg-[#F5F6FC]'
                  ]"
                  @click.stop="handleTimeUnitSelect({ unit: key })"
                >
                  <div class="text-sm text-[#000]">
                    {{ $t(`subscription.time_unit_${key}`) }}
                  </div>
                  <div class="mt-2 flex items-baseline">
                    <span class="text-lg font-bold text-[#000]">
                      {{ (active_subscription_info[`${key}_info`] || {}).currency_symbol || '￥' }}
                    </span>
                    <span class="text-2xl font-bold text-[#000] mx-1">
                      {{ (active_subscription_info[`${key}_info`] || {}).amount || 0 }}
                    </span>
                    <span class="text-sm text-[#333]">/ {{ $t(`subscription.${key}`) }}</span>
                  </div>
                </li>
              </template>
            </ul>
          </template>
        </ElMain>
        <!-- </ElScrollbar> -->
        <ElAside
          class="flex flex-col flex-none max-w-[306px] relative bg-[#FAFAFA] border-l border-[#F2F0F0] !py-8 !pl-6 !pr-12 max-md:hidden w-[30%]"
        >
          <ElIcon class="cursor-pointer absolute top-6 right-5 text-[#666]" :size="20">
            <Close @click.stop="close" />
          </ElIcon>
          <h2 class="mt-10 text-2xl font-semibold text-[#000]">
            {{ $t('subscription.aside_title') }}
          </h2>
          <p class="mt-2 text-sm text-[#333]">{{ $t('subscription.aside_desc') }}</p>
          <div class="w-full flex items-center gap-2 mt-3">
            <img class="flex-none size-8 rounded-full overflow-hidden object-cover" :src="userStore.info.avatar" />
            <div class="flex-1 text-sm font-medium text-[#333]">
              {{ userStore.info.nickname }}
            </div>
          </div>
          <ElDivider class="!my-5 !border-[#E7ECF7]" />
          <div class="w-full flex-1 md:min-h-[108px]">
            <div class="w-full flex items-center justify-between text-lg text-[#000]">
              <span>{{ active_subscription_info.group_name || '- -' }}</span>
              <span>{{ active_time_info.currency_symbol || '￥' }}{{ active_time_info.amount || 0 }}</span>
            </div>
          </div>
          <template v-if="+active_time_info.amount">
            <ElDivider class="!my-5 !border-[#E7ECF7]" />
            <template v-if="payment_options.length > 1">
              <h2 class="text-lg text-[#000]">{{ $t('subscription.payment') }}</h2>
              <ElRadioGroup v-model="active_payment" class="mt-2" :disabled="pay_disabled">
                <ElRadio v-for="opt in payment_options" :key="opt.pay_type" :value="opt.pay_type">
                  <span class="text-[#333]">
                    {{ $t(`${opt.label}`) }}
                  </span>
                </ElRadio>
              </ElRadioGroup>
            </template>
            <div class="text-[#000] text-xl mt-8 flex items-center justify-between gap-2">
              <span>
                {{ $t('subscription.total') }}
              </span>
              <span>
                {{ active_time_info.currency_symbol || '￥' }}
                {{ Number(active_time_info.amount || 0).toFixed(2) }}
              </span>
            </div>
            <ElButton
              class="w-full mt-4"
              type="primary"
              size="large"
              :disabled="pay_disabled"
              @click.stop="handleQrcodeOpen({ defaultDisabled: true })"
            >
              {{ $t('action.pay') }}
            </ElButton>
          </template>
        </ElAside>
        <ElFooter
          v-if="+active_time_info.amount"
          class="flex flex-col flex-none relative bg-[#FAFAFA] border-l border-[#F2F0F0] !py-2 !pl-6 !pr-12 shadow shadow-slate-700 min-h-max md:hidden"
        >
          <template v-if="+active_time_info.amount">
            <div class="min-h-1 pb-1" :class="{ 'pb-0': display_payment_detail }">
              <div v-if="display_payment_detail">
                <h2 class="mt-10 text-2xl font-semibold text-[#000]">
                  {{ $t('subscription.aside_title') }}
                </h2>
                <p class="mt-2 text-sm text-[#333]">{{ $t('subscription.aside_desc') }}</p>
                <div class="w-full flex items-center gap-2 mt-3">
                  <img class="flex-none size-8 rounded-full overflow-hidden object-cover" :src="userStore.info.avatar" />
                  <div class="flex-1 text-sm font-medium text-[#333]">
                    {{ userStore.info.nickname }}
                  </div>
                </div>
                <ElDivider class="!my-5 !border-[#E7ECF7]" />
                <div class="w-full flex-1 md:min-h-[108px] min-h-[48px] max-md:border-b">
                  <div class="w-full flex items-center justify-between text-lg text-[#000]">
                    <span>{{ active_subscription_info.group_name || '- -' }}</span>
                    <span>{{ active_time_info.currency_symbol || '￥' }}{{ active_time_info.amount || 0 }}</span>
                  </div>
                </div>
              </div>
              <div v-if="payment_options.length > 1" class="flex flex-row items-center gap-4 min-h-max md:pb-12">
                <h2 class="text-lg text-[#000] whitespace-nowrap">{{ $t('subscription.payment') }}</h2>
                <ElRadioGroup v-model="active_payment" class="" :disabled="pay_disabled">
                  <ElRadio v-for="opt in payment_options" :key="opt.pay_type" :value="opt.pay_type">
                    <span class="text-[#333]">
                      {{ $t(`${opt.label}`) }}
                    </span>
                  </ElRadio>
                </ElRadioGroup>
              </div>
              <div class="flex flex-row items-end gap-4">
                <div class="mt-4 flex gap-1 items-end">
                  <span class="text-lg text-[#000]">
                    {{ active_time_info.currency_symbol || '￥' }}
                  </span>
                  <span class="text-4xl font-semibold text-[#000]">
                    {{ Number(active_time_info.amount || 0).toFixed(2) }}
                  </span>
                </div>

                <div class="whitespace-nowrap min-w-max flex items-center gap-1" @click="handlePaymentDetailToggle">
                  <span class="text-sm text-[#333]">明细</span>
                  <ElIcon class="text-sm text-[#333]" :class="{ 'rotate-180': display_payment_detail }">
                    <ArrowUp />
                  </ElIcon>
                </div>

                <ElButton
                  class="w-full mt-4"
                  type="primary"
                  size="large"
                  :disabled="pay_disabled"
                  @click.stop="handleQrcodeOpen({ defaultDisabled: true })"
                >
                  {{ $t('action.pay') }}
                </ElButton>
              </div>
            </div>
          </template>
        </ElFooter>
      </ElContainer>
    </ElDialog>

    <PaymentQrcode ref="payment_qrcode_ref" @success="updateUserGroup({ refresh: true })" />
  </div>
</template>

<script setup lang="ts">
import { ArrowLeft, ArrowRight, ArrowUp, Close } from '@element-plus/icons-vue'
import { ref, computed, nextTick, provide } from 'vue'
import { ElScrollbar } from 'element-plus'
import PaymentQrcode from './payment-qrcode.vue'
import subscriptionApi from '@/api/modules/subscription'
import { PAYMENT_TYPE, PAYMENT_TYPE_LABEL_MAP, PaymentType } from '@/constants/payment'
import paymentApi from '@/api/modules/payment'
import { useUserStore } from '@/stores/modules/user'

const userStore = useUserStore()

const version_scrollbar_ref = ref()
const payment_qrcode_ref = ref()
const visible = ref(false)
const loading = ref(false)
const subscription_options = ref([])
const active_group_id = ref('')
const active_time_unit = ref('month')
const scroll_left = ref(0)
const scroll_left_limit = ref(0)
const scroll_left_distance = ref(0)
const active_payment = ref(PAYMENT_TYPE.WECHAT)
const payment_options = ref([])
const display_payment_detail = ref(false)

const active_subscription_info = computed(() => subscription_options.value.find((item) => item.group_id === active_group_id.value) || {})
const active_time_info = computed(() => active_subscription_info.value[`${active_time_unit.value}_info`] || {})
const pay_disabled = computed(() => !Number(active_time_info.value.amount || 0))

provide('active_subscription_info', active_subscription_info)
provide('active_time_info', active_time_info)
provide('active_payment', active_payment)

async function open() {
  visible.value = true
  loadPaymentSettingData()
  await loadSubscriptionData()
  active_group_id.value = (subscription_options.value[0] || {}).group_id
  await nextTick()
  const scrollbar_el = version_scrollbar_ref.value.$el
  const scrollbar_view_el = scrollbar_el.querySelector('.el-scrollbar__view')
  const scrollbar_view_children_el_list = scrollbar_view_el.children
  if (scrollbar_view_children_el_list.length > 3) {
    const child_el_width = scrollbar_view_children_el_list[0].getBoundingClientRect().width
    scroll_left_distance.value = child_el_width + 16
    scroll_left_limit.value = scroll_left_distance.value * (scrollbar_view_children_el_list.length - 3)
  }
}
function close() {
  visible.value = false
}
async function loadSubscriptionData() {
  const { access_token } = userStore.info
  if (!access_token) return Promise.resolve(subscription_options.value)
  if (subscription_options.value.length) return Promise.resolve(subscription_options.value)
  loading.value = true
  const { list = [] } = await subscriptionApi.list().finally(() => {
    loading.value = false
  })
  subscription_options.value = list
  updateUserGroup()
  return Promise.resolve(subscription_options.value)
}
async function handleVersionSelect({ data = {} } = {}) {
  active_group_id.value = data.group_id || active_group_id.value
  await nextTick()
  active_time_unit.value = 'month'
  if (!+(active_subscription_info.value.month_info || {}).amount && +(active_subscription_info.value.year_info || {}).amount) {
    active_time_unit.value = 'year'
  }
}
function handleTimeUnitSelect({ unit = 'month' }) {
  active_time_unit.value = unit
}
function handleAIAssistantOpen() {
  // window.open('#', '_blank')
}
function handleScrollLeft() {
  for (let i = 0; i < scroll_left_distance.value; i += 6) {
    setTimeout(() => {
      scroll_left.value -= 6
      version_scrollbar_ref.value.setScrollLeft(scroll_left.value)
    }, 1)
  }
}
function handleScrollRight() {
  for (let i = 0; i < scroll_left_distance.value; i += 6) {
    setTimeout(() => {
      scroll_left.value += 6
      version_scrollbar_ref.value.setScrollLeft(scroll_left.value)
    }, 1)
  }
}
async function loadPaymentSettingData({ defaultDisabled = false } = {}) {
  const list = await paymentApi.getAvailableList()
  payment_options.value = list
    .filter((item = {}) => [PAYMENT_TYPE.WECHAT, PAYMENT_TYPE.ALIPAY, PAYMENT_TYPE.MANUAL].includes(item.pay_type) && item.enabled && item.configured)
    .map((item) => {
      return {
        ...item,
        label: PAYMENT_TYPE_LABEL_MAP[item.pay_type as PaymentType]
      }
    })
  if (!defaultDisabled) active_payment.value = (payment_options.value[0] || {}).pay_type || PAYMENT_TYPE.WECHAT
}
async function validateUpgrade() {
  await loadSubscriptionData()
  if (!subscription_options.value.length) return false
  // 逻辑调整, 支付配置放在支付按钮判断
  // await loadPaymentSettingData()
  // if (!payment_options.value.length) return false
  return true
}
async function updateUserGroup({ refresh = false } = {}) {
  if (refresh) await userStore.getUserInfo()
  const subscription_data = subscription_options.value.find((item = {}) => item.group_id === userStore.info.group_id)
  if (subscription_data) {
    userStore.setGroupName(subscription_data.group_name || userStore.info.group_name)
    userStore.setGroupIcon(subscription_data.logo_url || userStore.info.group_icon)
  }
}

function handlePaymentDetailToggle() {
  display_payment_detail.value = !display_payment_detail.value
}
async function handleQrcodeOpen({ defaultDisabled = false } = {}) {
  await loadPaymentSettingData({ defaultDisabled })
  if (!payment_options.value.length) {
    ElMessage.warning(window.$t('authority.payment_not_setting'))
  } else {
    payment_qrcode_ref.value.open()
  }
}

defineExpose({
  open,
  close,
  validateUpgrade
})
</script>

<style scoped>
.version-scrollbar::after {
  content: '';
  width: 100%;
  height: 2.8rem;
  position: absolute;
  left: 0;
  bottom: 0;
  background: linear-gradient(to top, rgb(255 255 255 / 80%), transparent);
}

.version-scrollbar ::v-deep(.el-scrollbar__bar.is-horizontal) {
  display: none;
}
</style>
