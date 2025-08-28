import { createI18n } from 'vue-i18n'
import HubUiX from 'hub-ui-x'

// element-plus 中的语言配置
import elementZhLocale from 'element-plus/es/locale/lang/zh-cn'
import elementZhTwLocale from 'element-plus/es/locale/lang/zh-tw'
import elementEnLocale from 'element-plus/es/locale/lang/en'
import elementJaLocale from 'element-plus/es/locale/lang/ja'
// 自己的语言配置
import zhLocale from './zh-cn.json'
import twLocale from './zh-tw.json'
import enLocale from './en.json'
import enLocaleV2 from './en-v2.json'
import jaLocale from './ja.json'
import { assign } from '@/utils'

assign(enLocale, enLocaleV2)

// 语言配置整合
const messages = {
  'zh-cn': {
    ...zhLocale,
    elementZhLocale,
  },
  'zh-tw': {
    ...twLocale,
    elementZhTwLocale,
  },
  en: {
    ...enLocale,
    elementEnLocale,
  },
  ja: {
    ...jaLocale,
    elementJaLocale,
  },
}

const browser_lang = navigator.language.toLowerCase()
// const is_zh = /^zh\b/.test(browser_lang)
const is_en = /^en\b/.test(browser_lang)
const is_ja = /^ja\b/.test(browser_lang)
const is_tw = /^tw\b/.test(browser_lang)
console.info('browser_lang: ', browser_lang)

const locale =
  localStorage.getItem('default_lang') || (is_ja ? 'ja' : is_en ? 'en' : is_tw ? 'zh-tw' : 'zh-cn')

// 创建 i18n
const i18n = createI18n({
  legacy: false,
  fallbackLocale: 'zh-cn',
  globalInjection: true, // 全局模式，可以直接使用 $t
  locale,
  messages,
})

HubUiX.setLang(locale)

// 注册到全局，方便js上使用
// @ts-expect-error
window.$t = i18n.global.t
export default i18n
