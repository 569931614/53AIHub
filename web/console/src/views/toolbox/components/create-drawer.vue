<template>
  <ElDrawer
    v-model="visible"
    :title="editable ? $t('action_edit') : $t('action_add')"
    size="840px"
    destroy-on-close
    append-to-body
    :close-on-click-modal="false"
  >
    <ElForm ref="formRef" :model="formData" label-position="top">
      <div class="text-base text-[#1D1E1F] font-medium mb-4">
        {{ $t('basic_info') }}
      </div>
      <div class="flex items-center gap-4">
        <ElFormItem
          class="flex-1"
          :label="$t('name')"
          prop="name"
          :rules="generateInputRules({ message: 'form_input_placeholder' })"
        >
          <ElInput
            v-model="formData.name"
            size="large"
            show-word-limit
            maxlength="20"
            :placeholder="$t('form_input_placeholder')"
          />
        </ElFormItem>
        <ElFormItem
          :label="$t('group')"
          class="flex-1"
          prop="group_id"
          :rules="generateInputRules({ message: 'form_select_placeholder' })"
        >
          <ElSelect v-model="formData.group_id" size="large">
            <ElOption
              v-for="item in showGroupOptions"
              :key="item.group_id"
              :value="item.group_id"
              :label="$t(item.group_name)"
            />
          </ElSelect>
        </ElFormItem>
      </div>
      <ElFormItem
        label="URL"
        prop="url"
        :rules="generateInputRules({ message: 'form_input_placeholder', validator: ['url'] })"
      >
        <ElInput v-model="formData.url" size="large" placeholder="http://" />
      </ElFormItem>

      <div class="flex items-center justify-between gap-2 mb-2">
        <div class="text-sm text-[#4F5052]">{{ $t('shared_account') }}</div>
        <el-button link size="large" class="!text-blue-500" @click="handleAdd"> +{{ $t('action_add') }} </el-button>
      </div>
      <ElFormItem>
        <TablePlus
          header-row-class-name="rounded overflow-hidden"
          header-cell-class-name="!bg-[#F6F7F8] !h-[60px] !border-none"
          :data="accountList"
          :pagination="false"
        >
          <ElTableColumn :label="$t('account')" min-width="140" show-overflow-tooltip>
            <template #default="{ row = {} }">
              {{ row.account }}
            </template>
          </ElTableColumn>
          <ElTableColumn :label="$t('password')" min-width="140" show-overflow-tooltip>
            <template #default="{ row = {} }">
              {{ row.password }}
            </template>
          </ElTableColumn>
          <ElTableColumn :label="$t('remark')" min-width="140" show-overflow-tooltip>
            <template #default="{ row = {} }">
              {{ row.remark || '--' }}
            </template>
          </ElTableColumn>
          <ElTableColumn :label="$t('operation')" width="120" align="left" fixed="right">
            <template #default="{ row }">
              <div class="flex">
                <el-button type="primary" link @click="onEdit(row)">
                  <SvgIcon name="edit" class="text-[#606266]"></SvgIcon>
                </el-button>
                <el-button type="primary" link @click="onDelete(row)">
                  <SvgIcon name="del" class="text-[#606266]"></SvgIcon>
                </el-button>
              </div>
            </template>
          </ElTableColumn>
        </TablePlus>
      </ElFormItem>
      <ElFormItem :label="$t('description')">
        <ElInput
          v-model="formData.description"
          type="textarea"
          :rows="3"
          resize="none"
          show-word-limit
          maxlength="200"
        />
      </ElFormItem>
      <ElFormItem :label="$t('avatar')" prop="logo" :rules="generateInputRules({ message: 'form_upload_placeholder' })">
        <UploadImage v-model="formData.logo" class="w-12 h-12" />
      </ElFormItem>

      <UseGroup
        :user-group="userGroup"
        :editable="editable"
        :subscription-group="subscriptionGroup"
        @change="onChange"
      />
    </ElForm>
    <template #footer>
      <div class="flex border-t pt-5 justify-end w-full">
        <ElButton size="large" :loading="submitting" @click="close">
          {{ $t('action_cancel') }}
        </ElButton>
        <ElButton v-debounce type="primary" size="large" @click="handleSave">
          {{ $t('action_confirm') }}
        </ElButton>
      </div>
    </template>
  </ElDrawer>

  <AddAccount ref="addRef" :account-list="accountList" @success="handleAddSuccess" />
