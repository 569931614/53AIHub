<template>
  <li
    v-loading="provider.channelLoading"
    class="flex-none min-w-[246px] w-[24%] h-[178px] flex flex-col border rounded box-border overflow-hidden"
  >
    <div class="flex items-center gap-4 p-5 box-border">
      <img
        class="flex-none size-10 overflow-hidden"
        :src="$getRealPath({ url: `/images/platform/${provider.icon}.png` })"
      />
      <div class="text-[#1B2B51] font-semibold">
        {{ $t(provider.label) }}
      </div>
    </div>
    <div class="text-xs text-[#4F5052] px-5 box-border">
      <template v-if="!provider.auth">
        {{ $t('connecting_agent_total', { total: provider.agentTotal }) }}
      </template>
      <template v-else-if="provider.connected && provider.authed_time">
        {{ $t('connected') }} · {{ $t('authorized_at') }} {{ provider.authed_time.slice(0, 16) }}
      </template>
      <template v-else-if="provider.connected">
        {{ $t('connecting') }}
      </template>
      <template v-else>
        {{ $t('not_connected') }}
      </template>
    </div>
    <div class="flex-1 w-full" />
    <div class="w-full h-11 flex border-t box-border">
      <template v-if="!provider.auth">
        <ElButton
          class="flex-1 h-[46px] text-[#3664EF] !border-none !outline-none rounded-none"
          type="text"
          size="default"
          @click.stop="$emit('authorize', { data: provider })"
        >
          {{ $t('action_manage') }}
        </ElButton>
        <ElDivider class="!h-full" direction="vertical" />
        <ElButton
          v-version="{ module: VERSION_MODULE.AGENT, count: allTotal, content: $t('version.agent_limit') }"
          class="flex-1 h-[46px] !border-none !outline-none rounded-none"
          link
          size="default"
          @click.stop="$emit('add', { data: provider })"
        >
          {{ $t('action_add') }}
        </ElButton>
      </template>
      <template v-else-if="provider.connected">
        <ElButton
          class="flex-1 h-[46px] !border-none !outline-none rounded-none"
          link
          type="primary"
          size="default"
          @click.stop="$emit('authorize', { data: provider })"
        >
          {{ $t('action_edit') }}
        </ElButton>
        <ElDivider class="!h-full" direction="vertical" />
        <ElButton
          class="flex-1 h-[46px] text-[#919499] !border-none !outline-none rounded-none"
          link
          size="default"
          @click.stop="$emit('delete', { data: provider })"
        >
          {{ $t('action_delete') }}
        </ElButton>
      </template>
      <template v-else>
        <ElButton
          class="flex-1 h-[46px] bg-[#F3F6FE] text-[#3664EF] !border-none !outline-none rounded-none"
          type="default"
          size="default"
          @click.stop="$emit('authorize', { data: provider })"
        >
          {{ $t('action_authorize') }}
        </ElButton>
      </template>
    </div>
  </li>
</template>

<script setup lang="ts">
import { VERSION_MODULE } from '@/constants/enterprise'

interface ProviderOption {
  id: number
  icon: string
  label: string
  auth: boolean
  connected: boolean
  authed_time: string
  agentTotal: number
  channelLoading: boolean
}

defineProps<{
  provider: ProviderOption
  allTotal: number
}>()

defineEmits<{
  authorize: [{ data: ProviderOption }]
  add: [{ data: ProviderOption }]
  delete: [{ data: ProviderOption }]
}>()
</script>
