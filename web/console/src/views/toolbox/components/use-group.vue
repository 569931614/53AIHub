<template>
  <div class="text-base text-[#1D1E1F] font-medium mt-10 mb-4">
    {{ $t('permission_setting') }}
  </div>
  <div>
    <ElFormItem
      :hidden="!(enterprise.info.is_independent || enterprise.info.is_industry)"
      :label="$t('register_user.title')"
    >
      <GroupSelect
        v-model="props.subscriptionGroup"
        type="checkbox"
        :group-type="GROUP_TYPE.USER"
        multiple
        :default-all="!props.editable"
        @change="handleSubscriptionChange"
      />
    </ElFormItem>
    <ElFormItem
      :hidden="!(enterprise.info.is_enterprise || enterprise.info.is_industry)"
      :label="$t('internal_user.title')"
      prop="user_group_ids"
    >
      <GroupSelect
        v-model="props.userGroup"
        type="picker"
        :group-type="GROUP_TYPE.INTERNAL_USER"
        multiple
        @change="handleUserGroupChange"
      />
    </ElFormItem>
  </div>
</template>

<script setup lang="ts">
import { GROUP_TYPE } from '@/constants/group';
import { useEnterpriseStore } from '@/stores/modules/enterprise';

const enterprise = useEnterpriseStore()

const props = withDefaults(
  defineProps<{
    userGroup: number[]
    subscriptionGroup: number[]
    editable: boolean
  }>(),
  {
    userGroup: () => [],
    subscriptionGroup: () => [],
    editable: true,
  }
)

const emit = defineEmits(['change'])

const handleUserGroupChange = (data: number[]) => {
  emit('change', { groupType: GROUP_TYPE.INTERNAL_USER, data })
}

const handleSubscriptionChange = (data: number[]) => {
  emit('change', { groupType: GROUP_TYPE.USER, data })
}
</script>

<style scoped></style>
