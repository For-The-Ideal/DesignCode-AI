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
      />
    </AppSidebar>
    
    <main class="detail-main">

      <!-- ═══ 顶部操作栏 ═══ -->
      <div class="top-bar">
        <button class="btn-back" @click="goBack">
          <i class="fas fa-arrow-left"></i>
          <span>返回列表</span>
        </button>

        <!-- 右侧操作：重新生成 + 质量检查 -->
        <div class="top-actions">
          <!-- 重新生成：把任务详情回填到工作台（图片/树/素材/页面关系/配置）再跳转 -->
          <button v-if="canRegenerate" class="btn-regen" type="button" @click="regenerate">
            <i class="fas fa-rotate-right"></i>
            <span>重新生成</span>
          </button>

          <!-- 质量检查：收进顶栏，点开看明细（不再挤在预览下方） -->
          <div v-if="currentQuality" class="qc-wrap">
            <button
              class="qc-chip"
              :class="`qc-chip--${currentQuality.status}`"
              type="button"
              @click="qcOpen = !qcOpen"
            >
              <i class="fas qc-icon" :class="qcIcon"></i>
              <span class="qc-title">质量检查</span>
              <span class="qc-badge">{{ qcLabel }}</span>
              <span class="qc-count">{{ currentQuality.passed }}/{{ currentQuality.total }}</span>
              <i class="fas qc-arrow" :class="qcOpen ? 'fa-chevron-up' : 'fa-chevron-down'"></i>
            </button>

            <!-- 透明遮罩：点页面空白处关闭浮层 -->
            <div v-if="qcOpen" class="qc-backdrop" @click="qcOpen = false"></div>
            <Transition name="qc-fade">
              <div v-if="qcOpen" class="qc-pop" :class="`qc-pop--${currentQuality.status}`">
                <div class="qc-pop-head">
                  <span class="qc-pop-title">质量检查</span>
                  <span class="qc-pop-count">{{ currentQuality.passed }}/{{ currentQuality.total }} 项通过</span>
                </div>
                <div v-if="!currentQuality.issues?.length" class="qc-ok">
                  <i class="fas fa-circle-check"></i> 所有规则均通过
                </div>
                <div
                  v-for="(issue, i) in currentQuality.issues"
                  :key="i"
                  class="qc-item"
                  :class="`qc-item--${issue.severity}`"
                >
                  <span class="qc-rule">{{ issue.rule }}</span>
                  <span class="qc-detail">{{ issue.detail }}</span>
                </div>
              </div>
            </Transition>
          </div>
        </div>
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
              <!-- 多页任务：页签切换预览（与右栏代码折叠双向联动） -->
              <div v-if="pages.length > 1 && !taskDetail?.iframeUrl" class="page-tabs">
                <button
                  v-for="(p, i) in pages"
                  :key="p.page_id"
                  class="page-tab"
                  :class="{ active: i === activePage }"
                  type="button"
                  @click="selectPage(i)"
                >
                  <span class="pt-idx">{{ String(i + 1).padStart(2, '0') }}</span>
                  <span class="pt-name">{{ p.page_name || p.page_id }}</span>
                </button>
              </div>
              <div class="pv-body">
                <!-- 生成中 / 排队中：走预览模板内置的生成动画（html 为空即 loading 态） -->
                <PreviewTemplate
                  v-if="isLive"
                  html=""
                  :showBottomNav="task.platform === 'mobile'"
                />
                <!-- 失败：失败占位样式 -->
                <div v-else-if="isFailed" class="pv-failed">
                  <div class="fail-ring"><i class="fas fa-triangle-exclamation"></i></div>
                  <p class="fail-title">任务生成失败</p>
                  <p class="fail-desc">{{ failReason || '生成过程中出现异常，该任务没有可用的预览与代码' }}</p>
                  <span class="fail-meta" v-if="typeof task.progress === 'number'">
                    失败前进度 {{ task.progress }}%<template v-if="task.currentStep"> · {{ task.currentStep }}</template>
                  </span>
                  <button class="fail-btn" @click="goBack">
                    <i class="fas fa-rotate-left"></i> 返回列表重新生成
                  </button>
                </div>
                <!-- Flutter：有构建预览 URL 时走 iframe 真机预览 -->
                <iframe
                  v-else-if="task.framework === 'Flutter' && taskDetail?.iframeUrl"
                  :src="taskDetail.iframeUrl"
                  class="pv-iframe"
                  frameborder="0"
                  allow="autoplay; fullscreen"
                ></iframe>
                <!-- 其它平台（vue/react）及 flutter 未构建时的兜底：HTML 内联预览（多页时只展示当前页） -->
                <PreviewTemplate
                  v-else-if="currentPreview"
                  :html="currentPreview"
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
             
              <div class="code-body" :class="{ 'code-body--pages': pages.length > 1 }">
                <!-- 多页：每页一个可折叠代码面板（同一时间只展开一页，展开才挂载编辑器） -->
                <div v-if="pages.length > 1" class="page-code-list">
                  <div
                    v-for="(p, i) in pages"
                    :key="p.page_id"
                    class="page-code-item"
                    :class="{ open: i === activePage }"
                  >
                    <button class="pc-head" type="button" @click="selectPage(i)">
                      <i class="fas" :class="i === activePage ? 'fa-chevron-down' : 'fa-chevron-right'"></i>
                      <span class="pc-idx">{{ String(i + 1).padStart(2, '0') }}</span>
                      <span class="pc-name">{{ p.page_name || p.page_id }}</span>
                      <span v-if="p.route" class="pc-route">{{ p.route }}</span>
                      <span class="pc-lines">{{ lineCount(p.code) }} 行</span>
                      <span class="pc-copy" role="button" title="复制本页代码" @click.stop="copyPage(p.code)">
                        <i class="fa-regular fa-copy"></i>
                      </span>
                    </button>
                    <div v-if="i === activePage" class="pc-body">
                      <CodeEditor
                        :modelValue="p.code"
                        :language="codeLang(task.framework)"
                        :label="task.framework"
                        :readonly="true"
                        height="100%"
                      />
                    </div>
                  </div>
                </div>
                <CodeEditor
                  v-else-if="taskDetail"
                  :modelValue="taskDetail.code"
                  :language="codeLang(task.framework)"
                  :label="task.framework"
                  :readonly="true"
                  height="100%"
                />
                <div v-else class="code-empty" :class="{ 'code-empty--failed': isFailed }">
                  <i class="fas" :class="isLive ? 'fa-spinner fa-spin' : isFailed ? 'fa-triangle-exclamation' : 'fa-code'"></i>
                  <span>{{ isLive ? '代码生成中，完成后自动展示' : isFailed ? '任务失败，未生成代码' : '暂无代码数据' }}</span>
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
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import AppSidebar from '~/components/layout/AppSidebar.vue'
import { useMenuListStore } from '~/stores/menuList'
import { useUserStore } from '~/stores/user'
import TasksSidebar from '~/components/tasks/TasksSidebar.vue'
import CodeEditor from '~/components/code/CodeEditor.vue'
import PreviewTemplate from '~/components/previewTempLate/index.vue'
import { taskApi } from '~/api/task'
import { ElMessageBox } from 'element-plus'
import { useCodeStore } from '~/stores/code'

