/**
 * 支付设置API接口
 * 只包含HTTP请求调用，业务逻辑在transform中处理
 */
import service from '@/api/config'
import { handleError } from '@/api/errorHandler'
import type {
  PaymentSettingListResponse,
  SavePaymentSettingRequest,
  UpdatePaymentStatusRequest,
} from './types'

export const paymentApi = {
  /**
   * 获取支付设置列表
   * @returns 支付设置列表响应
   */
  getPaymentSettings(): Promise<PaymentSettingListResponse> {
    return service
      .get('/api/pay_settings')
      .then(res => res.data)
      .catch(handleError)
  },

  /**
   * 保存支付设置
   * @param data 支付设置数据
   * @returns API响应
   */
  savePaymentSetting(data: SavePaymentSettingRequest) {
    const { pay_setting_id, ...requestData } = data

    return pay_setting_id
      ? service
          .patch(`/api/pay_settings/${pay_setting_id}/config`, {
            pay_config: requestData.pay_config,
            extra_config: requestData.extra_config,
          })
          .catch(handleError)
      : service.post('/api/pay_settings', requestData).catch(handleError)
  },

  /**
   * 更新支付状态
   * @param pay_setting_id 支付设置ID
   * @param data 状态更新数据
   * @returns API响应
   */
  updatePaymentStatus(pay_setting_id: number, data: UpdatePaymentStatusRequest) {
    return service.patch(`/api/pay_settings/${pay_setting_id}/status`, data).catch(handleError)
  },
}

export default paymentApi
