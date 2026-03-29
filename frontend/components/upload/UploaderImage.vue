<template>
  <div class="upload-section">
    <div class="tech-card">
      <div class="card-title">
        <i class="fas fa-cloud-upload-alt"></i>
        <span>上传设计稿</span>
      </div>
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

      <!-- 预览区域（每行最多3张） -->
      <div class="preview-area" v-if="images.length > 0">
        <div class="preview-header">
          <span><i class="fas fa-images"></i> 已上传 ({{ images.length }}/{{ maxLen }})</span>
          <button class="clear-btn" @click="clearAll">清空全部</button>
        </div>
        <div class="preview-grid">
          <div v-for="(img, idx) in images" :key="img.id" class="preview-item">
            <img :src="img.preview" alt="预览" />
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

      <!-- AI 智能分析提示 -->
      <div v-if="images.length > 0" class="ai-tip">
        <i class="fas fa-robot"></i>
        <span>AI 将结合图片和描述生成更精准的代码，描述越详细效果越好</span>
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
        <div class="framework-btn" :class="{ active: framework === 'vue' }" @click="framework = 'vue'">
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

    <!-- 生成按钮 -->
    <button class="generate-btn" @click="generateCode" :disabled="images.length === 0">
      <i class="fas fa-play"></i> 开始生成代码
    </button>

    <!-- 评分卡片 -->
    <div class="score-card" v-if="score > 0">
      <div class="score-header">
        <div>
          <h3>质量评分</h3>
          <p class="score-sub">AI 多维评估</p>
        </div>
        <div class="score-value">{{ score }}</div>
      </div>
      <div class="score-dimensions">
        <div v-for="dim in scoreDimensions" :key="dim.name" class="dimension-item">
          <div class="dimension-label">
            <span><i :class="dim.icon"></i> {{ dim.name }}</span>
            <span>{{ dim.score }}</span>
          </div>
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: dim.score + '%' }"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- 描述编辑弹窗 -->
    <DescEditorModal
      :visible="descModalVisible"
      :images="images"
      :focus-index="editingIndex"
      @update:visible="descModalVisible = $event"
      @save="handleSaveDescriptions"
      @close="closeDescModal"
    />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import DescEditorModal from './DescEditorModal.vue'
import {useToastNotificationStore} from "~/stores/toastNotification"
import { identifyApi } from "/api/identify"
const toastNotificationStore = useToastNotificationStore()

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

// 描述弹窗
const descModalVisible = ref(false)
const editingIndex = ref(null)


// 计算是否达到上限
const isMaxReached = computed(() => images.value.length >= maxLen.value)

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

// 添加图片
const addImages = (files) => {
  const imageFiles = files.filter(f => f.type.startsWith('image/'))
  const remaining = maxLen.value - images.value.length
  const toAdd = imageFiles.slice(0, remaining)
  
  toAdd.forEach(file => {
    const reader = new FileReader()
    reader.onload = (e) => {
      images.value.push({
        file: file,
        preview: e.target.result,
        description: '',
        type: detectImageType(file.name)
      })
    }
    reader.readAsDataURL(file)
  })
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
  descModalVisible.value = true
}

// 关闭描述弹窗
const closeDescModal = () => {
  descModalVisible.value = false
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
  try {
    
    // 构建设计稿数组 - 图片和描述在同一个对象里
    const designs = await Promise.all(
      images.value.map(async (img) => ({
        image: await toBase64(img.file),
        type: img.type,
        description: img.description || ''
      }))
    )
    
    // 请求参数 - 只有一个 designs 数组
    const payload = {
      designs: designs,
      framework: framework.value,
      quality: qualityValue.value
    }
    toastNotificationStore.info('正在生成代码...')
    console.log('payload:', payload)
    await identifyApi.sendFile(payload)
    // await uploaderStore.setUserInfo(payload)
    // const response = await fetch('/api/generate', {
    //   method: 'POST',
    //   headers: { 'Content-Type': 'application/json' },
    //   body: JSON.stringify(payload)
    // })
    
    // const result = await response.json()
    
    // if (result.code) {
    //   console.log('生成的代码:', result.code)
    //   score.value = result.score || 85
    //   scoreDimensions.value = result.dimensions || [
    //     { name: '视觉还原度', score: 85, icon: 'fas fa-palette' },
    //     { name: '代码质量', score: 82, icon: 'fas fa-code' },
    //     { name: '响应式设计', score: 78, icon: 'fas fa-mobile-alt' },
    //     { name: '性能优化', score: 75, icon: 'fas fa-tachometer-alt' }
    //   ]
    // }
    
  } catch (error) {
    console.error('生成失败:', error)
   
  } finally {

  }
}


 const toBase64 = (file) => {
      return new Promise((resolve, reject) => {
        const reader = new FileReader()
        reader.readAsDataURL(file)
        reader.onload = () => resolve(reader.result.split(',')[1])
        reader.onerror = reject
      })
    }

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

/* 拖拽上传区 */
.dropzone {
  border: 2px dashed rgba(0, 255, 255, 0.4);
  border-radius: 20px;
  padding: 40px;
  text-align: center;
  background: rgba(0, 0, 0, 0.3);
  cursor: pointer;
  transition: all 0.3s;
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
  width: 70px;
  height: 70px;
  background: linear-gradient(135deg, rgba(0, 255, 255, 0.2), rgba(255, 0, 255, 0.2));
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
}

.dropzone-icon i {
  font-size: 28px;
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
  margin-top: 20px;
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
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

@media (max-width: 768px) {
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
  overflow: hidden;
  background: rgba(0, 0, 0, 0.3);
  transition: transform 0.2s;
}

.preview-item:hover {
  transform: translateY(-2px);
}

.preview-item img {
  width: 100%;
  aspect-ratio: 1;
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

/* 评分卡片 */
.score-card {
  background: rgba(15, 20, 30, 0.6);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(0, 255, 255, 0.2);
  border-radius: 24px;
  overflow: hidden;
  margin: 24px 0px;
}

.score-header {
  background: linear-gradient(135deg, rgba(0, 255, 255, 0.2), rgba(255, 0, 255, 0.2));
  padding: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.score-header h3 {
  font-size: 18px;
}

.score-sub {
  font-size: 12px;
  opacity: 0.7;
  margin-top: 4px;
}

.score-value {
  font-size: 48px;
  font-weight: 700;
  color: #00ffff;
}

.score-dimensions {
  padding: 20px;
}

.dimension-item {
  margin-bottom: 16px;
}

.dimension-label {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
  margin-bottom: 8px;
  color: #aaa;
}

.progress-bar {
  width: 100%;
  height: 6px;
  background: #333;
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #00ffff, #ff00ff);
  border-radius: 3px;
}
</style>