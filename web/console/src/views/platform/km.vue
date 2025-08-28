<template>
  <Layout class="px-[60px] py-8">
    <Header :title="$t('module.platform')" />
    <div class="flex-1 flex flex-col bg-white p-6 mt-3 box-border max-h-[calc(100vh-100px)] overflow-auto">
      <!-- 大模型列表 -->
      <h2 class="w-full flex items-center font-semibold text-[#1D1E1F] mb-6">
        <div class="flex-1">
          {{ $t('module.platform_model') }}
        </div>
      </h2>
      <ul v-loading="channelLoading" class="w-full flex flex-col gap-4 mb-8">
        <ModelGroupV2
          v-for="group in channelList"
          :key="group.channel_type"
          :group="group"
          @edit="handleModelEdit"
          @delete="handleModelDelete"
          @model-edit="onModelEdit"
        />
        <ElButton class="flex-none !border-none w-[106px]" type="primary" plain size="large" @click="handleModelSelect">
          + {{ $t('action_add') }}
        </ElButton>
      </ul>
    </div>
  </Layout>

  <ModelSaveDialog ref="modelSaveRef" :model-list="modelList" @success="loadModelList" />
  <ModelSelectDialog ref="modelSelectRef" :list="channelList" :model-list="modelList" @add="handleModelAdd" />
  <ModelSettingDialog ref="modelSettingRef" @success="loadModelList" />
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import ModelGroupV2 from './componentsv2/model-group.vue'

import ModelSaveDialog from './componentsv2/model-save-dialog.vue'
import ModelSelectDialog from './componentsv2/model-select-dialog.vue'
import ModelSettingDialog from './components/model-setting-dialog.vue'

import channelApi, { transformChannelData, transformModelList, type ModelOption } from '@/api/modules/channel/index'

interface ChannelGroup {
  label: string
  icon: string
  channel_type: number
  multiple: boolean
  data: any
  children: any[]
}

// 状态管理
const channelList = ref<ChannelGroup[]>([])
const channelLoading = ref(false)

const modelList = ref<ModelOption[]>([])

// 组件引用
const modelSaveRef = ref()
const modelSelectRef = ref()
const modelSettingRef = ref()

const createChannelGroup = (list: any) => {
  return list.reduce((acc: ChannelGroup[], item: any) => {
    const group = {
      name: item.name,
      label: item.label,
      icon: item.icon,
      channel_type: item.channel_type,
      multiple: false,
      data: item,
      modelValue: item.group.map(item => item.type.toString()),
      children: [item],
    }
    acc.push(group)
    return acc
  }, [])
}

const loadModelList = async () => {
  channelLoading.value = true
  try {
    const list = await channelApi.listv2()
    channelList.value = createChannelGroup(list.map(item => transformChannelData(item)))
  } finally {
    channelLoading.value = false
  }
}

const handleModelSelect = () => modelSelectRef.value.open()

const handleModelAdd = (data: ModelOption) => {
  modelSaveRef.value.open({ channel_type: data.channel_type })
}

const handleModelEdit = (data: any) => modelSaveRef.value.open(data)

const handleModelDelete = async (data: any, model: any) => {
  await ElMessageBox.confirm(window.$t('module.platform_model_delete_confirm'))
  const isChildRemove = model && data.modelOptions.length > 1
  if (isChildRemove) {
    await channelApi.update(data.channel_id, {
      channel_id: data.channel_id,
      key: data.api_key,
      base_url: data.base_url,
      config: data.config || {},
      models: data.modelOptions
        ?.map((item: any) => item.value)
        .filter((item: any) => item !== model.value)
        .join(','),
      name: data.name,
      type: data.channel_type,
      custom_config: data.custom_config,
    })
  } else {
    await channelApi.delete(data.channel_id)
  }

  ElMessage.success(window.$t('action_delete_success'))
  loadModelList()
}

const onModelEdit = ({ data, parentData }: { data: any; parentData: any }) => {
  modelSettingRef.value.open({ data: { ...parentData, ...data, id: data.value } })
}

// 初始化
const refresh = () => {
  loadModelList()

  channelApi.models.config().then(res => {
    modelList.value = transformModelList(res)
  })
}

onMounted(() => refresh())
</script>

<style scoped>
::v-deep(.el-collapse-item__arrow) {
  margin-left: 6px;
}
</style>
