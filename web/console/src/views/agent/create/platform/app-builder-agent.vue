<template>
  <div :class="[showChannelConfig ? '' : 'py-7']">
    <ElForm ref="formRef" :model="store.form_data" label-width="104px" label-position="top">
      <template v-if="showChannelConfig">
        <div class="text-base text-[#1D1E1F] font-medium mb-4">
          {{ $t('app_builder') }}
        </div>

        <el-form-item :label="$t('module.website_info_name')">
          <el-select v-model="store.form_data.custom_config.provider_id" size="large" @change="onProviderChange">
            <el-option v-for="item in providers" :key="item.provider_id" :label="item.name" :value="item.provider_id" />
          </el-select>
        </el-form-item>

        <ElFormItem
          prop="custom_config.app_builder_bot_id"
          :label="$t('select_agent')"
          :rules="generateInputRules({ message: 'form_select_placeholder' })"
        >
          <SelectPlus
            v-model="store.form_data.custom_config.app_builder_bot_id"
            :use-i18n="false"
            size="large"
            :options="bots"
          />
        </ElFormItem>

        <div class="text-base text-[#1D1E1F] font-medium mb-4">
          {{ $t('basic_info') }}
        </div>
        <AgentInfo v-model="store.form_data" />
      </template>

      <template v-else>
        <BaseConfig />
        <RelateApp />
        <ExpandConfig />
        <UseScope />
      </template>
    </ElForm>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import AgentInfo from '../components/agent-info.vue'
import BaseConfig from '../components/base-config.vue'
import ExpandConfig from '../components/expand-config.vue'
import UseScope from '../components/use-scope.vue'
import RelateApp from '../components/relate-agents.vue'
import { useAgentFormStore } from '../store'
import { generateInputRules } from '@/utils/form-rule'

import { PROVIDER_VALUES } from '@/constants/platform/config'
import providersApi from '@/api/modules/providers/index'
import { transformProviderList } from '@/api/modules/providers/transform'
import agentApi, { AppBuilderBotItem, transformAppBuilderBotItem } from '@/api/modules/agent'
import { ProviderItem } from '@/api/modules/providers/types'

const props = defineProps({
  showChannelConfig: {
    type: Boolean,
    default: false,
  },
})

const store = useAgentFormStore()

const providers = ref<ProviderItem[]>([])
const bots = ref<AppBuilderBotItem[]>([])

const formRef = ref()

const loadBots = async () => {
  const list = await agentApi.appbuilder.bots_list({
    provider_id: store.form_data.custom_config.provider_id,
  })
  bots.value = list.map(transformAppBuilderBotItem)
}

const loadProviders = async () => {
  const list = await providersApi.list({
    providerType: PROVIDER_VALUES.APP_BUILDER,
  })
  providers.value = transformProviderList(list)

  if (providers.value.length && !store.form_data.custom_config.provider_id) {
    store.form_data.custom_config.provider_id = providers.value[0].provider_id
  }
  loadBots()
}

const onProviderChange = () => {
  loadBots()
  store.form_data.custom_config.app_builder_bot_id = ''
}

onMounted(() => {
  if (props.showChannelConfig) {
    loadProviders()
  }
})

const validateForm = async () => formRef.value.validate()

defineExpose({
  validateForm,
})
</script>

<style scoped></style>