// 复用同一页面实例（按参数切换而非重建），避免换任务时整页 remount 闪烁
definePageMeta({ key: false })

const route = useRoute()
const router = useRouter()
const sidebarOpen = ref(true)
const sidebarExpandedId = ref(null)

const switchTask = (id) => router.push(`/detail/${id}`)

// ═══ 导航 ═══
const menuStore = useMenuListStore()
// 当前路径在 setup 里取，交给 store getter 算 active（getter 里不能 useRouter）
const navItems = computed(() => menuStore.navItemsWithActive(route.path))

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
const taskDetail = ref(null)   // 生成结果（code + preview + iframeUrl + 多页 pages）
const activePage = ref(0)      // 多页模式当前页下标（左页签 / 右折叠共用）
const taskList = ref([])       // 侧栏任务列表
const loading = ref(false)
const failReason = ref('')     // 失败原因（results.error_msg；后端在 failed 时一并返回）

// ═══ 状态派生：生成中/排队中 → 生成动画；失败 → 失败样式 ═══
const isLive = computed(() => ['running', 'pending'].includes(task.value?.status))
const isFailed = computed(() => task.value?.status === 'failed')

// 拉取任务列表（侧栏）
async function fetchTaskList() {
  const res = await taskApi.getTaskList()
  if (res?.code === 200 && Array.isArray(res.data)) taskList.value = res.data
}

