import { ElMessageBox } from 'element-plus'
import service from '../config'
import { handleError } from '../errorHandler'
import {
  CHANNEL_TYPE,
  CHANNEL_TYPE_ICON_MAP,
  CHANNEL_TYPE_LABEL_MAP,
  CHANNEL_TYPE_OPTIONS,
  CHANNEL_TYPE_VALUE_MAP,
  type ChannelType,
} from '@/constants/platform'

export {
  CHANNEL_TYPE,
  CHANNEL_TYPE_LABEL_MAP,
  CHANNEL_TYPE_ICON_MAP,
  CHANNEL_TYPE_VALUE_MAP,
  CHANNEL_TYPE_OPTIONS,
  type ChannelType,
}

export interface ModelOption {
  categories: Array<{
    model_count: number
    model_type: number
    models: Array<{
      model_id: string
      model_name: string
    }>
  }>
  platform_id: string
  platform_name: string
  channel_type: number
  can_multiple: boolean
}

export interface ChannelRequestData {
  base_url: string
  config: string
  key: string
  custom_config: string
  model_mapping: string
  model_type: number
  models: string
  name: string
  other: string
  priority: number
  provider_id: number
  type: number
  weight: number
}

export interface RawChannelItem {
  channel_id: number
  eid: number
  type: number
  key: string
  weight: number
  name: string
  models: string
  config: string
  custom_config: string
  other: string
  model_mapping: string
  priority: number
  base_url: string
  used_quota: number
  status: number
  balance: number
  balance_updated_time: number
  test_time: number
  response_time: number
  provider_id: number
  created_time: number
  updated_time: number
}

export interface ChannelItem extends RawChannelItem {
  icon: string
  custom_config: Record<string, any>
  config: Record<string, any>
  channel_id: number
  models: string[]
}

interface ModelConfig {
  id: string
  vision?: boolean
  provider: ChannelType
}

interface ChannelData {
  value?: number
  channel_id?: number
  channel_type?: number | string
  type?: string | number
  name?: string
  id?: string
  label?: string
  icon?: string
  key?: string
  api_key?: string
  base_url?: string
  config?: any
  models?: string[] | string
  model_options?: Array<{
    value: string
    label: string
    icon: string
    vision: boolean
  }>
  modelOptions?: Array<{
    value: string
    label: string
    icon: string
    vision: boolean
  }>
  owned_by?: string
  priority?: number
  weight?: number
  model_mapping?: string
}

interface ModelItem {
  id: string
  icon?: string
  label?: string
  owned_by?: string
  value?: string
}

export const OPEN_AI_MODEL_LIST: ModelConfig[] = [
  // 'gpt-4',
  { provider: CHANNEL_TYPE.OPENAI, id: 'gpt-4.1', vision: true },
  { provider: CHANNEL_TYPE.OPENAI, id: 'gpt-4.1-mini', vision: true },
  { provider: CHANNEL_TYPE.OPENAI, id: 'gpt-4o', vision: true },
  { provider: CHANNEL_TYPE.OPENAI, id: 'gpt-4o-mini', vision: true },
  { provider: CHANNEL_TYPE.OPENAI, id: 'o1', vision: true },
  { provider: CHANNEL_TYPE.OPENAI, id: 'o1-mini' },
]

export const SILICON_FLOW_MODEL_LIST: ModelConfig[] = [
  { provider: CHANNEL_TYPE.SILICONFLOW, id: 'Qwen/Qwen3-8B', vision: false },
  { provider: CHANNEL_TYPE.SILICONFLOW, id: 'Qwen/Qwen2.5-7B-Instruct' },
  { provider: CHANNEL_TYPE.SILICONFLOW, id: 'deepseek-ai/DeepSeek-R1' },
  { provider: CHANNEL_TYPE.SILICONFLOW, id: 'deepseek-ai/DeepSeek-V3', vision: false },
  { provider: CHANNEL_TYPE.SILICONFLOW, id: 'THUDM/glm-4-9b-chat' },
]

export const DEEP_SEEK_MODEL_LIST: ModelConfig[] = [
  { provider: CHANNEL_TYPE.DEEPSEEK, id: 'deepseek-chat', vision: false },
  { provider: CHANNEL_TYPE.DEEPSEEK, id: 'deepseek-reasoner' },
]

