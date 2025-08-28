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
  }
}

export default enterprise
