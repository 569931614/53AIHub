<template>
  <div class="h-full flex flex-col px-[60px] py-8">
    <Header :title="$t('module.homepage')">
      <template #right>
        <LanguageDropdown />
      </template>
    </Header>

    <div
      v-loading="loading"
      class="flex-1 overflow-y-auto flex bg-white mt-4 py-[67px] box-border px-[76px] rounded-lg"
    >
      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-4">
          <img class="h-[70px] object-contain" :src="enterpriseInfo.logo" alt="" />
          <h3 class="text-2xl text-[#1D1E1F] font-semibold truncate">
            {{ enterpriseInfo.name }}
          </h3>
          <SvgIcon
            v-router.push="{ name: 'Info' }"
            class="cursor-pointer hover:opacity-60"
            name="edit"
            color="#2563EB"
          />
        </div>
        <div class="text-sm text-[#9A9A9A] mt-3">
          {{ enterpriseInfo.description || '' }}
        </div>
        <div class="mt-10 flex flex-col gap-6">
          <!-- #ifndef KM -->
          <div class="flex items-center gap-8">
            <div class="flex-none w-[64px] text-base text-[#9A9A9A]">
              {{ isOpLocalEnv ? $t('website_home') : $t('website_domain') }}
            </div>
            <div v-if="isSaasLogin || isOpLocalEnv" class="flex items-center gap-2">
              <ElLink
                style="--el-link-text-color: #1d1e1f; --el-link-font-size: 16px"
                target="_blank"
                :href="`${domainUrl}?access_token=${userInfo.access_token}&eid=${enterpriseInfo.eid}`"
              >
                {{ domainUrl }}
                <SvgIcon class="cursor-pointer ml-1" name="blank" width="16" color="#2563EB" />
              </ElLink>
              <SvgIcon
                v-if="!isOpLocalEnv"
                v-router.push="{ name: 'Domain' }"
                class="cursor-pointer hover:opacity-60"
                name="edit"
                width="16"
                color="#2563EB"
              />
            </div>
          </div>
          <!-- #endif -->
          <!-- #ifdef KM -->
          <div class="flex items-center gap-8">
            <div class="flex-none w-[64px] text-base text-[#9A9A9A]">
              {{ $t('website_home') }}
            </div>
            <div v-if="isSaasLogin || isOpLocalEnv" class="flex items-center gap-2">
              <ElLink
                style="--el-link-text-color: #1d1e1f; --el-link-font-size: 16px"
                target="_blank"
                :href="`${defaultDomain}`"
              >
                {{ defaultDomain }}
                <SvgIcon class="cursor-pointer ml-1" name="blank" width="16" color="#2563EB" />
              </ElLink>
            </div>
          </div>
          <!-- #endif -->

          <div v-if="isSaasLogin" class="flex items-center gap-8">
            <div class="flex-none w-[64px] text-base text-[#9A9A9A]">
              {{ $t('version.title') }}
            </div>
            <div class="flex items-center gap-3">
              <div class="text-base text-[#1D1E1F]">
                {{ enterpriseStore.version.name }}
              </div>
              <el-button
                v-if="enterpriseInfo.version <= 3"
                link
                type="primary"
                class="text-base"
                @click="handleService('upgrade')"
              >
                {{ $t('action_upgrade') }}
              </el-button>
            </div>
          </div>
          <div class="flex items-center gap-8">
            <div class="flex-none w-[64px] text-base text-[#9A9A9A]">
              {{ $t('module.agent') }}
            </div>
            <div class="flex items-center gap-3">
              <div class="text-base text-[#1D1E1F]">
                {{ formatNumber(+indexInfo.agent_count || 0) }} / {{ getModuleMax(VERSION_MODULE.AGENT) }}
              </div>
            </div>
          </div>
          <div class="flex items-center gap-8">
            <div class="flex-none w-[64px] text-base text-[#9A9A9A]">
              {{ $t('module.prompt') }}
            </div>
            <div class="flex items-center gap-3">
              <div class="text-base text-[#1D1E1F]">{{ formatNumber(+indexInfo.prompt_count || 0) }} / ∞</div>
            </div>
          </div>
          <div class="flex items-center gap-8">
            <div class="flex-none w-[64px] text-base text-[#9A9A9A]">
              {{ $t('module.ai_toolbox') }}
            </div>
            <div class="flex items-center gap-3">
              <div class="text-base text-[#1D1E1F]">{{ formatNumber(+indexInfo.ai_link_count || 0) }} / ∞</div>
            </div>
          </div>
          <div class="flex items-center gap-8">
            <div class="flex-none w-[64px] text-base text-[#9A9A9A]">
              {{ $t('register_user.title') }}
            </div>
            <div class="flex items-center gap-3">
              <div class="text-base text-[#1D1E1F]">
                {{ formatNumber(+indexInfo.user_count || 0) }} / {{ getModuleMax(VERSION_MODULE.REGISTERED_USER) }}
              </div>
            </div>
          </div>
          <div class="flex items-center gap-8">
            <div class="flex-none w-[64px] text-base text-[#9A9A9A]">
              {{ $t('create_time') }}
            </div>
            <div class="flex items-center gap-3">
              <div class="text-base text-[#1D1E1F]">
                {{ (enterpriseInfo.created_time || '').substr(0, 16) }}
              </div>
            </div>
          </div>
          <div v-if="isSaasLogin" class="flex items-center gap-8">
            <div class="flex-none w-[64px] text-base text-[#9A9A9A]">
              {{ $t('service_expired_time') }}
            </div>
            <div class="flex items-center gap-3">
              <div class="text-base text-[#1D1E1F]">
                {{ enterpriseInfo.expired_time || $t('apply.expired_time_forever') }}
              </div>
              <el-button
                v-if="enterpriseInfo.expired_time"
                link
                type="primary"
                class="text-base"
                @click="handleService('renew')"
              >
                {{ $t('action_renew_v2') }}
              </el-button>
            </div>
          </div>
        </div>
        <div class="w-full h-[67px]" />
      </div>
      <div class="ml-[48px] pt-[68px] flex flex-col items-end">
        <div class="flex flex-col gap-3 w-[238px]">
          <div class="flex justify-between text-sm text-[#4F5052]">
            <span>{{ $t('configuration_completion') }}</span>
            <span>100%</span>
          </div>
          <el-progress :percentage="100" :show-text="false" />
        </div>
        <div class="w-[312px] h-[220px] mt-[74px] flex flex-col items-center justify-center bg-[#F5F7FA] group">
          <img
            class="w-10 h-8 mb-8 transition-all duration-300 ease-in-out group-hover:hidden"
            :src="$getRealPath({ url: '/images/index/wechat.png' })"
            alt=""
          />
          <img
            class="w-[120px] h-[120px] mb-2 transition-all duration-300 ease-in-out group-hover:block hidden"
            :src="$getRealPath({ url: '/images/index/qrcode.png' })"
            alt=""
          />
          <h6 class="text-xl text-[#1D1E1F] font-medium mb-2">
            {{ $t('join_group') }}
          </h6>
          <p class="text-sm text-[#1D1E1F]">
            {{ $t('join_group_desc') }}
          </p>
        </div>
        <!-- <img class="w-[312px] h-[220px] mt-[74px]" :src="$getRealPath({ url: '/images/intro.png' })" alt=""> -->
      </div>
    </div>
    <div class="text-sm text-[#9A9A9A] text-center py-11">
      {{ $t('copyright_desc', { year }) }}
    </div>

    <ServiceDialog
      v-model:visible="service.visible"
      :title="service.type === 'upgrade' ? $t('action_upgrade') : $t('version.scan_consult')"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import ServiceDialog from '@/components/ServiceDialog/index.vue'

import { useEnterpriseStore, useUserStore } from '@/stores'
import { useEnv } from '@/hooks/useEnv'
import { VERSION_MODULE } from '@/constants/enterprise'

const { isOpLocalEnv } = useEnv()
const enterpriseStore = useEnterpriseStore()
const userStore = useUserStore()
const year = new Date().getFullYear()
const loading = ref(false)
const indexInfo = ref({})
const service = ref({
  type: 'upgrade',
  visible: false,
})

const userInfo = computed(() => userStore.info)
const enterpriseInfo = computed(() => enterpriseStore.info)
const isSaasLogin = computed(() => userStore.is_saas_login)
const defaultDomain = computed(() => `${window.location.origin}/#/index`)
const domainUrl = computed(() => (isOpLocalEnv.value ? defaultDomain.value : enterpriseInfo.value.domain))

const formatNumber = (num = 0) => num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')

const getModuleMax = (module: string) => {
  return enterpriseStore.version?.features[module]?.max || '∞'
}

const handleService = (type: string) => {
  service.value.type = type
  service.value.visible = true
}

onMounted(async () => {
  loading.value = true
  indexInfo.value = await enterpriseStore.loadHomeInfo().finally(() => {
    loading.value = false
  })
})
</script>

<style scoped></style>
