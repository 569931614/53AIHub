<template>
  <!-- 登录弹窗 -->
  <el-dialog
    v-model="isVisible"
    class="login-dialog"
    destroy-on-close
    :style="{
      backgroundImage: `url(${$getPublicPath('/images/login_bg.png')})`
    }"
    :close-on-click-modal="false"
    :show-close="login_way !== 'bind_mobile'"
    center
  >
    <!-- 顶栏显示 -->
    <div v-if="login_way === 'bind_mobile'" class="flex justify-center mt-5">
      <h4 class="text-xl text-primary font-semibold w-full flex items-center">
        <ElIcon class="mr-1 cursor-pointer" size="18" color="#4E4F51" @click="handleLOGIN_WAY(LOGIN_WAY.wechat_login)">
          <ArrowLeft />
        </ElIcon>
        {{ $t('login.bind_mobile') }}
      </h4>
    </div>
    <div v-else class="flex justify-center mt-5">
      <h4 class="text-xl text-[#1D1E1F] font-bold text-center">
        {{ $t('login.' + login_way + '_title') }}
      </h4>
    </div>

    <template v-if="login_way === LOGIN_WAY.wechat_login || login_way === LOGIN_WAY.wecom_login">
      <WechatView v-if="login_way === LOGIN_WAY.wechat_login" height="292px" @oauth-success="handleOauthSuccess" />
      <WecomView v-else height="292px" />

      <div class="text-xs text-[#9A9A9A] text-center">
        {{ $t('login.agree') }}
        <a class="text-[#4F5052] cursor-pointer underline">{{ $t('login.terms_of_service') }}</a>
        {{ $t('action.and') }}
        <a class="text-[#4F5052] cursor-pointer underline">{{ $t('login.privacy_policy') }}</a>
      </div>
    </template>
    <template v-else>
      <!-- 密码登陆的form表单 -->
      <el-form
        v-if="login_way === LOGIN_WAY.password_login"
        ref="formRef"
        label-position="top"
        :model="form"
        :rules="rules"
        class="mt-7"
        @keyup.enter="handleSubmit"
      >
        <el-form-item :label="$t('form.account')" prop="username">
          <el-input
            v-model="form.username"
            v-trim
            size="large"
            class="el-input--main"
            :placeholder="$t('form.input_placeholder') + $t('form.email')"
            clearable
          />
        </el-form-item>

        <!-- 密码的输入框 -->
        <el-form-item :label="$t('form.password')" prop="password">
          <el-input
            v-model="form.password"
            v-trim
            show-password
            size="large"
            class="el-input--main"
            :placeholder="$t('form.input_placeholder') + $t('form.password')"
          ></el-input>
        </el-form-item>
      </el-form>

      <!-- 短信登陆的form表单 -->
      <el-form
        v-else-if="login_way === LOGIN_WAY.message_login || login_way === LOGIN_WAY.bind_mobile"
        ref="formRef"
        label-position="top"
        :model="form"
        :rules="rules"
        class="mt-7"
        @keyup.enter="handleSubmit"
      >
        <el-form-item :label="$t('form.mobile')" prop="username">
          <el-input
            v-model="form.username"
            v-trim
            size="large"
            class="el-input--main"
            :placeholder="$t('form.input_placeholder') + $t('form.account_alias')"
            clearable
          />
        </el-form-item>

        <el-form-item :label="login_way === LOGIN_WAY.bind_mobile ? '' : $t('form.verify_code')" prop="verify_code">
          <div class="flex items-center" style="width: 100%">
            <el-input
              v-model="form.verify_code"
              v-trim
              size="large"
              class="el-input--main no-right-radius w-80"
              :placeholder="$t('form.input_placeholder') + $t('form.verify_code')"
            ></el-input>
            <el-button
              v-debounce
              :disabled="isSending || !isMobile"
              class="!bg-[#f5f5f5] border-0 px-3 h-[44px] w-29 no-left-radius"
              @click.stop="handleGetCode"
            >
              <div :class="['text-sm', 'pl-5', 'border-l', 'pr-1', 'text-[#2563EB]', { 'text-[#9A9A9A]': isSending || !isMobile }]">
                {{ codeCount ? `${codeCount}s` : $t('form.get_verify_code') }}
              </div>
            </el-button>
          </div>
        </el-form-item>
      </el-form>

      <div v-if="login_way !== 'bind_mobile'" class="flex items-center justify-between mt-3 max-md:flex-col max-md:gap-2">
        <!-- 底部协议 -->
        <div v-if="['password_login', 'message_login', 'wechat_login'].includes(login_way)" class="text-xs text-[#9A9A9A] flex">
          {{ $t('login.agree') }}
          <a class="text-[#4F5052] cursor-pointer underline">{{ $t('login.terms_of_service') }}</a>
          {{ $t('action.and') }}
          <a class="text-[#4F5052] cursor-pointer underline">{{ $t('login.privacy_policy') }}</a>
        </div>
        <div class="flex items-center justify-end">
          <el-button class="mr-1" link type="primary" @click="handleRegister">
            {{ $t('action.user_register') }}
          </el-button>
          <template v-if="!isOpLocalEnv || (isOpLocalEnv && openSMTP)">
            <div class="border-l border-[#E6E8EB] mr-1 h-4"></div>
            <el-button link type="primary" @click="handleForgetPassword">
              {{ $t('action.forget_password') }}
            </el-button>
          </template>
        </div>
      </div>

      <template v-if="login_way === 'bind_mobile'">
        <div class="flex items-center justify-end">
          <ElButton v-debounce type="primary" :disabled="!form.verify_code && !form.username" size="large" class="min-w-[96px]" @click="handleSubmit">
            {{ $t('action.ok') }}
          </ElButton>
        </div>
      </template>
      <el-button v-else v-debounce type="primary" size="large" round class="w-full mt-5" @click="handleSubmit">
        {{ $t('action.login') }}
      </el-button>
    </template>

    <!-- 其他登录方式 -->
    <template v-if="!isOpLocalEnv && login_way !== 'bind_mobile'">
      <ElDivider class="my-8">
        <span class="text-placeholder text-sm">
          {{ $t('login.other_login_way') }}
        </span>
      </ElDivider>
      <div class="flex items-center justify-center mt-5">
        <template v-for="item in LOGIN_WAYList" :key="item.value">
          <div
            class="flex-1 flex flex-col items-center justify-center gap-3 cursor-pointer"
            :class="item.value === login_way ? 'text-theme' : 'text-regular'"
            @click="handleLOGIN_WAY(item.value)"
          >
            <div class="size-6">
              <svg-icon :name="item.icon" size="24" :stroke="item.value !== 'wecom_login'" />
            </div>
            <p class="text-sm">
              {{ $t(`login.${item.label}`) }}
            </p>
          </div>
        </template>
      </div>
    </template>
  </el-dialog>

  <!-- 注册弹窗 -->
  <el-dialog
    v-model="registerVisible"
    class="login-dialog"
    destroy-on-close
    :style="{
      backgroundImage: `url(${$getPublicPath('/images/login_bg.png')})`
    }"
    :close-on-click-modal="false"
    center
  >
    <Register ref="registerRef" :open-s-m-t-p="openSMTP" @success="handleClose" @close="handleClose"></Register>
  </el-dialog>

  <!-- 重置密码弹窗 -->
  <el-dialog
    v-model="forgetPasswordVisible"
    class="login-dialog"
    destroy-on-close
    :style="{
      backgroundImage: `url(${$getPublicPath('/images/login_bg.png')})`
    }"
    :close-on-click-modal="false"
    center
  >
    <ElButton class="absolute top-8 left-8 !text-[#B9BEC2]" type="info" link @click="handleClose">
      <ElIcon class="mr-1">
        <ArrowLeft />
      </ElIcon>
      {{ $t('action.back') }}
    </ElButton>

    <div class="pb-2">
      <h4 class="text-xl text-[#1D1E1F] font-bold text-center pb-8">
        {{ $t('form.reset_password') }}
      </h4>
      <ForgetPassword ref="forgetPasswordRef" @success="handleClose" @close="handleClosePaw" />
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ArrowLeft } from '@element-plus/icons-vue'
import { ref, reactive, watch, computed, onMounted } from 'vue'
import type { FormInstance } from 'element-plus'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/modules/user'
import { useAgentStore } from '@/stores/modules/agent'
import { useEnterpriseStore } from '@/stores/modules/enterprise'
import { getMobileRules, getAccountOrEmailRules, getEmailRules, getPasswordRules } from '@/utils/form-rules'
import SvgIcon from '@/components/SvgIcon.vue'

