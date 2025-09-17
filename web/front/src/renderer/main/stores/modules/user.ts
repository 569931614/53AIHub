import { defineStore } from 'pinia'
import userApi from '@/api/modules/user'
import { subscriptionApi } from '@/api/modules/subscription'

import eventBus from '@/utils/event-bus'
import { getSimpleDateFormatString } from '@/utils/moment'
import { EVENT_NAMES } from '@/constants/events'
import useEnv from '@/hooks/useEnv'

const { isOpLocalEnv } = useEnv()

export const DEFAULT_GROUP_NAME = '免费版'
export const DEFAULT_GROUP_ICON = 'vip-1'

interface UpdateData {
  nickname?: string
  avatar?: string
}

const TOKEN_KEY = 'access_token'

const ONE_DAY_MS = 1000 * 60 * 60 * 24

const DEFAULT_USER: User.Info = {
  access_token: localStorage.getItem(TOKEN_KEY) || '',
  user_id: '',
  openid: '',
  username: '',
  nickname: '',
  avatar: '',
  email: '',
  eid: 0,
  role: 0,
  mobile: '',
  group_id: 0,
  group_ids: [],
  group_name: DEFAULT_GROUP_NAME,
  group_icon: DEFAULT_GROUP_ICON,
  group_expire_day: 99,
  group_isexpired: false,
  group_expire_time: '',
  is_internal: false
}

export const useUserStore = defineStore('user-store', {
  state: (): {
    info: User.Info
    is_login: boolean
    subscriptions: Subscription.State[]
  } => ({
    info: { ...DEFAULT_USER },
    is_login: !!DEFAULT_USER.access_token,
    subscriptions: []
  }),
  actions: {
    async login(data: User.LoginForm) {
      const res = await userApi.login(data)
      this.setAccessToken(res.data.access_token)
      await this.getUserInfo()
      eventBus.emit(EVENT_NAMES.LOGIN_SUCCESS)
    },

    async sms_login(data: User.SmsLoginForm) {
      const res = await userApi.sms_login(data)
      this.setAccessToken(res.data.access_token)
      await this.getUserInfo()
      eventBus.emit(EVENT_NAMES.LOGIN_SUCCESS)
    },

    async wechat_login(params: { unionid?: string }) {
      const res = await userApi.wechat_login(params).catch(() => ({ data: { access_token: '' } }))
      if (!res.data.user.access_token) return Promise.reject(new Error('access_token is empty'))
      this.setAccessToken(res.data.user.access_token)
      await this.getUserInfo()
      eventBus.emit(EVENT_NAMES.LOGIN_SUCCESS)
      return res.data
    },
    async bind_wechat(data: User.BindWechatForm) {
      const res = await userApi.bind_wechat(data)
      const isCreated = Boolean(res.data.access_token && data.mobile)
      if (isCreated) this.setAccessToken(res.data.access_token)
      await this.getUserInfo()
      if (isCreated) eventBus.emit(EVENT_NAMES.LOGIN_SUCCESS)
    },

    async unbind_wechat() {
      await userApi.unbind_wechat()
      await this.getUserInfo()
    },

    async register(data: User.RegisterForm) {
      const registerData = {
        ...data,
        nickname: data.nickname || data.username
      }
      const res = await userApi.register(registerData)
      this.setAccessToken(res.data.access_token)
      await this.getUserInfo()
    },

    async reset_password(data: User.ResetPasswordForm) {
      const newPasswordData = {
        ...data
      }
      await userApi.reset_password(newPasswordData)
      await this.getUserInfo()
    },

    async change_mobile(data: User.ChangeMobileForm, id: string) {
      const newMobileData = {
        ...data
      }
      await userApi.change_mobile(newMobileData, id)
      await this.getUserInfo()
    },

    async update(data: UpdateData) {
      const res = await userApi.update(data)
      Object.assign(this.info, {
        avatar: res.data.avatar,
        nickname: res.data.nickname
      })
    },

    async getUserInfo() {
      if (!localStorage.getItem(TOKEN_KEY)) return
      try {
        const [res, { list: subscription_list = [] }] = await Promise.all([
          userApi.me(),
          subscriptionApi.list()
        ])
        const info = {
          access_token: res.data.access_token || '',
          user_id: res.data.user_id || '',
          openid: res.data.openid || '',
          username: res.data.username || '',
          nickname: res.data.nickname || '',
          avatar:
            res.data.avatar.replace(/^(\/\/)/, 'http://') ||
            'https://chat.53ai.com/images/robot_avatar.png',
          email: res.data.email || '',
          eid: res.data.eid || 0,
          role: res.data.role || 0,
          mobile: res.data.mobile || '',
          group_id: res.data.group_id || 0,
          group_ids: res.data.group_ids || [],
          group_name: res.data.group_name || DEFAULT_GROUP_NAME,
          group_icon: res.data.group_icon || DEFAULT_GROUP_ICON,
          group_expire_time: res.data.expired_time
            ? getSimpleDateFormatString({
                date: res.data.expired_time,
                format: 'YYYY-MM-DD hh:mm'
              })
            : '',
          group_expire_day: res.data.expired_time
            ? Math.max(
                Math.ceil((new Date(res.data.expired_time).getTime() - Date.now()) / ONE_DAY_MS),
                0
              )
            : 99,
          group_isexpired: res.data.expired_time ? res.data.expired_time < Date.now() : false,
          is_internal: res.data.type === 2
        }
        this.subscriptions = subscription_list
        const subscription_data = subscription_list.find(
          (item = {}) => item.group_id === info.group_id
        )
        if (info.is_internal || !subscription_data) {
          info.group_expire_time = ''
          info.group_isexpired = false
          info.group_expire_day = 99
        } else {
          info.group_name = subscription_data.group_name || DEFAULT_GROUP_NAME
          info.group_icon = subscription_data.logo_url || DEFAULT_GROUP_ICON
          if (subscription_data.is_default) {
            info.group_expire_time = ''
            info.group_isexpired = false
            info.group_expire_day = 99
          }
        }
        this.info = info
        this.is_login = true
        if (window.$chat53ai)
          window.$chat53ai.$win({ type: 'agenthub_login', data: JSON.stringify({ ...this.info }) })
      } catch (error: any) {
        const response = error.response || {}
        const data = response.data || error || {}
        const { message } = data
        if (['token expired', 'forbidden'].includes(message)) {
          this.logout({ redirectDisabled: true })
        }
        throw error
      }
    },
    setGroupName(group_name: string) {
      this.info.group_name = group_name || DEFAULT_GROUP_NAME
    },
    setGroupIcon(group_icon: string) {
      this.info.group_icon = group_icon || DEFAULT_GROUP_ICON
    },
    setAccessToken(token: string) {
      localStorage.setItem(TOKEN_KEY, token)
      this.info.access_token = token
    },
    updateInfo(data) {
      this.info = {
        ...this.info,
        ...data
      }
      this.is_login = true
    },
    async logout({ redirectDisabled = false }: { redirectDisabled?: boolean } = {}) {
      this.info = { ...DEFAULT_USER }
      this.is_login = false
      if (!isOpLocalEnv.value) {
        await userApi.logout()
      }
      localStorage.removeItem(TOKEN_KEY)
      eventBus.clearCache(EVENT_NAMES.LOGIN_SUCCESS)
      setTimeout(() => {
        if (!redirectDisabled && !window.$isElectron) {
          window.location.href = '/'
        }
      }, 800)
    }
  }
})
