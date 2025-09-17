<template>
  <el-upload
    :id="uuid"
    ref="uploadRef"
    :class="[hide ? 'absolute right-0 top-0' : '', full ? 'el-upload--full' : '']"
    :data="extraData"
    :accept="accept"
    :name="name"
    :drag="drag"
    :disabled="disabled"
    :action="action"
    :show-file-list="false"
    :limit="limit"
    :multiple="multiple"
    :auto-upload="autoUpload"
    :on-success="handleFileSuccess"
    :on-error="handleFileError"
    :on-progress="handleFileProgress"
    :before-upload="handleBeforeUpload"
    :on-exceed="handleExceed"
    :http-request="httpRequest"
    v-bind="$attrs"
  >
    <slot />
  </el-upload>
  <!-- :on-change="handleFileChange" -->
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { genFileId } from 'element-plus'
import type { UploadFile, UploadInstance, UploadProps, UploadRawFile } from 'element-plus'
import { API_HOST } from '@/api/host'
import uploadApi from '@/api/modules/upload'
import { checkPermission } from '@/utils/permission'

const props = withDefaults(
  defineProps<{
    // .pdf, .md, markdown,.txt
    accept: string
    name?: string
    // 大小为M.
    size?: number
    hide?: boolean
    full?: boolean
    drag?: boolean
    multiple?: boolean
    limit?: number
    disabled?: boolean
    autoUpload?: boolean
    extraData?: any
    action?: string
  }>(),
  {
    accept: '',
    name: 'file',
    size: 15,
    hide: false,
    full: false,
    drag: false,
    multiple: false,
    disabled: false,
    autoUpload: true,
    extraData: () => ({})
  }
)
const emits = defineEmits<{
  (e: 'error', data: any): void
  (e: 'success', data: { id: string; size: string; icon: string; name: string }): void
  (e: 'before', file: File): void
  (e: 'progress', file: any, value: number): void
  (
    e: 'change',
    args?: {
      file: any
      fileList: any[]
    }
  ): void
}>()
const uuid = `upload_${Math.random()}`

const uploadRef = ref<UploadInstance>()
const uploadingCount = ref(0)

let firstType = ''
const isOverLimit = (data: File) => {
  if (!firstType) firstType = data.type

  const { name } = data
  if (!new RegExp(`(${props.accept.split(',').join('|')})$`).test(name)) {
    ElMessage.warning(
      window.$t('file.support_format', {
        format: props.accept
          .replace(/\./g, '')
          .split(',')
          .map((item) => item.toUpperCase())
          .join('、')
      })
    )
    return true
  }
  if (data.size === 0) {
    ElMessage.warning(window.$t('file.file_empty'))
    return true
  }
  if (props.size && data.size / 1024 / 1024 > props.size) {
    ElMessage.warning(window.$t('file.file_exceed', { size: props.size, name }))
    return true
  }

  uploadingCount.value++
  return false
}
const handleBeforeUpload = async (data: File) => {
  if (isOverLimit(data)) return false
  const isLogin = checkPermission()
  if (!isLogin) {
    ElMessage.warning(window.$t('authority.login_not_permission'))
    return false
  }
  // data.loading = true
  emits('before', data)
  return true
}
const handleFileProgress = async (e, data) => {
  if (!data.vid) {
    Object.assign(data, {
      id: '',
      vid: data.raw.uid,
      name: data.raw.name,
      url: '',
      error_msg: '',
      job_id: ''
    })
  } else emits('progress', data, e.percent.toFixed(2))
}

const checkAllUploadsCompleted = () => {
  uploadingCount.value--
  if (uploadingCount.value === 0) firstType = ''
}
const handleFileSuccess = (e, data) => {
  data.loading = false
  data.id = e.id
  data.error_msg = e.error_msg
  data.url = e.url
  emits('success', data)
  checkAllUploadsCompleted()
}
const handleFileError = (e, data) => {
  data.loading = false
  emits('error', {
    error_msg: '上传失败'
  })
  checkAllUploadsCompleted()
}

const httpRequest = (data: any) => {
  return uploadApi.upload(data.file).then((res) => {
    return {
      id: res.data.id,
      url: `${API_HOST}/api/preview/${res.data.preview_key || ''}`,
      size: res.data.size,
      name: res.data.file_name,
      mime_type: res.data.mime_type
    }
  })
}

const handleExceed: UploadProps['onExceed'] = (files) => {
  if (props.limit === 1) {
    uploadRef.value!.clearFiles()
    const file = files[0] as UploadRawFile
    file.uid = genFileId()
    uploadRef.value!.handleStart(file)
    if (props.autoUpload) uploadRef.value!.submit()
  } else {
    ElMessage.warning(window.$t('file.file_exceed_limit', { limit: props.limit }))
  }
}

defineExpose({
  abort: (file: UploadFile) => uploadRef.value!.abort(file),
  submit: () => uploadRef.value!.submit(),
  clearFiles: () => uploadRef.value!.clearFiles(),
  trigger() {
    const parent = document.getElementById(uuid)
    const input = parent?.querySelector('.el-upload__input')
    input?.dispatchEvent(new MouseEvent('click'))
    // input.click()
  },
  handleStart: (file: File) => {
    file.uid = genFileId()
    uploadRef.value.handleStart(file)
  },
  handleRemove: (...args) => {
    uploadRef.value!.handleRemove(...args)
  }
})
</script>

<style scoped>
.el-upload--full {
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  opacity: 0;
}

::v-deep(.el-upload) {
  width: 100%;
  height: 100%;
}
::v-deep(.el-upload-dragger) {
  width: 100%;
  height: 100%;
  background: transparent;
  border: none;
  padding: 0;
  display: inline-flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  outline: 0;
}
</style>
