<template>
  <!-- AI 整页操作进行中：遮罩整块编辑区域，拦截一切点击，避免边生成边改标注 -->
  <transition name="ai-fade">
    <div v-if="show" class="ai-page-loading" role="status" aria-live="polite">
      <div class="ai-loading-box">
        <div class="ai-orbit">
          <span class="ai-orbit-ring"></span>
          <span class="ai-orbit-core"><i class="fas fa-magic"></i></span>
        </div>
        <p class="ai-loading-text">{{ text }}</p>
        <p v-if="sub" class="ai-loading-sub">{{ sub }}</p>
      </div>
    </div>
  </transition>
</template>

<script setup>
defineProps({
  // 是否显示遮罩
  show: { type: Boolean, default: false },
  // 主文案
  text: { type: String, default: 'AI 正在生成整页区域描述…' },
  // 副文案（进度/回填说明），空则不显示
  sub: { type: String, default: '' },
})
</script>

<style scoped>
/* 遮罩相对最近的定位父级铺满：调用方容器需 position: relative */
.ai-page-loading {
  position: absolute;
  inset: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(5, 6, 12, 0.6);
  backdrop-filter: blur(3px);
  cursor: progress;
  /* 明确拦截指针，防止点到下层画布/工具栏 */
  pointer-events: auto;
}
.ai-loading-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
  padding: 32px 48px;
  border-radius: 16px;
  background: rgba(10, 14, 22, 0.85);
  border: 1px solid rgba(0, 255, 255, 0.25);
  box-shadow: 0 0 40px rgba(0, 255, 255, 0.12);
}
.ai-orbit { position: relative; width: 64px; height: 64px; }
.ai-orbit-ring {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  border: 2px solid rgba(0, 255, 255, 0.15);
  border-top-color: #00ffff;
  animation: ai-spin 0.9s linear infinite;
}
.ai-orbit-core {
  position: absolute;
  inset: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: rgba(0, 255, 255, 0.1);
  color: #00ffff;
  font-size: 18px;
  animation: ai-pulse 1.4s ease-in-out infinite;
}
.ai-loading-text { margin: 0; font-size: 15px; font-weight: 600; color: rgba(255, 255, 255, 0.9); }
.ai-loading-sub { margin: 0; font-size: 12px; color: rgba(0, 255, 255, 0.6); }
@keyframes ai-spin { to { transform: rotate(360deg); } }
@keyframes ai-pulse { 0%, 100% { transform: scale(1); opacity: 0.85; } 50% { transform: scale(1.08); opacity: 1; } }
.ai-fade-enter-active, .ai-fade-leave-active { transition: opacity 0.2s ease; }
.ai-fade-enter-from, .ai-fade-leave-to { opacity: 0; }
</style>
