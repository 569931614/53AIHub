/// <reference types="vite/client" />
/// <reference types="vue/macros-global" />

// Vue模板类型支持
declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

// 全局类型声明
declare global {
  interface ImportMetaEnv {
    readonly VITE_APP_TITLE: string
    readonly VITE_APP_API_BASE_URL: string
    readonly VITE_PLATFORM: string
    readonly VITE_BASE_PATH: string
    readonly VITE_GLOB_APP_PWA: string
  }

  interface ImportMeta {
    readonly env: ImportMetaEnv
  }
}
