/**
 * 导航类型枚举
 */
export const NAVIGATION_TYPE = {
  /** 系统导航 */
  SYSTEM: 1,
  /** 外部链接 */
  EXTERNAL: 2,
  /** 自定义页面 */
  CUSTOM: 3,
} as const

export type NavigationType = (typeof NAVIGATION_TYPE)[keyof typeof NAVIGATION_TYPE]

/**
 * 导航打开方式枚举
 */
export const NAVIGATION_TARGET = {
  /** 当前窗口打开 */
  SELF: 1,
  /** 新窗口打开 */
  BLANK: 2,
} as const

export type NavigationTarget = (typeof NAVIGATION_TARGET)[keyof typeof NAVIGATION_TARGET]

/**
 * 导航类型标签映射
 */
export const NAVIGATION_TYPE_LABEL_MAP = new Map<NavigationType, string>([
  [NAVIGATION_TYPE.SYSTEM, 'navigation.type.system'],
  [NAVIGATION_TYPE.EXTERNAL, 'navigation.type.external'],
  [NAVIGATION_TYPE.CUSTOM, 'navigation.type.custom'],
])

/**
 * 导航打开方式标签映射
 */
export const NAVIGATION_TARGET_LABEL_MAP = new Map<NavigationTarget, string>([
  [NAVIGATION_TARGET.SELF, 'navigation.target.self'],
  [NAVIGATION_TARGET.BLANK, 'navigation.target.blank'],
])

/**
 * 默认导航配置
 */
const createDefaultConfig = () =>
  JSON.stringify({
    target: NAVIGATION_TARGET.SELF,
    seo_title: '',
    seo_keywords: '',
    seo_description: '',
  })

/**
 * 默认初始化数据
 */
export const NAVIGATION_INIT_DATA = [
  {
    jump_path: '/index',
    name: '首页',
    sort: 9999,
    config: createDefaultConfig(),
  },
  {
    jump_path: '/agent',
    name: '智能体',
    sort: 9998,
    config: createDefaultConfig(),
  },
  {
    jump_path: '/prompt',
    name: '提示词',
    sort: 9997,
    config: createDefaultConfig(),
  },
  {
    jump_path: '/toolkit',
    name: 'AI工具',
    sort: 9996,
    config: createDefaultConfig(),
  },
] as const

/**
 * 表单验证规则配置
 */
export const NAVIGATION_FORM_RULES = {
  NAME_REQUIRED: { required: true, message: 'form_input_placeholder' },
  PATH_REQUIRED: { required: true, message: 'form_input_placeholder' },
} as const

/**
 * 导航相关常量
 */
export const NAVIGATION_CONSTANTS = {
  /** 最大导航项数量 */
  MAX_ITEMS: 8,
  /** 默认排序值 */
  DEFAULT_SORT: 9999,
} as const
