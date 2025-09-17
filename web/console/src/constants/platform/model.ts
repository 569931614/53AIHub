import { MODEL_VALUES } from './config'

export interface FormConfig {
  label: string
  prop: string
  type: 'input' | 'select' | 'url'
  rules?: unknown[]
  placeholder?: string
  required?: boolean
  size?: 'large' | 'default' | 'small'
  multiple?: boolean
  default?: string
}

// 深度搜索表单配置
export const DEEPSEEK_FORM_CONFIG: FormConfig[] = [
  {
    label: 'module.platform_model_api_endpoint',
    prop: 'base_url',
    type: 'url',
    placeholder: 'module.platform_model_api_endpoint_placeholder',
    required: true,
    default: 'https://api.deepseek.com',
  },
  {
    label: 'module.platform_tool_api_key',
    prop: 'key',
    type: 'input',
    placeholder: 'module.platform_tool_api_key_placeholder',
    required: true,
  },
  {
    label: 'module.platform_model_models',
    prop: 'models',
    type: 'select',
    multiple: true,
    placeholder: 'module.platform_model_models_placeholder',
    required: true,
  },
]

// OpenAI表单配置
export const OPENAI_FORM_CONFIG: FormConfig[] = [
  {
    label: 'module.platform_model_base_url',
    prop: 'base_url',
    type: 'url',
    placeholder: 'module.platform_model_base_url_placeholder',
    required: true,
  },
  {
    label: 'module.platform_tool_api_key',
    prop: 'key',
    type: 'input',
    placeholder: 'module.platform_tool_api_key_placeholder',
    required: true,
  },
  {
    label: 'module.platform_model_models',
    prop: 'models',
    type: 'select',
    multiple: true,
    placeholder: 'module.platform_model_models_placeholder',
    required: true,
  },
]

export const SILICONFLOW_FORM_CONFIG: FormConfig[] = [
  {
    label: 'module.platform_tool_api_key',
    prop: 'key',
    type: 'input',
    placeholder: 'module.platform_tool_api_key_placeholder',
    required: true,
  },
  {
    label: 'module.platform_model_models',
    prop: 'models',
    type: 'select',
    multiple: true,
    placeholder: 'module.platform_model_models_placeholder',
    required: true,
  },
]

// Azure表单配置
export const AZURE_FORM_CONFIG: FormConfig[] = [
  {
    label: 'module.platform_model_name',
    prop: 'name',
    type: 'input',
    placeholder: 'module.platform_model_name_placeholder',
    required: true,
  },
  {
    label: 'module.platform_model_base_url_azure',
    prop: 'base_url',
    type: 'url',
    placeholder: 'module.platform_model_base_url_azure_placeholder',
    required: true,
  },
  {
    label: 'module.platform_tool_api_key',
    prop: 'key',
    type: 'input',
    placeholder: 'module.platform_tool_api_key_placeholder',
    required: true,
  },
  {
    label: 'module.platform_model_other',
    prop: 'other',
    type: 'input',
    placeholder: 'module.platform_model_other_placeholder',
    required: true,
  },
  {
    label: 'module.platform_model_models_azure',
    prop: 'models',
    type: 'select',
    placeholder: 'module.platform_model_models_azure_placeholder',
    required: true,
  },
]

export const BAILIAN_FORM_CONFIG: FormConfig[] = [
  {
    label: 'module.platform_model_api_endpoint',
    prop: 'base_url',
    type: 'url',
    placeholder: 'module.platform_model_base_url_placeholder',
    required: true,
    // default: 'https://example.com',
  },
  {
    label: 'module.platform_tool_api_key',
    prop: 'key',
    type: 'input',
    placeholder: 'module.platform_tool_api_key_placeholder',
    required: true,
  },
  {
    label: 'module.platform_model_models',
    prop: 'models',
    type: 'select',
    multiple: true,
    placeholder: 'module.platform_model_models_placeholder',
    required: true,
  },
]

// 获取表单配置
export const getFormConfig = (channel_type: number): FormConfig[] => {
  switch (channel_type) {
    case MODEL_VALUES.DEEPSEEK:
      return DEEPSEEK_FORM_CONFIG
    case MODEL_VALUES.OPENAI:
      return OPENAI_FORM_CONFIG
    case MODEL_VALUES.AZURE:
      return AZURE_FORM_CONFIG
    case MODEL_VALUES.SILICONFLOW:
      return SILICONFLOW_FORM_CONFIG
    case MODEL_VALUES.BAILIAN:
    case MODEL_VALUES.VOLCENGINE:
    case MODEL_VALUES.QIANFAN:
      return BAILIAN_FORM_CONFIG
    default:
      return []
  }
}
