<template>
  <el-aside v-show="!siderHidden" width="232px">
    <div class="flex flex-col flex-1 h-full">
      <div class="flex-none pl-7 pt-8 pb-5">
        <!-- #ifndef KM -->
        <img class="w-[136px] object-contain" :src="$getRealPath({ url: '/images/logo_2.png' })" alt="" />
        <!-- #endif -->

        <!-- #ifdef KM -->
        <img class="w-[136px] object-contain" :src="$getRealPath({ url: '/images/km-logo.png' })" alt="" />
        <!-- #endif -->
      </div>
      <div class="flex-1 border-t overflow-y-auto scrollbar--none">
        <el-menu
          router
          :default-active="activeName"
          :default-openeds="[openedMenu, '/agent']"
          active-text-color="#2563EB"
          background-color="#F6F7F8"
          class="border-none mx-4"
          text-color="#4C4D4E"
          style="--el-menu-item-height: 54px; --el-menu-item-font-size: 16px; --el-menu-hover-bg-color: #eeeff0"
        >
          <el-menu-item index="/index">
            <svg-icon name="home" width="18px" class="mr-2" />
            <span>{{ $t('module.homepage') }}</span>
          </el-menu-item>

          <!-- #ifndef KM -->
          <el-sub-menu index="/agent">
            <template #title>
              <svg-icon name="app" width="18px" class="mr-2" />
              <span>{{ $t('module.app_management') }}</span>
            </template>
            <el-menu-item-group>
              <el-menu-item index="/agent">
                {{ $t('module.agent') }}
              </el-menu-item>
              <el-menu-item index="/prompt">
                {{ $t('prompt.title') }}
              </el-menu-item>
              <el-menu-item index="/toolbox">
                {{ $t('module.ai_toolbox') }}
              </el-menu-item>
            </el-menu-item-group>
          </el-sub-menu>
          <!-- #endif -->

          <!-- #ifdef KM -->
          <el-sub-menu index="/space">
            <template #title>
              <svg-icon name="app" width="18px" class="mr-2" />
              <span>{{ $t('module.knowledge') }}</span>
            </template>
            <el-menu-item-group>
              <el-menu-item index="/space">
                {{ $t('module.space') }}
              </el-menu-item>
            </el-menu-item-group>
          </el-sub-menu>
          <el-sub-menu index="/chunk">
            <template #title>
              <svg-icon name="app" width="18px" class="mr-2" />
              <span>{{ $t('module.library') }}</span>
            </template>
            <el-menu-item-group>
              <el-menu-item index="/chunk">
                {{ $t('module.chunk_setting') }}
              </el-menu-item>
            </el-menu-item-group>

            <el-menu-item-group>
              <el-menu-item index="/model"> 模型设置 </el-menu-item>
            </el-menu-item-group>
          </el-sub-menu>
          <!-- #endif -->

          <el-sub-menu index="3">
            <template #title>
              <svg-icon name="operate" width="18px" class="mr-2" />
              <span>{{ $t('module.operation_management') }}</span>
            </template>
            <el-menu-item-group>
              <!-- #ifndef KM -->
              <el-menu-item v-if="enterprise_info.is_independent || enterprise_info.is_industry" index="/order">
                {{ $t('module.operation_order') }}
              </el-menu-item>
              <el-menu-item v-if="enterprise_info.is_independent || enterprise_info.is_industry" index="/user/register">
                {{ $t('register_user.title') }}
              </el-menu-item>

              <el-menu-item v-if="!isOpLocalEnv && !enterprise_info.is_independent" index="/user/internal">
                {{ $t('internal_user.title') }}
              </el-menu-item>

              <el-menu-item index="/subscription">
                {{ $t('module.subscription') }}
              </el-menu-item>
              <!-- #endif -->

              <!-- #ifdef KM -->
              <el-menu-item index="/user/internal">
                {{ $t('internal_user.title') }}
              </el-menu-item>
              <!-- #endif -->

              <el-menu-item index="/user/admin">
                {{ $t('admin_user.title') }}
              </el-menu-item>
              <!-- <el-menu-item index="/visit">
								{{ $t('module.operation_visit') }}
							</el-menu-item> -->
            </el-menu-item-group>
          </el-sub-menu>

          <!-- #ifndef KM -->
          <el-sub-menu index="5">
            <template #title>
              <svg-icon name="decoration" width="18px" class="mr-2" />
              <span>{{ $t('module.site_decoration') }}</span>
            </template>
            <el-menu-item-group>
              <el-menu-item index="/template-style">
                {{ $t('module.template_style') }}
              </el-menu-item>
              <el-menu-item index="/banner">
                {{ $t('module.banner_diagram') }}
              </el-menu-item>
              <el-menu-item index="/navigation">
                {{ $t('navigation.title') }}
              </el-menu-item>
            </el-menu-item-group>
          </el-sub-menu>
          <!-- #endif -->
          <el-sub-menu index="6">
            <template #title>
              <svg-icon name="setting" width="18px" class="mr-2" />
              <span>{{ $t('module.site_config') }}</span>
            </template>
            <el-menu-item-group>
              <el-menu-item index="/info">
                {{ $t('module.website_info') }}
              </el-menu-item>
              <el-menu-item v-if="!isOpLocalEnv && !enterprise_info.is_independent" index="/sso">
                {{ $t('sso.title') }}
              </el-menu-item>

              <!-- #ifdef KM -->
              <el-menu-item index="/sso">
                {{ $t('sso.title') }}
              </el-menu-item>
              <!-- #endif -->

              <el-menu-item index="/platform">
                {{ $t('module.platform') }}
              </el-menu-item>
              <!-- #ifndef KM -->
              <el-menu-item index="/payment">
                {{ $t('module.payment') }}
              </el-menu-item>
              <el-menu-item v-if="!isOpLocalEnv" index="/domain">
                {{ $t('module.domain') }}
              </el-menu-item>
              <el-menu-item index="/statistics">
                {{ $t('module.statistics') }}
              </el-menu-item>
              <el-menu-item v-if="isOpLocalEnv" index="/SMTP">
                {{ $t('module.SMTP') }}
              </el-menu-item>
              <!-- #endif -->
              <el-menu-item v-if="!isWorkEnv" index="/system-log">
                {{ $t('module.system_log') }}
              </el-menu-item>
            </el-menu-item-group>
          </el-sub-menu>
        </el-menu>
      </div>

      <div class="flex-none flex flex-col pt-2 px-4 pb-6 border-t">
        <div
          class="h-14 flex items-center gap-2 px-4 rounded-lg cursor-pointer hover:bg-[#EEEFF0]"
          @click="handleFunctionUpdate"
        >
          <div class="w-6 h-6 flex-center">
            <svg-icon name="generate" width="18" />
          </div>
          <div class="flex-1 text-[#4C4D4E] text-sm">
            {{ $t('function_update') }}
          </div>
          <el-icon size="12px" color="#707172">
            <ArrowRight />
          </el-icon>
        </div>
        <!-- <div
					class="h-14 flex items-center gap-2 px-4 rounded-lg cursor-pointer hover:bg-[#EEEFF0] opacity-60 pointer-events-none">
					<div class="w-6 h-6 flex-center">
						<svg-icon name="remind" width="18" />
					</div>
					<div class="flex-1 text-[#4C4D4E] text-sm">
						{{ $t('notification_center') }}
					</div>
					<el-icon size="12px" color="#707172">
						<ArrowRight />
					</el-icon>
				</div> -->
        <div
          class="h-14 flex items-center gap-2 px-4 rounded-lg cursor-pointer hover:bg-[#EEEFF0]"
          @click.stop="onExit"
        >
          <div class="w-6 h-6 flex-center">
            <svg-icon name="exit" width="18" />
          </div>
          <div class="flex-1 text-[#4C4D4E] text-sm">
            {{ $t('action_exit') }}
          </div>
        </div>
        <div class="flex items-center justify-start gap-2 px-5 mt-4 text-[#1D1E1F]" style="font-size: 10px">
          <span>POWERED BY</span>
          <!-- #ifndef KM -->
          <img class="w-[58px] object-contain" :src="$getRealPath({ url: '/images/logo_2.png' })" alt="" />
          <!-- #endif -->
          <!-- #ifdef KM -->
          <img class="w-[58px] object-contain" :src="$getRealPath({ url: '/images/km-logo.png' })" alt="" />
          <!-- #endif -->
        </div>
      </div>
    </div>
  </el-aside>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watchEffect } from 'vue'
