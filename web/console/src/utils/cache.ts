/**
 * 缓存管理工具类
 * 提供内存缓存功能，支持过期时间和泛型数据类型
 */

/** 数据获取器类型定义 */
type Fetcher<T> = (() => Promise<T>) | (() => T) | T

/** 缓存项接口定义 */
interface CacheItem<T> {
  /** 缓存的数据 */
  data: T
  /** 过期时间戳 */
  expireTime: number
}

/**
 * 缓存管理器类
 * 使用单例模式，提供全局统一的缓存管理
 */
export class CacheManager {
  /** 单例实例 */
  // eslint-disable-next-line no-use-before-define
  private static instance: CacheManager | null = null

  /** 缓存数据存储 Map */
  private cacheMap: Map<string, CacheItem<unknown>> = new Map()

  /**
   * 获取单例实例
   * @returns CacheManager 实例
   */
  static getInstance(): CacheManager {
    if (!CacheManager.instance) {
      CacheManager.instance = new CacheManager()
    }
    return CacheManager.instance
  }

  /**
   * 判断值是否为 Promise
   * @param value 待检查的值
   * @returns 是否为 Promise
   */
  private isPromise<T>(value: unknown): value is Promise<T> {
    return value != null && typeof (value as { then?: unknown }).then === 'function'
  }

  /**
   * 设置缓存
   * @param key 缓存键
   * @param value 缓存值
   * @param expireMinutes 过期时间（分钟），默认 1 分钟
   */
  set<T>(key: string, value: T, expireMinutes = 1): void {
    const expireTime = Date.now() + expireMinutes * 60 * 1000
    this.cacheMap.set(key, {
      data: value,
      expireTime,
    })
  }

  /**
   * 获取缓存
   * @param key 缓存键
   * @returns 缓存值，如果不存在或已过期则返回 null
   */
  get<T>(key: string): T | null {
    const now = Date.now()
    const cacheItem = this.cacheMap.get(key)

    if (cacheItem && now < cacheItem.expireTime) {
      return cacheItem.data as T
    }

    // 清理过期缓存
    this.cacheMap.delete(key)
    return null
  }

  /**
   * 获取缓存或执行获取函数
   * 如果缓存存在且未过期，直接返回缓存值
   * 否则执行获取函数并缓存结果
   *
   * @param key 缓存键
   * @param fetcher 数据获取器，可以是函数、异步函数或直接的值
   * @param expireMinutes 过期时间（分钟），默认 2 分钟
   * @returns 缓存或获取的数据
   *
   * @example
   * ```typescript
   * // 缓存API请求结果
   * const userData = await cache.getOrFetch(
   *   'user:123',
   *   () => fetch('/api/user/123').then(res => res.json()),
   *   5 // 5分钟过期
   * )
   *
   * // 缓存计算结果
   * const computed = await cache.getOrFetch(
   *   'expensive-calc',
   *   () => expensiveCalculation(),
   *   30 // 30分钟过期
   * )
   * ```
   */
  async getOrFetch<T>(key: string, fetcher: Fetcher<T>, expireMinutes = 2): Promise<T> {
    // 检查缓存
    const cachedValue = this.get<T>(key)
    if (cachedValue !== null) return cachedValue

    // 处理不同类型的 fetcher
    let result: T
    if (typeof fetcher === 'function') {
      const fetchResult = (fetcher as () => T | Promise<T>)()
      if (this.isPromise<T>(fetchResult)) {
        result = await fetchResult
      } else {
        result = fetchResult
      }
    } else {
      result = fetcher
    }

    // 存储结果并返回
    this.set(key, result, expireMinutes)
    return result
  }

  /**
   * 删除指定缓存
   * @param key 缓存键
   */
  delete(key: string): void {
    this.cacheMap.delete(key)
  }

  /**
   * 清空所有缓存
   */
  clear(): void {
    this.cacheMap.clear()
  }

  /**
   * 获取缓存状态信息
   * @returns 缓存统计信息
   */
  getStats(): {
    /** 缓存项总数 */
    total: number
    /** 有效缓存项数量 */
    valid: number
    /** 过期缓存项数量 */
    expired: number
  } {
    const now = Date.now()
    let validCount = 0
    let expiredCount = 0

    for (const [, item] of this.cacheMap.entries()) {
      if (now < item.expireTime) {
        validCount++
      } else {
        expiredCount++
      }
    }

    return {
      total: this.cacheMap.size,
      valid: validCount,
      expired: expiredCount,
    }
  }

  /**
   * 清理所有过期缓存
   * @returns 清理的缓存项数量
   */
  cleanExpired(): number {
    const now = Date.now()
    let cleanedCount = 0

    for (const [key, item] of this.cacheMap.entries()) {
      if (now >= item.expireTime) {
        this.cacheMap.delete(key)
        cleanedCount++
      }
    }

    return cleanedCount
  }
}

/** 默认导出缓存管理器单例实例 */
export default CacheManager.getInstance()