export const ALL_MODEL_LIST = [
  ...OPEN_AI_MODEL_LIST,
  ...SILICON_FLOW_MODEL_LIST,
  ...DEEP_SEEK_MODEL_LIST,
]

export const MODEL_ALIAS_MAP = new Map([
  ['deepseek-chat', 'DeepSeek-V3'],
  ['deepseek-reasoner', 'DeepSeek-R1'],
])

export const getModelIcon = ({ value = '' }: { value: string }) => {
  let icon = ''
  if (/deepseek/i.test(value)) icon = 'deepseek'
  else if (/tongyi|qwen/i.test(value)) icon = 'tongyi'
  else if (/thudm/i.test(value)) icon = 'zhipu'
  else if (/ai\/yi/i.test(value)) icon = 'yi'
  else if (/internlm/i.test(value)) icon = 'internlm'
  else if (/baai/i.test(value)) icon = 'baai'
  else if (/google/i.test(value)) icon = 'google'
  else if (/mistralai/i.test(value)) icon = 'mistralai'
  else if (/llama/i.test(value)) icon = 'llama'

  return icon
}

export const getFormatChannelData = (data: ChannelData = {}) => {
  const channelId = data.channel_id || 0
  data.value = channelId
  data.channel_id = channelId
  data.channel_type = data.channel_type || data.type || 0
  data.type = data.name || data.id || ''
  data.channel_type =
    data.channel_type || Number(CHANNEL_TYPE_VALUE_MAP.get(data.type || '') || 0) || 0
  data.name = data.name || ''
  data.label = CHANNEL_TYPE_LABEL_MAP.get(data.channel_type) || data.name || ''
  data.icon = CHANNEL_TYPE_ICON_MAP.get(data.channel_type) || ''
  data.api_key = data.key || ''
  data.base_url = data.base_url || ''
  try {
    data.config = typeof data.config === 'string' ? JSON.parse(data.config) : data.config || {}
  } catch (error) {
    data.config = {}
  }
  data.models = typeof data.models === 'string' ? data.models.split(',') : data.models || []
  const model_alias_map = (data.config?.model_alias_map || {}) as Record<string, string>

  data.model_options = (Array.isArray(data.models) ? data.models : []).map(value => ({
    value,
    label: model_alias_map[value] || MODEL_ALIAS_MAP.get(value) || value,
    icon: getModelIcon({ value }) || data.icon || '',
    vision: ALL_MODEL_LIST.find(item => item.id === value)?.vision || false,
  }))
  data.modelOptions = (Array.isArray(data.models) ? data.models : []).map(value => ({
    value,
    label: model_alias_map[value] || MODEL_ALIAS_MAP.get(value) || value,
    icon: getModelIcon({ value }) || data.icon || '',
    vision: ALL_MODEL_LIST.find(item => item.id === value)?.vision || false,
  }))

  return data
}