import useEnv from '@/hooks/useEnv'
import useMobile from '@/hooks/useMobile'

import commonApi from '@/api/modules/common'
import enterpriseApi from '@/api/modules/enterprise'
import Register from './register.vue'
import ForgetPassword from './forgetPassword.vue'
import WechatView from './wechat.vue'
import WecomView from './wecom.vue'

const emit = defineEmits(['login', 'close'])

const { isOpLocalEnv } = useEnv()
const registerVisible = ref(false)
const forgetPasswordVisible = ref(false)

const userStore = useUserStore()
const agentStore = useAgentStore()
const enterpriseStore = useEnterpriseStore()

const { sendcode, codeRule, codeCount } = useMobile()

const formRef = ref<FormInstance>()

const openSMTP = ref(false)

const isVisible = ref(false)

const isSending = ref(true)

const isMobile = computed(() => {
  return /^1[3-9]\d{9}$/.test(form.username)
})
const LOGIN_WAY = {
  password_login: 'password_login',
  message_login: 'message_login',
  wechat_login: 'wechat_login',
  wecom_login: 'wecom_login',
  bind_mobile: 'bind_mobile'
} as const

type LoginWay = (typeof LOGIN_WAY)[keyof typeof LOGIN_WAY]

const login_way = ref<LoginWay>(LOGIN_WAY.password_login as LoginWay)

