<template>
  <div class="flex items-center justify-between mt-6">
    <h2 class="text-primary font-semibold">{{ $t('profile.info') }}</h2>
    <div v-if="!userStore.info.is_internal" class="flex items-center gap-1">
      <div class="h-6 flex-center gap-1 px-2 text-sm text-placeholder whitespace-nowrap" :title="userStore.info.group_name">
        <img
          :src="
            !/\.png$/.test(userStore.info.group_icon)
              ? $getPublicPath(`/images/subscription/${userStore.info.group_icon}.png`)
              : userStore.info.group_icon
          "
          class="w-4 h-4 object-cover"
        />
        <p class="max-w-[5em] truncate">{{ userStore.info.group_name }}</p>
      </div>
      <div class="flex-none w-px h-3 bg-[#E6E8EB] mx-1"></div>
      <div class="text-sm text-secondary flex items-center gap-1">
        <span>{{ $t('subscription.expire_time') }}：</span>
        <span class="text-primary">{{ userStore.info.group_expire_time || '永久有效' }}</span>
      </div>
    </div>
  </div>

  <div class="flex items-center gap-2 w-full p-4 box-border cursor-pointer mt-4 rounded-lg overflow-hidden bg-[#F7F7F7]" @click="handleEdit">
    <ElAvatar class="flex-none" :size="48" :src="userStore.info.avatar" style="--el-avatar-bg-color: transparent" />
    <div class="flex-1 w-0">
      <div class="text-primary font-semibold">
        {{ userStore.info.nickname }}
      </div>
      <p class="text-sm text-secondary mt-1 flex items-center gap-1">
        {{ userStore.info.mobile || '- -' }}
        <!-- <svg-icon name="copy" v-copy="userStore.info.mobile"></svg-icon> -->
      </p>
    </div>
    <ElIcon class="text-regular" size="16">
      <ArrowRight />
    </ElIcon>
  </div>
  <div class="flex items-center justify-center">
    <h2 class="text-primary font-semibold mt-6">{{ $t('profile.bind_accounts') }}</h2>
  </div>
  <div class="flex flex-col mt-4 bg-[#F7F7F7] rounded-lg box-border overflow-hidden">
    <div class="p-4 box-border flex items-center gap-2.5 account-item">
      <img class="size-6" :src="$getPublicPath('/images/profile/mobile.png')" alt="" />
      <div class="flex-none text-base text-primary">
        {{ $t('profile.bind_mobile') }}
      </div>
      <div class="flex-1 w-0 text-sm text-placeholder invisible md:visible">
        {{ userStore.info.mobile || $t('profile.unbind_account') }}
      </div>
      <el-button v-if="!isOpLocalEnv" type="info" class="border-none !text-[#586D9A]" link plain @click="change_Mobile_Visible = true">
        {{ userStore.info.mobile ? $t('profile.change') : $t('profile.bind') }}
      </el-button>
    </div>
    <div class="p-4 box-border flex items-center gap-2.5 account-item">
      <img class="size-6" :src="$getPublicPath('/images/profile/password.png')" alt="" />
      <div class="flex-1 w-0 text-base text-primary">{{ $t('profile.login_password') }}</div>
      <el-button
        v-if="!isOpLocalEnv || (isOpLocalEnv && openSMTP)"
        type="info"
        class="!border-none !text-[#586D9A] !bg-transparent"
        link
        plain
        @click="reset_Password_Visible = true"
      >
        {{ $t('form.change') }}
      </el-button>
    </div>
    <div class="p-4 box-border flex items-center gap-2.5 account-item">
      <img class="size-6" :src="$getPublicPath('/images/profile/email.png')" alt="" />
      <div class="flex-none text-base text-primary">{{ $t('profile.bind_email') }}</div>
      <div class="flex-1 w-0 text-sm text-placeholder invisible md:visible">
        {{ userStore.info.email || $t('profile.unbind_account') }}
      </div>
      <el-button
        v-if="!isOpLocalEnv || (isOpLocalEnv && openSMTP)"
        type="info"
        class="border-none !text-[#586D9A]"
        link
        plain
        @click="emailVisible = true"
      >
        {{ userStore.info.email ? $t('profile.change') : $t('profile.bind') }}
      </el-button>
    </div>
    <div class="p-4 box-border flex items-center gap-2.5 account-item">
      <img class="size-6" :src="$getPublicPath('/images/profile/wechat.png')" alt="" />
      <div class="flex-none text-base text-primary">
        {{ $t('profile.bind_wechat') }}
      </div>
      <div class="flex-1 w-0 text-sm text-placeholder invisible md:visible">
        {{ userStore.info.openid || $t('profile.unbind_account') }}
      </div>
      <template v-if="userStore.info.openid">
        <ElButton type="info" class="border-none !text-[#586D9A]" link plain @click="handleUnbindWechat">
          {{ $t('profile.unbind') }}
        </ElButton>
        <ElButton type="info" class="border-none !text-[#586D9A] !ml-0" link plain @click="handleBindWechat">
          {{ $t('profile.change') }}
        </ElButton>
      </template>
      <template v-else>
        <ElButton :disabled="isOpLocalEnv" type="info" class="border-none !text-[#586D9A]" link plain @click="handleBindWechat">
          {{ $t('profile.bind') }}
        </ElButton>
      </template>
    </div>
    <div class="p-4 box-border flex items-center gap-2.5 account-item">
      <img class="size-6" :src="$getPublicPath('/images/profile/google.png')" alt="" />
      <div class="flex-none text-base text-primary">
        {{ $t('profile.bind_google') }}
      </div>
      <div class="flex-1 w-0 text-sm text-placeholder invisible md:visible">
        {{ userStore.info.wechatid || $t('profile.unbind_account') }}
      </div>
      <el-button type="info" class="border-none text-placeholder-important" link plain disabled>{{ $t('profile.bind') }}</el-button>
    </div>
  </div>
  <h2 class="text-primary font-semibold mt-6">{{ $t('profile.subscription_info') }}</h2>
  <router-link :to="{ name: 'Order' }" class="flex flex-col mt-4 bg-[#F7F7F7] rounded-lg box-border overflow-hidden">
    <div class="p-4 box-border flex items-center gap-2.5 account-item">
      <img class="size-6" :src="$getPublicPath('/images/profile/order.png')" alt="" />
      <div class="flex-none text-base text-primary">
        {{ $t('profile.order_info') }}
      </div>
      <div class="flex-1 w-0 text-sm text-placeholder invisible md:visible"></div>
      <ElIcon class="text-regular" size="16">
        <ArrowRight />
      </ElIcon>
    </div>
  </router-link>
  <h2 class="text-primary font-semibold mt-6">{{ $t('profile.logined_devices') }}</h2>
  <div class="flex flex-col gap-4 mt-4">
    <div class="h-[110px] px-6 rounded-lg flex items-center bg-[#F8F8F9]">
      <div class="flex-1">
        <div class="text-base text-primary flex items-center gap-2">
          <svg-icon name="network"></svg-icon>
          WindowsChrome 网页版
        </div>
        <div class="text-sm text-secondary ml-6 mt-1">2025年3月22日上午10点51分</div>
      </div>
      <div class="text-sm text-regular">{{ $t('profile.used_devices') }}</div>
    </div>
  </div>

  <el-dialog v-model="profileVisible" class="update-user-dialog" :title="$t('action.edit') + $t('profile.info')" width="500">
    <div class="">
      <div class="w-[73px] text-[#182B50] text-opacity-80">{{ $t('form.avatar') }}</div>
      <div class="mt-3 relative rounded-full w-[60px] h-[60px] overflow-hidden cursor-pointer" @click="handleChangeAvatar">
        <el-avatar style="--el-avatar-size: 60px" :src="profileForm.avatar" />
        <div
          class="absolute bottom-0 left-0 right-0 h-5 flex justify-center items-center bg-[#000000] bg-opacity-60 text-white text-opacity-80 text-xs"
        >
          {{ $t('form.change') }}
        </div>
      </div>
    </div>
    <div class="mt-6">
      <div class="w-[73px] text-[#182B50] text-opacity-80">{{ $t('form.nickname') }}</div>
      <el-input v-model="profileForm.nickname" v-trim size="large" class="mt-3 flex-1" maxlength="15" show-word-limit />
    </div>
    <template #footer>
      <el-button size="large" @click="profileVisible = false">
        {{ $t('action.cancel') }}
      </el-button>
      <el-button v-debounce size="large" type="primary" :disabled="!profileForm.nickname.trim()" @click="handleSaveProfile">
        {{ $t('action.confirm') }}
      </el-button>
    </template>
    <CropperDialog ref="corpperRef" :upload-avatar="true" :limit-size="2" @confirm="handleSuccessCropper" />
  </el-dialog>

  <!-- 绑定/更换手机号弹窗 -->
  <el-dialog
    v-model="change_Mobile_Visible"
    class="update-user-dialog"
    :title="userStore.info.mobile ? $t('profile.change') + $t('form.mobile') : $t('profile.bind') + $t('form.mobile')"
    width="500"
    @close="handleClose"
  >
    <div v-if="userStore.info.mobile" class="mb-4">
      {{ $t('form.verify_old_mobile') + userStore.info.mobile }}
    </div>

    <changeMobile ref="mobileRef" @success="onChangeMobileSuccess" @close="change_Mobile_Visible = false"></changeMobile>
  </el-dialog>

  <!-- 修改密码弹窗 -->
  <el-dialog v-model="reset_Password_Visible" class="update-user-dialog" :title="$t('profile.change_password')" width="500" @close="handleClose">
    <ResetPassword ref="passwordRef" @success="reset_Password_Visible = false"></ResetPassword>
  </el-dialog>

  <!-- 绑定/更换邮箱弹窗 -->
  <el-dialog
    v-model="emailVisible"
    class="update-user-dialog"
    :title="userStore.info.email ? $t('profile.change') + $t('form.email') : $t('profile.bind_email')"
    width="500"
    @close="handleClose"
  >
    <Email ref="emailRef" @success="onChangeEmailSuccess" @close="emailVisible = false"></Email>
  </el-dialog>

  <ElDialog
    v-model="unbindWechatVisible"
    :show-close="false"
    class="!rounded-xl"
    header-class="px-4 py-3"
    body-class="px-4 pb-3"
    footer-class="px-4 pb-2"
    destroy-on-close
    append-to-body
    width="480"
  >
    <template #header>
      <div class="flex items-center">
        <ElIcon color="#FF9500" size="24">
          <WarnTriangleFilled />
        </ElIcon>
        <h2 class="flex-1 w-0 text-primary font-semibold text-lg ml-2">
          {{ $t('profile.unbind_wechat_confirm_title') }}
        </h2>
        <ElIcon color="#909399" size="18" class="cursor-pointer" @click="unbindWechatVisible = false">
          <Close />
        </ElIcon>
      </div>
    </template>
    {{ $t('profile.unbind_wechat_confirm_desc') }}
    <template #footer>
      <ElButton class="bg-[#F6F7F9] !h-[36px]" plain size="large" @click="unbindWechatVisible = false">
        {{ $t('profile.unbind_wechat_confirm_cancel') }}
      </ElButton>
      <ElButton class="!h-[36px]" type="danger" size="large" @click="handleUnbindWechat({ confirm: true })">
        {{ $t('profile.unbind_wechat_confirm_ok') }}
      </ElButton>
    </template>
  </ElDialog>
  <ElDialog
    v-model="bindWechatVisible"
    :show-close="false"
    class="!rounded-xl"
    header-class="px-4 py-4"
    body-class="px-4 pb-3"
    footer-class="px-4 pb-2"
    destroy-on-close
    append-to-body
    width="500"
  >
    <template #header>
      <div class="flex items-center justify-center relative">
        <img class="size-6" :src="$getPublicPath('/images/profile/wechat.png')" alt="" />
        <h2 class="text-primary font-semibold text-lg ml-2">
          {{ $t('profile.bind_wechat_title') }}
        </h2>
        <ElIcon color="#909399" size="20" class="cursor-pointer absolute right-0 -top-2" @click="bindWechatVisible = false">
          <Close />
        </ElIcon>
      </div>
    </template>
    <WechatView v-if="qrcodeRender" class="h-[280px] overflow-hidden" @oauth-success="handleOauthSuccess" />
  </ElDialog>
