<template>
  <ElDrawer
    v-model="visible"
    :title="$t('module.platform_model_add')"
    :close-on-click-modal="false"
    size="700px"
    destroy-on-close
    append-to-body
  >
    <ul class="flex flex-col gap-3">
      <li
        v-for="opt in channel_options"
        :key="opt.platform_id"
        class="h-[72px] flex items-center gap-4 py-5 px-6 rounded bg-[#F8F9FA]"
      >
        <img class="flex-none size-10 object-contain" :src="opt.icon" />
        <div class="flex-1 text-[#1B2B51] font-semibold">
          {{ $t(opt.platform_name) }}
        </div>
        <ElButton
          class="flex-none !border-none"
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
  </ElDrawer>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

import type { ModelOption } from '@/api/modules/channel/index'

const props = withDefaults(
  defineProps<{
    list: unknown[]
    modelList: ModelOption[]
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
  return props.modelList.map(item => {
    return {
      ...item,
      isAdd: item.can_multiple ? false : props.list.some(a => a.channel_type === item.channel_type),
    }
  })
})

const handleAdd = (opt: unknown) => {
  emits('add', opt)
  visible.value = false
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
