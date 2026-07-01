<template>
  <div class="code-page">
    <CodeSidebar />

    <main class="code-main">
      <div class="code-header">
        <div class="header-left">
          <h2 class="header-title">AI 智能代码生成</h2>
          <p class="header-sub">上传设计稿，AI 自动识别并生成高质量代码</p>
        </div>
        <div class="header-stats">
          <div class="stat-card stat-pending">
            <i class="stat-icon fa-solid fa-clock"></i>
            <div class="stat-body">
              <span class="stat-num" :key="taskCounts.pending">{{ taskCounts.pending }}</span>
              <span class="stat-label">排队中</span>
            </div>
          </div>
          <div class="stat-card stat-running">
            <i class="stat-icon fa-solid fa-spinner fa-spin"></i>
            <div class="stat-body">
              <span class="stat-num" :key="taskCounts.running">{{ taskCounts.running }}</span>
              <span class="stat-label">运行中</span>
            </div>
          </div>
          <div class="stat-card stat-success">
            <i class="stat-icon fa-solid fa-circle-check"></i>
            <div class="stat-body">
              <span class="stat-num" :key="taskCounts.success">{{ taskCounts.success }}</span>
              <span class="stat-label">已完成</span>
            </div>
          </div>
          <div class="stat-card stat-failed">
            <i class="stat-icon fa-solid fa-circle-xmark  "></i>
            <div class="stat-body">
              <span class="stat-num" :key="taskCounts.failed">{{ taskCounts.failed }}</span>
              <span class="stat-label">已失败</span>
            </div>
          </div>
        </div>
      </div>

      <section class="upload-section">
        <div class="upload-section-inner">
          <div class="upload-left">
            <div class="glow-section">
              <div class="section-label">
                <span class="label-num">1</span> 
                <span class="section-label-text">上传设计稿</span>
              </div>
              <UploadZone />
            </div>
          </div>

          <div class="upload-right">
            <div class="glow-section">
              <ConfigPanel />
            </div>
          </div>
        </div>
      </section>

      <div class="glow-section">
        <FlowSteps :activeStep="cycleStep" />
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { storeToRefs } from 'pinia'
import CodeSidebar from '~/components/code/CodeSidebar.vue'
import UploadZone from '~/components/upload/UploadZone.vue'
import ConfigPanel from '~/components/code/ConfigPanel.vue'
import FlowSteps from '~/components/code/FlowSteps.vue'
import { useCodeStore } from '~/stores/code'
import { useUserStore } from '~/stores/user'

const store = useCodeStore()
const userStore = useUserStore()
const { taskCounts } = storeToRefs(userStore)

// ═══ 步骤循环高亮 ═══
const cycleStep = ref(0)
let cycleTimer = null

onMounted(() => {
  // 步骤循环高亮
  cycleTimer = setInterval(() => {
    cycleStep.value = (cycleStep.value + 1) % 5
  }, 1800)
})

onUnmounted(() => {
  if (cycleTimer) clearInterval(cycleTimer)
})

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
  top: -1px; left: -1px; right: -1px; bottom: -1px;
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

.code-header {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-radius: 16px;
  gap: 16px;
}
.header-title { font-size: 28px; font-weight: 800; color: #e2e8f0; letter-spacing: -0.3px; }
.header-sub { font-size: 13px; color: rgba(255,255,255,0.45); margin-top: 4px; }

.header-stats {
  display: flex;
  gap: 12px;
  flex-shrink: 0;
}
.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 24px;
  border-radius: 14px;
  border: 1px solid;
  background: rgba(5, 8, 18, 0.85);
  backdrop-filter: blur(8px);
  position: relative;
  overflow: hidden;
}
/* 顶部发光线 */
.stat-card::before {
  content: '';
  position: absolute;
  top: 0; left: 16px; right: 16px;
  height: 1px;
  opacity: 0.6;
}
.stat-pending::before { background: linear-gradient(90deg, transparent, #f59e0b, transparent); }
.stat-running::before { background: linear-gradient(90deg, transparent, #00ffff, transparent); }
.stat-success::before { background: linear-gradient(90deg, transparent, #22c55e, transparent); }

.stat-icon {
  font-size: 24px;
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: rgba(255,255,255,0.04);
  flex-shrink: 0;
}
.stat-body {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.stat-num {
  font-family: 'Courier New', 'Consolas', 'Source Code Pro', monospace;
  font-size: 32px;
  font-weight: 700;
  line-height: 1;
  animation: tickIn 0.35s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.stat-label {
  font-size: 13px;
  font-weight: 500;
  letter-spacing: 1px;
}

@keyframes tickIn {
  0%   { transform: scale(0.5); opacity: 0; }
  60%  { transform: scale(1.2); }
  100% { transform: scale(1); opacity: 1; }
}

/* ── 排队中 ── */
.stat-pending { border-color: rgba(245,158,11,0.18); }
.stat-pending .stat-icon { color: #f59e0b; box-shadow: 0 0 14px rgba(245,158,11,0.25); }
.stat-pending .stat-num { color: #fbbf24; text-shadow: 0 0 10px rgba(251,191,36,0.4); }
.stat-pending .stat-label { color: rgba(245,158,11,0.6); }

/* ── 运行中 ── */
.stat-running { border-color: rgba(0,255,255,0.18); }
.stat-running .stat-icon { color: #00ffff; box-shadow: 0 0 14px rgba(0,255,255,0.25); }
.stat-running .stat-num { color: #22d3ee; text-shadow: 0 0 10px rgba(34,211,238,0.4); }
.stat-running .stat-label { color: rgba(0,255,255,0.6); }

/* ── 已完成 ── */
.stat-success { border-color: rgba(34,197,94,0.18); }
.stat-success .stat-icon { color: #22c55e; box-shadow: 0 0 14px rgba(34,197,94,0.25); }
.stat-success .stat-num { color: #4ade80; text-shadow: 0 0 10px rgba(74,222,128,0.4); }
.stat-success .stat-label { color: rgba(34,197,94,0.6); }

.upload-section-inner { display: flex; gap: 24px; }
.upload-left { flex: 1.5; display: flex; flex-direction: column; gap: 12px; }
.upload-left .glow-section,
.upload-right .glow-section { flex: 1; display: flex; flex-direction: column; padding: 20px; }

.section-label {
  font-size: 14px;
  font-weight: 600;
  color: #60a5fa;
  text-transform: uppercase;
  letter-spacing: 1.5px;
  margin-bottom: 16px;
  display: inline-block;
  padding-left: 14px;
}
.section-label::before {
  content: '';
  position: absolute;
  left: 0; top: 2px; bottom: 2px;
  width: 2px;
  background: linear-gradient(135deg, #60a5fa, #818cf8);
  border-radius: 1px;
}
.label-num {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px; height: 24px;
  border-radius: 8px;
  font-size: 12px;
  background: #60a5fa;
  color: #fff;
  font-weight: 700;
  margin-right: 6px;
}
.section-label-text { font-weight: 700; color: #fff; }
.upload-right { flex: 1; display: flex; flex-direction: column; gap: 16px; }

@media (max-width: 1100px) {
  .upload-section-inner { flex-direction: column; }
}
</style>
