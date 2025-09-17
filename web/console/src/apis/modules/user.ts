export default {
  saas_login: {
    url: '/saas/auth/login',
    method: 'POST',
  },
  saas_sms_login: {
    url: '/saas/auth/sms_login',
    method: 'POST',
  },
  saas_logout: {
    url: '/saas/auth/logout',
    method: 'POST',
  },
  logout: {
    url: '/logout',
    method: 'POST',
  },
  register: {
    url: '/register',
    method: 'POST',
  },
  self_info: {
    url: '/users/me',
    method: 'GET',
  },
  reset_password: {
    url: '/saas/auth/reset_password',
    method: 'POST',
  },
  list: {
    url: '/users',
    method: 'GET',
  },
  delete: {
    url: '/users/${user_id}',
    method: 'DELETE',
  },
  update: {
    url: '/users/${user_id}',
    method: 'PUT',
  },
}
