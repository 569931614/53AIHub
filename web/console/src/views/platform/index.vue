<template>
  <Layout class="px-[60px] py-8">
    <Header :title="$t('module.platform')" />
    <div class="flex-1 flex flex-col bg-white p-6 mt-3 box-border max-h-[calc(100vh-100px)] overflow-auto">
      <!-- 平台列表 -->
      <!-- #ifndef KM -->
      <template v-for="group in providerGroupList" :key="group.label">
        <h2 class="font-semibold text-[#1D1E1F] mb-6">
          {{ $t(group.label) }}
        </h2>
        <ul v-loading="providerLoading" class="flex flex-wrap gap-4 mb-8">
          <ProviderCard
            v-for="provider in group.children"
            :key="provider.id"
            :provider="provider"
            :all-total="allTotal"
            @authorize="handleProviderAuthorize"
            @add="handleAgentAdd"
            @delete="handleProviderDelete"
          />
        </ul>
      </template>
      <!-- #endif -->

      <!-- 大模型列表 -->
      <h2 class="w-full flex items-center font-semibold text-[#1D1E1F] mb-6">
        <div class="flex-1">
          {{ $t('module.platform_model') }}
        </div>
      </h2>
      <ul v-loading="channelLoading" class="w-full flex flex-col gap-4 mb-8">
        <ModelGroup
          v-for="group in channelList"
          :key="group.channelType"
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

  <ModelSaveDialog ref="modelSaveRef" @success="loadModelList" />
  <ModelSelectDialog ref="modelSelectRef" :list="channelList" @add="handleModelAdd" />
  <ProviderAuthorizeDialog ref="authorizeRef" @success="loadProviderList" />
  <ModelSettingDialog ref="modelSettingRef" @success="loadModelList" />
  <AgentListDrawer ref="agentListDrawerRef" @change="onAgentListChange" />
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import AgentListDrawer from './components/agent-list-drawer.vue'

import ModelSaveDialog from './components/model-save-dialog.vue'
import ModelSelectDialog from './components/model-select-dialog.vue'
import ModelSettingDialog from './components/model-setting-dialog.vue'

import ProviderAuthorizeDialog from './components/provider-authorize-dialog.vue'
import ProviderCard from './components/provider-card.vue'
import ModelGroup from './components/model-group.vue'

import { PROVIDER_VALUE } from '@/constants/platform/provider'
import type { ModelConfig, ProviderConfig } from '@/constants/platform/config'
import { getModelByChannelType, getModelChannelTypes, getProvidersByAuth } from '@/constants/platform/config'
import { isInternalNetwork } from '@/utils'
import { agentApi, channelApi, providerApi } from '@/api'
import TipConfirm from '@/components/TipConfirm/setup'

// 类型定义
interface ProviderOption extends ProviderConfig {
  connected: boolean
  authed_time: string
  client_id: string
  client_secret: string
  agentTotal: number
  channelLoading: boolean
  provider_id?: number
}

interface ChannelGroup {
  label: string
  icon: string
  channelType: number
  multiple: boolean
  data: any
  children: any[]
}

interface ProviderGroup {
  label: string
  children: ProviderOption[]
}

// 工具函数
const createProviderOption = (item: ProviderConfig): ProviderOption => ({
  ...item,
  connected: false,
  authed_time: '',
  client_id: '',
  client_secret: '',
  agentTotal: 0,
  channelLoading: !item.auth,
})

// 状态管理
const authProviders = ref<ProviderOption[]>(getProvidersByAuth(true).map(createProviderOption))
const agentProviders = ref<ProviderOption[]>(getProvidersByAuth(false).map(createProviderOption))
const channelList = ref<ChannelGroup[]>([])
const allTotal = ref(0)
const providerLoading = ref(false)
const channelLoading = ref(false)

// 组件引用
const authorizeRef = ref()
const modelSaveRef = ref()
const modelSelectRef = ref()
const modelSettingRef = ref()
const agentListDrawerRef = ref()

// 计算属性
const providerGroupList = computed<ProviderGroup[]>(() => {
  const list = [...authProviders.value, ...agentProviders.value]
  return list.reduce((acc: ProviderGroup[], item) => {
    let group = acc.find(row => row.label === item.category)
    if (!group) {
      group = { label: item.category, children: [] }
      acc.push(group)
    }
    group.children.push(item)
    return acc
  }, [])
})

