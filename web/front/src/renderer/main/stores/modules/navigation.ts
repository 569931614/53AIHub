import { defineStore } from 'pinia'
import navigationApi from '@/api/modules/navigation'
import { NAVIGATION_TYPE, NAVIGATION_TARGET } from '@/constants/navigation'
import { cacheManager } from '@/utils/cache'
import { IMG_HOST } from '@/api/host'

const getFormatData = (data: Navigation.State) => {
  try {
    data.config = typeof data.config === 'string' ? JSON.parse(data.config) : data.config
  } catch (error) {
    data.config = {}
  }
  data.type = +data.type || +data.config.type || NAVIGATION_TYPE.SYSTEM
  // data.type_label = NAVIGATION_TYPE_LABEL_MAP.get(data.type)
  data.target = +data.target || +data.config.target || NAVIGATION_TARGET.SELF
  // data.target_label = NAVIGATION_TARGET_LABEL_MAP.get(data.target)
  data.url = data.menu_path = data.jump_path = data.jump_path || ''
  if ([NAVIGATION_TYPE.CUSTOM, NAVIGATION_TYPE.SYSTEM].includes(data.type)) {
    if (data.menu_path !== '/index') data.menu_path = `/index${data.jump_path}`
    data.url = `${location.origin}/#${data.menu_path}`
  }

  // 系统页面图标固定
  if (data.type === 1) {
    const pathIconMap: Record<string, number> = {
      '/index': 1,
      '/agent': 2,
      '/prompt': 3,
      '/toolkit': 4
    }
    data.icon = `${IMG_HOST}/navigation/icon${pathIconMap[data.jump_path]}.png`
  }
  return data
}

const CACHE_KEYS = {
  NAVIGATION_LIST: 'navigation_list'
}

const cacheNavigations = JSON.parse(localStorage.getItem(CACHE_KEYS.NAVIGATION_LIST) || '[]')

export const useNavigationStore = defineStore('navigation', {
  state: (): {
    navigations: Navigation.State[]
    agentNavigation: Partial<Navigation.State>
    promptNavigation: Partial<Navigation.State>
    toolkitNavigation: Partial<Navigation.State>
    homeNavigation: Partial<Navigation.State>
    loading: boolean
  } => ({
    navigations: cacheNavigations,
    agentNavigation: cacheNavigations.find((item) => item.jump_path === '/agent') || {},
    promptNavigation: cacheNavigations.find((item) => item.jump_path === '/prompt') || {},
    toolkitNavigation: cacheNavigations.find((item) => item.jump_path === '/toolkit') || {},
    homeNavigation: cacheNavigations.find((item) => item.jump_path === '/index') || {},
    loading: false
  }),

  getters: {},
  actions: {
    async fetchNavigations(): Promise<Navigation.State[]> {
      this.loading = true
      const fetchData = async () => {
        const { list = [] } = await navigationApi.list().catch(() => {
          this.loading = false
        })
        if (!list.length) {
          await navigationApi.init()
          return fetchData()
        }
        return list
      }
      const list = await cacheManager.getOrFetch(CACHE_KEYS.NAVIGATION_LIST, fetchData)
      this.loading = false
      this.navigations = list.filter((item) => +item.status).map((item) => getFormatData(item))
      this.agentNavigation = this.navigations.find((item) => item.jump_path === '/agent') || {}
      this.promptNavigation = this.navigations.find((item) => item.jump_path === '/prompt') || {}
      this.toolkitNavigation = this.navigations.find((item) => item.jump_path === '/toolkit') || {}
      this.homeNavigation = this.navigations.find((item) => item.jump_path === '/index') || {}
      localStorage.setItem(CACHE_KEYS.NAVIGATION_LIST, JSON.stringify(this.navigations))
      return this.navigations
    }
  }
})
