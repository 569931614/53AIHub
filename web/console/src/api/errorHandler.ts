import { ElMessage } from 'element-plus'

import type { ResponseStatus } from './code'
import {
  ERROR_MESSAGES,
  RESPONSE_CODE_MESSAGE_MAP,
  RESPONSE_DATA_MESSAGE_MAP,
  RESPONSE_MESSAGE_MAP,
  ResponseCode,
} from './code'

// 定义错误响应接口
interface ErrorResponse {
  status?: ResponseStatus
  response?: {
    status?: ResponseStatus
    data?: {
      code?: ResponseCode
      message?: string
      [key: string]: unknown
    }
  }
  message?: string
}

// 消息缓存，用于防止重复显示
const messageCache = new Map<string, number>()
const messageTimers = new Map<string, NodeJS.Timeout>() // 存储定时器ID
const CACHE_DURATION = 3000 // 3秒内相同消息不重复显示

// 显示消息（带去重功能）
function showMessage(message: string) {
  const now = Date.now()

  // 检查是否在缓存期内
  if (messageCache.has(message)) {
    const lastShowTime = messageCache.get(message)!
    if (now - lastShowTime < CACHE_DURATION) {
      return // 跳过显示
    }

    // 清除之前的定时器
    const existingTimer = messageTimers.get(message)
    if (existingTimer) {
      clearTimeout(existingTimer)
    }
  }

  // 显示消息并更新缓存
  ElMessage.warning(message)
  messageCache.set(message, now)

  // 设置新的定时器并存储ID
  const timer = setTimeout(() => {
    messageCache.delete(message)
    messageTimers.delete(message)
  }, CACHE_DURATION)

  messageTimers.set(message, timer)
}

// 统一错误处理
export function handleError(error: ErrorResponse): Promise<never> {
  const status = error.response?.status || 500
  let resData = error.response?.data
  try {
    if (resData && typeof resData === 'string') {
      resData = JSON.parse(resData)
    }
  } catch (error) {
    resData = {}
  }
  const code = resData?.code
  const data = resData?.data
  let message = resData?.message

  const messageMatch = RESPONSE_MESSAGE_MAP.get(message || '')
  if (messageMatch) {
    if (messageMatch === 'not_tip') message = ''
    else message = window.$t(messageMatch)
  } else {
    // 优化消息获取逻辑
    // 先看错误集里有没有转换的信息提示，
    // 然后再是错误信息显示，
    // 如果没有则Code相关的错误提示，最后是错误信
    message =
      (data !== undefined && RESPONSE_DATA_MESSAGE_MAP.get(data)
        ? window.$t(RESPONSE_DATA_MESSAGE_MAP.get(data))
        : '') ||
      message ||
      error.message ||
      (code !== undefined && RESPONSE_CODE_MESSAGE_MAP.get(code)
        ? window.$t(RESPONSE_CODE_MESSAGE_MAP.get(code))
        : '') ||
      (ERROR_MESSAGES.get(status) ? window.$t(ERROR_MESSAGES.get(status)) : '') ||
      error.message ||
      window.$t('response_message.unknown_error')
  }

  // 使用带去重功能的消息显示
  if (message) {
    showMessage(message)
  }

  if (code === ResponseCode.TOKEN_EXPIRED_ERROR) {
    localStorage.removeItem('access_token')
    window.location.reload(true)
  }
  return Promise.reject(error)
}

// 清理所有消息缓存和定时器（可选，用于应用退出时清理）
export function clearMessageCache(): void {
  // 清理所有定时器
  for (const timer of messageTimers.values()) {
    clearTimeout(timer)
  }

  // 清空缓存
  messageCache.clear()
  messageTimers.clear()
}
