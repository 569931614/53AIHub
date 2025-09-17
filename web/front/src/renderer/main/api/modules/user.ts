import service from '../config'
import { handleError } from '../errorHandler'
import useEnv from '@/hooks/useEnv'

const { isOpLocalEnv } = useEnv()

const userApi = {
  login(data: { username: string; password: string }) {
    return service.post('/api/login', data).catch(handleError)
  },
  logout() {
    return service.post('/api/logout').catch(handleError)
  },
  sms_login(data: { mobile: string; verify_code: string }) {
    return service.post('/api/sms_login', data).catch(handleError)
  },
  wechat_login(params: { unionid: string }) {
    // return service.get('/api/saas/wechat/user', { params }).catch(handleError)
    return service.get('/api/saas/wechat/user', { params })
  },
  bind_wechat(data: {
    mobile?: string
    verify_code?: string
    openid: string
    unionid?: string
    nickname?: string
  }) {
    let api_url = '/api/saas/wechat/bind'
    if (data.mobile) api_url = '/api/saas/wechat/user'
    return service.post(api_url, data).catch(handleError)
  },
  unbind_wechat() {
    return service.post('/api/saas/wechat/unbind').catch(handleError)
  },
  register(data: { username: string; password: string; nickname: string; verfiy_code?: string }) {
    return service.post('/api/register', data).catch(handleError)
  },
  reset_password(data: User.ResetPasswordForm) {
    if (!isOpLocalEnv.value) {
      return service.post('/api/saas/auth/reset_password', data).catch(handleError)
    }
    return service.post('/api/reset_password', data).catch(handleError)
  },
  change_mobile(data: User.ChangeMobileForm, id: string) {
    return service.patch(`/api/users/${id}/mobile`, data).catch(handleError)
  },
  me() {
    return service.get('/api/users/me').catch(handleError)
  },
  update(data: { nickname?: string; avatar?: string }) {
    return service.put(`/api/users/me`, data).catch(handleError)
  },
  updatePassword(data: { password: string; newPassword: string }) {
    return service.put(`/api/users/password`, data).catch(handleError)
  },
  checkUsername(account: string) {
    return service
      .post(`/api/check_account`, {
        account
      })
      .catch(handleError)
  },
  update_default_subscription(user_id: number) {
    return service.put(`/api/users/${user_id}/default_subscription`).catch(handleError)
  }
}

export default userApi

export type { userApi }
