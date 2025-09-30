import { createApp } from 'vue'
import HubUiX from 'hub-ui-x'
import i18n from './locales/index'

import App from './App.vue'
import { setupStore } from './stores'
import { setupRouter } from './router'
import { setupAssets } from './plugins/index'
import { setupDirective } from './directive'
import { setupFilter, setupGlobal } from './global'

import 'element-plus/theme-chalk/el-message.css'
import 'element-plus/theme-chalk/el-notification.css'
import 'element-plus/theme-chalk/el-message-box.css'
import 'element-plus/theme-chalk/el-overlay.css'

import 'hub-ui-x/dist/hub-ui-x.css'

import 'virtual:svg-icons-register'

async function bootstrap() {
  const app = createApp(App)

  setupAssets()
  setupStore(app)

  setupGlobal(app)
  setupFilter(app)
  setupDirective(app)

  await setupRouter(app)
  app.use(i18n)
  app.use(HubUiX)
  app.mount('#app')
}

bootstrap()
