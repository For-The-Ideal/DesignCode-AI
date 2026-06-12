<template>
  <div class="upload-section">
    <div class="tech-card">
      <div class="card-title">
        <i class="fas fa-cloud-upload-alt"></i>
        <span>上传设计稿</span>
      </div>
      <div class="upload-content">
        <div 
          class="dropzone" 
          :class="{ 'dropzone-disabled': isMaxReached, 'drag-over': isDragging }"
          @click="!isMaxReached && triggerFileInput()" 
          @dragover.prevent="!isMaxReached && (isDragging = true)" 
          @dragleave.prevent="isDragging = false" 
          @drop.prevent="!isMaxReached && handleDrop($event)">
          <div class="dropzone-icon">
            <i class="fas fa-file-image"></i>
          </div>
          <p v-if="!isMaxReached">点击或拖拽上传设计稿</p>
          <p v-else class="max-warning">已达到上传上限 ({{ maxLen }}/{{ maxLen }})</p>
          <span>支持 PNG、JPG、JPEG，最多 {{ maxLen }} 张</span>
          <input type="file" ref="fileInput" :disabled="isMaxReached" multiple 
          accept="image/*" style="display: none" @change="handleFileSelect" />
        </div>

        <!-- 预览区域 -->
        <div class="preview-area" v-if="images.length > 0">
          <div class="preview-header">
            <span><i class="fas fa-images"></i> 已上传 ({{ images.length }}/{{ maxLen }})</span>
            <span v-if="uploadedCount < images.length" class="upload-progress">{{ uploadedCount }}/{{ images.length }} 已上传完成</span>
            <button class="clear-btn" @click="clearAll">清空全部</button>
          </div>
          <div class="preview-grid">
            <div v-for="(img, idx) in images" :key="img.id" class="preview-item" :class="{ 'upload-error': img.uploadError }">
              <img :src="img.preview" alt="预览" />
      
              <!-- 上传状态 -->
              <div v-if="img.uploading" class="upload-status uploading">
                <i class="fas fa-spinner fa-spin"></i>
                <span>上传中...</span>
              </div>
              <div v-else-if="img.cosUrl" class="upload-status done">
                <i class="fas fa-check-circle"></i>
                <span>已上传</span>
              </div>
              <div v-else-if="img.uploadError" class="upload-status failed">
                <i class="fas fa-exclamation-circle"></i>
                <span>上传失败</span>
              </div>
              <div class="preview-remove" @click="removeImage(idx)">×</div>
              <!-- 描述显示区域 -->
              <div class="image-description" v-if="img.description" @click="openDescModal(idx)">
                <div class="desc-content">
                  <i class="fas fa-quote-left"></i>
                  <span class="desc-text">{{ truncateText(img.description, 40) }}</span>
                </div>
                <div class="desc-footer">
                  <span class="desc-type">{{ img.type || '设计稿' }}</span>
                  <span class="desc-edit"><i class="fas fa-pen"></i> 编辑</span>
                </div>
              </div>
              <div class="image-description empty" v-else @click="openDescModal(idx)">
                <i class="fas fa-plus-circle"></i>
                <span>添加描述</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 配置卡片 -->
    <div class="tech-card config-card">
      <div class="card-title">
        <i class="fas fa-sliders-h"></i>
        <span>生成配置</span>
      </div>

      <div class="framework-selector">
        <div class="framework-btn" :class="{ active: framework === 'flutter' }" @click="framework = 'flutter'">
          <i class="fab fa-flutter"></i> Flutter
        </div>
        <div class="framework-btn" :class="{ active: framework === 'react' }" @click="framework = 'react'">
          <i class="fab fa-react"></i> React
        </div>
        <div class="framework-btn" :class="{ active: framework === 'vue3' }" @click="framework = 'vue3'">
          <i class="fab fa-vuejs"></i> Vue
        </div>
      </div>

      <div class="slider-container">
        <div class="slider-label">
          <span>质量要求</span>
          <span class="quality-value">{{ qualityValue }}</span>
        </div>
        <div class="slider-wrapper">
          <input type="range" v-model.number="qualityValue" min="60" max="100" step="1" class="quality-slider" />
          <div class="slider-marks">
            <span v-for="mark in [60, 65, 70, 75, 80, 85, 90, 95, 100]" 
                  :key="mark" 
                  class="slider-mark"
                  :class="{ active: qualityValue >= mark }"
                  @click="qualityValue = mark">
              {{ mark }}
            </span>
          </div>
        </div>
      </div>
    </div>
    
      <!-- AI 智能分析提示 -->
      <div v-if="images.length > 0" class="ai-tip">
        <i class="fas fa-robot"></i>
        <span>AI 将结合图片和描述生成更精准的代码，描述越详细效果越好</span>
      </div>
    </div>

    

    <!-- 生成按钮 -->
    <button class="generate-btn" @click="generateCode" :disabled="images.length === 0 || !allUploaded || generating || (taskProgress > 0 && taskProgress < 100)">
      <i v-if="taskProgress > 0 && taskProgress < 100" class="fas fa-hourglass-half"></i>
      <i v-else-if="generating" class="fas fa-spinner fa-spin"></i>
      <i v-else class="fas fa-play"></i>
      {{ (taskProgress > 0 && taskProgress < 100) ? '任务执行中...' : generating ? '正在生成...' : '开始生成' }}
    </button>

    <!-- 描述编辑弹窗 -->
    <DescEditorModal
      ref="descEditorRef"
      :images="images"
      :focus-index="editingIndex"
      @save="handleSaveDescriptions"
      @close="closeDescModal"
    />
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import DescEditorModal from './DescEditorModal.vue'
import { commonApi } from "~/api/common.js"
import { ElMessage } from 'element-plus'
import { useSSE } from '~/composables/useSSE'
import { useGeneration } from '~/composables/useGeneration'
import { fileToBase64 } from '~/utils/index.js'

