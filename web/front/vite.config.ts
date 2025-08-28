import { resolve } from 'path'
import { defineConfig, loadEnv } from 'vite'

import { commonPlugins, commonServer, commonResolve, commonCss, commonBuild } from './vite.common'

export default defineConfig(({ command, mode }) => {
  const env = loadEnv(mode, process.cwd())
  return {
    root: resolve('src/renderer'),
    publicDir: resolve('src/renderer/public/'),
    server: commonServer,
    resolve: commonResolve,
    css: commonCss,
    build: {
      ...commonBuild,
      outDir: resolve('dist')
    },
    plugins: commonPlugins,
    define: {
      'process.env': env
    }
  }
})
