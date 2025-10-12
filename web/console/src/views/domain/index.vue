<template>
  <Layout class="px-[60px] py-8">
    <Header :title="$t('module.domain')" />

    <div class="flex-1 flex flex-col bg-white p-6 mt-3 box-border">
      <div v-loading="loading" class="flex-1 max-h-[calc(100vh-100px)] overflow-auto">
        <!-- 专属域名部分 -->
        <section class="mb-8">
          <h1 class="font-semibold text-[#1D1E1F]">
            {{ $t('module.domain_exclusive') }}
          </h1>
          <div class="mt-4 border rounded overflow-hidden p-6">
            <label class="text-[#1D1E1F] text-sm">{{ $t('module.domain_exclusive_label') }}</label>
            <div class="w-full mt-4 flex items-center gap-3">
              <ElInput
                v-model="exclusiveDomainUrl"
                class="!max-w-[600px]"
                link
                :placeholder="$t('form_input_placeholder')"
                disabled
                size="large"
              />
              <ElButton
                class="flex-none text-[#3664EF]"
                type="default"
                size="large"
                @click="handleCopyDomain(exclusiveDomainUrl)"
              >
                <ElIcon :size="16" class="mr-2" color="#3664EF">
                  <CopyDocument />
                </ElIcon>
                {{ $t('action_copy') }}
              </ElButton>
              <div class="flex-1 h-2" />
              <ElButton class="flex-none text-[#5A6D9E] !p-0" link size="large" @click="handleOpenExclusiveSetting">
                <ElIcon :size="16" class="mr-2" color="#5A6D9E">
                  <Setting />
                </ElIcon>
                {{ $t('action_setting') }}
              </ElButton>
            </div>
          </div>
        </section>

        <!-- 独立域名部分 -->
        <section>
          <h1 class="font-semibold text-[#1D1E1F]">
            {{ $t('module.domain_independent') }}
          </h1>
          <div
            v-version="{ module: VERSION_MODULE.INDEPENDENT_DOMAIN, mode: 'tooltip' }"
            class="mt-4 border rounded overflow-hidden p-6"
          >
            <label class="text-[#1D1E1F] text-sm">
              {{ $t('module.domain_independent_label') }}
              <template v-if="independentDomainUrl">
                <ElTag class="ml-3 !border-none !bg-[#E3F6E0] !text-[#09BB07]" type="success" size="default">
                  {{ $t('effective') }}
                </ElTag>
                <ElTag
                  v-if="independentDomainInfo.httpsEnabled"
                  class="ml-3 !border-none !bg-[#E3F6E0] !text-[#09BB07] flex items-center"
                  type="success"
                  size="default"
                >
                  <SvgIcon class="!inline-block" name="global" width="12" height="12" />
                  {{ $t('https_enabled') }}
                </ElTag>
              </template>
            </label>
            <div class="w-full mt-4 flex items-center gap-3">
              <template v-if="independentDomainUrl">
                <ElInput
                  v-model="independentDomainUrl"
                  class="!max-w-[600px]"
                  link
                  :placeholder="$t('form_input_placeholder')"
                  disabled
                  size="large"
                />
                <ElButton
                  class="flex-none text-[#3664EF]"
                  type="default"
                  size="large"
                  @click="handleCopyDomain(independentDomainUrl)"
                >
                  <ElIcon :size="16" class="mr-2" color="#3664EF">
                    <CopyDocument />
                  </ElIcon>
                  {{ $t('action_copy') }}
                </ElButton>
              </template>
              <div class="flex-1 text-sm text-[#9A9A9A]">
                {{ independentDomainUrl ? '' : $t('module.domain_independent_desc') }}
              </div>
              <ElButton class="flex-none text-[#5A6D9E] !p-0" link size="large" @click="handleOpenIndependentSetting">
                <ElIcon :size="16" class="mr-2" color="#5A6D9E">
                  <Setting />
                </ElIcon>
                {{ $t('action_setting') }}
              </ElButton>
              <ElButton
                v-if="independentDomainUrl"
                class="flex-none text-[#5A6D9E] !p-0 !ml-0"
                link
                size="large"
                @click="handleDeleteIndependentDomain"
              >
                <ElIcon :size="16" class="mr-2" color="#5A6D9E">
                  <Delete />
                </ElIcon>
                {{ $t('action_delete') }}
              </ElButton>
            </div>
          </div>
        </section>
      </div>
    </div>
  </Layout>

  <!-- 弹窗组件 -->
  <ExclusiveSettingDialog ref="exclusiveSettingRef" @success="loadDomainData" />
  <IndependentSettingDialog ref="independentSettingRef" @success="loadDomainData" />
</template>

