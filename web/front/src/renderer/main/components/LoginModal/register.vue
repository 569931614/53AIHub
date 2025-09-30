<template>
  <!-- 顶栏显示 -->
  <div class="flex justify-center mt-5">
    <template v-for="way in filteredRegisterWays" :key="way.value">
      <el-button
        :class="[
          'bg-transparent border-0 hover:bg-white hover:border-[#2563EB]',
          'px-0 rounded-none',
          way.value === REGISTER_WAY.email ? '!ml-7.5' : '',
          registerWay === way.value ? 'border-b-0.75 border-[#2563EB]' : 'border-b-0 border-[#1D1E1F]'
        ]"
        @click="handleRegisterWay(way.value)"
      >
        <h4 class="text-xl text-center mb-3" :class="[registerWay === way.value ? 'text-[#1D1E1F] font-bold' : 'text-[#94959B]']">
          {{ $t(`form.${way.value}`) + $t('action.register') }}
        </h4>
      </el-button>
    </template>
  </div>

  <el-form ref="formRef" label-position="top" :model="form" :rules="[]" class="px-2 mt-7" @keyup.enter="handleSubmit">
    <el-form-item :label="$t(`form.${registerWay}`)" prop="username" :rules="[getUsernameRules(), usernameCheck]">
      <el-input
        v-model="form.username"
        v-trim
        size="large"
        class="el-input--main"
        :placeholder="$t('form.input_placeholder') + $t(`form.${registerWay}`)"
        clearable
      />
      <template #error>
        <div v-if="existingAccount" class="text-xs text-[#f56c6c] absolute" style="top: 100%; left: 0">
          {{ $t(`form.existing_${registerWay}`) }}
          <button type="button" class="text-xs text-[#2563EB] underline" @click="handleClose">
            {{ $t('action.login') }}
          </button>
        </div>
      </template>
    </el-form-item>

    <el-form-item v-if="!isOpLocalEnv || (isOpLocalEnv && openSMTP)" :label="$t('form.verify_code')" prop="verify_code" :rules="[getCodeRules()]">
      <div class="flex items-center" style="width: 100%">
        <el-input
          v-model="form.verify_code"
          v-trim
          size="large"
          class="el-input--main w-80 no-right-radius flex-1"
          :placeholder="$t('form.input_placeholder') + $t('form.verify_code')"
        >
          <template #append>
            <el-button
              v-debounce
              :disabled="!isRegister || Boolean(getCodeCount())"
              class="!bg-[#f5f5f5] border-0 w-[100px] no-left-radius"
              @click.stop="handleGetCode"
            >
              <div :class="['text-[#2563EB]', { 'text-[#9A9A9A]': !isRegister || Boolean(getCodeCount()) }]">
                {{ getCodeCount() ? `${getCodeCount()}s` : $t('form.get_verify_code') }}
              </div>
            </el-button>
          </template>
        </el-input>
      </div>
    </el-form-item>

    <!-- 密码的输入框 -->
    <el-form-item :label="$t('form.password')" prop="password" :rules="[getPasswordRules()]">
      <el-input
        v-model="form.password"
        v-trim
        show-password
        size="large"
        class="el-input--main"
        :placeholder="$t('form.input_placeholder') + $t('form.password')"
      ></el-input>
    </el-form-item>
  </el-form>

  <!-- 已有账号立即登录 -->
  <div class="flex justify-end items-center">
    {{ $t('status.existing_account') }},
    <el-button link type="primary" @click="handleClose">
      {{ $t('action.login_directly') }}
    </el-button>
  </div>

  <!-- 注册按钮  -->
  <el-button v-debounce type="primary" round class="w-full mt-5 !h-10" @click="handleSubmit">
    {{ $t('action.register') }}
  </el-button>

  <!-- 底部协议 -->
  <div class="text-xs text-[#9A9A9A] text-center mt-5">
    {{ $t('register.agree') }}
    <a class="text-[#4F5052] cursor-pointer underline">{{ $t('register.terms_of_service') }}</a>
    {{ $t('action.and') }}
    <a class="text-[#4F5052] cursor-pointer underline">{{ $t('register.privacy_policy') }}</a>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, nextTick } from 'vue'
