<template>
  <ElDialog
    v-model="visible"
    :title="$t(editable ? 'action_edit' : 'action_add')"
    :close-on-click-modal="false"
    width="600px"
    destroy-on-close
    append-to-body
    @close="close"
  >
    <ElForm ref="form_ref" :model="form" label-position="top">
      <h1 class="text-sm text-[#1D1E1F]">
        {{ $t('user_info') }}
      </h1>
      <!-- :rules="generateInputRules({ message: 'form_upload_placeholder' })" -->
      <ElFormItem class="mt-5" prop="avatar">
        <template #label>
          <span class="text-[#9A9A9A]">{{ $t('avatar') }}</span>
        </template>
        <UploadImage v-model="form.avatar" />
      </ElFormItem>
      <ElFormItem
        :label="$t('user')"
        prop="nickname"
        :rules="generateInputRules({ message: 'form_upload_placeholder' })"
      >
        <ElInput
          v-model="form.nickname"
          size="large"
          maxlength="20"
          show-word-limit
          :placeholder="$t('form_input_placeholder')"
        />
      </ElFormItem>
      <!-- <ElFormItem
        :label="$t('login_password')"
        prop="password"
        :rules="
          !form.password.length
            ? []
            : [
                ...generateInputRules({ message: 'login.password_placeholder', validator: ['password'] }),
                {
                  min: 8,
                  max: 20,
                  message: $t('login.password_length'),
                  trigger: 'blur',
                },
              ]
        "
      >
        <ElInput
          v-model="form.password"
          size="large"
          :placeholder="$t('empty_to_not_change')"
          clearable
          @blur="form_ref.validateField('password')"
        />
      </ElFormItem> -->
      <ElFormItem
        :label="$t('subscription.title')"
        prop="group_id"
        :rules="[
          {
            validator: (rule, value, callback) => {
              if (!value) return callback(new Error($t('form_select_placeholder')))
              callback()
            },
            trigger: 'blur',
          },
        ]"
      >
        <ElSelect v-model="form.group_id" size="large" :placeholder="$t('form_input_placeholder')">
          <ElOption v-for="item in subscription_options" :key="item.value" :label="item.label" :value="item.value" />
        </ElSelect>
      </ElFormItem>
      <ElFormItem :label="$t('subscription.end_at')" prop="expired_time" :rules="[]">
        <ElDatePicker
          v-model="form.expired_time"
          format="YYYY-MM-DD HH:mm"
          type="datetime"
          value-format="YYYY-MM-DD HH:mm"
          size="large"
          :placeholder="$t('permanent_effect')"
        />
      </ElFormItem>
      <h1 class="text-sm text-[#1D1E1F] mt-6">
        {{ $t('bind_account') }}
      </h1>
      <ul class="mb-4">
        <li class="flex items-center gap-2 mt-5 text-sm text-[#1D1E1F]">
          <SvgIcon class="flex-none" name="mobile-circle" />
          <div class="flex-none w-[88px]">
            {{ $t('mobile_v2') }}
          </div>
          <div class="text-[#9A9A9A]">
            {{ origin_data.mobile || $t('not_bound') }}
          </div>
        </li>
        <li class="flex items-center gap-2 mt-5 text-sm text-[#1D1E1F]">
          <SvgIcon class="flex-none" name="email-circle" />
          <div class="flex-none w-[88px]">
            {{ $t('email') }}
          </div>
          <div class="text-[#9A9A9A]">
            {{ origin_data.email || $t('not_bound') }}
          </div>
        </li>
        <li v-if="false" class="flex items-center gap-2 mt-5 text-sm text-[#1D1E1F]">
          <SvgIcon class="flex-none" name="wechat-circle" />
          <div class="flex-none w-[88px]">
            {{ $t('open_id') }}
          </div>
          <div class="text-[#9A9A9A]">
            {{ origin_data.open_id || $t('not_bound') }}
          </div>
        </li>
        <li v-if="false" class="flex items-center gap-2 mt-5 text-sm text-[#1D1E1F]">
          <SvgIcon class="flex-none" name="google" />
          <div class="flex-none w-[88px]">
            {{ $t('google_account') }}
          </div>
          <div class="text-[#9A9A9A]">
            {{ origin_data.google_account || $t('not_bound') }}
          </div>
        </li>
      </ul>
    </ElForm>
    <template #footer>
      <div class="py-4 flex items-center justify-center">
        <ElButton class="w-[96px] h-[36px]" type="primary" :loading="submitting" @click="handleSave">
          {{ $t('action_save') }}
        </ElButton>
        <ElButton class="w-[96px] h-[36px] text-[#1D1E1F]" type="info" plain @click.stop="close">
          {{ $t('action_cancel') }}
        </ElButton>
      </div>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import UploadImage from '@/components/Upload/image.vue'

import { getTimeStamp } from '@/utils/moment'
import { generateInputRules } from '@/utils/form-rule'
import { useUserStore } from '@/stores'

const emits = defineEmits<{
  (e: 'success'): void
}>()

const user_store = useUserStore()

const form_ref = ref()
const visible = ref(false)
const editable = ref(false)
const form = reactive({
  avatar: '',
  nickname: '',
  password: '',
  group_id: 0,
  expired_time: '',
})
const origin_data = ref({})
const subscription_options = ref([])
const submitting = ref(false)

const open = ({ data = {}, subscription_options: _subscription_options = [] } = {}) => {
  form.avatar = data.avatar || ''
  form.nickname = data.nickname || ''
  form.password = data.password || ''
  form.group_id = data.group_id || (_subscription_options[0] || {}).value || 0
  form.expired_time = data.expired_time || ''
  editable.value = !!data.user_id
  origin_data.value = data
  subscription_options.value = _subscription_options.filter(item => item.value !== 0)
  if (!_subscription_options.find(item => item.value === form.group_id)) form.group_id = ''
  visible.value = true
}
const close = () => {
  visible.value = false
  reset()
}
const reset = () => {
  form.avatar = ''
  form.nickname = ''
  form.password = ''
}
const handleSave = async () => {
  const valid = await form_ref.value.validate()
  if (!valid) return
  submitting.value = true
  await user_store
    .save({
      data: {
        user_id: origin_data.value.user_id,
        avatar: form.avatar,
        nickname: form.nickname,
        password: form.password,
        group_id: form.group_id,
        expired_time: (form.expired_time && getTimeStamp(form.expired_time)) || 0,
      },
    })
    .finally(() => {
      submitting.value = false
    })
  emits('success')
  ElMessage.success(window.$t('action_save_success'))
  close()
}

defineExpose({
  open,
  close,
  reset,
})
</script>

<style scoped></style>
