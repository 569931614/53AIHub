<template>
  <div ref="scrollRef" class="relative h-full flex flex-col overflow-y-auto overflow-x-hidden">
    <MainHeader type="homepage" :need-login="false" :sider-button="false" :back="false" sticky :main-class="['w-11/12 lg:w-4/5']">
      <template #before_suffix>
        <a href="/" class="flex flex-none items-center gap-2 overflow-hidden">
          <img :alt="enterpriseStore.display_name" :title="enterpriseStore.display_name" class="min-w-11 h-11 rounded" :src="enterpriseStore.logo" />
          <span class="text-2xl font-semibold nav-text truncate">
            {{ enterpriseStore.display_name }}
          </span>
        </a>
        <div class="flex-1 w-0 menu overflow-hidden">
          <ElSkeleton :loading="navigationStore.loading && false" animated>
            <template #template>
              <ElSkeletonItem v-for="i in 4" :key="i" variant="text" class="ml-4 mt-2 !w-[82px] !h-[42px]" />
            </template>
            <template #default>
              <ElMenu router mode="horizontal" :default-active="activeMenuItem">
                <ElMenuItem
                  v-for="item in navigationStore.navigations"
                  :key="item.navigation_id"
                  class="relative !cursor-pointer !opacity-100 hover-text-theme-important"
                  :index="item.menu_path"
                  :disabled="item.target === NAVIGATION_TARGET.BLANK || item.type === NAVIGATION_TYPE.EXTERNAL"
                >
                  {{ item.name }}
                  <div
                    v-if="item.target === NAVIGATION_TARGET.BLANK || item.type === NAVIGATION_TYPE.EXTERNAL"
                    class="absolute top-0 left-0 w-full h-full bg-transparent z-[10]"
                    @click="handleNavigationClick(item)"
                  />
                </ElMenuItem>
              </ElMenu>
            </template>
          </ElSkeleton>
        </div>
      </template>
    </MainHeader>
    <div
      v-if="
        !['Index', 'HomeIndex', 'HomePromptDetail', 'HomeAgentDetail', 'HomeChat'].includes(route.name) &&
        enterpriseStore.banner_info.url_list &&
        !!enterpriseStore.banner_info.url_list.length
      "
      class="w-full flex-none"
    >
      <ElCarousel
        class="!w-full"
        :arrow="enterpriseStore.banner_info.url_list.length > 1 ? 'always' : 'never'"
        :indicator-position="enterpriseStore.banner_info.url_list.length > 1 ? 'outside' : 'none'"
        :interval="enterpriseStore.banner_info.interval ? parseInt(enterpriseStore.banner_info.interval * 1000) : 5000"
      >
        <ElCarouselItem v-for="url in enterpriseStore.banner_info.url_list || []" :key="url" class="w-full">
          <img :src="url" class="object-cover w-full h-full" />
        </ElCarouselItem>
      </ElCarousel>
    </div>
    <div class="flex-1" :class="{ 'overflow-hidden': ['HomeChat'].includes(route.name) }">
      <RouterView v-slot="{ Component, route }" v-loading="is_redirect">
        <component :is="Component" :key="route.path" />
      </RouterView>
    </div>
    <!-- 底部布局 -->
    <mFooter />
  </div>
</template>

<script setup lang="ts">
import { nextTick, onMounted, ref, watch, provide } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import MainHeader from '@/layout/header.vue'
import { useEnterpriseStore } from '@/stores/modules/enterprise'
import { useNavigationStore } from '@/stores/modules/navigation'
import { NAVIGATION_TYPE, NAVIGATION_TARGET } from '@/constants/navigation'
import mFooter from '@/layout/m-footer.vue'

// const mainRef = inject<Ref<HTMLElement>>('mainRef')

const enterpriseStore = useEnterpriseStore()
const navigationStore = useNavigationStore()
const route = useRoute()
const router = useRouter()

const activeMenuItem = ref()

