import channelApi, { transformSelectData } from '@/api/modules/channel/index'

export const loadModels = (type?: string): Promise<any[]> => {
  return channelApi.listv2().then(res => {
    const modelList = res.map(item => transformSelectData(item, type))
    return modelList
  })
}