const form = reactive({
  username: '',
  password: '',
  verify_code: ''
})

const rules = computed(() => {
  return {
    username: [isOpLocalEnv.value ? getEmailRules() : login_way.value === LOGIN_WAY.password_login ? getAccountOrEmailRules() : getMobileRules()],
    password: [getPasswordRules()],
    verify_code: [codeRule]
  }
})

const LOGIN_WAYList = computed(() => {
  return [
    {
      icon: 'wechat',
      label: 'wechat_login',
      value: LOGIN_WAY.wechat_login
    },
    {
      icon: 'wecom',
      label: 'wecom_login',
      value: LOGIN_WAY.wecom_login
    },
    {
      icon: 'safe',
      label: 'password_login',
      value: LOGIN_WAY.password_login
    },
    {
      icon: 'iphone',
      label: 'message_login',
      value: LOGIN_WAY.message_login
    }
  ]
    .filter((item) => {
      return isOpLocalEnv.value ? item.value === LOGIN_WAY.password_login : true
    })
    .filter((item) => {
      return item.value === LOGIN_WAY.wecom_login ? enterpriseStore.is_install_wecom : true
    })
})

const handleGetCode = () => {
  if (!isMobile.value) return
  sendcode(form.username)
}

const accountLogin = () => {
  return userStore.login({
    username: form.username,
    password: form.password
  })
}

// const smsLogin = () => {
//   return commonApi.verifycode({
//     mobile: form.username,
//     verifycode: form.verify_code,
//     type: '1'
//   })
// }

