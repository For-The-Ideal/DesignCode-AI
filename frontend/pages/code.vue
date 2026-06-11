<!-- pages/code.vue -->
<template>
  <div class="page-container">
    <main class="main-content">
      <div class="page-header">
        <h1>AI 代码生成</h1>
        <p>上传设计稿 · AI 智能识别 · 秒级生成高质量代码</p>
      </div>

      <!-- ═══ 上传区 ═══ -->
      <UploaderImage @generated="onGenerated" :isBusy="isBusy" />

      <!-- ═══ 未生成时的占位 ═══ -->
      <div v-if="!hasGenerated && !isBusy" class="empty-state">
        <div class="empty-icon">
          <i class="fas fa-code"></i>
        </div>
        <h3>等待生成代码</h3>
        <p>上传设计稿并点击「开始生成代码」，AI 将自动为你生成高质量代码</p>
        <div class="feature-hints">
          <div class="hint-item">
            <i class="fas fa-image"></i>
            <span>支持多张设计稿上传</span>
          </div>
          <div class="hint-item">
            <i class="fas fa-edit"></i>
            <span>为每张图片添加描述让 AI 更懂你</span>
          </div>
          <div class="hint-item">
            <i class="fas fa-mobile-alt"></i>
            <span>一键生成 Flutter / React / Vue 代码</span>
          </div>
        </div>
      </div>

      <!-- ═══ 任务进行中（等待 SSE 数据） ═══ -->
      <div v-if="isBusy && !hasGenerated" class="busy-state">
        <div class="busy-spinner"><i class="fas fa-spinner fa-spin"></i></div>
        <h3>AI 生成中...</h3>
        <div class="busy-progress-bar">
          <div class="busy-progress-fill" :style="{ width: taskProgress + '%' }"></div>
        </div>
        <p class="busy-step">{{ taskCurrentStep || '正在连接...' }}</p>
        <p class="busy-hint">生成完成后下方将展示代码与预览，请勿关闭页面</p>
      </div>

      <!-- ═══ 生成结果展示 ═══ -->
      <template v-if="hasGenerated">
        <div class="core-layout">
          <div class="panel editor-panel">
            <div class="panel-header">
              <div class="panel-title">
                <i class="fas fa-code"></i>
                <span>{{ generatedLang }} 代码</span>
              </div>
              <div class="panel-actions">
                <button class="act-btn" @click="handleCopy"><i class="fas fa-copy"></i><span>复制</span></button>
                <button class="act-btn" @click="handleFormat"><i class="fas fa-magic"></i><span>格式化</span></button>
                <button class="act-btn ghost" @click="handleDownload"><i class="fas fa-download"></i><span>下载</span></button>
              </div>
            </div>
            <div class="panel-body editor-body">
              <CodeEditor
                v-model="template.templateCode"
                :language="codeLanguage"
                :readonly="false"
                :auto-scroll="sseStatus === 'streaming'"
                height="700px"
                placeholder="// AI 生成的代码将在这里展示..."
              />
            </div>
          </div>

          <div class="panel preview-panel">
            <div class="panel-header">
              <div class="panel-title">
                <i class="fas fa-mobile-alt"></i>
                <span>实时预览</span>
              </div>
              <span class="device-badge" v-if="generatedFramework">iPhone 15 Pro</span>
            </div>
            <div class="panel-body preview-body">
              <FlutterTemplate :html="template.previewCode" :showBottomNav="false" />
            </div>
          </div>
        </div>
      </template>
    </main>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted } from 'vue'
import CodeEditor from '~/components/code/CodeEditor.vue'
import FlutterTemplate from '@/components/common/FlutterTemplate.vue'
import UploaderImage from '~/components/upload/UploaderImage.vue'
import { useGeneration } from '~/composables/useGeneration'

import { ElMessage } from 'element-plus'
// ═══ 生成流程 composable（SSE + 流式渲染 统一入口）═══
const {
  template,
  sseStatus,
  connectSSE,
  disconnectSSE,
  cleanup,
  isBusy,
  taskStatus,
  taskProgress,
  taskCurrentStep,
  restoreTask,
} = useGeneration()

// ═══ 页面状态 ═══
const hasGenerated = ref(false)
const generatedLang = ref('Dart')
const generatedFramework = ref('')
const codeLanguage = ref('dart')

const langMap = {
  flutter: { label: 'Dart', lang: 'dart' },
  react: { label: 'TypeScript', lang: 'typescript' },
  vue: { label: 'Vue', lang: 'html' },
}

