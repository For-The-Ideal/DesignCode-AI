<template>
  <section class="screen">
    <div class="glow-layer glow-1"></div>
    <div class="glow-layer glow-2"></div>
    <div class="glow-layer glow-3"></div>
    <div class="scan-line"></div>
    <div class="hero-content">
      <div class="badge-wrapper">
        <div class="badge-glow"></div>
        <div class="badge">
          <i class="fas fa-circle dot-left"></i>
          <span>AI POWERED · IMAGE SCANNING · REAL-TIME</span>
          <i class="fas fa-circle dot-right"></i>
        </div>
      </div>
      <h1 class="hero-title">
        <span class="title-text">设计稿秒变代码</span>
        <span class="title-shadow">设计稿秒变代码</span>
      </h1>
      <p class="hero-description">上传设计稿，AI自动扫描分析，生成Flutter、React、Vue代码，像素级还原</p>
      <div class="hero-stats">
        <div class="stat-item"><div class="stat-number">98%</div><div class="stat-label">准确率</div></div>
        <div class="stat-item"><div class="stat-number">2.3s</div><div class="stat-label">平均生成</div></div>
        <div class="stat-item"><div class="stat-number">3+</div><div class="stat-label">框架支持</div></div>
      </div>
      <div class="hero-actions">
        <button class="btn-primary" @click="$emit('getStarted')"><i class="fas fa-arrow-right"></i> 立即开始</button>
        <button class="btn-secondary" @click="$emit('watchDemo')"><i class="fas fa-play"></i> 观看演示</button>
      </div>
    </div>
    <div class="scroll-hint"><span>▼ 滚动探索 AI 代码生成 ▼</span></div>
  </section>
</template>

<script setup>
defineEmits(['getStarted', 'watchDemo'])
</script>

<style scoped>
.screen {
  scroll-snap-align: start;
  min-height: 90vh;
  width: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  position: relative;
}
.hero-content {
  text-align: center;
  max-width: 900px;
  margin: 0 auto;
}
.glow-layer {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  pointer-events: none;
  animation: floatGlow 8s ease-in-out infinite;
  opacity: 0.4;
  z-index: 0;
}
.glow-1 {
  top: 10%;
  left: -10%;
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(0, 255, 255, 0.3), transparent);
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
  0%,100% { transform: translate(0,0) scale(1); opacity:0.3; }
  50% { transform: translate(30px,-20px) scale(1.1); opacity:0.6; }
}
.scan-line {
  position: absolute;
  top: -100%;
  left: 0;
  width: 100%;
  height: 3px;
  background: linear-gradient(90deg, transparent, #00ffff, #ff00ff, #00ffff, transparent);
  animation: scanMove 10s linear infinite;
  pointer-events: none;
  opacity: 0.5;
  z-index: 1;
}
@keyframes scanMove {
  0% { top: -100%; }
  100% { top: 200%; }
}
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
  0%,100% { opacity:0; transform: translate(-50%,-50%) scale(0.95);}
  50% { opacity:0.5; transform: translate(-50%,-50%) scale(1.05);}
}
.badge {
  position: relative;
  display: inline-flex;
  align-items: center;
  gap: 12px;
  padding: 8px 24px;
  background: rgba(10,15,25,0.7);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(0,255,255,0.4);
  border-radius: 60px;
  font-size: 14px;
  font-weight: 500;
  color: #00ffff;
  z-index: 1;
}
.hero-title {
  margin-bottom: 20px;
  font-size: 64px;
  font-weight: 800;
}
.title-text {
  background: linear-gradient(135deg, #00ffff, #ff00ff, #00ffff);
  background-size: 200% auto;
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  animation: gradientFlow 3s linear infinite;
}
@keyframes gradientFlow {
  0% { background-position: 0% 50%;}
  100% { background-position: 200% 50%;}
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
}
.hero-description {
  font-size: 18px;
  color: rgba(255,255,255,0.7);
  max-width: 700px;
  margin: 0 auto 32px;
}
.hero-stats {
  display: flex;
  justify-content: center;
  gap: 60px;
  margin-bottom: 40px;
}
.stat-number {
  font-size: 36px;
  font-weight: 800;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}
.stat-label {
  font-size: 14px;
  color: #aaa;
}
.hero-actions {
  display: flex;
  justify-content: center;
  gap: 20px;
}
.btn-primary, .btn-secondary {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 14px 32px;
  border-radius: 40px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
}
.btn-primary {
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  border: none;
  color: white;
  box-shadow: 0 4px 15px rgba(0,255,255,0.3);
}
.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0,255,255,0.4);
}
.btn-secondary {
  background: transparent;
  border: 1px solid rgba(0,255,255,0.3);
  color: #00ffff;
}
.btn-secondary:hover {
  background: rgba(0,255,255,0.1);
  transform: translateY(-2px);
}
.scroll-hint {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 0.8rem;
  color: #6f9eff;
  animation: bounce 1.8s infinite;
  z-index: 10;
}
@keyframes bounce {
  0%,100% { transform: translateX(-50%) translateY(0);}
  50% { transform: translateX(-50%) translateY(8px);}
}
@media (max-width: 768px) {
  .hero-title { font-size: 42px; }
  .hero-description { font-size: 16px; }
  .stat-number { font-size: 28px; }
  .hero-stats { gap: 30px; }
}
</style>