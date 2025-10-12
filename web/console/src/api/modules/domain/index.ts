import service from '@/api/config'
import { handleError } from '@/api/errorHandler'

import type { RawDomainListResponse, ExclusiveDomainData, IndependentDomainData } from './types'

/**
 * 域名管理 API
 */
export const domainApi = {
  /**
   * 获取域名列表
   */
  list(): Promise<RawDomainListResponse> {
    return service
      .get('/api/saas/domains')
      .then(res => res.data)
      .catch(handleError)
  },

  /**
   * 创建专属域名
   */
  createExclusive(data: ExclusiveDomainData) {
    return service.post('/api/saas/domains/exclusive', data).catch(handleError)
  },

  /**
   * 更新专属域名
   */
  updateExclusive(domainId: number, data: ExclusiveDomainData) {
    return service.put(`/api/saas/domains/exclusive/${domainId}`, data).catch(handleError)
  },

  /**
   * 创建独立域名
   */
  createIndependent(data: IndependentDomainData) {
    return service.post('/api/saas/domains/independent', data).catch(handleError)
  },

  /**
   * 更新独立域名
   */
  updateIndependent(domainId: number, data: IndependentDomainData) {
    return service.put(`/api/saas/domains/independent/${domainId}`, data).catch(handleError)
  },

  /**
   * 删除独立域名
   */
  deleteIndependent(domainId: number) {
    return service.delete(`/api/saas/domains/independent/${domainId}`).catch(handleError)
  },
}

export default domainApi
