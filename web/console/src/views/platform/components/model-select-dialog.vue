<template>
  <ElDialog
    v-model="visible"
    :title="$t('module.platform_model_add')"
    :close-on-click-modal="false"
    width="700px"
    destroy-on-close
    append-to-body
  >
    <ul class="flex flex-col gap-3 mb-8 pb-1 overflow-auto max-h-[60vh] pr-1">
      <li
        v-for="opt in channel_options"
        :key="opt.id"
        class="flex items-center gap-4 py-5 px-6 border rounded box-border"
      >
        <img
          class="flex-none w-[40px] h-[40px] object-contain"
          :src="$getRealPath({ url: `/images/platform/${opt.icon}.png` })"
        />
        <div class="flex-1 text-[#1B2B51] font-semibold">
          {{ $t(opt.label) }}
        </div>
        <ElButton
          class="flex-none !border-none !px-5 text-[#3664EF] bg-[#F6F7F8]"
          :class="opt.isAdd ? '!text-[#9A9A9A] !bg-[#F6F7F8]' : ''"
          type="primary"
          plain
          size="large"
          :disabled="opt.isAdd"
          @click="handleAdd(opt)"
        >
          {{ $t(opt.isAdd ? 'action_add_success' : 'action_add') }}
        </ElButton>
      </li>
    </ul>
  </ElDialog>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';

import { models } from '@/constants/platform/config';

const props = withDefaults(
  defineProps<{
    list: unknown[]
  }>(),
  {
    list: () => [],
  }
)

const emits = defineEmits<{
  (e: 'add', opt: unknown): void
}>()

const visible = ref(false)

const channel_options = computed(() => {
  return Object.values(models).map(item => {
    return {
      ...item,
      isAdd: item.multiple ? false : props.list.some(a => a.channel_type === item.channelType),
    }
  })
})

const handleAdd = (opt: unknown) => {
  emits('add', opt)
}

const open = () => {
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

<style></style>
