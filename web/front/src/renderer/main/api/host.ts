/**
 * 格式化URL，确保返回完整的URL
 * @param url 需要格式化的URL
 * @returns 格式化后的完整URL
 */
const formatUrl = (url: string): string => {
  if (!url) return ''

  const { origin } = window.location

  // 如果已经是完整的URL，直接返回
  if (url.startsWith('http')) {
    return url
  }

  return origin
}

/**
 * 获取环境变量或window对象中的值
 * @param windowKey window对象中的键名
 * @param envKey 环境变量键名
 * @param defaultValue 默认值
 * @returns 配置值
 */
const getConfigValue = (
  windowKey: keyof Window,
  envKey: string,
  defaultValue: string = ''
): string => {
  return (window[windowKey] as string) || process.env[envKey] || defaultValue
}

// 导出配置常量
export const API_HOST = formatUrl(getConfigValue('api_host', 'VITE_GLOB_API_HOST'))
export const QYY_HOST = formatUrl(getConfigValue('qyy_host', 'VITE_GLOB_QYY_HOST'))
export const ADMIN_URL = getConfigValue('admin_url', 'VITE_GLOB_ADMIN_URL', '/console')

// 基于API_HOST的派生配置
export const IMG_HOST = `${API_HOST}/api/images`
export const LIB_HOST = `${API_HOST}/api/libs`
