<template>
  <div class="detail-page">
    <AppSidebar
      v-model="sidebarOpen"
      brand="任务中心"
      :navItems="navItems"
      @navClick="handleNavClick">
      <TasksSidebar
        :tasks="taskList"
        v-model:expandedId="sidebarExpandedId"
        :sidebarOpen="true"
        :showTitle="false"
        :showDetail="false"
        :collapsible="false"
        @navigate="switchTask"
        @blocked="handleBlocked"
      />
    </AppSidebar>
    
    <main class="detail-main">

      <!-- ═══ 顶部操作栏 ═══ -->
      <div class="top-bar">
        <button class="btn-back" @click="goBack">
          <i class="fas fa-arrow-left"></i>
          <span>返回列表</span>
        </button>
      </div>

      <!-- ═══ 未找到（仅非加载中且确实无数据时） ═══ -->
      <div v-if="!task && !loading" class="not-found">
        <div class="nf-icon"><i class="fas fa-ghost"></i></div>
        <h2>任务不存在</h2>
        <p>未找到 ID 为 "{{ route.params.id }}" 的任务</p>
        <button class="btn-back-lg" @click="goBack">
          <i class="fas fa-arrow-left"></i> 返回任务列表
        </button>
      </div>

      <template v-if="task">
        <!-- ═══ 内容区：预览 + 代码并排 ═══ -->
        <div class="content-row" :class="{ 'is-switching': loading }">
          <!-- 切换任务时的加载遮罩：覆盖旧内容，避免整页闪烁 -->
          <div v-if="loading" class="switch-loading">
            <div class="switch-spinner"></div>
            <span>正在加载任务…</span>
          </div>
          <div class="content-left">
            <div class="pv-card">
             
              <div class="pv-body">
                <!-- Flutter：有构建预览 URL 时走 iframe 真机预览 -->
                <iframe
                  v-if="task.framework === 'Flutter' && taskDetail?.iframeUrl"
                  :src="taskDetail.iframeUrl"
                  class="pv-iframe"
                  frameborder="0"
                  allow="autoplay; fullscreen"
                ></iframe>
                <!-- 其它平台（vue/react）及 flutter 未构建时的兜底：HTML 内联预览 -->
                <PreviewTemplate
                  v-else-if="taskDetail?.previewHtml"
                  :html="taskDetail.previewHtml"
                  :showBottomNav="task.platform === 'mobile'"
                />
                <div v-else class="pv-empty">
                  <i class="fas fa-eye-slash"></i>
                  <span>暂无预览数据</span>
                </div>
              </div>
            </div>
          </div>
          <div class="content-right">
            <div class="code-card">
             
              <div class="code-body">
                <CodeEditor
                  v-if="taskDetail"
                  :modelValue="taskDetail.code"
                  :language="codeLang(task.framework)"
                  :readonly="true"
                  height="100%"
                />
                <div v-else class="code-empty">
                  <i class="fas fa-code"></i>
                  <span>暂无代码数据</span>
                </div>
              </div>
            </div>
          </div>
        </div>

      </template>
    </main>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AppSidebar from '~/components/layout/AppSidebar.vue'
import { useMenuListStore } from '~/stores/menuList'
import TasksSidebar from '~/components/tasks/TasksSidebar.vue'
import CodeEditor from '~/components/code/CodeEditor.vue'
import PreviewTemplate from '~/components/previewTempLate/index.vue'
import { taskApi } from '~/api/task'
import { ElMessage } from 'element-plus'

// 复用同一页面实例（按参数切换而非重建），避免换任务时整页 remount 闪烁
definePageMeta({ key: false })

const route = useRoute()
const router = useRouter()
const sidebarOpen = ref(true)
const sidebarExpandedId = ref(null)

const switchTask = (id) => router.push(`/detail/${id}`)

// 排队中/失败：不可进入，仅提示
const handleBlocked = (task) => {
  if (task.status === 'failed') ElMessage.warning('任务生成失败，无法查看详情')
  else if (task.status === 'pending') ElMessage.info('任务正在排队中，暂不可查看详情')
}

