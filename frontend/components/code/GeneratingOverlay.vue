<!-- components/code/GeneratingOverlay.vue -->
<template>
  <div class="generating-overlay">
    <div class="generating-content">
      <!-- 旋转的芯片图标 -->
      <div class="generating-icon">
        <div class="icon-ring">
          <i class="fas fa-microchip"></i>
        </div>
        <div class="ring ring-1"></div>
        <div class="ring ring-2"></div>
      </div>

      <!-- 进度百分比 -->
      <div class="generating-percent">{{ progress }}%</div>

      <!-- 状态文字 -->
      <div class="generating-status">{{ currentStep }}</div>

      <!-- 进度条 -->
      <div class="generating-progress-bar">
        <div class="progress-fill" :style="{ width: progress + '%' }">
          <div class="progress-glow"></div>
        </div>
      </div>

      <!-- 步骤指示器 -->
      <div class="generating-steps">
        <div
          v-for="(step, index) in steps"
          :key="index"
          class="step"
          :class="{
            active: currentStepIndex >= index,
            completed: currentStepIndex > index,
          }"
        >
          <div class="step-dot">
            <i v-if="currentStepIndex > index" class="fas fa-check"></i>
            <span v-else>{{ index + 1 }}</span>
          </div>
          <div class="step-label">{{ step.label }}</div>
        </div>
      </div>

      <!-- 动态消息 -->
      <div class="generating-messages">
        <div
          v-for="(msg, idx) in messages"
          :key="idx"
          class="message-item"
          :class="{ new: idx === messages.length - 1 }"
        >
          <i class="fas fa-chevron-right"></i>
          <span>{{ msg }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from "vue";

const props = defineProps({
  // 是否显示
  visible: {
    type: Boolean,
    default: false,
  },
  // 初始进度
  initialProgress: {
    type: Number,
    default: 0,
  },
  // 自定义步骤
  customSteps: {
    type: Array,
    default: null,
  },
});

const emit = defineEmits(["complete", "progress"]);

// 默认步骤
const defaultSteps = [
  { label: "图像识别", progress: 20, message: "正在识别设计稿中的UI元素..." },
  { label: "布局分析", progress: 40, message: "分析页面布局结构..." },
  { label: "代码生成", progress: 60, message: "生成代码框架..." },
  { label: "优化输出", progress: 80, message: "优化代码质量..." },
  { label: "完成", progress: 100, message: "生成完成！" },
];

const steps = computed(() => props.customSteps || defaultSteps);

const progress = ref(props.initialProgress);
const currentStepIndex = ref(0);
const currentStep = ref("准备就绪");
const messages = ref([]);
let timer = null;

// 更新进度
const updateProgress = (newProgress) => {
  progress.value = Math.min(newProgress, 100);
  emit("progress", progress.value);

  // 根据进度更新当前步骤
  for (let i = 0; i < steps.value.length; i++) {
    if (progress.value >= steps.value[i].progress) {
      currentStepIndex.value = i;
      currentStep.value = steps.value[i].label;

      // 添加消息（避免重复添加）
      const lastMsg = messages.value[messages.value.length - 1];
      if (!lastMsg || lastMsg !== steps.value[i].message) {
        messages.value.push(steps.value[i].message);
        // 自动滚动到最新消息
        setTimeout(() => {
          const container = document.querySelector(".generating-messages");
          if (container) container.scrollTop = container.scrollHeight;
        }, 50);
      }
    }
  }

  // 完成时触发事件
  if (progress.value >= 100) {
    setTimeout(() => {
      emit("complete");
    }, 500);
  }
};

// 开始生成（模拟进度）
const start = (duration = 3000) => {
  progress.value = 0;
  messages.value = [];
  currentStepIndex.value = 0;

  const stepDuration = duration / 100;
  let currentProgress = 0;

  if (timer) clearInterval(timer);

  timer = setInterval(() => {
    currentProgress += 1;
    if (currentProgress <= 100) {
      updateProgress(currentProgress);
    } else {
      clearInterval(timer);
      timer = null;
    }
  }, stepDuration);
};

// 设置进度（手动控制）
const setProgress = (value) => {
  updateProgress(value);
};

// 重置
const reset = () => {
  if (timer) clearInterval(timer);
  progress.value = 0;
  messages.value = [];
  currentStepIndex.value = 0;
  currentStep.value = "准备就绪";
  timer = null;
};

// 监听 visible 变化
watch(
  () => props.visible,
  (newVal) => {
    if (newVal && !timer && progress.value === 0) {
      start();
    }
  },
);

onUnmounted(() => {
  if (timer) clearInterval(timer);
});

defineExpose({
  start,
  setProgress,
  reset,
  updateProgress,
});
</script>

<style scoped>
.generating-overlay {
  background: rgba(10, 10, 15, 0.95);
  backdrop-filter: blur(12px);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
  animation: fadeIn 0.3s ease;
}

.generating-content {
  text-align: center;
  max-width: 400px;
  width: 90%;
  padding: 32px;
}

/* 图标动画 */
.generating-icon {
  position: relative;
  width: 80px;
  height: 80px;
  margin: 0 auto 24px;
}

.icon-ring {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: pulse 1s ease-in-out infinite;
}

.icon-ring i {
  font-size: 36px;
  color: white;
}

.ring {
  position: absolute;
  inset: -4px;
  border-radius: 50%;
  border: 2px solid transparent;
  animation: ringRotate 1.5s linear infinite;
}

.ring-1 {
  border-top-color: #00ffff;
  border-right-color: #ff00ff;
}

.ring-2 {
  inset: -8px;
  border-bottom-color: #00ffff;
  border-left-color: #ff00ff;
  animation-duration: 2s;
  animation-direction: reverse;
}

/* 进度百分比 */
.generating-percent {
  font-size: 48px;
  font-weight: 800;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  margin-bottom: 8px;
  font-family: monospace;
}

/* 状态文字 */
.generating-status {
  font-size: 14px;
  color: #00ffff;
  margin-bottom: 20px;
  letter-spacing: 1px;
}

/* 进度条 */
.generating-progress-bar {
  width: 100%;
  height: 4px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 24px;
}

.progress-fill {
  position: relative;
  height: 100%;
  background: linear-gradient(90deg, #00ffff, #ff00ff);
  border-radius: 4px;
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

/* 步骤指示器 */
.generating-steps {
  display: flex;
  justify-content: space-between;
  margin-bottom: 24px;
  gap: 8px;
}

.step {
  flex: 1;
  text-align: center;
}

.step-dot {
  width: 28px;
  height: 28px;
  margin: 0 auto 6px;
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

.step.active .step-dot {
  background: rgba(0, 255, 255, 0.2);
  border-color: #00ffff;
  color: #00ffff;
  box-shadow: 0 0 10px rgba(0, 255, 255, 0.3);
}

.step.completed .step-dot {
  background: #00ffff;
  border-color: #00ffff;
  color: #0a0a0f;
}

.step-label {
  font-size: 10px;
  color: #666;
}

.step.active .step-label {
  color: #00ffff;
}

.step.completed .step-label {
  color: #0f0;
}

/* 消息列表 */
.generating-messages {
  max-height: 100px;
  overflow-y: auto;
  text-align: left;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 12px;
  padding: 12px;
  font-size: 12px;
}

.message-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 0;
  color: #888;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  animation: slideIn 0.2s ease;
}

.message-item i {
  font-size: 10px;
  color: #00ffff;
}

.message-item.new {
  color: #00ffff;
}

.message-item:last-child {
  border-bottom: none;
}

/* 动画 */
@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes ringRotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

@keyframes pulse {
  0%,
  100% {
    transform: scale(0.95);
    opacity: 0.8;
  }
  50% {
    transform: scale(1.05);
    opacity: 1;
  }
}

@keyframes shimmer {
  0% {
    transform: translateX(-20px);
  }
  100% {
    transform: translateX(40px);
  }
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(-10px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

/* 滚动条 */
.generating-messages::-webkit-scrollbar {
  width: 4px;
}

.generating-messages::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 2px;
}

.generating-messages::-webkit-scrollbar-thumb {
  background: rgba(0, 255, 255, 0.3);
  border-radius: 2px;
}
</style>
