import { lib_host } from './config'

/**
 * 库项类型定义
 * 每个库包含基本信息和动态加载的 Promise
 */
type LibItem = {
  /** 库的唯一标识符 */
  id: string
  /** 库的 JavaScript 文件路径 */
  src: string
  /** 库加载完成后的回调函数，用于执行额外的初始化操作 */
  callback: () => void
  /** 库加载的 Promise，用于缓存加载状态 */
  _promise?: Promise<void>
}

/**
 * 支持的库配置
 * 使用 Record 类型确保类型安全
 */
const libs: Record<string, LibItem> = {
  /** Vditor 富文本编辑器库 */
  vditor: {
    id: 'vditor-lib',
    src: `${lib_host}/js/vditor/dist/index.min.js`,
    callback() {
      // 动态加载 Vditor 的 CSS 样式文件
      const css = document.createElement('link')
      css.rel = 'stylesheet'
      css.href = `${lib_host}/js/vditor/dist/index.css`
      document.head.appendChild(css)
    },
  },
  /** UEditor 富文本编辑器库 */
  ueditor: {
    id: 'ueditor-lib',
    src: `${lib_host}/js/UEditor/ueditor.all.min.js`,
    callback() {
      // 动态加载 UEditor 的配置文件
      const script = document.createElement('script')
      script.src = `${lib_host}/js/UEditor/ueditor.config.js`
      script.id = 'ueditor-config'
      document.head.appendChild(script)
    },
  },
} as const

/** 支持的库名称类型 */
type LibName = keyof typeof libs

/** 支持的库名称 */
export const LIB_NAME = Object.keys(libs) as LibName[]

/**
 * 动态加载外部 JavaScript 库
 * 支持缓存机制，避免重复加载
 *
 * @param name 库名称，必须是预定义的库之一
 * @returns Promise<void> 加载完成后的 Promise
 *
 * @example
 * ```typescript
 * // 加载 Vditor 编辑器
 * await loadLib('vditor')
 *
 * // 加载 UEditor 编辑器
 * await loadLib('ueditor')
 * ```
 */
export default (name: LibName): Promise<void> => {
  // 检查库是否存在
  if (!libs[name]) return Promise.reject(new Error(`Library ${name} not found`))

  // 如果库已经加载过，直接返回缓存的 Promise
  if (!libs[name]._promise) {
    libs[name]._promise = new Promise<void>((resolve, reject) => {
      const { src, id, callback } = libs[name]

      // 创建 script 标签
      const script = document.createElement('script')
      script.src = src
      script.id = id

      // 处理 ES 模块
      if (src.endsWith('.mjs')) script.type = 'module'

      // 加载成功回调
      script.onload = () => {
        // 执行库特定的初始化回调
        if (callback) callback()

        // 延迟 100ms 确保库完全初始化
        setTimeout(() => {
          resolve()
        }, 100)
      }

      // 加载失败回调
      script.onerror = () => {
        reject(new Error(`Failed to load library ${name}`))
      }

      // 将 script 标签添加到页面
      document.body.appendChild(script)
    })
  }

  // 返回加载 Promise（使用非空断言，因为此时 _promise 一定存在）
  return libs[name]._promise!
}
