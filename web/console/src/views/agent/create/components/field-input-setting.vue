<template>
  <el-dialog
    v-model="visible"
    width="600"
    destroy-on-close
    :title="widgetForm.id ? $t('action.edit') : $t('action.add')"
  >
    <el-form
      ref="formRef"
      :model="widgetForm"
      label-width="100px"
      label-position="top"
      require-asterisk-position="right"
      style="max-height: 520px; overflow-y: auto"
    >
      <el-form-item
        :label="$t('agent.variable_name')"
        prop="variable"
        required
        :rules="generateInputRules({ message: 'form.input_placeholder', validator: ['variable'] })"
      >
        <el-input
          v-model="widgetForm.variable"
          :disabled="widgetForm.is_system"
          size="large"
          maxlength="30"
          show-word-limit
          :placeholder="$t('form.input_placeholder') + $t('agent.variable_name')"
        />
      </el-form-item>
      <el-form-item
        :label="$t('agent.variable_label')"
        prop="label"
        required
        :rules="generateInputRules({ message: 'form.input_placeholder' })"
      >
        <el-input
          v-model="widgetForm.label"
          size="large"
          maxlength="30"
          show-word-limit
          :placeholder="$t('form.input_placeholder') + $t('agent.variable_label')"
        />
      </el-form-item>
      <el-form-item :label="$t('agent.variable_type')">
        <template v-if="widgetForm.is_system">
          <el-input :model-value="typeLabel" size="large" disabled :placeholder="$t('form.input_placeholder')" />
        </template>
        <div v-else class="flex flex-wrap gap-2">
          <template v-for="item in typeList" :key="item.type">
            <div
              class="w-[100px] h-10 border rounded flex-center gap-1 cursor-pointer"
              :class="[
                widgetForm.type === item.type
                  ? 'border-[#2563EB] text-[#2563EB] bg-[#2563EB] bg-opacity-[8%]'
                  : 'text-[#182B50] bg-[#F9FAFC]',
              ]"
              @click="handleType(item)"
            >
              <span class="text-sm">{{ item.label }}</span>
            </div>
          </template>
        </div>
      </el-form-item>
      <el-form-item v-if="['date'].includes(widgetForm.type)" :label="$t('form.select_placeholder')">
        <el-select v-model="widgetForm.date_format" class="w-full" size="large" placeholder="请选择">
          <el-option label="年" value="y" />
          <el-option label="年/月" value="y-m" />
          <el-option label="年/月/日" value="y-m-d" />
          <el-option label="时/分" value="h-m" />
          <el-option label="时间范围" value="daterange" />
        </el-select>
      </el-form-item>
      <el-form-item v-if="type === 'input'" :label="$t('agent.variable_desc')">
        <el-input
          v-model="widgetForm.desc"
          size="large"
          maxlength="1000"
          show-word-limit
          :placeholder="$t('form.input_placeholder')"
        />
      </el-form-item>
      <el-form-item v-if="type === 'input'" :label="$t('agent.variable_required')">
        <el-switch v-model="widgetForm.required" />
      </el-form-item>
      <el-form-item
        v-if="['text', 'textarea'].includes(widgetForm.type) && type === 'input'"
        :label="$t('agent.variable_max_length')"
      >
        <el-input-number
          v-model="widgetForm.max_length"
          style="width: 100%"
          class="el-input-number--left"
          :precision="0"
          :controls="false"
          size="large"
          :maxlength="256"
          :placeholder="$t('form.input_placeholder')"
        />
      </el-form-item>
      <el-form-item
        v-if="['text', 'textarea'].includes(widgetForm.type) && type === 'input'"
        :label="$t('agent.variable_show_word_limit')"
      >
        <el-switch v-model="widgetForm.show_word_limit" />
      </el-form-item>

      <el-form-item v-if="widgetForm.type === 'select'" :label="$t('agent.variable_options')">
        <div class="flex flex-col gap-3 w-full">
          <template v-for="(item, index) in widgetForm.options" :key="item.value">
            <div class="flex items-center">
              <el-input v-model="item.label" class="flex-1" size="large" :placeholder="$t('form.input_placeholder')" />
              <div class="px-2">
                <el-icon class="cursor-pointer" color="#999999" @click="handleDelOption(index)">
                  <Delete />
                </el-icon>
              </div>
            </div>
          </template>
        </div>
        <div
          class="w-full h-10 leading-10 rounded text-center border border-dashed border-[#DCDFE6] cursor-pointer text-sm text-[#182B50] text-opacity-80 mt-3"
          @click="handleAddOption"
        >
          + {{ $t('action.add') }}
        </div>
      </el-form-item>
      <el-form-item v-if="widgetForm.type === 'select'" label="模式">
        <el-radio-group v-model="widgetForm.multiple">
          <el-radio size="large" :label="false"> 单选 </el-radio>
          <el-radio size="large" :label="true"> 多选 </el-radio>
        </el-radio-group>
      </el-form-item>
      <template v-if="widgetForm.type === 'file'">
        <el-form-item label="上传文件类型">
          <el-select
            v-model="widgetForm.file_type"
            class="w-full"
            size="large"
            placeholder="请选择"
            @change="handleFileTypeChange($event, widgetForm)"
          >
            <el-option label="不限格式" value="all" />
            <el-option label="自定义" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item
          v-if="widgetForm.file_type === 'custom'"
          label="支持文件格式"
          required
          :rules="generateInputRules({ message: 'form.select_placeholder' })"
        >
          <el-select v-model="widgetForm.file_accept" multiple class="w-full" size="large" placeholder="请选择">
            <el-option label="docx" value="docx" />
            <el-option label="pdf" value="pdf" />
            <el-option label="xlsx" value="xlsx" />
            <el-option label="csv" value="csv" />
            <el-option label="txt" value="txt" />
            <el-option label="png" value="png" />
            <el-option label="jpg" value="jpg" />
            <el-option label="bmp" value="bmp" />
            <el-option label="markdown" value="md" />
            <el-option label="tiff" value="tiff" />
            <el-option label="html" value="html" />
          </el-select>
        </el-form-item>
        <el-form-item label="单个文件上限">
          <div class="w-full flex items-center gap-5 overflow-hidden">
            <div class="flex-1">
              <el-slider v-model="widgetForm.file_size" :min="1" :max="300" />
            </div>
            <span class="text-sm text-[#182B50]">{{ widgetForm.file_size }}M</span>
          </div>
        </el-form-item>
        <el-form-item label="上传最大数量">
          <div class="flex items-center gap-2">
            <el-input-number
              v-model="widgetForm.file_limit"
              :precision="0"
              :min="1"
              :max="6"
              size="large"
              placeholder="请输入"
            />
            <span class="text-sm text-[#182B50]">个</span>
          </div>
        </el-form-item>
      </template>
    </el-form>
    <template #footer>
      <el-button size="large" @click="visible = false"> {{ $t('action.cancel') }} </el-button>
      <el-button v-debounce type="primary" size="large" @click="handleSave"> {{ $t('action.save') }} </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Delete } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import { generateInputRules } from '@/utils/form-rule'
