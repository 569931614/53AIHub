import service from '../config'
import { handleError } from '../errorHandler'

import type { WebsiteType } from '@/constants/enterprise'

export const saasApi = {
  product: {
    list() {
      return service.get(`/api/saas/products`).catch(handleError)
    },
    find(version: WebsiteType) {
      return service.get(`/api/saas/products/${version}`).catch(handleError)
    },
  },
  wechat_login(params: { unionid: string; from?: string }) {
    return service.get('/api/saas/wechat/user', { params }).catch(handleError)
  },
  bind_wechat(data: {
    mobile?: string
    verify_code?: string
    openid: string
    unionid?: string
    nickname?: string
    from?: string
  }) {
    let api_url = '/api/saas/wechat/bind'
    if (data.mobile) api_url = '/api/saas/wechat/user'
    return service.post(api_url, data).catch(handleError)
  },
}

export default saasApi
