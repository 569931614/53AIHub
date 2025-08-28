import service from '@/api/config'
import { handleError } from '@/api/errorHandler'

import type { SpaceListResponse, SpaceListRequest, SpaceCreateRequest, SpaceItem } from './types'

export const spacesApi = {
  list(params: SpaceListRequest): Promise<SpaceListResponse> {
    return service
      .get('/api/spaces', { params })
      .then(res => res.data)
      .catch(handleError)
  },

  create(data: SpaceCreateRequest) {
    return service.post('/api/spaces', data).catch(handleError)
  },

  update(space_id: number, data: SpaceCreateRequest) {
    return service.put(`/api/spaces/${space_id}`, data).catch(handleError)
  },

  delete(space_id: number) {
    return service.delete(`/api/spaces/${space_id}`).catch(handleError)
  },

  detail(space_id: number): Promise<SpaceItem> {
    return service
      .get(`/api/spaces/${space_id}`)
      .then(res => res.data)
      .catch(handleError)
  },
}

export default spacesApi
