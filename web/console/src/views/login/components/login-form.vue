<template>
  <div v-if="is_login">
    <h4 class="text-3xl text-[#1D1E1F] font-bold text-center mb-8">
      {{ $t('login.select_enterprise') }}
    </h4>
    <ul v-loading="loading" class="w-[400px] max-h-[440px] pr-1 box-border overflow-auto flex flex-col gap-4">
      <!-- <ElEmpty v-if="!enterprise_list.length" class="mt-10" :description="$t('no_data')" /> -->
      <li
        v-for="item in enterprise_list"
        :key="item.apply_id"
        class="flex items-center border rounded-sm p-4 cursor-pointer hover:border-[#3664EF]"
        @click="handleEnterpriseSelect({ data: item })"
      >
        <img class="flex-none mr-4 w-[60px] h-[60px] object-cover rounded" :src="item.logo" />
        <div class="flex-1">
          <div class="max-w-64 text-base text-[#182B50] truncate">
            {{ item.name || '- -' }}
          </div>
          <div class="text-sm text-[#9A9A9A] mt-2">
            {{ item.domain || '- -' }}
          </div>
        </div>
        <ElTag v-if="item.is_process" type="warning" effect="light">
          {{ $t('apply.process') }}
        </ElTag>
        <ElTag v-else-if="item.is_expired" type="info" effect="light">
          {{ $t('apply.expired') }}
        </ElTag>
        <ElTag v-else-if="item.is_reject" type="danger" effect="light">
          {{ $t('apply.reject') }}
        </ElTag>
        <ElIcon v-else-if="item.is_loading" size="18" color="#999" class="animate-spin">
          <Loading />
        </ElIcon>
        <ElIcon v-else size="18">
          <ArrowRight />
        </ElIcon>
      </li>
      <li
        class="w-full min-h-[88px] flex-center border border-dashed border-[#E7EEFF] bg-[#F4F7FE] text-[#3664EF] rounded-sm cursor-pointer hover:border-[#3664EF]"
        @click="onApply"
      >
        + {{ $t('apply.create_site') }}
      </li>
    </ul>
    <ElButton
      type="primary"
      text
      class="relative mt-4 !bg-transparent left-1/2 -translate-x-1/2"
      @click.stop="loadEnterpriseList"
    >
      {{ $t('apply.refresh_list') }}
    </ElButton>
    <div
      class="flex items-center text-sm text-[#5B6A91] cursor-pointer w-max mx-auto mt-2"
      @click="service_visible = true"
    >
      <SvgIcon name="service" width="14px" height="14px" class="mr-2" />
      {{ $t('apply.contact_customer_service') }}
    </div>
  </div>
  <ElForm v-else ref="form_ref" class="relative max-w-[440px] w-full" :model="form" label-position="top">
    <h4 class="text-3xl text-[#1D1E1F] font-bold text-center mb-10">
      {{ $t(`login.${form.type}_login`) }}
    </h4>
    <template v-if="form.type === 'mobile'">
      <ElFormItem
        prop="username"
        :rules="generateInputRules({ message: 'login.mobile_placeholder', validator: ['text', 'mobile'] })"
      >
        <template #label>
          <span class="text-[#1D1E1F]">{{ $t('mobile') }}</span>
        </template>
        <ElInput
          v-model="form.username"
          style="--el-input-bg-color: #f1f2f3; --el-input-border-color: transparent; --el-input-height: 44px"
          autocomplete="new-username"
          name="prevent_autofill_username"
          size="large"
          :placeholder="$t('login.mobile_placeholder')"
          clearable
          @input="checkAccountValidation"
          @blur="checkAccountValidation"
        />
      </ElFormItem>
      <ElFormItem
        class="relative"
        prop="verify_code"
        :rules="generateInputRules({ message: 'verification_code_placeholder' })"
      >
        <template #label>
          <span class="text-[#1D1E1F]">{{ $t('verification_code') }}</span>
        </template>
        <VerificationCodeInput
          ref="verify_code_input_ref"
          v-model="form.verify_code"
          :account="form.username"
          :disabled="!isAccountValid"
        />
      </ElFormItem>
      <ElButton
        type="primary"
        round
        class="w-full mt-8 !h-10"
        :disabled="!form.username || !form.verify_code"
        :loading="submitting"
        @click="onLogin"
      >
        {{ $t('action_login') }}
      </ElButton>
    </template>
    <template v-if="form.type === 'password'">
      <ElFormItem
        prop="username"
        :rules="generateInputRules({ message: 'login.account_placeholder', validator: ['text'] })"
      >
        <template #label>
          <span class="text-[#1D1E1F]">{{ $t('account') }}</span>
        </template>
        <ElInput
          v-model="form.username"
          style="--el-input-bg-color: #f1f2f3; --el-input-border-color: transparent; --el-input-height: 44px"
          size="large"
          :placeholder="$t('login.account_placeholder')"
          clearable
        />
      </ElFormItem>
      <ElFormItem
        class="relative"
        prop="password"
        :rules="[
          ...generateInputRules({ message: 'login.password_placeholder', validator: ['password'] }),
          {
            min: 8,
            max: 20,
            message: $t('login.password_length_v2'),
            trigger: 'blur',
          },
        ]"
      >
        <template #label>
          <span class="text-[#1D1E1F]">{{ $t('password') }}</span>
        </template>
        <ElInput
          v-model="form.password"
          style="--el-input-bg-color: #f1f2f3; --el-input-border-color: transparent; --el-input-height: 44px"
          size="large"
          type="password"
          show-password
          clearable
          :placeholder="$t('login.password_placeholder')"
          @keyup.enter="onLogin"
        />
        <div class="absolute right-0 -bottom-9">
          <ElButton type="text" class="bg-transparent text-sm" @click="onRegister">
            {{ $t('user_register') }}
          </ElButton>
          <ElDivider direction="vertical" />
          <ElButton type="text" class="bg-transparent text-sm" @click="onForgetPassword">
            {{ $t('login.forget_password') }}
          </ElButton>
        </div>
      </ElFormItem>
      <ElButton
        type="primary"
        round
        class="w-full mt-8 !h-10"
        :disabled="!form.username || !form.password"
        :loading="submitting"
        @click="onLogin"
      >
        {{ $t('action_login') }}
      </ElButton>
    </template>
    <template v-if="form.type === 'wechat'">
      <WeChat @oauth-success="handleOauthSuccess" />
    </template>
    <template v-if="form.type === 'bind_mobile'">
      <ElFormItem
        prop="username"
        class="mt-5"
        label-position="top"
        :rules="generateInputRules({ message: 'login.mobile_placeholder', validator: ['text', 'mobile'] })"
      >
        <template #label>
          <span class="text-[#1D1E1F]">{{ $t('mobile') }}</span>
        </template>
        <ElInput
          v-model="form.username"
          style="--el-input-bg-color: #f1f2f3; --el-input-border-color: transparent; --el-input-height: 44px"
          autocomplete="new-username"
          name="prevent_autofill_username"
          size="large"
          :placeholder="$t('login.mobile_placeholder')"
          clearable
        />
      </ElFormItem>
      <ElFormItem
        class="relative"
        prop="verify_code"
        :rules="generateInputRules({ message: 'verification_code_placeholder' })"
      >
        <VerificationCodeInput ref="verify_code_input_ref" v-model="form.verify_code" :account="form.username" />
      </ElFormItem>
      <ElButton
        type="primary"
        round
        class="w-full mt-8 !h-10"
        :disabled="!form.username || !form.verify_code"
        :loading="submitting"
        @click="onLogin"
      >
        {{ $t('action_confirm') }}
      </ElButton>
    </template>
    <ElDivider class="!w-[80%] !mx-auto">
      <span class="text-[#9A9A9A]">{{ $t('or') }}</span>
    </ElDivider>
    <div class="flex flex-col gap-4">
      <ElButton
        v-if="form.type !== 'password'"
        type="default"
        round
        class="w-full !h-10 !ml-0 font-medium"
        @click="onPasswordLogin"
      >
        <SvgIcon class="mr-2" name="password" width="28" />
        {{ $t('login.password_login_v2') }}
      </ElButton>
      <ElButton
        v-if="form.type !== 'mobile'"
        type="default"
        round
        class="w-full !h-10 !ml-0 font-medium"
        @click="onMobileLogin"
      >
        <SvgIcon class="mr-2" name="mobile" width="28" />
        {{ $t('login.mobile_login_v2') }}
      </ElButton>
      <ElButton type="default" round class="w-full !h-10 !ml-0 font-medium" @click="onWechatLogin">
        <SvgIcon class="mr-2" name="wechat" width="28" />
        {{ $t('login.wechat_login_v2') }}
      </ElButton>
      <ElButton type="default" round class="w-full !h-10 !ml-0 font-medium" @click="onGoogleLogin">
        <SvgIcon class="mr-2" name="google" width="28" />
        {{ $t('login.google_login_v2') }}
      </ElButton>
      <div
        class="text-xs text-[#9A9A9A] text-center"
        v-html="
          $t('login.agree_and_policy', {
            agree: `<span class=\'agree-hook cursor-pointer text-[#4F5052] text-xs mx-1 -mt-0.5 underline underline-offset-4\'>${$t('login.agree')}</span>`,
            policy: `<span class=\'policy-hook cursor-pointer text-[#4F5052] text-xs mx-1 -mt-0.5 underline underline-offset-4\'>${$t('login.policy')}</span>`,
          })
        "
      />
    </div>
  </ElForm>
  <ServiceDialog v-model:visible="service_visible" :title="$t('apply.contact_customer_service_v2')" />
