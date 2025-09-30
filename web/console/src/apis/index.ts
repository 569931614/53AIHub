// import type { AxiosProgressEvent, GenericAbortSignal } from 'axios'
// import cache from '@/utils/cache'
import user from './modules/user'
import group from './modules/group'
import conversation from './modules/conversation'
import enterprise from './modules/enterprise'
import setting from './modules/setting'
import qyy from './modules/qyy'
import { isObject } from '@/utils'
import { del, get, patch, post, put } from '@/utils/request'

const api_config = {
  user,
  group,
  conversation,
  enterprise,
  setting,
  qyy,
}

const REQUEST_MAP = new Map([
  ['GET', get],
  ['POST', post],
  ['PATCH', patch],
  ['PUT', put],
  ['DELETE', del],
])

// export function fetchChatAPIProcess<T = any>(
//   params: {
//     data: {
//       conversationid: string
//       eid: number
//       question: string
//       modelid?: string
//       datasetid?: string
//       contentid?: number
//       characterid?: string
//       app_model_config_id?: string
//       othersaveparams: any
//     }
//     // prompt: string
//     // options?: { conversationId?: string; parentMessageId?: string }
//     signal?: GenericAbortSignal
//     onDownloadProgress?: (progressEvent: AxiosProgressEvent) => void
//   },
// ) {
//   // const settingStore = useSettingStore()
//   // const authStore = useAuthStore()

//   const data: Record<string, any> = {
//     conversationid: params.data.conversationid,
//     eid: params.data.eid,
//     question: params.data.question,
//     modelid: params.data.modelid,
//     datasetid: params.data.datasetid,
//     characterid: params.data.characterid,
//     app_model_config_id: params.data.app_model_config_id,
//     othersaveparams: params.data.othersaveparams,
//   }
//   if (typeof params.data.contentid !== 'undefined')
//     data.contentid = params.data.contentid

//   // if (authStore.isChatGPTAPI) {
//   //   data = {
//   //     ...data,
//   //     systemMessage: settingStore.systemMessage,
//   //     temperature: settingStore.temperature,
//   //     top_p: settingStore.top_p,
//   //   }
//   // }

//   return post<T>({
//     url: `${apiHost}/chat/send`,
//     data,
//     signal: params.signal,
//     onDownloadProgress: params.onDownloadProgress,
//   })
// }

// export function downloadUrl<T = any>(url: string, data?: any) {
//   return get<T>({
//     url,
//     data,
//     responseType: 'blob',
//   })
// }

const findEach = (source: any = {}, url = '') => {
  Object.keys(source).forEach(key => {
    const data = source[key]
    if (data.method && !isObject(data.method)) {
      let method = data.method || ''
      Object.keys(data).forEach(data_key => {
        if (!Object.keys(data).includes('method'))
          findEach(data[data_key], `${url}/${key}/${data_key}`)
      })
      source[key] = (config: any = {}) => {
        method = (config.method || method).toUpperCase()
        if (['DELETE', 'DEL'].includes(method)) method = 'DELETE'
        const _request = REQUEST_MAP.get(method) || REQUEST_MAP.get('GET')
        return _request({
          ...data,
          method,
          url: data.url || `${url}/${key}`,
          ...config,
        })
      }
    } else {
      findEach(data, `${url}/${key}`)
    }
  })
}
findEach(api_config)

export default {
  ...api_config,
}
