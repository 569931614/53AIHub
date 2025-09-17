export const inputTypeList = [
  {
    label: window.$t('variable_type.text'),
    type: 'text',
  },
  {
    label: window.$t('variable_type.textarea'),
    type: 'textarea',
  },
  {
    label: window.$t('variable_type.inputNumber'),
    type: 'inputNumber',
  },
  {
    label: window.$t('variable_type.select'),
    type: 'select',
  },
  { label: window.$t('variable_type.date'), type: 'date', allowed: ['53ai_workflow'] },
  { label: window.$t('variable_type.tag'), type: 'tag', allowed: ['53ai_workflow'] },
  { label: window.$t('variable_type.file'), type: 'file' },
]

export const outputTypeList = [
  { label: window.$t('variable_type.textarea'), type: 'textarea' },
  { label: window.$t('variable_type.image'), type: 'image' },
  { label: window.$t('variable_type.audio'), type: 'audio' },
  { label: window.$t('variable_type.video'), type: 'video' },
  { label: window.$t('variable_type.markdown'), type: 'markdown' },
  { label: window.$t('variable_type.array_text'), type: 'array_text' },
  { label: window.$t('variable_type.array_image'), type: 'array_image' },
  { label: window.$t('variable_type.array_audio'), type: 'array_audio' },
  { label: window.$t('variable_type.array_video'), type: 'array_video' },
]

export const outputDefaultField: Agent.Field = {
  id: '',
  variable: '',
  label: '',
  type: 'text',
  desc: '',
  required: false,
  max_length: 0,
  is_system: false,
  options: [],
  date_format: '',
  multiple: false,
  show_word_limit: false,
  file_type: 'all',
  file_accept: [],
  file_limit: 1,
  file_size: 30,
}
