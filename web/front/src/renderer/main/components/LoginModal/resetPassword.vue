<template>
  <div v-if="!isOpLocalEnv" class="mb-2">
    <h3>{{ $t('form.reset_password_method') }}</h3>
    <el-radio-group v-model="verify_way" @change="handleVerifyWayChange">
      <el-radio value="email_verify" size="large" :disabled="!userStore.info.email">{{ $t('form.email_verify') }}</el-radio>
      <el-radio value="mobile_verify" size="large" :disabled="!userStore.info.mobile">{{ $t('form.mobile_verify') }}</el-radio>
    </el-radio-group>
  </div>

  <el-form ref="formRef" label-position="top" :model="form" :rules="rules" @keyup.enter="handleSubmit">
    <el-form-item :label="$t('form.verify_code')" prop="verify_code" :rules="[getCodeRules()]">
      <div class="flex items-center" style="width: 100%">
        <el-input
          v-model="form.verify_code"
          v-trim
          size="large"
          class="min-w-80 no-right-radius flex-1"
          :placeholder="$t('form.input_placeholder') + $t('form.verify_code')"
        >
          <template #append>
            <el-button v-debounce :disabled="isSending" class="w-29 no-left-radius" @click.stop="handleGetCode">
              <div :class="['text-[#2563EB]', { 'text-[#9A9A9A]': isSending }]">
                {{ getCodeCount() ? `${getCodeCount()}s` : $t('form.get_verify_code') }}
              </div>
            </el-button>
          </template>
        </el-input>
      </div>
    </el-form-item>
    <el-form-item :label="$t('form.new_password')" prop="new_password" :rules="[getPasswordRules()]">
      <el-input v-model="form.new_password" v-trim show-password size="large" :placeholder="$t('form.new_password_placeholder')"></el-input>
    </el-form-item>

    <el-form-item
      :label="$t('form.new_password_confirm')"
      prop="confirm_password"
      :rules="[getPasswordRules(), getConfirmPasswordRules(form, 'new_password')]"
    >
      <el-input
        v-model="form.confirm_password"
        v-trim
        show-password
        size="large"
        :placeholder="$t('form.new_password_confirm_placeholder')"
      ></el-input>
    </el-form-item>

    <!-- 修改按钮 -->
    <el-button v-debounce type="primary" round class="w-full mt-3 !h-10" @click="handleSubmit">
      {{ $t('action.update_password') }}
    </el-button>
  </el-form>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import commonApi from '@/api/modules/common'
import useEnv from '@/hooks/useEnv'

import { useUserStore } from '@/stores/modules/user'
import useEmail from '@/hooks/useEmail'
import useMobile from '@/hooks/useMobile'
import { getPasswordRules, getConfirmPasswordRules } from '@/utils/form-rules'

const { isOpLocalEnv } = useEnv()

const emits = defineEmits(['success'])

const userStore = useUserStore()
const { emailCodeRule, sendEmailCode, emailCodeCount } = useEmail()
const { sendcode, codeRule, codeCount } = useMobile()

const formRef = ref<FormInstance>()

const form = reactive({
  verify_code: '',
  new_password: '',
  confirm_password: ''
})

const verify_way = ref(userStore.info.email ? 'email_verify' : 'mobile_verify')
const isSending = ref(true)

// 计算属性
const rules = computed(() => ({
  new_password: [getPasswordRules()],
  confirm_password: [getPasswordRules(), getConfirmPasswordRules(form, 'new_password')],
  verify_code: [getCodeRules()]
}))

// 工具函数
const getCodeRules = () => {
  return verify_way.value === 'email_verify' ? emailCodeRule : codeRule
}

const getCodeCount = () => {
  return verify_way.value === 'email_verify' ? emailCodeCount.value : codeCount.value
}
// 验证码发送
const handleGetCode = () => {
  const sendCodeFn = verify_way.value === 'email_verify' ? sendEmailCode : sendcode
  const target = verify_way.value === 'email_verify' ? userStore.info.email : userStore.info.mobile

  sendCodeFn(target)
  isSending.value = Boolean(getCodeCount())
}

// 提交表单
const handleSubmit = () => {
  return formRef.value?.validate().then(async (valid) => {
    if (!valid) return

    try {
      if (verify_way.value === 'email_verify') {
        await performEmailReset()
      } else {
        await performMobileReset()
      }

      ElMessage.success(window.$t('status.update_success'))
      emits('success')
      resetForm()
    } catch (error) {
      ElMessage.error()
    }
  })
}

// 邮箱重置密码
const performEmailReset = async () => {
  await userStore.reset_password({
    email: userStore.info.email,
    verify_code: form.verify_code,
    new_password: form.new_password,
    confirm_password: form.confirm_password
  })
}

// 手机号重置密码
const performMobileReset = async () => {
  await commonApi.verifycode({
    mobile: userStore.info.mobile,
    verifycode: form.verify_code,
    type: '1'
  })
  await userStore.reset_password({
    mobile: userStore.info.mobile,
    verify_code: form.verify_code,
    new_password: form.new_password,
    confirm_password: form.confirm_password
  })
}

// 验证方式切换
const handleVerifyWayChange = () => {
  resetForm()
}

// 重置表单
const resetForm = () => {
  Object.assign(form, {
    verify_code: '',
    new_password: '',
    confirm_password: ''
  })
  formRef.value?.clearValidate()
}

defineExpose({
  resetForm
})

// 监听验证码倒计时状态
watch(
  [() => codeCount.value, () => emailCodeCount.value],
  ([mobileCount, emailCount]) => {
    isSending.value = mobileCount > 0 || emailCount > 0
  },
  {
    immediate: true
  }
)
</script>

<style scoped></style>