</template>

<script setup lang="ts">
import { ArrowRight, Close, WarnTriangleFilled } from '@element-plus/icons-vue'
import { ref, reactive, nextTick, onMounted } from 'vue'
import enterpriseApi from '@/api/modules/enterprise'
import CropperDialog from '@/components/CropperDialog/index.vue'
// import Password from './password.vue'
import Email from '@/components/LoginModal/email.vue'
import ResetPassword from '@/components/LoginModal/resetPassword.vue'
import changeMobile from './changeMobile.vue'
import WechatView from '@/components/LoginModal/wechat.vue'

import useEnv from '@/hooks/useEnv'

import { useUserStore } from '@/stores/modules/user'

const { isOpLocalEnv } = useEnv()

const userStore = useUserStore()

const corpperRef = ref()
const mobileRef = ref()
const passwordRef = ref()
const emailRef = ref()

const profileVisible = ref(false)
// const passwordVisible = ref(false)
const emailVisible = ref(false)
const reset_Password_Visible = ref(false)
const change_Mobile_Visible = ref(false)
const openSMTP = ref(false)

const handleClose = () => {
  passwordRef.value?.resetForm?.()
  mobileRef.value?.resetForm?.()
  emailRef.value?.resetForm?.()
}

const profileForm = reactive({
  nickname: '',
  avatar: ''
})

