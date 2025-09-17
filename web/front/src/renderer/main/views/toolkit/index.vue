<template>
  <div class="">
    <MainHeader sticky>
      <template #before_suffix>
        <div class="text-base text-primary font-bold line-clamp-1 max-md:flex-1 max-md:text-center">
          {{ navigationStore.toolkitNavigation.name || $t('module.toolbox') }}
        </div>
      </template>
    </MainHeader>
    <div class="w-11/12 lg:w-4/5 mx-auto min-h-[80vh]">
      <ToolkitGroupList :sticky-offset="70" />
    </div>
    <MFooter />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import MainHeader from '@/layout/header.vue'
import MFooter from '@/layout/m-footer.vue'
import ToolkitGroupList from './components/group-list.vue'

import { useNavigationStore } from '@/stores/modules/navigation'

const navigationStore = useNavigationStore()

onMounted(() => {})
</script>

<style scoped></style>
