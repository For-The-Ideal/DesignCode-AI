<!-- components/code/GeneratingOverlay.vue -->
<!-- 集成了 CodeEditor + 生成进度覆盖层的复合组件 -->
<template>
  <div class="generating-wrapper">
    <!-- 底层：代码编辑器 -->
    <CodeEditor
      :model-value="modelValue"
      :language="language"
      :readonly="readonly"
      :auto-scroll="autoScroll"
      :height="height"
      :placeholder="placeholder"
      @update:model-value="emit('update:modelValue', $event)"
    />
    <!-- 上层：进度覆盖层 -->
    <div v-if="visible" class="generating-overlay">
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
        <div class="generating-percent">{{ displayProgress }}%</div>

        <!-- 状态文字 -->
        <div class="generating-status">{{ displayStep }}</div>

        <!-- 进度条 -->
        <div class="generating-progress-bar">
          <div class="progress-fill" :style="{ width: displayProgress + '%' }">
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
  </div>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import CodeEditor from "./CodeEditor.vue";

const props = defineProps({
  // ── 编辑器相关 ──
  modelValue: { type: String, default: '' },
  language: { type: String, default: 'dart' },
  readonly: { type: Boolean, default: false },
  autoScroll: { type: Boolean, default: false },
  height: { type: [String, Number], default: '700px' },
  placeholder: { type: String, default: '// AI 生成的代码将在这里展示...' },

  // ── 覆盖层相关 ──
  visible: { type: Boolean, default: false },
  progress: { type: Number, default: 0 },
  currentStep: { type: String, default: '' },
  customSteps: { type: Array, default: null },
});

const emit = defineEmits(["update:modelValue", "complete"]);

// 默认步骤
const defaultSteps = [
  { label: "图像识别", progress: 20, message: "正在识别设计稿中的UI元素..." },
  { label: "布局分析", progress: 40, message: "分析页面布局结构..." },
  { label: "代码生成", progress: 60, message: "AI 正在生成代码框架..." },
  { label: "优化输出", progress: 80, message: "优化代码质量中..." },
  { label: "完成", progress: 100, message: "生成完成！" },
];

const steps = computed(() => props.customSteps || defaultSteps);
const currentStepIndex = ref(0);
const messages = ref([]);

// 显示用状态文字：优先使用真实 currentStep，回退为根据进度推算
const displayStep = computed(() => {
  if (props.currentStep) {
    const stepMap = {
      DownloadImages: '下载图片素材',
      VisionAnalyzeSkill: '视觉分析设计稿',
      FlutterGenerateSkill: 'Flutter 代码生成',
      Vue3GenerateSkill: 'Vue3 代码生成',
      ReactGenerateSkill: 'React 代码生成',
      GenerateSkill: '代码生成',
      SaveResult: '保存生成结果',
      Done: '生成完成',
    }
    return stepMap[props.currentStep] || props.currentStep
  }
  for (let i = steps.value.length - 1; i >= 0; i--) {
    if (props.progress >= steps.value[i].progress) {
      return steps.value[i].label
    }
  }
  return '准备就绪'
})

// 显示用进度（直接使用外部 progress）
const displayProgress = computed(() => props.progress)

// 监听外部 progress
watch(() => props.progress, (val) => {
  if (val == null) return
  for (let i = 0; i < steps.value.length; i++) {
    if (val >= steps.value[i].progress) {
      currentStepIndex.value = i
    }
  }
  for (let i = 0; i < steps.value.length; i++) {
    if (val >= steps.value[i].progress) {
      const lastMsg = messages.value[messages.value.length - 1]
      if (!lastMsg || lastMsg !== steps.value[i].message) {
        messages.value.push(steps.value[i].message)
      }
    }
  }
  if (val >= 100) {
    emit('complete')
  }
})

// visible 从 false → true 时重置消息
watch(() => props.visible, (val) => {
  if (val) {
    messages.value = []
    currentStepIndex.value = 0
  }
})
</script>

<style scoped>
.generating-wrapper {
  position: relative;
  width: 100%;
  height: 100%;
}

.generating-overlay {
  position: absolute;
  inset: 0;
  background: rgba(10, 10, 15, 0.92);
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
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes ringRotate {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

@keyframes pulse {
  0%, 100% { transform: scale(0.95); opacity: 0.8; }
  50% { transform: scale(1.05); opacity: 1; }
}

@keyframes shimmer {
  0% { transform: translateX(-20px); }
  100% { transform: translateX(40px); }
}

@keyframes slideIn {
  from { opacity: 0; transform: translateX(-10px); }
  to { opacity: 1; transform: translateX(0); }
}

/* 滚动条 */
.generating-messages::-webkit-scrollbar { width: 4px; }
.generating-messages::-webkit-scrollbar-track { background: rgba(255, 255, 255, 0.05); border-radius: 2px; }
.generating-messages::-webkit-scrollbar-thumb { background: rgba(0, 255, 255, 0.3); border-radius: 2px; }
</style>