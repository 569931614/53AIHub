<template>
  <div class="h-full overflow-y-auto">
    <div class="md:w-[90%] sm:w-[90%] lg:w-[65%] mx-auto pt-20 px-4">
      <div class="flex flex-col gap-3 items-center">
        <h3 class="text-[32px] font-bold">{{ $t('guide.title') }}</h3>
        <p class="text-sm text-[#2029459e]">{{ $t('guide.description') }}</p>
      </div>
      <el-steps class="mt-12" :active="steps" align-center>
        <el-step :title="$t('guide.website_info')" />
        <el-step :title="$t('guide.website_setting')" />
        <el-step :title="$t('guide.website_success')" />
      </el-steps>

      <ElForm ref="formRef" :model="form" :rules="rules" label-position="top" class="flex flex-col mt-10">
        <template v-if="isFirstStep">
          <ElFormItem :label="$t('guide.website_info_name')" prop="name">
            <ElInput
              v-model="form.name"
              class="max-w-full h-11"
              :placeholder="$t('guide.website_info_name_placeholder')"
              size="large"
              clearable
              maxlength="120"
              show-word-limit
            />
          </ElFormItem>

          <ElFormItem :label="$t('guide.website_info_logo')" prop="logo">
            <div class="mt-4 w-full flex items-center gap-4">
              <ElImage
                v-if="form.logo"
                class="h-[70px] w-[70px] rounded overflow-hidden"
                :src="form.logo"
                :preview-src-list="[form.logo]"
                fit="contain"
              />
              <UploadLogo
                v-model="form.logo"
                class="w-auto h-auto"
                show-text
                :text="$t(form.logo ? 'guide.website_info_logo_change' : 'guide.website_info_logo_upload')"
              />
            </div>
            <div class="mt-2 w-full text-sm text-[#9A9A9A]">
              {{ $t('guide.website_info_logo_tip') }}
            </div>
          </ElFormItem>

          <ElFormItem :label="$t('guide.website_style')">
            <ul class="flex flex-wrap gap-4">
              <li
                v-for="value in [WEBSITE_STYLE.WEBSITE, WEBSITE_STYLE.SOFTWARE]"
                :key="value"
                class="w-[172px] p-1.5 bg-[#F5F5F5] flex relative flex-col cursor-pointer items-center border rounded box-border overflow-hidden text-sm hover:border-[#3664EF] hover:text-[#3664EF]"
                :class="[form.type === value ? 'border-[#3664EF] text-[#3664EF]' : 'text-[#4F5052]']"
                @click.stop="form.type = value"
              >
                <div v-if="form.type === value">
                  <div class="right-angle-triangle"></div>
                  <SvgIcon name="tick" stroke="true" class="absolute !w-4 h-2 top-0 right-0 text-[#ffffff]" />
                </div>
                <div class="text-sm p-1.5">
                  {{ $t(WEBSITE_STYLE_LABEL_MAP.get(value)) }}
                </div>
                <ElImage class="w-full mt-2" :src="WEBSITE_STYLE_DEMO_MAP.get(value)" fit="contain" />
              </li>
            </ul>
          </ElFormItem>

          <ElFormItem :label="$t('guide.website_info_language')">
            <ElSelect v-model="form.language" class="h-11" size="large" @change="handleLanguageChange">
              <ElOption v-for="item in languageOptions" :key="item.value" :label="item.label" :value="item.value" />
            </ElSelect>
          </ElFormItem>
        </template>

        <template v-else-if="isSecondStep">
          <ElFormItem :label="$t('form.account')" prop="account" :rules="[getAccountOrEmailRules()]">
            <ElInput v-model="form.account" class="h-11" :placeholder="$t('form.account_format')" size="large" clearable />
          </ElFormItem>

          <el-form-item :label="$t('form.password')" prop="password" :rules="[getPasswordRules()]">
            <el-input v-model="form.password" v-trim show-password size="large" :placeholder="$t('form.password_placeholder')"></el-input>
          </el-form-item>

          <el-form-item
            :label="$t('guide.confirm_password')"
            prop="confirm_password"
            :rules="[getPasswordRules(), getConfirmPasswordRules(form, 'password')]"
          >
            <el-input
              v-model="form.confirm_password"
              v-trim
              show-password
              size="large"
              :placeholder="$t('guide.confirm_password_placeholder')"
            ></el-input>
          </el-form-item>
        </template>

        <template v-else-if="isThirdStep">
          <div class="flex flex-col items-center justify-center py-10">
            <h4 class="text-xl font-semibold text-[#202945] mb-2">{{ $t('guide.init_success') }}</h4>
            <p
              class="text-sm text-[#2029459e]"
              v-html="
                $t('guide.jump_tip', {
                  count: `<span style='color: #3664EF; font-weight: 500'>${countdown}</span>`
                })
              "
            ></p>
            <el-button type="text" class="mt-4 text-[#3664EF]" @click="handleManualJump">{{ $t('guide.jump_now') }}</el-button>
          </div>
        </template>
      </ElForm>

      <ElButton v-if="isFirstStep || isSecondStep" type="primary" class="w-full h-11 my-6" round @click="handleNext">
        {{ isFirstStep ? $t('guide.next') : $t('guide.init') }}
      </ElButton>
    </div>
  </div>
