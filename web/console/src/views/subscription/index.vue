<template>
  <Layout class="px-[60px] !overflow-y-auto relative">
    <Header class="sticky top-0 left-0 right-0 z-10 pt-8 pb-4 bg-[#F6F7F8]" :title="$t('module.subscription')" />
    <div class="flex flex-col bg-white px-6 pt-6 box-border overflow-hidden">
      <h1 class="flex-none text-sm text-[#4F5052]">{{ $t('module.subscription_header_title') }}</h1>
      <div class="overflow-x-auto mt-6 pb-1 pr-1">
        <ElForm
          ref="form_ref"
          v-loading="loading"
          class="min-h-[70vh] overflow-hidden flex gap-4 w-max"
          :model="subscription_list"
          label-position="top"
        >
          <div
            v-for="(item, index) in subscription_list"
            :key="index"
            class="flex-none w-[334px] bg-[#F9F9FC] p-4 box-border"
          >
            <ElFormItem>
              <div class="w-full flex items-center gap-2.5">
                <ElDropdown placement="bottom" trigger="click" @command="onIconChange($event, item)">
                  <img
                    :src="
                      !/\.png$/.test(item.logo_url)
                        ? $getRealPath({ url: `/images/subscription/${item.logo_url}.png` })
                        : item.logo_url
                    "
                    class="flex-none w-[36px] h-[36px] object-cover cursor-pointer"
                    alt="订阅图标"
                  />
                  <template #dropdown>
                    <ElDropdownMenu>
                      <ElDropdownItem v-for="i in [1, 2, 3, 4, 5]" :key="i" :command="i">
                        <img
                          :src="img_host + `/subscription/vip-${i}.png`"
                          class="w-[36px] h-[36px] object-cover"
                          :alt="`VIP ${i}`"
                        />
                      </ElDropdownItem>
                    </ElDropdownMenu>
                  </template>
                </ElDropdown>
                <ElFormItem
                  class="flex-1"
                  :prop="`${index}.group_name`"
                  :rules="generateInputRules({ message: 'form_input_placeholder' })"
                >
                  <ElInput
                    v-model="item.group_name"
                    size="large"
                    :placeholder="$t('form_input_placeholder')"
                    show-word-limit
                    maxlength="10"
                  />
                </ElFormItem>
                <ElIcon
                  class="flex-none cursor-pointer"
                  :class="[item.is_default ? 'invisible' : '']"
                  size="16"
                  color="#F04F4D"
                  @click="handleRemove({ data: item, index })"
                >
                  <Delete />
                </ElIcon>
              </div>
            </ElFormItem>
            <ElFormItem :label="$t('module.subscription_charge')">
              <div class="w-full flex items-center gap-2">
                <ElSelect v-model="item.year_info.currency" class="flex-none !w-[86px]" size="large">
                  <ElOption
                    v-for="option in unit_options"
                    :key="option.value"
                    :label="option.label"
                    :value="option.value"
                  />
                </ElSelect>
                <ElInputNumber
                  v-model="item.year_info.amount"
                  class="charge-point-input flex-1"
                  size="large"
                  :controls="false"
                  :precision="2"
                  :min="0"
                  :disabled="item.is_default"
                  :max="999999999999"
                  :placeholder="$t('form_input_placeholder')"
                  @keydown="numberInputKeydownHandler"
                />
                <div class="flex-none text-sm text-[#1D1E1F]">/{{ $t('yearly') }}</div>
              </div>
              <div class="mt-3 w-full flex items-center gap-2">
                <ElSelect v-model="item.month_info.currency" class="flex-none !w-[86px]" size="large">
                  <ElOption
                    v-for="option in unit_options"
                    :key="option.value"
                    :label="option.label"
                    :value="option.value"
                  />
                </ElSelect>
                <ElInputNumber
                  v-model="item.month_info.amount"
                  class="charge-point-input flex-1"
                  size="large"
                  :controls="false"
                  :precision="2"
                  :min="0"
                  :disabled="item.is_default"
                  :max="999999999999"
                  :placeholder="$t('form_input_placeholder')"
                  @keydown="numberInputKeydownHandler"
                />
                <div class="flex-none text-sm text-[#1D1E1F]">/{{ $t('monthly') }}</div>
              </div>
            </ElFormItem>
            <ElFormItem :label="$t('module.subscription_points')">
              <div class="w-full flex items-center gap-2">
                <ElInputNumber
                  v-model="item.point_month_info.amount"
                  class="charge-point-input flex-1"
                  size="large"
                  :controls="false"
                  :precision="0"
                  :min="0"
                  :max="999999999999"
                  :placeholder="$t('form_input_placeholder')"
                  @keydown="numberInputKeydownHandler"
                />
                <div class="flex-none text-sm text-[#1D1E1F]">/{{ $t('monthly') }}</div>
              </div>
            </ElFormItem>
            <ElDivider />
            <template v-if="item.agents && item.agents.length">
              <ElFormItem class="!mb-0" :label="$t('module.subscription_agent_bots')">
                <template v-for="(agent, agent_index) in item.agents" :key="agent_index">
                  <div class="w-full flex items-center gap-2 mb-3">
                    <img
                      :src="agent.logo"
                      class="flex-none w-[18px] h-[18px] object-contain rounded-full overflow-hidden"
                      :alt="agent.name"
                    />
                    <div class="flex-1 truncate text-sm text-[#4F5052]">{{ agent.name }}</div>
                  </div>
                </template>
              </ElFormItem>
              <ElDivider />
            </template>
            <ElFormItem class="!mb-0">
              <template #label>
                <div class="flex items-center gap-1">
                  {{ $t('module.subscription_ai_assistant') }}
                  <ElTooltip :content="$t('module.subscription_ai_assistant_tip')" placement="bottom">
                    <SvgIcon class="flex-none" name="question" color="#A0A7B8" width="14" height="14" />
                  </ElTooltip>
                  <div class="flex-1"></div>
                  <ElSwitch v-model="item.ai_enabled" class="ml-2" size="small" />
                </div>
              </template>
              <div class="flex flex-wrap items-center gap-2">
                <ElButton class="!p-2" type="default" plain size="default" :disabled="!item.ai_enabled">
                  <SvgIcon
                    class="mr-1.5"
                    :class="{ 'opacity-50': !item.ai_enabled }"
                    name="windows"
                    width="14"
                    height="14"
                  />
                  <span class="text-xs">Windows</span>
                </ElButton>
                <ElButton class="!p-2 !ml-0" type="default" plain size="default" :disabled="!item.ai_enabled">
                  <SvgIcon
                    class="mr-1.5"
                    :class="{ 'opacity-50': !item.ai_enabled }"
                    name="ios"
                    width="14"
                    height="14"
                  />
                  <span class="text-xs">macOS</span>
                </ElButton>
                <ElButton class="!p-2 !ml-0" type="default" plain size="default" :disabled="!item.ai_enabled">
                  <SvgIcon
                    class="mr-1.5"
                    :class="{ 'opacity-50': !item.ai_enabled }"
                    name="chrome"
                    width="14"
                    height="14"
                  />
                  <span class="text-xs">Google</span>
                </ElButton>
              </div>
            </ElFormItem>
          </div>
          <div
            v-if="subscription_list.length < 5"
            class="flex-none w-[200px] bg-[#F9F9FC] flex items-center justify-center p-4 box-border cursor-pointer hover:opacity-70 transition-opacity"
            @click="onAdd"
          >
            <ElIcon size="16" color="#4F5052">
              <Plus />
            </ElIcon>
          </div>
        </ElForm>
      </div>
    </div>
    <div v-if="!loading" class="sticky bottom-0 left-0 right-0 z-10 bg-[#F6F7F8] pb-8 box-border">
      <div class="bg-white px-6 pt-2 pb-6">
        <ElDivider />
        <ElButton v-debounce class="w-[96px]" type="primary" size="large" @click="handleSave">{{
          $t('action_save')
        }}</ElButton>
      </div>
    </div>
  </Layout>

  <el-dialog
    v-model="transfer_dialog_visible"
    class="el-dialog--center"
    :title="$t('subscription.transfer_title')"
    width="480px"
  >
    <p class="text-sm text-[#4F5052]">{{ $t('subscription.transfer_desc') }}</p>
    <template v-for="item in deleted_subscription_list" :key="item.id">
      <div class="flex items-center gap-2 mt-4">
        <div class="flex-1">
          <SelectPlus
            v-model="item.group_id"
            size="large"
            :options="formatSelectOptions(deleted_subscription_list)"
            disabled
          />
        </div>
        <div class="text-sm text-[#4F5052] mx-3">{{ $t('subscription.transfer_to') }}</div>
        <div class="flex-1">
          <SelectPlus v-model="item.target_group_id" size="large" :options="formatSelectOptions(subscription_list)" />
        </div>
      </div>
    </template>

    <div class="h-5"></div>
    <template #footer>
      <ElButton v-debounce size="large" type="primary" @click="handleSubmit">{{ $t('action_confirm') }}</ElButton>
      <ElButton v-debounce size="large" @click="handleCancel">{{ $t('action_cancel') }}</ElButton>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { Delete, Plus } from '@element-plus/icons-vue'
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { generateInputRules } from '@/utils/form-rule'
import { subscriptionApi } from '@/api/modules/subscription'
import { img_host } from '@/utils/config'
import { createNewSubscriptionItem, transformSubscriptionItemForSave } from './utils'
import type { SubscriptionItem, UnitOption } from '@/types/subscription'

