import type { LibraryItem, LibraryDisplayItem } from './types'
import { getSimpleDateFormatString } from '@/utils/moment'

export const transformLibraryItem = (item: LibraryItem): LibraryDisplayItem => ({
  ...item,
  created_time: getSimpleDateFormatString({
    date: item.created_time,
    format: 'YYYY-MM-DD hh:mm',
  }),
  updated_time: getSimpleDateFormatString({
    date: item.updated_time,
    format: 'YYYY-MM-DD hh:mm',
  }),
})

export const transformLibraryList = (items: LibraryItem[]): LibraryDisplayItem[] => {
  return items.map(transformLibraryItem)
}

export const getDefaultLibraryRequest = (space_id: number) => ({
  space_id,
  offset: 0,
  limit: 10,
  keyword: '',
})
