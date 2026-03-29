<template>
  <section class="hero-section">
    <!-- 动态光晕层 -->
    <div class="glow-layer glow-1"></div>
    <div class="glow-layer glow-2"></div>
    <div class="glow-layer glow-3"></div>

    <!-- 扫描线效果 -->
    <div class="scan-line"></div>

    <!-- 粒子漂浮效果 -->
    <div class="particles-container">
      <div
        v-for="i in 100"
        :key="i"
        class="particle"
        :style="getParticleStyle(i)"
      ></div>
    </div>

    <div class="hero-container">
      <!-- 左侧：内容区域 -->
      <div class="hero-content">
        <!-- 徽章区域 -->
        <div class="badge-wrapper">
          <div class="badge-glow"></div>
          <div class="badge">
            <i class="fas fa-circle dot-left"></i>
            <span>{{ badgeText }}</span>
            <i class="fas fa-circle dot-right"></i>
          </div>
        </div>

        <!-- 主标题 -->
        <h1 class="hero-title">
          <span class="title-text">{{ title }}</span>
          <span class="title-shadow">{{ title }}</span>
        </h1>

        <!-- 副标题/描述 -->
        <p class="hero-description">{{ description }}</p>

        <!-- 统计数据 -->
        <div class="hero-stats">
          <div class="stat-item" v-for="stat in stats" :key="stat.label">
            <div class="stat-number">{{ stat.value }}</div>
            <div class="stat-label">{{ stat.label }}</div>
          </div>
        </div>

        <!-- CTA 按钮组 -->
        <div class="hero-actions">
          <button class="btn-primary" @click="$emit('get-started')">
            <i class="fas fa-arrow-right"></i>
            立即开始
          </button>
          <button class="btn-secondary" @click="$emit('watch-demo')">
            <i class="fas fa-play"></i>
            观看演示
          </button>
        </div>
      </div>

      <!-- 右侧：卡片区域 -->
      <CardsArea />
    </div>

    <!-- 装饰线 -->
    <div class="decorative-line"></div>
  </section>
</template>

<script setup>
import { onMounted } from 'vue'
import CardsArea from './CardsArea.vue'

const props = defineProps({
  badgeText: {
    type: String,
    default: "AI POWERED · IMAGE SCANNING · REAL-TIME",
  },
  title: {
    type: String,
    default: "设计稿秒变代码",
  },
  description: {
    type: String,
    default:
      "上传设计稿，AI 自动扫描分析，生成 Flutter、React、Vue 代码，像素级还原",
  },
  stats: {
    type: Array,
    default: () => [
      { value: "98%", label: "准确率" },
      { value: "2.3s", label: "平均生成" },
      { value: "3+", label: "框架支持" },
    ],
  },
})

const emit = defineEmits(["get-started", "watch-demo"])

const getParticleStyle = () => {
  const size = Math.random() * 3 + 1
  const left = Math.random() * 100
  const duration = Math.random() * 20 + 10
  const delay = Math.random() * 10
  return {
    width: `${size}px`,
    height: `${size}px`,
    left: `${left}%`,
    animationDuration: `${duration}s`,
    animationDelay: `${delay}s`,
  }
}
</script>

<style scoped>
.hero-section {
  position: relative;
  padding: 60px 24px 40px;
  overflow: hidden;
  background: transparent;
}

.hero-container {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 60px;
  position: relative;
  z-index: 2;
}

/* ========== 左侧内容区域 ========== */
.hero-content {
  flex: 1;
  text-align: left;
}

/* ========== 光晕层 ========== */
.glow-layer {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  pointer-events: none;
  animation: floatGlow 8s ease-in-out infinite;
  opacity: 0.4;
}

.glow-1 {
  top: 10%;
  left: -10%;
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(0, 255, 255, 0.3), transparent);
  animation-delay: 0s;
}

.glow-2 {
  bottom: 0;
  right: -10%;
  width: 500px;
  height: 500px;
  background: radial-gradient(circle, rgba(255, 0, 255, 0.25), transparent);
  animation-delay: -2s;
}

.glow-3 {
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 600px;
  height: 600px;
  background: radial-gradient(circle, rgba(0, 255, 255, 0.1), transparent);
  animation-delay: -4s;
}

@keyframes floatGlow {
  0%, 100% {
    transform: translate(0, 0) scale(1);
    opacity: 0.3;
  }
  50% {
    transform: translate(30px, -20px) scale(1.1);
    opacity: 0.6;
  }
}

/* ========== 扫描线 ========== */
.scan-line {
  position: absolute;
  top: -100%;
  left: 0;
  width: 100%;
  height: 3px;
  background: linear-gradient(
    90deg,
    transparent,
    #00ffff,
    #ff00ff,
    #00ffff,
    transparent
  );
  animation: scanMove 10s linear infinite;
  pointer-events: none;
  opacity: 0.5;
}

@keyframes scanMove {
  0% {
    top: -100%;
  }
  100% {
    top: 200%;
  }
}

/* ========== 粒子效果 ========== */
.particles-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  overflow: hidden;
}

.particle {
  position: absolute;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  border-radius: 50%;
  opacity: 0;
  animation: floatParticle linear infinite;
}

@keyframes floatParticle {
  0% {
    transform: translateY(100vh) rotate(0deg);
    opacity: 0;
  }
  10% {
    opacity: 0.5;
  }
  90% {
    opacity: 0.5;
  }
  100% {
    transform: translateY(-100vh) rotate(360deg);
    opacity: 0;
  }
}

