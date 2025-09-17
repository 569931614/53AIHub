<template>
  <ElForm ref="form_ref" class="relative max-w-[440px] w-full" :model="form" label-position="top">
    <h4 class="text-3xl text-[#1D1E1F] font-bold text-center mb-10">
      {{ $t('account_register') }}
    </h4>

    <!-- <ElFormItem>
      <template #label>
        <span class="text-[#1D1E1F]">{{ $t('login.select_register_way') }}</span>
      </template>
      <ElRadioGroup v-model="form.username_type" @change="onUsernameTypeChange">
        <ElRadio label="email">
          {{ $t('login.email_validate') }}
        </ElRadio>
        <ElRadio label="mobile">
          {{ $t('login.mobile_validate') }}
        </ElRadio>
      </ElRadioGroup>
    </ElFormItem> -->

    <ElFormItem
      prop="username"
      class="relative"
      :rules="
        generateInputRules({
          message: `login.${form.username_type}_placeholder`,
          validator: ['text', form.username_type],
        })
      "
    >
      <template #label>
        <span class="text-[#1D1E1F]">{{ $t('account') }}</span>
      </template>
      <ElInput
        v-model="form.username"
        style="--el-input-bg-color: #f1f2f3; --el-input-border-color: transparent; --el-input-height: 44px"
        size="large"
        autocomplete="new-username"
        name="prevent_autofill_username"
        :placeholder="$t(`login.${form.username_type}_placeholder`)"
        clearable
        @blur="checkAccount"
        @input="checkAccountValidation"
      />
      <div v-if="account_exists" class="absolute -bottom-5 left-0 text-red-500 text-xs">
        {{ $t(`login.${form.username_type}_exists`) }}
        <ElButton type="primary" link size="small" class="!p-0 !bg-transparent -ml-1" @click="onLogin">
          {{ $t('action_login') }}
        </ElButton>
      </div>
    </ElFormItem>
    <ElFormItem
      class="relative"
      prop="verification_code"
      :rules="generateInputRules({ message: 'verification_code_placeholder' })"
    >
      <template #label>
        <span class="text-[#1D1E1F]">{{ $t('verification_code') }}</span>
      </template>
      <VerificationCodeInput
        ref="verification_code_input_ref"
        v-model="form.verification_code"
        :account="form.username"
        :account-type="form.username_type"
        :disabled="account_exists || !isAccountValid"
        :maxlength="form.username_type === 'mobile' ? 4 : 6"
      />
    </ElFormItem>
    <ElFormItem
      class="relative"
      prop="password"
      :rules="[
        ...generateInputRules({ message: 'login.password_placeholder', validator: ['password'] }),
        {
          min: 8,
          max: 20,
          message: $t('login.password_length'),
          trigger: 'blur',
        },
      ]"
    >
      <template #label>
        <span class="text-[#1D1E1F]">{{ $t('password') }}</span>
      </template>
      <ElInput
        v-model="form.password"
        style="--el-input-bg-color: #f1f2f3; --el-input-border-color: transparent; --el-input-height: 44px"
        size="large"
        type="password"
        show-password
        clearable
        autocomplete="new-password"
        name="prevent_autofill_password"
        :placeholder="$t('login.password_placeholder')"
      />
    </ElFormItem>

    <ElButton
      type="primary"
      round
      :loading="submitting"
      size="large"
      class="w-full mt-6 !h-10"
      :disabled="!form.username || !form.password || !form.verification_code"
      @click.stop="onRegister"
    >
      {{ $t('action_register') }}
    </ElButton>
    <div class="w-full flex-center mt-4 text-sm text-[#9A9A9A]">
      {{ $t('login.has_account') }}
      <ElButton type="primary" text class="!p-0 !bg-transparent" @click="onLogin">
        {{ $t('login.direct_login') }}
      </ElButton>
    </div>
  </ElForm>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, reactive, ref, nextTick } from 'vue'
import { useEnterpriseStore, useUserStore } from '@/stores'
import { generateInputRules } from '@/utils/form-rule'
import { authApi } from '@/api/modules/auth'
import eventBus from '@/utils/event-bus'
import { validateFormField } from '@/utils/form-validator'

interface FormType {
  username: string
  username_type: 'email' | 'mobile'
  password: string
  verification_code: string
}

const emits = defineEmits<{
  (e: 'login'): void
}>()

const user_store = useUserStore()
const enterprise_store = useEnterpriseStore()
const form_ref = ref()
const verification_code_input_ref = ref()
const form = reactive<FormType>({
  username_type: 'mobile',
  username: '',
  verification_code: '',
  password: '',
})
const submitting = ref(false)
const account_exists = ref(false)
const isAccountValid = ref(false)

// 检查账号验证状态
const checkAccountValidation = async () => {
  isAccountValid.value = await validateFormField(form_ref, 'username')
}

onMounted(() => {
  const { access_token } = user_store.info
  eventBus.on('language-change', onLanguageChange)
})
onUnmounted(() => {
  eventBus.off('language-change', onLanguageChange)
})

const onLanguageChange = () => {
  if (form_ref.value) form_ref.value.clearValidate()
  // form_ref.value.validate()
}
const reset = () => {
  form.username = ''
  form.verification_code = ''
  form.password = ''
}
const onRegister = async () => {
  let valid = await form_ref.value.validate()
  if (!valid) return
  submitting.value = true
  await checkAccount()
  if (account_exists.value) {
    ElMessage.warning(window.$t('login.account_exists'))
    return (submitting.value = false)
  }
  if (form.username_type === 'mobile') {
    valid = await verification_code_input_ref.value.validateCode()
    if (!valid) return (submitting.value = false)
  }
  await user_store
    .login({
      data: {
        username: form.username,
        password: form.password,
        verify_code: form.verification_code,
      },
      hideError: false,
    })
    .catch(err => {
      submitting.value = false
      return Promise.reject(err)
    })
  const { list = [] } = await enterprise_store.loadListData({ data: { status: 0 } }).catch(err => {
    submitting.value = false
    return Promise.reject(err)
  })
  submitting.value = false
  if (list.length > 0) return ElMessage.warning($t('login.apply_open_repetition'))
  // ElMessage.success($t('action_save_success'))
  ElMessage.success(window.$t('action_login_success'))
  reset()
  onLogin()
}
const onLogin = () => {
  emits('login')
}
const checkAccount = async () => {
  account_exists.value = false
  const valid = await form_ref.value.validateField('username')
  if (!valid) return
  const { exists = false } = await authApi.checkAccount({ data: { account: form.username } })
  account_exists.value = exists
  return exists
}

const onUsernameTypeChange = async () => {
  form.username = ''
  verification_code_input_ref.value.reset()
  await nextTick()
  form_ref.value.clearValidate('username')
  form_ref.value.clearValidate('verification_code')
}

defineExpose({
  reset,
})
</script>

<style scoped lang="scss">
::v-deep(.el-step__line) {
  left: 70% !important;
  right: -30% !important;
  background-color: #e9ebee !important;
}
::v-deep(.el-step__title) {
  margin-top: 8px;
  font-weight: 400 !important;
}
::v-deep(.el-step__head.is-finish .el-step__line) {
  background-color: #3664ef !important;
  // background-color: #82A2F6 !important;
}
::v-deep(.el-step__head.is-process .el-step__line) {
  // background-color: #3664EF !important;
}
</style>
