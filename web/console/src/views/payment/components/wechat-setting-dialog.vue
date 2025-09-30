<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import UploadCertificate from '@/components/Upload/certificate.vue'

import { generateInputRules } from '@/utils/form-rule'
import { paymentApi } from '@/api/modules/payment'
import { prepareSavePaymentSettingData } from '@/api/modules/payment/transform'
import { PAYMENT_TYPE } from '@/constants/payment'

const emits = defineEmits(['success'])

const form_ref = ref()
const visible = ref(false)
const form = reactive({
  mchId: '',
  appId: '',
  apiV3Key: '',
  serialNo: '',
  certPath: '',
  certName: '',
  privateKeyPath: '',
  privateKeyName: '',
  notifyUrl: '',
  platformCertPath: ''
})
const origin_data = ref({})
const submitting = ref(false)

const open = ({ data = {} } = {}) => {
  const config = data.pay_config || {}
  const extra_config = data.extra_config || {}
  form.mchId = config.mchId || ''
  form.appId = config.appId || ''
  form.apiV3Key = config.apiV3Key || ''
  form.serialNo = config.serialNo || ''
  form.certPath = config.certPath || ''
  form.certName = config.certName || extra_config.certName || ''
  form.privateKeyPath = config.privateKeyPath || ''
  form.privateKeyName = config.privateKeyName || extra_config.privateKeyName || ''
  form.notifyUrl = config.notifyUrl || ''
  form.platformCertPath = config.platformCertPath || ''
  origin_data.value = data
  visible.value = true
}
const reset = () => {
  form.mchId = ''
  form.appId = ''
  form.apiV3Key = ''
  form.serialNo = ''
  form.certPath = ''
  form.certName = ''
  form.privateKeyPath = ''
  form.privateKeyName = ''
  form.notifyUrl = ''
  form.platformCertPath = ''
}

const close = () => {
  visible.value = false
  reset()
}