const scrollRef = ref<HTMLElement>()
provide('mainScrollRef', scrollRef)

const fetchNavigationData = async () => {
  await navigationStore.fetchNavigations()
  const customNavigations = navigationStore.navigations.filter((item) => item.type === NAVIGATION_TYPE.CUSTOM)
  const indexRoute = router.getRoutes().find((item) => item.name === 'Index')
  customNavigations.forEach((item) => {
    if (indexRoute) {
      indexRoute.children.push({
        path: item.menu_path,
        name: `Home${item.name}`,
        component: () => import(`@/views/custom/index.vue`)
      })
    }
  })
  router.addRoute(indexRoute)
}
const handleNavigationClick = (data: any) => {
  if (data.type === NAVIGATION_TYPE.EXTERNAL) {
    if (data.target === NAVIGATION_TARGET.BLANK) {
      window.open(data.url, '_blank')
    } else {
      window.location.href = data.url
    }
  } else if (data.target === NAVIGATION_TARGET.BLANK) {
    window.open(data.url, '_blank')
  }
}
const is_redirect = ref(false)
const handleRedirect = async () => {
  const from_home = route.query.from_home as string
  let redirect = route.query.redirect as string
  await fetchNavigationData()
  if (navigationStore.navigations.length) {
    if (+from_home || redirect) {
      is_redirect.value = true
      await nextTick()
      if (+from_home && !redirect) {
        const defaultNavigation = navigationStore.navigations.filter((item) => item.type !== NAVIGATION_TYPE.EXTERNAL)[0]
        if (defaultNavigation && defaultNavigation.menu_path !== route.path) redirect = defaultNavigation.menu_path
        else redirect = '/index'
      }
      if (redirect) await router.replace({ path: redirect as string })
      is_redirect.value = false
    }
  }
}
const setMeta = ({ key = '', value = '' }) => {
  const meta = document.querySelector(`meta[name="${key}"]`)
  if (meta) {
    meta.setAttribute('content', value)
  } else {
    const meta = document.createElement('meta')
    meta.setAttribute('name', key)
    meta.setAttribute('content', value)
    document.head.appendChild(meta)
  }
}

onMounted(async () => {
  // mainRef.value?.scrollTo(0, 0)
  // handleRedirect()
})

watch(
  () => route.query,
  () => {
    handleRedirect()
  },
  {
    immediate: true,
    deep: true
  }
)
watch(
  () => route.path,
  (path) => {
    const currentNavigation = navigationStore.navigations.find((item) => item.menu_path === path) || {}
    const config = currentNavigation.config || {}
    if (config.seo_title) setMeta({ key: 'title', value: config.seo_title })
    if (config.seo_keywords) setMeta({ key: 'keywords', value: config.seo_keywords })
    if (config.seo_description) setMeta({ key: 'description', value: config.seo_description })
    if (path.includes('/chat') || path.includes('/agent')) {
      activeMenuItem.value = navigationStore.navigations[1].menu_path
    } else if (path.includes('/prompt')) {
      activeMenuItem.value = navigationStore.navigations[2].menu_path
    } else if (path.includes('/toolkit')) {
      activeMenuItem.value = navigationStore.navigations[3].menu_path
    } else {
      activeMenuItem.value = navigationStore.navigations[0].menu_path
    }
  },
  {
    immediate: true
  }
)
</script>

<style scoped>
.router-link-active {
  @apply ring-2 ring-blue-500;
}

::v-deep(.el-menu) {
  background-color: transparent;
  border-bottom: none;
}

::v-deep(.el-menu-item),
::v-deep(.el-menu-item.is-active) {
  background-color: transparent !important;
  border-bottom: none;
}

::v-deep(.el-carousel__container) {
  height: 256px;
}

@media (min-width: 768px) {
  ::v-deep(.el-carousel__container) {
    height: 300px;
  }
}

@media (min-width: 1024px) {
  ::v-deep(.el-carousel__container) {
    height: 385px;
  }
}
</style>
