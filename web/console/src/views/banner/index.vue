<template>
  <Layout class="px-[60px] py-8">
    <Header :title="$t('module.banner_diagram')"></Header>
    <div v-loading="loading" class="mt-5 flex-1 flex flex-col gap-4 bg-white py-6 px-8 box-border">
      <ElForm ref="formRef" class="flex-1 max-h-[calc(100vh-288px)] overflow-auto" :model="form" label-position="top">
        <ElFormItem :label="$t('banner.upload_image')">
          <div class="text-[#939499] text-xs w-full">{{ $t('banner.upload_image_tip') }}</div>
          <ul v-if="form.url_list.length > 0" class="mt-4 w-full flex flex-col gap-4">
            <li v-for="(url, index) in form.url_list" :key="url" class="w-full relative">
              <UploadImage
                v-model="form.url_list[index]"
                class="!w-full !h-[14vw]"
                :cropper-disabled="true"
                @confirm="onUploadImageConfirm($event, index)"
              >
                <template #mask-text>
                  <ElIcon class="cursor-pointer" size="24" color="#fff">
                    <EditPen />
                  </ElIcon>
                  <ElIcon class="cursor-pointer" size="24" color="#fff" @click="onImageDelete(index)">
                    <Delete />
                  </ElIcon>
                </template>
              </UploadImage>
            </li>
          </ul>
          <UploadImage
            class="w-auto h-auto"
            :cropper-disabled="true"
            :disabled="addDisabled"
            @confirm="onUploadImageConfirm"
          >
            <ElButton class="mt-4 !border-none !outline-none" type="primary" plain size="large" :disabled="addDisabled">
              + {{ $t('action_add') }}({{ form.url_list.length }}/{{ UPLOAD_COUNT_LIMIT }})
            </ElButton>
          </UploadImage>
        </ElFormItem>
        <ElFormItem :label="$t('banner.interval')">
          <div class="text-[#939499] text-xs w-full">{{ $t('banner.interval_tip') }}</div>
          <ElInputNumber
            v-model="form.interval"
            class="mt-4 !w-[300px] el-input-number--left"
            size="large"
            :controls="false"
            :min="1"
            :max="1000"
          >
            <template #suffix>
              <span>{{ $t('second') }}</span>
            </template>
          </ElInputNumber>
        </ElFormItem>
        <div class="w-full h-8"></div>
      </ElForm>
      <ElDivider class="!mt-3" />
      <ElButton v-debounce class="h-[36px] w-[96px]" type="primary" size="large" @click="handleSave"
        >{{ $t('action_save') }}
      </ElButton>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Delete, EditPen } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import UploadImage from '@/components/Upload/image.vue'

import { getDefaultBanner, transformBanner } from '@/api/modules/banner/transform'
import { bannerApi } from '@/api/modules/banner/index'
import type { Banner } from '@/api/modules/banner/types'
import { BANNER_CONFIG } from '@/constants/banner'

const UPLOAD_COUNT_LIMIT = BANNER_CONFIG.MAX_IMAGES

const formRef = ref<FormInstance>()

const loading = ref(true)
const submitting = ref(false)
const form = reactive<Banner>(getDefaultBanner())
const addDisabled = computed(() => {
  return form.url_list.length >= UPLOAD_COUNT_LIMIT
})

const fetchBannerData = async () => {
  loading.value = true
  const data = await bannerApi
    .get()
    .then(transformBanner)
    .finally(() => {
      loading.value = false
    })
  form.url_list = data.url_list || []
  form.interval = data.interval || BANNER_CONFIG.DEFAULT_INTERVAL
}

const handleSave = async () => {
  const valid = await formRef.value?.validate()
  if (!valid) return
  submitting.value = true
  await bannerApi
    .save({
      ...form,
    })
    .finally(() => {
      submitting.value = false
    })
  ElMessage.success(window.$t('action_save_success'))
  fetchBannerData()
}

const onUploadImageConfirm = ({ url = '' }, index?: number) => {
  if (index !== undefined && form.url_list[index]) form.url_list[index] = url
  else form.url_list.push(url)
}
const onImageDelete = (index: number) => {
  form.url_list.splice(index, 1)
}

onMounted(() => {
  fetchBannerData()
})
</script>

<style scoped lang="scss"></style>
