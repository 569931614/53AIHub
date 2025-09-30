/**
 * 域名解析类型常量
 */
export const INDEPENDENT_RESOLVE_TYPE = {
  CNAME: 1,
  CUSTOM: 2,
} as const

export type IndependentResolveType =
  (typeof INDEPENDENT_RESOLVE_TYPE)[keyof typeof INDEPENDENT_RESOLVE_TYPE]

/**
 * 独立域名 SSL 证书类型
 */
export const INDEPENDENT_SSL_CERT_TYPE = {
  '53AI': 1,
  CUSTOM: 2,
} as const

export type IndependentSslCertType =
  (typeof INDEPENDENT_SSL_CERT_TYPE)[keyof typeof INDEPENDENT_SSL_CERT_TYPE]

/**
 * 域名配置默认值
 */
export const DOMAIN_CONFIG = {
  DEFAULT_ENABLE_HTTPS: true,
  DEFAULT_FORCE_HTTPS: false,
  DEFAULT_USE_SUBDIR: false,
  DEFAULT_SUBDIR: '',
} as const

/**
 * 域名类型
 */
export const DOMAIN_TYPE = {
  EXCLUSIVE: 'exclusive',
  INDEPENDENT: 'independent',
} as const

export type DomainType = (typeof DOMAIN_TYPE)[keyof typeof DOMAIN_TYPE]