<script setup lang="ts">
import { CopyDocument, Delete, Setting } from '@element-plus/icons-vue'
import { onMounted, ref, computed } from 'vue'
import ExclusiveSettingDialog from './components/exclusive-setting-dialog.vue'
import IndependentSettingDialog from './components/independent-setting-dialog.vue'
import { VERSION_MODULE } from '@/constants/enterprise'
import { copyToClip } from '@/utils/copy'
import { useEnv } from '@/hooks/useEnv'
import { domainApi } from '@/api/modules/domain/index'

// 类型定义
type DomainConfig = {
  enable_https?: string | number
  [key: string]: unknown
}

type DomainInfo = {
  id?: number
  domain?: string
  domain_name?: string
  config?: string | DomainConfig
  [key: string]: unknown
}

type IndependentDomainInfo = {
  httpsEnabled: boolean
  domainName: string
  rawData: DomainInfo
}

// 环境变量
const { isDevEnv } = useEnv()

// 组件引用
const exclusiveSettingRef = ref<InstanceType<typeof ExclusiveSettingDialog>>()
const independentSettingRef = ref<InstanceType<typeof IndependentSettingDialog>>()

// 响应式数据
const loading = ref(false)
const exclusiveDomainInfo = ref<DomainInfo>({})
const independentDomainInfo = ref<IndependentDomainInfo>({
  httpsEnabled: false,
  domainName: '',
  rawData: {},
})

// 计算属性
const exclusiveDomainUrl = computed(() => {
  const domainName = exclusiveDomainInfo.value.domain_name || ''
  if (!domainName) return ''
  return `https://${domainName}${isDevEnv.value ? '.hub' : ''}.53ai.com`
})

const independentDomainUrl = computed(() => {
  const { domainName, httpsEnabled } = independentDomainInfo.value
  if (!domainName) return ''
  return `http${httpsEnabled ? 's' : ''}://${domainName}`
})

// 数据处理方法（需要在 loadDomainData 之前定义）
const processExclusiveDomainData = (domainData: DomainInfo) => {
  exclusiveDomainInfo.value = domainData

  if (domainData.domain) {
    let domainName = domainData.domain
      .trim()
      .replace(/^https?:\/\//, '')
      .replace(/\.53ai\.com$/, '')

    if (isDevEnv.value) {
      domainName = domainName.replace(/\.hub$/, '')
    }

    exclusiveDomainInfo.value.domain_name = domainName
  }
}

const processIndependentDomainData = (domainData: DomainInfo) => {
  const rawData = { ...domainData }

  // 解析配置
  let config: DomainConfig = {}
  if (domainData.config) {
    try {
      config = typeof domainData.config === 'string' ? JSON.parse(domainData.config) : domainData.config
    } catch (error) {
      console.error('解析独立域名配置失败:', error)
      config = {}
    }
  }

  rawData.config = config

  // 处理域名信息
  const domainName = (domainData.domain || '').trim().replace(/^https?:\/\//, '')
  const httpsEnabled = Boolean(Number(config.enable_https))

  independentDomainInfo.value = {
    httpsEnabled,
    domainName,
    rawData,
  }
}

// 主要方法定义
const loadDomainData = async () => {
  loading.value = true

  try {
    const { exclusive_domains = [], independent_domains = [] } = await domainApi.list()

    // 处理专属域名数据
    processExclusiveDomainData(exclusive_domains[0] || {})

    // 处理独立域名数据
    processIndependentDomainData(independent_domains[0] || {})
  } catch (error) {
    console.error('加载域名数据失败:', error)
    ElMessage.error('加载域名数据失败')
  } finally {
    loading.value = false
  }
}

const handleCopyDomain = async (domainUrl: string) => {
  if (!domainUrl) {
    ElMessage.warning('没有可复制的域名')
    return
  }

  try {
    await copyToClip(domainUrl)
    ElMessage.success(window.$t('action_copy_success'))
  } catch (error) {
    console.error('复制失败:', error)
    ElMessage.error('复制失败')
  }
}

const handleOpenExclusiveSetting = () => {
  const settingData = {
    ...exclusiveDomainInfo.value,
    domain: exclusiveDomainUrl.value,
  }
  exclusiveSettingRef.value?.open({ data: settingData })
}

const handleOpenIndependentSetting = () => {
  independentSettingRef.value?.open({ data: independentDomainInfo.value.rawData })
}

const handleDeleteIndependentDomain = async () => {
  try {
    await ElMessageBox.confirm(window.$t('module.domain_independent_delete_confirm'))

    const domainId = independentDomainInfo.value.rawData.id
    if (!domainId) {
      ElMessage.error('域名ID不存在')
      return
    }

    await domainApi.deleteIndependent(domainId)

    // 重置独立域名信息
    independentDomainInfo.value = {
      httpsEnabled: false,
      domainName: '',
      rawData: {},
    }

    ElMessage.success(window.$t('action_delete_success'))
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除独立域名失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 生命周期
onMounted(() => {
  loadDomainData()
})
</script>

<style scoped lang="scss">
section {
  &:not(:last-child) {
    margin-bottom: 2rem;
  }
}
</style>
