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
          <div
            v-for="card in statCards"
            :key="card.key"
            :class="['stat-card', `stat-${card.key}`]"
          >
            <span class="stat-icon" v-html="card.svg"></span>
            <div class="stat-body">
              <span class="stat-num" :key="card.count">{{ card.count }}</span>
              <span class="stat-label">{{ card.label }}</span>
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
import { ref, computed, onMounted, onUnmounted } from 'vue'
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

// ═══ 统计卡片 ═══
const statCards = computed(() => [
  { key: 'pending', label: '排队中', count: taskCounts.value.pending,
    svg: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><path d="M12 2L20 7v10l-8 5-8-5V7z"/><path d="M8 12h4v4" stroke-linecap="round"/><circle cx="12" cy="12" r="1.5" fill="currentColor" stroke="none"/><circle cx="16" cy="8" r="1" fill="currentColor" stroke="none" opacity=".4"/><circle cx="16" cy="16" r="1" fill="currentColor" stroke="none" opacity=".4"/></svg>' },
  { key: 'running', label: '运行中', count: taskCounts.value.running,
    svg: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><path d="M12 2L20 7v10l-8 5-8-5V7z"/><circle class="run-ring" cx="12" cy="12" r="4.5" stroke-dasharray="10 18" stroke-linecap="round"/><circle cx="12" cy="12" r="1.5" fill="currentColor" stroke="none"/></svg>' },
  { key: 'success', label: '已完成', count: taskCounts.value.success,
    svg: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><path d="M12 2L20 7v10l-8 5-8-5V7z"/><path d="M8 12l2.5 2.5L16 9" stroke-linecap="round" stroke-linejoin="round"/></svg>' },
  { key: 'failed', label: '已失败', count: taskCounts.value.failed,
    svg: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><path d="M12 2L20 7v10l-8 5-8-5V7z"/><path d="M9 9l6 6M15 9l-6 6" stroke-linecap="round"/></svg>' },
])

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

/* ═══ HUD 数据面板 ═══ */
.stat-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 20px;
  min-width: 150px;
  background: rgba(8, 10, 20, 0.9);
  backdrop-filter: blur(12px);
  position: relative;
  overflow: hidden;
  clip-path: polygon(
    8px 0, calc(100% - 8px) 0,
    100% 8px, 100% calc(100% - 8px),
    calc(100% - 8px) 100%, 8px 100%,
    0 calc(100% - 8px), 0 8px
  );
}
/* 外框发光 */
.stat-card::before {
  content: '';
  position: absolute;
  inset: 1px;
  clip-path: polygon(
    8px 0, calc(100% - 8px) 0,
    100% 8px, 100% calc(100% - 8px),
    calc(100% - 8px) 100%, 8px 100%,
    0 calc(100% - 8px), 0 8px
  );
  z-index: -1;
}
.stat-pending::before { background: linear-gradient(135deg, rgba(245,158,11,0.06), rgba(245,158,11,0.01)); }
.stat-running::before { background: linear-gradient(135deg, rgba(0,255,255,0.08), rgba(0,255,255,0.01)); }
.stat-success::before { background: linear-gradient(135deg, rgba(34,197,94,0.06), rgba(34,197,94,0.01)); }
.stat-failed::before { background: linear-gradient(135deg, rgba(239,68,68,0.06), rgba(239,68,68,0.01)); }

/* 顶部扫描线 */
.stat-card::after {
  content: '';
  position: absolute;
  top: 0; left: 20px; right: 20px;
  height: 1px;
  opacity: 0;
  transition: opacity 0.4s;
}
.stat-card:hover::after { opacity: 1; }
.stat-pending::after { background: linear-gradient(90deg, transparent, #f59e0b, transparent); box-shadow: 0 0 6px rgba(245,158,11,0.5); }
.stat-running::after { background: linear-gradient(90deg, transparent, #00ffff, transparent); box-shadow: 0 0 6px rgba(0,255,255,0.5); }
.stat-success::after { background: linear-gradient(90deg, transparent, #22c55e, transparent); box-shadow: 0 0 6px rgba(34,197,94,0.5); }
.stat-failed::after { background: linear-gradient(90deg, transparent, #ef4444, transparent); box-shadow: 0 0 6px rgba(239,68,68,0.5); }

.stat-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  position: relative;
}
.stat-icon svg {
  width: 100%;
  height: 100%;
}
/* v-html 子元素需非 scoped 样式才能命中，见下方 <style> 块 */
.stat-icon::after {
  content: '';
  position: absolute;
  inset: -4px;
  border-radius: 4px;
  opacity: 0.15;
}

.stat-body {
  display: flex;
  flex-direction: column;
  gap: 2px;
  align-items: flex-end;
}
.stat-num {
  font-family: 'JetBrains Mono', 'Fira Code', 'Consolas', monospace;
  font-size: 28px;
  font-weight: 700;
  line-height: 1;
  font-variant-numeric: tabular-nums;
  position: relative;
}
.stat-label {
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 2px;
  text-transform: uppercase;
}

/* ═══ 排队中 ─── */
.stat-pending {
  border: 1px solid rgba(245,158,11,0.2);
  box-shadow: 0 0 20px rgba(245,158,11,0.06), inset 0 0 20px rgba(245,158,11,0.03);
}
.stat-pending .stat-icon { color: #f59e0b; text-shadow: 0 0 12px rgba(245,158,11,0.5); }
.stat-pending .stat-icon::after { border: 1px dashed rgba(245,158,11,0.2); }
.stat-pending .stat-num { color: #fbbf24; text-shadow: 0 0 12px rgba(251,191,36,0.5); }
.stat-pending .stat-label { color: rgba(245,158,11,0.7); }

/* ═══ 运行中 ─── */
.stat-running {
  border: 1px solid rgba(0,255,255,0.2);
  box-shadow: 0 0 24px rgba(0,255,255,0.08), inset 0 0 20px rgba(0,255,255,0.04);
}
.stat-running .stat-icon { color: #00ffff; text-shadow: 0 0 14px rgba(0,255,255,0.6); animation: iconPulse 2s ease-in-out infinite; }
.stat-running .stat-icon::after { border: 1px dashed rgba(0,255,255,0.25); animation: dashSpin 4s linear infinite; }
.stat-running .stat-num { color: #22d3ee; text-shadow: 0 0 14px rgba(34,211,238,0.6); }
.stat-running .stat-label { color: rgba(0,255,255,0.7); }

/* ═══ 已完成 ─── */
.stat-success {
  border: 1px solid rgba(34,197,94,0.2);
  box-shadow: 0 0 20px rgba(34,197,94,0.06), inset 0 0 20px rgba(34,197,94,0.03);
}
.stat-success .stat-icon { color: #22c55e; text-shadow: 0 0 12px rgba(34,197,94,0.5); }
.stat-success .stat-icon::after { border: 1px dashed rgba(34,197,94,0.2); }
.stat-success .stat-num { color: #4ade80; text-shadow: 0 0 12px rgba(74,222,128,0.5); }
.stat-success .stat-label { color: rgba(34,197,94,0.7); }

/* ═══ 已失败 ─── */
.stat-failed {
  border: 1px solid rgba(239,68,68,0.2);
  box-shadow: 0 0 20px rgba(239,68,68,0.06), inset 0 0 20px rgba(239,68,68,0.03);
}
.stat-failed .stat-icon { color: #ef4444; text-shadow: 0 0 12px rgba(239,68,68,0.5); }
.stat-failed .stat-icon::after { border: 1px dashed rgba(239,68,68,0.2); }
.stat-failed .stat-num { color: #f87171; text-shadow: 0 0 12px rgba(248,113,113,0.5); }
.stat-failed .stat-label { color: rgba(239,68,68,0.7); }

/* ═══ 动效 ═══ */
@keyframes tickIn {
  0%   { transform: scale(0.6) translateY(6px); opacity: 0; filter: blur(4px); }
  60%  { transform: scale(1.08); opacity: 1; filter: blur(0); }
  100% { transform: scale(1); opacity: 1; filter: blur(0); }
}
@keyframes iconPulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50%      { opacity: 0.7; transform: scale(1.1); }
}
@keyframes dashSpin {
  0%   { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

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

<style>
/* v-html 渲染的 SVG 图标，需非 scoped 样式 */
.stat-running .stat-icon .run-ring {
  animation: ringSpin 2s linear infinite;
  transform-origin: 12px 12px;
}
@keyframes ringSpin {
  to { transform: rotate(360deg); }
}
</style>