/* ========== 徽章 ========== */
.badge-wrapper {
  position: relative;
  display: inline-block;
  margin-bottom: 24px;
}

.badge-glow {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: calc(100% + 8px);
  height: calc(100% + 8px);
  background: linear-gradient(90deg, #00ffff, #ff00ff);
  border-radius: 60px;
  filter: blur(12px);
  opacity: 0;
  animation: badgePulse 2s ease-in-out infinite;
}

@keyframes badgePulse {
  0%, 100% {
    opacity: 0;
    transform: translate(-50%, -50%) scale(0.95);
  }
  50% {
    opacity: 0.5;
    transform: translate(-50%, -50%) scale(1.05);
  }
}

.badge {
  position: relative;
  display: inline-flex;
  align-items: center;
  gap: 12px;
  padding: 8px 24px;
  background: rgba(10, 15, 25, 0.7);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(0, 255, 255, 0.4);
  border-radius: 60px;
  font-size: 12px;
  font-weight: 500;
  letter-spacing: 1px;
  color: #00ffff;
  z-index: 1;
  transition: all 0.3s;
}

.badge:hover {
  border-color: rgba(0, 255, 255, 0.8);
  box-shadow: 0 0 20px rgba(0, 255, 255, 0.2);
}

.badge i {
  font-size: 8px;
  animation: dotPulse 1.2s ease-in-out infinite;
}

.dot-left {
  animation-delay: 0s;
}
.dot-right {
  animation-delay: 0.6s;
}

@keyframes dotPulse {
  0%, 100% {
    opacity: 0.3;
    transform: scale(0.8);
  }
  50% {
    opacity: 1;
    transform: scale(1.2);
    text-shadow: 0 0 8px #00ffff;
  }
}

/* ========== 主标题 ========== */
.hero-title {
  position: relative;
  margin-bottom: 20px;
  font-size: 48px;
  font-weight: 800;
  line-height: 1.2;
  letter-spacing: -0.02em;
}

.title-text {
  position: relative;
  background: linear-gradient(135deg, #00ffff, #ff00ff, #00ffff);
  background-size: 200% auto;
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  animation: gradientFlow 3s linear infinite;
  display: inline-block;
}

.title-shadow {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  filter: blur(20px);
  opacity: 0.3;
  pointer-events: none;
  animation: shadowGlow 2s ease-in-out infinite;
}

@keyframes gradientFlow {
  0% {
    background-position: 0% 50%;
  }
  100% {
    background-position: 200% 50%;
  }
}

@keyframes shadowGlow {
  0%, 100% {
    opacity: 0.2;
    filter: blur(15px);
  }
  50% {
    opacity: 0.5;
    filter: blur(25px);
  }
}

/* ========== 描述文字 ========== */
.hero-description {
  max-width: 100%;
  margin-bottom: 32px;
  font-size: 16px;
  line-height: 1.6;
  color: rgba(255, 255, 255, 0.7);
}

/* ========== 统计数据 ========== */
.hero-stats {
  display: flex;
  gap: 40px;
  margin-bottom: 32px;
}

.stat-item {
  text-align: left;
}

.stat-number {
  font-size: 28px;
  font-weight: 800;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  line-height: 1.2;
}

.stat-label {
  font-size: 12px;
  color: #888;
  margin-top: 4px;
}

/* ========== CTA 按钮 ========== */
.hero-actions {
  display: flex;
  gap: 16px;
  margin-bottom: 32px;
  flex-wrap: wrap;
}

.btn-primary {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 28px;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  border: none;
  border-radius: 40px;
  color: white;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 4px 15px rgba(0, 255, 255, 0.3);
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 255, 255, 0.4);
}

.btn-secondary {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 28px;
  background: transparent;
  border: 1px solid rgba(0, 255, 255, 0.3);
  border-radius: 40px;
  color: #00ffff;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-secondary:hover {
  background: rgba(0, 255, 255, 0.1);
  transform: translateY(-2px);
}

/* ========== 装饰线 ========== */
.decorative-line {
  width: 100px;
  height: 2px;
  background: linear-gradient(
    90deg,
    transparent,
    #00ffff,
    #ff00ff,
    transparent
  );
  margin: 60px auto 0;
  animation: lineExpand 2s ease-in-out infinite;
}

@keyframes lineExpand {
  0%, 100% {
    width: 80px;
    opacity: 0.4;
  }
  50% {
    width: 140px;
    opacity: 1;
  }
}

/* ========== 响应式 ========== */
@media (max-width: 968px) {
  .hero-container {
    flex-direction: column;
    gap: 60px;
  }

  .hero-content {
    text-align: center;
  }

  .hero-description {
    margin-left: auto;
    margin-right: auto;
  }

  .hero-stats {
    justify-content: center;
  }

  .stat-item {
    text-align: center;
  }

  .hero-actions {
    justify-content: center;
  }

  .hero-title {
    font-size: 36px;
  }
}

@media (max-width: 768px) {
  .hero-section {
    padding: 40px 20px 30px;
  }

  .hero-title {
    font-size: 28px;
  }

  .hero-description {
    font-size: 14px;
  }

  .stat-number {
    font-size: 20px;
  }

  .hero-stats {
    gap: 24px;
  }
}
</style>