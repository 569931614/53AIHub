<template>
  <ElForm ref="form_ref" v-loading="loading" class="relative max-w-[440px] w-full" :model="form" label-position="top">
    <h4 class="text-3xl text-[#1D1E1F] font-bold text-center mb-10">
      {{ $t('apply.create_site') }}
    </h4>
    <!-- 暂时去掉，现在后台自动审核成功了。需要跳过审核的 -->
    <ElSteps v-if="false" class="mb-8" :active="active_step" align-center>
      <ElStep
        v-for="(step_label, step_index) in ['module.website_info', 'apply.waiting_audit', 'apply.create_success']"
        :key="step_index"
      >
        <template #icon>
          <div class="w-9 h-9 bg-[#F2F3F3] rounded-full flex items-center justify-center">
            <div
              class="w-7 h-7 rounded-full text-white text-sm flex items-center justify-center"
              :class="[
                step_index < active_step
                  ? 'bg-[#82A2F7]'
                  : active_step === step_index
                    ? 'bg-[#3664F0]'
                    : 'bg-[#CFD1D6]',
              ]"
            >
              {{ step_index + 1 }}
            </div>
          </div>
        </template>
        <template #title>
          <span
            :class="[
              step_index < active_step
                ? 'text-[#B3C7FA]'
                : active_step === step_index
                  ? 'text-[#3664EF]'
                  : 'text-[#9A9A9A]',
            ]"
          >
            {{ $t(step_label) }}
          </span>
        </template>
      </ElStep>
    </ElSteps>

    <template v-if="active_step === 0">
      <ElFormItem
        class="is-required"
        prop="website_name"
        :rules="generateInputRules({ message: 'login.website_name_placeholder' })"
      >
        <template #label>
          <span class="text-[#1D1E1F]">{{ $t('login.website_name') }}</span>
        </template>
        <ElInput
          v-model="form.website_name"
          style="--el-input-bg-color: #f1f2f3; --el-input-border-color: transparent; --el-input-height: 44px"
          size="large"
          :placeholder="$t(`login.website_name_placeholder`)"
          clearable
        />
      </ElFormItem>
      <ElFormItem
        class="is-required"
        prop="contact_name"
        :rules="generateInputRules({ message: 'login.contact_name_placeholder' })"
      >
        <template #label>
          <span class="text-[#1D1E1F]">{{ $t('login.contact_name') }}</span>
        </template>
        <ElInput
          v-model="form.contact_name"
          style="--el-input-bg-color: #f1f2f3; --el-input-border-color: transparent; --el-input-height: 44px"
          size="large"
          :placeholder="$t(`login.contact_name_placeholder`)"
          clearable
        />
      </ElFormItem>
      <ElButton
        type="primary"
        round
        class="w-full mt-6 !h-10"
        :loading="submitting"
        :disabled="!form.website_name || !form.contact_name"
        @click.stop="onNextStep"
      >
        {{ $t('action_confirm') }}
      </ElButton>
    </template>
    <template v-else-if="active_step === 1">
      <div
        class="h-[424px] p-10 box-border bg-[#EFF9FF] rounded-lg flex flex-col items-center justify-center text-center"
      >
        <div class="flex items-center justify-center gap-2">
          <ElIcon color="#4CBF65" size="28">
            <CircleCheckFilled />
          </ElIcon>
          <span class="text-[#1D1E1F] text-2xl font-bold">{{ $t('apply.waiting_audit') }}</span>
        </div>
        <div class="text-[#666] text-sm mt-4">
          {{ $t('apply_success_desc') }}
        </div>
        <img class="w-[148px] object-contain mt-14" src="//chat.53ai.com/images/upgrade-qrcode.png" />
      </div>
    </template>
    <template v-else />
    <ElButton
      type="primary"
      text
      class="!p-0 mt-4 !mx-auto relative left-1/2 -translate-x-1/2 !bg-transparent"
      @click="onLogin"
    >
      {{ $t('login.back_to_login') }}
    </ElButton>
  </ElForm>
</template>

<script setup lang="ts">
import { CircleCheckFilled } from '@element-plus/icons-vue'

import { nextTick, onMounted, onUnmounted, reactive, ref } from 'vue'
import { useEnterpriseStore, useUserStore } from '@/stores'
import { generateInputRules } from '@/utils/form-rule'
import eventBus from '@/utils/event-bus'

const emits = defineEmits<{
  (e: 'login'): void
}>()

const user_store = useUserStore()
const enterprise_store = useEnterpriseStore()
const form_ref = ref()
const active_step = ref(0)
const form = reactive({
  website_name: '',
  contact_name: '',
})
const submitting = ref(false)
const loading = ref(false)
onMounted(async () => {
  const { access_token } = user_store.info
  if (access_token) {
    loading.value = true
    const { list = [] } = await enterprise_store.loadListData({ data: { status: 0 } }).finally(() => {
      loading.value = false
    })
    if (list.length > 0) active_step.value = 1
  }
  eventBus.on('language-change', onLanguageChange)
})
onUnmounted(() => {
  eventBus.off('language-change', onLanguageChange)
})

const onLanguageChange = () => {
  if (form_ref.value) form_ref.value.clearValidate()
  // form_ref.value.validate()
}
const reset = () => {
  form.website_name = ''
  form.contact_name = ''
}
const onNextStep = async () => {
  const valid = await form_ref.value.validate()
  if (!valid) return
  switch (active_step.value) {
    case 0:
      submitting.value = true
      await enterprise_store
        .apply({
          data: {
            contact_name: form.contact_name,
            enterprise_name: form.website_name,
            phone: user_store.info.username,
            email: '',
          },
        })
        .then(() => {
          nextTick(() => {
            emits('login')
          })
        })
        .finally(() => {
          submitting.value = false
        })
      // active_step.value++
      ElMessage.success($t('apply.create_success'))
      break
    case 1:
      break
    default:
      break
  }
  reset()
}
const onLogin = () => {
  emits('login')
}

defineExpose({
  reset,
})
</script>

<style scoped lang="scss">
::v-deep(.el-step__line) {
  left: 70% !important;
  right: -30% !important;
  background-color: #e9ebee !important;
}
::v-deep(.el-step__title) {
  margin-top: 8px;
  font-weight: 400 !important;
}
::v-deep(.el-step__head.is-finish .el-step__line) {
  background-color: #3664ef !important;
  // background-color: #82A2F6 !important;
}
::v-deep(.el-step__head.is-process .el-step__line) {
  // background-color: #3664EF !important;
}
</style>
