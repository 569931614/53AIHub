<template>
  <Layout class="px-[60px] py-8">
    <Header :title="$t(route.meta.title)" />
    <div class="bg-white rounded-md overflow-hidden py-8 px-10 mt-5 box-border space-y-10 flex-1 overflow-y-auto">
      <div v-for="item in ssoList" :key="item.type">
        <h1 class="text-lg text-[#1D1E1F] flex items-center gap-2">
          <img :src="$getRealPath({ url: `/images/sso/${item.type}.png` })" class="size-6 object-cover" />
          {{ $t(`sso.${item.type}.title`) }}
        </h1>
        <div class="mt-4 flex flex-wrap items-center gap-4">
          <div class="flex-none w-[300px] min-h-[176px] rounded overflow-hidden border p-5 box-border">
            <h2 class="font-semibold text-[#1D1E1F] h-[26px]">
              {{ $t(`sso.${item.type}.access_title`) }}
              <ElTag v-if="installedMap[item.type]" class="ml-2 -mt-1" type="success">
                {{ $t('connected') }}
              </ElTag>
            </h2>
            <p class="text-xs text-[#939499] mt-3">
              {{
                installedMap[item.type]
                  ? `${$t('sso.authorized_enterprise')}： ${enterpriseStore.info.wecom_info.corp_name}`
                  : $t(`sso.${item.type}.access_desc`)
              }}
            </p>
            <div class="mt-12">
              <ElButton v-if="installedMap[item.type]" @click="handleUnauthorized({ type: item.type })">
                {{ $t('sso.unauthorized') }}
              </ElButton>
              <template v-else>
                <ElButton
                  type="primary"
                  :disabled="['dingtalk', 'ad_ldap', 'feishu'].includes(item.type)"
                  @click="handleAuthorized({ type: item.type })"
                >
                  {{ $t('sso.auth_access') }}
                </ElButton>
              </template>
            </div>
          </div>
          <div
            v-if="item.type === 'wecom'"
            class="relative flex-none w-[300px] min-h-[176px] rounded overflow-hidden border p-5 box-border"
          >
            <h2 class="font-semibold text-[#1D1E1F]">
              {{ $t(`sso.organization_sync`) }}
            </h2>
            <p class="text-xs text-[#939499] mt-3">
              {{ $t(`sso.${item.type}.sync_desc`) }}
            </p>
            <div v-if="installedMap[item.type]" class="absolute right-5 top-4">
              <ElSwitch
                v-debounce
                :model-value="syncValue.value === item.sync_value"
                @change="handleSyncToggle($event, item.sync_value)"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
    <AccessDialog ref="accessRef" />
  </Layout>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import AccessDialog from './components/access-dialog.vue'

import { useEnterpriseStore } from '@/stores'

import { useSso } from '@/hooks/useSso'

const route = useRoute()
const enterpriseStore = useEnterpriseStore()
const { syncValue, loadSyncSetting, saveSyncSetting } = useSso()

const accessRef = ref<InstanceType<typeof AccessDialog>>()

const ssoList = ref([
  {
    type: 'wecom',
    title: 'sso.wecom.title',
    access_title: 'sso.wecom.access_title',
    access_desc: 'sso.wecom.access_desc',
    sync_desc: 'sso.wecom.sync_desc',
    sync_value: '1',
  },
  {
    type: 'dingtalk',
    title: 'sso.dingtalk.title',
    access_title: 'sso.dingtalk.access_title',
    access_desc: 'sso.dingtalk.access_desc',
    sync_desc: 'sso.dingtalk.sync_desc',
  },
  {
    type: 'ad_ldap',
    title: 'sso.ad_ldap.title',
    access_title: 'sso.ad_ldap.access_title',
    access_desc: 'sso.ad_ldap.access_desc',
    sync_desc: 'sso.ad_ldap.sync_desc',
  },
  {
    type: 'feishu',
    title: 'sso.feishu.title',
    access_title: 'sso.feishu.access_title',
    access_desc: 'sso.feishu.access_desc',
    sync_desc: 'sso.feishu.sync_desc',
  },
])

const installedMap = computed(() => {
  return {
    wecom: enterpriseStore.info.is_install_wecom,
  }
})

const handleAuthorized = ({ type }: { type: string }) => {
  accessRef.value?.open({ type })
}

const handleUnauthorized = ({ type }: { type: string }) => {
  let url = ''
  switch (type) {
    case 'wecom':
      url = 'https://work.weixin.qq.com/login'
      break
    case 'dingtalk':
      url = 'https://open-dev.dingtalk.com/fe/app?hash=%23%2Fisv%2Fapp#/isv/app'
      break
  }
  window.open(url, '_blank')
}
const handleSyncToggle = async (checked: boolean, value: string) => {
  await saveSyncSetting(checked ? value : '0')
}

onMounted(() => {
  loadSyncSetting()
})
onUnmounted(() => {})
</script>

<style scoped></style>