const handleEdit = () => {
  profileForm.avatar = userStore.info.avatar
  profileForm.nickname = userStore.info.nickname

  profileVisible.value = true
}

const handleChangeAvatar = () => {
  corpperRef.value.uploadFile()
}
const handleSuccessCropper = (data) => {
  profileForm.avatar = data.url
}

const handleSaveProfile = () => {
  return userStore.update(profileForm).then(() => {
    profileVisible.value = false
    ElMessage.success($t('status.updated'))
  })
}

const onChangeEmailSuccess = () => {
  userStore.getUserInfo()
  emailVisible.value = false
}

const onChangeMobileSuccess = () => {
  userStore.getUserInfo()
  change_Mobile_Visible.value = false
}
const unbindWechatVisible = ref(false)
const handleUnbindWechat = async ({ confirm = false } = {}) => {
  if (confirm) {
    await userStore.unbind_wechat()
    ElMessage.success(window.$t('profile.unbind_success'))
  }
  unbindWechatVisible.value = !confirm
}
const bindWechatVisible = ref(false)
const handleBindWechat = () => {
  bindWechatVisible.value = true
}
const qrcodeRender = ref(true)
const handleOauthSuccess = async (data: any) => {
  await userStore.bind_wechat({ openid: data.openid, unionid: data.unionid }).catch(async (err) => {
    qrcodeRender.value = false
    await nextTick()
    qrcodeRender.value = true
    return Promise.reject(err)
  })
  ElMessage.success(window.$t(userStore.info.openid ? 'profile.change_success' : 'profile.bind_success'))
  bindWechatVisible.value = false
}

const loadSMTP = async () => {
  const { data } = await enterpriseApi.getSMTPInfo('smtp')
  openSMTP.value = data
}

onMounted(() => {
  loadSMTP()
})
</script>

<style scoped>
.account-item {
  position: relative;
}

.account-item::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 48px;
  width: 100%;
  border-bottom: 1px solid #edeeef;
}

.account-item:last-child::after {
  display: none;
}
</style>
<style>
.update-user-dialog {
  --el-dialog-border-radius: 12px;
  --el-dialog-padding-primary: 30px;
}

.update-user-dialog .el-dialog__headerbtn {
  --el-color-info: #939499;

  top: 16px;
  right: 16px;
  font-size: 18px;
}
</style>