// 拉取单个任务详情：GET /api/v1/task/:id → TaskStatusResponse
// silent=true：失败复核用的静默刷新，不触发切换遮罩，请求失败也不清空当前页
async function loadTask(id, silent = false) {
  if (!id) { task.value = null; taskDetail.value = null; return }
  if (!silent) loading.value = true
  try {
    const res = await taskApi.getTaskById(id)
    // 请求期间又切了任务：丢弃旧响应，避免覆盖新任务（表现为“任务不存在”）
    if (id !== route.params.id) return
    // 被同键去重取消（request.js 对 abort 的 Promise resolve(null)）：保留当前页面
    if (!res) return
    // 只有服务端明确 404 才判定任务不存在；其它错误不清空已展示内容
    if (res.code === 404) { task.value = null; taskDetail.value = null; return }
    const d = res.code === 200 ? res.data : null
    if (!d) return
    task.value = {
      id: d.task_id,
      target: d.target,
      framework: FW_LABEL[d.target] || d.target,
      platform: d.platform,
      status: d.status,
      progress: d.progress,
      currentStep: d.current_step,
      images: d.images || [],
      // 重新生成回填：页面关系 + 生成配置（图片/树已在 images 中）
      page_links: d.page_links || [],
      options: d.options || [],
      advanced: d.advanced || [],
      component_lib: d.component_lib || '',
    }
    const r = d.result
    // 失败：不挂载空产物（否则右栏会出现空编辑器），只保留原因供失败卡片展示
    if (d.status === 'failed') {
      taskDetail.value = null
      failReason.value = r?.error_msg || ''
      // 首次拉取就遇到 failed 也要复核：任务可能仍由另一个执行者跑着
      if (!silent) scheduleFailCheck()
    } else {
      failReason.value = ''
      taskDetail.value = r
        ? {
            code: r.code,
            preview: r.preview,
            iframeUrl: r.iframe_url,
            // 多页产物：pages 为空/长度 1 时自动走原有单编辑器逻辑
            pages: r.pages || [],
            appPreview: r.app_preview || '',
            navMap: r.nav_map || [],
            qualityReport: r.quality_report || null,
          }
        : null
    }
    activePage.value = 0
  } catch {
    // 网络异常：不清空已展示的任务，避免误判“任务不存在”
  } finally {
    if (!silent) loading.value = false
  }
}

// ═══ 失败状态复核 ═══
// SSE / 接口可能携带竞态写入的 failed（如另一个实例启动时的残留清理错过），
// 而任务实际仍由执行者跑着并会写回 success。失败不直接采信：
// 延迟 2 秒回查 DB，连续 2 次仍为 failed 才渲染失败卡片，避免“生成中突然失败”。
const FAIL_CHECK_DELAY = 2000
const FAIL_CHECK_TIMES = 2
let failTimer = null
let failChecks = 0

function scheduleFailCheck () {
  if (failTimer || failChecks >= FAIL_CHECK_TIMES) return
  failTimer = setTimeout(async () => {
    failTimer = null
    failChecks++
    await loadTask(route.params.id, true)
    if (isFailed.value && failChecks < FAIL_CHECK_TIMES) scheduleFailCheck()
  }, FAIL_CHECK_DELAY)
}

function resetFailCheck () {
  if (failTimer) clearTimeout(failTimer)
  failTimer = null
  failChecks = 0
}

// ═══ 多页：左页签预览 / 右折叠代码双向联动 ═══
const pages = computed(() => taskDetail.value?.pages || [])
const currentPreview = computed(() => {
  if (pages.value.length) return pages.value[activePage.value]?.preview || ''
  return taskDetail.value?.preview || ''
})
const selectPage = (i) => { activePage.value = i }

