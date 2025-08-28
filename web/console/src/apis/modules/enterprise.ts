export default {
  saas_list: {
    url: '/saas/enterprise/applies',
    method: 'GET',
  },
  saas_apply: {
    url: '/saas/enterprise/apply',
    method: 'POST',
  },
  saas_detail: {
    url: '/saas/enterprise/${eid}',
    method: 'GET',
  },
  is_saas: {
    url: '/enterprises/is_saas',
    method: 'GET',
  },
  saas_self_info: {
    url: '/saas/enterprise/current',
    method: 'GET',
  },
  self_info: {
    url: '/enterprises/current',
    method: 'GET',
  },
  home_info: {
    url: '/enterprises/homepage',
    method: 'GET',
  },
  update: {
    url: '/enterprises/${eid}',
    method: 'PUT',
  },
  smtp_config: {
    url: '/enterprise-configs',
    method: 'GET',
  },
  smtp_detail: {
    url: '/enterprise-configs/${type}',
    method: 'GET',
  },
  smtp_save: {
    url: '/enterprise-configs/${type}',
    method: 'POST',
  },
  smtp_send: {
    url: '/email/send_test',
    method: 'POST',
  },
}
