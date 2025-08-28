<template>
  <ElDialog v-model="visible" :title="isEdit ? $t('action.edit') : $t('action.add')" width="600">
    <ElForm ref="formRef" :model="accountData" label-position="top">
      <ElFormItem :label="$t('account')" prop="account" :rules="accountRules">
        <ElInput v-model="accountData.account" :placeholder="$t('form_input_placeholder')" />
      </ElFormItem>
      <ElFormItem
        :label="$t('password')"
        prop="password"
        :rules="generateInputRules({ message: 'form_input_placeholder', validator: ['password'] })"
      >
        <ElInput v-model="accountData.password" :placeholder="$t('form_input_placeholder')" />
      </ElFormItem>
      <ElFormItem :label="$t('remark')" :rules="generateInputRules({ message: 'form_input_placeholder' })">
        <ElInput
          v-model="accountData.remark"
          :placeholder="$t('form_input_placeholder')"
          type="textarea"
          :rows="3"
          resize="none"
          show-word-limit
          maxlength="200"
        />
      </ElFormItem>
    </ElForm>
    <div class="w-full flex justify-center">
      <ElButton size="large" @click="close">{{ $t('action_cancel') }}</ElButton>
      <ElButton size="large" type="primary" @click="confirm">{{ $t('action_confirm') }}</ElButton>
    </div>
  </ElDialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'

import { generateInputRules } from '@/utils/form-rule'

const emits = defineEmits(['success'])

interface AccountItem {
  account: string
  password: string
  remark: string
}

const props = withDefaults(
  defineProps<{
    accountList: AccountItem[]
  }>(),
  {
    accountList: () => [],
  }
)

const visible = ref(false)
const isEdit = ref(false)

const formRef = ref()
const accountData = reactive({
  account: '',
  password: '',
  remark: '',
})

// 检查共享账号是否已添加
const checkAccountExists = (account: string): boolean => {
  return props.accountList.some(item => item.account === account)
}

const accountRules = computed(() => {
  const baseRules = generateInputRules({ message: 'form_input_placeholder', validator: ['account'] })
  const duplicateValidator = (rule: any, value: any, callback: any) => {
    if (!isEdit.value && value && checkAccountExists(value)) {
      callback(new Error($t('form_account_exit')))
    } else {
      callback()
    }
  }

  return [...baseRules, { validator: duplicateValidator, trigger: 'blur' }]
})

const clear = () => {
  accountData.account = ''
  accountData.password = ''
  accountData.remark = ''
}

const open = (data: AccountItem) => {
  clear()
  if (data) {
    isEdit.value = true
    accountData.account = data.account
    accountData.password = data.password
    accountData.remark = data.remark
  } else {
    isEdit.value = false
  }
  visible.value = true
}

const close = () => {
  clear()
  visible.value = false
}

const confirm = async () => {
  const valid = await formRef.value.validate()
  if (!valid) return
  ElMessage.success(isEdit.value ? $t('action_save_success') : $t('action_add_success'))
  emits('success', {
    account: accountData.account,
    password: accountData.password,
    remark: accountData.remark,
  })
  close()
}

defineExpose({
  open,
  close,
})
</script>

<style scoped></style>
