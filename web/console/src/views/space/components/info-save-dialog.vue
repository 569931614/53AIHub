<template>
  <el-dialog
    v-model="visible"
    class="el-dialog--footer-center"
    :title="originData.id ? $t('action.edit') : $t('action.create')"
    destroy-on-close
    :close-on-click-modal="false"
    width="500px"
  >
    <el-form ref="formRef" :model="formData" label-width="108px" label-position="left">
      <el-form-item :label="$t('space.icon')" prop="icon" :rules="rules.icon" required>
        <UploadImage v-model="formData.icon" class="w-12 h-12" />
      </el-form-item>
      <el-form-item :label="$t('space.name')" prop="name" :rules="rules.name" required>
        <el-input
          v-model="formData.name"
          size="large"
          :placeholder="$t('space.name_placeholder')"
          maxlength="20"
          show-word-limit
        />
      </el-form-item>
      <el-form-item :label="$t('space.description')" prop="description">
        <el-input
          v-model="formData.description"
          type="textarea"
          resize="none"
          :rows="5"
          :placeholder="$t('space.description_placeholder')"
          maxlength="200"
          show-word-limit
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <ElButton v-debounce class="w-24 h-9" type="primary" @click.stop="handleSave">
        {{ $t('action_save') }}
      </ElButton>
      <ElButton class="w-24 h-9 text-[#1D1E1F]" type="info" plain @click.stop="close">
        {{ $t('action_cancel') }}
      </ElButton>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import UploadImage from '@/components/Upload/image.vue'

import { spacesApi } from '@/api/modules/spaces/index'
import type { SpaceItem, SpaceCreateRequest } from '@/api/modules/spaces/types'

import { generateInputRules } from '@/utils/form-rule'

const emit = defineEmits(['refresh'])

const formRef = ref()
const visible = ref(false)
const originData = ref<Partial<SpaceItem>>({})

const formData = reactive<SpaceCreateRequest>({
  name: '',
  description: '',
  icon: '',
})

const rules = {
  name: generateInputRules({ validator: ['text'] }),
  icon: generateInputRules({ validator: ['image'] }),
}

const open = (data: SpaceItem = {} as SpaceItem) => {
  formData.name = data.name || ''
  formData.description = data.description || ''
  formData.icon = data.icon || ''
  originData.value = data || {}
  visible.value = true
}

const close = () => {
  visible.value = false
}

const handleSave = async () => {
  const valid = await formRef.value.validate()
  if (!valid) return

  if (originData.value.id) {
    await spacesApi.update(originData.value.id, formData)
    ElMessage.success(window.$t('message_status.save_success'))
  } else {
    await spacesApi.create(formData)
    ElMessage.success(window.$t('message_status.create_success'))
  }
  emit('refresh')
  close()
}

defineExpose({
  open,
  close,
})
</script>