</template>

<script setup lang="ts">
import { ArrowRight, Loading } from '@element-plus/icons-vue'

import { computed, onMounted, onUnmounted, reactive, ref } from 'vue'
import ServiceDialog from '@/components/ServiceDialog/index.vue'
import WeChat from './wechat.vue'
import { useEnterpriseStore, useUserStore } from '@/stores'
import { generateInputRules } from '@/utils/form-rule'
import { validateFormField } from '@/utils/form-validator'
import eventBus from '@/utils/event-bus'
import { sleep } from '@/utils'
import systemLogApi from '@/api/modules/system-log'
import { SYSTEM_LOG_ACTION } from '@/constants/system-log'

const emits = defineEmits<{
  (e: 'forget'): void
  (e: 'register'): void
  (e: 'apply'): void
  (e: 'login-success'): void
}>()

const user_store = useUserStore()
const enterprise_store = useEnterpriseStore()
const form_ref = ref()
const verify_code_input_ref = ref()
const form = reactive({
  type: 'password',
  username: '',
  password: '',
  verify_code: '',
})
const submitting = ref(false)
const enterprise_list = ref([])
const loading = ref(false)
const service_visible = ref(false)
const oauth_data = ref<any>({})
const isAccountValid = ref(false)

const is_login = computed(() => {
  const { access_token } = user_store.info
  return Boolean(access_token)
})