import { generateRandomId } from '@/utils'
import { inputTypeList, outputTypeList, outputDefaultField } from '@/constants/agent'

const props = defineProps<{
  type: 'input' | 'output'
  agentType: string
}>()

const emit = defineEmits<{
  (e: 'save', value: Agent.Field): void
}>()

const formRef = ref<FormInstance>()
const visible = ref(false)

const widgetForm = ref<Agent.Field>({
  ...outputDefaultField,
  file_accept: ['doc', 'docx', 'pdf', 'xlsx', 'csv', 'txt', 'png', 'jpg', 'bmp', 'md', 'tiff', 'html'],
})

const typeList = computed(() => {
  return (props.type === 'input' ? inputTypeList : outputTypeList).filter(item =>
    item.allowed ? item.allowed.includes(props.agentType) : true
  )
})

const typeLabel = computed(() => {
  return [...inputTypeList, ...outputTypeList].find(item => item.type === widgetForm.value.type)?.label
})

const handleType = item => {
  widgetForm.value.type = item.type
}

const handleAddOption = () => {
  widgetForm.value.options.push({
    id: '',
    label: '',
    value: '',
  })
}

const handleDelOption = index => {
  widgetForm.value.options.splice(index, 1)
}

const handleFileTypeChange = (fileType: string, form: Agent.Field) => {
  if (fileType === 'all')
    form.file_accept = ['doc', 'docx', 'pdf', 'xlsx', 'csv', 'txt', 'png', 'jpg', 'bmp', 'md', 'tiff', 'html']
}
const handleSave = async () => {
  const is_valid = await formRef.value?.validate()
  if (!is_valid) return
  emit('save', {
    ...widgetForm.value,
    id: widgetForm.value.id || generateRandomId(10),
  })
  visible.value = false
}

defineExpose({
  open(data: Agent.Field = {} as Agent.Field) {
    widgetForm.value.id = data.id || ''
    widgetForm.value.variable = data.variable || ''
    widgetForm.value.label = data.label || ''
    widgetForm.value.type = data.type || typeList.value[0].type
    widgetForm.value.desc = data.desc || ''
    widgetForm.value.required = data.required || false
    widgetForm.value.max_length = data.max_length || 0
    widgetForm.value.show_word_limit = data.show_word_limit || false
    widgetForm.value.options = data.options || []
    widgetForm.value.multiple = data.multiple || false
    widgetForm.value.file_type = data.file_type || 'all'
    widgetForm.value.file_accept = data.file_accept || []
    widgetForm.value.file_limit = data.file_limit || 1
    widgetForm.value.file_size = data.file_size || 30
    widgetForm.value.is_system = data.is_system || false

    visible.value = true
  },
})
</script>

<style scoped>
/* ::v-deep(.el-input-number .el-input__inner) {
  text-align: left;
} */
</style>
