/**
 * 表单验证规则
 */

export const textValidator = ({ rule, value, callback, message } = {}) => {
  value = (String(value) || '').trim()
  if (!value) return callback(new Error(window.$t(message)))
  callback()
}

export const getPasswordRules = () => {
  return {
    validator: (rule: any, value: string) => {
      if (value.length < 8 || value.length > 20) {
        return Promise.reject(window.$t('form.password_length'))
      }
      if (/[\u4e00-\u9fa5]/.test(value)) {
        return Promise.reject(window.$t('form.password_no_chinese'))
      }
      return Promise.resolve()
    },
    trigger: 'blur'
  }
}

export const getMobileRules = () => {
  return {
    validator: (rule: any, value: string) => {
      // 手机号正则（中国大陆手机号）
      const phoneRegex = /^1[3-9]\d{9}$/

      if (!phoneRegex.test(value)) {
        return Promise.reject(window.$t('form.mobile_format'))
      }
      return Promise.resolve()
    },
    trigger: 'blur'
  }
}

export const getEmailRules = () => {
  return {
    validator: (rule: any, value: string) => {
      if (!value) {
        return Promise.reject(window.$t('form.email_format'))
      }
      // 邮箱正则
      const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
      if (!emailRegex.test(value)) {
        return Promise.reject(window.$t('form.email_format'))
      }
      return Promise.resolve()
    },
    trigger: 'blur'
  }
}

export const getAccountOrEmailRules = () => {
  return {
    validator: (rule: any, value: string) => {
      if (!value) {
        return Promise.reject(window.$t('form.account_format'))
      }
      const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
      if (emailRegex.test(value)) {
        return Promise.resolve()
      }
      const mobileRegex = /^1[3-9]\d{9}$/
      if (mobileRegex.test(value)) {
        return Promise.resolve()
      }
      return Promise.reject(window.$t('form.account_format'))
    },
    trigger: 'blur'
  }
}

export const getAccountRules = () => {
  return {
    validator: (rule: any, value: string) => {
      value = value.trim()
      if (value.length < 5 || value.length > 20) {
        return Promise.reject(window.$t('form.username_length'))
      }
      const reg = /^[a-zA-Z0-9_]+$/
      if (!reg.test(value)) {
        return Promise.reject(window.$t('form.account_format'))
      }
      return Promise.resolve()
    },
    trigger: 'blur'
  }
}

export const getConfirmPasswordRules = (form: any, passwordField: string) => {
  return {
    validator: (rule: any, value: string) => {
      if (value !== form[passwordField]) {
        return Promise.reject(window.$t('form.password_not_match'))
      }
      return Promise.resolve()
    },
    trigger: 'blur'
  }
}

export const getRequiredRules = (message: string, trigger: string | string[] = 'blur') => {
  return {
    required: true,
    message: message || window.$t('form.input_placeholder'),
    trigger
  }
}

export const generateInputRules = ({
  message = 'form_input_placeholder',
  trigger = ['blur', 'change'],
  validator = ['text']
} = {}) => {
  const rules = []
  if (validator.includes('text'))
    rules.push({
      validator: (rule, value, callback) => textValidator({ rule, value, callback, message }),
      trigger
    })
  return rules
}
