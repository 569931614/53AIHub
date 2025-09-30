<template>
  <el-form ref="formRef" label-position="top" :model="form" :rules="rules" :validate-on-rule-change="true" @keyup.enter="handleSubmit">
    <el-form-item v-if="userStore.info.mobile" prop="old_code">
      <el-input v-model="form.old_code" v-trim size="large" :placeholder="$t('form.input_placeholder') + $t('form.verify_code')">
        <template #append>
          <el-button :disabled="isSendingOld" @click.stop="handleGetOldCode">
            <div :class="codeCount ? 'text-[#9A9A9A]' : 'text-[#2563EB]'">
              {{ codeCount ? `${codeCount}s` : $t('form.get_verify_code') }}
            </div>
          </el-button>
        </template>
      </el-input>
    </el-form-item>

    <el-form-item :label="getMobileLabel()" prop="new_mobile">
      <el-input v-model="form.new_mobile" v-trim size="large" :placeholder="getMobilePlaceholder()" clearable />
    </el-form-item>

    <el-form-item :label="$t('form.verify_code')" prop="new_code">
      <el-input v-model="form.new_code" v-trim size="large" :placeholder="$t('form.input_placeholder') + $t('form.verify_code')">
        <template #append>
          <el-button :disabled="isSendingNew" @click.stop="handleGetNewCode">
            <div :class="newCodeCount ? 'text-[#9A9A9A]' : 'text-[#2563EB]'">
              {{ newCodeCount ? `${newCodeCount}s` : $t('form.get_verify_code') }}
            </div>
          </el-button>
        </template>
      </el-input>
    </el-form-item>

    <!-- 更换/绑定按钮 -->
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
import { getMobileRules } from '@/utils/form-rules'
import commonApi from '@/api/modules/common'

import { useUserStore } from '@/stores/modules/user'
import useMobile from '@/hooks/useMobile'

const emits = defineEmits(['success', 'close'])

const userStore = useUserStore()
const { sendcode, codeRule, codeCount } = useMobile()
const { sendcode: newSendCode, codeRule: newCodeRule, codeCount: newCodeCount } = useMobile()

const formRef = ref<FormInstance>()

const form = reactive({
  old_code: '',
  new_mobile: '',
  new_code: ''
})

const rules = reactive({
  new_mobile: [getMobileRules()],
  old_code: [codeRule],
  new_code: [newCodeRule]
})

const isSendingOld = ref(false)
const isSendingNew = ref(false)

// 计算属性
const isMobile = computed(() => /^1[3-9]\d{9}$/.test(form.new_mobile))

// 工具函数
const getMobileLabel = () => {
  return userStore.info.mobile ? window.$t('form.new_mobile') : window.$t('form.mobile')
}

const getMobilePlaceholder = () => {
  return userStore.info.mobile
    ? window.$t('form.input_placeholder') + window.$t('form.new_mobile')
    : window.$t('form.input_placeholder') + window.$t('form.mobile')
}

const handleGetOldCode = () => {
  sendcode(userStore.info.mobile)
  isSendingOld.value = Boolean(codeCount.value)
}

const handleGetNewCode = () => {
  if (!isMobile.value) return
  newSendCode(form.new_mobile)
  isSendingNew.value = Boolean(newCodeCount.value)
}

// 重置表单
const resetForm = () => {
  Object.assign(form, {
    old_code: '',
    new_mobile: '',
    new_code: ''
  })
  formRef.value?.resetFields()
}

const handleClose = () => {
  resetForm()
  emits('close')
}

// 更换手机号逻辑
const performChangeMobile = async (id) => {
  await commonApi.verifycode({
    mobile: userStore.info.mobile,
    verifycode: form.old_code,
    type: '1'
  })

  await commonApi.verifycode({
    mobile: form.new_mobile,
    verifycode: form.new_code,
    type: '1'
  })

  await userStore.change_mobile(
    {
      new_code: form.new_code,
      new_mobile: form.new_mobile,
      old_code: form.old_code
    },
    id
  )
  ElMessage.success(window.$t('profile.change') + window.$t('status.success'))
}

// 绑定手机号逻辑
const performBindMobile = async (id) => {
  await commonApi.verifycode({
    mobile: form.new_mobile,
    verifycode: form.new_code,
    type: '1'
  })

  await userStore.change_mobile(
    {
      new_code: form.new_code,
      new_mobile: form.new_mobile
    },
    id
  )
  ElMessage.success(window.$t('profile.bind') + window.$t('status.success'))
}
const handleSubmit = () => {
  return formRef.value?.validate().then(async (valid) => {
    if (!valid) return

    try {
      await userStore.getUserInfo()
      const id = userStore.info.user_id

      if (userStore.info.mobile) {
        await performChangeMobile(id)
      } else {
        await performBindMobile(id)
      }

      resetForm()
      emits('success')
    } catch (error) {}
  })
}

defineExpose({
  resetForm
})
</script>

<style scoped></style>
