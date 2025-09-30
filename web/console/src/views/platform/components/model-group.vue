<template>
  <li class="w-full p-5 border rounded box-border overflow-hidden bg-[#F8F9FA]">
    <div class="flex items-center gap-4">
      <img
        class="flex-none w-[40px] h-[40px] object-contain rounded-full overflow-hidden"
        :src="$getRealPath({ url: `/images/platform/${group.icon}.png` })"
      />
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
        <ElCollapse :model-value="['models']" class="w-full !border-none">
          <ElCollapseItem name="models">
            <template #title>
              <span class="text-[#4F5052] text-sm">
                {{
                  $t('module.platform_model_models_total', {
                    total: group.data.modelOptions.length,
                  })
                }}
              </span>
            </template>
            <ul class="w-full flex flex-col gap-y-5 box-border overflow-auto">
              <li
                v-for="model in group.data.modelOptions"
                :key="model.value"
                class="w-full box-border flex items-center gap-2 group"
              >
                <img
                  v-if="model.icon"
                  class="flex-none w-[20px] h-[20px] object-contain"
                  :src="$getRealPath({ url: `/images/platform/${model.icon}.png` })"
                />
                <label class="flex-none text-sm text[#1D1E1F]">{{ model.label }}</label>
                <ElIcon
                  class="flex-none cursor-pointer hover:opacity-70"
                  size="14"
                  @click="$emit('model-edit', { data: model, parentData: group.data })"
                >
                  <Setting />
                </ElIcon>
                <template v-if="getTestResult(group.data, model)">
                  <template v-if="getTestResult(group.data, model).loading"></template>
                  <el-tag v-else-if="getTestResult(group.data, model).success" type="success">
                    {{ $t('action_test_success') }}
                  </el-tag>
                  <el-tag v-else type="danger">{{ $t('action_test_failed') }}</el-tag>
                </template>
                <div class="flex-1" />
                <el-button
                  :loading="getTestResult(group.data, model)?.loading"
                  class="group-hover:visible invisible"
                  type="primary"
                  link
                  @click="handleTest(model, group.data)"
                >
                  {{ $t('action_test') }}
                </el-button>
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
                  $t('module.platform_model_models_total', {
                    total: group.children.length,
                  })
                }}
              </span>
            </template>
            <ul class="w-full flex flex-col gap-y-5 box-border overflow-auto">
              <div v-for="channel in group.children" :key="channel.channelId">
                <li
                  v-for="model in channel.modelOptions"
                  :key="model.value"
                  class="w-full box-border flex items-center gap-2 group"
                >
                  <img
                    v-if="model.icon"
                    class="flex-none w-[20px] h-[20px] object-contain"
                    :src="$getRealPath({ url: `/images/platform/${model.icon}.png` })"
                  />
                  <label class="flex-none text-sm text[#1D1E1F]">{{ model.label }}</label>
                  <ElIcon class="flex-none cursor-pointer hover:opacity-70" size="14" @click="$emit('edit', channel)">
                    <Setting />
                  </ElIcon>

                  <template v-if="getTestResult(channel, model)">
                    <template v-if="getTestResult(channel, model).loading"></template>
                    <el-tag v-else-if="getTestResult(channel, model).success" type="success">
                      {{ $t('action_test_success') }}
                    </el-tag>
                    <el-tag v-else type="danger">{{ $t('action_test_failed') }}</el-tag>
                  </template>

                  <div class="flex-1" />

                  <el-button
                    :loading="getTestResult(channel, model)?.loading"
                    class="group-hover:visible invisible"
                    type="primary"
                    link
                    @click="handleTest(model, channel)"
                  >
                    {{ $t('action_test') }}
                  </el-button>
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
import { ref } from 'vue'
import { Delete, Remove, Setting } from '@element-plus/icons-vue'
import channelApi from '@/api/modules/channel/index'

interface ChannelGroup {
  label: string
  icon: string
  channelType: number
  multiple: boolean
  data: any
  children: any[]
}

withDefaults(
  defineProps<{
    group: ChannelGroup
  }>(),
  {
    group: () => ({
      label: '',
      icon: '',
      channelType: 0,
      multiple: false,
      data: {},
      children: [],
    }),
  }
)

defineEmits<{
  edit: [data: any]
  delete: [data: any, model: any]
  'model-edit': [{ data: any; parentData: any }]
}>()

const testMap = ref<Record<string, { loading: boolean; success: boolean }>>({})

const getTestKey = (data: any, model: any) => {
  return `${data.channel_id}-${model.value}`
}

const getTestResult = (data: any, model: any) => {
  return testMap.value[getTestKey(data, model)]
}

const handleTest = (model: any, data: any) => {
  const key = getTestKey(data, model)
  testMap.value[key] = { loading: true, success: false }
  return channelApi
    .test(data.channel_id, {
      model: model.value,
    })
    .then(res => {
      testMap.value[key] = { loading: false, success: res ? res.success : false }
    })
    .catch(() => {
      testMap.value[key] = { loading: false, success: false }
    })
    .finally(() => {
      if (testMap.value[key].success) {
        ElMessage.success(window.$t('platform.model_test_success', { platform: data.name }))
      } else {
        ElMessage.error(window.$t('platform.model_test_failed'))
      }
    })
}
</script>
