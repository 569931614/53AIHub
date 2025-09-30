import type { PAYMENT_TYPE } from '@/constants/payment'

declare namespace Payment {
  type PaymentType = typeof PAYMENT_TYPE[keyof typeof PAYMENT_TYPE]
}