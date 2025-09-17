<template>
  <!-- 顶栏显示 -->
  <div class="flex justify-center mt-5">
    <el-button
      v-if="!isOpLocalEnv"
      :class="[
        'bg-transparent border-0 hover:bg-white hover:border-[#2563EB]',
        'px-0',
        'rounded-none',
        register_way === 'mobile' ? 'border-b-0.75 border-[#2563EB]' : 'border-b-0 border-[#1D1E1F]'
      ]"
      @click="handleRegisterWay('mobile')"
    >
      <h4 class="text-xl text-center mb-3" :class="[register_way === 'mobile' ? 'text-[#1D1E1F] font-bold' : 'text-[#94959B]']">
        {{ $t('form.mobile') + $t('action.register') }}
      </h4>
    </el-button>

    <el-button
      :class="[
        'bg-transparent border-0 hover:bg-white hover:border-[#2563EB]',
        'px-0',
        'rounded-none',
        '!ml-7.5',
        register_way === 'email' ? 'border-b-0.75 border-[#2563EB]' : 'border-b-0 border-[#1D1E1F]'
      ]"
      @click="handleRegisterWay('email')"
    >
      <h4 class="text-xl text-center mb-3" :class="[register_way === 'email' ? 'text-[#1D1E1F] font-bold' : 'text-[#94959B]']">
        {{ $t('form.email') + $t('action.register') }}
      </h4>
    </el-button>
  </div>

  <el-form ref="formRef" label-position="top" :model="form" :rules="[]" class="px-2 mt-7" @keyup.enter="handleSubmit">
    <el-form-item
      :label="register_way === 'mobile' ? $t('form.mobile') : $t('form.email')"
      prop="username"
      :rules="[
        register_way === 'email' ? getEmailRules() : getMobileRules(),
        {
          validator: async (rule, value, callback) => {
            try {
              // 跳过空值触发和格式不符合触发
              if (form.username.trim() === '' || !isFormatCorrect) {
                return
              }
              // 等待验证完成
              await onUsernameBlur()

              // 检查完成后，根据isRegister的值进行验证
              if (register_way === 'mobile' && !isRegister) {
                existing_mobile = true
                callback(new Error($t('form.mobile') + $t('register.unregistered')))
              } else if (register_way === 'email' && !isRegister) {
                existing_email = true
                callback(new Error($t('form.email') + $t('register.unregistered')))
              }
            } catch (error) {}
          },
          trigger: 'blur'
        }
      ]"
    >
      <el-input
        v-if="register_way === 'mobile'"
        v-model="form.username"
        v-trim
        size="large"
        class="el-input--main"
        :placeholder="$t('form.input_placeholder') + $t('form.mobile')"
        clearable
        onblur="onUsernameBlur"
      />
      <el-input
        v-else
        v-model="form.username"
        v-trim
        size="large"
        class="el-input--main"
        :placeholder="$t('form.input_placeholder') + $t('form.email')"
        clearable
        onblur="onUsernameBlur"
      />
      <template #error>
        <div v-if="existing_mobile || existing_email" class="text-xs text-[#f56c6c] absolute" style="top: 100%; left: 0">
          {{ register_way === 'mobile' ? $t('form.existing_mobile') : $t('form.existing_email') }}
          <button type="button" class="text-xs text-[#2563EB] underline" @click="handleClose">
            {{ $t('action.login') }}
          </button>
        </div>
      </template>
    </el-form-item>

    <el-form-item
      v-if="!isOpLocalEnv || (isOpLocalEnv && openSMTP)"
      :label="$t('form.verify_code')"
      prop="verify_code"
      :rules="[register_way === 'email' ? emailCodeRule : codeRule]"
    >
      <div class="flex items-center" style="width: 100%">
        <el-input
          v-model="form.verify_code"
          v-trim
          size="large"
          class="el-input--main w-80 no-right-radius flex-1"
          :placeholder="$t('form.input_placeholder') + $t('form.verify_code')"
        ></el-input>
        <el-button
          v-if="register_way === 'email'"
          v-debounce
          :disabled="!isRegister || isSending"
          class="!bg-[#f5f5f5] border-0 h-[44px] w-[100px] no-left-radius"
          @click.stop="handleGetCode"
        >
          <div :class="['text-sm', 'pl-4', 'border-l', 'pr-3', 'text-[#2563EB]', { 'text-[#9A9A9A]': !isRegister || isSending }]">
            {{ emailCodeCount ? `${emailCodeCount}s` : $t('form.get_verify_code') }}
          </div>
        </el-button>
        <el-button
          v-else
          v-debounce
          :disabled="!isRegister || isSending"
          class="!bg-[#f5f5f5] border-0 h-[44px] w-[100px] no-left-radius"
          @click.stop="handleGetCode"
        >
          <div :class="['text-sm', 'pl-4', 'border-l', 'pr-3', 'text-[#2563EB]', { 'text-[#9A9A9A]': !isRegister || isSending }]">
            {{ codeCount ? `${codeCount}s` : $t('form.get_verify_code') }}
          </div>
        </el-button>
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
import { ref, reactive, computed, nextTick, watch } from 'vue'
import type { FormInstance } from 'element-plus'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/modules/user'
import { getMobileRules, getEmailRules, getPasswordRules } from '@/utils/form-rules'

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