const emit = defineEmits(['generated'])

const props = defineProps({
  /** 外部传入的已上传图片列表（页面刷新恢复时用），格式: [{ url: 'cos_url', desc: '' }] */
  initialImages: {
    type: Array,
    default: () => [],
  },
})

// 图片列表
const images = ref([])
const isDragging = ref(false)
const fileInput = ref(null)
// 配置
const framework = ref('flutter')
const qualityValue = ref(90)
const maxLen = ref(5) // 最多5张
// 评分数据
const score = ref(0)
const scoreDimensions = ref([])

const { isAvailable, isAlive, status, connect } = useSSE()
const { saveActiveTask, taskProgress } = useGeneration()
// 本地按钮锁（点击→SSE connected 之间防连点）
const generating = ref(false)

// 描述弹窗
const descEditorRef = ref(null)
const editingIndex = ref(null)


// 计算是否达到上限
const isMaxReached = computed(() => images.value.length >= maxLen.value)

// 已上传完成的图片数
const uploadedCount = computed(() => images.value.filter(img => !!img.cosUrl).length)

// 是否所有图片都上传完成
const allUploaded = computed(() => images.value.length > 0 && images.value.every(img => !!img.cosUrl))


// 触发文件选择
const triggerFileInput = () => {
  if (!isMaxReached.value) {
    fileInput.value?.click()
  }
}

// 处理文件选择
const handleFileSelect = (e) => {
  const files = Array.from(e.target.files)
  addImages(files)
  if (fileInput.value) fileInput.value.value = ''
}

// 处理拖拽
const handleDrop = (e) => {
  isDragging.value = false
  const files = Array.from(e.dataTransfer.files)
  addImages(files)
}

