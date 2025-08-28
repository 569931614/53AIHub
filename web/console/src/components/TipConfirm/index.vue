<template>
  <ElDialog
    v-model="visible"
    header-class="hidden"
    width="520px"
    :style="{
      '--el-dialog-border-radius': '16px',
    }"
    :close-on-click-modal="false"
    :show-close="false"
    destroy-on-close
    append-to-body
    center
    align-center
  >
    <header class="p-2 flex items-center justify-between">
      <h1 class="text-lg font-semibold text-[#1D1E1F]">{{ title }}</h1>
      <ElIcon class="font-bold" :size="20" color="#9A9A9A" @click="close">
        <Close />
      </ElIcon>
    </header>
    <section class="p-2 text-base text-[#535456]">
      {{ content }}
    </section>
    <template #footer>
      <div class="flex items-center justify-center my-2">
        <ElButton v-if="showCancelButton" size="large" @click="close">
          {{ cancelButtonText }}
        </ElButton>
        <ElButton v-if="showConfirmButton" class="!px-8" type="primary" size="large" @click="close">
          {{ confirmButtonText }}
        </ElButton>
      </div>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
import { Close } from '@element-plus/icons-vue';

import { ref, onMounted, onUnmounted } from 'vue';

withDefaults(
  defineProps<{
    title: string
    content: string
    confirmButtonText: string
    showConfirmButton: boolean
    cancelButtonText: string
    showCancelButton: boolean
  }>(),
  {
    title: '',
    content: '',
    confirmButtonText: window.$t('action_confirm'),
    showConfirmButton: true,
    cancelButtonText: window.$t('action_cancel'),
    showCancelButton: true,
  }
)

const visible = ref(false)

onMounted(() => {})
onUnmounted(() => {})

const open = async () => {
  visible.value = true
}
const close = () => {
  visible.value = false
}

defineExpose({
  open,
  close,
})
</script>

<style scoped lang="scss"></style>
