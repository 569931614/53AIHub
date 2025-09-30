<template>
  <Layout class="px-[60px] py-8">
    <Header :title="$t('module.payment')" />

    <div class="flex-1 flex flex-col bg-white py-8 px-6 mt-3">
      <h1 class="font-semibold text-[#1D1E1F]">CNY</h1>
      <!-- CNY支付方式 -->
      <div class="mt-5 grid grid-cols-2 gap-5">
        <!-- 微信支付 -->
        <PaymentCard :setting-info="wechat_setting_info" type="wechat" @command="handleCommand" />

        <!-- 支付宝 -->
        <PaymentCard :setting-info="alipay_setting_info" type="alipay" @command="handleCommand" />

        <!-- 手动支付 -->
        <PaymentCard :setting-info="manual_setting_info" type="manual" @command="handleCommand" />
      </div>
      <!-- USD支付方式 -->
      <h1 class="mt-10 font-semibold text-[#1D1E1F] opacity-60">USD</h1>
      <div class="mt-5 grid grid-cols-2 gap-5 opacity-60">
        <!-- PayPal支付（暂未开放） -->
        <PaymentCard :setting-info="getPaymentSettingInfo('paypal')" type="paypal" @command="handleCommand" />

        <!-- 预留位置 -->
        <div class="flex-1 rounded-lg p-5 pb-8 group" />
      </div>
    </div>
  </Layout>

  <WechatSettingDialog ref="wechat_setting_ref" @success="refresh" />
  <ManualSettingDialog ref="manual_setting_ref" @success="refresh" />
  <AlipaySettingDialog ref="alipay_setting_ref" @success="refresh" />
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import WechatSettingDialog from './components/wechat-setting-dialog.vue'
import ManualSettingDialog from './components/manual-setting-dialog.vue'
import AlipaySettingDialog from './components/alipay-setting-dialog.vue'
import PaymentCard from './components/payment-card.vue'

// 使用新的支付API模块
import { paymentApi } from '@/api/modules/payment'
import { transformPaymentSettingList, transformToPaymentSettingMap } from '@/api/modules/payment/transform'
import { PAYMENT_COMMAND, PAYMENT_STATUS } from '@/constants/payment'
import type { PaymentSettingMap, PaymentSetting } from '@/api/modules/payment/types'
import TipConfirm from '@/components/TipConfirm/setup'
import { isInternalNetwork } from '@/utils'

// 对话框引用
const wechat_setting_ref = ref()
const manual_setting_ref = ref()
const alipay_setting_ref = ref()

// 支付设置数据
const paymentSettings = ref<PaymentSettingMap>({
  wechat: {} as PaymentSetting,
  alipay: {} as PaymentSetting,
  manual: {} as PaymentSetting,
  paypal: {} as PaymentSetting,
})

/**
 * 获取支付设置数据
 */
const refresh = async () => {
  try {
    const rawData = await paymentApi.getPaymentSettings()
    const settingsList = transformPaymentSettingList(rawData)
    paymentSettings.value = transformToPaymentSettingMap(settingsList)
  } catch (error) {
    console.error('获取支付设置失败:', error)
    ElMessage.error('获取支付设置失败')
  }
}

/**
 * 获取指定类型的支付设置信息
 * @param type 支付类型
 * @returns 支付设置信息
 */
const getPaymentSettingInfo = (type: keyof PaymentSettingMap) => {
  return paymentSettings.value[type] || ({} as PaymentSetting)
}

/**
 * 获取对话框引用
 * @param type 支付类型
 * @returns 对话框引用
 */
const getDialogRef = (type: keyof PaymentSettingMap) => {
  const dialogMap: Record<string, any> = {
    wechat: wechat_setting_ref,
    alipay: alipay_setting_ref,
    manual: manual_setting_ref,
  }
  return dialogMap[type]
}

/**
 * 更新支付状态
 * @param pay_setting_id 支付设置ID
 * @param pay_status 支付状态
 */
const updatePaymentStatus = async (pay_setting_id: number, pay_status: boolean) => {
  try {
    await paymentApi.updatePaymentStatus(pay_setting_id, { pay_status })
    const statusText = pay_status ? 'enabled' : 'disabled'
    ElMessage.success(window.$t(statusText))
    await refresh()
  } catch (error) {
    console.error('更新支付状态失败:', error)
    ElMessage.error('更新支付状态失败')
  }
}

/**
 * 处理支付操作命令
 * @param command 操作命令
 * @param type 支付类型
 */
const handleCommand = async (command: string, type: string = '') => {
  // PayPal功能暂未开放
  if (type === 'paypal') {
    ElMessage.warning(window.$t('feature_coming_soon'))
    return
  }

  // 内网环境限制（手动支付除外）
  if (isInternalNetwork() && type !== 'manual') {
    TipConfirm({
      title: window.$t('local_config_limited_tip'),
      content: window.$t('local_config_limited_desc', { url: window.location.href }),
      confirmButtonText: window.$t('know_it'),
      showCancelButton: false,
    }).open()
    return
  }

  const settingInfo = getPaymentSettingInfo(type as keyof PaymentSettingMap)

  switch (command) {
    case PAYMENT_COMMAND.SETTING:
      // 打开对应的设置对话框
      const dialogRef = getDialogRef(type as keyof PaymentSettingMap)
      if (dialogRef) {
        dialogRef.value.open({ data: settingInfo })
      }
      break

    case PAYMENT_COMMAND.ENABLE:
      await updatePaymentStatus(settingInfo.pay_setting_id, PAYMENT_STATUS.ENABLED)
      break

    case PAYMENT_COMMAND.DISABLE:
      await updatePaymentStatus(settingInfo.pay_setting_id, PAYMENT_STATUS.DISABLED)
      break
  }
}

// 计算属性：获取各支付类型的设置信息
const wechat_setting_info = computed(() => getPaymentSettingInfo('wechat'))
const alipay_setting_info = computed(() => getPaymentSettingInfo('alipay'))
const manual_setting_info = computed(() => getPaymentSettingInfo('manual'))

// 组件挂载时获取数据
onMounted(() => {
  refresh()
})
</script>

<style scoped lang="scss"></style>
