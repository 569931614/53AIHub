<script setup lang="ts">
import { ref, computed, reactive } from 'vue'
import { useConversationStore } from '@/stores/modules/conversation'

const emits = defineEmits<{
  (e: 'new'): void
}>()

const convStore = useConversationStore()
const currentConv = computed(() => convStore.currentConversation)

const visible = ref(false)

const state = reactive({
  editVisible: false
})
const convForm = reactive({
  conversation_id: 0,
  title: ''
})

const handleCreate = () => {
  emits('new')
  visible.value = false
}

const handleEditConv = async () => {
  await convStore.editConversation(convForm)
  state.editVisible = false
}
const delConversation = async (conv: Conversation.Info) => {
  await ElMessageBox.confirm(
    window.$t('chat.conversation_confirm_delete'),
    window.$t('action.del'),
    {
      confirmButtonText: window.$t('action.del'),
      cancelButtonText: window.$t('action.cancel'),
      type: 'warning'
    }
  )
  convStore.delConversation(conv)
}
const handleCommandConv = (event: string, conv: Conversation.Info) => {
  if (event === 'del') {
    delConversation(conv)
  } else if (event === 'edit') {
    convForm.conversation_id = conv.conversation_id
    convForm.title = conv.title
    state.editVisible = true
  }
}

defineExpose({
  open() {
    visible.value = true
  }
})
</script>

<template>
  <el-drawer v-model="visible" :title="$t('chat.history')" size="300px">
    <el-button class="w-full border-none -mt-4" type="primary" size="large" plain @click="handleCreate">
      + {{ $t('chat.new_conversation') }}
    </el-button>
    <div class="flex flex-col gap-2 mt-4">
      <template v-for="item in convStore.conversations" :key="item.conversation_id">
        <div class="group p-3 rounded cursor-pointer hover:bg-[#F5F6FA]"
          :class="[currentConv.conversation_id === item.conversation_id ? 'bg-[#F5F6FA]' : '']"
          @click="convStore.setCurrentState(item.agent_id, item.conversation_id)">
          <div class="flex items-center gap-2">
            <div class="flex-1 text-sm text-primary truncate">{{ item.title }}</div>
            <el-dropdown @command="handleCommandConv($event, item)">
              <div class="size-7 flex-center cursor-pointer invisible group-hover:visible">
                <svg-icon name="more-h" />
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="edit">
                    <svg-icon name="edit" class="mr-1"></svg-icon>
                    {{ $t('action.rename') }}
                  </el-dropdown-item>
                  <el-dropdown-item command="del">
                    <span class="text-[#FA5151] flex-center">
                      <svg-icon name="del" class="mr-1"></svg-icon>
                      {{ $t('action.del') }}
                    </span>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
          <div class="mt-2 text-xs text-secondary">{{ item.created_at }}</div>
        </div>
      </template>
    </div>
  </el-drawer>

  <el-dialog v-model="state.editVisible" :title="$t('chat.edit_conversation')" width="480">
    <el-input size="large" v-model="convForm.title" v-trim :placeholder="$t('form.input_placeholder')" :maxlength="20"
      show-word-limit></el-input>

    <template #footer>
      <el-button @click="state.editVisible = false">{{ $t('action.cancel') }}</el-button>
      <el-button v-debounce type="primary" :disabled="!convForm.title.trim()" @click="handleEditConv">{{
        $t('action.confirm') }}</el-button>
    </template>
  </el-dialog>
</template>

<style scoped></style>
