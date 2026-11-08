<template>
  <div class="code-page">
    <CodeSidebar />

    <main class="code-main code-main--workbench">
      <!-- 单图 / 多图统一工作台：左栏页面列表、中栏画布、右栏区域 + 页面关系 + 生成配置 -->
      <WorkbenchView />
    </main>
  </div>
</template>

<script setup>
import { watch } from 'vue'
import CodeSidebar from '~/components/code/CodeSidebar.vue'
import WorkbenchView from '~/components/workbench/WorkbenchView.vue'
import { useCodeStore, normalizeMode } from '~/stores/code'

const route = useRoute()
const router = useRouter()
const store = useCodeStore()

// URL → store：带参进入或浏览器前进/后退时同步单图/多图
// 无参视为默认单图，保证后退时状态与地址栏一致
watch(
  () => route.query.mode,
  (val) => store.setMode(normalizeMode(val))
)

// store → URL：页面内 tab 切换后写回地址栏（replace，不增加历史记录）
// 默认单图不写进 URL，避免首次进入/后退时被补成 ?mode=single
watch(
  () => store.mode,
  (val) => {
    if (!import.meta.client) return
    if (route.query.mode === val) return
    if (val === 'single' && !route.query.mode) return
    router.replace({ path: route.path, query: { ...route.query, mode: val } })
  }
)
</script>

<style scoped>
.code-page {
  display: flex;
  min-height: calc(100vh - 140px);
  background: #0a0a0f;
  position: relative;
}

.code-main {
  flex: 1;
  overflow-y: auto;
  padding: 20px 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}
/* 工作台：几乎占满视口（只让出顶部导航），内部各栏自滚动，不再整页滚动 */
.code-main--workbench {
  overflow: hidden;
  height: calc(100vh - 84px);
  padding: 8px 12px;
}
</style>