</template>

<script setup lang="ts">
import { computed, inject, reactive, ref, nextTick } from 'vue'
import UploadImage from '@/components/Upload/image.vue'
import AddAccount from './add-account.vue'
import UseGroup from './use-group.vue'

import { generateInputRules } from '@/utils/form-rule'
import { aiLinkApi } from '@/api/modules/ai-link'
import { GROUP_TYPE } from '@/constants/group'

interface accountItem {
  account: string
  password: string
  remark: string
}

const emits = defineEmits<{
  (e: 'success'): any
}>()

const formRef = ref()
const addRef = ref()

const groupOptions = inject('groupOptions', [])

const visible = ref(false)
const editable = ref(false)
const accountEdit = ref(false)
const submitting = ref(false)
const originInfo = ref({})
const userGroup = ref([])
const subscriptionGroup = ref([])
const editingIndex = ref(-1)
const sort = ref(0)

const accountList = ref<accountItem[]>([])

const formData = reactive({
  logo: '',
  name: '',
  url: '',
  description: '',
  group_id: '',
})

const showGroupOptions = computed(() => groupOptions.value.filter(item => +item.group_id > 0))

const reset = () => {
  formData.logo = ''
  formData.name = ''
  formData.url = ''
  formData.description = ''
  formData.group_id = (showGroupOptions.value[0] || {}).group_id || ''
  submitting.value = false
}

const handleAdd = () => {
  accountEdit.value = false
  addRef.value.open()
}

const handleAddSuccess = (accountData: accountItem) => {
  if (accountEdit.value) {
    accountList.value.splice(editingIndex.value, 1, accountData)
  } else {
    accountList.value = [...accountList.value, accountData]
  }
}

const onDelete = async row => {
  await ElMessageBox.confirm($t('form_delete_confirm'), $t('action.delete'))
  const index = accountList.value.findIndex(item => item.account === row.account)
  if (index !== -1) {
    accountList.value.splice(index, 1)
    ElMessage.success($t('action_delete_success'))
  }
}

const onEdit = row => {
  accountEdit.value = true
  editingIndex.value = accountList.value.findIndex(item => item.account === row.account)
  addRef.value.open(row)
}

const onChange = item => {
  if (item.groupType === GROUP_TYPE.USER) {
    subscriptionGroup.value = item.data.value?.length ? item.data.value : []
  } else {
    userGroup.value = item.data.value?.length ? item.data.value : []
  }
}

const open = async ({ data = {} } = {}) => {
  reset()
  await nextTick()
  sort.value = data.sort || 0
  userGroup.value = data.user_group_ids || []
  subscriptionGroup.value = data.user_group_ids || []
  editable.value = !!+data.ai_link_id
  accountList.value = []
  if (editable.value === true) {
    const detail = await aiLinkApi.detail(data.ai_link_id)
    accountList.value = detail.data.shared_account ? JSON.parse(detail.data.shared_account) : []
  }
  formData.logo = data.logo || ''
  formData.name = data.name || ''
  formData.url = data.url || ''
  formData.description = data.description || ''
  formData.group_id = data.group_id || formData.group_id || ''
  originInfo.value = data
  visible.value = true
}
const close = async () => {
  visible.value = false
}
const handleSave = async () => {
  if (submitting.value) return
  const valid = await formRef.value.validate()
  if (!valid) return
  submitting.value = true

  const requestData = {
    ...formData,
    sort: sort.value,
    shared_account: accountList.value?.length ? JSON.stringify(accountList.value) : '',
    subscription_group_ids: subscriptionGroup.value,
    user_group_ids: userGroup.value,
    ai_link_id: originInfo.value.ai_link_id,
  }
  await aiLinkApi
    .save({
      data: requestData,
    })
    .catch(() => {
      submitting.value = false
    })
  emits('success')
  ElMessage.success(window.$t('action_save_success'))
  close()
}

defineExpose({
  open,
  close,
  reset,
})
</script>

<style scoped></style>
