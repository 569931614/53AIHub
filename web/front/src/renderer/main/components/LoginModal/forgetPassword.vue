<template>
  <div v-if="!isOpLocalEnv" class="mb-2">
    <h3>{{ $t('form.reset_password_method') }}</h3>
    <el-radio-group v-model="verify_way" @change="handleVerifyWayChange">
      <el-radio value="email_verify" size="large">{{ $t('form.email_verify') }}</el-radio>
      <el-radio value="mobile_verify" size="large">{{ $t('form.mobile_verify') }}</el-radio>
    </el-radio-group>
  </div>

  <el-form ref="formRef" label-position="top" :model="form" :rules="[]" @keyup.enter="handleSubmit">
    <el-form-item
      :label="$t(`form.${verify_way === 'email_verify' ? 'email' : 'mobile'}`)"
      prop="username"
      :rules="[getUsernameRules(), usernameCheck]"
    >
      <el-input
        v-model="form.username"
        v-trim
        size="large"
        :placeholder="$t('form.input_placeholder') + $t(`form.${verify_way === 'email_verify' ? 'email' : 'mobile'}`)"
        clearable
        @blur="onUsernameBlur"
      />
      <template #error>
        <div v-if="!existingAccount" class="text-xs text-[#f56c6c] absolute" style="top: 100%; left: 0">
          {{ $t('status.not_found_account') }}
          <button type="button" class="text-xs text-[#2563EB] underline" @click="handleClose">
            {{ $t('action.register') }}
          </button>
        </div>
      </template>
    </el-form-item>

    <el-form-item :label="$t('form.verify_code')" prop="verify_code" :rules="[getCodeRules()]">
      <div class="flex items-center" style="width: 100%">
        <el-input
          v-model="form.verify_code"
          v-trim
          size="large"
          class="no-right-radius flex-1"
          :placeholder="$t('form.input_placeholder') + $t('form.verify_code')"
        >
          <template #append>
            <el-button v-debounce :disabled="isRegister || isSending" class="!bg-[#f5f5f5] border-0 w-29 no-left-radius" @click.stop="handleGetCode">
              <div :class="['text-[#2563EB]', { 'text-[#9A9A9A]': isRegister || isSending }]">
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
import { ref, reactive, computed, nextTick, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import useEnv from '@/hooks/useEnv'
import commonApi from '@/api/modules/common'
import userApi from '@/api/modules/user'

import { useUserStore } from '@/stores/modules/user'
import useEmail from '@/hooks/useEmail'
import useMobile from '@/hooks/useMobile'
import { getMobileRules, getEmailRules, getPasswordRules, getConfirmPasswordRules } from '@/utils/form-rules'

const emits = defineEmits(['success', 'close'])

const userStore = useUserStore()
const { emailCodeRule, sendEmailCode, emailCodeCount } = useEmail()
const { sendcode, codeRule, codeCount } = useMobile()
const { isOpLocalEnv } = useEnv()

const formRef = ref<FormInstance>()

const form = reactive({
  username: '',
  verify_code: '',
  new_password: '',
  confirm_password: ''
})

const verify_way = ref('email_verify')
const isSending = ref(true)
const existingAccount = ref(true)
const usernameCache = reactive(new Map())
const isRegister = ref(true)

// 获取用户名验证规则
const getUsernameRules = () => {
  return verify_way.value === 'email_verify' ? getEmailRules() : getMobileRules()
}

// 获取验证码规则
const getCodeRules = () => {
  return verify_way.value === 'email_verify' ? emailCodeRule : codeRule
}

// 获取验证码倒计时
const getCodeCount = () => {
  return verify_way.value === 'email_verify' ? emailCodeCount.value : codeCount.value
}

const isFormatCorrect = computed(() => {
  const patterns = {
    mobile_verify: /^1[3-9]\d{9}$/,
    email_verify: /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
  }
  return patterns[verify_way.value].test(form.username)
})

// 重置表单
const resetForm = () => {
  Object.assign(form, {
    username: '',
    verify_code: '',
    new_password: '',
    confirm_password: ''
  })
  existingAccount.value = true
  isRegister.value = true
  nextTick(() => {
    formRef.value?.clearValidate()
  })
}

const handleVerifyWayChange = () => {
  resetForm()
}

// 用户名验证规则
const usernameCheck = {
  validator: async (_rule, _value, callback) => {
    try {
      if (form.username.trim() === '' || !isFormatCorrect.value) {
        return callback()
      }

      await onUsernameBlur()

      if (isRegister.value) {
        existingAccount.value = false
        return callback(new Error(window.$t(`form.${verify_way.value === 'email_verify' ? 'email' : 'mobile'}`) + window.$t('register.unregistered')))
      }
      return callback()
    } catch (error) {
      return callback()
    }
  },
  trigger: 'blur'
}

const handleGetCode = () => {
  const sendCodeFn = verify_way.value === 'email_verify' ? sendEmailCode : sendcode
  sendCodeFn(form.username)
}

const onUsernameBlur = async () => {
  if (!isFormatCorrect.value) return Promise.resolve()

  if (usernameCache.has(form.username)) {
    const cachedResult = usernameCache.get(form.username)
    if (Date.now() - cachedResult.timestamp < 2 * 60 * 1000) {
      isRegister.value = !cachedResult.exists
      return Promise.resolve()
    }
  }

  // 返回Promise确保外部可以await
  return userApi.checkUsername(form.username).then((res) => {
    isRegister.value = !res.data.exists
    usernameCache.set(form.username, {
      exists: res.data.exists,
      timestamp: Date.now()
    })
  })
}

const handleClose = () => {
  resetForm()
  emits('close')
}

const handleSubmit = () => {
  return formRef.value?.validate().then(async (valid) => {
    if (!valid) return

    try {
      const resetData = {
        verify_code: form.verify_code,
        new_password: form.new_password,
        confirm_password: form.confirm_password
      }

      if (verify_way.value === 'email_verify') {
        await userStore.reset_password({
          email: form.username,
          ...resetData
        })
      } else {
        await commonApi.verifycode({
          mobile: form.username,
          verifycode: form.verify_code,
          type: '1'
        })
        await userStore.reset_password({
          mobile: form.username,
          ...resetData
        })
      }

      ElMessage.success(window.$t('status.update_success'))
      emits('success')
      resetForm()
    } catch (error) {
      ElMessage.error()
    }
  })
}
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
