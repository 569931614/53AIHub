// 验证器函数类型
type ValidatorFn = (value: string, message: string) => string | null

// 验证规则类型
type ValidationRule = {
  validator: (rule: unknown, value: unknown, callback: (error?: Error) => void) => void
  trigger: string[]
}

// 基础验证器工厂函数
const createValidator = (validatorFn: ValidatorFn) => {
  return ({
    value,
    callback,
    message,
  }: {
    value: unknown
    callback: (error?: Error) => void
    message: string
  }): void => {
    const trimmedValue = String(value || '').trim()
    const error = validatorFn(trimmedValue, message)
    callback(error ? new Error(error) : undefined)
  }
}

// 各种验证器实现
const validators = {
  // 必填验证
  required: (value: string, message: string) => (!value ? message : null),

  // 链接验证
  link: (value: string, message: string) => {
    if (!value) return message
    const pattern =
      /^(https?:\/\/)?((([\w.-]+)(\.[\w.-]+)+)|((\d{1,3}\.){3}\d{1,3}))(:\d+)?([\/#\?].*)?$/
    return pattern.test(value) ? null : window.$t('form_link_validator')
  },

  // 手机号验证
  mobile: (value: string, message: string) => {
    if (!value) return message
    const pattern = /^(\+86)?(13[0-9]|14[0-9]|15[0-9]|16[0-9]|17[0-9]|18[0-9]|19[0-9])\d{8}$/
    return pattern.test(value) ? null : window.$t('form_mobile_validator')
  },

  // 邮箱验证
  email: (value: string, message: string) => {
    if (!value) return message
    const pattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
    return pattern.test(value) ? null : window.$t('form_email_validator')
  },

  // 手机号或邮箱验证
  mobileOrEmail: (value: string, message: string) => {
    if (!value) return message
    const mobilePattern = /^(13[0-9]|14[0-9]|15[0-9]|16[0-9]|17[0-9]|18[0-9]|19[0-9])\d{8}$/
    const emailPattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
    return mobilePattern.test(value) || emailPattern.test(value)
      ? null
      : window.$t('form_mobile_or_email_validator')
  },

  // 密码验证
  password: (value: string, message: string) => {
    if (!value) return message
    if (/[\u4e00-\u9fa5\s]/.test(value)) return window.$t('form_password_validator')
    return null
  },

  // URL验证
  url: (value: string, message: string) => {
    if (!value) return message
    const pattern = /^(https?:\/\/)?([\w.-]+)(\.[\w.-]+)+([\/#\?].*)?$/
    return pattern.test(value) ? null : window.$t('form_url_validator')
  },

  // 路径验证
  path: (value: string, message: string) => {
    if (!value) return message
    const pattern = /^(\/[\w-]+)+$/
    return pattern.test(value) ? null : window.$t('form_path_validator')
  },

  // 图片URL验证
  image: (value: string, message: string) => {
    if (!value) return message
    const pattern = /^https?:\/\/.+\.(jpg|jpeg|png|gif|bmp|webp)$/
    return pattern.test(value) ? null : window.$t('form_image_validator')
  },

  // 变量名验证
  variable: (value: string, message: string) => {
    if (!value) return message
    const pattern = /^[a-zA-Z_][a-zA-Z0-9_]*$/
    return pattern.test(value) ? null : window.$t('form_variable_validator')
  },
} as const

export type ValidatorType = keyof typeof validators

// 导出各个验证器
export const textValidator = createValidator(validators.required)
export const linkValidator = createValidator(validators.link)
export const mobileValidator = createValidator(validators.mobile)
export const emailValidator = createValidator(validators.email)
export const mobileOrEmailValidator = createValidator(validators.mobileOrEmail)
export const passwordValidator = createValidator(validators.password)
export const urlValidator = createValidator(validators.url)
export const pathValidator = createValidator(validators.path)
export const imageValidator = createValidator(validators.image)
export const variableValidator = createValidator(validators.variable)

// 生成表单验证规则
export const generateFormRules = ({
  message = window.$t('form_input_placeholder'),
  trigger = ['blur', 'change'],
  validator = ['required'],
}: {
  message?: string
  trigger?: string[]
  validator?: string[]
} = {}): ValidationRule[] => {
  const rules: ValidationRule[] = []

  // 处理验证器
  validator.forEach(v => {
    const fn = validators[v as ValidatorType]
    if (fn) {
      rules.push({
        validator: (rule: unknown, value: unknown, callback: (error?: Error) => void) => {
          const err = fn(String(value || '').trim(), message)
          callback(err ? new Error(err) : undefined)
        },
        trigger,
      })
    }
  })

  return rules
}