import type { FormInstance } from 'element-plus'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/modules/user'
import { getPasswordRules, getEmailRules, getMobileRules } from '@/utils/form-rules'

import useMobile from '@/hooks/useMobile'
import useEmail from '@/hooks/useEmail'
import useEnv from '@/hooks/useEnv'

import userApi from '@/api/modules/user'
import commonApi from '@/api/modules/common'

withDefaults(
  defineProps<{
    openSMTP: boolean
  }>(),
  {
    openSMTP: false
  }
)

const emits = defineEmits(['success', 'close'])

const REGISTER_WAY = {
  mobile: 'mobile',
  email: 'email'
} as const

const registerWays = [
  { value: REGISTER_WAY.mobile, label: 'mobile' },
  { value: REGISTER_WAY.email, label: 'email' }
]

const userStore = useUserStore()
const { isOpLocalEnv } = useEnv()
const { emailCodeRule, sendEmailCode, emailCodeCount } = useEmail()
const { sendcode, codeRule, codeCount } = useMobile()

const formRef = ref<FormInstance>()

const form = reactive({
  username: '',
  password: '',
  verify_code: ''
})

const registerWay = ref(isOpLocalEnv.value ? REGISTER_WAY.email : REGISTER_WAY.mobile)
const existingAccount = ref(false)
const usernameCache = reactive(new Map())
const isRegister = ref(false)

// 过滤注册方式
const filteredRegisterWays = computed(() => {
  return registerWays.filter((way) => !isOpLocalEnv.value || way.value === REGISTER_WAY.email)
})

const isFormatCorrect = computed(() => {
  const patterns = {
    [REGISTER_WAY.mobile]: /^1[3-9]\d{9}$/,
    [REGISTER_WAY.email]: /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
  }
  return patterns[registerWay.value].test(form.username)
})

// 获取用户名验证规则
const getUsernameRules = () => {
  if (registerWay.value === REGISTER_WAY.email) {
    return getEmailRules()
  }
  return getMobileRules()
}

// 获取验证码规则
const getCodeRules = () => {
  return registerWay.value === REGISTER_WAY.email ? emailCodeRule : codeRule
}

// 获取验证码倒计时
const getCodeCount = () => {
  return registerWay.value === REGISTER_WAY.email ? emailCodeCount.value : codeCount.value
}

const usernameCheck = {
  validator: async (_rule, _value, callback) => {
    try {
      if (form.username.trim() === '' || !isFormatCorrect.value) {
        return
      }

      await onUsernameBlur()

      if (!isRegister.value) {
        existingAccount.value = true
        callback(new Error(window.$t(`form.${registerWay.value}`) + window.$t('register.unregistered')))
      }
    } catch (error) {}
  },
  trigger: 'blur'
}

// 重置表单
const resetForm = () => {
  Object.assign(form, {
    username: '',
    verify_code: '',
    password: ''
  })
  existingAccount.value = false
  isRegister.value = false
  nextTick(() => {
    formRef.value?.clearValidate()
  })
}

const handleRegisterWay = (way) => {
  resetForm()
  registerWay.value = way
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

const handleGetCode = () => {
  const sendCodeFn = registerWay.value === REGISTER_WAY.email ? sendEmailCode : sendcode
  sendCodeFn(form.username)
}

const handleClose = () => {
  resetForm()
  emits('close')
}

const handleSubmit = () => {
  return formRef.value?.validate().then(async (valid) => {
    if (!valid) return

    try {
      // 手机号注册需要先验证验证码
      if (registerWay.value === REGISTER_WAY.mobile) {
        await commonApi.verifycode({
          mobile: form.username,
          verifycode: form.verify_code,
          type: '1'
        })
      }

      await userStore.register({
        username: form.username,
        password: form.password,
        verify_code: form.verify_code
      })

      ElMessage.success(window.$t('action.register') + window.$t('status.success'))
      emits('success')
    } catch (error) {}
  })
}
</script>
