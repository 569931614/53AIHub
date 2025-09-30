import {
  NAVIGATION_TYPE,
  NAVIGATION_TYPE_LABEL_MAP,
  NAVIGATION_TARGET,
  NAVIGATION_TARGET_LABEL_MAP,
} from '@/constants/navigation'
import type { NavigationItem, RawNavigationItem } from './types'
import { img_host } from '@/utils/config'

/**
 * 获取默认的导航项配置
 */
export const getDefaultNavigationItem = (): Partial<NavigationItem> => ({
  type: NAVIGATION_TYPE.EXTERNAL,
  target: NAVIGATION_TARGET.SELF,
  config: {
    target: NAVIGATION_TARGET.SELF,
    seo_title: '',
    seo_keywords: '',
    seo_description: '',
  },
  status: 1,
  sort: 0,
})

/**
 * 转换单个导航项数据
 */
export function transformNavigationItem(rawItem: RawNavigationItem): NavigationItem {
  try {
    // 解析config字段
    let config: any = {}
    if (typeof rawItem.config === 'string') {
      config = JSON.parse(rawItem.config)
    } else if (typeof rawItem.config === 'object' && rawItem.config !== null) {
      config = rawItem.config
    }

    // 确定type值
    const type = (Number(rawItem.type) || Number(config.type) || NAVIGATION_TYPE.SYSTEM) as any

    // 确定target值
    const target = (Number(rawItem.target) ||
      Number(config.target) ||
      NAVIGATION_TARGET.SELF) as any

    // 获取标签
    const type_label = NAVIGATION_TYPE_LABEL_MAP.get(type)
    const target_label = NAVIGATION_TARGET_LABEL_MAP.get(target)

    let icon = rawItem.icon
    // 系统页面图标固定
    if (rawItem.type === 1) {
      const pathIconMap: Record<string, number> = {
        '/index': 1,
        '/agent': 2,
        '/prompt': 3,
        '/toolkit': 4,
      }
      icon = `${img_host}/navigation/icon${pathIconMap[rawItem.jump_path]}.png`
    }

    return {
      ...rawItem,
      type,
      type_label,
      target,
      target_label,
      config: {
        target,
        seo_title: config.seo_title || '',
        seo_keywords: config.seo_keywords || '',
        seo_description: config.seo_description || '',
        agent_id: config.agent_id,
        agent_class_id: config.agent_class_id,
      },
      icon,
    }
  } catch (error) {
    console.error('转换导航项数据失败:', error)
    return {
      ...rawItem,
      ...getDefaultNavigationItem(),
      type: NAVIGATION_TYPE.SYSTEM,
      target: NAVIGATION_TARGET.SELF,
      type_label: NAVIGATION_TYPE_LABEL_MAP.get(NAVIGATION_TYPE.SYSTEM),
      target_label: NAVIGATION_TARGET_LABEL_MAP.get(NAVIGATION_TARGET.SELF),
    } as NavigationItem
  }
}

/**
 * 转换导航列表数据
 */
export function transformNavigationList(rawList: RawNavigationItem[]): NavigationItem[] {
  try {
    return rawList.map(item => transformNavigationItem(item))
  } catch (error) {
    console.error('转换导航列表数据失败:', error)
    return []
  }
}