// 添加图片并立即上传到 COS
const addImages = async (files) => {
  const imageFiles = files.filter(f => f.type.startsWith('image/'))
  const remaining = maxLen.value - images.value.length
  const toAdd = imageFiles.slice(0, remaining)
  
  for (const file of toAdd) {
    const preview = await fileToPreview(file)
    const idx = images.value.length
    images.value.push({
      file,
      preview,
      cosUrl: '',
      uploading: true,
      uploadError: '',
      description: '',
      type: detectImageType(file.name)
    })
    // 立即上传到 COS
    uploadOneImage(idx, file)
  }
}

// 单张图片转预览 dataURL
function fileToPreview(file) {
  return new Promise((resolve) => {
    const reader = new FileReader()
    reader.onload = (e) => resolve(e.target.result)
    reader.readAsDataURL(file)
  })
}

// 单张图片上传到 COS
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

// 检测图片类型
const detectImageType = (filename) => {
  const name = filename.toLowerCase()
  if (name.includes('main') || name.includes('首页')) return '主设计稿'
  if (name.includes('detail') || name.includes('细节')) return '细节说明'
  if (name.includes('flow') || name.includes('流程')) return '流程图'
  if (name.includes('spec') || name.includes('标注')) return '标注规范'
  return '设计稿'
}

// 移除图片
const removeImage = (idx) => {
  images.value.splice(idx, 1)
}

// 清空所有
const clearAll = () => {
    images.value = []
    editingIndex.value = null
}

// 打开描述弹窗（可指定索引）
const openDescModal = (idx = null) => {
  editingIndex.value = idx
  descEditorRef.value?.open()
}

// 关闭描述弹窗
const closeDescModal = () => {
  descEditorRef.value?.close()
  editingIndex.value = null
}

// 保存描述
const handleSaveDescriptions = (updatedDescriptions) => {
  images.value.forEach((img, idx) => {
    img.description = updatedDescriptions[idx] || ''
  })
}

// 截取文本
const truncateText = (text, maxLen) => {
  if (!text) return ''
  return text.length > maxLen ? text.slice(0, maxLen) + '...' : text
}

const generateCode = async () => {
  if (generating.value) return

  // SSE 连接状态前置校验
  if (!isAvailable()) {
    ElMessage.error('SSE 连接不可用，请刷新页面后重试')
    return
  }

  generating.value = true
  try {
    // 直接从组件状态取已上传完成的 COS URL
    const allUploaded = images.value.every(img => !!img.cosUrl)
    if (!allUploaded) {
      throw new Error('存在未上传完成的图片，请稍后重试')
    }
    const payload = {
      target: framework.value,
      images: images.value.map(img => ({
        url: img.cosUrl,
        desc: img.description || ''
      })),
      quality: qualityValue.value,
    }

    console.log('[Uploader] payload:', payload)
    console.log('[Uploader] SSE isAvailable:', isAvailable())

    // 通知后端开始生成 → 拿到 task_id → 连接 SSE
    const result = await commonApi.generateUi(payload)
    console.log('[Uploader] result:', result)

    if (!result || !result.data || !result.data.task_id) {
      generating.value = false
      ElMessage.error('AI 启动失败，请稍后重试')
      return
    }

    // 立即保存任务信息到 localStorage（先于 connect 执行，确保刷新可恢复）
    saveActiveTask(result.data.task_id, framework.value)

    // 通知父组件（code.vue）框架类型
    emit('generated', { framework: framework.value })

    // 再连接 SSE（可能较慢，不阻塞 localStorage 写入）
    console.log('[Uploader] 连接 SSE, task_id:', result.data.task_id)
    await connect(result.data.task_id)

    // 连接成功后,转换按钮状态: "正在生成..." → "任务执行中..."
    generating.value = false

  } catch (error) {
    console.error('[Uploader] 生成失败:', error)
    ElMessage.error(error.message || '生成失败，请稍后重试')
  } finally {
    generating.value = false
  }
}

// 外部传入的已上传图片 → 回显到预览区
watch(() => props.initialImages, (imgs) => {
  if (imgs && imgs.length > 0) {
    images.value = imgs.map((img, idx) => ({
      file: null,
      preview: img.url,           // COS URL 直接当预览图
      cosUrl: img.url,
      uploading: false,
      uploadError: '',
      description: img.desc || '',
      type: '设计稿',
      id: `restored_${idx}`,
    }))
  }
}, { immediate: true })

