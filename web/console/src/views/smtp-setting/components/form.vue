<template>
  <div v-if="props.type === 'email'" class="mt-5 w-3/5">
    <ElForm ref="formRef" :model="formData" size="large" label-position="top">
      <ElFormItem
        :label="$t('module.SMTP_server')"
        prop="smtp_host"
        :rules="generateInputRules({ message: 'form_input_placeholder' })"
      >
        <ElInput v-model="formData.smtp_host" :placeholder="$t('form_input_placeholder')" clearable />
      </ElFormItem>

      <ElFormItem
        :label="$t('module.SMTP_port')"
        prop="smtp_port"
        :rules="generateInputRules({ message: 'form_input_placeholder', validator: ['port'] })"
      >
        <ElInput v-model="formData.smtp_port" :placeholder="$t('form_input_placeholder')" clearable />
      </ElFormItem>

      <ElFormItem
        :label="$t('module.SMTP_email_account')"
        prop="smtp_username"
        :rules="generateInputRules({ message: 'form_input_placeholder', validator: ['email'] })"
      >
        <ElInput v-model="formData.smtp_username" :placeholder="$t('form_input_placeholder')" clearable />
      </ElFormItem>

      <ElFormItem
        :label="$t('module.SMTP_email_password')"
        prop="smtp_password"
        :rules="generateInputRules({ message: 'form_input_placeholder', validator: ['password'] })"
      >
        <ElInput
          v-model="formData.smtp_password"
          type="password"
          :placeholder="$t('form_input_placeholder')"
          clearable
          show-password
        />
      </ElFormItem>

      <ElFormItem
        :label="$t('module.SMTP_addresser_email')"
        prop="smtp_from"
        :rules="generateInputRules({ message: 'form_input_placeholder', validator: ['email'] })"
      >
        <ElInput v-model="formData.smtp_from" :placeholder="$t('form_input_placeholder')" clearable />
      </ElFormItem>

      <ElFormItem :label="$t('module.SMTP_openTLS')">
        <ElSwitch v-model="formData.smtp_is_ssl"></ElSwitch>
      </ElFormItem>

      <ElFormItem
        :label="$t('module.SMTP_receiver_email')"
        prop="smtp_to"
        :rules="generateInputRules({ message: 'form_input_placeholder', validator: ['email'] })"
      >
        <div class="w-full flex gap-3">
          <ElInput v-model="formData.smtp_to" :placeholder="$t('form_input_placeholder')" clearable />
          <ElButton :loading="isSending" type="primary" plain :disabled="countDown > 0" @click="handleSendEmail">
            {{ countDown > 0 ? `${countDown}s` : $t('module.SMTP_send_email') }}
          </ElButton>
        </div>
      </ElFormItem>
    </ElForm>

    <div class="flex">
      <ElButton :loading="isSaving" type="primary" class="w-24 h-9" @click="handleSave">{{
        $t('action.save')
      }}</ElButton>
      <ElButton class="w-24 h-9" @click="handleReset('email')">{{ $t('action_reset') }}</ElButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue';
import { generateInputRules } from '@/utils/form-rule';
import { useEnterpriseStore } from '@/stores';

const props = withDefaults(
  defineProps<{
    type: string
  }>(),
  {
    type: '',
  }
)

const isSaving = ref(false)
const isSending = ref(false)

const store = useEnterpriseStore()

const formRef = ref()

const countDown = ref(0)
let timer: NodeJS.Timeout | null = null

const formData = reactive({
  smtp_host: '',
  smtp_port: '',
  smtp_username: '',
  smtp_password: '',
  smtp_from: '',
  smtp_to: '',
  smtp_is_ssl: true,
})

const handleSendEmail = async () => {
  try {
    isSending.value = true
    const config = {
      from: formData.smtp_from,
      host: formData.smtp_host,
      is_ssl: formData.smtp_is_ssl,
      password: formData.smtp_password,
      port: Number(formData.smtp_port),
      to: formData.smtp_to,
      username: formData.smtp_username,
    }
    await store.sendTestEmail({ data: config })
    isSending.value = false
    ElMessage.success(window.$t('action_send_success'))
    countDown.value = 60
    timer = setInterval(() => {
      countDown.value--
      if (countDown.value <= 0) {
        clearInterval(timer as NodeJS.Timeout)
        timer = null
        countDown.value = 0
      }
    }, 1000)
  } catch (error) {
    console.log(error)
    isSending.value = false
  }
}

const handleReset = (data: string) => {
  if (data === 'email') {
    formData.smtp_host = ''
    formData.smtp_port = ''
    formData.smtp_username = ''
    formData.smtp_password = ''
    formData.smtp_from = ''
    formData.smtp_to = ''
    formData.smtp_is_ssl = true
  }
}

const handleSave = async () => {
  const valid = await formRef.value.validate()
  if (!valid) return
  isSaving.value = true
  const requestConfig = {
    content: JSON.stringify(formData),
    enabled: true,
    type: 'smtp',
  }
  try {
    await store.saveSMTPInfo({ data: requestConfig })
    ElMessage.success(window.$t('action_save_success'))
  } catch (error) {
    console.log(error)
  } finally {
    isSaving.value = false
  }
}

const getData = () => {
  store.loadSMTPDetail({ data: { type: 'smtp' } }).then(data => {
    const res = JSON.parse(data.content)
    formData.smtp_host = res.smtp_host || ''
    formData.smtp_port = res.smtp_port || ''
    formData.smtp_username = res.smtp_username || ''
    formData.smtp_password = res.smtp_password || ''
    formData.smtp_to = res.smtp_to || ''
    formData.smtp_from = res.smtp_from || ''
    formData.smtp_is_ssl = res.smtp_is_ssl
  })
  return { ...formData }
}

defineExpose({
  getData,
})

onMounted(() => {
  getData()
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped></style>
