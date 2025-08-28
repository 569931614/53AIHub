import service from '@/api/config'
import { handleError } from '@/api/errorHandler'
import { getModelByChannelType, ModelValue } from '@/constants/platform/config'

export interface RawModelOption {
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

export interface ModelOption extends RawModelOption {
  icon: string
  categories: Array<{
    icon: string
    model_type: number
    model_type_name: string
    models: Array<{ icon: string; model_id: string; model_name: string }>
  }>
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
  type: ModelValue
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
  name: string
  label: string
  custom_config: Record<string, any>
  config: Record<string, any>
  channel_type: number
  models: string[]
  group: Array<{
    type: number
    typeName: string
    options: Array<{
      value: string
      type: number
      typeName: string
      label: string
      icon: string
    }>
  }>
  aliasMap: Record<string, string>
  options: Array<{
    value: string
    label: string
    icon: string
  }>
}

export const getModelIcon = (value: string) => {
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

  return icon ? window.$getRealPath({ url: `/images/platform/${icon}.png` }) : ''
}

export const transformChannelData = (data: RawChannelItem): ChannelItem => {
  const model = getModelByChannelType(data.type) || {}

  const config = (typeof data.config === 'string' ? JSON.parse(data.config) : data.config) || {}
  const custom_config =
    typeof data.custom_config === 'string' && data.custom_config
      ? JSON.parse(data.custom_config)
      : data.custom_config || {}

  const models = typeof data.models === 'string' ? data.models.split(',') : data.models || []
  const model_alias_map = (config.model_alias_map || {}) as Record<string, string>

  const options = models.map(value => {
    const type = custom_config[value] || 0
    return {
      value,
      type,
      typeName: type === '1' ? '推理' : type === '2' ? '嵌入' : '重排序',
      label: model_alias_map[value] || value,
      icon:
        getModelIcon(value) || window.$getRealPath({ url: `/images/platform/${model.icon}.png` }),
    }
  })
  const group = options.reduce((acc, item) => {
    const data = acc.find(row => row.type === item.type)
    if (data) {
      data.options.push(item)
    } else {
      acc.push({
        type: item.type,
        typeName: item.typeName,
        options: [item],
      })
    }
    return acc
  }, [])
  return {
    ...data,
    icon: window.$getRealPath({ url: `/images/platform/${model.icon}.png` }),
    // name: model.name,
    label: window.$t(`provider_platform.${model.name}`),
    channel_type: data.type,
    custom_config,
    config,
    models,
    aliasMap: model_alias_map,
    group,
  }
}

export const transformSelectData = (data: RawChannelItem, type?: string): ChannelItem => {
  const result = transformChannelData(data)
  let models = result.models
  if (type) {
    models = models.filter(item => result.custom_config[item] === type)
  }
  return {
    ...result,
    options: models.map(value => {
      return {
        value: `${data.channel_id}_${value}`,
        label: result.aliasMap[value] || value,
        icon: getModelIcon(value) || result.icon,
      }
    }),
  }
}

export const transformModelList = (data: RawModelOption[]): ModelOption[] => {
  return data.map(item => {
    return {
      ...item,
      icon: window.$getRealPath({ url: `/images/platform/${item.platform_id}.png` }),
      channel_type: item.channel_type,
      categories: item.categories.map(cate => {
        return {
          ...cate,
          model_type_name:
            cate.model_type === 1
              ? '逻辑推理模型'
              : cate.model_type === 2
                ? '向量嵌入模型'
                : '重排序模型',
          models: cate.models.map(model => {
            return {
              ...model,
              icon: getModelIcon(model.model_id),
            }
          }),
        }
      }),
    }
  })
}

export const channelApi = {
  models: {
    config(): Promise<RawModelOption[]> {
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

  delete(channel_id: number) {
    return service.delete(`/api/channels/${channel_id}`).catch(handleError)
  },
}

export default channelApi
