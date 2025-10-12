<script setup lang="ts">
import { Plus } from '@element-plus/icons-vue'
import { reactive, ref } from 'vue'
import UploadImage from '@/components/Upload/image.vue'

import { paymentApi } from '@/api/modules/payment'
import { prepareSavePaymentSettingData } from '@/api/modules/payment/transform'
import { PAYMENT_TYPE } from '@/constants/payment'

const emits = defineEmits(['success'])

const form_ref = ref()
const visible = ref(false)
const form = reactive({
  pay_qrcode: '',
  pay_desc: ''
})
const origin_data = ref({})
const submitting = ref(false)

const textValidator = ({ rule, value, callback, message } = {}) => {
  value = (value || '').trim()
  if (!value) {
    callback(new Error(window.$t(message)))
    return
  }
  callback()
}

const open = ({ data = {} } = {}) => {
  const extra_config = data.extra_config || {}
  form.pay_qrcode = extra_config.pay_qrcode || ''
  form.pay_desc = extra_config.pay_desc || ''
  origin_data.value = data
  visible.value = true
}
const reset = () => {
  form.pay_qrcode = ''
  form.pay_desc = ''
}

const close = () => {
  visible.value = false
  reset()
}
const handleConfirm = async () => {
  const valid = await form_ref.value.validate()
  if (!valid) return
  submitting.value = true
  const { preparedData, pay_setting_id } = prepareSavePaymentSettingData({
    pay_setting_id: origin_data.value.pay_setting_id,
    pay_config: {},
    extra_config: {
      pay_qrcode: form.pay_qrcode,
      pay_desc: form.pay_desc
    },
    pay_type: PAYMENT_TYPE.MANUAL
  })
  
  await paymentApi
    .savePaymentSetting({ pay_setting_id, ...preparedData })
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
  reset
})
</script>

<template>
  <ElDialog
    v-model="visible"
    :title="$t('payment.type.manual')"
    :close-on-click-modal="false"
    width="700px"
    destroy-on-close
    append-to-body
    @close="close"
  >
    <ElForm ref="form_ref" :model="form" label-position="top">
      <ol class="w-full flex flex-col gap-3 bg-[#F6F9FC] p-5 mb-6 box-border text-sm text-[#4F5052]">
        <li>{{ $t('payment.manual_tip') }}</li>
      </ol>
      <ElFormItem
        :label="$t('payment.manual.qrcode')"
        prop="pay_qrcode"
        :rules="[
          {
            validator: (_rule, value, callback) => textValidator({ rule: _rule, value, callback, message: 'payment.manual.qrcode_placeholder' }),
            trigger: 'blur'
          }
        ]"
      >
        <UploadImage v-model="form.pay_qrcode" class="!w-[148px] !h-[148px] border bg-[#F7F8FA] rounded-md overflow-hidden">
          <div class="w-full h-full flex flex-col items-center justify-center gap-2">
            <img v-if="form.pay_qrcode" :src="form.pay_qrcode" class="absolute top-0 left-0 w-full h-full object-cover" />
            <template v-else>
              <ElIcon color="#9A9A9A" size="20">
                <Plus />
              </ElIcon>
              <span class="text-[#9A9A9A] text-sm">{{ $t('action_upload') }}</span>
            </template>
          </div>
        </UploadImage>
      </ElFormItem>
      <ElFormItem :label="$t('payment.manual.desc')" prop="pay_desc">
        <ElInput
          v-model="form.pay_desc"
          size="large"
          resize="none"
          type="textarea"
          :rows="4"
          :maxlength="30"
          show-word-limit
          clearable
          :placeholder="$t('payment.manual.desc_placeholder')"
        />
      </ElFormItem>
    </ElForm>
    <template #footer>
      <div class="py-4 flex items-center justify-center">
        <ElButton class="w-[96px] h-[36px]" type="primary" :loading="submitting" @click="handleConfirm">
          {{ $t('action_confirm') }}
        </ElButton>
        <ElButton class="w-[96px] h-[36px] text-[#1D1E1F]" type="info" plain @click.stop="close">
          {{ $t('action_cancel') }}
        </ElButton>
      </div>
    </template>
  </ElDialog>
</template>

<style scoped>
::v-deep(.el-upload-dragger) {
  position: relative;
  width: 148px;
  height: 148px;
  box-sizing: border-box;
  border-style: solid;
  background-color: #f7f9fa;
}
</style>