// ═══ 质量检查：当前展示页的报告（多页优先取本页，单页回退入口页） ═══
const qcOpen = ref(false)
const currentQuality = computed(() => {
  const p = pages.value[activePage.value]
  return p?.quality_report || taskDetail.value?.qualityReport || null
})
const qcLabel = computed(() => ({ pass: '通过', warning: '有告警', fail: '未通过' }[currentQuality.value?.status] || ''))
const qcIcon = computed(() => ({ pass: 'fa-circle-check', warning: 'fa-triangle-exclamation', fail: 'fa-circle-xmark' }[currentQuality.value?.status] || 'fa-shield-halved'))
const lineCount = (code) => (code || '').split('\n').length
// 复制文本：优先 Clipboard API，非安全上下文（http）回退 execCommand
async function writeClipboard(text) {
  if (!text) return
  try {
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(text)
      return
    }
  } catch {}
  try {
    const ta = document.createElement('textarea')
    ta.value = text
    ta.setAttribute('readonly', '')
    ta.style.cssText = 'position:fixed;top:-9999px;left:-9999px;opacity:0'
    document.body.appendChild(ta)
    ta.select()
    ta.setSelectionRange(0, text.length)
    document.execCommand('copy')
    document.body.removeChild(ta)
  } catch {}
}
const copyPage = (code) => writeClipboard(code)

// 1) watch 路由参数变化：同步侧栏展开项 + 加载新任务（实例复用后统一走此处，不再另设 onBeforeRouteUpdate 以免重复请求）
watch(() => route.params.id, (id) => {
  resetFailCheck()
  qcOpen.value = false
  sidebarExpandedId.value = id
  loadTask(id)
})

// ═══ 生成中跟随用户级 SSE：进度实时同步，完成后自动拉取生成结果 ═══
const userStore = useUserStore()
const { taskProgress } = storeToRefs(userStore)
watch(taskProgress, (map) => {
  const live = map?.[route.params.id]
  if (!live || !task.value) return
  if (typeof live.progress === 'number') task.value.progress = live.progress
  if (live.status && live.status !== 'failed') {
    resetFailCheck()
    task.value.status = live.status
    // 状态推进到完成且尚未拿到结果时重新拉取（失败无需结果）
    if (live.status === 'success' && !taskDetail.value) loadTask(route.params.id)
  }
  // failed 不直接采信：回查 DB，仍失败才展示
  if (live.status === 'failed') scheduleFailCheck()
}, { immediate: true })

// ═══ 重新生成：任务详情回填工作台后跳转单图/多图 ═══
// 恢复会整桶替换目标模式的数据，工作台已有未提交内容时需要用户确认
const codeStore = useCodeStore()
const canRegenerate = computed(() => !!task.value?.images?.length && !isLive.value)

async function regenerate() {
  const t = task.value
  if (!t?.images?.length) return
  const mode = t.images.length > 1 ? 'multi' : 'single'
  const existing = codeStore.modes[mode]
  const dirty = existing.images.length || existing.assets.length || existing.pageLinks.length
  if (dirty) {
    try {
      await ElMessageBox.confirm(
        `「${mode === 'multi' ? '多图' : '单图'}」工作台还有未提交的内容，重新生成会用本任务的图片与结构覆盖它。`,
        '覆盖当前编辑',
        { confirmButtonText: '覆盖并继续', cancelButtonText: '取消', type: 'warning' },
      )
    } catch {
      return
    }
  }
  const restored = codeStore.restoreFromTask(t)
  if (!restored) return
  router.push({ path: '/code', query: { mode: restored } })
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
.top-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 0 8px;
  flex-shrink: 0;
}
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

