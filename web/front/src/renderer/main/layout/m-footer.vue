<template>
  <template v-if="enterpriseStore.isSoftStyle">
    <div
      v-if="enterpriseStore.copyright.toLowerCase() !== 'true'"
      class="w-full flex justify-center items-center gap-1.5 text-xs text-placeholder my-4"
    >
      <span>本网站由</span>
      <img :src="$getPublicPath(`/images/53ai-hub-v2.png`)" class="flex-none w-[72px] object-cover" />
      <span>提供技术支持</span>
    </div>
    <div v-else class="h-12"></div>
    <div v-if="fixed" class="h-14 md:hidden"></div>
    <div class="bg-white flex items-center border-t md:hidden" :class="{ 'fixed bottom-0 left-0 right-0 z-[9]': fixed }">
      <template v-for="item in footerList" :key="item.text">
        <div class="flex-1 flex items-center justify-center">
          <router-link
            :to="{ name: item.name }"
            class="h-14 px-2 flex flex-col items-center justify-center gap-1 cursor-pointer"
            :class="[$route.name === item.name ? 'text-theme' : 'text-primary']"
          >
            <svg-icon size="18" :name="item.icon" :stroke="item.stroke"></svg-icon>
            <p class="text-sm leading-none max-w-[100px] truncate">{{ item.title || $t(item.text) }}</p>
          </router-link>
        </div>
      </template>
    </div>
  </template>

  <template v-else>
    <div v-if="enterpriseStore.copyright.toLowerCase() !== 'true'" class="mt-auto relative py-8 md:py-10 lg:py-12 page-footer-bg page-footer-text">
      <!-- <div  class="w-11/12 lg:w-4/5 mx-auto">
        <div class="grid grid-cols-2 md:grid-cols-3 lg:flex gap-6 md:gap-4">
          <div class="flex-1">
            <h4 class="text-lg md:text-xl font-medium mb-3 md:mb-4">产品服务</h4>
            <ul class="space-y-2 text-sm md:text-base text-white text-opacity-50">
              <li><a href="/product/work-ai">工作AI</a></li>
              <li><a href="/product/life-ai">生活AI</a></li>
              <li><a href="/product/enterprise-ai">企业AI</a></li>
              <li><a href="/product/custom-model">大模型定制</a></li>
              <li><a href="/product/magic-workshop">魔法工坊</a></li>
            </ul>
          </div>
          <div class="flex-1">
            <h4 class="text-lg md:text-xl font-medium mb-3 md:mb-4">解决方案</h4>
            <ul class="space-y-2 text-sm md:text-base text-white text-opacity-50">
              <li><a href="/solution/application">【应用场景】解决方案</a></li>
              <li><a href="/solution/function">【功能场景】解决方案</a></li>
              <li><a href="/solution/industry">【行业场景】解决方案</a></li>
            </ul>
          </div>
          <div class="flex-1">
            <h4 class="text-lg md:text-xl font-medium mb-3 md:mb-4">AI知识库</h4>
            <ul class="space-y-2 text-sm md:text-base text-white text-opacity-50">
              <li><a href="/knowledge/technology">前沿技术</a></li>
              <li><a href="/knowledge/application">应用场景</a></li>
              <li><a href="/knowledge/product">产品世界</a></li>
              <li><a href="/knowledge/practice">企业实践</a></li>
            </ul>
          </div>
          <div class="flex-1">
            <h4 class="text-lg md:text-xl font-medium mb-3 md:mb-4">AI百宝箱</h4>
            <ul class="space-y-2 text-sm md:text-base text-white text-opacity-50">
              <li><a href="/download/windows">Windows版本</a></li>
              <li><a href="/download/mac">Mac OS版本</a></li>
              <li><a href="/download/mobile">移动版本</a></li>
            </ul>
          </div>
          <div class="flex-1">
            <h4 class="text-lg md:text-xl font-medium mb-3 md:mb-4">关于我们</h4>
            <ul class="space-y-2 text-sm md:text-base text-white text-opacity-50">
              <li><a href="/about/company">公司介绍</a></li>
              <li><a href="/about/contact">联系合作</a></li>
            </ul>
          </div>
        </div>
      </div>
      <div  class="border-t border-dashed border-white border-opacity-10 my-6 md:my-8 lg:my-10"></div> -->
      <div class="w-11/12 lg:w-4/5 mx-auto flex flex-col md:flex-row items-center">
        <div class="flex-1 w-full md:w-auto mb-6 md:mb-0">
          <div v-if="false" class="flex flex-col md:flex-row md:items-center mb-4">
            <span class="text-base text-white mb-2 md:mb-0">友情链接：</span>
            <div class="flex flex-wrap gap-3 md:space-x-4 text-sm md:text-base text-white text-opacity-50">
              <a href="/legal/statement">法律声明</a>
              <a href="/legal/privacy">隐私政策</a>
              <a href="/community/blackroom">小黑屋</a>
              <a href="/partner/recruit">代理招募</a>
            </div>
          </div>

          <!-- 版权信息 -->
          <!-- <div class="text-sm md:text-base">
            {{ enterpriseStore.copyright }}
            Copyright © 2012-{{ new Date().getFullYear() }} 深圳市某某科技有限公司 粤ICP备12345678号
          </div> -->

          <!-- 地址信息 -->
          <div v-if="false" class="mt-4 text-sm md:text-base text-white text-opacity-50 space-y-2">
            <p>广州：广州市科学城科学大道科学大学科技园B栋5楼 联系电话：020-1234-5678</p>
            <p>深圳：深圳市南山区某某大厦1234室 联系电话：0755-1234-5678</p>
            <p>上海：上海市某某区某某路1234号 联系电话：021-1234-5678</p>
          </div>
        </div>
        <div v-if="false" class="text-center">
          <img alt="企业AI专家微信公众号二维码" title="扫码关注公众号" class="size-24 md:size-28 mx-auto mb-2" />
          <p class="text-sm md:text-base text-white text-opacity-50">微信扫码</p>
          <p class="text-sm md:text-base text-white text-opacity-50">和创始人交个朋友</p>
        </div>
      </div>
      <div class="w-full flex justify-center items-center gap-1.5 text-xs absolute top-1/2 -translate-y-1/2 left-0 right-0">
        <span>本网站由</span>
        <img :src="$getPublicPath(`/images/53ai-hub.png`)" class="flex-none w-[72px] object-cover" />
        <span>提供技术支持</span>
      </div>
    </div>
  </template>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useNavigationStore } from '@/stores/modules/navigation'
import { useEnterpriseStore } from '@/stores/modules/enterprise'

const navigationStore = useNavigationStore()

const enterpriseStore = useEnterpriseStore()

// 移动端底部导航栏
const footerList = ref([
  {
    name: 'Agent',
    text: 'module.agent',
    title: '',
    icon: 'app',
    stroke: true
  },
  {
    name: 'Prompt',
    text: 'module.prompt',
    title: '',
    icon: 'prompt',
    stroke: true
  },
  {
    name: 'Toolkit',
    text: 'module.toolbox',
    title: '',
    icon: 'toolkit',
    stroke: true
  }
])

onMounted(async () => {
  await navigationStore.fetchNavigations()
  footerList.value[0].title = navigationStore.agentNavigation.name || ''
  footerList.value[1].title = navigationStore.promptNavigation.name || ''
  footerList.value[2].title = navigationStore.toolkitNavigation.name || ''
})

withDefaults(
  defineProps<{
    fixed?: boolean
  }>(),
  {
    fixed: true
  }
)
</script>

<style scoped></style>
