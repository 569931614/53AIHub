<template>
  <ElDialog
    v-model="visible"
    :title="$t(editable ? 'action_edit' : 'action_create')"
    :close-on-click-modal="false"
    width="680px"
    append-to-body
    @close="close"
  >
    <ElForm ref="form_ref" :model="form" label-position="left" label-width="78px">
      <ElFormItem :label="$t('module.nav_type')">
        <ElRadioGroup v-model="form.type" size="large">
          <ElRadio v-for="item in nav_type_options" :key="item.value" :label="item.value">{{ item.label }}</ElRadio>
        </ElRadioGroup>
      </ElFormItem>
      <ElFormItem v-if="form.type !== NAVIGATION_TYPE.AGENT" :label="$t('module.nav_name')">
        <ElSelect
          v-model="form.nav_id"
          size="large"
          :placeholder="$t('module.nav_name_placeholder')"
          @change="handleNavChange"
        >
          <ElOption v-for="item in nav_options" :key="item.value" :label="item.label" :value="item.value" />
        </ElSelect>
      </ElFormItem>
      <ElFormItem v-else :label="$t('action_select')">
        <div class="w-full flex items-center gap-2">
          <ElSelect
            v-model="form.agent_class_id"
            class="flex-none !w-[160px]"
            size="large"
            :placeholder="$t('module.nav_agent_class_placeholder')"
          >
            <ElOption
              v-for="item in nav_agent_class_options"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </ElSelect>
          <ElSelect
            v-model="form.agent_id"
            class="flex-1"
            size="large"
            :placeholder="$t('module.nav_agent_placeholder')"
            @change="handleAgentChange"
          >
            <ElOption v-for="item in nav_agent_options" :key="item.value" :label="item.label" :value="item.value" />
          </ElSelect>
        </div>
      </ElFormItem>
      <ElFormItem :label="$t('module.nav_url')">
        <ElInput v-model="form.url" disabled size="large" :placeholder="$t('form_select_placeholder')" />
      </ElFormItem>
      <ElFormItem :label="$t('module.nav_target')">
        <ElSelect
          v-model="form.target"
          :disabled="form.type !== NAVIGATION_TYPE.AGENT"
          size="large"
          :placeholder="$t('module.nav_target_placeholder')"
        >
          <ElOption v-for="item in nav_target_options" :key="item.value" :label="item.label" :value="item.value" />
        </ElSelect>
      </ElFormItem>
    </ElForm>
    <template #footer>
      <div class="py-4 flex items-center justify-center">
        <ElButton class="w-[96px] h-[36px]" type="primary" @click="handleConfirm">{{ $t('action_confirm') }}</ElButton>
        <ElButton class="w-[96px] h-[36px] text-[#1D1E1F]" type="info" plain @click.stop="close">
          {{ $t('action_cancel') }}
        </ElButton>
      </div>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
import { reactive, ref, computed } from 'vue'

import { useEnterpriseStore } from '@/stores/modules/enterprise'
import { useAgentStore } from '@/stores/modules/agent'
import {
  NAVIGATION_TYPE,
  NAVIGATION_TARGET,
  NAVIGATION_TYPE_LABEL_MAP,
  NAVIGATION_TARGET_LABEL_MAP,
} from '@/constants/navigation'

const enterpriseStore = useEnterpriseStore()
const agentStore = useAgentStore()

// Types
type NavigationFormData = {
  type: string
  nav_id: string
  url: string
  target: string
  agent_class_id: string
  agent_id: string
  name: string
}

// Emits
const emits = defineEmits<{
  confirm: [data: { data: NavigationFormData }]
}>()

// Refs
const visible = ref(false)
const editable = ref(false)
const form_ref = ref()

// Form data
const form = reactive<NavigationFormData>({
  type: NAVIGATION_TYPE.HOMEPAGE,
  nav_id: '',
  url: '',
  target: NAVIGATION_TARGET.BLANK,
  agent_class_id: '',
  agent_id: '',
  name: '',
})

// Options
const nav_type_options = [
  { label: NAVIGATION_TYPE_LABEL_MAP.get(NAVIGATION_TYPE.HOMEPAGE), value: NAVIGATION_TYPE.HOMEPAGE },
  { label: NAVIGATION_TYPE_LABEL_MAP.get(NAVIGATION_TYPE.AGENT), value: NAVIGATION_TYPE.AGENT },
]

const nav_target_options = [
  { label: NAVIGATION_TARGET_LABEL_MAP.get(NAVIGATION_TARGET.SELF), value: NAVIGATION_TARGET.SELF },
  { label: NAVIGATION_TARGET_LABEL_MAP.get(NAVIGATION_TARGET.BLANK), value: NAVIGATION_TARGET.BLANK },
]

// Computed
const nav_options = computed(() => {
  const { nav_list = [] } = enterpriseStore.enterpriseInfo
  return nav_list.map(item => ({
    label: item.name,
    value: item.navigation_id,
    url: item.jump_path,
  }))
})

const nav_agent_class_options = computed(() => {
  const { agent_class_list = [] } = agentStore.agentInfo
  return agent_class_list.map(item => ({
    label: item.name,
    value: item.agent_class_id,
  }))
})

const nav_agent_options = computed(() => {
  const { agent_list = [] } = agentStore.agentInfo
  return agent_list
    .filter(item => item.agent_class_id === form.agent_class_id)
    .map(item => ({
      label: item.name,
      value: item.agent_id,
      url: `#/agent/${item.agent_id}`,
    }))
})

// Methods
const reset = () => {
  Object.assign(form, {
    type: NAVIGATION_TYPE.HOMEPAGE,
    nav_id: '',
    url: '',
    target: NAVIGATION_TARGET.BLANK,
    agent_class_id: '',
    agent_id: '',
    name: '',
  })
}

const open = ({ data = {} } = {}) => {
  visible.value = true
  editable.value = !!data.navigation_id

  Object.assign(form, {
    type: data.type || NAVIGATION_TYPE.HOMEPAGE,
    nav_id: data.nav_id || '',
    url: data.url || '',
    target: data.target || NAVIGATION_TARGET.BLANK,
    agent_class_id: data.agent_class_id || '',
    agent_id: data.agent_id || '',
    name: data.name || '',
  })
}

const close = () => {
  visible.value = false
  reset()
}

const handleNavChange = () => {
  const selectedNav = nav_options.value.find(item => item.value === form.nav_id)
  if (selectedNav) {
    form.url = selectedNav.url
    form.name = selectedNav.label
  }
}

const handleAgentChange = () => {
  const selectedAgent = nav_agent_options.value.find(item => item.value === form.agent_id)
  if (selectedAgent) {
    form.url = selectedAgent.url
    form.name = selectedAgent.label
  }
}

const handleConfirm = async () => {
  try {
    const valid = await form_ref.value.validate()
    if (!valid) return

    ElMessage.success(window.$t('action_save_success'))
    emits('confirm', { data: form })
    close()
  } catch (error) {
    console.error('Form validation failed:', error)
  }
}

defineExpose({
  open,
  close,
  reset,
})
</script>

<style scoped lang="scss"></style>
