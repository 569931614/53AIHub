export interface LibraryItem {
  id: number
  eid: number
  space_id: number
  creator_id: number
  name: string
  icon: string
  description: string
  sort: number
  status: number
  created_time: number
  updated_time: number
}

export interface LibraryListResponse {
  list: LibraryItem[]
  total: number
}

export interface LibraryListRequest {
  space_id: number
  offset: number
  limit: number
  keyword?: string
}

export interface LibraryCreateRequest {
  space_id: number
  name: string
  icon: string
  description: string
}

export interface LibraryDisplayItem extends Omit<LibraryItem, 'created_time' | 'updated_time'> {
  created_time: string
  updated_time: string
}
