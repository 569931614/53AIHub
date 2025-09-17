<template>
  <div class="flex flex-col px-4">
    <div class="flex items-center justify-between mb-5">
      <div class="text-base text-[#1D1E1F]">
        {{ $t('debug_preview') }}
      </div>
      <div class="flex-center gap-1 cursor-pointer" @click="handleRestart">
        <el-icon><RefreshRight /></el-icon>
        <span class="text-sm text-[#1D1E1F]">
          {{ $t('restart') }}
        </span>
      </div>
    </div>

    <div class="flex-1 overflow-y-auto">
      <div v-if="showError" class="flex flex-col items-center gap-4">
        <img :src="$getRealPath({ url: '/images/chat/test_error.png' })" class="w-10" />
        <p class="text-sm text-[#1D1E1F] mx-10 text-center">
          {{ errorMessage }}
        </p>
      </div>

      <div v-else-if="showResult">
        <x-bubble-assistant class="!mb-0" :streaming="loading"></x-bubble-assistant>
        <div
          v-if="store.form_data.settings.output_fields.length === 0 && !loading && !initialHasOutputFields"
          class="flex flex-col items-center"
        >
          <div
            class="border prompt-input-wrapper rounded w-full h-full flex flex-col !bg-[#F8F9FA] relative overflow-y-auto"
          >
            <div class="min-h-10 pl-3 pr-2 border-b flex items-center justify-between rounded-t bg-[#F8F9FA]">
              <div class="flex-1 text-sm text-[#4F5052] truncate">JSON</div>
            </div>
            <PromptInput
              v-model="resultString"
              style="flex: none; min-height: 40vh; height: max-content"
              show-line
              show-token
              :word-wrap="true"
            />
          </div>
          <ElButton type="primary" class="mt-8" @click="handleSyncVariables">{{ $t('sync_output_variable') }}</ElButton>
        </div>

        <template v-for="item in result" :key="item.id">
          <div class="text-sm text-[#1D1E1F] mt-2">
            <!-- <x-md-renderer :content="result" /> -->
            <x-bubble-assistant v-if="item.type === 'markdown'" :content="item.value" :streaming="loading">
            </x-bubble-assistant>
            <div v-else-if="item.type.includes('image')" class="overflow-hidden flex flex-col gap-5">
              <img
                v-for="(src, index) in Array.isArray(item.value) ? item.value : [item.value]"
                :key="index"
                :src="src"
                class="max-w-full h-auto object-contain rounded"
              />
            </div>
            <div v-else-if="item.type.includes('video')" class="overflow-hidden flex flex-col gap-5">
              <video
                v-for="(src, index) in Array.isArray(item.value) ? item.value : [item.value]"
                :key="index"
                :src="getSrc(src, item.id)"
                controls
                class="max-w-full h-auto"
              ></video>
            </div>
            <div v-else-if="item.type.includes('audio')" class="overflow-hidden flex flex-col gap-5">
              <audio
                v-for="(src, index) in Array.isArray(item.value) ? item.value : [item.value]"
                :key="index"
                :src="getSrc(src, item.id)"
                controls
                class="max-w-full"
              ></audio>
            </div>
            <div v-else-if="item.type.includes('text')">
              <p
                v-for="(text, index) in Array.isArray(item.value) ? item.value : [item.value]"
                :key="index"
                class="whitespace-pre-wrap break-all"
              >
                {{ text }}
              </p>
            </div>
          </div>
        </template>
      </div>

      <template v-else>
        <el-form ref="formRef" :model="form" label-position="top" require-asterisk-position="right" @submit.prevent>
          <template v-for="(item, index) in form" :key="item.id">
            <el-form-item
              v-if="item.type === 'text'"
              :prop="`${index}.value`"
              :label="item.label"
              :required="item.required"
              :rules="[
                { required: item.required, message: $t('form.input_placeholder') + item.label, trigger: 'blur' },
              ]"
            >
              <el-input
                v-model="item.value"
                size="large"
                :placeholder="$t('form.input_placeholder')"
                :maxlength="item.max_length ? item.max_length : -1"
                :show-word-limit="item.show_word_limit"
              />
              <div v-if="item.desc" class="text-xs text-[#182b50] text-opacity-30 mt-1">
                {{ item.desc }}
              </div>
            </el-form-item>
            <el-form-item
              v-else-if="item.type === 'textarea'"
              :prop="`${index}.value`"
              :label="item.label"
              :required="item.required"
              :rules="[
                { required: item.required, message: $t('form.input_placeholder') + item.label, trigger: 'blur' },
              ]"
            >
              <el-input
                v-model="item.value"
                size="large"
                type="textarea"
                :rows="4"
                resize="none"
                :placeholder="$t('form.input_placeholder')"
                :maxlength="item.max_length ? item.max_length : -1"
                :show-word-limit="item.show_word_limit"
              />
              <div v-if="item.desc" class="text-xs text-[#182b50] text-opacity-30 mt-1">
                {{ item.desc }}
              </div>
            </el-form-item>
            <el-form-item
              v-else-if="item.type === 'inputNumber'"
              :prop="`${index}.value`"
              :label="item.label"
              :required="item.required"
              :rules="[
                { required: item.required, message: $t('form.input_placeholder') + item.label, trigger: 'blur' },
              ]"
            >
              <div>
                <el-input-number
                  v-model="item.value"
                  size="large"
                  :min="1"
                  :placeholder="$t('form.input_placeholder')"
                />
                <div v-if="item.desc" class="text-xs text-[#182b50] text-opacity-30 mt-1">
                  {{ item.desc }}
                </div>
              </div>
            </el-form-item>
            <el-form-item
              v-else-if="item.type === 'select'"
              :prop="`${index}.value`"
              :label="item.label"
              :required="item.required"
              :rules="[
                { required: item.required, message: $t('form.input_placeholder') + item.label, trigger: 'change' },
              ]"
            >
              <el-select
                v-model="item.value"
                class="w-full"
                size="large"
                :multiple="item.multiple"
                :placeholder="$t('form.select_placeholder')"
              >
                <el-option
                  v-for="option in item.options"
                  :key="option.value"
                  :label="option.label"
                  :value="option.label"
                />
              </el-select>
              <div v-if="item.desc" class="text-xs text-[#182b50] text-opacity-30 mt-1">
                {{ item.desc }}
              </div>
            </el-form-item>

            <el-form-item
              v-else-if="item.type === 'date'"
              :prop="`${index}.value`"
              :label="item.label"
              :required="item.required"
              :rules="[
                {
                  required: item.required,
                  message: $t('form.select_placeholder') + item.label,
                  trigger: ['change', 'blur'],
                },
              ]"
            >
              <el-time-picker
                v-if="item.date_format === 'h-m'"
                v-model="item.value"
                format="HH:mm"
                size="large"
                value-format="HH:mm"
                :placeholder="$t('form.select_placeholder')"
              />
              <el-date-picker
                v-else-if="item.date_format === 'y'"
                v-model="item.value"
                size="large"
                type="year"
                value-format="YYYY"
                :placeholder="$t('form.select_placeholder')"
              />
              <el-date-picker
                v-else-if="item.date_format === 'y-m'"
                v-model="item.value"
                size="large"
                type="month"
                value-format="YYYY-MM"
                :placeholder="$t('form.select_placeholder')"
              />
              <el-date-picker
                v-else-if="item.date_format === 'y-m-d'"
                v-model="item.value"
                size="large"
                type="date"
                value-format="YYYY-MM-DD"
                :placeholder="$t('form.select_placeholder')"
              />
              <el-date-picker
                v-else-if="item.date_format === 'y-m-d-h'"
                v-model="item.value"
                size="large"
                type="datetime"
                format="YYYY-MM-DD HH"
                time-format="HH"
                value-format="YYYY-MM-DD HH"
                :placeholder="$t('form.select_placeholder')"
              />
              <el-date-picker
                v-else-if="item.date_format === 'daterange'"
                v-model="item.value"
                size="large"
                type="daterange"
                value-format="YYYY-MM-DD HH:mm"
              />

              <div v-if="item.desc" class="text-xs text-[#182b50] text-opacity-30 mt-1">
                {{ item.desc }}
              </div>
            </el-form-item>

            <el-form-item
              v-else-if="item.type === 'tag'"
              :prop="`${index}.value`"
              :label="item.label"
              :required="item.required"
              :rules="[{ validator: validator(item), trigger: 'change' }]"
            >
              <div>
                <div class="flex flex-wrap gap-3">
                  <template v-for="(tag, childIndex) in item.value" :key="childIndex">
                    <div
                      class="border bordre-[#B0B7C3] rounded-sm min-h-[32px] inline-flex items-center px-3 py-1 text-xs text-[#182B50] text-opacity-80 break-all"
                    >
                      {{ tag }}
                      <el-icon class="cursor-pointer ml-1" color="#d2d5dc" @click="handleDelTag(item, childIndex)">
                        <Close />
                      </el-icon>
                    </div>
                  </template>
                  <el-input
                    v-if="item.focus"
                    v-model="item.temp"
                    autofocus
                    style="width: 104px"
                    class="h-8"
                    :placeholder="$t('form_input_placeholder')"
                    @keypress.enter="handleAddTag(item)"
                    @blur="handleAddTag(item)"
                  />
                  <div
                    v-else
                    class="border bordre-[#B0B7C3] border-dashed rounded-sm h-8 inline-flex items-center px-3 cursor-pointer"
                    @click="handleFocusTag(item)"
                  >
                    <span class="text-xs text-[#182B50] text-opacity-80">+ {{ $t('action_add') }}</span>
                  </div>
                </div>
                <div v-if="item.desc" class="text-xs text-[#182b50] text-opacity-30 mt-1">
                  {{ item.desc }}
                </div>
              </div>
            </el-form-item>

            <el-form-item
              v-else-if="item.type === 'file'"
              :prop="`${index}.value`"
              :label="item.label"
              :required="item.required"
              :rules="[{ validator: validator(item), trigger: ['change', 'blur'] }]"
            >
              <div class="w-full">
                <div v-show="item.file_limit !== item.value.length">
                  <FileUpload
                    ref="uploadRef"
                    v-model:file-list="item.value"
                    class="w-20"
                    drag
                    :accept="item.file_accept.map(item => `.${item}`).join(',')"
                    :limit="item.file_limit"
                    :multiple="item.file_limit === 1 ? false : true"
                    :size="item.file_size"
                    :show-file-list="false"
                  >
                    <div class="w-20 h-20 border border-dashed rounded-sm flex-center flex-col">
                      <!-- <img class="w-4 h-4" src="/images/upload.png" /> -->

                      <div class="text-xs text-[#182B5066] mt-2">点击上传</div>
                    </div>
                  </FileUpload>
                </div>
                <template v-for="file in item.value" :key="file.uid">
                  <div class="h-9 px-2 border rounded mt-3 flex items-center gap-2">
                    <div class="flex-1 text-sm text-[#182B50] truncate">
                      {{ file.name }}
                    </div>
                    <div v-if="file.status === 'success'" class="flex items-center">
                      <el-button type="primary" link @click="handleViewFile(file)"> 查看 </el-button>
                      <div class="w-px h-4 mx-1 bg-[#E3E5EA]" />
                      <el-button type="danger" link @click="handleDelFile(file, item)"> 删除 </el-button>
                    </div>
                    <div v-else class="flex items-center">
                      <el-icon class="animate-rotate">
                        <Loading />
                      </el-icon>
                    </div>
                  </div>
                </template>

                <div class="flex items-center gap-1 mt-2">
                  <el-icon size="14" color="#182B50">
                    <Warning />
                  </el-icon>
                  <span class="text-xs text-[#182B50CC]">单个文件大小不超过{{ item.file_size }}MB</span>
                </div>
                <div>
                  <span class="text-xs text-[#182B50CC]">支持格式：{{ item.file_accept.join('、') }}</span>
                </div>
              </div>
            </el-form-item>
          </template>
        </el-form>
        <div>
          <el-button v-debounce type="primary" size="large" @click="handleStartRunning">
            {{ $t('start_running') }}
          </el-button>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { watch, ref } from 'vue'