// ═══ 监听 SSE 流式输出 → 驱动视图更新 ═══

/** 监听代码内容变化，自动标记为已生成 */
watch(() => template.templateCode, (code) => {
  if (code && code.length > 0 && !hasGenerated.value) {
    hasGenerated.value = true
  }
})

/** 监听 SSE 状态变化 */
watch(sseStatus, (val) => {
  if (val === 'done' && template.templateCode && template.templateCode.length > 0) {
    hasGenerated.value = true
  }
})

// ═══ 上传组件回调 ═══

const onGenerated = async (result) => {
  hasGenerated.value = true
  generatedFramework.value = result.framework
  if (result.framework && langMap[result.framework]) {
    generatedLang.value = langMap[result.framework].label
    codeLanguage.value = langMap[result.framework].lang
  }
}

// ═══ 生命周期 ═══

onMounted(async () => {
  // 检查是否有之前未完成的任务，恢复进度
  const restored = await restoreTask()
  if (restored) {
    if (restored.status === 'pending' || restored.status === 'running') {
      // 有任务正在执行中，复现生成的框架配置
      if (restored.framework && langMap[restored.framework]) {
        generatedFramework.value = restored.framework
        generatedLang.value = langMap[restored.framework].label
        codeLanguage.value = langMap[restored.framework].lang
      }
      // SSE 已由 restoreTask 内部重连，等待推送即可
      ElMessage.info('检测到正在执行中的任务，已恢复进度监听')
    } else if (restored.status === 'success') {
      // 已完成的任务，展示结果
      hasGenerated.value = true
      if (restored.framework && langMap[restored.framework]) {
        generatedFramework.value = restored.framework
        generatedLang.value = langMap[restored.framework].label
        codeLanguage.value = langMap[restored.framework].lang
      }
    } else if (restored.status === 'failed') {
      ElMessage.error('上轮任务执行失败，请重新生成')
    }
  }
  // 没有已存在的任务 → 保持空状态等待用户操作
})

onUnmounted(() => {
  cleanup()
  disconnectSSE()
})

// ═══ 工具栏操作 ═══

const handleCopy = async () => {
  try {
    await navigator.clipboard.writeText(template.templateCode)
    ElMessage.success('已复制到剪贴板')
  } catch {
    ElMessage.error('复制失败')
  }
}

const handleFormat = () => {
  let indent = 0
  template.templateCode = template.templateCode.split('\n').map(line => {
    const t = line.trim()
    if (!t) return ''
    if (/^[})]/.test(t)) indent = Math.max(0, indent - 1)
    const r = '  '.repeat(indent) + t
    if (/[({]$/.test(t)) indent++
    return r
  }).join('\n')
  ElMessage.success('代码已格式化')
}

const handleDownload = () => {
  const ext = codeLanguage.value === 'dart' ? 'dart' : codeLanguage.value === 'typescript' ? 'tsx' : 'vue'
  const blob = new Blob([template.templateCode], { type: 'text/plain' })
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = `generated.${ext}`
  a.click()
  URL.revokeObjectURL(a.href)
  ElMessage.success('已下载')
}
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #0a0a0f 0%, #0f1a1f 100%);
}
.main-content {
  max-width: 1600px;
  margin: 0 auto;
  padding: 40px 24px 60px;
}

.page-header {
  text-align: center;
  margin-bottom: 32px;
}
.page-header h1 {
  font-size: 40px;
  font-weight: 800;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 10px;
}
.page-header p {
  color: #6b7280;
  font-size: 15px;
}

/* ═══ SSE 状态栏 ═══ */
.sse-status-bar {
  display: flex;
  justify-content: center;
  padding: 12px 0;
  margin-top: -8px;
}
.sse-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 20px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 500;
}
.sse-badge.connecting {
  background: rgba(0, 255, 255, 0.08);
  border: 1px solid rgba(0, 255, 255, 0.25);
  color: #00cfff;
}
.sse-badge.streaming {
  background: rgba(0, 255, 100, 0.08);
  border: 1px solid rgba(0, 255, 100, 0.3);
  color: #00ff88;
}
.sse-badge.error {
  background: rgba(255, 60, 60, 0.08);
  border: 1px solid rgba(255, 60, 60, 0.3);
  color: #ff5555;
}

