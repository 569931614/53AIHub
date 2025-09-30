import {
  DOMAIN_CONFIG,
  INDEPENDENT_RESOLVE_TYPE,
  INDEPENDENT_SSL_CERT_TYPE,
} from '@/constants/domain'
import type {
  DomainListResponse,
  RawDomainListResponse,
  IndependentDomainConfig,
  IndependentDomainData,
  ExclusiveDomainData,
} from './types'

/**
 * 获取默认的独立域名配置
 */
export const getDefaultIndependentConfig = (): IndependentDomainConfig => ({
  resolve_type: INDEPENDENT_RESOLVE_TYPE.CNAME,
  enable_https: DOMAIN_CONFIG.DEFAULT_ENABLE_HTTPS,
  force_https: DOMAIN_CONFIG.DEFAULT_FORCE_HTTPS,
  ssl_cert_type: INDEPENDENT_SSL_CERT_TYPE['53AI'],
  ssl_certificate: '',
  ssl_private_key: '',
  use_subdir: DOMAIN_CONFIG.DEFAULT_USE_SUBDIR,
  subdir: DOMAIN_CONFIG.DEFAULT_SUBDIR,
})

/**
 * 获取默认的独立域名数据
 */
export const getDefaultIndependentDomain = (): IndependentDomainData => ({
  domain: '',
  config: getDefaultIndependentConfig(),
})

/**
 * 获取默认的专属域名数据
 */
export const getDefaultExclusiveDomain = (): ExclusiveDomainData => ({
  domain: '',
})

/**
 * 转换域名列表响应数据
 */
export function transformDomainList(rawData: RawDomainListResponse): DomainListResponse {
  try {
    return {
      exclusive_domains: rawData.exclusive_domains || [],
      independent_domains: rawData.independent_domains || [],
    }
  } catch (error) {
    console.error('转换域名列表数据失败:', error)
    return {
      exclusive_domains: [],
      independent_domains: [],
    }
  }
}

/**
 * 验证独立域名配置
 */
export function validateIndependentConfig(config: Partial<IndependentDomainConfig>): boolean {
  // 如果启用 HTTPS 且使用自定义证书，则必须提供证书和私钥
  if (config.enable_https && config.ssl_cert_type === INDEPENDENT_SSL_CERT_TYPE.CUSTOM) {
    return !!(config.ssl_certificate && config.ssl_private_key)
  }

  // 如果使用子目录，则必须提供子目录名称
  if (config.use_subdir) {
    return !!(config.subdir && config.subdir.trim())
  }

  return true
}

/**
 * 格式化域名（移除协议前缀等）
 */
export function formatDomain(domain: string): string {
  return domain
    .replace(/^https?:\/\//, '') // 移除协议前缀
    .replace(/\/$/, '') // 移除末尾斜杠
    .toLowerCase() // 转为小写
    .trim()
}