// ═══ 导航 ═══
const menuStore = useMenuListStore()
const navItems = menuStore.navItems.map(item => ({
  ...item,
  active: route.path === item.to,
}))

const handleNavClick = (item) => {
  if (item.active) return
  router.push(item.to)
}

const goBack = () => router.push('/tasks')

// ═══ 工具函数 ═══
// target(后端小写) → 展示用框架名（CodeEditor / codeLang 依赖此大小写）
const FW_LABEL = { flutter: 'Flutter', vue3: 'Vue3', react: 'React' }
const codeLang = (fw) => ({ Flutter: 'dart', React: 'typescript', Vue3: 'html', Vue: 'html' }[fw] || 'dart')

// ═══ 数据（真实接口，非 Mock） ═══
const task = ref(null)         // 当前任务概要（驱动预览/代码渲染）
const taskDetail = ref(null)   // 生成结果（code + previewHtml + iframeUrl）
const taskList = ref([])       // 侧栏任务列表
const loading = ref(false)

// 拉取任务列表（侧栏）
async function fetchTaskList() {
  const res = await taskApi.getTaskList()
  if (res?.code === 200 && Array.isArray(res.data)) taskList.value = res.data
}

// 拉取单个任务详情：GET /api/v1/task/:id → TaskStatusResponse
async function loadTask(id) {
  if (!id) { task.value = null; taskDetail.value = null; return }
  loading.value = true
  try {
    const res = await taskApi.getTaskById(id)
    const d = res?.code === 200 ? res.data : null
    if (!d) { task.value = null; taskDetail.value = null; return }
    task.value = {
      id: d.task_id,
      framework: FW_LABEL[d.target] || d.target,
      platform: d.platform,
      status: d.status,
      progress: d.progress,
      images: d.images || [],
    }
    const r = d.result
    taskDetail.value = r ? { code: r.code, previewHtml: r.preview, iframeUrl: r.iframe_url } : null
  } catch {
    task.value = null
    taskDetail.value = null
  } finally {
    loading.value = false
  }
}

// 1) watch 路由参数变化：同步侧栏展开项 + 加载新任务（实例复用后统一走此处，不再另设 onBeforeRouteUpdate 以免重复请求）
watch(() => route.params.id, (id) => { sidebarExpandedId.value = id; loadTask(id) })

const copyCode = () => {
  if (!taskDetail.value?.code) return
  navigator.clipboard.writeText(taskDetail.value.code)
}

const downloadCode = () => {
  if (!taskDetail.value?.code) return
  const lang = codeLang(task.value?.framework || 'Flutter')
  const ext = { dart: '.dart', typescript: '.tsx', html: '.vue' }[lang] || '.txt'
  const blob = new Blob([taskDetail.value.code], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `generated${ext}`
  a.click()
  URL.revokeObjectURL(url)
}

// 初始加载：侧栏列表 + 当前详情
onMounted(() => {
  sidebarExpandedId.value = route.params.id
  fetchTaskList()
  loadTask(route.params.id)
})
</script>

<style scoped>
/* ═══ 布局 ═══ */
.detail-page {
  display: flex;
  height: 100vh;
  background: #0a0a0f;
  color: #e2e8f0;
  overflow: hidden;
}

.detail-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 4px 28px 28px;
  overflow-y: auto;
  min-width: 0;
}
.detail-main::-webkit-scrollbar { width: 4px; }
.detail-main::-webkit-scrollbar-track { background: transparent; }
.detail-main::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.06); border-radius: 2px; }

/* ═══ 顶部返回 ═══ */
.top-bar { padding: 14px 0 8px; }
.btn-back {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 16px;
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 8px;
  background: transparent;
  color: rgba(255,255,255,0.35);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
  font-family: inherit;
}
.btn-back:hover {
  border-color: rgba(0,255,255,0.2);
  color: #e2e8f0;
  background: rgba(255,255,255,0.02);
}

/* ═══ 信息卡片 ═══ */
.info-card {
  background: rgba(15, 22, 32, 0.75);
  border: 1px solid rgba(0,255,255,0.06);
  border-radius: 14px;
  padding: 0;
  margin-bottom: 18px;
  overflow: hidden;
}
.info-head {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 24px 28px 18px;
}
.info-thumb {
  width: 46px; height: 46px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}
