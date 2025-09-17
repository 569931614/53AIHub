import service from '../config'
import { handleError } from '../errorHandler'

export const enterprise = {
  current() {
    return service.get('/api/enterprises/current').catch(handleError)
  },
  get(id: string) {
    return service.get(`/api/enterprises/${id}`).catch(handleError)
  },
  getSMTPInfo(type: string) {
    return service.get(`/api/enterprise-configs/${type}/enabled`).catch(handleError)
  },
  async update(
    id: number,
    data: {
      display_name: string
      logo: string
      language: string
      template_type: string
    }
  ) {
    return service.put(`/api/enterprises/${id}`, data).catch(handleError)
  }
}

export default enterprise
