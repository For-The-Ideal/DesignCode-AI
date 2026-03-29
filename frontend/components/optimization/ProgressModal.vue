<!-- components/optimization/ProgressModal.vue -->
<template>
  <Teleport to="body">
    <div class="modal-overlay" v-if="visible" @click.self="handleClose">
      <div class="modal-content progress-modal">
        <div class="modal-header">
          <div class="header-icon">
            <i class="fas fa-microchip"></i>
          </div>
          <h3>AI 智能优化</h3>
          <button class="close-btn" @click="handleClose" v-if="!optimizing">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <div class="modal-body">
          <!-- 进度信息 -->
          <div class="progress-section">
            <div class="progress-info">
              <span class="progress-label">{{ progressMsg }}</span>
              <span class="progress-percent">{{ progress }}%</span>
            </div>
            <div class="progress-bar-wrapper">
              <div class="progress-bar">
                <div class="progress-fill" :style="{ width: progress + '%' }">
                  <div class="progress-glow"></div>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 步骤指示器 -->
          <div class="steps-container">
            <div 
              v-for="(step, index) in steps" 
              :key="index"
              class="step-item"
              :class="{ 
                active: progress >= step.progress, 
                completed: progress > step.progress 
              }"
            >
              <div class="step-dot">
                <i v-if="progress > step.progress" class="fas fa-check"></i>
                <span v-else>{{ index + 1 }}</span>
              </div>
              <div class="step-label">{{ step.label }}</div>
              <div class="step-line" v-if="index < steps.length - 1"></div>
            </div>
          </div>
          
          <!-- 加载动画 -->
          <div class="loading-animation" v-if="optimizing && progress < 100">
            <div class="loader"></div>
            <p>AI 正在分析优化中...</p>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  progress: {
    type: Number,
    default: 0
  },
  progressMsg: {
    type: String,
    default: '正在分析代码...'
  },
  optimizing: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['close', 'complete'])

const steps = [
  { label: '分析代码', progress: 20 },
  { label: '识别问题', progress: 40 },
  { label: '应用优化', progress: 60 },
  { label: '验证结果', progress: 80 },
  { label: '完成', progress: 100 }
]

const handleClose = () => {
  emit('close')
}

// 监听进度完成
watch(() => props.progress, (newVal) => {
  if (newVal >= 100 && props.optimizing === false) {
    setTimeout(() => {
      emit('complete')
    }, 500)
  }
})
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.85);
  backdrop-filter: blur(12px);
  z-index: 1000;
  display: flex;
  justify-content: center;
  align-items: center;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal-content {
  background: rgba(10, 20, 30, 0.95);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(0, 255, 255, 0.3);
  border-radius: 28px;
  width: 500px;
  max-width: 90%;
  animation: slideUp 0.3s ease;
  overflow: hidden;
}

@keyframes slideUp {
  from {
    transform: translateY(50px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.modal-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 20px 24px;
  background: linear-gradient(135deg, rgba(0, 255, 255, 0.1), rgba(255, 0, 255, 0.1));
  border-bottom: 1px solid rgba(0, 255, 255, 0.2);
}

.header-icon {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-icon i {
  font-size: 18px;
  color: white;
}

.modal-header h3 {
  flex: 1;
  font-size: 18px;
  font-weight: 600;
  color: #fff;
  margin: 0;
}

.close-btn {
  width: 32px;
  height: 32px;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 8px;
  color: #888;
  cursor: pointer;
  transition: all 0.3s;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
}

.modal-body {
  padding: 32px 24px;
}

/* 进度区域 */
.progress-section {
  margin-bottom: 32px;
}

.progress-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.progress-label {
  font-size: 13px;
  color: #aaa;
}

.progress-percent {
  font-size: 20px;
  font-weight: 700;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}

.progress-bar-wrapper {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  padding: 2px;
}

.progress-bar {
  height: 6px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  position: relative;
  height: 100%;
  background: linear-gradient(90deg, #00ffff, #ff00ff);
  border-radius: 3px;
  transition: width 0.3s ease;
}

.progress-glow {
  position: absolute;
  top: 0;
  right: 0;
  width: 20px;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.5));
  animation: shimmer 1s infinite;
}

@keyframes shimmer {
  0% {
    transform: translateX(-20px);
  }
  100% {
    transform: translateX(40px);
  }
}

/* 步骤指示器 */
.steps-container {
  display: flex;
  justify-content: space-between;
  margin-bottom: 32px;
  position: relative;
}

.step-item {
  flex: 1;
  text-align: center;
  position: relative;
  z-index: 2;
}

.step-dot {
  width: 32px;
  height: 32px;
  margin: 0 auto 8px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  color: #888;
  transition: all 0.3s;
}

.step-item.active .step-dot {
  background: rgba(0, 255, 255, 0.2);
  border-color: #00ffff;
  color: #00ffff;
  box-shadow: 0 0 10px rgba(0, 255, 255, 0.3);
}

.step-item.completed .step-dot {
  background: #00ffff;
  border-color: #00ffff;
  color: #0a0a0f;
}

.step-label {
  font-size: 11px;
  color: #666;
}

.step-item.active .step-label {
  color: #00ffff;
}

.step-item.completed .step-label {
  color: #0f0;
}

/* 加载动画 */
.loading-animation {
  text-align: center;
  padding-top: 16px;
}

.loader {
  width: 40px;
  height: 40px;
  margin: 0 auto 12px;
  border: 2px solid rgba(0, 255, 255, 0.2);
  border-top-color: #00ffff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.loading-animation p {
  font-size: 12px;
  color: #888;
}
</style>