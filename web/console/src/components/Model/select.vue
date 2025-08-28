<template>
  <SelectPlus v-model="value" size="large" :use-i18n="false" :options="options" @change="onModelChange">
    <template #item_after="{ data }">
      <el-tooltip v-if="(data as any).vision" :content="$t('support_image')" placement="top">
        <div class="flex-center inline-flex align-middle ml-1 w-4 h-4 bg-[#FDF8EB] rounded-sm">
          <el-icon size="10px" color="#F0A105">
            <View />
          </el-icon>
        </div>
      </el-tooltip>
    </template>
  </SelectPlus>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { View } from '@element-plus/icons-vue'
import { loadModels } from './index'
import { ModelUseType } from '@/constants/platform/config'

// 定义类型，兼容 SelectPlus 组件的接口
interface ModelOption {
  value: string
  label: string
  icon?: string
  vision?: boolean
}

interface ChannelOption {
  value: string
  label: string
  icon?: string
  options: ModelOption[]
}

const props = withDefaults(
  defineProps<{
    modelValue: string
    type?: ModelUseType
  }>(),
  {
    modelValue: '',
  }
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'change', value: string, option: ChannelOption | ModelOption): void
}>()

const options = ref<ChannelOption[]>([])

const value = computed({
  get: () => props.modelValue,
  set: value => emit('update:modelValue', value),
})

const onModelChange = (result: { value: string; option: any }) => {
  emit('update:modelValue', result.value)
  emit('change', result.value, result.option)
}

const loadChannelOptions = async () => {
  try {
    const modelList = await loadModels(props.type)

    options.value = modelList
    // 如果modelValue 有值， 需要找modelList是否存在， 如果存在， 则设置value = ''
    if (props.modelValue) {
      const options = modelList.map(item => item.options).flat()
      const option = options.find(item => item.value === props.modelValue)
      if (!option) {
        value.value = ''
      }
    }
  } catch (error) {
    console.error('Failed to load channel options:', error)
  }
}

onMounted(() => {
  loadChannelOptions()
})
</script>
