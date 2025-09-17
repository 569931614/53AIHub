/// <reference types="vite/client" />

declare global {
  interface Window {
    $t: (key: string, params?: Record<string, any>) => string
    $chat53ai?: any
    $getPublicPath: (path: string) => string
    electron?: any
    $isElectron?: boolean
    $noop?: () => void

    api_host: string
    qyy_host: string
    admin_url: string
  }
}

export {}

// declare module '@vue/runtime-core' {
//   interface ComponentCustomProperties {
//     $getPublicPath: (path: string) => string
//     // is electron
//     $isElectron?: boolean
//     $noop?: () => void
//   }
// }
