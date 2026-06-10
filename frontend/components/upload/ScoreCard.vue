<template>
  <div class="score-card" v-if="score > 0">
    <div class="score-header">
      <div>
        <h3>质量评分</h3>
        <p class="score-sub">AI 多维评估</p>
      </div>
      <div class="score-value">{{ score }}</div>
    </div>
    <div class="score-dimensions">
      <div v-for="dim in scoreDimensions" :key="dim.name" class="dimension-item">
        <div class="dimension-label">
          <span><i :class="dim.icon"></i> {{ dim.name }}</span>
          <span>{{ dim.score }}</span>
        </div>
        <div class="progress-bar">
          <div class="progress-fill" :style="{ width: dim.score + '%' }"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  score: { type: Number, default: 0 },
  scoreDimensions: { type: Array, default: () => [] },
})
</script>

<style scoped>
.score-card {
  background: rgba(15, 20, 30, 0.6);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(0, 255, 255, 0.2);
  border-radius: 24px;
  overflow: hidden;
  margin: 24px 0;
}

.score-header {
  background: linear-gradient(135deg, rgba(0, 255, 255, 0.2), rgba(255, 0, 255, 0.2));
  padding: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.score-header h3 {
  font-size: 18px;
}

.score-sub {
  font-size: 12px;
  opacity: 0.7;
  margin-top: 4px;
}

.score-value {
  font-size: 48px;
  font-weight: 700;
  color: #00ffff;
}

.score-dimensions {
  padding: 20px;
}

.dimension-item {
  margin-bottom: 16px;
}

.dimension-label {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
  margin-bottom: 8px;
  color: #aaa;
}

.progress-bar {
  width: 100%;
  height: 6px;
  background: #333;
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #00ffff, #ff00ff);
  border-radius: 3px;
}
</style>