import { Loading, RefreshRight, Warning, Close } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import FileUpload from '@/components/Upload/index.vue'
import PromptInput from '@/components/Prompt/input.vue'

import conversationApi from '@/api/modules/conversation'

import { useAgentFormStore } from '../store'
import { useConversationStore } from '@/stores'

import { generateRandomId } from '@/utils'
import { isUrl } from '@/utils/url'
import { outputDefaultField } from '@/constants/agent'

const store = useAgentFormStore()
const conversationStore = useConversationStore()

const formRef = ref<FormInstance>()
const showResult = ref(false)
const loading = ref(false)
const result = ref([])
const resultString = ref('')
const showError = ref(false)
const initialHasOutputFields = ref(false)
const errorMessage = ref('')

interface FormItem extends Agent.Field {
  value: string | string[]
  temp?: string
  focus?: boolean
}

const form = ref<FormItem[]>([])
const abortController = ref<AbortController | null>(null)

const validator = (item: FormItem) => {
  return (rule: any, value: any, callback: any) => {
    if (item.required) {
      let hasVal = false

      if (item.type === 'file') {
        // 处理文件类型
        hasVal = Array.isArray(item.value) && item.value.length > 0
      } else if (Array.isArray(item.value)) {
        // 处理数组类型
        hasVal = item.value.some(val => val && String(val).trim().length > 0)
      } else {
        // 处理字符串类型
        hasVal = item.value && String(item.value).trim().length > 0
      }

      if (hasVal) callback()
      else callback(new Error(`请添加${item.label}`))
    } else {
      callback()
    }
  }
}

