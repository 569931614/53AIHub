<template>
  <div class="border rounded-lg p-5 pb-8 group">
    <!-- 支付方式头部 -->
    <div class="relative w-full flex items-center gap-3">
      <!-- 支付图标 -->
      <SvgIcon :name="iconName" width="24" />

      <!-- 支付方式名称 -->
      <label class="font-semibold text-[#1D1E1F]">{{ paymentLabel }}</label>

      <!-- 启用状态标签 -->
      <ElTag
        v-if="settingInfo.pay_status"
        class="!border-none !bg-[#E3F6E0] !text-[#09BB07]"
        type="success"
        size="default"
      >
        {{ $t('enabled') }}
      </ElTag>

      <div class="flex-1" />

      <!-- 操作下拉菜单 -->
      <ElDropdown placement="bottom" @command="handleCommand">
        <div
          class="!border-none !outline-none p-1 cursor-pointer rounded overflow-hidden invisible group-hover:visible hover:bg-[#F0F0F0]"
        >
          <ElIcon class="rotate-90" size="16">
            <MoreFilled />
          </ElIcon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <!-- 设置选项 -->
            <el-dropdown-item :command="PAYMENT_COMMAND.SETTING">
              {{ $t('action_setting') }}
            </el-dropdown-item>

            <!-- 启用/禁用选项（仅当已配置时显示） -->
            <template v-if="settingInfo.pay_setting_id">
              <el-dropdown-item v-if="settingInfo.pay_status" :command="PAYMENT_COMMAND.DISABLE">
                {{ $t('action_disable') }}
              </el-dropdown-item>
              <el-dropdown-item v-else :command="PAYMENT_COMMAND.ENABLE">
                {{ $t('action_enable') }}
              </el-dropdown-item>
            </template>
          </el-dropdown-menu>
        </template>
      </ElDropdown>
    </div>

    <!-- 配置状态信息 -->
    <div class="mt-3 text-sm text-[#4F5052]">
      <template v-if="settingInfo.pay_setting_id">
        {{ $t('setting') }} · {{ $t('updated_at') }} {{ formatUpdateTime }}
      </template>
      <template v-else>
        {{ $t('not_setting') }}
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { MoreFilled } from '@element-plus/icons-vue'
import { computed } from 'vue'
import { PAYMENT_TYPE_ICON_MAP, PAYMENT_TYPE_LABEL_MAP, PAYMENT_COMMAND } from '@/constants/payment'
import type { PaymentType } from '@/constants/payment'
import type { PaymentSetting } from '@/api/modules/payment/types'

// 组件属性
interface Props {
  settingInfo: Partial<PaymentSetting>
  type: string
}

// 组件事件
interface Emits {
  (e: 'command', command: string, type: string): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

/**
 * 获取支付方式图标名称
 */
const iconName = computed(() => {
  return PAYMENT_TYPE_ICON_MAP.get(props.settingInfo.pay_type as any) || 'default'
})

/**
 * 安全获取支付方式标签
 */
const getPaymentLabel = (payType: PaymentType) => {
  const label = PAYMENT_TYPE_LABEL_MAP.get(payType)
  return label || ''
}

/**
 * 获取支付方式标签
 */
const paymentLabel = computed(() => {
  return getPaymentLabel(props.settingInfo.pay_type as PaymentType)
})

/**
 * 格式化更新时间
 */
const formatUpdateTime = computed(() => {
  return props.settingInfo.updated_time?.slice(0, 16) || ''
})

/**
 * 处理下拉菜单命令
 * @param command 命令
 */
const handleCommand = (command: string) => {
  emit('command', command, props.type)
}
</script>

<style scoped lang="scss">
// 组件样式
</style>
