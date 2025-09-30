import service from '@/api/config'
import { handleError } from '@/api/errorHandler'

import type { Banner, RawBanner } from './types'

export const bannerApi = {
  get(): Promise<RawBanner> {
    return service
      .get('/api/enterprises/banner')
      .then(res => res.data.banner)
      .catch(handleError)
  },
  async save(data: Banner) {
    return service
      .put('/api/enterprises/banner', { banner: JSON.stringify(data) })
      .catch(handleError)
  },
}

export default bannerApi
