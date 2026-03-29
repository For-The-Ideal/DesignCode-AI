<!-- components/optimization/CompareModal.vue -->
<template>
  <Teleport to="body">
    <div class="modal-overlay" v-if="visible" @click.self="close">
      <div class="modal-content compare-modal">
        <div class="modal-header">
          <div class="header-icon">
            <i class="fas fa-code-branch"></i>
          </div>
          <h3>优化前后对比</h3>
          <button class="close-btn" @click="close">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <div class="modal-body">
          <!-- 对比内容 -->
          <div class="compare-container">
            <div class="compare-grid">
              <!-- 优化前 -->
              <div class="compare-box original">
                <div class="compare-title">
                  <i class="fas fa-history"></i>
                  <span>优化前</span>
                  <span class="score-badge">{{ originalScore }}分</span>
                </div>
                <div class="compare-code-wrapper">
                  <pre class="compare-code">{{ formattedOriginalCode }}</pre>
                </div>
              </div>
              
              <!-- 优化后 -->
              <div class="compare-box optimized">
                <div class="compare-title">
                  <i class="fas fa-magic"></i>
                  <span>优化后</span>
                  <span class="score-badge improved">{{ optimizedScore }}分</span>
                </div>
                <div class="compare-code-wrapper">
                  <pre class="compare-code optimized">{{ formattedOptimizedCode }}</pre>
                </div>
              </div>
            </div>
            
            <!-- 差异统计 -->
            <div class="diff-stats">
              <div class="stat-item">
                <i class="fas fa-plus-circle"></i>
                <span>新增 <strong>{{ diffStats.added }}</strong> 行</span>
              </div>
              <div class="stat-item">
                <i class="fas fa-minus-circle"></i>
                <span>删除 <strong>{{ diffStats.removed }}</strong> 行</span>
              </div>
              <div class="stat-item">
                <i class="fas fa-chart-line"></i>
                <span>质量提升 <strong>+{{ scoreImprovement }}</strong> 分</span>
              </div>
            </div>
          </div>
        </div>
        
        <div class="modal-footer">
          <button class="btn-reject" @click="reject">
            <i class="fas fa-times"></i>
            取消
          </button>
          <button class="btn-accept" @click="accept">
            <i class="fas fa-check"></i>
            应用优化
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  originalCode: {
    type: String,
    default: ''
  },
  optimizedCode: {
    type: String,
    default: ''
  },
  originalScore: {
    type: Number,
    default: 78
  },
  optimizedScore: {
    type: Number,
    default: 92
  },
  diffStats: {
    type: Object,
    default: () => ({
      added: 12,
      removed: 5
    })
  }
})

const emit = defineEmits(['close', 'accept', 'reject'])

const scoreImprovement = computed(() => {
  return props.optimizedScore - props.originalScore
})

const formattedOriginalCode = computed(() => {
  if (!props.originalCode) return '暂无代码'
  const lines = props.originalCode.split('\n')
  return lines.map((line, index) => {
    const lineNum = (index + 1).toString().padStart(3, ' ')
    return `${lineNum} | ${line}`
  }).join('\n')
})

const formattedOptimizedCode = computed(() => {
  if (!props.optimizedCode) return '暂无代码'
  const lines = props.optimizedCode.split('\n')
  return lines.map((line, index) => {
    const lineNum = (index + 1).toString().padStart(3, ' ')
    return `${lineNum} | ${line}`
  }).join('\n')
})

const close = () => {
  emit('close')
}

const accept = () => {
  emit('accept', props.optimizedCode)
  close()
}

const reject = () => {
  emit('reject')
  close()
}
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
  width: 900px;
  max-width: 90%;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  animation: slideUp 0.3s ease;
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
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

/* 对比区域 */
.compare-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 24px;
}

.compare-box {
  background: rgba(0, 0, 0, 0.3);
  border-radius: 16px;
  overflow: hidden;
  border: 1px solid rgba(0, 255, 255, 0.2);
}

.compare-title {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.4);
  border-bottom: 1px solid rgba(0, 255, 255, 0.1);
  font-size: 13px;
  font-weight: 500;
}

.compare-title i {
  color: #00ffff;
}

.score-badge {
  margin-left: auto;
  padding: 2px 8px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  font-size: 11px;
  color: #888;
}

.score-badge.improved {
  background: rgba(0, 255, 0, 0.15);
  color: #0f0;
}

.compare-code-wrapper {
  max-height: 300px;
  overflow: auto;
}

.compare-code {
  padding: 16px;
  font-family: 'Fira Code', monospace;
  font-size: 12px;
  line-height: 1.6;
  color: #aaa;
  white-space: pre-wrap;
  margin: 0;
}

.compare-code.optimized {
  color: #00ffff;
}

/* 差异统计 */
.diff-stats {
  display: flex;
  justify-content: center;
  gap: 32px;
  padding: 16px;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 40px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #aaa;
}

.stat-item i {
  font-size: 14px;
}

.stat-item i.fa-plus-circle {
  color: #0f0;
}

.stat-item i.fa-minus-circle {
  color: #f00;
}

.stat-item i.fa-chart-line {
  color: #00ffff;
}

.stat-item strong {
  color: #fff;
}

/* 底部按钮 */
.modal-footer {
  display: flex;
  gap: 16px;
  justify-content: flex-end;
  padding: 20px 24px;
  border-top: 1px solid rgba(0, 255, 255, 0.1);
}

.btn-accept,
.btn-reject {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 28px;
  border-radius: 40px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-accept {
  background: linear-gradient(135deg, #00ffff, #00aa99);
  border: none;
  color: white;
}

.btn-accept:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 20px rgba(0, 255, 255, 0.3);
}

.btn-reject {
  background: rgba(255, 68, 68, 0.15);
  border: 1px solid rgba(255, 68, 68, 0.3);
  color: #ff6666;
}

.btn-reject:hover {
  background: rgba(255, 68, 68, 0.25);
}

/* 滚动条 */
.modal-body::-webkit-scrollbar {
  width: 6px;
}

.modal-body::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 3px;
}

.modal-body::-webkit-scrollbar-thumb {
  background: rgba(0, 255, 255, 0.3);
  border-radius: 3px;
}

/* 响应式 */
@media (max-width: 768px) {
  .compare-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .diff-stats {
    flex-direction: column;
    align-items: center;
    gap: 12px;
  }
}
</style>