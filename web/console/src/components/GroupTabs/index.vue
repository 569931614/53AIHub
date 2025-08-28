<template>
  <ElSelect
    v-if="type === 'dropdown'"
    ref="selectRef"
    v-model="selectedValue"
    multiple
    collapse-tags
    collapse-tags-tooltip
    :placeholder="$t('all')"
    size="large"
    class="max-w-[160px]"
    @visible-change="handleDropdownVisibleChange"
  >
    <template #prefix>
      <span class="max-w-16 text-sm truncate">
        {{ selectedOptions.length ? selectedOptions[0].group_name || '--' : $t('group') + ':' }}
      </span>
    </template>
    <ElOption v-for="opt in tabOptions" :key="opt.group_id" :label="opt.group_name" :value="opt.group_id" @click.stop />
    <template #empty>
      <div class="text-center text-[#A4AABA] text-sm py-8">
        {{ $t('no_data') }}
      </div>
    </template>
    <template #footer>
      <div class="w-full">
        <div class="w-full flex items-center justify-between gap-5">
          <div class="cursor-pointer text-[#5A6D9E] text-sm hover:opacity-80" @click="groupRef.open">
            {{ $t('group_management') }}
          </div>
          <div>
            <ElButton class="w-12 h-7" @click="handleDropdownCancel">
              {{ $t('action_cancel') }}
            </ElButton>
            <ElButton class="w-12 h-7" type="primary" @click="handleDropdownConfirm">
              {{ $t('action_confirm') }}
            </ElButton>
          </div>
        </div>
      </div>
    </template>
  </ElSelect>
  <ElTabs v-else-if="type === 'tabs-pure'" v-model="activeTab" class="group-tabs" type="card" @tab-change="onTabChange">
    <ElTabPane
      v-for="opt in options"
      :key="opt.group_id"
      :label="opt.group_name"
      :name="opt.group_id"
      :disabled="disabled || opt.disabled"
    />
  </ElTabs>
  <ElTabs v-else v-model="activeTab" class="group-tabs" type="card" @tab-change="onTabChange">
    <ElTabPane
      v-for="opt in tabOptions"
      :key="opt.group_id"
      :label="$t(opt.group_name)"
      :name="opt.group_id"
      :disabled="disabled || opt.disabled"
    />
  </ElTabs>
  <GroupDialog v-if="groupType" ref="groupRef" :group-type="groupType" @change="onGroupChange" />
</template>

<script setup lang="ts">
import { computed, nextTick, ref } from 'vue';
import type { Group } from '@/api/modules/group';
import type { GroupType } from '@/constants/group';

const props = withDefaults(
  defineProps<{
    type?: 'tabs' | 'dropdown' | 'tabs-pure'
    groupType?: GroupType
    modelValue?: string | number | string[] | number[] | null
    options?: Group[]
    disabled?: boolean
  }>(),
  {
    type: 'tabs',
    modelValue: null,
    groupType: 1,
    options: [],
    disabled: false,
  }
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'change', result: { value: string }): void
  (e: 'get-options', result: { options: Group[] }): void
}>()

const groupRef = ref<InstanceType<typeof GroupDialog>>()
const selectRef = ref<InstanceType<typeof ElSelect>>()
const activeTab = ref(props.modelValue)
const tabOptions = ref<Group[]>(props.options)
const selectedValue = ref(props.modelValue)

const selectedOptions = computed(() => {
  return tabOptions.value.filter(item => selectedValue.value.includes(item.group_id))
})

const setModelValue = (value: string | number | string[] | number[]) => {
  value = JSON.parse(JSON.stringify(value))
  emit('update:modelValue', value)
  emit('change', { value })
}
const onTabChange = () => {
  setModelValue(activeTab.value)
}
const onGroupChange = async ({ value = [] } = {}) => {
  if (props.type === 'tabs-pure') return
  tabOptions.value = []
  await nextTick()
  tabOptions.value = [...value]
  if (['tabs'].includes(props.type)) {
    tabOptions.value.unshift({ group_id: '-1', group_name: 'all' })
    if (!tabOptions.value.find(item => item.group_id === activeTab.value)) {
      activeTab.value = tabOptions.value[0].group_id
      onTabChange()
    }
  }
  emit('get-options', { options: tabOptions.value })
}
const handleDropdownCancel = () => {
  selectedValue.value = JSON.parse(JSON.stringify(props.modelValue || []))
  selectRef.value?.blur()
}
const handleDropdownConfirm = async () => {
  setModelValue(selectedValue.value)
  selectRef.value?.blur()
  await nextTick()
}
const handleDropdownVisibleChange = (visible: boolean) => {
  if (!visible) handleDropdownCancel()
}

defineExpose({
  open: () => {
    groupRef.value?.open()
  },
  getData: () => {
    return tabOptions.value
  },
})
</script>

<style lang="scss" scoped>
.group-tabs {
  ::v-deep(.el-tabs__header) {
    margin-bottom: 0;
  }

  ::v-deep(.el-tabs__content) {
    display: none;
  }

  ::v-deep(.el-tabs__header),
  ::v-deep(.el-tabs__nav),
  ::v-deep(.el-tabs__item) {
    border: none;
  }
  ::v-deep(.el-tabs__item) {
    border-radius: 4px;
    background-color: #fff;
    color: #1d1e1f;
  }

  ::v-deep(.el-tabs__item + .el-tabs__item) {
    margin-left: 12px;
  }

  ::v-deep(.el-tabs__item.is-active) {
    background-color: #e0eaff;
    color: #2563eb;
  }
  ::v-deep(.el-tabs__item.is-disabled) {
    opacity: 0.7;
  }
}

::v-deep(.el-select__selected-item:nth-child(1)) .el-tag {
  display: none !important;
}

::v-deep(.el-select__placeholder) {
  color: black !important;
}
</style>