// 监听 SSE 状态: 空闲时确保按钮可点击
watch(status, (val) => {
  if (val === 'idle' || val === 'error' || val === 'maxRetries') {
    generating.value = false
  }
})

</script>

<style scoped lang="scss">
/* 科技卡片 */
.tech-card {
  background: rgba(15, 20, 30, 0.6);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(0, 255, 255, 0.2);
  border-radius: 24px;
  padding: 24px;
  margin-bottom: 24px;
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 10px;
  color: #00ffff;
}

/* 上传区 + 预览区 左右布局 */
.upload-content {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}

/* 拖拽上传区 */
.dropzone {
  width: 260px;
  min-height: 260px;
  flex-shrink: 0;
  border: 2px dashed rgba(0, 255, 255, 0.4);
  border-radius: 20px;
  padding: 32px 20px;
  text-align: center;
  background: rgba(0, 0, 0, 0.3);
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.dropzone.drag-over {
  border-color: #00ffff;
  background: rgba(0, 255, 255, 0.05);
}

/* 达到上限时的禁用样式 */
.dropzone-disabled {
  opacity: 0.5;
  cursor: not-allowed;
  border-color: rgba(255, 0, 0, 0.4);
  background: rgba(0, 0, 0, 0.5);
}

.dropzone-disabled:hover {
  border-color: rgba(255, 0, 0, 0.4);
  transform: none;
}

.max-warning {
  color: #ff4444;
  font-weight: 500;
}

.dropzone-icon {
  width: 56px;
  height: 56px;
  background: linear-gradient(135deg, rgba(0, 255, 255, 0.2), rgba(255, 0, 255, 0.2));
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
}

.dropzone-icon i {
  font-size: 24px;
  color: #00ffff;
}

.dropzone p {
  margin-bottom: 8px;
}

.dropzone span {
  font-size: 12px;
  color: #666;
}

/* 预览区域 */
.preview-area {
  flex: 1;
  min-width: 0;
  max-height: 420px;
  overflow-y: auto;
}
.preview-area::-webkit-scrollbar { width: 4px; }
.preview-area::-webkit-scrollbar-track { background: transparent; }
.preview-area::-webkit-scrollbar-thumb {
  background: rgba(0, 255, 255, 0.2);
  border-radius: 2px;
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  font-size: 13px;
  color: #00ffff;
  flex-wrap: wrap;
  gap: 10px;
}

.clear-btn {
  background: none;
  border: none;
  color: #ff4444;
  cursor: pointer;
  font-size: 12px;
}

.preview-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
}

@media (max-width: 768px) {
  .upload-content {
    flex-direction: column;
  }
  .dropzone {
    width: 100%;
    min-height: 160px;
  }
  .preview-area {
    max-height: none;
  }
  .preview-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 480px) {
  .preview-grid {
    grid-template-columns: 1fr;
  }
}

.preview-item {
  position: relative;
  border-radius: 12px;
  max-width: 200px;
  overflow: hidden;
  background: rgba(0, 0, 0, 0.3);
  transition: transform 0.2s;
}

.preview-item:hover {
  transform: translateY(-2px);
}

.preview-item img {
  width: 100%;
  max-height: 185px;
  object-fit: cover;
}

.preview-remove {
  position: absolute;
  top: 4px;
  right: 4px;
  width: 22px;
  height: 22px;
  background: rgba(255, 0, 0, 0.8);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.3s;
  z-index: 2;
}

.preview-item:hover .preview-remove {
  opacity: 1;
}

