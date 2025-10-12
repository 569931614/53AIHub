import service from '../config'
import { handleError } from '../errorHandler'

export type DefaultLinkItem = {
  name: string
  logo: string
  url: string
  description: string
  sort: number
}

export type DefaultLinkReuest = {
  links: {
    ai_link: DefaultLinkItem[]
    delete: boolean
  }
}

export const settingApi = {
  list() {
    return service.get('/api/settings').catch(handleError)
  },
  get(key: string) {
    return service.get(`/api/settings/key/${key}`).catch(handleError)
  },
  detail(group_name: string) {
    return service.get(`/api/settings/group/${group_name}`).catch(handleError)
  },
  create(data: { key: string; value: string }) {
    return service.post('/api/settings', data).catch(handleError)
  },
  update(setting_id: number, data: { key: string; value: string }) {
    return service.put(`/api/settings/${setting_id}`, data).catch(handleError)
  },
  default_links: {
    list() {
      return service.get('/api/settings/default_links').catch(handleError)
    },
    save(data: DefaultLinkReuest) {
      return service.post('/api/settings/default_links', data).catch(handleError)
    },
  },
}

export default settingApi
