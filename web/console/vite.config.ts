import path from 'path'
import fs from 'fs'
import type { PluginOption } from 'vite'
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { VitePWA } from 'vite-plugin-pwa'

import wasm from 'vite-plugin-wasm'
import topLevelAwait from 'vite-plugin-top-level-await'

import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'
import conditionalCompilation from './vite-plugins/conditional-compilation'

// import { visualizer } from 'rollup-plugin-visualizer'

// 读取 package.json 获取版本号
const packageJson = JSON.parse(fs.readFileSync('./package.json', 'utf-8'))
const { version } = packageJson

// 创建生成 version.txt 的插件
const versionPlugin = () => {
  return {
    name: 'version-txt',
    writeBundle: {
      sequential: true,
      order: 'post',
      handler: async (options: { dir?: string }) => {
        const outDir = options.dir || 'dist'
        fs.writeFileSync(path.join(outDir, 'version.txt'), version)
      },
    },
  }
}

function setupPlugins(env: ImportMetaEnv): PluginOption[] {
  return [
    // 条件编译插件 - 必须在其他插件之前执行
    conditionalCompilation({
      platform: env.VITE_PLATFORM,
      debug: true, // 强制开启调试模式
    }),
    vue({
      template: {
        compilerOptions: {
          isCustomElement: tag => tag.startsWith('ww-open-data') || tag.startsWith('dt-open-data'),
        },
      },
      script: {
        defineModel: true,
        propsDestructure: true,
      },
    }),
    env.VITE_GLOB_APP_PWA === 'true' &&
      VitePWA({
        injectRegister: 'auto',
        manifest: {
          name: '53aiHub',
          short_name: '53aiHub',
          icons: [
            { src: 'pwa-140x140.png', sizes: '140x140', type: 'image/png' },
            { src: 'pwa-210x210.png', sizes: '210x210', type: 'image/png' },
          ],
        },
      }),
    AutoImport({
      resolvers: [ElementPlusResolver()],
    }),
    Components({
      resolvers: [ElementPlusResolver()],
    }),
    createSvgIconsPlugin({
      iconDirs: [path.resolve(process.cwd(), 'src/icons')],
      symbolId: 'icon-[name]',
    }),
    wasm(),
    topLevelAwait(),
    versionPlugin(),
  ]
}

export default defineConfig(env => {
  const viteEnv = loadEnv(env.mode, process.cwd()) as unknown as ImportMetaEnv
  console.log(viteEnv)
  return {
    base: viteEnv.VITE_BASE_PATH || '/console',
    resolve: {
      alias: {
        '@': path.resolve(process.cwd(), 'src'),
      },
    },
    plugins: setupPlugins(viteEnv),
    server: {
      host: '0.0.0.0',
      port: viteEnv.VITE_PLATFORM === 'km' ? 8003 : 8002,
      open: false,
      proxy: {
        '/api': {
          target: viteEnv.VITE_APP_API_BASE_URL,
          changeOrigin: true, // 允许跨域
          rewrite: urlPath => urlPath.replace('/api/', '/'),
        },
      },
      allowedHosts: ['hubtest.53ai.com', 'hub.53ai.com', 'kmtest.53ai.com', 'km.53ai.com'],
    },
    build: {
      outDir: '../../api/static/console',
      reportCompressedSize: false,
      sourcemap: false,
      commonjsOptions: {
        ignoreTryCatch: false,
      },
      assetsDir: 'static/images/',
      rollupOptions: {
        output: {
          chunkFileNames: 'static/js/[name]-[hash].js',
          entryFileNames: 'static/js/[name]-[hash].js',
          assetFileNames: 'static/[ext]/[name]-[hash].[ext]',
          // 解决打包时Some chunks are larger警告
          // manualChunks(id) {
          //   if (id.includes('node_modules')) {
          //     return id
          //       .toString()
          //       .split('node_modules/')[1]
          //       .split('/')[0]
          //       .toString()
          //   }
          // },
        },
      },
    },
  }
})