const handleSubmit = () => {
  return formRef.value?.validate().then(async (valid) => {
    if (!valid) return
    try {
      if (['message_login', 'bind_mobile'].includes(login_way.value)) {
        await commonApi.verifycode({
          mobile: form.username,
          verifycode: form.verify_code,
          type: '1'
        })
        if (login_way.value === 'bind_mobile') {
          await userStore.bind_wechat({
            mobile: form.username,
            verify_code: form.verify_code,
            openid: oauth_data.value.openid,
            unionid: oauth_data.value.unionid,
            nickname: oauth_data.value.nickname
          })
        } else {
          await userStore.sms_login({
            mobile: form.username,
            verify_code: form.verify_code
          })
        }
      } else {
        await accountLogin()
      }
      ElMessage.success(window.$t('status.login_success'))
      agentStore.loadAgentList()
      close()
    } catch (error) {
      const response = error.response || {}
      const data = response.data || {}
      const message = data.message || ''
      if (message.includes('record not found')) {
        // 本地版自动注册
        if (isOpLocalEnv.value && !openSMTP.value) {
          await userStore.register({
            username: form.username,
            password: form.password
          })
          ElMessage.success(window.$t('status.login_success'))
          agentStore.loadAgentList()
          isVisible.value = false
        } else {
          ElMessage.warning(window.$t('status.not_found_account'))
        }
      }
      console.log(error)
    }
  })
}

// 添加重置函数
const resetForm = () => {
  form.username = ''
  form.verify_code = ''
  form.password = ''
  clearFormValidation()
}

const oauth_data = ref<any>({})

const handleLOGIN_WAY = (value: LoginWay) => {
  login_way.value = value
  resetForm()
}

const handleClose = () => {
  isVisible.value = true
  registerVisible.value = false
  forgetPasswordVisible.value = false
  form.verify_code = ''
  login_way.value = LOGIN_WAY.password_login
}

const handleClosePaw = () => {
  forgetPasswordVisible.value = false
  registerVisible.value = true
}

const handleRegister = () => {
  isVisible.value = false
  registerVisible.value = true
  resetForm()
}

const handleForgetPassword = () => {
  isVisible.value = false
  forgetPasswordVisible.value = true
  resetForm()
}

// 添加清除表单验证的方法
const clearFormValidation = () => {
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

const open = (data: { way?: LoginWay; openid?: string; unionid?: string } = {}) => {
  // login_way.value = data.way || LOGIN_WAY.password_login
  if (data.way === LOGIN_WAY.wechat_login && data.openid) {
    handleOauthSuccess({
      openid: data.openid,
      unionid: data.unionid
    })
  }
  oauth_data.value = data
  isVisible.value = true
  clearFormValidation()
}

const close = () => {
  isVisible.value = false
  clearFormValidation()
  emit('close')
}

const handleOauthSuccess = async (data: any) => {
  await userStore.wechat_login({ unionid: data.unionid }).catch((err) => {
    oauth_data.value = data
    login_way.value = LOGIN_WAY.bind_mobile
    return Promise.reject(err)
  })
  ElMessage.success(window.$t('status.login_success'))
  agentStore.loadAgentList()
  close()
}

const loadSMTP = async () => {
  const { data } = await enterpriseApi.getSMTPInfo('smtp')
  openSMTP.value = data
}

watch(
  () => codeCount.value,
  (newVal) => {
    isSending.value = newVal > 0
  },
  {
    immediate: true
  }
)

defineExpose({
  open,
  close
})

onMounted(() => {
  loadSMTP()
})
</script>

<style>
.login-dialog {
  --el-dialog-border-radius: 16px;
  --el-dialog-padding-primary: 32px;
  --el-dialog-width: 500px;

  background: linear-gradient(180deg, #eaf3ff 0%, #fff 20%) !important;
  background-size: cover;
}

.login-dialog .el-dialog__header {
  padding: 0;
}

.login-dialog .el-dialog__headerbtn {
  --el-color-info: #939499;

  top: 18px;
  right: 16px;
  font-size: 18px;
}

.no-right-radius .el-input__wrapper {
  border-top-right-radius: 0 !important;
  border-bottom-right-radius: 0 !important;
}

.no-left-radius {
  border-top-left-radius: 0 !important;
  border-bottom-left-radius: 0 !important;
}
</style>