// 响应式数据
const form_ref = ref<FormInstance>()
const loading = ref(false)
const subscription_list = ref<SubscriptionItem[]>([])
const deleted_subscription_list = ref<SubscriptionItem[]>([])

const transfer_dialog_visible = ref(false)

const _original_subscription_list = ref<SubscriptionItem[]>([])

// 货币选项
const unit_options = ref<UnitOption[]>([
  { value: 'CNY', label: 'CNY' },
  { value: 'USD', label: 'USD' },
])

// 刷新数据
const refresh = async (): Promise<void> => {
  loading.value = true
  try {
    subscription_list.value = await subscriptionApi.list()
    _original_subscription_list.value = subscription_list.value.map((item: SubscriptionItem) => ({ ...item }))
  } catch (error) {
    console.error('Failed to fetch subscription list:', error)
  } finally {
    loading.value = false
  }
}

// 添加新订阅项
const onAdd = (): void => {
  const template = subscription_list.value[0]
  if (!template) return

  const nextIndex = subscription_list.value.length + 1
  const newSubscriptionData = createNewSubscriptionItem(template, nextIndex)
  subscription_list.value.push(newSubscriptionData)
}

// 删除订阅项
const handleRemove = (params: { data: SubscriptionItem; index: number }) => {
  const { data, index } = params
  if (data.is_default) return

  data.delete = true
  if (data.group_id) {
    deleted_subscription_list.value.push({
      ...data,
      target_group_id: subscription_list.value[0].group_id,
    })
  }

  subscription_list.value.splice(index, 1)
}

