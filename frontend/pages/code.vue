<template>
  <div class="code-page">
    <CodeSidebar />

    <main class="code-main">
      <!-- 顶部标题 -->
      <div class="code-header">
        <div class="header-left">
          <h2 class="header-title">AI 智能代码生成</h2>
          <p class="header-sub">上传设计稿，AI 自动识别并生成高质量代码</p>
        </div>
      </div>

      <!-- 上传 + 配置 -->
      <section class="upload-section">
        <div class="upload-section-inner">
          <!-- 左栏：上传 -->
          <div class="upload-left">
            <div class="glow-section">
              <div class="section-label"><span class="label-num">1</span> <span class="section-label-text">上传设计稿</span></div>
              <UploadZone
                @files-selected="handleFilesSelected"
                :disabled="isMaxReached || isGenerating"
                :hint="uploadHint"
                :images="images"
                :isGenerating="isGenerating"
                @remove="removeImage"
                @clear="clearAll"
                @reorder="images = $event"
              />
            </div>
          </div>

          <!-- 右栏：配置 -->
          <div class="upload-right">
            <div class="glow-section">
              <ConfigPanel v-model="config" :canGenerate="images.length > 0 && allUploaded" @generate="handleGenerate" />
            </div>
          </div>
        </div>
      </section>

      <!-- 生成流程 -->
      <div class="glow-section">
        <FlowSteps :activeStep="activeStep" />
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import CodeSidebar from '~/components/code/CodeSidebar.vue'
import UploadZone from '~/components/upload/UploadZone.vue'
import ConfigPanel from '~/components/code/ConfigPanel.vue'
import FlowSteps from '~/components/code/FlowSteps.vue'
import { commonApi } from '~/api/common'
import { useGeneration } from '~/composables/useGeneration'
import { fileToBase64, handleCopy, handleDownload } from '~/utils/index.js'
import { ElMessage } from 'element-plus'

// ═══ 框架映射 ═══
const targetMap = { Flutter: 'flutter', React: 'react', Vue: 'vue3' }
const langMap = {
  flutter: { label: 'Dart', lang: 'dart' },
  react:   { label: 'TypeScript', lang: 'typescript' },
  vue3:    { label: 'Vue', lang: 'html' },
}
const deviceMap = { mobile: '手机端', desktop: '桌面端', tablet: '平板端' }
const maxImages = 3

// ═══ 配置 ═══
const config = ref({
  framework: 'Flutter',
  platform: 'mobile',
  options: ['responsive'],
  advanced: [],
  componentLib: '',
})

// ═══ 生成 composable ═══
const {
  template, taskStatus, taskProgress, taskCurrentStep,
  saveActiveTask, clearActiveTask, restoreTask,
  sseStatus, connectSSE, disconnectSSE, isAvailable,
  isLogin, credits, openLoginModal, startGenerating, finishGenerating, refreshCredits,
} = useGeneration()

// ═══ 本地状态 ═══
const images = ref([])
const hasGenerated = ref(false)

// ── 计算属性 ──
const activeStep = computed(() => {
  if (taskProgress.value === 0) return 0
  if (taskProgress.value < 20) return 0
  if (taskProgress.value < 50) return 1
  if (taskProgress.value < 75) return 2
  if (taskProgress.value < 100) return 3
  return 4
})

const isGenerating = computed(() =>
  taskStatus.value === 'running' || taskStatus.value === 'pending'
)

const target = computed(() => targetMap[config.value.framework] || 'flutter')

const langLabel = computed(() => langMap[target.value]?.label || 'Dart')
const codeLanguage = computed(() => langMap[target.value]?.lang || 'dart')

const deviceLabel = computed(() => deviceMap[config.value.platform] || '')

const allUploaded = computed(() =>
  images.value.length > 0 && images.value.every(img => !!img.cosUrl)
)

const isMaxReached = computed(() => images.value.length >= maxImages)

const uploadHint = computed(() => {
  if (isGenerating.value) return '正在生成中，暂不支持操作图片'
  if (isMaxReached.value) return `已达到上传上限（${maxImages}/${maxImages}），请先移除再上传`
  return `支持 PNG、JPG、JPEG，最多 ${maxImages} 张，单张不超过 10MB`
})

const isGenerateDisabled = computed(() =>
  images.value.length === 0 || !allUploaded.value || isGenerating.value
)

// ═══ 文件操作 ═══
const handleFilesSelected = (files) => {
  if (isGenerating.value) return
  if (!isLogin.value) {
    openLoginModal()
    return
  }
  if (images.value.length >= maxImages) {
    ElMessage.warning(`最多上传 ${maxImages} 张图片`)
    return
  }
  addImages(files)
}

