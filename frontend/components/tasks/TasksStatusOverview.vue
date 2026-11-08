<template>
  <div class="slot-content">
    <div class="section-title">
      <i class="fas fa-chart-pie"></i>
      <span>任务概览</span>
      <span v-if="total" class="title-extra">成功率 {{ successRate }}%</span>
    </div>

    <div class="overview">
      <!-- 环形图：按状态占比分段，中心是任务总数 -->
      <div class="donut-wrap">
        <svg class="donut" viewBox="0 0 80 80" aria-hidden="true">
          <g transform="rotate(-90 40 40)">
            <circle class="donut-track" cx="40" cy="40" :r="R" />
            <circle
              v-for="seg in segments"
              :key="seg.key"
              class="donut-seg"
              cx="40"
              cy="40"
              :r="R"
              :stroke="seg.color"
              :stroke-dasharray="seg.dasharray"
              :stroke-dashoffset="seg.dashoffset"
              :style="{ filter: `drop-shadow(0 0 3px ${seg.glow})` }"
            />
          </g>
        </svg>
        <div class="donut-center">
          <span class="donut-total">{{ total }}</span>
          <span class="donut-caption">总任务</span>
        </div>
      </div>

      <!-- 图例：状态点 + 名称 + 数量（颜色与环形图一致） -->
      <ul class="legend">
        <li
          v-for="s in statusList"
          :key="s.key"
          class="legend-row"
          :class="'legend-row--' + s.key"
        >
          <span class="dot"></span>
          <span class="legend-label">{{ s.label }}</span>
          <span class="legend-count">{{ s.count }}</span>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  statusList: { type: Array, default: () => [] },
})

// 环形图与图例共用同一套状态色（与任务卡状态徽章一致）
const COLORS = {
  success: '#34d399',
  running: '#facc15',
  pending: '#60a5fa',
  failed: '#f87171',
}
const colorOf = (key) => COLORS[key] || '#94a3b8'

const total = computed(() => props.statusList.reduce((sum, it) => sum + (Number(it.count) || 0), 0))

// 成功率：成功数 / 总数（标题右侧）
const successRate = computed(() => {
  if (!total.value) return 0
  const ok = props.statusList.find((s) => s.key === 'success')
  return Math.round(((Number(ok?.count) || 0) / total.value) * 100)
})

// 环形图分段：圆周长按占比切段；段间留 2.5 的缝，只有一段时整圈闭合
const R = 33
const CIRC = 2 * Math.PI * R
const segments = computed(() => {
  if (!total.value) return []
  const filled = props.statusList.filter((s) => (Number(s.count) || 0) > 0)
  const gap = filled.length > 1 ? 2.5 : 0
  let acc = 0
  return filled.map((s) => {
    const raw = ((Number(s.count) || 0) / total.value) * CIRC
    const len = Math.max(raw - gap, 0.6)
    const seg = {
      key: s.key,
      color: colorOf(s.key),
      glow: colorOf(s.key) + '66',
      dasharray: `${len} ${CIRC - len}`,
      dashoffset: -acc,
    }
    acc += raw
    return seg
  })
})
</script>

<style scoped>
.slot-content { padding: 0 2px; }

/* 标题行：图标 + 标题 + 右侧成功率（更高特异性覆盖 .section-title span） */
.section-title {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 10px;
}
.section-title i { font-size: 11px; color: rgba(255, 255, 255, 0.75); }
.section-title span { font-size: 12px; font-weight: 600; color: rgba(255, 255, 255, 0.75); }
.section-title .title-extra {
  margin-left: auto;
  font-size: 10px;
  font-weight: 500;
  color: #fff;
  font-variant-numeric: tabular-nums;
}

/* 主体：环形图在上居中，图例在下方两列排布 */
.overview {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}
.donut-wrap {
  position: relative;
  width: 112px;
  height: 112px;
  flex-shrink: 0;
}
.donut { width: 100%; height: 100%; display: block; }
/* stroke-width 是 viewBox 单位，随圆半径同比缩放 */
.donut-track, .donut-seg { fill: none; stroke-width: 6.5; }
.donut-track { stroke: rgba(255, 255, 255, 0.06); }
/* 数字变化时环形图平滑过渡 */
.donut-seg { transition: stroke-dasharray 0.45s ease, stroke-dashoffset 0.45s ease; }
.donut-center {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  pointer-events: none;
}
.donut-total {
  font-size: 28px;
  font-weight: 700;
  line-height: 1;
  color: #f1f5f9;
  font-variant-numeric: tabular-nums;
}
.donut-caption {
  margin-top: 4px;
  font-size: 9px;
  letter-spacing: 0.5px;
  color: rgba(255, 255, 255, 0.3);
}

/* 图例：圆下方 2×2 状态格，颜色取自 --c */
.legend {
  width: 100%;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 6px;
  margin: 0;
  padding: 0;
  list-style: none;
}
.legend-row {
  --c: 148, 163, 184;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 8px;
  border-radius: 8px;
  background: rgba(var(--c), 0.06);
  border: 1px solid rgba(var(--c), 0.14);
  transition: border-color 0.2s, background 0.2s;
}
.legend-row:hover {
  border-color: rgba(var(--c), 0.35);
  background: rgba(var(--c), 0.1);
}
.legend-row--success { --c: 52, 211, 153; }
.legend-row--running { --c: 250, 204, 21; }
.legend-row--pending { --c: 96, 165, 250; }
.legend-row--failed  { --c: 248, 113, 113; }
.dot {
  width: 6px; height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
  background: rgb(var(--c));
  box-shadow: 0 0 5px rgba(var(--c), 0.55);
}
.legend-label {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.45);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.legend-count {
  margin-left: auto;
  font-size: 11px;
  font-weight: 700;
  color: rgb(var(--c));
  font-variant-numeric: tabular-nums;
}

/* 生成中：呼吸灯强调“正在进行” */
@keyframes statPulse {
  0%, 100% { box-shadow: 0 0 4px rgba(var(--c), 0.5); }
  50% { box-shadow: 0 0 10px rgba(var(--c), 0.9); }
}
.legend-row--running .dot { animation: statPulse 1.4s ease-in-out infinite; }
</style>
