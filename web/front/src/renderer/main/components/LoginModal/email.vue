<template>
  <el-form ref="formRef" label-position="top" :model="form" :rules="rules" :validate-on-rule-change="true" @keyup.enter="handleSubmit">
    <el-form-item :label="$t('form.new_email')" prop="email">
      <el-input v-model="form.email" v-trim size="large" :placeholder="$t('form.input_placeholder') + $t('form.email')" clearable />
    </el-form-item>
    <el-form-item :label="$t('form.verify_code')" prop="verify_code">
      <el-input v-model="form.verify_code" v-trim size="large" :placeholder="$t('form.input_placeholder') + $t('form.verify_code')">
        <template #append>
          <el-button :disabled="isSending || Boolean(emailCodeCount)" @click.stop="handleGetCode">
            <div :class="emailCodeCount ? 'text-[#9A9A9A]' : 'text-[#2563EB]'">
              {{ emailCodeCount ? `${emailCodeCount}s` : $t('form.get_verify_code') }}
            </div>
          </el-button>
        </template>
      </el-input>
    </el-form-item>

    <!-- 更换按钮 -->
    <div class="flex justify-end mt-7.5">
      <el-button class="w-24 h-9" @click="handleClose">
        {{ $t('action.cancel') }}
      </el-button>
      <el-button type="primary" class="w-24 h-9" @click="handleSubmit">
        {{ $t('action.ok') }}
      </el-button>
    </div>
  </el-form>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { getEmailRules } from '@/utils/form-rules'
import commonApi from '@/api/modules/common'

import { useUserStore } from '@/stores/modules/user'
import useEmail from '@/hooks/useEmail'

const emits = defineEmits(['success', 'close'])

const userStore = useUserStore()
const { sendEmailCode, emailCodeRule, emailCodeCount } = useEmail()

const formRef = ref<FormInstance>()

const form = reactive({
  email: '',
  verify_code: ''
})

const isSending = ref(false)

// 计算属性
const isEmail = computed(() => /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/.test(form.email))

const rules = computed(() => ({
  email: [getEmailRules()],
  verify_code: [emailCodeRule]
}))

// 验证码发送
const handleGetCode = () => {
  if (!isEmail.value) return
  isSending.value = true
  sendEmailCode(form.email).finally(() => {
    isSending.value = false
  })
}

// 关闭弹窗
const handleClose = () => {
  emits('close')
}

// 重置表单
const resetForm = () => {
  Object.assign(form, {
    email: '',
    verify_code: ''
  })
  formRef.value?.resetFields()
}

// 提交表单
const handleSubmit = () => {
  return formRef.value?.validate().then(async (valid) => {
    if (!valid) return

    try {
      await userStore.getUserInfo()
      const id = userStore.info.user_id

      await commonApi.verifyEmailcode(
        {
          email: form.email,
          code: form.verify_code
        },
        id
      )

      const message = userStore.info.email
        ? window.$t('profile.bind') + window.$t('status.success')
        : window.$t('profile.change') + window.$t('status.success')

      ElMessage.success(message)
      emits('success')
    } catch (error) {}
  })
}

defineExpose({
  resetForm
})
</script>

<style scoped></style>
