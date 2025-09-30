<template>
  <ElDialog
    v-model="visible"
    :title="$t('payment.type.alipay')"
    :align-center="true"
    :close-on-click-modal="false"
    width="700px"
    destroy-on-close
    append-to-body
    @close="close"
  >
    <div ref="guide_ref" class="gap-3 bg-[#F6F9FC] p-5 mb-4 box-border text-sm text-[#4F5052]">
      <div
        class="whitespace-pre-wrap leading-6"
        v-html="$t('payment.alipay_guide_html', { callback_url: callback_url })"
      />
      <ElIcon
        ref="copy_ref"
        v-copy="callback_url"
        class="cursor-pointer ml-1 mt-1 text-[#4F5052] hover:text-[#3664EF]"
        :size="14"
      >
        <CopyDocument />
      </ElIcon>
    </div>

    <ElForm ref="form_ref" :model="form" label-position="top">
      <ElFormItem
        :label="$t('payment.alipay_app_id')"
        prop="appId"
        :rules="generateInputRules({ message: 'payment.alipay_app_id_placeholder' })"
      >
        <template #label>
          <span>{{ $t('payment.alipay_app_id') }}</span>
          <span
            class="text-[#9A9A9A] gap-1 ml-2 text-sm hover:opacity-80 cursor-pointer"
            @click="onGuideOpen({ mode: 'app' })"
          >
            <svg-icon class="inline translate-y-0.5" name="help" width="14" color="#999" />
            {{ $t('how_get') }}
          </span>
        </template>
        <ElInput v-model="form.appId" size="large" clearable :placeholder="$t('form_input_placeholder')" />
      </ElFormItem>
      <ElFormItem
        :label="$t('payment.alipay_mch_id')"
        prop="privateKey"
        :rules="generateInputRules({ message: 'payment.alipay_mch_id_placeholder' })"
      >
        <ElInput v-model="form.privateKey" size="large" clearable :placeholder="$t('form_input_placeholder')" />
      </ElFormItem>
      <ElFormItem
        :label="$t('payment.alipay_api_secret')"
        prop="alipayPublicKey"
        :rules="generateInputRules({ message: 'payment.alipay_api_secret_placeholder' })"
      >
        <ElInput v-model="form.alipayPublicKey" size="large" clearable :placeholder="$t('form_input_placeholder')" />
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
    <ElDialog
      v-model="guideVisible"
      :title="$t(guideTitle)"
      :align-center="true"
      width="860px"
      destroy-on-close
      append-to-body
    >
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

<script setup lang="ts">
import { computed, reactive, ref, shallowRef, nextTick } from 'vue'
import { ElIcon } from 'element-plus'
import { CopyDocument } from '@element-plus/icons-vue'

import { generateInputRules } from '@/utils/form-rule'
import { paymentApi } from '@/api/modules/payment'
import { prepareSavePaymentSettingData } from '@/api/modules/payment/transform'
import { PAYMENT_TYPE } from '@/constants/payment'
import { useUserStore } from '@/stores/modules/user'
import { api_host } from '@/utils/config'

const emits = defineEmits(['success'])

const user_store = useUserStore()

const copy_ref = shallowRef<InstanceType<typeof ElIcon> | null>(null)
const guide_ref = ref()

const form_ref = ref()
const visible = ref(false)
const form = reactive({
  appId: '',
  privateKey: '',
  alipayPublicKey: '',
})
const origin_data = ref({})
const submitting = ref(false)

const callback_url = computed(() => {
  return `${api_host}/api/payment/alipay/notify/${user_store.info.user_id}`
})

const open = ({ data = {} } = {}) => {
  const config = (data as any).pay_config || {}
  form.appId = config.appId || ''
  form.privateKey = config.privateKey || ''
  form.alipayPublicKey = config.alipayPublicKey || ''
  origin_data.value = data
  visible.value = true
  nextTick(() => {
    const copy_hook_el = guide_ref.value?.querySelector('.copy-hook')
    copy_hook_el?.appendChild(copy_ref.value?.$el)
  })
}

const reset = () => {
  form.appId = ''
  form.privateKey = ''
  form.alipayPublicKey = ''
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
  const { preparedData, pay_setting_id } = prepareSavePaymentSettingData({
    pay_setting_id: (origin_data.value as any).pay_setting_id,
    pay_config,
    extra_config: {},
    pay_type: PAYMENT_TYPE.ALIPAY,
  })

  await paymentApi.savePaymentSetting({ pay_setting_id, ...preparedData }).finally(() => {
    submitting.value = false
  })
  emits('success')
  ElMessage.success(window.$t('action_save_success'))
  close()
}

const guideVisible = ref(false)
const guideTitle = computed(() => {
  const guideTitleMap = new Map([['app', 'wechat_payment.app_guide.title']])
  return guideTitleMap.get('app')
})
const guideList = computed(() => {
  const guideListMap = new Map([
    [
      'app',
      [
        {
          title: window.$t('alipay_payment.app_guide.step_1'),
          imageList: ['/images/alipay-payment/app-guide-1.png'],
        },
      ],
    ],
  ])
  return guideListMap.get('app')
})

const onGuideOpen = ({ mode }: { mode: string }) => {
  guideVisible.value = true
}

defineExpose({
  open,
  close,
  reset,
})
</script>

<style scoped></style>