const setFormatForm = () => {
  initialHasOutputFields.value = store.form_data.settings.output_fields.length > 0
  form.value = (store.form_data.settings.input_fields || []).map(item => {
    if (['tag', 'file'].includes(item.type)) {
      return {
        ...item,
        value: [],
      }
    }
    return {
      ...item,
      value: item.type === 'select' && item.multiple ? [] : '',
    }
  })
}

// 从对象中获取url
const getSrc = (value: any, id: string) => {
  if (typeof value === 'object' && value !== null) {
    for (const key in value) {
      if (Object.prototype.hasOwnProperty.call(value, key)) {
        const val = value[key]
        if (typeof val === 'string' && isUrl(val)) {
          return val
        }
      }
    }
    result.value = result.value.filter(item => item.id !== id)
    ElMessage.error(window.$t('not_found_url'))
  }
  return value
}

const handleRestart = () => {
  setFormatForm()
  showResult.value = false
  loading.value = false
  result.value = []
  showError.value = false
  errorMessage.value = ''
  if (abortController.value) {
    abortController.value.abort()
    abortController.value = null
  }
}
const handleFocusTag = (item: FormItem) => {
  item.temp = ''
  item.focus = true
}
const handleAddTag = (item: FormItem) => {
  const temp = item.temp?.trim()
  if (temp) {
    item.value.unshift(temp)
    item.temp = ''
  }
  item.focus = false
}
const handleDelTag = (item: FormItem, index: number) => {
  item.value.splice(index, 1)
}

