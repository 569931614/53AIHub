import { RawProviderItem, ProviderItem } from './types'
import { getProviderByProviderType } from '@/constants/platform/config'

export const transformProviderItem = (item: RawProviderItem): ProviderItem => {
  const provider = getProviderByProviderType(item.provider_type)
  return {
    ...item,
    provider_icon: window.$getRealPath({ url: `/images/platform/${provider.icon}.png` }),
    provider_label: provider.label,
    configs: typeof item.configs === 'string' ? JSON.parse(item.configs) : item.configs || {},
  }
}

export const transformProviderList = (list: RawProviderItem[]): ProviderItem[] => {
  return list.map(transformProviderItem)
}
