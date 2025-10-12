<template>
  <div class="flex items-center gap-4">
    <ElFormItem
      class="flex-1"
      :label="$t('name')"
      prop="name"
      :rules="generateInputRules({ message: 'form_input_placeholder' })"
    >
      <ElInput
        v-model="store.form_data.name"
        size="large"
        show-word-limit
        maxlength="20"
        :placeholder="$t('form_input_placeholder')"
      />
    </ElFormItem>
    <ElFormItem
      :label="$t('group')"
      class="flex-1"
      prop="group_id"
      :rules="generateInputRules({ message: 'form_select_placeholder' })"
    >
      <ElSelect v-model="store.form_data.group_id" size="large">
        <ElOption v-for="item in store.group_options" :key="item.value" :value="item.value" :label="item.label" />
      </ElSelect>
    </ElFormItem>
  </div>
  <ElFormItem :label="$t('description')">
    <ElInput
      v-model="store.form_data.description"
      type="textarea"
      :rows="3"
      resize="none"
      show-word-limit
      maxlength="200"
    />
  </ElFormItem>
  <ElFormItem :label="$t('avatar')" prop="logo" :rules="generateInputRules({ message: 'form_upload_placeholder' })">
    <UploadImage v-model="store.form_data.logo" class="w-12 h-12" />
  </ElFormItem>
  <ElFormItem
    :label="$t('action_sort')"
    prop="sort"
    :rules="generateInputRules({ message: 'form_input_placeholder', validator: ['number'], min: 0, max: 99999999 })"
  >
    <ElInputNumber
      v-model="store.form_data.sort"
      class="!w-[300px] el-input-number--left"
      size="large"
      :controls="false"
      :precision="0"
      :min="0"
      :max="99999999"
      :placeholder="$t('form_input_placeholder')"
      @keydown="$numberInputKeydownHandler"
    />
    <div class="w-full text-sm text-[#9A9A9A]">
      {{ $t('module.agent_sort_desc') }}
    </div>
  </ElFormItem>
</template>

<script setup lang="ts">
import { useAgentFormStore } from '../store'
import UploadImage from '@/components/Upload/image.vue'

import { generateInputRules } from '@/utils/form-rule'

const store = useAgentFormStore()
</script>

<style lang="scss" scoped>
::v-deep(.el-input__inner) {
  text-align: left;
}
</style>