</template>

<script setup>
import { ElButton, ElForm } from 'element-plus'
import { ref, reactive, onUnmounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import UploadLogo from '@/components/Upload/image.vue'
import { WEBSITE_STYLE, WEBSITE_STYLE_LABEL_MAP, WEBSITE_STYLE_DEMO_MAP } from '@/constants/website'
import { getAccountOrEmailRules, getPasswordRules, getConfirmPasswordRules, generateInputRules } from '@/utils/form-rules'
import { enterprise } from '@/api/modules/enterprise'
import userApi from '@/api/modules/user'

// 常量定义
const DEFAULT_LOCAL_EID = 1
const COUNTDOWN_DURATION = 2
const COUNTDOWN_INTERVAL = 1000

// 获取 i18n 实例和路由
const { locale } = useI18n()
const router = useRouter()

// 响应式数据
const steps = ref(1)
const formRef = ref()
const countdown = ref(COUNTDOWN_DURATION)
let jumpTimer = null

// 表单数据
const form = reactive({
  name: '',
  logo: 'https://hub.53ai.com/console/images/default_website_logo.png',
  type: WEBSITE_STYLE.WEBSITE,
  language: 'zh-cn',
  account: '',
  password: '',
  confirm_password: ''
})

// 语言选项配置
const languageOptions = [
  { label: '中文-CN', value: 'zh-cn' },
  { label: '中文-TW', value: 'zh-tw' },
  { label: '英文-EN', value: 'en' },
  { label: '日文-JP', value: 'jp' }
]

// 默认模板配置
const defaultTemplateConfig = {
  style_type: 'software',
  theme_color: '#3664EF',
  text_color: '#333333',
  nav_bg_color: '#ffffff',
  nav_text_color: '#333333',
  page_footer_bg_color: '#18191F',
  page_footer_text_color: '#F2F2F2'
}

// 表单验证规则
const rules = reactive({
  name: generateInputRules({ message: 'guide.website_info_name_placeholder' }),
  logo: generateInputRules({ message: 'guide.website_info_logo_placeholder' })
})

// 计算属性
const isFirstStep = computed(() => steps.value === 1)
const isSecondStep = computed(() => steps.value === 2)
const isThirdStep = computed(() => steps.value === 3)

// 语言切换处理函数
const handleLanguageChange = () => {
  locale.value = form.language
}

// 清理定时器
const clearJumpTimer = () => {
  if (jumpTimer) {
    clearInterval(jumpTimer)
    jumpTimer = null
  }
}

// 开始倒计时
const startJumpCountdown = () => {
  clearJumpTimer()
  countdown.value = COUNTDOWN_DURATION

  jumpTimer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      clearJumpTimer()
      router.push('/index')
    }
  }, COUNTDOWN_INTERVAL)
}

// 处理第一步：更新企业信息
const handleFirstStep = async () => {
  try {
    await enterprise.update(DEFAULT_LOCAL_EID, {
      display_name: form.name,
      logo: form.logo,
      template_type: JSON.stringify({
        ...defaultTemplateConfig,
        style_type: form.type
      }),
      language: form.language
    })
    steps.value++
  } catch (error) {
    console.error('更新企业信息失败:', error)
    throw error
  }
}

// 处理第二步：用户注册
const handleSecondStep = async () => {
  try {
    await userApi.register({
      username: form.account,
      password: form.password,
      nickname: form.account
    })
    steps.value++
    startJumpCountdown()
  } catch (error) {
    console.error('用户注册失败:', error)
    throw error
  }
}

// 下一步处理
const handleNext = async () => {
  try {
    const valid = await formRef.value.validate()
    if (!valid) return

    if (isFirstStep.value) {
      await handleFirstStep()
    } else if (isSecondStep.value) {
      await handleSecondStep()
    }
  } catch (error) {
    console.error('操作失败:', error)
  }
}

// 手动跳转
const handleManualJump = () => {
  clearJumpTimer()
  router.push('/index')
}

// 组件卸载时清理定时器
onUnmounted(() => {
  clearJumpTimer()
})
</script>

<style scoped>
.right-angle-triangle {
  width: 0;
  height: 0;
  border-top: 31px solid #3664ef;
  border-left: 29px solid transparent;
  position: absolute;
  top: 0;
  right: 0;
}

:deep(.el-step__icon.is-text) {
  width: 36px;
  height: 36px;
  border-radius: 50%;
}

:deep(.el-step__title.is-wait) {
  color: #c0c2c4;
  font-weight: bold;
}

:deep(.el-step__title.is-process) {
  color: #c0c2c4;
}

:deep(.el-step__head.is-process > .el-step__icon) {
  background: #c0c2c4;
  color: #fff;
  margin-left: 12px;
}

:deep(.el-step__head.is-wait > .el-step__icon) {
  background: #c0c2c4;
  color: #fff;
  margin-left: 12px;
}

:deep(.el-step__head.is-finish > .el-step__icon) {
  background: #409eff;
  color: #fff;
  margin-left: 12px;
}

:deep(.el-step__line) {
  background-color: #dcdfe6;
  height: 2px;
  margin-top: 5px;
}
</style>
