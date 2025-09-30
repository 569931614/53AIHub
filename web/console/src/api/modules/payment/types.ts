/**
 * 支付设置相关类型定义
 */

// 支付类型枚举
export type PaymentType = 1 | 2 | 3 | 4 | -1

// 原始支付设置数据（从API返回的原始格式）
export interface RawPaymentSetting {
  pay_setting_id: string | number
  pay_type: string | number
  pay_status: string | number
  pay_config: string | object
  extra_config: string | object
  created_time: string | number
  updated_time: string | number
}

// 处理后的支付设置数据
export interface PaymentSetting {
  pay_setting_id: number
  pay_type: PaymentType
  pay_label: string
  pay_status: boolean
  pay_config: Record<string, any>
  extra_config: Record<string, any>
  created_time: string
  updated_time: string
}

// 支付设置列表响应
export interface PaymentSettingListResponse {
  pay_settings: RawPaymentSetting[]
}

// 保存支付设置的请求参数
export interface SavePaymentSettingRequest {
  pay_setting_id: number
  pay_config: Record<string, any>
  extra_config?: Record<string, any>
  pay_status?: boolean
  pay_type: PaymentType
}

// 更新支付状态的请求参数
export interface UpdatePaymentStatusRequest {
  pay_status: boolean
}

// 支付设置映射
export interface PaymentSettingMap {
  wechat: PaymentSetting
  alipay: PaymentSetting
  manual: PaymentSetting
  paypal: PaymentSetting
}
