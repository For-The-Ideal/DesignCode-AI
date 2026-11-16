<template>
  <div class="home-page">
    <AppSidebar
      v-model="sidebarOpen"
      brand="首页"
      expandedWidth="225px"
      :navItems="navItems"
      @navClick="handleNavClick"
    />

    <main class="home-main">
      <div class="screens-container">
        <!-- 第一屏 -->
        <HeroSection @getStarted="scrollToSecond" @watchDemo="playDemo" />

        <!-- 第二屏 -->
        <SSEGenerator />

        <!-- 第三屏 -->
        <DiagnosticBubbles />
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { ElMessage } from 'element-plus'
import AppSidebar from '~/components/layout/AppSidebar.vue'
import HeroSection from '~/components/landing/HeroSection.vue'
import SSEGenerator from '~/components/landing/SSEGenerator.vue'
import DiagnosticBubbles from '~/components/landing/DiagnosticBubbles.vue'
import { useMenuListStore } from '~/stores/menuList'
import { useCommonStore } from '~/stores/common'
import { useUserStore } from '~/stores/user'

const route = useRoute()
const router = useRouter()
const sidebarOpen = ref(true)

const { isLogin } = storeToRefs(useUserStore())
const commonStore = useCommonStore()

// ═══ 侧边栏导航 ═══
const menuStore = useMenuListStore()
// 当前路径在 setup 里取，交给 store getter 算 active（getter 里不能 useRouter）
const navItems = computed(() => menuStore.navItemsWithActive(route.path))
const handleNavClick = (item) => {
  if (item.active) return
  // 未登录只允许停留在首页，其余入口统一点登录弹窗（与工作台侧栏行为一致）
  if (!isLogin.value) return commonStore.setLoginModalVisible(true)
  router.push(item.to)
}

// ═══ 落地页交互 ═══
// 锁定页面级滚动，只让容器内部吸附滚动。组件卸载时 Nuxt 会自动移除该 class，
// 不会像之前那样把 html/body 写成非 scoped 全局样式、污染其它页面的滚动。
useHead({ bodyAttrs: { class: 'landing-lock' } })

const scrollToSecond = () => {
  const container = document.querySelector('.screens-container')
  // 容器比视口矮（扣掉了 AppHeader），按容器自身高度翻页才能对齐吸附点
  if (container) container.scrollBy({ top: container.clientHeight, behavior: 'smooth' })
}

const playDemo = () => {
  ElMessage.warning('演示视频即将上线')
}
</script>

<style scoped>
.home-page {
  display: flex;
  /* 扣掉顶部 AppHeader：padding 12px × 2 + logo 内容 40px + 底边框 1px */
  height: calc(100vh - 65px);
  background: #0a0a0f;
}
.home-main {
  flex: 1;
  min-width: 0;
  overflow: hidden;
}
.screens-container {
  height: 100%;
  scroll-snap-type: y mandatory;
  overflow-y: scroll;
  scroll-behavior: smooth;
}
</style>

<style>
/* 仅首页（落地页）生效：class 由上面的 useHead 在本页挂载/卸载时增删 */
body.landing-lock {
  overflow: hidden;
}
body.landing-lock::-webkit-scrollbar {
  display: none;
}
</style>
