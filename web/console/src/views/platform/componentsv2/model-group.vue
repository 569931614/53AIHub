<template>
  <li class="w-full p-5 border rounded box-border overflow-hidden bg-[#F8F9FA]">
    <div class="flex items-center gap-4">
      <img class="flex-none w-[40px] h-[40px] object-contain rounded-full overflow-hidden" :src="group.icon" />
      <div class="flex-1 text-[#1B2B51] font-semibold">
        {{ group.label }}
      </div>
      <template v-if="!group.multiple">
        <ElButton class="flex-none !px-5" type="default" size="large" @click.stop="$emit('edit', group.data)">
          {{ $t('action_setting') }}
        </ElButton>
        <ElButton
          v-debounce
          class="flex-none !ml-0 !px-3"
          type="default"
          size="large"
          @click.stop="$emit('delete', group.data, null)"
        >
          <ElIcon size="14">
            <Delete />
          </ElIcon>
        </ElButton>
      </template>
    </div>
    <div class="w-full flex flex-col bg-white rounded overflow-hidden gap-4 mt-6 px-5">
      <template v-if="!group.multiple">
        <ElCollapse :model-value="group.modelValue" class="w-full !border-none">
          <ElCollapseItem v-for="item in group.data.group" :key="item.type" :name="item.type">
            <template #title>
              <span class="text-[#4F5052] text-sm">
                {{
                  $t('platform.model_tip', {
                    total: item.options.length,
                    type: item.typeName,
                  })
                }}
              </span>
            </template>
            <ul class="w-full flex flex-col gap-y-5 box-border overflow-auto">
              <li v-for="model in item.options" :key="model.value" class="w-full box-border flex items-center gap-2">
                <img v-if="model.icon" class="flex-none size-5 object-contain" :src="model.icon" />
                <label class="flex-none text-sm text[#1D1E1F]">{{ model.label }}</label>
                <ElIcon
                  class="flex-none cursor-pointer hover:opacity-70"
                  size="14"
                  @click="$emit('model-edit', { data: model, parentData: group.data })"
                >
                  <Setting />
                </ElIcon>
                <div class="flex-1" />
                <ElIcon
                  class="flex-none cursor-pointer hover:opacity-70"
                  size="16"
                  color="#F04F4D"
                  @click="$emit('delete', group.data, model)"
                >
                  <Remove />
                </ElIcon>
              </li>
            </ul>
          </ElCollapseItem>
        </ElCollapse>
      </template>
      <template v-else>
        <!-- 同一个模型下有多个模型， azure 可以添加多个 -->
        <ElCollapse :model-value="['models']" class="w-full !border-none">
          <ElCollapseItem name="models">
            <template #title>
              <span class="text-[#4F5052] text-sm">
                {{
                  $t('platform.model_tip', {
                    total: group.children.length,
                    type: group.typeName,
                  })
                }}
              </span>
            </template>
            <ul class="w-full flex flex-col gap-y-5 box-border overflow-auto">
              <div v-for="channel in group.children" :key="channel.channelId">
                <li v-for="model in channel.group" :key="model.value" class="w-full box-border flex items-center gap-2">
                  <img v-if="model.icon" class="flex-none size-5 object-contain" :src="model.icon" />
                  <label class="flex-none text-sm text[#1D1E1F]">{{ model.label }}</label>
                  <ElIcon class="flex-none cursor-pointer hover:opacity-70" size="14" @click="$emit('edit', channel)">
                    <Setting />
                  </ElIcon>
                  <div class="flex-1" />
                  <ElIcon
                    class="flex-none cursor-pointer hover:opacity-70"
                    size="16"
                    color="#F04F4D"
                    @click="$emit('delete', channel, null)"
                  >
                    <Remove />
                  </ElIcon>
                </li>
              </div>
            </ul>
          </ElCollapseItem>
        </ElCollapse>
      </template>
    </div>
  </li>
</template>

<script setup lang="ts">
import { Delete, Remove, Setting } from '@element-plus/icons-vue'

interface ChannelGroup {
  label: string
  icon: string
  channelType: number
  multiple: boolean
  data: any
  children: any[]
}

defineProps<{
  group: ChannelGroup
}>()

defineEmits<{
  edit: [data: any]
  delete: [data: any, model: any]
  'model-edit': [{ data: any; parentData: any }]
}>()
</script>
