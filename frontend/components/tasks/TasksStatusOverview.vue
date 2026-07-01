<template>
  <div class="slot-content">
    <div class="section-title">
      <i class="fas fa-chart-pie"></i>
      <span>任务概览</span>
    </div>
    <div class="stats">
      <div
        v-for="s in statusList"
        :key="s.key"
        class="stat-row"
      >
        <div class="stat-left">
          <span class="dot" :class="'dot--' + s.key"></span>
          <span class="stat-label">{{ s.label }}</span>
        </div>
        <span class="stat-num">{{ s.count }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  statusList: { type: Array, default: () => [] },
})


</script>

<style scoped>
.slot-content { padding: 0 2px; }
.section-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 10px;
  font-weight: 600;
  color: rgba(255,255,255,0.25);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 10px;
}
.section-title i { font-size: 11px; }
.stats {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.stat-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 8px;
  border-radius: 8px;
  transition: background 0.2s;
}
.stat-row:hover { background: rgba(255,255,255,0.03); }
.stat-left { display: flex; align-items: center; gap: 8px; }
.dot {
  width: 7px; height: 7px;
  border-radius: 50%;
  flex-shrink: 0;
}
.dot--success { background: #34d399; }
.dot--running { background: #facc15; }
.dot--pending { background: #60a5fa; }
.dot--failed  { background: #f87171; }
.stat-label {
  font-size: 12px;
  color: rgba(255,255,255,0.45);
}
.stat-num {
  font-size: 13px;
  font-weight: 700;
  color: rgba(255,255,255,0.6);
  font-variant-numeric: tabular-nums;
}
</style>
