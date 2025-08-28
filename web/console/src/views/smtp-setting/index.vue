<template>
  <Layout class="px-[60px] py-8">
    <Header :title="$t('module.SMTP')" />

    <div class="flex-1 flex flex-col gap-4 bg-white p-6 mt-3 box-border overflow-y-auto">
      <div class="w-full px-6 py-4 border rounded hover:shadow">
        <div class="h-8 flex justify-between items-center">
          <p>{{ $t('module.SMTP_email_log') }}</p>
          <div>
            <ElSwitch v-model="openEmail" @change="handleOpenEmail"></ElSwitch>
            <span class="ml-2">{{ openEmail ? $t('action_enable') : $t('action_close') }}</span>
          </div>
        </div>

        <Form v-if="openEmail" ref="formRef" :type="'email'" />
      </div>

      <div class="w-full px-6 border rounded hover:shadow">
        <div class="h-16 flex justify-between items-center">
          <p>{{ $t('module.SMTP_mobile_log') }}</p>
          <div>
            <ElSwitch v-model="openMobile" disabled @change="handleOpenMobile" @click="handleClick"></ElSwitch>
            <span class="ml-2">{{ openMobile ? $t('action_enable') : $t('action_close') }}</span>
          </div>
        </div>

        <!-- <Form v-if="isOpLocalEnv && openMobile" /> -->
      </div>
    </div>
  </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useEnterpriseStore } from '@/stores'
import Form from './components/form.vue'

const store = useEnterpriseStore()

const openEmail = ref(false)
const openMobile = ref(false)

const formRef = ref()

const handleOpenEmail = async () => {
  if (!openEmail.value) {
    if (!formRef.value) return
    const formData = formRef.value.getData()
    if (!formData || Object.keys(formData).length === 0) return
    try {
      const requestConfig = {
        content: JSON.stringify(formData),
        enabled: false,
        type: 'smtp',
      }
      await store.saveSMTPInfo({ data: requestConfig })
    } catch (error) {
      console.log(error)
    }
  }
}

const handleOpenMobile = () => {
  ElMessage.warning($t('feature_coming_soon'))
}

const handleClick = () => {
  ElMessage.warning($t('feature_coming_soon'))
}
onMounted(() => {
  store.loadSMTPInfo().then(data => {
    openEmail.value = data[0].enabled
  })
})
</script>

<style scoped></style>
