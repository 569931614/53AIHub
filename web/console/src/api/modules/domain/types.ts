import type { IndependentResolveType, IndependentSslCertType } from '@/constants/domain'

/**
 * 域名基础信息
 */
export interface DomainInfo {
  domain_id?: number
  domain: string
  status?: number
  created_time?: string
  updated_time?: string
}

/**
 * 独立域名配置
 */
export interface IndependentDomainConfig {
  resolve_type: IndependentResolveType
  enable_https: boolean
  force_https: boolean
  ssl_cert_type: IndependentSslCertType
  ssl_certificate: string
  ssl_private_key: string
  use_subdir: boolean
  subdir: string
}

/**
 * 独立域名数据
 */
export interface IndependentDomainData {
  domain: string
  config: IndependentDomainConfig
}

/**
 * 专属域名数据
 */
export interface ExclusiveDomainData {
  domain: string
}

/**
 * 域名列表响应
 */
export interface DomainListResponse {
  exclusive_domains?: DomainInfo[]
  independent_domains?: DomainInfo[]
}

/**
 * 原始域名列表响应（从 API 返回）
 */
export type RawDomainListResponse = any