.thumb--success { background: rgba(52,211,153,0.1); color: #34d399; }
.thumb--running { background: rgba(0,216,255,0.1); color: #00d8ff; }
.thumb--pending { background: rgba(255,255,255,0.03); color: rgba(255,255,255,0.3); }
.thumb--failed  { background: rgba(255,90,108,0.1); color: #ff5a6c; }

.info-main { flex: 1; min-width: 0; }
.info-title {
  font-size: 22px;
  font-weight: 700;
  color: #f1f5f9;
  margin: 0 0 10px;
  line-height: 1.3;
}
.info-meta {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  color: rgba(255,255,255,0.3);
}
.fw-tag { font-size: 11px !important; }
.meta-sep { color: rgba(255,255,255,0.1); font-size: 10px; }
.meta-plat i { margin-right: 4px; }
.meta-id {
  font-family: 'Fira Code', monospace;
  font-size: 11px;
  color: rgba(255,255,255,0.2);
  background: rgba(255,255,255,0.03);
  padding: 2px 6px;
  border-radius: 4px;
}

.info-badge { flex-shrink: 0; }
.status-label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 14px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  white-space: nowrap;
}
.slabel--success { background: rgba(52,211,153,0.08); color: #34d399; border: 1px solid rgba(52,211,153,0.15); }
.slabel--running { background: rgba(0,216,255,0.08); color: #00d8ff; border: 1px solid rgba(0,216,255,0.15); }
.slabel--pending { background: rgba(255,255,255,0.02); color: rgba(255,255,255,0.35); border: 1px solid rgba(255,255,255,0.06); }
.slabel--failed  { background: rgba(255,90,108,0.08); color: #ff5a6c; border: 1px solid rgba(255,90,108,0.15); }
.status-dot {
  width: 6px; height: 6px;
  border-radius: 50%;
  display: inline-block;
}
.slabel--success .status-dot { background: #34d399; }
.slabel--running .status-dot { background: #00d8ff; animation: dotPulse 1.5s ease-in-out infinite; }
.slabel--pending .status-dot { background: rgba(255,255,255,0.3); }
.slabel--failed .status-dot  { background: #ff5a6c; }
@keyframes dotPulse { 0%,100%{opacity:1} 50%{opacity:0.3} }

/* 状态栏 */
.info-stats {
  display: flex;
  align-items: center;
  gap: 0;
  padding: 14px 28px;
  border-top: 1px solid rgba(255,255,255,0.04);
  background: rgba(0,0,0,0.12);
}
.istat {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 24px;
}
.istat:first-child { padding-left: 0; }
.istat-lbl {
  font-size: 10px;
  color: rgba(255,255,255,0.25);
  text-transform: uppercase;
  letter-spacing: 0.3px;
  white-space: nowrap;
}
.istat-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 120px;
}
.istat-bar .el-progress { flex: 1; }
:deep(.dp .el-progress-bar__outer) { background: rgba(255,255,255,0.05) !important; }
.istat-pct { font-size: 12px; color: rgba(255,255,255,0.35); width: 28px; text-align: right; flex-shrink: 0; }
.istat-val { font-size: 13px; color: rgba(255,255,255,0.55); font-weight: 500; white-space: nowrap; }
.istat-divider {
  width: 1px; height: 18px;
  background: rgba(255,255,255,0.06);
}

/* ═══ 内容区：预览 + 代码并排 ═══ */
.content-row {
  position: relative;
  display: flex;
  gap: 18px;
  flex: 1;
  min-height: 0;
  margin-bottom: 18px;
}
/* 切换加载遮罩：旧内容压暗 + 居中转圈 */
.content-row.is-switching .content-left,
.content-row.is-switching .content-right { opacity: 0.25; filter: blur(1px); }
.content-row .content-left,
.content-row .content-right { transition: opacity .2s, filter .2s; }
.switch-loading {
  position: absolute; inset: 0; z-index: 5;
  display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 12px;
  color: rgba(255,255,255,0.55); font-size: 13px;
}
.switch-spinner {
  width: 34px; height: 34px; border-radius: 50%;
  border: 3px solid rgba(0, 255, 255, 0.18);
  border-top-color: #00ffff;
  animation: switch-spin .8s linear infinite;
}
@keyframes switch-spin { to { transform: rotate(360deg); } }
.content-left {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}
.content-right {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

/* 代码面板（在 content-right 内） */
.code-card {
  background: rgba(15, 22, 32, 0.6);
  border: 1px solid rgba(255,255,255,0.04);
  border-radius: 12px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
}

/* 预览卡片 */
.pv-card {
  height: 100%;
  background: rgba(15, 22, 32, 0.6);
  border: 1px solid rgba(255,255,255,0.04);
  border-radius: 12px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
.pv-header {
  padding: 10px 16px;
  background: rgba(0,0,0,0.2);
  border-bottom: 1px solid rgba(255,255,255,0.04);
}
.pv-title-row {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: rgba(255,255,255,0.35);
}
.pv-title-row i { font-size: 13px; }
.pv-title { font-weight: 600; }
.pv-device {
  margin-left: auto;
  font-size: 10px;
  color: rgba(255,255,255,0.2);
  text-transform: uppercase;
  letter-spacing: 0.3px;
}
.pv-body {
  flex: 1;
  display: flex;
  overflow: hidden;
}
.pv-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: rgba(255,255,255,0.06);
}
.pv-empty i { font-size: 28px; }
.pv-empty span { font-size: 12px; color: rgba(255,255,255,0.15); }

/* Flutter iframe 真机预览：撑满预览容器 */
.pv-iframe {
  flex: 1;
  width: 100%;
  height: 100%;
  border: 0;
  background: #0a0a0f;
}

.code-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  background: rgba(0,0,0,0.25);
  border-bottom: 1px solid rgba(255,255,255,0.04);
}
.code-dots { display: flex; gap: 6px; }
.pdot { width: 9px; height: 9px; border-radius: 50%; }
.pdot-red    { background: #ff5f56; }
.pdot-yellow { background: #ffbd2e; }
.pdot-green  { background: #27c93f; }

.code-lang {
  font-size: 11px;
  color: rgba(255,255,255,0.25);
  text-transform: uppercase;
  letter-spacing: 0.3px;
}
.code-actions { margin-left: auto; display: flex; gap: 4px; }
.code-act {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 6px;
  background: transparent;
  color: rgba(255,255,255,0.3);
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s;
  font-family: inherit;
}
.code-act:hover {
  background: rgba(255,255,255,0.03);
  border-color: rgba(0,255,255,0.15);
  color: #e2e8f0;
}
.code-body {
  flex: 1;
  min-height: 0;
}
:deep(.code-editor-wrapper) {
  border: none !important;
  border-radius: 0 !important;
}
:deep(.code-editor-wrapper:hover) { box-shadow: none !important; }
.code-empty {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: rgba(255,255,255,0.06);
}
.code-empty i { font-size: 28px; }
.code-empty span { font-size: 12px; color: rgba(255,255,255,0.15); }

/* ═══ Not Found ═══ */
.not-found {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-top: -60px;
}
.nf-icon { font-size: 52px; color: rgba(255,255,255,0.05); margin-bottom: 4px; animation: ghostFloat 3s ease-in-out infinite; }
@keyframes ghostFloat { 0%,100%{transform:translateY(0)} 50%{transform:translateY(-8px)} }
.not-found h2 { font-size: 20px; color: rgba(255,255,255,0.3); margin: 0; }
.not-found p { font-size: 13px; color: rgba(255,255,255,0.15); margin: 0; }
.btn-back-lg {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 20px;
  border: 1px solid rgba(0,255,255,0.12);
  border-radius: 8px;
  background: transparent;
  color: rgba(255,255,255,0.4);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
  margin-top: 6px;
  font-family: inherit;
}
.btn-back-lg:hover {
  background: rgba(255,255,255,0.03);
  border-color: rgba(0,255,255,0.25);
  color: #e2e8f0;
}

/* ═══ 响应式 ═══ */
@media (max-width: 1100px) {
  .content-row { flex-direction: column; }
  .detail-main { padding: 4px 16px 16px; }
}
</style>
