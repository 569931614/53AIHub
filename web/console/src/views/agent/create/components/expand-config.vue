<template>
  <template v-if="store.support_file || store.support_image">
    <div class="text-base text-[#1D1E1F] font-medium mb-4 mt-10">
      {{ $t('expand_setting') }}
    </div>
    <div v-if="store.support_file" class="flex items-center gap-2">
      <div class="flex-1">
        <div class="text-sm text-[#1D1E1F]">
          {{ $t('agent_file_parse') }}
        </div>
      </div>
      <div class="flex-none text-sm text-[#1D1E1F]">
        {{ store.form_data.settings.file_parse.enable ? $t('action_open') : $t('action_close') }}
        <el-switch v-model="store.form_data.settings.file_parse.enable" />
      </div>
    </div>
    <div v-if="store.support_image" class="flex items-center gap-2 mt-4">
      <div class="flex-1">
        <div class="text-sm text-[#1D1E1F]">
          {{ $t('agent_image_parse') }}
        </div>
      </div>
      <div class="flex-none text-sm text-[#1D1E1F]">
        {{ store.form_data.settings.image_parse.enable ? $t('action_open') : $t('action_close') }}
        <el-switch v-model="store.form_data.settings.image_parse.enable" />
      </div>
    </div>
  </template>

  <!-- Prompt 类型上下文轮数设置（总在 Prompt 下显示，不依赖文件/图片开关） -->
  <div v-if="store.agent_type === 'prompt'" class="mt-10">
    <div class="text-base text-[#1D1E1F] font-medium mb-4">
      {{ tOr('context_setting', '上下文设置') }}
    </div>
    <div class="flex items-center gap-2">
      <div class="flex-1">
        <div class="text-sm text-[#1D1E1F]">
          {{ tOr('context_rounds', '上下文轮数') }}
        </div>
      </div>
      <div class="flex-none text-sm text-[#1D1E1F]">
        <el-input-number
          v-model="store.form_data.configs.chat.history_pairs"
          class="!w-[160px] el-input-number--left"
          size="large"
          :controls="false"
          :precision="0"
          :min="1"

        />
      </div>
    </div>
    <div class="text-xs text-[#9A9A9A] mt-1">
      {{ tOr('context_rounds_desc', '每次请求回带的历史问答对数量，建议 6。') }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useAgentFormStore } from '../store';

const store = useAgentFormStore()
const { t } = useI18n()

const tOr = (key: string, fallback: string) => {
  const v = t(key) as unknown as string
  return v && v !== key ? v : fallback
}
</script>

<style scoped></style>