// 图标改变
const onIconChange = (value: string, data: SubscriptionItem): void => {
  data.logo_url = `${img_host}/subscription/vip-${value}.png`
}

const formatSelectOptions = (list: SubscriptionItem[]): GroupOptionItem[] => {
  return list.map((item: SubscriptionItem) => ({
    value: item.group_id,
    label: item.group_name,
    icon: item.logo_url,
  }))
}

const handleSubmit = async () => {
  try {
    const items = subscription_list.value.map((item: SubscriptionItem, index: number) =>
      transformSubscriptionItemForSave(item, index, subscription_list.value.length)
    )

    if (deleted_subscription_list.value.length) {
      items.push(...deleted_subscription_list.value)
    }
    transfer_dialog_visible.value = false
    await subscriptionApi.save({ data: { items } })
    ElMessage.success(window.$t('action_save_success'))
    deleted_subscription_list.value = []
    await refresh()
  } catch (error) {
    console.error('Save failed:', error)
    ElMessage.error(window.$t('action_save_failed'))
  }
}

// 处理保存
const handleSave = (): void => {
  if (!form_ref.value) return

  form_ref.value.validate(async (valid: boolean) => {
    if (!valid) return

    if (deleted_subscription_list.value.length) {
      transfer_dialog_visible.value = true
      return
    }
    handleSubmit()
  })
}

// 数字输入处理
const numberInputKeydownHandler = (event: KeyboardEvent): void => {
  // 允许的按键：数字、退格、删除、Tab、Enter、左右箭头、小数点
  const allowedKeys = ['Backspace', 'Delete', 'Tab', 'Enter', 'ArrowLeft', 'ArrowRight', '.']

  if (
    !allowedKeys.includes(event.key) &&
    !(event.key >= '0' && event.key <= '9') &&
    !(event.ctrlKey || event.metaKey)
  ) {
    event.preventDefault()
  }
}

const handleCancel = (): void => {
  transfer_dialog_visible.value = false
  deleted_subscription_list.value = []
  subscription_list.value = _original_subscription_list.value.map((item: SubscriptionItem) => ({ ...item }))
}

// 生命周期钩子
onMounted(() => {
  refresh()
})

onUnmounted(() => {})
</script>

<style scoped>
::v-deep .charge-point-input .el-input__inner {
  text-align: left !important;
}

:deep(.el-divider--horizontal) {
  margin: 20px 0;
}

:deep(.el-form-item__label) {
  width: 100%;
  padding-right: 0;
}
</style>