// API 调用函数
const loadProviderList = async () => {
  providerLoading.value = true
  try {
    const list = await providerApi.list()
    authProviders.value = authProviders.value.map(item => {
      const providerData = list.find((row: any) => item.id === row.provider_type)

      return providerData
        ? {
            ...item,
            ...providerData,
            connected: true,
            client_id: providerData.configs?.client_id || '',
            client_secret: providerData.configs?.client_secret || '',
          }
        : item
    })
  } finally {
    providerLoading.value = false
  }
}

const loadAllTotal = async () => {
  const { count = 0 } = await agentApi.list({
    params: { group_id: '-1', keyword: '', offset: 0, limit: 1 },
  })
  allTotal.value = count
}

const loadAgentListCount = async () => {
  loadAllTotal()
  const promises = agentProviders.value.map(async provider => {
    const { count = 0 } = await agentApi.list({
      params: { channel_types: provider.id.toString(), limit: 1 },
    })
    provider.agentTotal = count
    provider.channelLoading = false
  })
  await Promise.all(promises)
}

const loadModelList = async () => {
  channelLoading.value = true
  try {
    const list = await channelApi.list()
    channelList.value = list
      .filter((item: any) => getModelChannelTypes().includes(item.channel_type as any))
      .reduce((acc: ChannelGroup[], item: any) => {
        let group = acc.find(row => row.channelType === item.channel_type)
        if (!group) {
          const model = getModelByChannelType(item.channel_type)
          group = {
            label: item.label,
            icon: item.icon,
            channelType: item.channel_type,
            multiple: model.multiple,
            data: item,
            children: [],
          }
          acc.push(group)
        }
        group.children.push(item)
        return acc
      }, [])
  } finally {
    channelLoading.value = false
  }
}

// 事件处理函数
const handleProviderAuthorize = ({ data }: { data: ProviderOption }): void => {
  if ([PROVIDER_VALUE.COZE_CN, PROVIDER_VALUE.COZE_OSV].includes(data.id) && isInternalNetwork()) {
    TipConfirm({
      title: window.$t('local_config_limited_tip'),
      content: window.$t('local_config_limited_desc', { url: window.location.href }),
      confirmButtonText: window.$t('know_it'),
      showCancelButton: false,
    }).open()
    return
  }
  if (data.auth) {
    authorizeRef.value.open({ data })
  } else {
    agentListDrawerRef.value.open({ data, type: data.id })
  }
}

const handleAgentAdd = ({ data }: { data: ProviderOption }) => {
  agentListDrawerRef.value.create({ data, type: data.id })
}

const handleProviderDelete = async ({ data }: { data: ProviderOption }) => {
  if (!data.provider_id) return

  await ElMessageBox.confirm(window.$t('module.platform_delete_confirm'))
  await providerApi.delete({ data: { provider_id: data.provider_id } })
  ElMessage.success(window.$t('action_delete_success'))
  setTimeout(() => {
    authProviders.value = getProvidersByAuth(true).map(createProviderOption)
    loadProviderList()
  }, 1000)
}

const handleModelSelect = () => modelSelectRef.value.open()

const handleModelAdd = (data: ModelConfig) => {
  modelSaveRef.value.open({ channel_type: data.channelType })
}

const handleModelEdit = (data: any) => modelSaveRef.value.open(data)

const handleModelDelete = async (data: any, model: any) => {
  await ElMessageBox.confirm(window.$t('module.platform_model_delete_confirm'))
  const isChildRemove = model && data.modelOptions.length > 1

  if (isChildRemove) {
    await channelApi.save({
      data: {
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
      },
    })
  } else {
    await channelApi.delete({ data: { channel_id: data.channel_id } })
  }

  ElMessage.success(window.$t('action_delete_success'))
  loadModelList()
}

const onModelEdit = ({ data, parentData }: { data: any; parentData: any }) => {
  modelSettingRef.value.open({ data: { ...parentData, ...data, id: data.value } })
}

const onAgentListChange = ({ data, count }: { data: ProviderOption; count: number }) => {
  const provider = agentProviders.value.find(item => item.id === data.id)
  if (provider) provider.agentTotal = count
  loadAgentListCount()
}

// 初始化
const refresh = () => {
  loadModelList()
  // #ifndef KM
  loadProviderList()
  loadAgentListCount()
  // #endif
}

onMounted(() => refresh())
</script>

<style scoped>
::v-deep(.el-collapse-item__arrow) {
  margin-left: 6px;
}
</style>
