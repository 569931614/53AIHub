import { defineStore } from 'pinia'
import { useEnterpriseStore } from './enterprise'

import { domainApi } from '@/api/modules/domain'
import {
  transformDomainList,
  formatDomain,
  validateIndependentConfig,
  getDefaultExclusiveDomain,
  getDefaultIndependentDomain,
} from '@/api/modules/domain/transform'
import type {
  DomainListResponse,
  ExclusiveDomainData,
  IndependentDomainData,
  DomainInfo,
} from '@/api/modules/domain/types'

export const useDomainStore = defineStore('domain-store', {
  state: () => ({
    domainList: null as DomainListResponse | null,
    loading: false,
  }),

  getters: {
    exclusiveDomains: state => state.domainList?.exclusive_domains || [],
    independentDomains: state => state.domainList?.independent_domains || [],
    totalDomains: state => {
      const exclusive = state.domainList?.exclusive_domains?.length || 0
      const independent = state.domainList?.independent_domains?.length || 0
      return exclusive + independent
    },
  },

  actions: {
    /**
     * 加载域名列表数据
     */
    async loadListData(): Promise<DomainListResponse> {
      this.loading = true
      try {
        const rawData = await domainApi.list()
        const transformedData = transformDomainList(rawData)
        this.domainList = transformedData
        return transformedData
      } finally {
        this.loading = false
      }
    },

    /**
     * 保存专属域名
     */
    async saveExclusiveDomain(data: { domain_id?: number; domain: string }): Promise<any> {
      const formattedDomain = formatDomain(data.domain)
      const domainData: ExclusiveDomainData = { domain: formattedDomain }

      let result
      if (data.domain_id) {
        result = await domainApi.updateExclusive(data.domain_id, domainData)
      } else {
        result = await domainApi.createExclusive(domainData)
      }

      // 刷新企业信息和域名列表
      await this._refreshAfterSave()
      return result
    },

    /**
     * 保存独立域名
     */
    async saveIndependentDomain(
      data: { domain_id?: number } & IndependentDomainData
    ): Promise<any> {
      // 验证配置
      if (!validateIndependentConfig(data.config)) {
        throw new Error('域名配置验证失败')
      }

      const formattedData: IndependentDomainData = {
        domain: formatDomain(data.domain),
        config: data.config,
      }

      let result
      if (data.domain_id) {
        result = await domainApi.updateIndependent(data.domain_id, formattedData)
      } else {
        result = await domainApi.createIndependent(formattedData)
      }

      // 刷新企业信息和域名列表
      await this._refreshAfterSave()
      return result
    },

    /**
     * 删除独立域名
     */
    async deleteIndependentDomain(domainId: number): Promise<void> {
      await domainApi.deleteIndependent(domainId)

      // 刷新企业信息和域名列表
      await this._refreshAfterSave()
    },

    /**
     * 根据域名 ID 查找域名信息
     */
    findDomainById(domainId: number, type: 'exclusive' | 'independent'): DomainInfo | undefined {
      const domains = type === 'exclusive' ? this.exclusiveDomains : this.independentDomains
      return domains.find(domain => domain.domain_id === domainId)
    },

    /**
     * 检查域名是否已存在
     */
    isDomainExists(domain: string, excludeId?: number): boolean {
      const formattedDomain = formatDomain(domain)
      const allDomains = [...this.exclusiveDomains, ...this.independentDomains]

      return allDomains.some(
        d => formatDomain(d.domain) === formattedDomain && d.domain_id !== excludeId
      )
    },

    /**
     * 获取默认的专属域名数据
     */
    getDefaultExclusiveDomain,

    /**
     * 获取默认的独立域名数据
     */
    getDefaultIndependentDomain,

    /**
     * 保存后刷新相关数据
     */
    async _refreshAfterSave(): Promise<void> {
      const enterpriseStore = useEnterpriseStore()
      await Promise.all([enterpriseStore.loadSelfInfo(), this.loadListData()])
    },

    /**
     * 重置状态
     */
    resetState(): void {
      this.domainList = null
      this.loading = false
    },
  },
})
