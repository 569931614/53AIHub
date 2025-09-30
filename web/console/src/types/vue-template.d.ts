import type { GlobalFilters } from '@/global/filters'
import type { GlobalMethods } from '@/global/methods'

// Vue模板类型声明
declare module '@vue/runtime-dom' {
  interface ComponentCustomProperties extends GlobalMethods {
    $filters: GlobalFilters
    $t: (key: string, ...args: unknown[]) => string
  }
}

// 确保Vue 3模板中的类型支持
declare module 'vue' {
  interface ComponentCustomProperties extends GlobalMethods {
    $filters: GlobalFilters
    $t: (key: string, ...args: unknown[]) => string
  }
}