/* ═══ 空状态 ═══ */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  text-align: center;
  background: rgba(15, 20, 30, 0.35);
  border: 1px solid rgba(0, 255, 255, 0.1);
  border-radius: 24px;
  margin-top: 24px;
}
.empty-icon {
  width: 88px;
  height: 88px;
  border-radius: 50%;
  background: linear-gradient(135deg, rgba(0,255,255,0.12), rgba(255,0,255,0.12));
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24px;
}
.empty-icon i {
  font-size: 36px;
  color: #00ffff;
}
.empty-state h3 {
  font-size: 22px;
  color: #ccc;
  margin-bottom: 10px;
}
.empty-state > p {
  color: #6b7280;
  font-size: 14px;
  margin-bottom: 36px;
}
.feature-hints {
  display: flex;
  gap: 32px;
  flex-wrap: wrap;
  justify-content: center;
}
.hint-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 20px;
  background: rgba(0,255,255,0.05);
  border: 1px solid rgba(0,255,255,0.12);
  border-radius: 14px;
  color: #888;
  font-size: 13px;
}
.hint-item i {
  color: #00cfff;
  font-size: 16px;
}

/* ═══ 进行中状态 ═══ */
.busy-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
  background: rgba(15, 20, 30, 0.35);
  border: 1px solid rgba(0, 255, 255, 0.1);
  border-radius: 24px;
  margin-top: 24px;
}
.busy-spinner {
  font-size: 36px;
  color: #00ffff;
  margin-bottom: 16px;
}
.busy-state h3 {
  font-size: 20px;
  color: #ccc;
  margin-bottom: 16px;
}
.busy-progress-bar {
  width: 280px;
  max-width: 80%;
  height: 6px;
  background: rgba(255,255,255,0.08);
  border-radius: 3px;
  overflow: hidden;
  margin-bottom: 12px;
}
.busy-progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #00ffff, #00ff88);
  border-radius: 3px;
  transition: width 0.5s ease;
}
.busy-step {
  font-size: 13px;
  color: #00cfff;
  margin-bottom: 8px;
}
.busy-hint {
  font-size: 12px;
  color: #6b7280;
}

/* ═══ 核心区 ═══ */
.core-layout {
  display: flex;
  gap: 24px;
  margin-top: 24px;
  margin-bottom: 32px;
}

.panel {
  background: rgba(15, 20, 30, 0.6);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(0, 255, 255, 0.18);
  border-radius: 18px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 18px;
  background: rgba(0, 0, 0, 0.3);
  border-bottom: 1px solid rgba(0, 255, 255, 0.1);
  flex-shrink: 0;
}
.panel-title {
  display: flex;
  align-items: center;
  gap: 9px;
  font-size: 14px;
  font-weight: 600;
  color: #00cfff;
}
.panel-title i { font-size: 15px; }
.panel-actions { display: flex; gap: 6px; }

.act-btn {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 5px 12px;
  background: rgba(0, 255, 255, 0.07);
  border: 1px solid rgba(0, 255, 255, 0.18);
  border-radius: 18px;
  color: #00cfff;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
  font-family: system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif;
}
.act-btn i { font-size: 12px; }
.act-btn:hover { background: rgba(0, 255, 255, 0.14); border-color: #00ffff; }
.act-btn.ghost { background: rgba(255,255,255,0.04); border-color: rgba(255,255,255,0.12); color: #aaa; }
.act-btn.ghost:hover { background: rgba(255,255,255,0.08); color: #fff; }

.device-badge {
  font-size: 11px;
  color: #6b7280;
  letter-spacing: 0.3px;
  background: rgba(255,255,255,0.04);
  padding: 3px 10px;
  border-radius: 20px;
}

.panel-body { overflow: hidden; flex: 1; min-height: 0; }

.editor-panel { flex: 1.55; min-width: 0; }
.editor-body { flex: 1; min-height: 0; }

.preview-panel { flex: 1; min-width: 0; }
.preview-body {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding: 16px;
  background: rgba(0, 0, 0, 0.15);
  overflow-y: auto;
}

@media (max-width: 1100px) {
  .core-layout { flex-direction: column; height: auto; }
  .editor-panel { height: 480px; }
  .preview-panel { height: auto; }
  .preview-body { justify-content: center; }
  .feature-hints { flex-direction: column; align-items: center; }
}

@media (max-width: 768px) {
  .main-content { padding: 20px 14px 40px; }
  .page-header h1 { font-size: 28px; }
  .editor-panel { height: 380px; }
  .act-btn span { display: none; }
}
</style>