const addImages = async (files) => {
  const imageFiles = files.filter(f => f.type.startsWith('image/'))
  const remaining = Math.min(maxImages - images.value.length, imageFiles.length)
  const toAdd = imageFiles.slice(0, remaining)

  for (const file of toAdd) {
    const { preview, width, height } = await fileToPreview(file)
    const idx = images.value.length
    images.value.push({
      file,
      preview,
      naturalWidth: width,
      naturalHeight: height,
      cosUrl: '',
      uploading: true,
      uploadError: '',
      description: '',
    })
    uploadOneImage(idx, file)
  }
}

function fileToPreview(file) {
  return new Promise((resolve) => {
    const reader = new FileReader()
    reader.onload = (e) => {
      const dataUrl = e.target.result
      // 读取图片自然尺寸
      const img = new Image()
      img.onload = () => {
        resolve({ preview: dataUrl, width: img.naturalWidth, height: img.naturalHeight })
        URL.revokeObjectURL(img.src)
      }
      img.src = dataUrl
    }
    reader.readAsDataURL(file)
  })
}

const uploadOneImage = async (idx, file) => {
  try {
    const base64 = await fileToBase64(file)
    const base64Data = base64.split(',')[1]
    const res = await commonApi.uploadImage(base64Data, file.name)
    if (res.code === 200 && res.data?.url) {
      images.value[idx].cosUrl = res.data.url
      images.value[idx].uploading = false
    } else {
      throw new Error(res.message || '上传失败')
    }
  } catch (e) {
    images.value[idx].uploading = false
    images.value[idx].uploadError = e.message || '上传失败'
    ElMessage.error(`"${file.name}" 上传失败: ${e.message}`)
  }
}

const removeImage = (idx) => {
  if (isGenerating.value) return
  images.value.splice(idx, 1)
}

const clearAll = () => {
  if (isGenerating.value) return
  images.value = []
  // 上一轮已完成 → 重置为 idle，避免按钮卡在"生成完成"
  if (taskStatus.value === 'success' || taskStatus.value === 'failed') {
    taskStatus.value = 'idle'
  }
}

// ═══ 生成逻辑 ═══
const handleGenerate = async () => {
  if (!isLogin.value) {
    ElMessage.warning('请先登录后再生成代码')
    openLoginModal()
    return
  }
  if (images.value.length === 0) {
    ElMessage.warning('请先上传设计稿')
    return
  }
  if (!allUploaded.value) {
    ElMessage.warning('请等待图片上传完成')
    return
  }
  if (credits.value < images.value.length) {
    ElMessage.warning(`积分不足，需要 ${images.value.length} 积分，当前剩余 ${credits.value} 积分`)
    return
  }

  startGenerating()

  const payload = {
    target: target.value,
    platform: config.value.platform,
    options: config.value.options,
    advanced: config.value.advanced,
    component_lib: config.value.options.includes('component') ? config.value.componentLib : '',
    images: images.value.map((img, i) => ({
      url: img.cosUrl,
      desc: img.description || '',
      sort_order: i + 1,
    })),
  }

  try {
    const result = await commonApi.generateUi(payload)
    if (!result || !result.data || !result.data.task_id) {
      ElMessage.error('AI 启动失败，请稍后重试')
      return
    }

    const taskId = result.data.task_id
    saveActiveTask(taskId, target.value)

    // 连接 SSE
    if (isAvailable()) {
      await connectSSE(taskId)
    }
  } catch (error) {
    console.error('[code.vue] 生成失败:', error)
    ElMessage.error(error.message || '生成失败，请稍后重试')
    finishGenerating()
  }
}

// 任务状态变化时同步全局 isGenerating
watch(taskStatus, (status) => {
  if (status === 'idle' || status === 'success' || status === 'failed') {
    finishGenerating()
  }
  if (status === 'success' || status === 'failed') {
    refreshCredits()
  }
})

// ═══ 代码操作 ═══
const handleCopyCode = () => {
  handleCopy(template.templateCode)
}

const handleDownloadCode = () => {
  handleDownload(template.templateCode, codeLanguage.value)
}

// ═══ 任务恢复 ═══
const doRestoreTask = async () => {
  const restored = await restoreTask()
  if (!restored) return

  // 回显图片
  if (restored.images && restored.images.length > 0) {
    images.value = restored.images.map((img, idx) => ({
      file: null,
      preview: img.url,
      cosUrl: img.url,
      uploading: false,
      uploadError: '',
      description: img.desc || '',
    }))
  }

  if (restored.status === 'success' || restored.status === 'failed') {
    hasGenerated.value = true
  }
}

