import service from '@/api/config'
import { handleError } from '@/api/errorHandler'

import type { LibraryListResponse, LibraryListRequest, LibraryCreateRequest } from './types'

export const librariesApi = {
  list(params: LibraryListRequest): Promise<LibraryListResponse> {
    return service
      .get('/api/libraries', { params })
      .then(res => res.data)
      .catch(handleError)
  },

  create(data: LibraryCreateRequest) {
    return service.post('/api/libraries', data).catch(handleError)
  },

  update(library_id: number, data: LibraryCreateRequest) {
    return service.put(`/api/libraries/${library_id}`, data).catch(handleError)
  },

  delete(library_id: number) {
    return service.delete(`/api/libraries/${library_id}`).catch(handleError)
  },
}

export default librariesApi
