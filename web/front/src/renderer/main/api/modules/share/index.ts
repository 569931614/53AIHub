import service from '../../config'
import { handleError } from '../../errorHandler'
import type { ShareCreateRequest, ShareCreateResponse, ShareFindReponse } from './types'

export const sharesApi = {
  create(data: ShareCreateRequest): Promise<ShareCreateResponse> {
    return service
      .post('/api/shares', data)
      .then((res) => res.data)
      .catch(handleError)
  },
  find(id: ShareCreateResponse['share_id']): Promise<ShareFindReponse> {
    return service.get(`/api/shares/${id}`).then((res) => res.data)
  }
}

export default sharesApi
