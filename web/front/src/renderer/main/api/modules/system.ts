import service from '../config'
import { handleError } from '../errorHandler'

export const system = {
  init() {
    return service.get('/api/is_init').catch(handleError)
  }
}

export default system