.top-actions {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 10px;
}
.btn-regen {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border-radius: 8px;
  border: 1px solid rgba(0, 255, 255, 0.2);
  background: rgba(0, 255, 255, 0.06);
  color: #7ff3f3;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
  font-family: inherit;
}
.btn-regen:hover {
  background: rgba(0, 255, 255, 0.12);
  border-color: rgba(0, 255, 255, 0.4);
  color: #d8ffff;
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
  flex: 1;
  min-height: 0;
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

/* ═══ 质量检查：顶栏胶囊 + 浮层明细（不占预览区高度） ═══ */
.qc-wrap {
  position: relative;
  margin-left: 0;
  z-index: 30;
}
.qc-chip {
  /* 高于透明遮罩：胶囊本身仍可点击/悬浮，点一下同样收起浮层 */
  position: relative;
  z-index: 22;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 14px;
  border-radius: 8px;
  border: 1px solid rgba(255,255,255,0.06);
  background: rgba(15, 22, 32, 0.6);
  color: rgba(255,255,255,0.55);
  font-size: 12px;
  font-family: inherit;
  cursor: pointer;
  transition: all 0.18s;
}
.qc-chip:hover { border-color: rgba(0,255,255,0.2); color: #e2e8f0; }
.qc-chip i { font-size: 13px; }
.qc-chip--pass { border-color: rgba(52, 211, 153, 0.2); }
.qc-chip--warning { border-color: rgba(251, 191, 36, 0.22); }
.qc-chip--fail { border-color: rgba(248, 113, 113, 0.24); }
.qc-chip--pass .qc-icon { color: #34d399; }
.qc-chip--warning .qc-icon { color: #fbbf24; }
.qc-chip--fail .qc-icon { color: #f87171; }
.qc-title { font-weight: 600; }
.qc-badge {
  padding: 1px 7px;
  border-radius: 999px;
  font-size: 10px;
}
.qc-chip--pass .qc-badge { background: rgba(52, 211, 153, 0.12); color: #34d399; }
.qc-chip--warning .qc-badge { background: rgba(251, 191, 36, 0.12); color: #fbbf24; }
.qc-chip--fail .qc-badge { background: rgba(248, 113, 113, 0.14); color: #f87171; }
.qc-count {
  color: rgba(255,255,255,0.3);
  font-size: 11px;
  font-variant-numeric: tabular-nums;
}
.qc-arrow { color: rgba(255,255,255,0.25); font-size: 10px; }

/* 透明遮罩：点页面空白处关闭浮层（浮层 z-index 高于遮罩） */
.qc-backdrop {
  position: fixed;
  inset: 0;
  z-index: 20;
}
.qc-pop {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  z-index: 21;
  width: 400px;
  max-width: 90vw;
  max-height: 60vh;
  overflow-y: auto;
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 7px;
  background: rgba(13, 19, 28, 0.98);
  border: 1px solid rgba(255,255,255,0.07);
  border-radius: 12px;
  box-shadow: 0 18px 48px rgba(0, 0, 0, 0.55);
}
.qc-pop--pass { border-color: rgba(52, 211, 153, 0.18); }
.qc-pop--warning { border-color: rgba(251, 191, 36, 0.2); }
.qc-pop--fail { border-color: rgba(248, 113, 113, 0.22); }
.qc-pop::-webkit-scrollbar { width: 4px; }
.qc-pop::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.08); border-radius: 2px; }
.qc-pop-head {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-bottom: 7px;
  border-bottom: 1px solid rgba(255,255,255,0.05);
}
.qc-pop-title { font-size: 12px; font-weight: 600; color: rgba(255,255,255,0.7); }
.qc-pop-count { margin-left: auto; font-size: 11px; color: rgba(255,255,255,0.3); }
.qc-fade-enter-active,
.qc-fade-leave-active { transition: opacity 0.16s ease, transform 0.16s ease; }
.qc-fade-enter-from,
.qc-fade-leave-to { opacity: 0; transform: translateY(-4px); }
.qc-ok {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #34d399;
  font-size: 12px;
}
.qc-item {
  display: flex;
  gap: 8px;
  font-size: 11px;
  line-height: 1.5;
  color: rgba(255,255,255,0.5);
}
.qc-rule {
  flex-shrink: 0;
  padding: 0 6px;
  border-radius: 4px;
  font-size: 10px;
  height: 18px;
  line-height: 18px;
}
.qc-item--critical .qc-rule { background: rgba(248, 113, 113, 0.14); color: #f87171; }
.qc-item--warning .qc-rule { background: rgba(251, 191, 36, 0.12); color: #fbbf24; }

/* ═══ 多页页签（左栏） ═══ */
.page-tabs {
  display: flex;
  gap: 6px;
  padding: 10px 12px 0;
  overflow-x: auto;
  flex-shrink: 0;
}
.page-tabs::-webkit-scrollbar { height: 3px; }
.page-tabs::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.08); border-radius: 2px; }
.page-tab {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 12px;
  border-radius: 8px;
  border: 1px solid rgba(255,255,255,0.08);
  background: rgba(255,255,255,0.03);
  color: rgba(255,255,255,0.5);
  font-size: 12px;
  white-space: nowrap;
  cursor: pointer;
  transition: all 0.18s ease;
  font-family: inherit;
}
.page-tab:hover { color: #e2e8f0; border-color: rgba(0,255,255,0.25); }
.page-tab.active {
  color: #06121a;
  background: linear-gradient(135deg, #22d3ee, #06b6d4);
  border-color: transparent;
  font-weight: 700;
}
.pt-idx { font-variant-numeric: tabular-nums; opacity: 0.75; }
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

/* ═══ 生成失败占位 ═══ */
.pv-failed {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 28px;
  text-align: center;
}
.fail-ring {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 26px;
  color: #ff5a6c;
  background: rgba(255, 90, 108, 0.1);
  border: 1px solid rgba(255, 90, 108, 0.35);
  animation: failPulse 2.4s ease-in-out infinite;
}
@keyframes failPulse {
  0%, 100% { box-shadow: 0 0 20px rgba(255, 90, 108, 0.18); }
  50%      { box-shadow: 0 0 34px rgba(255, 90, 108, 0.32); }
}
.fail-title { margin: 6px 0 0; font-size: 17px; font-weight: 700; color: #ff8f9c; }
.fail-desc { margin: 0; font-size: 12px; color: rgba(255,255,255,0.35); line-height: 1.7; }
.fail-meta { font-size: 11px; color: rgba(255, 90, 108, 0.6); font-family: 'Fira Code', monospace; }
.fail-btn {
  margin-top: 6px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 7px 18px;
  border: 1px solid rgba(255, 90, 108, 0.35);
  border-radius: 8px;
  background: rgba(255, 90, 108, 0.08);
  color: #ff8f9c;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
  font-family: inherit;
}
.fail-btn:hover { background: rgba(255, 90, 108, 0.16); border-color: rgba(255, 90, 108, 0.55); }
.code-empty--failed i { color: rgba(255, 90, 108, 0.4); }
.code-empty--failed span { color: rgba(255, 90, 108, 0.55); }

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

/* ═══ 多页代码折叠（右栏）：同时只展开一页，展开才挂载编辑器 ═══ */
.code-body--pages { display: flex; flex-direction: column; overflow: hidden; }
.page-code-list { flex: 1; min-height: 0; display: flex; flex-direction: column; overflow-y: auto; }
.page-code-list::-webkit-scrollbar { width: 4px; }
.page-code-list::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.08); border-radius: 2px; }
.page-code-item {
  flex-shrink: 0;
  border-bottom: 1px solid rgba(255,255,255,0.06);
}
.page-code-item.open {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}
.pc-head {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 10px 14px;
  background: transparent;
  border: none;
  cursor: pointer;
  color: rgba(255,255,255,0.6);
  font-size: 12px;
  font-family: inherit;
  text-align: left;
  transition: background 0.18s ease, color 0.18s ease;
}
.pc-head:hover { background: rgba(255,255,255,0.03); }
.page-code-item.open .pc-head { color: #22d3ee; background: rgba(34,211,238,0.06); }
.pc-idx { font-variant-numeric: tabular-nums; opacity: 0.7; flex-shrink: 0; }
.pc-name { font-weight: 600; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.pc-route { font-family: 'Fira Code', monospace; font-size: 11px; color: rgba(255,255,255,0.3); flex-shrink: 0; }
.pc-lines { margin-left: auto; font-size: 11px; color: rgba(255,255,255,0.25); flex-shrink: 0; }
.pc-copy { padding: 2px 6px; border-radius: 4px; color: rgba(255,255,255,0.35); flex-shrink: 0; }
.pc-copy:hover { color: #22d3ee; background: rgba(34,211,238,0.12); }
.pc-body { flex: 1; min-height: 0; display: flex; }
.pc-body > * { flex: 1; min-width: 0; }
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
