import { defineStore } from 'pinia'
import enterpriseApi from '@/api/modules/enterprise'
import { cacheManager } from '@/utils/cache'

const getDefaultLogo = () => window.$getPublicPath('/images/default_logo.png')

const CACHE_KEYS = {
  ENTERPRISE_INFO: 'enterprise_info'
}

export const useEnterpriseStore = defineStore('enterprise-store', {
  state: (): Enterprise.State => ({
    id: 0,
    type: '',
    banner: '',
    timezone: '',
    domain: '',
    slogan: '',
    status: 0,
    template_type: '',
    layout_type: '',
    created_time: 0,
    updated_time: 0,
    logo: getDefaultLogo(),
    ico: getDefaultLogo(),
    display_name: '',
    language: 'zh-cn',
    copyright: '',
    keywords: [],
    description: '',
    banner_info: {
      url_list: [],
      interval: ''
    },
    template_style_info: {
      style_type: 'software'
    },
    is_independent: false,
    is_enterprise: false,
    is_industry: false,
    is_install_wecom: false
  }),
  getters: {
    isSoftStyle: (state) => state.template_style_info.style_type === 'software'
  },
  actions: {
    initTemplateStyle() {
      let {
        theme_color,
        text_color,
        nav_bg_color,
        nav_text_color,
        page_footer_bg_color,
        page_footer_text_color
      } = this.template_style_info
      theme_color = theme_color || '#2563eb'
      text_color = text_color || '#333333'
      nav_bg_color = nav_bg_color || '#ffffff'
      nav_text_color = nav_text_color || '#333333'
      page_footer_bg_color = page_footer_bg_color || '#18191f'
      page_footer_text_color = page_footer_text_color || '#f2f2f2'

      document.documentElement.style.setProperty('--el-color-primary', theme_color)
      document.documentElement.style.setProperty('--el-color-primary-light-3', `${theme_color}cc`)
      document.documentElement.style.setProperty('--el-color-primary-light-5', `${theme_color}aa`)
      document.documentElement.style.setProperty('--el-color-primary-light-7', `${theme_color}99`)
      document.documentElement.style.setProperty('--el-color-primary-light-8', `${theme_color}66`)
      document.documentElement.style.setProperty('--el-color-primary-light-9', `${theme_color}33`)
      document.documentElement.style.setProperty('--el-text-color-primary', text_color)
      document.documentElement.style.setProperty('--el-text-color-regular', `${text_color}cc`)
      document.documentElement.style.setProperty('--el-text-color-secondary', `${text_color}99`)
      document.documentElement.style.setProperty('--el-text-color-placeholder', `${text_color}66`)
      document.documentElement.style.setProperty('--el-color-nav-bg', nav_bg_color)
      document.documentElement.style.setProperty('--el-color-nav-text', nav_text_color)
      document.documentElement.style.setProperty('--el-color-page-footer-bg', page_footer_bg_color)
      document.documentElement.style.setProperty(
        '--el-color-page-footer-text',
        page_footer_text_color
      )
    },
    setMetaDescription(description: string = '') {
      this.description = description
      if (!description) return
      const meta = document.querySelector('meta[name="description"]')
      if (meta) {
        meta.setAttribute('content', description)
      } else {
        const meta = document.createElement('meta')
        meta.setAttribute('name', 'description')
        meta.setAttribute('content', description)
        document.head.appendChild(meta)
      }
    },
    setMetaKeywords(keywords: string[] = []) {
      this.keywords = keywords
      if (!keywords.length) return
      const meta = document.querySelector('meta[name="keywords"]')
      if (meta) {
        meta.setAttribute('content', keywords.join(', '))
      } else {
        const meta = document.createElement('meta')
        meta.setAttribute('name', 'keywords')
        meta.setAttribute('content', keywords.join(', '))
        document.head.appendChild(meta)
      }
    },
    setDocumentTitleAndIcon(title: string, iconUrl: string) {
      document.title = title
      const link = document.querySelector('link[rel="icon"]') || document.createElement('link')
      link.rel = 'icon'
      link.href = iconUrl || getDefaultLogo()
      if (!document.querySelector('link[rel="icon"]')) {
        document.head.appendChild(link)
      }
    },
    setAppLanguage(language: string) {
      const lang = language === 'En' ? 'en' : language
      this.language = lang
    },
    saveToStorage() {
      localStorage.setItem(
        'enterprise',
        JSON.stringify({
          logo: this.logo,
          display_name: this.display_name,
          language: this.language
        })
      )
    },
    async loadFromStorage(): Promise<Enterprise.State> {
      const storedEnterprise = localStorage.getItem('enterprise')
      if (storedEnterprise) {
        try {
          const parsedEnterprise = JSON.parse(storedEnterprise)
          this.logo = parsedEnterprise.logo || getDefaultLogo()
          this.display_name = parsedEnterprise.display_name
          this.language = parsedEnterprise.language

          this.setDocumentTitleAndIcon(parsedEnterprise.display_name, parsedEnterprise.logo)
          this.setAppLanguage(parsedEnterprise.language)
        } catch (error) {
          console.error('解析localStorage中的企业信息失败', error)
        }
      }
      return this
    },
    async loadInfo(): Promise<Enterprise.State> {
      const fetchInfo = async () => {
        const res = await enterpriseApi.current()
        return res.data.enterprise
      }

      try {
        const info = await cacheManager.getOrFetch(CACHE_KEYS.ENTERPRISE_INFO, fetchInfo)
        const {
          display_name,
          logo,
          language,
          copyright,
          ico,
          keywords,
          description,
          banner,
          template_type,
          type,
          wecom_install_info
        } = info

        this.logo = logo || getDefaultLogo()
        this.ico = ico || this.logo || getDefaultLogo()
        this.display_name = display_name
        this.copyright = copyright
        this.description = description || ''
        this.is_independent = type === 'independent'
        this.is_enterprise = type === 'enterprise'
        this.is_industry = type === 'industry'
        this.is_install_wecom = wecom_install_info?.install_wecom_app || false
        // 兼容都为false的情况
        if (!this.is_independent && !this.is_enterprise && !this.is_industry)
          this.is_independent = true

        try {
          this.keywords = JSON.parse(keywords || '[]')
        } catch (error) {
          // console.error(error)
          this.keywords = []
        }
        try {
          this.banner_info = JSON.parse(banner || '{}')
        } catch (error) {
          // console.error(error)
          this.banner_info = {
            url_list: [],
            interval: ''
          }
        }
        try {
          this.template_style_info = JSON.parse(template_type || '{}')
        } catch (error) {
          // console.error(error)
          this.template_style_info = {}
        }
        if (!['website', 'software'].includes(this.template_style_info.style_type))
          this.template_style_info.style_type = 'website'
        // 屏幕宽度小于768或electron环境 默认软件风格
        const isMobile = window.innerWidth < 768 || window.electron
        // 移动端默认软件风格
        if (isMobile) this.template_style_info.style_type = 'software'

        this.setAppLanguage(language)
        this.setDocumentTitleAndIcon(display_name, this.ico)
        this.setMetaKeywords(this.keywords)
        this.setMetaDescription(this.description)
        this.saveToStorage()
        this.initTemplateStyle()
      } catch (error) {
        console.error('获取企业信息失败', error)
      }
      return this
    }
  }
})
