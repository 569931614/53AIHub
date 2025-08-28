<template>
  <ElDialog v-model="visible" :title="$t('toolbox.account_access')" width="600" align-center>
    <div class="text-[#999999]">{{ $t('toolbox.account_text') }}</div>

    <div v-loading="loading" class="h-72 overflow-y-auto flex flex-col gap-3 mt-4">
      <div v-for="(item, index) in accountList" :key="index" class="bg-[#F2F7FF] flex flex-col gap-7 p-5 rounded">
        <div class="flex gap-6">
          <span class="text-[#999999] flex-none">{{ $t('form.account') }}</span>
          <span class="text-[#1D1E1F] break-words whitespace-pre-wrap min-w-0">{{ item.account }}</span>
          <el-link v-copy="item.account" :underline="false">
            <SvgIcon name="copy" />
          </el-link>
        </div>
        <div class="flex gap-6">
          <span class="text-[#999999] flex-none">{{ $t('form.password') }}</span>
          <span class="text-[#1D1E1F] break-words whitespace-pre-wrap min-w-0">{{ item.password }}</span>
          <el-link v-copy="item.password" :underline="false">
            <SvgIcon name="copy" />
          </el-link>
        </div>
        <div v-if="item.remark" class="w-full flex gap-6">
          <span class="text-[#999999] flex-none">{{ $t('form.remark') }}</span>
          <span class="text-[#1D1E1F] flex-1 break-words whitespace-pre-wrap min-w-0">{{ item.remark }}</span>
        </div>
      </div>
    </div>

    <div class="flex justify-center mt-4">
      <ElButton type="primary" size="large" @click="handleVisit">{{ $t('toolbox.click_access') }}</ElButton>
      <ElButton size="large" @click="close">{{ $t('action.cancel') }}</ElButton>
    </div>
  </ElDialog>
</template>

<script setup>
import { ref } from 'vue'

import { ElButton } from 'element-plus'
import linksApi from '@/api/modules/links'
import SvgIcon from '@/components/SvgIcon.vue'

const visible = ref(false)
const loading = ref(false)

const accountList = ref()
const url = ref()

const open = async (item) => {
  visible.value = true
  loading.value = true
  url.value = item.url
  try {
    const info = await linksApi.detail(item.id)
    accountList.value = JSON.parse(info.data.shared_account)
  } catch (error) {
    console.log(error)
  } finally {
    loading.value = false
  }
}

const close = () => {
  visible.value = false
}

const handleVisit = () => {
  window.open(url.value, '_blank')
}

defineExpose({
  open,
  close
})
</script>

<style scoped></style>
