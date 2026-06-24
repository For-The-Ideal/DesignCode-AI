<template>
  <section class="flow-section">
    <h3 class="flow-title">生成流程</h3>
    <div class="flow-steps">
      <template v-for="(step, i) in steps" :key="i">
        <!-- 箭头 -->
        <div v-if="i > 0" class="flow-arrow">→</div>
        <!-- 步骤卡片 -->
        <div
          class="flow-step"
          :class="{ 'flow-step--active': i === activeStep }"
        >
          <span class="step-num">{{ i + 1 }}</span>
          <div>
            <div class="step-title">{{ step.title }}</div>
            <div class="step-desc">{{ step.desc }}</div>
          </div>
        </div>
      </template>
    </div>
  </section>
</template>

<script setup>
defineProps({
  steps: {
    type: Array,
    default: () => [
      { title: '上传设计稿', desc: '多图上传，描述说明' },
      { title: 'AI识别分析', desc: '智能识别页面元素' },
      { title: '生成代码', desc: '生成高质量代码' },
      { title: '预览效果', desc: '实时预览生成效果' },
      { title: '导出代码', desc: '下载或部署项目' },
    ],
  },
  activeStep: { type: Number, default: 0 },
})
</script>

<style scoped>
.flow-section {
  border-radius: 16px;
  padding: 20px 24px;
  background: rgba(15, 20, 30, 0.5);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(0, 255, 255, 0.08);
}

.flow-title {
  font-size: 15px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
  margin-bottom: 16px;
}

.flow-steps {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  gap: 8px;
}

.flow-arrow {
  display: flex;
  align-items: center;
  padding: 0 4px;
  font-size: 18px;
  color: rgba(255, 255, 255, 0.2);
}

.flow-step {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 14px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(0, 255, 255, 0.05);
  transition: all 0.3s;
}
.flow-step--active {
  background: rgba(0, 255, 255, 0.05);
  border-color: rgba(0, 255, 255, 0.15);
}

.step-num {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  font-size: 11px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.45);
}
.flow-step--active .step-num {
  background: rgba(0, 255, 255, 0.15);
  color: #00ffff;
}

.step-title {
  font-size: 13px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.55);
}
.flow-step--active .step-title {
  color: rgba(255, 255, 255, 0.85);
}

.step-desc {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.35);
  margin-top: 1px;
}
.flow-step--active .step-desc {
  color: rgba(255, 255, 255, 0.45);
}
</style>