export const channelApi = {
  async list(
    {
      params = {},
    }: {
      params: {
        provider_id?: number
        channel_types?: string
      }
    } = { params: {} }
  ) {
    params = JSON.parse(JSON.stringify(params))
    const { data = [] } = await service.get('/api/channels', { params }).catch(handleError)
    const list = data.map((item: ChannelData = {}) => getFormatChannelData(item))
    return list
  },

  async save({ data = {} }: { data: ChannelData } = { data: {} }) {
    data = {
      channel_id: 0,
      key: '',
      base_url: '',
      config: {},
      model_mapping: '',
      models: '',
      name: '',
      priority: 0,
      type: 0,
      weight: 0,
      ...data,
    }
    if (!data.type) {
      data.type = CHANNEL_TYPE_VALUE_MAP.get(data.name || '') || 0
      if (!data.type) return Promise.reject(new Error('Invalid channel type'))
    }
    const { channel_id } = data
    delete data.channel_id
    if (typeof data.config === 'object') {
      if (data.config.model_alias_map) {
        const models = typeof data.models === 'string' ? data.models.split(',') : data.models || []
        Object.keys(data.config.model_alias_map).forEach(key => {
          if (!models.includes(key)) delete data.config.model_alias_map[key]
        })
        if (!Object.keys(data.config.model_alias_map).length) delete data.config.model_alias_map
      }
      data.config = JSON.stringify(data.config)
    }
    if (Array.isArray(data.models)) data.models = data.models.join(',')
    const { data: resultData = {} } = await service[channel_id ? 'put' : 'post'](
      `/api/channels${channel_id ? `/${channel_id}` : ''}`,
      data
    ).catch(handleError)
    return getFormatChannelData(resultData)
  },

  async delete({ data: { channel_id } }: { data: { channel_id: number } }) {
    return service.delete(`/api/channels/${channel_id}`).catch(handleError)
  },

  async modelList({ params: { owner = 'all' } }: { params: { owner: string } }) {
    const { data: { models = [] } = {} } = await service
      .get('/api/channels/models', { params: { owner } })
      .catch(handleError)
    let list = models.map((item: ModelItem) => {
      const id = item.id || item.value || ''
      return {
        ...item,
        id,
        icon: item.icon || getModelIcon({ value: id }) || '',
        label: item.label || MODEL_ALIAS_MAP.get(id) || '',
      }
    })
    // #ifndef KM
    list = list.filter(
      item =>
        ![CHANNEL_TYPE.OPENAI, CHANNEL_TYPE.SILICONFLOW, CHANNEL_TYPE.DEEPSEEK].includes(
          item.owned_by as ChannelType
        )
    )
    const manual_list = [
      ...OPEN_AI_MODEL_LIST.map(model => ({
        value: model.id,
        id: model.id,
        label: model.id,
        icon: CHANNEL_TYPE.OPENAI,
        owned_by: CHANNEL_TYPE.OPENAI,
        vision: model.vision,
      })),
      ...SILICON_FLOW_MODEL_LIST.map(model => ({
        value: model.id,
        id: model.id,
        label: model.id,
        icon: getModelIcon({ value: model.id }) || CHANNEL_TYPE.SILICONFLOW || '',
        owned_by: CHANNEL_TYPE.SILICONFLOW,
        vision: model.vision,
      })),
      ...DEEP_SEEK_MODEL_LIST.map(model => ({
        value: model.id,
        id: model.id,
        label: MODEL_ALIAS_MAP.get(model.id) || model.id,
        icon: CHANNEL_TYPE.DEEPSEEK,
        owned_by: CHANNEL_TYPE.DEEPSEEK,
        vision: model.vision,
      })),
    ]
    list.push(...manual_list)
    // #endif
    return list
  },

  async cozeWorkspaceList() {
    const { data = [] } = await service.get('/api/coze/workspaces').catch((error: any = {}) => {
      const { response: { data: { code } = {} } = {} } = error
      if ([13].includes(+code)) {
        ElMessageBox.confirm(
          (window as any).$t('coze_cn_not_auth', (window as any).$t('tip')),
          (window as any).$t('tip'),
          {
            confirmButtonText: (window as any).$t('action_go'),
            type: 'warning',
            closeOnClickModal: false,
            showClose: false,
            showCancelButton: false,
          }
        ).then(() => {
          location.hash = '#/platform'
        })
        return Promise.reject(error)
      }
      return handleError(error)
    })
    return data
  },

  async cozeBotList(
    {
      params: { workspace_id } = { workspace_id: 0 },
    }: {
      params: { workspace_id: number }
    } = { params: { workspace_id: 0 } }
  ) {
    const { data = [] } = await service
      .get(`/api/coze/workspaces/${workspace_id}/bots`)
      .catch(handleError)
    return data
  },

  async appBuilderBotList() {
    const { data = [] } = await service.get('/api/appbuilder/bots').catch(handleError)
    return data
  },

  async chat53aiAgentList() {
    const { data = [] } = await service.get('/api/53ai/bots').catch(handleError)
    return data
  },
  async chat53aiWorkflowList() {
    const { data = [] } = await service.get('/api/53ai/workflows').catch(handleError)
    return data
  },
  models: {
    config(): Promise<ModelOption[]> {
      return service
        .get('/api/channels/km/models')
        .then(res => res.data.platforms)
        .catch(handleError)
    },
  },
  listv2(): Promise<RawChannelItem[]> {
    return service
      .get('/api/channels')
      .then(res => res.data)
      .catch(handleError)
  },
  create(data: ChannelRequestData) {
    return service.post('/api/channels', data).catch(handleError)
  },
  update(channel_id: number, data: ChannelRequestData) {
    return service.put(`/api/channels/${channel_id}`, data).catch(handleError)
  },
}

export default channelApi