import { useRoute } from 'vue-router'
import { ArrowRight } from '@element-plus/icons-vue'
import { useEnterpriseStore, useUserStore } from '@/stores'
import { useEnv } from '@/hooks/useEnv'

withDefaults(
  defineProps<{
    siderHidden?: boolean
  }>(),
  {
    siderHidden: false,
  }
)

const { isOpLocalEnv, isWorkEnv } = useEnv()
const route = useRoute()
const user_store = useUserStore()
const enterprise_store = useEnterpriseStore()
const activeName = ref('')
const openedMenu = ref('') // 新增：存储当前展开的子菜单

const enterprise_info = computed(() => enterprise_store.info)

watchEffect(() => {
  const paths = route.path.match(/\/[^/]+/g) || ['']
  activeName.value = paths[0]
  openedMenu.value = paths[0]
})

onMounted(() => {})

const onExit = () => {
  user_store.logoff({ show_confirm: true, back_to_login: true })
}
const handleFunctionUpdate = () => {
  window.open('https://doc.53ai.com/%E5%85%A5%E9%97%A8/%E4%BA%A7%E5%93%81%E8%B7%AF%E7%BA%BF%E5%9B%BE.html', '_blank')
}
</script>

<style scoped>
:deep(.el-menu-item-group__title:empty) {
  display: none;
}

/* :deep(.el-menu-item)  {
	margin-bottom: 6px;
} */
:deep(.el-menu-item),
:deep(.el-sub-menu__title) {
  padding: 0 16px !important;
}

:deep(.el-menu-item:hover),
:deep(.el-sub-menu__title:hover) {
  border-radius: 8px;
}

:deep(.el-sub-menu .el-menu-item) {
  --el-menu-sub-item-height: 40px;
  --el-menu-base-level-padding: 0px;
  --el-menu-level-padding: 16px;

  font-size: 14px;
}

:deep(.el-menu-item-group) {
  padding-left: 30px;
}

:deep(.el-sub-menu .el-menu-item:hover) {
  background-color: #eeeff0;
  border-radius: 8px;
}
</style>