const defaultDomain = computed(() => `${window.location.origin}/#/index`)

// 检查账号验证状态
const checkAccountValidation = async () => {
  isAccountValid.value = await validateFormField(form_ref, 'username')
}

onMounted(() => {
  if (is_login.value) return loadEnterpriseList()
  const agree_hook_el = form_ref.value.$el.querySelector('.agree-hook')
  agree_hook_el.onclick = onAgree
  const policy_hook_el = form_ref.value.$el.querySelector('.policy-hook')
  policy_hook_el.onclick = onPolicy
  eventBus.on('language-change', onLanguageChange)
})
onUnmounted(() => {
  eventBus.off('language-change', onLanguageChange)
})

const onLanguageChange = () => {
  if (form_ref.value) form_ref.value.clearValidate()
  // form_ref.value.validate()
}

const loadEnterpriseList = async () => {
  loading.value = true
  const { list = [] } = await enterprise_store.loadListData({ data: { status: -1 } }).finally(() => {
    loading.value = false
  })
  enterprise_list.value = list
}
const handleEnterpriseSelect = async ({ data = {} } = {}) => {
  if (data.is_process) return ElMessage.warning(window.$t('apply.process'))
  if (data.is_reject) return ElMessage.warning(data.reject_reason || window.$t('apply.reject'))
  if (data.is_expired) return (service_visible.value = true)
  if (data.is_loading) return

  data.is_loading = true
  if (!data.eid) {
    // 考虑申请流程存在异步情况，需要多次请求
    const { apply_id } = data
    let request_count = 0
    const refreshData = async () => {
      const { list = [] } = await enterprise_store.loadListData({ data: { status: -1 } })
      request_count++
      const apply_data = list.find(item => item.apply_id == apply_id)
      if (!apply_data.eid && request_count < 5) {
        await sleep(1)
        await refreshData()
      }
      data.eid = apply_data.eid || 0
      return apply_data.eid
    }
    await refreshData()
    data.is_loading = false
    if (!data.eid) return ElMessage.warning('Invalid eid')
  }

  await enterprise_store.loadDetailData({ data: { eid: data.eid } })
  await systemLogApi.create({
    action: SYSTEM_LOG_ACTION.LOGIN,
    content: '登录',
  })

  // 重置加载状态
  data.is_loading = false

  // 构建目标URL
  // #ifdef KM
  if (window.parent) {
    window.parent.postMessage(
      {
        action: 'saas-login-success',
        eid: data.eid,
        access_token: user_store.info.access_token,
      },
      '*'
    )
  } else {
    window.location.reload()
  }
  // #endif

  // #ifndef KM
  const targetUrl = `${data.domain}?access_token=${user_store.info.access_token}&eid=${data.eid}`
  if (window.parent) {
    window.parent.postMessage({
      action: 'saas-login-redirect',
      url: targetUrl,
    })
  }
  // #endif
}
const onApply = () => {
  emits('apply')
}
const reset = () => {
  form.type = 'password'
  form.username = ''
  form.password = ''
  form.verify_code = ''
}
const onLogin = async () => {
  const valid = await form_ref.value.validate()
  if (!valid) return
  submitting.value = true
  if (form.type.includes('mobile')) {
    const valid = await verify_code_input_ref.value.validateCode()
    if (!valid) return (submitting.value = false)
  }
  if (form.type === 'bind_mobile') {
    await user_store.bind_wechat({
      mobile: form.username,
      verify_code: form.verify_code,
      openid: oauth_data.value.openid,
      unionid: oauth_data.value.unionid,
      nickname: oauth_data.value.nickname,
      from: 'saas',
    })
  } else {
    const data = await user_store
      .login({ type: form.type, data: form, hideError: true })
      .catch(err => {
        // if (err.code == RESPONSE_CODE_UNAUTHORIZED_ERROR && err.origin_message == 'unauthorized: user not found') {
        // 	onRegister()
        // } else ElMessage.warning(window.$t(err.message))
        ElMessage.warning(
          window.$t(
            err.origin_message === 'unauthorized'
              ? 'response_message.user_not_found'
              : 'response_message.username_or_password_is_incorrect'
          )
        )
        return Promise.reject(err)
      })
      .finally(() => {
        submitting.value = false
      })
  }
  loadEnterpriseList()
  ElMessage.success(window.$t('action_login_success'))
  reset()
}
// const onRegister = () => {
// 	form_ref.value.validate(async (valid) => {
// 		if (!valid) return
// 		await user_store.register({ data: form })
// 		// ElMessage.success(window.$t('action_register_success'))
// 		ElMessage.success(window.$t('action_login_success'))
// 		reset()
// 	})
// }

const handleOauthSuccess = async (data: any) => {
  await user_store.wechat_login({ unionid: data.unionid, from: 'saas' }).catch(err => {
    oauth_data.value = data
    form.type = 'bind_mobile'
    return Promise.reject(err)
  })
  loadEnterpriseList()
  ElMessage.success(window.$t('action_login_success'))
  reset()
}

const onForgetPassword = () => {
  emits('forget')
}
const onRegister = () => {
  emits('register')
}
const onPasswordLogin = () => {
  form.username = ''
  form.password = ''
  form.verify_code = ''
  form.type = 'password'
}
const onMobileLogin = () => {
  form.username = ''
  form.password = ''
  form.verify_code = ''
  form.type = 'mobile'
}
const onWechatLogin = () => {
  form.type = 'wechat'
}
const onGoogleLogin = () => {
  ElMessage.warning(window.$t('feature_coming_soon'))
}
const onAgree = () => {
  ElMessage.warning(window.$t('feature_coming_soon'))
}
const onPolicy = () => {
  ElMessage.warning(window.$t('feature_coming_soon'))
}

defineExpose({
  reset,
})
</script>

<style>
.dialog .el-dialog__header {
  border: none !important;
}
</style>
