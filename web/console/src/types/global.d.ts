import 'vue'

import type { GlobalMethods } from '@/global/methods'

// 类型声明，避免ts错误
declare global {
  interface Window extends GlobalMethods {
    $t: (key: string, ...args: unknown[]) => string
  }
}