/* 上传状态 */
.upload-status {
  position: absolute;
  top: 4px;
  left: 4px;
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  gap: 4px;
  z-index: 2;
}
.upload-status i { font-size: 10px; }
.upload-status.uploading {
  background: rgba(0, 255, 255, 0.2);
  color: #00ffff;
}
.upload-status.done {
  background: rgba(0, 255, 0, 0.15);
  color: #00ff88;
}
.upload-status.failed {
  background: rgba(255, 0, 0, 0.2);
  color: #ff4757;
}
.preview-item.upload-error {
  box-shadow: 0 0 0 1px rgba(255, 71, 87, 0.5);
}

.config-card{
  margin-top: 20px;
}

/* 描述显示区域 */
.image-description {
  padding: 8px 10px;
  background: rgba(0, 0, 0, 0.7);
  border-top: 1px solid rgba(0, 255, 255, 0.2);
  cursor: pointer;
  transition: all 0.3s;
}

.image-description:hover {
  background: rgba(0, 0, 0, 0.85);
}

.image-description.empty {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  color: #666;
  font-size: 12px;
  padding: 12px;
}

.image-description.empty i {
  font-size: 14px;
  color: #00ffff;
}

.image-description.empty:hover {
  color: #00ffff;
}

.desc-content {
  display: flex;
  align-items: flex-start;
  gap: 6px;
  margin-bottom: 6px;
}

.desc-content i {
  font-size: 10px;
  color: #00ffff;
  margin-top: 2px;
}

.desc-text {
  font-size: 11px;
  color: #ccc;
  line-height: 1.4;
  flex: 1;
}

.desc-footer {
  display: flex;
  justify-content: space-between;
  font-size: 9px;
  color: #888;
}

.desc-type {
  padding: 2px 6px;
  background: rgba(0, 255, 255, 0.15);
  border-radius: 10px;
}

.desc-edit {
  cursor: pointer;
  transition: color 0.3s;
}

.desc-edit:hover {
  color: #00ffff;
}

/* AI 提示 */
.ai-tip {
  margin-top: 16px;
  padding: 10px 12px;
  background: rgba(0, 255, 255, 0.05);
  border-radius: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 11px;
  color: #888;
  border: 1px solid rgba(0, 255, 255, 0.15);
}

.ai-tip i {
  color: #00ffff;
  font-size: 14px;
}

/* 框架选择器 */
.framework-selector {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 24px;
}

.framework-btn {
  padding: 12px;
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid rgba(0, 255, 255, 0.3);
  border-radius: 12px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
}

.framework-btn.active {
  background: linear-gradient(135deg, rgba(0, 255, 255, 0.2), rgba(255, 0, 255, 0.2));
  border-color: #00ffff;
  color: #00ffff;
}

/* 滑块区域 */
.slider-container {
  margin-top: 8px;
}

.slider-label {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
  font-size: 14px;
  color: #aaa;
}

.quality-value {
  color: #00ffff;
  font-weight: 600;
}

.slider-wrapper {
  position: relative;
  padding-bottom: 24px;
}

.quality-slider {
  width: 100%;
  height: 4px;
  -webkit-appearance: none;
  background: #333;
  border-radius: 2px;
  margin: 0;
}

.quality-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 16px;
  height: 16px;
  background: #00ffff;
  border-radius: 50%;
  cursor: pointer;
  box-shadow: 0 0 8px #00ffff;
}

.slider-marks {
  display: flex;
  justify-content: space-between;
  margin-top: 8px;
  padding: 0 2px;
}

.slider-mark {
  font-size: 10px;
  color: #555;
  cursor: pointer;
  transition: all 0.2s;
  text-align: center;
  flex: 1;
}

.slider-mark:hover {
  color: #00ffff;
}

.slider-mark.active {
  color: #00ffff;
  font-weight: 600;
}

/* 生成按钮 */
.generate-btn {
  width: 100%;
  padding: 16px;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  border: none;
  border-radius: 16px;
  font-size: 16px;
  font-weight: 600;
  color: white;
  cursor: pointer;
  transition: all 0.3s;
}

.generate-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 10px 30px rgba(0, 255, 255, 0.3);
}

.generate-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>