const handleConfirm = async () => {
  const valid = await form_ref.value.validate()
  if (!valid) return
  submitting.value = true
  const pay_config = JSON.parse(JSON.stringify(form))
  const extra_config = {
    ...JSON.parse(JSON.stringify(origin_data.value.extra_config || {})),
    certName: form.certName,
    privateKeyName: form.privateKeyName
  }
  delete pay_config.certName
  delete pay_config.privateKeyName
  const { preparedData, pay_setting_id } = prepareSavePaymentSettingData({
    pay_setting_id: origin_data.value.pay_setting_id,
    pay_config,
    extra_config,
    pay_type: PAYMENT_TYPE.WECHAT
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

const guideVisible = ref(false)
const guideMode = ref('mch')
const guideTitle = computed(() => {
  const guideTitleMap = new Map([
    ['mch', 'wechat_payment.mch_guide.title'],
    ['app', 'wechat_payment.app_guide.title'],
    ['api', 'wechat_payment.api_guide.title'],
    ['cert', 'wechat_payment.cert_guide.title'],
    ['serial', 'wechat_payment.serial_guide.title']
  ])
  return guideTitleMap.get(guideMode.value)
})
const guideList = computed(() => {
  const guideListMap = new Map([
    [
      'mch',
      [
        {
          title: window.$t('wechat_payment.mch_guide.step_1', {
            url: '<a style="color: #586D9A;" href="https://pay.weixin.qq.com/" target="_blank">https://pay.weixin.qq.com/</a>'
          }),
          imageList: []
        },
        {
          title: window.$t('wechat_payment.mch_guide.step_2'),
          imageList: ['/images/wechat-payment/mch-guide-1.png']
        },
        {
          title: window.$t('wechat_payment.mch_guide.step_3'),
          imageList: ['/images/wechat-payment/mch-guide-2.png', '/images/wechat-payment/mch-guide-3.png']
        }
      ]
    ],
    [
      'app',
      [
        {
          title: window.$t('wechat_payment.app_guide.step_1'),
          imageList: ['/images/wechat-payment/app-guide-1.png']
        },
        {
          title: window.$t('wechat_payment.app_guide.step_2'),
          imageList: ['/images/wechat-payment/app-guide-2.png']
        },
        {
          title: window.$t('wechat_payment.app_guide.step_3', {
            url: '<a style="color: #586D9A;" href="https://mp.weixin.qq.com" target="_blank">https://mp.weixin.qq.com</a>'
          }),
          imageList: ['/images/wechat-payment/app-guide-3.png']
        }
      ]
    ],
    [
      'api',
      [
        {
          title: window.$t('wechat_payment.api_guide.step_1'),
          imageList: ['/images/wechat-payment/api-guide-1.png']
        },
        {
          title: window.$t('wechat_payment.api_guide.step_2'),
          imageList: []
        }
      ]
    ],
    [
      'cert',
      [
        {
          title: window.$t('wechat_payment.cert_guide.step_1'),
          imageList: ['/images/wechat-payment/cert-guide-1.png', '/images/wechat-payment/cert-guide-2.png']
        },
        {
          title: window.$t('wechat_payment.cert_guide.step_2'),
          imageList: []
        }
      ]
    ],
    [
      'serial',
      [
        {
          title: window.$t('wechat_payment.serial_guide.step_1'),
          imageList: ['/images/wechat-payment/serial-guide-1.png']
        }
      ]
    ]
  ])
  return guideListMap.get(guideMode.value)
})

const onGuideOpen = ({ mode }: { mode: string }) => {
  guideMode.value = mode
  guideVisible.value = true
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
    :title="$t('payment.type.wechat')"
    :align-center="true"
    :close-on-click-modal="false"
    width="700px"
    destroy-on-close
    append-to-body
    @close="close"
  >
    <ElForm ref="form_ref" :model="form" label-position="top">
      <ElFormItem :label="$t('payment.wechat_mch_id')" prop="mchId" :rules="generateInputRules({ message: 'payment.wechat_mch_id_placeholder' })">
        <template #label>
          <span>{{ $t('payment.wechat_mch_id') }}</span>
          <span class="text-[#9A9A9A] gap-1 ml-2 text-sm hover:opacity-80 cursor-pointer" @click="onGuideOpen({ mode: 'mch' })">
            <svg-icon class="inline translate-y-0.5" name="help" width="14" color="#999" />
            {{ $t('how_get') }}
          </span>
        </template>
        <ElInput v-model="form.mchId" size="large" clearable :placeholder="$t('form_input_placeholder')" />
      </ElFormItem>
      <ElFormItem :label="$t('payment.wechat_app_id')" prop="appId" :rules="generateInputRules({ message: 'payment.wechat_app_id_placeholder' })">
        <template #label>
          <span>{{ $t('payment.wechat_app_id') }}</span>
          <span class="text-[#9A9A9A] gap-1 ml-2 text-sm hover:opacity-80 cursor-pointer" @click="onGuideOpen({ mode: 'app' })">
            <svg-icon class="inline translate-y-0.5" name="help" width="14" color="#999" />
            {{ $t('how_get') }}
          </span>
        </template>
        <ElInput v-model="form.appId" size="large" clearable :placeholder="$t('form_input_placeholder')" />
      </ElFormItem>
      <ElFormItem
        :label="$t('payment.wechat_api_secret')"
        prop="apiV3Key"
        :rules="generateInputRules({ message: 'payment.wechat_api_secret_placeholder' })"
      >
        <template #label>
          <span>{{ $t('payment.wechat_api_secret') }}</span>
          <span class="text-[#9A9A9A] gap-1 ml-2 text-sm hover:opacity-80 cursor-pointer" @click="onGuideOpen({ mode: 'api' })">
            <svg-icon class="inline translate-y-0.5" name="help" width="14" color="#999" />
            {{ $t('how_get') }}
          </span>
        </template>
        <ElInput v-model="form.apiV3Key" size="large" clearable :placeholder="$t('form_input_placeholder')" />
      </ElFormItem>
      <ElFormItem :label="$t('payment.wechat_cert')" prop="certPath" :rules="generateInputRules({ message: 'payment.wechat_cert_placeholder' })">
        <template #label>
          <span>{{ $t('payment.wechat_cert') }}</span>
          <span class="text-[#9A9A9A] gap-1 ml-2 text-sm hover:opacity-80 cursor-pointer" @click="onGuideOpen({ mode: 'cert' })">
            <svg-icon class="inline translate-y-0.5" name="help" width="14" color="#999" />
            {{ $t('how_get') }}
          </span>
        </template>
        <UploadCertificate v-model="form.certPath" v-model:file-name="form.certName" @change="$refs.form_ref.validateField('certPath')" />
      </ElFormItem>
      <ElFormItem
        :label="$t('payment.wechat_private_key')"
        prop="privateKeyPath"
        :rules="generateInputRules({ message: 'payment.wechat_private_key_placeholder' })"
      >
        <template #label>
          <span>{{ $t('payment.wechat_private_key') }}</span>
          <span class="text-[#9A9A9A] gap-1 ml-2 text-sm hover:opacity-80 cursor-pointer" @click="onGuideOpen({ mode: 'cert' })">
            <svg-icon class="inline translate-y-0.5" name="help" width="14" color="#999" />
            {{ $t('how_get') }}
          </span>
        </template>
        <UploadCertificate
          v-model="form.privateKeyPath"
          v-model:file-name="form.privateKeyName"
          @change="$refs.form_ref.validateField('privateKeyPath')"
        />
      </ElFormItem>
      <ElFormItem
        :label="$t('payment.wechat_serial_no')"
        prop="serialNo"
        :rules="generateInputRules({ message: 'payment.wechat_serial_no_placeholder' })"
      >
        <template #label>
          <span>{{ $t('payment.wechat_serial_no') }}</span>
          <span class="text-[#9A9A9A] gap-1 ml-2 text-sm hover:opacity-80 cursor-pointer" @click="onGuideOpen({ mode: 'serial' })">
            <svg-icon class="inline translate-y-0.5" name="help" width="14" color="#999" />
            {{ $t('how_get') }}
          </span>
        </template>
        <ElInput v-model="form.serialNo" size="large" clearable :placeholder="$t('form_input_placeholder')" />
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
    <ElDialog v-model="guideVisible" :title="$t(guideTitle)" :align-center="true" width="860px" destroy-on-close append-to-body>
      <ul class="flex flex-col gap-4 pb-4 box-border max-h-[84vh] overflow-y-auto">
        <li v-for="(item, index) in guideList" :key="index" class="flex flex-col gap-2 text-[#1D1E1F] text-sm">
          <div class="text-wrap break-words whitespace-pre-wrap" v-html="item.title" />
          <div v-for="(image, imageIndex) in item.imageList" :key="imageIndex" class="w-full">
            <ElImage
              :src="$getRealPath({ url: image })"
              class="w-full"
              fit="contain"
              :preview-teleported="true"
              :preview-src-list="[$getRealPath({ url: image })]"
            />
          </div>
        </li>
      </ul>
    </ElDialog>
  </ElDialog>
</template>

<style scoped></style>
