export interface SpaceItem {
  id: number
  eid: number
  name: string
  description: string
  icon: string
  owner_id: number
  sort: number
  status: number
  library_count: number
  is_default: boolean
  created_time: number
  updated_time: number
  owner_info: {
    nickname: string
  }
}

export interface SpaceListRequest {
  offset: number
  limit: number
  name?: string
}
export interface SpaceListResponse {
  spaces: SpaceItem[]
  total: number
}

export interface SpaceCreateRequest {
  name: string
  description: string
  icon: string
}

export interface SpaceDisplayItem extends Omit<SpaceItem, 'created_time' | 'updated_time'> {
  created_time: string
  updated_time: string
}
