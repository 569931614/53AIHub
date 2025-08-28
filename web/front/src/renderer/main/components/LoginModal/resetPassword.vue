<template>
  <div v-if="!isOpLocalEnv" class="mb-2">
    <h3>{{ $t('form.reset_password_method') }}</h3>
    <el-radio-group v-model="verify_way" @change="handleVerifyWayChange">
      <el-radio value="email_verify" size="large" :disabled="!userStore.info.email">{{ $t('form.email_verify') }}</el-radio>
      <el-radio value="mobile_verify" size="large" :disabled="!userStore.info.mobile">{{ $t('form.mobile_verify') }}</el-radio>
    </el-radio-group>
  </div>

  <el-form ref="formRef" label-position="top" :model="form" :rules="rules" @keyup.enter="handleSubmit">
    <el-form-item :label="$t('form.verify_code')" prop="verify_code" :rules="[verify_way === 'email_verify' ? emailCodeRule : codeRule]">
      <div class="flex items-center" style="width: 100%">
        <el-input
          v-model="form.verify_code"
          v-trim
          size="large"
          class="input_style min-w-80 no-right-radius flex-1"
          :placeholder="$t('form.input_placeholder') + $t('form.verify_code')"
        ></el-input>
        <el-button
          v-if="verify_way === 'email_verify'"
          v-debounce
          :disabled="isSending"
          class="!bg-[#f5f5f5] border-0 h-[44px] w-29 no-left-radius"
          @click.stop="handleGetCode"
        >
          <div :class="['text-sm', 'pl-5', 'border-l', 'pr-1', 'text-[#2563EB]', { 'text-[#9A9A9A]': isSending }]">
            {{ emailCodeCount ? `${emailCodeCount}s` : $t('form.get_verify_code') }}
          </div>
        </el-button>
        <el-button v-else v-debounce :disabled="isSending" class="!bg-[#f5f5f5] border-0 h-[44px] w-29 no-left-radius" @click.stop="handleGetCode">
          <div :class="['text-sm', 'pl-5', 'border-l', 'pr-1', 'text-[#2563EB]', { 'text-[#9A9A9A]': isSending }]">
            {{ codeCount ? `${codeCount}s` : $t('form.get_verify_code') }}
          </div>
        </el-button>
      </div>
    </el-form-item>
    <el-form-item :label="$t('form.new_password')" prop="new_password" :rules="[getPasswordRules()]">
      <el-input
        v-model="form.new_password"
        v-trim
        show-password
        size="large"
        class="input_style"
        :placeholder="$t('form.new_password_placeholder')"
      ></el-input>
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
        class="input_style"
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
import { getMobileRules, getEmailRules, getPasswordRules, getConfirmPasswordRules } from '@/utils/form-rules'

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

const verify_way = ref('email_verify')
if (!userStore.info.email) {
  verify_way.value = 'mobile_verify'
}

const rules = computed(() => {
  return {
    username: [verify_way.value === 'email_verify' ? getEmailRules() : getMobileRules()],
    new_password: [getPasswordRules()],
    confirm_password: [getPasswordRules(), getConfirmPasswordRules(form, 'new_password')],
    verify_code: [verify_way.value === 'email_verify' ? emailCodeRule : codeRule]
  }
})

const isSending = ref(true)
const handleGetCode = () => {
  if (verify_way.value === 'email_verify') {
    sendEmailCode(userStore.info.email)
  } else {
    sendcode(userStore.info.mobile)
  }
  if (codeCount || emailCodeCount) {
    isSending.value = true
  } else {
    isSending.value = false
  }
}

const handleSubmit = () => {
  return formRef.value?.validate().then(async (valid) => {
    if (!valid) return
    console.log('valid')
    try {
      if (verify_way.value === 'email_verify') {
        await userStore.reset_password({
          email: userStore.info.email,
          verify_code: form.verify_code,
          new_password: form.new_password,
          confirm_password: form.confirm_password
        })
      } else {
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
      ElMessage.success(window.$t('status.update_success'))
      emits('success')
      // 重置表单内容
      resetForm()
    } catch (error) {
      ElMessage.error()
    }
  })
}

const handleVerifyWayChange = () => {
  resetForm()
}

const resetForm = () => {
  // 重置字段值
  form.verify_code = ''
  form.new_password = ''
  form.confirm_password = ''
  clearFormValidation()
}

// 添加清除表单验证的方法
const clearFormValidation = () => {
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

defineExpose({
  resetForm
})

watch(
  () => codeCount.value,
  (newVal) => {
    isSending.value = newVal > 0
  },
  {
    immediate: true
  }
)
watch(
  () => emailCodeCount.value,
  (newVal) => {
    isSending.value = newVal > 0
  },
  {
    immediate: true
  }
)
</script>

<style scoped></style>
