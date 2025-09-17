<template>
  <div
    class="w-[50px] h-[50px] rounded overflow-hidden relative cursor-pointer group"
    :class="[disabled ? 'cursor-not-allowed' : '']"
    @click.stop="onSelectFile"
  >
    <slot>
      <ElButton v-if="showText" type="text">{{ $t(text) }}</ElButton>
      <template v-else>
        <img v-if="modelValue" class="w-full h-full object-cover" :src="modelValue" alt="logo" />
        <div v-else class="w-full h-full flex items-center justify-center border rounded upload-image-placeholder">
          <ElIcon size="16" color="#9A9A9A">
            <Plus />
          </ElIcon>
        </div>
        <div
          class="hidden group-hover:flex absolute top-0 right-0 bottom-0 left-0 bg-black bg-opacity-40 items-center justify-center gap-6 text-xs text-white"
        >
          <slot name="mask-text">
            {{ $t(text) }}
          </slot>
        </div>
      </template>
    </slot>
    <CropperDialog
      ref="cropper_ref"
      action="python"
      :cropper-disabled="cropperDisabled"
      :allow-type-list="allowTypeList"
      :limit-size="2"
      @confirm="onConfirm"
    />
  </div>
</template>

<script setup lang="ts">
import { Plus } from '@element-plus/icons-vue'
import { ref } from 'vue'
import CropperDialog from '@/components/CropperDialog/index.vue'

const props = withDefaults(
  defineProps<{
    modelValue?: string
    text?: string
    showText?: boolean
    cropperDisabled?: boolean
    allowTypeList?: string[]
    disabled?: boolean
  }>(),
  {
    modelValue: '',
    text: 'action_replace',
    showText: false,
    cropperDisabled: false,
    disabled: false
  }
)
const emits = defineEmits<{
  (e: 'update:modelValue', data: string): void
  (e: 'confirm', result: { url: string }): void
}>()

const cropper_ref = ref()

const onSelectFile = () => {
  if (props.disabled) return
  cropper_ref.value.uploadFile()
}
const onConfirm = (data: { url: string }) => {
  emits('update:modelValue', data.url)
  emits('confirm', { url: data.url })
}
</script>

<style>
.el-form-item.is-error .upload-image-placeholder {
  border-color: #f56c6c;
}
</style>
