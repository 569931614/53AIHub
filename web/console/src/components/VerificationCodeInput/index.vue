<template>
  <ElInput
    v-model="input_value"
    :style="{ '--el-input-bg-color': bgColor, '--el-input-border-color': 'transparent', '--el-input-height': height }"
    :clearable="clearable"
    :size="size"
    :disabled="disabled"
    :maxlength="maxlength"
    :placeholder="$t(placeholder)"
    @change="onChange"
  >
    <template #append>
      <ElButton class="!bg-transparent" type="primary" text :disabled="send_disabled" @click="onSend">
        <span :class="[send_disabled ? 'cursor-not-allowed' : 'cursor-pointer hover:opacity-70 text-[#3664EF]']">
          {{ send_countdown > 0 ? `${$t('action_send_success')}(${send_countdown}s)` : $t('get_verification_code') }}
        </span>
      </ElButton>
    </template>
  </ElInput>
</template>

<script setup lang="ts">
import { computed, onUnmounted, ref } from 'vue';
import api from '@/apis';
import { commonApi } from '@/api';

const props = withDefaults(
  defineProps<{
    account?: string
    accountType?: 'email' | 'mobile'
    modelValue?: string
    bgColor?: string
    height?: string
    size?: 'large' | 'default' | 'small'
    clearable?: boolean
    disabled?: boolean
    countdown?: number
    maxlength?: number
    placeholder?: string
  }>(),
  {
    account: '',
    accountType: 'mobile',
    modelValue: '',
    bgColor: '#F1F2F3',
    height: '44px',
    size: 'large',
    clearable: true,
    disabled: false,
    countdown: 60,
    maxlength: 4,
    placeholder: 'verification_code_placeholder',
  }
)
const emits = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const input_value = ref(props.modelValue)
const send_countdown = ref(0)

const send_disabled = computed(() => {
  return props.disabled || !props.account || send_countdown.value > 0
})
const real_account_type = computed(() => {
  return props.accountType || /^(13[0-9]|14[0-9]|15[0-9]|16[0-9]|17[0-9]|18[0-9]|19[0-9])\d{8}$/.test(props.account)
    ? 'mobile'
    : 'email'
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})

const onChange = (value: string) => {
  emits('update:modelValue', value)
}

// function isChinaMobile(phone: string): boolean {
//   // 清洗并处理手机号
//   let cleaned = phone.replace(/\D/g, '') // 去除非数字字符
//   if (cleaned.length > 11)
//     cleaned = cleaned.slice(-11) // 处理带国际区号的情况

//   // 基础验证
//   if (cleaned.length !== 11 || !cleaned.startsWith('1'))
//     return false

//   // 中国移动号段正则（2023年最新版）
//   const mobilePattern = /^1(3[5-9]|34[0-8]|440|4(7|8[0-9])|5[0-27-9]|7[28]|8[2-47-8]|9[58]|20)/
//   return mobilePattern.test(cleaned)
// }

let timer: NodeJS.Timeout | null = null
const onSend = async () => {
  if (!props.account) {
    ElMessage.warning(window.$t(`login.${real_account_type.value}_placeholder`))
    return
  }
  if (real_account_type.value !== 'mobile') {
    ElMessage.warning(window.$t('form_mobile_validator'))
    return
  }

  if (props.accountType === 'mobile') {
    await api.qyy.send_code({
      data: {
        mobile: props.account,
        source: 'companyibos',
      },
    })
  } else {
    await commonApi.sendEmailCode({ email: props.account })
  }
  send_countdown.value = props.countdown
  ElMessage.success(window.$t('action_send_success'))
  timer = setInterval(() => {
    send_countdown.value--
    if (send_countdown.value < 0) clearInterval(timer)
  }, 1000)
}
const reset = () => {
  if (timer) clearInterval(timer)
  input_value.value = ''
  send_countdown.value = 0
}
const validateCode = async () => {
  return commonApi
    .verifycode({
      mobile: props.account,
      verifycode: input_value.value,
      type: '1',
    })
    .catch(() => Promise.resolve(false))
}

defineExpose({
  reset,
  validateCode,
})
</script>

<style scoped lang="scss">
::v-deep(.el-input-group__append) {
  position: relative;
  background-color: var(--el-input-bg-color);
}
::v-deep(.el-input-group__append::before) {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  transform: translateY(-50%);
  display: block;
  width: 1px;
  height: 14px;
  background-color: #c7c7c7;
}
</style>
