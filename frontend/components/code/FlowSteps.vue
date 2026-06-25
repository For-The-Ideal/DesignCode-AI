<template>
  <section class="flow-section">
    <h3 class="flow-title">生成流程</h3>
    <div class="flow-steps">
      <template v-for="(step, i) in steps" :key="i">
        <!-- 箭头 -->
        <div v-if="i > 0" class="flow-arrow">
          <span class="flow-arrow-line"></span>
          <i class="fa-solid fa-angle-right"></i>
        </div>
        <!-- 步骤卡片 -->
        <div
          class="flow-step"
          :class="{ 'flow-step--active': i === activeStep }"
        >
          <i class="step-icon" :class="step.icon"></i>
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
      { title: '上传设计稿', desc: '多图上传，描述说明', icon: 'fa-solid fa-cloud-arrow-up' },
      { title: 'AI识别分析', desc: '智能识别页面元素', icon: 'fa-solid fa-microchip' },
      { title: '生成代码', desc: '生成高质量代码', icon: 'fa-solid fa-code' },
      { title: '预览效果', desc: '实时预览生成效果', icon: 'fa-solid fa-display' },
      { title: '导出代码', desc: '下载或部署项目', icon: 'fa-solid fa-download' },
    ],
  },
  activeStep: { type: Number, default: 0 },
})
</script>

<style scoped>
.flow-section {
  border-radius: 16px;
  padding: 28px 24px;
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
  align-self: center;
  gap: 2px;
  padding: 0 2px;
  color: rgba(255, 255, 255, 0.25);
}
.flow-arrow-line {
  width: 20px;
  height: 1px;
  border-top: 1px dashed currentColor;
}
.flow-arrow i {
  font-size: 10px;
}

.flow-step {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 20px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(0, 255, 255, 0.06);
  transition: all 0.3s;
}
.flow-step--active {
  background: rgba(0, 255, 255, 0.06);
  border-color: rgba(0, 255, 255, 0.25);
  box-shadow: 0 0 16px rgba(0, 255, 255, 0.06);
}

.step-icon {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  font-size: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: rgba(255, 255, 255, 0.04);
  color: rgba(255, 255, 255, 0.35);
  border: 1px solid rgba(0, 255, 255, 0.08);
  transition: all 0.3s;
}
.flow-step--active .step-icon {
  background: rgba(0, 255, 255, 0.12);
  color: #00ffff;
  border-color: rgba(0, 255, 255, 0.3);
  box-shadow: 0 0 10px rgba(0, 255, 255, 0.1);
}

.step-title {
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.55);
}
.flow-step--active .step-title {
  color: rgba(255, 255, 255, 0.85);
}

.step-desc {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.35);
  margin-top: 2px;
}
.flow-step--active .step-desc {
  color: rgba(255, 255, 255, 0.45);
}
</style>