const formRef = ref<FormInstance>()

const userStore = useUserStore()

const emits = defineEmits(['success', 'close'])

const { isOpLocalEnv } = useEnv()
const { emailCodeRule, sendEmailCode, emailCodeCount } = useEmail()
const { sendcode, codeRule, codeCount } = useMobile()

const isSending = ref(true)

const form = reactive({
  username: '',
  password: '',
  verify_code: ''
})

const register_way = ref(isOpLocalEnv.value ? 'email' : 'mobile')

// 添加重置函数
const resetForm = () => {
  form.username = ''
  form.verify_code = ''
  form.password = ''
  existing_mobile.value = false
  existing_email.value = false
  isRegister.value = false
  nextTick(() => {
    clearFormValidation()
  })
}

const handleRegisterWay = (way) => {
  resetForm()
  register_way.value = way
}

const isFormatCorrect = computed(() => {
  if (register_way.value === 'mobile') {
    return /^1[3-9]\d{9}$/.test(form.username)
  }
  return /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/.test(form.username)
})

const existing_mobile = ref(false)
const existing_email = ref(false)
const rules = computed(() => {
  return {
    username: [
      register_way.value === 'email' ? getEmailRules() : getMobileRules(),
      {
        validator: async (rule, value, callback) => {
          try {
            // 跳过空值触发和格式不符合触发
            if (form.username.trim() === '' || !isFormatCorrect.value) {
              return
            }
            // 等待验证完成
            await onUsernameBlur()

            // 检查完成后，根据isRegister的值进行验证
            if (register_way.value === 'mobile' && !isRegister.value) {
              existing_mobile.value = true
              callback(new Error(window.$t('form.mobile') + window.$t('register.unregistered')))
            } else if (register_way.value === 'email' && !isRegister.value) {
              existing_email.value = true
              callback(new Error(window.$t('form.email') + window.$t('register.unregistered')))
            }
          } catch (error) {}
        },
        trigger: 'blur'
      }
    ],
    password: [getPasswordRules()],
    verify_code: [register_way.value === 'email' ? emailCodeRule : codeRule]
  }
})

// 添加清除表单验证的方法
const clearFormValidation = () => {
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

// 添加用户名验证缓存
const usernameCache = reactive(new Map())
const isRegister = ref(false)

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
  if (register_way.value === 'email') {
    sendEmailCode(form.username)
  } else {
    sendcode(form.username)
  }
}

const handleClose = () => {
  // 重置表单
  form.username = ''
  form.verify_code = ''
  form.password = ''
  resetForm()
  emits('close')
}

const handleSubmit = () => {
  return formRef.value?.validate().then(async (valid) => {
    if (!valid) return
    console.log('valid')
    try {
      if (register_way.value === 'mobile') {
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
