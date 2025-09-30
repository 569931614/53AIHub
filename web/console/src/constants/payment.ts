/**
 * 支付相关常量配置
 */

// 支付类型枚举
export const PAYMENT_TYPE = {
  ALL: -1,
  WECHAT: 1,
  MANUAL: 2,
  PAYPAL: 3,
  ALIPAY: 4,
} as const

export type PaymentType = (typeof PAYMENT_TYPE)[keyof typeof PAYMENT_TYPE]

// 支付类型标签映射
export const PAYMENT_TYPE_LABEL_MAP = new Map([
  [PAYMENT_TYPE.ALL, window.$t('payment.type.all')],
  [PAYMENT_TYPE.WECHAT, window.$t('payment.type.wechat')],
  [PAYMENT_TYPE.MANUAL, window.$t('payment.type.manual')],
  [PAYMENT_TYPE.PAYPAL, window.$t('payment.type.paypal')],
  [PAYMENT_TYPE.ALIPAY, window.$t('payment.type.alipay')],
])

// 支付类型对应的图标名称
export const PAYMENT_TYPE_ICON_MAP = new Map([
  [PAYMENT_TYPE.WECHAT, 'wechat'],
  [PAYMENT_TYPE.ALIPAY, 'alipay'],
  [PAYMENT_TYPE.MANUAL, 'manual-pay'],
  [PAYMENT_TYPE.PAYPAL, 'paypal'],
])

// 支付类型对应的键名（用于组件引用）
export const PAYMENT_TYPE_KEY_MAP = new Map([
  [PAYMENT_TYPE.WECHAT, 'wechat'],
  [PAYMENT_TYPE.ALIPAY, 'alipay'],
  [PAYMENT_TYPE.MANUAL, 'manual'],
  [PAYMENT_TYPE.PAYPAL, 'paypal'],
])

// 支持的支付类型列表
export const SUPPORTED_PAYMENT_TYPES = [
  PAYMENT_TYPE.WECHAT,
  PAYMENT_TYPE.ALIPAY,
  PAYMENT_TYPE.MANUAL,
  PAYMENT_TYPE.PAYPAL,
] as const

// 默认支付配置
export const DEFAULT_PAYMENT_CONFIG = {
  pay_setting_id: 0,
  pay_config: {},
  extra_config: {},
  pay_status: true,
  pay_type: PAYMENT_TYPE.WECHAT,
} as const

// 支付状态相关常量
export const PAYMENT_STATUS = {
  ENABLED: true,
  DISABLED: false,
} as const

// 支付操作命令
export const PAYMENT_COMMAND = {
  SETTING: 'setting',
  ENABLE: 'enable',
  DISABLE: 'disable',
} as const