// ═══ 监测 SSE 完成 ═══
watch(sseStatus, (val) => {
  if (val === 'done' || val === 'idle') {
    if (template.templateCode && template.templateCode.length > 0) {
      hasGenerated.value = true
    }
  }
})

watch(() => template.templateCode, (code) => {
  if (code && code.length > 0 && !hasGenerated.value) {
    hasGenerated.value = true
  }
})

// ═══ 生命周期 ═══
// onMounted(() => {
//   doRestoreTask()
// })
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

/* ── 科幻发光区块 ── */
.glow-section {
  position: relative;
  border-radius: 16px;
  border: 1px solid rgba(0, 255, 255, 0.08);
  background: rgba(10, 14, 23, 0.7);
  backdrop-filter: blur(16px);
  padding: 20px 24px;
  transition: all 0.4s;
  overflow: hidden;
}
.glow-section::before {
  content: '';
  position: absolute;
  top: -1px;
  left: -1px;
  right: -1px;
  bottom: -1px;
  border-radius: 17px;
  background: linear-gradient(135deg, rgba(0, 255, 255, 0.12), transparent 30%, transparent 70%, rgba(255, 0, 255, 0.08));
  z-index: 0;
  pointer-events: none;
}
.glow-section:hover {
  border-color: rgba(0, 255, 255, 0.18);
  box-shadow: 0 0 30px rgba(0, 255, 255, 0.04), inset 0 0 30px rgba(0, 255, 255, 0.02);
}
.glow-section > * { position: relative; z-index: 1; }

/* 区块标题 */
.section-label {
  font-size: 14px;
  font-weight: 600;
  color: #60a5fa;
  text-transform: uppercase;
  letter-spacing: 1.5px;
  margin-bottom: 16px;
  position: relative;
  display: inline-block;
  padding-left: 14px;
}
.section-label::before {
  content: '';
  position: absolute;
  left: 0;
  top: 2px;
  bottom: 2px;
  width: 2px;
  background: linear-gradient(135deg, #60a5fa, #818cf8);
  border-radius: 1px;
}

/* 顶部标题 */
.code-header {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  justify-content: space-between;
  padding: 20px 24px;
  border-radius: 16px;
}
.header-title {
  font-size: 28px;
  font-weight: 800;
  color: #e2e8f0;
  letter-spacing: -0.3px;
}
.header-sub {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.45);
  margin-top: 4px;
}

.upload-section-inner { display: flex; gap: 24px; }

.upload-left {
  flex: 1.5;
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.upload-left .glow-section,
.upload-right .glow-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 20px;
}


.label-num {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 8px;
  font-size: 12px;
  background: #60a5fa;
  color: #ffffff;
  font-weight: 700;
  margin-right: 6px;
}

.section-label-text { font-weight: 700; color: #ffffff; }

.upload-right {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* ── 结果展示区 ── */
.result-section {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  padding: 20px;
}
.result-inner {
  display: flex;
  gap: 24px;
  flex: 1;
  min-height: 480px;
}
.result-editor {
  flex: 1.55;
  min-width: 0;
  display: flex;
  flex-direction: column;
  background: rgba(15, 20, 30, 0.6);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(0, 255, 255, 0.18);
  border-radius: 18px;
  overflow: hidden;
}
.result-preview {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  background: rgba(15, 20, 30, 0.6);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(0, 255, 255, 0.18);
  border-radius: 18px;
  overflow: hidden;
}
.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 18px;
  background: rgba(0, 0, 0, 0.3);
  border-bottom: 1px solid rgba(0, 255, 255, 0.1);
  font-size: 14px;
  font-weight: 600;
  color: #00cfff;
}
.result-header i { font-size: 15px; margin-right: 6px; }
.result-actions { display: flex; gap: 6px; }

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
  font-family: inherit;
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

@media (max-width: 1100px) {
  .upload-section-inner { flex-direction: column; }
  .upload-right { width: 100%; }
  .result-inner { flex-direction: column; }
  .result-editor { height: 480px; }
  .result-preview { height: auto; }
}

/* ── 生成结果弹窗 ── */
:deep(.result-dialog) {
  .el-dialog__body {
    padding: 0;
    height: calc(90vh - 120px);
    overflow: hidden;
  }
}

:deep(.result-dialog) .result-inner {
  height: 100%;
  padding: 20px 24px;
  gap: 24px;
}

.dialog-title {
  font-size: 16px;
  font-weight: 600;
  color: #00cfff;
  display: flex;
  align-items: center;
  gap: 8px;
}
.dialog-title i { font-size: 15px; }
</style>