const handleViewFile = (file: any) => {
  window.open(file.url, '_blank')
}
const handleDelFile = (file: any, item: FormItem) => {
  item.value = item.value.filter(item => item.id !== file.id)
}

const getInputs = () => {
  const inputs = form.value.reduce((result, item) => {
    if (item.type === 'file') {
      result[`${item.variable}`] = Array.isArray(item.value)
        ? item.value.map(item => `file_id:${item.id}`).join(',')
        : `file_id:${item.value}`
    } else {
      result[`${item.variable}`] =
        item.type === 'select' && !item.multiple
          ? item.value
          : Array.isArray(item.value)
            ? item.value.join(',')
            : String(item.value)
    }
    return result
  }, {})
  // Object.keys(inputs).forEach(key => {
  //   if (!inputs[key]) {
  //     delete inputs[key]
  //   }
  // })
  return inputs
}

// 没有任何输入时也应正常运行
const getQuestion = inputs => {
  const keys = Object.keys(inputs)
  for (let index = 0; index < keys.length; index++) {
    const key = keys[index]
    const value = inputs[key]
    if (value === undefined) continue // 跳过未定义值

    if (typeof value === 'string' && value.includes('file_id:')) {
      return 'image'
    }
    if (value !== undefined) {
      return String(value).slice(0, 20)
    }
  }
  return '' // 所有值均处理完毕后返回空字符串
}

const handleStartRunning = async () => {
  const isValid = await formRef.value?.validate()
  if (!isValid) return
  const inputs = getInputs()
  loading.value = true

  const conv = await conversationStore.save({
    data: { agent_id: store.agent_data.agent_id, title: getQuestion(inputs) },
  })

  const data = {
    conversation_id: conv.conversation_id,
    model: `agent-${store.agent_data.agent_id}`,
    parameters: inputs,
    stream: true,
  }
  abortController.value = new AbortController()
  showResult.value = true
  conversationApi.workflow
    .run(data, {
      responseType: 'stream',
      onDownloadProgress: e => {
        // processedLength = processStreamData(e, processedLength)
      },
      signal: abortController.value.signal,
    })
    .then(response => {
      const res = JSON.parse(response)
      if (store.form_data.settings.output_fields.length > 0) {
        const output: Record<string, string> = store.form_data.settings.output_fields.reduce((result, item) => {
          result.push({
            id: item.id,
            label: item.label,
            type: item.type,
            variable: item.variable,
            value: res.data.workflow_output_data[item.variable] || '',
          })
          return result
        }, [])
        result.value = output
      } else {
        resultString.value = JSON.stringify(res.data.workflow_output_data, null, 2)
      }
    })
    .catch(res => {
      const resData = JSON.parse(res.response.data)
      showError.value = true
      errorMessage.value = resData.message
    })
    .finally(() => {
      loading.value = false
      abortController.value = null
    })
}

// 同步输出变量
const handleSyncVariables = () => {
  if (!resultString.value) return

  // 格式化result为对象类型
  let resultData: Record<string, any>
  try {
    resultData = JSON.parse(resultString.value)
  } catch (e) {
    console.error('解析result失败:', e)
    return
  }

  // 生成输出字段列表
  const fields: Agent.Field[] = Object.keys(resultData).map(key => {
    return { ...outputDefaultField, id: generateRandomId(10), variable: key, label: key, type: 'textarea' }
  })

  // 更新存储中的输出字段
  store.form_data.settings.output_fields = fields

  showResult.value = false
  setTimeout(() => {
    showResult.value = true
    handleRestart()
  }, 0)
}

watch(
  () => store.form_data.settings.input_fields,
  () => {
    setFormatForm()
  },
  { deep: true, immediate: true }
)
</script>

<style scoped></style>
