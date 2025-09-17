<template>
  <ElContainer class="w-screen h-screen bg-white">
    <ElAside
      width="55%"
      class="relative bg-[url('/images/login/background.png')] bg-cover bg-center bg-no-repeat max-md:hidden"
    >
      <!-- #ifndef KM -->
      <img
        class="w-[24%] object-contain absolute top-8 left-10"
        :src="$getRealPath({ url: '/images/logo_2.png' })"
        alt=""
      />
      <img
        class="w-[48%] object-contain absolute top-[25%] left-1/2 -translate-x-1/2"
        :src="$getRealPath({ url: '/images/login/title.png' })"
        alt=""
      />
      <img
        class="w-[70%] object-contain absolute top-[32%] left-1/2 -translate-x-1/2"
        :src="$getRealPath({ url: '/images/login/demo.png' })"
        alt=""
      />
      <!-- #endif -->

      <!-- #ifdef KM -->
      <img
        class="w-[24%] object-contain absolute top-8 left-10"
        :src="$getRealPath({ url: '/images/km-logo.png' })"
        alt=""
      />
      <img
        class="w-[68%] object-contain absolute top-[25%] left-1/2 -translate-x-1/2"
        :src="$getRealPath({ url: '/images/login/km-title.png' })"
        alt=""
      />
      <img
        class="w-[80%] object-contain absolute top-[32%] left-1/2 -translate-x-1/2"
        :src="$getRealPath({ url: '/images/login/km-demo.png' })"
        alt=""
      />
      <!-- #endif -->
    </ElAside>
    <ElMain class="relative flex flex-col justify-center items-center pt-10 px-6 box-border overflow-auto">
      <div class="absolute top-6 right-6 flex gap-2 text-[#777] text-base">
        <!-- <SvgIcon class="cursor-pointer hover:opacity-70" name="github" width="24" @click="onGithubOpen" /> -->
        <LanguageDropdown />
        <ElButton
          class="!p-0 !border-none !outline-none !bg-transparent !leading-1 !h-auto hover:opacity-70"
          type="default"
          text
          @click="onOfficialWebsiteOpen"
        >
          <SvgIcon class="mr-1" name="home_v2" width="16" />
          {{ $t('official_website') }}
        </ElButton>
        <!-- <ElButton class="!p-0 !border-none !outline-none !bg-transparent !leading-1 !h-auto hover:opacity-70"
					type="default" text @click="onApplyOpen">
					<SvgIcon class="mr-1" name="apply" width="16" />
					{{ $t('apply_open') }}
				</ElButton> -->
      </div>
      <LoginForm
        v-if="form_type === FORM_TYPE.LOGIN"
        @forget="onForgetOpen"
        @register="onRegisterOpen"
        @apply="onApplyOpen"
      />
      <ApplyForm v-else-if="form_type === FORM_TYPE.APPLY" @login="onLoginOpen" />
      <ForgetForm v-else-if="form_type === FORM_TYPE.FORGET" @login="onLoginOpen" @register="onRegisterOpen" />
      <RegisterForm v-else-if="form_type === FORM_TYPE.REGISTER" @login="onLoginOpen" />
    </ElMain>
  </ElContainer>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import LoginForm from './components/login-form.vue'
import ForgetForm from './components/forget-form.vue'
import ApplyForm from './components/apply-form.vue'
import RegisterForm from './components/register-form.vue'

const FORM_TYPE = {
  LOGIN: 'login',
  APPLY: 'apply',
  FORGET: 'forget',
  REGISTER: 'register',
} as const

type FormType = (typeof FORM_TYPE)[keyof typeof FORM_TYPE]

const form_type = ref<FormType>(FORM_TYPE.LOGIN)

const onOfficialWebsiteOpen = () => {
  window.open('https://hub.53ai.com', '_blank')
}

const onApplyOpen = async () => {
  form_type.value = FORM_TYPE.APPLY
}
const onForgetOpen = () => {
  form_type.value = FORM_TYPE.FORGET
}
const onLoginOpen = () => {
  form_type.value = FORM_TYPE.LOGIN
}
const onRegisterOpen = () => {
  form_type.value = FORM_TYPE.REGISTER
}
</script>

<style lang="scss" scoped></style>
